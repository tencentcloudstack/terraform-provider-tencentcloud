package vpc

import (
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDnats() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDnatsRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the VPC.",
			},
			"nat_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the NAT gateway.",
			},
			"elastic_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIp,
				Description:  "Network address of the EIP.",
			},
			"elastic_port": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidatePort,
				Description:  "Port of the EIP.",
			},
			"private_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIp,
				Description:  "Network address of the backend service.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the NAT forward.",
			},
			"private_port": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidatePort,
				Description:  "Port of intranet.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"dnat_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the DNATs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the VPC.",
						},
						"nat_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the NAT.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the network protocol. Valid values: `TCP` and `UDP`.",
						},
						"elastic_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network address of the EIP.",
						},
						"elastic_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port of the EIP.",
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
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of the NAT forward.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDnatsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dnats.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
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
	if v, ok := d.GetOk("description"); ok {
		params["description"] = v.(string)
	}
	request.Filters = make([]*vpc.Filter, 0, len(params))
	for k, v := range params {
		filter := &vpc.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, filter)
	}
	var response *vpc.DescribeNatGatewayDestinationIpPortTranslationNatRulesResponse

	offset := uint64(0)
	request.Offset = &offset
	result := make([]*vpc.NatGatewayDestinationIpPortTranslationNatRule, 0)
	limit := uint64(NAT_DESCRIBE_LIMIT)
	request.Limit = &limit
	for {
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeNatGatewayDestinationIpPortTranslationNatRules(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read DNAT failed, reason:%s\n", logId, err.Error())
			return err
		} else {
			result = append(result, response.Response.NatGatewayDestinationIpPortTranslationNatRuleSet...)
			if len(response.Response.NatGatewayDestinationIpPortTranslationNatRuleSet) < NAT_DESCRIBE_LIMIT {
				break
			} else {
				offset = offset + limit
			}
		}
	}
	ids := make([]string, 0, len(result))
	dnatList := make([]map[string]interface{}, 0, len(result))
	for _, dnat := range result {
		mapping := map[string]interface{}{
			"nat_id":       *dnat.NatGatewayId,
			"vpc_id":       *dnat.VpcId,
			"elastic_ip":   *dnat.PublicIpAddress,
			"elastic_port": strconv.Itoa(int(*dnat.PublicPort)),
			"private_ip":   *dnat.PrivateIpAddress,
			"private_port": strconv.Itoa(int(*dnat.PrivatePort)),
			"protocol":     *dnat.IpProtocol,
			"description":  *dnat.Description,
		}
		dnatList = append(dnatList, mapping)
		var entry = &vpc.DestinationIpPortTranslationNatRule{}
		entry.IpProtocol = dnat.IpProtocol
		entry.PublicIpAddress = dnat.PublicIpAddress
		entry.PublicPort = dnat.PublicPort
		ids = append(ids, buildDnatId(entry, *dnat.VpcId, *dnat.NatGatewayId))
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("dnat_list", dnatList); e != nil {
		log.Printf("[CRITAL]%s provider set DNAT list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dnatList); e != nil {
			return e
		}
	}

	return nil

}
