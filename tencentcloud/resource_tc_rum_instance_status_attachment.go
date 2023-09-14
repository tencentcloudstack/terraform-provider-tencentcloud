/*
Provides a resource to create a rum instance_status_attachment

Example Usage

```hcl
resource "tencentcloud_rum_instance_status_attachment" "instance_status_attachment" {
  instance_id = "rum-pasZKEI3RLgakj"
  operate     = "stop"
}
```

Import

rum instance_status_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_rum_instance_status_attachment.instance_status_attachment instance_id
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

func resourceTencentCloudRumInstanceStatusAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRumInstanceStatusAttachmentCreate,
		Read:   resourceTencentCloudRumInstanceStatusAttachmentRead,
		Update: resourceTencentCloudRumInstanceStatusAttachmentUpdate,
		Delete: resourceTencentCloudRumInstanceStatusAttachmentDelete,
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

func resourceTencentCloudRumInstanceStatusAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_instance_status_attachment.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRumInstanceStatusAttachmentUpdate(d, meta)
}

func resourceTencentCloudRumInstanceStatusAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_instance_status_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()
	instanceStatusAttachment, err := service.DescribeRumInstanceStatusAttachmentById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instanceStatusAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RumInstanceStatusAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instanceStatusAttachment.InstanceId != nil {
		_ = d.Set("instance_id", instanceStatusAttachment.InstanceId)
	}

	if instanceStatusAttachment.InstanceStatus != nil {
		_ = d.Set("instance_status", instanceStatusAttachment.InstanceStatus)
	}

	return nil
}

func resourceTencentCloudRumInstanceStatusAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_instance_status_attachment.update")()
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

	return resourceTencentCloudRumInstanceStatusAttachmentRead(d, meta)
}

func resourceTencentCloudRumInstanceStatusAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_instance_status_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
