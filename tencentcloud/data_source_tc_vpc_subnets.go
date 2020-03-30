/*
 Use this data source to query vpc subnets information.

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "guagua_vpc_instance_test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "guagua_vpc_subnet_test"
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        =  "10.0.20.0/28"
  is_multicast      =  false

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_vpc_subnets" "id_instances" {
  subnet_id = tencentcloud_subnet.subnet.id
}

data "tencentcloud_vpc_subnets" "name_instances" {
  name = tencentcloud_subnet.subnet.name
}

data "tencentcloud_vpc_subnets" "tags_instances" {
  tags = tencentcloud_subnet.subnet.tags
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVpcSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcSubnetsRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the VPC to be queried.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the subnet to be queried.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the subnet to be queried.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zone of the subnet to be queried.",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filter default or no default subnets.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the subnet to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"instance_list": {Type: schema.TypeList,
				Computed:    true,
				Description: "List of subnets.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The availability zone of the subnet.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the VPC.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the subnet.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the subnet.",
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A network address block of the subnet.",
						},
						"is_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether it is the default subnet of the VPC for this region.",
						},
						"is_multicast": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether multicast is enabled.",
						},
						"route_table_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the routing table.",
						},
						"available_ip_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of available IPs.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the subnet resource.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags of the subnet resource.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudVpcSubnetsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_subnets.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
	region := meta.(*TencentCloudClient).apiV3Conn.Region

	var (
		vpcId            string
		subnetId         string
		name             string
		availabilityZone string
		isDefault        *bool
	)
	if temp, ok := d.GetOk("vpc_id"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			vpcId = tempStr
		}
	}
	if temp, ok := d.GetOk("subnet_id"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			subnetId = tempStr
		}
	}
	if temp, ok := d.GetOk("name"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			name = tempStr
		}
	}
	if temp, ok := d.GetOk("availability_zone"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			availabilityZone = tempStr
		}
	}

	if temp, ok := d.GetOkExists("is_default"); ok {
		isDefault = helper.Bool(temp.(bool))
	}

	tags := helper.GetTags(d, "tags")

	infos, err := vpcService.DescribeSubnets(ctx, subnetId, vpcId, name, availabilityZone, tags, isDefault)
	if err != nil {
		return err
	}

	var infoList = make([]map[string]interface{}, 0, len(infos))

	for _, item := range infos {
		respTags, err := tagService.DescribeResourceTags(ctx, "vpc", "subnet", region, item.subnetId)
		if err != nil {
			return err
		}

		var infoMap = make(map[string]interface{})

		infoMap["availability_zone"] = item.zone
		infoMap["vpc_id"] = item.vpcId
		infoMap["subnet_id"] = item.subnetId
		infoMap["name"] = item.name
		infoMap["cidr_block"] = item.cidr
		infoMap["is_default"] = item.isDefault
		infoMap["is_multicast"] = item.isMulticast
		infoMap["route_table_id"] = item.routeTableId
		infoMap["available_ip_count"] = item.availableIpCount
		infoMap["create_time"] = item.createTime
		infoMap["tags"] = respTags

		infoList = append(infoList, infoMap)
	}

	if err := d.Set("instance_list", infoList); err != nil {
		log.Printf("[CRITAL]%s provider set  subnet instances fail, reason:%s\n ", logId, err.Error())
		return err
	}

	key := "vpc_subnet" + subnetId + "_" + name

	if isDefault != nil {
		key += "_" + fmt.Sprintf("%v", *isDefault)
	}
	if availabilityZone != "" {
		key += "_" + availabilityZone
	}

	d.SetId(key)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), infoList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}
	return nil
}
