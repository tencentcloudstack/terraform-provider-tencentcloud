package dnspod

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

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
			"domain_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Domain ID to bind to the package.",
			},

			"resource_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Package resource ID.",
			},

			// computed
			"domain": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},

			"grade": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Package grade code.",
			},

			"grade_title": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Package grade title.",
			},

			"vip_start_at": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "VIP start time.",
			},

			"vip_end_at": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "VIP end time.",
			},

			"vip_auto_renew": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "VIP auto renew status. YES: enabled, NO: disabled, DEFAULT: default.",
			},

			"remain_times": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Remaining domain bind/change times for the package.",
			},

			"grade_level": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Domain grade level.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Package binding status.",
			},

			"is_grace_period": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Whether the package is in grace period.",
			},

			"downgrade": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether the package is downgraded.",
			},
		},
	}
}

func resourceTencentCloudDnspodPackageDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_package_domain.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		resourceId string
		domainId   uint64
	)

	if v, ok := d.GetOk("resource_id"); ok {
		resourceId = v.(string)
	}

	if v, ok := d.GetOk("domain_id"); ok {
		domainId = uint64(v.(int))
	}

	request := dnspod.NewModifyPackageDomainRequest()
	request.Operation = helper.String("bind")
	request.ResourceId = helper.String(resourceId)
	request.NewDomainId = helper.IntUint64(int(domainId))

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().ModifyPackageDomainWithContext(context.TODO(), request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dnspod package_domain failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{resourceId, strconv.FormatUint(domainId, 10)}, tccommon.FILED_SP))

	return resourceTencentCloudDnspodPackageDomainRead(d, meta)
}

func resourceTencentCloudDnspodPackageDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_package_domain.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	resourceId := idSplit[0]
	domainIdStr := idSplit[1]
	domainId, err := strconv.ParseUint(domainIdStr, 10, 64)
	if err != nil {
		return fmt.Errorf("id is broken, domain_id parse error for id=%s: %v", d.Id(), err)
	}

	var (
		packageItem *dnspod.PackageListItem
	)

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		request := dnspod.NewDescribeDomainVipListRequest()
		request.ResourceIdList = []*string{helper.String(resourceId)}
		limit := uint64(100)
		request.Limit = &limit
		request.Offset = helper.IntUint64(0)

		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().DescribeDomainVipListWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			log.Printf("[CRUD] dnspod package_domain id=%s, response is nil", d.Id())
			return nil
		}

		for _, item := range result.Response.PackageList {
			if item != nil && item.DomainId != nil && *item.DomainId == domainId {
				packageItem = item
				break
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read dnspod package_domain failed, reason:%+v", logId, err)
		return err
	}

	if packageItem == nil {
		log.Printf("[WARN]%s resource `dnspod package_domain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("resource_id", resourceId)
	_ = d.Set("domain_id", int(domainId))

	if packageItem.Domain != nil {
		_ = d.Set("domain", packageItem.Domain)
	}
	if packageItem.Grade != nil {
		_ = d.Set("grade", packageItem.Grade)
	}
	if packageItem.GradeTitle != nil {
		_ = d.Set("grade_title", packageItem.GradeTitle)
	}
	if packageItem.VipStartAt != nil {
		_ = d.Set("vip_start_at", packageItem.VipStartAt)
	}
	if packageItem.VipEndAt != nil {
		_ = d.Set("vip_end_at", packageItem.VipEndAt)
	}
	if packageItem.VipAutoRenew != nil {
		_ = d.Set("vip_auto_renew", packageItem.VipAutoRenew)
	}
	if packageItem.RemainTimes != nil {
		_ = d.Set("remain_times", int(*packageItem.RemainTimes))
	}
	if packageItem.GradeLevel != nil {
		_ = d.Set("grade_level", int(*packageItem.GradeLevel))
	}
	if packageItem.Status != nil {
		_ = d.Set("status", packageItem.Status)
	}
	if packageItem.IsGracePeriod != nil {
		_ = d.Set("is_grace_period", packageItem.IsGracePeriod)
	}
	if packageItem.Downgrade != nil {
		_ = d.Set("downgrade", *packageItem.Downgrade)
	}

	return nil
}

func resourceTencentCloudDnspodPackageDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_package_domain.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	resourceId := idSplit[0]
	oldDomainIdStr := idSplit[1]
	oldDomainId, err := strconv.ParseUint(oldDomainIdStr, 10, 64)
	if err != nil {
		return fmt.Errorf("id is broken, domain_id parse error for id=%s: %v", d.Id(), err)
	}

	if d.HasChange("domain_id") {
		oldId, newId := d.GetChange("domain_id")
		_ = oldId

		newDomainId := uint64(newId.(int))

		request := dnspod.NewModifyPackageDomainRequest()
		request.Operation = helper.String("change")
		request.ResourceId = helper.String(resourceId)
		request.DomainId = helper.IntUint64(int(oldDomainId))
		request.NewDomainId = helper.IntUint64(int(newDomainId))

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().ModifyPackageDomainWithContext(context.TODO(), request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update dnspod package_domain failed, reason:%+v", logId, err)
			return err
		}

		d.SetId(strings.Join([]string{resourceId, strconv.FormatUint(newDomainId, 10)}, tccommon.FILED_SP))
	}

	return resourceTencentCloudDnspodPackageDomainRead(d, meta)
}

func resourceTencentCloudDnspodPackageDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_package_domain.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	resourceId := idSplit[0]
	domainIdStr := idSplit[1]
	domainId, err := strconv.ParseUint(domainIdStr, 10, 64)
	if err != nil {
		return fmt.Errorf("id is broken, domain_id parse error for id=%s: %v", d.Id(), err)
	}

	request := dnspod.NewModifyPackageDomainRequest()
	request.Operation = helper.String("unbind")
	request.ResourceId = helper.String(resourceId)
	request.DomainId = helper.IntUint64(int(domainId))

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().ModifyPackageDomainWithContext(context.TODO(), request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete dnspod package_domain failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
