package datatype

type ScanSuccess struct {
	Name   string
	Output string
}

var _ Renderable = (*ScanSuccess)(nil)

func (s ScanSuccess) Headers() []string {
	return []string{"name", "output"}
}

func (s ScanSuccess) Rows() []string {
	return []string{
		boldGreen.Sprint(s.Name),
		s.Output,
	}
}
