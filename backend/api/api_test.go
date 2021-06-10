package api_test

import (
	"github.com/nais/dp/backend/api"
	"github.com/nais/dp/backend/config"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

import (
	"context"
	"github.com/nais/dp/backend/firestore"
	"github.com/nais/dp/backend/iam"
)

type apiTest struct {
	Name          string
	Method        string
	Path          string
	Body          string
	ExpectedCode  int
	Authenticated bool
}

var createDPloggedIn = apiTest{
	Name:          "Create dataproduct",
	Method:        "POST",
	Path:          "/api/v1/dataproducts",
	Body:          `{"name": "cdp", "description": "desc", "datastore": [], "team": "team"}`,
	ExpectedCode:  http.StatusCreated,
	Authenticated: true,
}

var apiTests = map[string][]apiTest{
	"Create only": {
		createDPloggedIn,
	},
	// "fetch": {
	// 	createDPloggedIn,
	// },
}

func TestAPIHappy(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	f, err := firestore.New(ctx, "aura-dev-d9f5", "dp", "au")
	assert.NoError(t, err)

	iam := iam.New(ctx)
	mux := api.New(f, iam, config.Config{
		DevMode: true,
	}, nil, nil)

	testServer := httptest.NewServer(mux)
	defer testServer.Close()
	client := testServer.Client()

	for name, tcs := range apiTests {
		t.Run(name, func(t *testing.T) {
			// t.Parallel()

			for _, tc := range tcs {
				t.Run(tc.Name, func(t *testing.T) {
					var body io.Reader
					if tc.Body != "" {
						body = strings.NewReader(tc.Body)
					}
					req, err := http.NewRequest(tc.Method, testServer.URL+tc.Path, body)
					if err != nil {
						t.Error(err)
					}

					resp, err := client.Do(req)
					if err != nil {
						t.Error(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != tc.ExpectedCode {
						t.Errorf("expected %v, got %v (%v)", tc.ExpectedCode, resp.StatusCode, http.StatusText(resp.StatusCode))
					}
				})
			}
		})
	}
}
