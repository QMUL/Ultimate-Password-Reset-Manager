// Package prm provides all the logic and glue for changing passwords.
// It makes use of ldap primarily, but also calls email, ntmlgen and
// crypto modules
package prm

// TODO
// * Potentially put consts and such in a text file for better integration
// * TLS and certs stuff (should be in ldap already)
// * cracklib equivalent check

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"gopkg.in/ldap.v2"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// PRM is a global struct object doofus - holds the secret config
type PRM struct {
	Config *PRMConfig
}

// Conn - exported functions we can perfom on our LDAP
type Conn interface {
	Bind(username, password string) error
	PasswordModify(passwordModifyRequest *ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error)
	Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error)
	Modify(modifyRequest *ldap.ModifyRequest) error
}

// Return status types
const (
	Success                = 1
	ErrorLDAP              = 2
	ErrorNoUser            = 3
	ErrorTimeOut           = 4
	ErrorPasswordMatch     = 5
	ErrorPasswordStrength  = 6
	ErrorPasswordLength    = 7
	ErrorPasswordIncorrect = 8
	ErrorNotImplemented    = 9
	ErrorOTP               = 10
	ErrorOTPError          = 11
	ErrorFatal             = 12
	ErrorDeclined          = 13
	ErrorOTPExpired        = 14
	SuccessFinished        = 15
)

// ResultMap is a map to provide useful strings for the errors and successes.
var ResultMap = map[int]string{
	Success:                "Success",
	ErrorLDAP:              "Error with LDAP Call",
	ErrorNoUser:            "Error; your username or password is incorrect",
	ErrorTimeOut:           "Error; time limit exceeded",
	ErrorPasswordMatch:     "Error; passwords provided do not match",
	ErrorPasswordStrength:  "Error; your password is not strong enough",
	ErrorPasswordLength:    "Error; your password is not long enough",
	ErrorPasswordIncorrect: "Error; your username or password is incorrect",
	ErrorNotImplemented:    "Error; this function has not been implemented",
	ErrorOTP:               "Error; One-time account unlocking code failed. Please contact its-research-support@qmul.ac.uk for a new code.",
	ErrorOTPExpired:        "Error; your one-time unlocking code has expired. Please contact its-research-support@qmul.ac.uk for a new code.",
	ErrorDeclined:          "Error; you must accept the terms and conditions to continue",
	SuccessFinished:        "Success: your password has been changed",
}

// Result is simply an int code from the return status types given above.
type Result struct {
	Message int
}

// ToString probably needs a better name; it just converts int errors to their human strings.
func (r *Result) ToString() string {
	return ResultMap[r.Message]
}

// parseForm takes a http.Request and looks for the form data and extracts it.
func (prm *PRM) parseForm(r *http.Request) (string, string, string, string, string, string, string, string, string) {
	r.ParseForm()
	username := strings.Join(r.Form["user"], "")
	p0 := strings.Join(r.Form["p0"], "")
	p1 := strings.Join(r.Form["p1"], "")
	p2 := strings.Join(r.Form["p2"], "")
	otp := strings.Join(r.Form["otp"], "")
	puffer := strings.Join(r.Form["puffer"], "")
	wuffer := strings.Join(r.Form["wuffer"], "")
	tuffer := strings.Join(r.Form["tuffer"], "")
	verb := strings.Join(r.Form["verb"], "")
	return username, p0, p1, p2, otp, puffer, wuffer, tuffer, verb
}

// ProcessSkipped takes the first form and immediately skips the terms and conditions
// and tries to change the password
func (prm *PRM) ProcessSkipped(r *http.Request) (result Result, data map[string]string) {

	username, p0, p1, _, _, _, _, _, _ := prm.parseForm(r)

	conn, err := prm.ldapConnect()
	err = prm.ldapBindAdmin(conn)

	defer conn.Close()

	if err != nil {
		prm.LogPRM(err.Error(), LOG_ERROR)
		return Result{ErrorFatal}, nil
	}

	newpassword := p1

	// Check that the username passed is legit to stop attacks on the hash
	entry := prm.SearchUsername(username, conn)
	if entry == nil {
		return Result{ErrorFatal}, nil
	}

	// Double check the existing password
	if !prm.CheckPasswordCorrect(username, p0, conn) {
		return Result{ErrorPasswordIncorrect}, nil
	}

	// Rebind as admin due to the check password correct above
	err = prm.ldapBindAdmin(conn)

	if err != nil {
		prm.LogPRM(err.Error(), LOG_ERROR)
		return Result{ErrorFatal}, nil
	}

	if !prm.ChangeLDAPPassword(username, newpassword, conn) {
		return Result{ErrorLDAP}, nil
	}

	// Check for a Linux password
	if !prm.ChangeLinuxPassword(username, newpassword, conn) {
		return Result{ErrorFatal}, nil
	}

	// Check for the Samba Password
	if !prm.ChangeSambaPassword(username, newpassword, conn) {
		return Result{ErrorFatal}, nil
	}

	// If all is well, send the email
	name, addy := prm.GetEmailDeets(username, conn)
	SendEmail(name, addy, prm.Config)

	m := make(map[string]string)
	m["username"] = username

	return Result{SuccessFinished}, m
}

