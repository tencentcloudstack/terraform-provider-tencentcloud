/*
Provides a resource to create a dnspod download_snapshot
Example Usage
```hcl
resource "tencentcloud_dnspod_download_snapshot_operation" "download_snapshot" {
  domain = "dnspod.cn"
  snapshot_id = "456"
}
```
*/
package tencentcloud

import (
	"strings"

	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDnspodDownloadSnapshotOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodDownloadSnapshotOperationCreate,
		Read:   resourceTencentCloudDnspodDownloadSnapshotOperationRead,
		Delete: resourceTencentCloudDnspodDownloadSnapshotOperationDelete,
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

			"cos_url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Snapshot download url.",
			},
		},
	}
}

func resourceTencentCloudDnspodDownloadSnapshotOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_download_snapshot_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dnspod.NewDownloadSnapshotRequest()
		response   = dnspod.NewDownloadSnapshotResponse()
		domain     string
		snapshotId string
		cosUrl     string
	)
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("snapshot_id"); ok {
		snapshotId = v.(string)
		request.SnapshotId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDnsPodClient().DownloadSnapshot(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dnspod download_snapshot failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.CosUrl != nil {
		cosUrl = *response.Response.CosUrl
		_ = d.Set("cos_url", cosUrl)
	}

	d.SetId(strings.Join([]string{domain, snapshotId}, FILED_SP))

	return resourceTencentCloudDnspodDownloadSnapshotOperationRead(d, meta)
}

func resourceTencentCloudDnspodDownloadSnapshotOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_download_snapshot_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDnspodDownloadSnapshotOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_download_snapshot_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
