package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func resourceTencentCloudSsmSecretVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSsmSecretVersionCreate,
		Read:   resourceTencentCloudSsmSecretVersionRead,
		Update: resourceTencentCloudSsmSecretVersionUpdate,
		Delete: resourceTencentCloudSsmSecretVersionDelete,
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
			"version_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Version of secret. The maximum length is 64 bytes. The version_id can only contain English letters, numbers, underscore and hyphen '-'. The first character must be a letter or number.",
			},
			"secret_binary": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"secret_string"},
				Description:  "The base64-encoded binary secret. secret_binary and secret_string must be set only one, and the maximum support is 4096 bytes. When secret status is `Disabled`, this field will not update anymore.",
			},
			"secret_string": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"secret_binary"},
				Description:  "The string text of secret. secret_binary and secret_string must be set only one, and the maximum support is 4096 bytes. When secret status is `Disabled`, this field will not update anymore.",
			},
		},
	}
}

func resourceTencentCloudSsmSecretVersionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_secret_version.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ssmService := SsmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	param := make(map[string]interface{})
	param["secret_name"] = d.Get("secret_name").(string)
	param["version_id"] = d.Get("version_id").(string)
	if v, ok := d.GetOk("secret_binary"); ok {
		param["secret_binary"] = v.(string)
	}
	if v, ok := d.GetOk("secret_string"); ok {
		param["secret_string"] = v.(string)
	}

	var outErr, inErr error
	var secretName, versionId string
	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		secretName, versionId, inErr = ssmService.PutSecretValue(ctx, param)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}
	d.SetId(strings.Join([]string{secretName, versionId}, FILED_SP))

	return resourceTencentCloudSsmSecretVersionRead(d, meta)
}

func resourceTencentCloudSsmSecretVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_secret_version.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ssmService := SsmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("SSM secret version id can't read, id is borken, id is %s", d.Id())
	}
	secretName := ids[0]
	versionId := ids[1]

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
		_ = d.Set("secret_name", secretVersionInfo.secretName)
		_ = d.Set("version_id", secretVersionInfo.versionId)
		_ = d.Set("secret_binary", secretVersionInfo.secretBinary)
		_ = d.Set("secret_string", secretVersionInfo.secretString)
	}

	return nil
}

func resourceTencentCloudSsmSecretVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_secret_version.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ssmService := SsmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("SSM secret version id can't read, id is borken, id is %s", d.Id())
	}
	secretName := ids[0]
	versionId := ids[1]

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
		d.Partial(true)

		param := make(map[string]interface{})
		param["secret_name"] = secretName
		param["version_id"] = versionId
		if v, ok := d.GetOk("secret_binary"); ok {
			param["secret_binary"] = v.(string)
		}
		if v, ok := d.GetOk("secret_string"); ok {
			param["secret_string"] = v.(string)
		}

		if d.HasChange("secret_binary") || d.HasChange("secret_string") {
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
			d.SetPartial("secret_binary")
			d.SetPartial("secret_string")
		}

		d.Partial(false)
	}

	return resourceTencentCloudSsmSecretVersionRead(d, meta)
}

func resourceTencentCloudSsmSecretVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_secret_version.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ssmService := SsmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("SSM secret version id can't read, id is borken, id is %s", d.Id())
	}
	secretName := ids[0]
	versionId := ids[1]

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := ssmService.DeleteSecretVersion(ctx, secretName, versionId)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete SSM secret version failed, reason:%+v", logId, err)
		return err
	}

	return resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, e := ssmService.DescribeSecretVersion(ctx, secretName, versionId)
		if e != nil {
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok && sdkError.Code == "ResourceNotFound" {
				return nil
			}
			return retryError(err)
		}

		return resource.RetryableError(fmt.Errorf("delete fail"))
	})
}
