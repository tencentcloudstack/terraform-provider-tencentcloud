/*
Provides a resource to create a lighthouse key_pair

Example Usage

```hcl
resource "tencentcloud_lighthouse_key_pair" "key_pair" {
  key_name = "key_name_test"
}
```

Import

lighthouse key_pair can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_key_pair.key_pair key_pair_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudLighthouseKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseKeyPairCreate,
		Read:   resourceTencentCloudLighthouseKeyPairRead,
		Delete: resourceTencentCloudLighthouseKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Key pair name, which can contain up to 25 digits, letters, and underscores.",
			},
			"public_key": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Public key content of the key pair, which is in the OpenSSH RSA format.",
			},
			"private_key": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Key to private key.",
			},
			"created_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Creation time. Expressed according to the ISO8601 standard, and using UTC time. Format: YYYY-MM-DDThh:mm:ssZ.",
			},
		},
	}
}

func createKeyPair(ctx context.Context, d *schema.ResourceData, meta interface{}) (keyId string, err error) {
	logId := getLogId(ctx)
	request := lighthouse.NewCreateKeyPairRequest()
	response := lighthouse.NewCreateKeyPairResponse()
	request.KeyName = helper.String(d.Get("key_name").(string))
	innerErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().CreateKeyPair(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if innerErr != nil {
		log.Printf("[CRITAL]%s create lighthouse keyPair failed, reason:%+v", logId, err)
		err = innerErr
		return
	}
	if response == nil || response.Response == nil || response.Response.KeyPair == nil {
		err = fmt.Errorf("Response is nil")
		return
	}
	keyId = *response.Response.KeyPair.KeyId
	return
}

func createKeyPairByImportPublicKey(ctx context.Context, d *schema.ResourceData, meta interface{}) (keyId string, err error) {
	logId := getLogId(ctx)
	request := lighthouse.NewImportKeyPairRequest()
	response := lighthouse.NewImportKeyPairResponse()
	request.KeyName = helper.String(d.Get("key_name").(string))
	request.PublicKey = helper.String(d.Get("public_key").(string))
	innerErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ImportKeyPair(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if innerErr != nil {
		log.Printf("[CRITAL]%s create lighthouse keyPair by import failed, reason:%+v", logId, err)
		err = innerErr
		return
	}
	if response == nil || response.Response == nil {
		err = fmt.Errorf("Response is nil")
		return
	}

	keyId = *response.Response.KeyId
	return
}

func resourceTencentCloudLighthouseKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_key_pair.create")()
	defer inconsistentCheck(d, meta)()

	var (
		keyId string
		err   error
	)
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	if _, ok := d.GetOk("public_key"); ok {
		keyId, err = createKeyPairByImportPublicKey(ctx, d, meta)
	} else {
		keyId, err = createKeyPair(ctx, d, meta)
	}
	if err != nil {
		return err
	}

	d.SetId(keyId)
	return resourceTencentCloudLighthouseKeyPairRead(d, meta)
}

func resourceTencentCloudLighthouseKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_key_pair.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	keyId := d.Id()

	keyPair, err := service.DescribeLighthouseKeyPairById(ctx, keyId)
	if err != nil {
		return err
	}

	if keyPair == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseKeyPair` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("key_name", keyPair.KeyName)
	_ = d.Set("public_key", keyPair.PublicKey)
	_ = d.Set("private_key", keyPair.PrivateKey)
	_ = d.Set("created_time", keyPair.CreatedTime)

	return nil
}

func resourceTencentCloudLighthouseKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_key_pair.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}
	keyId := d.Id()

	if err := service.DeleteLighthouseKeyPairById(ctx, keyId); err != nil {
		return err
	}

	return nil
}
