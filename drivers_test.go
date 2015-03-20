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
}
