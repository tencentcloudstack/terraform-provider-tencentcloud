/*
Provides a resource to create a pts job

Example Usage

```hcl
resource "tencentcloud_pts_job" "job" {
  scenario_id = &lt;nil&gt;
  job_owner = &lt;nil&gt;
  project_id = &lt;nil&gt;
  debug = &lt;nil&gt;
  note = &lt;nil&gt;
                                                        }
```

Import

pts job can be imported using the id, e.g.

```
terraform import tencentcloud_pts_job.job job_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudPtsJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPtsJobCreate,
		Read:   resourceTencentCloudPtsJobRead,
		Update: resourceTencentCloudPtsJobUpdate,
		Delete: resourceTencentCloudPtsJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"scenario_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Pts scenario id.",
			},

			"job_owner": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Job owner.",
			},

			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project ID.",
			},

			"debug": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to debug.",
			},

			"note": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Note.",
			},

			"load": {
				Computed:    true,
				Type:        schema.TypeList,
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
				Computed:    true,
				Type:        schema.TypeList,
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
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The running status of the task; `0`: JobUnknown, `1`: JobCreated, `2`: JobPending, `3`: JobPreparing, `4`: JobSelectClustering, `5`: JobCreateTasking, `6`: JobSyncTasking, `11`: JobRunning, `12`: JobFinished, `13`: JobPrepareException, `14`: JobFinishException, `15`: JobAborting, `16`: JobAborted, `17`: JobAbortException, `18`: JobDeleted, `19`: JobSelectClusterException, `20`: JobCreateTaskException, `21`: JobSyncTaskException.",
			},

			"start_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Start time of the job.",
			},

			"end_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "End time of the job.",
			},

			"max_virtual_user_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Maximum number of VU for the job.",
			},

			"error_rate": {
				Computed:    true,
				Description: "Percentage of error rate.",
			},

			"duration": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Job duration.",
			},

			"max_requests_per_second": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Maximum requests per second.",
			},

			"request_total": {
				Computed:    true,
				Description: "Total number of requests.",
			},

			"requests_per_second": {
				Computed:    true,
				Description: "Average number of requests per second.",
			},

			"response_time_average": {
				Computed:    true,
				Description: "Average response time.",
			},

			"response_time_p99": {
				Computed:    true,
				Description: "99th percentile response time.",
			},

			"response_time_p95": {
				Computed:    true,
				Description: "95th percentile response time.",
			},

			"response_time_p90": {
				Computed:    true,
				Description: "90th percentile response time.",
			},

			"response_time_max": {
				Computed:    true,
				Description: "Maximum response time.",
			},

			"response_time_min": {
				Computed:    true,
				Description: "Minimum response time.",
			},

			"load_source_infos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Host message of generating voltage.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"i_p": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP of host.",
						},
						"pod_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Pod of host.",
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Region to which it belongs.",
						},
					},
				},
			},

			"test_scripts": {
				Computed:    true,
				Type:        schema.TypeList,
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
				Computed:    true,
				Type:        schema.TypeList,
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
				Computed:    true,
				Type:        schema.TypeList,
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
				Computed:    true,
				Type:        schema.TypeList,
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
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scheduled job ID.",
			},

			"type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scene Type.",
			},

			"domain_name_config": {
				Computed:    true,
				Type:        schema.TypeList,
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
									"i_p": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The IP address to be bound.",
									},
								},
							},
						},
						"d_n_s_config": {
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
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Cause of interruption.",
			},

			"created_at": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Creation time of the job.",
			},

			"notification_hooks": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Notification event callback.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"events": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Notification event.",
						},
						"u_r_l": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Callback URL.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudPtsJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_job.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = pts.NewStartJobRequest()
		response = pts.NewStartJobResponse()
		jobId    string
	)
	if v, ok := d.GetOk("scenario_id"); ok {
		request.ScenarioId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("job_owner"); ok {
		request.JobOwner = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("debug"); ok {
		request.Debug = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("note"); ok {
		request.Note = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().StartJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create pts job failed, reason:%+v", logId, err)
		return err
	}

	jobId = *response.Response.JobId
	d.SetId(jobId)

	return resourceTencentCloudPtsJobRead(d, meta)
}

func resourceTencentCloudPtsJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_job.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	jobId := d.Id()

	job, err := service.DescribePtsJobById(ctx, jobId)
	if err != nil {
		return err
	}

	if job == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PtsJob` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

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

					concurrencyMap["stages"] = []interface{}{stagesList}
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

			loadMap["geo_regions_load_distribution"] = []interface{}{geoRegionsLoadDistributionList}
		}

		_ = d.Set("load", []interface{}{loadMap})
	}

	if job.Datasets != nil {
		datasetsList := []interface{}{}
		for _, datasets := range job.Datasets {
			datasetsMap := map[string]interface{}{}

			if job.Datasets.Name != nil {
				datasetsMap["name"] = job.Datasets.Name
			}

			if job.Datasets.Split != nil {
				datasetsMap["split"] = job.Datasets.Split
			}

			if job.Datasets.HeaderInFile != nil {
				datasetsMap["header_in_file"] = job.Datasets.HeaderInFile
			}

			if job.Datasets.HeaderColumns != nil {
				datasetsMap["header_columns"] = job.Datasets.HeaderColumns
			}

			if job.Datasets.LineCount != nil {
				datasetsMap["line_count"] = job.Datasets.LineCount
			}

			if job.Datasets.UpdatedAt != nil {
				datasetsMap["updated_at"] = job.Datasets.UpdatedAt
			}

			if job.Datasets.Size != nil {
				datasetsMap["size"] = job.Datasets.Size
			}

			if job.Datasets.HeadLines != nil {
				datasetsMap["head_lines"] = job.Datasets.HeadLines
			}

			if job.Datasets.TailLines != nil {
				datasetsMap["tail_lines"] = job.Datasets.TailLines
			}

			if job.Datasets.Type != nil {
				datasetsMap["type"] = job.Datasets.Type
			}

			if job.Datasets.FileId != nil {
				datasetsMap["file_id"] = job.Datasets.FileId
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

	if job.MaxRequestsPerSecond != nil {
		_ = d.Set("max_requests_per_second", job.MaxRequestsPerSecond)
	}

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

	if job.LoadSourceInfos != nil {
		loadSourceInfosList := []interface{}{}
		for _, loadSourceInfos := range job.LoadSourceInfos {
			loadSourceInfosMap := map[string]interface{}{}

			if job.LoadSourceInfos.IP != nil {
				loadSourceInfosMap["i_p"] = job.LoadSourceInfos.IP
			}

			if job.LoadSourceInfos.PodName != nil {
				loadSourceInfosMap["pod_name"] = job.LoadSourceInfos.PodName
			}

			if job.LoadSourceInfos.Region != nil {
				loadSourceInfosMap["region"] = job.LoadSourceInfos.Region
			}

			loadSourceInfosList = append(loadSourceInfosList, loadSourceInfosMap)
		}

		_ = d.Set("load_source_infos", loadSourceInfosList)

	}

	if job.TestScripts != nil {
		testScriptsList := []interface{}{}
		for _, testScripts := range job.TestScripts {
			testScriptsMap := map[string]interface{}{}

			if job.TestScripts.Name != nil {
				testScriptsMap["name"] = job.TestScripts.Name
			}

			if job.TestScripts.Size != nil {
				testScriptsMap["size"] = job.TestScripts.Size
			}

			if job.TestScripts.Type != nil {
				testScriptsMap["type"] = job.TestScripts.Type
			}

			if job.TestScripts.UpdatedAt != nil {
				testScriptsMap["updated_at"] = job.TestScripts.UpdatedAt
			}

			if job.TestScripts.EncodedContent != nil {
				testScriptsMap["encoded_content"] = job.TestScripts.EncodedContent
			}

			if job.TestScripts.EncodedHttpArchive != nil {
				testScriptsMap["encoded_http_archive"] = job.TestScripts.EncodedHttpArchive
			}

			if job.TestScripts.LoadWeight != nil {
				testScriptsMap["load_weight"] = job.TestScripts.LoadWeight
			}

			testScriptsList = append(testScriptsList, testScriptsMap)
		}

		_ = d.Set("test_scripts", testScriptsList)

	}

	if job.Protocols != nil {
		protocolsList := []interface{}{}
		for _, protocols := range job.Protocols {
			protocolsMap := map[string]interface{}{}

			if job.Protocols.Name != nil {
				protocolsMap["name"] = job.Protocols.Name
			}

			if job.Protocols.Size != nil {
				protocolsMap["size"] = job.Protocols.Size
			}

			if job.Protocols.Type != nil {
				protocolsMap["type"] = job.Protocols.Type
			}

			if job.Protocols.UpdatedAt != nil {
				protocolsMap["updated_at"] = job.Protocols.UpdatedAt
			}

			if job.Protocols.FileId != nil {
				protocolsMap["file_id"] = job.Protocols.FileId
			}

			protocolsList = append(protocolsList, protocolsMap)
		}

		_ = d.Set("protocols", protocolsList)

	}

	if job.RequestFiles != nil {
		requestFilesList := []interface{}{}
		for _, requestFiles := range job.RequestFiles {
			requestFilesMap := map[string]interface{}{}

			if job.RequestFiles.Name != nil {
				requestFilesMap["name"] = job.RequestFiles.Name
			}

			if job.RequestFiles.Size != nil {
				requestFilesMap["size"] = job.RequestFiles.Size
			}

			if job.RequestFiles.Type != nil {
				requestFilesMap["type"] = job.RequestFiles.Type
			}

			if job.RequestFiles.UpdatedAt != nil {
				requestFilesMap["updated_at"] = job.RequestFiles.UpdatedAt
			}

			if job.RequestFiles.FileId != nil {
				requestFilesMap["file_id"] = job.RequestFiles.FileId
			}

			requestFilesList = append(requestFilesList, requestFilesMap)
		}

		_ = d.Set("request_files", requestFilesList)

	}

	if job.Plugins != nil {
		pluginsList := []interface{}{}
		for _, plugins := range job.Plugins {
			pluginsMap := map[string]interface{}{}

			if job.Plugins.Name != nil {
				pluginsMap["name"] = job.Plugins.Name
			}

			if job.Plugins.Size != nil {
				pluginsMap["size"] = job.Plugins.Size
			}

			if job.Plugins.Type != nil {
				pluginsMap["type"] = job.Plugins.Type
			}

			if job.Plugins.UpdatedAt != nil {
				pluginsMap["updated_at"] = job.Plugins.UpdatedAt
			}

			if job.Plugins.FileId != nil {
				pluginsMap["file_id"] = job.Plugins.FileId
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
					hostAliasesMap["i_p"] = hostAliases.IP
				}

				hostAliasesList = append(hostAliasesList, hostAliasesMap)
			}

			domainNameConfigMap["host_aliases"] = []interface{}{hostAliasesList}
		}

		if job.DomainNameConfig.DNSConfig != nil {
			dNSConfigMap := map[string]interface{}{}

			if job.DomainNameConfig.DNSConfig.Nameservers != nil {
				dNSConfigMap["nameservers"] = job.DomainNameConfig.DNSConfig.Nameservers
			}

			domainNameConfigMap["d_n_s_config"] = []interface{}{dNSConfigMap}
		}

		_ = d.Set("domain_name_config", []interface{}{domainNameConfigMap})
	}

	if job.AbortReason != nil {
		_ = d.Set("abort_reason", job.AbortReason)
	}

	if job.CreatedAt != nil {
		_ = d.Set("created_at", job.CreatedAt)
	}

	if job.NotificationHooks != nil {
		notificationHooksList := []interface{}{}
		for _, notificationHooks := range job.NotificationHooks {
			notificationHooksMap := map[string]interface{}{}

			if job.NotificationHooks.Events != nil {
				notificationHooksMap["events"] = job.NotificationHooks.Events
			}

			if job.NotificationHooks.URL != nil {
				notificationHooksMap["u_r_l"] = job.NotificationHooks.URL
			}

			notificationHooksList = append(notificationHooksList, notificationHooksMap)
		}

		_ = d.Set("notification_hooks", notificationHooksList)

	}

	return nil
}

func resourceTencentCloudPtsJobUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_job.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := pts.NewUpdateJobRequest()

	jobId := d.Id()

	request.JobId = &jobId

	immutableArgs := []string{"scenario_id", "job_owner", "project_id", "debug", "note", "load", "datasets", "status", "start_time", "end_time", "max_virtual_user_count", "error_rate", "duration", "max_requests_per_second", "request_total", "requests_per_second", "response_time_average", "response_time_p99", "response_time_p95", "response_time_p90", "response_time_max", "response_time_min", "load_source_infos", "test_scripts", "protocols", "request_files", "plugins", "cron_id", "type", "domain_name_config", "abort_reason", "created_at", "notification_hooks"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("scenario_id") {
		if v, ok := d.GetOk("scenario_id"); ok {
			request.ScenarioId = helper.String(v.(string))
		}
	}

	if d.HasChange("project_id") {
		if v, ok := d.GetOk("project_id"); ok {
			request.ProjectId = helper.String(v.(string))
		}
	}

	if d.HasChange("note") {
		if v, ok := d.GetOk("note"); ok {
			request.Note = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().UpdateJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update pts job failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPtsJobRead(d, meta)
}

func resourceTencentCloudPtsJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_job.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PtsService{client: meta.(*TencentCloudClient).apiV3Conn}
	jobId := d.Id()

	if err := service.DeletePtsJobById(ctx, jobId); err != nil {
		return err
	}

	return nil
}
