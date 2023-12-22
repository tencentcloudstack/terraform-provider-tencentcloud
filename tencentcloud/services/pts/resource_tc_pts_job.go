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

func ResourceTencentCloudPtsJob() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudPtsJobRead,
		Create: resourceTencentCloudPtsJobCreate,
		Update: resourceTencentCloudPtsJobUpdate,
		Delete: resourceTencentCloudPtsJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"scenario_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Pts scenario id.",
			},

			"job_owner": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Job owner.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID.",
			},

			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to debug.",
			},

			"note": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Note.",
			},

			"load": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Pressure configuration of job.",
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
				Computed:    true,
				Description: "Dataset file for the job.",
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

			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The running status of the task; `0`: JobUnknown, `1`: JobCreated, `2`: JobPending, `3`: JobPreparing, `4`: JobSelectClustering, `5`: JobCreateTasking, `6`: JobSyncTasking, `11`: JobRunning, `12`: JobFinished, `13`: JobPrepareException, `14`: JobFinishException, `15`: JobAborting, `16`: JobAborted, `17`: JobAbortException, `18`: JobDeleted, `19`: JobSelectClusterException, `20`: JobCreateTaskException, `21`: JobSyncTaskException.",
			},

			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Start time of the job.",
			},

			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "End time of the job.",
			},

			"max_virtual_user_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum number of VU for the job.",
			},

			"error_rate": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Percentage of error rate.",
			},

			"duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Job duration.",
			},

			"max_requests_per_second": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum requests per second.",
			},

			"request_total": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Total number of requests.",
			},

			"requests_per_second": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Average number of requests per second.",
			},

			"response_time_average": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Average response time.",
			},

			"response_time_p99": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "99th percentile response time.",
			},

			"response_time_p95": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "95th percentile response time.",
			},

			"response_time_p90": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "90th percentile response time.",
			},

			"response_time_max": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Maximum response time.",
			},

			"response_time_min": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Minimum response time.",
			},

			// "load_source_infos": {
			// 	Type:        schema.TypeList,
			// 	Computed:    true,
			// 	Description: "Host message of generating voltage.",
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"ip": {
			// 				Type:        schema.TypeString,
			// 				Optional:    true,
			// 				Description: "IP of host.",
			// 			},
			// 			"pod_name": {
			// 				Type:        schema.TypeString,
			// 				Optional:    true,
			// 				Description: "Pod of host.",
			// 			},
			// 			"region": {
			// 				Type:        schema.TypeString,
			// 				Optional:    true,
			// 				Description: "Region to which it belongs.",
			// 			},
			// 		},
			// 	},
			// },

			"test_scripts": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Test script information.",
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
				Computed:    true,
				Description: "Protocol script information.",
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
				Computed:    true,
				Description: "Request file information.",
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

			"plugins": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Expansion package file information.",
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

			"cron_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Scheduled job ID.",
			},

			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Scene Type.",
			},

			"domain_name_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Domain name binding configuration.",
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

			"abort_reason": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Cause of interruption.",
			},

			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job Id.",
			},

			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the job.",
			},

			// "notification_hooks": {
			// 	Type:        schema.TypeList,
			// 	Computed:    true,
			// 	Description: "Notification event callback.",
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"events": {
			// 				Type: schema.TypeSet,
			// 				Elem: &schema.Schema{
			// 					Type: schema.TypeString,
			// 				},
			// 				Optional:    true,
			// 				Description: "Notification event.",
			// 			},
			// 			"url": {
			// 				Type:        schema.TypeString,
			// 				Optional:    true,
			// 				Description: "Callback URL.",
			// 			},
			// 		},
			// 	},
			// },
		},
	}
}

func resourceTencentCloudPtsJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_job.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = pts.NewStartJobRequest()
		response   *pts.StartJobResponse
		projectId  string
		jobId      string
		scenarioId string
	)

	if v, ok := d.GetOk("scenario_id"); ok {
		scenarioId = v.(string)
		request.ScenarioId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("job_owner"); ok {
		request.JobOwner = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		request.ProjectId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("debug"); v != nil {
		request.Debug = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("note"); ok {
		request.Note = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePtsClient().StartJob(request)
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
		log.Printf("[CRITAL]%s create pts job failed, reason:%+v", logId, err)
		return err
	}

	jobId = *response.Response.JobId

	d.SetId(projectId + tccommon.FILED_SP + scenarioId + tccommon.FILED_SP + jobId)
	return resourceTencentCloudPtsJobRead(d, meta)
}

func resourceTencentCloudPtsJobRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_job.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := PtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	scenarioId := idSplit[1]
	jobId := idSplit[2]

	job, err := service.DescribePtsJob(ctx, projectId, scenarioId, jobId)

	if err != nil {
		return err
	}

	if job == nil {
		d.SetId("")
		return fmt.Errorf("resource `job` %s does not exist", jobId)
	}

	_ = d.Set("job_id", jobId)

	if job.ScenarioId != nil {
		_ = d.Set("scenario_id", job.ScenarioId)
	}

	if job.JobOwner != nil {
		_ = d.Set("job_owner", job.JobOwner)
	}

	if job.ProjectId != nil {
		_ = d.Set("project_id", job.ProjectId)
	}

	if job.Debug != nil {
		_ = d.Set("debug", job.Debug)
	}

	if job.Note != nil {
		_ = d.Set("note", job.Note)
	}

	if job.Load != nil {
		loadMap := map[string]interface{}{}
		if job.Load.LoadSpec != nil {
			loadSpecMap := map[string]interface{}{}
			if job.Load.LoadSpec.Concurrency != nil {
				concurrencyMap := map[string]interface{}{}
				if job.Load.LoadSpec.Concurrency.Stages != nil {
					stagesList := []interface{}{}
					for _, stages := range job.Load.LoadSpec.Concurrency.Stages {
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
				if job.Load.LoadSpec.Concurrency.IterationCount != nil {
					concurrencyMap["iteration_count"] = job.Load.LoadSpec.Concurrency.IterationCount
				}
				if job.Load.LoadSpec.Concurrency.MaxRequestsPerSecond != nil {
					concurrencyMap["max_requests_per_second"] = job.Load.LoadSpec.Concurrency.MaxRequestsPerSecond
				}
				if job.Load.LoadSpec.Concurrency.GracefulStopSeconds != nil {
					concurrencyMap["graceful_stop_seconds"] = job.Load.LoadSpec.Concurrency.GracefulStopSeconds
				}

				loadSpecMap["concurrency"] = []interface{}{concurrencyMap}
			}
			if job.Load.LoadSpec.RequestsPerSecond != nil {
				requestsPerSecondMap := map[string]interface{}{}
				if job.Load.LoadSpec.RequestsPerSecond.MaxRequestsPerSecond != nil {
					requestsPerSecondMap["max_requests_per_second"] = job.Load.LoadSpec.RequestsPerSecond.MaxRequestsPerSecond
				}
				if job.Load.LoadSpec.RequestsPerSecond.DurationSeconds != nil {
					requestsPerSecondMap["duration_seconds"] = job.Load.LoadSpec.RequestsPerSecond.DurationSeconds
				}
				if job.Load.LoadSpec.RequestsPerSecond.Resources != nil {
					requestsPerSecondMap["resources"] = job.Load.LoadSpec.RequestsPerSecond.Resources
				}
				if job.Load.LoadSpec.RequestsPerSecond.StartRequestsPerSecond != nil {
					requestsPerSecondMap["start_requests_per_second"] = job.Load.LoadSpec.RequestsPerSecond.StartRequestsPerSecond
				}
				if job.Load.LoadSpec.RequestsPerSecond.TargetRequestsPerSecond != nil {
					requestsPerSecondMap["target_requests_per_second"] = job.Load.LoadSpec.RequestsPerSecond.TargetRequestsPerSecond
				}
				if job.Load.LoadSpec.RequestsPerSecond.GracefulStopSeconds != nil {
					requestsPerSecondMap["graceful_stop_seconds"] = job.Load.LoadSpec.RequestsPerSecond.GracefulStopSeconds
				}

				loadSpecMap["requests_per_second"] = []interface{}{requestsPerSecondMap}
			}
			if job.Load.LoadSpec.ScriptOrigin != nil {
				scriptOriginMap := map[string]interface{}{}
				if job.Load.LoadSpec.ScriptOrigin.MachineNumber != nil {
					scriptOriginMap["machine_number"] = job.Load.LoadSpec.ScriptOrigin.MachineNumber
				}
				if job.Load.LoadSpec.ScriptOrigin.MachineSpecification != nil {
					scriptOriginMap["machine_specification"] = job.Load.LoadSpec.ScriptOrigin.MachineSpecification
				}
				if job.Load.LoadSpec.ScriptOrigin.DurationSeconds != nil {
					scriptOriginMap["duration_seconds"] = job.Load.LoadSpec.ScriptOrigin.DurationSeconds
				}

				loadSpecMap["script_origin"] = []interface{}{scriptOriginMap}
			}

			loadMap["load_spec"] = []interface{}{loadSpecMap}
		}
		if job.Load.VpcLoadDistribution != nil {
			vpcLoadDistributionMap := map[string]interface{}{}
			if job.Load.VpcLoadDistribution.RegionId != nil {
				vpcLoadDistributionMap["region_id"] = job.Load.VpcLoadDistribution.RegionId
			}
			if job.Load.VpcLoadDistribution.Region != nil {
				vpcLoadDistributionMap["region"] = job.Load.VpcLoadDistribution.Region
			}
			if job.Load.VpcLoadDistribution.VpcId != nil {
				vpcLoadDistributionMap["vpc_id"] = job.Load.VpcLoadDistribution.VpcId
			}
			if job.Load.VpcLoadDistribution.SubnetIds != nil {
				vpcLoadDistributionMap["subnet_ids"] = job.Load.VpcLoadDistribution.SubnetIds
			}

			loadMap["vpc_load_distribution"] = []interface{}{vpcLoadDistributionMap}
		}
		if job.Load.GeoRegionsLoadDistribution != nil {
			geoRegionsLoadDistributionList := []interface{}{}
			for _, geoRegionsLoadDistribution := range job.Load.GeoRegionsLoadDistribution {
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

	if job.Datasets != nil {
		datasetsList := []interface{}{}
		for _, datasets := range job.Datasets {
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

	if job.Status != nil {
		_ = d.Set("status", job.Status)
	}

	if job.StartTime != nil {
		_ = d.Set("start_time", job.StartTime)
	}

	if job.EndTime != nil {
		_ = d.Set("end_time", job.EndTime)
	}

	if job.MaxVirtualUserCount != nil {
		_ = d.Set("max_virtual_user_count", job.MaxVirtualUserCount)
	}

	if job.ErrorRate != nil {
		_ = d.Set("error_rate", job.ErrorRate)
	}

	if job.Duration != nil {
		_ = d.Set("duration", job.Duration)
	}

	// if job.MaxRequestsPerSecond != nil {
	// 	_ = d.Set("max_requests_per_second", job.MaxRequestsPerSecond)
	// }

	if job.RequestTotal != nil {
		_ = d.Set("request_total", job.RequestTotal)
	}

	if job.RequestsPerSecond != nil {
		_ = d.Set("requests_per_second", job.RequestsPerSecond)
	}

	if job.ResponseTimeAverage != nil {
		_ = d.Set("response_time_average", job.ResponseTimeAverage)
	}

	if job.ResponseTimeP99 != nil {
		_ = d.Set("response_time_p99", job.ResponseTimeP99)
	}

	if job.ResponseTimeP95 != nil {
		_ = d.Set("response_time_p95", job.ResponseTimeP95)
	}

	if job.ResponseTimeP90 != nil {
		_ = d.Set("response_time_p90", job.ResponseTimeP90)
	}

	if job.ResponseTimeMax != nil {
		_ = d.Set("response_time_max", job.ResponseTimeMax)
	}

	if job.ResponseTimeMin != nil {
		_ = d.Set("response_time_min", job.ResponseTimeMin)
	}

	// if job.LoadSourceInfos != nil {
	// 	loadSourceInfosList := []interface{}{}
	// 	for _, loadSourceInfos := range job.LoadSourceInfos {
	// 		loadSourceInfosMap := map[string]interface{}{}
	// 		if loadSourceInfos.IP != nil {
	// 			loadSourceInfosMap["ip"] = loadSourceInfos.IP
	// 		}
	// 		if loadSourceInfos.PodName != nil {
	// 			loadSourceInfosMap["pod_name"] = loadSourceInfos.PodName
	// 		}
	// 		if loadSourceInfos.Region != nil {
	// 			loadSourceInfosMap["region"] = loadSourceInfos.Region
	// 		}

	// 		loadSourceInfosList = append(loadSourceInfosList, loadSourceInfosMap)
	// 	}
	// 	_ = d.Set("load_source_infos", loadSourceInfosList)
	// }

	if job.TestScripts != nil {
		testScriptsList := []interface{}{}
		for _, testScripts := range job.TestScripts {
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
				testScriptsMap["encoded_http_archive"] = testScripts.EncodedHttpArchive
			}
			if testScripts.LoadWeight != nil {
				testScriptsMap["load_weight"] = testScripts.LoadWeight
			}

			testScriptsList = append(testScriptsList, testScriptsMap)
		}
		_ = d.Set("test_scripts", testScriptsList)
	}

	if job.Protocols != nil {
		protocolsList := []interface{}{}
		for _, protocols := range job.Protocols {
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

	if job.RequestFiles != nil {
		requestFilesList := []interface{}{}
		for _, requestFiles := range job.RequestFiles {
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

	if job.Plugins != nil {
		pluginsList := []interface{}{}
		for _, plugins := range job.Plugins {
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

	if job.CronId != nil {
		_ = d.Set("cron_id", job.CronId)
	}

	if job.Type != nil {
		_ = d.Set("type", job.Type)
	}

	if job.DomainNameConfig != nil {
		domainNameConfigMap := map[string]interface{}{}
		if job.DomainNameConfig.HostAliases != nil {
			hostAliasesList := []interface{}{}
			for _, hostAliases := range job.DomainNameConfig.HostAliases {
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
		if job.DomainNameConfig.DNSConfig != nil {
			dNSConfigMap := map[string]interface{}{}
			if job.DomainNameConfig.DNSConfig.Nameservers != nil {
				dNSConfigMap["nameservers"] = job.DomainNameConfig.DNSConfig.Nameservers
			}

			domainNameConfigMap["dns_config"] = []interface{}{dNSConfigMap}
		}

		_ = d.Set("domain_name_config", []interface{}{domainNameConfigMap})
	}

	if job.AbortReason != nil {
		_ = d.Set("abort_reason", job.AbortReason)
	}

	if job.CreatedAt != nil {
		_ = d.Set("created_at", job.CreatedAt)
	}

	// if job.NotificationHooks != nil {
	// 	notificationHooksList := []interface{}{}
	// 	for _, notificationHooks := range job.NotificationHooks {
	// 		notificationHooksMap := map[string]interface{}{}
	// 		if notificationHooks.Events != nil {
	// 			notificationHooksMap["events"] = notificationHooks.Events
	// 		}
	// 		if notificationHooks.URL != nil {
	// 			notificationHooksMap["url"] = notificationHooks.URL
	// 		}

	// 		notificationHooksList = append(notificationHooksList, notificationHooksMap)
	// 	}
	// 	_ = d.Set("notification_hooks", notificationHooksList)
	// }

	return nil
}

func resourceTencentCloudPtsJobUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_job.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := pts.NewUpdateJobRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	scenarioId := idSplit[1]
	jobId := idSplit[2]

	request.ProjectId = &projectId
	request.ScenarioId = &scenarioId
	request.JobId = &jobId

	if d.HasChange("job_owner") {
		return fmt.Errorf("`job_owner` do not support change now.")
	}

	if d.HasChange("debug") {
		return fmt.Errorf("`debug` do not support change now.")
	}

	if d.HasChange("note") {
		if v, ok := d.GetOk("note"); ok {
			request.Note = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePtsClient().UpdateJob(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create pts job failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPtsJobRead(d, meta)
}

func resourceTencentCloudPtsJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_job.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := PtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	scenarioId := idSplit[1]
	jobId := idSplit[2]

	if err := service.DeletePtsJobById(ctx, projectId, scenarioId, jobId); err != nil {
		return err
	}

	return nil
}
