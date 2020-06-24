/*
Use this data source to query DB resources for the specific SQLServer instance.

Example Usage

```hcl
resource "tencentcloud_sqlserver_db" "mysqlserver_db" {
  instance_id = "mssql-XXXXXX"
  name = "sqlserver_db_terraform"
  charset = "Chinese_PRC_BIN"
  remark = "test-remark"
}

data "tencentcloud_sqlserver_db" "mysqlserver" {
  instance_id = tencentcloud_sqlserver_db.mysqlserver_db.instance_id
  name        = tencentcloud_sqlserver_db.mysqlserver_db.name
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceTencentSqlserverDB() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentSqlserverDBRead,
		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SQLServer instance ID which DB belongs to.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of DB.",
			},
			// Computed
			"charset": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Character set DB uses, could be `Chinese_PRC_CI_AS`, `Chinese_PRC_CS_AS`, `Chinese_PRC_BIN`, `Chinese_Taiwan_Stroke_CI_AS`, `SQL_Latin1_General_CP1_CI_AS`, and `SQL_Latin1_General_CP1_CS_AS`. Default value is `Chinese_PRC_CI_AS`.",
			},
			"remark": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Remark of the DB.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Database creation time.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Database status. Valid values are `creating`, `running`, `modifying`, `dropping`.",
			},
		},
	}
}

func dataSourceTencentSqlserverDBRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencent_sqlserver_db.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	// precheck
	instanceId := d.Get("instance_id").(string)
	_, has, err := sqlserverService.DescribeSqlserverInstanceById(ctx, instanceId)
	if err != nil {
		return fmt.Errorf("[CRITAL]%s DescribeSqlserverInstanceById fail, reason:%s\n", logId, err)
	}
	if !has {
		return fmt.Errorf("[CRITAL]%s SQLServer instance %s dose not exist", logId, instanceId)
	}

	dbName := d.Get("name").(string)
	id := instanceId + FILED_SP + dbName
	dbInfo, has, err := sqlserverService.DescribeDBDetailsById(ctx, id)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("[CRITAL]%s DB %s doesn't exist for SQLServer instance %s", logId, dbName, instanceId)
	}
	_ = d.Set("instance_id", instanceId)
	_ = d.Set("name", dbName)
	_ = d.Set("charset", dbInfo.Charset)
	_ = d.Set("remark", dbInfo.Remark)
	_ = d.Set("create_time", dbInfo.CreateTime)
	_ = d.Set("status", SQLSERVER_DB_STATUS[*dbInfo.Status])
	d.SetId(id)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), id); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
		}
	}
	return nil
}
