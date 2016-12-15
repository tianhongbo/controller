package main

import (
        "log"
        "os/exec"
	"strings"
)

func createEmulators() {

	//init cmd emulatorwaitboot.sh emulator-5566
	cmdStr2 := "/home/ubuntu/controller/bin/create_emulators.sh"

	parts2 := strings.Fields(cmdStr2)
	head2 := parts2[0]
	parts2 = parts2[1:len(parts2)]
	cmd := exec.Command(head2, parts2...)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("emulators are created.")
}
