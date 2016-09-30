/*
prmserver / prmserver.fcgi

This is essentially just an FCGI Server that directs HTTP traffic to the correct
functions in the prm module. It is designed to be run by apache and not
directly by the user however there is one option that is useful.

Command-line interface:

To return the current version number:

		prm_server -v
*/
package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/fcgi"
	"os"
	"pass.hpc.qmul.ac.uk/prm"
	"path/filepath"
	"strings"
)

// Page is a type for passing messages to HTML Templates
type Page struct {
	Title   string
	Message string
	Wuffer  string
	Puffer  string
	Tuffer  string
}

// FastCGIServer is our basic struct for state on the server
type FastCGIServer struct {
	PRMHandler prm.PRM
}

// renderIndex renders the first index page reading in the index.html template and setting it
func renderIndex(w http.ResponseWriter, r *http.Request, p *prm.PRM) {
	config := p.Config
	title := r.URL.Path[len("/"):]
	g := &Page{Title: title, Message: ""}
	t, _ := template.ParseFiles(config.TemplatePath + "index.html")
	t.Execute(w, g)
}

// processForm handles the form on the first page, passing the form to the prm module
func processForm(w http.ResponseWriter, r *http.Request, p *prm.PRM) {
	config := p.Config
	result, data := p.ProcessForm(r)
	if result.Message != prm.Success {
		g := &Page{Title: "Error", Message: result.ToString()}
		t, _ := template.ParseFiles(config.TemplatePath + "error.html")
		p.LogPRM(result.ToString(), prm.LOG_DEBUG)

		t.Execute(w, g)
		return
	}
	// If we have an otp show terms and conditions otherwise dont
	if len(data["otp"]) > 0 {
		g := &Page{Title: "Terms and conditions", Message: result.ToString(), Wuffer: data["wuffer"], Puffer: data["puffer"], Tuffer: data["tuffer"]}
		t, _ := template.ParseFiles(config.TemplatePath + "terms.html")
		t.Execute(w, g)
	} else {

		// Make the attempt to change the data now
		result, data = p.ProcessSkipped(r)

		if result.Message != prm.SuccessFinished {
			g := &Page{Title: "Error", Message: result.ToString()}
			t, _ := template.ParseFiles(config.TemplatePath + "error.html")
			p.LogPRM(result.ToString(), prm.LOG_ERROR)

			t.Execute(w, g)
		} else {

			// Log the successful user

			username := data["username"]
			p.LogPRM(username+" with IP "+r.RemoteAddr+" successfully set password", prm.LOG_INFO)

			g := &Page{Title: "Success", Message: result.ToString()}
			t, _ := template.ParseFiles(config.TemplatePath + "success.html")
			t.Execute(w, g)
		}
	}
}

// processTerms processes the terms and conditions acceptance page
// Essentially the same as above for now
func processTerms(w http.ResponseWriter, r *http.Request, p *prm.PRM) {

	config := p.Config
	result, data := p.ProcessTerms(r)
	if result.Message != prm.SuccessFinished {
		g := &Page{Title: "Error", Message: result.ToString()}
		t, _ := template.ParseFiles(config.TemplatePath + "error.html")
		p.LogPRM(result.ToString(), prm.LOG_DEBUG)

		t.Execute(w, g)
		return
	}

	// Log the successful user
	username := data["username"]
	p.LogPRM(username+" with IP "+r.RemoteAddr+" successfully set password", prm.LOG_INFO)

	g := &Page{Title: "Success", Message: result.ToString()}
	t, _ := template.ParseFiles(config.TemplatePath + "success.html")
	t.Execute(w, g)
	return

}

// processPassword processes a password in an ajax style. It is here for the
// cracklib check which is sent by jquery everytime the user enters a new password
func processPassword(w http.ResponseWriter, r *http.Request, p *prm.PRM) {
	r.ParseForm()
	password := strings.Join(r.Form["password"], "")
	fmt.Fprintf(w, prm.TestPassword(password))
}

// ServeHTTP deals with the URLs, providing the correct response given the URL
func (s *FastCGIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//http.HandleFunc("/", renderIndex)
	//http.HandleFunc("/change", processForm)

	if r.URL.Path == "/" {
		renderIndex(w, r, &s.PRMHandler)
		return
	} else if r.URL.Path == "/change" {
		processForm(w, r, &s.PRMHandler)
		return
	} else if r.URL.Path == "/accept" {
		processTerms(w, r, &s.PRMHandler)
		return
	} else if r.URL.Path == "/check" {
		processPassword(w, r, &s.PRMHandler)
		return
	}

	s.PRMHandler.LogPRM("Path not found: "+r.URL.Path, prm.LOG_DEBUG)
	http.NotFound(w, r)
	return
}

