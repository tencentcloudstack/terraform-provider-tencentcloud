/*
Provides a resource to create a sqlserver config_database_c_t

Example Usage

```hcl
resource "tencentcloud_sqlserver_config_database_c_t" "config_database_c_t" {
  d_b_names =
  modify_type = "enable"
  instance_id = "mssql-i1z41iwd"
  change_retention_day = 7
}
```

Import

sqlserver config_database_c_t can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_database_c_t.config_database_c_t config_database_c_t_id
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

func resourceTencentCloudSqlserverConfigDatabaseCT() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverConfigDatabaseCTCreate,
		Read:   resourceTencentCloudSqlserverConfigDatabaseCTRead,
		Update: resourceTencentCloudSqlserverConfigDatabaseCTUpdate,
		Delete: resourceTencentCloudSqlserverConfigDatabaseCTDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"d_b_names": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Array of database names.",
			},

			"modify_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Enable or disable CT. Valid values: enable, disable.",
			},

			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"change_retention_day": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Retention period (in days) of change tracking information when CT is enabled. Value range: 3-30. Default value: 3.",
			},
		},
	}
}

func resourceTencentCloudSqlserverConfigDatabaseCTCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_database_c_t.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverConfigDatabaseCTUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigDatabaseCTRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_database_c_t.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	configDatabaseCTId := d.Id()

	configDatabaseCT, err := service.DescribeSqlserverConfigDatabaseCTById(ctx, instanceId)
	if err != nil {
		return err
	}

	if configDatabaseCT == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverConfigDatabaseCT` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if configDatabaseCT.DBNames != nil {
		_ = d.Set("d_b_names", configDatabaseCT.DBNames)
	}

	if configDatabaseCT.ModifyType != nil {
		_ = d.Set("modify_type", configDatabaseCT.ModifyType)
	}

	if configDatabaseCT.InstanceId != nil {
		_ = d.Set("instance_id", configDatabaseCT.InstanceId)
	}

	if configDatabaseCT.ChangeRetentionDay != nil {
		_ = d.Set("change_retention_day", configDatabaseCT.ChangeRetentionDay)
	}

	return nil
}

func resourceTencentCloudSqlserverConfigDatabaseCTUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_database_c_t.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewModifyDatabaseCTRequest()

	configDatabaseCTId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"d_b_names", "modify_type", "instance_id", "change_retention_day"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyDatabaseCT(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver configDatabaseCT failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverConfigDatabaseCTRead(d, meta)
}

func resourceTencentCloudSqlserverConfigDatabaseCTDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_database_c_t.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
