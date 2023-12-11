package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSqlserverInstanceHa() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverInstanceHaCreate,
		Read:   resourceTencentCloudSqlserverInstanceHaRead,
		Update: resourceTencentCloudSqlserverInstanceHaUpdate,
		Delete: resourceTencentCloudSqlserverInstanceHaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"wait_switch": {
				Optional:    true,
				Type:        schema.TypeInt,
				Default:     0,
				Description: "Switch the execution mode, 0-execute immediately, 1-execute within the time window, the default value is 0.",
			},
		},
	}
}

func resourceTencentCloudSqlserverInstanceHaCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_instance_ha.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverInstanceHaUpdate(d, meta)
}

func resourceTencentCloudSqlserverInstanceHaRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_instance_ha.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId = d.Id()
	)

	instanceHa, err := service.DescribeSqlserverInstanceHaById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instanceHa == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverInstanceHa` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instanceHa.InstanceId != nil {
		_ = d.Set("instance_id", instanceHa.InstanceId)
	}

	if instanceHa.WaitSwitch != nil {
		_ = d.Set("wait_switch", instanceHa.WaitSwitch)
	}

	return nil
}

func resourceTencentCloudSqlserverInstanceHaUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_instance_ha.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = sqlserver.NewSwitchCloudInstanceHARequest()
		instanceId = d.Id()
	)

	request.InstanceId = &instanceId

	if v, ok := d.GetOkExists("wait_switch"); ok {
		request.WaitSwitch = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().SwitchCloudInstanceHA(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver instanceHa failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverInstanceHaRead(d, meta)
}

func resourceTencentCloudSqlserverInstanceHaDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_instance_ha.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
