/*
Provides a resource to create a apigateway upstream

Example Usage

```hcl
resource "tencentcloud_apigateway_upstream" "upstream" {
  scheme = ""
  algorithm = ""
  uniq_vpc_id = ""
  upstream_name = ""
  upstream_description = ""
  upstream_type = ""
  retries =
  upstream_host = ""
  nodes {
		host = ""
		port =
		weight =
		vm_instance_id = ""
		tags =
		healthy = ""
		service_name = ""
		name_space = ""
		cluster_id = ""
		source = ""
		unique_service_name = ""

  }
  health_checker {
		enable_active_check =
		enable_passive_check =
		healthy_http_status = ""
		unhealthy_http_status = ""
		tcp_failure_threshold =
		timeout_threshold =
		http_failure_threshold =
		active_check_http_path = ""
		active_check_timeout =
		active_check_interval =
		active_request_header =
		unhealthy_timeout =

  }
  k8s_service {
		weight =
		cluster_id = ""
		namespace = ""
		service_name = ""
		port =
		extra_labels {
			key = ""
			value = ""
		}
		name = ""

  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

apigateway upstream can be imported using the id, e.g.

```
terraform import tencentcloud_apigateway_upstream.upstream upstream_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudApigatewayUpstream() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApigatewayUpstreamCreate,
		Read:   resourceTencentCloudApigatewayUpstreamRead,
		Update: resourceTencentCloudApigatewayUpstreamUpdate,
		Delete: resourceTencentCloudApigatewayUpstreamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"scheme": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Backend protocol, value range: HTTP, HTTPS.",
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
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Backend access type, value range: IP_ PORT, K8S.",
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
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Dye labelNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"healthy": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The node health status does not need to be passed during creation or editing. OFF: Off, HEALTHY: Healthy, UNHEALTHY: Abnormal, NO_ DATA: Data not reported. Currently, only VPC channels are supported.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
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
							Description: "Detect the requested path during active health checks. The default is&amp;#39;/&amp;#39;.",
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
						"active_request_header": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The request header for detecting requests during active health checks.",
						},
						"unhealthy_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The automatic recovery time of abnormal node status, in seconds. When only passive checking is enabled, it must be set to a value&amp;gt;0, otherwise the passive exception node will not be able to recover. The default is 30 seconds.",
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
							Description: "Weight.",
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

func resourceTencentCloudApigatewayUpstreamCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_upstream.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = apigateway.NewCreateUpstreamRequest()
		response   = apigateway.NewCreateUpstreamResponse()
		upstreamId string
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
			if v, ok := dMap["healthy"]; ok {
				upstreamNode.Healthy = helper.String(v.(string))
			}
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
		if v, ok := dMap["active_request_header"]; ok {
			for _, item := range v.([]interface{}) {
				activeRequestHeaderMap := item.(map[string]interface{})
				upstreamHealthCheckerReqHeaders := apigateway.UpstreamHealthCheckerReqHeaders{}
				upstreamHealthChecker.ActiveRequestHeader = append(upstreamHealthChecker.ActiveRequestHeader, &upstreamHealthCheckerReqHeaders)
			}
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseApigatewayClient().CreateUpstream(request)
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
	d.SetId(upstreamId)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::apigw:%s:uin/:upstreamId/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudApigatewayUpstreamRead(d, meta)
}

func resourceTencentCloudApigatewayUpstreamRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_upstream.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

	upstreamId := d.Id()

	upstream, err := service.DescribeApigatewayUpstreamById(ctx, upstreamId)
	if err != nil {
		return err
	}

	if upstream == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ApigatewayUpstream` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if upstream.Scheme != nil {
		_ = d.Set("scheme", upstream.Scheme)
	}

	if upstream.Algorithm != nil {
		_ = d.Set("algorithm", upstream.Algorithm)
	}

	if upstream.UniqVpcId != nil {
		_ = d.Set("uniq_vpc_id", upstream.UniqVpcId)
	}

	if upstream.UpstreamName != nil {
		_ = d.Set("upstream_name", upstream.UpstreamName)
	}

	if upstream.UpstreamDescription != nil {
		_ = d.Set("upstream_description", upstream.UpstreamDescription)
	}

	if upstream.UpstreamType != nil {
		_ = d.Set("upstream_type", upstream.UpstreamType)
	}

	if upstream.Retries != nil {
		_ = d.Set("retries", upstream.Retries)
	}

	if upstream.UpstreamHost != nil {
		_ = d.Set("upstream_host", upstream.UpstreamHost)
	}

	if upstream.Nodes != nil {
		nodesList := []interface{}{}
		for _, nodes := range upstream.Nodes {
			nodesMap := map[string]interface{}{}

			if upstream.Nodes.Host != nil {
				nodesMap["host"] = upstream.Nodes.Host
			}

			if upstream.Nodes.Port != nil {
				nodesMap["port"] = upstream.Nodes.Port
			}

			if upstream.Nodes.Weight != nil {
				nodesMap["weight"] = upstream.Nodes.Weight
			}

			if upstream.Nodes.VmInstanceId != nil {
				nodesMap["vm_instance_id"] = upstream.Nodes.VmInstanceId
			}

			if upstream.Nodes.Tags != nil {
				nodesMap["tags"] = upstream.Nodes.Tags
			}

			if upstream.Nodes.Healthy != nil {
				nodesMap["healthy"] = upstream.Nodes.Healthy
			}

			if upstream.Nodes.ServiceName != nil {
				nodesMap["service_name"] = upstream.Nodes.ServiceName
			}

			if upstream.Nodes.NameSpace != nil {
				nodesMap["name_space"] = upstream.Nodes.NameSpace
			}

			if upstream.Nodes.ClusterId != nil {
				nodesMap["cluster_id"] = upstream.Nodes.ClusterId
			}

			if upstream.Nodes.Source != nil {
				nodesMap["source"] = upstream.Nodes.Source
			}

			if upstream.Nodes.UniqueServiceName != nil {
				nodesMap["unique_service_name"] = upstream.Nodes.UniqueServiceName
			}

			nodesList = append(nodesList, nodesMap)
		}

		_ = d.Set("nodes", nodesList)

	}

	if upstream.HealthChecker != nil {
		healthCheckerMap := map[string]interface{}{}

		if upstream.HealthChecker.EnableActiveCheck != nil {
			healthCheckerMap["enable_active_check"] = upstream.HealthChecker.EnableActiveCheck
		}

		if upstream.HealthChecker.EnablePassiveCheck != nil {
			healthCheckerMap["enable_passive_check"] = upstream.HealthChecker.EnablePassiveCheck
		}

		if upstream.HealthChecker.HealthyHttpStatus != nil {
			healthCheckerMap["healthy_http_status"] = upstream.HealthChecker.HealthyHttpStatus
		}

		if upstream.HealthChecker.UnhealthyHttpStatus != nil {
			healthCheckerMap["unhealthy_http_status"] = upstream.HealthChecker.UnhealthyHttpStatus
		}

		if upstream.HealthChecker.TcpFailureThreshold != nil {
			healthCheckerMap["tcp_failure_threshold"] = upstream.HealthChecker.TcpFailureThreshold
		}

		if upstream.HealthChecker.TimeoutThreshold != nil {
			healthCheckerMap["timeout_threshold"] = upstream.HealthChecker.TimeoutThreshold
		}

		if upstream.HealthChecker.HttpFailureThreshold != nil {
			healthCheckerMap["http_failure_threshold"] = upstream.HealthChecker.HttpFailureThreshold
		}

		if upstream.HealthChecker.ActiveCheckHttpPath != nil {
			healthCheckerMap["active_check_http_path"] = upstream.HealthChecker.ActiveCheckHttpPath
		}

		if upstream.HealthChecker.ActiveCheckTimeout != nil {
			healthCheckerMap["active_check_timeout"] = upstream.HealthChecker.ActiveCheckTimeout
		}

		if upstream.HealthChecker.ActiveCheckInterval != nil {
			healthCheckerMap["active_check_interval"] = upstream.HealthChecker.ActiveCheckInterval
		}

		if upstream.HealthChecker.ActiveRequestHeader != nil {
			healthCheckerMap["active_request_header"] = upstream.HealthChecker.ActiveRequestHeader
		}

		if upstream.HealthChecker.UnhealthyTimeout != nil {
			healthCheckerMap["unhealthy_timeout"] = upstream.HealthChecker.UnhealthyTimeout
		}

		_ = d.Set("health_checker", []interface{}{healthCheckerMap})
	}

	if upstream.K8sService != nil {
		k8sServiceList := []interface{}{}
		for _, k8sService := range upstream.K8sService {
			k8sServiceMap := map[string]interface{}{}

			if upstream.K8sService.Weight != nil {
				k8sServiceMap["weight"] = upstream.K8sService.Weight
			}

			if upstream.K8sService.ClusterId != nil {
				k8sServiceMap["cluster_id"] = upstream.K8sService.ClusterId
			}

			if upstream.K8sService.Namespace != nil {
				k8sServiceMap["namespace"] = upstream.K8sService.Namespace
			}

			if upstream.K8sService.ServiceName != nil {
				k8sServiceMap["service_name"] = upstream.K8sService.ServiceName
			}

			if upstream.K8sService.Port != nil {
				k8sServiceMap["port"] = upstream.K8sService.Port
			}

			if upstream.K8sService.ExtraLabels != nil {
				extraLabelsList := []interface{}{}
				for _, extraLabels := range upstream.K8sService.ExtraLabels {
					extraLabelsMap := map[string]interface{}{}

					if extraLabels.Key != nil {
						extraLabelsMap["key"] = extraLabels.Key
					}

					if extraLabels.Value != nil {
						extraLabelsMap["value"] = extraLabels.Value
					}

					extraLabelsList = append(extraLabelsList, extraLabelsMap)
				}

				k8sServiceMap["extra_labels"] = []interface{}{extraLabelsList}
			}

			if upstream.K8sService.Name != nil {
				k8sServiceMap["name"] = upstream.K8sService.Name
			}

			k8sServiceList = append(k8sServiceList, k8sServiceMap)
		}

		_ = d.Set("k8s_service", k8sServiceList)

	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "apigw", "upstreamId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudApigatewayUpstreamUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_upstream.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := apigateway.NewModifyUpstreamRequest()

	upstreamId := d.Id()

	request.UpstreamId = &upstreamId

	immutableArgs := []string{"scheme", "algorithm", "uniq_vpc_id", "upstream_name", "upstream_description", "upstream_type", "retries", "upstream_host", "nodes", "health_checker", "k8s_service"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

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

	if d.HasChange("upstream_type") {
		if v, ok := d.GetOk("upstream_type"); ok {
			request.UpstreamType = helper.String(v.(string))
		}
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
				upstreamNode := apigateway.UpstreamNode{}
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
				if v, ok := dMap["healthy"]; ok {
					upstreamNode.Healthy = helper.String(v.(string))
				}
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
			if v, ok := dMap["active_request_header"]; ok {
				for _, item := range v.([]interface{}) {
					activeRequestHeaderMap := item.(map[string]interface{})
					upstreamHealthCheckerReqHeaders := apigateway.UpstreamHealthCheckerReqHeaders{}
					upstreamHealthChecker.ActiveRequestHeader = append(upstreamHealthChecker.ActiveRequestHeader, &upstreamHealthCheckerReqHeaders)
				}
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseApigatewayClient().ModifyUpstream(request)
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
		resourceName := BuildTagResourceName("apigw", "upstreamId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudApigatewayUpstreamRead(d, meta)
}

func resourceTencentCloudApigatewayUpstreamDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_upstream.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
	upstreamId := d.Id()

	if err := service.DeleteApigatewayUpstreamById(ctx, upstreamId); err != nil {
		return err
	}

	return nil
}
