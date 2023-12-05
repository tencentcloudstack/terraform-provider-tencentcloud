package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDtsMigrateJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDtsMigrateJobsRead,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "job id.",
			},

			"job_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "job name.",
			},

			"status": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "migrate status.",
			},

			"src_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "source instance id.",
			},

			"src_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "source region.",
			},

			"src_database_type": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "source database type.",
			},

			"src_access_type": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "source access type.",
			},

			"dst_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "source instance id.",
			},

			"dst_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "destination region.",
			},

			"dst_database_type": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "destination database type.",
			},

			"dst_access_type": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "destination access type.",
			},

			"run_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "run mode.",
			},

			"order_seq": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "order by, default by create time.",
			},

			"tag_filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "tag filters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tag value.",
						},
					},
				},
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "migration job list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "job id.",
						},
						"job_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "job name.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "update time.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "start time.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "end time.",
						},
						"brief_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "brief message for migrate error.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "status.",
						},
						"run_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "run mode, optional value is immediate or Timed.",
						},
						"expect_run_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "expected run time.",
						},
						"action": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "action info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"all_action": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "all action list.",
									},
									"allowed_action": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "allowed action list.",
									},
								},
							},
						},
						"step_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "step info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"step_all": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "number of all steps.",
									},
									"step_now": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "current step.",
									},
									"master_slave_distance": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "master slave distance.",
									},
									"seconds_behind_master": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "seconds behind master.",
									},
									"step_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "step infos.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"step_no": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "step number.",
												},
												"step_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "step name.",
												},
												"step_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "step id.",
												},
												"status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "current status.",
												},
												"start_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "start time.",
												},
												"step_message": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "step message.",
												},
												"percent": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "the percent of miragtion progress.",
												},
												"errors": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "error list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"message": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "message.",
															},
															"solution": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "solution.",
															},
															"help_doc": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "help document.",
															},
														},
													},
												},
												"warnings": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "warning list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"message": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "message.",
															},
															"solution": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "solution.",
															},
															"help_doc": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "help document.",
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
						"src_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "source info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "region.",
									},
									"access_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "access type.",
									},
									"database_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "database type.",
									},
									"node_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "node type.",
									},
									"info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "db info.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"role": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "node role.",
												},
												"db_kernel": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "database kernel.",
												},
												"host": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "host.",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "port.",
												},
												"user": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "user.",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "password.",
												},
												"cvm_instance_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "cvm instance id.",
												},
												"uniq_vpn_gw_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "vpn gateway id.",
												},
												"instance_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "instance id.",
												},
												"ccn_gw_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ccn gateway id.",
												},
												"vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "vpc id.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "subnet id.",
												},
												"engine_version": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "engine version.",
												},
												"account": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "account.",
												},
												"account_role": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "account role.",
												},
												"account_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "account mode.",
												},
												"tmp_secret_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "temporary secret id.",
												},
												"tmp_secret_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "temporary secret key.",
												},
												"tmp_token": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "temporary token.",
												},
											},
										},
									},
									"supplier": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "supplier.",
									},
									"extra_attr": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "extra attributes.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "key.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "value.",
												},
											},
										},
									},
								},
							},
						},
						"dst_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "destination info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "region.",
									},
									"access_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "access type.",
									},
									"database_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "database type.",
									},
									"node_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "node type.",
									},
									"info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "db info.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"role": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "node role.",
												},
												"db_kernel": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "database kernel.",
												},
												"host": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "host.",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "port.",
												},
												"user": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "user.",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "password.",
												},
												"cvm_instance_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "cvm instance id.",
												},
												"uniq_vpn_gw_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "vpn gateway id.",
												},
												"instance_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "instance id.",
												},
												"ccn_gw_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ccn gateway id.",
												},
												"vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "vpc id.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "subnet id.",
												},
												"engine_version": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "engine version.",
												},
												"account": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "account.",
												},
												"account_role": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "account role.",
												},
												"account_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "account mode.",
												},
												"tmp_secret_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "temporary secret id.",
												},
												"tmp_secret_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "temporary secret key.",
												},
												"tmp_token": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "temporary token.",
												},
											},
										},
									},
								},
							},
						},
						"compare_task": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "compare task info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"compare_task_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "compare task id.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "status.",
									},
								},
							},
						},
						"trade_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "trade info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deal_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "deal name.",
									},
									"last_deal_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "last deal name.",
									},
									"instance_class": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance class.",
									},
									"trade_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "trade status.",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "expired time.",
									},
									"offline_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "offline time.",
									},
									"isolate_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "isolate time.",
									},
									"offline_reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "offline reason.",
									},
									"isolate_reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "isolate reason.",
									},
									"pay_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "pay type.",
									},
									"billing_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "billing type.",
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "tag list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "tag value.",
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

