package iam

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/iam"
)

func UpdateDatasetAccessControl(entity, projectID, datasetID string) error {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	ds := client.Dataset(datasetID)
	meta, err := ds.Metadata(ctx)
	if err != nil {
		return fmt.Errorf("Getting metadata for %v in %v: %v", datasetID, projectID, err)
	}

	// Append a new access control entry to the existing access list.
	update := bigquery.DatasetMetadataToUpdate{
		Access: append(meta.Access, &bigquery.AccessEntry{
			Role:       bigquery.ReaderRole,
			EntityType: bigquery.UserEmailEntity,
			Entity:     entity},
		),
	}

	// Leverage the ETag for the update to assert there's been no modifications to the
	// dataset since the metadata was originally read.
	if _, err := ds.Update(ctx, update, meta.ETag); err != nil {
		return fmt.Errorf("Updating %v in %v: %v", datasetID, projectID, err)
	}
	return nil
}

func UpdateBigqueryTableAccessControl(ctx context.Context, projectID, datasetID, tableID, member string) error {
	bqClient, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer bqClient.Close()

	policy, err := getPolicy(ctx, bqClient, datasetID, tableID)
	if err != nil {
		return fmt.Errorf("Getting policy for %v.%v in %v: %v", datasetID, tableID, projectID, err)
	}

	// no support for V3 for BigQuery yet, and no support for conditions
	role := "roles/bigquery.dataViewer"
	policy.Add(member, iam.RoleName(role))

	bqTable := bqClient.Dataset(datasetID).Table(tableID)
	bqTable.IAM().SetPolicy(ctx, policy)

	return nil
}

func RemoveMemberFromBigQueryTable(ctx context.Context, projectID, datasetID, tableID, member string) error {
	bqClient, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer bqClient.Close()

	policy, err := getPolicy(ctx, bqClient, datasetID, tableID)
	if err != nil {
		return fmt.Errorf("Getting policy for %v.%v in %v: %v", datasetID, tableID, projectID, err)
	}

	// no support for V3 for BigQuery yet, and no support for conditions
	role := "roles/bigquery.dataViewer"
	policy.Remove(member, iam.RoleName(role))

	bqTable := bqClient.Dataset(datasetID).Table(tableID)
	bqTable.IAM().SetPolicy(ctx, policy)

	return nil
}

func getPolicy(ctx context.Context, bqclient *bigquery.Client, datasetID, tableID string) (*iam.Policy, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	var dataset = bqclient.Dataset(datasetID)
	var table *bigquery.Table = dataset.Table(tableID)
	policy, err := table.IAM().Policy(ctx)
	if err != nil {
		return nil, fmt.Errorf("Getting policy for %v.%v: %v", datasetID, tableID, err)
	}

	return policy, nil
}

func CheckAccessInBigQueryTable(ctx context.Context, projectID, datasetID, tableID, member string) (bool, error) {
	bqClient, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return false, fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer bqClient.Close()

	policy, err := getPolicy(ctx, bqClient, datasetID, tableID)
	if err != nil {
		return false, fmt.Errorf("Getting policy for %v.%v in %v: %v", datasetID, tableID, projectID, err)
	}

	// no support for V3 for BigQuery yet, and no support for conditions
	role := "roles/bigquery.dataViewer"
	return policy.HasRole(member, iam.RoleName(role)), nil
}
