package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudMysqlRollbackStop() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlRollbackStopCreate,
		Read:   resourceTencentCloudMysqlRollbackStopRead,
		Delete: resourceTencentCloudMysqlRollbackStopDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cloud database instance ID.",
			},
		},
	}
}

func resourceTencentCloudMysqlRollbackStopCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_rollback_stop.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	asyncRequestId, err := service.DeleteMysqlRollbackById(ctx, instanceId)
	if err != nil {
		return err
	}

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("%s delete mysql rollback status is %s", instanceId, taskStatus))
		}
		err = fmt.Errorf("%s delete mysql rollback status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete mysql rollback fail, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudMysqlRollbackStopRead(d, meta)
}

func resourceTencentCloudMysqlRollbackStopRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_rollback_stop.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMysqlRollbackStopDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_rollback_Stop.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
