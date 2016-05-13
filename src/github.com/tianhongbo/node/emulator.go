package main

import (
	"fmt"
	"strings"
	"time"
	"os/exec"
	"strconv"
	"os"
	"bytes"
)

type FreeEmulatorPortPool []int

type Emulator struct {
	Id			string 		`json:"id"`
	Name		string 		`json:"name"`
	Status		string       `json:"status"` //"processing": from create to completely boot, "running": VNC & SSH available
	ConnectedHostname string `json: "connected_hostname"`
	VNCPort 	int 		`json: "vnc_port"`
	SSHPort 	int 		`json: "ssh_port"`
	ADBName 	string 		`json: "-"`
	EmulatorPort int 		`json:"port"`
	StartTime	time.Time	`json:"start_time"`
	StopTime	time.Time 	`json:"stop_time"`
	Cmd			*exec.Cmd 	`json:"-"`
	CmdWaitBoot		*exec.Cmd 	`json:"-"`	//Initializing Script exectued after emulator start
	CmdInit		*exec.Cmd 	`json:"-"`	//Initializing Script exectued after emulator start
	CmdAttach		*exec.Cmd 	`json:"-"`	//Enable network connection Script
	CmdDetach		*exec.Cmd 	`json:"-"`	//Disable network connection script
}


type Emulators []Emulator


type ApiStartEmulatorResponse struct {
	Id	uint64	`json:"id"`
	Port	int	`json:"port"`
	
}


type ApiShowEmulatorResponse struct {
	Id		string		`json:"id"`
	Name		string		`json:"name"`
	Port		int		`json:"port"`
	StartTime	time.Time	`json:"starttime"`
	StopTime	time.Time	`json:"stoptime"`
}


func (emu *Emulator) startEmulator() {
	emu.Cmd.Start()

	fmt.Println("Start to create the emulator: ADB name = ", emu.ADBName)

	if err := emu.Cmd.Wait(); err != nil {
		fmt.Println("Emulator Cmd process is killed. Error code = ", err)
	} else {
		fmt.Println("Emulator Cmd process finished.")
	}

}

func (emu *Emulator) startEmulatorWaitBoot() {
	emu.CmdWaitBoot.Start()

	fmt.Println("Start to wait the emulator booting: ADB name = ", emu.ADBName)

	if err := emu.CmdWaitBoot.Wait(); err != nil {
		fmt.Println("Emulator CmdWaitBoot process is killed. Error code = ", err)
	} else {
		fmt.Println("Emulator CmdWaitBoot process finished.")
	}
	//update emulator status from "processing" to "running"
	RepoUpdateEmulatorStatus("running", emu.Id)


}
func (emu *Emulator) startInitEmulator() {
	emu.CmdInit.Start()

	fmt.Println("Start to initialize the emulator: ADB name = ", emu.ADBName)

	if err := emu.CmdInit.Wait(); err != nil {
		fmt.Println("Emulator Cmd Init process is killed. Error code = ", err)
	} else {
		fmt.Println("Emulator Cmd Init process finished.")
	}

}

func (emu *Emulator) stop() {
	var err error

	emu.StopTime = time.Now()

	//err = emu.Cmd.Process.Signal(os.Kill)
	//err = emu.Cmd.Wait()
	//fmt.Println(emu.Cmd.Stdout)

	if err = emu.Cmd.Process.Signal(os.Kill); err != nil {
		fmt.Println("Warning! Kill Emulator Cmd process faild. Error code = ", err)
	} else {
		fmt.Println("Emulator Cmd process is being killed... ID = ", emu.Id)
	}

	if err = emu.CmdWaitBoot.Process.Signal(os.Kill); err != nil {
		fmt.Println("Warning! Kill Emulator CmdWaitBoot process faild. Error code = ", err)
	} else {
		fmt.Println("Emulator CmdWaitBoot process is being killed... ID = ", emu.Id)
	}


	//Here os.Interrupt is necessary because CmdInit need to clean up before exit
	if err = emu.CmdInit.Process.Signal(os.Interrupt); err != nil {
		fmt.Println("Warning! Kill Emulator Cmd Init process faild. Error code = ", err)
	} else {
		fmt.Println("Emulator CMD Init process is being killed... ID = ", emu.Id)
	}

//	fmt.Println(emu.CmdInit.Stdout)
	//err = syscall.Kill(emu.CmdInit.Process.Pid, 9)


}

func (emu *Emulator) initCmd() {


	//init cmd emulator64-arm -avd myandroid -no-window -verbose -no-boot-anim -noskin
	//cmdStr := "emulator64-arm -avd android-api-10-"+strconv.Itoa(int(emu.EmulatorPort))+" -wipe-data -no-window -no-boot-anim -noskin -port " + strconv.Itoa(int(emu.EmulatorPort))
	cmdStr := "emulator64-arm -avd android-api-10-"+strconv.Itoa(int(emu.EmulatorPort))+" -wipe-data -no-window -port " + strconv.Itoa(int(emu.EmulatorPort))
	parts := strings.Fields(cmdStr)
	head := parts[0]
	parts = parts[1:len(parts)]
	emu.Cmd = exec.Command(head, parts...)
	randomBytes := &bytes.Buffer{}
	emu.Cmd.Stdout = randomBytes

	//init cmd emulatorwaitboot.sh emulator-5566
	cmdStr2 := INSTALL_SCRIPT_PATH + "emulatorwaitboot.sh " + emu.ADBName

	parts2 := strings.Fields(cmdStr2)
	head2 := parts2[0]
	parts2 = parts2[1:len(parts2)]
	emu.CmdWaitBoot = exec.Command(head2, parts2...)

	//init cmd install.sh emulator-5566 vnc_port ssh_port
	cmdStr3 := INSTALL_SCRIPT_PATH + "install.sh " + emu.ADBName + " " + strconv.Itoa(emu.VNCPort) + " " + strconv.Itoa(emu.SSHPort) + " " + strconv.Itoa(emu.EmulatorPort)

	parts3 := strings.Fields(cmdStr3)
	head3 := parts3[0]
	parts3 = parts3[1:len(parts3)]
	emu.CmdInit = exec.Command(head3, parts3...)

}

func (emu *Emulator) attach() {

	cmdStr := INSTALL_SCRIPT_PATH + "attach.sh " + emu.ADBName
	parts := strings.Fields(cmdStr)
	head := parts[0]
	parts = parts[1:len(parts)]
	cmd := exec.Command(head, parts...)
	//randomBytes := &bytes.Buffer{}
	//cmd.Stdout = randomBytes
	err := cmd.Run()
	fmt.Println("The emualator was attached, error code = ", err)
}

func (emu *Emulator) detach() {

	cmdStr := INSTALL_SCRIPT_PATH + "detach.sh " + emu.ADBName
	parts := strings.Fields(cmdStr)
	head := parts[0]
	parts = parts[1:len(parts)]
	cmd := exec.Command(head, parts...)
	//randomBytes := &bytes.Buffer{}
	//cmd.Stdout = randomBytes
	err := cmd.Run()
	fmt.Println("The emualator was detached, error code = ", err)
}