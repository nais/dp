package iam_test

import (
	"github.com/nais/dp/backend/iam"
	"testing"
)

func TestAddMemberToBucket(t *testing.T) {

	bucketName := "container_resource_usage"
	member := "user:christine.teig@nav.no"

	err := iam.UpdateBucketAccessControl(bucketName, member)
	if err != nil {
		return
	}
}

func TestRemoveMemberFromBucket(t *testing.T) {

	bucketName := "container_resource_usage"
	member := "user:christine.teig@nav.no"

	err := iam.RemoveMemberFromBucket(bucketName, member)
	if err != nil {
		return
	}
}
