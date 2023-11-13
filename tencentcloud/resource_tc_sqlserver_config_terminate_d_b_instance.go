/*
Provides a resource to create a sqlserver config_terminate_d_b_instance

Example Usage

```hcl
resource "tencentcloud_sqlserver_config_terminate_d_b_instance" "config_terminate_d_b_instance" {
  instance_id_set =
}
```

Import

sqlserver config_terminate_d_b_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_terminate_d_b_instance.config_terminate_d_b_instance config_terminate_d_b_instance_id
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

func resourceTencentCloudSqlserverConfigTerminateDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverConfigTerminateDBInstanceCreate,
		Read:   resourceTencentCloudSqlserverConfigTerminateDBInstanceRead,
		Update: resourceTencentCloudSqlserverConfigTerminateDBInstanceUpdate,
		Delete: resourceTencentCloudSqlserverConfigTerminateDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id_set": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of instance IDs to be actively destroyed, in the format: [mssql-3l3fgqn7]. Same as the instance ID displayed in the cloud database console page.",
			},
		},
	}
}

func resourceTencentCloudSqlserverConfigTerminateDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_terminate_d_b_instance.create")()
	defer inconsistentCheck(d, meta)()

	d.SetId()

	return resourceTencentCloudSqlserverConfigTerminateDBInstanceUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigTerminateDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_terminate_d_b_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	configTerminateDBInstanceId := d.Id()

	configTerminateDBInstance, err := service.DescribeSqlserverConfigTerminateDBInstanceById(ctx, instanceIdSet)
	if err != nil {
		return err
	}

	if configTerminateDBInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverConfigTerminateDBInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if configTerminateDBInstance.InstanceIdSet != nil {
		_ = d.Set("instance_id_set", configTerminateDBInstance.InstanceIdSet)
	}

	return nil
}

func resourceTencentCloudSqlserverConfigTerminateDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_terminate_d_b_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewTerminateDBInstanceRequest()

	configTerminateDBInstanceId := d.Id()

	request.InstanceIdSet = &instanceIdSet

	immutableArgs := []string{"instance_id_set"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().TerminateDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver configTerminateDBInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverConfigTerminateDBInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverConfigTerminateDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_terminate_d_b_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
