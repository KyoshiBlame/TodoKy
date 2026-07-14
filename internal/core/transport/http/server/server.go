package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/KyoshiBlame/TodoKy/internal/core/logger"
	core_http_middleware "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	//мультиплексер по входящему http запросу распознает через какие middleware тому следует пройти и в какой обрботчик его нужно направить
	mux    *http.ServeMux
	config Config
	log    *core_logger.Logger

	middleware []core_http_middleware.Middleware
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {

	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware...)
	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	//запуск через горутину чтобы создать GraceFullShutDown через context
	go func() {

		defer close(ch)

		s.log.Warn("Start HTTP server", zap.String("addr", s.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP error: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("ShutDown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutDownTimeout,
		)

		defer cancel()
		//остановка принятия http запросов, но не отменяет обработку старых
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutDown HTTP server: %w ", err)
		}

		s.log.Warn("HTTP server stopped")
	}

	return nil

}

func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(
			//убираем префик, потому что фичи о нём знать не обязательно
			prefix+"/",
			http.StripPrefix(prefix, router.WithMiddleware()),
		)
	}
}
