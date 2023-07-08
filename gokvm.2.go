package main

https://habr.com/ru/companies/ruvds/articles/684300/

/*

./gokvm.2 --virtual-machine-state --id node190

*/
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"

//	"uncloudzone/libtars/tylibvirt"

	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
)

// VirState represents current lifecycle state of a machine
// Pending = VM was just created and there is no state yet
// Running = VM is running
// Blocked = Blocked on resource
// Paused = VM is paused
// Shutdown = VM is being shut down
// Shutoff = VM is shut off
// Crashed = Most likely VM crashed on startup cause something is missing.
// Hybernating = Virtual Machine is hybernating usually due to guest machine request
// TODO:
type VirtState string

const (
	VirtStatePending     = VirtState("Pending")     // VM was just created and there is no state yet
	VirtStateRunning     = VirtState("Running")     // VM is running
	VirtStateBlocked     = VirtState("Blocked")     // VM Blocked on resource
	VirtStatePaused      = VirtState("Paused")      // VM is paused
	VirtStateShutdown    = VirtState("Shutdown")    // VM is being shut down
	VirtStateShutoff     = VirtState("Shutoff")     // VM is shut off
	VirtStateCrashed     = VirtState("Crashed")     // Most likely VM crashed on startup cause something is missing.
	VirtStateHybernating = VirtState("Hybernating") // VM is hybernating usually due to guest machine request
)

type VirtualMachineStatus string

const (
	VirtualMachineStatusDeleted      = VirtualMachineStatus("deleted")
	VirtualMachineStatusCreated      = VirtualMachineStatus("created")
	VirtualMachineStatusReady        = VirtualMachineStatus("ready")
	VirtualMachineStatusStarting     = VirtualMachineStatus("starting")
	VirtualMachineStatusImaging      = VirtualMachineStatus("imaging")
	VirtualMachineStatusRunning      = VirtualMachineStatus("running")
	VirtualMachineStatusOff          = VirtualMachineStatus("off")
	VirtualMachineStatusShuttingDown = VirtualMachineStatus("shutting_down")
)

// Versions - originally created for testing purposes, not actually something we would need.
// var libvirtVersion = *pflag.Bool("libvirt-version", false, "Returns result with version of libvirt populated")
// var virshVersion = *pflag.Bool("virsh-version", false, "Returns result with version of virsh populated")
// var tarsvirtVersion = *pflag.Bool("tarsvirt-version", false, "Returns result with version of tarsvirt populated")

// VirtualMachine commands
var virtualMachineState = pflag.Bool("virtual-machine-state", false, "Returns result with a current machine state")
var virtualMachineSoftReboot = pflag.Bool("virtual-machine-soft-reboot", false, "reboots a machine gracefully, as chosen by hypervisor. Returns result with a current machine state")
var virtualMachineHardReboot = pflag.Bool("virtual-machine-hard-reboot", false, "sends a VM into hard-reset mode. This is damaging to all ongoing file operations. Returns result with a current machine state")
var virtualMachineShutdown = pflag.Bool("virtual-machine-shutdown", false, "gracefully shuts down the VM. Returns result with a current machine state")
var virtualMachineShutoff = pflag.Bool("virtual-machine-shutoff", false, "kills running VM. Equivalent to pulling a plug out of a computer. Returns result with a current machine state")
var virtualMachineStart = pflag.Bool("virtual-machine-start", false, "starts up a VM. Returns result with a current machine state")
var virtualMachinePause = pflag.Bool("virtual-machine-pause", false, "stops the execution of the VM. CPU is not used, but memory is still occupied. Returns result with a current machine state")
var virtualMachineResume = pflag.Bool("virtual-machine-resume", false, "called after Pause, to resume the invocation of the VM. Returns result with a current machine state")
var virtualMachineCreate = pflag.Bool("virtual-machine-create", false, "creates a new machine. Requires --xml-template parameter. Returns result with a current machine state")
var virtualMachineDelete = pflag.Bool("virtual-machine-delete", false, "deletes an existing machine.")

var id = pflag.String("id", "", "id of the machine to work with")
var xmlTemplate = pflag.String("xml-template", "", "path to an xml template file that describes a machine. See qemu docs on xml templates.")

var v *libvirt.Libvirt

type VirtualMachine struct {
	CPUCount uint16
	CPUTime  uint64
	MemoryBytes uint64
	MaxMemoryBytes uint64
        State VirtState
}

// TODO: cool things you can do with Domain, but do not know how to:
// virDomainInterfaceAddresses - gets data about an IP addresses on a current interfaces. Mega-tool.
// virDomainGetGuestInfo - full data about a config of the guest OS
// virDomainGetState - provides the data about an actual domain state. Why is it shutoff or hybernating. Requires copious amount of magic fuckery to find out the actual reason with multiplication and matrix transforms, but can be translated into a redable form.
func main() {

	pflag.Parse()

	virtinit()

	switch {
	case *virtualMachineState:
		VirtualMachineState(*id)
	case *virtualMachineSoftReboot:
		VirtualMachineSoftReboot(*id)
	case *virtualMachineHardReboot:
		VirtualMachineHardReboot(*id)
	case *virtualMachineShutdown:
		VirtualMachineShutdown(*id)
	case *virtualMachineShutoff:
		VirtualMachineShutoff(*id)
	case *virtualMachineStart:
		VirtualMachineStart(*id)
	case *virtualMachinePause:
		VirtualMachinePause(*id)
	case *virtualMachineResume:
		VirtualMachineResume(*id)
	case *virtualMachineCreate:
		VirtualMachineCreate(*xmlTemplate)
	case *virtualMachineDelete:
		VirtualMachineDelete(*id)
	}

}

