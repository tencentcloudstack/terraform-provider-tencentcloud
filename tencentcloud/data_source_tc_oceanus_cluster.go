/*
Use this data source to query detailed information of oceanus cluster

Example Usage

```hcl
data "tencentcloud_oceanus_cluster" "cluster" {
  cluster_ids =
  order_type = 1
  filters {
		name = "name"
		values =

  }
  work_space_id = "space-1239"
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudOceanusCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusClusterRead,
		Schema: map[string]*schema.Schema{
			"cluster_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Query one or more clusters by their ID. The maximum number of clusters that can be queried at once is 100.",
			},

			"order_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The sorting rule of the cluster information results. Possible values are 1 (sort by time in descending order), 2 (sort by time in ascending order), and 3 (sort by status).",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "The filtering rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The field to be filtered.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "The filtering values of the field.",
						},
					},
				},
			},

			"work_space_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Workspace SerialId.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudOceanusClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_oceanus_cluster.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_ids"); ok {
		clusterIdsSet := v.(*schema.Set).List()
		paramMap["ClusterIds"] = helper.InterfacesStringsPoint(clusterIdsSet)
	}

	if v, _ := d.GetOk("order_type"); v != nil {
		paramMap["OrderType"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*oceanus.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := oceanus.Filter{}
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
		paramMap["filters"] = tmpSet
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		paramMap["WorkSpaceId"] = helper.String(v.(string))
	}

	service := OceanusService{client: meta.(*TencentCloudClient).apiV3Conn}

	var clusterSet []*oceanus.Cluster

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusClusterByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		clusterSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(clusterSet))
	tmpList := make([]map[string]interface{}, 0, len(clusterSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
