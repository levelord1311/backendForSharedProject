package shutdown

import (
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/logging"
	"io"
	"os"
	"os/signal"
)

func Graceful(signals []os.Signal, closeItems ...io.Closer) {
	logger := logging.GetLogger()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, signals...)
	sig := <-sigCh
	logger.Infof("Caught signal %s. Shutting down...", sig)

	for _, closer := range closeItems {
		if err := closer.Close(); err != nil {
			logger.Errorf("failed to close %v: %v", closer, err)
		}
	}
}
