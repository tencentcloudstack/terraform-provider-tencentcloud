/*
Provides a resource to create a apiGateway update_api_app_key

Example Usage

```hcl
resource "tencentcloud_api_gateway_update_api_app_key" "example" {
  api_app_id  = "app-krljp4wn"
  api_app_key = "APID6JmG21yRCc03h4z16hlsTqj1wpO3dB3ZQcUP"
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiGateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudApiGatewayUpdateApiAppKey() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_api_gateway_update_api_app_key.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAPIGatewayClient().UpdateApiAppKey(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_api_gateway_update_api_app_key.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudApiGatewayUpdateApiAppKeyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_update_api_app_key.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
