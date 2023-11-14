/*
Provides a resource to create a apigateway api_app

Example Usage

```hcl
resource "tencentcloud_apigateway_api_app" "api_app" {
  api_app_name = ""
  api_app_desc = ""
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

apigateway api_app can be imported using the id, e.g.

```
terraform import tencentcloud_apigateway_api_app.api_app api_app_id
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

func resourceTencentCloudApigatewayApiApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApigatewayApiAppCreate,
		Read:   resourceTencentCloudApigatewayApiAppRead,
		Update: resourceTencentCloudApigatewayApiAppUpdate,
		Delete: resourceTencentCloudApigatewayApiAppDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"api_app_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "User defined application name.",
			},

			"api_app_desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Application description.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudApigatewayApiAppCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_api_app.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = apigateway.NewCreateApiAppRequest()
		response = apigateway.NewCreateApiAppResponse()
		apiAppId string
	)
	if v, ok := d.GetOk("api_app_name"); ok {
		request.ApiAppName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("api_app_desc"); ok {
		request.ApiAppDesc = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseApigatewayClient().CreateApiApp(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create apigateway apiApp failed, reason:%+v", logId, err)
		return err
	}

	apiAppId = *response.Response.ApiAppId
	d.SetId(apiAppId)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::apigw:%s:uin/:apiAppId/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudApigatewayApiAppRead(d, meta)
}

func resourceTencentCloudApigatewayApiAppRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_api_app.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

	apiAppId := d.Id()

	apiApp, err := service.DescribeApigatewayApiAppById(ctx, apiAppId)
	if err != nil {
		return err
	}

	if apiApp == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ApigatewayApiApp` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if apiApp.ApiAppName != nil {
		_ = d.Set("api_app_name", apiApp.ApiAppName)
	}

	if apiApp.ApiAppDesc != nil {
		_ = d.Set("api_app_desc", apiApp.ApiAppDesc)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "apigw", "apiAppId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudApigatewayApiAppUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_api_app.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := apigateway.NewModifyApiAppRequest()

	apiAppId := d.Id()

	request.ApiAppId = &apiAppId

	immutableArgs := []string{"api_app_name", "api_app_desc"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("api_app_name") {
		if v, ok := d.GetOk("api_app_name"); ok {
			request.ApiAppName = helper.String(v.(string))
		}
	}

	if d.HasChange("api_app_desc") {
		if v, ok := d.GetOk("api_app_desc"); ok {
			request.ApiAppDesc = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseApigatewayClient().ModifyApiApp(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update apigateway apiApp failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("apigw", "apiAppId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudApigatewayApiAppRead(d, meta)
}

func resourceTencentCloudApigatewayApiAppDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_api_app.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
	apiAppId := d.Id()

	if err := service.DeleteApigatewayApiAppById(ctx, apiAppId); err != nil {
		return err
	}

	return nil
}
