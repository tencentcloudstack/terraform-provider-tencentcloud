/*
Provides a resource to configure(enable/disable) the css domain.

Example Usage

```hcl
resource "tencentcloud_css_domain_config" "enable_domain" {
  domain_name = "your_domain_name"
  enable_domain = true
}
```

```hcl
resource "tencentcloud_css_domain_config" "forbid_domain" {
  domain_name = "your_domain_name"
  enable_domain = false
}
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

func resourceTencentCloudCssDomainConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssDomainConfigCreate,
		Read:   resourceTencentCloudCssDomainConfigRead,
		Update: resourceTencentCloudCssDomainConfigUpdate,
		Delete: resourceTencentCloudCssDomainConfigDelete,
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain Name.",
			},

			"enable_domain": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Switch. true: enable the specified domain, false: disable the specified domain.",
			},
		},
	}
}

func resourceTencentCloudCssDomainConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_domain_config.create")()
	defer inconsistentCheck(d, meta)()
	var name string
	if v, ok := d.GetOk("domain_name"); ok {
		name = v.(string)
	}

	d.SetId(name)
	return resourceTencentCloudRedisBackupConfigUpdate(d, meta)
}

func resourceTencentCloudCssDomainConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_domain_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	domainName := d.Id()

	domainConfig, err := service.DescribeCssDomainById(ctx, domainName)
	if err != nil {
		return err
	}

	if domainConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssDomainConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if domainConfig.Name != nil {
		_ = d.Set("domain_name", domainConfig.Name)
	}

	if domainConfig.Status != nil {
		var enable *bool
		status := helper.UInt64Int64(*domainConfig.Status)

		switch *status {
		case CSS_DOMAIN_STATUS_ACTIVATED:
			enable = helper.Bool(true)
		case CSS_DOMAIN_STATUS_DEACTIVATED:
			enable = helper.Bool(false)
		default:
		}
		_ = d.Set("enable_domain", enable)
	}

	return nil
}

func resourceTencentCloudCssDomainConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_domain_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		enable        *bool
		enableRequest = css.NewEnableLiveDomainRequest()
		forbidRequest = css.NewForbidLiveDomainRequest()
	)

	if d.HasChange("enable_domain") {
		if v, _ := d.GetOk("enable_domain"); v != nil {
			enable = helper.Bool(v.(bool))
		}
	}

	if enable != nil {
		if *enable {
			enableRequest.DomainName = helper.String(d.Id())
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().EnableLiveDomain(enableRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, enableRequest.GetAction(), enableRequest.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s EnableLiveDomain failed, reason:%+v", logId, err)
				return err
			}
		} else {
			forbidRequest.DomainName = helper.String(d.Id())
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().ForbidLiveDomain(forbidRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, forbidRequest.GetAction(), forbidRequest.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s ForbidLiveDomain failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudCssDomainConfigRead(d, meta)
}

func resourceTencentCloudCssDomainConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_domain_config.delete")()
	defer inconsistentCheck(d, meta)()
	// do nothing
	return nil
}
