package utils

import (
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func WaitClose(log *zap.Logger, f func()) {
	ch := make(chan os.Signal, 10) // 10 为了IDE调试的时候不卡
	signal.Notify(ch, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	sig := <-ch
	log.Warn("--------------------------------------------------------------------")
	log.Warn("Receive Signal,server is shutting down...", zap.String("signal", sig.String()))
	log.Warn("--------------------------------------------------------------------")
	f()
}
