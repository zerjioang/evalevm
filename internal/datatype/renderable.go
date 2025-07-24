package datatype

type Renderable interface {
	Headers() []string
	Rows() []string
}
