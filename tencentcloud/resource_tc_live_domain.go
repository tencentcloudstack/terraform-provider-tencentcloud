/*
Provides a resource to create a live domain

Example Usage

```hcl
resource "tencentcloud_live_domain" "domain" {
  domain_name = ""
  domain_type =
  play_type =
  is_delay_live =
  is_mini_program_live =
  verify_owner_type = ""
}
```

Import

live domain can be imported using the id, e.g.

```
terraform import tencentcloud_live_domain.domain domain_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudLiveDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLiveDomainCreate,
		Read:   resourceTencentCloudLiveDomainRead,
		Update: resourceTencentCloudLiveDomainUpdate,
		Delete: resourceTencentCloudLiveDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain name.",
			},

			"domain_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Domain name type.0: push domain name.1: playback domain name.",
			},

			"play_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Pull domain name type:1: Mainland China.2: global.3: outside Mainland China.Default value: 1.",
			},

			"is_delay_live": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether it is LCB:0: LVB,1: LCB.Default value: 0.",
			},

			"is_mini_program_live": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether it is LVB on Mini Program.0: LVB.1: LVB on Mini Program.Default value: 0.",
			},

			"verify_owner_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The domain verification type.Valid values (the value of this parameter must be the same as `VerifyType` of the `AuthenticateDomainOwner` API):dnsCheck: Check immediately whether the verification DNS record has been added successfully. If so, record this verification result.fileCheck: Check immediately whether the verification HTML file has been uploaded successfully. If so, record this verification result.dbCheck: Check whether the domain has already been verified.If you do not pass a value, `dbCheck` will be used.",
			},
		},
	}
}

func resourceTencentCloudLiveDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_domain.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = live.NewAddLiveDomainRequest()
		response   = live.NewAddLiveDomainResponse()
		domainName string
	)
	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
		request.DomainName = helper.String(v.(string))
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().AddLiveDomain(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create live domain failed, reason:%+v", logId, err)
		return err
	}

	domainName = *response.Response.DomainName
	d.SetId(domainName)

	return resourceTencentCloudLiveDomainRead(d, meta)
}

func resourceTencentCloudLiveDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_domain.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	domainId := d.Id()

	domain, err := service.DescribeLiveDomainById(ctx, domainName)
	if err != nil {
		return err
	}

	if domain == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LiveDomain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if domain.DomainName != nil {
		_ = d.Set("domain_name", domain.DomainName)
	}

	if domain.DomainType != nil {
		_ = d.Set("domain_type", domain.DomainType)
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

	if domain.VerifyOwnerType != nil {
		_ = d.Set("verify_owner_type", domain.VerifyOwnerType)
	}

	return nil
}

func resourceTencentCloudLiveDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_domain.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := live.NewModifyLivePlayDomainRequest()

	domainId := d.Id()

	request.DomainName = &domainName

	immutableArgs := []string{"domain_name", "domain_type", "play_type", "is_delay_live", "is_mini_program_live", "verify_owner_type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("domain_name") {
		if v, ok := d.GetOk("domain_name"); ok {
			request.DomainName = helper.String(v.(string))
		}
	}

	if d.HasChange("play_type") {
		if v, ok := d.GetOkExists("play_type"); ok {
			request.PlayType = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().ModifyLivePlayDomain(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update live domain failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLiveDomainRead(d, meta)
}

func resourceTencentCloudLiveDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_domain.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}
	domainId := d.Id()

	if err := service.DeleteLiveDomainById(ctx, domainName); err != nil {
		return err
	}

	return nil
}
