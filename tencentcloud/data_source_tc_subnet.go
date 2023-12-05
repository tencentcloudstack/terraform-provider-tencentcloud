package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTencentCloudSubnet() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source has been deprecated in Terraform TencentCloud provider version 1.10.0. Please use 'tencentcloud_vpc_subnets' instead.",
		Read:               dataSourceTencentCloudSubnetRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPC ID.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Subnet.",
			},
			"cidr_block": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CIDR block of the Subnet.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The AZ for the subnet.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name for the Subnet.",
			},
			"route_table_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Route Table ID.",
			},
		},
	}
}

func dataSourceTencentCloudSubnetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_subnet.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	vpcId := d.Get("vpc_id").(string)
	subnetId := d.Get("subnet_id").(string)

	infos, err := vpcService.DescribeSubnets(ctx, subnetId, vpcId, "", "", map[string]string{}, nil, nil, "", "")
	if err != nil {
		return err
	}

	if len(infos) == 0 {
		return fmt.Errorf("subnet vpc_id=%s, subnet_id=%s not found", vpcId, subnetId)
	}

	subnet := infos[0]
	d.SetId(subnet.subnetId)
	_ = d.Set("cidr_block", subnet.cidr)
	_ = d.Set("name", subnet.name)
	_ = d.Set("route_table_id", subnet.routeTableId)
	_ = d.Set("availability_zone", subnet.zone)

	return nil
}
