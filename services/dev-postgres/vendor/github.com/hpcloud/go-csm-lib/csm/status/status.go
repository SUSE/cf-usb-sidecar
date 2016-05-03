package status

type Status int

const (
	None Status = 1 + iota
	Unknown
	Successful
	Failed
)
