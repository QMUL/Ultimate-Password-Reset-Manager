package main_test

import (
	"fmt"
	"pass.hpc.qmul.ac.uk/prm"
	"testing"
)

// Basic tests for the main prm_server section

// Test reading the config correctly
func TestConfig(t *testing.T) {
	config := main.ReadConfig()

	fmt.Println("Config File Params:", config.LDAPHost, config.LDAPPort, config.BindPassword, config.CertFilePath)

	if config.LDAPHost != "192.168.33.4" {
		t.Errorf("192.168.33.4")
	}

	if config.LDAPPort != 389 {
		t.Errorf("LDAP Port not 389")
	}

	if config.BindPassword != "prm" {
		t.Errorf("BindPassword != prm")
	}

	if config.CertFilePath != "/etc/pki/tls/certs/vagrantCA-cert.pem" {
		t.Errorf("CertFilePath /etc/pki/tls/certs/vagrantCA-cert.pem")
	}

}
