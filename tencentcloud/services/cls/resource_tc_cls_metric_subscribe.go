package cls

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clsv20201016 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsMetricSubscribe() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsMetricSubscribeCreate,
		Read:   resourceTencentCloudClsMetricSubscribeRead,
		Update: resourceTencentCloudClsMetricSubscribeUpdate,
		Delete: resourceTencentCloudClsMetricSubscribeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subscribe task name, up to 64 characters, start with a letter, support 0-9, a-z, A-Z, _, -, Chinese characters.",
			},

			"topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Log topic id.",
			},

			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud product namespace.",
			},

			"metrics": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Metric config info list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Metric name.",
						},
						"periods": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "Statistical period, unit: second(s).",
						},
						"metric_labels": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Custom metric labels.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Metric label name.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Metric label content.",
									},
								},
							},
						},
					},
				},
			},

			"instance_info": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Instance config info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_dimension": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Instance dimension.",
						},
						"instances": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Instance value list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Instance info value list.",
									},
								},
							},
						},
					},
				},
			},

			"enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Task switch, 1: pause, 2: enable.",
			},

			// computed
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subscribe task id.",
			},

			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Subscribe task running status. 0: creating, 1: paused, 2: running, 3: abnormal.",
			},

			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Creation time (second-level timestamp).",
			},

			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Update time (second-level timestamp).",
			},
		},
	}
}

func resourceTencentCloudClsMetricSubscribeCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_metric_subscribe.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = clsv20201016.NewCreateMetricSubscribeRequest()
		response = clsv20201016.NewCreateMetricSubscribeResponse()
		topicId  string
		taskId   string
	)

	if v, ok := d.GetOk("topic_id"); ok {
		request.TopicId = helper.String(v.(string))
		topicId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		request.Namespace = helper.String(v.(string))
	}

	if v, ok := d.GetOk("metrics"); ok {
		for _, item := range v.([]interface{}) {
			metricsMap := item.(map[string]interface{})
			metricConfig := clsv20201016.MetricConfig{}
			if v, ok := metricsMap["metric_name"].(string); ok && v != "" {
				metricConfig.MetricName = helper.String(v)
			}

			if v, ok := metricsMap["periods"].([]interface{}); ok && len(v) > 0 {
				for _, period := range v {
					metricConfig.Periods = append(metricConfig.Periods, helper.IntUint64(period.(int)))
				}
			}

			if v, ok := metricsMap["metric_labels"].([]interface{}); ok && len(v) > 0 {
				for _, labelItem := range v {
					labelMap := labelItem.(map[string]interface{})
					metricLabel := clsv20201016.MetricLabel{}
					if v, ok := labelMap["key"].(string); ok && v != "" {
						metricLabel.Key = helper.String(v)
					}

					if v, ok := labelMap["value"].(string); ok && v != "" {
						metricLabel.Value = helper.String(v)
					}

					metricConfig.MetricLabels = append(metricConfig.MetricLabels, &metricLabel)
				}
			}

			request.Metrics = append(request.Metrics, &metricConfig)
		}
	}

	if v, ok := d.GetOk("instance_info"); ok {
		for _, item := range v.([]interface{}) {
			instanceInfoMap := item.(map[string]interface{})
			instanceConfig := clsv20201016.InstanceConfig{}
			if v, ok := instanceInfoMap["instance_dimension"].([]interface{}); ok && len(v) > 0 {
				for _, dimension := range v {
					instanceConfig.InstanceDimension = append(instanceConfig.InstanceDimension, helper.String(dimension.(string)))
				}
			}

			if v, ok := instanceInfoMap["instances"].([]interface{}); ok && len(v) > 0 {
				for _, instanceItem := range v {
					instanceMap := instanceItem.(map[string]interface{})
					instance := clsv20201016.Instance{}
					if v, ok := instanceMap["values"].([]interface{}); ok && len(v) > 0 {
						for _, value := range v {
							instance.Values = append(instance.Values, helper.String(value.(string)))
						}
					}

					instanceConfig.Instances = append(instanceConfig.Instances, &instance)
				}
			}

			request.InstanceInfo = &instanceConfig
			break
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().CreateMetricSubscribeWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create cls metric subscribe failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create cls metric subscribe failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[CRUD] cls_metric_subscribe logId=%s, topicId=%s, taskId=%+v", logId, topicId, response.Response.TaskId)

	if response.Response.TaskId == nil || *response.Response.TaskId == "" {
		return fmt.Errorf("Create cls metric subscribe failed, TaskId is empty.")
	}

	taskId = *response.Response.TaskId
	d.SetId(strings.Join([]string{topicId, taskId}, tccommon.FILED_SP))
	return resourceTencentCloudClsMetricSubscribeRead(d, meta)
}

func resourceTencentCloudClsMetricSubscribeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_metric_subscribe.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	topicId := idSplit[0]
	taskId := idSplit[1]

	respData, err := service.DescribeClsMetricSubscribeById(ctx, topicId, taskId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] cls_metric_subscribe id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.TopicId != nil {
		_ = d.Set("topic_id", respData.TopicId)
	}

	if respData.Namespace != nil {
		_ = d.Set("namespace", respData.Namespace)
	}

	if respData.Metrics != nil && len(respData.Metrics) > 0 {
		metricsList := make([]map[string]interface{}, 0, len(respData.Metrics))
		for _, metricConfig := range respData.Metrics {
			metricMap := map[string]interface{}{}
			if metricConfig.MetricName != nil {
				metricMap["metric_name"] = *metricConfig.MetricName
			}

			if metricConfig.Periods != nil && len(metricConfig.Periods) > 0 {
				periodsList := make([]interface{}, 0, len(metricConfig.Periods))
				for _, period := range metricConfig.Periods {
					if period != nil {
						periodsList = append(periodsList, int(*period))
					}
				}
				if len(periodsList) > 0 {
					metricMap["periods"] = periodsList
				}
			}

			if metricConfig.MetricLabels != nil && len(metricConfig.MetricLabels) > 0 {
				metricLabelsList := make([]map[string]interface{}, 0, len(metricConfig.MetricLabels))
				for _, metricLabel := range metricConfig.MetricLabels {
					metricLabelMap := map[string]interface{}{}
					if metricLabel.Key != nil {
						metricLabelMap["key"] = *metricLabel.Key
					}

					if metricLabel.Value != nil {
						metricLabelMap["value"] = *metricLabel.Value
					}

					metricLabelsList = append(metricLabelsList, metricLabelMap)
				}
				metricMap["metric_labels"] = metricLabelsList
			}

			metricsList = append(metricsList, metricMap)
		}

		_ = d.Set("metrics", metricsList)
	}

	if respData.InstanceInfo != nil {
		instanceInfoList := make([]map[string]interface{}, 0, 1)
		instanceInfoMap := map[string]interface{}{}
		if respData.InstanceInfo.InstanceDimension != nil && len(respData.InstanceInfo.InstanceDimension) > 0 {
			instanceDimensionList := make([]interface{}, 0, len(respData.InstanceInfo.InstanceDimension))
			for _, dimension := range respData.InstanceInfo.InstanceDimension {
				if dimension != nil {
					instanceDimensionList = append(instanceDimensionList, *dimension)
				}
			}
			if len(instanceDimensionList) > 0 {
				instanceInfoMap["instance_dimension"] = instanceDimensionList
			}
		}

		if respData.InstanceInfo.Instances != nil && len(respData.InstanceInfo.Instances) > 0 {
			instancesList := make([]map[string]interface{}, 0, len(respData.InstanceInfo.Instances))
			for _, instance := range respData.InstanceInfo.Instances {
				instanceMap := map[string]interface{}{}
				if instance.Values != nil && len(instance.Values) > 0 {
					valuesList := make([]interface{}, 0, len(instance.Values))
					for _, value := range instance.Values {
						if value != nil {
							valuesList = append(valuesList, *value)
						}
					}
					if len(valuesList) > 0 {
						instanceMap["values"] = valuesList
					}
				}

				instancesList = append(instancesList, instanceMap)
			}
			instanceInfoMap["instances"] = instancesList
		}

		instanceInfoList = append(instanceInfoList, instanceInfoMap)
		_ = d.Set("instance_info", instanceInfoList)
	}

	if respData.Enable != nil {
		_ = d.Set("enable", int(*respData.Enable))
	}

	if respData.TaskId != nil {
		_ = d.Set("task_id", respData.TaskId)
	}

	if respData.Status != nil {
		_ = d.Set("status", int(*respData.Status))
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", int(*respData.CreateTime))
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", int(*respData.UpdateTime))
	}

	return nil
}

func resourceTencentCloudClsMetricSubscribeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_metric_subscribe.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	topicId := idSplit[0]
	taskId := idSplit[1]

	needChange := false
	mutableArgs := []string{"name", "namespace", "metrics", "instance_info", "enable"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := clsv20201016.NewModifyMetricSubscribeRequest()
		request.TopicId = helper.String(topicId)
		request.TaskId = helper.String(taskId)

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("namespace"); ok {
			request.Namespace = helper.String(v.(string))
		}

		if v, ok := d.GetOk("metrics"); ok {
			for _, item := range v.([]interface{}) {
				metricsMap := item.(map[string]interface{})
				metricConfig := clsv20201016.MetricConfig{}
				if v, ok := metricsMap["metric_name"].(string); ok && v != "" {
					metricConfig.MetricName = helper.String(v)
				}

				if v, ok := metricsMap["periods"].([]interface{}); ok && len(v) > 0 {
					for _, period := range v {
						metricConfig.Periods = append(metricConfig.Periods, helper.IntUint64(period.(int)))
					}
				}

				if v, ok := metricsMap["metric_labels"].([]interface{}); ok && len(v) > 0 {
					for _, labelItem := range v {
						labelMap := labelItem.(map[string]interface{})
						metricLabel := clsv20201016.MetricLabel{}
						if v, ok := labelMap["key"].(string); ok && v != "" {
							metricLabel.Key = helper.String(v)
						}

						if v, ok := labelMap["value"].(string); ok && v != "" {
							metricLabel.Value = helper.String(v)
						}

						metricConfig.MetricLabels = append(metricConfig.MetricLabels, &metricLabel)
					}
				}

				request.Metrics = append(request.Metrics, &metricConfig)
			}
		}

		if v, ok := d.GetOk("instance_info"); ok {
			for _, item := range v.([]interface{}) {
				instanceInfoMap := item.(map[string]interface{})
				instanceConfig := clsv20201016.InstanceConfig{}
				if v, ok := instanceInfoMap["instance_dimension"].([]interface{}); ok && len(v) > 0 {
					for _, dimension := range v {
						instanceConfig.InstanceDimension = append(instanceConfig.InstanceDimension, helper.String(dimension.(string)))
					}
				}

				if v, ok := instanceInfoMap["instances"].([]interface{}); ok && len(v) > 0 {
					for _, instanceItem := range v {
						instanceMap := instanceItem.(map[string]interface{})
						instance := clsv20201016.Instance{}
						if v, ok := instanceMap["values"].([]interface{}); ok && len(v) > 0 {
							for _, value := range v {
								instance.Values = append(instance.Values, helper.String(value.(string)))
							}
						}

						instanceConfig.Instances = append(instanceConfig.Instances, &instance)
					}
				}

				request.InstanceInfo = &instanceConfig
				break
			}
		}

		if v, ok := d.GetOk("enable"); ok {
			request.Enable = helper.IntUint64(v.(int))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().ModifyMetricSubscribeWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify cls metric subscribe failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update cls metric subscribe failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudClsMetricSubscribeRead(d, meta)
}

func resourceTencentCloudClsMetricSubscribeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_metric_subscribe.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = clsv20201016.NewDeleteMetricSubscribeRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	topicId := idSplit[0]
	taskId := idSplit[1]

	request.TopicId = helper.String(topicId)
	request.TaskId = helper.String(taskId)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().DeleteMetricSubscribeWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete cls metric subscribe failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete cls metric subscribe failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
