package config

import (
        "fmt"
        "os"
        "gopkg.in/yaml.v2"
)

type Log struct {
                DEBUG struct {
                        ENABLE bool `yaml:"ENABLE"`
                } `yaml:"DEBUG"`
                SYSLOG struct {
                        ENABLE bool   `yaml:"ENABLE"`
                        IP     string `yaml:"IP"`
                        PORT   int    `yaml:"PORT"`
                } `yaml:"SYSLOG"`
                FILE struct {
                        NAME   string `yaml:"NAME"`
                } `yaml:"FILE"`
        };

type Banner struct {
                LoginBanner    string `yaml:"login_banner"`
                PasswordBanner string `yaml:"password_banner"`
                Banner         string `yaml:"banner"`
                BannerAccept   string `yaml:"banner_accept"`
                BannerReject   string `yaml:"banner_reject"`
        };

type Ldap struct {
                Base         string   `yaml:"Base"`
                Host         string   `yaml:"Host"`
                Port         int      `yaml:"Port"`
                UseSSL       bool     `yaml:"UseSSL"`
                BindDN       string   `yaml:"BindDN"`
                BindPassword string   `yaml:"BindPassword"`
                UserFilter   string   `yaml:"UserFilter"`
                GroupFilter  string   `yaml:"GroupFilter"`
                Attributes   []string `yaml:"Attributes"`
        };

type UserLDAP struct {
                Base       	string   `yaml:"Base"`
                Group      	string   `yaml:"Group"`
                PrivLvl   	int      `yaml:"priv-lvl"`
                Description     string   `yaml:"description"`
                IPAccess        []string `yaml:"IPAccess"`
                PERMIT     	[]string `yaml:"PERMIT,omitempty"`
		USERInfo	map[string][]string
        };

type Device struct {
                Network string `yaml:"network"`
                Token   string `yaml:"token"`
        };

type User struct {
                Login    	string   `yaml:"login"`
                Password 	string   `yaml:"password"`
                PrivLvl  	int      `yaml:"priv-lvl"`
                Description  	string   `yaml:"description"`
                IPAccess 	[]string `yaml:"IPAccess"`
                PERMIT   	[]string `yaml:"PERMIT"`
        };

type Config struct {
	PID  	string `yaml:"PID"`
	BIND 	string `yaml:"BIND"`
	PORT 	int    `yaml:"PORT"`
	LOG  	Log `yaml:"LOG"`
	LDAP   	Ldap `yaml:"LDAP"`
	ACCESS 	[]UserLDAP `yaml:"ACCESS"`
	Banner 	Banner `yaml:"banner"`
	DEVICE 	[] Device`yaml:"DEVICE"`
	USER	[]User `yaml:"USER"`
};

func Var_dump(expression ...interface{} ) {
	fmt.Println(fmt.Sprintf("%#v", expression))
}

func NewConfig(configPath string) (*Config, error) {
    // Create config structure
    config := &Config{} 

    // Open config file
    file, err := os.Open(configPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Init new YAML decode
    d := yaml.NewDecoder(file)

    // Start YAML decoding from file
    if err := d.Decode(&config); err != nil {
        return nil, err
    }

    return config, nil
}
