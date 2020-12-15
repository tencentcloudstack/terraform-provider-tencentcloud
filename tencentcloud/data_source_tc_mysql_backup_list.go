/*
Use this data source to query the list of backup databases.

Example Usage

```hcl
data "tencentcloud_mysql_backup_list" "default" {
  mysql_id           = "my-test-database"
  max_number         = 10
  result_output_file = "mytestpath"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentMysqlBackupList() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceTencentMysqlBackupListRead,
		Schema: map[string]*schema.Schema{
			"mysql_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance ID, such as `cdb-c1nl9rpv`. It is identical to the instance ID displayed in the database console page.",
			},
			"max_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validateIntegerInRange(1, 10000),
				Description:  "The latest files to list, rang from 1 to 10000. And the default value is `10`.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of MySQL backup. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The earliest time at which the backup starts. For example, `2` indicates 2:00 am.",
						},
						"finish_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time at which the backup finishes.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the size of backup file.",
						},
						"backup_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of Backup task.",
						},
						"backup_model": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup method. Supported values include: `physical` - physical backup, and `logical` - logical backup.",
						},
						"intranet_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL for downloads internally.",
						},
						"internet_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL for downloads externally.",
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The owner of the backup files.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentMysqlBackupListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mysql_backup_list.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	max_number, _ := d.Get("max_number").(int)
	backInfoItems, err := mysqlService.DescribeBackupsByMysqlId(ctx, d.Get("mysql_id").(string), int64(max_number))

	if err != nil {
		return fmt.Errorf("api[DescribeBackups]fail, return %s", err.Error())
	}

	var itemShemas []map[string]interface{}
	var ids = make([]string, len(backInfoItems))

	for index, item := range backInfoItems {
		mapping := map[string]interface{}{
			"time":         *item.Date,
			"finish_time":  *item.FinishTime,
			"size":         *item.Size,
			"backup_id":    *item.BackupId,
			"backup_model": *item.Type,
			"intranet_url": strings.Replace(*item.IntranetUrl, "\u0026", "&", -1),
			"internet_url": strings.Replace(*item.InternetUrl, "\u0026", "&", -1),
			"creator":      *item.Creator,
		}
		ids[index] = fmt.Sprintf("%d", *item.BackupId)
		itemShemas = append(itemShemas, mapping)
	}

	if err := d.Set("list", itemShemas); err != nil {
		log.Printf("[CRITAL]%s provider set itemShemas fail, reason:%s\n ", logId, err.Error())
	}
	d.SetId(helper.DataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {

		if err := writeToFile(output.(string), itemShemas); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail,  reason[%s]\n",
				logId, output.(string), err.Error())
		}

	}
	return nil
}
