package group

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	soracom "github.com/ks6088ts/soracom-sdk-go/openapiclient"
	sdk "github.com/soracom/soracom-sdk-go"

	"github.com/ks6088ts/terraform-provider-soracom/internal/conns"
)

func ResourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer:      nil, // fixme: impl
		Schema: map[string]*schema.Schema{
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

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	tags_ := d.Get("tags").(map[string]interface{})
	tags := make(map[string]string)
	for key, value := range tags_ {
		tags[key] = value.(string)
	}

	createGroupRequest := soracom.CreateGroupRequest{
		Tags: &tags,
	}

	client := meta.(*conns.SoracomClient)
	apiResponse, httpResponse, err := client.Client.GroupApi.CreateGroup(client.GetContext(ctx)).CreateGroupRequest(createGroupRequest).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if httpResponse.StatusCode != 201 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	d.SetId(*apiResponse.GroupId)

	return diags
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	groupId := d.Id()

	ac := sdk.NewAPIClient(nil)
	ac.APIKey = *meta.(*conns.SoracomClient).AuthResponse.ApiKey
	ac.Token = *meta.(*conns.SoracomClient).AuthResponse.Token
	group, err := ac.GetGroup(groupId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tags", group.Tags); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "failed to set tags",
			Detail:   fmt.Sprintf("tags=%v, err=%v", group.Tags, err),
		})
	}

	return diags
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChange("tags") {
		return resourceGroupRead(ctx, d, meta)
	}

	tags := d.Get("tags").(map[string]interface{})
	var tagUpdateRequest []soracom.TagUpdateRequest
	for key, value := range tags {
		tagUpdateRequest = append(tagUpdateRequest, soracom.TagUpdateRequest{
			TagName:  key,
			TagValue: value.(string),
		})
	}

	groupId := d.Id()

	client := meta.(*conns.SoracomClient)
	_, httpResponse, err := client.Client.GroupApi.PutGroupTags(client.GetContext(ctx), groupId).TagUpdateRequest(tagUpdateRequest).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if httpResponse.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	groupId := d.Id()
	client := meta.(*conns.SoracomClient)
	httpResponse, err := client.Client.GroupApi.DeleteGroup(client.GetContext(ctx), groupId).Execute()
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
