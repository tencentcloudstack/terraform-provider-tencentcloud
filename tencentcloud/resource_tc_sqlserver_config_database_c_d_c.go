/*
Provides a resource to create a sqlserver config_database_c_d_c

Example Usage

```hcl
resource "tencentcloud_sqlserver_config_database_c_d_c" "config_database_c_d_c" {
  d_b_names =
  modify_type = "enable"
  instance_id = "mssql-i1z41iwd"
}
```

Import

sqlserver config_database_c_d_c can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_database_c_d_c.config_database_c_d_c config_database_c_d_c_id
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

func resourceTencentCloudSqlserverConfigDatabaseCDC() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverConfigDatabaseCDCCreate,
		Read:   resourceTencentCloudSqlserverConfigDatabaseCDCRead,
		Update: resourceTencentCloudSqlserverConfigDatabaseCDCUpdate,
		Delete: resourceTencentCloudSqlserverConfigDatabaseCDCDelete,
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
				Description: "Enable or disable CDC. Valid values: enable, disable.",
			},

			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
		},
	}
}

func resourceTencentCloudSqlserverConfigDatabaseCDCCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_database_c_d_c.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverConfigDatabaseCDCUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigDatabaseCDCRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_database_c_d_c.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	configDatabaseCDCId := d.Id()

	configDatabaseCDC, err := service.DescribeSqlserverConfigDatabaseCDCById(ctx, instanceId)
	if err != nil {
		return err
	}

	if configDatabaseCDC == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverConfigDatabaseCDC` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if configDatabaseCDC.DBNames != nil {
		_ = d.Set("d_b_names", configDatabaseCDC.DBNames)
	}

	if configDatabaseCDC.ModifyType != nil {
		_ = d.Set("modify_type", configDatabaseCDC.ModifyType)
	}

	if configDatabaseCDC.InstanceId != nil {
		_ = d.Set("instance_id", configDatabaseCDC.InstanceId)
	}

	return nil
}

func resourceTencentCloudSqlserverConfigDatabaseCDCUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_database_c_d_c.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewModifyDatabaseCDCRequest()

	configDatabaseCDCId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"d_b_names", "modify_type", "instance_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyDatabaseCDC(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver configDatabaseCDC failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverConfigDatabaseCDCRead(d, meta)
}

func resourceTencentCloudSqlserverConfigDatabaseCDCDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_database_c_d_c.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
