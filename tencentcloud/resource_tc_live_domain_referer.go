/*
Provides a resource to create a live domain_referer

Example Usage

```hcl
resource "tencentcloud_live_domain_referer" "domain_referer" {
  domain_name = "5000.liveplay.myqcloud.com"
  enable = 1
  type = 1
  allow_empty = 1
  rules = ""
}
```

Import

live domain_referer can be imported using the id, e.g.

```
terraform import tencentcloud_live_domain_referer.domain_referer domain_referer_id
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

func resourceTencentCloudLiveDomainReferer() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLiveDomainRefererCreate,
		Read:   resourceTencentCloudLiveDomainRefererRead,
		Update: resourceTencentCloudLiveDomainRefererUpdate,
		Delete: resourceTencentCloudLiveDomainRefererDelete,
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
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Whether to enable the referer blacklist authentication of the current domain name.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "List type: 0: blacklist, 1: whitelist.",
			},

			"allow_empty": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Allow blank referers, 0: not allowed, 1: allowed.",
			},

			"rules": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The list of referers to; separate.",
			},
		},
	}
}

func resourceTencentCloudLiveDomainRefererCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_domain_referer.create")()
	defer inconsistentCheck(d, meta)()

	var domainName string
	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
	}

	d.SetId(domainName)

	return resourceTencentCloudLiveDomainRefererUpdate(d, meta)
}

func resourceTencentCloudLiveDomainRefererRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_domain_referer.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	domainRefererId := d.Id()

	domainReferer, err := service.DescribeLiveDomainRefererById(ctx, domainName)
	if err != nil {
		return err
	}

	if domainReferer == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LiveDomainReferer` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if domainReferer.DomainName != nil {
		_ = d.Set("domain_name", domainReferer.DomainName)
	}

	if domainReferer.Enable != nil {
		_ = d.Set("enable", domainReferer.Enable)
	}

	if domainReferer.Type != nil {
		_ = d.Set("type", domainReferer.Type)
	}

	if domainReferer.AllowEmpty != nil {
		_ = d.Set("allow_empty", domainReferer.AllowEmpty)
	}

	if domainReferer.Rules != nil {
		_ = d.Set("rules", domainReferer.Rules)
	}

	return nil
}

func resourceTencentCloudLiveDomainRefererUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_domain_referer.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := live.NewModifyLiveDomainRefererRequest()

	domainRefererId := d.Id()

	request.DomainName = &domainName

	immutableArgs := []string{"domain_name", "enable", "type", "allow_empty", "rules"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().ModifyLiveDomainReferer(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update live domainReferer failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLiveDomainRefererRead(d, meta)
}

func resourceTencentCloudLiveDomainRefererDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_domain_referer.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
