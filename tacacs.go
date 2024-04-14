package main

import (
	"flag"
	"github.com/vorlon001/tacacs-server-on-golang/pkg/module/go-daemon"
	configure "github.com/vorlon001/tacacs-server-on-golang/pkg/module/config"
	"log"
	"os"
	"fmt"
	"syscall"
)


var (

	config = &configure.Config{}
        stop = make(chan struct{})
        done = make(chan struct{})
        reload = make(chan configure.Config) 
	PidFileName = "tacacs.pid"
	LogFileName = "tacacs.log"
	PID_NAME    = "[tacacs+ daemon]"
	signal = flag.String("h", "", `Send signal to the daemon:
  quit — graceful shutdown
  stop — fast shutdown
  reload — reloading the configuration file
  testconfig - test config
`)
)

var Version string

func VersionBuild(version string) {
        fmt.Printf("https://github.com/vorlon001, (C) Vorlon001\n")
        fmt.Printf("Tacacs+ Server version:%s\n",version)
}



func main() {

	VersionBuild(Version)

	config_path  := "tacacs.yml"
	flag.Parse()
        fmt.Println("config:", config_path)
	fmt.Println("ARG signal:", flag.Args())

	if len(flag.Args())==1 {
               	if flag.Args()[0]=="testconfig" {
			cfg, _ := configure.NewConfig(config_path)
		        log.Println("CONFIG:\n%T %v\n\n", cfg, cfg )
                        os.Exit(0)
		} else if flag.Args()[0]=="stop" {
		        cfg, config_err := configure.NewConfig(config_path)
		        if config_err != nil {
		            log.Fatalf("ERROR LOAD CONFIG: %s", config_err.Error())
		            os.Exit(1)
		        }
		        var config = cfg
		        PidFileName = config.PID;
			fmt.Println(get_pig_daemon(PidFileName));
			fmt.Println(syscall.Kill(get_pig_daemon(PidFileName), syscall.SIGQUIT))
			os.Exit(0)
		} else  if flag.Args()[0]=="quit" {
		        cfg, config_err := configure.NewConfig(config_path)
		        if config_err != nil {
		            log.Fatalf("ERROR LOAD CONFIG: %s", config_err.Error())
		            os.Exit(1)
		        }
		        var config = cfg
		        PidFileName = config.PID;
                        fmt.Println(get_pig_daemon(PidFileName));
                        fmt.Println(syscall.Kill(get_pig_daemon(PidFileName), syscall.SIGQUIT))
                        os.Exit(0)
                } else  if flag.Args()[0]=="reload" {

		        cfg, config_err := configure.NewConfig(config_path)
		        if config_err != nil {
		            log.Fatalf("ERROR LOAD CONFIG: %s", config_err.Error())
		            os.Exit(1)
		        }
		        var config = cfg

		        PidFileName = config.PID;

                        fmt.Println(get_pig_daemon(PidFileName));
                        fmt.Println(syscall.Kill(get_pig_daemon(PidFileName), syscall.SIGHUP))
                        os.Exit(0)
                }
	}

        cfg, config_err := configure.NewConfig(config_path)
        if config_err != nil {
            log.Fatalf("ERROR LOAD CONFIG: %s", config_err.Error())
            os.Exit(1)
        }
	var config = cfg

	if len(config.PID)>0 {
        	PidFileName = config.PID;
	}
	if len(config.LOG.FILE.NAME)>0 {
	        LogFileName = config.LOG.FILE.NAME;
	}

	daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
        daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGKILL, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)

	cntxt := &daemon.Context{
		PidFileName: PidFileName,
		PidFilePerm: 0644,
		LogFileName: LogFileName,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{PID_NAME},
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalf("Unable send signal to the daemon: %s", err.Error())
                        os.Exit(1)
		}
		daemon.SendCommands(d)
		return
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Println("- - - - - - - - - - - - - - -")
	log.Println("daemon started")

	go worker(config)

	err = daemon.ServeSignals()
	if err != nil {
		log.Printf("Error: %s", err.Error())
                os.Exit(1)
	}

	log.Println("daemon terminated")
}

