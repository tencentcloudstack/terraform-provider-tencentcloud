package tke

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTkeBackupStorageLocation() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_backup_storage_location.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := TkeService{client: client}

	request := genCreateBackupStorageLocationRequest(d)
	err := service.createBackupStorageLocation(ctx, request)
	if err != nil {
		return err
	}

	// wait for status ok
	err = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		locations, errRet := service.DescribeBackupStorageLocations(ctx, []string{d.Get("name").(string)})
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
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
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_backup_storage_location.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := TkeService{client: client}

	locations, err := service.DescribeBackupStorageLocations(ctx, []string{d.Id()})
	if err != nil {
		return err
	}
	for _, location := range locations {
		if *location.Name == d.Id() {
			_ = d.Set("name", location.Name)
			_ = d.Set("storage_region", location.StorageRegion)
			_ = d.Set("bucket", location.Bucket)
			_ = d.Set("path", location.Path)
			_ = d.Set("state", location.State)
			_ = d.Set("message", location.Message)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceTencentCloudTkeBackupStorageLocationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_backup_storage_location.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := TkeService{client: client}

	err := service.DeleteBackupStorageLocation(ctx, d.Id())
	if err != nil {
		return err
	}

	// wait until location is deleted
	err = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		locations, errRet := service.DescribeBackupStorageLocations(ctx, []string{d.Id()})
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if len(locations) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("location %s is still not deleted", d.Id()))
	})
	if err != nil {
		return err
	}

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
