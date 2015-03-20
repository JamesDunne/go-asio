package asio

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
