/*
Provides a resource to create a dnspod snapshot

Example Usage

```hcl
resource "tencentcloud_dnspod_snapshot" "snapshot" {
  domain = "dnspod.cn"
}
```

Import

dnspod snapshot can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_snapshot.snapshot domain#snapshot_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDnspodSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodSnapshotCreate,
		Read:   resourceTencentCloudDnspodSnapshotRead,
		Delete: resourceTencentCloudDnspodSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Domain name.",
			},
		},
	}
}

func resourceTencentCloudDnspodSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_snapshot.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = dnspod.NewCreateSnapshotRequest()
		// response   = dnspod.NewCreateSnapshotResponse()
		// domain string
		// snapshotId string
	)
	if v, ok := d.GetOk("domain"); ok {
		// domain = v.(string)
		request.Domain = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDnsPodClient().CreateSnapshot(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		// response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dnspod snapshot failed, reason:%+v", logId, err)
		return err
	}

	// snapshotId = *response.Response.SnapshotId
	// d.SetId(strings.Join([]string{domain, snapshotId}, FILED_SP))

	return resourceTencentCloudDnspodSnapshotRead(d, meta)
}

func resourceTencentCloudDnspodSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_snapshot.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DnspodService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_dnspod_snapshot id is broken, id is %s", d.Id())
	}
	domain := idSplit[0]
	snapshotId := idSplit[1]

	snapshot, err := service.DescribeDnspodSnapshotById(ctx, domain, snapshotId)
	if err != nil {
		return err
	}

	if snapshot == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DnspodSnapshot` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("domain", domain)
	_ = d.Set("snapshot_id", snapshotId)

	// if snapshot.DomainId != nil {
	// 	_ = d.Set("domain_id", snapshot.DomainId)
	// }

	return nil
}

func resourceTencentCloudDnspodSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_snapshot.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DnspodService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_dnspod_snapshot id is broken, id is %s", d.Id())
	}
	// domain := idSplit[0]
	snapshotId := idSplit[1]

	if err := service.DeleteDnspodSnapshotById(ctx, snapshotId); err != nil {
		return err
	}

	return nil
}
