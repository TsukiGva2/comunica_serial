package main

import (
	"fmt"
	"log"
	"sync/atomic"
)

type PCData struct {
	Tags       atomic.Int32
	UniqueTags atomic.Int32
	CommStatus atomic.Bool
	WifiStatus atomic.Bool
	Lte4Status atomic.Bool
	RfidStatus atomic.Bool
	SysVersion atomic.Int32
	Backups    atomic.Int32
	Envios     atomic.Int32
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (pd *PCData) format() string {
	return fmt.Sprintf("<%d;%d;%d;%d;%d;%d;%d;%d;%d>",
		pd.Tags.Load(), pd.UniqueTags.Load(), boolToInt(pd.CommStatus.Load()), boolToInt(pd.WifiStatus.Load()),
		boolToInt(pd.Lte4Status.Load()), boolToInt(pd.RfidStatus.Load()), pd.SysVersion.Load(), pd.Backups.Load(), pd.Envios.Load())
}

func (pd *PCData) Send(sender *SerialSender) {
	data := pd.format()
	log.Println("Sending data:", data)
	sender.SendData(data)
}
