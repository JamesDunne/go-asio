package asio

import (
	"bytes"
	"syscall"
	"unsafe"
)

// Special ASIO error values:
const (
	ASE_OK      = 0          // This value will be returned whenever the call succeeded
	ASE_SUCCESS = 0x3f4847a0 // unique success return value for ASIOFuture calls
)

// Known ASIO error values:
const (
	ASE_NotPresent       = -1000 + iota // hardware input or output is not present or available
	ASE_HWMalfunction                   // hardware is malfunctioning (can be returned by any ASIO function)
	ASE_InvalidParameter                // input parameter invalid
	ASE_InvalidMode                     // hardware is in a bad mode or used in a bad mode
	ASE_SPNotAdvancing                  // hardware is not running when sample position is inquired
	ASE_NoClock                         // sample clock or rate cannot be determined or is not present
	ASE_NoMemory                        // not enough memory for completing the request
)

type Error struct {
	errno int32
	msg   string
}

// Fixed instances of errors:
var (
	ErrorNotPresent       = &Error{errno: ASE_NotPresent, msg: "hardware input or output is not present or available"}
	ErrorHWMalfunction    = &Error{errno: ASE_HWMalfunction, msg: "hardware is malfunctioning (can be returned by any ASIO function)"}
	ErrorInvalidParameter = &Error{errno: ASE_InvalidParameter, msg: "input parameter invalid"}
	ErrorInvalidMode      = &Error{errno: ASE_InvalidMode, msg: "hardware is in a bad mode or used in a bad mode"}
	ErrorSPNotAdvancing   = &Error{errno: ASE_SPNotAdvancing, msg: "hardware is not running when sample position is inquired"}
	ErrorNoClock          = &Error{errno: ASE_NoClock, msg: "sample clock or rate cannot be determined or is not present"}
	ErrorNoMemory         = &Error{errno: ASE_NoMemory, msg: "not enough memory for completing the request"}
)

// Mapping of known ASIO error values to Errors:
var knownErrors map[int32]*Error = map[int32]*Error{
	ASE_NotPresent:       ErrorNotPresent,
	ASE_HWMalfunction:    ErrorHWMalfunction,
	ASE_InvalidParameter: ErrorInvalidParameter,
	ASE_InvalidMode:      ErrorInvalidMode,
	ASE_SPNotAdvancing:   ErrorSPNotAdvancing,
	ASE_NoClock:          ErrorNoClock,
	ASE_NoMemory:         ErrorNoMemory,
}

func asError(drv *IASIO, ase uintptr) *Error {
	errno := int32(ase)

	switch errno {
	case ASE_OK:
		return nil
	case ASE_SUCCESS:
		return nil
	}
	if err, ok := knownErrors[errno]; ok {
		return err
	}

	// This rarely seems to return anything useful
	return &Error{errno: errno, msg: drv.GetErrorMessage()}
}

func (err *Error) Error() string {
	return err.msg
}

type ASIOCallbacks struct {
	//void (*bufferSwitch) (long doubleBufferIndex, ASIOBool directProcess);
	//	// bufferSwitch indicates that both input and output are to be processed.
	//	// the current buffer half index (0 for A, 1 for B) determines
	//	// - the output buffer that the host should start to fill. the other buffer
	//	//   will be passed to output hardware regardless of whether it got filled
	//	//   in time or not.
	//	// - the input buffer that is now filled with incoming data. Note that
	//	//   because of the synchronicity of i/o, the input always has at
	//	//   least one buffer latency in relation to the output.
	//	// directProcess suggests to the host whether it should immedeately
	//	// start processing (directProcess == ASIOTrue), or whether its process
	//	// should be deferred because the call comes from a very low level
	//	// (for instance, a high level priority interrupt), and direct processing
	//	// would cause timing instabilities for the rest of the system. If in doubt,
	//	// directProcess should be set to ASIOFalse.
	//	// Note: bufferSwitch may be called at interrupt time for highest efficiency.

	//void (*sampleRateDidChange) (ASIOSampleRate sRate);
	//	// gets called when the AudioStreamIO detects a sample rate change
	//	// If sample rate is unknown, 0 is passed (for instance, clock loss
	//	// when externally synchronized).

	//long (*asioMessage) (long selector, long value, void* message, double* opt);
	//	// generic callback for various purposes, see selectors below.
	//	// note this is only present if the asio version is 2 or higher

	//ASIOTime* (*bufferSwitchTimeInfo) (ASIOTime* params, long doubleBufferIndex, ASIOBool directProcess);
	//	// new callback with time info. makes ASIOGetSamplePosition() and various
	//	// calls to ASIOGetSampleRate obsolete,
	//	// and allows for timecode sync etc. to be preferred; will be used if
	//	// the driver calls asioMessage with selector kAsioSupportsTimeInfo.
}

