package server

import (
	"os"
	"os/signal"
	"syscall"
)

// OnTerminate registers a callback function to be called when the program receives a termination signal.
func OnTerminate(fn func(sig os.Signal)) chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go terminationSignalHandler(sigChan, fn)
	return sigChan
}

// terminationSignalHandler waits for a signal on the provided channel and calls the provided function when a signal is received.
func terminationSignalHandler(c chan os.Signal, f func(sig os.Signal)) {
	sig := <-c

	f(sig)

	signal.Stop(c)
	close(c)
}
