/*
Use this data source to query detailed information of NAT gateways.

Example Usage

```hcl
data "tencentcloud_nat_gateways" "foo"{
	name = "main"
	vpc_id = "vpc-xfqag"
	id = "nat-xfaq1"
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func dataSourceTencentCloudNatGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudNatGatewaysRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the vpc.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the nat gateway.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the nat gateway.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"nats": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the dedicated tunnels.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the nat gateway.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the vpc.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the nat gateway.",
						},
						"state": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "State of the nat gateway.",
						},
						"max_concurrent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The upper limit of concurrent connection of nat gateway, the available values include: 1000000,3000000,10000000. Default is 1000000.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum public network output bandwidth of nat gateway (unit: Mbps), the available values include: 20,50,100,200,500,1000,2000,5000. Default is 100.",
						},
						"assigned_eip_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "EIP arrays bound to the gateway. The value of at least 1.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the nat gateway.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudNatGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_nat_gateways.read")()

	logId := getLogId(contextNil)
	request := vpc.NewDescribeNatGatewaysRequest()

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
			Name:   stringToPointer(k),
			Values: []*string{stringToPointer(v)},
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
		mapping := map[string]interface{}{
			"id":               *nat.NatGatewayId,
			"vpc_id":           *nat.VpcId,
			"name":             *nat.NatGatewayName,
			"max_concurrent":   *nat.MaxConcurrentConnection,
			"bandwidth":        *nat.InternetMaxBandwidthOut,
			"state":            *nat.State,
			"assigned_eip_set": flattenAddressList((*nat).PublicIpAddressSet),
			"create_time":      *nat.CreatedTime,
		}
		natList = append(natList, mapping)
		ids = append(ids, *nat.NatGatewayId)
	}
	d.SetId(dataResourceIdsHash(ids))
	if e := d.Set("nats", natList); e != nil {
		log.Printf("[CRITAL]%s provider set clb list fail, reason:%s\n ", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), natList); e != nil {
			return e
		}
	}

	return nil

}
