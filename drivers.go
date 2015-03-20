package asio

import (
	"fmt"
	"syscall"
	"unsafe"
)

type winUTF16string struct {
	utf16  []uint16
	length uint32
}

func (utfstring *winUTF16string) String() string {
	return syscall.UTF16ToString(utfstring.utf16[:utfstring.length])
}

func (utfstring *winUTF16string) Addr() *uint16 {
	return &utfstring.utf16[0]
}

func RegOpenKey(key syscall.Handle, subkey string, desiredAccess uint32) (handle syscall.Handle, err error) {
	var subkeyUTF16 *uint16
	subkeyUTF16, err = syscall.UTF16PtrFromString(subkey)
	if err != nil {
		return syscall.InvalidHandle, err
	}

	err = syscall.RegOpenKeyEx(
		key,
		subkeyUTF16,
		uint32(0),
		desiredAccess,
		&handle,
	)
	return
}

type ASIODriver struct {
	Name  string
	CLSID string
	GUID  *GUID

	com   *IDispatch
	iasio *IASIO
}

type IASIO struct {
	//virtual ASIOBool init(void *sysHandle) = 0;
	pInit uintptr
	//virtual void getDriverName(char *name) = 0;
	pGetDriverName uintptr

	//virtual long getDriverVersion() = 0;
	//virtual void getErrorMessage(char *string) = 0;
	//virtual ASIOError start() = 0;
	//virtual ASIOError stop() = 0;
	//virtual ASIOError getChannels(long *numInputChannels, long *numOutputChannels) = 0;
	//virtual ASIOError getLatencies(long *inputLatency, long *outputLatency) = 0;
	//virtual ASIOError getBufferSize(long *minSize, long *maxSize, long *preferredSize, long *granularity) = 0;
	//virtual ASIOError canSampleRate(ASIOSampleRate sampleRate) = 0;
	//virtual ASIOError getSampleRate(ASIOSampleRate *sampleRate) = 0;
	//virtual ASIOError setSampleRate(ASIOSampleRate sampleRate) = 0;
	//virtual ASIOError getClockSources(ASIOClockSource *clocks, long *numSources) = 0;
	//virtual ASIOError setClockSource(long reference) = 0;
	//virtual ASIOError getSamplePosition(ASIOSamples *sPos, ASIOTimeStamp *tStamp) = 0;
	//virtual ASIOError getChannelInfo(ASIOChannelInfo *info) = 0;
	//virtual ASIOError createBuffers(ASIOBufferInfo *bufferInfos, long numChannels, long bufferSize, ASIOCallbacks *callbacks) = 0;
	//virtual ASIOError disposeBuffers() = 0;
	//virtual ASIOError controlPanel() = 0;
	//virtual ASIOError future(long selector,void *opt) = 0;
	//virtual ASIOError outputReady() = 0;
}

func (drv *ASIODriver) Open() error {
	unk, err := CreateInstance(drv.GUID, drv.GUID)
	if err != nil {
		return err
	}
	drv.com = unk

	drv.iasio = (*IASIO)(unsafe.Pointer(unk.lpVtbl))
	return nil
}

func (drv *ASIODriver) Close() {
	drv.com.Release()
}

func (drv *ASIODriver) GetDriverName() string {
	fmt.Println(drv.iasio.pGetDriverName)

	name := [128]byte{}
	// TODO: fixme!
	syscall.Syscall(drv.iasio.pGetDriverName, 2,
		uintptr(unsafe.Pointer(drv.iasio)),
		uintptr(unsafe.Pointer(&name)),
		uintptr(0))
	return string(name[:])
}

func newDriver(key syscall.Handle, keynameUTF16 winUTF16string) (drv *ASIODriver, err error) {
	var subkey syscall.Handle
	err = syscall.RegOpenKeyEx(key, keynameUTF16.Addr(), 0, syscall.KEY_READ, &subkey)
	if err != nil {
		return nil, err
	}
	defer syscall.RegCloseKey(subkey)

	clsidName, err := syscall.UTF16PtrFromString("clsid")
	if err != nil {
		return nil, err
	}

	// Get CLSID of driver impl:
	clsidUTF16, datatype, datasize := make([]uint16, 128, 128), uint32(syscall.REG_SZ), uint32(256)
	err = syscall.RegQueryValueEx(subkey, clsidName, nil, &datatype, (*byte)(unsafe.Pointer(&clsidUTF16[0])), &datasize)
	if err != nil {
		return nil, err
	}

	// Convert the subkey name from UTF-16 to a string:
	keyname := keynameUTF16.String()
	drv = &ASIODriver{
		Name:  keyname,
		CLSID: syscall.UTF16ToString(clsidUTF16),
	}

	drv.GUID, err = CLSIDFromStringUTF16(&clsidUTF16[0])
	if err != nil {
		return nil, err
	}

	return drv, nil
}

// Enumerate list of ASIO drivers registered on the system
func ListDrivers() (drivers map[string]*ASIODriver, err error) {
	var key syscall.Handle
	key, err = RegOpenKey(syscall.HKEY_LOCAL_MACHINE, "Software\\ASIO", syscall.KEY_ENUMERATE_SUB_KEYS)
	if err != nil {
		return
	}
	defer syscall.RegCloseKey(key)

	drivers = make(map[string]*ASIODriver)

	// Enumerate subkeys:
	index := uint32(0)
	for err == nil {
		keynameUTF16 := winUTF16string{
			utf16:  make([]uint16, 128),
			length: uint32(128),
		}

		// Get next subkey:
		err = syscall.RegEnumKeyEx(key, index, keynameUTF16.Addr(), &keynameUTF16.length, nil, nil, nil, nil)
		// Determine when to stop:
		if err != nil {
			if errno, ok := err.(syscall.Errno); ok {
				// 259 is "No more data" error; aka end of enumeration.
				if uintptr(errno) == uintptr(259) {
					err = nil
					break
				}
			}
			fmt.Println(err)
			return
		}

		index++

		// Create an ASIODriver based on the key:
		drv, err := newDriver(key, keynameUTF16)
		if err != nil {
			continue
		}

		drivers[drv.Name] = drv
	}

	return drivers, nil
}
