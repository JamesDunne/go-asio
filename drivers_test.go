package asio

import (
	"testing"
)
import "fmt"

func TestListDrivers(t *testing.T) {
	drivers, err := ListDrivers()
	if err != nil {
		t.Error(err)
		return
	}

	for _, drv := range drivers {
		fmt.Printf("%s: %s\n", drv.CLSID, drv.Name)
	}

	{
		fmt.Printf("CoInitialize(0)\n")
		CoInitialize(0)
		defer fmt.Printf("CoUninitialize()\n")
		defer CoUninitialize()

		ua1000 := drivers["UA-1000"]

		fmt.Printf("ua1000.Open()\n")
		err = ua1000.Open()
		defer fmt.Printf("ua1000.Close()\n")
		defer ua1000.Close()
		if err != nil {
			t.Error(err)
			return
		}

		drv := ua1000.ASIO

		fmt.Printf("UA-1000 opened.\n")

		fmt.Printf("getDriverName():      '%s'\n", drv.GetDriverName())
		fmt.Printf("getDriverVersion():   %d\n", drv.GetDriverVersion())
		fmt.Printf("getErrorMessage():    '%s'\n", drv.GetErrorMessage())

		// ASIO startup procedure:

		// getChannels
		in, out, err := drv.GetChannels()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			t.Error(err)
			return
		}
		fmt.Printf("getChannels():        %d, %d\n", in, out)

		// getBufferSize
		// getSampleRate
		// outputReady
		// createBuffers (set callbacks)
		// getChannelInfo (for N)
		// getLatencies
		// start

		// wait

		// stop
		// disposeBuffers
		// obj->Release()
	}
}
