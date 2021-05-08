package iam_test

import (
	"github.com/nais/dp/backend/iam"
	"testing"
)

func TestBigqueryIam(t *testing.T) {
	member := "user:christine.teig@nav.no"
	projectID := "aura-dev-d9f5"
	datasetID := "container_resource_usage"
	iam.UpdateDatasetAccessControl(member, projectID, datasetID)
}

func TestBigqueryTableIam(t *testing.T) {

	member := "user:christine.teig@nav.no"
	projectID := "aura-dev-d9f5"
	datasetID := "container_resource_usage"
	tableID := "container_resource_usage"
	iam.UpdateBigqueryTableAccessControl(member, projectID, datasetID, tableID)
}

func TestBigqueryViewIam(t *testing.T) {

	member := "user:johnny.horvi@nav.no"
	projectID := "aura-dev-d9f5"
	datasetID := "container_resource_usage"
	viewID := "container_resource_usage_aura"
	iam.UpdateBigqueryViewAccessControl(member, projectID, datasetID, viewID)
}
