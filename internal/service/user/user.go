package user

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	soracom "github.com/ks6088ts/soracom-sdk-go/openapiclient"
	"github.com/ks6088ts/terraform-provider-soracom/internal/conns"
)

func ResourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer:      nil, // fixme: impl
		Schema: map[string]*schema.Schema{
			"operator_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	operatorId := d.Get("operator_id").(string)
	userName := d.Get("user_name").(string)
	description := d.Get("description").(string)
	request := soracom.CreateUserRequest{
		Description: &description,
	}

	client := meta.(*conns.SoracomClient)
	httpResponse, err := client.Client.UserApi.CreateUser(client.GetContext(ctx), operatorId, userName).CreateUserRequest(request).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error: %v, %v", httpResponse, err))
	}
	if httpResponse.StatusCode != 201 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	// fixme: should parse response
	d.SetId(userName)

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	userName := d.Id()
	operatorId := d.Get("operator_id").(string)

	client := meta.(*conns.SoracomClient)
	apiResponse, httpResponse, err := client.Client.UserApi.GetUser(client.GetContext(ctx), operatorId, userName).Execute()
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
	if err := d.Set("user_name", apiResponse.UserName); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "failed to set user_name",
			Detail:   fmt.Sprintf("user_name=%v, err=%v", apiResponse.UserName, err),
		})
	}
	if err := d.Set("description", apiResponse.Description); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "failed to set description",
			Detail:   fmt.Sprintf("description=%v, err=%v", apiResponse.Description, err),
		})
	}

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	userName := d.Id()
	operatorId := d.Get("operator_id").(string)
	description := d.Get("description").(string)
	request := soracom.UpdateUserRequest{
		Description: &description,
	}

	client := meta.(*conns.SoracomClient)
	httpResponse, err := client.Client.UserApi.UpdateUser(client.GetContext(ctx), operatorId, userName).UpdateUserRequest(request).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error: %v, %v", httpResponse, err))
	}
	if httpResponse.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	// fixme: should parse response
	d.SetId(userName)
	return resourceUserRead(ctx, d, meta)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	userName := d.Id()
	operatorId := d.Get("operator_id").(string)
	client := meta.(*conns.SoracomClient)
	httpResponse, err := client.Client.UserApi.DeleteUser(client.GetContext(ctx), operatorId, userName).Execute()
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
