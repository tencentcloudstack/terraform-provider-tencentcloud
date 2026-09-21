package teo

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoInferenceServiceDeploymentRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoInferenceServiceDeploymentRecordsRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Inference service ID.",
			},

			"sort_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "create-time",
				Description: "Sort field, value: create-time. Default value: create-time.",
			},

			"sort_order": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "desc",
				Description: "Sort method, value: asc / desc. Default value: desc.",
			},

			"record_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Inference service deployment record list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment record ID.",
						},
						"operation": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment operation type. Value: create/update/resume/stop.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment status. Value: processing/succeeded/failed.",
						},
						"duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Deployment duration, unit: seconds.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment initiation time, using ISO date format.",
						},
						"active_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the deployment configuration is the currently effective configuration. Value: active/inactive.",
						},
						"inference_service_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Inference service deployment configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"listen_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port that the model service needs to listen on.",
									},
									"request_paths": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Request path list of the inference service.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"containers": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Container configuration list of the inference service.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"image_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Image type. Value: TCR.",
												},
												"tcr_repository_config": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "TCR image repository configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"tcr_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "TCR service type. Value: Personal/Enterprise.",
															},
															"image": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Image address.",
															},
															"registry_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Image repository instance ID. Required when TCRType = Enterprise.",
															},
															"region_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Region name.",
															},
														},
													},
												},
												"startup_command": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Container startup command.",
												},
												"environment_variables": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Container runtime environment variable list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Variable name.",
															},
															"value": {
																Type:        schema.TypeString,
																Computed:    true,
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
										Computed:    true,
										Description: "Resource configuration of the inference service.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"scaling_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Scaling mode. Value: Auto/Manual.",
												},
												"hardware_spec": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Hardware specification identifier. Deprecated.",
												},
												"hardware_spec_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Hardware specification unique identifier ID.",
												},
												"hardware_config": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Hardware configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"gpu_num": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Number of GPU cards allocated to a single instance of the inference service.",
															},
															"cpu_num": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Number of CPU cores allocated to a single instance of the inference service.",
															},
															"mem_size": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Memory size allocated to a single instance of the inference service. Unit: MB.",
															},
															"disk_size": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Temporary disk size allocated to a single instance of the inference service. Unit: MB.",
															},
														},
													},
												},
												"auto_scaling_config": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Automatic scaling configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"min_instance_count": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum number of instances.",
															},
															"scaling_policies": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Scaling policy list.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"policy_name": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Policy name.",
																		},
																		"policy_type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Policy type. Value: ScheduledScaling.",
																		},
																		"scheduled_scaling_policy": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Scheduled scaling configuration.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"scheduled_actions": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Scheduled scaling action list.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"cron_expression": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Cron expression.",
																								},
																								"min_instance_count": {
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Minimum instance count after the scheduled scaling action is triggered.",
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
														},
													},
												},
												"manual_instance_config": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Manual instance configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"fixed_instance_count": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Fixed instance count.",
															},
														},
													},
												},
												"concurrency": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Concurrency of a single instance.",
												},
											},
										},
									},
									"affinity_config": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Affinity configuration of the inference service.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Inference service affinity switch. Value: On/Off.",
												},
												"affinity_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Inference service affinity mode. Value: SessionId.",
												},
												"session_id_affinity_config": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Session ID affinity configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"source": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Location where the session ID parameter is passed. Value: Header.",
															},
															"header_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The request header name used to pass the session ID.",
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
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTeoInferenceServiceDeploymentRecordsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_inference_service_deployment_records.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(nil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId    string
		serviceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
		zoneId = v.(string)
	}
	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
		serviceId = v.(string)
	}
	if v, ok := d.GetOk("sort_by"); ok {
		paramMap["SortBy"] = helper.String(v.(string))
	}
	if v, ok := d.GetOk("sort_order"); ok {
		paramMap["SortOrder"] = helper.String(v.(string))
	}

	var respData []*teov20220901.InferenceServiceDeploymentRecord
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTeoInferenceServiceDeploymentRecordsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || len(result) == 0 {
			log.Printf("[DATASOURCE] read empty, skip SetId")
			return resource.NonRetryableError(e)
		}
		respData = result
		return nil
	})
	if reqErr != nil {
		log.Printf("[CRITAL]%s read tencentcloud_teo_inference_service_deployment_records failed, reason:%s\n", logId, reqErr.Error())
		return reqErr
	}

	ids := make([]string, 0, len(respData))
	recordSetList := make([]map[string]interface{}, 0, len(respData))
	for _, record := range respData {
		recordMap := map[string]interface{}{}

		if record.RecordId != nil {
			recordMap["record_id"] = record.RecordId
			ids = append(ids, *record.RecordId)
		}

		if record.Operation != nil {
			recordMap["operation"] = record.Operation
		}

		if record.Status != nil {
			recordMap["status"] = record.Status
		}

		if record.Duration != nil {
			recordMap["duration"] = int(*record.Duration)
		}

		if record.CreateTime != nil {
			recordMap["create_time"] = record.CreateTime
		}

		if record.ActiveStatus != nil {
			recordMap["active_status"] = record.ActiveStatus
		}

		if record.InferenceServiceConfig != nil {
			inferenceServiceConfig := record.InferenceServiceConfig
			inferenceServiceConfigMap := map[string]interface{}{}

			if inferenceServiceConfig.ListenPort != nil {
				inferenceServiceConfigMap["listen_port"] = int(*inferenceServiceConfig.ListenPort)
			}

			if inferenceServiceConfig.RequestPaths != nil {
				inferenceServiceConfigMap["request_paths"] = inferenceServiceConfig.RequestPaths
			}

			if inferenceServiceConfig.Containers != nil {
				containersList := make([]map[string]interface{}, 0, len(inferenceServiceConfig.Containers))
				for _, container := range inferenceServiceConfig.Containers {
					containerMap := map[string]interface{}{}

					if container.ImageType != nil {
						containerMap["image_type"] = container.ImageType
					}

					if container.TcrRepositoryConfig != nil {
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
						containerMap["tcr_repository_config"] = []interface{}{tcrRepositoryConfigMap}
					}

					if container.StartupCommand != nil {
						containerMap["startup_command"] = container.StartupCommand
					}

					if container.EnvironmentVariables != nil {
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
				inferenceServiceConfigMap["containers"] = containersList
			}

			if inferenceServiceConfig.ResourceConfig != nil {
				resourceConfig := inferenceServiceConfig.ResourceConfig
				resourceConfigMap := map[string]interface{}{}

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
					hardwareConfigMap := map[string]interface{}{}
					if resourceConfig.HardwareConfig.GPUNum != nil {
						hardwareConfigMap["gpu_num"] = int(*resourceConfig.HardwareConfig.GPUNum)
					}
					if resourceConfig.HardwareConfig.CPUNum != nil {
						hardwareConfigMap["cpu_num"] = int(*resourceConfig.HardwareConfig.CPUNum)
					}
					if resourceConfig.HardwareConfig.MemSize != nil {
						hardwareConfigMap["mem_size"] = int(*resourceConfig.HardwareConfig.MemSize)
					}
					if resourceConfig.HardwareConfig.DiskSize != nil {
						hardwareConfigMap["disk_size"] = int(*resourceConfig.HardwareConfig.DiskSize)
					}
					resourceConfigMap["hardware_config"] = []interface{}{hardwareConfigMap}
				}

				if resourceConfig.AutoScalingConfig != nil {
					autoScalingConfigMap := map[string]interface{}{}

					if resourceConfig.AutoScalingConfig.MinInstanceCount != nil {
						autoScalingConfigMap["min_instance_count"] = int(*resourceConfig.AutoScalingConfig.MinInstanceCount)
					}

					if resourceConfig.AutoScalingConfig.ScalingPolicies != nil {
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
								scheduledScalingPolicyMap := map[string]interface{}{}

								if policy.ScheduledScalingPolicy.ScheduledActions != nil {
									scheduledActionsList := make([]map[string]interface{}, 0, len(policy.ScheduledScalingPolicy.ScheduledActions))
									for _, action := range policy.ScheduledScalingPolicy.ScheduledActions {
										actionMap := map[string]interface{}{}
										if action.CronExpression != nil {
											actionMap["cron_expression"] = action.CronExpression
										}
										if action.MinInstanceCount != nil {
											actionMap["min_instance_count"] = int(*action.MinInstanceCount)
										}
										scheduledActionsList = append(scheduledActionsList, actionMap)
									}
									scheduledScalingPolicyMap["scheduled_actions"] = scheduledActionsList
								}

								policyMap["scheduled_scaling_policy"] = []interface{}{scheduledScalingPolicyMap}
							}

							scalingPoliciesList = append(scalingPoliciesList, policyMap)
						}
						autoScalingConfigMap["scaling_policies"] = scalingPoliciesList
					}

					resourceConfigMap["auto_scaling_config"] = []interface{}{autoScalingConfigMap}
				}

				if resourceConfig.ManualInstanceConfig != nil {
					manualInstanceConfigMap := map[string]interface{}{}
					if resourceConfig.ManualInstanceConfig.FixedInstanceCount != nil {
						manualInstanceConfigMap["fixed_instance_count"] = int(*resourceConfig.ManualInstanceConfig.FixedInstanceCount)
					}
					resourceConfigMap["manual_instance_config"] = []interface{}{manualInstanceConfigMap}
				}

				if resourceConfig.Concurrency != nil {
					resourceConfigMap["concurrency"] = int(*resourceConfig.Concurrency)
				}

				inferenceServiceConfigMap["resource_config"] = []interface{}{resourceConfigMap}
			}

			if inferenceServiceConfig.AffinityConfig != nil {
				affinityConfig := inferenceServiceConfig.AffinityConfig
				affinityConfigMap := map[string]interface{}{}

				if affinityConfig.Switch != nil {
					affinityConfigMap["switch"] = affinityConfig.Switch
				}

				if affinityConfig.AffinityMode != nil {
					affinityConfigMap["affinity_mode"] = affinityConfig.AffinityMode
				}

				if affinityConfig.SessionIdAffinityConfig != nil {
					sessionIdAffinityConfigMap := map[string]interface{}{}
					if affinityConfig.SessionIdAffinityConfig.Source != nil {
						sessionIdAffinityConfigMap["source"] = affinityConfig.SessionIdAffinityConfig.Source
					}
					if affinityConfig.SessionIdAffinityConfig.HeaderName != nil {
						sessionIdAffinityConfigMap["header_name"] = affinityConfig.SessionIdAffinityConfig.HeaderName
					}
					affinityConfigMap["session_id_affinity_config"] = []interface{}{sessionIdAffinityConfigMap}
				}

				inferenceServiceConfigMap["affinity_config"] = []interface{}{affinityConfigMap}
			}

			recordMap["inference_service_config"] = []interface{}{inferenceServiceConfigMap}
		}

		recordSetList = append(recordSetList, recordMap)
	}

	_ = d.Set("record_set", recordSetList)

	d.SetId(zoneId + tccommon.FILED_SP + serviceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), recordSetList); e != nil {
			return e
		}
	}

	return nil
}