// interface IASIO : public IUnknown {
type pIASIOVtbl struct {
	// v-tables are flattened in memory for simple direct cases like this.
	pIUnknownVtbl

	//virtual ASIOBool init(void *sysHandle) = 0;
	pInit uintptr
	//virtual void getDriverName(char *name) = 0;
	pGetDriverName uintptr
	//virtual long getDriverVersion() = 0;
	pGetDriverVersion uintptr
	//virtual void getErrorMessage(char *string) = 0;
	pGetErrorMessage uintptr

	//virtual ASIOError start() = 0;
	pStart uintptr
	//virtual ASIOError stop() = 0;
	pStop uintptr
	//virtual ASIOError getChannels(long *numInputChannels, long *numOutputChannels) = 0;
	pGetChannels uintptr
	//virtual ASIOError getLatencies(long *inputLatency, long *outputLatency) = 0;
	pGetLatencies uintptr
	//virtual ASIOError getBufferSize(long *minSize, long *maxSize, long *preferredSize, long *granularity) = 0;
	pGetBufferSize uintptr
	//virtual ASIOError canSampleRate(ASIOSampleRate sampleRate) = 0;
	pCanSampleRate uintptr
	//virtual ASIOError getSampleRate(ASIOSampleRate *sampleRate) = 0;
	pGetSampleRate uintptr
	//virtual ASIOError setSampleRate(ASIOSampleRate sampleRate) = 0;
	pSetSampleRate uintptr
	//virtual ASIOError getClockSources(ASIOClockSource *clocks, long *numSources) = 0;
	pGetClockSources uintptr
	//virtual ASIOError setClockSource(long reference) = 0;
	pSetClockSource uintptr
	//virtual ASIOError getSamplePosition(ASIOSamples *sPos, ASIOTimeStamp *tStamp) = 0;
	pGetSamplePosition uintptr
	//virtual ASIOError getChannelInfo(ASIOChannelInfo *info) = 0;
	pGetChannelInfo uintptr
	//virtual ASIOError createBuffers(ASIOBufferInfo *bufferInfos, long numChannels, long bufferSize, ASIOCallbacks *callbacks) = 0;
	pCreateBuffers uintptr
	//virtual ASIOError disposeBuffers() = 0;
	pDisposeBuffers uintptr
	//virtual ASIOError controlPanel() = 0;
	pControlPanel uintptr
	//virtual ASIOError future(long selector,void *opt) = 0;
	pFuture uintptr
	//virtual ASIOError outputReady() = 0;
	pOutputReady uintptr
}

type IASIO struct {
	vtbl_asio *pIASIOVtbl
}

func (drv *IASIO) AsIUnknown() *IUnknown { return (*IUnknown)(unsafe.Pointer(drv)) }

func (drv *IASIO) Init(sysHandle uintptr) (ok bool) {
	r1, _, _ := syscall.Syscall(drv.vtbl_asio.pInit, 2,
		uintptr(unsafe.Pointer(drv)),
		sysHandle,
		uintptr(0))
	ok = (r1 != 0)
	return
}

func (drv *IASIO) GetDriverName() string {
	name := [128]byte{0}
	syscall.Syscall(drv.vtbl_asio.pGetDriverName, 2,
		uintptr(unsafe.Pointer(drv)),
		uintptr(unsafe.Pointer(&name[0])),
		uintptr(0))

	lz := bytes.IndexByte(name[:], byte(0))
	return string(name[:lz])
}

func (drv *IASIO) GetDriverVersion() int32 {
	r1, _, _ := syscall.Syscall(drv.vtbl_asio.pGetDriverVersion, 2,
		uintptr(unsafe.Pointer(drv)),
		uintptr(0),
		uintptr(0))
	return int32(r1)
}

func (drv *IASIO) GetErrorMessage() string {
	str := [128]byte{0}

	_, _, _ = syscall.Syscall(drv.vtbl_asio.pGetErrorMessage, 2,
		uintptr(unsafe.Pointer(drv)),
		uintptr(unsafe.Pointer(&str[0])),
		uintptr(0))

	lz := bytes.IndexByte(str[:], byte(0))
	return string(str[:lz])
}

func (drv *IASIO) Start() (err error) {
	errno, _, _ := syscall.Syscall(drv.vtbl_asio.pStart, 1,
		uintptr(unsafe.Pointer(drv)),
		uintptr(0),
		uintptr(0))
	if derr := asError(drv, errno); derr != nil {
		return derr
	}
	return nil
}

func (drv *IASIO) Stop() (err error) {
	errno, _, _ := syscall.Syscall(drv.vtbl_asio.pStop, 1,
		uintptr(unsafe.Pointer(drv)),
		uintptr(0),
		uintptr(0))
	if derr := asError(drv, errno); derr != nil {
		return derr
	}
	return nil
}

