package apigateway

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAPIGatewayUpdateService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayUpdateServiceCreate,
		Read:   resourceTencentCloudAPIGatewayUpdateServiceRead,
		Delete: resourceTencentCloudAPIGatewayUpdateServiceDelete,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Service ID.",
			},
			"environment_name": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"test", "prepub", "release"}),
				Description:  "The name of the environment to be switched, currently supporting three environments: test (test environment), prepub (pre release environment), and release (release environment).",
			},
			"version_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The version number of the switch.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayUpdateServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_update_service.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		request   = apigateway.NewUpdateServiceRequest()
		serviceId string
	)

	if v, ok := d.GetOk("service_id"); ok {
		request.ServiceId = helper.String(v.(string))
		serviceId = v.(string)
	}

	if v, ok := d.GetOk("environment_name"); ok {
		request.EnvironmentName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("version_name"); ok {
		request.VersionName = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAPIGatewayClient().UpdateService(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate apigateway updateService failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(serviceId)
	return resourceTencentCloudAPIGatewayUpdateServiceRead(d, meta)
}

func resourceTencentCloudAPIGatewayUpdateServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_update_service.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudAPIGatewayUpdateServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_update_service.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
