/*
Provides a resource to create a dasb bind_device_account_private_key

Example Usage

```hcl
resource "tencentcloud_dasb_bind_device_account_private_key" "example" {
  device_account_id    = 16
  private_key          = "MIICXAIBAAKBgQCqGKukO1De7zhZj6+H0qtjTkVxwTCpvKe4eCZ0FPqri0cb2JZfXJ/DgYSF6vUpwmJG8wVQZKjeGcjDOL5UlsuusFncCzWBQ7RKNUSesmQRMSGkVb1/3j+skZ6UtW+5u09lHNsj6tQ51s1SPrCBkedbNf0Tp0GbMJDyR4e9T04ZZwIDAQABAoGAFijko56+qGyN8M0RVyaRAXz++xTqHBLh"
  private_key_password = "TerraformPassword"
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDasbBindDeviceAccountPrivateKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDasbBindDeviceAccountPrivateKeyCreate,
		Read:   resourceTencentCloudDasbBindDeviceAccountPrivateKeyRead,
		Delete: resourceTencentCloudDasbBindDeviceAccountPrivateKeyDelete,

		Schema: map[string]*schema.Schema{
			"device_account_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Host account ID.",
			},
			"private_key": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Host account private key, the latest length is 128 bytes, the maximum length is 8192 bytes.",
			},
			"private_key_password": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Host account private key password, maximum length 256 bytes.",
			},
		},
	}
}

func resourceTencentCloudDasbBindDeviceAccountPrivateKeyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_bind_device_account_private_key.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId           = getLogId(contextNil)
		request         = dasb.NewBindDeviceAccountPrivateKeyRequest()
		deviceAccountId string
	)

	if v, ok := d.GetOkExists("device_account_id"); ok {
		request.Id = helper.IntUint64(v.(int))
		deviceAccountId = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("private_key"); ok {
		request.PrivateKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("private_key_password"); ok {
		request.PrivateKeyPassword = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDasbClient().BindDeviceAccountPrivateKey(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dasb bindDeviceAccountPrivateKey failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(deviceAccountId)

	return resourceTencentCloudDasbBindDeviceAccountPrivateKeyRead(d, meta)
}

func resourceTencentCloudDasbBindDeviceAccountPrivateKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_bind_device_account_private_key.read")()
	defer inconsistentCheck(d, meta)()

	if v, ok := d.GetOkExists("device_account_id"); ok {
		_ = d.Set("device_account_id", v.(int))
	}

	if v, ok := d.GetOk("private_key"); ok {
		_ = d.Set("private_key", v.(string))
	}

	if v, ok := d.GetOk("private_key_password"); ok {
		_ = d.Set("private_key_password", v.(string))
	}

	return nil
}

func resourceTencentCloudDasbBindDeviceAccountPrivateKeyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_bind_device_account_private_key.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId           = getLogId(contextNil)
		ctx             = context.WithValue(context.TODO(), logIdKey, logId)
		service         = DasbService{client: meta.(*TencentCloudClient).apiV3Conn}
		deviceAccountId = d.Id()
	)

	if err := service.DeleteDasbBindDeviceAccountPrivateKeyById(ctx, deviceAccountId); err != nil {
		return err
	}

	return nil
}
