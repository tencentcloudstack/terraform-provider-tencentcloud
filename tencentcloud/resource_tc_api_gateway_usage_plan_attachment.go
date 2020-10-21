/*
Use this resource to attach API gateway usage plan to service.

Example Usage

```hcl
resource "tencentcloud_api_gateway_usage_plan" "plan" {
	usage_plan_name         = "my_plan"
	usage_plan_desc         = "nice plan"
	max_request_num         = 100
	max_request_num_pre_sec = 10
}

resource "tencentcloud_api_gateway_service" "service" {
  	service_name = "niceservice"
  	protocol     = "http&https"
  	service_desc = "your nice service"
  	net_type     = ["INNER", "OUTER"]
  	ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "api" {
    service_id            = tencentcloud_api_gateway_service.service.id
    api_name              = "hello_update"
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
        desc          = "your email please?"
        default_value = "tom@qq.com"
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
    	code           = 10
        msg            = "system error"
       	desc           = "system error code"
       	converted_code = -10
        need_convert   = true
    }
}

resource "tencentcloud_api_gateway_usage_plan_attachment" "attach_service" {
	usage_plan_id  = tencentcloud_api_gateway_usage_plan.plan.id
  	service_id     = tencentcloud_api_gateway_service.service.id
	environment    = "release"
	bind_type      = "API"
   	api_id         = tencentcloud_api_gateway_api.api.id
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

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Description:  "Environment to be bound `test`,`prepub` or `release`.",
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
	)

	if v, ok := d.GetOk("api_id"); ok {
		apiId = v.(string)
	}

	if bindType == API_GATEWAY_TYPE_API && apiId == "" {
		return fmt.Errorf("parameter `api_ids` is required when `bind_type` is `API`")
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

	d.SetId(strings.Join([]string{usagePlanId, serviceId, environment, bindType, apiId}, FILED_SP))

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
	if len(ids) != 5 {
		return fmt.Errorf("id is broken, id is %s", id)
	}

	var (
		usagePlanId = ids[0]
		serviceId   = ids[1]
		environment = ids[2]
		bindType    = ids[3]
		apiId       = ids[4]
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
	if len(ids) != 5 {
		return fmt.Errorf("id is broken, id is %s", id)
	}

	var (
		usagePlanId = ids[0]
		serviceId   = ids[1]
		environment = ids[2]
		bindType    = ids[3]
		apiId       = ids[4]
	)

	if usagePlanId == "" || serviceId == "" || environment == "" || bindType == "" {
		return fmt.Errorf("id is broken, id is %s", id)
	}
	if bindType == API_GATEWAY_TYPE_API && apiId == "" {
		return fmt.Errorf("id is broken, id is %s", id)
	}

	// BindEnvironment
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
