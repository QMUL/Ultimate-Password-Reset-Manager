# ultimate-password-reset-manager

The password reset manager redux! This is not your grandfather's password manager! :D In all seriousness, this is a modern version of the classic Perl script. Features and changes include:

- Written in Golang - lots of things included, plus its compiled
- Communicates via fcgi protocol, widely supported
- Data, Logic and Presentation are seperated
- Simple templates for changing the 'look' without recompiling
- Better integration with cracklib
- Use of bootstrap and jquery for a cleaner interface, complete with live feedback


## Building

To build this password manager, you will need the Go language installed and setup as per the instructions [on the Go webpage](https://golang.org/doc/install)

There are two dependencies that can be installed as follows:

    go get gopkg.in/yaml.v2
    go get gopkg.in/ldap.v2

### Requirements
Before you attempt to build this application you need to install the following packages:

#### RHEL/CenOS

```bash
yum install golang golang-cover golang-godoc cracklib-devel rpmbuild
```

Please note that ```golang-cover``` and ```golang-godoc``` come from the epel repository.

#### Debian/Ubuntu

### Building with Go

To build, do the following.

1. Create the necessary directory structure with ${GOPATH} defining the root directory of that structure:

    <pre>
    mkdir ultimate-password-reset-manager
    cd !$
    mkdir -p {bin,pkg,src}
    git clone &lt;prm github repository url&gt; src/pass.hpc.qmul.ac.uk
    </pre>

2. Checkout the repository into your ${GOPATH}/src directory

    <pre>
    cd ${GOPATH}
    git clone &lt;prm github repository url&gt; src/pass.hpc.qmul.ac.uk
    </pre>

3. Run the following command 

    <pre>
    go install pass.hpc.qmul.ac.uk/prmserver
    </pre>

### Building with CMake

CMake is also an option. Create a directory for building as outlined above, then type a command such as this

    cmake <path to pass.hpc.qmul.ac.uk> -DCMAKE_INSTALL_PREFIX=../install

then type:

* make 
* make tests
* make install
* make package
 
where:
* ```<path to pass.hpc.qmul.ac.uk>``` should point at the top level pass directory that contains the main CMakeLists.txt
* ```-DCMAKE_INSTALL_PREFIX=``` should point to a directory which will contain the deployable structure

Make package will generate an rpm package. Thw package version should correspond to the git tag for the package. At this point in time the package version is hardcoded in the _cmake/PrmVersion.cmake_ file and should be changed whenever new tag is pushed. You can override the version information by passing an appropriate macros to the cmake invocation:

* VERSION_MAJOR=X
* VERSION_MINOR=Y
* VERSION_PATCH=Z
* VERSION_RELEASE=V

### cracklib

If the deictionaries are not provided via the package manager one needs to generate the cracklib dictionary from another dictionary, usually /usr/share/dict/words, although there are others

    sudo /usr/sbin/create-cracklib-dict /usr/share/dict/words

This creates the default dictionary location in

    /usr/share/cracklib/pw_dict

Centos provides default dictionaries with its cracklib package (via cracklib-dicts package), so generation may not be necessary. The location is the same in both cases.

### Building the documentation

The documentation is built on godoc - [https://blog.golang.org/godoc-documenting-go-code](https://blog.golang.org/godoc-documenting-go-code).

To build the documentation one would use godoc once all the PATH variables and such have been set you can run

    godoc pass.hpc.qmul.ac.uk/prm 

Because the main package is a binary, we only get the command documentation. We can do that via

    godoc cmd/pass.hpc.qmul.ac.uk/prmserver

To generate nice html that you can pipe out somewhere use the switch *-html*

## Deployment under Apache

### Pre-requisites
You need apache, mod-fcgid and cracklib installed:

    yum install httpd mod-fcgid cracklib

For development you also need cracklib-devel package:

    yum install cracklib-devel

### Settings

Sample setting can be found in the _passwordmanager/config.yml.template_ file. This file location is passed to the application via the UPRM_CONFIG_FILE environment variable

In apache this can be set using the following directives [https://httpd.apache.org/mod_fcgid/mod/mod_fcgid.html#fcgidinitialenv https://httpd.apache.org/mod_fcgid/mod/mod_fcgid.html#fcgidinitialenv]

Currently the following options are configurable (displayed with their default values):

    ---
    templatepath: ../templates/
    ldaphost: ldaphost-test
    ldapport: 389
    bindpassword: prm
    ldapinsecureskipverify: true
    binddn: <your bind dn>
    basedn: <your base dn>
    certfilepath: <cert file path>
    uffer: DD23AA67833BCDEF
    loglevel: DEBUG
    passwordmodifyldap: uid=%v,ou=People
    userfieldldap: uid
    orgfieldldap: ou=People   
    emailsub: "Email Subject"
    emailmsg: | 
     Dear %NAME% 
       Here is an email   
    ---

The value *loglevel* takes the following options in order of decreasing verbosity:

    DEBUG
    INFO
    WARN
    ERROR 

Note that the *uffer* value is an AES key and needs to be 16 characters long (or any acceptable key length for the golang implementation of the aes cipher). This *must* be changed to a random string on deployment. If it is not set users will be able to change their passwords bypassing the terms and conditions and aes errors will appear in the logs.

The *ldap* fields are set for our local install. You can alter these for your ldap install. *passwordmodifyldap* refers to the search fields for finding the user, whose password you wish to modify. *userfieldldap* refers to the name of the user identification field and the *orgfieldldap* refers to the organisation you are looking within.

### Using standard io and apache controls

This method is very similar to classic CGI scripting, where Apache controls the launching of the executable. Note that this current master branch supports this method only. You'd need to adjust the listener code and recompile in order to use the mod_proxy method below.

To setup Apache, you need a config a little like this:


    ScriptAlias /passwordmanager/ "/srv/www/password/passwordmanager"
    
    Alias /static /srv/www/password/static
    Alias /check /srv/www/password/passwordmanager/prm_server.fcgi/check
    Alias /change /srv/www/password/passwordmanager/prm_server.fcgi/change
    Alias /accept /srv/www/password/passwordmanager/prm_server.fcgi/accept

    Alias / /srv/www/password/passwordmanager/prm_server.fcgi

    #Set the location of the config file
    FcgidInitialEnv UPRM_CONFIG_FILE /vagrant_data/config.yml

    <Directory "/srv/html/passwordmanager">
        SetHandler fcgid-script
        AllowOverride None
        Options +ExecCGI
        Allow from all
    </Directory>

    <location /static>
       Options None
        Order deny,allow
        Allow from all
    </location>

    <location /check>
        Order deny,allow
        Allow from all
    </location>


    <location /change>
        Order deny,allow
        Allow from all
    </location>

    <location /accept>
        Order deny,allow
        Allow from all
    </location>


This config is probably overkill but it works on the Vagrant PRM machine. Likely, someone who knows more about Apache can come up with a better one. Basically, there are only 4 URLS this system needs, with the static path being for all the images, css and the like.

Note that the executable is renamed to **prm_server.fcgi** - it is unclear whether or not this is just convention, or Apache insists on such a suffix.

### SELinux
In order to deploy this in SELinux enabled environment the following keys have to be set (using _setsebool -P key on_):

* httpd\_can\_network\_connect
* httpd\_enable\_cgi

Also the following custom policy (_prm\_cracklib.te_) has to be enabled for the cracklib to work:

    module prm_cracklib 1.0;
    
    require {
        type crack_db_t;
        type httpd_sys_script_t;
        type sysctl_net_t;
        class dir { search };
        class file { read getattr open };
    }
    
    #============= httpd_t ==============
    allow httpd_sys_script_t crack_db_t:dir search;
    allow httpd_sys_script_t crack_db_t:file { read getattr open };
    
    allow httpd_sys_script_t sysctl_net_t:dir search;
    allow httpd_sys_script_t sysctl_net_t:file { read open };

To add this module execute the following:

* checkmodule -M -m -o prm\_cracklib.mod prm\_cracklib.te
* semodule\_package -o prm\_cracklib.pp -m prm\_cracklib.mod
* semodule -i prm\_cracklib.pp


### Using mod_proxy and fcgi

Please note that this step is generally not required and is only included for the completensss.

There are several ways to do this, depending on other factors of your setup. Out of the box, apache comes with two modules: **mod_proxy** and **mod_proxy_fcgi**. You can see the setup here

http://httpd.apache.org/docs/2.4/mod/mod_proxy_fcgi.html

The config looks like this:

    Alias /static /srv/www//password/static

    <location /static>
        Order deny,allow
        Allow from all
    </location>

    ProxyPassMatch ^/$ fcgi://127.0.0.1:9001
    ProxyPassMatch ^/check$ fcgi://127.0.0.1:9001/check
    ProxyPassMatch ^/change$ fcgi://127.0.0.1:9001/change
    ProxyPassMatch ^/accept$ fcgi://127.0.0.1:9001/accept
 
To run the program, enter the base directory (**/vagrant_prm_data** on the vagrant box) and run
    cd passwordmanager
    ./prm_server

Where you run from affects where the templates are loaded from. Check the **config.yml** for the TemplatePath parameter. At present, its relative to where the executable is run from.                                        
