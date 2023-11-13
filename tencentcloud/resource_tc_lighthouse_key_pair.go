/*
Provides a resource to create a lighthouse key_pair

Example Usage

```hcl
resource "tencentcloud_lighthouse_key_pair" "key_pair" {
  key_name = "key_name_test"
  public_key = "public key content"
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudLighthouseKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseKeyPairCreate,
		Read:   resourceTencentCloudLighthouseKeyPairRead,
		Update: resourceTencentCloudLighthouseKeyPairUpdate,
		Delete: resourceTencentCloudLighthouseKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Key pair name, which can contain up to 25 digits, letters, and underscores.",
			},

			"public_key": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Public key content of the key pair, which is in the OpenSSH RSA format.",
			},
		},
	}
}

func resourceTencentCloudLighthouseKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_key_pair.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		createKeyPairRequest  = lighthouse.NewCreateKeyPairRequest()
		createKeyPairResponse = lighthouse.NewCreateKeyPairResponse()
	)
	if v, ok := d.GetOk("key_name"); ok {
		request.KeyName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("public_key"); ok {
		request.PublicKey = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().CreateKeyPair(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create lighthouse keyPair failed, reason:%+v", logId, err)
		return err
	}

	keyId = *response.Response.KeyId
	d.SetId(keyId)

	return resourceTencentCloudLighthouseKeyPairRead(d, meta)
}

func resourceTencentCloudLighthouseKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_key_pair.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	keyPairId := d.Id()

	keyPair, err := service.DescribeLighthouseKeyPairById(ctx, keyId)
	if err != nil {
		return err
	}

	if keyPair == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseKeyPair` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if keyPair.KeyName != nil {
		_ = d.Set("key_name", keyPair.KeyName)
	}

	if keyPair.PublicKey != nil {
		_ = d.Set("public_key", keyPair.PublicKey)
	}

	return nil
}

func resourceTencentCloudLighthouseKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_key_pair.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"key_name", "public_key"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudLighthouseKeyPairRead(d, meta)
}

func resourceTencentCloudLighthouseKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_key_pair.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}
	keyPairId := d.Id()

	if err := service.DeleteLighthouseKeyPairById(ctx, keyId); err != nil {
		return err
	}

	return nil
}
