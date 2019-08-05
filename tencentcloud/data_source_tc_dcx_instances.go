/*
Use this data source to query detailed information of dedicated tunnels instances.

Example Usage

```hcl
data "tencentcloud_dcx_instances" "name_select"{
    name = "main"
}

data "tencentcloud_dcx_instances"  "id" {
    dcx_id = "dcx-3ikuw30k"
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
	"strings"
)

func dataSourceTencentCloudDcxInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcxInstancesRead,

		Schema: map[string]*schema.Schema{
			"dcx_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "ID of the dedicated tunnels to be queried.",
			},
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Name of the dedicated tunnels to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the dedicated tunnels.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dcx_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the dedicated tunnel.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the dedicated tunnel.",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the network, and available values include VPC, BMVPC and CCN. The default value is VPC.",
						},
						"dcg_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the DC Gateway. Currently only new in the console.",
						},
						"network_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region of the dedicated tunnel.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the VPC or BMVPC.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bandwidth of the DC.",
						},
						"route_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the route, and available values include BGP and STATIC. The default value is BGP.",
						},
						"bgp_asn": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "BGP ASN of the user.",
						},
						"bgp_auth_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "BGP key of the user.",
						},
						"route_filter_prefixes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Static route, the network address of the user IDC.",
						},
						"vlan": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Vlan of the dedicated tunnels, and the range of values is [0-3000]. '0' means that only one tunnel can be created for the physical connect.",
						},
						"tencent_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Interconnect IP of the DC within Tencent.",
						},
						"customer_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Interconnect IP of the DC within client.",
						},
						"dc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the DC.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State of the dedicated tunnels, and available values include PENDING, ALLOCATING, ALLOCATED, ALTERING, DELETING, DELETED, COMFIRMING and REJECTED.",
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

func dataSourceTencentCloudDcxInstancesRead(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)

	defer LogElapsed(logId + "data_source.tencentcloud_dcx_instances.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		id   = ""
		name = ""
	)
	if temp, ok := d.GetOk("dcx_id"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			id = tempStr
		}
	}
	if temp, ok := d.GetOk("name"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			name = tempStr
		}
	}

	var infos, err = service.DescribeDirectConnectTunnels(ctx, id, name)

	if err != nil {
		return err
	}
	var instanceList = make([]map[string]interface{}, 0, len(infos))

	for _, item := range infos {

		var infoMap = make(map[string]interface{})
		infoMap["dcx_id"] = *item.DirectConnectTunnelId
		infoMap["name"] = *item.DirectConnectTunnelName
		infoMap["network_type"] = strings.ToUpper(service.strPt2str(item.NetworkType))

		infoMap["network_region"] = service.strPt2str(item.NetworkRegion)
		infoMap["vpc_id"] = service.strPt2str(item.VpcId)
		infoMap["bandwidth"] = service.int64Pt2int64(item.Bandwidth)

		infoMap["route_type"] = strings.ToUpper(service.strPt2str(item.RouteType))

		if item.BgpPeer == nil {
			infoMap["bgp_asn"] = 0
			infoMap["bgp_auth_key"] = ""
		} else {
			infoMap["bgp_asn"] = service.int64Pt2int64(item.BgpPeer.Asn)
			infoMap["bgp_auth_key"] = service.strPt2str(item.BgpPeer.AuthKey)
		}

		infoMap["vlan"] = service.int64Pt2int64(item.Vlan)
		infoMap["tencent_address"] = service.strPt2str(item.TencentAddress)
		infoMap["customer_address"] = service.strPt2str(item.CustomerAddress)
		infoMap["dcg_id"] = service.strPt2str(item.DirectConnectGatewayId)

		infoMap["dc_id"] = service.strPt2str(item.DirectConnectId)
		infoMap["state"] = strings.ToUpper(service.strPt2str(item.State))
		infoMap["create_time"] = service.strPt2str(item.CreatedTime)

		var routeFilterPrefixes = make([]string, 0, len(item.RouteFilterPrefixes))
		for _, v := range item.RouteFilterPrefixes {
			if v.Cidr != nil {
				routeFilterPrefixes = append(routeFilterPrefixes, *v.Cidr)
			}
		}
		infoMap["route_filter_prefixes"] = routeFilterPrefixes

		instanceList = append(instanceList, infoMap)
	}

	if err := d.Set("instance_list", instanceList); err != nil {
		log.Printf("[CRITAL]%s provider set  dcx instances fail, reason:%s\n ", logId, err.Error())
		return err
	}

	m := md5.New()
	m.Write([]byte(id + "_" + name))
	d.SetId(fmt.Sprintf("%x", m.Sum(nil)))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), instanceList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}
	return nil
}
