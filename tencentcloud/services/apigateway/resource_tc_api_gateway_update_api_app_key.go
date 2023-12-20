package apigateway

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiGateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudApiGatewayUpdateApiAppKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApiGatewayUpdateApiAppKeyCreate,
		Read:   resourceTencentCloudApiGatewayUpdateApiAppKeyRead,
		Delete: resourceTencentCloudApiGatewayUpdateApiAppKeyDelete,

		Schema: map[string]*schema.Schema{
			"api_app_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Application unique ID.",
			},
			"api_app_key": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Key of the application.",
			},
			//"api_app_secret": {
			//	Optional:    true,
			//	ForceNew:    true,
			//	Type:        schema.TypeString,
			//	Description: "Application Secret.",
			//},
		},
	}
}

func resourceTencentCloudApiGatewayUpdateApiAppKeyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_update_api_app_key.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = apiGateway.NewUpdateApiAppKeyRequest()
		apiAppId string
	)

	if v, ok := d.GetOk("api_app_id"); ok {
		request.ApiAppId = helper.String(v.(string))
		apiAppId = v.(string)
	}

	if v, ok := d.GetOk("api_app_key"); ok {
		request.ApiAppKey = helper.String(v.(string))
	}

	//if v, ok := d.GetOk("api_app_secret"); ok {
	//	request.ApiAppSecret = helper.String(v.(string))
	//}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAPIGatewayClient().UpdateApiAppKey(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate apiGateway updateApiAppKey failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(apiAppId)
	return resourceTencentCloudApiGatewayUpdateApiAppKeyRead(d, meta)
}

func resourceTencentCloudApiGatewayUpdateApiAppKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_update_api_app_key.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudApiGatewayUpdateApiAppKeyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_update_api_app_key.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
