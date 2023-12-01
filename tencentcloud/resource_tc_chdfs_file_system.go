/*
Provides a resource to create a chdfs file_system

Example Usage

```hcl
resource "tencentcloud_chdfs_file_system" "file_system" {
  capacity_quota           = 10995116277760
  description              = "file system for terraform test"
  enable_ranger            = true
  file_system_name         = "terraform-test"
  posix_acl                = false
  ranger_service_addresses = [
    "127.0.0.1:80",
    "127.0.0.1:8000",
  ]
  super_users              = [
    "terraform",
    "iac",
  ]
}
```

Import

chdfs file_system can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_file_system.file_system file_system_id
```
*/
package tencentcloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	chdfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/chdfs/v20201112"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudChdfsFileSystem() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudChdfsFileSystemCreate,
		Read:   resourceTencentCloudChdfsFileSystemRead,
		Update: resourceTencentCloudChdfsFileSystemUpdate,
		Delete: resourceTencentCloudChdfsFileSystemDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"file_system_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "file system name.",
			},

			"capacity_quota": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "file system capacity. min 1GB, max 1PB, CapacityQuota is N * 1073741824.",
			},

			"posix_acl": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "check POSIX ACL or not.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "desc of the file system.",
			},

			"super_users": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "super users of the file system, default empty.",
			},

			"enable_ranger": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "check the ranger address or not.",
			},

			"ranger_service_addresses": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "ranger address list, default empty.",
			},
		},
	}
}

func resourceTencentCloudChdfsFileSystemCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_file_system.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = chdfs.NewCreateFileSystemRequest()
		response     = chdfs.NewCreateFileSystemResponse()
		fileSystemId string
	)
	if v, ok := d.GetOk("file_system_name"); ok {
		request.FileSystemName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("capacity_quota"); v != nil {
		request.CapacityQuota = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("posix_acl"); v != nil {
		request.PosixAcl = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("super_users"); ok {
		superUsersSet := v.(*schema.Set).List()
		for i := range superUsersSet {
			superUsers := superUsersSet[i].(string)
			request.SuperUsers = append(request.SuperUsers, &superUsers)
		}
	}

	if v, _ := d.GetOk("enable_ranger"); v != nil {
		request.EnableRanger = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("ranger_service_addresses"); ok {
		rangerServiceAddressesSet := v.(*schema.Set).List()
		for i := range rangerServiceAddressesSet {
			rangerServiceAddresses := rangerServiceAddressesSet[i].(string)
			request.RangerServiceAddresses = append(request.RangerServiceAddresses, &rangerServiceAddresses)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseChdfsClient().CreateFileSystem(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create chdfs fileSystem failed, reason:%+v", logId, err)
		return err
	}

	fileSystemId = *response.Response.FileSystem.FileSystemId
	d.SetId(fileSystemId)

	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"2"}, 2*readRetryTimeout, time.Second, service.ChdfsFileSystemStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudChdfsFileSystemRead(d, meta)
}

func resourceTencentCloudChdfsFileSystemRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_file_system.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}

	fileSystemId := d.Id()

	fileSystem, err := service.DescribeChdfsFileSystemById(ctx, fileSystemId)
	if err != nil {
		return err
	}

	if fileSystem == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ChdfsFileSystem` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if fileSystem.FileSystemName != nil {
		_ = d.Set("file_system_name", fileSystem.FileSystemName)
	}

	if fileSystem.CapacityQuota != nil {
		_ = d.Set("capacity_quota", fileSystem.CapacityQuota)
	}

	if fileSystem.PosixAcl != nil {
		_ = d.Set("posix_acl", fileSystem.PosixAcl)
	}

	if fileSystem.Description != nil {
		_ = d.Set("description", fileSystem.Description)
	}

	if fileSystem.SuperUsers != nil {
		_ = d.Set("super_users", fileSystem.SuperUsers)
	}

	if fileSystem.EnableRanger != nil {
		_ = d.Set("enable_ranger", fileSystem.EnableRanger)
	}

	if fileSystem.RangerServiceAddresses != nil {
		_ = d.Set("ranger_service_addresses", fileSystem.RangerServiceAddresses)
	}

	return nil
}

func resourceTencentCloudChdfsFileSystemUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_file_system.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := chdfs.NewModifyFileSystemRequest()

	fileSystemId := d.Id()

	request.FileSystemId = &fileSystemId

	if d.HasChange("file_system_name") {
		if v, ok := d.GetOk("file_system_name"); ok {
			request.FileSystemName = helper.String(v.(string))
		}
	}

	if d.HasChange("capacity_quota") {
		if v, _ := d.GetOk("capacity_quota"); v != nil {
			request.CapacityQuota = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("posix_acl") {
		if v, _ := d.GetOk("posix_acl"); v != nil {
			request.PosixAcl = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("super_users") {
		if v, ok := d.GetOk("super_users"); ok {
			superUsersSet := v.(*schema.Set).List()
			for i := range superUsersSet {
				superUsers := superUsersSet[i].(string)
				request.SuperUsers = append(request.SuperUsers, &superUsers)
			}
		}
	}

	if d.HasChange("enable_ranger") {
		if v, _ := d.GetOk("enable_ranger"); v != nil {
			request.EnableRanger = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("ranger_service_addresses") {
		if v, ok := d.GetOk("ranger_service_addresses"); ok {
			rangerServiceAddressesSet := v.(*schema.Set).List()
			for i := range rangerServiceAddressesSet {
				rangerServiceAddresses := rangerServiceAddressesSet[i].(string)
				request.RangerServiceAddresses = append(request.RangerServiceAddresses, &rangerServiceAddresses)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseChdfsClient().ModifyFileSystem(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update chdfs fileSystem failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudChdfsFileSystemRead(d, meta)
}

func resourceTencentCloudChdfsFileSystemDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_file_system.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}
	fileSystemId := d.Id()

	if err := service.DeleteChdfsFileSystemById(ctx, fileSystemId); err != nil {
		return err
	}

	return nil
}
