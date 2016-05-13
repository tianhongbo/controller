package main
import (
	"time"
	"os/exec"
	"fmt"
	"os"
	"strconv"
	"strings"
	"bytes"
)

/*
"id": 6001
"imei": "353188020901357",
"name": "Samsung S2",
"connected_hostname": "host100",
"vnc_port": 5901,
"ssh_port": 6000,
"adb_name": "G1002de5a082",
"start_time": "2015-07-20T11:12:59.625592627-07:00",
"stop_time": "0001-01-01T00:00:00Z"
*/

//try

type Device struct {
	IMEI string `json: "imei"`
	Id string    `json:"id"`
	Status string    `json:"status"` //"available": free, "processing": from create to completely boot, "running": VNC & SSH available
	Name string	`json: "name"`
	ConnectedHostname string `json: "connected_hostname"`
	VNCPort int `json: "vnc_port"`
	SSHPort int `json: "ssh_port"`
	ADBName string `json: "adb_name"`
	IP string `json:"-"`
	StartTime time.Time `json:"start_time"`
	StopTime time.Time `json:"stop_time"`

	CmdInit		*exec.Cmd 	`json:"-"`	//Initializing Script exectued after device allocated
	CmdAttach		*exec.Cmd 	`json:"-"`	//Enable network connection Script
	CmdDetach		*exec.Cmd 	`json:"-"`	//Disable network connection script

}

type Devices []Device

func (d *Device) start() {

	d.CmdInit.Start()

	fmt.Println("A device is successfully launched at IP: ", d.IP)

}

func (d *Device) stop() {
	var err error

	d.StopTime = time.Now()

	err = d.CmdInit.Process.Signal(os.Kill)
	err = d.CmdInit.Wait()
	fmt.Println(d.CmdInit.Stdout)
	//err = syscall.Kill(d.CmdInit.Process.Pid, 9)


	fmt.Println("A device is successfully stopped.", err)
}

func (d *Device) initCmd() {


	//init cmd "deviceinstall.sh 4e640638 192.168.1.16 vnc_port ssh_port
	cmdStr := INSTALL_SCRIPT_PATH + "deviceinstall.sh " + d.ADBName + " " + d.IP + " " + strconv.Itoa(d.VNCPort) + " " + strconv.Itoa(d.SSHPort)

	parts := strings.Fields(cmdStr)
	head := parts[0]
	parts = parts[1:len(parts)]
	d.CmdInit = exec.Command(head, parts...)
	randomBytes := &bytes.Buffer{}
	d.CmdInit.Stdout = randomBytes

}

func (d *Device) attach() {

	cmdStr := INSTALL_SCRIPT_PATH + "deviceattach.sh " + d.ADBName
	parts := strings.Fields(cmdStr)
	head := parts[0]
	parts = parts[1:len(parts)]
	cmd := exec.Command(head, parts...)
	randomBytes := &bytes.Buffer{}
	cmd.Stdout = randomBytes
	cmd.Run()
}

func (d *Device) detach() {

	cmdStr := INSTALL_SCRIPT_PATH + "devicedetach.sh " + d.ADBName
	parts := strings.Fields(cmdStr)
	head := parts[0]
	parts = parts[1:len(parts)]
	cmd := exec.Command(head, parts...)
	randomBytes := &bytes.Buffer{}
	cmd.Stdout = randomBytes
	cmd.Run()
}