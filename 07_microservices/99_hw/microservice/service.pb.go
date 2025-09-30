// This file is a stub. To be replaced by protoc generated file
package main

type Event struct {
	Timestamp int64
	Consumer  string
	Method    string
	Host      string
}

func (Event) GetHost() string { return "not implemented yet" }

type Stat struct {
	Timestamp  int64
	ByMethod   map[string]uint64
	ByConsumer map[string]uint64
}

type StatInterval struct {
	IntervalSeconds uint64
}

type Nothing struct {
	Dummy bool
}
