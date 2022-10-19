package group

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ks6088ts/terraform-provider-soracom/internal/conns"
	sdk "github.com/soracom/soracom-sdk-go"
)

func ResourceGroupConfigurationAir() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupConfigurationAirCreate,
		ReadContext:   resourceGroupConfigurationAirRead,
		UpdateContext: resourceGroupConfigurationAirUpdate,
		DeleteContext: resourceGroupConfigurationAirDelete,
		Importer:      nil, // fixme: impl
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"use_custom_dns": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dns_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"meta_data": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"read_only": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"allow_origin": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceGroupConfigurationAirCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	groupId := d.Get("group_id").(string)
	useCustomDns := d.Get("use_custom_dns").(bool)
	_dnsServers := d.Get("dns_servers").([]interface{})
	var dnsServers []string
	for _, dnsServer := range _dnsServers {
		dnsServers = append(dnsServers, dnsServer.(string))
	}
	_metaData := d.Get("meta_data").([]interface{})
	metaData := _metaData[0].(map[string]interface{})
	metaDataEnabled := metaData["enabled"].(bool)
	metaDataReadOnly := metaData["read_only"].(bool)
	metaDataAllowOrigin := metaData["allow_origin"].(string)
	userData := d.Get("user_data").(string)
	airConfig1 := &sdk.AirConfig{
		UseCustomDNS: useCustomDns,
		DNSServers:   dnsServers,
		MetaData: sdk.MetaData{
			Enabled:     metaDataEnabled,
			ReadOnly:    metaDataReadOnly,
			AllowOrigin: metaDataAllowOrigin,
		},
		UserData: userData,
	}
	ac := sdk.NewAPIClient(nil)
	ac.APIKey = *meta.(*conns.SoracomClient).AuthResponse.ApiKey
	ac.Token = *meta.(*conns.SoracomClient).AuthResponse.Token
	response, err := ac.UpdateAirConfig(groupId, airConfig1)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(response.GroupID)

	return diags
}

func resourceGroupConfigurationAirRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// fixme: impl
	return diags
}

func resourceGroupConfigurationAirUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceGroupConfigurationAirRead(ctx, d, meta)
}

func resourceGroupConfigurationAirDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	groupId := d.Id()
	client := meta.(*conns.SoracomClient)
	httpResponse, err := client.Client.GroupApi.DeleteConfigurationNamespace(client.GetContext(ctx), groupId, "SoracomAir").Execute()
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
