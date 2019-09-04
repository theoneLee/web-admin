package e

type SelfError struct {
	Code int
	Msg  string
}

func Error() string {
	return ""
}
