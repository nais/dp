package api_test

import (
	"github.com/nais/dp/backend/api"
	"testing"

	"github.com/nais/dp/backend/iam"
	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	missingType := map[string]string{"no": "type"}
	err := api.validateDatastore(missingType)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "no type defined")

	invalidType := map[string]string{"type": "nonexistent"}
	err = api.validateDatastore(invalidType)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "unknown datastore type: nonexistent")

	validBucket := map[string]string{"type": iam.BucketType, "project_id": "ap", "bucket_id": "x"}
	assert.NoError(t, api.validateDatastore(validBucket))
	invalidBucket := map[string]string{"type": iam.BucketType, "project_id": "ap"}
	err = api.validateDatastore(invalidBucket)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "missing key: bucket_id")

	validBigQuery := map[string]string{"type": iam.BigQueryType, "project_id": "pi", "dataset_id": "did", "resource_id": "rid"}
	assert.NoError(t, api.validateDatastore(validBigQuery))

	invalidBigQuery := map[string]string{"type": iam.BigQueryType, "project_id": "pi", "dataset_id": "did"}
	err = api.validateDatastore(invalidBigQuery)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "missing key: resource_id")
}
