package mconfigClient

type MconfigClient interface {
	Config
	Adapter
	WatchChange
}

type Option func(*Options)
