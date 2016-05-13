package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"os/exec"
	"strings"
	"bytes"
	"time"
	"github.com/gorilla/mux"

)

func Index(w http.ResponseWriter, r *http.Request) {

	cmdStr := "emulator -avd  Android_2.2 -no-window -verbose -no-boot-anim -noskin"
	//cmdStr := "android list target"
	parts := strings.Fields(cmdStr)
	command := parts[0]
	args := parts[1:len(parts)]
	fmt.Println(command, args)

	//parts := strings.Fields(cmd)
	//head := parts[0]
	//parts = parts[1:len(parts)]
	//fmt.Println(head, parts)

	cmd := exec.Command(command, args...)
	randomBytes := &bytes.Buffer{}
	cmd.Stdout = randomBytes

	// Start command asynchronously
	err := cmd.Start()
	fmt.Println("A emulator is launched: %s \n %s", err, randomBytes.String())
	fmt.Fprintf(w, "Successful!\n")
}

// Emulator
// /emulators GET
func EmulatorIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(emulators); err != nil {
		panic(err)
	}
}

// /emulators POST
func EmulatorCreate(w http.ResponseWriter, r *http.Request) {
	var e Emulator

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &e); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		fmt.Println("Error! Can't unmarshal Json from create emulator request.", body)
		return
	}

	if e.EmulatorPort, err = RepoAllocateEmulatorPort(); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		fmt.Println("Error! Can't allocate Emulator port for create emulator request.")

		return
	}

	//Assign ADB name according Android SDK naming convention "emulator-'port'". e.g. emulator-5554
	e.ADBName = "emulator-" + strconv.Itoa(e.EmulatorPort)

	//Allocate SSH port
	if e.SSHPort, err = RepoAllocateSSHPort(); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		fmt.Println("Error! Can't allocate SSH port for create emulator request.")
		return
	}

	//Allocate VNC port
	if e.VNCPort, err = RepoAllocateVNCPort(); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		fmt.Println("Error! Can't allocate VNC port for create emulator request.")
		return
	}

	e.StartTime = time.Now()	//startTime	time.Time
	e.StopTime = time.Time{}	//stopTime	time.Time

	e.ConnectedHostname = THIS_HOST_NAME
	e.Status = "processing"

	e.initCmd()
	go e.startEmulator()
	go e.startEmulatorWaitBoot()
	go e.startInitEmulator()

	RepoCreateEmulator(e)

	fmt.Println("Emulator was created successfully. ID = ", e.Id)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		panic(err)
	}

	return
}


// /emulators/{id} DELETE
func EmulatorDestroy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string

	id = vars["id"]

	if e,err := RepoFindEmulator(id); err == nil {

		fmt.Println("Delete emulator successfully: ", e.CmdInit)
		e.stop()
		RepoFreeEmulatorPort(e.EmulatorPort)
		RepoFreeSSHPort(e.SSHPort)
		RepoFreeVNCPort(e.VNCPort)
		RepoDestroyEmulator(id)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	return
}

// /emulators/{id} GET
func EmulatorShow(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var id string
	var err error
	var e Emulator

	id = vars["id"]

	if e,err = RepoFindEmulator(id); err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(e); err != nil {
			panic(err)
		}
	}
	return
}


func HubIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(hubs); err != nil {
		panic(err)
	}
}

func HubShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var hubId string
	var err error

	hubId = vars["id"]

	hub,err := RepoFindHub(hubId)

	if err == nil {
		//Find the hub
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(hub); err != nil {
			panic(err)
		}
	} else {
		// If we didn't find a hub, 404
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}
	}


	return

}

func HubDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var hubId string
	var err error

	hubId = vars["id"]
	err = RepoDestroyHub(hubId)
	if err == nil {
		// Find the hub and deleted successfully
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

	} else {
		// If we didn't find it, 404
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}

	}

	return
}


/*
Test with this curl command:

curl -H "Content-Type: application/json" -d '{"name":"New Hub"}' http://localhost:8080/hubs

*/
func HubCreate(w http.ResponseWriter, r *http.Request) {
	var hub Hub
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &hub); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	//Initialize the connections

	for i := 0; i < hub.PortNum; i++ {
		hub.Connections = append(hub.Connections, Connection{i, "", ""})
	}
	//Initialize the start time
	hub.StartTime = time.Now()
	

	t := RepoCreateHub(hub)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

