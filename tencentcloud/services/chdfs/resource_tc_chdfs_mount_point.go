package chdfs

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	chdfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/chdfs/v20201112"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudChdfsMountPoint() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_chdfs_mount_point.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseChdfsClient().CreateMountPoint(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_chdfs_mount_point.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ChdfsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

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
	defer tccommon.LogElapsed("resource.tencentcloud_chdfs_mount_point.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseChdfsClient().ModifyMountPoint(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_chdfs_mount_point.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ChdfsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	mountPointId := d.Id()

	if err := service.DeleteChdfsMountPointById(ctx, mountPointId); err != nil {
		return err
	}

	return nil
}
