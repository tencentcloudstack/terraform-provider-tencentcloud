/*
Provides a resource to create a live play_auth_key

Example Usage

```hcl
resource "tencentcloud_live_play_auth_key" "play_auth_key" {
  domain_name = "5000.livepush.myqcloud.com"
  enable = 1
  auth_key = "xx"
  auth_delta = 60
  auth_back_key = "xx"
}
```

Import

live play_auth_key can be imported using the id, e.g.

```
terraform import tencentcloud_live_play_auth_key.play_auth_key play_auth_key_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"log"
)

func resourceTencentCloudLivePlayAuthKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLivePlayAuthKeyCreate,
		Read:   resourceTencentCloudLivePlayAuthKeyRead,
		Update: resourceTencentCloudLivePlayAuthKeyUpdate,
		Delete: resourceTencentCloudLivePlayAuthKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain Name.",
			},

			"enable": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Enable or not, 0: Close, 1: Enable. No transfer means that the current value is not modified.",
			},

			"auth_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Authentication key. No transfer means that the current value is not modified.",
			},

			"auth_delta": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Valid time, unit: second. No transfer means that the current value is not modified.",
			},

			"auth_back_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Alternate key for authentication. No transfer means that the current value is not modified.",
			},
		},
	}
}

func resourceTencentCloudLivePlayAuthKeyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_play_auth_key.create")()
	defer inconsistentCheck(d, meta)()

	var domainName string
	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
	}

	d.SetId(domainName)

	return resourceTencentCloudLivePlayAuthKeyUpdate(d, meta)
}

func resourceTencentCloudLivePlayAuthKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_play_auth_key.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	playAuthKeyId := d.Id()

	playAuthKey, err := service.DescribeLivePlayAuthKeyById(ctx, domainName)
	if err != nil {
		return err
	}

	if playAuthKey == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LivePlayAuthKey` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if playAuthKey.DomainName != nil {
		_ = d.Set("domain_name", playAuthKey.DomainName)
	}

	if playAuthKey.Enable != nil {
		_ = d.Set("enable", playAuthKey.Enable)
	}

	if playAuthKey.AuthKey != nil {
		_ = d.Set("auth_key", playAuthKey.AuthKey)
	}

	if playAuthKey.AuthDelta != nil {
		_ = d.Set("auth_delta", playAuthKey.AuthDelta)
	}

	if playAuthKey.AuthBackKey != nil {
		_ = d.Set("auth_back_key", playAuthKey.AuthBackKey)
	}

	return nil
}

func resourceTencentCloudLivePlayAuthKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_play_auth_key.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := live.NewModifyLivePlayAuthKeyRequest()

	playAuthKeyId := d.Id()

	request.DomainName = &domainName

	immutableArgs := []string{"domain_name", "enable", "auth_key", "auth_delta", "auth_back_key"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().ModifyLivePlayAuthKey(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update live playAuthKey failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLivePlayAuthKeyRead(d, meta)
}

func resourceTencentCloudLivePlayAuthKeyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_play_auth_key.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
