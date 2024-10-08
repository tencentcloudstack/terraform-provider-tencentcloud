// Code generated by iacg; DO NOT EDIT.
package tke

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudKubernetesBackupStorageLocation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesBackupStorageLocationCreate,
		Read:   resourceTencentCloudKubernetesBackupStorageLocationRead,
		Delete: resourceTencentCloudKubernetesBackupStorageLocationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

func resourceTencentCloudKubernetesBackupStorageLocationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_backup_storage_location.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		name string
	)
	var (
		request  = tkev20180525.NewCreateBackupStorageLocationRequest()
		response = tkev20180525.NewCreateBackupStorageLocationResponse()
	)

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().CreateBackupStorageLocationWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create kubernetes backup storage location failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	if err := resourceTencentCloudKubernetesBackupStorageLocationCreatePostHandleResponse0(ctx, response); err != nil {
		return err
	}

	d.SetId(name)

	return resourceTencentCloudKubernetesBackupStorageLocationRead(d, meta)
}

func resourceTencentCloudKubernetesBackupStorageLocationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_backup_storage_location.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	name := d.Id()

	_ = d.Set("name", name)

	respData, err := service.DescribeKubernetesBackupStorageLocationById(ctx, name)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_backup_storage_location` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.StorageRegion != nil {
		_ = d.Set("storage_region", respData.StorageRegion)
	}

	if respData.Bucket != nil {
		_ = d.Set("bucket", respData.Bucket)
	}

	if respData.Path != nil {
		_ = d.Set("path", respData.Path)
	}

	if respData.State != nil {
		_ = d.Set("state", respData.State)
	}

	if respData.Message != nil {
		_ = d.Set("message", respData.Message)
	}

	return nil
}

func resourceTencentCloudKubernetesBackupStorageLocationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_backup_storage_location.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	name := d.Id()

	var (
		request  = tkev20180525.NewDeleteBackupStorageLocationRequest()
		response = tkev20180525.NewDeleteBackupStorageLocationResponse()
	)

	request.Name = helper.String(name)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DeleteBackupStorageLocationWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete kubernetes backup storage location failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	if err := resourceTencentCloudKubernetesBackupStorageLocationDeletePostHandleResponse0(ctx, response); err != nil {
		return err
	}

	return nil
}
