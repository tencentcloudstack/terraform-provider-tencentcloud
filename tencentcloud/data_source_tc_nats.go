package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func dataSourceTencentCloudNats() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.18.0. Please use 'tencentcloud_nat_gateways' instead.",
		Read:               dataSourceTencentCloudNatsRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
			},
			"state": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_concurrent": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			// Computed values
			"nats": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_concurrent": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"assigned_eip_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
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
	request.Offset = uint64Pt(0)
	request.Limit = uint64Pt(100)

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

	d.SetId(dataResourceIdsHash(ids))
	if e := d.Set("nats", natList); e != nil {
		log.Printf("[CRITAL]%s provider set clb list fail, reason:%s\n ", logId, e.Error())
		return e
	}

	return nil
}