// ProcessTerms deals with the acceptance of the terms and conditions form which is
// the second page after a correct series of inputs from the user.
func (prm *PRM) ProcessTerms(r *http.Request) (result Result, data map[string]string) {
	_, _, _, _, _, puffer, wuffer, tuffer, verb := prm.parseForm(r)
	conn, err := prm.ldapConnect()
	err = prm.ldapBindAdmin(conn)

	defer conn.Close()

	if err != nil {
		prm.LogPRM(err.Error(), LOG_ERROR)
		return Result{ErrorFatal}, nil
	}

	if verb != "Accept" {
		return Result{ErrorDeclined}, nil
	}

	username := decryptUffer(wuffer, prm.Config.Uffer)
	newpassword := decryptUffer(puffer, prm.Config.Uffer)
	starttime, err := strconv.ParseInt(decryptUffer(tuffer, prm.Config.Uffer), 10, 64)

	// Check that the username passed is legit to stop attacks on the hash
	entry := prm.SearchUsername(username, conn)
	if entry == nil {
		return Result{ErrorFatal}, nil
	}

	// Check to see if time has expired
	now := time.Now().Unix()
	if now-starttime > 1800 {
		return Result{ErrorTimeOut}, nil
	}

	if !prm.ChangeLDAPPassword(username, newpassword, conn) {
		return Result{ErrorLDAP}, nil
	}

	// Check for a Linux password
	if !prm.ChangeLinuxPassword(username, newpassword, conn) {
		return Result{ErrorFatal}, nil
	}

	// Check for the Samba Password
	if !prm.ChangeSambaPassword(username, newpassword, conn) {
		return Result{ErrorFatal}, nil
	}

	// If all is well, send the email
	name, addy := prm.GetEmailDeets(username, conn)
	SendEmail(name, addy, prm.Config)

	m := make(map[string]string)
	m["username"] = username

	return Result{SuccessFinished}, m
}

// ProcessForm deals with the intial form, doing all the various checks and calls to ldap
// it returns a Result which is then checked, directing the flow to pass or fail.
func (prm *PRM) ProcessForm(r *http.Request) (result Result, data map[string]string) {
	r.ParseForm()
	conn, err := prm.ldapConnect()
	err = prm.ldapBindAdmin(conn)

	if err != nil {
		prm.LogPRM(err.Error(), LOG_ERROR)
		return Result{ErrorFatal}, nil
	}

	defer conn.Close()

	username, p0, p1, p2, otp, _, _, _, _ := prm.parseForm(r)

	// Find the user
	entry := prm.SearchUsername(username, conn)
	if entry == nil {
		return Result{ErrorNoUser}, nil
	}

	// Check new passwords match..
	if !(p1 == p2) {
		return Result{ErrorPasswordMatch}, nil
	}

	// ... or are too weak
	if !prm.passwordStrength(p1) {
		return Result{ErrorPasswordStrength}, nil
	}

	// ... or fail the #cracklib check
	cracklibMsg := TestPassword(p1)
	if cracklibMsg != "GOOD" {
		return Result{ErrorPasswordStrength}, nil
	}

	// ... or are too short
	if len(p1) < 9 {
		return Result{ErrorPasswordLength}, nil
	}

	m := make(map[string]string)
	m["wuffer"] = createUffer(username, prm.Config.Uffer)
	m["puffer"] = createUffer(p1, prm.Config.Uffer)
	m["tuffer"] = createUffer(fmt.Sprintf("%d", time.Now().Unix()), prm.Config.Uffer)
	m["username"] = username
	m["otp"] = otp

	// Decide upon one-time-code or password as the way to go
	if len(otp) > 0 {
		result, code := prm.CheckOTP(username, otp, conn)
		if !result {
			return Result{code}, nil
		}
		return Result{Success}, m
	}

	if !prm.CheckPasswordCorrect(username, p0, conn) {
		return Result{ErrorPasswordIncorrect}, nil
	}

	return Result{Success}, m

}

// passwordStrength is a null function that can probably be removed
func (prm *PRM) passwordStrength(password string) (result bool) {
	return true
}

