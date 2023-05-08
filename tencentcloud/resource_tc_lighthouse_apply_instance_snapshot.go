/*
Provides a resource to create a lighthouse apply_instance_snapshot

Example Usage

```hcl
resource "tencentcloud_lighthouse_apply_instance_snapshot" "apply_instance_snapshot" {
  instance_id = "lhins-123456"
  snapshot_id = "lhsnap-123456"
}
```
*/
package tencentcloud

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudLighthouseApplyInstanceSnapshot() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_lighthouse_apply_instance_snapshot.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ApplyInstanceSnapshot(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate lighthouse applyInstanceSnapshot failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + FILED_SP + snapshotId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*readRetryTimeout, time.Second, service.LighthouseApplySnapshotStateRefreshFunc(snapshotId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseApplyInstanceSnapshotRead(d, meta)
}

func resourceTencentCloudLighthouseApplyInstanceSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_apply_instance_snapshot.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLighthouseApplyInstanceSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_apply_instance_snapshot.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
