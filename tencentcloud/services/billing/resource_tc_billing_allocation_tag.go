package billing

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	billingv20180709 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBillingAllocationTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBillingAllocationTagCreate,
		Read:   resourceTencentCloudBillingAllocationTagRead,
		Delete: resourceTencentCloudBillingAllocationTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tag_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cost allocation tag key.",
			},

			// computed
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Tag type, 0 normal tag, 1 account tag.",
			},
		},
	}
}

func resourceTencentCloudBillingAllocationTagCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_billing_allocation_tag.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = billingv20180709.NewCreateAllocationTagRequest()
		tagKey  string
	)

	if v, ok := d.GetOk("tag_key"); ok {
		tagKey = v.(string)
		request.TagKey = append(request.TagKey, helper.String(tagKey))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBillingV20180709Client().CreateAllocationTagWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create billing allocation tag failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create billing allocation tag failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(tagKey)
	return resourceTencentCloudBillingAllocationTagRead(d, meta)
}

func resourceTencentCloudBillingAllocationTagRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_billing_allocation_tag.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = BillingService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		tagKey  = d.Id()
	)

	respData, err := service.DescribeBillingAllocationTagById(ctx, tagKey)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_billing_allocation_tag` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.TagKey != nil {
		_ = d.Set("tag_key", respData.TagKey)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	return nil
}

func resourceTencentCloudBillingAllocationTagDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_billing_allocation_tag.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = billingv20180709.NewDeleteAllocationTagRequest()
		tagKey  = d.Id()
	)

	request.TagKey = append(request.TagKey, helper.String(tagKey))
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBillingV20180709Client().DeleteAllocationTagWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete billing allocation tag failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