func HubAttach(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var hubId string
	var connection Connection
	var err error

	hubId = vars["id"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &connection); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	if i ,err := RepoAttachHub(hubId, connection); err == nil {
		// attach successfully

		// manipulate the device network
		fmt.Println("Turn on the Intenet connection ")
		switch connection.ResourceType {
		case "emulator":
			if e,err := RepoFindEmulator(connection.ResourceId); err == nil {
				e.attach();
			}
			break;
		case "device":
			if e,err := RepoFindDeviceById(connection.ResourceId); err == nil {
				e.attach();
			}
			break;
		default:
			fmt.Println("The resource type is not support.", connection.ResourceType)
			break;
		}

		// respond
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		t := AttachRsp{connection.ResourceType, connection.ResourceId, hubId, i}
		if err := json.NewEncoder(w).Encode(t); err != nil {
			panic(err)
		}

	} else {
		// attach failed
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(jsonErr{Code: 422, Text: "No available port."}); err != nil {
			panic(err)
		}

	}

	return

}

func HubDetach(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var hubId string
	var err error

	connection := Connection{}

	hubId = vars["id"]


	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &connection); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	if err := RepoDetachHub(hubId, connection); err == nil {
		// Detach successfully

		// manipulate the device network
		fmt.Println("Turn off the Intenet connection ")
		switch connection.ResourceType {
		case "emulator":
			if e,err := RepoFindEmulator(connection.ResourceId); err == nil {
				e.detach()
			}
			break;
		case "device":
			if e,err := RepoFindDeviceById(connection.ResourceId); err == nil {
				e.detach()
			}
			break;
		default:
			fmt.Println("The resource type is not support.", connection.ResourceType)
			break;
		}

		// respond
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

	} else {
		// detach failed
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}

	}

	return
}

/*
Device Operation
 */

// /devices GET
func DeviceIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(devices); err != nil {
		panic(err)
	}
}

// /device POST
func DeviceCreate(w http.ResponseWriter, r *http.Request) {
	var d Device

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &d); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		fmt.Println("Error! Can't unmarshal JSON data from create device request.")
		return
	}

	// confirm the device is existing with IMEI
	if  _, err = RepoFindDevice(d.IMEI); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		fmt.Println("Error! Can't find the device for create device request. IMEI = ", d.IMEI)
		return
	}

	if d.SSHPort, err = RepoAllocateSSHPort(); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		fmt.Println("Error! Can't allocate SSH port resource for create device request. IMEI = ", d.IMEI)
		return

	}

	if d.VNCPort, err = RepoAllocateVNCPort(); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		fmt.Println("Error! Can't allocate VNC port resource for create device request. IMEI = ", d.IMEI)
		return

	}

	//Initialize the start time
	d.StartTime = time.Now()

	//save to repository
	RepoAllocateDevice(d)

	d,_ = RepoFindDevice(d.IMEI)

	// should be start() after find because some initialization is done inside "RepoAllocateDevice(d)"

	d.start()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(d); err != nil {
		panic(err)
	}

	return
}

// /devices/{id} GET
func DeviceShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var imei string

	imei = vars["imei"]

	if d,err := RepoFindDevice(imei); err == nil {
		//Find the device
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(d); err != nil {
			panic(err)
		}
	} else {
		// If we didn't find a device, 404
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}
	}

	return

}

// /devices/{imei} DELETE
func DeviceDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var imei string
	var err error
	var d Device

	imei = vars["imei"];

	//Free resources alloacted to this device
	if d, err = RepoFindDevice(imei); err != nil {
		// If we didn't find it, 404
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}
		return
	}

	if d.Status == "available" {
		// Find the device, but it is in free status
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		return
	}

	//Free resources
	RepoFreeSSHPort(d.SSHPort)
	RepoFreeVNCPort(d.VNCPort)

	//stop commands
	d.stop()

	//Free devcie
	if err = RepoFreeDevice(imei); err != nil {
		// if it is failed to free the device, then return "404"
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}
		return
	}

	// device is deleted successfully
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	return
}

