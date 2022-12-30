/*
Provides a resource to create a cfs auto_snapshot_policy_attachment

Example Usage

```hcl
resource "tencentcloud_cfs_auto_snapshot_policy_attachment" "auto_snapshot_policy_attachment" {
  auto_snapshot_policy_id = "asp-basic"
  file_system_ids         = "cfs-4xzkct19,cfs-iobiaxtj"
}
```

Import

cfs auto_snapshot_policy_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cfs_auto_snapshot_policy_attachment.auto_snapshot_policy_attachment auto_snapshot_policy_id#file_system_ids
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCfsAutoSnapshotPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfsAutoSnapshotPolicyAttachmentCreate,
		Read:   resourceTencentCloudCfsAutoSnapshotPolicyAttachmentRead,
		Delete: resourceTencentCloudCfsAutoSnapshotPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auto_snapshot_policy_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of the snapshot to be unbound.",
			},

			"file_system_ids": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "List of IDs of the file systems to be unbound, separated by comma.",
			},
		},
	}
}

func resourceTencentCloudCfsAutoSnapshotPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfs_auto_snapshot_policy_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request              = cfs.NewBindAutoSnapshotPolicyRequest()
		autoSnapshotPolicyId string
		fileSystemIds        string
	)
	if v, ok := d.GetOk("auto_snapshot_policy_id"); ok {
		autoSnapshotPolicyId = v.(string)
		request.AutoSnapshotPolicyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("file_system_ids"); ok {
		fileSystemIds = v.(string)
		request.FileSystemIds = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfsClient().BindAutoSnapshotPolicy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cfs autoSnapshotPolicyAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(autoSnapshotPolicyId + FILED_SP + fileSystemIds)

	return resourceTencentCloudCfsAutoSnapshotPolicyAttachmentRead(d, meta)
}

func resourceTencentCloudCfsAutoSnapshotPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfs_auto_snapshot_policy_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CfsService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	autoSnapshotPolicyId := idSplit[0]
	fileSystemIds := idSplit[1]

	autoSnapshotPolicyAttachment, err := service.DescribeCfsAutoSnapshotPolicyAttachmentById(ctx, autoSnapshotPolicyId, fileSystemIds)
	if err != nil {
		return err
	}

	if autoSnapshotPolicyAttachment == nil {
		d.SetId("")
		return fmt.Errorf("resource `tencentcloud_cfs_auto_snapshot_policy_attachment` %s does not exist", d.Id())
	}

	if autoSnapshotPolicyAttachment.AutoSnapshotPolicyId != nil {
		_ = d.Set("auto_snapshot_policy_id", autoSnapshotPolicyId)
	}

	if autoSnapshotPolicyAttachment.FileSystems != nil {
		_ = d.Set("file_system_ids", fileSystemIds)
	}

	return nil
}

func resourceTencentCloudCfsAutoSnapshotPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfs_auto_snapshot_policy_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CfsService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	autoSnapshotPolicyId := idSplit[0]
	fileSystemIds := idSplit[1]

	if err := service.DeleteCfsAutoSnapshotPolicyAttachmentById(ctx, autoSnapshotPolicyId, fileSystemIds); err != nil {
		return err
	}

	return nil
}
