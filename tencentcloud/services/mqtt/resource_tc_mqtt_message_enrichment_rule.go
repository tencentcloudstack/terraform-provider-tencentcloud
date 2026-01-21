package mqtt

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mqttv20240516 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mqtt/v20240516"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMqttMessageEnrichmentRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMqttMessageEnrichmentRuleCreate,
		Read:   resourceTencentCloudMqttMessageEnrichmentRuleRead,
		Update: resourceTencentCloudMqttMessageEnrichmentRuleUpdate,
		Delete: resourceTencentCloudMqttMessageEnrichmentRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "MQTT instance ID.",
			},

			"rule_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule name, 3-64 characters, supports Chinese, letters, numbers, `-` and `_`.",
			},

			"condition": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Rule matching condition.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "User name.",
						},
						"client_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Client ID.",
						},
						"topic": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Topic.",
						},
					},
				},
			},

			"actions": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Rule execution actions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message_expiry_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Message expiration interval.",
						},
						"response_topic": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Response Topic.",
						},
						"correlation_data": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Correlation Data.",
						},
						"user_property": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "User Properties.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Value.",
									},
								},
							},
						},
					},
				},
			},

			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Rule priority, smaller number means higher priority.",
			},

			"status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Policy status, 0: undefined; 1: active; 2: inactive, default is 2.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remark information. not exceeding 128 characters in length.",
			},

			// computed
			"rule_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rule ID.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time, millisecond timestamp.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time, millisecond timestamp.",
			},
		},
	}
}

