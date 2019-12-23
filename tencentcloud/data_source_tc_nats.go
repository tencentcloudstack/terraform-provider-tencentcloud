/*
The NATs data source lists a number of NATs resource information owned by an TencentCloud account.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_nat_gateways.

Example Usage

```hcl
# Query the NAT gateway by ID
data "tencentcloud_nats" "anat" {
  id = "nat-k6ualnp2"
}

# Query the list of normal NAT gateways
data "tencentcloud_nats" "nat_state" {
  state = 0
}

# Multi conditional query NAT gateway list
data "tencentcloud_nats" "multi_nat" {
  name           = "terraform test"
  vpc_id         = "vpc-ezij4ltv"
  max_concurrent = 3000000
  bandwidth      = 500
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudNats() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.18.0. Please use 'tencentcloud_nat_gateways' instead.",
		Read:               dataSourceTencentCloudNatsRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID for NAT Gateway.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The VPC ID for NAT Gateway.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "The name for NAT Gateway.",
			},
			"state": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "NAT gateway status, 0: Running, 1: Unavailable, 2: Be in arrears and out of service.",
			},
			"max_concurrent": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The upper limit of concurrent connection of NAT gateway, for example: 1000000, 3000000, 10000000.",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The maximum public network output bandwidth of the gateway (unit: Mbps), for example: 10, 20, 50, 100, 200, 500, 1000, 2000, 5000.",
			},
			"nats": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the dedicated tunnels.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID for NAT Gateway.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC ID for NAT Gateway.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for NAT Gateway.",
						},
						"state": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "NAT gateway status, 0: Running, 1: Unavailable, 2: Be in arrears and out of service.",
						},
						"max_concurrent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The upper limit of concurrent connection of NAT gateway, for example: 1000000, 3000000, 10000000.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum public network output bandwidth of the gateway (unit: Mbps), for example: 10, 20, 50, 100, 200, 500, 1000, 2000, 5000.",
						},
						"assigned_eip_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Elastic IP arrays bound to the gateway.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The create time of the NAT gateway.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudNatsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_nats.read")()

	logId := getLogId(contextNil)
	request := vpc.NewDescribeNatGatewaysRequest()
	request.Offset = helper.Uint64(0)
	request.Limit = helper.Uint64(100)

	params := make(map[string]string)
	if v, ok := d.GetOk("id"); ok {
		params["nat-gateway-id"] = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		params["nat-gateway-name"] = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		params["vpc-id"] = v.(string)
	}

	request.Filters = make([]*vpc.Filter, 0, len(params))
	for k, v := range params {
		filter := &vpc.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, filter)
	}

	var response *vpc.DescribeNatGatewaysResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeNatGateways(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read nat gateway failed, reason:%s\n ", logId, err.Error())
		return err
	}

	ids := make([]string, 0, len(response.Response.NatGatewaySet))
	natList := make([]map[string]interface{}, 0, len(response.Response.NatGatewaySet))
	for _, nat := range response.Response.NatGatewaySet {
		networkState := 0
		switch *nat.NetworkState {
		case "AVAILABLE":
			networkState = 0
		case "UNAVAILABLE":
			networkState = 1
		case "INSUFFICIENT":
			networkState = 2
		}

		if state, ok := d.GetOk("state"); ok && networkState != state.(int) {
			continue
		}
		if max_concurrent, ok := d.GetOk("max_concurrent"); ok && *nat.MaxConcurrentConnection != uint64(max_concurrent.(int)) {
			continue
		}
		if bandwidth, ok := d.GetOk("bandwidth"); ok && *nat.InternetMaxBandwidthOut != uint64(bandwidth.(int)) {
			continue
		}

		mapping := map[string]interface{}{
			"id":               *nat.NatGatewayId,
			"vpc_id":           *nat.VpcId,
			"name":             *nat.NatGatewayName,
			"max_concurrent":   *nat.MaxConcurrentConnection,
			"bandwidth":        *nat.InternetMaxBandwidthOut,
			"state":            networkState,
			"assigned_eip_set": flattenAddressList((*nat).PublicIpAddressSet),
			"create_time":      *nat.CreatedTime,
		}
		natList = append(natList, mapping)
		ids = append(ids, *nat.NatGatewayId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("nats", natList); e != nil {
		log.Printf("[CRITAL]%s provider set clb list fail, reason:%s\n ", logId, e.Error())
		return e
	}

	return nil
}
