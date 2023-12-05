package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudHaVips() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudHaVipsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "Name of the HA VIP. The length of character is limited to 1-60.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the HA VIP to be queried.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VPC id of the HA VIP to be queried.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet id of the HA VIP to be queried.",
			},
			"address_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIp,
				Description:  "EIP of the HA VIP to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"ha_vip_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the dedicated HA VIPs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the HA VIP.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the HA VIP.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC id.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet id.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual IP address, it must not be occupied and in this VPC network segment. If not set, it will be assigned after resource created automatically.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State of the HA VIP. Valid values: `AVAILABLE`, `UNBIND`.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance id that is associated.",
						},
						"network_interface_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network interface id that is associated.",
						},
						"address_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "EIP that is associated.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the HA VIP.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudHaVipsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ha_vips.read")()

	logId := getLogId(contextNil)

	request := vpc.NewDescribeHaVipsRequest()

	params := make(map[string]string)
	if v, ok := d.GetOk("id"); ok {
		params["havip-id"] = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		params["havip-name"] = v.(string)
	}
	if v, ok := d.GetOk("address_ip"); ok {
		params["address-ip"] = v.(string)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		params["subnet-id"] = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		params["vpc-ip"] = v.(string)
	}
	request.Filters = make([]*vpc.Filter, 0, len(params))
	for k, v := range params {
		filter := &vpc.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, filter)
	}
	offset := uint64(0)
	request.Offset = &offset
	result := make([]*vpc.HaVip, 0)
	limit := uint64(HAVIP_DESCRIBE_LIMIT)
	request.Limit = &limit
	for {
		var response *vpc.DescribeHaVipsResponse
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeHaVips(request)
			if e != nil {
				return retryError(errors.WithStack(e))
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read HA VIP failed, reason:%+v", logId, err)
			return err
		} else {
			result = append(result, response.Response.HaVipSet...)
			if len(response.Response.HaVipSet) < HAVIP_DESCRIBE_LIMIT {
				break
			} else {
				offset = offset + limit
			}
		}
	}
	ids := make([]string, 0, len(result))
	haVipList := make([]map[string]interface{}, 0, len(result))
	for _, haVip := range result {
		mapping := map[string]interface{}{
			"id":          *haVip.HaVipId,
			"vip":         *haVip.Vip,
			"name":        *haVip.HaVipName,
			"state":       *haVip.State,
			"vpc_id":      *haVip.VpcId,
			"subnet_id":   *haVip.SubnetId,
			"create_time": *haVip.CreatedTime,
		}
		if haVip.NetworkInterfaceId != nil {
			mapping["network_interface_id"] = *haVip.NetworkInterfaceId
		}
		if haVip.AddressIp != nil {
			mapping["address_ip"] = *haVip.AddressIp
		}
		if haVip.InstanceId != nil {
			mapping["instance_id"] = *haVip.InstanceId
		}
		haVipList = append(haVipList, mapping)
		ids = append(ids, *haVip.HaVipId)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("ha_vip_list", haVipList); e != nil {
		log.Printf("[CRITAL]%s provider set haVip list fail, reason:%s\n", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), haVipList); e != nil {
			return e
		}
	}

	return nil

}
