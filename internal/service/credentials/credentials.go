package credentials

import (
	"context"
	"fmt"

	soracom "github.com/ks6088ts/soracom-sdk-go/openapiclient"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ks6088ts/terraform-provider-soracom/internal/conns"
)

func ResourceCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCredentialsCreate,
		ReadContext:   resourceCredentialsRead,
		UpdateContext: resourceCredentialsUpdate,
		DeleteContext: resourceCredentialsDelete,
		Importer:      nil, // fixme: impl
		Schema: map[string]*schema.Schema{
			"credentials_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"credentials": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	credentialsId := d.Get("credentials_id").(string)
	description := d.Get("description").(string)
	type_ := d.Get("type").(string)
	credentials := d.Get("credentials").(map[string]interface{})
	createAndUpdateCredentialsModel := soracom.CreateAndUpdateCredentialsModel{
		Description: &description,
		Type:        &type_,
		Credentials: credentials,
	}

	client := meta.(*conns.SoracomClient)
	apiResponse, httpResponse, err := client.Client.CredentialApi.CreateCredential(client.GetContext(ctx), credentialsId).CreateAndUpdateCredentialsModel(createAndUpdateCredentialsModel).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if httpResponse.StatusCode != 201 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	d.SetId(*apiResponse.CredentialsId)

	return diags
}

func resourceCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*conns.SoracomClient)
	apiResponse, httpResponse, err := client.Client.CredentialApi.ListCredentials(client.GetContext(ctx)).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if httpResponse.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	credentialsId := d.Id()
	for _, credentialsModel := range apiResponse {
		if *credentialsModel.CredentialsId != credentialsId {
			continue
		}
		if err := d.Set("description", credentialsModel.Description); err != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "failed to set description",
				Detail:   fmt.Sprintf("description=%v, err=%v", credentialsModel.Description, err),
			})
		}
		if err := d.Set("type", credentialsModel.Type); err != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "failed to set type",
				Detail:   fmt.Sprintf("type=%v, err=%v", credentialsModel.Type, err),
			})
		}
		// ignore credentials since credentials cannot be obtained via API
	}
	return diags
}

func resourceCredentialsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceCredentialsRead(ctx, d, meta)
}

func resourceCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	credentialsId := d.Id()
	client := meta.(*conns.SoracomClient)
	httpResponse, err := client.Client.CredentialApi.DeleteCredential(client.GetContext(ctx), credentialsId).Execute()
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
