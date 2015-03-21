package asio

import (
	"bytes"
	"syscall"
	"unsafe"
)

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

func (obj *IASIO) AsIUnknown() *IUnknown { return (*IUnknown)(unsafe.Pointer(obj)) }

func (obj *IASIO) Init(sysHandle uintptr) (ok bool) {
	r1, _, _ := syscall.Syscall(obj.vtbl_asio.pInit, 2,
		uintptr(unsafe.Pointer(obj)),
		sysHandle,
		uintptr(0))
	ok = (r1 != 0)
	return
}

func (obj *IASIO) GetDriverName() string {
	name := [128]byte{0}
	syscall.Syscall(obj.vtbl_asio.pGetDriverName, 2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&name[0])),
		uintptr(0))

	lz := bytes.IndexByte(name[:], byte(0))
	return string(name[:lz])
}

func (obj *IASIO) GetDriverVersion() int32 {
	r1, _, _ := syscall.Syscall(obj.vtbl_asio.pGetDriverVersion, 2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(0),
		uintptr(0))
	return int32(r1)
}
