package group

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ks6088ts/terraform-provider-soracom/internal/conns"
	sdk "github.com/soracom/soracom-sdk-go"
)

func ResourceGroupConfigurationFunk() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupConfigurationFunkCreate,
		ReadContext:   resourceGroupConfigurationFunkRead,
		UpdateContext: resourceGroupConfigurationFunkUpdate,
		DeleteContext: resourceGroupConfigurationFunkDelete,
		Importer:      nil, // fixme: impl
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"destination": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider": {
							Type:     schema.TypeString,
							Required: true,
						},
						"service": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_url": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"credentials_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGroupConfigurationFunkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	groupId := d.Get("group_id").(string)
	enabled := d.Get("enabled").(bool)
	_destination := d.Get("destination").([]interface{})
	destination := _destination[0].(map[string]interface{})
	provider := destination["provider"].(string)
	service := destination["service"].(string)
	resourceUrl := destination["resource_url"].(string)
	credentialsId := d.Get("credentials_id").(string)
	contentType := d.Get("content_type").(string)

	ac := sdk.NewAPIClient(nil)
	ac.APIKey = *meta.(*conns.SoracomClient).AuthResponse.ApiKey
	ac.Token = *meta.(*conns.SoracomClient).AuthResponse.Token

	type funkDestination struct {
		Provider    string `json:"provider"`
		Service     string `json:"service"`
		ResourceUrl string `json:"resourceUrl"`
	}
	type funkCredentialsId struct {
		CredentialsId string `json:"$credentialsId"`
	}
	response, err := ac.UpdateGroupConfigurations(groupId, "SoracomFunk", []sdk.GroupConfig{
		{
			Key:   "enabled",
			Value: enabled,
		},
		{
			Key: "destination",
			Value: funkDestination{
				Provider:    provider,
				Service:     service,
				ResourceUrl: resourceUrl,
			},
		},
		{
			Key: "credentialsId",
			Value: funkCredentialsId{
				CredentialsId: credentialsId,
			},
		},
		{
			Key:   "contentType",
			Value: contentType,
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(response.GroupID)

	return diags
}

func resourceGroupConfigurationFunkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// fixme: impl
	return diags
}

func resourceGroupConfigurationFunkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceGroupConfigurationFunkRead(ctx, d, meta)
}

func resourceGroupConfigurationFunkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	groupId := d.Id()
	client := meta.(*conns.SoracomClient)
	httpResponse, err := client.Client.GroupApi.DeleteConfigurationNamespace(client.GetContext(ctx), groupId, "SoracomFunk").Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if httpResponse.StatusCode != 204 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
