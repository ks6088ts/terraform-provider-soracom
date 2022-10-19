package role

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ks6088ts/terraform-provider-soracom/internal/conns"

	soracom "github.com/ks6088ts/soracom-sdk-go/openapiclient"
)

func ResourceRoleAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleAttachmentCreate,
		ReadContext:   resourceRoleAttachmentRead,
		UpdateContext: resourceRoleAttachmentUpdate,
		DeleteContext: resourceRoleAttachmentDelete,
		Importer:      nil, // fixme: impl
		Schema: map[string]*schema.Schema{
			"operator_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRoleAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	operatorId := d.Get("operator_id").(string)
	roleId := d.Get("role_id").(string)
	userName := d.Get("user_name").(string)
	request := soracom.AttachRoleRequest{
		RoleId: &roleId,
	}

	client := meta.(*conns.SoracomClient)
	httpResponse, err := client.Client.RoleApi.AttachRole(client.GetContext(ctx), operatorId, userName).AttachRoleRequest(request).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error: %v, %v", httpResponse, err))
	}
	if httpResponse.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	d.SetId(roleId)

	return diags
}

func resourceRoleAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	roleId := d.Id()
	operatorId := d.Get("operator_id").(string)
	userName := d.Get("user_name").(string)

	client := meta.(*conns.SoracomClient)
	apiResponse, httpResponse, err := client.Client.RoleApi.ListRoleAttachedUsers(client.GetContext(ctx), operatorId, roleId).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if httpResponse.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	for _, response := range apiResponse {
		if userName != *response.UserName {
			continue
		}
		if err := d.Set("user_name", *response.UserName); err != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "failed to set user_name",
				Detail:   fmt.Sprintf("user_name=%v, err=%v", *response.UserName, err),
			})
		}
		if err := d.Set("operator_id", operatorId); err != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "failed to set operator_id",
				Detail:   fmt.Sprintf("operator_id=%v, err=%v", operatorId, err),
			})
		}
		if err := d.Set("role_id", roleId); err != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "failed to set role_id",
				Detail:   fmt.Sprintf("role_id=%v, err=%v", roleId, err),
			})
		}
	}

	return diags
}

func resourceRoleAttachmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// fixme: impl
	return resourceRoleAttachmentRead(ctx, d, meta)
}

func resourceRoleAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	operatorId := d.Get("operator_id").(string)
	roleId := d.Get("role_id").(string)
	userName := d.Get("user_name").(string)

	client := meta.(*conns.SoracomClient)
	httpResponse, err := client.Client.RoleApi.DetachRole(client.GetContext(ctx), operatorId, userName, roleId).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if httpResponse.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
