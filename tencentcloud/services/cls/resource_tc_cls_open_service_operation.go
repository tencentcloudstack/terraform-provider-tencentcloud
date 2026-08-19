package cls

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clsv20201016 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsOpenServiceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsOpenServiceOperationCreate,
		Read:   resourceTencentCloudClsOpenServiceOperationRead,
		Delete: resourceTencentCloudClsOpenServiceOperationDelete,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Account service status. `0`: service opened, `1`: service not opened.",
			},
		},
	}
}

func resourceTencentCloudClsOpenServiceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_open_service_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = clsv20201016.NewOpenClsServiceRequest()
	)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().OpenClsServiceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Open cls service failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s open cls service failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(helper.BuildToken())
	return resourceTencentCloudClsOpenServiceOperationRead(d, meta)
}

func resourceTencentCloudClsOpenServiceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_open_service_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = clsv20201016.NewGetClsServiceRequest()
		response = clsv20201016.NewGetClsServiceResponse()
	)

	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().GetClsServiceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Get cls service failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read cls service failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Status != nil {
		_ = d.Set("status", response.Response.Status)
	}

	return nil
}

func resourceTencentCloudClsOpenServiceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_open_service_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
