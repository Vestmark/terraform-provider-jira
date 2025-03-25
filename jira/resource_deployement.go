package jira

import {
	"fmt"
	"log"
	"time"
	"context"
	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
}

const jiraBaseURL   =  "https://vestmark.atlassian.net"
const jiraAPIToken  =  "api_token_from_vault"
const jiraUserEmail =	"jenkinsadmin@vestmark.com"

type JiraDeployment struct {
	EnvironmentID string `json:"environmentId`
	EnvironmentName string `json:"environmentName`
	EnvironmentType string `json:environmentType`	
	IssueKeys []string		`json:issueKeys`
}

func sendDeploymentToJira(deployment JiraDeployment) error {
	url := fmt.Sprintf("%s/rest/deployments/0.1/bulk", jiraBaseURL)
	jsonData, err := json.Marshal([]JiraDeployment{deployment})
	if err!= nil {
		return err
	}

	req, err:= http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil{
		return err
	}

	req.Header.Set("Authorization", "Basic " + jiraAPIToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.client{}
	resp, err := client.Do(req)
	if err !=nil{
		return err
	} 
	 
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK{
		return fmt.Errorf("failed to send deployment info to JIRA, status code: %d", resp.StatusCode)
	}
	return nil
}

func resourceCreateDeployment(d *schema.ResourceData, m interface{}) error {
	environmentId := d.Get("environmentId").(string)
	environmentName := d.Get("environmentName").(string)
	environmentType := d.Get("environmentType").(string)
	issueKeys       := d.Get("issueKeys").([]string)

	
	jiraDeployment := JiraDeployment{
		EnvironmentID := environmentId
		EnvironmentName := environmentName
		EnvironmentType := environmentType
	}

 
	err = sendDeploymentToJira(jiraDeployment)
	if err != nil{
		return fmt.Errorf("Failed to send the deployment: %d", err)
	}

	return nil
}

func resourceDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateDeployment,

		Schema: map[string]*schema.Schema{
			"environmentId": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environmentName": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environmentType": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issueKeys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

