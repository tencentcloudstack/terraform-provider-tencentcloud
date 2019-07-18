/*
Use this data source to query detailed information of CLB

Example Usage

```hcl
data "tencentcloud_clb" "clblab" {
    clb_id             = "lb-k2zjp9lv"
    network_type           = "OPEN"
    clb_name           = "myclb"
    project_id         = "Default Project"
    result_output_file = "mytestpath"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudClbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbInstancesRead,

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: " ID of the CLB to be queried.",
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CLB_NETWORK_TYPE),
				Description:  "Type of CLB instance, and available values include 'OPEN' and 'INTERNAL'",
			},
			"clb_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The available zone that the CBS instance locates at.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the project within the CLB.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
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
							Description: "Creation time of the CLB",
						},
						"status_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Latest state transition time of CLB.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the VPC",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the subnet",
						},
						"security_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "ID of the security groups.",
						},
						"target_region_info": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Information of backend service are attached the CLB.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The available tags within this CLB.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudClbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	params := make(map[string]interface{})
	if v, ok := d.GetOk("clb_id"); ok {
		params["clb_id"] = v.(string)
	}
	if v, ok := d.GetOk("clb_name"); ok {
		params["clb_name"] = v.(string)
	}
	if v, ok := d.GetOk("project_id"); ok {
		params["project_id"] = v.(int)
	}
	if v, ok := d.GetOk("network_type"); ok {
		params["network_type"] = v.(string)
	}

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	clbs, err := clbService.DescribeLoadBalancerByFilter(ctx, params)
	if err != nil {
		return err
	}

	clbList := make([]map[string]interface{}, 0, len(clbs))
	ids := make([]string, 0, len(clbs))
	for _, clb := range clbs {
		mapping := map[string]interface{}{
			"clb_id":             *clb.LoadBalancerId,
			"clb_name":           *clb.LoadBalancerName,
			"network_type":       *clb.LoadBalancerType,
			"status":             *clb.Status,
			"create_time":        *clb.CreateTime,
			"status_time":        *clb.StatusTime,
			"project_id":         *clb.ProjectId,
			"vpc_id":             *clb.VpcId,
			"subnet_id":          *clb.SubnetId,
			"clb_vips":           flattenStringList(clb.LoadBalancerVips),
			"target_region_info": flattenClbTargetRegionInfoMappings(clb.TargetRegionInfo),
			"security_groups":    flattenStringList(clb.SecureGroups),
			"tags":               flattenClbTagsMapping(clb.Tags),
		}

		clbList = append(clbList, mapping)
		ids = append(ids, *clb.LoadBalancerId)
	}

	d.SetId(dataResourceIdsHash(ids))
	if err = d.Set("clb_list", clbList); err != nil {
		log.Printf("[CRITAL]%s provider set clb list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), clbList); err != nil {
			return err
		}
	}

	return nil
}