func resourceTencentCloudMqttMessageEnrichmentRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_message_enrichment_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = mqttv20240516.NewCreateMessageEnrichmentRuleRequest()
		response   = mqttv20240516.NewCreateMessageEnrichmentRuleResponse()
		instanceId string
		ruleId     string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("rule_name"); ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("condition"); ok {
		conditionList := v.([]interface{})
		for _, item := range conditionList {
			if item == nil {
				continue
			}

			conditionMap := item.(map[string]interface{})
			conditionData := map[string]interface{}{}
			if username, exists := conditionMap["username"]; exists {
				conditionData["username"] = username.(string)
			}

			if clientId, exists := conditionMap["client_id"]; exists {
				conditionData["clientId"] = clientId.(string)
			}

			if topic, exists := conditionMap["topic"]; exists {
				conditionData["topic"] = topic.(string)
			}

			conditionJSON, err := json.Marshal(conditionData)
			if err != nil {
				return fmt.Errorf("failed to marshal condition: %v", err)
			}

			conditionBase64 := base64.StdEncoding.EncodeToString(conditionJSON)
			request.Condition = helper.String(conditionBase64)
		}
	}

	if v, ok := d.GetOk("actions"); ok {
		actionsList := v.([]interface{})
		for _, item := range actionsList {
			if item == nil {
				continue
			}

			actionsMap := item.(map[string]interface{})
			actionsData := map[string]interface{}{}
			if messageExpiryInterval, exists := actionsMap["message_expiry_interval"]; exists {
				actionsData["messageExpiryInterval"] = messageExpiryInterval.(int)
			}

			if responseTopic, exists := actionsMap["response_topic"]; exists {
				actionsData["responseTopic"] = responseTopic.(string)
			}

			if correlationData, exists := actionsMap["correlation_data"]; exists {
				actionsData["correlationData"] = correlationData.(string)
			}

			if userProperty, exists := actionsMap["user_property"]; exists {
				userPropertyList := userProperty.([]interface{})
				var userProperties []map[string]interface{}
				for _, prop := range userPropertyList {
					propMap := prop.(map[string]interface{})
					userProperties = append(userProperties, map[string]interface{}{
						"key":   propMap["key"].(string),
						"value": propMap["value"].(string),
					})
				}

				if len(userProperties) > 0 {
					actionsData["userProperty"] = userProperties
				}
			}

			actionsJSON, err := json.Marshal(actionsData)
			if err != nil {
				return fmt.Errorf("failed to marshal actions: %v", err)
			}

			actionsBase64 := base64.StdEncoding.EncodeToString(actionsJSON)
			request.Actions = helper.String(actionsBase64)
		}
	}

	if v, ok := d.GetOkExists("priority"); ok {
		request.Priority = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("status"); ok {
		request.Status = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().CreateMessageEnrichmentRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create mqtt message enrichment rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create mqtt message enrichment rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Id == nil {
		return fmt.Errorf("Id is nil.")
	}

	ruleId = helper.Int64ToStr(*response.Response.Id)
	d.SetId(instanceId + tccommon.FILED_SP + ruleId)

	return resourceTencentCloudMqttMessageEnrichmentRuleRead(d, meta)
}

func resourceTencentCloudMqttMessageEnrichmentRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_message_enrichment_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = MqttService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	ruleId := idSplit[1]

	respData, err := service.DescribeMqttMessageEnrichmentRuleById(ctx, instanceId, ruleId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_mqtt_message_enrichment_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.InstanceId != nil {
		_ = d.Set("instance_id", respData.InstanceId)
	}

	if respData.RuleName != nil {
		_ = d.Set("rule_name", respData.RuleName)
	}

	if respData.Condition != nil {
		// Decode Base64 encoded JSON condition
		conditionJSON, err := base64.StdEncoding.DecodeString(*respData.Condition)
		if err != nil {
			log.Printf("[WARN] Failed to decode condition Base64: %v", err)
		} else {
			var conditionData map[string]interface{}
			if err := json.Unmarshal(conditionJSON, &conditionData); err != nil {
				log.Printf("[WARN] Failed to unmarshal condition JSON: %v", err)
			} else {
				conditionList := []map[string]interface{}{{}}
				if username, exists := conditionData["username"]; exists {
					conditionList[0]["username"] = username
				}

				if clientId, exists := conditionData["clientId"]; exists {
					conditionList[0]["client_id"] = clientId
				}

				if topic, exists := conditionData["topic"]; exists {
					conditionList[0]["topic"] = topic
				}

				_ = d.Set("condition", conditionList)
			}
		}
	}

	if respData.Actions != nil {
		// Decode Base64 encoded JSON actions
		actionsJSON, err := base64.StdEncoding.DecodeString(*respData.Actions)
		if err != nil {
			log.Printf("[WARN] Failed to decode actions Base64: %v", err)
		} else {
			var actionsData map[string]interface{}
			if err := json.Unmarshal(actionsJSON, &actionsData); err != nil {
				log.Printf("[WARN] Failed to unmarshal actions JSON: %v", err)
			} else {
				actionsList := []map[string]interface{}{{}}
				if messageExpiryInterval, exists := actionsData["messageExpiryInterval"]; exists {
					if val, ok := messageExpiryInterval.(float64); ok {
						actionsList[0]["message_expiry_interval"] = int(val)
					}
				}

				if responseTopic, exists := actionsData["responseTopic"]; exists {
					actionsList[0]["response_topic"] = responseTopic
				}

				if correlationData, exists := actionsData["correlationData"]; exists {
					actionsList[0]["correlation_data"] = correlationData
				}

				if userProperty, exists := actionsData["userProperty"]; exists {
					if userProps, ok := userProperty.([]interface{}); ok {
						var userPropertyList []map[string]interface{}
						for _, prop := range userProps {
							if propMap, ok := prop.(map[string]interface{}); ok {
								userPropertyList = append(userPropertyList, map[string]interface{}{
									"key":   propMap["key"],
									"value": propMap["value"],
								})
							}
						}

						if len(userPropertyList) > 0 {
							actionsList[0]["user_property"] = userPropertyList
						}
					}
				}

				_ = d.Set("actions", actionsList)
			}
		}
	}

	if respData.Priority != nil {
		_ = d.Set("priority", respData.Priority)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.Remark != nil {
		_ = d.Set("remark", respData.Remark)
	}

	if respData.Id != nil {
		_ = d.Set("rule_id", respData.Id)
	}

	if respData.CreatedTime != nil {
		_ = d.Set("create_time", helper.Int64ToStr(*respData.CreatedTime))
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", helper.Int64ToStr(*respData.UpdateTime))
	}

	return nil
}

func resourceTencentCloudMqttMessageEnrichmentRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_message_enrichment_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = mqttv20240516.NewModifyMessageEnrichmentRuleRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	ruleId := idSplit[1]

	needChange := false
	if d.HasChange("rule_name") {
		needChange = true
		if v, ok := d.GetOk("rule_name"); ok {
			request.RuleName = helper.String(v.(string))
		}
	}

	if d.HasChange("condition") {
		needChange = true
		if v, ok := d.GetOk("condition"); ok {
			conditionList := v.([]interface{})
			for _, item := range conditionList {
				if item == nil {
					continue
				}

				conditionMap := item.(map[string]interface{})
				conditionData := map[string]interface{}{}
				if username, exists := conditionMap["username"]; exists {
					conditionData["username"] = username.(string)
				}

				if clientId, exists := conditionMap["client_id"]; exists {
					conditionData["clientId"] = clientId.(string)
				}

				if topic, exists := conditionMap["topic"]; exists {
					conditionData["topic"] = topic.(string)
				}

				conditionJSON, err := json.Marshal(conditionData)
				if err != nil {
					return fmt.Errorf("failed to marshal condition: %v", err)
				}

				conditionBase64 := base64.StdEncoding.EncodeToString(conditionJSON)
				request.Condition = helper.String(conditionBase64)
			}
		}
	}

	if d.HasChange("actions") {
		needChange = true
		if v, ok := d.GetOk("actions"); ok {
			actionsList := v.([]interface{})
			for _, item := range actionsList {
				if item == nil {
					continue
				}

				actionsMap := item.(map[string]interface{})
				actionsData := map[string]interface{}{}
				if messageExpiryInterval, exists := actionsMap["message_expiry_interval"]; exists {
					actionsData["messageExpiryInterval"] = messageExpiryInterval.(int)
				}

				if responseTopic, exists := actionsMap["response_topic"]; exists {
					actionsData["responseTopic"] = responseTopic.(string)
				}

				if correlationData, exists := actionsMap["correlation_data"]; exists {
					actionsData["correlationData"] = correlationData.(string)
				}

				if userProperty, exists := actionsMap["user_property"]; exists {
					userPropertyList := userProperty.([]interface{})
					var userProperties []map[string]interface{}
					for _, prop := range userPropertyList {
						propMap := prop.(map[string]interface{})
						userProperties = append(userProperties, map[string]interface{}{
							"key":   propMap["key"].(string),
							"value": propMap["value"].(string),
						})
					}

					if len(userProperties) > 0 {
						actionsData["userProperty"] = userProperties
					}
				}

				actionsJSON, err := json.Marshal(actionsData)
				if err != nil {
					return fmt.Errorf("failed to marshal actions: %v", err)
				}

				actionsBase64 := base64.StdEncoding.EncodeToString(actionsJSON)
				request.Actions = helper.String(actionsBase64)
			}
		}
	}

	if d.HasChange("priority") {
		needChange = true
		if v, ok := d.GetOkExists("priority"); ok {
			request.Priority = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("status") {
		needChange = true
		if v, ok := d.GetOkExists("status"); ok {
			request.Status = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("remark") {
		needChange = true
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	if needChange {
		request.InstanceId = helper.String(instanceId)
		request.Id = helper.StrToInt64Point(ruleId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().ModifyMessageEnrichmentRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update mqtt message enrichment rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudMqttMessageEnrichmentRuleRead(d, meta)
}

func resourceTencentCloudMqttMessageEnrichmentRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_message_enrichment_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = MqttService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	ruleId := idSplit[1]

	if err := service.DeleteMqttMessageEnrichmentRuleById(ctx, instanceId, ruleId); err != nil {
		return err
	}

	return nil
}
