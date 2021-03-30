package mconfig

// MconfigClient ...
type MconfigClient interface {
	Config
	Adapter
	WatchChange
}

// Option ...
type Option func(*Options)
