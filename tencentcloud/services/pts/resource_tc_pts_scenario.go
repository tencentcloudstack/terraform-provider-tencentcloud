package pts

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPtsScenario() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudPtsScenarioRead,
		Create: resourceTencentCloudPtsScenarioCreate,
		Update: resourceTencentCloudPtsScenarioUpdate,
		Delete: resourceTencentCloudPtsScenarioDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Pts Scenario name.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Pressure test engine type.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project id.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Pts Scenario Description.",
			},

			"load": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Pressure allocation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_spec": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Pressure allocation.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"concurrency": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Configuration of concurrent pressure mode.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"stages": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Multi-phase configuration array.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"duration_seconds": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Pressure time.",
															},
															"target_virtual_users": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Number of virtual users.",
															},
														},
													},
												},
												"iteration_count": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Number of runs.",
												},
												"max_requests_per_second": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Maximum RPS.",
												},
												"graceful_stop_seconds": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Wait time for graceful termination of the task.",
												},
											},
										},
									},
									"requests_per_second": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Configuration of RPS pressure mode.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"max_requests_per_second": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Maximum RPS.",
												},
												"duration_seconds": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Pressure time.",
												},
												"resources": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Number of resources.",
												},
												"start_requests_per_second": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Initial RPS.",
												},
												"target_requests_per_second": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Target RPS, invalid input parameter.",
												},
												"graceful_stop_seconds": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Elegant shutdown waiting time.",
												},
											},
										},
									},
									"script_origin": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Built-in stress mode in script.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"machine_number": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Number of machines.",
												},
												"machine_specification": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Machine specification.",
												},
												"duration_seconds": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Pressure testing time.",
												},
											},
										},
									},
								},
							},
						},
						"vpc_load_distribution": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Source of stress.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Regional ID.",
									},
									"region": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Region.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "VPC ID.",
									},
									"subnet_ids": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Subnet ID list.",
									},
								},
							},
						},
						"geo_regions_load_distribution": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Pressure distribution.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Regional ID.",
									},
									"region": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Region.",
									},
									"percentage": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Percentage.",
									},
								},
							},
						},
					},
				},
			},

			"datasets": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Test data set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The file name where the test dataset is located.",
						},
						"split": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Test whether the dataset is fragmented.",
						},
						"header_in_file": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether the first line is the parameter name.",
						},
						"header_columns": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Parameter name array.",
						},
						"line_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of file lines.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Update time.",
						},
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of file bytes.",
						},
						"head_lines": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Header data row.",
						},
						"tail_lines": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Trailing data row.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File type.",
						},
						"file_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File ID.",
						},
					},
				},
			},

			"extensions": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "deprecated.",
			},

			"cron_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "cron job ID.",
			},

			"test_scripts": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Test script file information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File name.",
						},
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "File size.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File type.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Update time.",
						},
						"encoded_content": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Base64 encoded file content.",
						},
						"encoded_http_archive": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Base64 encoded har structure.",
						},
						"load_weight": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Script weight, range 1-100.",
						},
					},
				},
			},

			"protocols": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Protocol file path.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Protocol name.",
						},
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "File name.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File type.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Update time.",
						},
						"file_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File ID.",
						},
					},
				},
			},

			"request_files": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Request file path.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File name.",
						},
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "File size.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File type.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Update time.",
						},
						"file_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File id.",
						},
					},
				},
			},

			"sla_policy": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "SLA strategy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sla_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "SLA rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Pressure test index.",
									},
									"aggregation": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Aggregation method of pressure test index.",
									},
									"condition": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Pressure test index condition judgment symbol.",
									},
									"value": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Threshold value.",
									},
									"label_filter": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "tag.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"label_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Label name.",
												},
												"label_value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Label value.",
												},
											},
										},
									},
									"abort_flag": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to stop the stress test task.",
									},
									"for": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "duraion.",
									},
								},
							},
						},
						"alert_channel": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Alarm notification channel.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"notice_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Notification template ID.",
									},
									"amp_consumer_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "AMP consumer ID.",
									},
								},
							},
						},
					},
				},
			},

			"plugins": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "SLA strategy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File name.",
						},
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "File size.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File type.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Update time.",
						},
						"file_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File id.",
						},
					},
				},
			},

			"domain_name_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Domain name resolution configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_aliases": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Domain name binding configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host_names": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "List of domain names to be bound.",
									},
									"ip": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The IP address to be bound.",
									},
								},
							},
						},
						"dns_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "DNS configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nameservers": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "DNS IP List.",
									},
								},
							},
						},
					},
				},
			},

			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Scene statu Note: this field may return null, indicating that a valid value cannot be obtained.",
			},

			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Scene creation time.",
			},

			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Scene modification time.",
			},

			"app_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "App ID Note: this field may return null, indicating that a valid value cannot be obtained.",
			},

			"uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User ID Note: this field may return null, indicating that a valid value cannot be obtained.",
			},

			"sub_account_uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sub-user ID Note: this field may return null, indicating that a valid value cannot be obtained.",
			},

			// "notification_hooks": {
			// 	Type:        schema.TypeList,
			// 	Computed:    true,
			// 	Description: "Notification event callback Note: this field may return null, indicating that a valid value cannot be obtained.",
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"events": {
			// 				Type: schema.TypeSet,
			// 				Elem: &schema.Schema{
			// 					Type: schema.TypeString,
			// 				},
			// 				Optional:    true,
			// 				Description: "Notification event Note: this field may return null, indicating that a valid value cannot be obtained.",
			// 			},
			// 			"url": {
			// 				Type:        schema.TypeString,
			// 				Optional:    true,
			// 				Description: "Callback URL Note: this field may return null, indicating that a valid value cannot be obtained.",
			// 			},
			// 		},
			// 	},
			// },
		},
	}
}

func resourceTencentCloudPtsScenarioCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_scenario.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = pts.NewCreateScenarioRequest()
		requestUp  = pts.NewUpdateScenarioRequest()
		response   *pts.CreateScenarioResponse
		projectId  string
		scenarioId string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
		requestUp.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
		requestUp.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		request.ProjectId = helper.String(v.(string))
		requestUp.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
		requestUp.Description = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "load"); ok {
		load := pts.Load{}
		if LoadSpecMap, ok := helper.InterfaceToMap(dMap, "load_spec"); ok {
			loadSpec := pts.LoadSpec{}
			if ConcurrencyMap, ok := helper.InterfaceToMap(LoadSpecMap, "concurrency"); ok {
				concurrency := pts.Concurrency{}
				if v, ok := ConcurrencyMap["stages"]; ok {
					for _, item := range v.([]interface{}) {
						StagesMap := item.(map[string]interface{})
						stage := pts.Stage{}
						if v, ok := StagesMap["duration_seconds"]; ok {
							stage.DurationSeconds = helper.Int64(int64(v.(int)))
						}
						if v, ok := StagesMap["target_virtual_users"]; ok {
							stage.TargetVirtualUsers = helper.Int64(int64(v.(int)))
						}
						concurrency.Stages = append(concurrency.Stages, &stage)
					}
				}
				if v, ok := ConcurrencyMap["iteration_count"]; ok {
					concurrency.IterationCount = helper.Int64(int64(v.(int)))
				}
				if v, ok := ConcurrencyMap["max_requests_per_second"]; ok {
					concurrency.MaxRequestsPerSecond = helper.Int64(int64(v.(int)))
				}
				if v, ok := ConcurrencyMap["graceful_stop_seconds"]; ok {
					concurrency.GracefulStopSeconds = helper.Int64(int64(v.(int)))
				}
				loadSpec.Concurrency = &concurrency
			}
			if RequestsPerSecondMap, ok := helper.InterfaceToMap(LoadSpecMap, "requests_per_second"); ok {
				requestsPerSecond := pts.RequestsPerSecond{}
				if v, ok := RequestsPerSecondMap["max_requests_per_second"]; ok {
					requestsPerSecond.MaxRequestsPerSecond = helper.Int64(int64(v.(int)))
				}
				if v, ok := RequestsPerSecondMap["duration_seconds"]; ok {
					requestsPerSecond.DurationSeconds = helper.Int64(int64(v.(int)))
				}
				if v, ok := RequestsPerSecondMap["resources"]; ok {
					requestsPerSecond.Resources = helper.Int64(int64(v.(int)))
				}
				if v, ok := RequestsPerSecondMap["start_requests_per_second"]; ok {
					requestsPerSecond.StartRequestsPerSecond = helper.Int64(int64(v.(int)))
				}
				if v, ok := RequestsPerSecondMap["target_requests_per_second"]; ok {
					requestsPerSecond.TargetRequestsPerSecond = helper.Int64(int64(v.(int)))
				}
				if v, ok := RequestsPerSecondMap["graceful_stop_seconds"]; ok {
					requestsPerSecond.GracefulStopSeconds = helper.Int64(int64(v.(int)))
				}
				loadSpec.RequestsPerSecond = &requestsPerSecond
			}
			if ScriptOriginMap, ok := helper.InterfaceToMap(LoadSpecMap, "script_origin"); ok {
				scriptOrigin := pts.ScriptOrigin{}
				if v, ok := ScriptOriginMap["machine_number"]; ok {
					scriptOrigin.MachineNumber = helper.Int64(int64(v.(int)))
				}
				if v, ok := ScriptOriginMap["machine_specification"]; ok {
					scriptOrigin.MachineSpecification = helper.String(v.(string))
				}
				if v, ok := ScriptOriginMap["duration_seconds"]; ok {
					scriptOrigin.DurationSeconds = helper.Int64(int64(v.(int)))
				}
				loadSpec.ScriptOrigin = &scriptOrigin
			}
			load.LoadSpec = &loadSpec
		}
		if VpcLoadDistributionMap, ok := helper.InterfaceToMap(dMap, "vpc_load_distribution"); ok {
			vpcLoadDistribution := pts.VpcLoadDistribution{}
			if v, ok := VpcLoadDistributionMap["region_id"]; ok {
				vpcLoadDistribution.RegionId = helper.Int64(int64(v.(int)))
			}
			if v, ok := VpcLoadDistributionMap["region"]; ok {
				vpcLoadDistribution.Region = helper.String(v.(string))
			}
			if v, ok := VpcLoadDistributionMap["vpc_id"]; ok {
				vpcLoadDistribution.VpcId = helper.String(v.(string))
			}
			if v, ok := VpcLoadDistributionMap["subnet_ids"]; ok {
				subnetIdsSet := v.(*schema.Set).List()
				for i := range subnetIdsSet {
					subnetIds := subnetIdsSet[i].(string)
					vpcLoadDistribution.SubnetIds = append(vpcLoadDistribution.SubnetIds, &subnetIds)
				}
			}
			load.VpcLoadDistribution = &vpcLoadDistribution
		}
		if v, ok := dMap["geo_regions_load_distribution"]; ok {
			for _, item := range v.([]interface{}) {
				GeoRegionsLoadDistributionMap := item.(map[string]interface{})
				geoRegionsLoadItem := pts.GeoRegionsLoadItem{}
				if v, ok := GeoRegionsLoadDistributionMap["region_id"]; ok {
					geoRegionsLoadItem.RegionId = helper.Int64(int64(v.(int)))
				}
				if v, ok := GeoRegionsLoadDistributionMap["region"]; ok {
					geoRegionsLoadItem.Region = helper.String(v.(string))
				}
				if v, ok := GeoRegionsLoadDistributionMap["percentage"]; ok {
					geoRegionsLoadItem.Percentage = helper.Int64(int64(v.(int)))
				}
				load.GeoRegionsLoadDistribution = append(load.GeoRegionsLoadDistribution, &geoRegionsLoadItem)
			}
		}

		request.Load = &load
		requestUp.Load = &load
	}

	if v, ok := d.GetOk("datasets"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			testData := pts.TestData{}
			if v, ok := dMap["name"]; ok {
				testData.Name = helper.String(v.(string))
			}
			if v, ok := dMap["split"]; ok {
				testData.Split = helper.Bool(v.(bool))
			}
			if v, ok := dMap["header_in_file"]; ok {
				testData.HeaderInFile = helper.Bool(v.(bool))
			}
			if v, ok := dMap["header_columns"]; ok {
				headerColumnsSet := v.(*schema.Set).List()
				for i := range headerColumnsSet {
					headerColumns := headerColumnsSet[i].(string)
					testData.HeaderColumns = append(testData.HeaderColumns, &headerColumns)
				}
			}
			if v, ok := dMap["line_count"]; ok {
				testData.LineCount = helper.Int64(int64(v.(int)))
			}
			if v, ok := dMap["updated_at"]; ok {
				testData.UpdatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["size"]; ok {
				testData.Size = helper.Int64(int64(v.(int)))
			}
			if v, ok := dMap["head_lines"]; ok {
				headLinesSet := v.(*schema.Set).List()
				for i := range headLinesSet {
					headLines := headLinesSet[i].(string)
					testData.HeadLines = append(testData.HeadLines, &headLines)
				}
			}
			if v, ok := dMap["tail_lines"]; ok {
				tailLinesSet := v.(*schema.Set).List()
				for i := range tailLinesSet {
					tailLines := tailLinesSet[i].(string)
					testData.TailLines = append(testData.TailLines, &tailLines)
				}
			}
			if v, ok := dMap["type"]; ok {
				testData.Type = helper.String(v.(string))
			}
			if v, ok := dMap["file_id"]; ok {
				testData.FileId = helper.String(v.(string))
			}

			request.Datasets = append(request.Datasets, &testData)
			requestUp.Datasets = append(requestUp.Datasets, &testData)
		}
	}

	if v, ok := d.GetOk("extensions"); ok {
		extensionsSet := v.(*schema.Set).List()
		for i := range extensionsSet {
			extensions := extensionsSet[i].(string)
			request.Extensions = append(request.Extensions, &extensions)
			requestUp.Extensions = append(requestUp.Extensions, &extensions)
		}
	}

	if v, ok := d.GetOk("cron_id"); ok {
		request.CronId = helper.String(v.(string))
		requestUp.CronId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("test_scripts"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			scriptInfo := pts.ScriptInfo{}
			if v, ok := dMap["name"]; ok {
				scriptInfo.Name = helper.String(v.(string))
			}
			if v, ok := dMap["size"]; ok {
				scriptInfo.Size = helper.Int64(int64(v.(int)))
			}
			if v, ok := dMap["type"]; ok {
				scriptInfo.Type = helper.String(v.(string))
			}
			if v, ok := dMap["updated_at"]; ok {
				scriptInfo.UpdatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["encoded_content"]; ok {
				scriptInfo.EncodedContent = helper.String(tccommon.StringToBase64(v.(string)))
			}
			if v, ok := dMap["encoded_http_archive"]; ok {
				scriptInfo.EncodedHttpArchive = helper.String(tccommon.StringToBase64(v.(string)))
			}
			if v, ok := dMap["load_weight"]; ok {
				scriptInfo.LoadWeight = helper.Int64(int64(v.(int)))
			}

			request.TestScripts = append(request.TestScripts, &scriptInfo)
			requestUp.TestScripts = append(requestUp.TestScripts, &scriptInfo)
		}
	}

	if v, ok := d.GetOk("protocols"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			protocolInfo := pts.ProtocolInfo{}
			if v, ok := dMap["name"]; ok {
				protocolInfo.Name = helper.String(v.(string))
			}
			if v, ok := dMap["size"]; ok {
				protocolInfo.Size = helper.Int64(int64(v.(int)))
			}
			if v, ok := dMap["type"]; ok {
				protocolInfo.Type = helper.String(v.(string))
			}
			if v, ok := dMap["updated_at"]; ok {
				protocolInfo.UpdatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["file_id"]; ok {
				protocolInfo.FileId = helper.String(v.(string))
			}

			request.Protocols = append(request.Protocols, &protocolInfo)
			requestUp.Protocols = append(requestUp.Protocols, &protocolInfo)
		}
	}

	if v, ok := d.GetOk("request_files"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			fileInfo := pts.FileInfo{}
			if v, ok := dMap["name"]; ok {
				fileInfo.Name = helper.String(v.(string))
			}
			if v, ok := dMap["size"]; ok {
				fileInfo.Size = helper.Int64(int64(v.(int)))
			}
			if v, ok := dMap["type"]; ok {
				fileInfo.Type = helper.String(v.(string))
			}
			if v, ok := dMap["updated_at"]; ok {
				fileInfo.UpdatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["file_id"]; ok {
				fileInfo.FileId = helper.String(v.(string))
			}

			request.RequestFiles = append(request.RequestFiles, &fileInfo)
			requestUp.RequestFiles = append(requestUp.RequestFiles, &fileInfo)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "sla_policy"); ok {
		sLAPolicy := pts.SLAPolicy{}
		if v, ok := dMap["sla_rules"]; ok {
			for _, item := range v.([]interface{}) {
				SLARulesMap := item.(map[string]interface{})
				sLARule := pts.SLARule{}
				if v, ok := SLARulesMap["metric"]; ok {
					sLARule.Metric = helper.String(v.(string))
				}
				if v, ok := SLARulesMap["aggregation"]; ok {
					sLARule.Aggregation = helper.String(v.(string))
				}
				if v, ok := SLARulesMap["condition"]; ok {
					sLARule.Condition = helper.String(v.(string))
				}
				if v, ok := SLARulesMap["value"]; ok {
					sLARule.Value = helper.Float64(v.(float64))
				}
				if v, ok := SLARulesMap["label_filter"]; ok {
					for _, item := range v.([]interface{}) {
						LabelFilterMap := item.(map[string]interface{})
						sLALabel := pts.SLALabel{}
						if v, ok := LabelFilterMap["label_name"]; ok {
							sLALabel.LabelName = helper.String(v.(string))
						}
						if v, ok := LabelFilterMap["label_value"]; ok {
							sLALabel.LabelValue = helper.String(v.(string))
						}
						sLARule.LabelFilter = append(sLARule.LabelFilter, &sLALabel)
					}
				}
				if v, ok := SLARulesMap["abort_flag"]; ok {
					sLARule.AbortFlag = helper.Bool(v.(bool))
				}
				if v, ok := SLARulesMap["for"]; ok {
					sLARule.For = helper.String(v.(string))
				}
				sLAPolicy.SLARules = append(sLAPolicy.SLARules, &sLARule)
			}
		}
		if AlertChannelMap, ok := helper.InterfaceToMap(dMap, "alert_channel"); ok {
			alertChannel := pts.AlertChannel{}
			if v, ok := AlertChannelMap["notice_id"]; ok {
				alertChannel.NoticeId = helper.String(v.(string))
			}
			if v, ok := AlertChannelMap["amp_consumer_id"]; ok {
				alertChannel.AMPConsumerId = helper.String(v.(string))
			}
			sLAPolicy.AlertChannel = &alertChannel
		}

		request.SLAPolicy = &sLAPolicy
		requestUp.SLAPolicy = &sLAPolicy
	}

	if v, ok := d.GetOk("plugins"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			fileInfo := pts.FileInfo{}
			if v, ok := dMap["name"]; ok {
				fileInfo.Name = helper.String(v.(string))
			}
			if v, ok := dMap["size"]; ok {
				fileInfo.Size = helper.Int64(int64(v.(int)))
			}
			if v, ok := dMap["type"]; ok {
				fileInfo.Type = helper.String(v.(string))
			}
			if v, ok := dMap["updated_at"]; ok {
				fileInfo.UpdatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["file_id"]; ok {
				fileInfo.FileId = helper.String(v.(string))
			}

			request.Plugins = append(request.Plugins, &fileInfo)
			requestUp.Plugins = append(requestUp.Plugins, &fileInfo)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "domain_name_config"); ok {
		domainNameConfig := pts.DomainNameConfig{}
		if v, ok := dMap["host_aliases"]; ok {
			for _, item := range v.([]interface{}) {
				HostAliasesMap := item.(map[string]interface{})
				hostAlias := pts.HostAlias{}
				if v, ok := HostAliasesMap["host_names"]; ok {
					hostNamesSet := v.(*schema.Set).List()
					for i := range hostNamesSet {
						hostNames := hostNamesSet[i].(string)
						hostAlias.HostNames = append(hostAlias.HostNames, &hostNames)
					}
				}
				if v, ok := HostAliasesMap["ip"]; ok {
					hostAlias.IP = helper.String(v.(string))
				}
				domainNameConfig.HostAliases = append(domainNameConfig.HostAliases, &hostAlias)
			}
		}
		if DNSConfigMap, ok := helper.InterfaceToMap(dMap, "dns_config"); ok {
			dNSConfig := pts.DNSConfig{}
			if v, ok := DNSConfigMap["nameservers"]; ok {
				nameserversSet := v.(*schema.Set).List()
				for i := range nameserversSet {
					nameservers := nameserversSet[i].(string)
					dNSConfig.Nameservers = append(dNSConfig.Nameservers, &nameservers)
				}
			}
			domainNameConfig.DNSConfig = &dNSConfig
		}

		request.DomainNameConfig = &domainNameConfig
		requestUp.DomainNameConfig = &domainNameConfig
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePtsClient().CreateScenario(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create pts scenario failed, reason:%+v", logId, err)
		return err
	}

	scenarioId = *response.Response.ScenarioId

	requestUp.ScenarioId = &scenarioId
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePtsClient().UpdateScenario(requestUp)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create pts scenario failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(projectId + tccommon.FILED_SP + scenarioId)
	return resourceTencentCloudPtsScenarioRead(d, meta)
}

func resourceTencentCloudPtsScenarioRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_scenario.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := PtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	scenarioId := idSplit[1]

	scenario, err := service.DescribePtsScenario(ctx, projectId, scenarioId)

	if err != nil {
		return err
	}

	if scenario == nil {
		d.SetId("")
		return fmt.Errorf("resource `scenario` %s does not exist", scenarioId)
	}

	if scenario.Name != nil {
		_ = d.Set("name", scenario.Name)
	}

	if scenario.Type != nil {
		_ = d.Set("type", scenario.Type)
	}

	if scenario.ProjectId != nil {
		_ = d.Set("project_id", scenario.ProjectId)
	}

	if scenario.Description != nil {
		_ = d.Set("description", scenario.Description)
	}

	if scenario.Load != nil {
		loadMap := map[string]interface{}{}
		if scenario.Load.LoadSpec != nil {
			loadSpecMap := map[string]interface{}{}
			if scenario.Load.LoadSpec.Concurrency != nil {
				concurrencyMap := map[string]interface{}{}
				if scenario.Load.LoadSpec.Concurrency.Stages != nil {
					stagesList := []interface{}{}
					for _, stages := range scenario.Load.LoadSpec.Concurrency.Stages {
						stagesMap := map[string]interface{}{}
						if stages.DurationSeconds != nil {
							stagesMap["duration_seconds"] = stages.DurationSeconds
						}
						if stages.TargetVirtualUsers != nil {
							stagesMap["target_virtual_users"] = stages.TargetVirtualUsers
						}

						stagesList = append(stagesList, stagesMap)
					}
					concurrencyMap["stages"] = stagesList
				}
				if scenario.Load.LoadSpec.Concurrency.IterationCount != nil {
					concurrencyMap["iteration_count"] = scenario.Load.LoadSpec.Concurrency.IterationCount
				}
				if scenario.Load.LoadSpec.Concurrency.MaxRequestsPerSecond != nil {
					concurrencyMap["max_requests_per_second"] = scenario.Load.LoadSpec.Concurrency.MaxRequestsPerSecond
				}
				if scenario.Load.LoadSpec.Concurrency.GracefulStopSeconds != nil {
					concurrencyMap["graceful_stop_seconds"] = scenario.Load.LoadSpec.Concurrency.GracefulStopSeconds
				}

				loadSpecMap["concurrency"] = []interface{}{concurrencyMap}
			}
			if scenario.Load.LoadSpec.RequestsPerSecond != nil {
				requestsPerSecondMap := map[string]interface{}{}
				if scenario.Load.LoadSpec.RequestsPerSecond.MaxRequestsPerSecond != nil {
					requestsPerSecondMap["max_requests_per_second"] = scenario.Load.LoadSpec.RequestsPerSecond.MaxRequestsPerSecond
				}
				if scenario.Load.LoadSpec.RequestsPerSecond.DurationSeconds != nil {
					requestsPerSecondMap["duration_seconds"] = scenario.Load.LoadSpec.RequestsPerSecond.DurationSeconds
				}
				if scenario.Load.LoadSpec.RequestsPerSecond.Resources != nil {
					requestsPerSecondMap["resources"] = scenario.Load.LoadSpec.RequestsPerSecond.Resources
				}
				if scenario.Load.LoadSpec.RequestsPerSecond.StartRequestsPerSecond != nil {
					requestsPerSecondMap["start_requests_per_second"] = scenario.Load.LoadSpec.RequestsPerSecond.StartRequestsPerSecond
				}
				if scenario.Load.LoadSpec.RequestsPerSecond.TargetRequestsPerSecond != nil {
					requestsPerSecondMap["target_requests_per_second"] = scenario.Load.LoadSpec.RequestsPerSecond.TargetRequestsPerSecond
				}
				if scenario.Load.LoadSpec.RequestsPerSecond.GracefulStopSeconds != nil {
					requestsPerSecondMap["graceful_stop_seconds"] = scenario.Load.LoadSpec.RequestsPerSecond.GracefulStopSeconds
				}

				loadSpecMap["requests_per_second"] = []interface{}{requestsPerSecondMap}
			}
			if scenario.Load.LoadSpec.ScriptOrigin != nil {
				scriptOriginMap := map[string]interface{}{}
				if scenario.Load.LoadSpec.ScriptOrigin.MachineNumber != nil {
					scriptOriginMap["machine_number"] = scenario.Load.LoadSpec.ScriptOrigin.MachineNumber
				}
				if scenario.Load.LoadSpec.ScriptOrigin.MachineSpecification != nil {
					scriptOriginMap["machine_specification"] = scenario.Load.LoadSpec.ScriptOrigin.MachineSpecification
				}
				if scenario.Load.LoadSpec.ScriptOrigin.DurationSeconds != nil {
					scriptOriginMap["duration_seconds"] = scenario.Load.LoadSpec.ScriptOrigin.DurationSeconds
				}

				loadSpecMap["script_origin"] = []interface{}{scriptOriginMap}
			}

			loadMap["load_spec"] = []interface{}{loadSpecMap}
		}
		if scenario.Load.VpcLoadDistribution != nil {
			vpcLoadDistributionMap := map[string]interface{}{}
			if scenario.Load.VpcLoadDistribution.RegionId != nil {
				vpcLoadDistributionMap["region_id"] = scenario.Load.VpcLoadDistribution.RegionId
			}
			if scenario.Load.VpcLoadDistribution.Region != nil {
				vpcLoadDistributionMap["region"] = scenario.Load.VpcLoadDistribution.Region
			}
			if scenario.Load.VpcLoadDistribution.VpcId != nil {
				vpcLoadDistributionMap["vpc_id"] = scenario.Load.VpcLoadDistribution.VpcId
			}
			if scenario.Load.VpcLoadDistribution.SubnetIds != nil {
				vpcLoadDistributionMap["subnet_ids"] = scenario.Load.VpcLoadDistribution.SubnetIds
			}

			loadMap["vpc_load_distribution"] = []interface{}{vpcLoadDistributionMap}
		}
		if scenario.Load.GeoRegionsLoadDistribution != nil {
			geoRegionsLoadDistributionList := []interface{}{}
			for _, geoRegionsLoadDistribution := range scenario.Load.GeoRegionsLoadDistribution {
				geoRegionsLoadDistributionMap := map[string]interface{}{}
				if geoRegionsLoadDistribution.RegionId != nil {
					geoRegionsLoadDistributionMap["region_id"] = geoRegionsLoadDistribution.RegionId
				}
				if geoRegionsLoadDistribution.Region != nil {
					geoRegionsLoadDistributionMap["region"] = geoRegionsLoadDistribution.Region
				}
				if geoRegionsLoadDistribution.Percentage != nil {
					geoRegionsLoadDistributionMap["percentage"] = geoRegionsLoadDistribution.Percentage
				}

				geoRegionsLoadDistributionList = append(geoRegionsLoadDistributionList, geoRegionsLoadDistributionMap)
			}
			loadMap["geo_regions_load_distribution"] = geoRegionsLoadDistributionList
		}

		_ = d.Set("load", []interface{}{loadMap})
	}

	if scenario.Datasets != nil {
		datasetsList := []interface{}{}
		for _, datasets := range scenario.Datasets {
			datasetsMap := map[string]interface{}{}
			if datasets.Name != nil {
				datasetsMap["name"] = datasets.Name
			}
			if datasets.Split != nil {
				datasetsMap["split"] = datasets.Split
			}
			if datasets.HeaderInFile != nil {
				datasetsMap["header_in_file"] = datasets.HeaderInFile
			}
			if datasets.HeaderColumns != nil {
				datasetsMap["header_columns"] = datasets.HeaderColumns
			}
			if datasets.LineCount != nil {
				datasetsMap["line_count"] = datasets.LineCount
			}
			if datasets.UpdatedAt != nil {
				datasetsMap["updated_at"] = datasets.UpdatedAt
			}
			if datasets.Size != nil {
				datasetsMap["size"] = datasets.Size
			}
			if datasets.HeadLines != nil {
				datasetsMap["head_lines"] = datasets.HeadLines
			}
			if datasets.TailLines != nil {
				datasetsMap["tail_lines"] = datasets.TailLines
			}
			if datasets.Type != nil {
				datasetsMap["type"] = datasets.Type
			}
			if datasets.FileId != nil {
				datasetsMap["file_id"] = datasets.FileId
			}

			datasetsList = append(datasetsList, datasetsMap)
		}
		_ = d.Set("datasets", datasetsList)
	}

	if scenario.Extensions != nil {
		_ = d.Set("extensions", scenario.Extensions)
	}

	if scenario.CronId != nil {
		_ = d.Set("cron_id", scenario.CronId)
	}

	if scenario.TestScripts != nil {
		testScriptsList := []interface{}{}
		for _, testScripts := range scenario.TestScripts {
			testScriptsMap := map[string]interface{}{}
			if testScripts.Name != nil {
				testScriptsMap["name"] = testScripts.Name
			}
			if testScripts.Size != nil {
				testScriptsMap["size"] = testScripts.Size
			}
			if testScripts.Type != nil {
				testScriptsMap["type"] = testScripts.Type
			}
			if testScripts.UpdatedAt != nil {
				testScriptsMap["updated_at"] = testScripts.UpdatedAt
			}
			if testScripts.EncodedContent != nil {
				content, err := tccommon.Base64ToString(*testScripts.EncodedContent)
				if err != nil {
					return fmt.Errorf("`testScripts.EncodedContent` %s does not be decoded to string", *testScripts.EncodedContent)
				}
				testScriptsMap["encoded_content"] = content
			}
			if testScripts.EncodedHttpArchive != nil {
				archive, err := tccommon.Base64ToString(*testScripts.EncodedHttpArchive)
				if err != nil {
					return fmt.Errorf("`testScripts.EncodedHttpArchive` %s does not be decoded to string", *testScripts.EncodedHttpArchive)
				}
				testScriptsMap["encoded_http_archive"] = archive
			}
			if testScripts.LoadWeight != nil {
				testScriptsMap["load_weight"] = testScripts.LoadWeight
			}

			testScriptsList = append(testScriptsList, testScriptsMap)
		}
		_ = d.Set("test_scripts", testScriptsList)
	}

	if scenario.Protocols != nil {
		protocolsList := []interface{}{}
		for _, protocols := range scenario.Protocols {
			protocolsMap := map[string]interface{}{}
			if protocols.Name != nil {
				protocolsMap["name"] = protocols.Name
			}
			if protocols.Size != nil {
				protocolsMap["size"] = protocols.Size
			}
			if protocols.Type != nil {
				protocolsMap["type"] = protocols.Type
			}
			if protocols.UpdatedAt != nil {
				protocolsMap["updated_at"] = protocols.UpdatedAt
			}
			if protocols.FileId != nil {
				protocolsMap["file_id"] = protocols.FileId
			}

			protocolsList = append(protocolsList, protocolsMap)
		}
		_ = d.Set("protocols", protocolsList)
	}

	if scenario.RequestFiles != nil {
		requestFilesList := []interface{}{}
		for _, requestFiles := range scenario.RequestFiles {
			requestFilesMap := map[string]interface{}{}
			if requestFiles.Name != nil {
				requestFilesMap["name"] = requestFiles.Name
			}
			if requestFiles.Size != nil {
				requestFilesMap["size"] = requestFiles.Size
			}
			if requestFiles.Type != nil {
				requestFilesMap["type"] = requestFiles.Type
			}
			if requestFiles.UpdatedAt != nil {
				requestFilesMap["updated_at"] = requestFiles.UpdatedAt
			}
			if requestFiles.FileId != nil {
				requestFilesMap["file_id"] = requestFiles.FileId
			}

			requestFilesList = append(requestFilesList, requestFilesMap)
		}
		_ = d.Set("request_files", requestFilesList)
	}

	if scenario.SLAPolicy != nil {
		sLAPolicyMap := map[string]interface{}{}
		if scenario.SLAPolicy.SLARules != nil {
			sLARulesList := []interface{}{}
			for _, sLARules := range scenario.SLAPolicy.SLARules {
				sLARulesMap := map[string]interface{}{}
				if sLARules.Metric != nil {
					sLARulesMap["metric"] = sLARules.Metric
				}
				if sLARules.Aggregation != nil {
					sLARulesMap["aggregation"] = sLARules.Aggregation
				}
				if sLARules.Condition != nil {
					sLARulesMap["condition"] = sLARules.Condition
				}
				if sLARules.Value != nil {
					sLARulesMap["value"] = sLARules.Value
				}
				if sLARules.LabelFilter != nil {
					labelFilterList := []interface{}{}
					for _, labelFilter := range sLARules.LabelFilter {
						labelFilterMap := map[string]interface{}{}
						if labelFilter.LabelName != nil {
							labelFilterMap["label_name"] = labelFilter.LabelName
						}
						if labelFilter.LabelValue != nil {
							labelFilterMap["label_value"] = labelFilter.LabelValue
						}

						labelFilterList = append(labelFilterList, labelFilterMap)
					}
					sLARulesMap["label_filter"] = labelFilterList
				}
				if sLARules.AbortFlag != nil {
					sLARulesMap["abort_flag"] = sLARules.AbortFlag
				}
				if sLARules.For != nil {
					sLARulesMap["for"] = sLARules.For
				}

				sLARulesList = append(sLARulesList, sLARulesMap)
			}
			sLAPolicyMap["sla_rules"] = sLARulesList
		}
		if scenario.SLAPolicy.AlertChannel != nil {
			alertChannelMap := map[string]interface{}{}
			if scenario.SLAPolicy.AlertChannel.NoticeId != nil {
				alertChannelMap["notice_id"] = scenario.SLAPolicy.AlertChannel.NoticeId
			}
			if scenario.SLAPolicy.AlertChannel.AMPConsumerId != nil {
				alertChannelMap["amp_consumer_id"] = scenario.SLAPolicy.AlertChannel.AMPConsumerId
			}

			sLAPolicyMap["alert_channel"] = []interface{}{alertChannelMap}
		}

		_ = d.Set("sla_policy", []interface{}{sLAPolicyMap})
	}

	if scenario.Plugins != nil {
		pluginsList := []interface{}{}
		for _, plugins := range scenario.Plugins {
			pluginsMap := map[string]interface{}{}
			if plugins.Name != nil {
				pluginsMap["name"] = plugins.Name
			}
			if plugins.Size != nil {
				pluginsMap["size"] = plugins.Size
			}
			if plugins.Type != nil {
				pluginsMap["type"] = plugins.Type
			}
			if plugins.UpdatedAt != nil {
				pluginsMap["updated_at"] = plugins.UpdatedAt
			}
			if plugins.FileId != nil {
				pluginsMap["file_id"] = plugins.FileId
			}

			pluginsList = append(pluginsList, pluginsMap)
		}
		_ = d.Set("plugins", pluginsList)
	}

	if scenario.DomainNameConfig != nil {
		domainNameConfigMap := map[string]interface{}{}
		if scenario.DomainNameConfig.HostAliases != nil {
			hostAliasesList := []interface{}{}
			for _, hostAliases := range scenario.DomainNameConfig.HostAliases {
				hostAliasesMap := map[string]interface{}{}
				if hostAliases.HostNames != nil {
					hostAliasesMap["host_names"] = hostAliases.HostNames
				}
				if hostAliases.IP != nil {
					hostAliasesMap["ip"] = hostAliases.IP
				}

				hostAliasesList = append(hostAliasesList, hostAliasesMap)
			}
			domainNameConfigMap["host_aliases"] = hostAliasesList
		}
		if scenario.DomainNameConfig.DNSConfig != nil {
			dNSConfigMap := map[string]interface{}{}
			if scenario.DomainNameConfig.DNSConfig.Nameservers != nil {
				dNSConfigMap["nameservers"] = scenario.DomainNameConfig.DNSConfig.Nameservers
			}

			domainNameConfigMap["dns_config"] = []interface{}{dNSConfigMap}
		}

		_ = d.Set("domain_name_config", []interface{}{domainNameConfigMap})
	}

	if scenario.Status != nil {
		_ = d.Set("status", scenario.Status)
	}

	if scenario.CreatedAt != nil {
		_ = d.Set("created_at", scenario.CreatedAt)
	}

	if scenario.UpdatedAt != nil {
		_ = d.Set("updated_at", scenario.UpdatedAt)
	}

	if scenario.AppId != nil {
		_ = d.Set("app_id", scenario.AppId)
	}

	if scenario.Uin != nil {
		_ = d.Set("uin", scenario.Uin)
	}

	if scenario.SubAccountUin != nil {
		_ = d.Set("sub_account_uin", scenario.SubAccountUin)
	}

	// if scenario.NotificationHooks != nil {
	// 	notificationHooksList := []interface{}{}
	// 	for _, notificationHooks := range scenario.NotificationHooks {
	// 		notificationHooksMap := map[string]interface{}{}
	// 		if notificationHooks.Events != nil {
	// 			notificationHooksMap["events"] = notificationHooks.Events
	// 		}
	// 		if notificationHooks.URL != nil {
	// 			notificationHooksMap["url"] = notificationHooks.URL
	// 		}

	// 		notificationHooksList = append(notificationHooksList, notificationHooksMap)
	// 	}
	// 	err = d.Set("notification_hooks", notificationHooksList)
	// 	if err != nil {
	// 		return fmt.Errorf("set error")
	// 	}
	// }

	return nil
}

func resourceTencentCloudPtsScenarioUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_scenario.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := pts.NewUpdateScenarioRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	scenarioId := idSplit[1]

	request.ProjectId = &projectId
	request.ScenarioId = &scenarioId

	if d.HasChange("project_id") {
		return fmt.Errorf("`project_id` do not support change now.")
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "load"); ok {
		load := pts.Load{}
		if LoadSpecMap, ok := helper.InterfaceToMap(dMap, "load_spec"); ok {
			loadSpec := pts.LoadSpec{}
			if ConcurrencyMap, ok := helper.InterfaceToMap(LoadSpecMap, "concurrency"); ok {
				concurrency := pts.Concurrency{}
				if v, ok := ConcurrencyMap["stages"]; ok {
					for _, item := range v.([]interface{}) {
						StagesMap := item.(map[string]interface{})
						stage := pts.Stage{}
						if v, ok := StagesMap["duration_seconds"]; ok {
							stage.DurationSeconds = helper.Int64(int64(v.(int)))
						}
						if v, ok := StagesMap["target_virtual_users"]; ok {
							stage.TargetVirtualUsers = helper.Int64(int64(v.(int)))
						}
						concurrency.Stages = append(concurrency.Stages, &stage)
					}
				}
				if v, ok := ConcurrencyMap["iteration_count"]; ok {
					concurrency.IterationCount = helper.Int64(int64(v.(int)))
				}
				if v, ok := ConcurrencyMap["max_requests_per_second"]; ok {
					concurrency.MaxRequestsPerSecond = helper.Int64(int64(v.(int)))
				}
				if v, ok := ConcurrencyMap["graceful_stop_seconds"]; ok {
					concurrency.GracefulStopSeconds = helper.Int64(int64(v.(int)))
				}
				loadSpec.Concurrency = &concurrency
			}
			if RequestsPerSecondMap, ok := helper.InterfaceToMap(LoadSpecMap, "requests_per_second"); ok {
				requestsPerSecond := pts.RequestsPerSecond{}
				if v, ok := RequestsPerSecondMap["max_requests_per_second"]; ok {
					requestsPerSecond.MaxRequestsPerSecond = helper.Int64(int64(v.(int)))
				}
				if v, ok := RequestsPerSecondMap["duration_seconds"]; ok {
					requestsPerSecond.DurationSeconds = helper.Int64(int64(v.(int)))
				}
				if v, ok := RequestsPerSecondMap["resources"]; ok {
					requestsPerSecond.Resources = helper.Int64(int64(v.(int)))
				}
				if v, ok := RequestsPerSecondMap["start_requests_per_second"]; ok {
					requestsPerSecond.StartRequestsPerSecond = helper.Int64(int64(v.(int)))
				}
				if v, ok := RequestsPerSecondMap["target_requests_per_second"]; ok {
					requestsPerSecond.TargetRequestsPerSecond = helper.Int64(int64(v.(int)))
				}
				if v, ok := RequestsPerSecondMap["graceful_stop_seconds"]; ok {
					requestsPerSecond.GracefulStopSeconds = helper.Int64(int64(v.(int)))
				}
				loadSpec.RequestsPerSecond = &requestsPerSecond
			}
			if ScriptOriginMap, ok := helper.InterfaceToMap(LoadSpecMap, "script_origin"); ok {
				scriptOrigin := pts.ScriptOrigin{}
				if v, ok := ScriptOriginMap["machine_number"]; ok {
					scriptOrigin.MachineNumber = helper.Int64(int64(v.(int)))
				}
				if v, ok := ScriptOriginMap["machine_specification"]; ok {
					scriptOrigin.MachineSpecification = helper.String(v.(string))
				}
				if v, ok := ScriptOriginMap["duration_seconds"]; ok {
					scriptOrigin.DurationSeconds = helper.Int64(int64(v.(int)))
				}
				loadSpec.ScriptOrigin = &scriptOrigin
			}
			load.LoadSpec = &loadSpec
		}
		if VpcLoadDistributionMap, ok := helper.InterfaceToMap(dMap, "vpc_load_distribution"); ok {
			vpcLoadDistribution := pts.VpcLoadDistribution{}
			if v, ok := VpcLoadDistributionMap["region_id"]; ok {
				vpcLoadDistribution.RegionId = helper.Int64(int64(v.(int)))
			}
			if v, ok := VpcLoadDistributionMap["region"]; ok {
				vpcLoadDistribution.Region = helper.String(v.(string))
			}
			if v, ok := VpcLoadDistributionMap["vpc_id"]; ok {
				vpcLoadDistribution.VpcId = helper.String(v.(string))
			}
			if v, ok := VpcLoadDistributionMap["subnet_ids"]; ok {
				subnetIdsSet := v.(*schema.Set).List()
				for i := range subnetIdsSet {
					subnetIds := subnetIdsSet[i].(string)
					vpcLoadDistribution.SubnetIds = append(vpcLoadDistribution.SubnetIds, &subnetIds)
				}
			}
			load.VpcLoadDistribution = &vpcLoadDistribution
		}
		if v, ok := dMap["geo_regions_load_distribution"]; ok {
			for _, item := range v.([]interface{}) {
				GeoRegionsLoadDistributionMap := item.(map[string]interface{})
				geoRegionsLoadItem := pts.GeoRegionsLoadItem{}
				if v, ok := GeoRegionsLoadDistributionMap["region_id"]; ok {
					geoRegionsLoadItem.RegionId = helper.Int64(int64(v.(int)))
				}
				if v, ok := GeoRegionsLoadDistributionMap["region"]; ok {
					geoRegionsLoadItem.Region = helper.String(v.(string))
				}
				if v, ok := GeoRegionsLoadDistributionMap["percentage"]; ok {
					geoRegionsLoadItem.Percentage = helper.Int64(int64(v.(int)))
				}
				load.GeoRegionsLoadDistribution = append(load.GeoRegionsLoadDistribution, &geoRegionsLoadItem)
			}
		}

		request.Load = &load
	}

	if v, ok := d.GetOk("datasets"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			testData := pts.TestData{}
			if v, ok := dMap["name"]; ok {
				testData.Name = helper.String(v.(string))
			}
			if v, ok := dMap["split"]; ok {
				testData.Split = helper.Bool(v.(bool))
			}
			if v, ok := dMap["header_in_file"]; ok {
				testData.HeaderInFile = helper.Bool(v.(bool))
			}
			if v, ok := dMap["header_columns"]; ok {
				headerColumnsSet := v.(*schema.Set).List()
				for i := range headerColumnsSet {
					headerColumns := headerColumnsSet[i].(string)
					testData.HeaderColumns = append(testData.HeaderColumns, &headerColumns)
				}
			}
			if v, ok := dMap["line_count"]; ok {
				testData.LineCount = helper.Int64(int64(v.(int)))
			}
			if v, ok := dMap["updated_at"]; ok {
				testData.UpdatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["size"]; ok {
				testData.Size = helper.Int64(int64(v.(int)))
			}
			if v, ok := dMap["head_lines"]; ok {
				headLinesSet := v.(*schema.Set).List()
				for i := range headLinesSet {
					headLines := headLinesSet[i].(string)
					testData.HeadLines = append(testData.HeadLines, &headLines)
				}
			}
			if v, ok := dMap["tail_lines"]; ok {
				tailLinesSet := v.(*schema.Set).List()
				for i := range tailLinesSet {
					tailLines := tailLinesSet[i].(string)
					testData.TailLines = append(testData.TailLines, &tailLines)
				}
			}
			if v, ok := dMap["type"]; ok {
				testData.Type = helper.String(v.(string))
			}
			if v, ok := dMap["file_id"]; ok {
				testData.FileId = helper.String(v.(string))
			}

			request.Datasets = append(request.Datasets, &testData)
		}
	}

	if v, ok := d.GetOk("extensions"); ok {
		extensionsSet := v.(*schema.Set).List()
		for i := range extensionsSet {
			extensions := extensionsSet[i].(string)
			request.Extensions = append(request.Extensions, &extensions)
		}
	}

	if v, ok := d.GetOk("cron_id"); ok {
		request.CronId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("test_scripts"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			scriptInfo := pts.ScriptInfo{}
			if v, ok := dMap["name"]; ok {
				scriptInfo.Name = helper.String(v.(string))
			}
			if v, ok := dMap["size"]; ok {
				scriptInfo.Size = helper.Int64(int64(v.(int)))
			}
			if v, ok := dMap["type"]; ok {
				scriptInfo.Type = helper.String(v.(string))
			}
			if v, ok := dMap["updated_at"]; ok {
				scriptInfo.UpdatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["encoded_content"]; ok {
				scriptInfo.EncodedContent = helper.String(tccommon.StringToBase64(v.(string)))
			}
			if v, ok := dMap["encoded_http_archive"]; ok {
				scriptInfo.EncodedHttpArchive = helper.String(tccommon.StringToBase64(v.(string)))
			}
			if v, ok := dMap["load_weight"]; ok {
				scriptInfo.LoadWeight = helper.Int64(int64(v.(int)))
			}

			request.TestScripts = append(request.TestScripts, &scriptInfo)
		}
	}

	if v, ok := d.GetOk("protocols"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			protocolInfo := pts.ProtocolInfo{}
			if v, ok := dMap["name"]; ok {
				protocolInfo.Name = helper.String(v.(string))
			}
			if v, ok := dMap["size"]; ok {
				protocolInfo.Size = helper.Int64(int64(v.(int)))
			}
			if v, ok := dMap["type"]; ok {
				protocolInfo.Type = helper.String(v.(string))
			}
			if v, ok := dMap["updated_at"]; ok {
				protocolInfo.UpdatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["file_id"]; ok {
				protocolInfo.FileId = helper.String(v.(string))
			}

			request.Protocols = append(request.Protocols, &protocolInfo)
		}
	}

	if v, ok := d.GetOk("request_files"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			fileInfo := pts.FileInfo{}
			if v, ok := dMap["name"]; ok {
				fileInfo.Name = helper.String(v.(string))
			}
			if v, ok := dMap["size"]; ok {
				fileInfo.Size = helper.Int64(int64(v.(int)))
			}
			if v, ok := dMap["type"]; ok {
				fileInfo.Type = helper.String(v.(string))
			}
			if v, ok := dMap["updated_at"]; ok {
				fileInfo.UpdatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["file_id"]; ok {
				fileInfo.FileId = helper.String(v.(string))
			}

			request.RequestFiles = append(request.RequestFiles, &fileInfo)
		}
	}

	if _, ok := d.GetOk("sla_policy"); ok {
		if dMap, ok := helper.InterfacesHeadMap(d, "sla_policy"); ok {
			sLAPolicy := pts.SLAPolicy{}
			if v, ok := dMap["sla_rules"]; ok {
				for _, item := range v.([]interface{}) {
					SLARulesMap := item.(map[string]interface{})
					sLARule := pts.SLARule{}
					if v, ok := SLARulesMap["metric"]; ok {
						sLARule.Metric = helper.String(v.(string))
					}
					if v, ok := SLARulesMap["aggregation"]; ok {
						sLARule.Aggregation = helper.String(v.(string))
					}
					if v, ok := SLARulesMap["condition"]; ok {
						sLARule.Condition = helper.String(v.(string))
					}
					if v, ok := SLARulesMap["value"]; ok {
						sLARule.Value = helper.Float64(v.(float64))
					}
					if v, ok := SLARulesMap["label_filter"]; ok {
						for _, item := range v.([]interface{}) {
							LabelFilterMap := item.(map[string]interface{})
							sLALabel := pts.SLALabel{}
							if v, ok := LabelFilterMap["label_name"]; ok {
								sLALabel.LabelName = helper.String(v.(string))
							}
							if v, ok := LabelFilterMap["label_value"]; ok {
								sLALabel.LabelValue = helper.String(v.(string))
							}
							sLARule.LabelFilter = append(sLARule.LabelFilter, &sLALabel)
						}
					}
					if v, ok := SLARulesMap["abort_flag"]; ok {
						sLARule.AbortFlag = helper.Bool(v.(bool))
					}
					if v, ok := SLARulesMap["for"]; ok {
						sLARule.For = helper.String(v.(string))
					}
					sLAPolicy.SLARules = append(sLAPolicy.SLARules, &sLARule)
				}
			}
			if AlertChannelMap, ok := helper.InterfaceToMap(dMap, "alert_channel"); ok {
				alertChannel := pts.AlertChannel{}
				if v, ok := AlertChannelMap["notice_id"]; ok {
					alertChannel.NoticeId = helper.String(v.(string))
				}
				if v, ok := AlertChannelMap["amp_consumer_id"]; ok {
					alertChannel.AMPConsumerId = helper.String(v.(string))
				}
				sLAPolicy.AlertChannel = &alertChannel
			}

			request.SLAPolicy = &sLAPolicy
		}
	}

	if v, ok := d.GetOk("plugins"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			fileInfo := pts.FileInfo{}
			if v, ok := dMap["name"]; ok {
				fileInfo.Name = helper.String(v.(string))
			}
			if v, ok := dMap["size"]; ok {
				fileInfo.Size = helper.Int64(int64(v.(int)))
			}
			if v, ok := dMap["type"]; ok {
				fileInfo.Type = helper.String(v.(string))
			}
			if v, ok := dMap["updated_at"]; ok {
				fileInfo.UpdatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["file_id"]; ok {
				fileInfo.FileId = helper.String(v.(string))
			}

			request.Plugins = append(request.Plugins, &fileInfo)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "domain_name_config"); ok {
		domainNameConfig := pts.DomainNameConfig{}
		if v, ok := dMap["host_aliases"]; ok {
			for _, item := range v.([]interface{}) {
				HostAliasesMap := item.(map[string]interface{})
				hostAlias := pts.HostAlias{}
				if v, ok := HostAliasesMap["host_names"]; ok {
					hostNamesSet := v.(*schema.Set).List()
					for i := range hostNamesSet {
						hostNames := hostNamesSet[i].(string)
						hostAlias.HostNames = append(hostAlias.HostNames, &hostNames)
					}
				}
				if v, ok := HostAliasesMap["ip"]; ok {
					hostAlias.IP = helper.String(v.(string))
				}
				domainNameConfig.HostAliases = append(domainNameConfig.HostAliases, &hostAlias)
			}
		}
		if DNSConfigMap, ok := helper.InterfaceToMap(dMap, "dns_config"); ok {
			dNSConfig := pts.DNSConfig{}
			if v, ok := DNSConfigMap["nameservers"]; ok {
				nameserversSet := v.(*schema.Set).List()
				for i := range nameserversSet {
					nameservers := nameserversSet[i].(string)
					dNSConfig.Nameservers = append(dNSConfig.Nameservers, &nameservers)
				}
			}
			domainNameConfig.DNSConfig = &dNSConfig
		}

		request.DomainNameConfig = &domainNameConfig
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePtsClient().UpdateScenario(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create pts scenario failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPtsScenarioRead(d, meta)
}

func resourceTencentCloudPtsScenarioDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_scenario.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := PtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	scenarioId := idSplit[1]

	if err := service.DeletePtsScenarioById(ctx, projectId, scenarioId); err != nil {
		return err
	}

	return nil
}
