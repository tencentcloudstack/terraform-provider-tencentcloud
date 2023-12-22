package lighthouse

import (
	"context"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudLighthouseDiskAttachment() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_disk_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().AttachDisks(request)
		if e != nil {
			return tccommon.RetryError(e)
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

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"ATTACHED"}, 20*tccommon.ReadRetryTimeout, time.Second, service.LighthouseDiskStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseDiskAttachmentRead(d, meta)
}

func resourceTencentCloudLighthouseDiskAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_disk_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	diskAttachment, err := service.DescribeLighthouseDiskById(ctx, d.Id())
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

	return nil
}

func resourceTencentCloudLighthouseDiskAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_disk_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	if err := service.DeleteLighthouseDiskAttachmentById(ctx, d.Id()); err != nil {
		return err
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"UNATTACHED"}, 20*tccommon.ReadRetryTimeout, time.Second, service.LighthouseDiskStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
