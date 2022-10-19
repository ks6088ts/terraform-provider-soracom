package sim

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ks6088ts/terraform-provider-soracom/internal/conns"

	soracom "github.com/ks6088ts/soracom-sdk-go/openapiclient"
)

func ResourceConfigurationGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConfigurationGroupCreate,
		ReadContext:   resourceConfigurationGroupRead,
		UpdateContext: resourceConfigurationGroupUpdate,
		DeleteContext: resourceConfigurationGroupDelete,
		Importer:      nil, // fixme: impl
		Schema: map[string]*schema.Schema{
			"sim_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sim_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceConfigurationGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	simId := d.Get("sim_id").(string)
	simGroupId := d.Get("sim_group_id").(string)

	request := soracom.SetGroupRequest{
		GroupId: &simGroupId,
	}

	client := meta.(*conns.SoracomClient)
	apiResponse, httpResponse, err := client.Client.SimApi.SetSimGroup(client.GetContext(ctx), simId).SetGroupRequest(request).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if httpResponse.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	d.SetId(*apiResponse.SimId)

	return diags
}

func resourceConfigurationGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	simId := d.Id()
	client := meta.(*conns.SoracomClient)
	apiResponse, httpResponse, err := client.Client.SimApi.GetSim(client.GetContext(ctx), simId).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if httpResponse.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}
	simGroupId := *apiResponse.GroupId

	if err := d.Set("sim_id", simId); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "failed to set sim_id",
			Detail:   fmt.Sprintf("sim_id=%v, err=%v", simId, err),
		})
	}
	if err := d.Set("sim_group_id", simGroupId); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "failed to set sim_group_id",
			Detail:   fmt.Sprintf("sim_group_id=%v, err=%v", simGroupId, err),
		})
	}

	return diags
}

func resourceConfigurationGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	simId := d.Id()
	simGroupId := d.Get("sim_group_id").(string)

	request := soracom.SetGroupRequest{
		GroupId: &simGroupId,
	}

	client := meta.(*conns.SoracomClient)
	_, httpResponse, err := client.Client.SimApi.SetSimGroup(client.GetContext(ctx), simId).SetGroupRequest(request).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if httpResponse.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	return resourceConfigurationGroupRead(ctx, d, meta)
}

func resourceConfigurationGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	simId := d.Id()
	simGroupId := d.Get("sim_group_id").(string)
	client := meta.(*conns.SoracomClient)
	_, httpResponse, err := client.Client.SimApi.UnsetSimGroup(client.GetContext(ctx), simId).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "failed to unset sim group",
			Detail:   fmt.Sprintf("simId=%v, simGroupId=%v, err=%v", simId, simGroupId, err),
		})
	}

	if httpResponse.StatusCode != 204 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "failed to unset sim group",
			Detail:   fmt.Sprintf("simId=%v, simGroupId=%v, statusCode=%v", simId, simGroupId, httpResponse.StatusCode),
		})
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
