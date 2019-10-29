package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudRouteTable() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source has been deprecated in Terraform TencentCloud provider version 1.10.0. Please use 'tencentcloud_vpc_route_tables' instead.",
		Read:               dataSourceTencentCloudRouteTableRead,

		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) < 2 || len(value) > 60 {
						errors = append(errors, fmt.Errorf("%s Cannot be less than 2 characters and cannot be longer than 60 characters", k))
					}
					return
				},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hub": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTencentCloudRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_route_table.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		routeTableId string
		name         string
	)
	if temp, ok := d.GetOk("route_table_id"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			routeTableId = tempStr
		}
	}
	if temp, ok := d.GetOk("name"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			name = tempStr
		}
	}

	var infos, err = service.DescribeRouteTables(ctx, routeTableId, name, "", map[string]string{})
	if err != nil {
		return err
	}

	if len(infos) == 0 {
		return fmt.Errorf("route table route_table_id=%s, name=%s not found", routeTableId, name)
	}

	routetable := infos[0]
	d.SetId(routetable.routeTableId)
	d.Set("route_table_id", routetable.routeTableId)
	d.Set("vpc_id", routetable.vpcId)
	d.Set("name", routetable.name)
	d.Set("subnet_num", len(routetable.entryInfos))
	d.Set("create_time", routetable.createTime)

	routes := make([]map[string]interface{}, 0, len(routetable.entryInfos))
	for _, r := range routetable.entryInfos {
		if strings.ToUpper(r.nextBub) == "LOCAL" {
			continue
		}
		m := make(map[string]interface{})
		nextType := ""
		for k, v := range routeTypeNewMap {
			if v == r.nextType {
				nextType = k
				break
			}
		}
		m["next_type"] = nextType
		m["next_hub"] = r.nextBub
		m["cidr_block"] = r.destinationCidr
		m["description"] = r.description
		routes = append(routes, m)
	}

	if err := d.Set("routes", routes); err != nil {
		return err
	}

	return nil
}
