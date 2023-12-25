package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcUsedIpAddress() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcUsedIpAddressRead,
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "VPC instance ID.",
			},

			"subnet_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Subnet instance ID.",
			},

			"ip_addresses": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "IPs to query.",
			},

			"ip_address_states": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Information of resources bound with the queried IPs Note: This parameter may return null, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC instance ID.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet instance ID.",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP address.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type.",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource ID.",
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

func dataSourceTencentCloudVpcUsedIpAddressRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_used_ip_address.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("vpc_id"); ok {
		paramMap["VpcId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		paramMap["SubnetId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_addresses"); ok {
		ipAddressesSet := v.(*schema.Set).List()
		paramMap["IpAddresses"] = helper.InterfacesStringsPoint(ipAddressesSet)
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var ipAddressStates []*vpc.IpAddressStates

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcUsedIpAddressByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		ipAddressStates = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(ipAddressStates))
	tmpList := make([]map[string]interface{}, 0, len(ipAddressStates))

	if ipAddressStates != nil {
		for _, ipAddressStates := range ipAddressStates {
			ipAddressStatesMap := map[string]interface{}{}

			if ipAddressStates.VpcId != nil {
				ipAddressStatesMap["vpc_id"] = ipAddressStates.VpcId
			}

			if ipAddressStates.SubnetId != nil {
				ipAddressStatesMap["subnet_id"] = ipAddressStates.SubnetId
			}

			if ipAddressStates.IpAddress != nil {
				ipAddressStatesMap["ip_address"] = ipAddressStates.IpAddress
			}

			if ipAddressStates.ResourceType != nil {
				ipAddressStatesMap["resource_type"] = ipAddressStates.ResourceType
			}

			if ipAddressStates.ResourceId != nil {
				ipAddressStatesMap["resource_id"] = ipAddressStates.ResourceId
			}

			ids = append(ids, *ipAddressStates.VpcId)
			tmpList = append(tmpList, ipAddressStatesMap)
		}

		_ = d.Set("ip_address_states", tmpList)
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
