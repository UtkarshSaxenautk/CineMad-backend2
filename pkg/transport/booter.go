package transport

type Booter interface {
	Initialize()
	Run()
	Shutdown()
}
