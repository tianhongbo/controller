package main

import (
	"log"
	"net/http"
)

/*
 * Bring up all emulators
 * This routine will be called before main()
 */
func init() {
	createEmulators()
}
func main() {

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":80", router))
}
