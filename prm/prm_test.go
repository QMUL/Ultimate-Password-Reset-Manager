package prm

import (
	"errors"
	"gopkg.in/ldap.v2"
	"net/http"
	"net/url"
	"testing"
)

type TestConn struct {
}

// Override the bind for the ldap
func (l *TestConn) Bind(username, password string) error {
	return nil
}

// Override the Modify function from LDAP
func (l *TestConn) Modify(modifyRequest *ldap.ModifyRequest) error {
	return errors.New("Error")
}

// PasswordModify override
func (l *TestConn) PasswordModify(passwordModifyRequest *ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error) {
	result := &ldap.PasswordModifyResult{}
	return result, nil

}

// Search override
func (l *TestConn) Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error) {
	result := &ldap.SearchResult{}
	return result, nil
}

func TestCheckPasswordCorrect(t *testing.T) {

	var prm = new(PRM)
	prm.Config = new(PRMConfig)

	var conn = new(TestConn)
	result := prm.CheckPasswordCorrect("user", "password", conn)

	if result == false {
		t.Errorf("Error in prm.CheckPasswordCorrect")
	}

}

// Test parsing some forms

func TestFormParse(t *testing.T) {

	var prm = new(PRM)
	req := &http.Request{Method: "GET"}
	username, p0, p1, p2, otp, puffer, wuffer, tuffer, verb := prm.parseForm(req)
	if username != "" {
		t.Errorf("Error in parseForm - username:" + username)
	}
	if p0 != "" {
		t.Errorf("Error in parseForm - p0:" + p0)
	}
	if p1 != "" {
		t.Errorf("Error in parseForm - p1:" + p1)
	}
	if p2 != "" {
		t.Errorf("Error in parseForm - p2:" + p2)
	}
	if otp != "" {
		t.Errorf("Error in parseForm - otp:" + otp)
	}
	if puffer != "" {
		t.Errorf("Error in parseForm - puffer:" + puffer)
	}
	if wuffer != "" {
		t.Errorf("Error in parseForm - wuffer:" + wuffer)
	}
	if tuffer != "" {
		t.Errorf("Error in parseForm - tuffer:" + tuffer)
	}
	if verb != "" {
		t.Errorf("Error in parseForm - verb:" + verb)
	}

	values := url.Values{}

	values.Add("user", "username")
	values.Add("p0", "password0")
	values.Add("p1", "password1")
	values.Add("p2", "password2")
	values.Add("otp", "otp")
	values.Add("puffer", "puffer")
	values.Add("wuffer", "wuffer")
	values.Add("tuffer", "tuffer")
	values.Add("verb", "verb")

	req = &http.Request{Method: "POST", Form: values}

	username, p0, p1, p2, otp, puffer, wuffer, tuffer, verb = prm.parseForm(req)

	if username != "username" {
		t.Errorf("Error in username(2) - username:" + username)
	}

	if p0 != "password0" {
		t.Errorf("Error in parseForm(2) - p0:" + p0)
	}

	if p1 != "password1" {
		t.Errorf("Error in parseForm(2) - p1:" + p1)
	}

	if p2 != "password2" {
		t.Errorf("Error in parseForm(2) - p2:" + p2)
	}

	if otp != "otp" {
		t.Errorf("Error in parseForm(2) - otp:" + otp)
	}

	if puffer != "puffer" {
		t.Errorf("Error in parseForm(2) - puffer:" + puffer)
	}

	if wuffer != "wuffer" {
		t.Errorf("Error in parseForm(2) - wuffer:" + wuffer)
	}

	if tuffer != "tuffer" {
		t.Errorf("Error in parseForm(2) - tuffer:" + tuffer)
	}

	if verb != "verb" {
		t.Errorf("Error in parseForm(2) - verb:" + verb)
	}

}
