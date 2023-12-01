/*
Use this resource to create API gateway access key.

Example Usage

Automatically generate key for API gateway access key.

```hcl
resource "tencentcloud_api_gateway_api_key" "example_auto" {
  secret_name = "tf_example_auto"
  status      = "on"
}
```

Manually generate a secret key for API gateway access key.

```hcl
resource "tencentcloud_api_gateway_api_key" "example_manual" {
  secret_name       = "tf_example_manual"
  status            = "on"
  access_key_type   = "manual"
  access_key_id     = "28e287e340507fa147b2c8284dab542f"
  access_key_secret = "0198a4b8c3105080f4acd9e507599eff"
}
```
Import

API gateway access key can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_api_key.test AKIDMZwceezso9ps5p8jkro8a9fwe1e7nzF2k50B
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func resourceTencentCloudAPIGatewayAPIKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayAPIKeyCreate,
		Read:   resourceTencentCloudAPIGatewayAPIKeyRead,
		Update: resourceTencentCloudAPIGatewayAPIKeyUpdate,
		Delete: resourceTencentCloudAPIGatewayAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"secret_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Custom key name.",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      API_GATEWAY_KEY_ENABLED,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_KEYS),
				Description:  "Key status. Valid values: `on`, `off`.",
			},
			"access_key_type": {
				Optional:     true,
				Default:      API_GATEWAY_KEY_TYPE_AUTO,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_KEYS_TYPE),
				Type:         schema.TypeString,
				Description:  "Key type, supports both auto and manual (custom keys), defaults to auto.",
			},
			"access_key_id": {
				Optional:     true,
				Computed:     true,
				Type:         schema.TypeString,
				ValidateFunc: validateStringLengthInRange(5, 50),
				Description:  "User defined key ID, required when access_key_type is manual. The length is 5-50 characters, consisting of letters, numbers, and English underscores.",
			},
			"access_key_secret": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateStringLengthInRange(10, 50),
				Description:  "The user-defined key must be passed when the access_key_type is manual. The length is 10-50 characters, consisting of letters, numbers, and English underscores.",
			},
			// Computed values.
			"modify_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayAPIKeyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api_key.create")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		request           = apigateway.NewCreateApiKeyRequest()
		response          = apigateway.NewCreateApiKeyResponse()
		statusStr         string
		accessKeyType     string
		accessKeyId       string
		accessKeySecret   string
	)

	if v, ok := d.GetOk("secret_name"); ok {
		request.SecretName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		statusStr = v.(string)
	}

	if v, ok := d.GetOk("access_key_type"); ok {
		request.AccessKeyType = helper.String(v.(string))
		accessKeyType = v.(string)
	}

	if accessKeyType == API_GATEWAY_KEY_TYPE_MANUAL {
		if v, ok := d.GetOk("access_key_id"); ok {
			accessKeyId = v.(string)
		}

		if v, ok := d.GetOk("access_key_secret"); ok {
			accessKeySecret = v.(string)
		}

		if accessKeyId == "" || accessKeySecret == "" {
			errRet := fmt.Errorf("`access_key_id`, `access_key_secret` required when access_key_type is `manual`")
			return errRet
		}

		request.AccessKeyId = &accessKeyId
		request.AccessKeySecret = &accessKeySecret
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAPIGatewayClient().CreateApiKey(request)
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

	//set status to disable
	if statusStr == API_GATEWAY_KEY_DISABLED {
		if err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			if err = apiGatewayService.DisableApiKey(ctx, accessKeyId); err != nil {
				return retryError(err)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	d.SetId(*response.Response.Result.AccessKeyId)

	return resourceTencentCloudAPIGatewayAPIKeyRead(d, meta)
}

func resourceTencentCloudAPIGatewayAPIKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api_key.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		accessKeyId       = d.Id()
		apiKey            *apigateway.ApiKey
		err               error
		has               bool
	)

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		apiKey, has, err = apiGatewayService.DescribeApiKey(ctx, accessKeyId)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("secret_name", apiKey.SecretName)
	_ = d.Set("status", API_GATEWAY_KEY_INT2STRS[*apiKey.Status])
	_ = d.Set("access_key_type", apiKey.AccessKeyType)
	_ = d.Set("access_key_id", apiKey.AccessKeyId)
	_ = d.Set("access_key_secret", apiKey.AccessKeySecret)
	_ = d.Set("modify_time", apiKey.ModifiedTime)
	_ = d.Set("create_time", apiKey.CreatedTime)

	return nil
}

func resourceTencentCloudAPIGatewayAPIKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api_key.update")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		request           = apigateway.NewUpdateApiKeyRequest()
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		accessKeyId       = d.Id()
	)

	immutableFields := []string{"access_key_id", "access_key_type"}
	for _, f := range immutableFields {
		if d.HasChange(f) {
			return fmt.Errorf("cannot update argument `%s`", f)
		}
	}

	if d.HasChange("access_key_secret") {
		if d.Get("access_key_type") == API_GATEWAY_KEY_TYPE_AUTO {
			errRet := fmt.Errorf("`access_key_id`, `access_key_secret` updated when access_key_type is `auto`")
			return errRet
		}

		request.AccessKeyId = &accessKeyId
		if v, ok := d.GetOk("access_key_secret"); ok {
			request.AccessKeySecret = helper.String(v.(string))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseAPIGatewayClient().UpdateApiKey(request)
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
	}

	if d.HasChange("status") {
		var (
			statusStr = d.Get("status").(string)
			err       error
		)

		if err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			if statusStr == API_GATEWAY_KEY_DISABLED {
				err = apiGatewayService.DisableApiKey(ctx, accessKeyId)
			} else {
				err = apiGatewayService.EnableApiKey(ctx, accessKeyId)
			}
			if err != nil {
				return retryError(err)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return resourceTencentCloudAPIGatewayAPIKeyRead(d, meta)
}

func resourceTencentCloudAPIGatewayAPIKeyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api_key.delete")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		accessKeyId       = d.Id()
	)

	//set status to disable before delete
	if d.Get("status") != API_GATEWAY_KEY_DISABLED {
		if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			if err := apiGatewayService.DisableApiKey(ctx, accessKeyId); err != nil {
				return retryError(err)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		inErr := apiGatewayService.DeleteApiKey(ctx, accessKeyId)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
}
