package main

import "time"

type AttachRsp struct {
	ResourceType	string		`json:"resource_type"`
	ResourceId	string		`json:"resource_id"`
	HubId		string		`json:"hub_id"`
	Port		int		`json:"hub_port"`
}

type Connection struct {
	Port		int		`json:"port"`
	ResourceType	string		`json:"resource_type"`
	ResourceId	string		`json:"resource_id"`

}
type Hub struct {
	Id        	string     	`json:"id"`
	Name      	string	    	`json:"name"`
	NetworkProvider	string		`json:"network_provider"`
	NetworkType	string		`json:"network_type"`
	PortNum		int		`json:"ports_num"`
	Connections	[]Connection	`json:"connections"`
	
	StartTime	time.Time	`json:"startTime"`
	StopTime 	time.Time	`json:"stopTime"`
}

type Hubs []Hub
