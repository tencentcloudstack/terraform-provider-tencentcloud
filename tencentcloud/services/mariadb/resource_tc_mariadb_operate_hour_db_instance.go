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
)

func ResourceTencentCloudMariadbOperateHourDbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbActivateHourDbInstanceCreate,
		Read:   resourceTencentCloudMariadbActivateHourDbInstanceRead,
		Update: resourceTencentCloudMariadbActivateHourDbInstanceUpdate,
		Delete: resourceTencentCloudMariadbActivateHourDbInstanceDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"operate": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Operation, `activate`- activate the hour db instance, `isolate`- isolate the hour db instance.",
			},
		},
	}
}

func resourceTencentCloudMariadbActivateHourDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_activate_hour_db_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMariadbActivateHourDbInstanceUpdate(d, meta)
}

func resourceTencentCloudMariadbActivateHourDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_activate_hour_db_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
		operate    string
	)

	if v, ok := d.GetOk("operate"); ok {
		operate = v.(string)
	}

	err := resource.Retry(10*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDBInstanceDetailById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if operate == "activate" {
			if *result.Status == MARIADB_STATUS_RUNNING {
				return nil
			} else if *result.Status == MARIADB_STATUS_FLOW {
				return resource.RetryableError(fmt.Errorf("mariadb accountPrivileges status is flow"))
			} else {
				e = fmt.Errorf("mariadb accountPrivileges status illegal")
				return resource.NonRetryableError(e)
			}
		} else if operate == "isolate" {
			if *result.Status == MARIADB_STATUS_ISOLATE {
				return nil
			} else if *result.Status == MARIADB_STATUS_FLOW {
				return resource.RetryableError(fmt.Errorf("mariadb accountPrivileges status is flow"))
			} else {
				e = fmt.Errorf("mariadb accountPrivileges status illegal")
				return resource.NonRetryableError(e)
			}
		} else {
			e = fmt.Errorf("[CRITAL]%s operate type error, %s", logId, operate)
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s update mariadb accountPrivileges task failed, reason:%+v", logId, err)
		return err
	}

	_ = d.Set("operate", operate)

	return nil
}

func resourceTencentCloudMariadbActivateHourDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_activate_hour_db_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		instanceId = d.Id()
	)

	if v, ok := d.GetOk("operate"); ok {
		operate := v.(string)
		if operate == "activate" {
			request := mariadb.NewActivateHourDBInstanceRequest()
			response := mariadb.NewActivateHourDBInstanceResponse()
			request.InstanceIds = common.StringPtrs([]string{instanceId})
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMariadbClient().ActivateHourDBInstance(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				response = result
				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s operate mariadb activateHourDbInstance failed, reason:%+v", logId, err)
				return err
			}

			if response == nil {
				return fmt.Errorf("operate mariadb activateHourDbInstance not found")
			}

		} else if operate == "isolate" {
			request := mariadb.NewIsolateHourDBInstanceRequest()
			response := mariadb.NewIsolateHourDBInstanceResponse()
			request.InstanceIds = common.StringPtrs([]string{instanceId})
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMariadbClient().IsolateHourDBInstance(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				response = result
				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s operate mariadb isolateHourInstance failed, reason:%+v", logId, err)
				return err
			}

			if response == nil {
				return fmt.Errorf("operate mariadb isolateHourInstance not found")
			}

		} else {
			return fmt.Errorf("[CRITAL]%s operate type error, %s", logId, operate)
		}
	}

	return resourceTencentCloudMariadbActivateHourDbInstanceRead(d, meta)
}

func resourceTencentCloudMariadbActivateHourDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_activate_hour_db_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
