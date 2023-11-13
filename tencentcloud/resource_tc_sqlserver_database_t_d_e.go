/*
Provides a resource to create a sqlserver database_t_d_e

Example Usage

```hcl
resource "tencentcloud_sqlserver_database_t_d_e" "database_t_d_e" {
  instance_id = "mssql-i1z41iwd"
  d_b_t_d_e_encrypt {
		d_b_name = ""
		encryption = ""

  }
}
```

Import

sqlserver database_t_d_e can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_database_t_d_e.database_t_d_e database_t_d_e_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"log"
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

			"d_b_t_d_e_encrypt": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Enable and disable database TDE encryption.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"d_b_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database.",
						},
						"encryption": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable - enable encryption, disable - disable encryption.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSqlserverDatabaseTDECreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_database_t_d_e.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverDatabaseTDEUpdate(d, meta)
}

func resourceTencentCloudSqlserverDatabaseTDERead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_database_t_d_e.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	databaseTDEId := d.Id()

	databaseTDE, err := service.DescribeSqlserverDatabaseTDEById(ctx, instanceId)
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

	if databaseTDE.DBTDEEncrypt != nil {
		dBTDEEncryptList := []interface{}{}
		for _, dBTDEEncrypt := range databaseTDE.DBTDEEncrypt {
			dBTDEEncryptMap := map[string]interface{}{}

			if databaseTDE.DBTDEEncrypt.DBName != nil {
				dBTDEEncryptMap["d_b_name"] = databaseTDE.DBTDEEncrypt.DBName
			}

			if databaseTDE.DBTDEEncrypt.Encryption != nil {
				dBTDEEncryptMap["encryption"] = databaseTDE.DBTDEEncrypt.Encryption
			}

			dBTDEEncryptList = append(dBTDEEncryptList, dBTDEEncryptMap)
		}

		_ = d.Set("d_b_t_d_e_encrypt", dBTDEEncryptList)

	}

	return nil
}

func resourceTencentCloudSqlserverDatabaseTDEUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_database_t_d_e.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewModifyDBEncryptAttributesRequest()

	databaseTDEId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "d_b_t_d_e_encrypt"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyDBEncryptAttributes(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver databaseTDE failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverDatabaseTDERead(d, meta)
}

func resourceTencentCloudSqlserverDatabaseTDEDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_database_t_d_e.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
