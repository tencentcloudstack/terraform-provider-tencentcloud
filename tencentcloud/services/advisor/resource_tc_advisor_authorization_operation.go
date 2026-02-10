package advisor

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	advisorv20200721 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/advisor/v20200721"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAdvisorAuthorizationOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAdvisorAuthorizationOperationCreate,
		Read:   resourceTencentCloudAdvisorAuthorizationOperationRead,
		Delete: resourceTencentCloudAdvisorAuthorizationOperationDelete,
		Schema: map[string]*schema.Schema{},
	}
}

func resourceTencentCloudAdvisorAuthorizationOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_advisor_authorization_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = advisorv20200721.NewCreateAdvisorAuthorizationRequest()
	)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAdvisorV20200721Client().CreateAdvisorAuthorizationWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create advisor authorization operation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(helper.BuildToken())
	return resourceTencentCloudAdvisorAuthorizationOperationRead(d, meta)
}

func resourceTencentCloudAdvisorAuthorizationOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_advisor_authorization_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudAdvisorAuthorizationOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_advisor_authorization_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
