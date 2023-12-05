package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTseGatewayServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseGatewayServicesRead,
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway ID.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "filter conditions, valid value:name,upstreamType.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "filter name.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "filter value.",
						},
					},
				},
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "result.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "service list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "service ID.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "service name.",
									},
									"tags": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "tag list.",
									},
									"upstream_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "upstream information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"host": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "an IP address or domain name.",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "port.",
												},
												"source_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "service source ID.",
												},
												"namespace": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "namespace.",
												},
												"service_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "the name of the service in registry or kubernetes.",
												},
												"targets": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "provided when service type is IPList.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"host": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Host.",
															},
															"port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "port.",
															},
															"weight": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "weight.",
															},
															"health": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "health.",
															},
															"created_time": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "created time.",
															},
															"source": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "source of target.",
															},
														},
													},
												},
												"source_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "source service type.",
												},
												"scf_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "scf lambda type.",
												},
												"scf_namespace": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "scf lambda namespace.",
												},
												"scf_lambda_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "scf lambda name.",
												},
												"scf_lambda_qualifier": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "scf lambda version.",
												},
												"slow_start": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "slow start time, unit:second,when it&#39;s enabled, weight of the node is increased from 1 to the target value gradually.",
												},
												"algorithm": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "load balance algorithm,default:round-robin,least-connections and consisten_hashing also support.",
												},
												"auto_scaling_group_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "auto scaling group ID of cvm.",
												},
												"auto_scaling_cvm_port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "auto scaling group port of cvm.",
												},
												"auto_scaling_tat_cmd_status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "tat cmd status in auto scaling group of cvm.",
												},
												"auto_scaling_hook_status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "hook status in auto scaling group of cvm.",
												},
												"source_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "the name of source service.",
												},
												"real_source_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "exact source service type.",
												},
											},
										},
									},
									"upstream_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "service type.",
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "created time.",
									},
									"editable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "editable status.",
									},
								},
							},
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "total count.",
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

func dataSourceTencentCloudTseGatewayServicesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tse_gateway_services.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("gateway_id"); ok {
		paramMap["GatewayId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tse.ListFilter, 0, len(filtersSet))

		for _, item := range filtersSet {
			listFilter := tse.ListFilter{}
			listFilterMap := item.(map[string]interface{})

			if v, ok := listFilterMap["key"]; ok {
				listFilter.Key = helper.String(v.(string))
			}
			if v, ok := listFilterMap["value"]; ok {
				listFilter.Value = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &listFilter)
		}
		paramMap["filters"] = tmpSet
	}

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result *tse.KongServices
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTseGatewayServicesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result.ServiceList))
	kongServicesMap := map[string]interface{}{}
	if result != nil {
		if result.ServiceList != nil {
			serviceListList := []interface{}{}
			for _, serviceList := range result.ServiceList {
				serviceListMap := map[string]interface{}{}

				if serviceList.ID != nil {
					serviceListMap["id"] = serviceList.ID
				}

				if serviceList.Name != nil {
					serviceListMap["name"] = serviceList.Name
				}

				if serviceList.Tags != nil {
					serviceListMap["tags"] = serviceList.Tags
				}

				if serviceList.UpstreamInfo != nil {
					upstreamInfoMap := map[string]interface{}{}

					if serviceList.UpstreamInfo.Host != nil {
						upstreamInfoMap["host"] = serviceList.UpstreamInfo.Host
					}

					if serviceList.UpstreamInfo.Port != nil {
						upstreamInfoMap["port"] = serviceList.UpstreamInfo.Port
					}

					if serviceList.UpstreamInfo.SourceID != nil {
						upstreamInfoMap["source_id"] = serviceList.UpstreamInfo.SourceID
					}

					if serviceList.UpstreamInfo.Namespace != nil {
						upstreamInfoMap["namespace"] = serviceList.UpstreamInfo.Namespace
					}

					if serviceList.UpstreamInfo.ServiceName != nil {
						upstreamInfoMap["service_name"] = serviceList.UpstreamInfo.ServiceName
					}

					if serviceList.UpstreamInfo.Targets != nil {
						targetsList := []interface{}{}
						for _, targets := range serviceList.UpstreamInfo.Targets {
							targetsMap := map[string]interface{}{}

							if targets.Host != nil {
								targetsMap["host"] = targets.Host
							}

							if targets.Port != nil {
								targetsMap["port"] = targets.Port
							}

							if targets.Weight != nil {
								targetsMap["weight"] = targets.Weight
							}

							if targets.Health != nil {
								targetsMap["health"] = targets.Health
							}

							if targets.CreatedTime != nil {
								targetsMap["created_time"] = targets.CreatedTime
							}

							if targets.Source != nil {
								targetsMap["source"] = targets.Source
							}

							targetsList = append(targetsList, targetsMap)
						}

						upstreamInfoMap["targets"] = targetsList
					}

					if serviceList.UpstreamInfo.SourceType != nil {
						upstreamInfoMap["source_type"] = serviceList.UpstreamInfo.SourceType
					}

					if serviceList.UpstreamInfo.ScfType != nil {
						upstreamInfoMap["scf_type"] = serviceList.UpstreamInfo.ScfType
					}

					if serviceList.UpstreamInfo.ScfNamespace != nil {
						upstreamInfoMap["scf_namespace"] = serviceList.UpstreamInfo.ScfNamespace
					}

					if serviceList.UpstreamInfo.ScfLambdaName != nil {
						upstreamInfoMap["scf_lambda_name"] = serviceList.UpstreamInfo.ScfLambdaName
					}

					if serviceList.UpstreamInfo.ScfLambdaQualifier != nil {
						upstreamInfoMap["scf_lambda_qualifier"] = serviceList.UpstreamInfo.ScfLambdaQualifier
					}

					if serviceList.UpstreamInfo.SlowStart != nil {
						upstreamInfoMap["slow_start"] = serviceList.UpstreamInfo.SlowStart
					}

					if serviceList.UpstreamInfo.Algorithm != nil {
						upstreamInfoMap["algorithm"] = serviceList.UpstreamInfo.Algorithm
					}

					if serviceList.UpstreamInfo.AutoScalingGroupID != nil {
						upstreamInfoMap["auto_scaling_group_id"] = serviceList.UpstreamInfo.AutoScalingGroupID
					}

					if serviceList.UpstreamInfo.AutoScalingCvmPort != nil {
						upstreamInfoMap["auto_scaling_cvm_port"] = serviceList.UpstreamInfo.AutoScalingCvmPort
					}

					if serviceList.UpstreamInfo.AutoScalingTatCmdStatus != nil {
						upstreamInfoMap["auto_scaling_tat_cmd_status"] = serviceList.UpstreamInfo.AutoScalingTatCmdStatus
					}

					if serviceList.UpstreamInfo.AutoScalingHookStatus != nil {
						upstreamInfoMap["auto_scaling_hook_status"] = serviceList.UpstreamInfo.AutoScalingHookStatus
					}

					if serviceList.UpstreamInfo.SourceName != nil {
						upstreamInfoMap["source_name"] = serviceList.UpstreamInfo.SourceName
					}

					if serviceList.UpstreamInfo.RealSourceType != nil {
						upstreamInfoMap["real_source_type"] = serviceList.UpstreamInfo.RealSourceType
					}

					serviceListMap["upstream_info"] = []interface{}{upstreamInfoMap}
				}

				if serviceList.UpstreamType != nil {
					serviceListMap["upstream_type"] = serviceList.UpstreamType
				}

				if serviceList.CreatedTime != nil {
					serviceListMap["created_time"] = serviceList.CreatedTime
				}

				if serviceList.Editable != nil {
					serviceListMap["editable"] = serviceList.Editable
				}

				serviceListList = append(serviceListList, serviceListMap)
				ids = append(ids, *serviceList.ID)
			}

			kongServicesMap["service_list"] = serviceListList
		}

		if result.TotalCount != nil {
			kongServicesMap["total_count"] = result.TotalCount
		}

		_ = d.Set("result", []interface{}{kongServicesMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), kongServicesMap); e != nil {
			return e
		}
	}
	return nil
}
