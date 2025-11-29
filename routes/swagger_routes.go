package routes

import (
	"net/http"

	_ "app2_http_api_database/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterSwaggerRoutes(mux *http.ServeMux) {
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
}
