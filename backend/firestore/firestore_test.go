package firestore_test

import (
	"context"
	"github.com/nais/dp/backend/firestore"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFirestoreWithEmulator(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("Skipping integration test")
	}

	f, err := firestore.New(context.Background(), "aura-dev-d9f5", "dp", "au")
	assert.NoError(t, err)

	testDataproduct := firestore.Dataproduct{
		Name:        "dp",
		Description: "desc",
		Datastore:   []map[string]string{{}},
		Team:        "team",
		Access:      nil,
	}

	dpID, err := f.CreateDataproduct(context.Background(), testDataproduct)

	assert.NoError(t, err)

	dataproduct, err := f.GetDataproduct(context.Background(), dpID)
	assert.Equal(t, testDataproduct.Name, dataproduct.Dataproduct.Name)
}
