/*
Provides a resource to create a apigateway api_key

Example Usage

```hcl
resource "tencentcloud_apigateway_api_key" "api_key" {
  secret_name = ""
  access_key_type = ""
  access_key_id = ""
  access_key_secret = ""
}
```

Import

apigateway api_key can be imported using the id, e.g.

```
terraform import tencentcloud_apigateway_api_key.api_key api_key_id
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

func resourceTencentCloudApigatewayApiKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApigatewayApiKeyCreate,
		Read:   resourceTencentCloudApigatewayApiKeyRead,
		Update: resourceTencentCloudApigatewayApiKeyUpdate,
		Delete: resourceTencentCloudApigatewayApiKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"secret_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "User defined key name.",
			},

			"access_key_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Key type, supports both auto and manual (custom keys), defaults to auto.",
			},

			"access_key_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User defined key ID, required when AccessKeyType is manual. The length is 5-50 characters, consisting of letters, numbers, and English underscores.",
			},

			"access_key_secret": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The user-defined key must be passed when the AccessKeyType is manual. The length is 10-50 characters, consisting of letters, numbers, and English underscores.",
			},
		},
	}
}

func resourceTencentCloudApigatewayApiKeyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_api_key.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = apigateway.NewCreateApiKeyRequest()
		response    = apigateway.NewCreateApiKeyResponse()
		accessKeyId string
	)
	if v, ok := d.GetOk("secret_name"); ok {
		request.SecretName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("access_key_type"); ok {
		request.AccessKeyType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("access_key_id"); ok {
		accessKeyId = v.(string)
		request.AccessKeyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("access_key_secret"); ok {
		request.AccessKeySecret = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseApigatewayClient().CreateApiKey(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create apigateway apiKey failed, reason:%+v", logId, err)
		return err
	}

	accessKeyId = *response.Response.AccessKeyId
	d.SetId(accessKeyId)

	return resourceTencentCloudApigatewayApiKeyRead(d, meta)
}

func resourceTencentCloudApigatewayApiKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_api_key.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

	apiKeyId := d.Id()

	apiKey, err := service.DescribeApigatewayApiKeyById(ctx, accessKeyId)
	if err != nil {
		return err
	}

	if apiKey == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ApigatewayApiKey` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if apiKey.SecretName != nil {
		_ = d.Set("secret_name", apiKey.SecretName)
	}

	if apiKey.AccessKeyType != nil {
		_ = d.Set("access_key_type", apiKey.AccessKeyType)
	}

	if apiKey.AccessKeyId != nil {
		_ = d.Set("access_key_id", apiKey.AccessKeyId)
	}

	if apiKey.AccessKeySecret != nil {
		_ = d.Set("access_key_secret", apiKey.AccessKeySecret)
	}

	return nil
}

func resourceTencentCloudApigatewayApiKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_api_key.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := apigateway.NewUpdateApiKeyRequest()

	apiKeyId := d.Id()

	request.AccessKeyId = &accessKeyId

	immutableArgs := []string{"secret_name", "access_key_type", "access_key_id", "access_key_secret"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("access_key_id") {
		if v, ok := d.GetOk("access_key_id"); ok {
			request.AccessKeyId = helper.String(v.(string))
		}
	}

	if d.HasChange("access_key_secret") {
		if v, ok := d.GetOk("access_key_secret"); ok {
			request.AccessKeySecret = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseApigatewayClient().UpdateApiKey(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update apigateway apiKey failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudApigatewayApiKeyRead(d, meta)
}

func resourceTencentCloudApigatewayApiKeyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_api_key.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
	apiKeyId := d.Id()

	if err := service.DeleteApigatewayApiKeyById(ctx, accessKeyId); err != nil {
		return err
	}

	return nil
}
