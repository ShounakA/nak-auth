package models

type from[To any] interface {
	From() To
}

type DataModel interface {
	CreateTable() string
	Equals(interface{}) bool
}
