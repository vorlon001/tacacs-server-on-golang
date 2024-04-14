package main

import (
        "regexp"
	"github.com/vorlon001/tacacs-server-on-golang/pkg/module/ldap-client"
	configure "github.com/vorlon001/tacacs-server-on-golang/pkg/module/config"
	"log"
	"os"
	"fmt"
)

func auth(login string, password string,config *configure.Config) interface{} {

        if len(config.LDAP.Base)==0 {
                log.Printf("Error: not found config.LDAP.Base")
                os.Exit(1)
        } else if len(config.LDAP.Host)==0 {
                log.Printf("Error: not found config.LDAP.Host")
                os.Exit(1)
        } else if config.LDAP.Port==0 {
                config.LDAP.Port = 389
                log.Printf("Error: not found config.LDAP.Port")
        } else if len(config.LDAP.BindDN)==0 {
                log.Printf("Error: not found config.LDAP.BindDN")
                os.Exit(1)
        } else if len(config.LDAP.BindPassword)==0 {
                log.Printf("Error: not found config.LDAP.BindPassword")
                os.Exit(1)
        } else if len(config.LDAP.UserFilter)==0 {
                log.Printf("Error: not found config.LDAP.UserFilter")
                os.Exit(1)
        } else if len(config.LDAP.GroupFilter)==0 {
                log.Printf("Error: not found config.LDAP.GroupFilter")
                os.Exit(1)
        }

        if config.LOG.DEBUG.ENABLE==true {
            log.Println(">> auth()",config)
        }

        client := &ldap.LDAPClient{
                Base:         config.LDAP.Base,
                Host:         config.LDAP.Host,
                Port:         config.LDAP.Port,
                UseSSL:       config.LDAP.UseSSL,
                BindDN:       config.LDAP.BindDN,
                BindPassword: config.LDAP.BindPassword,
                UserFilter:   config.LDAP.UserFilter,
                GroupFilter:  config.LDAP.GroupFilter,
                Attributes:   []string{"cn", "mail","memberOf","mobile","postOfficeBox","ipPhone","mail","userPrincipalName","displayName"},
        }
        defer client.Close()

        var status_group interface{}
        var GROUP_NAME  string

        status_group = interface{}(nil)
        GROUP_NAME  = "-"



        if config.LOG.DEBUG.ENABLE==true {
                log.Println("VERIFY2:", status_group , GROUP_NAME)
        }

        if config.LOG.DEBUG.ENABLE==true {
                log.Println("INIT LOCAL AUTH:", login, password,config.USER)
        }

        if len(config.USER)>0 {
                for _,w := range config.USER {
                        if len(w.Login)>0 && len(w.Password)>0{
                                if w.Login==login && w.Password==password {
                                    if config.LOG.DEBUG.ENABLE==true {
                                        log.Printf("config.USER -> %#v \n",w);
                                    }
                                    status_group = w;
                                    if config.LOG.DEBUG.ENABLE==true {
                                        log.Println("VERIFY3:", login , password ,status_group);
                                    }
                                    return status_group;

                                }
                        } else {
                            if config.LOG.DEBUG.ENABLE==true {
                                log.Printf("config.USER -> %#v not found login or password\n",w);
                            }
                        }
                }
        }

        if config.LOG.DEBUG.ENABLE==true {
                log.Println("INIT LDAP AUTH:", login)
        }

        ok, user, err := client.Authenticate(login,password)
        if err != nil {
                if config.LOG.DEBUG.ENABLE==true {
                    log.Println("Error authenticating user %s: %+v", "username", err)
                }
		return status_group
        }
        if !ok {
                if config.LOG.DEBUG.ENABLE==true {
                    log.Println("Authenticating failed for user %s", "username")
                }
                return status_group
        }

	if config.LOG.DEBUG.ENABLE==true {
        	log.Printf("displayName: %#v\n", user["displayName"])
	        log.Printf("cn: %#v\n", user["cn"])
	        log.Printf("userPrincipalName: %#v\n", user["userPrincipalName"])
	        log.Printf("ipPhone: %#v\n", user["ipPhone"])
	        log.Printf("mobile: %#v\n", user["mobile"])
	}

	verify := func  ( group string, base string , URL string ) (bool , string) {
	        r := regexp.MustCompile(fmt.Sprintf("CN=(%s),(%s)",group,base))
	        var status bool;
	        var GROUP  string;
	        status = false;
	        GROUP  = "";
	        res := r.FindStringSubmatch(URL)
	        if len(res)==3 {
        	        status = true
	                GROUP  = res[1]
                        if config.LOG.DEBUG.ENABLE==true {
	                    log.Printf(res[1],res[2],len(res))
                        }
	        }
	        return  status , GROUP
	}

	for _,v := range config.ACCESS {
                if len(v.Base)==0 {
                      if config.LOG.DEBUG.ENABLE==true {
                          log.Printf("Error: not found in config.ACCESS BASE in %v",v)
                      }
                      os.Exit(1)
                }
                if len(v.Group)==0 {
                      if config.LOG.DEBUG.ENABLE==true {
                          log.Printf("Error: not found in config.ACCESS Group in %v",v)
                      }
                      os.Exit(1)
                }
		if config.LOG.DEBUG.ENABLE==true {
			log.Println("VERIFY0:",  v)
        	        log.Printf("VERIFY0: %T %#v",  v ,v)
                	log.Printf("VERIFY0: %T %#v",  v.Base ,v.Base)
	                log.Printf("VERIFY0: %T %#v",  v.Group ,v.Group)
		}

		for _,w := range user["memberOf"] {
			if status, GROUP := verify(v.Group,v.Base,w); status {
				v.USERInfo = user
				status_group = v
				GROUP_NAME  = GROUP
                                if config.LOG.DEBUG.ENABLE==true {
				    log.Println("VERIFY1:", status_group , GROUP_NAME)
                                }
				return status_group
			}
		}
	}

        return status_group;
}



