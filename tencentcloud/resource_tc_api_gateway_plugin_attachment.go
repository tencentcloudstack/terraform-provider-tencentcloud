/*
Provides a resource to create a apiGateway plugin_attachment

Example Usage

```hcl
resource "tencentcloud_api_gateway_plugin" "example" {
  plugin_name = "tf-example"
  plugin_type = "IPControl"
  plugin_data = jsonencode({
    "type" : "white_list",
    "blocks" : "1.1.1.1",
  })
  description = "desc."
}

resource "tencentcloud_api_gateway_service" "example" {
  service_name = "tf_example_service"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "example" {
  service_id            = tencentcloud_api_gateway_service.example.id
  api_name              = "tf_example_api"
  api_desc              = "desc."
  auth_type             = "APP"
  protocol              = "HTTP"
  enable_cors           = true
  request_config_path   = "/user/info"
  request_config_method = "GET"

  request_parameters {
    name          = "name"
    position      = "QUERY"
    type          = "string"
    desc          = "desc."
    default_value = "terraform"
    required      = true
  }

  service_config_type      = "HTTP"
  service_config_timeout   = 15
  service_config_url       = "https://www.qq.com"
  service_config_path      = "/user"
  service_config_method    = "GET"
  response_type            = "HTML"
  response_success_example = "success"
  response_fail_example    = "fail"

  response_error_codes {
    code           = 400
    msg            = "system error msg."
    desc           = "system error desc."
    converted_code = 407
    need_convert   = true
  }
}

resource "tencentcloud_api_gateway_service_release" "example" {
  service_id       = tencentcloud_api_gateway_api.example.service_id
  environment_name = "release"
  release_desc     = "desc."
}

resource "tencentcloud_api_gateway_plugin_attachment" "example" {
  plugin_id        = tencentcloud_api_gateway_plugin.example.id
  service_id       = tencentcloud_api_gateway_service_release.example.service_id
  api_id           = tencentcloud_api_gateway_api.example.id
  environment_name = "release"
}
```

Import

apiGateway plugin_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_api_gateway_plugin_attachment.example plugin-hnqntalp#service-q3f533ja#release#api-62ud9woa
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiGateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAPIGatewayPluginAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayPluginAttachmentCreate,
		Read:   resourceTencentCloudAPIGatewayPluginAttachmentRead,
		Delete: resourceTencentCloudAPIGatewayPluginAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"plugin_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Id of Plugin.",
			},
			"service_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Id of Service.",
			},
			"environment_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Name of Environment.",
			},
			"api_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Id of API.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayPluginAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_plugin_attachment.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId           = getLogId(contextNil)
		request         = apiGateway.NewAttachPluginRequest()
		pluginId        string
		serviceId       string
		environmentName string
		apiId           string
	)

	if v, ok := d.GetOk("plugin_id"); ok {
		pluginId = v.(string)
		request.PluginId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_id"); ok {
		serviceId = v.(string)
		request.ServiceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("environment_name"); ok {
		environmentName = v.(string)
		request.EnvironmentName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("api_id"); ok {
		apiId = v.(string)
		request.ApiIds = []*string{helper.String(v.(string))}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAPIGatewayClient().AttachPlugin(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if !*result.Response.Result {
			e = fmt.Errorf(" create apigateway pluginAttachment result: false.")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create apiGateway pluginAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{pluginId, serviceId, environmentName, apiId}, FILED_SP))
	return resourceTencentCloudAPIGatewayPluginAttachmentRead(d, meta)
}

func resourceTencentCloudAPIGatewayPluginAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_plugin_attachment.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	pluginId := idSplit[0]
	serviceId := idSplit[1]
	environmentName := idSplit[2]
	apiId := idSplit[3]

	pluginAttachment, err := service.DescribeApiGatewayPluginAttachmentById(ctx, pluginId, serviceId, environmentName, apiId)
	if err != nil {
		return err
	}

	if pluginAttachment == nil {
		d.SetId("")
		return fmt.Errorf("resource `APIGatewayPluginAttachment` %s does not exist", d.Id())
	}

	_ = d.Set("plugin_id", pluginId)

	if pluginAttachment.ServiceId != nil {
		_ = d.Set("service_id", pluginAttachment.ServiceId)
	}

	if pluginAttachment.Environment != nil {
		_ = d.Set("environment_name", pluginAttachment.Environment)
	}

	if pluginAttachment.ApiId != nil {
		_ = d.Set("api_id", pluginAttachment.ApiId)
	}

	return nil
}

func resourceTencentCloudAPIGatewayPluginAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_plugin_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	pluginId := idSplit[0]
	serviceId := idSplit[1]
	environmentName := idSplit[2]
	apiId := idSplit[3]

	if err := service.DeleteApiGatewayPluginAttachmentById(ctx, pluginId, serviceId, environmentName, apiId); err != nil {
		return err
	}

	return nil
}