// ReadConfig reads in the YAML config file, setting the PRMConfig used
// throughout the system
func ReadConfig(p *prm.PRM) {

	yamlConfig := prm.YamlConfig{}
	// Attempt an env read then /usr/local/secret first then locally
	filename, _ := filepath.Abs(os.Getenv("UPRM_CONFIG_FILE"))
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		panic(err)
	}

	// Now convert
	p.Config = new(prm.PRMConfig)
	config := p.Config

	// Potentially we should set all to defaults and check for missing
	config.TemplatePath = yamlConfig.TemplatePath
	config.ListenAddress = yamlConfig.ListenAddress
	config.LDAPHost = yamlConfig.LDAPHost
	config.LDAPPort = yamlConfig.LDAPPort
	config.BindPassword = yamlConfig.BindPassword
	config.CertFilePath = yamlConfig.CertFilePath
	config.BaseDN = yamlConfig.BaseDN
	config.BindDN = yamlConfig.BindDN
	config.EmailMsg = yamlConfig.EmailMsg
	config.EmailSub = yamlConfig.EmailSub
	config.LogLevel = prm.LOG_ERROR
	config.Uffer = yamlConfig.Uffer
	config.LDAPInsecureSkipVerify = yamlConfig.LDAPInsecureSkipVerify
	config.PasswordModifyLDAP = yamlConfig.PasswordModifyLDAP
	config.ORGFieldLDAP = yamlConfig.ORGFieldLDAP
	config.UserFieldLDAP = yamlConfig.UserFieldLDAP

	if yamlConfig.LogLevel == "DEBUG" {
		config.LogLevel = prm.LOG_DEBUG
	}

	if yamlConfig.LogLevel == "INFO" {
		config.LogLevel = prm.LOG_INFO
	}

	if yamlConfig.LogLevel == "WARN" {
		config.LogLevel = prm.LOG_DEBUG
	}

	if yamlConfig.LogLevel == "ERROR" {
		config.LogLevel = prm.LOG_ERROR
	}

	//fmt.Printf("Listening on: %#v\n", config.ListenAddress)

	p.Config = config

	p.LogPRM("Path to Templates: "+config.TemplatePath, prm.LOG_INFO)
	p.LogPRM("Path to CertFile: "+config.CertFilePath, prm.LOG_INFO)
	p.LogPRM("Log level: "+prm.LogLevelToString(config.LogLevel), prm.LOG_INFO)
	p.LogPRM("Listen address: "+config.ListenAddress, prm.LOG_DEBUG)
	p.LogPRM("LDAP Host address: "+config.LDAPHost, prm.LOG_DEBUG)
	p.LogPRM("LDAP Port: "+string(config.LDAPPort), prm.LOG_DEBUG)
	p.LogPRM("BaseDN: "+config.BaseDN, prm.LOG_DEBUG)
	p.LogPRM("BindDN: "+config.BindDN, prm.LOG_DEBUG)
	p.LogPRM("PasswordModifyLDAP: "+config.PasswordModifyLDAP, prm.LOG_DEBUG)
	p.LogPRM("ORGFieldLDAP: "+config.ORGFieldLDAP, prm.LOG_DEBUG)
	p.LogPRM("UserFieldLDAP: "+config.UserFieldLDAP, prm.LOG_DEBUG)

	if config.LDAPInsecureSkipVerify {
		p.LogPRM("LDAP insecure skip verify: true", prm.LOG_DEBUG)
	} else {
		p.LogPRM("LDAP insecure skip verify: false", prm.LOG_DEBUG)
	}
	p.LogPRM("uffer: "+config.Uffer, prm.LOG_DEBUG)
	p.LogPRM("Email message: "+config.EmailMsg, prm.LOG_DEBUG)
	p.LogPRM("Email subject: "+config.EmailSub, prm.LOG_DEBUG)

}

func main() {
	// Test for the version flag
	var ip = flag.Bool("v", false, "display the version and quit")
	flag.Parse()
	if *ip == true {
		fmt.Println(prm.GetVersionString())
		return
	}

	fmt.Println("Welcome to the Password Manager - The Next Generation!")
	prmHandler := new(prm.PRM)
	ReadConfig(prmHandler)
	//listener, err := net.Listen("unix", config.ListenAddress)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer listener.Close()

	// create a server object with the config and PRM handler
	srv := new(FastCGIServer)

	srv.PRMHandler = *prmHandler

	err := fcgi.Serve(nil, srv)

	if err != nil {
		fmt.Println("error on serve")
		log.Fatal("[prm:error]", err)
	}
}
