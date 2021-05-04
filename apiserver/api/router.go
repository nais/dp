package api

import (
	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"github.com/nais/dp/apiserver/middleware"
	"net/http"
)

func New(client *firestore.Client) chi.Router {
	api := api{client}

	latencyHistBuckets := []float64{.001, .005, .01, .025, .05, .1, .5, 1, 3, 5}
	prometheusMiddleware := middleware.PrometheusMiddleware("apiserver", latencyHistBuckets...)
	prometheusMiddleware.Initialize("/api/v1/", http.MethodGet, http.StatusOK)

	r := chi.NewRouter()

	r.Use(prometheusMiddleware.Handler())

	r.Get("/dataproducts", api.dataproducts)

	return r
}
