package cdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudMysqlRollbackStop() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_rollback_stop.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	asyncRequestId, err := service.DeleteMysqlRollbackById(ctx, instanceId)
	if err != nil {
		return err
	}

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
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
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_rollback_stop.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMysqlRollbackStopDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_rollback_Stop.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
