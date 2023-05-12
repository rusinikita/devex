package tests

type Coverage struct {
	Package string
	Percent uint8
}

type Mutant struct {
	Package string
	File    string
	Type    string
	Status  string
}
