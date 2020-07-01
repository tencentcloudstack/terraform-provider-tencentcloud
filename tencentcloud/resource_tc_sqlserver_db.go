/*
Provides a SQL Server DB resource belongs to SQL Server instance.

Example Usage

```hcl
resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  name = "example"
  charset = "Chinese_PRC_BIN"
  remark = "test-remark"
}
```

Import

SQL Server DB can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_db.foo mssql-3cdq7kx5#db_name
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudSqlserverDB() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverDBCreate,
		Read:   resourceTencentCloudSqlserverDBRead,
		Update: resourceTencentCloudSqlserverDBUpdate,
		Delete: resourceTencentCloudSqlserverDBDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SQLServer instance ID which DB belongs to.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of SQL Server DB. The DataBase name must be unique and must be composed of numbers, letters and underlines, and the first one can not be underline.",
			},
			"charset": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Chinese_PRC_CI_AS",
				ForceNew:    true,
				Description: "Character set DB uses. Valid values: `Chinese_PRC_CI_AS`, `Chinese_PRC_CS_AS`, `Chinese_PRC_BIN`, `Chinese_Taiwan_Stroke_CI_AS`, `SQL_Latin1_General_CP1_CI_AS`, and `SQL_Latin1_General_CP1_CS_AS`. Default value is `Chinese_PRC_CI_AS`.",
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
				Description: "Database status, could be `creating`, `running`, `modifying` which means changing the remark, and `deleting`.",
			},
		},
	}
}

func resourceTencentCloudSqlserverDBCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_db.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Get("instance_id").(string)
	_, has, err := sqlserverService.DescribeSqlserverInstanceById(ctx, instanceId)
	if err != nil {
		return fmt.Errorf("[CRITAL]%s DescribeSqlserverInstanceById fail, reason:%s\n", logId, err)
	}
	if !has {
		return fmt.Errorf("[CRITAL]%s SQL Server instance %s dose not exist for DB creation", logId, instanceId)
	}

	dbName := d.Get("name").(string)
	charset := d.Get("charset").(string)
	remark := d.Get("remark").(string)

	if err := sqlserverService.CreateSqlserverDB(ctx, instanceId, dbName, charset, remark); err != nil {
		return err
	}

	d.SetId(instanceId + FILED_SP + dbName)
	return resourceTencentCloudSqlserverDBRead(d, meta)
}

func resourceTencentCloudSqlserverDBRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_db.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	id := d.Id()
	dbInfo, has, err := sqlserverService.DescribeDBDetailsById(ctx, id)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}
	idItem := strings.Split(id, FILED_SP)
	if len(idItem) < 2 {
		return fmt.Errorf("broken ID %s of SQL Server DB", id)
	}
	instanceId := idItem[0]
	dbName := idItem[1]
	_ = d.Set("instance_id", instanceId)
	_ = d.Set("name", dbName)
	_ = d.Set("charset", dbInfo.Charset)
	_ = d.Set("remark", dbInfo.Remark)
	_ = d.Set("create_time", dbInfo.CreateTime)
	_ = d.Set("status", SQLSERVER_DB_STATUS[*dbInfo.Status])

	return nil
}

func resourceTencentCloudSqlserverDBUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_db.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := d.Get("instance_id").(string)

	if d.HasChange("remark") {
		if err := sqlserverService.ModifySqlserverDBRemark(ctx, instanceId, d.Get("name").(string), d.Get("remark").(string)); err != nil {
			return err
		}
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
	_, has, err := sqlserverService.DescribeSqlserverInstanceById(ctx, instanceId)
	if err != nil {
		return fmt.Errorf("[CRITAL]%s DescribeSqlserverInstanceById when deleting SQL Server DB fail, reason:%s\n", logId, err)
	}
	if !has {
		return nil
	}
	id := d.Id()
	_, has, err = sqlserverService.DescribeDBDetailsById(ctx, id)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	return sqlserverService.DeleteSqlserverDB(ctx, instanceId, name)
}
