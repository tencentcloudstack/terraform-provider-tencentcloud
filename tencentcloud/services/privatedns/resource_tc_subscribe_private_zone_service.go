package privatedns

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func ResourceTencentCloudSubscribePrivateZoneService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSubscribePrivateZoneServiceCreate,
		Read:   resourceTencentCloudSubscribePrivateZoneServiceRead,
		Delete: resourceTencentCloudSubscribePrivateZoneServiceDelete,
		Schema: map[string]*schema.Schema{
			"service_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private domain resolution service activation status.",
			},
		},
	}
}

func resourceTencentCloudSubscribePrivateZoneServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_subscribe_private_zone_service.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = privatedns.NewSubscribePrivateZoneServiceRequest()
		response = privatedns.NewSubscribePrivateZoneServiceResponse()
	)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivateDnsClient().SubscribePrivateZoneServiceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create subscribe private zone service failed. Response is nil"))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create subscribe private zone service failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.ServiceStatus != nil {
		_ = d.Set("service_status", response.Response.ServiceStatus)
	}

	d.SetId(*response.Response.RequestId)

	return resourceTencentCloudSubscribePrivateZoneServiceRead(d, meta)
}

func resourceTencentCloudSubscribePrivateZoneServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_subscribe_private_zone_service.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSubscribePrivateZoneServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_subscribe_private_zone_service.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
