package dnspod

import (
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDnspodDomainLock() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodDomainLockCreate,
		Read:   resourceTencentCloudDnspodDomainLockRead,
		Delete: resourceTencentCloudDnspodDomainLockDelete,
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

			"lock_code": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Domain unlock code, can be obtained through the ModifyDomainLock interface.",
			},
		},
	}
}

func resourceTencentCloudDnspodDomainLockCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_domain_lock.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = dnspod.NewModifyDomainLockRequest()
		response = dnspod.NewModifyDomainLockResponse()
		domain   string
		lockCode string
	)
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("lock_days"); ok {
		request.LockDays = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().ModifyDomainLock(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dnspod domain_lock failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.LockInfo != nil && response.Response.LockInfo.LockCode != nil {
		lockCode = *response.Response.LockInfo.LockCode
	}

	d.SetId(strings.Join([]string{domain, lockCode}, tccommon.FILED_SP))
	_ = d.Set("lock_code", lockCode)

	return resourceTencentCloudDnspodDomainLockRead(d, meta)
}

func resourceTencentCloudDnspodDomainLockRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_domain_lock.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDnspodDomainLockDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_domain_lock.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = dnspod.NewModifyDomainUnlockRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_dnspod_domain_lock id is broken, id is %s", d.Id())
	}
	request.Domain = helper.String(idSplit[0])
	request.LockCode = helper.String(idSplit[1])

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().ModifyDomainUnlock(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dnspod domain_unlock failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
