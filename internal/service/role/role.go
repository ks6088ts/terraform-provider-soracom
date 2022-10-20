package role

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ks6088ts/terraform-provider-soracom/internal/conns"

	soracom "github.com/ks6088ts/soracom-sdk-go/generated/api"
)

func ResourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,
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
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permission": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	operatorId := d.Get("operator_id").(string)
	roleId := d.Get("role_id").(string)
	description := d.Get("description").(string)
	permission := d.Get("permission").(string)
	request := soracom.CreateOrUpdateRoleRequest{
		Description: &description,
		Permission:  permission,
	}

	client := meta.(*conns.SoracomClient)
	apiResponse, httpResponse, err := client.Client.RoleApi.CreateRole(client.GetContext(ctx), operatorId, roleId).CreateOrUpdateRoleRequest(request).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error: %v, %v, %v", apiResponse, httpResponse, err))
	}
	if httpResponse.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	d.SetId(*apiResponse.RoleId)

	return diags
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	roleId := d.Id()
	operatorId := d.Get("operator_id").(string)

	client := meta.(*conns.SoracomClient)
	apiResponse, httpResponse, err := client.Client.RoleApi.GetRole(client.GetContext(ctx), operatorId, roleId).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if httpResponse.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	if err := d.Set("operator_id", operatorId); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "failed to set operator_id",
			Detail:   fmt.Sprintf("operator_id=%v, err=%v", operatorId, err),
		})
	}
	if err := d.Set("role_id", apiResponse.RoleId); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "failed to set role_id",
			Detail:   fmt.Sprintf("role_id=%v, err=%v", apiResponse.RoleId, err),
		})
	}
	if err := d.Set("description", apiResponse.Description); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "failed to set description",
			Detail:   fmt.Sprintf("description=%v, err=%v", apiResponse.Description, err),
		})
	}
	if err := d.Set("permission", apiResponse.Permission); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "failed to set permission",
			Detail:   fmt.Sprintf("permission=%v, err=%v", apiResponse.Permission, err),
		})
	}

	return diags
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	operatorId := d.Get("operator_id").(string)
	roleId := d.Id()
	description := d.Get("description").(string)
	permission := d.Get("permission").(string)
	request := soracom.CreateOrUpdateRoleRequest{
		Description: &description,
		Permission:  permission,
	}

	client := meta.(*conns.SoracomClient)
	httpResponse, err := client.Client.RoleApi.UpdateRole(client.GetContext(ctx), operatorId, roleId).CreateOrUpdateRoleRequest(request).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error: %v, %v", httpResponse, err))
	}
	if httpResponse.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	d.SetId(roleId)
	return resourceRoleRead(ctx, d, meta)
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	roleId := d.Id()
	operatorId := d.Get("operator_id").(string)
	client := meta.(*conns.SoracomClient)
	httpResponse, err := client.Client.RoleApi.DeleteRole(client.GetContext(ctx), operatorId, roleId).Execute()
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
