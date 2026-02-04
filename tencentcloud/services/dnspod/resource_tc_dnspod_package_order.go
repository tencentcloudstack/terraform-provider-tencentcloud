package dnspod

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspodintl "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/dnspod/v20210323"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDnspodPackageOrder() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodPackageOrderCreate,
		Read:   resourceTencentCloudDnspodPackageOrderRead,
		Delete: resourceTencentCloudDnspodPackageOrderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain.",
			},

			"grade": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Valid options for the package version are as follows: `DPG_PROFESSIONAL`; `DPG_ENTERPRISE`; `DPG_ULTIMATE`.",
			},

			// computed
			"domain_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Domain ID.",
			},
		},
	}
}

func resourceTencentCloudDnspodPackageOrderCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_package_order.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = dnspodintl.NewCreatePackageOrderRequest()
		domain  string
	)

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("grade"); ok {
		request.Grade = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodIntlClient().CreatePackageOrderWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dnspod package order failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(domain)
	return resourceTencentCloudDnspodPackageOrderRead(d, meta)
}

func resourceTencentCloudDnspodPackageOrderRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_package_order.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DnspodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		domain  = d.Id()
	)

	respData, err := service.DescribeDnspodPackageOrderById(ctx, domain)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_dnspod_package_order` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Domain != nil {
		_ = d.Set("domain", respData.Domain)
	}

	if respData.Grade != nil {
		_ = d.Set("grade", respData.Grade)
	}

	if respData.DomainId != nil {
		_ = d.Set("domain_id", respData.DomainId)
	}

	return nil
}

func resourceTencentCloudDnspodPackageOrderDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_package_order.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = dnspodintl.NewDeletePackageOrderRequest()
		domain  = d.Id()
	)

	request.Domain = &domain
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodIntlClient().DeletePackageOrderWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete dnspod package order failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
