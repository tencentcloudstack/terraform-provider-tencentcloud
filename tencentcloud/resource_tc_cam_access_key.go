/*
Provides a resource to create a cam access_key

Example Usage

```hcl
resource "tencentcloud_cam_access_key" "access_key" {
  target_uin = 100033690181
}
```
Update
```hcl
resource "tencentcloud_cam_access_key" "access_key" {
  target_uin = 100033690181
  status = "Inactive"
}
```
Encrypted
```hcl
resource "tencentcloud_cam_access_key" "access_key" {
  target_uin = 100033690181
  pgp_key = "keybase:some_person_that_exists"
}
```
Import

cam access_key can be imported using the id, e.g.

```
terraform import tencentcloud_cam_access_key.access_key access_key_id
```
*/
package tencentcloud

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCamAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamAccessKeyCreate,
		Read:   resourceTencentCloudCamAccessKeyRead,
		Update: resourceTencentCloudCamAccessKeyUpdate,
		Delete: resourceTencentCloudCamAccessKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"target_uin": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Specify user Uin, if not filled, the access key is created for the current user by default.",
			},
			"access_key": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Access_key is the access key identification, required when updating.",
			},
			"pgp_key": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Either a base-64 encoded PGP public key, or a keybase username in the form keybase:some_person_that_exists, for use in the encrypted_secret output attribute. If providing a base-64 encoded PGP public key, make sure to provide the \"raw\" version and not the \"armored\" one (e.g. avoid passing the -a option to gpg --export).",
			},
			"secret_access_key": {
				Computed:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
				Description: "Access key (key is only visible when created, please keep it properly).",
			},
			"encrypted_secret_access_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Encrypted secret, base64 encoded, if pgp_key was specified. This attribute is not available for imported resources. The encrypted secret may be decrypted using the command line, for example: terraform output -raw encrypted_secret | base64 --decode | keybase pgp decrypt.",
			},
			"key_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Fingerprint of the PGP key used to encrypt the secret. This attribute is not available for imported resources.",
			},
			"status": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Key status, activated (Active) or inactive (Inactive), required when updating.",
			},
		},
	}
}

func resourceTencentCloudCamAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_access_key.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cam.NewCreateAccessKeyRequest()
		response  = cam.NewCreateAccessKeyResponse()
		targetUin int64
	)
	if v, ok := d.GetOkExists("target_uin"); ok {
		targetUin = int64(v.(int))
		request.TargetUin = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().CreateAccessKey(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam AccessKey failed, reason:%+v", logId, err)
		return err
	}
	if response == nil || response.Response == nil || response.Response.AccessKey == nil || response.Response.AccessKey.SecretAccessKey == nil {
		return fmt.Errorf("CAM AccessKey id is nil")
	}
	d.SetId(helper.Int64ToStr(targetUin) + FILED_SP + *response.Response.AccessKey.AccessKeyId)

	if v, ok := d.GetOk("pgp_key"); ok {
		pgpKey := v.(string)

		encryptionKey, err := retrieveGPGKey(pgpKey)
		if err != nil {
			return fmt.Errorf("creating IAM Access Key (%v): %s", targetUin, err)
		}
		fingerprint, encrypted, err := encryptValue(encryptionKey, *response.Response.AccessKey.SecretAccessKey, "IAM Access Key Secret")
		if err != nil {
			return fmt.Errorf("creating IAM Access Key (%v): %s", targetUin, err)
		}

		_ = d.Set("key_fingerprint", fingerprint)
		_ = d.Set("encrypted_secret_access_key", encrypted)

	} else {
		_ = d.Set("secret_access_key", response.Response.AccessKey.SecretAccessKey)
	}

	return resourceTencentCloudCamAccessKeyRead(d, meta)
}

func resourceTencentCloudCamAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_access_key.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	uin := idSplit[0]
	accessKey := idSplit[1]

	AccessKey, err := service.DescribeCamAccessKeyById(ctx, helper.StrToUInt64(uin), accessKey)
	if err != nil {
		return err
	}

	if AccessKey == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamAccessKey` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if AccessKey.AccessKeyId != nil {
		_ = d.Set("access_key", AccessKey.AccessKeyId)
	}
	if AccessKey.Status != nil {
		_ = d.Set("status", AccessKey.Status)
	}
	_ = d.Set("target_uin", helper.StrToUInt64(uin))

	return nil
}

func resourceTencentCloudCamAccessKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_access_key.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cam.NewUpdateAccessKeyRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	uin := idSplit[0]
	accessKey := idSplit[1]
	request.TargetUin = helper.StrToUint64Point(uin)
	request.AccessKeyId = &accessKey

	if d.HasChange("status") {
		if v, ok := d.GetOkExists("status"); ok {
			request.Status = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().UpdateAccessKey(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cam AccessKey failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCamAccessKeyRead(d, meta)
}

func resourceTencentCloudCamAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_access_key.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	uin := idSplit[0]
	accessKey := idSplit[1]

	if err := service.DeleteCamAccessKeyById(ctx, uin, accessKey); err != nil {
		return err
	}

	return nil
}

func retrieveGPGKey(pgpKey string) (string, error) {

	encryptionKey := pgpKey
	if strings.HasPrefix(pgpKey, kbPrefix) {
		publicKeys, err := FetchKeybasePubkeys([]string{pgpKey})
		if err != nil {
			return "", fmt.Errorf("retrieving Public Key (%s): %w", pgpKey, err)
		}
		encryptionKey = publicKeys[pgpKey]
	}

	return encryptionKey, nil
}

func encryptValue(encryptionKey, value, description string) (string, string, error) {
	fingerprints, encryptedValue, err :=
		EncryptShares([][]byte{[]byte(value)}, []string{encryptionKey})
	if err != nil {
		return "", "", fmt.Errorf("encrypting %s: %w", description, err)
	}

	return fingerprints[0], base64.StdEncoding.EncodeToString(encryptedValue[0]), nil
}
