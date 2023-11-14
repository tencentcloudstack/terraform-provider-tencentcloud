/*
Provides a resource to create a rum instance_status_attachment

Example Usage

```hcl
resource "tencentcloud_rum_instance_status_attachment" "instance_status_attachment" {
  instance_id = "rum-xxx"
}
```

Import

rum instance_status_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_rum_instance_status_attachment.instance_status_attachment instance_status_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudRumInstanceStatusAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRumInstanceStatusAttachmentCreate,
		Read:   resourceTencentCloudRumInstanceStatusAttachmentRead,
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
		},
	}
}

func resourceTencentCloudRumInstanceStatusAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_instance_status_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = rum.NewResumeInstanceRequest()
		response   = rum.NewResumeInstanceResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRumClient().ResumeInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create rum instanceStatusAttachment failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudRumInstanceStatusAttachmentRead(d, meta)
}

func resourceTencentCloudRumInstanceStatusAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_instance_status_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceStatusAttachmentId := d.Id()

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

	return nil
}

func resourceTencentCloudRumInstanceStatusAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_instance_status_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceStatusAttachmentId := d.Id()

	if err := service.DeleteRumInstanceStatusAttachmentById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
