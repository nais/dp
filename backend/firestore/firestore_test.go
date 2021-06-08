package firestore_test

import (
	"context"
	"os"
	"testing"

	"github.com/nais/dp/backend/firestore"
	"github.com/stretchr/testify/assert"
)

func TestFirestoreCRUD(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("Skipping integration test")
	}

	f, err := firestore.New(context.Background(), "aura-dev-d9f5", "dp", "au")
	assert.NoError(t, err)

	t.Run("Create", func(t *testing.T) {
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
	})

	t.Run("Update", func(t *testing.T) {
		dps, err := f.GetDataproducts(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 1, len(dps))

		dp := dps[0]

		assert.Equal(t, "dp", dp.Dataproduct.Name)
		err = f.UpdateDataproduct(context.Background(), dp.ID, firestore.Dataproduct{Name: "dpdp"})
		assert.NoError(t, err)

		newDp, err := f.GetDataproduct(context.Background(), dp.ID)
		assert.NoError(t, err)

		assert.True(t, dp.Updated.Before(newDp.Updated))
		assert.True(t, newDp.Created.Before(newDp.Updated))
		assert.Equal(t, "dpdp", newDp.Dataproduct.Name)
	})

	t.Run("Delete", func(t *testing.T) {
		dps, err := f.GetDataproducts(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 1, len(dps))

		dp := dps[0]

		err = f.DeleteDataproduct(context.Background(), dp.ID)
		assert.NoError(t, err)

		dps, err = f.GetDataproducts(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 0, len(dps))
	})
}