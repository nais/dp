package iam

import (
	"context"
	"fmt"
	"time"
)

const (
	BucketType   = "bucket"
	BigQueryType = "bigquery"
)

func CheckDatastoreAccess(ctx context.Context, datastore map[string]string, subject string) (bool, error) {
	datastoreType := datastore["type"]
	if len(datastoreType) == 0 {
		return false, fmt.Errorf("no type defined")
	}

	switch datastoreType {
	case BucketType:
		return CheckAccessInBucket(ctx, datastore["bucket_id"], subject)
	case BigQueryType:
		return CheckAccessInBigQueryTable(ctx, datastore["project_id"], datastore["dataset_id"], datastore["resource_id"], subject)
	}

	return false, fmt.Errorf("unknown datastore type: %v", datastoreType)
}

func UpdateDatastoreAccess(ctx context.Context, datastore map[string]string, accessMap map[string]time.Time) error {
	datastoreType := datastore["type"]
	if len(datastoreType) == 0 {
		return fmt.Errorf("no type defined")
	}

	switch datastoreType {
	case BucketType:
		for subject, expiry := range accessMap {
			if err := UpdateBucketAccessControl(ctx, datastore["bucket_id"], subject, expiry); err != nil {
				return err
			}
			return nil
		}
	case BigQueryType:
		for subject := range accessMap {
			if err := UpdateBigqueryTableAccessControl(ctx, datastore["project_id"], datastore["dataset_id"], datastore["resource_id"], subject); err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("unknown datastore type: %v", datastoreType)
}

func RemoveDatastoreAccess(ctx context.Context, datastore map[string]string, subject string) error {
	datastoreType := datastore["type"]
	if len(datastoreType) == 0 {
		return fmt.Errorf("no type defined")
	}

	switch datastoreType {
	case BucketType:
		return RemoveMemberFromBucket(ctx, datastore["bucket_id"], subject)
	case BigQueryType:
		return RemoveMemberFromBigQueryTable(ctx, datastore["project_id"], datastore["dataset_id"], datastore["resource_id"], subject)
	}

	return fmt.Errorf("unknown datastore type: %v", datastoreType)
}
