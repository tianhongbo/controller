package main

import (
	"fmt"
	//"strings"
	//"bytes"
	//"os/exec"
	//"strconv"
	"errors"
	"time"
)


// Initiate the repository data
func init() {
	//intialize free VNC port pool
	for i := VNC_MIN_PORT; i <= VNC_MAX_PORT; i++ {
		freeVNCPortPool = append(freeVNCPortPool,i)
	}
	//intialize free SSH port pool
	for i := SSH_MIN_PORT; i <= SSH_MAX_PORT; i++ {
		freeSSHPortPool = append(freeSSHPortPool,i)
	}
	//intialize free emulator port pool
	for i := EMULATOR_MIN_PORT; i <= EMULATOR_MAX_PORT; i = i+2 {
		freeEmulatorPortPool = append(freeEmulatorPortPool,i)
	}

	//intialize devices list
	RepoCreateDevice(Device{IMEI:"357288042352104", Status:"available", ADBName:"4e640638", IP:"10.189.146.231", ConnectedHostname: THIS_HOST_NAME})
	//RepoCreateDevice(Device{IMEI:"353188020902633", Status:"available", ADBName:"G1002de5a083", IP:"192.168.1.17", ConnectedHostname: THIS_HOST_NAME})
	//RepoCreateDevice(Device{IMEI:"353188020902634", Status:"available", ADBName:"G1002de5a084", IP:"192.168.1.18", ConnectedHostname: THIS_HOST_NAME})


	//fmt.Println(inventories)
	//fmt.Println(freeVNCPortPool)
	//fmt.Println(freeSSHPortPool)
	//fmt.Println(freeEmulatorPortPool)
}

const THIS_HOST_NAME = "host101"
const EMULATOR_MIN_PORT = 5554
const EMULATOR_MAX_PORT = 5584
const EMULATOR_MAX_NUM = (EMULATOR_MAX_PORT-EMULATOR_MIN_PORT)/2

//For external use
const VNC_MIN_PORT = 5901
const VNC_MAX_PORT = 5920

//For localhost ADB mapping
//To avoid the conflict with external use VNC port
const VNC_INTERNAL_MIN_PORT = 5941
const VNC_INTERNAL_MAX_PORT = 5960

const SSH_MIN_PORT = 5921
const SSH_MAX_PORT = 5940

/*
 * For Mac Only
const INSTALL_SCRIPT_PATH  = "/Users/Scott/master/src/github.com/tianhongbo/node/"
const INSTALL_DATA_PATH  = "/Users/Scott/master/src/github.com/tianhongbo/node/"
*/

/*
 * For Ubuntu only
 */
const INSTALL_SCRIPT_PATH  = "/home/ubuntu2/controller/src/github.com/tianhongbo/node/"
const INSTALL_DATA_PATH  = "/home/ubuntu2/controller/src/github.com/tianhongbo/node/"

//VNC available port pool#
var freeVNCPortPool FreeVNCPortPool

func RepoAllocateVNCPort() (int, error) {
	for i,port := range freeVNCPortPool {
		freeVNCPortPool = append(freeVNCPortPool[:i], freeVNCPortPool[i+1:]...)
		return port, nil
	}
	return 0, errors.New("can't find the available emulator port resource.")
}

func RepoFreeVNCPort(port int) {
	if (port >= VNC_MIN_PORT && port <= VNC_MAX_PORT) {
		freeVNCPortPool = append(freeVNCPortPool, port)
	}
	return
}

//SSH available port pool#
var freeSSHPortPool FreeSSHPortPool

func RepoAllocateSSHPort() (int, error) {
	for i,port := range freeSSHPortPool {
		freeSSHPortPool = append(freeSSHPortPool[:i], freeSSHPortPool[i+1:]...)
		return port, nil
	}
	return 0, errors.New("can't find the available emulator port resource.")
}

func RepoFreeSSHPort(port int) {
	if (port >= SSH_MIN_PORT && port <= SSH_MAX_PORT) {
		freeSSHPortPool = append(freeSSHPortPool, port)
	}
	return
}


//Emulator available port pool #
var freeEmulatorPortPool FreeEmulatorPortPool

func RepoAllocateEmulatorPort() (int, error) {
	for i,port := range freeEmulatorPortPool {
		freeEmulatorPortPool = append(freeEmulatorPortPool[:i], freeEmulatorPortPool[i+1:]...)
		return port, nil
	}
	return 0, errors.New("can't find the available emulator port resource.")
}

func RepoFreeEmulatorPort(port int) {
	if (port >= EMULATOR_MIN_PORT && port <= EMULATOR_MAX_PORT) {
		freeEmulatorPortPool = append(freeEmulatorPortPool, port)
	}
	return
}

