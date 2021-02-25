package mconfig

type MconfigClient interface {
	Config
	Adapter
	WatchChange
}

type Option func(*Options)
