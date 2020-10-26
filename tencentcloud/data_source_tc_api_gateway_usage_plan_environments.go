/*
Used to query the environment list bound by the plan.

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

resource "tencentcloud_api_gateway_usage_plan_attachment" "attach_service" {
	usage_plan_id  = tencentcloud_api_gateway_usage_plan.plan.id
	service_id     = tencentcloud_api_gateway_service.service.id
	environment    = "test"
	bind_type      = "SERVICE"
}

data "tencentcloud_api_gateway_usage_plan_environments" "environment_test" {
	usage_plan_id = tencentcloud_api_gateway_usage_plan_attachment.attach_service.usage_plan_id
	bind_type     = "SERVICE"
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func dataSourceTencentCloudAPIGatewayUsagePlanEnvironments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudUsagePlanEnvironmentRead,

		Schema: map[string]*schema.Schema{
			"usage_plan_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the usage plan to be queried.",
			},
			"bind_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      API_GATEWAY_TYPE_SERVICE,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_TYPES),
				Description:  "Binding type. Valid values: `API`, `SERVICE`. Default value: `SERVICE`.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			// Computed values.
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of usage plan binding details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The service ID.",
						},
						"service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The service name.",
						},
						"api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The API ID, this value is empty if attach service.",
						},
						"api_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The API name, this value is empty if attach service.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The API path, this value is empty if attach service.",
						},
						"method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The API method, this value is empty if attach service.",
						},
						"environment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The environment name.",
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
				},
			},
		},
	}
}

func dataSourceTencentCloudUsagePlanEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_api_gateway_usage_plans.read")

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		usagePlanId       = d.Get("usage_plan_id").(string)
		bindType          = d.Get("bind_type").(string)
		infos             []*apigateway.UsagePlanEnvironment
		list              []map[string]interface{}
		err               error
	)

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		infos, err = apiGatewayService.DescribeUsagePlanEnvironments(ctx, usagePlanId, bindType)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	for _, info := range infos {
		list = append(list, map[string]interface{}{
			"service_id":   info.ServiceId,
			"service_name": info.ServiceName,
			"api_id":       info.ApiId,
			"api_name":     info.ApiName,
			"path":         info.Path,
			"method":       info.Method,
			"environment":  info.Environment,
			"modify_time":  info.ModifiedTime,
			"create_time":  info.CreatedTime,
		})
	}

	if err = d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
	}

	d.SetId(strings.Join([]string{usagePlanId, bindType}, FILED_SP))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}
	return nil
}
