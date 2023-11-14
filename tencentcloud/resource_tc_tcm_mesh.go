/*
Provides a resource to create a tcm mesh

Example Usage

```hcl
resource "tencentcloud_tcm_mesh" "mesh" {
  mesh_id = "mesh-xxxxxxxx"
  display_name = "test mesh"
  mesh_version = "1.8.1"
  type = "HOSTED"
  config {
		tracing {
			enable = true
			a_p_m {
				enable = true
				region = "ap-shanghai"
				instance_id = "apm-xxx"
			}
			sampling =
			zipkin {
				address = "10.10.10.10:9411"
			}
		}
		prometheus {
			vpc_id = "vpc-xxx"
			subnet_id = "subnet-xxx"
			region = "sh"
			instance_id = "prom-xxx"
			custom_prom {
				is_public_addr = false
				vpc_id = "vpc-xxx"
				url = "http://x.x.x.x:9090"
				auth_type = "none, basic"
				username = "test"
				password = "test"
			}
		}
		istio {
			outbound_traffic_policy = "ALLOW_ANY"
			disable_policy_checks = true
			enable_pilot_h_t_t_p = true
			disable_h_t_t_p_retry = true
			smart_d_n_s {
				istio_meta_d_n_s_capture = true
				istio_meta_d_n_s_auto_allocate = true
			}
		}
		inject {
			exclude_i_p_ranges =
			hold_application_until_proxy_starts = true
			hold_proxy_until_application_ends = true
		}
		sidecar_resources {
			limits {
				name = "cpu"
				quantity = "100m"
			}
			requests {
				name = "cpu"
				quantity = "100m"
			}
		}
		access_log {
			enable = true
			template = "istio"
			selected_range {
				items {
					namespace = "test"
					gateways =
				}
				all = true
			}
			c_l_s {
				enable = true
				log_set = "f832fd4a-2b57-4573-ab6c-c3c69caf84c9"
				topic = "1ad05336-8afc-4e56-91e5-28d8a8511761"
			}
			encoding = "JSON"
			format = ""
			address = "10.10.10.4:3398"
			enable_server = true
			enable_stdout = true
		}

  }
  tag_list {
		key = "key"
		value = "value"
		passthrough = true

  }
}
```

Import

tcm mesh can be imported using the id, e.g.

```
terraform import tencentcloud_tcm_mesh.mesh mesh_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcm/v20210413"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTcmMesh() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcmMeshCreate,
		Read:   resourceTencentCloudTcmMeshRead,
		Update: resourceTencentCloudTcmMeshUpdate,
		Delete: resourceTencentCloudTcmMeshDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mesh_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Mesh ID.",
			},

			"display_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Mesh name.",
			},

			"mesh_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Mesh version.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Mesh type.",
			},

			"config": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Mesh configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tracing": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Tracing config.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether enable tracing.",
									},
									"a_p_m": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "APM config.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether enable APM.",
												},
												"region": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Region.",
												},
												"instance_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Instance id of the APM.",
												},
											},
										},
									},
									"sampling": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Tracing sampling, 0.0-1.0.",
									},
									"zipkin": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Third party zipkin config.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Zipkin address.",
												},
											},
										},
									},
								},
							},
						},
						"prometheus": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Prometheus configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Vpc id.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Subnet id.",
									},
									"region": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Region.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Instance id.",
									},
									"custom_prom": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Custom prometheus.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_public_addr": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether it is public address, default false.",
												},
												"vpc_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Vpc id.",
												},
												"url": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Url of the prometheus.",
												},
												"auth_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Authentication type of the prometheus.",
												},
												"username": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Username of the prometheus, used in basic authentication type.",
												},
												"password": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Password of the prometheus, used in basic authentication type.",
												},
											},
										},
									},
								},
							},
						},
						"istio": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Istio configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"outbound_traffic_policy": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Outbound traffic policy, REGISTRY_ONLY or ALLOW_ANY, see https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/#MeshConfig-OutboundTrafficPolicy-Mode.",
									},
									"disable_policy_checks": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Disable policy checks.",
									},
									"enable_pilot_h_t_t_p": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Enable HTTP/1.0 support.",
									},
									"disable_h_t_t_p_retry": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Disable http retry.",
									},
									"smart_d_n_s": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "SmartDNS configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"istio_meta_d_n_s_capture": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Enable dns proxy.",
												},
												"istio_meta_d_n_s_auto_allocate": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Enable auto allocate address.",
												},
											},
										},
									},
								},
							},
						},
						"inject": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Sidecar inject configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"exclude_i_p_ranges": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "IP ranges that should not be proxied.",
									},
									"hold_application_until_proxy_starts": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Let istio-proxy(sidecar) start first, before app container.",
									},
									"hold_proxy_until_application_ends": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Let istio-proxy(sidecar) stop last, after app container.",
									},
								},
							},
						},
						"sidecar_resources": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Default sidecar requests and limits.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"limits": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Sidecar limits.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Resource type name.",
												},
												"quantity": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Resource quantity.",
												},
											},
										},
									},
									"requests": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Sidecar requests.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Resource type name.",
												},
												"quantity": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Resource quantity.",
												},
											},
										},
									},
								},
							},
						},
						"access_log": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Access log configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Enable access log.",
									},
									"template": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Access log template, support istio (default) and trace (more detailed).",
									},
									"selected_range": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Only enable access log for selected range.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"items": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Enable access log for selected items.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"namespace": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Select namespace name.",
															},
															"gateways": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Optional:    true,
																Description: "Select gateway names.",
															},
														},
													},
												},
												"all": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Select all.",
												},
											},
										},
									},
									"c_l_s": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Post log to CLS.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Enable post log to CLS.",
												},
												"log_set": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "CLS log set id.",
												},
												"topic": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "CLS log topic id.",
												},
											},
										},
									},
									"encoding": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Access log encoding, support TEXT and JSON.",
									},
									"format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Customized access log format, see https://istio.io/latest/docs/tasks/observability/logs/access-log/.",
									},
									"address": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Third-party grpc log server address.",
									},
									"enable_server": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Enable third-party grpc log server.",
									},
									"enable_stdout": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Enable stdout.",
									},
								},
							},
						},
					},
				},
			},

			"tag_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "A list of associated tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
						"passthrough": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Passthrough to other related product.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTcmMeshCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_mesh.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tcm.NewCreateMeshRequest()
		response = tcm.NewCreateMeshResponse()
		meshId   string
	)
	if v, ok := d.GetOk("mesh_id"); ok {
		meshId = v.(string)
		request.MeshId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("display_name"); ok {
		request.DisplayName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("mesh_version"); ok {
		request.MeshVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "config"); ok {
		meshConfig := tcm.MeshConfig{}
		if tracingMap, ok := helper.InterfaceToMap(dMap, "tracing"); ok {
			tracingConfig := tcm.TracingConfig{}
			if v, ok := tracingMap["enable"]; ok {
				tracingConfig.Enable = helper.Bool(v.(bool))
			}
			if aPMMap, ok := helper.InterfaceToMap(tracingMap, "a_p_m"); ok {
				aPM := tcm.APM{}
				if v, ok := aPMMap["enable"]; ok {
					aPM.Enable = helper.Bool(v.(bool))
				}
				if v, ok := aPMMap["region"]; ok {
					aPM.Region = helper.String(v.(string))
				}
				if v, ok := aPMMap["instance_id"]; ok {
					aPM.InstanceId = helper.String(v.(string))
				}
				tracingConfig.APM = &aPM
			}
			if v, ok := tracingMap["sampling"]; ok {
				tracingConfig.Sampling = helper.Float64(v.(float64))
			}
			if zipkinMap, ok := helper.InterfaceToMap(tracingMap, "zipkin"); ok {
				tracingZipkin := tcm.TracingZipkin{}
				if v, ok := zipkinMap["address"]; ok {
					tracingZipkin.Address = helper.String(v.(string))
				}
				tracingConfig.Zipkin = &tracingZipkin
			}
			meshConfig.Tracing = &tracingConfig
		}
		if prometheusMap, ok := helper.InterfaceToMap(dMap, "prometheus"); ok {
			prometheusConfig := tcm.PrometheusConfig{}
			if v, ok := prometheusMap["vpc_id"]; ok {
				prometheusConfig.VpcId = helper.String(v.(string))
			}
			if v, ok := prometheusMap["subnet_id"]; ok {
				prometheusConfig.SubnetId = helper.String(v.(string))
			}
			if v, ok := prometheusMap["region"]; ok {
				prometheusConfig.Region = helper.String(v.(string))
			}
			if v, ok := prometheusMap["instance_id"]; ok {
				prometheusConfig.InstanceId = helper.String(v.(string))
			}
			if customPromMap, ok := helper.InterfaceToMap(prometheusMap, "custom_prom"); ok {
				customPromConfig := tcm.CustomPromConfig{}
				if v, ok := customPromMap["is_public_addr"]; ok {
					customPromConfig.IsPublicAddr = helper.Bool(v.(bool))
				}
				if v, ok := customPromMap["vpc_id"]; ok {
					customPromConfig.VpcId = helper.String(v.(string))
				}
				if v, ok := customPromMap["url"]; ok {
					customPromConfig.Url = helper.String(v.(string))
				}
				if v, ok := customPromMap["auth_type"]; ok {
					customPromConfig.AuthType = helper.String(v.(string))
				}
				if v, ok := customPromMap["username"]; ok {
					customPromConfig.Username = helper.String(v.(string))
				}
				if v, ok := customPromMap["password"]; ok {
					customPromConfig.Password = helper.String(v.(string))
				}
				prometheusConfig.CustomProm = &customPromConfig
			}
			meshConfig.Prometheus = &prometheusConfig
		}
		if istioMap, ok := helper.InterfaceToMap(dMap, "istio"); ok {
			istioConfig := tcm.IstioConfig{}
			if v, ok := istioMap["outbound_traffic_policy"]; ok {
				istioConfig.OutboundTrafficPolicy = helper.String(v.(string))
			}
			if v, ok := istioMap["disable_policy_checks"]; ok {
				istioConfig.DisablePolicyChecks = helper.Bool(v.(bool))
			}
			if v, ok := istioMap["enable_pilot_h_t_t_p"]; ok {
				istioConfig.EnablePilotHTTP = helper.Bool(v.(bool))
			}
			if v, ok := istioMap["disable_h_t_t_p_retry"]; ok {
				istioConfig.DisableHTTPRetry = helper.Bool(v.(bool))
			}
			if smartDNSMap, ok := helper.InterfaceToMap(istioMap, "smart_d_n_s"); ok {
				smartDNSConfig := tcm.SmartDNSConfig{}
				if v, ok := smartDNSMap["istio_meta_d_n_s_capture"]; ok {
					smartDNSConfig.IstioMetaDNSCapture = helper.Bool(v.(bool))
				}
				if v, ok := smartDNSMap["istio_meta_d_n_s_auto_allocate"]; ok {
					smartDNSConfig.IstioMetaDNSAutoAllocate = helper.Bool(v.(bool))
				}
				istioConfig.SmartDNS = &smartDNSConfig
			}
			meshConfig.Istio = &istioConfig
		}
		if injectMap, ok := helper.InterfaceToMap(dMap, "inject"); ok {
			injectConfig := tcm.InjectConfig{}
			if v, ok := injectMap["exclude_i_p_ranges"]; ok {
				excludeIPRangesSet := v.(*schema.Set).List()
				for i := range excludeIPRangesSet {
					excludeIPRanges := excludeIPRangesSet[i].(string)
					injectConfig.ExcludeIPRanges = append(injectConfig.ExcludeIPRanges, &excludeIPRanges)
				}
			}
			if v, ok := injectMap["hold_application_until_proxy_starts"]; ok {
				injectConfig.HoldApplicationUntilProxyStarts = helper.Bool(v.(bool))
			}
			if v, ok := injectMap["hold_proxy_until_application_ends"]; ok {
				injectConfig.HoldProxyUntilApplicationEnds = helper.Bool(v.(bool))
			}
			meshConfig.Inject = &injectConfig
		}
		if sidecarResourcesMap, ok := helper.InterfaceToMap(dMap, "sidecar_resources"); ok {
			resourceRequirements := tcm.ResourceRequirements{}
			if v, ok := sidecarResourcesMap["limits"]; ok {
				for _, item := range v.([]interface{}) {
					limitsMap := item.(map[string]interface{})
					resource := tcm.Resource{}
					if v, ok := limitsMap["name"]; ok {
						resource.Name = helper.String(v.(string))
					}
					if v, ok := limitsMap["quantity"]; ok {
						resource.Quantity = helper.String(v.(string))
					}
					resourceRequirements.Limits = append(resourceRequirements.Limits, &resource)
				}
			}
			if v, ok := sidecarResourcesMap["requests"]; ok {
				for _, item := range v.([]interface{}) {
					requestsMap := item.(map[string]interface{})
					resource := tcm.Resource{}
					if v, ok := requestsMap["name"]; ok {
						resource.Name = helper.String(v.(string))
					}
					if v, ok := requestsMap["quantity"]; ok {
						resource.Quantity = helper.String(v.(string))
					}
					resourceRequirements.Requests = append(resourceRequirements.Requests, &resource)
				}
			}
			meshConfig.SidecarResources = &resourceRequirements
		}
		if accessLogMap, ok := helper.InterfaceToMap(dMap, "access_log"); ok {
			accessLogConfig := tcm.AccessLogConfig{}
			if v, ok := accessLogMap["enable"]; ok {
				accessLogConfig.Enable = helper.Bool(v.(bool))
			}
			if v, ok := accessLogMap["template"]; ok {
				accessLogConfig.Template = helper.String(v.(string))
			}
			if selectedRangeMap, ok := helper.InterfaceToMap(accessLogMap, "selected_range"); ok {
				selectedRange := tcm.SelectedRange{}
				if v, ok := selectedRangeMap["items"]; ok {
					for _, item := range v.([]interface{}) {
						itemsMap := item.(map[string]interface{})
						selectedItems := tcm.SelectedItems{}
						if v, ok := itemsMap["namespace"]; ok {
							selectedItems.Namespace = helper.String(v.(string))
						}
						if v, ok := itemsMap["gateways"]; ok {
							gatewaysSet := v.(*schema.Set).List()
							for i := range gatewaysSet {
								gateways := gatewaysSet[i].(string)
								selectedItems.Gateways = append(selectedItems.Gateways, &gateways)
							}
						}
						selectedRange.Items = append(selectedRange.Items, &selectedItems)
					}
				}
				if v, ok := selectedRangeMap["all"]; ok {
					selectedRange.All = helper.Bool(v.(bool))
				}
				accessLogConfig.SelectedRange = &selectedRange
			}
			if cLSMap, ok := helper.InterfaceToMap(accessLogMap, "c_l_s"); ok {
				cLS := tcm.CLS{}
				if v, ok := cLSMap["enable"]; ok {
					cLS.Enable = helper.Bool(v.(bool))
				}
				if v, ok := cLSMap["log_set"]; ok {
					cLS.LogSet = helper.String(v.(string))
				}
				if v, ok := cLSMap["topic"]; ok {
					cLS.Topic = helper.String(v.(string))
				}
				accessLogConfig.CLS = &cLS
			}
			if v, ok := accessLogMap["encoding"]; ok {
				accessLogConfig.Encoding = helper.String(v.(string))
			}
			if v, ok := accessLogMap["format"]; ok {
				accessLogConfig.Format = helper.String(v.(string))
			}
			if v, ok := accessLogMap["address"]; ok {
				accessLogConfig.Address = helper.String(v.(string))
			}
			if v, ok := accessLogMap["enable_server"]; ok {
				accessLogConfig.EnableServer = helper.Bool(v.(bool))
			}
			if v, ok := accessLogMap["enable_stdout"]; ok {
				accessLogConfig.EnableStdout = helper.Bool(v.(bool))
			}
			meshConfig.AccessLog = &accessLogConfig
		}
		request.Config = &meshConfig
	}

	if v, ok := d.GetOk("tag_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			tag := tcm.Tag{}
			if v, ok := dMap["key"]; ok {
				tag.Key = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				tag.Value = helper.String(v.(string))
			}
			if v, ok := dMap["passthrough"]; ok {
				tag.Passthrough = helper.Bool(v.(bool))
			}
			request.TagList = append(request.TagList, &tag)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcmClient().CreateMesh(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tcm mesh failed, reason:%+v", logId, err)
		return err
	}

	meshId = *response.Response.MeshId
	d.SetId(meshId)

	return resourceTencentCloudTcmMeshRead(d, meta)
}

func resourceTencentCloudTcmMeshRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_mesh.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcmService{client: meta.(*TencentCloudClient).apiV3Conn}

	meshId := d.Id()

	mesh, err := service.DescribeTcmMeshById(ctx, meshId)
	if err != nil {
		return err
	}

	if mesh == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TcmMesh` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if mesh.MeshId != nil {
		_ = d.Set("mesh_id", mesh.MeshId)
	}

	if mesh.DisplayName != nil {
		_ = d.Set("display_name", mesh.DisplayName)
	}

	if mesh.MeshVersion != nil {
		_ = d.Set("mesh_version", mesh.MeshVersion)
	}

	if mesh.Type != nil {
		_ = d.Set("type", mesh.Type)
	}

	if mesh.Config != nil {
		configMap := map[string]interface{}{}

		if mesh.Config.Tracing != nil {
			tracingMap := map[string]interface{}{}

			if mesh.Config.Tracing.Enable != nil {
				tracingMap["enable"] = mesh.Config.Tracing.Enable
			}

			if mesh.Config.Tracing.APM != nil {
				aPMMap := map[string]interface{}{}

				if mesh.Config.Tracing.APM.Enable != nil {
					aPMMap["enable"] = mesh.Config.Tracing.APM.Enable
				}

				if mesh.Config.Tracing.APM.Region != nil {
					aPMMap["region"] = mesh.Config.Tracing.APM.Region
				}

				if mesh.Config.Tracing.APM.InstanceId != nil {
					aPMMap["instance_id"] = mesh.Config.Tracing.APM.InstanceId
				}

				tracingMap["a_p_m"] = []interface{}{aPMMap}
			}

			if mesh.Config.Tracing.Sampling != nil {
				tracingMap["sampling"] = mesh.Config.Tracing.Sampling
			}

			if mesh.Config.Tracing.Zipkin != nil {
				zipkinMap := map[string]interface{}{}

				if mesh.Config.Tracing.Zipkin.Address != nil {
					zipkinMap["address"] = mesh.Config.Tracing.Zipkin.Address
				}

				tracingMap["zipkin"] = []interface{}{zipkinMap}
			}

			configMap["tracing"] = []interface{}{tracingMap}
		}

		if mesh.Config.Prometheus != nil {
			prometheusMap := map[string]interface{}{}

			if mesh.Config.Prometheus.VpcId != nil {
				prometheusMap["vpc_id"] = mesh.Config.Prometheus.VpcId
			}

			if mesh.Config.Prometheus.SubnetId != nil {
				prometheusMap["subnet_id"] = mesh.Config.Prometheus.SubnetId
			}

			if mesh.Config.Prometheus.Region != nil {
				prometheusMap["region"] = mesh.Config.Prometheus.Region
			}

			if mesh.Config.Prometheus.InstanceId != nil {
				prometheusMap["instance_id"] = mesh.Config.Prometheus.InstanceId
			}

			if mesh.Config.Prometheus.CustomProm != nil {
				customPromMap := map[string]interface{}{}

				if mesh.Config.Prometheus.CustomProm.IsPublicAddr != nil {
					customPromMap["is_public_addr"] = mesh.Config.Prometheus.CustomProm.IsPublicAddr
				}

				if mesh.Config.Prometheus.CustomProm.VpcId != nil {
					customPromMap["vpc_id"] = mesh.Config.Prometheus.CustomProm.VpcId
				}

				if mesh.Config.Prometheus.CustomProm.Url != nil {
					customPromMap["url"] = mesh.Config.Prometheus.CustomProm.Url
				}

				if mesh.Config.Prometheus.CustomProm.AuthType != nil {
					customPromMap["auth_type"] = mesh.Config.Prometheus.CustomProm.AuthType
				}

				if mesh.Config.Prometheus.CustomProm.Username != nil {
					customPromMap["username"] = mesh.Config.Prometheus.CustomProm.Username
				}

				if mesh.Config.Prometheus.CustomProm.Password != nil {
					customPromMap["password"] = mesh.Config.Prometheus.CustomProm.Password
				}

				prometheusMap["custom_prom"] = []interface{}{customPromMap}
			}

			configMap["prometheus"] = []interface{}{prometheusMap}
		}

		if mesh.Config.Istio != nil {
			istioMap := map[string]interface{}{}

			if mesh.Config.Istio.OutboundTrafficPolicy != nil {
				istioMap["outbound_traffic_policy"] = mesh.Config.Istio.OutboundTrafficPolicy
			}

			if mesh.Config.Istio.DisablePolicyChecks != nil {
				istioMap["disable_policy_checks"] = mesh.Config.Istio.DisablePolicyChecks
			}

			if mesh.Config.Istio.EnablePilotHTTP != nil {
				istioMap["enable_pilot_h_t_t_p"] = mesh.Config.Istio.EnablePilotHTTP
			}

			if mesh.Config.Istio.DisableHTTPRetry != nil {
				istioMap["disable_h_t_t_p_retry"] = mesh.Config.Istio.DisableHTTPRetry
			}

			if mesh.Config.Istio.SmartDNS != nil {
				smartDNSMap := map[string]interface{}{}

				if mesh.Config.Istio.SmartDNS.IstioMetaDNSCapture != nil {
					smartDNSMap["istio_meta_d_n_s_capture"] = mesh.Config.Istio.SmartDNS.IstioMetaDNSCapture
				}

				if mesh.Config.Istio.SmartDNS.IstioMetaDNSAutoAllocate != nil {
					smartDNSMap["istio_meta_d_n_s_auto_allocate"] = mesh.Config.Istio.SmartDNS.IstioMetaDNSAutoAllocate
				}

				istioMap["smart_d_n_s"] = []interface{}{smartDNSMap}
			}

			configMap["istio"] = []interface{}{istioMap}
		}

		if mesh.Config.Inject != nil {
			injectMap := map[string]interface{}{}

			if mesh.Config.Inject.ExcludeIPRanges != nil {
				injectMap["exclude_i_p_ranges"] = mesh.Config.Inject.ExcludeIPRanges
			}

			if mesh.Config.Inject.HoldApplicationUntilProxyStarts != nil {
				injectMap["hold_application_until_proxy_starts"] = mesh.Config.Inject.HoldApplicationUntilProxyStarts
			}

			if mesh.Config.Inject.HoldProxyUntilApplicationEnds != nil {
				injectMap["hold_proxy_until_application_ends"] = mesh.Config.Inject.HoldProxyUntilApplicationEnds
			}

			configMap["inject"] = []interface{}{injectMap}
		}

		if mesh.Config.SidecarResources != nil {
			sidecarResourcesMap := map[string]interface{}{}

			if mesh.Config.SidecarResources.Limits != nil {
				limitsList := []interface{}{}
				for _, limits := range mesh.Config.SidecarResources.Limits {
					limitsMap := map[string]interface{}{}

					if limits.Name != nil {
						limitsMap["name"] = limits.Name
					}

					if limits.Quantity != nil {
						limitsMap["quantity"] = limits.Quantity
					}

					limitsList = append(limitsList, limitsMap)
				}

				sidecarResourcesMap["limits"] = []interface{}{limitsList}
			}

			if mesh.Config.SidecarResources.Requests != nil {
				requestsList := []interface{}{}
				for _, requests := range mesh.Config.SidecarResources.Requests {
					requestsMap := map[string]interface{}{}

					if requests.Name != nil {
						requestsMap["name"] = requests.Name
					}

					if requests.Quantity != nil {
						requestsMap["quantity"] = requests.Quantity
					}

					requestsList = append(requestsList, requestsMap)
				}

				sidecarResourcesMap["requests"] = []interface{}{requestsList}
			}

			configMap["sidecar_resources"] = []interface{}{sidecarResourcesMap}
		}

		if mesh.Config.AccessLog != nil {
			accessLogMap := map[string]interface{}{}

			if mesh.Config.AccessLog.Enable != nil {
				accessLogMap["enable"] = mesh.Config.AccessLog.Enable
			}

			if mesh.Config.AccessLog.Template != nil {
				accessLogMap["template"] = mesh.Config.AccessLog.Template
			}

			if mesh.Config.AccessLog.SelectedRange != nil {
				selectedRangeMap := map[string]interface{}{}

				if mesh.Config.AccessLog.SelectedRange.Items != nil {
					itemsList := []interface{}{}
					for _, items := range mesh.Config.AccessLog.SelectedRange.Items {
						itemsMap := map[string]interface{}{}

						if items.Namespace != nil {
							itemsMap["namespace"] = items.Namespace
						}

						if items.Gateways != nil {
							itemsMap["gateways"] = items.Gateways
						}

						itemsList = append(itemsList, itemsMap)
					}

					selectedRangeMap["items"] = []interface{}{itemsList}
				}

				if mesh.Config.AccessLog.SelectedRange.All != nil {
					selectedRangeMap["all"] = mesh.Config.AccessLog.SelectedRange.All
				}

				accessLogMap["selected_range"] = []interface{}{selectedRangeMap}
			}

			if mesh.Config.AccessLog.CLS != nil {
				cLSMap := map[string]interface{}{}

				if mesh.Config.AccessLog.CLS.Enable != nil {
					cLSMap["enable"] = mesh.Config.AccessLog.CLS.Enable
				}

				if mesh.Config.AccessLog.CLS.LogSet != nil {
					cLSMap["log_set"] = mesh.Config.AccessLog.CLS.LogSet
				}

				if mesh.Config.AccessLog.CLS.Topic != nil {
					cLSMap["topic"] = mesh.Config.AccessLog.CLS.Topic
				}

				accessLogMap["c_l_s"] = []interface{}{cLSMap}
			}

			if mesh.Config.AccessLog.Encoding != nil {
				accessLogMap["encoding"] = mesh.Config.AccessLog.Encoding
			}

			if mesh.Config.AccessLog.Format != nil {
				accessLogMap["format"] = mesh.Config.AccessLog.Format
			}

			if mesh.Config.AccessLog.Address != nil {
				accessLogMap["address"] = mesh.Config.AccessLog.Address
			}

			if mesh.Config.AccessLog.EnableServer != nil {
				accessLogMap["enable_server"] = mesh.Config.AccessLog.EnableServer
			}

			if mesh.Config.AccessLog.EnableStdout != nil {
				accessLogMap["enable_stdout"] = mesh.Config.AccessLog.EnableStdout
			}

			configMap["access_log"] = []interface{}{accessLogMap}
		}

		_ = d.Set("config", []interface{}{configMap})
	}

	if mesh.TagList != nil {
		tagListList := []interface{}{}
		for _, tagList := range mesh.TagList {
			tagListMap := map[string]interface{}{}

			if mesh.TagList.Key != nil {
				tagListMap["key"] = mesh.TagList.Key
			}

			if mesh.TagList.Value != nil {
				tagListMap["value"] = mesh.TagList.Value
			}

			if mesh.TagList.Passthrough != nil {
				tagListMap["passthrough"] = mesh.TagList.Passthrough
			}

			tagListList = append(tagListList, tagListMap)
		}

		_ = d.Set("tag_list", tagListList)

	}

	return nil
}

func resourceTencentCloudTcmMeshUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_mesh.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tcm.NewModifyMeshRequest()

	meshId := d.Id()

	request.MeshId = &meshId

	immutableArgs := []string{"mesh_id", "display_name", "mesh_version", "type", "config", "tag_list"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("display_name") {
		if v, ok := d.GetOk("display_name"); ok {
			request.DisplayName = helper.String(v.(string))
		}
	}

	if d.HasChange("config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "config"); ok {
			meshConfig := tcm.MeshConfig{}
			if tracingMap, ok := helper.InterfaceToMap(dMap, "tracing"); ok {
				tracingConfig := tcm.TracingConfig{}
				if v, ok := tracingMap["enable"]; ok {
					tracingConfig.Enable = helper.Bool(v.(bool))
				}
				if aPMMap, ok := helper.InterfaceToMap(tracingMap, "a_p_m"); ok {
					aPM := tcm.APM{}
					if v, ok := aPMMap["enable"]; ok {
						aPM.Enable = helper.Bool(v.(bool))
					}
					if v, ok := aPMMap["region"]; ok {
						aPM.Region = helper.String(v.(string))
					}
					if v, ok := aPMMap["instance_id"]; ok {
						aPM.InstanceId = helper.String(v.(string))
					}
					tracingConfig.APM = &aPM
				}
				if v, ok := tracingMap["sampling"]; ok {
					tracingConfig.Sampling = helper.Float64(v.(float64))
				}
				if zipkinMap, ok := helper.InterfaceToMap(tracingMap, "zipkin"); ok {
					tracingZipkin := tcm.TracingZipkin{}
					if v, ok := zipkinMap["address"]; ok {
						tracingZipkin.Address = helper.String(v.(string))
					}
					tracingConfig.Zipkin = &tracingZipkin
				}
				meshConfig.Tracing = &tracingConfig
			}
			if prometheusMap, ok := helper.InterfaceToMap(dMap, "prometheus"); ok {
				prometheusConfig := tcm.PrometheusConfig{}
				if v, ok := prometheusMap["vpc_id"]; ok {
					prometheusConfig.VpcId = helper.String(v.(string))
				}
				if v, ok := prometheusMap["subnet_id"]; ok {
					prometheusConfig.SubnetId = helper.String(v.(string))
				}
				if v, ok := prometheusMap["region"]; ok {
					prometheusConfig.Region = helper.String(v.(string))
				}
				if v, ok := prometheusMap["instance_id"]; ok {
					prometheusConfig.InstanceId = helper.String(v.(string))
				}
				if customPromMap, ok := helper.InterfaceToMap(prometheusMap, "custom_prom"); ok {
					customPromConfig := tcm.CustomPromConfig{}
					if v, ok := customPromMap["is_public_addr"]; ok {
						customPromConfig.IsPublicAddr = helper.Bool(v.(bool))
					}
					if v, ok := customPromMap["vpc_id"]; ok {
						customPromConfig.VpcId = helper.String(v.(string))
					}
					if v, ok := customPromMap["url"]; ok {
						customPromConfig.Url = helper.String(v.(string))
					}
					if v, ok := customPromMap["auth_type"]; ok {
						customPromConfig.AuthType = helper.String(v.(string))
					}
					if v, ok := customPromMap["username"]; ok {
						customPromConfig.Username = helper.String(v.(string))
					}
					if v, ok := customPromMap["password"]; ok {
						customPromConfig.Password = helper.String(v.(string))
					}
					prometheusConfig.CustomProm = &customPromConfig
				}
				meshConfig.Prometheus = &prometheusConfig
			}
			if istioMap, ok := helper.InterfaceToMap(dMap, "istio"); ok {
				istioConfig := tcm.IstioConfig{}
				if v, ok := istioMap["outbound_traffic_policy"]; ok {
					istioConfig.OutboundTrafficPolicy = helper.String(v.(string))
				}
				if v, ok := istioMap["disable_policy_checks"]; ok {
					istioConfig.DisablePolicyChecks = helper.Bool(v.(bool))
				}
				if v, ok := istioMap["enable_pilot_h_t_t_p"]; ok {
					istioConfig.EnablePilotHTTP = helper.Bool(v.(bool))
				}
				if v, ok := istioMap["disable_h_t_t_p_retry"]; ok {
					istioConfig.DisableHTTPRetry = helper.Bool(v.(bool))
				}
				if smartDNSMap, ok := helper.InterfaceToMap(istioMap, "smart_d_n_s"); ok {
					smartDNSConfig := tcm.SmartDNSConfig{}
					if v, ok := smartDNSMap["istio_meta_d_n_s_capture"]; ok {
						smartDNSConfig.IstioMetaDNSCapture = helper.Bool(v.(bool))
					}
					if v, ok := smartDNSMap["istio_meta_d_n_s_auto_allocate"]; ok {
						smartDNSConfig.IstioMetaDNSAutoAllocate = helper.Bool(v.(bool))
					}
					istioConfig.SmartDNS = &smartDNSConfig
				}
				meshConfig.Istio = &istioConfig
			}
			if injectMap, ok := helper.InterfaceToMap(dMap, "inject"); ok {
				injectConfig := tcm.InjectConfig{}
				if v, ok := injectMap["exclude_i_p_ranges"]; ok {
					excludeIPRangesSet := v.(*schema.Set).List()
					for i := range excludeIPRangesSet {
						excludeIPRanges := excludeIPRangesSet[i].(string)
						injectConfig.ExcludeIPRanges = append(injectConfig.ExcludeIPRanges, &excludeIPRanges)
					}
				}
				if v, ok := injectMap["hold_application_until_proxy_starts"]; ok {
					injectConfig.HoldApplicationUntilProxyStarts = helper.Bool(v.(bool))
				}
				if v, ok := injectMap["hold_proxy_until_application_ends"]; ok {
					injectConfig.HoldProxyUntilApplicationEnds = helper.Bool(v.(bool))
				}
				meshConfig.Inject = &injectConfig
			}
			if sidecarResourcesMap, ok := helper.InterfaceToMap(dMap, "sidecar_resources"); ok {
				resourceRequirements := tcm.ResourceRequirements{}
				if v, ok := sidecarResourcesMap["limits"]; ok {
					for _, item := range v.([]interface{}) {
						limitsMap := item.(map[string]interface{})
						resource := tcm.Resource{}
						if v, ok := limitsMap["name"]; ok {
							resource.Name = helper.String(v.(string))
						}
						if v, ok := limitsMap["quantity"]; ok {
							resource.Quantity = helper.String(v.(string))
						}
						resourceRequirements.Limits = append(resourceRequirements.Limits, &resource)
					}
				}
				if v, ok := sidecarResourcesMap["requests"]; ok {
					for _, item := range v.([]interface{}) {
						requestsMap := item.(map[string]interface{})
						resource := tcm.Resource{}
						if v, ok := requestsMap["name"]; ok {
							resource.Name = helper.String(v.(string))
						}
						if v, ok := requestsMap["quantity"]; ok {
							resource.Quantity = helper.String(v.(string))
						}
						resourceRequirements.Requests = append(resourceRequirements.Requests, &resource)
					}
				}
				meshConfig.SidecarResources = &resourceRequirements
			}
			if accessLogMap, ok := helper.InterfaceToMap(dMap, "access_log"); ok {
				accessLogConfig := tcm.AccessLogConfig{}
				if v, ok := accessLogMap["enable"]; ok {
					accessLogConfig.Enable = helper.Bool(v.(bool))
				}
				if v, ok := accessLogMap["template"]; ok {
					accessLogConfig.Template = helper.String(v.(string))
				}
				if selectedRangeMap, ok := helper.InterfaceToMap(accessLogMap, "selected_range"); ok {
					selectedRange := tcm.SelectedRange{}
					if v, ok := selectedRangeMap["items"]; ok {
						for _, item := range v.([]interface{}) {
							itemsMap := item.(map[string]interface{})
							selectedItems := tcm.SelectedItems{}
							if v, ok := itemsMap["namespace"]; ok {
								selectedItems.Namespace = helper.String(v.(string))
							}
							if v, ok := itemsMap["gateways"]; ok {
								gatewaysSet := v.(*schema.Set).List()
								for i := range gatewaysSet {
									gateways := gatewaysSet[i].(string)
									selectedItems.Gateways = append(selectedItems.Gateways, &gateways)
								}
							}
							selectedRange.Items = append(selectedRange.Items, &selectedItems)
						}
					}
					if v, ok := selectedRangeMap["all"]; ok {
						selectedRange.All = helper.Bool(v.(bool))
					}
					accessLogConfig.SelectedRange = &selectedRange
				}
				if cLSMap, ok := helper.InterfaceToMap(accessLogMap, "c_l_s"); ok {
					cLS := tcm.CLS{}
					if v, ok := cLSMap["enable"]; ok {
						cLS.Enable = helper.Bool(v.(bool))
					}
					if v, ok := cLSMap["log_set"]; ok {
						cLS.LogSet = helper.String(v.(string))
					}
					if v, ok := cLSMap["topic"]; ok {
						cLS.Topic = helper.String(v.(string))
					}
					accessLogConfig.CLS = &cLS
				}
				if v, ok := accessLogMap["encoding"]; ok {
					accessLogConfig.Encoding = helper.String(v.(string))
				}
				if v, ok := accessLogMap["format"]; ok {
					accessLogConfig.Format = helper.String(v.(string))
				}
				if v, ok := accessLogMap["address"]; ok {
					accessLogConfig.Address = helper.String(v.(string))
				}
				if v, ok := accessLogMap["enable_server"]; ok {
					accessLogConfig.EnableServer = helper.Bool(v.(bool))
				}
				if v, ok := accessLogMap["enable_stdout"]; ok {
					accessLogConfig.EnableStdout = helper.Bool(v.(bool))
				}
				meshConfig.AccessLog = &accessLogConfig
			}
			request.Config = &meshConfig
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcmClient().ModifyMesh(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tcm mesh failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTcmMeshRead(d, meta)
}

func resourceTencentCloudTcmMeshDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_mesh.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcmService{client: meta.(*TencentCloudClient).apiV3Conn}
	meshId := d.Id()

	if err := service.DeleteTcmMeshById(ctx, meshId); err != nil {
		return err
	}

	return nil
}
