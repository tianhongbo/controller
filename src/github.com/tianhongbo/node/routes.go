package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"EmulatorIndex",
		"GET",
		"/emulators",
		EmulatorIndex,
	},
	Route{
		"EmulatorCreate",
		"POST",
		"/emulators",
		EmulatorCreate,
	},
	Route{
		"EmulatorDelete",
		"GET",
		"/emulators/{id}",
		EmulatorShow,
	},
	Route{
		"EmulatorDelete",
		"DELETE",
		"/emulators/{id}",
		EmulatorDestroy,
	},
	//devices

	Route{
		"Devices",
		"GET",
		"/devices",
		DeviceIndex,
	},
	Route{
		"DeviceCreate",
		"POST",
		"/devices",
		DeviceCreate,
	},
	Route{
		"DeviceShow",
		"GET",
		"/devices/{imei}",
		DeviceShow,
	},
	Route{
		"DeviceDelete",
		"DELETE",
		"/devices/{imei}",
		DeviceDelete,
	},

	//hubs
	Route{
		"Hubs",
		"GET",
		"/hubs",
		HubIndex,
	},
	Route{
		"HubCreate",
		"POST",
		"/hubs",
		HubCreate,
	},
	Route{
		"HubShow",
		"GET",
		"/hubs/{id}",
		HubShow,
	},
	Route{
		"HubDelete",
		"DELETE",
		"/hubs/{id}",
		HubDelete,
	},
	Route{
		"HubAttach",
		"POST",
		"/hubs/{id}/connections",
		HubAttach,
	},
	Route{
		"HubDetach",
		"DELETE",
		"/hubs/{id}/connections",
		HubDetach,
	},
}
