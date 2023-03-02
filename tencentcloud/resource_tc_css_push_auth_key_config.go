/*
Provides a resource to create a css push_auth_key_config

Example Usage

```hcl
resource "tencentcloud_css_push_auth_key_config" "push_auth_key_config" {
  domain_name = "5000.livepush.myqcloud.com"
  enable = 0
  master_auth_key = "xx"
  backup_auth_key = "xx"
  auth_delta = 60
}
```

Import

css push_auth_key_config can be imported using the id, e.g.

```
terraform import tencentcloud_css_push_auth_key_config.push_auth_key_config push_auth_key_config_id
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

func resourceTencentCloudCssPushAuthKeyConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssPushAuthKeyConfigCreate,
		Read:   resourceTencentCloudCssPushAuthKeyConfigRead,
		Update: resourceTencentCloudCssPushAuthKeyConfigUpdate,
		Delete: resourceTencentCloudCssPushAuthKeyConfigDelete,
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

			"master_auth_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Primary authentication key. No transfer means that the current value is not modified.",
			},

			"backup_auth_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Standby authentication key. No transfer means that the current value is not modified.",
			},

			"auth_delta": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Valid time, unit: second.",
			},
		},
	}
}

func resourceTencentCloudCssPushAuthKeyConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_push_auth_key_config.create")()
	defer inconsistentCheck(d, meta)()

	d.SetId(helper.UInt64ToStr(domainName))

	return resourceTencentCloudCssPushAuthKeyConfigRead(d, meta)
}

func resourceTencentCloudCssPushAuthKeyConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_push_auth_key_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	pushAuthKeyConfigId := d.Id()

	pushAuthKeyConfig, err := service.DescribeCssPushAuthKeyConfigById(ctx, domainName)
	if err != nil {
		return err
	}

	if pushAuthKeyConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssPushAuthKeyConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if pushAuthKeyConfig.DomainName != nil {
		_ = d.Set("domain_name", pushAuthKeyConfig.DomainName)
	}

	if pushAuthKeyConfig.Enable != nil {
		_ = d.Set("enable", pushAuthKeyConfig.Enable)
	}

	if pushAuthKeyConfig.MasterAuthKey != nil {
		_ = d.Set("master_auth_key", pushAuthKeyConfig.MasterAuthKey)
	}

	if pushAuthKeyConfig.BackupAuthKey != nil {
		_ = d.Set("backup_auth_key", pushAuthKeyConfig.BackupAuthKey)
	}

	if pushAuthKeyConfig.AuthDelta != nil {
		_ = d.Set("auth_delta", pushAuthKeyConfig.AuthDelta)
	}

	return nil
}

func resourceTencentCloudCssPushAuthKeyConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_push_auth_key_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := css.NewModifyLivePushAuthKeyRequest()

	pushAuthKeyConfigId := d.Id()

	request.DomainName = &domainName

	immutableArgs := []string{"domain_name", "enable", "master_auth_key", "backup_auth_key", "auth_delta"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().ModifyLivePushAuthKey(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update css pushAuthKeyConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCssPushAuthKeyConfigRead(d, meta)
}

func resourceTencentCloudCssPushAuthKeyConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_push_auth_key_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
