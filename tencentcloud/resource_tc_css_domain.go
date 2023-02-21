/*
Provides a resource to create a css domain

Example Usage

```hcl
resource "tencentcloud_css_domain" "domain" {
  domain_name = "iac-tf.cloud"
  domain_type = 0
  play_type = 1
  is_delay_live = 0
  is_mini_program_live = 0
  verify_owner_type = "dbCheck"
}
```

Import

css domain can be imported using the id, e.g.

```
terraform import tencentcloud_css_domain.domain domain_name
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCssDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssDomainCreate,
		Read:   resourceTencentCloudCssDomainRead,
		Update: resourceTencentCloudCssDomainUpdate,
		Delete: resourceTencentCloudCssDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain Name.",
			},

			"domain_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Domain type: `0`: push stream. `1`: playback.",
			},

			"play_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Play Type. This parameter is valid only if `DomainType` is 1. Available values: `1`: in Mainland China. `2`: global. `3`: outside Mainland China.",
			},

			"is_delay_live": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether it is LCB: `0`: LVB. `1`: LCB.",
			},

			"is_mini_program_live": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "`0`: LVB. `1`: LVB on Mini Program. Note: this field may return null, indicating that no valid values can be obtained.",
			},

			"verify_owner_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Domain name attribution verification type. `dnsCheck`, `fileCheck`, `dbCheck`. The default is `dbCheck`.",
			},
		},
	}
}

func resourceTencentCloudCssDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_domain.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = css.NewAddLiveDomainRequest()
		domainName string
	)
	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
		request.DomainName = helper.String(domainName)
	}

	if v, ok := d.GetOkExists("domain_type"); ok {
		request.DomainType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("play_type"); ok {
		request.PlayType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("is_delay_live"); ok {
		request.IsDelayLive = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_mini_program_live"); ok {
		request.IsMiniProgramLive = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("verify_owner_type"); ok {
		request.VerifyOwnerType = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().AddLiveDomain(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css domain failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(domainName)

	return resourceTencentCloudCssDomainRead(d, meta)
}

func resourceTencentCloudCssDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_domain.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	domainName := d.Id()

	domain, err := service.DescribeCssDomainById(ctx, domainName)
	if err != nil {
		return err
	}

	if domain == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssDomain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if domain.Name != nil {
		_ = d.Set("domain_name", domain.Name)
	}

	if domain.Type != nil {
		_ = d.Set("domain_type", domain.Type)
	}

	if domain.PlayType != nil {
		_ = d.Set("play_type", domain.PlayType)
	}

	if domain.IsDelayLive != nil {
		_ = d.Set("is_delay_live", domain.IsDelayLive)
	}

	if domain.IsMiniProgramLive != nil {
		_ = d.Set("is_mini_program_live", domain.IsMiniProgramLive)
	}

	return nil
}

func resourceTencentCloudCssDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_domain.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := css.NewModifyLivePlayDomainRequest()

	request.DomainName = helper.String(d.Id())

	immutableArgs := []string{"domain_name", "domain_type", "play_type", "is_delay_live", "is_mini_program_live", "verify_owner_type", "domain_infos"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if v, ok := d.GetOkExists("domain_type"); ok {
		dtype := v.(int)
		if v.(int) != CSS_DOMAIN_TYPE_PLAY_BACK {
			return fmt.Errorf("argument domain_type:[%v] does not support modify(ModifyLivePlayDomain)", dtype)
		}
	}

	if d.HasChange("domain_name") {
		if v, ok := d.GetOk("domain_name"); ok {
			request.DomainName = helper.String(v.(string))
		}
	}

	if d.HasChange("play_type") {
		if v, _ := d.GetOk("play_type"); v != nil {
			request.PlayType = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().ModifyLivePlayDomain(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update css domain failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCssDomainRead(d, meta)
}

func resourceTencentCloudCssDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_domain.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	var domainType *uint64

	if v, ok := d.GetOkExists("domain_type"); ok {
		domainType = helper.IntUint64(v.(int))
	}

	if err := service.DeleteCssDomainById(ctx, helper.String(d.Id()), domainType); err != nil {
		return err
	}

	return nil
}
