/*
Provides a resource to create a lighthouse disk_attachment

Example Usage

```hcl
resource "tencentcloud_lighthouse_disk_attachment" "disk_attachment" {
  disk_id = "lhdisk-xxxxxx"
  instance_id = "lhins-xxxxxx"
}
```

Import

lighthouse disk_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_disk_attachment.disk_attachment disk_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudLighthouseDiskAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseDiskAttachmentCreate,
		Read:   resourceTencentCloudLighthouseDiskAttachmentRead,
		Delete: resourceTencentCloudLighthouseDiskAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"disk_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Disk id.",
			},

			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
		},
	}
}

func resourceTencentCloudLighthouseDiskAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_disk_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = lighthouse.NewAttachDisksRequest()
		diskId     string
		instanceId string
	)
	if v, ok := d.GetOk("disk_id"); ok {
		diskId = v.(string)
		request.DiskIds = []*string{&diskId}
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().AttachDisks(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create lighthouse diskAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(diskId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"ATTACHED"}, 20*readRetryTimeout, time.Second, service.LighthouseDiskAttachmentStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseDiskAttachmentRead(d, meta)
}

func resourceTencentCloudLighthouseDiskAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_disk_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	diskAttachment, err := service.DescribeLighthouseDiskAttachmentById(ctx, d.Id())
	if err != nil {
		return err
	}

	if diskAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseDiskAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if diskAttachment.DiskId != nil {
		_ = d.Set("disk_id", diskAttachment.DiskId)
	}

	if diskAttachment.InstanceId != nil {
		_ = d.Set("instance_id", diskAttachment.InstanceId)
	}

	if diskAttachment.RenewFlag != nil {
		_ = d.Set("renew_flag", diskAttachment.RenewFlag)
	}

	return nil
}

func resourceTencentCloudLighthouseDiskAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_disk_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	if err := service.DeleteLighthouseDiskAttachmentById(ctx, d.Id()); err != nil {
		return err
	}

	conf := BuildStateChangeConf([]string{}, []string{"UNATTACHED"}, 20*readRetryTimeout, time.Second, service.LighthouseDiskAttachmentStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
