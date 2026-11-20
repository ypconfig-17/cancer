package cancer

type LifeCycle interface {
	start()
	isRunning() bool
	stop()
	shutdown()
}
