/*
Use this data source to query API gateway services.

Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "niceservice"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

data "tencentcloud_api_gateway_services" "name" {
    service_name = tencentcloud_api_gateway_service.service.service_name
}

data "tencentcloud_api_gateway_services" "id" {
    service_id = tencentcloud_api_gateway_service.service.id
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func dataSourceTencentCloudAPIGatewayServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayServicesRead,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Service name for query.",
			},
			"service_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Service ID for query.",
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
				Description: "A list of services.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom service ID.",
						},
						"service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom service name.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service frontend request type. Valid values: `http`, `https`, `http&https`.",
						},
						"service_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom service description.",
						},
						"exclusive_set_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Self-deployed cluster name, which is used to specify the self-deployed cluster where the service is to be created.",
						},
						"net_type": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Description: "Network type list, which is used to specify the supported network types. " +
								"Valid values: `INNER`, `OUTER`. " +
								"`INNER` indicates access over private network, and `OUTER` indicates access over public network.",
						},
						"ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP version number.",
						},
						"internal_sub_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private network access sub-domain name.",
						},
						"outer_sub_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public network access subdomain name.",
						},
						"inner_http_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port number for http access over private network.",
						},
						"inner_https_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port number for https access over private network.",
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
						"usage_plan_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of attach usage plans. Each element contains the following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"usage_plan_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the usage plan.",
									},
									"usage_plan_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the usage plan.",
									},
									"bind_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Binding type.",
									},
									"api_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the API.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAPIGatewayServicesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_api_gateway_services.read")()

	var (
		logId                  = getLogId(contextNil)
		ctx                    = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService      = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		services               []*apigateway.Service
		serviceName, serviceId string
		has                    bool
		err                    error
	)

	if v, ok := d.GetOk("service_name"); ok {
		serviceName = v.(string)
	}
	if v, ok := d.GetOk("service_id"); ok {
		serviceId = v.(string)
	}

	if outErr := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		services, err = apiGatewayService.DescribeServicesStatus(ctx, serviceId, serviceName)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); outErr != nil {
		return outErr
	}

	list := make([]map[string]interface{}, 0, len(services))

	for _, service := range services {
		var info apigateway.DescribeServiceResponse
		if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			info, has, err = apiGatewayService.DescribeService(ctx, *service.ServiceId)
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
		if !has {
			continue
		}

		var plans []*apigateway.ApiUsagePlan

		var planList = make([]map[string]interface{}, 0, len(info.Response.ApiIdStatusSet))
		var hasContains = make(map[string]bool, len(info.Response.ApiIdStatusSet))

		//from service
		if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			plans, err = apiGatewayService.DescribeServiceUsagePlan(ctx, *service.ServiceId)
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		}); err != nil {
			return err
		}

		for _, item := range plans {
			if hasContains[*item.UsagePlanId] {
				continue
			}
			hasContains[*item.UsagePlanId] = true
			planList = append(
				planList, map[string]interface{}{
					"usage_plan_id":   item.UsagePlanId,
					"usage_plan_name": item.UsagePlanName,
					"bind_type":       API_GATEWAY_TYPE_SERVICE,
					"api_id":          "",
				})
		}

		//from api
		if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			plans, err = apiGatewayService.DescribeApiUsagePlan(ctx, *service.ServiceId)
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
		for _, item := range plans {
			planList = append(
				planList, map[string]interface{}{
					"usage_plan_id":   item.UsagePlanId,
					"usage_plan_name": item.UsagePlanName,
					"bind_type":       API_GATEWAY_TYPE_API,
					"api_id":          item.ApiId,
				})
		}

		list = append(list, map[string]interface{}{
			"service_id":          info.Response.ServiceId,
			"service_name":        info.Response.ServiceName,
			"protocol":            info.Response.Protocol,
			"service_desc":        info.Response.ServiceDesc,
			"exclusive_set_name":  info.Response.ExclusiveSetName,
			"ip_version":          info.Response.IpVersion,
			"net_type":            info.Response.NetTypes,
			"internal_sub_domain": info.Response.InternalSubDomain,
			"outer_sub_domain":    info.Response.OuterSubDomain,
			"inner_http_port":     info.Response.InnerHttpPort,
			"inner_https_port":    info.Response.InnerHttpsPort,
			"modify_time":         info.Response.ModifiedTime,
			"create_time":         info.Response.CreatedTime,
			"usage_plan_list":     planList,
		})
	}

	if err = d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
	}

	d.SetId(strings.Join([]string{serviceName, serviceId}, FILED_SP))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}
	return nil
}
