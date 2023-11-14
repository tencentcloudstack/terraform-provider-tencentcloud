/*
Provides a resource to create a lighthouse disk_attachment

Example Usage

```hcl
resource "tencentcloud_lighthouse_disk_attachment" "disk_attachment" {
  disk_ids =
  instance_id = "lhins-123456"
  renew_flag = "NOTIFY_AND_MANUAL_RENEW"
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
	"time"
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
			"disk_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of cloud disk IDs.",
			},

			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"renew_flag": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Whether Auto-Renewal is enabled.",
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
		response   = lighthouse.NewAttachDisksResponse()
		diskIds    string
		instanceId string
	)
	if v, ok := d.GetOk("disk_ids"); ok {
		diskIdsSet := v.(*schema.Set).List()
		for i := range diskIdsSet {
			diskIds := diskIdsSet[i].(string)
			request.DiskIds = append(request.DiskIds, &diskIds)
		}
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("renew_flag"); ok {
		request.RenewFlag = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().AttachDisks(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create lighthouse diskAttachment failed, reason:%+v", logId, err)
		return err
	}

	diskIds = *response.Response.DiskIds
	d.SetId(strings.Join([]string{diskIds, instanceId}, FILED_SP))

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

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

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	diskIds := idSplit[0]
	instanceId := idSplit[1]

	diskAttachment, err := service.DescribeLighthouseDiskAttachmentById(ctx, diskIds, instanceId)
	if err != nil {
		return err
	}

	if diskAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseDiskAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if diskAttachment.DiskIds != nil {
		_ = d.Set("disk_ids", diskAttachment.DiskIds)
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

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	diskIds := idSplit[0]
	instanceId := idSplit[1]

	if err := service.DeleteLighthouseDiskAttachmentById(ctx, diskIds, instanceId); err != nil {
		return err
	}

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"UNATTACHED"}, 20*readRetryTimeout, time.Second, service.LighthouseDiskAttachmentStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
