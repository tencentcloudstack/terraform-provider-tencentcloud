/*
Provides a resource to create a css push_auth_key_config

Example Usage

```hcl
resource "tencentcloud_css_push_auth_key_config" "push_auth_key_config" {
  domain_name = "your_push_domain_name"
  enable = 1
  master_auth_key = "testmasterkey"
  backup_auth_key = "testbackkey"
  auth_delta = 1800
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
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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

	var domainName string
	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
	}

	d.SetId(domainName)

	return resourceTencentCloudCssPushAuthKeyConfigUpdate(d, meta)
}

func resourceTencentCloudCssPushAuthKeyConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_push_auth_key_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	domainName := d.Id()

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

	request.DomainName = helper.String(d.Id())

	if v, ok := d.GetOkExists("enable"); ok {
		request.Enable = helper.IntInt64(v.(int))
	}

	if d.HasChange("master_auth_key") {
		if v, ok := d.GetOk("master_auth_key"); ok {
			request.MasterAuthKey = helper.String(v.(string))
		}
	}

	if d.HasChange("backup_auth_key") {
		if v, ok := d.GetOk("backup_auth_key"); ok {
			request.BackupAuthKey = helper.String(v.(string))
		}
	}

	if d.HasChange("auth_delta") {
		if v, _ := d.GetOk("auth_delta"); v != nil {
			request.AuthDelta = helper.IntUint64(v.(int))
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
	//donothing
	return nil
}
