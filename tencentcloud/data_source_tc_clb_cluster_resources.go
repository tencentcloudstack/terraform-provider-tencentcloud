/*
Use this data source to query detailed information of clb cluster_resources

Example Usage

```hcl
data "tencentcloud_clb_cluster_resources" "cluster_resources" {
  filters {
    name = "idle"
    values = ["True"]
  }
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudClbClusterResources() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbClusterResourcesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions to query cluster. cluster-id - String - Required: No - (Filter condition) Filter by cluster ID, such as tgw-12345678. vip - String - Required: No - (Filter condition) Filter by loadbalancer vip, such as 192.168.0.1. loadblancer-id - String - Required: No - (Filter condition) Filter by loadblancer ID, such as lbl-12345678. idle - String - Required: No - (Filter condition) Filter by Whether load balancing is idle, such as True, False.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter name.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Filter values.",
						},
					},
				},
			},

			"cluster_resource_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Cluster resource set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster ID.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "vip.",
						},
						"load_balancer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Loadbalance Id.",
						},
						"idle": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Is it idle.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cluster name.",
						},
						"isp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Isp.",
						},
						"clusters_zone": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "clusters zone.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"master_zone": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Availability master zone where the cluster is located.",
									},
									"slave_zone": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Availability slave zone where the cluster is located.",
									},
								},
							},
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudClbClusterResourcesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_clb_cluster_resources.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*clb.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := clb.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var clusterResourceSet []*clb.ClusterResource

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbClusterResourcesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		clusterResourceSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(clusterResourceSet))
	tmpList := make([]map[string]interface{}, 0, len(clusterResourceSet))

	if clusterResourceSet != nil {
		for _, clusterResource := range clusterResourceSet {
			clusterResourceMap := map[string]interface{}{}

			if clusterResource.ClusterId != nil {
				clusterResourceMap["cluster_id"] = clusterResource.ClusterId
			}

			if clusterResource.Vip != nil {
				clusterResourceMap["vip"] = clusterResource.Vip
			}

			if clusterResource.LoadBalancerId != nil {
				clusterResourceMap["load_balancer_id"] = clusterResource.LoadBalancerId
			}

			if clusterResource.Idle != nil {
				clusterResourceMap["idle"] = clusterResource.Idle
			}

			if clusterResource.ClusterName != nil {
				clusterResourceMap["cluster_name"] = clusterResource.ClusterName
			}

			if clusterResource.Isp != nil {
				clusterResourceMap["isp"] = clusterResource.Isp
			}

			if clusterResource.ClustersZone != nil {
				clustersZoneMap := map[string]interface{}{}

				if clusterResource.ClustersZone.MasterZone != nil {
					clustersZoneMap["master_zone"] = clusterResource.ClustersZone.MasterZone
				}

				if clusterResource.ClustersZone.SlaveZone != nil {
					clustersZoneMap["slave_zone"] = clusterResource.ClustersZone.SlaveZone
				}

				clusterResourceMap["clusters_zone"] = []interface{}{clustersZoneMap}
			}

			ids = append(ids, *clusterResource.ClusterId)
			tmpList = append(tmpList, clusterResourceMap)
		}

		_ = d.Set("cluster_resource_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