// CheckPasswordCorrect checks with LDAP to make sure we can login as this user with this password
// It returns true if all the ldap details provided are correct or false otherwise
func (prm *PRM) CheckPasswordCorrect(username string, password string, conn Conn) (result bool) {
	err := conn.Bind(fmt.Sprintf(prm.Config.PasswordModifyLDAP+",%v", username, prm.Config.BaseDN), password)
	if err != nil {
		prm.LogPRM(err.Error(), LOG_ERROR)
		return false
	}

	return true
}

// ChangeLDAPPassword actually changes the LDAP Password - it appears this is somewhat messy in the original prm
func (prm *PRM) ChangeLDAPPassword(username string, newpassword string, conn Conn) (result bool) {
	passwordModifyRequest := ldap.NewPasswordModifyRequest(fmt.Sprintf(prm.Config.PasswordModifyLDAP+",%v", username, prm.Config.BaseDN), "", newpassword)

	_, err := conn.PasswordModify(passwordModifyRequest)

	if err != nil {
		prm.LogPRM(err.Error(), LOG_ERROR)
		return false
	}

	return true
}

// ChangeLinuxPassword change the Linux password
// So far it always returns true which we might need to fix
func (prm *PRM) ChangeLinuxPassword(username string, newpassword string, conn Conn) (result bool) {
	hash, _ := CreatePasswordHash(newpassword)
	modify := ldap.NewModifyRequest(fmt.Sprintf(prm.Config.PasswordModifyLDAP+",%v", username, prm.Config.BaseDN))
	modify.Replace("userPassword", []string{hash})

	return true
}

// ChangeSambaPassword checks to see if there is a samba password and changes it
// Returns true if successful and false if not
func (prm *PRM) ChangeSambaPassword(username string, newpassword string, conn Conn) (result bool) {

	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("%v,%v", prm.Config.ORGFieldLDAP, prm.Config.BaseDN),
		ldap.ScopeSingleLevel, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%v=%v)", prm.Config.UserFieldLDAP, username),
		nil,
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		prm.LogPRM("ChangeSambaPassword Error: "+err.Error(), LOG_ERROR)
		return false
	}

	// If we dont have a single entry we cant change anything but if its zero
	// the user doesnt have one which is fine
	if len(sr.Entries) != 1 {
		prm.LogPRM("ChangeSambaPassword Warning: entries = -1", LOG_WARN)
		if len(sr.Entries) == 0 {
			return true
		}
		return false
	}

	entrytexts := sr.Entries[0].GetAttributeValues("objectClass")

	for _, entrytext := range entrytexts {

		if strings.Contains(entrytext, "sambaSamAccount") {

			modify := ldap.NewModifyRequest(fmt.Sprintf(prm.Config.PasswordModifyLDAP+",%v", username, prm.Config.BaseDN))
			modify.Replace("sambaNTPassword", []string{Ntlmgen(newpassword)})
			err = conn.Modify(modify)

			if err != nil {
				prm.LogPRM("[prm:error] ChangeSambaPassword Error:"+err.Error(), LOG_ERROR)
				return false
			}

			modify = ldap.NewModifyRequest(fmt.Sprintf(prm.Config.PasswordModifyLDAP+",%v", username, prm.Config.BaseDN))
			modify.Replace("sambaPwdLastSet", []string{fmt.Sprintf("%d", time.Now().Unix())})
			err = conn.Modify(modify)

			if err != nil {
				prm.LogPRM("ChangeSambaPassword Error: "+err.Error(), LOG_ERROR)
				return false
			}

			modify = ldap.NewModifyRequest(fmt.Sprintf(prm.Config.PasswordModifyLDAP+",%v", username, prm.Config.BaseDN))
			modify.Replace("sambaAcctFlags", []string{"[UX         ]"})
			err = conn.Modify(modify)

			if err != nil {
				prm.LogPRM("ChangeSambaPassword Error: "+err.Error(), LOG_ERROR)
				return false
			}
			return true
		}
	}

	return false
}

// parse_otp takes a string, trims it and returns epoch and actual otp
func (prm *PRM) parseOtp(code string) (epoch int64, otp string) {
	otp = strings.TrimLeft(code[:9], "0")
	epoch, _ = strconv.ParseInt(code[9:], 10, 64)
	return epoch, otp
}

