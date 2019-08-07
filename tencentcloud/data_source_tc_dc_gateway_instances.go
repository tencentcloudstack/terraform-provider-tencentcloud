/*
Use this data source to query detailed information of direct connect gateway instances.

Example Usage

```hcl
resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource "tencentcloud_dc_gateway" "ccn_main" {
  name                = "ci-cdg-ccn-test"
  network_instance_id = "${tencentcloud_ccn.main.id}"
  network_type        = "CCN"
  gateway_type        = "NORMAL"
}

#You need to sleep for a few seconds because there is a cache on the server
data "tencentcloud_dc_gateway_instances" "name_select"{
  name = "${tencentcloud_dc_gateway.ccn_main.name}"
}

data "tencentcloud_dc_gateway_instances"  "id_select" {
  dcg_id = "${tencentcloud_dc_gateway.ccn_main.id}"
}
```
*/
package tencentcloud

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceTencentCloudDcGatewayInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcGatewayInstancesRead,
		Schema: map[string]*schema.Schema{
			"dcg_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the DCG to be queried.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the DCG to be queried.",
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
				Description: "Information list of the DCG.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dcg_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the DCG",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the DCG",
						},
						"dcg_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP of the DCG",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP of the DCG",
						},
						"network_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of associated network, the available value include 'VPC' and 'CCN'.",
						},
						"gateway_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the gateway, the available value include 'NORMAL' and 'NAT'. Default is 'NORMAL'.",
						},
						"cnn_route_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of CCN route, the available value include 'BGP' and 'STATIC'.",
						},
						"enable_bgp": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the BGP is enabled.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of resource.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDcGatewayInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dcgateway_instances.read")()

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		id   string = ""
		name string = ""
	)

	if temp, ok := d.GetOk("dcg_id"); ok {
		if tempStr := temp.(string); tempStr != "" {
			id = tempStr
		}
	}

	if temp, ok := d.GetOk("name"); ok {
		if tempStr := temp.(string); tempStr != "" {
			name = tempStr
		}
	}
	var infos, err = service.DescribeDirectConnectGateways(ctx, id, name)
	if err != nil {
		return err
	}

	var infoList = make([]map[string]interface{}, 0, len(infos))

	for _, item := range infos {
		var infoMap = make(map[string]interface{})
		infoMap["dcg_id"] = item.dcgId
		infoMap["name"] = item.name
		infoMap["dcg_ip"] = item.dcgIp
		infoMap["network_type"] = item.networkType
		infoMap["network_instance_id"] = item.networkInstanceId
		infoMap["gateway_type"] = item.gatewayType
		infoMap["cnn_route_type"] = item.cnnRouteType
		infoMap["create_time"] = item.createTime
		infoMap["enable_bgp"] = item.enableBGP
		infoList = append(infoList, infoMap)
	}
	if err := d.Set("instance_list", infoList); err != nil {
		log.Printf("[CRITAL]%s provider set  dcg instances fail, reason:%s\n ",
			logId,
			err.Error())
		return err
	}

	m := md5.New()
	m.Write([]byte("dcg_instances" + id + "_" + name))
	d.SetId(fmt.Sprintf("%x", m.Sum(nil)))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), infoList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId,
				output.(string),
				err.Error())
			return err
		}
	}
	return nil

}
