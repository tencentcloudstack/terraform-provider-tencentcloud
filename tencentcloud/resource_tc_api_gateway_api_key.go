/*
Use this resource to create API gateway access key.

Example Usage

```hcl
resource "tencentcloud_api_gateway_api_key" "test" {
  secret_name = "my_api_key"
  status      = "on"
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
			// Computed values.
			"access_key_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created API key.",
			},
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
		secretName        = d.Get("secret_name").(string)
		statusStr         = d.Get("status").(string)
		accessKeyId       string
		err               error
	)

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		accessKeyId, err = apiGatewayService.CreateApiKey(ctx, secretName)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(accessKeyId)

	//wait API key create ok
	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, inErr := apiGatewayService.DescribeApiKey(ctx, accessKeyId)
		if inErr != nil {
			return retryError(inErr, InternalError)
		}
		if !has {
			return resource.RetryableError(fmt.Errorf("accessKeyId %s not found on server", accessKeyId))
		}
		return nil

	}); err != nil {
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
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		accessKeyId       = d.Id()
	)

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
