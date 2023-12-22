package trabbit

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"

	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTdmqRabbitmqVipInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqRabbitmqVipInstanceRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "query condition filter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the filter parameter.",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "value.",
						},
					},
				},
			},
			// computed
			"instances": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Instance information list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance id.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance name.",
						},
						"instance_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance versionNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance status, 0 means creating, 1 means normal, 2 means isolating, 3 means destroyed, 4 - abnormal, 5 - delivery failed.",
						},
						"node_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of nodes.",
						},
						"config_display": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance configuration specification name.",
						},
						"max_tps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Peak TPS.",
						},
						"max_band_width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Peak bandwidth, in Mbps.",
						},
						"max_storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Storage capacity, in GB.",
						},
						"expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance expiration time, in milliseconds.",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Automatic renewal mark, 0 indicates the default state (the user has not set it, that is, the initial state is manual renewal), 1 indicates automatic renewal, 2 indicates that the automatic renewal is not specified (user setting).",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "0-postpaid, 1-prepaid.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "RemarksNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"spec_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Configuration ID.",
						},
						"exception_information": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cluster is abnormal.Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTdmqRabbitmqVipInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tdmq_rabbitmq_vip_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		instances []*tdmq.RabbitMQVipInstance
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tdmq.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := tdmq.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqRabbitmqVipInstanceByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		instances = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instances))
	tmpList := make([]map[string]interface{}, 0, len(instances))

	if instances != nil {
		for _, rabbitMQVipInstance := range instances {
			rabbitMQVipInstanceMap := map[string]interface{}{}

			if rabbitMQVipInstance.InstanceId != nil {
				rabbitMQVipInstanceMap["instance_id"] = rabbitMQVipInstance.InstanceId
			}

			if rabbitMQVipInstance.InstanceName != nil {
				rabbitMQVipInstanceMap["instance_name"] = rabbitMQVipInstance.InstanceName
			}

			if rabbitMQVipInstance.InstanceVersion != nil {
				rabbitMQVipInstanceMap["instance_version"] = rabbitMQVipInstance.InstanceVersion
			}

			if rabbitMQVipInstance.Status != nil {
				rabbitMQVipInstanceMap["status"] = rabbitMQVipInstance.Status
			}

			if rabbitMQVipInstance.NodeCount != nil {
				rabbitMQVipInstanceMap["node_count"] = rabbitMQVipInstance.NodeCount
			}

			if rabbitMQVipInstance.ConfigDisplay != nil {
				rabbitMQVipInstanceMap["config_display"] = rabbitMQVipInstance.ConfigDisplay
			}

			if rabbitMQVipInstance.MaxTps != nil {
				rabbitMQVipInstanceMap["max_tps"] = rabbitMQVipInstance.MaxTps
			}

			if rabbitMQVipInstance.MaxBandWidth != nil {
				rabbitMQVipInstanceMap["max_band_width"] = rabbitMQVipInstance.MaxBandWidth
			}

			if rabbitMQVipInstance.MaxStorage != nil {
				rabbitMQVipInstanceMap["max_storage"] = rabbitMQVipInstance.MaxStorage
			}

			if rabbitMQVipInstance.ExpireTime != nil {
				rabbitMQVipInstanceMap["expire_time"] = rabbitMQVipInstance.ExpireTime
			}

			if rabbitMQVipInstance.AutoRenewFlag != nil {
				rabbitMQVipInstanceMap["auto_renew_flag"] = rabbitMQVipInstance.AutoRenewFlag
			}

			if rabbitMQVipInstance.PayMode != nil {
				rabbitMQVipInstanceMap["pay_mode"] = rabbitMQVipInstance.PayMode
			}

			if rabbitMQVipInstance.Remark != nil {
				rabbitMQVipInstanceMap["remark"] = rabbitMQVipInstance.Remark
			}

			if rabbitMQVipInstance.SpecName != nil {
				rabbitMQVipInstanceMap["spec_name"] = rabbitMQVipInstance.SpecName
			}

			if rabbitMQVipInstance.ExceptionInformation != nil {
				rabbitMQVipInstanceMap["exception_information"] = rabbitMQVipInstance.ExceptionInformation
			}

			ids = append(ids, *rabbitMQVipInstance.InstanceId)
			tmpList = append(tmpList, rabbitMQVipInstanceMap)
		}

		_ = d.Set("instances", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
