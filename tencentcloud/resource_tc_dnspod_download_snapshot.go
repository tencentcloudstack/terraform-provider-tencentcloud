/*
Provides a resource to create a dnspod download_snapshot

Example Usage

```hcl
resource "tencentcloud_dnspod_download_snapshot" "download_snapshot" {
  domain = "dnspod.cn"
  snapshot_id = "456"
  domain_id = 123
}
```

Import

dnspod download_snapshot can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_download_snapshot.download_snapshot download_snapshot_id
```
*/
package tencentcloud

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudDnspodDownloadSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodDownloadSnapshotCreate,
		Read:   resourceTencentCloudDnspodDownloadSnapshotRead,
		Delete: resourceTencentCloudDnspodDownloadSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
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

func resourceTencentCloudDnspodDownloadSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_download_snapshot.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = dnspod.NewDownloadSnapshotRequest()
		// response = dnspod.NewDownloadSnapshotResponse()
		domain   string
		snapshotId string
	)
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("snapshot_id"); ok {
		snapshotId = v.(string)
		request.SnapshotId = helper.String(v.(string))
	}

	// if v,  ok := d.GetOkExists("domain_id"); ok {
	// 	request.DomainId = helper.IntUint64(v.(int))
	// }

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDnsPodClient().DownloadSnapshot(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		// response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dnspod download_snapshot failed, reason:%+v", logId, err)
		return err
	}

	// domain = *response.Response.Domain
	d.SetId(strings.Join([]string{domain, snapshotId}, FILED_SP))

	return resourceTencentCloudDnspodDownloadSnapshotRead(d, meta)
}

func resourceTencentCloudDnspodDownloadSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_download_snapshot.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDnspodDownloadSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_download_snapshot.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