// VirtualMachineState returns current state of a virtual machine.
func VirtualMachineState(id string) {
	var ret VirtualMachine //VirtualMachineState

	d, err := v.DomainLookupByName(id)
        fmt.Printf("\n>>>>>>> %#v %#v %#v\n", d, err, id)
	herr(err)

	state, maxmem, mem, ncpu, cputime, err := v.DomainGetInfo(d)
        fmt.Printf(">>>>>>>>>> %#v %#v %#v %#v %#v %#v \n", state, maxmem, mem, ncpu, cputime, err)
	herr(err)

	ret.CPUCount = ncpu
	ret.CPUTime = cputime
	// God only knows why they return memory in kilobytes.
	ret.MemoryBytes = mem * 1024
	ret.MaxMemoryBytes = maxmem * 1024
	temp := libvirt.DomainState(state)
	herr(err)

	switch temp {
	case libvirt.DomainNostate:
		ret.State = VirtStatePending
	case libvirt.DomainRunning:
		ret.State = VirtStateRunning
	case libvirt.DomainBlocked:
		ret.State = VirtStateBlocked
	case libvirt.DomainPaused:
		ret.State = VirtStatePaused
	case libvirt.DomainShutdown:
		ret.State = VirtStateShutdown
	case libvirt.DomainShutoff:
		ret.State = VirtStateShutoff
	case libvirt.DomainCrashed:
		ret.State = VirtStateCrashed
	case libvirt.DomainPmsuspended:
		ret.State = VirtStateHybernating
	}
	v.DomainGetState(d, 0)
	hret(ret)
}

// VirtualMachineCreate creates a new VM from an xml template file
func VirtualMachineCreate(xmlTemplate string) {

	xml, err := ioutil.ReadFile(xmlTemplate)
	herr(err)

	d, err := v.DomainDefineXML(string(xml))
	herr(err)

	hret(d)
}

// VirtualMachineDelete deletes a new VM from an xml template file
func VirtualMachineDelete(id string) {
	d, err := v.DomainLookupByName(id)
	herr(err)
	err = v.DomainUndefineFlags(d, libvirt.DomainUndefineKeepNvram)
	herr(err)
	hok(fmt.Sprintf("%v was deleted", id))
}

// VirtualMachineSoftReboot reboots a machine gracefully, as chosen by hypervisor.
func VirtualMachineSoftReboot(id string) {
	d, err := v.DomainLookupByName(id)
	herr(err)

	err = v.DomainReboot(d, libvirt.DomainRebootDefault)
	herr(err)

	hok(fmt.Sprintf("%v was soft-rebooted successfully", id))
}

// VirtualMachineHardReboot sends a VM into hard-reset mode. This is damaging to all ongoing file operations.
func VirtualMachineHardReboot(id string) {
	d, err := v.DomainLookupByName(id)
	herr(err)

	err = v.DomainReset(d, 0)
	herr(err)

	hok(fmt.Sprintf("%v was hard-rebooted successfully", id))
}

// VirtualMachineShutdown gracefully shuts down the VM.
func VirtualMachineShutdown(id string) {
	d, err := v.DomainLookupByName(id)
	herr(err)

	err = v.DomainShutdown(d)
	herr(err)

	hok(fmt.Sprintf("%v was shutdown successfully", id))
}

// VirtualMachineShutoff kills running VM. Equivalent to pulling a plug out of a computer.
func VirtualMachineShutoff(id string) {
	d, err := v.DomainLookupByName(id)
	herr(err)

	err = v.DomainDestroy(d)
	herr(err)

	hok(fmt.Sprintf("%v was shutoff successfully", id))
}

// VirtualMachineStart starts up a VM.
func VirtualMachineStart(id string) {
	d, err := v.DomainLookupByName(id)
	herr(err)

	//v.DomainRestore()
	//_, err = v.DomainCreateWithFlags(d, uint32(libvirt.DomainStartBypassCache))
	err = v.DomainCreate(d)

	herr(err)

	hok(fmt.Sprintf("%v was started", id))
}

// VirtualMachinePause stops the execution of the VM. CPU is not used, but memory is still occupied.
func VirtualMachinePause(id string) {
	d, err := v.DomainLookupByName(id)
	herr(err)

	err = v.DomainSuspend(d)
	herr(err)

	hok(fmt.Sprintf("%v is paused", id))
}

// VirtualMachineResume can be called after Pause, to resume the invocation of the VM.
func VirtualMachineResume(id string) {
	d, err := v.DomainLookupByName(id)
	herr(err)

	err = v.DomainResume(d)
	herr(err)

	hok(fmt.Sprintf("%v was resumed", id))
}

func herr(e error) {
	if e != nil {
		fmt.Printf(`{"error":"%v"}\n`, strings.ReplaceAll(e.Error(), "\"", ""))
		os.Exit(1)
	}
}

func hok(message string) {
	fmt.Printf(`{"ok":"%v"}\n`, strings.ReplaceAll(message, "\"", ""))
	os.Exit(0)
}

func hret(i interface{}) {
	ret, err := json.Marshal(i)
	herr(err)
	fmt.Printf("%s\n",string(ret))
	os.Exit(0)
}

func virtinit() {
	v = libvirt.NewWithDialer(dialers.NewLocal(dialers.WithLocalTimeout(time.Second * 2)))
	if err := v.Connect(); err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
}
