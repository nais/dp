package api

import (
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/nais/dp/apiserver/middleware"
	"gopkg.in/go-playground/validator.v8"
)

func New(client *firestore.Client, validate *validator.Validate) chi.Router {
	api := api{client, validate}

	latencyHistBuckets := []float64{.001, .005, .01, .025, .05, .1, .5, 1, 3, 5}
	prometheusMiddleware := middleware.PrometheusMiddleware("apiserver", latencyHistBuckets...)
	prometheusMiddleware.Initialize("/api/v1/", http.MethodGet, http.StatusOK)

	r := chi.NewRouter()

	r.Use(prometheusMiddleware.Handler())
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
	}))

	r.Get("/dataproducts", api.dataproducts)
	r.Post("/dataproducts", api.createDataproduct)
	r.Put("/dataproducts/{productID}", api.updateDataproduct)
	r.Get("/dataproducts/{productID}", api.getDataproduct)
	r.Delete("/dataproducts/{productID}", api.deleteDataproduct)

	return r
}
