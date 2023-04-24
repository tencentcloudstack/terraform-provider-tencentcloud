/*
Provide a resource to create tke backup storage location.

~> **NOTE:** To create this resource, you need to create a cos bucket with prefix "tke-backup" in advance.

Example Usage
```
resource "tencentcloud_kubernetes_backup_storage_location" "example_backup" {
  name            = "example-backup-1"
  storage_region  = "ap-guangzhou" # region of you pre-created COS bucket
  bucket          = "tke-backup-example-1" # bucket name of your pre-created COS bucket
}
```

Import

tke backup storage location can be imported, e.g.

```
$ terraform import tencentcloud_kubernetes_backup_storage_location.test xxx
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTkeBackupStorageLocation() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Create: resourceTencentCloudTkeBackupStorageLocationCreate,
		Read:   resourceTencentCloudTkeBackupStorageLocationRead,
		Delete: resourceTencentCloudTkeBackupStorageLocationDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the backup storage location.",
			},
			"storage_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region of the storage.",
			},
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the bucket.",
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Prefix of the bucket.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the backup storage location.",
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Message of the backup storage location.",
			},
		},
	}
}

func resourceTencentCloudTkeBackupStorageLocationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_backup_storage_location.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := TkeService{client: client}

	request := genCreateBackupStorageLocationRequest(d)
	err := service.createBackupStorageLocation(ctx, request)
	if err != nil {
		return err
	}

	// wait for status ok
	err = resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		locations, errRet := service.describeBackupStorageLocations(ctx, []string{d.Get("name").(string)})
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if len(locations) != 1 {
			resource.RetryableError(fmt.Errorf("more than 1 location returnen in api response, expected 1 but got %d", len(locations)))
		}
		if locations[0].State == nil {
			return resource.RetryableError(fmt.Errorf("location %s is still in state nil", d.Get("name").(string)))
		}
		if len(locations) == 1 && *locations[0].State == backupStorageLocationStateAvailable {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("location %s is still in state %s", d.Get("name").(string), *locations[0].State))
	})
	if err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))
	return resourceTencentCloudTkeBackupStorageLocationRead(d, meta)
}

func resourceTencentCloudTkeBackupStorageLocationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_backup_storage_location.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := TkeService{client: client}

	locations, err := service.describeBackupStorageLocations(ctx, []string{d.Id()})
	if err != nil {
		return err
	}
	has := false
	for _, location := range locations {
		if *location.Name == d.Id() {
			has = true
			_ = d.Set("name", location.Name)
			_ = d.Set("storage_region", location.StorageRegion)
			_ = d.Set("bucket", location.Bucket)
			_ = d.Set("path", location.Path)
			_ = d.Set("state", location.State)
			_ = d.Set("message", location.Message)
			return nil
		}
	}
	if !has {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceTencentCloudTkeBackupStorageLocationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_backup_storage_location.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := TkeService{client: client}

	err := service.deleteBackupStorageLocation(ctx, d.Id())
	if err != nil {
		return err
	}

	// wait until location is deleted
	err = resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		locations, errRet := service.describeBackupStorageLocations(ctx, []string{d.Id()})
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if len(locations) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("location %s is still not deleted", d.Id()))
	})

	return nil
}

func genCreateBackupStorageLocationRequest(d *schema.ResourceData) (request *tke.CreateBackupStorageLocationRequest) {
	request = tke.NewCreateBackupStorageLocationRequest()
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}
	if v, ok := d.GetOk("storage_region"); ok {
		request.StorageRegion = helper.String(v.(string))
	}
	if v, ok := d.GetOk("bucket"); ok {
		request.Bucket = helper.String(v.(string))
	}
	if v, ok := d.GetOk("path"); ok {
		request.Path = helper.String(v.(string))
	}
	return request
}
