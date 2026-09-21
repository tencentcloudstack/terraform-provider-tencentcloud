package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoInferenceService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoInferenceServiceCreate,
		Read:   resourceTencentCloudTeoInferenceServiceRead,
		Update: resourceTencentCloudTeoInferenceServiceUpdate,
		Delete: resourceTencentCloudTeoInferenceServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Inference service name, not accepted by ModifyInferenceService, hence ForceNew.",
			},

			"listen_port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Model service listen port, 1-65535.",
			},

			"containers": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Container config list, only support 1 container currently.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Image type, e.g. TCR.",
						},
						"tcr_repository_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "TCR repository config, required when image_type is TCR.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tcr_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "TCR service type: Personal or Enterprise.",
									},
									"image": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Image address.",
									},
									"registry_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Registry instance ID, required when tcr_type is Enterprise.",
									},
									"region_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Region name.",
									},
								},
							},
						},
						"startup_command": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Container startup command.",
						},
						"environment_variables": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Container environment variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Variable name.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Variable value.",
									},
								},
							},
						},
					},
				},
			},

			"resource_config": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Resource config of inference service.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scaling_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Scaling mode: Auto or Manual.",
						},
						"hardware_spec": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Hardware spec identifier, deprecated, use hardware_spec_id instead.",
						},
						"hardware_spec_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Hardware spec unique ID, only effective on create.",
						},
						"hardware_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Hardware config.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gpu_num": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "GPU card number, only effective on create.",
									},
									"cpu_num": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "CPU core number.",
									},
									"mem_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Memory size in MB.",
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Disk size in MB.",
									},
								},
							},
						},
						"auto_scaling_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Auto scaling config, required when scaling_mode is Auto.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_instance_count": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Min instance count.",
									},
									"scaling_policies": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Scaling policy list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"policy_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Policy name.",
												},
												"policy_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Policy type.",
												},
												"scheduled_scaling_policy": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Scheduled scaling policy.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"scheduled_actions": {
																Type:        schema.TypeList,
																Required:    true,
																Description: "Scheduled action list.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"cron_expression": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Cron expression.",
																		},
																		"min_instance_count": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Min instance count.",
																		},
																	},
																},
															},
															"effective_range": {
																Type:        schema.TypeList,
																Required:    true,
																MaxItems:    1,
																Description: "Effective range.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"effective_type": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Effective type: LongTerm or Custom.",
																		},
																		"start_date": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Start date, required when effective_type is Custom.",
																		},
																		"end_date": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "End date, required when effective_type is Custom.",
																		},
																	},
																},
															},
															"time_zone": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Time zone, e.g. UTC, Asia/Shanghai.",
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
						"manual_instance_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Manual instance config, required when scaling_mode is Manual.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fixed_instance_count": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Fixed instance count.",
									},
								},
							},
						},
						"concurrency": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Concurrency per instance, default 1.",
						},
					},
				},
			},

			"affinity_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Affinity config, write-only (not read back from DescribeInferenceServices).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Affinity switch: On or Off.",
						},
						"affinity_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Affinity mode, e.g. SessionId.",
						},
						"session_id_affinity_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Session ID affinity config, required when affinity_mode is SessionId.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Source of session id param, e.g. Header.",
									},
									"header_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Header name for session id.",
									},
								},
							},
						},
					},
				},
			},

			"request_paths": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Request path list, up to 20.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description, up to 60 chars.",
			},

			"service_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inference service ID.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inference service status.",
			},

			"scaling_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Scaling status.",
			},

			"current_instance_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current running instance count.",
			},

			"inference_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inference access URL.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update time.",
			},
		},
	}
}

func resourceTencentCloudTeoInferenceServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_service.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = teo.NewCreateInferenceServiceRequest()
		response = teo.NewCreateInferenceServiceResponse()
		zoneId   string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("listen_port"); ok {
		request.ListenPort = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("containers"); ok {
		request.Containers = buildTeoInferenceContainers(v.([]interface{}))
	}

	if v, ok := d.GetOk("resource_config"); ok {
		request.ResourceConfig = buildTeoInferenceResourceConfig(v.([]interface{}))
	}

	if v, ok := d.GetOk("affinity_config"); ok {
		request.AffinityConfig = buildTeoInferenceAffinityConfig(v.([]interface{}))
	}

	if v, ok := d.GetOk("request_paths"); ok {
		requestPaths := make([]*string, 0, len(v.([]interface{})))
		for _, item := range v.([]interface{}) {
			if item == nil {
				continue
			}
			requestPaths = append(requestPaths, helper.String(item.(string)))
		}
		request.RequestPaths = requestPaths
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateInferenceServiceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create inference_service failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create inference_service failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s inference_service logId=%s, d.Id()=%s", logId, logId, d.Id())

	if response.Response.ServiceId == nil || *response.Response.ServiceId == "" {
		return fmt.Errorf("ServiceId is nil or empty.")
	}

	serviceId := *response.Response.ServiceId
	d.SetId(strings.Join([]string{zoneId, serviceId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoInferenceServiceRead(d, meta)
}

func resourceTencentCloudTeoInferenceServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_service.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	serviceId := idSplit[1]

	respData, err := service.DescribeTeoInferenceServiceById(ctx, zoneId, serviceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] inference_service id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.ListenPort != nil {
		_ = d.Set("listen_port", respData.ListenPort)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.RequestPaths != nil {
		requestPathsList := make([]string, 0, len(respData.RequestPaths))
		for _, item := range respData.RequestPaths {
			if item == nil {
				continue
			}
			requestPathsList = append(requestPathsList, *item)
		}
		_ = d.Set("request_paths", requestPathsList)
	}

	if respData.Containers != nil && len(respData.Containers) > 0 {
		containersList := make([]map[string]interface{}, 0, len(respData.Containers))
		for _, container := range respData.Containers {
			containerMap := map[string]interface{}{}
			if container.ImageType != nil {
				containerMap["image_type"] = container.ImageType
			}
			if container.StartupCommand != nil {
				containerMap["startup_command"] = container.StartupCommand
			}
			if container.TcrRepositoryConfig != nil {
				tcrRepositoryConfigList := make([]map[string]interface{}, 0, 1)
				tcrRepositoryConfigMap := map[string]interface{}{}
				if container.TcrRepositoryConfig.TCRType != nil {
					tcrRepositoryConfigMap["tcr_type"] = container.TcrRepositoryConfig.TCRType
				}
				if container.TcrRepositoryConfig.Image != nil {
					tcrRepositoryConfigMap["image"] = container.TcrRepositoryConfig.Image
				}
				if container.TcrRepositoryConfig.RegistryId != nil {
					tcrRepositoryConfigMap["registry_id"] = container.TcrRepositoryConfig.RegistryId
				}
				if container.TcrRepositoryConfig.RegionName != nil {
					tcrRepositoryConfigMap["region_name"] = container.TcrRepositoryConfig.RegionName
				}
				tcrRepositoryConfigList = append(tcrRepositoryConfigList, tcrRepositoryConfigMap)
				containerMap["tcr_repository_config"] = tcrRepositoryConfigList
			}
			if container.EnvironmentVariables != nil && len(container.EnvironmentVariables) > 0 {
				environmentVariablesList := make([]map[string]interface{}, 0, len(container.EnvironmentVariables))
				for _, envVar := range container.EnvironmentVariables {
					envVarMap := map[string]interface{}{}
					if envVar.Key != nil {
						envVarMap["key"] = envVar.Key
					}
					if envVar.Value != nil {
						envVarMap["value"] = envVar.Value
					}
					environmentVariablesList = append(environmentVariablesList, envVarMap)
				}
				containerMap["environment_variables"] = environmentVariablesList
			}
			containersList = append(containersList, containerMap)
		}
		_ = d.Set("containers", containersList)
	}

	if respData.ResourceConfig != nil {
		resourceConfigList := make([]map[string]interface{}, 0, 1)
		resourceConfigMap := map[string]interface{}{}
		resourceConfig := respData.ResourceConfig
		if resourceConfig.ScalingMode != nil {
			resourceConfigMap["scaling_mode"] = resourceConfig.ScalingMode
		}
		if resourceConfig.HardwareSpec != nil {
			resourceConfigMap["hardware_spec"] = resourceConfig.HardwareSpec
		}
		if resourceConfig.HardwareSpecId != nil {
			resourceConfigMap["hardware_spec_id"] = resourceConfig.HardwareSpecId
		}
		if resourceConfig.HardwareConfig != nil {
			hardwareConfigList := make([]map[string]interface{}, 0, 1)
			hardwareConfigMap := map[string]interface{}{}
			if resourceConfig.HardwareConfig.GPUNum != nil {
				hardwareConfigMap["gpu_num"] = *resourceConfig.HardwareConfig.GPUNum
			}
			if resourceConfig.HardwareConfig.CPUNum != nil {
				hardwareConfigMap["cpu_num"] = *resourceConfig.HardwareConfig.CPUNum
			}
			if resourceConfig.HardwareConfig.MemSize != nil {
				hardwareConfigMap["mem_size"] = resourceConfig.HardwareConfig.MemSize
			}
			if resourceConfig.HardwareConfig.DiskSize != nil {
				hardwareConfigMap["disk_size"] = resourceConfig.HardwareConfig.DiskSize
			}
			hardwareConfigList = append(hardwareConfigList, hardwareConfigMap)
			resourceConfigMap["hardware_config"] = hardwareConfigList
		}
		if resourceConfig.AutoScalingConfig != nil {
			autoScalingConfigList := make([]map[string]interface{}, 0, 1)
			autoScalingConfigMap := map[string]interface{}{}
			if resourceConfig.AutoScalingConfig.MinInstanceCount != nil {
				autoScalingConfigMap["min_instance_count"] = resourceConfig.AutoScalingConfig.MinInstanceCount
			}
			if resourceConfig.AutoScalingConfig.ScalingPolicies != nil && len(resourceConfig.AutoScalingConfig.ScalingPolicies) > 0 {
				scalingPoliciesList := make([]map[string]interface{}, 0, len(resourceConfig.AutoScalingConfig.ScalingPolicies))
				for _, policy := range resourceConfig.AutoScalingConfig.ScalingPolicies {
					policyMap := map[string]interface{}{}
					if policy.PolicyName != nil {
						policyMap["policy_name"] = policy.PolicyName
					}
					if policy.PolicyType != nil {
						policyMap["policy_type"] = policy.PolicyType
					}
					if policy.ScheduledScalingPolicy != nil {
						scheduledScalingPolicyList := make([]map[string]interface{}, 0, 1)
						scheduledScalingPolicyMap := map[string]interface{}{}
						scheduledScalingPolicy := policy.ScheduledScalingPolicy
						if scheduledScalingPolicy.TimeZone != nil {
							scheduledScalingPolicyMap["time_zone"] = scheduledScalingPolicy.TimeZone
						}
						if scheduledScalingPolicy.ScheduledActions != nil && len(scheduledScalingPolicy.ScheduledActions) > 0 {
							scheduledActionsList := make([]map[string]interface{}, 0, len(scheduledScalingPolicy.ScheduledActions))
							for _, action := range scheduledScalingPolicy.ScheduledActions {
								actionMap := map[string]interface{}{}
								if action.CronExpression != nil {
									actionMap["cron_expression"] = action.CronExpression
								}
								if action.MinInstanceCount != nil {
									actionMap["min_instance_count"] = action.MinInstanceCount
								}
								scheduledActionsList = append(scheduledActionsList, actionMap)
							}
							scheduledScalingPolicyMap["scheduled_actions"] = scheduledActionsList
						}
						if scheduledScalingPolicy.EffectiveRange != nil {
							effectiveRangeList := make([]map[string]interface{}, 0, 1)
							effectiveRangeMap := map[string]interface{}{}
							if scheduledScalingPolicy.EffectiveRange.EffectiveType != nil {
								effectiveRangeMap["effective_type"] = scheduledScalingPolicy.EffectiveRange.EffectiveType
							}
							if scheduledScalingPolicy.EffectiveRange.StartDate != nil {
								effectiveRangeMap["start_date"] = scheduledScalingPolicy.EffectiveRange.StartDate
							}
							if scheduledScalingPolicy.EffectiveRange.EndDate != nil {
								effectiveRangeMap["end_date"] = scheduledScalingPolicy.EffectiveRange.EndDate
							}
							effectiveRangeList = append(effectiveRangeList, effectiveRangeMap)
							scheduledScalingPolicyMap["effective_range"] = effectiveRangeList
						}
						scheduledScalingPolicyList = append(scheduledScalingPolicyList, scheduledScalingPolicyMap)
						policyMap["scheduled_scaling_policy"] = scheduledScalingPolicyList
					}
					scalingPoliciesList = append(scalingPoliciesList, policyMap)
				}
				autoScalingConfigMap["scaling_policies"] = scalingPoliciesList
			}
			autoScalingConfigList = append(autoScalingConfigList, autoScalingConfigMap)
			resourceConfigMap["auto_scaling_config"] = autoScalingConfigList
		}
		if resourceConfig.ManualInstanceConfig != nil {
			manualInstanceConfigList := make([]map[string]interface{}, 0, 1)
			manualInstanceConfigMap := map[string]interface{}{}
			if resourceConfig.ManualInstanceConfig.FixedInstanceCount != nil {
				manualInstanceConfigMap["fixed_instance_count"] = resourceConfig.ManualInstanceConfig.FixedInstanceCount
			}
			manualInstanceConfigList = append(manualInstanceConfigList, manualInstanceConfigMap)
			resourceConfigMap["manual_instance_config"] = manualInstanceConfigList
		}
		if resourceConfig.Concurrency != nil {
			resourceConfigMap["concurrency"] = resourceConfig.Concurrency
		}
		resourceConfigList = append(resourceConfigList, resourceConfigMap)
		_ = d.Set("resource_config", resourceConfigList)
	}

	if respData.ServiceId != nil {
		_ = d.Set("service_id", respData.ServiceId)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.ScalingStatus != nil {
		_ = d.Set("scaling_status", respData.ScalingStatus)
	}

	if respData.CurrentInstanceCount != nil {
		_ = d.Set("current_instance_count", respData.CurrentInstanceCount)
	}

	if respData.InferenceURL != nil {
		_ = d.Set("inference_url", respData.InferenceURL)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	return nil
}

func resourceTencentCloudTeoInferenceServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_service.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	serviceId := idSplit[1]

	needChange := false
	mutableArgs := []string{"listen_port", "containers", "resource_config", "affinity_config", "request_paths", "description"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := teo.NewModifyInferenceServiceRequest()
		request.ZoneId = helper.String(zoneId)
		request.ServiceId = helper.String(serviceId)

		if v, ok := d.GetOk("listen_port"); ok {
			request.ListenPort = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("request_paths"); ok {
			requestPaths := make([]*string, 0, len(v.([]interface{})))
			for _, item := range v.([]interface{}) {
				if item == nil {
					continue
				}
				requestPaths = append(requestPaths, helper.String(item.(string)))
			}
			request.RequestPaths = requestPaths
		}

		if v, ok := d.GetOk("containers"); ok {
			request.Containers = buildTeoInferenceContainersForModify(v.([]interface{}))
		}

		if v, ok := d.GetOk("resource_config"); ok {
			request.ResourceConfig = buildTeoInferenceResourceConfigForModify(v.([]interface{}))
		}

		if v, ok := d.GetOk("affinity_config"); ok {
			request.AffinityConfig = buildTeoInferenceAffinityConfig(v.([]interface{}))
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyInferenceServiceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Update inference_service failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update inference_service failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTeoInferenceServiceRead(d, meta)
}

func resourceTencentCloudTeoInferenceServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_service.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	serviceId := idSplit[1]

	request := teo.NewOperateInferenceServiceRequest()
	request.ZoneId = helper.String(zoneId)
	request.ServiceId = helper.String(serviceId)
	request.Operation = helper.String("Delete")

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().OperateInferenceServiceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete inference_service failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete inference_service failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func buildTeoInferenceContainers(containersList []interface{}) []*teo.InferenceContainerConfig {
	containers := make([]*teo.InferenceContainerConfig, 0, len(containersList))
	for _, item := range containersList {
		if item == nil {
			continue
		}
		containerMap := item.(map[string]interface{})
		container := &teo.InferenceContainerConfig{}
		if v, ok := containerMap["image_type"].(string); ok && v != "" {
			container.ImageType = helper.String(v)
		}
		if v, ok := containerMap["tcr_repository_config"].([]interface{}); ok && len(v) > 0 {
			tcrMap := v[0].(map[string]interface{})
			tcrConfig := &teo.InferenceTCRRepositoryConfig{}
			if v, ok := tcrMap["tcr_type"].(string); ok && v != "" {
				tcrConfig.TCRType = helper.String(v)
			}
			if v, ok := tcrMap["image"].(string); ok && v != "" {
				tcrConfig.Image = helper.String(v)
			}
			if v, ok := tcrMap["registry_id"].(string); ok && v != "" {
				tcrConfig.RegistryId = helper.String(v)
			}
			if v, ok := tcrMap["region_name"].(string); ok && v != "" {
				tcrConfig.RegionName = helper.String(v)
			}
			container.TcrRepositoryConfig = tcrConfig
		}
		if v, ok := containerMap["startup_command"].(string); ok && v != "" {
			container.StartupCommand = helper.String(v)
		}
		if v, ok := containerMap["environment_variables"].([]interface{}); ok && len(v) > 0 {
			envVars := make([]*teo.InferenceEnvironmentVariable, 0, len(v))
			for _, envItem := range v {
				if envItem == nil {
					continue
				}
				envMap := envItem.(map[string]interface{})
				envVar := &teo.InferenceEnvironmentVariable{}
				if v, ok := envMap["key"].(string); ok && v != "" {
					envVar.Key = helper.String(v)
				}
				if v, ok := envMap["value"].(string); ok && v != "" {
					envVar.Value = helper.String(v)
				}
				envVars = append(envVars, envVar)
			}
			container.EnvironmentVariables = envVars
		}
		containers = append(containers, container)
	}
	return containers
}

func buildTeoInferenceContainersForModify(containersList []interface{}) []*teo.InferenceContainerConfigForModify {
	containers := make([]*teo.InferenceContainerConfigForModify, 0, len(containersList))
	for _, item := range containersList {
		if item == nil {
			continue
		}
		containerMap := item.(map[string]interface{})
		container := &teo.InferenceContainerConfigForModify{}
		if v, ok := containerMap["image_type"].(string); ok && v != "" {
			container.ImageType = helper.String(v)
		}
		if v, ok := containerMap["tcr_repository_config"].([]interface{}); ok && len(v) > 0 {
			tcrMap := v[0].(map[string]interface{})
			tcrConfig := &teo.InferenceTCRRepositoryConfig{}
			if v, ok := tcrMap["tcr_type"].(string); ok && v != "" {
				tcrConfig.TCRType = helper.String(v)
			}
			if v, ok := tcrMap["image"].(string); ok && v != "" {
				tcrConfig.Image = helper.String(v)
			}
			if v, ok := tcrMap["registry_id"].(string); ok && v != "" {
				tcrConfig.RegistryId = helper.String(v)
			}
			if v, ok := tcrMap["region_name"].(string); ok && v != "" {
				tcrConfig.RegionName = helper.String(v)
			}
			container.TcrRepositoryConfig = tcrConfig
		}
		if v, ok := containerMap["startup_command"].(string); ok && v != "" {
			container.StartupCommand = helper.String(v)
		}
		if v, ok := containerMap["environment_variables"].([]interface{}); ok && len(v) > 0 {
			envVars := make([]*teo.InferenceEnvironmentVariable, 0, len(v))
			for _, envItem := range v {
				if envItem == nil {
					continue
				}
				envMap := envItem.(map[string]interface{})
				envVar := &teo.InferenceEnvironmentVariable{}
				if v, ok := envMap["key"].(string); ok && v != "" {
					envVar.Key = helper.String(v)
				}
				if v, ok := envMap["value"].(string); ok && v != "" {
					envVar.Value = helper.String(v)
				}
				envVars = append(envVars, envVar)
			}
			container.EnvironmentVariables = envVars
		}
		containers = append(containers, container)
	}
	return containers
}

func buildTeoInferenceResourceConfig(resourceConfigList []interface{}) *teo.InferenceResourceConfig {
	if len(resourceConfigList) < 1 {
		return nil
	}
	resourceConfigMap := resourceConfigList[0].(map[string]interface{})
	resourceConfig := &teo.InferenceResourceConfig{}
	if v, ok := resourceConfigMap["scaling_mode"].(string); ok && v != "" {
		resourceConfig.ScalingMode = helper.String(v)
	}
	if v, ok := resourceConfigMap["hardware_spec"].(string); ok && v != "" {
		resourceConfig.HardwareSpec = helper.String(v)
	}
	if v, ok := resourceConfigMap["hardware_spec_id"].(string); ok && v != "" {
		resourceConfig.HardwareSpecId = helper.String(v)
	}
	if v, ok := resourceConfigMap["hardware_config"].([]interface{}); ok && len(v) > 0 {
		hardwareConfigMap := v[0].(map[string]interface{})
		hardwareConfig := &teo.InferenceHardwareConfig{}
		if v, ok := hardwareConfigMap["gpu_num"].(float64); ok && v != 0 {
			hardwareConfig.GPUNum = helper.Float64(v)
		}
		if v, ok := hardwareConfigMap["cpu_num"].(float64); ok && v != 0 {
			hardwareConfig.CPUNum = helper.Float64(v)
		}
		if v, ok := hardwareConfigMap["mem_size"].(int); ok && v != 0 {
			hardwareConfig.MemSize = helper.IntInt64(v)
		}
		if v, ok := hardwareConfigMap["disk_size"].(int); ok && v != 0 {
			hardwareConfig.DiskSize = helper.IntInt64(v)
		}
		resourceConfig.HardwareConfig = hardwareConfig
	}
	if v, ok := resourceConfigMap["auto_scaling_config"].([]interface{}); ok && len(v) > 0 {
		autoScalingConfigMap := v[0].(map[string]interface{})
		autoScalingConfig := &teo.InferenceAutoScalingConfig{}
		if v, ok := autoScalingConfigMap["min_instance_count"].(int); ok {
			autoScalingConfig.MinInstanceCount = helper.IntInt64(v)
		}
		if v, ok := autoScalingConfigMap["scaling_policies"].([]interface{}); ok && len(v) > 0 {
			scalingPolicies := make([]*teo.InferenceScalingPolicy, 0, len(v))
			for _, policyItem := range v {
				if policyItem == nil {
					continue
				}
				policyMap := policyItem.(map[string]interface{})
				policy := &teo.InferenceScalingPolicy{}
				if v, ok := policyMap["policy_name"].(string); ok && v != "" {
					policy.PolicyName = helper.String(v)
				}
				if v, ok := policyMap["policy_type"].(string); ok && v != "" {
					policy.PolicyType = helper.String(v)
				}
				if v, ok := policyMap["scheduled_scaling_policy"].([]interface{}); ok && len(v) > 0 {
					scheduledScalingPolicyMap := v[0].(map[string]interface{})
					scheduledScalingPolicy := &teo.InferenceScheduledScalingPolicy{}
					if v, ok := scheduledScalingPolicyMap["time_zone"].(string); ok && v != "" {
						scheduledScalingPolicy.TimeZone = helper.String(v)
					}
					if v, ok := scheduledScalingPolicyMap["scheduled_actions"].([]interface{}); ok && len(v) > 0 {
						scheduledActions := make([]*teo.InferenceScheduledScalingAction, 0, len(v))
						for _, actionItem := range v {
							if actionItem == nil {
								continue
							}
							actionMap := actionItem.(map[string]interface{})
							action := &teo.InferenceScheduledScalingAction{}
							if v, ok := actionMap["cron_expression"].(string); ok && v != "" {
								action.CronExpression = helper.String(v)
							}
							if v, ok := actionMap["min_instance_count"].(int); ok {
								action.MinInstanceCount = helper.IntInt64(v)
							}
							scheduledActions = append(scheduledActions, action)
						}
						scheduledScalingPolicy.ScheduledActions = scheduledActions
					}
					if v, ok := scheduledScalingPolicyMap["effective_range"].([]interface{}); ok && len(v) > 0 {
						effectiveRangeMap := v[0].(map[string]interface{})
						effectiveRange := &teo.InferenceScheduledScalingEffectiveRange{}
						if v, ok := effectiveRangeMap["effective_type"].(string); ok && v != "" {
							effectiveRange.EffectiveType = helper.String(v)
						}
						if v, ok := effectiveRangeMap["start_date"].(string); ok && v != "" {
							effectiveRange.StartDate = helper.String(v)
						}
						if v, ok := effectiveRangeMap["end_date"].(string); ok && v != "" {
							effectiveRange.EndDate = helper.String(v)
						}
						scheduledScalingPolicy.EffectiveRange = effectiveRange
					}
					policy.ScheduledScalingPolicy = scheduledScalingPolicy
				}
				scalingPolicies = append(scalingPolicies, policy)
			}
			autoScalingConfig.ScalingPolicies = scalingPolicies
		}
		resourceConfig.AutoScalingConfig = autoScalingConfig
	}
	if v, ok := resourceConfigMap["manual_instance_config"].([]interface{}); ok && len(v) > 0 {
		manualInstanceConfigMap := v[0].(map[string]interface{})
		manualInstanceConfig := &teo.InferenceManualInstanceConfig{}
		if v, ok := manualInstanceConfigMap["fixed_instance_count"].(int); ok {
			manualInstanceConfig.FixedInstanceCount = helper.IntInt64(v)
		}
		resourceConfig.ManualInstanceConfig = manualInstanceConfig
	}
	if v, ok := resourceConfigMap["concurrency"].(int); ok && v != 0 {
		resourceConfig.Concurrency = helper.IntInt64(v)
	}
	return resourceConfig
}

func buildTeoInferenceResourceConfigForModify(resourceConfigList []interface{}) *teo.InferenceResourceConfigForModify {
	if len(resourceConfigList) < 1 {
		return nil
	}
	resourceConfigMap := resourceConfigList[0].(map[string]interface{})
	resourceConfig := &teo.InferenceResourceConfigForModify{}
	if v, ok := resourceConfigMap["scaling_mode"].(string); ok && v != "" {
		resourceConfig.ScalingMode = helper.String(v)
	}
	if v, ok := resourceConfigMap["hardware_config"].([]interface{}); ok && len(v) > 0 {
		hardwareConfigMap := v[0].(map[string]interface{})
		hardwareConfig := &teo.InferenceHardwareConfigForModify{}
		if v, ok := hardwareConfigMap["cpu_num"].(float64); ok && v != 0 {
			hardwareConfig.CPUNum = helper.Float64(v)
		}
		if v, ok := hardwareConfigMap["mem_size"].(int); ok && v != 0 {
			hardwareConfig.MemSize = helper.IntInt64(v)
		}
		if v, ok := hardwareConfigMap["disk_size"].(int); ok && v != 0 {
			hardwareConfig.DiskSize = helper.IntInt64(v)
		}
		resourceConfig.HardwareConfig = hardwareConfig
	}
	if v, ok := resourceConfigMap["auto_scaling_config"].([]interface{}); ok && len(v) > 0 {
		autoScalingConfigMap := v[0].(map[string]interface{})
		autoScalingConfig := &teo.InferenceAutoScalingConfig{}
		if v, ok := autoScalingConfigMap["min_instance_count"].(int); ok {
			autoScalingConfig.MinInstanceCount = helper.IntInt64(v)
		}
		if v, ok := autoScalingConfigMap["scaling_policies"].([]interface{}); ok && len(v) > 0 {
			scalingPolicies := make([]*teo.InferenceScalingPolicy, 0, len(v))
			for _, policyItem := range v {
				if policyItem == nil {
					continue
				}
				policyMap := policyItem.(map[string]interface{})
				policy := &teo.InferenceScalingPolicy{}
				if v, ok := policyMap["policy_name"].(string); ok && v != "" {
					policy.PolicyName = helper.String(v)
				}
				if v, ok := policyMap["policy_type"].(string); ok && v != "" {
					policy.PolicyType = helper.String(v)
				}
				if v, ok := policyMap["scheduled_scaling_policy"].([]interface{}); ok && len(v) > 0 {
					scheduledScalingPolicyMap := v[0].(map[string]interface{})
					scheduledScalingPolicy := &teo.InferenceScheduledScalingPolicy{}
					if v, ok := scheduledScalingPolicyMap["time_zone"].(string); ok && v != "" {
						scheduledScalingPolicy.TimeZone = helper.String(v)
					}
					if v, ok := scheduledScalingPolicyMap["scheduled_actions"].([]interface{}); ok && len(v) > 0 {
						scheduledActions := make([]*teo.InferenceScheduledScalingAction, 0, len(v))
						for _, actionItem := range v {
							if actionItem == nil {
								continue
							}
							actionMap := actionItem.(map[string]interface{})
							action := &teo.InferenceScheduledScalingAction{}
							if v, ok := actionMap["cron_expression"].(string); ok && v != "" {
								action.CronExpression = helper.String(v)
							}
							if v, ok := actionMap["min_instance_count"].(int); ok {
								action.MinInstanceCount = helper.IntInt64(v)
							}
							scheduledActions = append(scheduledActions, action)
						}
						scheduledScalingPolicy.ScheduledActions = scheduledActions
					}
					if v, ok := scheduledScalingPolicyMap["effective_range"].([]interface{}); ok && len(v) > 0 {
						effectiveRangeMap := v[0].(map[string]interface{})
						effectiveRange := &teo.InferenceScheduledScalingEffectiveRange{}
						if v, ok := effectiveRangeMap["effective_type"].(string); ok && v != "" {
							effectiveRange.EffectiveType = helper.String(v)
						}
						if v, ok := effectiveRangeMap["start_date"].(string); ok && v != "" {
							effectiveRange.StartDate = helper.String(v)
						}
						if v, ok := effectiveRangeMap["end_date"].(string); ok && v != "" {
							effectiveRange.EndDate = helper.String(v)
						}
						scheduledScalingPolicy.EffectiveRange = effectiveRange
					}
					policy.ScheduledScalingPolicy = scheduledScalingPolicy
				}
				scalingPolicies = append(scalingPolicies, policy)
			}
			autoScalingConfig.ScalingPolicies = scalingPolicies
		}
		resourceConfig.AutoScalingConfig = autoScalingConfig
	}
	if v, ok := resourceConfigMap["manual_instance_config"].([]interface{}); ok && len(v) > 0 {
		manualInstanceConfigMap := v[0].(map[string]interface{})
		manualInstanceConfig := &teo.InferenceManualInstanceConfig{}
		if v, ok := manualInstanceConfigMap["fixed_instance_count"].(int); ok {
			manualInstanceConfig.FixedInstanceCount = helper.IntInt64(v)
		}
		resourceConfig.ManualInstanceConfig = manualInstanceConfig
	}
	if v, ok := resourceConfigMap["concurrency"].(int); ok && v != 0 {
		resourceConfig.Concurrency = helper.IntInt64(v)
	}
	return resourceConfig
}

func buildTeoInferenceAffinityConfig(affinityConfigList []interface{}) *teo.InferenceAffinityConfig {
	if len(affinityConfigList) < 1 {
		return nil
	}
	affinityConfigMap := affinityConfigList[0].(map[string]interface{})
	affinityConfig := &teo.InferenceAffinityConfig{}
	if v, ok := affinityConfigMap["switch"].(string); ok && v != "" {
		affinityConfig.Switch = helper.String(v)
	}
	if v, ok := affinityConfigMap["affinity_mode"].(string); ok && v != "" {
		affinityConfig.AffinityMode = helper.String(v)
	}
	if v, ok := affinityConfigMap["session_id_affinity_config"].([]interface{}); ok && len(v) > 0 {
		sessionIdMap := v[0].(map[string]interface{})
		sessionIdConfig := &teo.SessionIdAffinityConfig{}
		if v, ok := sessionIdMap["source"].(string); ok && v != "" {
			sessionIdConfig.Source = helper.String(v)
		}
		if v, ok := sessionIdMap["header_name"].(string); ok && v != "" {
			sessionIdConfig.HeaderName = helper.String(v)
		}
		affinityConfig.SessionIdAffinityConfig = sessionIdConfig
	}
	return affinityConfig
}
