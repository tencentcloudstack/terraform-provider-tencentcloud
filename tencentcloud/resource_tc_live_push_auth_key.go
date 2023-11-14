/*
Provides a resource to create a live push_auth_key

Example Usage

```hcl
resource "tencentcloud_live_push_auth_key" "push_auth_key" {
  domain_name = "5000.livepush.myqcloud.com"
  enable = 0
  master_auth_key = "xx"
  backup_auth_key = "xx"
  auth_delta = 60
}
```

Import

live push_auth_key can be imported using the id, e.g.

```
terraform import tencentcloud_live_push_auth_key.push_auth_key push_auth_key_id
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

func resourceTencentCloudLivePushAuthKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLivePushAuthKeyCreate,
		Read:   resourceTencentCloudLivePushAuthKeyRead,
		Update: resourceTencentCloudLivePushAuthKeyUpdate,
		Delete: resourceTencentCloudLivePushAuthKeyDelete,
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

func resourceTencentCloudLivePushAuthKeyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_push_auth_key.create")()
	defer inconsistentCheck(d, meta)()

	var domainName string
	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
	}

	d.SetId(domainName)

	return resourceTencentCloudLivePushAuthKeyUpdate(d, meta)
}

func resourceTencentCloudLivePushAuthKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_push_auth_key.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	pushAuthKeyId := d.Id()

	pushAuthKey, err := service.DescribeLivePushAuthKeyById(ctx, domainName)
	if err != nil {
		return err
	}

	if pushAuthKey == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LivePushAuthKey` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if pushAuthKey.DomainName != nil {
		_ = d.Set("domain_name", pushAuthKey.DomainName)
	}

	if pushAuthKey.Enable != nil {
		_ = d.Set("enable", pushAuthKey.Enable)
	}

	if pushAuthKey.MasterAuthKey != nil {
		_ = d.Set("master_auth_key", pushAuthKey.MasterAuthKey)
	}

	if pushAuthKey.BackupAuthKey != nil {
		_ = d.Set("backup_auth_key", pushAuthKey.BackupAuthKey)
	}

	if pushAuthKey.AuthDelta != nil {
		_ = d.Set("auth_delta", pushAuthKey.AuthDelta)
	}

	return nil
}

func resourceTencentCloudLivePushAuthKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_push_auth_key.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := live.NewModifyLivePushAuthKeyRequest()

	pushAuthKeyId := d.Id()

	request.DomainName = &domainName

	immutableArgs := []string{"domain_name", "enable", "master_auth_key", "backup_auth_key", "auth_delta"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().ModifyLivePushAuthKey(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update live pushAuthKey failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLivePushAuthKeyRead(d, meta)
}

func resourceTencentCloudLivePushAuthKeyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_push_auth_key.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
