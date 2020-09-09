package backend

type Property struct {
	DataType  string
	Namespace string
	Key       string
	Value     interface{}
}

type PropertyCollection struct {
	Properties []Property
}
