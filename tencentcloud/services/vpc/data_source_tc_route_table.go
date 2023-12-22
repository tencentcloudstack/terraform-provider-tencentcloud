package vpc

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudRouteTable() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source has been deprecated in Terraform TencentCloud provider version 1.10.0. Please use 'tencentcloud_vpc_route_tables' instead.",
		Read:               dataSourceTencentCloudRouteTableRead,

		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Route Table ID.",
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
				Description: "The Route Table name.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The VPC ID.",
			},
			"subnet_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of associated subnets.",
			},
			"routes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The information list of the VPC route table.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The RouteEntry's target network segment.",
						},
						"next_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The `next_hub` type.",
						},
						"next_hub": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The RouteEntry's next hub.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The RouteEntry's description.",
						},
					},
				},
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of routing table.",
			},
		},
	}
}

func dataSourceTencentCloudRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_route_table.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

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

	var infos, err = service.DescribeRouteTables(ctx, routeTableId, name, "", map[string]string{}, nil, "")
	if err != nil {
		return err
	}

	if len(infos) == 0 {
		return fmt.Errorf("route table route_table_id=%s, name=%s not found", routeTableId, name)
	}

	routetable := infos[0]
	d.SetId(routetable.routeTableId)
	_ = d.Set("route_table_id", routetable.routeTableId)
	_ = d.Set("vpc_id", routetable.vpcId)
	_ = d.Set("name", routetable.name)
	_ = d.Set("subnet_num", len(routetable.entryInfos))
	_ = d.Set("create_time", routetable.createTime)

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
