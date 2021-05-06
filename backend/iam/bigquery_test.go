package iam_test

import (
	"fmt"
	"github.com/nais/dp/backend/iam"
	"testing"
)

func TestBigqueryIam(t *testing.T) {
	fmt.Println("hello")
	fmt.Println(iam.UpdateDatasetAccessControl("christine.teig@nav.no", "aura-dev-d9f5", "container_resource_usage"))
}
