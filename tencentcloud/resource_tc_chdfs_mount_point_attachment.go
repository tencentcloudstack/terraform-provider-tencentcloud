/*
Provides a resource to create a chdfs mount_point_attachment

Example Usage

```hcl
resource "tencentcloud_chdfs_mount_point_attachment" "mount_point_attachment" {
  access_group_ids = [
    "ag-bvmzrbsm",
    "ag-lairqrgr",
  ]
  mount_point_id   = "f14mpfy5lh4e-KuiL"
}
```

Import

chdfs mount_point_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_mount_point_attachment.mount_point_attachment mount_point_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	chdfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/chdfs/v20201112"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudChdfsMountPointAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudChdfsMountPointAttachmentCreate,
		Read:   resourceTencentCloudChdfsMountPointAttachmentRead,
		Delete: resourceTencentCloudChdfsMountPointAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mount_point_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "associate mount point.",
			},

			"access_group_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "associate access group id.",
			},
		},
	}
}

func resourceTencentCloudChdfsMountPointAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_mount_point_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = chdfs.NewAssociateAccessGroupsRequest()
		mountPointId string
	)
	if v, ok := d.GetOk("mount_point_id"); ok {
		mountPointId = v.(string)
		request.MountPointId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("access_group_ids"); ok {
		accessGroupIdsSet := v.(*schema.Set).List()
		for i := range accessGroupIdsSet {
			accessGroupIds := accessGroupIdsSet[i].(string)
			request.AccessGroupIds = append(request.AccessGroupIds, &accessGroupIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseChdfsClient().AssociateAccessGroups(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create chdfs mountPointAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(mountPointId)

	return resourceTencentCloudChdfsMountPointAttachmentRead(d, meta)
}

func resourceTencentCloudChdfsMountPointAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_mount_point_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}

	mountPointId := d.Id()

	mountPointAttachment, err := service.DescribeChdfsMountPointById(ctx, mountPointId)
	if err != nil {
		return err
	}

	if mountPointAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ChdfsMountPointAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if mountPointAttachment.MountPointId != nil {
		_ = d.Set("mount_point_id", mountPointAttachment.MountPointId)
	}

	if mountPointAttachment.AccessGroupIds != nil {
		_ = d.Set("access_group_ids", mountPointAttachment.AccessGroupIds)
	}

	return nil
}

func resourceTencentCloudChdfsMountPointAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_mount_point_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}

	mountPointId := d.Id()

	mountPoint, err := service.DescribeChdfsMountPointById(ctx, mountPointId)
	if err != nil {
		return err
	}

	if err := service.DeleteChdfsMountPointAttachmentById(ctx, mountPointId, mountPoint.AccessGroupIds); err != nil {
		return err
	}

	return nil
}
