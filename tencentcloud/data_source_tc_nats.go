package tencentcloud

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	vpc "github.com/zqfan/tencentcloud-sdk-go/services/vpc/unversioned"
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

	conn := meta.(*TencentCloudClient).vpcConn
	args := vpc.NewDescribeNatGatewayRequest()
	args.Offset = common.IntPtr(0)
	args.Limit = common.IntPtr(50)

	if v, ok := d.GetOk("id"); ok {
		args.NatId = common.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		args.VpcId = common.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("name"); ok {
		args.NatName = common.StringPtr(v.(string))
	}

	response, err := conn.DescribeNatGateway(args)

	b, _ := json.Marshal(response)
	log.Printf("[DEBUG] conn.DescribeNatGateway response: %s", b)
	if _, ok := err.(*common.APIError); ok {
		return fmt.Errorf("conn.DescribeNatGateway error: %v", err)
	} else if err != nil {
		return err
	} else if response == nil || len(response.Data) == 0 {
		return fmt.Errorf("no matching NAT gateway found: %v", args)
	}

	var s []map[string]interface{}
	var ids []string

	for _, nat := range response.Data {

		if state, ok := d.GetOk("state"); ok && *nat.State != state.(int) {
			continue
		}
		if max_concurrent, ok := d.GetOk("max_concurrent"); ok && *nat.MaxConcurrent != max_concurrent.(int) {
			continue
		}
		if bandwidth, ok := d.GetOk("bandwidth"); ok && *nat.Bandwidth != bandwidth.(int) {
			continue
		}

		mapping := map[string]interface{}{
			"id":             *nat.NatId,
			"vpc_id":         *nat.UnVpcId,
			"name":           *nat.NatName,
			"state":          *nat.State,
			"max_concurrent": *nat.MaxConcurrent,
			"bandwidth":      *nat.Bandwidth,
			"create_time":    *nat.CreateTime,
		}

		var eips []string
		if len(nat.EipSet) > 0 {
			for _, eip := range nat.EipSet {
				eips = append(eips, *eip)
			}
		}
		mapping["assigned_eip_set"] = eips

		log.Printf("[DEBUG] tencentcloud_nat - adding nat: %v", mapping)
		s = append(s, mapping)
		ids = append(ids, *nat.NatId)
	}

	if len(ids) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	d.SetId(dataResourceIdsHash(ids))

	if err := d.Set("nats", s); err != nil {
		return err
	}
	return nil
}
