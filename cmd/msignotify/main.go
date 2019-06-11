package main

import (
	"flag"
	"fmt"

	msignotify "github.com/dfuse-io/example-push-notifications"
	"github.com/dfuse-io/example-push-notifications/storage"

	"github.com/anachronistic/apns"
)

var certPath = flag.String("cert", "", "TCP Listener addr for http")
var keyPath = flag.String("key", "", "TCP Listener addr for gRPC")
var apiKey = flag.String("api_key", "", "dfuse api key")

func main() {

	flag.Parse()
	if *certPath == "" || *keyPath == "" {
		panic(fmt.Sprintf("Need both cert and key parameters, got: %s %s", *certPath, *keyPath))
	}
	send := make(chan msignotify.Notification)

	storage := storage.NewMemoryStorage()
	storage.OptInDeviceToken("leslapinsdev", "BBF082487C7236F65F4B17645596A31A3234A304CF5AC4DB73A1B2C85A4D2445", msignotify.IOS)

	go func() {
		server := msignotify.NewServer(*apiKey, "dev1.api.dfuse.dev:443", storage)
		if err := server.Run(send); err != nil {
			panic(err)
		}
	}()

	apnsClient := apns.NewClient("gateway.sandbox.push.apple.com:2195", *certPath, *keyPath)

	for {
		notification := <-send

		if notification.DeviceToken.DeviceType == msignotify.IOS {
			payload := apns.NewPayload()
			payload.Alert = notification.Message
			payload.Badge = 1
			payload.Sound = "bingbong.aiff"

			pn := apns.NewPushNotification()
			pn.DeviceToken = notification.DeviceToken.Token
			pn.AddPayload(payload)

			resp := apnsClient.Send(pn)

			alert, _ := pn.PayloadString()
			fmt.Println("  Alert:", alert)
			fmt.Println("Success:", resp.Success)
			fmt.Println("  Error:", resp.Error)
		}
	}
}
