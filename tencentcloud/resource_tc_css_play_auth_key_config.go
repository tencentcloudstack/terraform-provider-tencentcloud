/*
Provides a resource to create a css play_auth_key_config

Example Usage

```hcl
resource "tencentcloud_css_play_auth_key_config" "play_auth_key_config" {
  domain_name = "5000.livepush.myqcloud.com"
  enable = 1
  auth_key = "xx"
  auth_delta = 60
  auth_back_key = "xx"
}
```

Import

css play_auth_key_config can be imported using the id, e.g.

```
terraform import tencentcloud_css_play_auth_key_config.play_auth_key_config play_auth_key_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCssPlayAuthKeyConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssPlayAuthKeyConfigCreate,
		Read:   resourceTencentCloudCssPlayAuthKeyConfigRead,
		Update: resourceTencentCloudCssPlayAuthKeyConfigUpdate,
		Delete: resourceTencentCloudCssPlayAuthKeyConfigDelete,
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

func resourceTencentCloudCssPlayAuthKeyConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_play_auth_key_config.create")()
	defer inconsistentCheck(d, meta)()

	d.SetId(helper.UInt64ToStr(domainName))

	return resourceTencentCloudCssPlayAuthKeyConfigRead(d, meta)
}

func resourceTencentCloudCssPlayAuthKeyConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_play_auth_key_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	playAuthKeyConfigId := d.Id()

	playAuthKeyConfig, err := service.DescribeCssPlayAuthKeyConfigById(ctx, domainName)
	if err != nil {
		return err
	}

	if playAuthKeyConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssPlayAuthKeyConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if playAuthKeyConfig.DomainName != nil {
		_ = d.Set("domain_name", playAuthKeyConfig.DomainName)
	}

	if playAuthKeyConfig.Enable != nil {
		_ = d.Set("enable", playAuthKeyConfig.Enable)
	}

	if playAuthKeyConfig.AuthKey != nil {
		_ = d.Set("auth_key", playAuthKeyConfig.AuthKey)
	}

	if playAuthKeyConfig.AuthDelta != nil {
		_ = d.Set("auth_delta", playAuthKeyConfig.AuthDelta)
	}

	if playAuthKeyConfig.AuthBackKey != nil {
		_ = d.Set("auth_back_key", playAuthKeyConfig.AuthBackKey)
	}

	return nil
}

func resourceTencentCloudCssPlayAuthKeyConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_play_auth_key_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := css.NewModifyLivePlayAuthKeyRequest()

	playAuthKeyConfigId := d.Id()

	request.DomainName = &domainName

	immutableArgs := []string{"domain_name", "enable", "auth_key", "auth_delta", "auth_back_key"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().ModifyLivePlayAuthKey(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update css playAuthKeyConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCssPlayAuthKeyConfigRead(d, meta)
}

func resourceTencentCloudCssPlayAuthKeyConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_play_auth_key_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
