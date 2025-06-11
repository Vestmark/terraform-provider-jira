package jira

import (
	"fmt"
	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	//"github.com/pkg/errors"
	"encoding/json"
	"net/http"
	"bytes"
)

type JiraDeployment struct {
	EnvironmentID string `json:"environmentId`
	EnvironmentName string `json:"environmentName`
	EnvironmentType string `json:environmentType`	
	IssueKeys []string	`json:issueKeys`
}

func sendDeploymentToJira(jiraClient *jira.Client, deployment JiraDeployment) error {
	const BaseURL = "https://vestmark.atlassian.net"
	url := BaseURL + "/rest/deployments/0.1/bulk" 			//fmt.Sprintf("/rest/deployments/0.1/bulk")
	jsonData, err := json.Marshal([]JiraDeployment{deployment})
	if err!= nil {
		return err
	}
	req, err:= jiraClient.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil{
		return err
	}
	//req.Header.Set("Authorization", "Basic " + jiraAPIToken)
	req.Header.Set("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := jiraClient.Do(req,nil)
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
	config := m.(*Config)
	environmentId 	:= d.Get("environmentId").(string)
	environmentName := d.Get("environmentName").(string)
	environmentType := d.Get("environmentType").(string)
	issueKeys       := d.Get("issueKeys").([]string)

	jiraDeployment := JiraDeployment{
		EnvironmentID : environmentId,
		EnvironmentName : environmentName,
		EnvironmentType : environmentType,
		IssueKeys : issueKeys,
	}
	err := sendDeploymentToJira(config.jiraClient, jiraDeployment)
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

