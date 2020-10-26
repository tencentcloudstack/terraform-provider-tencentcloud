/*
Use this data source to query API gateway IP strategy.

Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "ck"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_ip_strategy" "test"{
	service_id      = tencentcloud_api_gateway_service.service.id
	strategy_name	= "tf_test"
	strategy_type	= "BLACK"
	strategy_data	= "9.9.9.9"
}

data "tencentcloud_api_gateway_ip_strategies" "id" {
	service_id = tencentcloud_api_gateway_ip_strategy.test.service_id
}

data "tencentcloud_api_gateway_ip_strategies" "name" {
    service_id = tencentcloud_api_gateway_ip_strategy.test.service_id
	strategy_name = tencentcloud_api_gateway_ip_strategy.test.strategy_name
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

func dataSourceTencentCloudAPIGatewayIpStrategy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayIpStrategyRead,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The service ID to be queried.",
			},
			"strategy_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of IP policy.",
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
				Description: "A list of strategy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"strategy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The strategy ID.",
						},
						"strategy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the strategy.",
						},
						"strategy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the strategy.",
						},
						"ip_list": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The list of IP.",
						},
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The service ID.",
						},
						"bind_api_total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of API bound to the strategy.",
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
						"attach_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of bound API details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service ID.",
									},
									"api_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The API ID.",
									},
									"api_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API interface description.",
									},
									"api_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API name.",
									},
									"vpc_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "VPC ID.",
									},
									"uniq_vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC unique ID.",
									},
									"api_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API type. Valid values: `NORMAL`, `TSF`. `NORMAL` means common API, `TSF` means microservice API.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API protocol.",
									},
									"auth_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API authentication type. Valid values: `SECRET`, `NONE`, `OAUTH`. `SECRET` means key pair authentication, `NONE` means no authentication.",
									},
									"api_business_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of oauth API. This field is valid when the `auth_type` is `OAUTH`, and the values are `NORMAL` (business API) and `OAUTH` (authorization API).",
									},
									"auth_relation_api_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique ID of the associated authorization API, which takes effect when the authType is `OAUTH` and `ApiBusinessType` is normal. Identifies the unique ID of the oauth2.0 authorization API bound to the business API.",
									},
									"tags": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Computed:    true,
										Description: "The label information associated with the API.",
									},
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API path.",
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API request method.",
									},
									"relation_business_api_ids": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Computed:    true,
										Description: "List of business API associated with authorized API.",
									},
									"oauth_config": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: "OAUTH configuration information. It takes effect when authType is `OAUTH`.",
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
				},
			},
		},
	}
}

func dataSourceTencentCloudAPIGatewayIpStrategyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_api_gateway_ip_strategy.read")

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serviceId         = d.Get("service_id").(string)
		infos             []*apigateway.IPStrategy
		list              []map[string]interface{}
		strategyName      string
		err               error
	)
	if v, ok := d.GetOk("strategy_name"); ok {
		strategyName = v.(string)
	}

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		infos, err = apiGatewayService.DescribeIPStrategysStatus(ctx, serviceId, strategyName)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	for _, info := range infos {
		var attachListInfo []map[string]interface{}

		for _, env := range API_GATEWAY_SERVICE_ENVS {
			var strategy *apigateway.IPStrategy
			if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				strategy, err = apiGatewayService.DescribeIPStrategies(ctx, serviceId, *info.StrategyId, env)
				if err != nil {
					return retryError(err, InternalError)
				}
				return nil
			}); err != nil {
				return err
			}

			for _, api := range strategy.BindApis {
				attachListInfo = append(attachListInfo, map[string]interface{}{
					"service_id":                api.ServiceId,
					"api_id":                    api.ApiId,
					"api_desc":                  api.ApiDesc,
					"api_name":                  api.ApiName,
					"vpc_id":                    api.VpcId,
					"uniq_vpc_id":               api.UniqVpcId,
					"api_type":                  api.ApiType,
					"protocol":                  api.Protocol,
					"auth_type":                 api.AuthType,
					"api_business_type":         api.ApiBusinessType,
					"auth_relation_api_id":      api.AuthRelationApiId,
					"tags":                      api.Tags,
					"path":                      api.Path,
					"method":                    api.Method,
					"relation_business_api_ids": api.RelationBuniessApiIds,
					"oauth_config":              flattenOauthConfigMappings(api.OauthConfig),
					"modify_time":               api.ModifiedTime,
					"create_time":               api.CreatedTime,
				})
			}
		}

		infoMap := map[string]interface{}{
			"strategy_id":          info.StrategyId,
			"strategy_name":        info.StrategyName,
			"strategy_type":        info.StrategyType,
			"ip_list":              info.StrategyData,
			"service_id":           info.ServiceId,
			"bind_api_total_count": info.BindApiTotalCount,
			"modify_time":          info.ModifiedTime,
			"create_time":          info.CreatedTime,
			"attach_list":          attachListInfo,
		}

		list = append(list, infoMap)
	}

	if err = d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
	}

	d.SetId(strings.Join([]string{serviceId, strategyName}, FILED_SP))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}
	return nil
}
