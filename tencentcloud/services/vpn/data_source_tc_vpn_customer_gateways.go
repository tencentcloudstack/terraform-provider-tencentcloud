package vpn

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpnCustomerGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpnCustomerGatewaysRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name of the customer gateway. The length of character is limited to 1-60.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the VPN customer gateway.",
			},
			"public_ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIp,
				Description:  "Public ip address of the VPN customer gateway.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the VPN customer gateway to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"gateway_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the dedicated gateways.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the VPN customer gateway.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the VPN customer gateway.",
						},
						"public_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public ip address of the VPN customer gateway.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags of the VPN customer gateway.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the VPN customer gateway.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudVpnCustomerGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_customer_gateways.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region

	request := vpc.NewDescribeCustomerGatewaysRequest()

	params := make(map[string]string)
	if v, ok := d.GetOk("id"); ok {
		params["customer-gateway-id"] = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		params["customer-gateway-name"] = v.(string)
	}
	if v, ok := d.GetOk("public_ip_address"); ok {
		params["ip-address"] = v.(string)
	}
	tags := helper.GetTags(d, "tags")

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
	result := make([]*vpc.CustomerGateway, 0)
	limit := uint64(svcvpc.VPN_DESCRIBE_LIMIT)
	request.Limit = &limit
	for {
		var response *vpc.DescribeCustomerGatewaysResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeCustomerGateways(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read VPN customer gateway failed, reason:%s\n", logId, err.Error())
			return err
		} else {
			result = append(result, response.Response.CustomerGatewaySet...)
			if len(response.Response.CustomerGatewaySet) < svcvpc.VPN_DESCRIBE_LIMIT {
				break
			} else {
				offset = offset + limit
			}
		}
	}
	ids := make([]string, 0, len(result))
	gatewayList := make([]map[string]interface{}, 0, len(result))
	for _, gateway := range result {
		//tags
		respTags, err := tagService.DescribeResourceTags(ctx, "vpc", "cgw", region, *gateway.CustomerGatewayId)
		if err != nil {
			return err
		}
		if len(tags) > 0 {
			if !reflect.DeepEqual(respTags, tags) {
				continue
			}
		}

		mapping := map[string]interface{}{
			"id":                *gateway.CustomerGatewayId,
			"name":              *gateway.CustomerGatewayName,
			"public_ip_address": *gateway.IpAddress,
			"create_time":       *gateway.CreatedTime,
			"tags":              respTags,
		}
		gatewayList = append(gatewayList, mapping)
		ids = append(ids, *gateway.CustomerGatewayId)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("gateway_list", gatewayList); e != nil {
		log.Printf("[CRITAL]%s provider set gateway list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), gatewayList); e != nil {
			return e
		}
	}

	return nil

}
