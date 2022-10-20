package sim

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ks6088ts/terraform-provider-soracom/internal/conns"

	soracom "github.com/ks6088ts/soracom-sdk-go/generated/api"
)

func ResourceTags() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTagsCreate,
		ReadContext:   resourceTagsRead,
		UpdateContext: resourceTagsUpdate,
		DeleteContext: resourceTagsDelete,
		Importer:      nil, // fixme: impl
		Schema: map[string]*schema.Schema{
			"sim_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func resourceTagsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	simId := d.Get("sim_id").(string)
	tags := d.Get("tags").(map[string]interface{})
	tagUpdateRequests := []soracom.TagUpdateRequest{}
	for key, value := range tags {
		tagUpdateRequests = append(tagUpdateRequests, *soracom.NewTagUpdateRequest(key, value.(string)))
	}

	client := meta.(*conns.SoracomClient)
	apiResponse, httpResponse, err := client.Client.SimApi.PutSimTags(client.GetContext(ctx), simId).TagUpdateRequest(tagUpdateRequests).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if httpResponse.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	d.SetId(*apiResponse.SimId)

	return diags
}

func resourceTagsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	tags := make(map[string]interface{})
	for key, value := range *apiResponse.Tags {
		tags[key] = value
	}

	if err := d.Set("sim_id", simId); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "failed to set sim_id",
			Detail:   fmt.Sprintf("sim_id=%v, err=%v", simId, err),
		})
	}
	if err := d.Set("tags", tags); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "failed to set tags",
			Detail:   fmt.Sprintf("tags=%v, err=%v", tags, err),
		})
	}

	return diags
}

func resourceTagsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChange("tags") {
		return resourceTagsRead(ctx, d, meta)
	}

	simId := d.Id()
	tags := d.Get("tags").(map[string]interface{})
	tagUpdateRequests := []soracom.TagUpdateRequest{}
	for key, value := range tags {
		tagUpdateRequests = append(tagUpdateRequests, *soracom.NewTagUpdateRequest(key, value.(string)))
	}

	client := meta.(*conns.SoracomClient)
	_, httpResponse, err := client.Client.SimApi.PutSimTags(client.GetContext(ctx), simId).TagUpdateRequest(tagUpdateRequests).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if httpResponse.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	return resourceTagsRead(ctx, d, meta)
}

func resourceTagsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	simId := d.Id()
	tags := d.Get("tags").(map[string]interface{})
	client := meta.(*conns.SoracomClient)
	for key, value := range tags {
		httpResponse, err := client.Client.SimApi.DeleteSimTag(client.GetContext(ctx), simId, key).Execute()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "failed to delete a tag",
				Detail:   fmt.Sprintf("simId=%v, tag=%v:%v, err=%v", simId, key, value, err),
			})
			continue
		}
		if httpResponse.StatusCode != 204 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "failed to delete a tag",
				Detail:   fmt.Sprintf("simId=%v, tag=%v:%v, statusCode=%v", simId, key, value, httpResponse.StatusCode),
			})
		}
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
