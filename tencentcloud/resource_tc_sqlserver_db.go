/*
Provides a SQLServer DB resource to database belongs to SQLServer instance.

Example Usage

```hcl
resource "tencentcloud_sqlserver_db" "mysqlserver_db" {
  instance_id = "mssql-XXXXXX"
  name = "sqlserver_db_terraform"
  charset = "Chinese_PRC_BIN"
  remark = "test-remark"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudSqlserverDB() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverDBCreate,
		Read:   resourceTencentCloudSqlserverDBRead,
		Update: resourceTencentCloudSqlserverDBUpdate,
		Delete: resourceTencentCloudSqlserverDBDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SQL server instance ID which DB belongs to.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of DB. The DataBase name must be unique and must be composed of numbers, letters and underlines, and the first one can not be underline.",
			},
			"charset": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Chinese_PRC_CI_AS",
				ForceNew:     true,
				Description:  "Character set DB uses. Valid values: `Chinese_PRC_CI_AS`, `Chinese_PRC_CS_AS`, `Chinese_PRC_BIN`, `Chinese_Taiwan_Stroke_CI_AS`, `SQL_Latin1_General_CP1_CI_AS`, and `SQL_Latin1_General_CP1_CS_AS`. Default value is `Chinese_PRC_CI_AS`.",
				ValidateFunc: validateAllowedStringValue(SQLSERVER_CHARSET_LIST),
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remark of the DB.",
			},
			// Computed
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

func resourceTencentCloudSqlserverDBCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_db.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceID := d.Get("instance_id").(string)
	_, has, err := sqlserverService.DescribeSqlserverInstanceById(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("[CRITAL]%s DescribeSqlserverInstanceById fail, reason:%s\n", logId, err)
	}
	if !has {
		return fmt.Errorf("[CRITAL]%s Sqlserver instance %s dose not exist for DB creation", logId, instanceID)
	}

	dbName := d.Get("name").(string)
	charset := d.Get("charset").(string)
	remark := d.Get("remark").(string)

	if err := sqlserverService.CreateSqlserverDB(ctx, instanceID, dbName, charset, remark); err != nil {
		return err
	}

	return resourceTencentCloudSqlserverDBRead(d, meta)
}

func resourceTencentCloudSqlserverDBRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_db.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := d.Get("instance_id").(string)
	name := d.Get("name").(string)
	dbInfo, has, err := sqlserverService.DescribeDBDetailsByName(ctx, instanceId, name)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}
	_ = d.Set("charset", *dbInfo.Charset)
	_ = d.Set("remark", *dbInfo.Remark)
	_ = d.Set("create_time", *dbInfo.CreateTime)
	_ = d.Set("status", SQLSERVER_DB_STATUS[*dbInfo.Status])

	d.SetId(name)
	return nil
}

func resourceTencentCloudSqlserverDBUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_db.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := d.Get("instance_id").(string)

	if d.HasChange("name") {
		oldValue, newValue := d.GetChange("name")
		if err := sqlserverService.ModifySqlserverDBName(ctx, instanceId, oldValue.(string), newValue.(string)); err != nil {
			return err
		}
		d.SetPartial("name")
	}

	if d.HasChange("remark") {
		if err := sqlserverService.ModifySqlserverDBRemark(ctx, instanceId, d.Get("name").(string), d.Get("remark").(string)); err != nil {
			return err
		}
		d.SetPartial("remark")
	}

	return nil
}

func resourceTencentCloudSqlserverDBDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_db.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := d.Get("instance_id").(string)
	name := d.Get("name").(string)

	// precheck before delete
	_, has, err := sqlserverService.DescribeDBDetailsByName(ctx, instanceId, name)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	return sqlserverService.DeleteSqlserverDB(ctx, instanceId, name)
}
