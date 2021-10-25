/*
Use this data source to query detailed information of VPN gateways.

Example Usage

```hcl
data "tencentcloud_nat_gateway_snats" "snat" {
  nat_gateway_id     = tencentcloud_nat_gateway.my_nat.id
  subnet_id          = tencentcloud_nat_gateway_snat.my_subnet.id
  public_ip_addr     = ["50.29.23.234"]
  description        = "snat demo"
  result_output_file = "./snat.txt"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudNatGatewaySnats() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudNatGatewaySnatRead,

		Schema: map[string]*schema.Schema{
			"nat_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "NAT gateway ID.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet instance ID.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance ID.",
			},
			"public_ip_addr": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Elastic IP address pool.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"snat_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the nat gateway snat.",
				Elem: &schema.Resource{
					Schema: NatGatewaySnatPara(),
				},
			},
		},
	}
}

func dataSourceTencentCloudNatGatewaySnatRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_nat_gateway_snats.read")()

	var (
		logId        = getLogId(contextNil)
		ctx          = context.WithValue(context.TODO(), logIdKey, logId)
		natGatewayId = d.Get("nat_gateway_id").(string)
		vpcService   = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	params := make(map[string]string)
	if v, ok := d.GetOk("subnet_id"); ok {
		params["resource-id"] = v.(string)
	}
	if v, ok := d.GetOk("instance_id"); ok {
		params["resource-id"] = v.(string)
	}
	if v, ok := d.GetOk("public_ip_addr"); ok {
		params["public-ip-address"] = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		params["description"] = v.(string)
	}

	filters := make([]*vpc.Filter, 0, len(params))
	for k, v := range params {
		filter := &vpc.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		filters = append(filters, filter)
	}
	err, result := vpcService.DescribeNatGatewaySnats(ctx, natGatewayId, filters)
	if err != nil {
		log.Printf("[CRITAL]%s read nat gateway snat failed, reason:%s\n ", logId, err.Error())
		return err
	}
	ids := make([]string, len(result))
	snatList := make([]map[string]interface{}, len(result))
	for _, snat := range result {
		m := map[string]interface{}{}
		m["nat_gateway_id"] = snat.NatGatewayId
		m["resource_type"] = snat.ResourceType
		m["public_ip_addr"] = snat.PublicIpAddresses
		m["description"] = snat.Description
		m["snat_id"] = snat.NatGatewaySnatId
		m["create_time"] = snat.CreatedTime
		snatList = append(snatList, m)
		ids = append(ids, *snat.NatGatewaySnatId)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("snat_list", snatList); e != nil {
		log.Printf("[CRITAL]%s provider set nat gateway snat list fail, reason:%s\n ", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), snatList); e != nil {
			return e
		}
	}

	return nil
}
