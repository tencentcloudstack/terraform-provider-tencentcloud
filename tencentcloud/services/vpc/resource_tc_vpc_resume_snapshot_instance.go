package vpc

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcResumeSnapshotInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcResumeSnapshotInstanceCreate,
		Read:   resourceTencentCloudVpcResumeSnapshotInstanceRead,
		Delete: resourceTencentCloudVpcResumeSnapshotInstanceDelete,
		Schema: map[string]*schema.Schema{
			"snapshot_policy_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Snapshot policy Id.",
			},

			"snapshot_file_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Snapshot file Id.",
			},

			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "InstanceId.",
			},
		},
	}
}

func resourceTencentCloudVpcResumeSnapshotInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_resume_snapshot_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request          = vpc.NewResumeSnapshotInstanceRequest()
		snapshotPolicyId string
		snapshotFileId   string
		instanceId       string
	)
	if v, ok := d.GetOk("snapshot_policy_id"); ok {
		snapshotPolicyId = v.(string)
		request.SnapshotPolicyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("snapshot_file_id"); ok {
		snapshotFileId = v.(string)
		request.SnapshotFileId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ResumeSnapshotInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate vpc resumeSnapshotInstance failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(snapshotPolicyId + tccommon.FILED_SP + snapshotFileId + tccommon.FILED_SP + instanceId)

	return resourceTencentCloudVpcResumeSnapshotInstanceRead(d, meta)
}

func resourceTencentCloudVpcResumeSnapshotInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_resume_snapshot_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpcResumeSnapshotInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_resume_snapshot_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
