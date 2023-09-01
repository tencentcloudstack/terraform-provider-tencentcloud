/*
Provides a resource to create a tse cngw_route_rate_limit

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_tse_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_tse_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_tse_cngw_gateway" "cngw_gateway" {
  description                = "terraform test1"
  enable_cls                 = true
  engine_region              = "ap-guangzhou"
  feature_version            = "STANDARD"
  gateway_version            = "2.5.1"
  ingress_class_name         = "tse-nginx-ingress"
  internet_max_bandwidth_out = 0
  name                       = "terraform-gateway1"
  trade_type                 = 0
  type                       = "kong"

  node_config {
    number        = 2
    specification = "1c2g"
  }

  vpc_config {
    subnet_id = tencentcloud_subnet.subnet.id
    vpc_id    = tencentcloud_vpc.vpc.id
  }

  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tse_cngw_service" "cngw_service" {
  gateway_id    = tencentcloud_tse_cngw_gateway.cngw_gateway.id
  name          = "terraform-test"
  path          = "/test"
  protocol      = "http"
  retries       = 5
  timeout       = 60000
  upstream_type = "HostIP"

  upstream_info {
    algorithm             = "round-robin"
    auto_scaling_cvm_port = 0
    host                  = "arunma.cn"
    port                  = 8012
    slow_start            = 0
  }
}

resource "tencentcloud_tse_cngw_route" "cngw_route" {
  destination_ports = []
  force_https       = false
  gateway_id        = tencentcloud_tse_cngw_gateway.cngw_gateway.id
  hosts = [
    "192.168.0.1:9090",
  ]
  https_redirect_status_code = 426
  paths = [
    "/user",
  ]
  headers {
		key = "req"
		value = "terraform"
  }
  preserve_host = false
  protocols = [
    "http",
    "https",
  ]
  route_name = "terraform-route"
  service_id = tencentcloud_tse_cngw_service.cngw_service.service_id
  strip_path = true
}

resource "tencentcloud_tse_cngw_route_rate_limit" "cngw_route_rate_limit" {
    gateway_id = tencentcloud_tse_cngw_gateway.cngw_gateway.id
    route_id   = tencentcloud_tse_cngw_route.cngw_route.route_id

    limit_detail {
        enabled             = true
        header              = "req"
        hide_client_headers = true
        is_delay            = true
        limit_by            = "header"
        line_up_time        = 10
        policy              = "redis"
        response_type       = "default"

        qps_thresholds {
            max  = 10
            unit = "minute"
        }
    }
}
```

Import

tse cngw_route_rate_limit can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_route_rate_limit.cngw_route_rate_limit gatewayId#routeId
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTseCngwRouteRateLimit() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseCngwRouteRateLimitCreate,
		Read:   resourceTencentCloudTseCngwRouteRateLimitRead,
		Update: resourceTencentCloudTseCngwRouteRateLimitUpdate,
		Delete: resourceTencentCloudTseCngwRouteRateLimitDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway ID.",
			},

			"route_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Route id, or route name.",
			},

			"limit_detail": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "rate limit configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "status of service rate limit.",
						},
						"qps_thresholds": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "qps threshold.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"unit": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "qps threshold unit.Reference value:`second`,`minute`,`hour`,`day`,`month`,`year`.",
									},
									"max": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "the max threshold.",
									},
								},
							},
						},
						"limit_by": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "basis for service rate limit.Reference value:`ip`,`service`,`consumer`,`credential`,`path`,`header`.",
						},
						"response_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "response strategy.Reference value:`url`: forward request according to url,`text`: response configuration,`default`: return directly.",
						},
						"hide_client_headers": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "whether to hide the headers of client.",
						},
						"is_delay": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "whether to enable request queuing.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "request paths that require rate limit.",
						},
						"header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "request headers that require rate limit.",
						},
						"external_redis": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "external redis information, maybe null.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"redis_host": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "redis ip, maybe null.",
									},
									"redis_password": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "redis password, maybe null.",
									},
									"redis_port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "redis port, maybe null.",
									},
									"redis_timeout": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "redis timeout, unit: `ms`, maybe null.",
									},
								},
							},
						},
						"policy": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "counter policy.Reference value:`local`,`redis`,`external_redis`.",
						},
						"rate_limit_response": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "response configuration, the response strategy is text, maybe null.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"body": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "custom response body, maybe bull.",
									},
									"headers": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "headrs.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "key of header.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "value of header.",
												},
											},
										},
									},
									"http_status": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "http status code.",
									},
								},
							},
						},
						"rate_limit_response_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "request forwarding address, maybe null.",
						},
						"line_up_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "queue time.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTseCngwRouteRateLimitCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_route_rate_limit.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tse.NewCreateCloudNativeAPIGatewayRouteRateLimitRequest()
		gatewayId string
		routeId   string
	)
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
		request.GatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("route_id"); ok {
		routeId = v.(string)
		request.Id = helper.String(v.(string))
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
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tse cngwRouteRateLimit failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(gatewayId + FILED_SP + routeId)

	return resourceTencentCloudTseCngwRouteRateLimitRead(d, meta)
}

func resourceTencentCloudTseCngwRouteRateLimitRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_route_rate_limit.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	routeId := idSplit[1]

	cngwRouteRateLimit, err := service.DescribeTseCngwRouteRateLimitById(ctx, gatewayId, routeId)
	if err != nil {
		return err
	}

	if cngwRouteRateLimit == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseCngwRouteRateLimit` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("gateway_id", gatewayId)
	_ = d.Set("route_id", routeId)

	if cngwRouteRateLimit != nil {
		limitDetailMap := map[string]interface{}{}

		if cngwRouteRateLimit.Enabled != nil {
			limitDetailMap["enabled"] = cngwRouteRateLimit.Enabled
		}

		if cngwRouteRateLimit.QpsThresholds != nil {
			qpsThresholdsList := []interface{}{}
			for _, qpsThresholds := range cngwRouteRateLimit.QpsThresholds {
				qpsThresholdsMap := map[string]interface{}{}

				if qpsThresholds.Unit != nil {
					qpsThresholdsMap["unit"] = qpsThresholds.Unit
				}

				if qpsThresholds.Max != nil {
					qpsThresholdsMap["max"] = qpsThresholds.Max
				}

				qpsThresholdsList = append(qpsThresholdsList, qpsThresholdsMap)
			}

			limitDetailMap["qps_thresholds"] = qpsThresholdsList
		}

		if cngwRouteRateLimit.LimitBy != nil {
			limitDetailMap["limit_by"] = cngwRouteRateLimit.LimitBy
		}

		if cngwRouteRateLimit.ResponseType != nil {
			limitDetailMap["response_type"] = cngwRouteRateLimit.ResponseType
		}

		if cngwRouteRateLimit.HideClientHeaders != nil {
			limitDetailMap["hide_client_headers"] = cngwRouteRateLimit.HideClientHeaders
		}

		if cngwRouteRateLimit.IsDelay != nil {
			limitDetailMap["is_delay"] = cngwRouteRateLimit.IsDelay
		}

		if cngwRouteRateLimit.Path != nil {
			limitDetailMap["path"] = cngwRouteRateLimit.Path
		}

		if cngwRouteRateLimit.Header != nil {
			limitDetailMap["header"] = cngwRouteRateLimit.Header
		}

		if cngwRouteRateLimit.ExternalRedis != nil {
			externalRedisMap := map[string]interface{}{}

			if cngwRouteRateLimit.ExternalRedis.RedisHost != nil {
				externalRedisMap["redis_host"] = cngwRouteRateLimit.ExternalRedis.RedisHost
			}

			if cngwRouteRateLimit.ExternalRedis.RedisPassword != nil {
				externalRedisMap["redis_password"] = cngwRouteRateLimit.ExternalRedis.RedisPassword
			}

			if cngwRouteRateLimit.ExternalRedis.RedisPort != nil {
				externalRedisMap["redis_port"] = cngwRouteRateLimit.ExternalRedis.RedisPort
			}

			if cngwRouteRateLimit.ExternalRedis.RedisTimeout != nil {
				externalRedisMap["redis_timeout"] = cngwRouteRateLimit.ExternalRedis.RedisTimeout
			}

			limitDetailMap["external_redis"] = []interface{}{externalRedisMap}
		}

		if cngwRouteRateLimit.Policy != nil {
			limitDetailMap["policy"] = cngwRouteRateLimit.Policy
		}

		if cngwRouteRateLimit.RateLimitResponse != nil {
			rateLimitResponseMap := map[string]interface{}{}

			if cngwRouteRateLimit.RateLimitResponse.Body != nil {
				rateLimitResponseMap["body"] = cngwRouteRateLimit.RateLimitResponse.Body
			}

			if cngwRouteRateLimit.RateLimitResponse.Headers != nil {
				headersList := []interface{}{}
				for _, headers := range cngwRouteRateLimit.RateLimitResponse.Headers {
					headersMap := map[string]interface{}{}

					if headers.Key != nil {
						headersMap["key"] = headers.Key
					}

					if headers.Value != nil {
						headersMap["value"] = headers.Value
					}

					headersList = append(headersList, headersMap)
				}

				rateLimitResponseMap["headers"] = headersList
			}

			if cngwRouteRateLimit.RateLimitResponse.HttpStatus != nil {
				rateLimitResponseMap["http_status"] = cngwRouteRateLimit.RateLimitResponse.HttpStatus
			}

			limitDetailMap["rate_limit_response"] = []interface{}{rateLimitResponseMap}
		}

		if cngwRouteRateLimit.RateLimitResponseUrl != nil {
			limitDetailMap["rate_limit_response_url"] = cngwRouteRateLimit.RateLimitResponseUrl
		}

		if cngwRouteRateLimit.LineUpTime != nil {
			limitDetailMap["line_up_time"] = cngwRouteRateLimit.LineUpTime
		}

		_ = d.Set("limit_detail", []interface{}{limitDetailMap})
	}

	return nil
}

func resourceTencentCloudTseCngwRouteRateLimitUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_route_rate_limit.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tse.NewModifyCloudNativeAPIGatewayRouteRateLimitRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	routeId := idSplit[1]

	request.GatewayId = &gatewayId
	request.Id = &routeId

	immutableArgs := []string{"gateway_id", "route_id", "limit_detail"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
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
		log.Printf("[CRITAL]%s update tse cngwRouteRateLimit failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTseCngwRouteRateLimitRead(d, meta)
}

func resourceTencentCloudTseCngwRouteRateLimitDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_route_rate_limit.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	routeId := idSplit[1]

	if err := service.DeleteTseCngwRouteRateLimitById(ctx, gatewayId, routeId); err != nil {
		return err
	}

	return nil
}
