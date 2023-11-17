/*
Provides a resource to create a dnspod domain_unlock

Example Usage

```hcl
resource "tencentcloud_dnspod_domain_unlock" "domain_unlock" {
  domain = "dnspod.cn"
  lock_code = ""
  domain_id = 123
}
```

*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudDnspodDomainUnlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodDomainUnlockCreate,
		Read:   resourceTencentCloudDnspodDomainUnlockRead,
		Delete: resourceTencentCloudDnspodDomainUnlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Domain name.",
			},

			"lock_code": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Domain unlock code, can be obtained through the ModifyDomainLock interface.",
			},

			"domain_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Domain ID. The parameter DomainId has a higher priority than the parameter Domain. If the parameter DomainId is passed, the parameter Domain will be ignored. You can find all Domains and DomainIds through the DescribeDomainList interface.",
			},
		},
	}
}

func resourceTencentCloudDnspodDomainUnlockCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_domain_unlock.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = dnspod.NewModifyDomainUnlockRequest()
		domain   string
	)
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("lock_code"); ok {
		request.LockCode = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("domain_id"); ok {
		request.DomainId = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDnsPodClient().ModifyDomainUnlock(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dnspod domain_unlock failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(domain)

	return resourceTencentCloudDnspodDomainUnlockRead(d, meta)
}

func resourceTencentCloudDnspodDomainUnlockRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_domain_unlock.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDnspodDomainUnlockDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_domain_unlock.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
