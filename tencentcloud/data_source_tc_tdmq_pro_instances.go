/*
Use this data source to query detailed information of tdmq pro_instances

Example Usage

```hcl
data "tencentcloud_tdmq_pro_instances" "pro_instances_filter" {
  filters {
    name   = "InstanceName"
    values = ["keep"]
  }
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTdmqProInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqProInstancesRead,
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
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "value.",
						},
					},
				},
			},
			"instances": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Instance information list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"instance_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance version.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance status, 0-creating, 1-normal, 2-isolating, 3-destroyed, 4-abnormal, 5-delivery failure, 6-allocation change, 7-allocation failure.",
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
						"scalable_tps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Elastic TPS outside specificationNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the VPCNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet idNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"max_band_width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Peak bandwidth. Unit: mbps.",
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

func dataSourceTencentCloudTdmqProInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tdmq_pro_instances.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		service   = TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}
		instances []*tdmq.PulsarProInstance
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

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqProInstancesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
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
		for _, pulsarProInstance := range instances {
			pulsarProInstanceMap := map[string]interface{}{}

			if pulsarProInstance.InstanceId != nil {
				pulsarProInstanceMap["instance_id"] = pulsarProInstance.InstanceId
			}

			if pulsarProInstance.InstanceName != nil {
				pulsarProInstanceMap["instance_name"] = pulsarProInstance.InstanceName
			}

			if pulsarProInstance.InstanceVersion != nil {
				pulsarProInstanceMap["instance_version"] = pulsarProInstance.InstanceVersion
			}

			if pulsarProInstance.Status != nil {
				pulsarProInstanceMap["status"] = pulsarProInstance.Status
			}

			if pulsarProInstance.ConfigDisplay != nil {
				pulsarProInstanceMap["config_display"] = pulsarProInstance.ConfigDisplay
			}

			if pulsarProInstance.MaxTps != nil {
				pulsarProInstanceMap["max_tps"] = pulsarProInstance.MaxTps
			}

			if pulsarProInstance.MaxStorage != nil {
				pulsarProInstanceMap["max_storage"] = pulsarProInstance.MaxStorage
			}

			if pulsarProInstance.ExpireTime != nil {
				pulsarProInstanceMap["expire_time"] = pulsarProInstance.ExpireTime
			}

			if pulsarProInstance.AutoRenewFlag != nil {
				pulsarProInstanceMap["auto_renew_flag"] = pulsarProInstance.AutoRenewFlag
			}

			if pulsarProInstance.PayMode != nil {
				pulsarProInstanceMap["pay_mode"] = pulsarProInstance.PayMode
			}

			if pulsarProInstance.Remark != nil {
				pulsarProInstanceMap["remark"] = pulsarProInstance.Remark
			}

			if pulsarProInstance.SpecName != nil {
				pulsarProInstanceMap["spec_name"] = pulsarProInstance.SpecName
			}

			if pulsarProInstance.ScalableTps != nil {
				pulsarProInstanceMap["scalable_tps"] = pulsarProInstance.ScalableTps
			}

			if pulsarProInstance.VpcId != nil {
				pulsarProInstanceMap["vpc_id"] = pulsarProInstance.VpcId
			}

			if pulsarProInstance.SubnetId != nil {
				pulsarProInstanceMap["subnet_id"] = pulsarProInstance.SubnetId
			}

			if pulsarProInstance.MaxBandWidth != nil {
				pulsarProInstanceMap["max_band_width"] = pulsarProInstance.MaxBandWidth
			}

			ids = append(ids, *pulsarProInstance.InstanceId)
			tmpList = append(tmpList, pulsarProInstanceMap)
		}

		_ = d.Set("instances", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
