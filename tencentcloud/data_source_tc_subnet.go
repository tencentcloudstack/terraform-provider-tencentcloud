package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudSubnet() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source has been deprecated in Terraform TencentCloud provider version 1.10.0. Please use 'tencentcloud_vpc_subnets' instead.",
		Read:               dataSourceTencentCloudSubnetRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTencentCloudSubnetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_subnet.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	vpcId := d.Get("vpc_id").(string)
	subnetId := d.Get("subnet_id").(string)

	infos, err := vpcService.DescribeSubnets(ctx, subnetId, vpcId, "", "", map[string]string{})
	if err != nil {
		return err
	}

	if len(infos) == 0 {
		return fmt.Errorf("subnet vpc_id=%s, subnet_id=%s not found", vpcId, subnetId)
	}

	subnet := infos[0]
	d.SetId(subnet.subnetId)
	d.Set("cidr_block", subnet.cidr)
	d.Set("name", subnet.name)
	d.Set("route_table_id", subnet.routeTableId)
	d.Set("availability_zone", subnet.zone)

	return nil
}
