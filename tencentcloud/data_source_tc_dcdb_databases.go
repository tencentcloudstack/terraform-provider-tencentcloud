/*
Use this data source to query detailed information of dcdb databases

Example Usage

```hcl
data "tencentcloud_dcdb_databases" "databases" {
  instance_id = "your_dcdb_instance_id"
  }
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDcdbDatabases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbDatabasesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id.",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Database information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database Name.",
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

func dataSourceTencentCloudDcdbDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dcdb_databases.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	dcdbService := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var dbs []*dcdb.Database
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := dcdbService.DescribeDcdbDatabasesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		dbs = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Dcdb list failed, reason:%+v", logId, err)
		return err
	}

	ids := make([]string, 0, len(dbs))
	list := make([]map[string]interface{}, 0, len(dbs))
	if dbs != nil {
		for _, db := range dbs {
			dbMap := map[string]interface{}{}
			if db.DbName != nil {
				dbMap["db_name"] = db.DbName
			}
			ids = append(ids, *db.DbName)
			list = append(list, dbMap)
		}
		d.SetId(helper.DataResourceIdsHash(ids))
		_ = d.Set("list", list)
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), list); e != nil {
			return e
		}
	}

	return nil
}
