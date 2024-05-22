// Code generated by iacg; DO NOT EDIT.
package cvm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudReservedInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudReservedInstancesRead,
		Schema: map[string]*schema.Schema{
			"reserved_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the reserved instance to be query.",
			},

			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The available zone that the reserved instance locates at.",
			},

			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of reserved instance.",
			},

			"reserved_instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of reserved instance. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reserved_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the reserved instance.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of reserved instance.",
						},
						"instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of reserved instance.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone of the reserved instance.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start time of the reserved instance.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expiry time of the reserved instance.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the reserved instance.",
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

func dataSourceTencentCloudReservedInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_reserved_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	var respData []*cvm.ReservedInstances
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeReservedInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	var ids []string
	reservedInstancesSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, reservedInstancesSet := range respData {
			reservedInstancesSetMap := map[string]interface{}{}

			var reservedInstancesId string
			if reservedInstancesSet.ReservedInstancesId != nil {
				reservedInstancesSetMap["reserved_instance_id"] = reservedInstancesSet.ReservedInstancesId
				reservedInstancesId = *reservedInstancesSet.ReservedInstancesId
			}

			if reservedInstancesSet.InstanceType != nil {
				reservedInstancesSetMap["instance_type"] = reservedInstancesSet.InstanceType
			}

			if reservedInstancesSet.InstanceCount != nil {
				reservedInstancesSetMap["instance_count"] = reservedInstancesSet.InstanceCount
			}

			if reservedInstancesSet.Zone != nil {
				reservedInstancesSetMap["availability_zone"] = reservedInstancesSet.Zone
			}

			if reservedInstancesSet.StartTime != nil {
				reservedInstancesSetMap["start_time"] = reservedInstancesSet.StartTime
			}

			if reservedInstancesSet.EndTime != nil {
				reservedInstancesSetMap["end_time"] = reservedInstancesSet.EndTime
			}

			if reservedInstancesSet.State != nil {
				reservedInstancesSetMap["status"] = reservedInstancesSet.State
			}

			ids = append(ids, reservedInstancesId)
			reservedInstancesSetList = append(reservedInstancesSetList, reservedInstancesSetMap)
		}

		_ = d.Set("reserved_instance_list", reservedInstancesSetList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), reservedInstancesSetList); e != nil {
			return e
		}
	}

	return nil
}
