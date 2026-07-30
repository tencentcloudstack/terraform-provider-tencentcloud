package dnspod

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDnspodPackageDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodPackageDomainCreate,
		Read:   resourceTencentCloudDnspodPackageDomainRead,
		Update: resourceTencentCloudDnspodPackageDomainUpdate,
		Delete: resourceTencentCloudDnspodPackageDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Package resource ID.",
			},

			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Domain ID to bind to the package.",
			},

			// computed
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain.",
			},
		},
	}
}

func resourceTencentCloudDnspodPackageDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_package_domain.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = dnspod.NewModifyPackageDomainRequest()
		resourceId string
		domainId   uint64
	)

	if v, ok := d.GetOk("resource_id"); ok {
		resourceId = v.(string)
	}

	if v, ok := d.GetOkExists("domain_id"); ok {
		domainId = uint64(v.(int))
	}

	request.Operation = helper.String("bind")
	request.ResourceId = helper.String(resourceId)
	request.NewDomainId = helper.IntUint64(int(domainId))

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().ModifyPackageDomainWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Bind dnspod package domain failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dnspod package_domain failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(resourceId)

	// ModifyPackageDomain is async, poll DescribeDomainVipList until the bound
	// domain reports Status `enable`.
	if err := waitDnspodPackageDomainStatus(ctx, meta, resourceId, domainId, true); err != nil {
		return err
	}

	return resourceTencentCloudDnspodPackageDomainRead(d, meta)
}

func resourceTencentCloudDnspodPackageDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_package_domain.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = DnspodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		resourceId = d.Id()
	)

	respData, err := service.DescribeDnspodPackageDomainById(ctx, resourceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_dnspod_package_domain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("resource_id", resourceId)

	if respData.DomainId != nil {
		_ = d.Set("domain_id", int(*respData.DomainId))
	}

	if respData.Domain != nil {
		_ = d.Set("domain", respData.Domain)
	}

	return nil
}

func resourceTencentCloudDnspodPackageDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_package_domain.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		resourceId = d.Id()
	)

	needChange := false
	mutableArgs := []string{"domain_id"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		oldId, newId := d.GetChange("domain_id")
		oldDomainId := uint64(oldId.(int))
		newDomainId := uint64(newId.(int))

		request := dnspod.NewModifyPackageDomainRequest()
		request.Operation = helper.String("change")
		request.ResourceId = helper.String(resourceId)
		request.DomainId = helper.IntUint64(int(oldDomainId))
		request.NewDomainId = helper.IntUint64(int(newDomainId))

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().ModifyPackageDomainWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Change dnspod package domain failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update dnspod package_domain failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// ModifyPackageDomain is async, poll DescribeDomainVipList until the
		// newly bound domain reports Status `enable`.
		if err := waitDnspodPackageDomainStatus(ctx, meta, resourceId, newDomainId, true); err != nil {
			return err
		}
	}

	return resourceTencentCloudDnspodPackageDomainRead(d, meta)
}

func resourceTencentCloudDnspodPackageDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_package_domain.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = dnspod.NewModifyPackageDomainRequest()
		resourceId = d.Id()
		domainId   uint64
	)

	if v, ok := d.GetOkExists("domain_id"); ok {
		domainId = uint64(v.(int))
	}

	request.Operation = helper.String("unbind")
	request.ResourceId = helper.String(resourceId)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().ModifyPackageDomainWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Unbind dnspod package domain failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete dnspod package_domain failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// ModifyPackageDomain is async, poll DescribeDomainVipList until the
	// package reports Status `enable`.
	if err := waitDnspodPackageDomainStatus(ctx, meta, resourceId, domainId, false); err != nil {
		return err
	}

	return nil
}

// waitDnspodPackageDomainStatus polls DescribeDomainVipList (queried by the
// package resource id) until the package task finishes. ModifyPackageDomain is
// asynchronous: the domain stays in a transient status until the task
// completes, at which point its Status turns to `enable` (compared
// case-insensitively to tolerate `ENABLE`).
//
// requireDomainMatch controls whether the bound domain must equal domainId:
//   - true  (bind / change): waits until the package item whose DomainId equals
//     domainId reports Status `enable`.
//   - false (unbind): waits until the package item reports Status `enable`,
//     regardless of which domain it currently carries.
func waitDnspodPackageDomainStatus(ctx context.Context, meta interface{}, resourceId string, domainId uint64, requireDomainMatch bool) error {
	service := DnspodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	return resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		item, e := service.DescribeDnspodPackageDomainById(ctx, resourceId)
		if e != nil {
			return resource.NonRetryableError(e)
		}

		if item == nil || item.Status == nil {
			return resource.RetryableError(fmt.Errorf("dnspod package_domain async task not finished, resource_id=%s domain_id=%d: package not ready", resourceId, domainId))
		}

		if requireDomainMatch && (item.DomainId == nil || *item.DomainId != domainId) {
			return resource.RetryableError(fmt.Errorf("dnspod package_domain async task not finished, resource_id=%s domain_id=%d not bound yet", resourceId, domainId))
		}

		if strings.EqualFold(*item.Status, "enable") {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("dnspod package_domain async task not finished, resource_id=%s domain_id=%d current status=%q", resourceId, domainId, *item.Status))
	})
}
