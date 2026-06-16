package contract

type Consumer interface {
	Start()
	Stop()
	ListenErrs() <-chan error
}

type Processor interface {
	// Process this func must not be blocked
	Process([]byte)
}
