package jira

import (
	"fmt"
	"log"
	"os"
	"time"

	//jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	//"github.com/pkg/errors"
	"bytes"
	"encoding/json"
	"net/http"
)

type DeploymentPayload struct{
	Deployments []JiraDeployment `json:"deployments"`
}

type JiraDeployment struct {
	DeploymentSequenceNumber int64 	`json:"deploymentSequenceNumber"`
	UpdateSequenceNumber	 int64  `json:"updateSequenceNumber"`
	DisplayName				 string `json:"displayName"`
	URL						 string `json:"url"`
	State					 string `json:"state"`	 
	LastUpdated 			 time.Time `json:"lastUpdated"`	
	Pipeline				 Pipeline
	Environment				 Environment		
	IssueKeys []string	`json:"issueKeys"`
}

type Pipeline struct{
		ID 	 		string `json:"id"`
		DisplayName string `json:"displayName"`
		URL     	string `json:"url"`    

} 
type Environment struct {
		ID   		string `json:"id"`
		DisplayName string `json:"displayName"`
		Type 		string `json:"type"`
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
			"display_name":	{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"state": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"url":	{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCreateDeployment(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	displayName 	:= d.Get("display_name").(string)
	url				:= d.Get("url").(string)
	state 			:= d.Get("state").(string)
	environmentId 	:= d.Get("environment_id").(string)
	environmentName := d.Get("environment_name").(string)
	environmentType := d.Get("environment_type").(string)
	keys       		:= d.Get("issue_keys").([]interface{}) //To avoid runtime panic
	issueKeys 		:= make([]string, len(keys))
	for i, v := range keys {
		issueKeys[i] = v.(string)
	}

	jiraDeployment := JiraDeployment{
		DeploymentSequenceNumber: 1,
		UpdateSequenceNumber	: 1,
		DisplayName				: displayName ,
		URL						: url,
		State					: state,
		LastUpdated				: time.Now().UTC() ,
		Pipeline				: Pipeline{
			ID			: "terraform-jira-deployment",
			DisplayName	: "Terraform JIRA deployment",
			URL			: url,
		} ,
		Environment				: Environment{
			ID	: environmentId,
			DisplayName : environmentName,
			Type : environmentType,
		},
		IssueKeys 		: issueKeys,
	}

	deploymentPayload := DeploymentPayload{
		Deployments: []JiraDeployment{jiraDeployment},
	} 

	err := sendDeploymentToJira(config, deploymentPayload)
	if err != nil{
		return fmt.Errorf("failed to send the deployment: %d", err)
	}
	return nil
}

func sendDeploymentToJira(config *Config, deployment DeploymentPayload) error {
	const BaseURL = "https://vestmark.atlassian.net"
	url := BaseURL + "/rest/deployments/0.1/bulk" 			//fmt.Sprintf("/rest/deployments/0.1/bulk")
	jsonData, err := json.Marshal(deployment)
	if err!= nil {
		return err
	}
	log.SetOutput(os.Stderr)
	log.Println(deployment)
	log.Println("===DEBUG:PAYLOAD===")
	log.Println(string(jsonData))
	log.Println("======")
	req, err:= config.jiraClient.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil{
		log.Println("Failed to create new request",err)
		return err
	}
	req.Header.Set("Authorization", "Bearer" + config.token)
	req.Header.Set("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := config.jiraClient.Do(req,nil)
	if err !=nil{
		log.Println("Failed to post request",err, resp.StatusCode, resp.Body)
		return err
	} 
	log.Println("RESPONSE:-",resp.Body)
	log.Println("STATUS-CODE:-",resp.StatusCode)
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



