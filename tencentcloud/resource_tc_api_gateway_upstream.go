package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAPIGatewayUpstream() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayUpstreamCreate,
		Read:   resourceTencentCloudAPIGatewayUpstreamRead,
		Update: resourceTencentCloudAPIGatewayUpstreamUpdate,
		Delete: resourceTencentCloudAPIGatewayUpstreamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"scheme": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_UPSTREAM_SCHEME),
				Description:  "Backend protocol, value range: HTTP, HTTPS, gRPC, gRPCs.",
			},
			"algorithm": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Load balancing algorithm, value range: ROUND-ROBIN.",
			},
			"uniq_vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "VPC Unique ID.",
			},
			"upstream_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Backend channel name.",
			},
			"upstream_description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Backend channel description.",
			},
			"upstream_type": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_UPSTREAM_TYPE),
				Description:  "Backend access type, value range: IP_PORT, K8S.",
			},
			"retries": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Request retry count, default to 3 times.",
			},
			"upstream_host": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Host request header forwarded by gateway to backend.",
			},
			"nodes": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Backend nodes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP or domain name.",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Port [0, 65535].",
						},
						"weight": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Weight [0, 100], 0 is disabled.",
						},
						"vm_instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CVM instance IDNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"tags": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "Dye labelNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						//"healthy": {
						//	Type:        schema.TypeString,
						//	Optional:    true,
						//	Description: "The node health status does not need to be passed during creation or editing. OFF: Off, HEALTHY: Healthy, UNHEALTHY: Abnormal, NO_ DATA: Data not reported. Currently, only VPC channels are supported.Note: This field may return null, indicating that a valid value cannot be obtained.",
						//},
						"service_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "K8S container service nameNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"name_space": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "K8S namespaceNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the TKE clusterNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"source": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Source of Node, value range: K8SNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"unique_service_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Unique service name recorded internally by API gatewayNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},
			"health_checker": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Health check configuration, currently only supports VPC channels.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_active_check": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Identify whether active health checks are enabled.",
						},
						"enable_passive_check": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Identify whether passive health checks are enabled.",
						},
						"healthy_http_status": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The HTTP status code that determines a successful request during a health check.",
						},
						"unhealthy_http_status": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The HTTP status code that determines a failed request during a health check.",
						},
						"tcp_failure_threshold": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "TCP continuous error threshold. 0 indicates disabling TCP checking. Value range: [0, 254].",
						},
						"timeout_threshold": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Continuous timeout threshold. 0 indicates disabling timeout checking. Value range: [0, 254].",
						},
						"http_failure_threshold": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "HTTP continuous error threshold. 0 means HTTP checking is disabled. Value range: [0, 254].",
						},
						"active_check_http_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Detect the requested path during active health checks. The default is&#39;/&#39;.",
						},
						"active_check_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The detection request for active health check timed out in seconds. The default is 5 seconds.",
						},
						"active_check_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The time interval for active health checks is 5 seconds by default.",
						},
						"unhealthy_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The automatic recovery time of abnormal node status, in seconds. When only passive checking is enabled, it must be set to a value&gt;0, otherwise the passive exception node will not be able to recover. The default is 30 seconds.",
						},
					},
				},
			},
			"k8s_service": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Configuration of K8S container service.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"weight": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "weight.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "K8s cluster ID.",
						},
						"namespace": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Container namespace.",
						},
						"service_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the container service.",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Port of service.",
						},
						"extra_labels": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Additional Selected Pod Label.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Key of Label.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Value of Label.",
									},
								},
							},
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Customized service name, optional.",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayUpstreamCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_upstream.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId        = getLogId(contextNil)
		request      = apigateway.NewCreateUpstreamRequest()
		response     = apigateway.NewCreateUpstreamResponse()
		upstreamId   string
		upstreamType string
	)

	if v, ok := d.GetOk("scheme"); ok {
		request.Scheme = helper.String(v.(string))
	}

	if v, ok := d.GetOk("algorithm"); ok {
		request.Algorithm = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_vpc_id"); ok {
		request.UniqVpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("upstream_name"); ok {
		request.UpstreamName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("upstream_description"); ok {
		request.UpstreamDescription = helper.String(v.(string))
	}

	if v, ok := d.GetOk("upstream_type"); ok {
		request.UpstreamType = helper.String(v.(string))
		upstreamType = v.(string)
	}

	if v, ok := d.GetOkExists("retries"); ok {
		request.Retries = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("upstream_host"); ok {
		request.UpstreamHost = helper.String(v.(string))
	}

	if v, ok := d.GetOk("nodes"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			upstreamNode := apigateway.UpstreamNode{}
			if upstreamType == API_GATEWAY_UPSTREAM_TYPE_IP {
				if v, ok := dMap["host"]; ok {
					upstreamNode.Host = helper.String(v.(string))
				}

				if v, ok := dMap["port"]; ok {
					upstreamNode.Port = helper.IntUint64(v.(int))
				}

				if v, ok := dMap["weight"]; ok {
					upstreamNode.Weight = helper.IntUint64(v.(int))
				}

				if v, ok := dMap["vm_instance_id"]; ok {
					upstreamNode.VmInstanceId = helper.String(v.(string))
				}

				if v, ok := dMap["tags"]; ok {
					tagsSet := v.(*schema.Set).List()
					for i := range tagsSet {
						tags := tagsSet[i].(string)
						upstreamNode.Tags = append(upstreamNode.Tags, &tags)
					}
				}

				//if v, ok := dMap["healthy"]; ok {
				//	upstreamNode.Healthy = helper.String(v.(string))
				//}

			} else {
				if v, ok := dMap["service_name"]; ok {
					upstreamNode.ServiceName = helper.String(v.(string))
				}

				if v, ok := dMap["name_space"]; ok {
					upstreamNode.NameSpace = helper.String(v.(string))
				}

				if v, ok := dMap["cluster_id"]; ok {
					upstreamNode.ClusterId = helper.String(v.(string))
				}

				if v, ok := dMap["source"]; ok {
					upstreamNode.Source = helper.String(v.(string))
				}

				if v, ok := dMap["unique_service_name"]; ok {
					upstreamNode.UniqueServiceName = helper.String(v.(string))
				}
			}

			request.Nodes = append(request.Nodes, &upstreamNode)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "health_checker"); ok {
		upstreamHealthChecker := apigateway.UpstreamHealthChecker{}
		if v, ok := dMap["enable_active_check"]; ok {
			upstreamHealthChecker.EnableActiveCheck = helper.Bool(v.(bool))
		}

		if v, ok := dMap["enable_passive_check"]; ok {
			upstreamHealthChecker.EnablePassiveCheck = helper.Bool(v.(bool))
		}

		if v, ok := dMap["healthy_http_status"]; ok {
			upstreamHealthChecker.HealthyHttpStatus = helper.String(v.(string))
		}

		if v, ok := dMap["unhealthy_http_status"]; ok {
			upstreamHealthChecker.UnhealthyHttpStatus = helper.String(v.(string))
		}

		if v, ok := dMap["tcp_failure_threshold"]; ok {
			upstreamHealthChecker.TcpFailureThreshold = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["timeout_threshold"]; ok {
			upstreamHealthChecker.TimeoutThreshold = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["http_failure_threshold"]; ok {
			upstreamHealthChecker.HttpFailureThreshold = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["active_check_http_path"]; ok {
			upstreamHealthChecker.ActiveCheckHttpPath = helper.String(v.(string))
		}

		if v, ok := dMap["active_check_timeout"]; ok {
			upstreamHealthChecker.ActiveCheckTimeout = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["active_check_interval"]; ok {
			upstreamHealthChecker.ActiveCheckInterval = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["unhealthy_timeout"]; ok {
			upstreamHealthChecker.UnhealthyTimeout = helper.IntUint64(v.(int))
		}

		request.HealthChecker = &upstreamHealthChecker
	}

	if v, ok := d.GetOk("k8s_service"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			k8sService := apigateway.K8sService{}
			if v, ok := dMap["weight"]; ok {
				k8sService.Weight = helper.IntInt64(v.(int))
			}

			if v, ok := dMap["cluster_id"]; ok {
				k8sService.ClusterId = helper.String(v.(string))
			}

			if v, ok := dMap["namespace"]; ok {
				k8sService.Namespace = helper.String(v.(string))
			}

			if v, ok := dMap["service_name"]; ok {
				k8sService.ServiceName = helper.String(v.(string))
			}

			if v, ok := dMap["port"]; ok {
				k8sService.Port = helper.IntInt64(v.(int))
			}

			if v, ok := dMap["extra_labels"]; ok {
				for _, item := range v.([]interface{}) {
					extraLabelsMap := item.(map[string]interface{})
					k8sLabel := apigateway.K8sLabel{}
					if v, ok := extraLabelsMap["key"]; ok {
						k8sLabel.Key = helper.String(v.(string))
					}

					if v, ok := extraLabelsMap["value"]; ok {
						k8sLabel.Value = helper.String(v.(string))
					}

					k8sService.ExtraLabels = append(k8sService.ExtraLabels, &k8sLabel)
				}
			}

			if v, ok := dMap["name"]; ok {
				k8sService.Name = helper.String(v.(string))
			}

			request.K8sService = append(request.K8sService, &k8sService)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAPIGatewayClient().CreateUpstream(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create apigateway upstream failed, reason:%+v", logId, err)
		return err
	}

	upstreamId = *response.Response.UpstreamId

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::apigateway:%s:uin/:upstreamId/%s", region, upstreamId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	d.SetId(upstreamId)
	return resourceTencentCloudAPIGatewayUpstreamRead(d, meta)
}

func resourceTencentCloudAPIGatewayUpstreamRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_upstream.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId        = getLogId(contextNil)
		ctx          = context.WithValue(context.TODO(), logIdKey, logId)
		service      = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		upstreamId   = d.Id()
		upstreamType string
	)

	upstreamInfo, err := service.DescribeApigatewayUpstreamById(ctx, upstreamId)
	if err != nil {
		return err
	}

	if upstreamInfo == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ApigatewayUpstream` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if upstreamInfo.Scheme != nil {
		_ = d.Set("scheme", upstreamInfo.Scheme)
	}

	if upstreamInfo.Algorithm != nil {
		_ = d.Set("algorithm", upstreamInfo.Algorithm)
	}

	if upstreamInfo.UniqVpcId != nil {
		_ = d.Set("uniq_vpc_id", upstreamInfo.UniqVpcId)
	}

	if upstreamInfo.UpstreamName != nil {
		_ = d.Set("upstream_name", upstreamInfo.UpstreamName)
	}

	if upstreamInfo.UpstreamDescription != nil {
		_ = d.Set("upstream_description", upstreamInfo.UpstreamDescription)
	}

	if upstreamInfo.UpstreamType != nil {
		_ = d.Set("upstream_type", upstreamInfo.UpstreamType)
		upstreamType = *upstreamInfo.UpstreamType
	}

	if upstreamInfo.Retries != nil {
		_ = d.Set("retries", upstreamInfo.Retries)
	}

	if upstreamInfo.UpstreamHost != nil {
		_ = d.Set("upstream_host", upstreamInfo.UpstreamHost)
	}

	if upstreamInfo.Nodes != nil {
		nodesList := []interface{}{}
		for _, nodes := range upstreamInfo.Nodes {
			nodesMap := map[string]interface{}{}
			if upstreamType == API_GATEWAY_UPSTREAM_TYPE_IP {
				if nodes.Host != nil {
					nodesMap["host"] = nodes.Host
				}

				if nodes.Port != nil {
					nodesMap["port"] = nodes.Port
				}

				if nodes.Weight != nil {
					nodesMap["weight"] = nodes.Weight
				}

				if nodes.VmInstanceId != nil {
					nodesMap["vm_instance_id"] = nodes.VmInstanceId
				}

				if nodes.Tags != nil {
					nodesMap["tags"] = nodes.Tags
				}

				//if nodes.Healthy != nil {
				//	nodesMap["healthy"] = nodes.Healthy
				//}

			} else {
				if nodes.ServiceName != nil {
					nodesMap["service_name"] = nodes.ServiceName
				}

				if nodes.NameSpace != nil {
					nodesMap["name_space"] = nodes.NameSpace
				}

				if nodes.ClusterId != nil {
					nodesMap["cluster_id"] = nodes.ClusterId
				}

				if nodes.Source != nil {
					nodesMap["source"] = nodes.Source
				}

				if nodes.UniqueServiceName != nil {
					nodesMap["unique_service_name"] = nodes.UniqueServiceName
				}
			}

			nodesList = append(nodesList, nodesMap)
		}

		_ = d.Set("nodes", nodesList)

	}

	if upstreamInfo.HealthChecker != nil {
		healthCheckerMap := map[string]interface{}{}

		if upstreamInfo.HealthChecker.EnableActiveCheck != nil {
			healthCheckerMap["enable_active_check"] = upstreamInfo.HealthChecker.EnableActiveCheck
		}

		if upstreamInfo.HealthChecker.EnablePassiveCheck != nil {
			healthCheckerMap["enable_passive_check"] = upstreamInfo.HealthChecker.EnablePassiveCheck
		}

		if upstreamInfo.HealthChecker.HealthyHttpStatus != nil {
			healthCheckerMap["healthy_http_status"] = upstreamInfo.HealthChecker.HealthyHttpStatus
		}

		if upstreamInfo.HealthChecker.UnhealthyHttpStatus != nil {
			healthCheckerMap["unhealthy_http_status"] = upstreamInfo.HealthChecker.UnhealthyHttpStatus
		}

		if upstreamInfo.HealthChecker.TcpFailureThreshold != nil {
			healthCheckerMap["tcp_failure_threshold"] = upstreamInfo.HealthChecker.TcpFailureThreshold
		}

		if upstreamInfo.HealthChecker.TimeoutThreshold != nil {
			healthCheckerMap["timeout_threshold"] = upstreamInfo.HealthChecker.TimeoutThreshold
		}

		if upstreamInfo.HealthChecker.HttpFailureThreshold != nil {
			healthCheckerMap["http_failure_threshold"] = upstreamInfo.HealthChecker.HttpFailureThreshold
		}

		if upstreamInfo.HealthChecker.ActiveCheckHttpPath != nil {
			healthCheckerMap["active_check_http_path"] = upstreamInfo.HealthChecker.ActiveCheckHttpPath
		}

		if upstreamInfo.HealthChecker.ActiveCheckTimeout != nil {
			healthCheckerMap["active_check_timeout"] = upstreamInfo.HealthChecker.ActiveCheckTimeout
		}

		if upstreamInfo.HealthChecker.ActiveCheckInterval != nil {
			healthCheckerMap["active_check_interval"] = upstreamInfo.HealthChecker.ActiveCheckInterval
		}

		if upstreamInfo.HealthChecker.ActiveRequestHeader != nil {
			healthCheckerMap["active_request_header"] = upstreamInfo.HealthChecker.ActiveRequestHeader
		}

		if upstreamInfo.HealthChecker.UnhealthyTimeout != nil {
			healthCheckerMap["unhealthy_timeout"] = upstreamInfo.HealthChecker.UnhealthyTimeout
		}

		_ = d.Set("health_checker", []interface{}{healthCheckerMap})
	}

	if upstreamInfo.K8sServices != nil {
		k8sServiceList := []interface{}{}
		for _, k8sService := range upstreamInfo.K8sServices {
			k8sServiceMap := map[string]interface{}{}

			if k8sService.Weight != nil {
				k8sServiceMap["weight"] = k8sService.Weight
			}

			if k8sService.ClusterId != nil {
				k8sServiceMap["cluster_id"] = k8sService.ClusterId
			}

			if k8sService.Namespace != nil {
				k8sServiceMap["namespace"] = k8sService.Namespace
			}

			if k8sService.ServiceName != nil {
				k8sServiceMap["service_name"] = k8sService.ServiceName
			}

			if k8sService.Port != nil {
				k8sServiceMap["port"] = k8sService.Port
			}

			if k8sService.ExtraLabels != nil {
				extraLabelsList := []interface{}{}
				for _, extraLabels := range k8sService.ExtraLabels {
					extraLabelsMap := map[string]interface{}{}

					if extraLabels.Key != nil {
						extraLabelsMap["key"] = extraLabels.Key
					}

					if extraLabels.Value != nil {
						extraLabelsMap["value"] = extraLabels.Value
					}

					extraLabelsList = append(extraLabelsList, extraLabelsMap)
				}

				k8sServiceMap["extra_labels"] = extraLabelsList
			}

			if k8sService.Name != nil {
				k8sServiceMap["name"] = k8sService.Name
			}

			k8sServiceList = append(k8sServiceList, k8sServiceMap)
		}

		_ = d.Set("k8s_service", k8sServiceList)

	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "apigateway", "upstreamId", tcClient.Region, upstreamId)
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudAPIGatewayUpstreamUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_upstream.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId        = getLogId(contextNil)
		request      = apigateway.NewModifyUpstreamRequest()
		upstreamId   = d.Id()
		upstreamType string
	)

	request.UpstreamId = &upstreamId

	if d.HasChange("scheme") {
		if v, ok := d.GetOk("scheme"); ok {
			request.Scheme = helper.String(v.(string))
		}
	}

	if d.HasChange("algorithm") {
		if v, ok := d.GetOk("algorithm"); ok {
			request.Algorithm = helper.String(v.(string))
		}
	}

	if d.HasChange("uniq_vpc_id") {
		if v, ok := d.GetOk("uniq_vpc_id"); ok {
			request.UniqVpcId = helper.String(v.(string))
		}
	}

	if d.HasChange("upstream_name") {
		if v, ok := d.GetOk("upstream_name"); ok {
			request.UpstreamName = helper.String(v.(string))
		}
	}

	if d.HasChange("upstream_description") {
		if v, ok := d.GetOk("upstream_description"); ok {
			request.UpstreamDescription = helper.String(v.(string))
		}
	}

	if v, ok := d.GetOk("upstream_type"); ok {
		request.UpstreamType = helper.String(v.(string))
		upstreamType = v.(string)
	}

	if d.HasChange("retries") {
		if v, ok := d.GetOkExists("retries"); ok {
			request.Retries = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("upstream_host") {
		if v, ok := d.GetOk("upstream_host"); ok {
			request.UpstreamHost = helper.String(v.(string))
		}
	}

	if d.HasChange("nodes") {
		if v, ok := d.GetOk("nodes"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				upstreamNode := apigateway.UpstreamNode{}
				if upstreamType == API_GATEWAY_UPSTREAM_TYPE_IP {
					if v, ok := dMap["host"]; ok {
						upstreamNode.Host = helper.String(v.(string))
					}

					if v, ok := dMap["port"]; ok {
						upstreamNode.Port = helper.IntUint64(v.(int))
					}

					if v, ok := dMap["weight"]; ok {
						upstreamNode.Weight = helper.IntUint64(v.(int))
					}

					if v, ok := dMap["vm_instance_id"]; ok {
						upstreamNode.VmInstanceId = helper.String(v.(string))
					}

					if v, ok := dMap["tags"]; ok {
						tagsSet := v.(*schema.Set).List()
						for i := range tagsSet {
							tags := tagsSet[i].(string)
							upstreamNode.Tags = append(upstreamNode.Tags, &tags)
						}
					}

					//if v, ok := dMap["healthy"]; ok {
					//	upstreamNode.Healthy = helper.String(v.(string))
					//}

				} else {
					if v, ok := dMap["service_name"]; ok {
						upstreamNode.ServiceName = helper.String(v.(string))
					}

					if v, ok := dMap["name_space"]; ok {
						upstreamNode.NameSpace = helper.String(v.(string))
					}

					if v, ok := dMap["cluster_id"]; ok {
						upstreamNode.ClusterId = helper.String(v.(string))
					}

					if v, ok := dMap["source"]; ok {
						upstreamNode.Source = helper.String(v.(string))
					}

					if v, ok := dMap["unique_service_name"]; ok {
						upstreamNode.UniqueServiceName = helper.String(v.(string))
					}
				}

				request.Nodes = append(request.Nodes, &upstreamNode)
			}
		}
	}

	if d.HasChange("health_checker") {
		if dMap, ok := helper.InterfacesHeadMap(d, "health_checker"); ok {
			upstreamHealthChecker := apigateway.UpstreamHealthChecker{}
			if v, ok := dMap["enable_active_check"]; ok {
				upstreamHealthChecker.EnableActiveCheck = helper.Bool(v.(bool))
			}

			if v, ok := dMap["enable_passive_check"]; ok {
				upstreamHealthChecker.EnablePassiveCheck = helper.Bool(v.(bool))
			}

			if v, ok := dMap["healthy_http_status"]; ok {
				upstreamHealthChecker.HealthyHttpStatus = helper.String(v.(string))
			}

			if v, ok := dMap["unhealthy_http_status"]; ok {
				upstreamHealthChecker.UnhealthyHttpStatus = helper.String(v.(string))
			}

			if v, ok := dMap["tcp_failure_threshold"]; ok {
				upstreamHealthChecker.TcpFailureThreshold = helper.IntUint64(v.(int))
			}

			if v, ok := dMap["timeout_threshold"]; ok {
				upstreamHealthChecker.TimeoutThreshold = helper.IntUint64(v.(int))
			}

			if v, ok := dMap["http_failure_threshold"]; ok {
				upstreamHealthChecker.HttpFailureThreshold = helper.IntUint64(v.(int))
			}

			if v, ok := dMap["active_check_http_path"]; ok {
				upstreamHealthChecker.ActiveCheckHttpPath = helper.String(v.(string))
			}

			if v, ok := dMap["active_check_timeout"]; ok {
				upstreamHealthChecker.ActiveCheckTimeout = helper.IntUint64(v.(int))
			}

			if v, ok := dMap["active_check_interval"]; ok {
				upstreamHealthChecker.ActiveCheckInterval = helper.IntUint64(v.(int))
			}

			if v, ok := dMap["unhealthy_timeout"]; ok {
				upstreamHealthChecker.UnhealthyTimeout = helper.IntUint64(v.(int))
			}

			request.HealthChecker = &upstreamHealthChecker
		}
	}

	if d.HasChange("k8s_service") {
		if v, ok := d.GetOk("k8s_service"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				k8sService := apigateway.K8sService{}
				if v, ok := dMap["weight"]; ok {
					k8sService.Weight = helper.IntInt64(v.(int))
				}

				if v, ok := dMap["cluster_id"]; ok {
					k8sService.ClusterId = helper.String(v.(string))
				}

				if v, ok := dMap["namespace"]; ok {
					k8sService.Namespace = helper.String(v.(string))
				}

				if v, ok := dMap["service_name"]; ok {
					k8sService.ServiceName = helper.String(v.(string))
				}

				if v, ok := dMap["port"]; ok {
					k8sService.Port = helper.IntInt64(v.(int))
				}

				if v, ok := dMap["extra_labels"]; ok {
					for _, item := range v.([]interface{}) {
						extraLabelsMap := item.(map[string]interface{})
						k8sLabel := apigateway.K8sLabel{}
						if v, ok := extraLabelsMap["key"]; ok {
							k8sLabel.Key = helper.String(v.(string))
						}

						if v, ok := extraLabelsMap["value"]; ok {
							k8sLabel.Value = helper.String(v.(string))
						}

						k8sService.ExtraLabels = append(k8sService.ExtraLabels, &k8sLabel)
					}
				}

				if v, ok := dMap["name"]; ok {
					k8sService.Name = helper.String(v.(string))
				}

				request.K8sService = append(request.K8sService, &k8sService)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAPIGatewayClient().ModifyUpstream(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update apigateway upstream failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("apigateway", "upstreamId", tcClient.Region, upstreamId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudAPIGatewayUpstreamRead(d, meta)
}

func resourceTencentCloudAPIGatewayUpstreamDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_upstream.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		upstreamId = d.Id()
	)

	if err := service.DeleteApigatewayUpstreamById(ctx, upstreamId); err != nil {
		return err
	}

	return nil
}
