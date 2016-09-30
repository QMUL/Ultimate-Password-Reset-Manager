package prm

// There is probably a better way to do this but we
// are doing type conversion so unless we mess with
// yaml parser, loading yaml then converting seems best

type PRMConfig struct {
	TemplatePath           string
	ListenAddress          string
	LDAPHost               string
	LDAPPort               int
	BindPassword           string
	CertFilePath           string
	BaseDN                 string
	BindDN                 string
	LogLevel               int
	EmailMsg               string
	EmailSub               string
	LDAPInsecureSkipVerify bool
	Uffer                  string
	PasswordModifyLDAP     string
	ORGFieldLDAP           string
	UserFieldLDAP          string
}

type YamlConfig struct {
	TemplatePath           string
	ListenAddress          string
	LDAPHost               string
	LDAPPort               int
	BindPassword           string
	CertFilePath           string
	BaseDN                 string
	BindDN                 string
	LogLevel               string
	EmailMsg               string
	EmailSub               string
	LDAPInsecureSkipVerify bool
	Uffer                  string
	PasswordModifyLDAP     string
	ORGFieldLDAP           string
	UserFieldLDAP          string
}
