package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcNetworkInterfaceLimit() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcNetworkInterfaceLimitRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of a CVM instance or ENI to query.",
			},

			"eni_quantity": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Quota of ENIs mounted to a CVM instance in a standard way.",
			},

			"eni_private_ip_address_quantity": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Quota of IP addresses that can be allocated to each standard-mounted ENI.",
			},

			"extend_eni_quantity": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Quota of ENIs mounted to a CVM instance as an extensionNote: this field may return `null`, indicating that no valid values can be obtained.",
			},

			"extend_eni_private_ip_address_quantity": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Quota of IP addresses that can be allocated to each extension-mounted ENI.Note: this field may return `null`, indicating that no valid values can be obtained.",
			},

			"sub_eni_quantity": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The quota of relayed ENIsNote: This field may return `null`, indicating that no valid values can be obtained.",
			},

			"sub_eni_private_ip_address_quantity": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The quota of IPs that can be assigned to each relayed ENI.Note: This field may return `null`, indicating that no valid values can be obtained.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudVpcNetworkInterfaceLimitRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_network_interface_limit.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var instanceId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var networkInterfaceLimit *vpc.DescribeNetworkInterfaceLimitResponseParams

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcNetworkInterfaceLimit(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		networkInterfaceLimit = result
		return nil
	})
	if err != nil {
		return err
	}

	limitMap := map[string]interface{}{}

	if networkInterfaceLimit.EniQuantity != nil {
		_ = d.Set("eni_quantity", networkInterfaceLimit.EniQuantity)
		limitMap["eni_quantity"] = networkInterfaceLimit.EniQuantity
	}

	if networkInterfaceLimit.EniPrivateIpAddressQuantity != nil {
		_ = d.Set("eni_private_ip_address_quantity", networkInterfaceLimit.EniPrivateIpAddressQuantity)
		limitMap["eni_private_ip_address_quantity"] = networkInterfaceLimit.EniPrivateIpAddressQuantity
	}

	if networkInterfaceLimit.ExtendEniQuantity != nil {
		_ = d.Set("extend_eni_quantity", networkInterfaceLimit.ExtendEniQuantity)
		limitMap["extend_eni_quantity"] = networkInterfaceLimit.ExtendEniQuantity
	}

	if networkInterfaceLimit.ExtendEniPrivateIpAddressQuantity != nil {
		_ = d.Set("extend_eni_private_ip_address_quantity", networkInterfaceLimit.ExtendEniPrivateIpAddressQuantity)
		limitMap["extend_eni_private_ip_address_quantity"] = networkInterfaceLimit.ExtendEniPrivateIpAddressQuantity
	}

	if networkInterfaceLimit.SubEniQuantity != nil {
		_ = d.Set("sub_eni_quantity", networkInterfaceLimit.SubEniQuantity)
		limitMap["sub_eni_quantity"] = networkInterfaceLimit.SubEniQuantity
	}

	if networkInterfaceLimit.SubEniPrivateIpAddressQuantity != nil {
		_ = d.Set("sub_eni_private_ip_address_quantity", networkInterfaceLimit.SubEniPrivateIpAddressQuantity)
		limitMap["sub_eni_private_ip_address_quantity"] = networkInterfaceLimit.SubEniPrivateIpAddressQuantity
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), limitMap); e != nil {
			return e
		}
	}
	return nil
}
