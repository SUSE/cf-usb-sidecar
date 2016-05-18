package csm

type CSMConnection interface {
	Write(response CSMResponse) error
	WriteError(err error) error
}
