package cdb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCdbStartCpuExpandAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbStartCpuExpandAttachmentCreate,
		Read:   resourceTencentCloudCdbStartCpuExpandAttachmentRead,
		Update: resourceTencentCloudCdbStartCpuExpandAttachmentUpdate,
		Delete: resourceTencentCloudCdbStartCpuExpandAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID, which can be obtained from the DescribeDBInstances API.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Expansion type. Valid values: `auto`, `manual`, `timeInterval`, `period`.",
			},
			"expand_cpu": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "CPU cores to expand. Required when `type` is `manual`, `timeInterval`, or `period`.",
			},
			"auto_strategy": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Auto expansion strategy. Required when `type` is `auto`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expand_threshold": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Auto expansion threshold. Valid values: 40, 50, 60, 70, 80, 90.",
						},
						"shrink_threshold": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Auto shrink threshold. Valid values: 10, 20, 30.",
						},
						"expand_second_period": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Expansion observation period in seconds. Valid values: 15, 30, 45, 60, 180, 300, 600, 900, 1800.",
						},
						"shrink_second_period": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Shrink observation period in seconds. Valid values: 300, 600, 900, 1800.",
						},
					},
				},
			},
			"time_interval_strategy": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Time interval expansion strategy. Required when `type` is `timeInterval`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Start expansion time as integer timestamp in seconds.",
						},
						"end_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "End expansion time as integer timestamp in seconds.",
						},
					},
				},
			},
			"period_strategy": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Period expansion strategy. Required when `type` is `period`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time_cycle": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Weekly cycle configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"monday": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to expand on Monday.",
									},
									"tuesday": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to expand on Tuesday.",
									},
									"wednesday": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to expand on Wednesday.",
									},
									"thursday": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to expand on Thursday.",
									},
									"friday": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to expand on Friday.",
									},
									"saturday": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to expand on Saturday.",
									},
									"sunday": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to expand on Sunday.",
									},
								},
							},
						},
						"time_interval": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Daily time range configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Start time string.",
									},
									"end_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "End time string.",
									},
								},
							},
						},
					},
				},
			},
			"async_request_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Async request ID returned by Create/Delete APIs.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceTencentCloudCdbStartCpuExpandAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdb_start_cpu_expand.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = cdb.NewStartCpuExpandRequest()
		instanceId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("expand_cpu"); ok {
		request.ExpandCpu = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("auto_strategy"); ok {
		for _, item := range v.([]interface{}) {
			autoStrategyMap := item.(map[string]interface{})
			autoStrategy := cdb.AutoStrategy{}
			if v, ok := autoStrategyMap["expand_threshold"].(int); ok && v != 0 {
				autoStrategy.ExpandThreshold = helper.IntInt64(v)
			}
			if v, ok := autoStrategyMap["shrink_threshold"].(int); ok && v != 0 {
				autoStrategy.ShrinkThreshold = helper.IntInt64(v)
			}
			if v, ok := autoStrategyMap["expand_second_period"].(int); ok && v != 0 {
				autoStrategy.ExpandSecondPeriod = helper.IntInt64(v)
			}
			if v, ok := autoStrategyMap["shrink_second_period"].(int); ok && v != 0 {
				autoStrategy.ShrinkSecondPeriod = helper.IntInt64(v)
			}
			request.AutoStrategy = &autoStrategy
		}
	}

	if v, ok := d.GetOk("time_interval_strategy"); ok {
		for _, item := range v.([]interface{}) {
			timeIntervalStrategyMap := item.(map[string]interface{})
			timeIntervalStrategy := cdb.TimeIntervalStrategy{}
			if v, ok := timeIntervalStrategyMap["start_time"].(int); ok && v != 0 {
				timeIntervalStrategy.StartTime = helper.IntInt64(v)
			}
			if v, ok := timeIntervalStrategyMap["end_time"].(int); ok && v != 0 {
				timeIntervalStrategy.EndTime = helper.IntInt64(v)
			}
			request.TimeIntervalStrategy = &timeIntervalStrategy
		}
	}

	if v, ok := d.GetOk("period_strategy"); ok {
		for _, item := range v.([]interface{}) {
			periodStrategyMap := item.(map[string]interface{})
			periodStrategy := cdb.PeriodStrategy{}
			if v, ok := periodStrategyMap["time_cycle"]; ok {
				for _, cycleItem := range v.([]interface{}) {
					timeCycleMap := cycleItem.(map[string]interface{})
					timeCycle := cdb.TImeCycle{}
					if v, ok := timeCycleMap["monday"].(bool); ok {
						timeCycle.Monday = helper.Bool(v)
					}
					if v, ok := timeCycleMap["tuesday"].(bool); ok {
						timeCycle.Tuesday = helper.Bool(v)
					}
					if v, ok := timeCycleMap["wednesday"].(bool); ok {
						timeCycle.Wednesday = helper.Bool(v)
					}
					if v, ok := timeCycleMap["thursday"].(bool); ok {
						timeCycle.Thursday = helper.Bool(v)
					}
					if v, ok := timeCycleMap["friday"].(bool); ok {
						timeCycle.Friday = helper.Bool(v)
					}
					if v, ok := timeCycleMap["saturday"].(bool); ok {
						timeCycle.Saturday = helper.Bool(v)
					}
					if v, ok := timeCycleMap["sunday"].(bool); ok {
						timeCycle.Sunday = helper.Bool(v)
					}
					periodStrategy.TimeCycle = &timeCycle
				}
			}
			if v, ok := periodStrategyMap["time_interval"]; ok {
				for _, intervalItem := range v.([]interface{}) {
					timeIntervalMap := intervalItem.(map[string]interface{})
					timeInterval := cdb.TimeInterval{}
					if v, ok := timeIntervalMap["start_time"].(string); ok && v != "" {
						timeInterval.StartTime = helper.String(v)
					}
					if v, ok := timeIntervalMap["end_time"].(string); ok && v != "" {
						timeInterval.EndTime = helper.String(v)
					}
					periodStrategy.TimeInterval = &timeInterval
				}
			}
			request.PeriodStrategy = &periodStrategy
		}
	}

	var asyncRequestId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().StartCpuExpandWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create cdb_start_cpu_expand failed, Response is nil."))
		}
		if result.Response.AsyncRequestId == nil || *result.Response.AsyncRequestId == "" {
			log.Printf("[CRITAL]%s create cdb_start_cpu_expand, logId=%s, id=%s\n", logId, logId, d.Id())
			return resource.NonRetryableError(fmt.Errorf("Create cdb_start_cpu_expand failed, AsyncRequestId is nil or empty."))
		}
		asyncRequestId = *result.Response.AsyncRequestId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create cdb_start_cpu_expand failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	_ = d.Set("async_request_id", asyncRequestId)
	d.SetId(instanceId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		result, e := service.DescribeCdbStartCpuExpandAttachmentById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil || result.Type == nil {
			return resource.RetryableError(fmt.Errorf("cdb_start_cpu_expand strategy not active yet"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cdb_start_cpu_expand poll failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCdbStartCpuExpandAttachmentRead(d, meta)
}

func resourceTencentCloudCdbStartCpuExpandAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdb_start_cpu_expand.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	respData, err := service.DescribeCdbStartCpuExpandAttachmentById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil || respData.Type == nil {
		log.Printf("[CRUD] cdb_start_cpu_expand id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if respData.Type != nil {
		_ = d.Set("type", respData.Type)
	}

	if respData.ExpandCpu != nil {
		_ = d.Set("expand_cpu", respData.ExpandCpu)
	}

	if respData.AutoStrategy != nil {
		autoStrategyList := []map[string]interface{}{}
		autoStrategyMap := map[string]interface{}{}
		if respData.AutoStrategy.ExpandThreshold != nil {
			autoStrategyMap["expand_threshold"] = respData.AutoStrategy.ExpandThreshold
		}
		if respData.AutoStrategy.ShrinkThreshold != nil {
			autoStrategyMap["shrink_threshold"] = respData.AutoStrategy.ShrinkThreshold
		}
		if respData.AutoStrategy.ExpandSecondPeriod != nil {
			autoStrategyMap["expand_second_period"] = respData.AutoStrategy.ExpandSecondPeriod
		}
		if respData.AutoStrategy.ShrinkSecondPeriod != nil {
			autoStrategyMap["shrink_second_period"] = respData.AutoStrategy.ShrinkSecondPeriod
		}
		autoStrategyList = append(autoStrategyList, autoStrategyMap)
		_ = d.Set("auto_strategy", autoStrategyList)
	}

	if respData.TimeIntervalStrategy != nil {
		timeIntervalStrategyList := []map[string]interface{}{}
		timeIntervalStrategyMap := map[string]interface{}{}
		if respData.TimeIntervalStrategy.StartTime != nil {
			timeIntervalStrategyMap["start_time"] = respData.TimeIntervalStrategy.StartTime
		}
		if respData.TimeIntervalStrategy.EndTime != nil {
			timeIntervalStrategyMap["end_time"] = respData.TimeIntervalStrategy.EndTime
		}
		timeIntervalStrategyList = append(timeIntervalStrategyList, timeIntervalStrategyMap)
		_ = d.Set("time_interval_strategy", timeIntervalStrategyList)
	}

	if respData.PeriodStrategy != nil {
		periodStrategyList := []map[string]interface{}{}
		periodStrategyMap := map[string]interface{}{}
		if respData.PeriodStrategy.TimeCycle != nil {
			timeCycleList := []map[string]interface{}{}
			timeCycleMap := map[string]interface{}{}
			if respData.PeriodStrategy.TimeCycle.Monday != nil {
				timeCycleMap["monday"] = respData.PeriodStrategy.TimeCycle.Monday
			}
			if respData.PeriodStrategy.TimeCycle.Tuesday != nil {
				timeCycleMap["tuesday"] = respData.PeriodStrategy.TimeCycle.Tuesday
			}
			if respData.PeriodStrategy.TimeCycle.Wednesday != nil {
				timeCycleMap["wednesday"] = respData.PeriodStrategy.TimeCycle.Wednesday
			}
			if respData.PeriodStrategy.TimeCycle.Thursday != nil {
				timeCycleMap["thursday"] = respData.PeriodStrategy.TimeCycle.Thursday
			}
			if respData.PeriodStrategy.TimeCycle.Friday != nil {
				timeCycleMap["friday"] = respData.PeriodStrategy.TimeCycle.Friday
			}
			if respData.PeriodStrategy.TimeCycle.Saturday != nil {
				timeCycleMap["saturday"] = respData.PeriodStrategy.TimeCycle.Saturday
			}
			if respData.PeriodStrategy.TimeCycle.Sunday != nil {
				timeCycleMap["sunday"] = respData.PeriodStrategy.TimeCycle.Sunday
			}
			timeCycleList = append(timeCycleList, timeCycleMap)
			periodStrategyMap["time_cycle"] = timeCycleList
		}
		if respData.PeriodStrategy.TimeInterval != nil {
			timeIntervalList := []map[string]interface{}{}
			timeIntervalMap := map[string]interface{}{}
			if respData.PeriodStrategy.TimeInterval.StartTime != nil {
				timeIntervalMap["start_time"] = respData.PeriodStrategy.TimeInterval.StartTime
			}
			if respData.PeriodStrategy.TimeInterval.EndTime != nil {
				timeIntervalMap["end_time"] = respData.PeriodStrategy.TimeInterval.EndTime
			}
			timeIntervalList = append(timeIntervalList, timeIntervalMap)
			periodStrategyMap["time_interval"] = timeIntervalList
		}
		periodStrategyList = append(periodStrategyList, periodStrategyMap)
		_ = d.Set("period_strategy", periodStrategyList)
	}

	return nil
}

func resourceTencentCloudCdbStartCpuExpandAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdb_start_cpu_expand.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"type", "expand_cpu", "auto_strategy", "time_interval_strategy", "period_strategy"}
	for _, arg := range immutableArgs {
		if d.HasChange(arg) {
			return fmt.Errorf("cdb_start_cpu_expand argument `%s` has changed, this is an immutable argument which can only be changed by recreating the resource", arg)
		}
	}

	return resourceTencentCloudCdbStartCpuExpandAttachmentRead(d, meta)
}

func resourceTencentCloudCdbStartCpuExpandAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdb_start_cpu_expand.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	err := service.DeleteCdbStartCpuExpandAttachmentById(ctx, instanceId)
	if err != nil {
		return err
	}

	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		result, e := service.DescribeCdbStartCpuExpandAttachmentById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result != nil && result.Type != nil {
			return resource.RetryableError(fmt.Errorf("cdb_start_cpu_expand strategy still active"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete cdb_start_cpu_expand poll failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
