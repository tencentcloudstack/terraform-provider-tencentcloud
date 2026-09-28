package cls

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsClsDeliverTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsClsDeliverTaskCreate,
		Read:   resourceTencentCloudClsClsDeliverTaskRead,
		Update: resourceTencentCloudClsClsDeliverTaskUpdate,
		Delete: resourceTencentCloudClsClsDeliverTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"task_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Deliver task name.",
			},

			"source_topic_config": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Source topic config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_filter_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Topic filter type. 1: static select.",
						},
						"logset_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Source logset id.",
						},
						"topics": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Source topic list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"topic_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Log topic id.",
									},
								},
							},
						},
					},
				},
			},

			"target_topic_config": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Target topic config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Target account type. 1: current main account, 2: other main account.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Target region, e.g. ap-guangzhou.",
						},
						"logset_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Target logset id.",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Target topic id.",
						},
						"role_arn": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Role ARN, required when AccountType=2.",
						},
						"external_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "External ID, required when AccountType=2.",
						},
					},
				},
			},

			"deliver_rule": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Deliver rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_scope": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Data deliver scope. 1: history+new, 2: custom time range, 3: new only.",
						},
					},
				},
			},

			"compliance": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Compliance. 1: agree cross-domain data transfer terms.",
			},

			"has_services_log": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to deliver service log. 1: off, 2: on.",
			},

			"enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Task status. 0: run, 1: pause.",
			},

			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Deliver task id.",
			},

			"uin": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Main account id.",
			},

			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Task status. 0: running, 1: paused, 2: completed, 3: abnormal.",
			},

			"progress": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Task progress percentage.",
			},

			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Create time (unix seconds).",
			},

			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Update time (unix seconds).",
			},

			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total count of matched tasks.",
			},
		},
	}
}

func resourceTencentCloudClsClsDeliverTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cls_deliver_task.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = cls.NewCreateCLSDeliverTaskRequest()
		response = cls.NewCreateCLSDeliverTaskResponse()
	)

	if v, ok := d.GetOk("task_name"); ok {
		request.TaskName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_topic_config"); ok {
		sourceTopicConfigList := v.([]interface{})
		if len(sourceTopicConfigList) > 0 {
			sourceTopicConfig := cls.SourceTopicConfig{}
			sourceTopicConfigMap := sourceTopicConfigList[0].(map[string]interface{})
			if v, ok := sourceTopicConfigMap["topic_filter_type"].(int); ok {
				sourceTopicConfig.TopicFilterType = helper.IntUint64(v)
			}

			if v, ok := sourceTopicConfigMap["logset_id"].(string); ok && v != "" {
				sourceTopicConfig.LogsetId = helper.String(v)
			}

			if v, ok := sourceTopicConfigMap["topics"]; ok {
				for _, item := range v.([]interface{}) {
					topicsMap := item.(map[string]interface{})
					sourceTopicInfo := cls.SourceTopicInfo{}
					if v, ok := topicsMap["topic_id"].(string); ok && v != "" {
						sourceTopicInfo.TopicId = helper.String(v)
					}
					sourceTopicConfig.Topics = append(sourceTopicConfig.Topics, &sourceTopicInfo)
				}
			}

			request.SourceTopicConfig = &sourceTopicConfig
		}
	}

	if v, ok := d.GetOk("target_topic_config"); ok {
		targetTopicConfigList := v.([]interface{})
		if len(targetTopicConfigList) > 0 {
			targetTopicConfig := cls.TargetTopicConfig{}
			targetTopicConfigMap := targetTopicConfigList[0].(map[string]interface{})
			if v, ok := targetTopicConfigMap["account_type"].(int); ok {
				targetTopicConfig.AccountType = helper.IntUint64(v)
			}

			if v, ok := targetTopicConfigMap["region"].(string); ok && v != "" {
				targetTopicConfig.Region = helper.String(v)
			}

			if v, ok := targetTopicConfigMap["logset_id"].(string); ok && v != "" {
				targetTopicConfig.LogsetId = helper.String(v)
			}

			if v, ok := targetTopicConfigMap["topic_id"].(string); ok && v != "" {
				targetTopicConfig.TopicId = helper.String(v)
			}

			if v, ok := targetTopicConfigMap["role_arn"].(string); ok && v != "" {
				targetTopicConfig.RoleArn = helper.String(v)
			}

			if v, ok := targetTopicConfigMap["external_id"].(string); ok && v != "" {
				targetTopicConfig.ExternalId = helper.String(v)
			}

			request.TargetTopicConfig = &targetTopicConfig
		}
	}

	if v, ok := d.GetOk("deliver_rule"); ok {
		deliverRuleList := v.([]interface{})
		if len(deliverRuleList) > 0 {
			deliverRule := cls.DeliverRule{}
			deliverRuleMap := deliverRuleList[0].(map[string]interface{})
			if v, ok := deliverRuleMap["data_scope"].(int); ok {
				deliverRule.DataScope = helper.IntUint64(v)
			}
			request.DeliverRule = &deliverRule
		}
	}

	if v, ok := d.GetOk("compliance"); ok {
		request.Compliance = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("has_services_log"); ok {
		request.HasServicesLog = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateCLSDeliverTaskWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create cls_cls_deliver_task failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls_cls_deliver_task failed, reason:%+v", logId, err)
		return err
	}

	log.Printf("[CRUD] cls_cls_deliver_task logId=%s, d.Id()=%s, taskId=%s", logId, d.Id(), *response.Response.TaskId)

	if response.Response.TaskId == nil || *response.Response.TaskId == "" {
		return fmt.Errorf("Create cls_cls_deliver_task failed, TaskId is nil or empty.")
	}

	taskId := *response.Response.TaskId
	d.SetId(taskId)
	return resourceTencentCloudClsClsDeliverTaskRead(d, meta)
}

func resourceTencentCloudClsClsDeliverTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cls_deliver_task.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		taskId  = d.Id()
	)

	respData, err := service.DescribeClsClsDeliverTaskById(ctx, taskId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] cls_cls_deliver_task id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.TaskId != nil {
		_ = d.Set("task_id", respData.TaskId)
	}

	if respData.TaskName != nil {
		_ = d.Set("task_name", respData.TaskName)
	}

	if respData.Compliance != nil {
		_ = d.Set("compliance", respData.Compliance)
	}

	if respData.HasServicesLog != nil {
		_ = d.Set("has_services_log", respData.HasServicesLog)
	}

	if respData.Enable != nil {
		_ = d.Set("enable", respData.Enable)
	}

	if respData.Uin != nil {
		_ = d.Set("uin", respData.Uin)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.Progress != nil {
		_ = d.Set("progress", respData.Progress)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	if respData.SourceTopicConfig != nil {
		sourceTopicConfigMap := map[string]interface{}{}
		if respData.SourceTopicConfig.TopicFilterType != nil {
			sourceTopicConfigMap["topic_filter_type"] = respData.SourceTopicConfig.TopicFilterType
		}

		if respData.SourceTopicConfig.LogsetId != nil {
			sourceTopicConfigMap["logset_id"] = respData.SourceTopicConfig.LogsetId
		}

		if respData.SourceTopicConfig.Topics != nil {
			topicsList := make([]map[string]interface{}, 0, len(respData.SourceTopicConfig.Topics))
			for _, topics := range respData.SourceTopicConfig.Topics {
				topicsMap := map[string]interface{}{}
				if topics.TopicId != nil {
					topicsMap["topic_id"] = topics.TopicId
				}
				topicsList = append(topicsList, topicsMap)
			}
			sourceTopicConfigMap["topics"] = topicsList
		}

		_ = d.Set("source_topic_config", []interface{}{sourceTopicConfigMap})
	}

	if respData.TargetTopicConfig != nil {
		targetTopicConfigMap := map[string]interface{}{}
		if respData.TargetTopicConfig.AccountType != nil {
			targetTopicConfigMap["account_type"] = respData.TargetTopicConfig.AccountType
		}

		if respData.TargetTopicConfig.Region != nil {
			targetTopicConfigMap["region"] = respData.TargetTopicConfig.Region
		}

		if respData.TargetTopicConfig.LogsetId != nil {
			targetTopicConfigMap["logset_id"] = respData.TargetTopicConfig.LogsetId
		}

		if respData.TargetTopicConfig.TopicId != nil {
			targetTopicConfigMap["topic_id"] = respData.TargetTopicConfig.TopicId
		}

		if respData.TargetTopicConfig.RoleArn != nil {
			targetTopicConfigMap["role_arn"] = respData.TargetTopicConfig.RoleArn
		}

		if respData.TargetTopicConfig.ExternalId != nil {
			targetTopicConfigMap["external_id"] = respData.TargetTopicConfig.ExternalId
		}

		_ = d.Set("target_topic_config", []interface{}{targetTopicConfigMap})
	}

	if respData.DeliverRule != nil {
		deliverRuleMap := map[string]interface{}{}
		if respData.DeliverRule.DataScope != nil {
			deliverRuleMap["data_scope"] = respData.DeliverRule.DataScope
		}

		_ = d.Set("deliver_rule", []interface{}{deliverRuleMap})
	}

	return nil
}

func resourceTencentCloudClsClsDeliverTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cls_deliver_task.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		taskId = d.Id()
	)

	needChange := false
	mutableArgs := []string{"task_name", "source_topic_config", "target_topic_config", "deliver_rule", "enable", "has_services_log"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := cls.NewModifyCLSDeliverTaskRequest()
		request.TaskId = helper.String(taskId)

		if v, ok := d.GetOk("task_name"); ok {
			request.TaskName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("source_topic_config"); ok {
			sourceTopicConfigList := v.([]interface{})
			if len(sourceTopicConfigList) > 0 {
				sourceTopicConfig := cls.SourceTopicConfig{}
				sourceTopicConfigMap := sourceTopicConfigList[0].(map[string]interface{})
				if v, ok := sourceTopicConfigMap["topic_filter_type"].(int); ok {
					sourceTopicConfig.TopicFilterType = helper.IntUint64(v)
				}

				if v, ok := sourceTopicConfigMap["logset_id"].(string); ok && v != "" {
					sourceTopicConfig.LogsetId = helper.String(v)
				}

				if v, ok := sourceTopicConfigMap["topics"]; ok {
					for _, item := range v.([]interface{}) {
						topicsMap := item.(map[string]interface{})
						sourceTopicInfo := cls.SourceTopicInfo{}
						if v, ok := topicsMap["topic_id"].(string); ok && v != "" {
							sourceTopicInfo.TopicId = helper.String(v)
						}
						sourceTopicConfig.Topics = append(sourceTopicConfig.Topics, &sourceTopicInfo)
					}
				}

				request.SourceTopicConfig = &sourceTopicConfig
			}
		}

		if v, ok := d.GetOk("target_topic_config"); ok {
			targetTopicConfigList := v.([]interface{})
			if len(targetTopicConfigList) > 0 {
				targetTopicConfig := cls.TargetTopicConfig{}
				targetTopicConfigMap := targetTopicConfigList[0].(map[string]interface{})
				if v, ok := targetTopicConfigMap["account_type"].(int); ok {
					targetTopicConfig.AccountType = helper.IntUint64(v)
				}

				if v, ok := targetTopicConfigMap["region"].(string); ok && v != "" {
					targetTopicConfig.Region = helper.String(v)
				}

				if v, ok := targetTopicConfigMap["logset_id"].(string); ok && v != "" {
					targetTopicConfig.LogsetId = helper.String(v)
				}

				if v, ok := targetTopicConfigMap["topic_id"].(string); ok && v != "" {
					targetTopicConfig.TopicId = helper.String(v)
				}

				if v, ok := targetTopicConfigMap["role_arn"].(string); ok && v != "" {
					targetTopicConfig.RoleArn = helper.String(v)
				}

				if v, ok := targetTopicConfigMap["external_id"].(string); ok && v != "" {
					targetTopicConfig.ExternalId = helper.String(v)
				}

				request.TargetTopicConfig = &targetTopicConfig
			}
		}

		if v, ok := d.GetOk("deliver_rule"); ok {
			deliverRuleList := v.([]interface{})
			if len(deliverRuleList) > 0 {
				deliverRule := cls.DeliverRule{}
				deliverRuleMap := deliverRuleList[0].(map[string]interface{})
				if v, ok := deliverRuleMap["data_scope"].(int); ok {
					deliverRule.DataScope = helper.IntUint64(v)
				}
				request.DeliverRule = &deliverRule
			}
		}

		if v, ok := d.GetOkExists("enable"); ok {
			request.Enable = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("has_services_log"); ok {
			request.HasServicesLog = helper.IntUint64(v.(int))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyCLSDeliverTaskWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify cls_cls_deliver_task failed, Response is nil."))
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cls_cls_deliver_task failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudClsClsDeliverTaskRead(d, meta)
}

func resourceTencentCloudClsClsDeliverTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cls_deliver_task.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		taskId = d.Id()
	)

	request := cls.NewDeleteCLSDeliverTaskRequest()
	request.TaskId = helper.String(taskId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().DeleteCLSDeliverTaskWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete cls_cls_deliver_task failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cls_cls_deliver_task failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
