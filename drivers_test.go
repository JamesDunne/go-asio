package asio

import "testing"

func TestListDrivers(t *testing.T) {
	drivers, err := ListDrivers()
	if err != nil {
		t.Error(err)
		return
	}

	for _, drv := range drivers {
		t.Logf("%s: %s\n", drv.CLSID, drv.Name)
	}

	{
		t.Log("CoInitialize(0)")
		CoInitialize(0)
		defer t.Log("CoUninitialize()")
		defer CoUninitialize()

		ua1000 := drivers["UA-1000"]

		t.Log("ua1000.Open()")
		ua1000.Open()
		defer t.Log("ua1000.Close()")
		defer ua1000.Close()

		t.Log("UA-1000 opened.")
		t.Logf("%s\n", ua1000.GetDriverName())
	}
}
