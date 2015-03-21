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
			t.Error(err)
			return
		}
		fmt.Printf("getChannels():        %d, %d\n", in, out)

		// getBufferSize
		minSize, maxSize, preferredSize, granularity, err := drv.GetBufferSize()
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("getBufferSize():      %d, %d, %d, %d\n", minSize, maxSize, preferredSize, granularity)

		// canSampleRate
		ok := drv.CanSampleRate(44100.)
		fmt.Printf("canSampleRate(44100): %v\n", ok)

		// getSampleRate
		srate, err := drv.GetSampleRate()
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("getSampleRate():      %v\n", srate)

		// SetSampleRate
		err = drv.SetSampleRate(44100.)
		fmt.Printf("setSampleRate(44100): %v\n", err)

		// outputReady
		fmt.Printf("outputReady():        %v\n", drv.OutputReady())

		// open control panel:
		//drv.ControlPanel()

		// createBuffers (set callbacks)
		// getChannelInfo (for N)

		// getLatencies
		latin, latout, err := drv.GetLatencies()
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("getLatencies():       %d, %d\n", latin, latout)

		// start

		// wait

		// stop
		// disposeBuffers
		// obj->Release()
	}
}
