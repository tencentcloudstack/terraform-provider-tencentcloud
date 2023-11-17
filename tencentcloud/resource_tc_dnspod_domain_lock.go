/*
Provides a resource to create a dnspod domain_lock

Example Usage

```hcl
resource "tencentcloud_dnspod_domain_lock" "domain_lock" {
  domain = "dnspod.cn"
  lock_days = 30
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

func resourceTencentCloudDnspodDomainLock() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodDomainLockCreate,
		Read:   resourceTencentCloudDnspodDomainLockRead,
		Delete: resourceTencentCloudDnspodDomainLockDelete,
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

			"lock_days": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The number of max days to lock the domain+ Old packages: D_FREE 30 days, D_PLUS 90 days, D_EXTRA 30 days, D_EXPERT 60 days, D_ULTRA 365 days+ New packages: DP_FREE 365 days, DP_PLUS 365 days, DP_EXTRA 365 days, DP_EXPERT 365 days, DP_ULTRA 365 days.",
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

func resourceTencentCloudDnspodDomainLockCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_domain_lock.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = dnspod.NewModifyDomainLockRequest()
		domain   string
	)
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("lock_days"); ok {
		request.LockDays = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("domain_id"); ok {
		request.DomainId = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDnsPodClient().ModifyDomainLock(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dnspod domain_lock failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(domain)

	return resourceTencentCloudDnspodDomainLockRead(d, meta)
}

func resourceTencentCloudDnspodDomainLockRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_domain_lock.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDnspodDomainLockDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_domain_lock.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
