package iam_test

import (
	"testing"
	"time"

	"github.com/nais/dp/backend/iam"
)

func TestAddMemberToBucket(t *testing.T) {
	t.Skip()

	bucketName := "container_resource_usage"
	member := "user:christine.teig@nav.no"
	start := time.Now()
	end := time.Now().AddDate(1, 0, 0)

	err := iam.UpdateBucketAccessControl(bucketName, member, start, end)
	if err != nil {
		return
	}
}

func TestRemoveMemberFromBucket(t *testing.T) {
	t.Skip()

	bucketName := "container_resource_usage"
	member := "user:christine.teig@nav.no"

	err := iam.RemoveMemberFromBucket(bucketName, member)
	if err != nil {
		return
	}
}
