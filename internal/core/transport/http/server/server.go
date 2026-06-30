package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/KyoshiBlame/TodoKy/internal/core/logger"
	"go.uber.org/zap"
)


type HTTPServer struct {
	//мультиплексер по входящему http запросу распознает через какие middleware тому следует пройти и в какой обрботчик его нужно направить
	mux *http.ServeMux
	config Config
	log core_logger.Logger
}


func NewHTTPServer(
	config Config,
	log core_logger.Logger,
) *HTTPServer {
	return &HTTPServer{
		mux: http.NewServeMux(),
		config: config,
		log: log,
	}
}

func (h *HTTPServer) Run(ctx context.Context) error {
	server := &http.Server {
		Addr: h.config.Addr,
		Handler: h.mux,

	}

	ch := make(chan error, 1)

	//запуск через горутину чтобы создать GraceFullShutDown через context
	go func () {
		
		defer close(ch)

		h.log.Warn("Start HTTP server", zap.String("addr", h.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	} ()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP error: %w", err)
		}
	case <-ctx.Done() :
		h.log.Warn("ShutDown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			h.config.ShutDownTimeout,
		)

		defer cancel()
		//остановка принятия http запросов, но не отменяет обработку старых
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutDown HTTP server: %w ", err)
		}

		h.log.Warn("HTTP server stopped")
	}

	return nil

}