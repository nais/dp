package auth_test

import (
	"context"
	"github.com/nais/dp/backend/auth"
	"github.com/nais/dp/backend/config"
	"testing"
)

func TestValidation(t *testing.T) {
	var teamProjects map[string][]string
	auth.FetchTeamProjects(context.Background(), teamProjects, config.Config{
		TeamGCPProjectsDev:  "https://raw.githubusercontent.com/nais/teams/master/gcp-projects/dev/teams.tf",
		TeamGCPProjectsProd: "https://raw.githubusercontent.com/nais/teams/master/gcp-projects/prod/teams.tf",
		TeamsToken:          "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",
	})

	/*
		missingType := map[string]string{"no": "type"}
		err := api.ValidateDatastore(missingType)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "no type defined")
	*/
}
