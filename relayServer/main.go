package main

import (
	"fmt"
	"github.com/giskook/smarthome-interface/olderhc"
	//"os"
	//"io"
	"log"
	"net/http"
)

func main() {

	fmt.Println("Server Start")
	http.HandleFunc("/netgate/network/login", olderhc.LoginHandler)
	http.HandleFunc("/netgate/network/deldevice", olderhc.DelDeviceHandler)

	http.HandleFunc("/netgate/network/setname", olderhc.SetNameHandler)
	http.HandleFunc("/netgate/network/getdevattribute", olderhc.GetDevAttributeHandler)
	http.HandleFunc("/netgate/network/calldev", olderhc.CallDevHandler)
	http.HandleFunc("/netgate/network/devwarn", olderhc.DevWarnHandler)
	http.HandleFunc("/netgate/network/defencecancle", olderhc.DefenceCancleHandler)
	http.HandleFunc("/netgate/network/levelcontrol", olderhc.LevelControlHandler)
	http.HandleFunc("/netgate/network/onoffdev", olderhc.OnOffDevHandler)
	http.HandleFunc("/netgate/network/read_deployment_status", olderhc.ReadDeployMentStatusHandler)

	http.HandleFunc("/netgate/network/getzbnodes", olderhc.GetZBNodeHandler)
	http.HandleFunc("/netgate/network/read_onoff_status", olderhc.ReadOnOffStatusHandler)
	go olderhc.GetHttpRouter().Run()

	err := http.ListenAndServe(":"+olderhc.Port, nil)

	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}

}
