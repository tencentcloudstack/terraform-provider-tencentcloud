/*
Use this data source to query detailed information of CLB

Example Usage

```hcl
data "tencentcloud_clb_instances" "foo" {
  clb_id             = "lb-k2zjp9lv"
  network_type       = "OPEN"
  clb_name           = "myclb"
  project_id         = 0
  result_output_file = "mytestpath"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudClbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbInstancesRead,

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the CLB to be queried.",
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CLB_NETWORK_TYPE),
				Description:  "Type of CLB instance, and available values include `OPEN` and `INTERNAL`.",
			},
			"clb_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the CLB to be queried.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Project ID of the CLB.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"master_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Master available zone id.",
			},
			"clb_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of cloud load balancers. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"clb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of CLB.",
						},
						"clb_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CLB.",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types of CLB.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the project.",
						},
						"clb_vips": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The virtual service address table of the CLB.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The status of CLB.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the CLB.",
						},
						"status_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Latest state transition time of CLB.",
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
						"security_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "ID set of the security groups.",
						},
						"target_region_info_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region information of backend service are attached the CLB.",
						},
						"target_region_info_vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VpcId information of backend service are attached the CLB.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The available tags within this CLB.",
						},
						"address_ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP version, only applicable to open CLB. Valid values are `IPV4`, `IPV6` and `IPv6FullChain`.",
						},
						"vip_isp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network operator, only applicable to open CLB. Valid values are `CMCC`(China Mobile), `CTCC`(Telecom), `CUCC`(China Unicom) and `BGP`. If this ISP is specified, network billing method can only use the bandwidth package billing (BANDWIDTH_PACKAGE).",
						},
						"internet_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Internet charge type, only applicable to open CLB. Valid values are `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.",
						},
						"internet_bandwidth_max_out": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Max bandwidth out, only applicable to open CLB. Valid value ranges is [1, 2048]. Unit is MB.",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Available zone unique id(numerical representation), This field maybe null, means cannot get a valid value",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available zone unique id(string representation), This field maybe null, means cannot get a valid value",
						},
						"zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available zone name, This field maybe null, means cannot get a valid value",
						},
						"zone_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region that this available zone belong to, This field maybe null, means cannot get a valid value",
						},
						"local_zone": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this available zone is local zone, This field maybe null, means cannot get a valid value",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudClbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_clb_instances.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	params := make(map[string]interface{})
	if v, ok := d.GetOk("clb_id"); ok {
		params["clb_id"] = v.(string)
	}
	if v, ok := d.GetOk("clb_name"); ok {
		params["clb_name"] = v.(string)
	}
	if v, ok := d.GetOkExists("project_id"); ok {
		params["project_id"] = v.(int)
	}
	if v, ok := d.GetOk("network_type"); ok {
		params["network_type"] = v.(string)
	}
	if v, ok := d.GetOk("master_zone"); ok {
		params["master_zone"] = v.(string)
	}

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var clbs []*clb.LoadBalancer
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := clbService.DescribeLoadBalancerByFilter(ctx, params)
		if e != nil {
			return retryError(e)
		}
		clbs = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB instances failed, reason:%+v", logId, err)
		return err
	}
	clbList := make([]map[string]interface{}, 0, len(clbs))
	ids := make([]string, 0, len(clbs))
	for _, clbInstance := range clbs {
		mapping := map[string]interface{}{
			"clb_id":                    clbInstance.LoadBalancerId,
			"clb_name":                  clbInstance.LoadBalancerName,
			"network_type":              clbInstance.LoadBalancerType,
			"status":                    clbInstance.Status,
			"create_time":               clbInstance.CreateTime,
			"status_time":               clbInstance.StatusTime,
			"project_id":                clbInstance.ProjectId,
			"vpc_id":                    clbInstance.VpcId,
			"subnet_id":                 clbInstance.SubnetId,
			"clb_vips":                  helper.StringsInterfaces(clbInstance.LoadBalancerVips),
			"target_region_info_region": clbInstance.TargetRegionInfo.Region,
			"target_region_info_vpc_id": clbInstance.TargetRegionInfo.VpcId,
			"address_ip_version":        clbInstance.AddressIPVersion,
			"vip_isp":                   clbInstance.VipIsp,
			"security_groups":           helper.StringsInterfaces(clbInstance.SecureGroups),
		}
		if clbInstance.NetworkAttributes != nil {
			mapping["internet_charge_type"] = *clbInstance.NetworkAttributes.InternetChargeType
			mapping["internet_bandwidth_max_out"] = *clbInstance.NetworkAttributes.InternetMaxBandwidthOut
		}
		if clbInstance.MasterZone != nil {
			mapping["zone_id"] = *clbInstance.MasterZone.ZoneId
			mapping["zone"] = *clbInstance.MasterZone.Zone
			mapping["zone_name"] = *clbInstance.MasterZone.ZoneName
			mapping["zone_region"] = *clbInstance.MasterZone.ZoneRegion
			mapping["local_zone"] = *clbInstance.MasterZone.LocalZone
		}

		if clbInstance.Tags != nil {
			tags := make(map[string]interface{}, len(clbInstance.Tags))
			for _, t := range clbInstance.Tags {
				tags[*t.TagKey] = *t.TagValue
			}
			mapping["tags"] = tags
		}
		clbList = append(clbList, mapping)
		ids = append(ids, *clbInstance.LoadBalancerId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("clb_list", clbList); e != nil {
		log.Printf("[CRITAL]%s provider set CLB list fail, reason:%+v", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), clbList); e != nil {
			return e
		}
	}

	return nil
}
