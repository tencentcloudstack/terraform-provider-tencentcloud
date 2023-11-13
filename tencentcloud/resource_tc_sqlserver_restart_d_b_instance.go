/*
Provides a resource to create a sqlserver restart_d_b_instance

Example Usage

```hcl
resource "tencentcloud_sqlserver_restart_d_b_instance" "restart_d_b_instance" {
  instance_id = "mssql-i1z41iwd"
}
```

Import

sqlserver restart_d_b_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_restart_d_b_instance.restart_d_b_instance restart_d_b_instance_id
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

func resourceTencentCloudSqlserverRestartDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverRestartDBInstanceCreate,
		Read:   resourceTencentCloudSqlserverRestartDBInstanceRead,
		Update: resourceTencentCloudSqlserverRestartDBInstanceUpdate,
		Delete: resourceTencentCloudSqlserverRestartDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
		},
	}
}

func resourceTencentCloudSqlserverRestartDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_restart_d_b_instance.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverRestartDBInstanceUpdate(d, meta)
}

func resourceTencentCloudSqlserverRestartDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_restart_d_b_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	restartDBInstanceId := d.Id()

	restartDBInstance, err := service.DescribeSqlserverRestartDBInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if restartDBInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverRestartDBInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if restartDBInstance.InstanceId != nil {
		_ = d.Set("instance_id", restartDBInstance.InstanceId)
	}

	return nil
}

func resourceTencentCloudSqlserverRestartDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_restart_d_b_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewRestartDBInstanceRequest()

	restartDBInstanceId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().RestartDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver restartDBInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverRestartDBInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverRestartDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_restart_d_b_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
