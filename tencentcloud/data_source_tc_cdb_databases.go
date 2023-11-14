/*
Use this data source to query detailed information of cdb databases

Example Usage

```hcl
data "tencentcloud_cdb_databases" "databases" {
  instance_id = "cdb-c1nl9rpv"
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  database_regexp = &lt;nil&gt;
  total_count = &lt;nil&gt;
  items = &lt;nil&gt;
  database_list {
		database_name = &lt;nil&gt;
		character_set = &lt;nil&gt;

  }
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

func dataSourceTencentCloudCdbDatabases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbDatabasesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Page offset.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The number of single requests, the default value is 20, the minimum value is 1, and the maximum value is 100.",
			},

			"database_regexp": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Regular expression to match database library names.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "The total number of instances that meet the query condition.",
			},

			"items": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Returned instance information.",
			},

			"database_list": {
				Type:        schema.TypeList,
				Description: "Database name and character set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_name": {
							Type:        schema.TypeString,
							Description: "The name of database.",
						},
						"character_set": {
							Type:        schema.TypeString,
							Description: "Character set type.",
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

func dataSourceTencentCloudCdbDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_databases.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("database_regexp"); ok {
		paramMap["DatabaseRegexp"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("items"); ok {
		itemsSet := v.(*schema.Set).List()
		paramMap["Items"] = helper.InterfacesStringsPoint(itemsSet)
	}

	if v, ok := d.GetOk("database_list"); ok {
		databaseListSet := v.([]interface{})
		tmpSet := make([]*cdb.DatabasesWithCharacterLists, 0, len(databaseListSet))

		for _, item := range databaseListSet {
			databasesWithCharacterLists := cdb.DatabasesWithCharacterLists{}
			databasesWithCharacterListsMap := item.(map[string]interface{})

			if v, ok := databasesWithCharacterListsMap["database_name"]; ok {
				databasesWithCharacterLists.DatabaseName = helper.String(v.(string))
			}
			if v, ok := databasesWithCharacterListsMap["character_set"]; ok {
				databasesWithCharacterLists.CharacterSet = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &databasesWithCharacterLists)
		}
		paramMap["database_list"] = tmpSet
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var items []*string

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbDatabasesByFilter(ctx, paramMap)
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
