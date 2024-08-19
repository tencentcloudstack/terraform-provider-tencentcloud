package cdc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdc/v20201214"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCdcDedicatedClusterInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceTencentCloudCdcDedicatedClusterInstanceTypesRead,
		Schema: map[string]*schema.Schema{
			"dedicated_cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Dedicated Cluster ID.",
			},
			// computed
			"dedicated_cluster_instance_type_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Dedicated Cluster Supported InstanceType.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone Name.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Type.",
						},
						"network_card": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance Type.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance CPU.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance Memory.",
						},
						"instance_family": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Family.",
						},
						"type_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Type Name.",
						},
						"storage_block_amount": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance Storage Block Amount.",
						},
						"instance_bandwidth": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Instance Bandwidth.",
						},
						"instance_pps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance Pps.",
						},
						"cpu_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance CPU Type.",
						},
						"gpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance GPU.",
						},
						"fpga": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance Fpga.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Remark.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Status.",
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

func DataSourceTencentCloudCdcDedicatedClusterInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cdc_dedicated_cluster_instance_types.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                           = tccommon.GetLogId(tccommon.ContextNil)
		ctx                             = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service                         = CdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dedicatedClusterInstanceTypeSet []*cdc.DedicatedClusterInstanceType
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("dedicated_cluster_id"); ok {
		paramMap["DedicatedClusterId"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdcDedicatedClusterInstanceTypesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		dedicatedClusterInstanceTypeSet = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(dedicatedClusterInstanceTypeSet))
	tmpList := make([]map[string]interface{}, 0, len(dedicatedClusterInstanceTypeSet))

	if dedicatedClusterInstanceTypeSet != nil {
		for _, dedicatedClusterInstanceType := range dedicatedClusterInstanceTypeSet {
			dedicatedClusterInstanceTypeMap := map[string]interface{}{}

			if dedicatedClusterInstanceType.Zone != nil {
				dedicatedClusterInstanceTypeMap["zone"] = dedicatedClusterInstanceType.Zone
			}

			if dedicatedClusterInstanceType.InstanceType != nil {
				dedicatedClusterInstanceTypeMap["instance_type"] = dedicatedClusterInstanceType.InstanceType
			}

			if dedicatedClusterInstanceType.NetworkCard != nil {
				dedicatedClusterInstanceTypeMap["network_card"] = dedicatedClusterInstanceType.NetworkCard
			}

			if dedicatedClusterInstanceType.Cpu != nil {
				dedicatedClusterInstanceTypeMap["cpu"] = dedicatedClusterInstanceType.Cpu
			}

			if dedicatedClusterInstanceType.Memory != nil {
				dedicatedClusterInstanceTypeMap["memory"] = dedicatedClusterInstanceType.Memory
			}

			if dedicatedClusterInstanceType.InstanceFamily != nil {
				dedicatedClusterInstanceTypeMap["instance_family"] = dedicatedClusterInstanceType.InstanceFamily
			}

			if dedicatedClusterInstanceType.TypeName != nil {
				dedicatedClusterInstanceTypeMap["type_name"] = dedicatedClusterInstanceType.TypeName
			}

			if dedicatedClusterInstanceType.StorageBlockAmount != nil {
				dedicatedClusterInstanceTypeMap["storage_block_amount"] = dedicatedClusterInstanceType.StorageBlockAmount
			}

			if dedicatedClusterInstanceType.InstanceBandwidth != nil {
				dedicatedClusterInstanceTypeMap["instance_bandwidth"] = dedicatedClusterInstanceType.InstanceBandwidth
			}

			if dedicatedClusterInstanceType.InstancePps != nil {
				dedicatedClusterInstanceTypeMap["instance_pps"] = dedicatedClusterInstanceType.InstancePps
			}

			if dedicatedClusterInstanceType.CpuType != nil {
				dedicatedClusterInstanceTypeMap["cpu_type"] = dedicatedClusterInstanceType.CpuType
			}

			if dedicatedClusterInstanceType.Gpu != nil {
				dedicatedClusterInstanceTypeMap["gpu"] = dedicatedClusterInstanceType.Gpu
			}

			if dedicatedClusterInstanceType.Fpga != nil {
				dedicatedClusterInstanceTypeMap["fpga"] = dedicatedClusterInstanceType.Fpga
			}

			if dedicatedClusterInstanceType.Remark != nil {
				dedicatedClusterInstanceTypeMap["remark"] = dedicatedClusterInstanceType.Remark
			}

			if dedicatedClusterInstanceType.Status != nil {
				dedicatedClusterInstanceTypeMap["status"] = dedicatedClusterInstanceType.Status
			}

			ids = append(ids, *dedicatedClusterInstanceType.InstanceType)
			tmpList = append(tmpList, dedicatedClusterInstanceTypeMap)
		}

		_ = d.Set("dedicated_cluster_instance_type_set", tmpList)
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
