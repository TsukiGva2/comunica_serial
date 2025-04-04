package main

import (
	"log"
	"time"
)

func main() {
	sender, err := NewSerialSender(115200)
	if err != nil {
		log.Fatalf("Failed to initialize SerialSender: %v", err)
	}
	defer sender.Close()

	pcData := &PCData{}
	pcData.Tags.Store(0)
	pcData.UniqueTags.Store(0)
	pcData.CommStatus.Store(false)
	pcData.WifiStatus.Store(false)
	pcData.Lte4Status.Store(false)
	pcData.RfidStatus.Store(false)
	pcData.SysVersion.Store(414)
	pcData.Backups.Store(0)
	pcData.Envios.Store(0)

	pcData.Send(sender)
	<-time.After(time.Second * 2)

	ticker := time.NewTicker(120 * time.Millisecond)
	defer ticker.Stop()

	log.Println("Starting to send data...")
	for range ticker.C {
		pcData.Tags.Add(1)
		pcData.Send(sender)
	}
}
