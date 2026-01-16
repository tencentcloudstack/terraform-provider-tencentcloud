package apm

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apm/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudApmInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudApmInstancesRead,

		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter by instance ID list (exact match).",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by instance ID (fuzzy match).",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by instance name (fuzzy match).",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Filter by tags.",
			},
			"demo_instance_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to query official demo instances. 0: non-demo, 1: demo. Default is 0.",
			},
			"all_regions_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to query instances in all regions. 0: no, 1: yes. Default is 0.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "APM instance list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance description.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "App ID.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance status.",
						},
						"create_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creator UIN.",
						},
						"trace_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trace data retention duration.",
						},
						"span_daily_counters": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Daily span count quota.",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Billing mode.",
						},
						"free": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether it is free edition.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tag list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value.",
									},
								},
							},
						},
						"err_rate_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Error rate threshold.",
						},
						"sample_rate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sampling rate.",
						},
						"error_sample": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Error sampling switch.",
						},
						"service_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Service count.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudApmInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_apm_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := ApmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	params := make(map[string]interface{})

	if v, ok := d.GetOk("instance_ids"); ok {
		params["instance_ids"] = v.([]interface{})
	}
	if v, ok := d.GetOk("instance_id"); ok {
		params["instance_id"] = v.(string)
	}
	if v, ok := d.GetOk("instance_name"); ok {
		params["instance_name"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		params["tags"] = v.(map[string]interface{})
	}
	if v, ok := d.GetOkExists("demo_instance_flag"); ok {
		params["demo_instance_flag"] = v.(int)
	}
	if v, ok := d.GetOkExists("all_regions_flag"); ok {
		params["all_regions_flag"] = v.(int)
	}

	var instances []*apm.ApmInstanceDetail
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := service.DescribeApmInstances(ctx, params)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instances = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read APM instances failed, reason:%+v", logId, err)
		return err
	}

	instanceList := make([]map[string]interface{}, 0, len(instances))
	ids := make([]string, 0, len(instances))

	for _, instance := range instances {
		instanceMap := map[string]interface{}{}

		if instance.InstanceId != nil {
			instanceMap["instance_id"] = instance.InstanceId
			ids = append(ids, *instance.InstanceId)
		}
		if instance.Name != nil {
			instanceMap["name"] = instance.Name
		}
		if instance.Description != nil {
			instanceMap["description"] = instance.Description
		}
		if instance.Region != nil {
			instanceMap["region"] = instance.Region
		}
		if instance.AppId != nil {
			instanceMap["app_id"] = instance.AppId
		}
		if instance.Status != nil {
			instanceMap["status"] = instance.Status
		}
		if instance.CreateUin != nil {
			instanceMap["create_uin"] = instance.CreateUin
		}
		if instance.TraceDuration != nil {
			instanceMap["trace_duration"] = instance.TraceDuration
		}
		if instance.SpanDailyCounters != nil {
			instanceMap["span_daily_counters"] = instance.SpanDailyCounters
		}
		if instance.PayMode != nil {
			instanceMap["pay_mode"] = instance.PayMode
		}
		if instance.Free != nil {
			instanceMap["free"] = instance.Free
		}
		if instance.Tags != nil {
			tags := make([]map[string]interface{}, 0, len(instance.Tags))
			for _, tag := range instance.Tags {
				tagMap := map[string]interface{}{}
				if tag.Key != nil {
					tagMap["key"] = tag.Key
				}
				if tag.Value != nil {
					tagMap["value"] = tag.Value
				}
				tags = append(tags, tagMap)
			}
			instanceMap["tags"] = tags
		}
		if instance.ErrRateThreshold != nil {
			instanceMap["err_rate_threshold"] = instance.ErrRateThreshold
		}
		if instance.SampleRate != nil {
			instanceMap["sample_rate"] = instance.SampleRate
		}
		if instance.ErrorSample != nil {
			instanceMap["error_sample"] = instance.ErrorSample
		}
		if instance.ServiceCount != nil {
			instanceMap["service_count"] = instance.ServiceCount
		}

		instanceList = append(instanceList, instanceMap)
	}

	_ = d.Set("instance_list", instanceList)
	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), instanceList); e != nil {
			return e
		}
	}

	return nil
}
