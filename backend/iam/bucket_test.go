package iam_test

import (
	"github.com/nais/dp/backend/iam"
	"testing"
)

func TestBucketIam(t *testing.T) {

	bucketName := "container_resource_usage"
	member := "user:christine.teig@nav.no"

	err := iam.UpdateBucketAccessControl(bucketName, member)
	if err != nil {
		return
	}
}
