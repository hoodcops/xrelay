package api

import (
	"github.com/go-chi/chi"
)

// InitRoutes sets up all routes reachable from the server
func InitRoutes(logger *zap.Logger, conn *amqp.Connection) *chi.Mux {
	router := chi.NewRouter()
	router.Mount("/status", statusCheckHandler(logger))
	
	return router
}

func statusCheckHandler(logger *zap.Logger, conn *amqp.Connection) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Requeust){

	}
}