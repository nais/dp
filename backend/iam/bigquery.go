package iam

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/iam"
	"context"
	"fmt"
	"time"
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
		return err
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
		return err
	}
	return nil
}

func UpdateBigqueryTableAccessControl(projectID, datasetID, tableID, member string) error {
	ctx := context.Background()
	bqClient, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer bqClient.Close()

	policy, err := getPolicy(bqClient, datasetID, tableID)

	// no support for V3 for BigQuery yet, and no support for conditions
	role := "roles/bigquery.dataViewer"
	userMember := "user:" + member
	policy.Add(userMember, iam.RoleName(role))

	bqTable := bqClient.Dataset(datasetID).Table(tableID)
	bqTable.IAM().SetPolicy(ctx, policy)

	return nil
}

func getPolicy(bqclient *bigquery.Client, datasetID, tableID string) (*iam.Policy, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	var dataset = bqclient.Dataset(datasetID)
	var table *bigquery.Table = dataset.Table(tableID)
	policy, err := table.IAM().Policy(ctx)
	if err != nil {
		return nil, err
	}

	return policy, nil
}
