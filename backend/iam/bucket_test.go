package iam_test

import (
	"context"
	"testing"
	"time"

	"github.com/nais/dp/backend/iam"
)

func TestAddMemberToBucket(t *testing.T) {
	t.Skip()

	bucketName := "container_resource_usage"
	member := "user:christine.teig@nav.no"
	end := time.Now().AddDate(1, 0, 0)

	err := iam.UpdateBucketAccessControl(context.Background(), bucketName, member, end)
	if err != nil {
		return
	}
}

func TestRemoveMemberFromBucket(t *testing.T) {
	t.Skip()

	bucketName := "container_resource_usage"
	member := "user:christine.teig@nav.no"

	err := iam.RemoveMemberFromBucket(context.Background(), bucketName, member)
	if err != nil {
		return
	}
}
