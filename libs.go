package main

import (
	"net"
        "regexp"
	"fmt"
	"time"
	"strconv"
	"log"
)

func Paring_Args_Request(Args []string) map[string]string {

	var parsing  = func (re *regexp.Regexp, Args []string) map[string]string {
		s := make(map[string]string,len(Args ))
		for _,v := range Args  {
			if a := re.FindAllString(v, -1);len(a)==2 {
				if  a[0]=="start_time" {
					if s, err := strconv.ParseInt(a[1], 10, 64); err == nil {
						a[1] = fmt.Sprintf("%s\n", time.Unix(s, 0))
					}
				}
				s[string(a[0])]=string(a[1])
			}
		}
		return s
	}
    var re = regexp.MustCompile(`[\w\-\<\>\s]+`)
    return parsing (re, Args)
}


func Verify_Ip_Access( IPAccess  []string, clientip string ) bool  {

				if len(IPAccess)==0 {
					return true
				} else {
					ip := net.ParseIP(clientip)
					for _, v:= range IPAccess {
						_, subnet, err := net.ParseCIDR( v )
						if err != nil {
							log.Println(err)
						}
						if subnet.Contains(ip) {
							return true
						}
					}
				}
                return false
	}

func Verify_Cmd(permit []string, cmd string) bool {
	var Find_Cmd =func(rgx string, cmd string) bool{
		var re = regexp.MustCompile(rgx)
		if a := re.FindAllString(cmd, -1);len(a)==1 {
			return true
	    	}
		return false
	}
	var status = false
	for _,v:= range permit  {
    		if Find_Cmd(v,cmd)==true {
			status = true
		}
	}
	return status
}