// CheckOTP checks to see if the OTP exists and if so, does the one provided match. Returns true if
// everything checks out, or false and an errorcode if not
func (prm *PRM) CheckOTP(username string, userotp string, conn Conn) (result bool, errorcode int) {

	entry := prm.SearchUsername(username, conn)

	if entry == nil {
		prm.LogPRM("CheckOTP Info: entry = nil", LOG_INFO)
		return false, ErrorOTP
	}

	code := entry.GetAttributeValue("internationaliSDNNumber")

	if len(code) < 10 {
		prm.LogPRM("ChangeOTP Info: len(code) < 10", LOG_INFO)
		return false, ErrorOTP
	}

	epoch, storedOtp := prm.parseOtp(code)

	if time.Now().Unix() > epoch {

		prm.LogPRM("ChangeOTP Info: len(code) < 10", LOG_INFO)
		prm.LogPRM("ChangeOTP Info: expired", LOG_WARN)

		return false, ErrorOTPExpired
	}

	if storedOtp == userotp {
		// Success so delete the OTP
		modify := ldap.NewModifyRequest(fmt.Sprintf(prm.Config.PasswordModifyLDAP+",%v", username, prm.Config.BaseDN))
		modify.Delete("internationaliSDNNumber", []string{code})
		err := conn.Modify(modify)

		if err != nil {
			prm.LogPRM("in remove otp "+err.Error(), LOG_INFO)
			return false, ErrorOTP
		}

		return true, Success
	}

	prm.LogPRM("ChangeOTP Info: Failed", LOG_WARN)

	return false, ErrorOTP
}

// GetEmailDeets grabs the email details for a user out of LDAP
// TODO - no error is given here for the user - we just trundle along :S
func (prm *PRM) GetEmailDeets(username string, conn Conn) (givenName string, emailAddr string) {

	givenName = "test"
	emailAddr = "test"

	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf(prm.Config.ORGFieldLDAP+",%v", prm.Config.BaseDN),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("("+prm.Config.UserFieldLDAP+"=%v)", username),
		nil,
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		prm.LogPRM("GetEmailDeets Error:"+err.Error(), LOG_INFO)
		return
	}

	if len(sr.Entries) != 1 {
		return
	}

	givenName = sr.Entries[0].GetAttributeValue("givenName")
	emailAddr = sr.Entries[0].GetAttributeValue("mail")
	prm.LogPRM("GetEmailDeets: "+givenName+" "+emailAddr, LOG_INFO)
	return
}

// SearchUsername searches LDAP to find an entry given a username and connection
func (prm *PRM) SearchUsername(username string, conn Conn) (ldapentry *ldap.Entry) {

	prm.LogPRM("searchusername "+prm.Config.ORGFieldLDAP+",%v", LOG_DEBUG)

	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf(prm.Config.ORGFieldLDAP+",%v", prm.Config.BaseDN),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("("+prm.Config.UserFieldLDAP+"=%v)", username),
		nil,
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		if prm.Config.LogLevel <= LOG_ERROR {
			prm.LogPRM("SearchUsername Error: "+err.Error(), LOG_ERROR)
		}
		return nil
	}

	if len(sr.Entries) == 1 {
		return sr.Entries[0]
	}

	prm.LogPRM("SearchUsername Info : returning not 1", LOG_WARN)

	return nil
}

// ldapBindAdmin binds as the admin user for ldap operations
func (prm *PRM) ldapBindAdmin(conn Conn) error {
	err := conn.Bind(prm.Config.BindDN, prm.Config.BindPassword)
	return err
}

// ldapConnect connects to ldap as a basic non-admin user
func (prm *PRM) ldapConnect() (*ldap.Conn, error) {
	host := fmt.Sprintf("%v:%d", prm.Config.LDAPHost, prm.Config.LDAPPort)
	// Refer to https://golang.org/pkg/crypto/tls/#example_Dial when thinking about certs - it may not work with Vagrant

	certHandle, err := os.Open(prm.Config.CertFilePath)

	if err != nil {
		prm.LogPRM("ldap_connect Error: "+err.Error(), LOG_ERROR)
		return nil, err
	}

	fileInfo, err := certHandle.Stat()
	data := make([]byte, fileInfo.Size())
	_, err = certHandle.Read(data)

	if err != nil {
		prm.LogPRM("ldap_connect Error: "+err.Error(), LOG_ERROR)
		return nil, err
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(data)

	if !ok {
		prm.LogPRM("ldap_connect Error: Failed to parse root certificate", LOG_ERROR)
		return nil, nil
	}

	conn, err := ldap.Dial("tcp", host)
	if err != nil {
		prm.LogPRM("ldap_connect Error:"+err.Error(), LOG_ERROR)
		return nil, err
	}

	err = conn.StartTLS(&tls.Config{RootCAs: roots, InsecureSkipVerify: prm.Config.LDAPInsecureSkipVerify})
	if err != nil {
		prm.LogPRM("ldap_connect Error:"+err.Error(), LOG_ERROR)
	}

	return conn, err
}
