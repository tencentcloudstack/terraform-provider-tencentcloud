package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudEbEventTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEbEventTargetCreate,
		Read:   resourceTencentCloudEbEventTargetRead,
		Update: resourceTencentCloudEbEventTargetUpdate,
		Delete: resourceTencentCloudEbEventTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"event_bus_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "event bus id.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "target type.",
			},

			"target_description": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "target description.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_description": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "QCS resource six-stage format, more references [resource six-stage format](https://cloud.tencent.com/document/product/598/10606).",
						},
						"scf_params": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "cloud function parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"batch_timeout": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Maximum waiting time for bulk delivery.",
									},
									"batch_event_count": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Maximum number of events for batch delivery.",
									},
									"enable_batch_delivery": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Enable batch delivery.",
									},
								},
							},
						},
						"ckafka_target_params": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Ckafka parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"topic_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The ckafka topic to deliver to.",
									},
									"retry_policy": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Required:    true,
										Description: "retry strategy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"retry_interval": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Retry Interval Unit: Seconds.",
												},
												"max_retry_attempts": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Maximum number of retries.",
												},
											},
										},
									},
								},
							},
						},
						"es_target_params": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "ElasticSearch parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"net_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "network connection type.",
									},
									"index_prefix": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "index prefix.",
									},
									"rotation_interval": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "es log rotation granularity.",
									},
									"output_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "DTS event configuration.",
									},
									"index_suffix_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "DTS index configuration.",
									},
									"index_template_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "es template type.",
									},
								},
							},
						},
					},
				},
			},

			"rule_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "event rule id.",
			},
		},
	}
}

func resourceTencentCloudEbEventTargetCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_target.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = eb.NewCreateTargetRequest()
		response   = eb.NewCreateTargetResponse()
		eventBusId string
		ruleId     string
		targetId   string
	)
	if v, ok := d.GetOk("event_bus_id"); ok {
		eventBusId = v.(string)
		request.EventBusId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "target_description"); ok {
		targetDescription := eb.TargetDescription{}
		if v, ok := dMap["resource_description"]; ok {
			targetDescription.ResourceDescription = helper.String(v.(string))
		}
		if sCFParamsMap, ok := helper.InterfaceToMap(dMap, "scf_params"); ok {
			sCFParams := eb.SCFParams{}
			if v, ok := sCFParamsMap["batch_timeout"]; ok {
				sCFParams.BatchTimeout = helper.IntInt64(v.(int))
			}
			if v, ok := sCFParamsMap["batch_event_count"]; ok {
				sCFParams.BatchEventCount = helper.IntInt64(v.(int))
			}
			if v, ok := sCFParamsMap["enable_batch_delivery"]; ok {
				sCFParams.EnableBatchDelivery = helper.Bool(v.(bool))
			}
			targetDescription.SCFParams = &sCFParams
		}
		if ckafkaTargetParamsMap, ok := helper.InterfaceToMap(dMap, "ckafka_target_params"); ok {
			ckafkaTargetParams := eb.CkafkaTargetParams{}
			if v, ok := ckafkaTargetParamsMap["topic_name"]; ok {
				ckafkaTargetParams.TopicName = helper.String(v.(string))
			}
			if retryPolicyMap, ok := helper.InterfaceToMap(ckafkaTargetParamsMap, "retry_policy"); ok {
				retryPolicy := eb.RetryPolicy{}
				if v, ok := retryPolicyMap["retry_interval"]; ok {
					retryPolicy.RetryInterval = helper.IntUint64(v.(int))
				}
				if v, ok := retryPolicyMap["max_retry_attempts"]; ok {
					retryPolicy.MaxRetryAttempts = helper.IntUint64(v.(int))
				}
				ckafkaTargetParams.RetryPolicy = &retryPolicy
			}
			targetDescription.CkafkaTargetParams = &ckafkaTargetParams
		}
		if eSTargetParamsMap, ok := helper.InterfaceToMap(dMap, "es_target_params"); ok {
			eSTargetParams := eb.ESTargetParams{}
			if v, ok := eSTargetParamsMap["net_mode"]; ok {
				eSTargetParams.NetMode = helper.String(v.(string))
			}
			if v, ok := eSTargetParamsMap["index_prefix"]; ok {
				eSTargetParams.IndexPrefix = helper.String(v.(string))
			}
			if v, ok := eSTargetParamsMap["rotation_interval"]; ok {
				eSTargetParams.RotationInterval = helper.String(v.(string))
			}
			if v, ok := eSTargetParamsMap["output_mode"]; ok {
				eSTargetParams.OutputMode = helper.String(v.(string))
			}
			if v, ok := eSTargetParamsMap["index_suffix_mode"]; ok {
				eSTargetParams.IndexSuffixMode = helper.String(v.(string))
			}
			if v, ok := eSTargetParamsMap["index_template_type"]; ok {
				eSTargetParams.IndexTemplateType = helper.String(v.(string))
			}
			targetDescription.ESTargetParams = &eSTargetParams
		}
		request.TargetDescription = &targetDescription
	}

	if v, ok := d.GetOk("rule_id"); ok {
		ruleId = v.(string)
		request.RuleId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseEbClient().CreateTarget(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create eb eventTarget failed, reason:%+v", logId, err)
		return err
	}

	targetId = *response.Response.TargetId
	d.SetId(eventBusId + FILED_SP + ruleId + FILED_SP + targetId)

	return resourceTencentCloudEbEventTargetRead(d, meta)
}

func resourceTencentCloudEbEventTargetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_target.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := EbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	eventBusId := idSplit[0]
	ruleId := idSplit[1]
	targetId := idSplit[2]

	eventTarget, err := service.DescribeEbEventTargetById(ctx, eventBusId, ruleId, targetId)
	if err != nil {
		return err
	}

	if eventTarget == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `EbEventTarget` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if eventTarget.EventBusId != nil {
		_ = d.Set("event_bus_id", eventTarget.EventBusId)
	}

	if eventTarget.Type != nil {
		_ = d.Set("type", eventTarget.Type)
	}

	if eventTarget.TargetDescription != nil {
		targetDescriptionMap := map[string]interface{}{}

		if eventTarget.TargetDescription.ResourceDescription != nil {
			targetDescriptionMap["resource_description"] = eventTarget.TargetDescription.ResourceDescription
		}

		if eventTarget.TargetDescription.SCFParams != nil {
			sCFParamsMap := map[string]interface{}{}

			if eventTarget.TargetDescription.SCFParams.BatchTimeout != nil {
				sCFParamsMap["batch_timeout"] = eventTarget.TargetDescription.SCFParams.BatchTimeout
			}

			if eventTarget.TargetDescription.SCFParams.BatchEventCount != nil {
				sCFParamsMap["batch_event_count"] = eventTarget.TargetDescription.SCFParams.BatchEventCount
			}

			if eventTarget.TargetDescription.SCFParams.EnableBatchDelivery != nil {
				sCFParamsMap["enable_batch_delivery"] = eventTarget.TargetDescription.SCFParams.EnableBatchDelivery
			}

			targetDescriptionMap["scf_params"] = []interface{}{sCFParamsMap}
		}

		if eventTarget.TargetDescription.CkafkaTargetParams != nil {
			ckafkaTargetParamsMap := map[string]interface{}{}

			if eventTarget.TargetDescription.CkafkaTargetParams.TopicName != nil {
				ckafkaTargetParamsMap["topic_name"] = eventTarget.TargetDescription.CkafkaTargetParams.TopicName
			}

			if eventTarget.TargetDescription.CkafkaTargetParams.RetryPolicy != nil {
				retryPolicyMap := map[string]interface{}{}

				if eventTarget.TargetDescription.CkafkaTargetParams.RetryPolicy.RetryInterval != nil {
					retryPolicyMap["retry_interval"] = eventTarget.TargetDescription.CkafkaTargetParams.RetryPolicy.RetryInterval
				}

				if eventTarget.TargetDescription.CkafkaTargetParams.RetryPolicy.MaxRetryAttempts != nil {
					retryPolicyMap["max_retry_attempts"] = eventTarget.TargetDescription.CkafkaTargetParams.RetryPolicy.MaxRetryAttempts
				}

				ckafkaTargetParamsMap["retry_policy"] = []interface{}{retryPolicyMap}
			}

			targetDescriptionMap["ckafka_target_params"] = []interface{}{ckafkaTargetParamsMap}
		}

		if eventTarget.TargetDescription.ESTargetParams != nil {
			eSTargetParamsMap := map[string]interface{}{}

			if eventTarget.TargetDescription.ESTargetParams.NetMode != nil {
				eSTargetParamsMap["net_mode"] = eventTarget.TargetDescription.ESTargetParams.NetMode
			}

			if eventTarget.TargetDescription.ESTargetParams.IndexPrefix != nil {
				eSTargetParamsMap["index_prefix"] = eventTarget.TargetDescription.ESTargetParams.IndexPrefix
			}

			if eventTarget.TargetDescription.ESTargetParams.RotationInterval != nil {
				eSTargetParamsMap["rotation_interval"] = eventTarget.TargetDescription.ESTargetParams.RotationInterval
			}

			if eventTarget.TargetDescription.ESTargetParams.OutputMode != nil {
				eSTargetParamsMap["output_mode"] = eventTarget.TargetDescription.ESTargetParams.OutputMode
			}

			if eventTarget.TargetDescription.ESTargetParams.IndexSuffixMode != nil {
				eSTargetParamsMap["index_suffix_mode"] = eventTarget.TargetDescription.ESTargetParams.IndexSuffixMode
			}

			if eventTarget.TargetDescription.ESTargetParams.IndexTemplateType != nil {
				eSTargetParamsMap["index_template_type"] = eventTarget.TargetDescription.ESTargetParams.IndexTemplateType
			}

			targetDescriptionMap["es_target_params"] = []interface{}{eSTargetParamsMap}
		}

		_ = d.Set("target_description", []interface{}{targetDescriptionMap})
	}

	if eventTarget.RuleId != nil {
		_ = d.Set("rule_id", eventTarget.RuleId)
	}

	return nil
}

func resourceTencentCloudEbEventTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_target.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := eb.NewUpdateTargetRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	eventBusId := idSplit[0]
	ruleId := idSplit[1]
	targetId := idSplit[2]

	request.EventBusId = &eventBusId
	request.RuleId = &ruleId
	request.TargetId = &targetId

	immutableArgs := []string{"event_bus_id", "type", "target_description", "rule_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseEbClient().UpdateTarget(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update eb eventTarget failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudEbEventTargetRead(d, meta)
}

func resourceTencentCloudEbEventTargetDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_target.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := EbService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	eventBusId := idSplit[0]
	ruleId := idSplit[1]
	targetId := idSplit[2]

	if err := service.DeleteEbEventTargetById(ctx, eventBusId, ruleId, targetId); err != nil {
		return err
	}

	return nil
}
