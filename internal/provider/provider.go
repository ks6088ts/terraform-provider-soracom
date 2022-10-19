package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ks6088ts/terraform-provider-soracom/internal"
	"github.com/ks6088ts/terraform-provider-soracom/internal/conns"
	"github.com/ks6088ts/terraform-provider-soracom/internal/service/credentials"
	"github.com/ks6088ts/terraform-provider-soracom/internal/service/eventhandler"
	"github.com/ks6088ts/terraform-provider-soracom/internal/service/group"
	"github.com/ks6088ts/terraform-provider-soracom/internal/service/role"
	"github.com/ks6088ts/terraform-provider-soracom/internal/service/sim"
	"github.com/ks6088ts/terraform-provider-soracom/internal/service/user"
)

// Provider returns a *schema.Provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"profile": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SORACOM Profile",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			// add DataSource here
		},
		ResourcesMap: map[string]*schema.Resource{
			"soracom_sim_configuration_group":    sim.ResourceConfigurationGroup(),
			"soracom_sim_tags":                   sim.ResourceTags(),
			"soracom_group":                      group.ResourceGroup(),
			"soracom_group_configuration_air":    group.ResourceGroupConfigurationAir(),
			"soracom_group_configuration_funk":   group.ResourceGroupConfigurationFunk(),
			"soracom_group_configuration_funnel": group.ResourceGroupConfigurationFunnel(),
			"soracom_eventhandler_handlers":      eventhandler.ResourceHandlers(),
			"soracom_credentials":                credentials.ResourceCredentials(),
			"soracom_role":                       role.ResourceRole(),
			"soracom_role_attachment":            role.ResourceRoleAttachment(),
			"soracom_user":                       user.ResourceUser(),
		},
		ConfigureContextFunc: providerConfigureContext,
	}
}

func providerConfigureContext(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := conns.Config{
		Profile: d.Get("profile").(string),
	}

	var diags diag.Diagnostics

	c, err := config.Client()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to create SORACOM API client @ Revision=%v, Version=%v", internal.Revision, internal.Version),
			Detail:   fmt.Sprintf("Unable to create SORACOM API client, %v", err),
		})
		return nil, diags
	}
	return c, diags
}
