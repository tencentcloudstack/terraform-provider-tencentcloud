/*
Provides a resource to create a apigateway plugin

Example Usage

```hcl
resource "tencentcloud_apigateway_plugin" "plugin" {
  plugin_name = ""
  plugin_type = ""
  plugin_data = ""
  description = ""
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

apigateway plugin can be imported using the id, e.g.

```
terraform import tencentcloud_apigateway_plugin.plugin plugin_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudApigatewayPlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApigatewayPluginCreate,
		Read:   resourceTencentCloudApigatewayPluginRead,
		Update: resourceTencentCloudApigatewayPluginUpdate,
		Delete: resourceTencentCloudApigatewayPluginDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"plugin_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "User defined plugin name. Up to 50 characters, up to 2 characters, supports a-z, A-Z, 0-9, _, Must start with a letter and end with a letter or number.",
			},

			"plugin_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Plugin type. Currently supports IPControl, TrafficControl, Cors, CustomiReq, CustomiAuth, Routing, TrafficControlByParameter, CircuitBreaker, and ProxyCache.",
			},

			"plugin_data": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Plugin definition statements, supporting JSON.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Plugin description, limited to 200 words.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudApigatewayPluginCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_plugin.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = apigateway.NewCreatePluginRequest()
		response = apigateway.NewCreatePluginResponse()
		pluginId string
	)
	if v, ok := d.GetOk("plugin_name"); ok {
		request.PluginName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("plugin_type"); ok {
		request.PluginType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("plugin_data"); ok {
		request.PluginData = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseApigatewayClient().CreatePlugin(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create apigateway plugin failed, reason:%+v", logId, err)
		return err
	}

	pluginId = *response.Response.PluginId
	d.SetId(pluginId)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::apigw:%s:uin/:pluginId/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudApigatewayPluginRead(d, meta)
}

func resourceTencentCloudApigatewayPluginRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_plugin.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

	pluginId := d.Id()

	plugin, err := service.DescribeApigatewayPluginById(ctx, pluginId)
	if err != nil {
		return err
	}

	if plugin == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ApigatewayPlugin` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if plugin.PluginName != nil {
		_ = d.Set("plugin_name", plugin.PluginName)
	}

	if plugin.PluginType != nil {
		_ = d.Set("plugin_type", plugin.PluginType)
	}

	if plugin.PluginData != nil {
		_ = d.Set("plugin_data", plugin.PluginData)
	}

	if plugin.Description != nil {
		_ = d.Set("description", plugin.Description)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "apigw", "pluginId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudApigatewayPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_plugin.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := apigateway.NewModifyPluginRequest()

	pluginId := d.Id()

	request.PluginId = &pluginId

	immutableArgs := []string{"plugin_name", "plugin_type", "plugin_data", "description"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("plugin_name") {
		if v, ok := d.GetOk("plugin_name"); ok {
			request.PluginName = helper.String(v.(string))
		}
	}

	if d.HasChange("plugin_data") {
		if v, ok := d.GetOk("plugin_data"); ok {
			request.PluginData = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseApigatewayClient().ModifyPlugin(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update apigateway plugin failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("apigw", "pluginId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudApigatewayPluginRead(d, meta)
}

func resourceTencentCloudApigatewayPluginDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_plugin.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
	pluginId := d.Id()

	if err := service.DeleteApigatewayPluginById(ctx, pluginId); err != nil {
		return err
	}

	return nil
}
