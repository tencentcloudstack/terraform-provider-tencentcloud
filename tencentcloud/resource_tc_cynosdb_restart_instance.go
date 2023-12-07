package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCynosdbRestartInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbRestartInstanceCreate,
		Read:   resourceTencentCloudCynosdbRestartInstanceRead,
		Delete: resourceTencentCloudCynosdbRestartInstanceDelete,
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

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "instance state.",
			},
		},
	}
}

func resourceTencentCloudCynosdbRestartInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_restart_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = cynosdb.NewRestartInstanceRequest()
		response   = cynosdb.NewRestartInstanceResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().RestartInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb restartInstance failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	flowId := *response.Response.FlowId
	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeFlow(ctx, flowId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("create cynosdb clusterPasswordComplexity is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb clusterPasswordComplexity fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudCynosdbRestartInstanceRead(d, meta)
}

func resourceTencentCloudCynosdbRestartInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_restart_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	_, instance, has, err := service.DescribeInstanceById(ctx, id)
	if err != nil {
		return err
	}
	if !has {
		log.Printf("[WARN]%s resource `DescribeInstanceById` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if instance.Status != nil {
		_ = d.Set("status", instance.Status)
	}

	return nil
}

func resourceTencentCloudCynosdbRestartInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_restart_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
