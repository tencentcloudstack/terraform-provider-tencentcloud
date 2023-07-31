/*
Provides a resource to create a sqlserver database_tde

Example Usage

Open database tde encryption

```hcl
resource "tencentcloud_sqlserver_database_tde" "example" {
  instance_id = "mssql-qelbzgwf"
  db_names    = ["example_db1", "example_db2"]
  encryption  = "enable"
}
```

Close database tde encryption

```hcl
resource "tencentcloud_sqlserver_database_tde" "example" {
  instance_id = "mssql-qelbzgwf"
  db_names    = ["example_db1", "example_db2"]
  encryption  = "disable"
}
```

Import

sqlserver database_tde can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_database_tde.database_tde database_tde_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSqlserverDatabaseTDE() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverDatabaseTDECreate,
		Read:   resourceTencentCloudSqlserverDatabaseTDERead,
		Update: resourceTencentCloudSqlserverDatabaseTDEUpdate,
		Delete: resourceTencentCloudSqlserverDatabaseTDEDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"db_names": {
				Required:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Database name list.",
			},
			"encryption": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "`enable` - enable encryption, `disable` - disable encryption.",
			},
		},
	}
}

func resourceTencentCloudSqlserverDatabaseTDECreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_database_tde.create")()
	defer inconsistentCheck(d, meta)()

	var (
		instanceId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	dbNameList := make([]string, 0)
	if v, ok := d.GetOk("db_names"); ok {
		dbNames := v.(*schema.Set).List()
		for i := range dbNames {
			dbName := dbNames[i].(string)
			dbNameList = append(dbNameList, dbName)
		}
	}

	dbNameListStr := strings.Join(dbNameList, COMMA_SP)
	d.SetId(strings.Join([]string{instanceId, dbNameListStr}, FILED_SP))

	return resourceTencentCloudSqlserverDatabaseTDEUpdate(d, meta)
}

func resourceTencentCloudSqlserverDatabaseTDERead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_database_tde.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		encryption string
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]
	dbNameListStr := idSplit[1]
	dbNameList := strings.Split(dbNameListStr, COMMA_SP)

	databaseTDE, err := service.DescribeSqlserverDatabaseTDEById(ctx, instanceId, dbNameList)
	if err != nil {
		return err
	}

	if databaseTDE == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverDatabaseTDE` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if databaseTDE.InstanceId != nil {
		_ = d.Set("instance_id", databaseTDE.InstanceId)
	}

	if databaseTDE.DBDetails != nil {
		tmpList := make([]string, 0)
		checkEncryption := make(map[string]string, 0)
		for _, item := range databaseTDE.DBDetails {
			if item.Name != nil {
				tmpList = append(tmpList, *item.Name)
			}

			if item.Encryption != nil {
				encryption = *item.Encryption
				checkEncryption[encryption] = ""
			}
		}

		if len(checkEncryption) != 1 {
			return fmt.Errorf("sqlserver database tde encryption result is not normal, id is %s", d.Id())
		}

		_ = d.Set("db_names", tmpList)
		_ = d.Set("encryption", encryption)
	}

	return nil
}

func resourceTencentCloudSqlserverDatabaseTDEUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_database_tde.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		request    = sqlserver.NewModifyDBEncryptAttributesRequest()
		encryption string
		flowId     int64
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]
	dbNameListStr := idSplit[1]
	dbNameList := strings.Split(dbNameListStr, COMMA_SP)

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("encryption"); ok {
		encryption = v.(string)
	}

	for _, v := range dbNameList {
		parameter := sqlserver.DBTDEEncrypt{}
		parameter.DBName = helper.String(v)
		parameter.Encryption = helper.String(encryption)
		request.DBTDEEncrypt = append(request.DBTDEEncrypt, &parameter)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyDBEncryptAttributes(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("qlserver databaseTDE not exists")
			return resource.NonRetryableError(e)
		}

		flowId = *result.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver databaseTDE failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCloneStatusByFlowId(ctx, flowId)
		if e != nil {
			return retryError(e)
		}

		if result == nil {
			e = fmt.Errorf("sqlserver databaseTDE instanceId %s flowId %d not exists", instanceId, flowId)
			return resource.NonRetryableError(e)
		}

		if *result.Status == SQLSERVER_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("sqlserver databaseTDE task status is running"))
		}

		if *result.Status == SQLSERVER_TASK_SUCCESS {
			return nil
		}

		if *result.Status == SQLSERVER_TASK_FAIL {
			return resource.NonRetryableError(fmt.Errorf("sqlserver databaseTDE task status is failed"))
		}

		e = fmt.Errorf("sqlserver databaseTDE task status is %v, we won't wait for it finish", *result.Status)
		return resource.NonRetryableError(e)
	})

	if err != nil {
		log.Printf("[CRITAL]%s sqlserver databaseTDE task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudSqlserverDatabaseTDERead(d, meta)
}

func resourceTencentCloudSqlserverDatabaseTDEDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_database_tde.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