//virtual ASIOError getChannels(long *numInputChannels, long *numOutputChannels) = 0;
func (drv *IASIO) GetChannels() (numInputChannels, numOutputChannels int, err error) {
	var tmpInputChannels, tmpOutputChannels uintptr

	errno, _, _ := syscall.Syscall(drv.vtbl_asio.pGetChannels, 3,
		uintptr(unsafe.Pointer(drv)),
		uintptr(unsafe.Pointer(&tmpInputChannels)),
		uintptr(unsafe.Pointer(&tmpOutputChannels)))

	serr := asError(drv, errno)
	if serr != nil {
		return 0, 0, serr
	}

	return int(tmpInputChannels), int(tmpOutputChannels), nil
}

////virtual ASIOError getLatencies(long *inputLatency, long *outputLatency) = 0;
func (drv *IASIO) GetLatencies() (inputLatency, outputLatency int, err error) {
	var tmpInputLatency, tmpOutputLatency uintptr

	errno, _, _ := syscall.Syscall(drv.vtbl_asio.pGetLatencies, 3,
		uintptr(unsafe.Pointer(drv)),
		uintptr(unsafe.Pointer(&tmpInputLatency)),
		uintptr(unsafe.Pointer(&tmpOutputLatency)))

	serr := asError(drv, errno)
	if serr != nil {
		return 0, 0, serr
	}

	return int(tmpInputLatency), int(tmpOutputLatency), nil
}

////virtual ASIOError getBufferSize(long *minSize, long *maxSize, long *preferredSize, long *granularity) = 0;
func (drv *IASIO) GetBufferSize() (minSize, maxSize, preferredSize, granularity int, err error) {
	var tmpminSize, tmpmaxSize, tmppreferredSize, tmpgranularity uintptr

	errno, _, _ := syscall.Syscall6(drv.vtbl_asio.pGetBufferSize, 5,
		uintptr(unsafe.Pointer(drv)),
		uintptr(unsafe.Pointer(&tmpminSize)),
		uintptr(unsafe.Pointer(&tmpmaxSize)),
		uintptr(unsafe.Pointer(&tmppreferredSize)),
		uintptr(unsafe.Pointer(&tmpgranularity)),
		uintptr(0),
	)

	serr := asError(drv, errno)
	if serr != nil {
		return 0, 0, 0, 0, serr
	}

	return int(tmpminSize), int(tmpmaxSize), int(tmppreferredSize), int(tmpgranularity), nil
}

// typedef double ASIOSampleRate;
////virtual ASIOError canSampleRate(ASIOSampleRate sampleRate) = 0;
func (drv *IASIO) CanSampleRate(sampleRate float64) (can bool) {
	errno, _, _ := syscall.Syscall(drv.vtbl_asio.pCanSampleRate, 2,
		uintptr(unsafe.Pointer(drv)),
		uintptr(unsafe.Pointer(&sampleRate)),
		uintptr(0))

	return errno != uintptr(0)
}

////virtual ASIOError getSampleRate(ASIOSampleRate *sampleRate) = 0;
//pGetSampleRate uintptr
func (drv *IASIO) GetSampleRate() (sampleRate float64, err error) {
	errno, _, _ := syscall.Syscall(drv.vtbl_asio.pGetSampleRate, 2,
		uintptr(unsafe.Pointer(drv)),
		uintptr(unsafe.Pointer(&sampleRate)),
		uintptr(0))

	serr := asError(drv, errno)
	if serr != nil {
		return 0., serr
	}

	return sampleRate, nil
}

////virtual ASIOError setSampleRate(ASIOSampleRate sampleRate) = 0;
//pSetSampleRate uintptr
func (drv *IASIO) SetSampleRate(sampleRate float64) (err error) {
	errno, _, _ := syscall.Syscall(drv.vtbl_asio.pSetSampleRate, 2,
		uintptr(unsafe.Pointer(drv)),
		uintptr(unsafe.Pointer(&sampleRate)),
		uintptr(0))

	serr := asError(drv, errno)
	if serr != nil {
		return serr
	}

	return nil
}

////virtual ASIOError getClockSources(ASIOClockSource *clocks, long *numSources) = 0;
//pGetClockSources uintptr
////virtual ASIOError setClockSource(long reference) = 0;
//pSetClockSource uintptr
////virtual ASIOError getSamplePosition(ASIOSamples *sPos, ASIOTimeStamp *tStamp) = 0;
//pGetSamplePosition uintptr
////virtual ASIOError getChannelInfo(ASIOChannelInfo *info) = 0;
//pGetChannelInfo uintptr
////virtual ASIOError createBuffers(ASIOBufferInfo *bufferInfos, long numChannels, long bufferSize, ASIOCallbacks *callbacks) = 0;
//pCreateBuffers uintptr
////virtual ASIOError disposeBuffers() = 0;
//pDisposeBuffers uintptr
////virtual ASIOError controlPanel() = 0;
//pControlPanel uintptr
////virtual ASIOError future(long selector,void *opt) = 0;
//pFuture uintptr
////virtual ASIOError outputReady() = 0;
//pOutputReady uintptr
