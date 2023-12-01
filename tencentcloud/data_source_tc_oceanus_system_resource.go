/*
Use this data source to query detailed information of oceanus system_resource

Example Usage

```hcl
data "tencentcloud_oceanus_system_resource" "example" {
  resource_ids = ["resource-abd503yt"]
  filters {
    name   = "Name"
    values = ["tf_example"]
  }
  cluster_id    = "cluster-n8yaia0p"
  flink_version = "Flink-1.11"
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

func dataSourceTencentCloudOceanusSystemResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusSystemResourceRead,
		Schema: map[string]*schema.Schema{
			"resource_ids": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Array of resource IDs to be queried.",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Query the resource configuration list. If not specified, return all job configuration lists under ResourceIds.N.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to be filtered.",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Filter values for the field.",
						},
					},
				},
			},
			"cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
			"flink_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Query built-in connectors for the corresponding Flink version.",
			},
			"resource_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Collection of resource details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource name.",
						},
						"resource_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource type. 1 indicates JAR package, which is currently the only supported value.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource remarks.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region to which the resource belongs.",
						},
						"latest_resource_config_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Latest version of the resource.",
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

func dataSourceTencentCloudOceanusSystemResourceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_oceanus_system_resource.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId              = getLogId(contextNil)
		ctx                = context.WithValue(context.TODO(), logIdKey, logId)
		service            = OceanusService{client: meta.(*TencentCloudClient).apiV3Conn}
		systemResourceList []*oceanus.SystemResourceItem
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("resource_ids"); ok {
		resourceIdsSet := v.(*schema.Set).List()
		paramMap["ResourceIds"] = helper.InterfacesStringsPoint(resourceIdsSet)
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

		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("flink_version"); ok {
		paramMap["FlinkVersion"] = helper.String(v.(string))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusSystemResourceByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		systemResourceList = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(systemResourceList))
	tmpList := make([]map[string]interface{}, 0, len(systemResourceList))
	if systemResourceList != nil {
		for _, systemResourceItem := range systemResourceList {
			systemResourceItemMap := map[string]interface{}{}

			if systemResourceItem.ResourceId != nil {
				systemResourceItemMap["resource_id"] = systemResourceItem.ResourceId
			}

			if systemResourceItem.Name != nil {
				systemResourceItemMap["name"] = systemResourceItem.Name
			}

			if systemResourceItem.ResourceType != nil {
				systemResourceItemMap["resource_type"] = systemResourceItem.ResourceType
			}

			if systemResourceItem.Remark != nil {
				systemResourceItemMap["remark"] = systemResourceItem.Remark
			}

			if systemResourceItem.Region != nil {
				systemResourceItemMap["region"] = systemResourceItem.Region
			}

			if systemResourceItem.LatestResourceConfigVersion != nil {
				systemResourceItemMap["latest_resource_config_version"] = systemResourceItem.LatestResourceConfigVersion
			}

			ids = append(ids, *systemResourceItem.ResourceId)
			tmpList = append(tmpList, systemResourceItemMap)
		}

		_ = d.Set("resource_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
