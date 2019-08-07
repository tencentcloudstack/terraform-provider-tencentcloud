/*
Use this data source to query vpc instances' information.

Example Usage

```hcl
resource "tencentcloud_vpc" "foo" {
    name="guagua_vpc_instance_test"
    cidr_block="10.0.0.0/16"
}

data "tencentcloud_vpc_instances" "id_instances" {
	vpc_id="${tencentcloud_vpc.foo.id}"
}

data "tencentcloud_vpc_instances" "name_instances" {
	name="${tencentcloud_vpc.foo.name}"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudVpcInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcInstancesRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "ID of the VPC to be queried.",
			},
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Name of the VPC to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"instance_list": {Type: schema.TypeList,
				Computed:    true,
				Description: "The information list of the VPC.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the VPC.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the VPC.",
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A network address block of a VPC CIDR.",
						},
						"is_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether it is the default VPC for this region.",
						},
						"is_multicast": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether VPC multicast is enabled.",
						},
						"dns_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "A list of DNS servers which can be used within the VPC.",
						},
						"subnet_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "A ID list of subnets within this VPC.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of VPC.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudVpcInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_instances.read")()

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		vpcId string = ""
		name  string = ""
	)
	if temp, ok := d.GetOk("vpc_id"); ok {
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

	var vpcInfos, err = service.DescribeVpcs(ctx, vpcId, name)

	if err != nil {
		return err
	}
	var vpcInfoList = make([]map[string]interface{}, 0, len(vpcInfos))

	for _, item := range vpcInfos {

		var infoMap = make(map[string]interface{})
		infoMap["vpc_id"] = item.vpcId
		infoMap["name"] = item.name
		infoMap["cidr_block"] = item.cidr
		infoMap["is_default"] = item.isDefault
		infoMap["is_multicast"] = item.isMulticast
		infoMap["dns_servers"] = item.dnsServers
		infoMap["create_time"] = item.createTime

		subnetInfos, err := service.DescribeSubnets(ctx, "", item.vpcId, "", "")
		if err != nil {
			return err
		}
		subnetIds := make([]string, 0, len(subnetInfos))
		for _, v := range subnetInfos {
			subnetIds = append(subnetIds, v.subnetId)
		}

		infoMap["subnet_ids"] = subnetIds
		vpcInfoList = append(vpcInfoList, infoMap)
	}

	if err := d.Set("instance_list", vpcInfoList); err != nil {
		log.Printf("[CRITAL]%s provider set  vpc instances fail, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId("vpc_instances" + vpcId + "_" + name)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), vpcInfoList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}
	return nil
}
