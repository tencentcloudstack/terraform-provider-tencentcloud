/*
Use this data source to query detailed information of postgres read_only_group

Example Usage

```hcl
data "tencentcloud_postgres_read_only_group" "read_only_group" {
  filters {
		name = "db-master-instance-id"
		values =

  }
  page_size = 10
  page_number = 0
  order_by = "CreateTime"
  order_by_type = "asc"
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudPostgresReadOnlyGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresReadOnlyGroupRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter condition. The primary ID must be specified in the format of db-master-instance-id to filter results, or else null will be returned.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter name.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "One or more filter values.",
						},
					},
				},
			},

			"page_size": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The number of results per page. Default value:10.",
			},

			"page_number": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Specify which page is displayed. Default value:1 (the first page).",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting criterion. Valid values:ROGroupId, CreateTime, Name.",
			},

			"order_by_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting order. Valid values:desc, asc.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudPostgresReadOnlyGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_postgres_read_only_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*postgres.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := postgres.Filter{}
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

	if v, _ := d.GetOk("page_size"); v != nil {
		paramMap["PageSize"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("page_number"); v != nil {
		paramMap["PageNumber"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["OrderByType"] = helper.String(v.(string))
	}

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	var readOnlyGroupList []*postgres.ReadOnlyGroup

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresReadOnlyGroupByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		readOnlyGroupList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(readOnlyGroupList))
	tmpList := make([]map[string]interface{}, 0, len(readOnlyGroupList))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
