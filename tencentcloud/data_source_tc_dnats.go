/*
Use this data source to query detailed information of DNATs.

Example Usage

```hcl
data "tencentcloud_dnats" "foo"{
	name = "main"
	vpc_id = "vpc-xfqag"
	nat_id = "nat-xfaq1"
	elastic_ip = "123.207.115.136"
	elastic_port = "80"
	private_ip = "172.16.0.88"
	private_port = "9001"
}
```
*/
package tencentcloud

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func dataSourceTencentCloudDnats() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDnatsRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the vpc.",
			},
			"nat_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the nat.",
			},
			"elastic_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIp,
				Description:  "Network address of the eip.",
			},
			"elastic_port": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validatePort,
				Description:  " Port of the eip.",
			},
			"private_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIp,
				Description:  "Network address of the backend service.",
			},
			"private_port": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validatePort,
				Description:  "Port of intranet.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"dnat_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the dnats.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the vpc.",
						},
						"nat_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the nat.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the network protocol, the available values include： TCP and UDP.",
						},
						"elastic_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network address of the eip.",
						},
						"elastic_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: " Port of the eip.",
						},
						"private_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network address of the backend service.",
						},
						"private_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port of intranet.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDnatsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dnats.read")()

	logId := getLogId(contextNil)
	request := vpc.NewDescribeNatGatewayDestinationIpPortTranslationNatRulesRequest()
	params := make(map[string]string)
	if v, ok := d.GetOk("nat_id"); ok {
		params["nat-gateway-id"] = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		params["vpc-id"] = v.(string)
	}
	if v, ok := d.GetOk("elastic_ip"); ok {
		params["public-ip-address"] = v.(string)
	}
	if v, ok := d.GetOk("elastic_port"); ok {
		params["public-port"] = v.(string)
	}
	if v, ok := d.GetOk("private_ip"); ok {
		params["private-ip-address"] = v.(string)
	}
	if v, ok := d.GetOk("private_port"); ok {
		params["private-port"] = v.(string)
	}
	request.Filters = make([]*vpc.Filter, 0, len(params))
	for k, v := range params {
		filter := &vpc.Filter{
			Name:   stringToPointer(k),
			Values: []*string{stringToPointer(v)},
		}
		request.Filters = append(request.Filters, filter)
	}
	var response *vpc.DescribeNatGatewayDestinationIpPortTranslationNatRulesResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeNatGatewayDestinationIpPortTranslationNatRules(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read dnat  failed, reason:%s\n ", logId, err.Error())
		return err
	}
	dnats := response.Response.NatGatewayDestinationIpPortTranslationNatRuleSet
	ids := make([]string, len(dnats))
	dnatList := make([]map[string]interface{}, 0, len(dnats))
	for _, dnat := range dnats {
		mapping := map[string]interface{}{
			"nat_id":       *dnat.NatGatewayId,
			"vpc_id":       *dnat.VpcId,
			"elastic_ip":   *dnat.PublicIpAddress,
			"elastic_port": strconv.Itoa(int(*dnat.PublicPort)),
			"private_ip":   *dnat.PrivateIpAddress,
			"private_port": strconv.Itoa(int(*dnat.PrivatePort)),
			"protocol":     *dnat.IpProtocol,
		}
		dnatList = append(dnatList, mapping)
		var entry = &vpc.DestinationIpPortTranslationNatRule{}
		entry.IpProtocol = stringToPointer(*dnat.IpProtocol)
		entry.PublicIpAddress = dnat.PublicIpAddress
		entry.PublicPort = dnat.PublicPort
		ids = append(ids, buildDnatId(entry, *dnat.VpcId, *dnat.NatGatewayId))
	}
	d.SetId(dataResourceIdsHash(ids))
	if e := d.Set("dnat_list", dnatList); e != nil {
		log.Printf("[CRITAL]%s provider set clb list fail, reason:%s\n ", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), dnatList); e != nil {
			return e
		}
	}

	return nil

}
