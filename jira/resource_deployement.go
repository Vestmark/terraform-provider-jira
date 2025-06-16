package jira

import (
	"fmt"
	"log"
	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	//"github.com/pkg/errors"
	"bytes"
	"encoding/json"
	"net/http"
)

type JiraDeployment struct {
	EnvironmentId string `json:"environmentId"`
	EnvironmentName string `json:"environmentName"`
	EnvironmentType string `json:"environmentType"`	
	IssueKeys []string	`json:"issueKeys"`
}

func resourceDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateDeployment,
		Read: 	resourceReadDeployment,
		Delete: resourceDeleteDeployment,
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"issue_keys": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceCreateDeployment(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	environmentId 	:= d.Get("environment_id").(string)
	environmentName := d.Get("environment_name").(string)
	environmentType := d.Get("environment_type").(string)
	issueKeys       := d.Get("issue_keys").([]string)

	jiraDeployment := JiraDeployment{
		EnvironmentId 	: environmentId,
		EnvironmentName : environmentName,
		EnvironmentType : environmentType,
		IssueKeys 		: issueKeys,
	}
	err := sendDeploymentToJira(config.jiraClient, jiraDeployment)
	if err != nil{
		return fmt.Errorf("failed to send the deployment: %d", err)
	}
	return nil
}

func sendDeploymentToJira(jiraClient *jira.Client, deployment JiraDeployment) error {
	const BaseURL = "https://vestmark.atlassian.net"
	url := BaseURL + "/rest/deployments/0.1/bulk" 			//fmt.Sprintf("/rest/deployments/0.1/bulk")
	jsonData, err := json.Marshal([]JiraDeployment{deployment})
	if err!= nil {
		return err
	}
	log.Println("===DEBUG:PAYLOAD===")
	log.Println(string(jsonData))
	log.Println("======")
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
	fmt.Println(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK{
		return fmt.Errorf("failed to send deployment info to JIRA, status code: %d", resp.StatusCode)
	}
	return nil
}

func resourceReadDeployment(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDeleteDeployment(d *schema.ResourceData, m interface{}) error {
	return nil
}



