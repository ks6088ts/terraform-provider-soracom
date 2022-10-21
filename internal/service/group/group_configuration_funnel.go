package group

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ks6088ts/terraform-provider-soracom/internal/conns"
	sdk "github.com/soracom/soracom-sdk-go"
)

func ResourceGroupConfigurationFunnel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupConfigurationFunnelCreate,
		ReadContext:   resourceGroupConfigurationFunnelRead,
		UpdateContext: resourceGroupConfigurationFunnelUpdate,
		DeleteContext: resourceGroupConfigurationFunnelDelete,
		Importer:      nil, // fixme: impl
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target group.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Flag for the configuration status.",
			},
			"destination": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "provider",
						},
						"service": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "service",
						},
						"resource_url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "resource_url",
						},
					},
				},
				Description: "destination",
			},
			"credentials_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "credentials_id",
			},
			"content_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "content_type",
			},
			"add_sim_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "add_sim_id",
			},
		},
	}
}

func resourceGroupConfigurationFunnelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	groupId := d.Get("group_id").(string)
	enabled := d.Get("enabled").(bool)
	credentialsId := d.Get("credentials_id").(string)
	contentType := d.Get("content_type").(string)
	addSimId := d.Get("add_sim_id").(bool)
	_destination := d.Get("destination").([]interface{})
	destination := _destination[0].(map[string]interface{})
	provider := destination["provider"].(string)
	service := destination["service"].(string)
	resourceUrl := destination["resource_url"].(string)

	ac := sdk.NewAPIClient(nil)
	ac.APIKey = *meta.(*conns.SoracomClient).AuthResponse.ApiKey
	ac.Token = *meta.(*conns.SoracomClient).AuthResponse.Token

	type funnelDestination struct {
		Provider    string `json:"provider"`
		Service     string `json:"service"`
		ResourceUrl string `json:"resourceUrl"`
	}
	response, err := ac.UpdateGroupConfigurations(groupId, "SoracomFunnel", []sdk.GroupConfig{
		{
			Key:   "enabled",
			Value: enabled,
		},
		{
			Key: "destination",
			Value: funnelDestination{
				Provider:    provider,
				Service:     service,
				ResourceUrl: resourceUrl,
			},
		},
		{
			Key:   "credentialsId",
			Value: credentialsId,
		},
		{
			Key:   "contentType",
			Value: contentType,
		},
		{
			Key:   "addSimId",
			Value: addSimId,
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(response.GroupID)

	return diags
}

func resourceGroupConfigurationFunnelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// fixme: impl
	return diags
}

func resourceGroupConfigurationFunnelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceGroupConfigurationFunnelRead(ctx, d, meta)
}

func resourceGroupConfigurationFunnelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	groupId := d.Id()
	client := meta.(*conns.SoracomClient)
	httpResponse, err := client.Client.GroupApi.DeleteConfigurationNamespace(client.GetContext(ctx), groupId, "SoracomFunnel").Execute()
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
