package cls

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsTopicCreate,
		Read:   resourceTencentCloudClsTopicRead,
		Update: resourceTencentCloudClsTopicUpdate,
		Delete: resourceTencentCloudClsTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"logset_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Logset ID.",
			},
			"topic_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Log topic name.",
			},
			"partition_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of log topic partitions. Default value: 1. Maximum value: 10.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list. Up to 10 tag key-value pairs are supported and must be unique.",
			},
			"auto_split": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable automatic split. Default value: true.",
			},
			"max_split_partitions": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "Maximum number of partitions to split into for this topic if" +
					" automatic split is enabled. Default value: 50.",
			},
			"storage_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Log topic storage class. Valid values: hot: real-time storage; cold: offline storage. Default value: hot. If cold is passed in, " +
					"please contact the customer service to add the log topic to the allowlist first.",
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Lifecycle in days. Value range: 1~366. Default value: 30.",
			},
			"hot_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "0: Turn off log sinking. Non 0: The number of days of standard storage after enabling log settling. HotPeriod needs to be greater than or equal to 7 and less than Period. Only effective when StorageType is hot.",
			},
			"describes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Log Topic Description.",
			},
			"is_web_tracking": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "No authentication switch. False: closed; True: Enable. The default is false. After activation, anonymous access to the log topic will be supported for specified operations.",
			},
			"extends": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Log Subject Extension Information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"anonymous_access": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Log topic authentication free configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operations": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Operation list, supporting trackLog (JS/HTTP upload log) and realtimeProducer (kafka protocol upload log).",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"conditions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Operation list, supporting trackLog (JS/HTTP upload log) and realtimeProducer (kafka protocol upload log).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"attributes": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Condition attribute, currently only VpcID is supported.",
												},
												"rule": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Conditional rule, 1: equal, 2: not equal.",
												},
												"condition_value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The value of the corresponding conditional attribute.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClsTopicCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_topic.create")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		request       = cls.NewCreateTopicRequest()
		response      *cls.CreateTopicResponse
		isWebTracking bool
	)

	if v, ok := d.GetOk("logset_id"); ok {
		request.LogsetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_name"); ok {
		request.TopicName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("partition_count"); ok {
		request.PartitionCount = helper.IntInt64(v.(int))
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		for k, v := range tags {
			key := k
			value := v
			request.Tags = append(request.Tags, &cls.Tag{
				Key:   &key,
				Value: &value,
			})
		}
	}

	if v, ok := d.GetOkExists("auto_split"); ok {
		request.AutoSplit = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("max_split_partitions"); ok {
		request.MaxSplitPartitions = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("storage_type"); ok {
		request.StorageType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("hot_period"); ok {
		request.HotPeriod = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("describes"); ok {
		request.Describes = helper.String(v.(string))
	} else {
		request.Describes = helper.String("")
	}

	if v, ok := d.GetOkExists("is_web_tracking"); ok {
		request.IsWebTracking = helper.Bool(v.(bool))
		isWebTracking = v.(bool)
	}

	if isWebTracking {
		if dMap, ok := helper.InterfacesHeadMap(d, "extends"); ok {
			topicExtendInfo := cls.TopicExtendInfo{}
			if anonymousAccessMap, ok := helper.InterfaceToMap(dMap, "anonymous_access"); ok {
				anonymousInfo := cls.AnonymousInfo{}
				if v, ok := anonymousAccessMap["operations"]; ok {
					tmpList := make([]*string, 0)
					for _, operation := range v.([]interface{}) {
						tmpList = append(tmpList, helper.String(operation.(string)))
					}

					anonymousInfo.Operations = tmpList
				}

				if v, ok := anonymousAccessMap["conditions"]; ok {
					for _, condition := range v.([]interface{}) {
						conditionMap := condition.(map[string]interface{})
						conditionInfo := cls.ConditionInfo{}
						if v, ok := conditionMap["attributes"]; ok {
							conditionInfo.Attributes = helper.String(v.(string))
						}

						if v, ok := conditionMap["rule"]; ok {
							conditionInfo.Rule = helper.IntUint64(v.(int))
						}

						if v, ok := conditionMap["condition_value"]; ok {
							conditionInfo.ConditionValue = helper.String(v.(string))
						}

						anonymousInfo.Conditions = append(anonymousInfo.Conditions, &conditionInfo)
					}
				}

				topicExtendInfo.AnonymousAccess = &anonymousInfo
			}

			request.Extends = &topicExtendInfo
		}
	} else {
		if _, ok := helper.InterfacesHeadMap(d, "extends"); ok {
			return fmt.Errorf("If `is_web_tracking` is false, Not support set `extends`.\n.")
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateTopic(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("create cls topic failed")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls topic failed, reason:%+v", logId, err)
		return err
	}

	id := *response.Response.TopicId
	d.SetId(id)
	return resourceTencentCloudClsTopicRead(d, meta)
}

func resourceTencentCloudClsTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_topic.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id      = d.Id()
	)

	topic, err := service.DescribeClsTopicById(ctx, id)
	if err != nil {
		return err
	}

	if topic == nil {
		d.SetId("")
		return fmt.Errorf("resource `Topic` %s does not exist", id)
	}

	_ = d.Set("logset_id", topic.LogsetId)
	_ = d.Set("topic_name", topic.TopicName)
	_ = d.Set("partition_count", topic.PartitionCount)

	tags := make(map[string]string, len(topic.Tags))
	for _, tag := range topic.Tags {
		tags[*tag.Key] = *tag.Value
	}

	_ = d.Set("tags", tags)
	_ = d.Set("auto_split", topic.AutoSplit)
	_ = d.Set("max_split_partitions", topic.MaxSplitPartitions)
	_ = d.Set("storage_type", topic.StorageType)
	_ = d.Set("period", topic.Period)
	_ = d.Set("hot_period", topic.HotPeriod)
	_ = d.Set("describes", topic.Describes)
	_ = d.Set("is_web_tracking", topic.IsWebTracking)

	if *topic.IsWebTracking {
		if topic.Extends != nil {
			extendMap := map[string]interface{}{}
			if topic.Extends.AnonymousAccess != nil {
				anonymousAccessMap := map[string]interface{}{}
				if topic.Extends.AnonymousAccess.Operations != nil {
					operationList := make([]string, 0, len(topic.Extends.AnonymousAccess.Operations))
					for _, v := range topic.Extends.AnonymousAccess.Operations {
						operationList = append(operationList, *v)
					}

					anonymousAccessMap["operations"] = operationList
				}

				if topic.Extends.AnonymousAccess.Conditions != nil {
					conditionList := []interface{}{}
					for _, v := range topic.Extends.AnonymousAccess.Conditions {
						conditionMap := map[string]interface{}{}
						if v.Attributes != nil {
							conditionMap["attributes"] = *v.Attributes
						}

						if v.Rule != nil {
							conditionMap["rule"] = *v.Rule
						}

						if v.ConditionValue != nil {
							conditionMap["condition_value"] = *v.ConditionValue
						}

						conditionList = append(conditionList, conditionMap)
					}

					anonymousAccessMap["conditions"] = conditionList
				}

				extendMap["anonymous_access"] = []interface{}{anonymousAccessMap}
			}

			_ = d.Set("extends", []interface{}{extendMap})
		}
	}

	return nil
}

func resourceTencentCloudClsTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_topic.update")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		request       = cls.NewModifyTopicRequest()
		id            = d.Id()
		isWebTracking bool
	)

	immutableArgs := []string{"partition_count", "storage_type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.TopicId = helper.String(id)

	if d.HasChange("topic_name") {
		request.TopicName = helper.String(d.Get("topic_name").(string))
	}

	if d.HasChange("tags") {
		tags := d.Get("tags").(map[string]interface{})
		request.Tags = make([]*cls.Tag, 0, len(tags))
		for k, v := range tags {
			key := k
			value := v
			request.Tags = append(request.Tags, &cls.Tag{
				Key:   &key,
				Value: helper.String(value.(string)),
			})
		}
	}

	if d.HasChange("auto_split") {
		request.AutoSplit = helper.Bool(d.Get("auto_split").(bool))
	}

	if d.HasChange("max_split_partitions") {
		request.MaxSplitPartitions = helper.IntInt64(d.Get("max_split_partitions").(int))
	}

	if d.HasChange("period") {
		request.Period = helper.IntInt64(d.Get("period").(int))
	}

	if d.HasChange("hot_period") {
		request.HotPeriod = helper.IntUint64(d.Get("hot_period").(int))
	}

	if d.HasChange("describes") {
		request.Describes = helper.String(d.Get("describes").(string))
	}

	if v, ok := d.GetOkExists("is_web_tracking"); ok {
		request.IsWebTracking = helper.Bool(v.(bool))
		isWebTracking = v.(bool)
	}

	if isWebTracking {
		if dMap, ok := helper.InterfacesHeadMap(d, "extends"); ok {
			if anonymousAccessMap, ok := helper.InterfaceToMap(dMap, "anonymous_access"); ok {
				topicExtendInfo := cls.TopicExtendInfo{}
				anonymousInfo := cls.AnonymousInfo{}
				if v, ok := anonymousAccessMap["operations"]; ok {
					tmpList := make([]*string, 0)
					for _, operation := range v.([]interface{}) {
						tmpList = append(tmpList, helper.String(operation.(string)))
					}

					anonymousInfo.Operations = tmpList
				}

				if v, ok := anonymousAccessMap["conditions"]; ok {
					for _, condition := range v.([]interface{}) {
						conditionMap := condition.(map[string]interface{})
						conditionInfo := cls.ConditionInfo{}
						if v, ok := conditionMap["attributes"]; ok {
							conditionInfo.Attributes = helper.String(v.(string))
						}

						if v, ok := conditionMap["rule"]; ok {
							conditionInfo.Rule = helper.IntUint64(v.(int))
						}

						if v, ok := conditionMap["condition_value"]; ok {
							conditionInfo.ConditionValue = helper.String(v.(string))
						}

						anonymousInfo.Conditions = append(anonymousInfo.Conditions, &conditionInfo)
					}
				}

				topicExtendInfo.AnonymousAccess = &anonymousInfo
				request.Extends = &topicExtendInfo
			}
		} else {
			return fmt.Errorf("If `is_web_tracking` is true, Must set `extends` params.\n.")
		}
	} else {
		if _, ok := helper.InterfacesHeadMap(d, "extends"); ok {
			return fmt.Errorf("If `is_web_tracking` is false, Not support set `extends` params.\n.")
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyTopic(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudClsTopicRead(d, meta)
}

func resourceTencentCloudClsTopicDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_topic.delete")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id      = d.Id()
	)

	if err := service.DeleteClsTopic(ctx, id); err != nil {
		return err
	}

	return nil
}
