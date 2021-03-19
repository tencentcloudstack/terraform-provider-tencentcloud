package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSsmSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSsmSecretCreate,
		Read:   resourceTencentCloudSsmSecretRead,
		Update: resourceTencentCloudSsmSecretUpdate,
		Delete: resourceTencentCloudSsmSecretDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"secret_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Name of secret which cannot be repeated in the same region. The maximum length is 128 bytes. The name can only contain English letters, numbers, underscore and hyphen '-'. The first character must be a letter or number.",
			},
			"init_secret": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "The secret of initial version.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Version of secret. The maximum length is 64 bytes. The version_id can only contain English letters, numbers, underscore and hyphen '-'. The first character must be a letter or number.",
						},
						"secret_binary": {
							Type:         schema.TypeString,
							Optional:     true,
							ExactlyOneOf: []string{"init_secret.0.secret_string"},
							Description:  "The base64-encoded binary secret. secret_binary and secret_string must be set only one, and the maximum support is 4096 bytes. When secret status is `Disabled`, this field will not update anymore.",
						},
						"secret_string": {
							Type:         schema.TypeString,
							Optional:     true,
							ExactlyOneOf: []string{"init_secret.0.secret_binary"},
							Description:  "The string text of secret. secret_binary and secret_string must be set only one, and the maximum support is 4096 bytes. When secret status is `Disabled`, this field will not update anymore.",
						},
					},
				},
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Specify whether to enable secret. Default value is `true`.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of secret. The maximum is 2048 bytes.",
			},
			"recovery_window_in_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specify the scheduled deletion date. Default value is `0` that means to delete immediately. 1-30 means the number of days reserved, completely deleted after this date.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of secret.",
			},
			"kms_key_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "KMS keyId used to encrypt secret. If it is empty, it means that the CMK created by SSM for you by default is used for encryption. You can also specify the KMS CMK created by yourself in the same region for encryption.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of secret.",
			},
		},
	}
}

func resourceTencentCloudSsmSecretCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_secret.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ssmService := SsmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	param := make(map[string]interface{})
	param["secret_name"] = d.Get("secret_name").(string)
	if v, ok := d.GetOk("description"); ok {
		param["description"] = v.(string)
	}
	if v, ok := d.GetOk("kms_key_id"); ok {
		param["kms_key_id"] = v.(string)
	}

	initSecret := d.Get("init_secret").([]interface{})
	versionInfo := initSecret[0].(map[string]interface{})
	param["version_id"] = versionInfo["version_id"].(string)
	if v, ok := versionInfo["secret_binary"]; ok {
		param["secret_binary"] = v.(string)
	}
	if v, ok := versionInfo["secret_string"]; ok {
		param["secret_string"] = v.(string)
	}

	var outErr, inErr error
	var secretName string
	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		secretName, inErr = ssmService.CreateSecret(ctx, param)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}
	d.SetId(secretName)

	if isEnabled := d.Get("is_enabled").(bool); !isEnabled {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = ssmService.DisableSecret(ctx, secretName)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
	}

	var secretInfo *SecretInfo
	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		secretInfo, inErr = ssmService.DescribeSecretByName(ctx, secretName)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("ssm", "secret", tcClient.Region, secretInfo.resourceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	return resourceTencentCloudSsmSecretRead(d, meta)
}

func resourceTencentCloudSsmSecretRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_secret.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ssmService := SsmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	secretName := d.Id()

	var outErr, inErr error
	var secretInfo *SecretInfo
	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		secretInfo, inErr = ssmService.DescribeSecretByName(ctx, secretName)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	if secretInfo == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("secret_name", secretInfo.secretName)
	_ = d.Set("description", secretInfo.description)
	_ = d.Set("kms_key_id", secretInfo.kmsKeyId)
	_ = d.Set("status", secretInfo.status)

	if secretInfo.status == SSM_STATUS_ENABLED {
		_ = d.Set("is_enabled", true)

		secret := d.Get("init_secret").([]interface{})
		var versionId string

		// import secret will import the first version as init_secret
		if len(secret) == 0 {
			var versionIds []string
			outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				versionIds, inErr = ssmService.DescribeSecretVersionIdsByName(ctx, secretName)
				if inErr != nil {
					return retryError(inErr)
				}
				return nil
			})
			if outErr != nil {
				log.Printf("[CRITAL]%s read SSM secret versionId list failed, reason:%+v", logId, outErr)
				return outErr
			}
			if len(versionIds) != 0 {
				versionId = versionIds[0]
			}
		} else {
			versionInfo := secret[0].(map[string]interface{})
			versionId = versionInfo["version_id"].(string)
		}

		if versionId != "" {
			var secretVersionInfo *SecretVersionInfo
			outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				secretVersionInfo, inErr = ssmService.DescribeSecretVersion(ctx, secretName, versionId)
				if inErr != nil {
					return retryError(inErr)
				}
				return nil
			})
			if outErr != nil {
				return outErr
			}

			initSecret := make(map[string]interface{})
			initSecret["version_id"] = secretVersionInfo.versionId
			if secretVersionInfo.secretBinary != "" {
				initSecret["secret_binary"] = secretVersionInfo.secretBinary
			}
			if secretVersionInfo.secretString != "" {
				initSecret["secret_string"] = secretVersionInfo.secretString
			}
			_ = d.Set("init_secret", []map[string]interface{}{initSecret})
		}
	} else {
		_ = d.Set("is_enabled", false)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "ssm", "secret", tcClient.Region, secretInfo.resourceId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudSsmSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_secret.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ssmService := SsmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	d.Partial(true)
	secretName := d.Id()

	if d.HasChange("description") {
		description := d.Get("description").(string)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := ssmService.UpdateSecretDescription(ctx, secretName, description)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify SSM secret description failed, reason:%+v", logId, err)
			return err
		}
		d.SetPartial("description")
	}

	if d.HasChange("is_enabled") {
		isEnabled := d.Get("is_enabled").(bool)
		err := updateSecretIsEnabled(ctx, ssmService, secretName, isEnabled)
		if err != nil {
			log.Printf("[CRITAL]%s modify SSM secret status failed, reason:%+v", logId, err)
			return err
		}
		d.SetPartial("is_enabled")
	}

	var outErr, inErr error
	var secretInfo *SecretInfo
	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		secretInfo, inErr = ssmService.DescribeSecretByName(ctx, secretName)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	if secretInfo.status == SSM_STATUS_ENABLED {
		err := updateSecretVersionInfo(ctx, d, ssmService)
		if err != nil {
			log.Printf("[CRITAL]%s modify SSM secret version failed, reason:%+v", logId, err)
			return err
		}
		d.SetPartial("init_secret.0.version_id")
		d.SetPartial("init_secret.0.secret_binary")
		d.SetPartial("init_secret.0.secret_string")
	}

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		secretInfo, err := ssmService.DescribeSecretByName(ctx, secretName)
		if err != nil {
			return err
		}
		resourceName := BuildTagResourceName("ssm", "secret", tcClient.Region, secretInfo.resourceId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceTencentCloudSsmSecretRead(d, meta)
}

func resourceTencentCloudSsmSecretDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_secret.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ssmService := SsmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	secretName := d.Id()
	recoveryWindowInDays := d.Get("recovery_window_in_days").(int)
	isEnabled := d.Get("is_enabled").(bool)
	if isEnabled {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := ssmService.DisableSecret(ctx, secretName)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify SSM secret status failed, reason:%+v", logId, err)
			return err
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := ssmService.DeleteSecret(ctx, secretName, uint64(recoveryWindowInDays))
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete SSM secret failed, reason:%+v", logId, err)
		return err
	}

	return resource.Retry(readRetryTimeout, func() *resource.RetryError {
		secretInfo, e := ssmService.DescribeSecretByName(ctx, secretName)
		if e != nil {
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok && sdkError.Code == "ResourceNotFound" {
				return nil
			}
			return retryError(err)
		}
		if secretInfo.status == SSM_STATUS_PENDINGDELETE {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("delete fail"))
	})
}

func updateSecretIsEnabled(ctx context.Context, ssmService SsmService, secretName string, isEnabled bool) error {
	var err error
	if isEnabled {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := ssmService.EnableSecret(ctx, secretName)
			if e != nil {
				return retryError(e)
			}
			return nil
		})

	} else {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := ssmService.DisableSecret(ctx, secretName)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
	}
	return err
}

func updateSecretVersionInfo(ctx context.Context, d *schema.ResourceData, ssmService SsmService) error {
	logId := getLogId(ctx)

	param := make(map[string]interface{})
	param["secret_name"] = d.Get("secret_name").(string)
	param["version_id"] = d.Get("init_secret.0.version_id").(string)
	if v, ok := d.GetOk("init_secret.0.secret_binary"); ok {
		param["secret_binary"] = v.(string)
	}
	if v, ok := d.GetOk("init_secret.0.secret_string"); ok {
		param["secret_string"] = v.(string)
	}
	if d.HasChange("init_secret.0.version_id") {
		oldVersionId, newVersionId := d.GetChange("init_secret.0.version_id")
		if oldVersionId.(string) != "" {
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				e := ssmService.DeleteSecretVersion(ctx, d.Get("secret_name").(string), oldVersionId.(string))
				if e != nil {
					return retryError(e)
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s delete SSM secret version failed, reason:%+v", logId, err)
				return err
			}
		}

		if newVersionId.(string) != "" {
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				_, _, e := ssmService.PutSecretValue(ctx, param)
				if e != nil {
					return retryError(e)
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s add SSM secret version failed, reason:%+v", logId, err)
				return err
			}
		}
	} else if d.HasChange("init_secret.0.secret_binary") || d.HasChange("init_secret.0.secret_string") {
		versionId := d.Get("init_secret.0.version_id").(string)
		if versionId != "" {
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				e := ssmService.UpdateSecret(ctx, param)
				if e != nil {
					return retryError(e)
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s modify SSM secret content failed, reason:%+v", logId, err)
				return err
			}
		}
	}
	return nil
}
