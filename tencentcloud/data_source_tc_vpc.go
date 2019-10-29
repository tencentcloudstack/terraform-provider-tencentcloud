package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudVpc() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source has been deprecated in Terraform TencentCloud provider version 1.10.0. Please use 'tencentcloud_vpc_instances' instead.",
		Read:               dataSourceTencentCloudVpcRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_multicast": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceTencentCloudVpcRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		vpcId string
		name  string
	)
	if temp, ok := d.GetOk("id"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			vpcId = tempStr
		}
	}
	if temp, ok := d.GetOk("name"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			name = tempStr
		}
	}

	var vpcInfos, err = service.DescribeVpcs(ctx, vpcId, name, map[string]string{})
	if err != nil {
		return err
	}

	if len(vpcInfos) == 0 {
		return fmt.Errorf("vpc id=%s, name=%s not found", vpcId, name)
	}

	vpc := vpcInfos[0]
	d.SetId(vpc.vpcId)
	d.Set("name", vpc.name)
	d.Set("cidr_block", vpc.cidr)
	d.Set("is_default", vpc.isDefault)
	d.Set("is_multicast", vpc.isMulticast)

	return nil
}
