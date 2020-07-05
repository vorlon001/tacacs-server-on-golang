package main

import (
	"flag"
	"net"
	"module/go-daemon"
	"module/tacplus"
	configure "module/config"
	"log"
	"log/syslog"
	"os"
	"io"
	"fmt"
	"syscall"
	"time"
	"strconv"
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

func main() {

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

func worker(config *configure.Config) {

  if config.LOG.SYSLOG.ENABLE==true {
      if config.LOG.SYSLOG.PORT==0 {
         config.LOG.SYSLOG.PORT = 514
         log.Printf("Error: config.LOG.SYSLOG.PORT not yet set")
      }
      if len(config.LOG.SYSLOG.IP)==0 {
         log.Printf("Error: config.LOG.SYSLOG.IP not yet set")
         sysLog, err := syslog.New(syslog.LOG_WARNING|syslog.LOG_DAEMON, PID_NAME)
         if err != nil {
            log.Println("ERROR init SYSLOG",err)
            os.Exit(1)
         }
         a := log.Writer()
         multi := io.MultiWriter(sysLog,a)
         log.SetOutput(multi)
      } else {
         sysLog, err := syslog.Dial("udp", "localhost:514",
                       syslog.LOG_WARNING|syslog.LOG_DAEMON, PID_NAME)
         if err != nil {
            log.Println("ERROR init SYSLOG",err)
            os.Exit(1)
         }
         a := log.Writer()
         multi := io.MultiWriter(sysLog,a)
         log.SetOutput(multi)
      }
  }


  if config.LOG.DEBUG.ENABLE==true {
        log.SetPrefix( "[DEBUG MODE]" )
        log.SetFlags(log.Ldate|log.Ltime|log.Llongfile)
   } else {
        log.SetPrefix( "[AAA]" )
        log.SetFlags(log.Lshortfile)
  }


  if len(config.BIND)==0 {
	config.BIND = "0.0.0.0"
        log.Printf("Error: config.BIND not yet set")
  }

  if config.PORT==0 {
        config.PORT = 49
        log.Printf("Error: config.PORT not yet set")
  }



  sock, err := net.Listen("tcp", fmt.Sprintf("%v:%v",config.BIND,config.PORT)) //"0.0.0.0:49")
  if err != nil {
        log.Println("Can't listen... (%s)\n",fmt.Sprintf("%s:%s",config.BIND,config.PORT))
        os.Exit(3);
  }

  if config.LOG.DEBUG.ENABLE==true {
	log.Println("CONNECT TO",fmt.Sprintf("%v:%v",config.BIND,config.PORT),config.BIND,config.PORT)
	log.Println(config.DEVICE);
  }

  if len(config.DEVICE)==0 {
        log.Printf("Error: config.DEVICE not yet set")
	config.DEVICE[0] = configure.Device{  Network: "0.0.0.0", Token: "private" };
  }

  tU := tacasUserCached{}
  tH := &tacasHandler{Config: config, UserCached: tU}
  tC := tacplus.ConnConfig{ DEVICE: config.DEVICE, Mux: true, }

  handler := tacplus.ServerConnHandler{
    Handler: tH, //tacasHandler{Config: config},
    ConnConfig: tC,
  }
  server := &tacplus.Server{
    ServeConn: func (nc net.Conn) {
      handler.Serve(nc)
    },
  }

// Listen & Serve
//  signalChan := make(chan os.Signal, 1)
//  signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
  errChan := make(chan error)

  go func() {
    log.Println("Waiting for request...")
    err := server.Serve(sock)
    if err != nil {
      errChan <- err
    }
  }()

LOOP:
	for {
		time.Sleep(time.Second) // this is work to be done by worker.
		select {
		case <-stop:
			log.Println("Stopping server")
			sock.Close()
			break LOOP
		case err := <- errChan:
			log.Println("[error] %v", err.Error())
                case a:=<-reload:

                        if config.LOG.DEBUG.ENABLE==true {
                           log.Println("POINT 1 worker() reload OLD CONFIG ",tH.Config);
                        }
                        config = &a
                        handler.ConnConfig.DEVICE = config.DEVICE
                        tH.Config = &a
                        config.PID = "RELOADED"

                        if config.LOG.DEBUG.ENABLE==true {
                           log.Println("POINT 2 worker() reload NEW CONFIG ",tH.Config);
                           log.Println("POINT 3 worker() reload ",config,os.Getpid(),get_pig_daemon(PidFileName))
                           log.Println("POINT 4 worker() reload ",a);
	                   log.Println("POINT 5 worker() reload  tH.UserCached ",tH.UserCached);
                           log.Println("POINT 6 worker() reload ",config,os.Getpid(),get_pig_daemon(PidFileName))
                        }
		default:
		}
	}
	done <- struct{}{}
}

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
        log.Println(get_pig_daemon(PidFileName));
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {

        log.Println(">>",config,os.Getpid(),get_pig_daemon(PidFileName))
	config.PID = "RELOADED"
        log.Println(">>",config,os.Getpid(),get_pig_daemon(PidFileName))

        cfg, config_err := configure.NewConfig("./daemon_6.yml")
        if config_err != nil {
            log.Fatalf("ERROR LOAD CONFIG: %s", config_err.Error())
        } else {
            reload <- *cfg
        }
	log.Println("configuration reloaded")
	return nil
}


func get_pig_daemon(a string) int {
    f, err := os.Open(a)
    if err != nil {
        log.Println("File PID not found", a);
        os.Exit(1)
    }
    defer f.Close()

    b1 := make([]byte, 5)
    n1, err := f.Read(b1)
    check(err)
    i, err := strconv.Atoi(string(b1[:n1]))

    check(err)

    if config.LOG.DEBUG.ENABLE==true {
            log.Println("get_pig_daemon() %s %v\n", a, i )
    }

    return i
}


func check(e error) {
    if e != nil {
        panic(e)
    }
}

