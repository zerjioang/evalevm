package datatype

type ScanResult struct {
	Error error
}

type ScanErrorDetails struct {
	Name    string
	Message string
}

var _ Renderable = (*ScanErrorDetails)(nil)

func (s ScanErrorDetails) Headers() []string {
	return []string{"name", "error message"}
}

func (s ScanErrorDetails) Rows() []string {
	return []string{
		boldRed.Sprint(s.Name),
		s.Message,
	}
}
