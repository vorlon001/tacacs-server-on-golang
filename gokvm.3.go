package main

import (
	"fmt"
	"log"
	"net"
	"time"
	"strconv"
	"github.com/digitalocean/go-libvirt"
)

var DomainState = map[int]string{
	0: "DomainNostate",
	1: "DomainRunning",
	2: "DomainBlocked",
	3: "DomainPaused",
	4: "DomainShutdown",
	5: "DomainShutoff",
	6: "DomainCrashed",
	7: "DomainPmsuspended" }

func main() {
	// This dials libvirt on the local machine, but you can substitute the first
	// two parameters with "tcp", "<ip address>:<port>" to connect to libvirt on
	// a remote machine.
	c, err := net.DialTimeout("unix", "/var/run/libvirt/libvirt-sock", 2*time.Second)
	if err != nil {
		log.Fatalf("failed to dial libvirt: %v", err)
	}

	l := libvirt.New(c)
	if err := l.Connect(); err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	v, err := l.Version()
	if err != nil {
		log.Fatalf("failed to retrieve libvirt version: %v", err)
	}
	fmt.Println("Version:", v)

	domains, err := l.Domains()
	if err != nil {
		log.Fatalf("failed to retrieve domains: %v", err)
	}

	fmt.Println("ID\tName\t\tUUID")
	fmt.Printf("--------------------------------------------------------\n")
	for _, d := range domains {
		fmt.Printf("%d\t%s\t%x\n\t%#v\n", d.ID, d.Name, d.UUID, d)
                a, _ := l.DomainState(d.Name)
                c, _ := strconv.Atoi(fmt.Sprintf("%d",a))
                fmt.Printf("\t\t %#v %s\n", a, DomainState[c] )
	}

	if err := l.Disconnect(); err != nil {
		log.Fatalf("failed to disconnect: %v", err)
	}
}