//Emulators
var emulators Emulators

func RepoFindEmulator(id string) (Emulator, error) {
	for _, e := range emulators {
		if e.Id == id {
			return e, nil
		}
	}

	return Emulator{}, errors.New("can't find the emulator.")

}

func RepoCreateEmulator(e Emulator) Emulator {
	emulators = append(emulators, e)
	return e
}

func RepoUpdateEmulatorStatus(s string, id string) Emulator {
	for i, e := range emulators {
		if e.Id == id {
			emulators[i].Status = s
			return e
		}
	}
	return Emulator{}
}

func RepoDestroyEmulator(id string) error {
	for i, e := range emulators {
		if e.Id == id {
			emulators = append(emulators[:i], emulators[i+1:]...)
			return nil
		}
	}
	return errors.New("can't find the emulator.")
}

// Mobile hubs
var hubs Hubs

func RepoFindHub(id string) (Hub,error) {
	for _, t := range hubs {
		if t.Id == id {
			return t, nil
		}
	}
	// return empty Hub if not found
	return Hub{}, errors.New("can't find the hub.")
}

func RepoCreateHub(t Hub) Hub {
	hubs = append(hubs, t)
	return t
}

func RepoDestroyHub(id string) error {
	for i, t := range hubs {
		if t.Id == id {
			hubs = append(hubs[:i], hubs[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Hub with id of %d to delete", id)
}

func RepoAttachHub(id string, connection Connection) (int, error) {
	for i, t := range hubs {
		if t.Id == id {
			for j, c := range hubs[i].Connections {
				if c.ResourceId == "" && c.ResourceType == "" {
					hubs[i].Connections[j].ResourceId = connection.ResourceId
					hubs[i].Connections[j].ResourceType = connection.ResourceType
					return c.Port, nil
				}
			}
			return 0, fmt.Errorf("Could not find free port to attach in the hub id of %d", id)
		}

	}

	return 0, fmt.Errorf("Could not find Hub with id of %d to attach", id)
}

func RepoDetachHub(id string, connection Connection) error {
	for i, t := range hubs {
		if t.Id == id {
			for j, c := range hubs[i].Connections {
				if c.ResourceId == connection.ResourceId && c.ResourceType == connection.ResourceType {
					hubs[i].Connections[j].ResourceId = ""
					hubs[i].Connections[j].ResourceType = ""
					return nil
				}
			}
			return errors.New("can't find the resource id.")
		}

	}

	return errors.New("can't find the hub.")
}

/*
The following are devices
 */

var devices Devices

func RepoFindDeviceById(id string) (Device,error) {
	for _, d := range devices {
		if d.Id == id {
			return d, nil
		}
	}
	// return empty Hub if not found
	return Device{}, errors.New("can't find the device.")
}

func RepoFindDevice(imei string) (Device,error) {
	for _, d := range devices {
		if d.IMEI == imei {
			return d, nil
		}
	}
	// return empty Hub if not found
	return Device{}, errors.New("can't find the device.")
}

func RepoCreateDevice(d Device) Device {
	devices = append(devices, d)
	return d
}

func RepoDestroyDevice(imei string) error {
	for i, d := range devices {
		if d.IMEI == imei {
			devices = append(devices[:i], devices[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find device with IMEI of %d to destroy", imei)
}

func RepoAllocateDevice(device Device) (error) {
	for i, d := range devices {
		if d.IMEI == device.IMEI {
			if d.Status != "available" {
				fmt.Println("Warning! the device is not on available status when it is allocated.")
			}
			devices[i].Status = "running"
			devices[i].Id = device.Id
			devices[i].Name = device.Name
			devices[i].SSHPort = device.SSHPort
			devices[i].VNCPort = device.VNCPort
			devices[i].StartTime = device.StartTime
			devices[i].StopTime = device.StopTime
			devices[i].initCmd()
			fmt.Println("Device is allocated. IMEI = ", device.IMEI)
			return nil
		}
	}
	// return error if not found
	return errors.New("can't find the device in the devices.")
}

func RepoFreeDevice(imei string) error {
	for i, d := range devices {
		if d.IMEI == imei {
			fmt.Println("Device is free. IMEI = ", imei)
			devices[i].Status = "available"
			devices[i].Id = ""
			devices[i].Name = ""
			devices[i].SSHPort = 0
			devices[i].VNCPort = 0
			devices[i].StartTime = time.Time{}
			devices[i].StopTime = time.Time{}

			return nil

		}
	}
	// return error if not found
	return errors.New("can't find the device in the devices.")
}