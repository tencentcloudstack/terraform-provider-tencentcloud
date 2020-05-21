/*
Provides a resource to create a cloud file system(CFS).

Example Usage

```hcl
resource "tencentcloud_cfs_file_system" "foo" {
  name = "test_file_system"
  availability_zone = "ap-guangzhou-3"
  access_group_id = "pgroup-7nx89k7l"
  protocol = "NFS"
  vpc_id = "vpc-ah9fbkap"
  subnet_id = "subnet-9mu2t9iw"
}
```

Import

Cloud file system can be imported using the id, e.g.

```
$ terraform import tencentcloud_cfs_file_system.foo cfs-6hgquxmj
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudCfsFileSystem() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfsFileSystemCreate,
		Read:   resourceTencentCloudCfsFileSystemRead,
		Update: resourceTencentCloudCfsFileSystemUpdate,
		Delete: resourceTencentCloudCfsFileSystemDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of a file system.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The available zone that the file system locates at.",
			},
			"access_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of a access group.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CFS_PROTOCOL_NFS,
				ValidateFunc: validateAllowedStringValue(CFS_PROTOCOL),
				ForceNew:     true,
				Description:  "File service protocol. Valid values are `NFS` and `CIFS`, and the default is `NFS`.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of a VPC network.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of a subnet.",
			},
			"mount_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "IP of mount point.",
			},

			// computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the file system.",
			},
		},
	}
}

func resourceTencentCloudCfsFileSystemCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfs_file_system.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cfsService := CfsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	request := cfs.NewCreateCfsFileSystemRequest()
	request.Zone = helper.String(d.Get("availability_zone").(string))
	request.PGroupId = helper.String(d.Get("access_group_id").(string))
	request.Protocol = helper.String(d.Get("protocol").(string))
	request.VpcId = helper.String(d.Get("vpc_id").(string))
	request.SubnetId = helper.String(d.Get("subnet_id").(string))
	if v, ok := d.GetOk("name"); ok {
		request.FsName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("mount_ip"); ok {
		request.MountIP = helper.String(v.(string))
	}
	request.NetInterface = helper.String("VPC")
	request.StorageType = helper.String("SD")

	fsId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := meta.(*TencentCloudClient).apiV3Conn.UseCfsClient().CreateCfsFileSystem(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response.Response.FileSystemId == nil {
			err = fmt.Errorf("file system id is nil")
			return resource.NonRetryableError(err)
		}
		fsId = *response.Response.FileSystemId
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(fsId)

	// wait for success status
	err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		fileSystems, errRet := cfsService.DescribeFileSystem(ctx, fsId, "", "")
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if len(fileSystems) < 1 {
			return resource.RetryableError(fmt.Errorf("file system %s not exist", fsId))
		}
		if *fileSystems[0].LifeCycleState == CFS_FILE_SYSTEM_STATUS_CREATING {
			return resource.RetryableError(fmt.Errorf("file system status is %s, retry...", *fileSystems[0].LifeCycleState))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudCfsFileSystemRead(d, meta)
}

func resourceTencentCloudCfsFileSystemRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfs_file_system.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	fsId := d.Id()
	cfsService := CfsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var fileSystem *cfs.FileSystemInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		fileSystems, errRet := cfsService.DescribeFileSystem(ctx, fsId, "", "")
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if len(fileSystems) > 0 {
			fileSystem = fileSystems[0]
		}
		return nil
	})
	if err != nil {
		return err
	}
	if fileSystem == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", fileSystem.FsName)
	_ = d.Set("availability_zone", fileSystem.Zone)
	_ = d.Set("access_group_id", fileSystem.PGroup.PGroupId)
	_ = d.Set("protocol", fileSystem.Protocol)
	_ = d.Set("create_time", fileSystem.CreationTime)

	var mountTarget *cfs.MountInfo
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		targets, errRet := cfsService.DescribeMountTargets(ctx, fsId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if len(targets) > 0 {
			mountTarget = targets[0]
		}
		return nil
	})
	if err != nil {
		return err
	}
	if mountTarget != nil {
		_ = d.Set("vpc_id", mountTarget.VpcId)
		_ = d.Set("subnet_id", mountTarget.SubnetId)
		_ = d.Set("mount_ip", mountTarget.IpAddress)
	}

	return nil
}

func resourceTencentCloudCfsFileSystemUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfs_file_system.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	fsId := d.Id()
	cfsService := CfsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	d.Partial(true)

	if d.HasChange("name") {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := cfsService.ModifyFileSystemName(ctx, fsId, d.Get("name").(string))
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
		d.SetPartial("name")
	}

	if d.HasChange("access_group_id") {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := cfsService.ModifyFileSystemAccessGroup(ctx, fsId, d.Get("access_group_id").(string))
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
		d.SetPartial("access_group_id")
	}

	d.Partial(false)

	return resourceTencentCloudCfsFileSystemRead(d, meta)
}

func resourceTencentCloudCfsFileSystemDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfs_file_system.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	fsId := d.Id()
	cfsService := CfsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		errRet := cfsService.DeleteFileSystem(ctx, fsId)
		if errRet != nil {
			return retryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
