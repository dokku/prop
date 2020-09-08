package backend

type Property struct {
	DataType  string
	Namespace string
	Key       string
	Value     string
}

type PropertyCollection struct {
	Properties []Property
}
