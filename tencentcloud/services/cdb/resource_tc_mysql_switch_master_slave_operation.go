package cdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlSwitchMasterSlaveOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlSwitchMasterSlaveOperationCreate,
		Read:   resourceTencentCloudMysqlSwitchMasterSlaveOperationRead,
		Delete: resourceTencentCloudMysqlSwitchMasterSlaveOperationDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"dst_slave": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "target instance. Possible values: `first` - first standby; `second` - second standby. The default value is `first`, and only multi-AZ instances support setting it to `second`.",
			},

			"force_switch": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to force switch. Default is False. Note that if you set the mandatory switch to True, there is a risk of data loss on the instance, so use it with caution.",
			},

			"wait_switch": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to switch within the time window. The default is False, i.e. do not switch within the time window. Note that if the ForceSwitch parameter is set to True, this parameter will not take effect.",
			},
		},
	}
}

func resourceTencentCloudMysqlSwitchMasterSlaveOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_switch_master_slave_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request    = mysql.NewSwitchDBInstanceMasterSlaveRequest()
		response   = mysql.NewSwitchDBInstanceMasterSlaveResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_slave"); ok {
		request.DstSlave = helper.String(v.(string))
	}

	if v, _ := d.GetOk("force_switch"); v != nil {
		request.ForceSwitch = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOk("wait_switch"); v != nil {
		request.WaitSwitch = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().SwitchDBInstanceMasterSlave(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mysql switchMasterSlaveOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	asyncRequestId := *response.Response.AsyncRequestId
	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("%s operate mysql switchMasterSlaveOperation status is %s", instanceId, taskStatus))
		}
		err = fmt.Errorf("%s operate mysql switchMasterSlaveOperation status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mysql switchMasterSlaveOperation fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudMysqlSwitchMasterSlaveOperationRead(d, meta)
}

func resourceTencentCloudMysqlSwitchMasterSlaveOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_switch_master_slave_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMysqlSwitchMasterSlaveOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_switch_master_slave_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
