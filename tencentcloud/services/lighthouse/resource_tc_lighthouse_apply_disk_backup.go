package lighthouse

import (
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudLighthouseApplyDiskBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseApplyDiskBackupCreate,
		Read:   resourceTencentCloudLighthouseApplyDiskBackupRead,
		Delete: resourceTencentCloudLighthouseApplyDiskBackupDelete,
		Schema: map[string]*schema.Schema{
			"disk_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Disk ID.",
			},

			"disk_backup_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Disk backup ID.",
			},
		},
	}
}

func resourceTencentCloudLighthouseApplyDiskBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_apply_disk_backup.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request      = lighthouse.NewApplyDiskBackupRequest()
		diskBackupId string
		diskId       string
	)
	if v, ok := d.GetOk("disk_id"); ok {
		diskId = v.(string)
		request.DiskId = helper.String(diskId)
	}

	if v, ok := d.GetOk("disk_backup_id"); ok {
		diskBackupId = v.(string)
		request.DiskBackupId = helper.String(diskBackupId)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().ApplyDiskBackup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate lighthouse applyDiskBackup failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(diskId + tccommon.FILED_SP + diskBackupId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*tccommon.ReadRetryTimeout, time.Second, service.LighthouseApplyDiskBackupStateRefreshFunc(diskBackupId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseApplyDiskBackupRead(d, meta)
}

func resourceTencentCloudLighthouseApplyDiskBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_apply_disk_backup.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLighthouseApplyDiskBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_apply_disk_backup.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
