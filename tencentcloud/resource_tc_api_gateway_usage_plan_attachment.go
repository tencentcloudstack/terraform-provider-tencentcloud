/*
Use this resource to attach API gateway usage plan to service.

~> **NOTE:** If the `auth_type` parameter of API is not `SECRET`, it cannot be bound access key.

Example Usage

Normal creation

```hcl
resource "tencentcloud_api_gateway_usage_plan" "example" {
  usage_plan_name         = "tf_example"
  usage_plan_desc         = "desc."
  max_request_num         = 100
  max_request_num_pre_sec = 10
}

resource "tencentcloud_api_gateway_service" "example" {
  service_name = "tf_example"
  protocol     = "http&https"
  service_desc = "desc."
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "example" {
  service_id            = tencentcloud_api_gateway_service.example.id
  api_name              = "tf_example"
  api_desc              = "my hello api update"
  auth_type             = "SECRET"
  protocol              = "HTTP"
  enable_cors           = true
  request_config_path   = "/user/info"
  request_config_method = "POST"
  request_parameters {
    name          = "email"
    position      = "QUERY"
    type          = "string"
    desc          = "desc."
    default_value = "test@qq.com"
    required      = true
  }
  service_config_type      = "HTTP"
  service_config_timeout   = 10
  service_config_url       = "http://www.tencent.com"
  service_config_path      = "/user"
  service_config_method    = "POST"
  response_type            = "XML"
  response_success_example = "<note>success</note>"
  response_fail_example    = "<note>fail</note>"
  response_error_codes {
    code           = 500
    msg            = "system error"
    desc           = "system error code"
    converted_code = 5000
    need_convert   = true
  }
}

resource "tencentcloud_api_gateway_usage_plan_attachment" "example" {
  usage_plan_id = tencentcloud_api_gateway_usage_plan.example.id
  service_id    = tencentcloud_api_gateway_service.example.id
  environment   = "release"
  bind_type     = "API"
  api_id        = tencentcloud_api_gateway_api.example.id
}
```

Bind the key to a usage plan

```hcl
resource "tencentcloud_api_gateway_api_key" "example" {
  secret_name = "tf_example"
  status      = "on"
}

resource "tencentcloud_api_gateway_usage_plan_attachment" "example" {
  usage_plan_id = tencentcloud_api_gateway_usage_plan.example.id
  service_id    = tencentcloud_api_gateway_service.example.id
  environment   = "release"
  bind_type     = "API"
  api_id        = tencentcloud_api_gateway_api.example.id

  access_key_ids = [tencentcloud_api_gateway_api_key.example.id]
}
```

Import

API gateway usage plan attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_usage_plan_attachment.attach_service usagePlan-pe7fbdgn#service-kuqd6xqk#release#API#api-p8gtanvy
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func resourceTencentCloudAPIGatewayUsagePlanAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayUsagePlanAttachmentCreate,
		Read:   resourceTencentCloudAPIGatewayUsagePlanAttachmentRead,
		Delete: resourceTencentCloudAPIGatewayUsagePlanAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"usage_plan_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the usage plan.",
			},
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the service.",
			},
			"environment": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_SERVICE_ENVS),
				Description:  "The environment to be bound. Valid values: `test`, `prepub`, `release`.",
			},
			"bind_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      API_GATEWAY_TYPE_SERVICE,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_TYPES),
				Description:  "Binding type. Valid values: `API`, `SERVICE`. Default value is `SERVICE`.",
			},
			"api_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the API. This parameter will be required when `bind_type` is `API`.",
			},
			"access_key_ids": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Array of key IDs to be bound.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayUsagePlanAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_usage_plan_attachment.create")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		usagePlanId       = d.Get("usage_plan_id").(string)
		serviceId         = d.Get("service_id").(string)
		environment       = d.Get("environment").(string)
		bindType          = d.Get("bind_type").(string)
		apiId             string
		err               error
		accessKeys        []string
	)

	if v, ok := d.GetOk("api_id"); ok {
		apiId = v.(string)
	}

	if v, ok := d.GetOk("access_key_ids"); ok {
		accessKeySet := v.(*schema.Set).List()
		for i := range accessKeySet {
			accessKeyId := accessKeySet[i].(string)
			accessKeys = append(accessKeys, accessKeyId)
		}
	}

	if bindType == API_GATEWAY_TYPE_API && apiId == "" {
		return fmt.Errorf("parameter `api_ids` is required when `bind_type` is `API`")
	}

	if bindType == API_GATEWAY_TYPE_API && apiId != "" && len(accessKeys) != 0 {
		if err = checkApiAuthType(ctx, apiGatewayService, serviceId, apiId); err != nil {
			return err
		}
	}

	//check usage plan
	if err = checkUsagePlan(ctx, apiGatewayService, usagePlanId); err != nil {
		return err
	}

	//check service
	if err = checkService(ctx, apiGatewayService, serviceId); err != nil {
		return err
	}

	// BindEnvironment
	if err = apiGatewayService.BindEnvironment(ctx, serviceId, environment, bindType, usagePlanId, apiId); err != nil {
		return err
	}

	// BindSecretIds
	if bindType == API_GATEWAY_TYPE_API && apiId != "" && len(accessKeys) != 0 {
		if err = apiGatewayService.BindSecretIds(ctx, usagePlanId, accessKeys); err != nil {
			return err
		}
	}

	accessKeysStr := strings.Join(accessKeys, COMMA_SP)
	d.SetId(strings.Join([]string{usagePlanId, serviceId, environment, bindType, apiId, accessKeysStr}, FILED_SP))

	return resourceTencentCloudAPIGatewayUsagePlanAttachmentRead(d, meta)
}

func resourceTencentCloudAPIGatewayUsagePlanAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_usage_plan_attachment.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                = d.Id()
		err               error
	)
	ids := strings.Split(id, FILED_SP)
	if len(ids) != 6 {
		return fmt.Errorf("id is broken, id is %s", id)
	}

	var (
		usagePlanId   = ids[0]
		serviceId     = ids[1]
		environment   = ids[2]
		bindType      = ids[3]
		apiId         = ids[4]
		accessKeysStr = ids[5]
	)

	if usagePlanId == "" || serviceId == "" || environment == "" || bindType == "" {
		return fmt.Errorf("id is broken, id is %s", id)
	}
	if bindType == API_GATEWAY_TYPE_API && apiId == "" {
		return fmt.Errorf("id is broken, id is %s", id)
	}

	// check usage plan
	if err = checkUsagePlan(ctx, apiGatewayService, usagePlanId); err != nil {
		return err
	}

	//check service
	if err = checkService(ctx, apiGatewayService, serviceId); err != nil {
		return err
	}

	plans := make([]*apigateway.ApiUsagePlan, 0)
	if bindType == API_GATEWAY_TYPE_API {
		if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			plans, err = apiGatewayService.DescribeApiUsagePlan(ctx, serviceId)
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
	} else {
		if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			plans, err = apiGatewayService.DescribeServiceUsagePlan(ctx, serviceId)
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	var setData = func() error {
		_ = d.Set("usage_plan_id", usagePlanId)
		_ = d.Set("service_id", serviceId)
		_ = d.Set("environment", environment)
		_ = d.Set("bind_type", bindType)
		_ = d.Set("api_id", apiId)
		return nil
	}

	for _, plan := range plans {
		if *plan.UsagePlanId == usagePlanId && *plan.Environment == environment {
			if bindType == API_GATEWAY_TYPE_API {
				if plan.ApiId != nil && *plan.ApiId == apiId {
					return setData()
				}
			} else {
				return setData()
			}
		}
	}

	if bindType == API_GATEWAY_TYPE_API && apiId != "" && accessKeysStr != "" {
		var accessKeyList []*apigateway.UsagePlanBindSecret
		if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			accessKeyList, err = apiGatewayService.DescribeApiUsagePlanSecretIds(ctx, usagePlanId)
			if err != nil {
				return retryError(err, InternalError)
			}

			return nil
		}); err != nil {
			return err
		}

		if len(accessKeyList) != 0 {
			tmpList := make([]string, 0)
			for _, v := range accessKeyList {
				tmpList = append(tmpList, *v.AccessKeyId)
			}
			_ = d.Set("access_key_ids", tmpList)
		}
	}

	return nil
}

func resourceTencentCloudAPIGatewayUsagePlanAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_usage_plan_attachment.delete")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                = d.Id()
		err               error
	)

	ids := strings.Split(id, FILED_SP)
	if len(ids) != 6 {
		return fmt.Errorf("id is broken, id is %s", id)
	}

	var (
		usagePlanId   = ids[0]
		serviceId     = ids[1]
		environment   = ids[2]
		bindType      = ids[3]
		apiId         = ids[4]
		accessKeysStr = ids[5]
	)

	if usagePlanId == "" || serviceId == "" || environment == "" || bindType == "" {
		return fmt.Errorf("id is broken, id is %s", id)
	}
	if bindType == API_GATEWAY_TYPE_API && apiId == "" {
		return fmt.Errorf("id is broken, id is %s", id)
	}

	if bindType == API_GATEWAY_TYPE_API && apiId != "" && accessKeysStr != "" {
		// UnBindSecretIds
		accessKeyList := strings.Split(accessKeysStr, COMMA_SP)
		if err = apiGatewayService.UnBindSecretIds(ctx, usagePlanId, accessKeyList); err != nil {
			return err
		}
	}

	// UnBindEnvironment
	if err = apiGatewayService.UnBindEnvironment(ctx, serviceId, environment, bindType, usagePlanId, apiId); err != nil {
		return err
	}

	return nil
}

func checkUsagePlan(ctx context.Context, api APIGatewayService, usagePlanId string) error {
	var (
		err error
		has bool
	)
	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, err = api.DescribeUsagePlan(ctx, usagePlanId)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	if !has {
		return fmt.Errorf("usage plan %s not exist", usagePlanId)
	}

	return nil
}

func checkService(ctx context.Context, api APIGatewayService, serviceId string) error {
	var (
		err error
		has bool
	)
	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, err = api.DescribeService(ctx, serviceId)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	if !has {
		return fmt.Errorf("service %s not exist", serviceId)
	}

	return nil
}

func checkApiAuthType(ctx context.Context, api APIGatewayService, serviceId, apiId string) error {
	var (
		res apigateway.ApiInfo
		err error
		has bool
	)
	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		res, has, err = api.DescribeApi(ctx, serviceId, apiId)
		if err != nil {
			return retryError(err, InternalError)
		}

		return nil
	}); err != nil {
		return err
	}

	if !has {
		return fmt.Errorf("service %s not exist", serviceId)
	}

	if *res.AuthType != "SECRET" {
		return fmt.Errorf("the auth_type value of the current apiId %s is not `SECRET`", apiId)
	}

	return nil
}
