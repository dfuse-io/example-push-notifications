package main

import (
	"flag"
	"fmt"

	"github.com/anachronistic/apns"
)

var certPath = flag.String("cert", "", "TCP Listener addr for http")
var keyPath = flag.String("key", "", "TCP Listener addr for gRPC")

func main() {

	flag.Parse()
	if *certPath == "" || *keyPath == "" {
		panic(fmt.Sprintf("Need both cert and key parameters, got: %s %s", *certPath, *keyPath))
	}
	send := make(chan Notification)

	db := NewDatabase()
	db.OptInDeviceToken("lelapinnoir2", "bbf082487c7236f65f4b17645596a31a3234a304cf5ac4db73a1b2c85a4d2445", IOS)

	go func() {
		server := NewServer("server_88bb56ce30a09e547450d9dc84e55716", db)
		if err := server.Run(send); err != nil {
			panic(err)
		}
	}()

	apnsClient := apns.NewClient("gateway.sandbox.push.apple.com:2195", *certPath, *keyPath)

	for {
		notification := <-send

		payload := apns.NewPayload()
		payload.Alert = notification.Message
		payload.Badge = 1
		payload.Sound = "bingbong.aiff"

		pn := apns.NewPushNotification()
		pn.DeviceToken = notification.DeviceToken
		pn.AddPayload(payload)

		resp := apnsClient.Send(pn)

		alert, _ := pn.PayloadString()
		fmt.Println("  Alert:", alert)
		fmt.Println("Success:", resp.Success)
		fmt.Println("  Error:", resp.Error)
	}
}
