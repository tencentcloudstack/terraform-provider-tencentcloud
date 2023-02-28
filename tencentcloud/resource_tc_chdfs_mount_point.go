/*
Provides a resource to create a chdfs mount_point

Example Usage

```hcl
resource "tencentcloud_chdfs_mount_point" "mount_point" {
  file_system_id     = "f14mpfy5lh4e"
  mount_point_name   = "terraform-test"
  mount_point_status = 1
}
```

Import

chdfs mount_point can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_mount_point.mount_point mount_point_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	chdfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/chdfs/v20201112"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudChdfsMountPoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudChdfsMountPointCreate,
		Read:   resourceTencentCloudChdfsMountPointRead,
		Update: resourceTencentCloudChdfsMountPointUpdate,
		Delete: resourceTencentCloudChdfsMountPointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mount_point_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "mount point name.",
			},

			"file_system_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "file system id you want to mount.",
			},

			"mount_point_status": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "mount status 1:open, 2:close.",
			},
		},
	}
}

func resourceTencentCloudChdfsMountPointCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_mount_point.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = chdfs.NewCreateMountPointRequest()
		response     = chdfs.NewCreateMountPointResponse()
		mountPointId string
	)
	if v, ok := d.GetOk("mount_point_name"); ok {
		request.MountPointName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("file_system_id"); ok {
		request.FileSystemId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("mount_point_status"); v != nil {
		request.MountPointStatus = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseChdfsClient().CreateMountPoint(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create chdfs mountPoint failed, reason:%+v", logId, err)
		return err
	}

	mountPointId = *response.Response.MountPoint.MountPointId
	d.SetId(mountPointId)

	return resourceTencentCloudChdfsMountPointRead(d, meta)
}

func resourceTencentCloudChdfsMountPointRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_mount_point.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}

	mountPointId := d.Id()

	mountPoint, err := service.DescribeChdfsMountPointById(ctx, mountPointId)
	if err != nil {
		return err
	}

	if mountPoint == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ChdfsMountPoint` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if mountPoint.MountPointName != nil {
		_ = d.Set("mount_point_name", mountPoint.MountPointName)
	}

	if mountPoint.FileSystemId != nil {
		_ = d.Set("file_system_id", mountPoint.FileSystemId)
	}

	if mountPoint.Status != nil {
		_ = d.Set("mount_point_status", mountPoint.Status)
	}

	return nil
}

func resourceTencentCloudChdfsMountPointUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_mount_point.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := chdfs.NewModifyMountPointRequest()

	mountPointId := d.Id()

	request.MountPointId = &mountPointId

	if d.HasChange("mount_point_name") {
		if v, ok := d.GetOk("mount_point_name"); ok {
			request.MountPointName = helper.String(v.(string))
		}
	}

	if d.HasChange("mount_point_status") {
		if v, _ := d.GetOk("mount_point_status"); v != nil {
			request.MountPointStatus = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseChdfsClient().ModifyMountPoint(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update chdfs mountPoint failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudChdfsMountPointRead(d, meta)
}

func resourceTencentCloudChdfsMountPointDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_mount_point.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}
	mountPointId := d.Id()

	if err := service.DeleteChdfsMountPointById(ctx, mountPointId); err != nil {
		return err
	}

	return nil
}
