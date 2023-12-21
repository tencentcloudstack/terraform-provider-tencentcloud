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

func ResourceTencentCloudLighthouseApplyInstanceSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseApplyInstanceSnapshotCreate,
		Read:   resourceTencentCloudLighthouseApplyInstanceSnapshotRead,
		Delete: resourceTencentCloudLighthouseApplyInstanceSnapshotDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"snapshot_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Snapshot ID.",
			},
		},
	}
}

func resourceTencentCloudLighthouseApplyInstanceSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_apply_instance_snapshot.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = lighthouse.NewApplyInstanceSnapshotRequest()
		snapshotId string
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("snapshot_id"); ok {
		snapshotId = v.(string)
		request.SnapshotId = helper.String(snapshotId)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().ApplyInstanceSnapshot(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate lighthouse applyInstanceSnapshot failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + tccommon.FILED_SP + snapshotId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*tccommon.ReadRetryTimeout, time.Second, service.LighthouseApplySnapshotStateRefreshFunc(snapshotId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseApplyInstanceSnapshotRead(d, meta)
}

func resourceTencentCloudLighthouseApplyInstanceSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_apply_instance_snapshot.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLighthouseApplyInstanceSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_apply_instance_snapshot.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
