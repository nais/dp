package iam_test

import (
	"context"
	"testing"

	"github.com/nais/dp/backend/iam"
)

func TestBigqueryIam(t *testing.T) {
	t.Skip()

	projectID := "aura-dev-d9f5"
	datasetID := "container_resource_usage"
	member := "user:christine.teig@nav.no"

	err := iam.UpdateDatasetAccessControl(member, projectID, datasetID)
	if err != nil {
		return
	}
}

func TestBigqueryTableIam(t *testing.T) {
	t.Skip()

	projectID := "aura-dev-d9f5"
	datasetID := "container_resource_usage"
	tableID := "container_resource_usage"
	member := "user:christine.teig@nav.no"

	err := iam.UpdateBigqueryTableAccessControl(context.Background(), member, projectID, datasetID, tableID)
	if err != nil {
		return
	}
}

func TestBigqueryViewIam(t *testing.T) {
	t.Skip()

	projectID := "aura-dev-d9f5"
	datasetID := "container_resource_usage"
	viewID := "container_resource_usage_aura"
	member := "user:johnny.horvi@nav.no"

	err := iam.UpdateBigqueryTableAccessControl(context.Background(), member, projectID, datasetID, viewID)
	if err != nil {
		return
	}
}
