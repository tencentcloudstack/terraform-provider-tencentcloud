/*
Provides a resource to create a rum instance_status_config

Example Usage

```hcl
resource "tencentcloud_rum_instance_status_config" "instance_status_config" {
  instance_id = "rum-pasZKEI3RLgakj"
  operate     = "stop"
}
```

Import

rum instance_status_config can be imported using the id, e.g.

```
terraform import tencentcloud_rum_instance_status_config.instance_status_config instance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
)

func resourceTencentCloudRumInstanceStatusConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRumInstanceStatusConfigCreate,
		Read:   resourceTencentCloudRumInstanceStatusConfigRead,
		Update: resourceTencentCloudRumInstanceStatusConfigUpdate,
		Delete: resourceTencentCloudRumInstanceStatusConfigDelete,
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

			"instance_status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Instance status (`1`=creating, `2`=running, `3`=abnormal, `4`=restarting, `5`=stopping, `6`=stopped, `7`=deleted).",
			},

			"operate": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "`resume`, `stop`.",
			},
		},
	}
}

func resourceTencentCloudRumInstanceStatusConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_instance_status_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRumInstanceStatusConfigUpdate(d, meta)
}

func resourceTencentCloudRumInstanceStatusConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_instance_status_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()
	instanceStatusConfig, err := service.DescribeRumInstanceStatusConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instanceStatusConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RumInstanceStatusConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instanceStatusConfig.InstanceId != nil {
		_ = d.Set("instance_id", instanceStatusConfig.InstanceId)
	}

	if instanceStatusConfig.InstanceStatus != nil {
		_ = d.Set("instance_status", instanceStatusConfig.InstanceStatus)
	}

	return nil
}

func resourceTencentCloudRumInstanceStatusConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_instance_status_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	instanceId := d.Id()

	var operate string
	if v, ok := d.GetOk("operate"); ok {
		operate = v.(string)
	}

	if operate == "resume" {
		request := rum.NewResumeInstanceRequest()
		request.InstanceId = &instanceId
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseRumClient().ResumeInstance(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s resume rum instance failed, reason:%+v", logId, err)
			return err
		}
	} else if operate == "stop" {
		request := rum.NewStopInstanceRequest()
		request.InstanceId = &instanceId
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseRumClient().StopInstance(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s stop rum instance failed, reason:%+v", logId, err)
			return err
		}
	} else {
		return fmt.Errorf("[CRITAL]%s Operation type error", logId)
	}

	return resourceTencentCloudRumInstanceStatusConfigRead(d, meta)
}

func resourceTencentCloudRumInstanceStatusConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_instance_status_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
