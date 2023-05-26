/*
Provides a resource to create a mysql restart_db_instances_operation

Example Usage

```hcl
resource "tencentcloud_mysql_restart_db_instances_operation" "restart_db_instances_operation" {
  instance_id = "cdb-bohspx3j"
}
```

Import

mysql restart_db_instances_operation can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_restart_db_instances_operation.restart_db_instances_operation restart_db_instances_operation_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func resourceTencentCloudMysqlRestartDbInstancesOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlRestartDbInstancesOperationCreate,
		Read:   resourceTencentCloudMysqlRestartDbInstancesOperationRead,
		Delete: resourceTencentCloudMysqlRestartDbInstancesOperationDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "An array of instance ID in the format: cdb-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Instance status.",
			},
		},
	}
}

func resourceTencentCloudMysqlRestartDbInstancesOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_restart_db_instances_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = mysql.NewRestartDBInstancesRequest()
		response   = mysql.NewRestartDBInstancesResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceIds = []*string{&instanceId}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().RestartDBInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mysql restartDbInstancesOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	asyncRequestId := *response.Response.AsyncRequestId
	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("%s operate mysql restartDbInstancesOperation status is %s", instanceId, taskStatus))
		}
		err = fmt.Errorf("%s operate mysql restartDbInstancesOperation is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mysql restartDbInstancesOperation fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudMysqlRestartDbInstancesOperationRead(d, meta)
}

func resourceTencentCloudMysqlRestartDbInstancesOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_restart_db_instances_operation.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	isolateInstance, err := service.DescribeDBInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if isolateInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `restartDbInstancesOperation` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if isolateInstance.InstanceId != nil {
		_ = d.Set("instance_id", isolateInstance.InstanceId)
	}

	if isolateInstance.Status != nil {
		_ = d.Set("status", isolateInstance.Status)
	}

	return nil
}

func resourceTencentCloudMysqlRestartDbInstancesOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_restart_db_instances_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
