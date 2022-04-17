package model

type Data struct {
	Key   string
	Value string
}

type Event interface {
	IsEvent()
	String() string
	Data() Data
}
