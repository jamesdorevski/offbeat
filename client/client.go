package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

// Retrieve a list of accepted worklogs for a given user and date range.
func GetWorklogs(reqData *GetWorklogsRequest) (*GetWorklogsResponse, error) {
	url := fmt.Sprintf("https://api.tempo.io/4/worklogs/user/%s?from=%s&to=%s", viper.GetString("tempo.userId"), reqData.Start, reqData.End)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+viper.GetString("tempo.apiKey"))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response GetWorklogsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Retrieve the internal Jira issue id for a given Jira issue key.
func GetIssueId(key string) (string, error) {
	instance := viper.GetString("atlassian.instance")
	email := viper.GetString("atlassian.email")
	apiKey := viper.GetString("atlassian.apiKey")
	sEnc := base64.StdEncoding.EncodeToString([]byte(email + ":" + apiKey))

	req, err := http.NewRequest("GET", instance+"/rest/api/3/issue/"+key, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Basic "+sEnc)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response Issue
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.ID, err
}

func CreateWorklog(reqData *CreateWorklogRequest) error {
	url := "https://api.tempo.io/4/worklogs"

	body := map[string]interface{}{
		"authorAccountId":  viper.GetString("tempo.userId"),
		"issueId":          reqData.IssueId,
		"startDate":        reqData.StartDate,
		"startTime":        reqData.StartTime,
		"timeSpentSeconds": reqData.TimeSpentSeconds,
	}

	json, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+viper.GetString("tempo.apiKey"))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