func dataSourceTencentCloudDtsMigrateJobsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dts_migrate_jobs.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("job_id"); ok {
		paramMap["job_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("job_name"); ok {
		paramMap["job_name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		statusSet := v.(*schema.Set).List()
		tmpList := make([]*string, 0, len(statusSet))
		for i := range statusSet {
			status := statusSet[i].(string)
			tmpList = append(tmpList, &status)
		}
		paramMap["status"] = tmpList
	}

	if v, ok := d.GetOk("src_instance_id"); ok {
		paramMap["src_instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("src_region"); ok {
		paramMap["src_region"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("src_database_type"); ok {
		srcdbtypeSet := v.(*schema.Set).List()
		tmpList := make([]*string, 0, len(srcdbtypeSet))
		for i := range srcdbtypeSet {
			srcdbtype := srcdbtypeSet[i].(string)
			tmpList = append(tmpList, &srcdbtype)
		}
		paramMap["src_database_type"] = tmpList
	}

	if v, ok := d.GetOk("src_access_type"); ok {
		src_access_typeSet := v.(*schema.Set).List()
		tmpList := make([]*string, 0, len(src_access_typeSet))
		for i := range src_access_typeSet {
			src_access_type := src_access_typeSet[i].(string)
			tmpList = append(tmpList, &src_access_type)
		}
		paramMap["src_access_type"] = tmpList
	}

	if v, ok := d.GetOk("dst_instance_id"); ok {
		paramMap["dst_instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_region"); ok {
		paramMap["dst_region"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_database_type"); ok {
		dstdbtypeSet := v.(*schema.Set).List()
		tmpList := make([]*string, 0, len(dstdbtypeSet))
		for i := range dstdbtypeSet {
			dstdbtype := dstdbtypeSet[i].(string)
			tmpList = append(tmpList, &dstdbtype)
		}
		paramMap["dst_database_type"] = tmpList
	}

	if v, ok := d.GetOk("dst_access_type"); ok {
		dst_access_typeSet := v.(*schema.Set).List()
		tmpList := make([]*string, len(dst_access_typeSet))
		for i := range dst_access_typeSet {
			dst_access_type := dst_access_typeSet[i].(string)
			tmpList = append(tmpList, &dst_access_type)
		}
		paramMap["dst_access_type"] = tmpList
	}

	if v, ok := d.GetOk("run_mode"); ok {
		paramMap["run_mode"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_seq"); ok {
		paramMap["order_seq"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tag_filters"); ok {
		filterSet := v.([]interface{})
		tmpList := make([]*dts.TagFilter, 0, len(filterSet))

		for _, f := range filterSet {
			fMap := f.(map[string]interface{})

			filter := dts.TagFilter{
				TagKey:   helper.String(fMap["tag_key"].(string)),
				TagValue: []*string{helper.String(fMap["tag_value"].(string))},
			}
			tmpList = append(tmpList, &filter)
		}
		paramMap["tag_filters"] = tmpList
	}

	dtsService := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var jobItems []*dts.JobItem
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := dtsService.DescribeDtsMigrateJobsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		jobItems = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Dts jobList failed, reason:%+v", logId, err)
		return err
	}

	ids := make([]string, 0, len(jobItems))
	jobList := make([]map[string]interface{}, 0, len(jobItems))

	if jobItems != nil {
		for _, item := range jobItems {
			jobListMap := map[string]interface{}{}
			if item.JobId != nil {
				jobListMap["job_id"] = item.JobId
			}
			if item.JobName != nil {
				jobListMap["job_name"] = item.JobName
			}
			if item.CreateTime != nil {
				jobListMap["create_time"] = item.CreateTime
			}
			if item.UpdateTime != nil {
				jobListMap["update_time"] = item.UpdateTime
			}
			if item.StartTime != nil {
				jobListMap["start_time"] = item.StartTime
			}
			if item.EndTime != nil {
				jobListMap["end_time"] = item.EndTime
			}
			if item.BriefMsg != nil {
				jobListMap["brief_msg"] = item.BriefMsg
			}
			if item.Status != nil {
				jobListMap["status"] = item.Status
			}
			if item.RunMode != nil {
				jobListMap["run_mode"] = item.RunMode
			}
			if item.ExpectRunTime != nil {
				jobListMap["expect_run_time"] = item.ExpectRunTime
			}
			if item.Action != nil {
				actionMap := map[string]interface{}{}
				if item.Action.AllAction != nil {
					actionMap["all_action"] = item.Action.AllAction
				}
				if item.Action.AllowedAction != nil {
					actionMap["allowed_action"] = item.Action.AllowedAction
				}

				jobListMap["action"] = []interface{}{actionMap}
			}
			if item.StepInfo != nil {
				stepInfoMap := map[string]interface{}{}
				if item.StepInfo.StepAll != nil {
					stepInfoMap["step_all"] = item.StepInfo.StepAll
				}
				if item.StepInfo.StepNow != nil {
					stepInfoMap["step_now"] = item.StepInfo.StepNow
				}
				if item.StepInfo.MasterSlaveDistance != nil {
					stepInfoMap["master_slave_distance"] = item.StepInfo.MasterSlaveDistance
				}
				if item.StepInfo.SecondsBehindMaster != nil {
					stepInfoMap["seconds_behind_master"] = item.StepInfo.SecondsBehindMaster
				}
				if item.StepInfo.StepInfo != nil {
					stepInfoList := []interface{}{}
					for _, stepInfo := range item.StepInfo.StepInfo {
						stepInfoMap := map[string]interface{}{}
						if stepInfo.StepNo != nil {
							stepInfoMap["step_no"] = stepInfo.StepNo
						}
						if stepInfo.StepName != nil {
							stepInfoMap["step_name"] = stepInfo.StepName
						}
						if stepInfo.StepId != nil {
							stepInfoMap["step_id"] = stepInfo.StepId
						}
						if stepInfo.Status != nil {
							stepInfoMap["status"] = stepInfo.Status
						}
						if stepInfo.StartTime != nil {
							stepInfoMap["start_time"] = stepInfo.StartTime
						}
						if stepInfo.StepMessage != nil {
							stepInfoMap["step_message"] = stepInfo.StepMessage
						}
						if stepInfo.Percent != nil {
							stepInfoMap["percent"] = stepInfo.Percent
						}
						if stepInfo.Errors != nil {
							errorsList := []interface{}{}
							for _, errors := range stepInfo.Errors {
								errorsMap := map[string]interface{}{}
								if errors.Message != nil {
									errorsMap["message"] = errors.Message
								}
								if errors.Solution != nil {
									errorsMap["solution"] = errors.Solution
								}
								if errors.HelpDoc != nil {
									errorsMap["help_doc"] = errors.HelpDoc
								}

								errorsList = append(errorsList, errorsMap)
							}
							stepInfoMap["errors"] = errorsList
						}
						if stepInfo.Warnings != nil {
							warningsList := []interface{}{}
							for _, warnings := range stepInfo.Warnings {
								warningsMap := map[string]interface{}{}
								if warnings.Message != nil {
									warningsMap["message"] = warnings.Message
								}
								if warnings.Solution != nil {
									warningsMap["solution"] = warnings.Solution
								}
								if warnings.HelpDoc != nil {
									warningsMap["help_doc"] = warnings.HelpDoc
								}

								warningsList = append(warningsList, warningsMap)
							}
							stepInfoMap["warnings"] = warningsList
						}

						stepInfoList = append(stepInfoList, stepInfoMap)
					}
					stepInfoMap["step_info"] = stepInfoList
				}

				jobListMap["step_info"] = []interface{}{stepInfoMap}
			}
			if item.SrcInfo != nil {
				srcInfoMap := map[string]interface{}{}
				if item.SrcInfo.Region != nil {
					srcInfoMap["region"] = item.SrcInfo.Region
				}
				if item.SrcInfo.AccessType != nil {
					srcInfoMap["access_type"] = item.SrcInfo.AccessType
				}
				if item.SrcInfo.DatabaseType != nil {
					srcInfoMap["database_type"] = item.SrcInfo.DatabaseType
				}
				if item.SrcInfo.NodeType != nil {
					srcInfoMap["node_type"] = item.SrcInfo.NodeType
				}
				if item.SrcInfo.Info != nil {
					infoList := []interface{}{}
					for _, info := range item.SrcInfo.Info {
						infoMap := map[string]interface{}{}
						if info.Role != nil {
							infoMap["role"] = info.Role
						}
						if info.DbKernel != nil {
							infoMap["db_kernel"] = info.DbKernel
						}
						if info.Host != nil {
							infoMap["host"] = info.Host
						}
						if info.Port != nil {
							infoMap["port"] = info.Port
						}
						if info.User != nil {
							infoMap["user"] = info.User
						}
						if info.Password != nil {
							infoMap["password"] = info.Password
						}
						if info.CvmInstanceId != nil {
							infoMap["cvm_instance_id"] = info.CvmInstanceId
						}
						if info.UniqVpnGwId != nil {
							infoMap["uniq_vpn_gw_id"] = info.UniqVpnGwId
						}
						if info.InstanceId != nil {
							infoMap["instance_id"] = info.InstanceId
						}
						if info.CcnGwId != nil {
							infoMap["ccn_gw_id"] = info.CcnGwId
						}
						if info.VpcId != nil {
							infoMap["vpc_id"] = info.VpcId
						}
						if info.SubnetId != nil {
							infoMap["subnet_id"] = info.SubnetId
						}
						if info.EngineVersion != nil {
							infoMap["engine_version"] = info.EngineVersion
						}
						if info.Account != nil {
							infoMap["account"] = info.Account
						}
						if info.AccountRole != nil {
							infoMap["account_role"] = info.AccountRole
						}
						if info.AccountMode != nil {
							infoMap["account_mode"] = info.AccountMode
						}
						if info.TmpSecretId != nil {
							infoMap["tmp_secret_id"] = info.TmpSecretId
						}
						if info.TmpSecretKey != nil {
							infoMap["tmp_secret_key"] = info.TmpSecretKey
						}
						if info.TmpToken != nil {
							infoMap["tmp_token"] = info.TmpToken
						}

						infoList = append(infoList, infoMap)
					}
					srcInfoMap["info"] = infoList
				}
				if item.SrcInfo.Supplier != nil {
					srcInfoMap["supplier"] = item.SrcInfo.Supplier
				}
				if item.SrcInfo.ExtraAttr != nil {
					extraAttrList := []interface{}{}
					for _, extraAttr := range item.SrcInfo.ExtraAttr {
						extraAttrMap := map[string]interface{}{}
						if extraAttr.Key != nil {
							extraAttrMap["key"] = extraAttr.Key
						}
						if extraAttr.Value != nil {
							extraAttrMap["value"] = extraAttr.Value
						}

						extraAttrList = append(extraAttrList, extraAttrMap)
					}
					srcInfoMap["extra_attr"] = extraAttrList
				}

				jobListMap["src_info"] = []interface{}{srcInfoMap}
			}
			if item.DstInfo != nil {
				dstInfoMap := map[string]interface{}{}
				if item.DstInfo.Region != nil {
					dstInfoMap["region"] = item.DstInfo.Region
				}
				if item.DstInfo.AccessType != nil {
					dstInfoMap["access_type"] = item.DstInfo.AccessType
				}
				if item.DstInfo.DatabaseType != nil {
					dstInfoMap["database_type"] = item.DstInfo.DatabaseType
				}
				if item.DstInfo.NodeType != nil {
					dstInfoMap["node_type"] = item.DstInfo.NodeType
				}
				if item.DstInfo.Info != nil {
					infoList := []interface{}{}
					for _, info := range item.DstInfo.Info {
						infoMap := map[string]interface{}{}
						if info.Role != nil {
							infoMap["role"] = info.Role
						}
						if info.DbKernel != nil {
							infoMap["db_kernel"] = info.DbKernel
						}
						if info.Host != nil {
							infoMap["host"] = info.Host
						}
						if info.Port != nil {
							infoMap["port"] = info.Port
						}
						if info.User != nil {
							infoMap["user"] = info.User
						}
						if info.Password != nil {
							infoMap["password"] = info.Password
						}
						if info.CvmInstanceId != nil {
							infoMap["cvm_instance_id"] = info.CvmInstanceId
						}
						if info.UniqVpnGwId != nil {
							infoMap["uniq_vpn_gw_id"] = info.UniqVpnGwId
						}
						if info.InstanceId != nil {
							infoMap["instance_id"] = info.InstanceId
						}
						if info.CcnGwId != nil {
							infoMap["ccn_gw_id"] = info.CcnGwId
						}
						if info.VpcId != nil {
							infoMap["vpc_id"] = info.VpcId
						}
						if info.SubnetId != nil {
							infoMap["subnet_id"] = info.SubnetId
						}
						if info.EngineVersion != nil {
							infoMap["engine_version"] = info.EngineVersion
						}
						if info.Account != nil {
							infoMap["account"] = info.Account
						}
						if info.AccountRole != nil {
							infoMap["account_role"] = info.AccountRole
						}
						if info.AccountMode != nil {
							infoMap["account_mode"] = info.AccountMode
						}
						if info.TmpSecretId != nil {
							infoMap["tmp_secret_id"] = info.TmpSecretId
						}
						if info.TmpSecretKey != nil {
							infoMap["tmp_secret_key"] = info.TmpSecretKey
						}
						if info.TmpToken != nil {
							infoMap["tmp_token"] = info.TmpToken
						}

						infoList = append(infoList, infoMap)
					}
					dstInfoMap["info"] = infoList
				}

				jobListMap["dst_info"] = []interface{}{dstInfoMap}
			}
			if item.CompareTask != nil {
				compareTaskMap := map[string]interface{}{}
				if item.CompareTask.CompareTaskId != nil {
					compareTaskMap["compare_task_id"] = item.CompareTask.CompareTaskId
				}
				if item.CompareTask.Status != nil {
					compareTaskMap["status"] = item.CompareTask.Status
				}

				jobListMap["compare_task"] = []interface{}{compareTaskMap}
			}
			if item.TradeInfo != nil {
				tradeInfoMap := map[string]interface{}{}
				if item.TradeInfo.DealName != nil {
					tradeInfoMap["deal_name"] = item.TradeInfo.DealName
				}
				if item.TradeInfo.LastDealName != nil {
					tradeInfoMap["last_deal_name"] = item.TradeInfo.LastDealName
				}
				if item.TradeInfo.InstanceClass != nil {
					tradeInfoMap["instance_class"] = item.TradeInfo.InstanceClass
				}
				if item.TradeInfo.TradeStatus != nil {
					tradeInfoMap["trade_status"] = item.TradeInfo.TradeStatus
				}
				if item.TradeInfo.ExpireTime != nil {
					tradeInfoMap["expire_time"] = item.TradeInfo.ExpireTime
				}
				if item.TradeInfo.OfflineTime != nil {
					tradeInfoMap["offline_time"] = item.TradeInfo.OfflineTime
				}
				if item.TradeInfo.IsolateTime != nil {
					tradeInfoMap["isolate_time"] = item.TradeInfo.IsolateTime
				}
				if item.TradeInfo.OfflineReason != nil {
					tradeInfoMap["offline_reason"] = item.TradeInfo.OfflineReason
				}
				if item.TradeInfo.IsolateReason != nil {
					tradeInfoMap["isolate_reason"] = item.TradeInfo.IsolateReason
				}
				if item.TradeInfo.PayType != nil {
					tradeInfoMap["pay_type"] = item.TradeInfo.PayType
				}
				if item.TradeInfo.BillingType != nil {
					tradeInfoMap["billing_type"] = item.TradeInfo.BillingType
				}

				jobListMap["trade_info"] = []interface{}{tradeInfoMap}
			}
			if item.Tags != nil {
				tagsList := []interface{}{}
				for _, tags := range item.Tags {
					tagsMap := map[string]interface{}{}
					if tags.TagKey != nil {
						tagsMap["tag_key"] = tags.TagKey
					}
					if tags.TagValue != nil {
						tagsMap["tag_value"] = tags.TagValue
					}

					tagsList = append(tagsList, tagsMap)
				}
				jobListMap["tags"] = tagsList
			}
			ids = append(ids, *item.JobId)
			jobList = append(jobList, jobListMap)
		}
		d.SetId(helper.DataResourceIdsHash(ids))
		_ = d.Set("list", jobList)
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), jobList); e != nil {
			return e
		}
	}

	return nil
}
