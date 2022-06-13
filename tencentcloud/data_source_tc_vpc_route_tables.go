/*
Use this data source to query vpc route tables information.

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "guagua-ci-temp-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "route_table" {
  vpc_id = tencentcloud_vpc.foo.id
  name   = "ci-temp-test-rt"

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_vpc_route_tables" "id_instances" {
  route_table_id = tencentcloud_route_table.route_table.id
}

data "tencentcloud_vpc_route_tables" "name_instances" {
  name = tencentcloud_route_table.route_table.name
}

data "tencentcloud_vpc_route_tables" "vpc_default_instance" {
  vpc_id           = tencentcloud_vpc.foo.id
  association_main = true
}

data "tencentcloud_vpc_route_tables" "tags_instances" {
  tags = tencentcloud_route_table.route_table.tags
}
```
*/
package tencentcloud

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVpcRouteTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcRouteTablesRead,

		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the routing table to be queried.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the routing table to be queried.",
			},
			"tag_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter if routing table has this tag.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the VPC to be queried.",
			},
			"association_main": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filter the main routing table.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the routing table to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The information list of the VPC route table.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_table_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the routing table.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the routing table.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the VPC.",
						},
						"subnet_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "List of subnet IDs bound to the route table.",
						},
						"is_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether it is the default routing table.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the routing table.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags of the routing table.",
						},
						"route_entry_infos": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Detailed information of each entry of the route table.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"route_entry_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of a route table entry.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description information user defined for a route table rule.",
									},
									"destination_cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The destination address block.",
									},
									"next_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of next-hop, and available values include `CVM`, `VPN`, `DIRECTCONNECT`, `PEERCONNECTION`, `SSLVPN`, `NAT`, `NORMAL_CVM`, `EIP` and `CCN`.",
									},
									"next_hub": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of next-hop gateway. Note: when 'next_type' is EIP, GatewayId will fix the value `0`.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudVpcRouteTablesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_route_tables.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
	region := meta.(*TencentCloudClient).apiV3Conn.Region

	var (
		routeTableId    string
		vpcId           string
		name            string
		associationMain *bool
		tagKey          string
	)
	if temp, ok := d.GetOk("vpc_id"); ok {
		vpcId = temp.(string)
	}

	if temp, ok := d.GetOk("route_table_id"); ok {
		routeTableId = temp.(string)
	}

	if temp, ok := d.GetOk("name"); ok {
		name = temp.(string)
	}

	if temp, ok := d.GetOkExists("association_main"); ok {
		associationMain = helper.Bool(temp.(bool))
	}

	if temp, ok := d.GetOk("tag_key"); ok {
		tagKey = temp.(string)
	}

	var (
		tags  = helper.GetTags(d, "tags")
		infos []VpcRouteTableBasicInfo
		err   error
	)

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		infos, err = service.DescribeRouteTables(ctx, routeTableId, name, vpcId, tags, associationMain, tagKey)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})

	var infoList = make([]map[string]interface{}, 0, len(infos))

	for _, item := range infos {
		routeEntryInfos := make([]map[string]string, len(item.entryInfos))

		for _, v := range item.entryInfos {
			routeEntryInfo := make(map[string]string)
			routeEntryInfo["route_entry_id"] = fmt.Sprintf("%d.%s",
				v.routeEntryId, item.routeTableId)
			routeEntryInfo["description"] = v.description
			routeEntryInfo["destination_cidr_block"] = v.destinationCidr
			routeEntryInfo["next_type"] = v.nextType
			routeEntryInfo["next_hub"] = v.nextBub
			routeEntryInfos = append(routeEntryInfos, routeEntryInfo)
		}

		respTags, err := tagService.DescribeResourceTags(ctx, "vpc", "rtb", region, item.routeTableId)
		if err != nil {
			return err
		}

		var infoMap = make(map[string]interface{})

		infoMap["route_table_id"] = item.routeTableId
		infoMap["name"] = item.name
		infoMap["vpc_id"] = item.vpcId
		infoMap["is_default"] = item.isDefault
		infoMap["subnet_ids"] = item.subnetIds
		infoMap["route_entry_infos"] = routeEntryInfos
		infoMap["create_time"] = item.createTime
		infoMap["tags"] = respTags

		infoList = append(infoList, infoMap)
	}

	if err := d.Set("instance_list", infoList); err != nil {
		log.Printf("[CRITAL]%s provider set  route table instances fail, reason:%s\n ", logId, err.Error())
		return err
	}

	idBytes, err := json.Marshal(map[string]interface{}{
		"routeTableId":    routeTableId,
		"associationMain": associationMain,
		"vpcId":           vpcId,
		"name":            name,
		"tagKey":          tagKey,
		"tags":            tags,
	})
	if err != nil {
		log.Printf("[CRITAL]%s create data source id error, reason:%s\n ", logId, err.Error())
		return err
	}

	md := md5.New()
	_, _ = md.Write(idBytes)
	id := fmt.Sprintf("%x", md.Sum(nil))
	d.SetId(id)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), infoList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}
	return nil
}
