package comunica_serial

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

type PCData struct {
	Tags                atomic.Int64
	UniqueTags          atomic.Int32
	CommStatus          atomic.Bool
	WifiStatus          atomic.Bool
	Lte4Status          atomic.Bool
	RfidStatus          atomic.Bool
	UsbStatus           atomic.Bool
	SysVersion          atomic.Int32
	Backups             atomic.Int32
	PermanentUniqueTags atomic.Int32
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

/*
This is based on the following checksum function

	bool check_sum(SafeString &msg) {
	  int idxStar = msg.indexOf('*');

	  cSF(check_sum_hex, 2);

	  msg.substring(check_sum_hex, idxStar + 1);

	  long sum = 0;

	  if (!check_sum_hex.hexToLong(sum)) {
	    return false;
	  }

	  for (size_t i = 1; i < idxStar; i++) {
	    sum ^= msg[i];
	  }

	  return (sum == 0);
	}
*/
func withChecksum(data string) string {
	var checksum byte

	for i := range len(data) {
		checksum ^= data[i]
	}

	return fmt.Sprintf("$%s*%02X", data, checksum)
}

func (pd *PCData) format() string {
	currentEpoch := time.Now().Unix()

	f := fmt.Sprintf("MYTMP;%d;%d;%d;%d;%d;%d;%d;%d;%d;%d;%d",
		pd.Tags.Load(), pd.UniqueTags.Load(), boolToInt(pd.CommStatus.Load()), boolToInt(pd.WifiStatus.Load()),
		boolToInt(pd.Lte4Status.Load()), boolToInt(pd.RfidStatus.Load()), boolToInt(pd.UsbStatus.Load()),
		pd.SysVersion.Load(), pd.Backups.Load(), pd.PermanentUniqueTags.Load(), currentEpoch)

	return withChecksum(f)
}

func (pd *PCData) Send(sender *SerialSender) {
	data := pd.format()
	log.Println("Sending data:", data)
	sender.SendData(data)
}
