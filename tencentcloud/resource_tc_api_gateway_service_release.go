/*
Use this resource to create API gateway service release.

Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "myservice"
  protocol     = "http"
  service_desc = "my nice service"
  net_type     = ["INNER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "api" {
    service_id            = tencentcloud_api_gateway_service.service.id
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

resource "tencentcloud_api_gateway_service_release" "service" {
  service_id       = tencentcloud_api_gateway_api.api.service.id
  environment_name = "release"
  release_desc     = "test service release"
}
```

Import

API gateway service release can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_service_release.service service-jjt3fs3s#release#20201015121916d85fb161-eaec-4dda-a7e0-659aa5f401be
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

func resourceTencentCloudAPIGatewayServiceRelease() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayServiceReleaseCreate,
		Read:   resourceTencentCloudAPIGatewayServiceReleaseRead,
		Delete: resourceTencentCloudAPIGatewayServiceReleaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of API gateway service.",
			},
			"environment_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_SERVICE_ENVS),
				Description:  "API gateway service environment name to be released. Valid values: `test`, `prepub`, `release`.",
			},
			"release_desc": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This release description of the API gateway service.",
			},
			"release_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The release version.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayServiceReleaseCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_service_release.create")()

	var (
		logId                = getLogId(contextNil)
		ctx                  = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService    = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serviceId            = d.Get("service_id").(string)
		environmentName      = d.Get("environment_name").(string)
		releaseDesc          = d.Get("release_desc").(string)
		err                  error
		has, existEnv        bool
		releaseResponse      *apigateway.ReleaseServiceResponse
		serviceResponse      apigateway.DescribeServiceResponse
		checkServiceResponse apigateway.DescribeServiceResponse
	)

	//check API gateway serviceid and service contains api
	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		checkServiceResponse, has, err = apiGatewayService.DescribeService(ctx, serviceId)
		if err != nil {
			return retryError(err, InternalError)
		}
		if has {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("API gateway service %s not found on server", serviceId))
	}); err != nil {
		return err
	}
	if *checkServiceResponse.Response.ApiTotalCount < 1 {
		return fmt.Errorf("there is no API under the current service, please create an API before publishing")
	}

	// release API gateway service
	if releaseResponse, err = apiGatewayService.ReleaseService(ctx, serviceId, environmentName, releaseDesc); err != nil {
		return err
	}

	//wait service release ok
	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		serviceResponse, has, err = apiGatewayService.DescribeService(ctx, serviceId)
		if err != nil {
			return retryError(err, InternalError)
		}
		if has {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("API gateway service %s not found on server", serviceId))
	}); err != nil {
		return err
	}

	for _, v := range serviceResponse.Response.AvailableEnvironments {
		if *v == environmentName {
			existEnv = true
			break
		}
	}

	if !existEnv {
		return fmt.Errorf("API gateway service not release success")
	}

	d.SetId(strings.Join([]string{serviceId, environmentName, *releaseResponse.Response.Result.ReleaseVersion}, FILED_SP))

	return resourceTencentCloudAPIGatewayServiceReleaseRead(d, meta)
}

func resourceTencentCloudAPIGatewayServiceReleaseRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_service_release.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                = d.Id()
		info              []*apigateway.ServiceReleaseHistoryInfo
		err               error
	)

	ids := strings.Split(id, FILED_SP)
	if len(ids) != 3 {
		return fmt.Errorf("id is broken, id is %s", id)
	}
	var (
		serviceId  = ids[0]
		envName    = ids[1]
		envVersion = ids[2]
	)

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		info, _, err = apiGatewayService.DescribeServiceEnvironmentReleaseHistory(ctx, serviceId, envName)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	if len(info) == 0 {
		d.SetId("")
		return nil
	}

	for _, v := range info {
		if *v.VersionName == envVersion {
			_ = d.Set("service_id", serviceId)
			_ = d.Set("environment_name", envName)
			_ = d.Set("release_desc", v.VersionDesc)
			_ = d.Set("release_version", v.VersionName)
		}
	}

	return nil
}

func resourceTencentCloudAPIGatewayServiceReleaseDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_service_release.delete")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                = d.Id()
		err               error
	)

	ids := strings.Split(id, FILED_SP)
	if len(ids) != 3 {
		return fmt.Errorf("id is broken, id is %s", id)
	}
	var (
		serviceId = ids[0]
		envName   = ids[1]
	)

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if err = apiGatewayService.UnReleaseService(ctx, serviceId, envName); err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
