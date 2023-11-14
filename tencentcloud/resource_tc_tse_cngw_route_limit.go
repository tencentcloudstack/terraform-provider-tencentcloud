/*
Provides a resource to create a tse cngw_route_limit

Example Usage

```hcl
resource "tencentcloud_tse_cngw_route_limit" "cngw_route_limit" {
  gateway_id = "gateway-xxxxxx"
  limit_detail {
		enabled = true
		qps_thresholds {
			unit = "second"
			max = 50
		}
		limit_by = "ip"
		response_type = "default"
		hide_client_headers = false
		is_delay = false
		path = "/test"
		header = "auth"
		external_redis {
			redis_host = ""
			redis_password = ""
			redis_port =
			redis_timeout =
		}
		policy = "redis"
		rate_limit_response {
			body = ""
			headers {
				key = ""
				value = ""
			}
			http_status =
		}
		rate_limit_response_url = ""
		line_up_time =

  }
}
```

Import

tse cngw_route_limit can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_route_limit.cngw_route_limit cngw_route_limit_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTseCngwRouteLimit() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseCngwRouteLimitCreate,
		Read:   resourceTencentCloudTseCngwRouteLimitRead,
		Update: resourceTencentCloudTseCngwRouteLimitUpdate,
		Delete: resourceTencentCloudTseCngwRouteLimitDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Gateway ID.",
			},

			"limit_detail": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Rate limit configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Status of service rate limit.",
						},
						"qps_thresholds": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Qps threshold.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"unit": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Qps threshold unit.Reference value:- second- minute- hour- day- month- year.",
									},
									"max": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The max threshold.",
									},
								},
							},
						},
						"limit_by": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Basis for service rate limit.Reference value:- ip- service- consumer- credential- path- header.",
						},
						"response_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Response strategy.Reference value:- url, forward request according to url- text, response configuration- default, return directly.",
						},
						"hide_client_headers": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether to hide the headers of client.",
						},
						"is_delay": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether to enable request queuing.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Request paths that require rate limit.",
						},
						"header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Request headers that require rate limit.",
						},
						"external_redis": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "External redis information, maybe null.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"redis_host": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Redis ip, maybe null.",
									},
									"redis_password": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Redis password, maybe null.",
									},
									"redis_port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Redis port, maybe null.",
									},
									"redis_timeout": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Redis timeout, unit:ms, maybe null.",
									},
								},
							},
						},
						"policy": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Counter policy.Reference value:- local- redis- external_redis.",
						},
						"rate_limit_response": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Response configuration, the response strategy is text, maybe null.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"body": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Custom response body, maybe bull.",
									},
									"headers": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Headrs.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Key of header.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Value of header.",
												},
											},
										},
									},
									"http_status": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Http status code.",
									},
								},
							},
						},
						"rate_limit_response_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Request forwarding address, maybe null.",
						},
						"line_up_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Queue time.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTseCngwRouteLimitCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_route_limit.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tse.NewCreateCloudNativeAPIGatewayRouteRateLimitRequest()
		response  = tse.NewCreateCloudNativeAPIGatewayRouteRateLimitResponse()
		gatewayId string
	)
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
		request.GatewayId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "limit_detail"); ok {
		cloudNativeAPIGatewayRateLimitDetail := tse.CloudNativeAPIGatewayRateLimitDetail{}
		if v, ok := dMap["enabled"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.Enabled = helper.Bool(v.(bool))
		}
		if v, ok := dMap["qps_thresholds"]; ok {
			for _, item := range v.([]interface{}) {
				qpsThresholdsMap := item.(map[string]interface{})
				qpsThreshold := tse.QpsThreshold{}
				if v, ok := qpsThresholdsMap["unit"]; ok {
					qpsThreshold.Unit = helper.String(v.(string))
				}
				if v, ok := qpsThresholdsMap["max"]; ok {
					qpsThreshold.Max = helper.IntInt64(v.(int))
				}
				cloudNativeAPIGatewayRateLimitDetail.QpsThresholds = append(cloudNativeAPIGatewayRateLimitDetail.QpsThresholds, &qpsThreshold)
			}
		}
		if v, ok := dMap["limit_by"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.LimitBy = helper.String(v.(string))
		}
		if v, ok := dMap["response_type"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.ResponseType = helper.String(v.(string))
		}
		if v, ok := dMap["hide_client_headers"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.HideClientHeaders = helper.Bool(v.(bool))
		}
		if v, ok := dMap["is_delay"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.IsDelay = helper.Bool(v.(bool))
		}
		if v, ok := dMap["path"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.Path = helper.String(v.(string))
		}
		if v, ok := dMap["header"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.Header = helper.String(v.(string))
		}
		if externalRedisMap, ok := helper.InterfaceToMap(dMap, "external_redis"); ok {
			externalRedis := tse.ExternalRedis{}
			if v, ok := externalRedisMap["redis_host"]; ok {
				externalRedis.RedisHost = helper.String(v.(string))
			}
			if v, ok := externalRedisMap["redis_password"]; ok {
				externalRedis.RedisPassword = helper.String(v.(string))
			}
			if v, ok := externalRedisMap["redis_port"]; ok {
				externalRedis.RedisPort = helper.IntInt64(v.(int))
			}
			if v, ok := externalRedisMap["redis_timeout"]; ok {
				externalRedis.RedisTimeout = helper.IntInt64(v.(int))
			}
			cloudNativeAPIGatewayRateLimitDetail.ExternalRedis = &externalRedis
		}
		if v, ok := dMap["policy"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.Policy = helper.String(v.(string))
		}
		if rateLimitResponseMap, ok := helper.InterfaceToMap(dMap, "rate_limit_response"); ok {
			rateLimitResponse := tse.RateLimitResponse{}
			if v, ok := rateLimitResponseMap["body"]; ok {
				rateLimitResponse.Body = helper.String(v.(string))
			}
			if v, ok := rateLimitResponseMap["headers"]; ok {
				for _, item := range v.([]interface{}) {
					headersMap := item.(map[string]interface{})
					kVMapping := tse.KVMapping{}
					if v, ok := headersMap["key"]; ok {
						kVMapping.Key = helper.String(v.(string))
					}
					if v, ok := headersMap["value"]; ok {
						kVMapping.Value = helper.String(v.(string))
					}
					rateLimitResponse.Headers = append(rateLimitResponse.Headers, &kVMapping)
				}
			}
			if v, ok := rateLimitResponseMap["http_status"]; ok {
				rateLimitResponse.HttpStatus = helper.IntInt64(v.(int))
			}
			cloudNativeAPIGatewayRateLimitDetail.RateLimitResponse = &rateLimitResponse
		}
		if v, ok := dMap["rate_limit_response_url"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.RateLimitResponseUrl = helper.String(v.(string))
		}
		if v, ok := dMap["line_up_time"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.LineUpTime = helper.IntInt64(v.(int))
		}
		request.LimitDetail = &cloudNativeAPIGatewayRateLimitDetail
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().CreateCloudNativeAPIGatewayRouteRateLimit(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tse cngwRouteLimit failed, reason:%+v", logId, err)
		return err
	}

	gatewayId = *response.Response.GatewayId
	d.SetId(gatewayId)

	return resourceTencentCloudTseCngwRouteLimitRead(d, meta)
}

func resourceTencentCloudTseCngwRouteLimitRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_route_limit.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	cngwRouteLimitId := d.Id()

	cngwRouteLimit, err := service.DescribeTseCngwRouteLimitById(ctx, gatewayId)
	if err != nil {
		return err
	}

	if cngwRouteLimit == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseCngwRouteLimit` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if cngwRouteLimit.GatewayId != nil {
		_ = d.Set("gateway_id", cngwRouteLimit.GatewayId)
	}

	if cngwRouteLimit.LimitDetail != nil {
		limitDetailMap := map[string]interface{}{}

		if cngwRouteLimit.LimitDetail.Enabled != nil {
			limitDetailMap["enabled"] = cngwRouteLimit.LimitDetail.Enabled
		}

		if cngwRouteLimit.LimitDetail.QpsThresholds != nil {
			qpsThresholdsList := []interface{}{}
			for _, qpsThresholds := range cngwRouteLimit.LimitDetail.QpsThresholds {
				qpsThresholdsMap := map[string]interface{}{}

				if qpsThresholds.Unit != nil {
					qpsThresholdsMap["unit"] = qpsThresholds.Unit
				}

				if qpsThresholds.Max != nil {
					qpsThresholdsMap["max"] = qpsThresholds.Max
				}

				qpsThresholdsList = append(qpsThresholdsList, qpsThresholdsMap)
			}

			limitDetailMap["qps_thresholds"] = []interface{}{qpsThresholdsList}
		}

		if cngwRouteLimit.LimitDetail.LimitBy != nil {
			limitDetailMap["limit_by"] = cngwRouteLimit.LimitDetail.LimitBy
		}

		if cngwRouteLimit.LimitDetail.ResponseType != nil {
			limitDetailMap["response_type"] = cngwRouteLimit.LimitDetail.ResponseType
		}

		if cngwRouteLimit.LimitDetail.HideClientHeaders != nil {
			limitDetailMap["hide_client_headers"] = cngwRouteLimit.LimitDetail.HideClientHeaders
		}

		if cngwRouteLimit.LimitDetail.IsDelay != nil {
			limitDetailMap["is_delay"] = cngwRouteLimit.LimitDetail.IsDelay
		}

		if cngwRouteLimit.LimitDetail.Path != nil {
			limitDetailMap["path"] = cngwRouteLimit.LimitDetail.Path
		}

		if cngwRouteLimit.LimitDetail.Header != nil {
			limitDetailMap["header"] = cngwRouteLimit.LimitDetail.Header
		}

		if cngwRouteLimit.LimitDetail.ExternalRedis != nil {
			externalRedisMap := map[string]interface{}{}

			if cngwRouteLimit.LimitDetail.ExternalRedis.RedisHost != nil {
				externalRedisMap["redis_host"] = cngwRouteLimit.LimitDetail.ExternalRedis.RedisHost
			}

			if cngwRouteLimit.LimitDetail.ExternalRedis.RedisPassword != nil {
				externalRedisMap["redis_password"] = cngwRouteLimit.LimitDetail.ExternalRedis.RedisPassword
			}

			if cngwRouteLimit.LimitDetail.ExternalRedis.RedisPort != nil {
				externalRedisMap["redis_port"] = cngwRouteLimit.LimitDetail.ExternalRedis.RedisPort
			}

			if cngwRouteLimit.LimitDetail.ExternalRedis.RedisTimeout != nil {
				externalRedisMap["redis_timeout"] = cngwRouteLimit.LimitDetail.ExternalRedis.RedisTimeout
			}

			limitDetailMap["external_redis"] = []interface{}{externalRedisMap}
		}

		if cngwRouteLimit.LimitDetail.Policy != nil {
			limitDetailMap["policy"] = cngwRouteLimit.LimitDetail.Policy
		}

		if cngwRouteLimit.LimitDetail.RateLimitResponse != nil {
			rateLimitResponseMap := map[string]interface{}{}

			if cngwRouteLimit.LimitDetail.RateLimitResponse.Body != nil {
				rateLimitResponseMap["body"] = cngwRouteLimit.LimitDetail.RateLimitResponse.Body
			}

			if cngwRouteLimit.LimitDetail.RateLimitResponse.Headers != nil {
				headersList := []interface{}{}
				for _, headers := range cngwRouteLimit.LimitDetail.RateLimitResponse.Headers {
					headersMap := map[string]interface{}{}

					if headers.Key != nil {
						headersMap["key"] = headers.Key
					}

					if headers.Value != nil {
						headersMap["value"] = headers.Value
					}

					headersList = append(headersList, headersMap)
				}

				rateLimitResponseMap["headers"] = []interface{}{headersList}
			}

			if cngwRouteLimit.LimitDetail.RateLimitResponse.HttpStatus != nil {
				rateLimitResponseMap["http_status"] = cngwRouteLimit.LimitDetail.RateLimitResponse.HttpStatus
			}

			limitDetailMap["rate_limit_response"] = []interface{}{rateLimitResponseMap}
		}

		if cngwRouteLimit.LimitDetail.RateLimitResponseUrl != nil {
			limitDetailMap["rate_limit_response_url"] = cngwRouteLimit.LimitDetail.RateLimitResponseUrl
		}

		if cngwRouteLimit.LimitDetail.LineUpTime != nil {
			limitDetailMap["line_up_time"] = cngwRouteLimit.LimitDetail.LineUpTime
		}

		_ = d.Set("limit_detail", []interface{}{limitDetailMap})
	}

	return nil
}

func resourceTencentCloudTseCngwRouteLimitUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_route_limit.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tse.NewModifyCloudNativeAPIGatewayRouteRateLimitRequest()

	cngwRouteLimitId := d.Id()

	request.GatewayId = &gatewayId

	immutableArgs := []string{"gateway_id", "limit_detail"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("gateway_id") {
		if v, ok := d.GetOk("gateway_id"); ok {
			request.GatewayId = helper.String(v.(string))
		}
	}

	if d.HasChange("limit_detail") {
		if dMap, ok := helper.InterfacesHeadMap(d, "limit_detail"); ok {
			cloudNativeAPIGatewayRateLimitDetail := tse.CloudNativeAPIGatewayRateLimitDetail{}
			if v, ok := dMap["enabled"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.Enabled = helper.Bool(v.(bool))
			}
			if v, ok := dMap["qps_thresholds"]; ok {
				for _, item := range v.([]interface{}) {
					qpsThresholdsMap := item.(map[string]interface{})
					qpsThreshold := tse.QpsThreshold{}
					if v, ok := qpsThresholdsMap["unit"]; ok {
						qpsThreshold.Unit = helper.String(v.(string))
					}
					if v, ok := qpsThresholdsMap["max"]; ok {
						qpsThreshold.Max = helper.IntInt64(v.(int))
					}
					cloudNativeAPIGatewayRateLimitDetail.QpsThresholds = append(cloudNativeAPIGatewayRateLimitDetail.QpsThresholds, &qpsThreshold)
				}
			}
			if v, ok := dMap["limit_by"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.LimitBy = helper.String(v.(string))
			}
			if v, ok := dMap["response_type"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.ResponseType = helper.String(v.(string))
			}
			if v, ok := dMap["hide_client_headers"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.HideClientHeaders = helper.Bool(v.(bool))
			}
			if v, ok := dMap["is_delay"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.IsDelay = helper.Bool(v.(bool))
			}
			if v, ok := dMap["path"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.Path = helper.String(v.(string))
			}
			if v, ok := dMap["header"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.Header = helper.String(v.(string))
			}
			if externalRedisMap, ok := helper.InterfaceToMap(dMap, "external_redis"); ok {
				externalRedis := tse.ExternalRedis{}
				if v, ok := externalRedisMap["redis_host"]; ok {
					externalRedis.RedisHost = helper.String(v.(string))
				}
				if v, ok := externalRedisMap["redis_password"]; ok {
					externalRedis.RedisPassword = helper.String(v.(string))
				}
				if v, ok := externalRedisMap["redis_port"]; ok {
					externalRedis.RedisPort = helper.IntInt64(v.(int))
				}
				if v, ok := externalRedisMap["redis_timeout"]; ok {
					externalRedis.RedisTimeout = helper.IntInt64(v.(int))
				}
				cloudNativeAPIGatewayRateLimitDetail.ExternalRedis = &externalRedis
			}
			if v, ok := dMap["policy"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.Policy = helper.String(v.(string))
			}
			if rateLimitResponseMap, ok := helper.InterfaceToMap(dMap, "rate_limit_response"); ok {
				rateLimitResponse := tse.RateLimitResponse{}
				if v, ok := rateLimitResponseMap["body"]; ok {
					rateLimitResponse.Body = helper.String(v.(string))
				}
				if v, ok := rateLimitResponseMap["headers"]; ok {
					for _, item := range v.([]interface{}) {
						headersMap := item.(map[string]interface{})
						kVMapping := tse.KVMapping{}
						if v, ok := headersMap["key"]; ok {
							kVMapping.Key = helper.String(v.(string))
						}
						if v, ok := headersMap["value"]; ok {
							kVMapping.Value = helper.String(v.(string))
						}
						rateLimitResponse.Headers = append(rateLimitResponse.Headers, &kVMapping)
					}
				}
				if v, ok := rateLimitResponseMap["http_status"]; ok {
					rateLimitResponse.HttpStatus = helper.IntInt64(v.(int))
				}
				cloudNativeAPIGatewayRateLimitDetail.RateLimitResponse = &rateLimitResponse
			}
			if v, ok := dMap["rate_limit_response_url"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.RateLimitResponseUrl = helper.String(v.(string))
			}
			if v, ok := dMap["line_up_time"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.LineUpTime = helper.IntInt64(v.(int))
			}
			request.LimitDetail = &cloudNativeAPIGatewayRateLimitDetail
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().ModifyCloudNativeAPIGatewayRouteRateLimit(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tse cngwRouteLimit failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTseCngwRouteLimitRead(d, meta)
}

func resourceTencentCloudTseCngwRouteLimitDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_route_limit.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
	cngwRouteLimitId := d.Id()

	if err := service.DeleteTseCngwRouteLimitById(ctx, gatewayId); err != nil {
		return err
	}

	return nil
}
