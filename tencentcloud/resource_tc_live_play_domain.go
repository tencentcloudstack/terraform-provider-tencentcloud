/*
Provides a resource to create a live play_domain

Example Usage

```hcl
resource "tencentcloud_live_play_domain" "play_domain" {
  domain_name = "5000.livepush.myqcloud.com"
  play_type = 1
}
```

Import

live play_domain can be imported using the id, e.g.

```
terraform import tencentcloud_live_play_domain.play_domain play_domain_id
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

func resourceTencentCloudLivePlayDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLivePlayDomainCreate,
		Read:   resourceTencentCloudLivePlayDomainRead,
		Update: resourceTencentCloudLivePlayDomainUpdate,
		Delete: resourceTencentCloudLivePlayDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain Name.",
			},

			"play_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Type of streaming domain name. 1 - Domestic; 2 - Global; 3 - Overseas.",
			},
		},
	}
}

func resourceTencentCloudLivePlayDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_play_domain.create")()
	defer inconsistentCheck(d, meta)()

	var domainName string
	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
	}

	d.SetId(domainName)

	return resourceTencentCloudLivePlayDomainUpdate(d, meta)
}

func resourceTencentCloudLivePlayDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_play_domain.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	playDomainId := d.Id()

	playDomain, err := service.DescribeLivePlayDomainById(ctx, domainName)
	if err != nil {
		return err
	}

	if playDomain == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LivePlayDomain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if playDomain.DomainName != nil {
		_ = d.Set("domain_name", playDomain.DomainName)
	}

	if playDomain.PlayType != nil {
		_ = d.Set("play_type", playDomain.PlayType)
	}

	return nil
}

func resourceTencentCloudLivePlayDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_play_domain.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := live.NewModifyLivePlayDomainRequest()

	playDomainId := d.Id()

	request.DomainName = &domainName

	immutableArgs := []string{"domain_name", "play_type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
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
		log.Printf("[CRITAL]%s update live playDomain failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLivePlayDomainRead(d, meta)
}

func resourceTencentCloudLivePlayDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_play_domain.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
