package crs

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudRedisBackupEncrypt() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisBackupEncryptCreate,
		Read:   resourceTencentCloudRedisBackupEncryptRead,
		Update: resourceTencentCloudRedisBackupEncryptUpdate,
		Delete: resourceTencentCloudRedisBackupEncryptDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of redis instance.",
			},

			"backup_encrypt": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable backup encryption. true: enable; false: disable.",
			},

			"encrypt_key_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Encryption key type. 0: Tencent Cloud managed (default); 1: user-defined KMS key.",
			},

			"kms_key_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The KMS key ID.",
			},

			"kms_region": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The KMS region.",
			},
		},
	}
}

func resourceTencentCloudRedisBackupEncryptCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_encrypt.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisBackupEncryptUpdate(d, meta)
}

func resourceTencentCloudRedisBackupEncryptRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_encrypt.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()

	request := redis.NewDescribeBackupEncryptConfigRequest()
	request.InstanceId = helper.String(instanceId)

	var response *redis.DescribeBackupEncryptConfigResponse
	var notFound bool
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().DescribeBackupEncryptConfigWithContext(ctx, request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "ResourceNotFound.InstanceNotExists" {
					return resource.NonRetryableError(e)
				}
			}
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil {
			notFound = true
			return nil
		}
		response = result
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read redis_backup_encrypt failed, reason:%+v", logId, err)
		return err
	}

	if notFound {
		log.Printf("[CRUD] redis_backup_encrypt id=%s", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("instance_id", instanceId)

	if response.Response != nil {
		if response.Response.BackupEncrypt != nil {
			_ = d.Set("backup_encrypt", response.Response.BackupEncrypt)
		}
		if response.Response.EncryptKeyType != nil {
			_ = d.Set("encrypt_key_type", response.Response.EncryptKeyType)
		}
		if response.Response.KmsKeyId != nil {
			_ = d.Set("kms_key_id", response.Response.KmsKeyId)
		}
		if response.Response.KmsRegion != nil {
			_ = d.Set("kms_region", response.Response.KmsRegion)
		}
	}

	return nil
}

func resourceTencentCloudRedisBackupEncryptUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_encrypt.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()
	request := redis.NewModifyBackupEncryptRequest()
	request.InstanceId = helper.String(instanceId)

	if v, ok := d.GetOkExists("backup_encrypt"); ok {
		request.BackupEncrypt = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("encrypt_key_type"); ok {
		request.EncryptKeyType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		request.KmsKeyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kms_region"); ok {
		request.KmsRegion = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().ModifyBackupEncryptWithContext(ctx, request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "ResourceNotFound.InstanceNotExists" {
					return resource.NonRetryableError(e)
				}
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis_backup_encrypt failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudRedisBackupEncryptRead(d, meta)
}

func resourceTencentCloudRedisBackupEncryptDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_encrypt.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
