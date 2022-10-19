package eventhandler

import (
	"context"
	"fmt"

	soracom "github.com/ks6088ts/soracom-sdk-go/openapiclient"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ks6088ts/terraform-provider-soracom/internal/conns"
)

func ResourceHandlers() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHandlersCreate,
		ReadContext:   resourceHandlersRead,
		UpdateContext: resourceHandlersUpdate,
		DeleteContext: resourceHandlersDelete,
		Importer:      nil, // fixme: impl
		Schema: map[string]*schema.Schema{
			"action_config_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"properties": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"body": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"content_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"endpoint": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"execution_date_time_const": {
										Type:     schema.TypeString,
										Required: true,
									},
									"execution_offset_minutes": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"function_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"headers": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"http_method": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"message": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"parameter1": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"parameter2": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"parameter3": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"parameter4": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"parameter5": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"secret_access_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"speed_class": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"title": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"to": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"url": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"properties": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"inactive_timeout_date_const": {
										Type:     schema.TypeString,
										Required: true,
									},
									"inactive_timeout_offset_minutes": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"limit_total_amount": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"limit_total_traffic_mega_byte": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"run_once_among_target": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"target_ota_status": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"target_speed_class": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"target_status": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_imsi": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_operator_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_sim_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func str2addr(str string) *string {
	// return nil if value is empty
	if str == "" {
		return nil
	}
	return &str
}

func resourceHandlersCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	description := d.Get("description").(string)
	name := d.Get("name").(string)
	status := d.Get("status").(string)
	targetGroupId := d.Get("target_group_id").(string)
	targetImsi := d.Get("target_imsi").(string)
	targetOperatorId := d.Get("target_operator_id").(string)
	targetSimId := d.Get("target_sim_id").(string)

	actionConfigList := make([]soracom.ActionConfig, 0)
	for _, actionConfig := range d.Get("action_config_list").(*schema.Set).List() {
		ac := actionConfig.(map[string]interface{})
		properties := ac["properties"].([]interface{})
		props := properties[0].(map[string]interface{})

		accessKey := props["access_key"].(string)
		body := props["body"].(string)
		contentType := props["content_type"].(string)
		endpoint := props["endpoint"].(string)
		executionDateTimeConst := props["execution_date_time_const"].(string)
		executionOffsetMinutes := props["execution_offset_minutes"].(string)
		functionName := props["function_name"].(string)
		headers := make(map[string]interface{})
		for key, value := range props["headers"].(map[string]interface{}) {
			headers[key] = value
		}
		if len(headers) == 0 {
			headers = nil
		}
		httpMethod := props["http_method"].(string)
		message := props["message"].(string)
		parameter1 := props["parameter1"].(string)
		parameter2 := props["parameter2"].(string)
		parameter3 := props["parameter3"].(string)
		parameter4 := props["parameter4"].(string)
		parameter5 := props["parameter5"].(string)
		secretAccessKey := props["secret_access_key"].(string)
		speedClass := props["speed_class"].(string)
		title := props["title"].(string)
		to := props["to"].(string)
		url := props["url"].(string)

		actionConfigList = append(actionConfigList, soracom.ActionConfig{
			Type: ac["type"].(string),
			Properties: soracom.ActionConfigProperty{
				AccessKey:              str2addr(accessKey),
				Body:                   str2addr(body),
				ContentType:            str2addr(contentType),
				Endpoint:               str2addr(endpoint),
				ExecutionDateTimeConst: executionDateTimeConst,
				ExecutionOffsetMinutes: str2addr(executionOffsetMinutes),
				FunctionName:           str2addr(functionName),
				Headers:                headers,
				HttpMethod:             str2addr(httpMethod),
				Message:                str2addr(message),
				Parameter1:             str2addr(parameter1),
				Parameter2:             str2addr(parameter2),
				Parameter3:             str2addr(parameter3),
				Parameter4:             str2addr(parameter4),
				Parameter5:             str2addr(parameter5),
				SecretAccessKey:        str2addr(secretAccessKey),
				SpeedClass:             str2addr(speedClass),
				Title:                  str2addr(title),
				To:                     str2addr(to),
				Url:                    str2addr(url),
			},
		})
	}

	ruleConfig_ := d.Get("rule_config").([]interface{})
	rc := ruleConfig_[0].(map[string]interface{})
	properties := rc["properties"].([]interface{})
	props := properties[0].(map[string]interface{})

	inactiveTimeoutDateConst := props["inactive_timeout_date_const"].(string)
	inactiveTimeoutOffsetMinutes := props["inactive_timeout_offset_minutes"].(string)
	limitTotalAmount := props["limit_total_amount"].(string)
	limitTotalTrafficMegaByte := props["limit_total_traffic_mega_byte"].(string)
	runOnceAmongTarget := "false"
	if props["run_once_among_target"].(bool) {
		runOnceAmongTarget = "true"
	}
	targetOtaStatus := props["target_ota_status"].(string)
	targetSpeedClass := props["target_speed_class"].(string)
	targetStatus := props["target_status"].(string)

	ruleConfig := soracom.RuleConfig{
		Type: rc["type"].(string),
		Properties: soracom.RuleConfigProperty{
			InactiveTimeoutDateConst:     inactiveTimeoutDateConst,
			InactiveTimeoutOffsetMinutes: str2addr(inactiveTimeoutOffsetMinutes),
			LimitTotalAmount:             str2addr(limitTotalAmount),
			LimitTotalTrafficMegaByte:    str2addr(limitTotalTrafficMegaByte),
			RunOnceAmongTarget:           str2addr(runOnceAmongTarget),
			TargetOtaStatus:              str2addr(targetOtaStatus),
			TargetSpeedClass:             str2addr(targetSpeedClass),
			TargetStatus:                 str2addr(targetStatus),
		},
	}

	createEventHandlerRequest := soracom.CreateEventHandlerRequest{
		Description:      str2addr(description),
		ActionConfigList: actionConfigList,
		Name:             str2addr(name),
		RuleConfig:       ruleConfig,
		Status:           status,
		TargetGroupId:    str2addr(targetGroupId),
		TargetImsi:       str2addr(targetImsi),
		TargetOperatorId: str2addr(targetOperatorId),
		TargetSimId:      str2addr(targetSimId),
	}

	client := meta.(*conns.SoracomClient)
	apiResponse, httpResponse, err := client.Client.EventHandlerApi.CreateEventHandler(client.GetContext(ctx)).CreateEventHandlerRequest(createEventHandlerRequest).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if httpResponse.StatusCode != 201 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	d.SetId(apiResponse.HandlerId)

	return diags
}

func resourceHandlersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	handlerId := d.Id()
	client := meta.(*conns.SoracomClient)
	_, httpResponse, err := client.Client.EventHandlerApi.GetEventHandler(client.GetContext(ctx), handlerId).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if httpResponse.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	return diags
}

func resourceHandlersUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceHandlersRead(ctx, d, meta)
}

func resourceHandlersDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	handlerId := d.Id()
	client := meta.(*conns.SoracomClient)

	httpResponse, err := client.Client.EventHandlerApi.DeleteEventHandler(client.GetContext(ctx), handlerId).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "failed to delete an event handler",
			Detail:   fmt.Sprintf("handlerId=%v, err=%v", handlerId, err),
		})
	}
	if httpResponse.StatusCode != 204 {
		return diag.FromErr(fmt.Errorf("invalid status code: %v", httpResponse.StatusCode))
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
