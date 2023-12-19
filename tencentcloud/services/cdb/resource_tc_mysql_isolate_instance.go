package cdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlIsolateInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlIsolateInstanceCreate,
		Read:   resourceTencentCloudMysqlIsolateInstanceRead,
		Update: resourceTencentCloudMysqlIsolateInstanceUpdate,
		Delete: resourceTencentCloudMysqlIsolateInstanceDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, the format is: cdb-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page, and you can use the [query instance list] (https://cloud.tencent.com/document/api/236/15872) interface Gets the value of the field InstanceId in the output parameter.",
			},

			"operate": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"isolate", "recover"}),
				Description:  "Manipulate instance, `isolate` - isolate instance, `recover`- recover isolated instance.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Instance status.",
			},
		},
	}
}

func resourceTencentCloudMysqlIsolateInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_isolate_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMysqlIsolateInstanceUpdate(d, meta)
}

func resourceTencentCloudMysqlIsolateInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_isolate_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	isolateInstance, err := service.DescribeDBInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if isolateInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlIsolateInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

func resourceTencentCloudMysqlIsolateInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_isolate_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()

	var operate string
	if v, ok := d.GetOk("operate"); ok {
		operate = v.(string)
	}

	if operate == "isolate" {
		request := mysql.NewIsolateDBInstanceRequest()
		request.InstanceId = helper.String(instanceId)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().IsolateDBInstance(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create mysql isolateInstance failed, reason:%+v", logId, err)
			return err
		}

		service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		err = resource.Retry(7*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			mysqlInfo, err := service.DescribeDBInstanceById(ctx, instanceId)

			if err != nil {
				if _, ok := err.(*errors.TencentCloudSDKError); !ok {
					return resource.RetryableError(err)
				} else {
					return resource.NonRetryableError(err)
				}
			}
			if mysqlInfo == nil {
				return nil
			}
			if *mysqlInfo.Status == MYSQL_STATUS_ISOLATING || *mysqlInfo.Status == MYSQL_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("mysql isolating."))
			}
			if *mysqlInfo.Status == MYSQL_STATUS_ISOLATED {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("after IsolateDBInstance mysql Status is %d", *mysqlInfo.Status))
		})

		if err != nil {
			log.Printf("[CRITAL]%s Isolate mysql isolateInstance fail, reason:%s\n ", logId, err.Error())
			return err
		}
	} else if operate == "recover" {
		service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		if err := service.DeleteMysqlIsolateInstanceById(ctx, instanceId); err != nil {
			return err
		}
		err := resource.Retry(7*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			mysqlInfo, err := service.DescribeDBInstanceById(ctx, instanceId)

			if err != nil {
				if _, ok := err.(*errors.TencentCloudSDKError); !ok {
					return resource.RetryableError(err)
				} else {
					return resource.NonRetryableError(err)
				}
			}
			if mysqlInfo == nil {
				return nil
			}
			if *mysqlInfo.Status == MYSQL_STATUS_ISOLATED {
				return resource.RetryableError(fmt.Errorf("mysql recovering."))
			}
			if *mysqlInfo.Status == MYSQL_STATUS_RUNNING {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("after ReleaseIsolatedDBInstances mysql Status is %d", *mysqlInfo.Status))
		})

		if err != nil {
			log.Printf("[CRITAL]%s ReleaseIsolatedDBInstances mysql fail, reason:%s\n ", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudMysqlIsolateInstanceRead(d, meta)
}

func resourceTencentCloudMysqlIsolateInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_isolate_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
