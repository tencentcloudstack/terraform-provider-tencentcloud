/*
Use this data source to query detailed information of cdb tables

Example Usage

```hcl
data "tencentcloud_cdb_tables" "tables" {
  instance_id = "cdb-c1nl9rpv"
  database = &lt;nil&gt;
  offset = 0
  limit = 20
  table_regexp = &lt;nil&gt;
  total_count = &lt;nil&gt;
  items = &lt;nil&gt;
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCdbTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbTablesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"database": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The name of database.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Record offset, the default value is 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The number returned by a single request, the default value is 20, and the maximum value is 2000.",
			},

			"table_regexp": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Match the regular expression of the database table name, the rules are the same as MySQL official website.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "The total number of database tables that meet the query condition.",
			},

			"items": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Returned database table information.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCdbTablesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_tables.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("database"); ok {
		paramMap["Database"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("table_regexp"); ok {
		paramMap["TableRegexp"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("items"); ok {
		itemsSet := v.(*schema.Set).List()
		paramMap["Items"] = helper.InterfacesStringsPoint(itemsSet)
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var items []*string

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbTablesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(items))
	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
