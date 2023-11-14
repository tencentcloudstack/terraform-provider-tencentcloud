/*
Provides a resource to create a apiGateway plugin_attachment

Example Usage

```hcl
resource "tencentcloud_api_gateway_plugin_attachment" "plugin_attachment" {
  plugin_id = ""
  service_id = ""
  environment_name = ""
  api_id = ""
}
```

Import

apiGateway plugin_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_api_gateway_plugin_attachment.plugin_attachment plugin_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiGateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudApiGatewayPluginAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApiGatewayPluginAttachmentCreate,
		Read:   resourceTencentCloudApiGatewayPluginAttachmentRead,
		Delete: resourceTencentCloudApiGatewayPluginAttachmentDelete,
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

func resourceTencentCloudApiGatewayPluginAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_plugin_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = apiGateway.NewAttachPluginRequest()
		response        = apiGateway.NewAttachPluginResponse()
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
		request.ApiId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseApiGatewayClient().AttachPlugin(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create apiGateway pluginAttachment failed, reason:%+v", logId, err)
		return err
	}

	pluginId = *response.Response.PluginId
	d.SetId(strings.Join([]string{pluginId, serviceId, environmentName, apiId}, FILED_SP))

	return resourceTencentCloudApiGatewayPluginAttachmentRead(d, meta)
}

func resourceTencentCloudApiGatewayPluginAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_plugin_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ApiGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

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
		log.Printf("[WARN]%s resource `ApiGatewayPluginAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if pluginAttachment.PluginId != nil {
		_ = d.Set("plugin_id", pluginAttachment.PluginId)
	}

	if pluginAttachment.ServiceId != nil {
		_ = d.Set("service_id", pluginAttachment.ServiceId)
	}

	if pluginAttachment.EnvironmentName != nil {
		_ = d.Set("environment_name", pluginAttachment.EnvironmentName)
	}

	if pluginAttachment.ApiId != nil {
		_ = d.Set("api_id", pluginAttachment.ApiId)
	}

	return nil
}

func resourceTencentCloudApiGatewayPluginAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_plugin_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ApiGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
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
