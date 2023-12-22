package tcm

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcm/v20210413"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTcmMesh() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTcmMeshRead,
		Create: resourceTencentCloudTcmMeshCreate,
		Update: resourceTencentCloudTcmMeshUpdate,
		Delete: resourceTencentCloudTcmMeshDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mesh_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mesh ID.",
			},

			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Mesh name.",
			},

			"mesh_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Mesh version.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Mesh type.",
			},

			"config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Mesh configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tracing": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Tracing config.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Whether enable tracing.",
									},
									"apm": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Description: "APM config.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable": {
													Type:        schema.TypeBool,
													Optional:    true,
													Computed:    true,
													Description: "Whether enable APM.",
												},
												"region": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Region.",
												},
												"instance_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Instance id of the APM.",
												},
											},
										},
									},
									"sampling": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "Tracing sampling, 0.0-1.0.",
									},
									"zipkin": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
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
							Computed:    true,
							Description: "Prometheus configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Vpc id.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Subnet id.",
									},
									"region": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Region.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Instance id.",
									},
									"custom_prom": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Description: "Custom prometheus.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_public_addr": {
													Type:        schema.TypeBool,
													Optional:    true,
													Computed:    true,
													Description: "Whether it is public address, default false.",
												},
												"vpc_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
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
													Computed:    true,
													Description: "Username of the prometheus, used in basic authentication type.",
												},
												"password": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
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
									"tracing": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Tracing config(Deprecated, please use MeshConfig.Tracing for configuration).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sampling": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Computed:    true,
													Description: "Tracing sampling, 0.0-1.0.",
												},
												"enable": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether enable tracing.",
												},
												"apm": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "APM config.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"enable": {
																Type:        schema.TypeBool,
																Required:    true,
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
									"disable_policy_checks": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Disable policy checks.",
									},
									"enable_pilot_http": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Enable HTTP/1.0 support.",
									},
									"disable_http_retry": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Disable http retry.",
									},
									"smart_dns": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "SmartDNS configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"istio_meta_dns_capture": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Enable dns proxy.",
												},
												"istio_meta_dns_auto_allocate": {
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
									"exclude_ip_ranges": {
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
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "Sidecar limits.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Resource type name, `cpu/memory`.",
												},
												"quantity": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Resource quantity, example: cpu-`100m`, memory-`1Gi`.",
												},
											},
										},
									},
									"requests": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "Sidecar requests.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Resource type name, `cpu/memory`.",
												},
												"quantity": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Resource quantity, example: cpu-`100m`, memory-`1Gi`.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"tag_list": {
				Type:        schema.TypeList,
				Optional:    true,
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
	defer tccommon.LogElapsed("resource.tencentcloud_tcm_mesh.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = tcm.NewCreateMeshRequest()
		response = tcm.NewCreateMeshResponse()
		meshId   string
	)

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
			if apmMap, ok := helper.InterfaceToMap(tracingMap, "apm"); ok {
				apm := tcm.APM{}
				if v, ok := apmMap["enable"]; ok {
					apm.Enable = helper.Bool(v.(bool))
				}
				if v, ok := apmMap["region"]; ok {
					apm.Region = helper.String(v.(string))
				}
				if v, ok := apmMap["instance_id"]; ok {
					apm.InstanceId = helper.String(v.(string))
				}
				tracingConfig.APM = &apm
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
			if v, ok := istioMap["enable_pilot_http"]; ok {
				istioConfig.EnablePilotHTTP = helper.Bool(v.(bool))
			}
			if v, ok := istioMap["disable_http_retry"]; ok {
				istioConfig.DisableHTTPRetry = helper.Bool(v.(bool))
			}
			if smartDNSMap, ok := helper.InterfaceToMap(istioMap, "smart_dns"); ok {
				smartDNSConfig := tcm.SmartDNSConfig{}
				if v, ok := smartDNSMap["istio_meta_dns_capture"]; ok {
					smartDNSConfig.IstioMetaDNSCapture = helper.Bool(v.(bool))
				}
				if v, ok := smartDNSMap["istio_meta_dns_auto_allocate"]; ok {
					smartDNSConfig.IstioMetaDNSAutoAllocate = helper.Bool(v.(bool))
				}
				istioConfig.SmartDNS = &smartDNSConfig
			}
			if tracingMap, ok := helper.InterfaceToMap(istioMap, "tracing"); ok {
				tracingConfig := tcm.TracingConfig{}
				if v, ok := tracingMap["sampling"]; ok {
					tracingConfig.Sampling = helper.Float64(v.(float64))
				}
				if v, ok := tracingMap["enable"]; ok {
					tracingConfig.Enable = helper.Bool(v.(bool))
				}
				if apmMap, ok := helper.InterfaceToMap(tracingMap, "apm"); ok {
					apm := tcm.APM{}
					if v, ok := apmMap["enable"]; ok {
						apm.Enable = helper.Bool(v.(bool))
					}
					if v, ok := apmMap["region"]; ok {
						apm.Region = helper.String(v.(string))
					}
					if v, ok := apmMap["instance_id"]; ok {
						apm.InstanceId = helper.String(v.(string))
					}
					tracingConfig.APM = &apm
				}
				if zipkinMap, ok := helper.InterfaceToMap(tracingMap, "zipkin"); ok {
					tracingZipkin := tcm.TracingZipkin{}
					if v, ok := zipkinMap["address"]; ok {
						tracingZipkin.Address = helper.String(v.(string))
					}
					tracingConfig.Zipkin = &tracingZipkin
				}
				istioConfig.Tracing = &tracingConfig
			}
			meshConfig.Istio = &istioConfig
		}
		if injectMap, ok := helper.InterfaceToMap(dMap, "inject"); ok {
			injectConfig := tcm.InjectConfig{}
			if v, ok := injectMap["exclude_ip_ranges"]; ok {
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
				for _, item := range v.(*schema.Set).List() {
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
				for _, item := range v.(*schema.Set).List() {
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcmClient().CreateMesh(request)
		if e != nil {
			return tccommon.RetryError(e)
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

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := TcmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		mesh, errRet := service.DescribeTcmMesh(ctx, meshId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if *mesh.Mesh.State == "PENDING" || *mesh.Mesh.State == "CREATING" || *mesh.Mesh.State != "RUNNING" {
			return resource.RetryableError(fmt.Errorf("mesh status is %v, retry...", *mesh.Mesh.State))
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(meshId)
	return resourceTencentCloudTcmMeshRead(d, meta)
}

func resourceTencentCloudTcmMeshRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcm_mesh.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TcmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	meshId := d.Id()

	meshResponse, err := service.DescribeTcmMesh(ctx, meshId)
	if err != nil {
		return err
	}

	mesh := meshResponse.Mesh
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

	if mesh.Version != nil {
		_ = d.Set("mesh_version", mesh.Version)
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
				apmMap := map[string]interface{}{}
				if mesh.Config.Tracing.APM.Enable != nil {
					apmMap["enable"] = mesh.Config.Tracing.APM.Enable
				}
				if mesh.Config.Tracing.APM.Region != nil {
					apmMap["region"] = mesh.Config.Tracing.APM.Region
				}
				if mesh.Config.Tracing.APM.InstanceId != nil {
					apmMap["instance_id"] = mesh.Config.Tracing.APM.InstanceId
				}

				tracingMap["apm"] = []interface{}{apmMap}
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

			if mesh.Config.Istio.Tracing != nil {
				tracingMap := map[string]interface{}{}

				if mesh.Config.Istio.Tracing.Sampling != nil {
					tracingMap["sampling"] = mesh.Config.Istio.Tracing.Sampling
				}

				if mesh.Config.Istio.Tracing.Enable != nil {
					tracingMap["enable"] = mesh.Config.Istio.Tracing.Enable
				}

				if mesh.Config.Istio.Tracing.APM != nil {
					apmMap := map[string]interface{}{}

					if mesh.Config.Istio.Tracing.APM.Enable != nil {
						apmMap["enable"] = mesh.Config.Istio.Tracing.APM.Enable
					}

					if mesh.Config.Istio.Tracing.APM.Region != nil {
						apmMap["region"] = mesh.Config.Istio.Tracing.APM.Region
					}

					if mesh.Config.Istio.Tracing.APM.InstanceId != nil {
						apmMap["instance_id"] = mesh.Config.Istio.Tracing.APM.InstanceId
					}

					tracingMap["apm"] = []interface{}{apmMap}
				}

				if mesh.Config.Istio.Tracing.Zipkin != nil {
					zipkinMap := map[string]interface{}{}

					if mesh.Config.Istio.Tracing.Zipkin.Address != nil {
						zipkinMap["address"] = mesh.Config.Istio.Tracing.Zipkin.Address
					}

					tracingMap["zipkin"] = []interface{}{zipkinMap}
				}

				istioMap["tracing"] = []interface{}{tracingMap}
			}

			if mesh.Config.Istio.DisablePolicyChecks != nil {
				istioMap["disable_policy_checks"] = mesh.Config.Istio.DisablePolicyChecks
			}
			if mesh.Config.Istio.EnablePilotHTTP != nil {
				istioMap["enable_pilot_http"] = mesh.Config.Istio.EnablePilotHTTP
			}
			if mesh.Config.Istio.DisableHTTPRetry != nil {
				istioMap["disable_http_retry"] = mesh.Config.Istio.DisableHTTPRetry
			}
			if mesh.Config.Istio.SmartDNS != nil {
				smartDNSMap := map[string]interface{}{}
				if mesh.Config.Istio.SmartDNS.IstioMetaDNSCapture != nil {
					smartDNSMap["istio_meta_dns_capture"] = mesh.Config.Istio.SmartDNS.IstioMetaDNSCapture
				}
				if mesh.Config.Istio.SmartDNS.IstioMetaDNSAutoAllocate != nil {
					smartDNSMap["istio_meta_dns_auto_allocate"] = mesh.Config.Istio.SmartDNS.IstioMetaDNSAutoAllocate
				}

				istioMap["smart_dns"] = []interface{}{smartDNSMap}
			}

			configMap["istio"] = []interface{}{istioMap}
		}

		if mesh.Config.Inject != nil {
			injectMap := map[string]interface{}{}

			if mesh.Config.Inject.ExcludeIPRanges != nil {
				injectMap["exclude_ip_ranges"] = mesh.Config.Inject.ExcludeIPRanges
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

				sidecarResourcesMap["limits"] = limitsList
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

				sidecarResourcesMap["requests"] = requestsList
			}

			configMap["sidecar_resources"] = []interface{}{sidecarResourcesMap}
		}

		err = d.Set("config", []interface{}{configMap})
		if err != nil {
			return fmt.Errorf("set error, err: %v", err)
		}
	}

	if mesh.TagList != nil {
		tagListList := []interface{}{}
		for _, tagList := range mesh.TagList {
			tagListMap := map[string]interface{}{}
			if tagList.Key != nil {
				tagListMap["key"] = tagList.Key
			}
			if tagList.Value != nil {
				tagListMap["value"] = tagList.Value
			}
			if tagList.Passthrough != nil {
				tagListMap["passthrough"] = tagList.Passthrough
			}

			tagListList = append(tagListList, tagListMap)
		}
		_ = d.Set("tag_list", tagListList)
	}

	return nil
}

func resourceTencentCloudTcmMeshUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcm_mesh.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := tcm.NewModifyMeshRequest()

	meshId := d.Id()
	request.MeshId = &meshId

	immutableArgs := []string{"mesh_id", "mesh_version", "type", "tag_list"}

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
				if apmMap, ok := helper.InterfaceToMap(tracingMap, "apm"); ok {
					apm := tcm.APM{}
					if v, ok := apmMap["enable"]; ok {
						apm.Enable = helper.Bool(v.(bool))
					}
					if v, ok := apmMap["region"]; ok {
						apm.Region = helper.String(v.(string))
					}
					if v, ok := apmMap["instance_id"]; ok {
						apm.InstanceId = helper.String(v.(string))
					}
					tracingConfig.APM = &apm
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
				if v, ok := istioMap["enable_pilot_http"]; ok {
					istioConfig.EnablePilotHTTP = helper.Bool(v.(bool))
				}
				if v, ok := istioMap["disable_http_retry"]; ok {
					istioConfig.DisableHTTPRetry = helper.Bool(v.(bool))
				}
				if smartDNSMap, ok := helper.InterfaceToMap(istioMap, "smart_dns"); ok {
					smartDNSConfig := tcm.SmartDNSConfig{}
					if v, ok := smartDNSMap["istio_meta_dns_capture"]; ok {
						smartDNSConfig.IstioMetaDNSCapture = helper.Bool(v.(bool))
					}
					if v, ok := smartDNSMap["istio_meta_dns_auto_allocate"]; ok {
						smartDNSConfig.IstioMetaDNSAutoAllocate = helper.Bool(v.(bool))
					}
					istioConfig.SmartDNS = &smartDNSConfig
				}
				if tracingMap, ok := helper.InterfaceToMap(istioMap, "tracing"); ok {
					tracingConfig := tcm.TracingConfig{}
					if v, ok := tracingMap["sampling"]; ok {
						tracingConfig.Sampling = helper.Float64(v.(float64))
					}
					if v, ok := tracingMap["enable"]; ok {
						tracingConfig.Enable = helper.Bool(v.(bool))
					}
					if apmMap, ok := helper.InterfaceToMap(tracingMap, "apm"); ok {
						apm := tcm.APM{}
						if v, ok := apmMap["enable"]; ok {
							apm.Enable = helper.Bool(v.(bool))
						}
						if v, ok := apmMap["region"]; ok {
							apm.Region = helper.String(v.(string))
						}
						if v, ok := apmMap["instance_id"]; ok {
							apm.InstanceId = helper.String(v.(string))
						}
						tracingConfig.APM = &apm
					}
					if zipkinMap, ok := helper.InterfaceToMap(tracingMap, "zipkin"); ok {
						tracingZipkin := tcm.TracingZipkin{}
						if v, ok := zipkinMap["address"]; ok {
							tracingZipkin.Address = helper.String(v.(string))
						}
						tracingConfig.Zipkin = &tracingZipkin
					}
					istioConfig.Tracing = &tracingConfig
				}
				meshConfig.Istio = &istioConfig
			}
			if injectMap, ok := helper.InterfaceToMap(dMap, "inject"); ok {
				injectConfig := tcm.InjectConfig{}
				if v, ok := injectMap["exclude_ip_ranges"]; ok {
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
					for _, item := range v.(*schema.Set).List() {
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
					for _, item := range v.(*schema.Set).List() {
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
			request.Config = &meshConfig
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcmClient().ModifyMesh(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_tcm_mesh.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TcmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	meshId := d.Id()

	if err := service.DeleteTcmMeshById(ctx, meshId); err != nil {
		return err
	}

	err := resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		mesh, errRet := service.DescribeTcmMesh(ctx, meshId)
		if errRet != nil {
			if tccommon.IsExpectError(errRet, []string{"ResourceNotFound"}) {
				return nil
			}
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if mesh != nil {
			if *mesh.Mesh.State == "DELETING" {
				return resource.RetryableError(fmt.Errorf("mesh status is %v, retry...", *mesh.Mesh.State))
			}
			if *mesh.Mesh.State == "DELETE_FAILED" {
				return resource.NonRetryableError(fmt.Errorf("mesh status is %v, retry...", *mesh.Mesh.State))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
