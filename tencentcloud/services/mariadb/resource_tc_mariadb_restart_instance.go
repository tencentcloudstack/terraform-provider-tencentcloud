package mariadb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMariadbRestartInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbRestartInstanceCreate,
		Read:   resourceTencentCloudMariadbRestartInstanceRead,
		Delete: resourceTencentCloudMariadbRestartInstanceDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "instance ID.",
			},
			"restart_time": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "expected restart time.",
			},
		},
	}
}

func resourceTencentCloudMariadbRestartInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_restart_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request    = mariadb.NewRestartDBInstancesRequest()
		instanceId string
		flowId     int64
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceIds = common.StringPtrs([]string{v.(string)})
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("restart_time"); ok {
		request.RestartTime = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMariadbClient().RestartDBInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		flowId = *result.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb restartInstance failed, reason:%+v", logId, err)
		return err
	}

	// wait
	err = resource.Retry(10*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeFlowById(ctx, flowId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if *result.Status == MARIADB_TASK_SUCCESS {
			return nil
		} else if *result.Status == MARIADB_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("operate mariadb restartInstance status is running"))
		} else if *result.Status == MARIADB_TASK_FAIL {
			return resource.NonRetryableError(fmt.Errorf("operate mariadb restartInstance status is fail"))
		} else {
			e = fmt.Errorf("operate mariadb restartInstance status illegal")
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb restartInstance task failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudMariadbRestartInstanceRead(d, meta)
}

func resourceTencentCloudMariadbRestartInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_restart_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMariadbRestartInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_restart_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
