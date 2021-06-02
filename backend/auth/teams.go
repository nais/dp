package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/nais/dp/backend/config"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
	"net/http"
)

type teamModule struct {
	Name          *string `hcl:"name"`
	NaisProjectId *string `hcl:"naisProjectId"`
}
type tfFile struct {
	TeamModules []*teamModule `hcl:"module,block"`
}

func UpdateTeams(c context.Context, mapToUpdate map[string]string, teamsURL, teamsToken string, updateFrequency time.Duration) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := fetchTeams(c, mapToUpdate, teamsURL, teamsToken); err != nil {
				log.Errorf("Fetching teams from url: %v: %v", teamsURL, err)
			}
			ticker.Reset(updateFrequency)
		case <-c.Done():
			return
		}
	}
}

func UpdateTeamProjects(c context.Context, mapToUpdate map[string][]string, config config.Config, updateFrequency time.Duration) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := FetchTeamProjects(c, mapToUpdate, config); err != nil {
				log.Errorf("Fetching teams from url: %v: %v", config.TeamGCPProjectsProd, err)
			}
			ticker.Reset(updateFrequency)
		case <-c.Done():
			return
		}
	}
}

func FetchTeamProjects(c context.Context, teamProjects map[string][]string, config config.Config) error {
	request, err := http.NewRequestWithContext(c, http.MethodGet, config.TeamGCPProjectsDev, nil)
	if err != nil {
		log.Fatalf("Creating http request for teams: %v", err) // assumed unrecoverable
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", config.TeamsToken))
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("performing http request on teams json URL: %v: %w", config.TeamGCPProjectsDev, err)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	var teamModules tfFile
	if err := hclsimple.Decode("noetull.tf", bytes, nil, &teamModules); err != nil {
		fmt.Printf("sumtn wrong parsing: %v", err)
	}

	fmt.Printf("%v", teamModules)
	return nil
}

func fetchTeams(c context.Context, mapToUpdate map[string]string, teamsURL, teamsToken string) error {
	request, err := http.NewRequestWithContext(c, http.MethodGet, teamsURL, nil)
	if err != nil {
		log.Fatalf("Creating http request for teams: %v", err) // assumed unrecoverable
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", teamsToken))
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("performing http request on teams json URL: %v: %w", teamsURL, err)
	}

	if err := json.NewDecoder(response.Body).Decode(&mapToUpdate); err != nil {
		return fmt.Errorf("unmarshalling response from teams json URL: %v: %w", teamsURL, err)
	}

	log.Infof("Updated UUID mapping: %d teams from %v", len(mapToUpdate), teamsURL)
	return nil
}
