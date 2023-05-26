/*
Provides a resource to create a sqlserver restart_db_instance

Example Usage

```hcl
resource "tencentcloud_sqlserver_restart_db_instance" "restart_db_instance" {
  instance_id = "mssql-i1z41iwd"
}
```

Import

sqlserver restart_db_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_restart_db_instance.restart_db_instance restart_db_instance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
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

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId = d.Id()
	)

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

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		request    = sqlserver.NewRestartDBInstanceRequest()
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId = d.Id()
		flowId     uint64
	)

	request.InstanceId = &instanceId
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().RestartDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("sqlserver business restartDBInstance not exists")
			return resource.NonRetryableError(e)
		}

		flowId = *result.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver restartDBInstance failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCloneStatusByFlowId(ctx, int64(flowId))
		if e != nil {
			return retryError(e)
		}

		if result == nil {
			e = fmt.Errorf("sqlserver restartDBInstance instanceId %s flowId %d not exists", instanceId, flowId)
			return resource.NonRetryableError(e)
		}

		if *result.Status == SQLSERVER_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("restart sqlserver restartDBInstance task status is running"))
		}

		if *result.Status == SQLSERVER_TASK_SUCCESS {
			return nil
		}

		if *result.Status == SQLSERVER_TASK_FAIL {
			return resource.NonRetryableError(fmt.Errorf("restart sqlserver restartDBInstance task status is failed"))
		}

		e = fmt.Errorf("restart sqlserver restartDBInstance task status is %v, we won't wait for it finish", *result.Status)
		return resource.NonRetryableError(e)
	})

	if err != nil {
		log.Printf("[CRITAL]%s restart sqlserver restartDBInstance task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudSqlserverRestartDBInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverRestartDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_restart_d_b_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
