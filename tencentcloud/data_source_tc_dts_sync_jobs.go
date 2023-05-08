/*
Use this data source to query detailed information of dts syncJobs

Example Usage

```hcl
resource "tencentcloud_dts_sync_job" "job" {
	job_name = "tf_dts_test"
	pay_mode = "PostPay"
	src_database_type = "mysql"
	src_region = "ap-guangzhou"
	dst_database_type = "cynosdbmysql"
	dst_region = "ap-guangzhou"
	tags {
	  tag_key = "aaa"
	  tag_value = "bbb"
	}
	auto_renew = 0
   instance_class = "micro"
  }

data "tencentcloud_dts_sync_jobs" "sync_jobs" {
  job_id = tencentcloud_dts_sync_job.job.id
  job_name = "tf_dts_test"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDtsSyncJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDtsSyncJobsRead,
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

			"order": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "order field.",
			},

			"order_seq": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "order way, optional value is DESC or ASC.",
			},

			"status": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "status.",
			},

			"run_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "run mode, optional value is mmediate or Timed.",
			},

			"job_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "job type.",
			},

			"pay_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "pay mode, optional value is PrePay or PostPay.",
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
				Description: "sync job list.",
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
						"pay_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "pay mode.",
						},
						"run_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "run mode.",
						},
						"expect_run_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "expected run time.",
						},
						"all_actions": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "all action list.",
						},
						"actions": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "support action list for current status.",
						},
						"options": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "options.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"init_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "init type.",
									},
									"deal_of_exist_same_table": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "deal of exist same table.",
									},
									"conflict_handle_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "conflict handle type.",
									},
									"add_additional_column": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "add additional column.",
									},
									"op_types": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "operation types.",
									},
									"conflict_handle_option": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "conflict handle option.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"condition_column": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "condition column.",
												},
												"condition_operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "condition override operator.",
												},
												"condition_order_in_src_and_dst": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "condition order in source and destination.",
												},
											},
										},
									},
									"ddl_options": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "ddl options.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ddl_object": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ddl object.",
												},
												"ddl_value": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "ddl value.",
												},
											},
										},
									},
								},
							},
						},
						"objects": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "objects.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "object mode.",
									},
									"databases": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "database list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"db_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "database name.",
												},
												"new_db_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "new database name.",
												},
												"db_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "database mode.",
												},
												"schema_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "schema name.",
												},
												"new_schema_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "new schema name.",
												},
												"table_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "table mode.",
												},
												"tables": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "table list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"table_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "table name.",
															},
															"new_table_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "new table name.",
															},
															"filter_condition": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "filter condition.",
															},
														},
													},
												},
												"view_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "view mode.",
												},
												"views": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "view list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"view_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "view name.",
															},
															"new_view_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "new view name.",
															},
														},
													},
												},
												"function_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "function mode.",
												},
												"functions": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "functions.",
												},
												"procedure_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "procedure mode.",
												},
												"procedures": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "procedures.",
												},
											},
										},
									},
									"advanced_objects": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "advanced objects.",
									},
								},
							},
						},
						"specification": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "specification.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "expire time.",
						},
						"src_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "source region.",
						},
						"src_database_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "source database type.",
						},
						"src_access_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "source access type.",
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
									"db_kernel": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "database kernel.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance id.",
									},
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ip.",
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
									"db_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "database name.",
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
									"cvm_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "cvm instance id.",
									},
									"uniq_dcg_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "dedicated gateway id.",
									},
									"uniq_vpn_gw_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "vpn gateway id.",
									},
									"ccn_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ccn id.",
									},
									"supplier": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "supplier.",
									},
									"engine_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "engine version.",
									},
									"account_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "account mode.",
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
						"dst_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "destination region.",
						},
						"dst_database_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "destination database type.",
						},
						"dst_access_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "destination access type.",
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
									"db_kernel": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "database kernel.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance id.",
									},
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ip.",
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
									"db_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "database name.",
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
									"cvm_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "cvm instance id.",
									},
									"uniq_dcg_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "dedicated gateway id.",
									},
									"uniq_vpn_gw_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "vpn gateway id.",
									},
									"ccn_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ccn id.",
									},
									"supplier": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "supplier.",
									},
									"engine_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "engine version.",
									},
									"account_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "account mode.",
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
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create time.",
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
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "status.",
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
						"detail": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "tag list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"step_all": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "total step numbers.",
									},
									"step_now": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "current step number.",
									},
									"progress": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "progress.",
									},
									"current_step_progress": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "current step progress.",
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
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "message.",
									},
									"step_infos": {
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
												"errors": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "error list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"code": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "code.",
															},
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
													Description: "waring list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"code": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "code.",
															},
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
												"progress": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "current step progress.",
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

func dataSourceTencentCloudDtsSyncJobsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dts_sync_jobs.read")()
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

	if v, ok := d.GetOk("order"); ok {
		paramMap["order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_seq"); ok {
		paramMap["order_seq"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		statusSet := v.(*schema.Set).List()
		tmpList := make([]*string, 0, len(statusSet))
		for i := range statusSet {
			status := statusSet[i].(string)
			tmpList = append(tmpList, helper.String(status))
		}
		paramMap["status"] = tmpList
	}

	if v, ok := d.GetOk("run_mode"); ok {
		paramMap["run_mode"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("job_type"); ok {
		paramMap["job_type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pay_mode"); ok {
		paramMap["pay_mode"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tag_filters"); ok {
		vv := v.([]interface{})
		filters := make([]*dts.TagFilter, 0, len(vv))
		for _, item := range vv {
			dMap := item.(map[string]interface{})
			filter := dts.TagFilter{}
			if v, ok := dMap["tag_key"]; ok {
				filter.TagKey = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				filter.TagValue[0] = helper.String(v.(string))
			}

			filters = append(filters, &filter)
		}
		paramMap["tag_filters"] = filters
	}

	dtsService := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var jobInfos []*dts.SyncJobInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := dtsService.DescribeDtsSyncJobsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		jobInfos = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Dts jobList failed, reason:%+v", logId, err)
		return err
	}

	ids := make([]string, 0, len(jobInfos))
	jobList := make([]map[string]interface{}, 0, len(jobInfos))

	if jobInfos != nil {
		for _, info := range jobInfos {
			jobListMap := map[string]interface{}{}
			if info.JobId != nil {
				jobListMap["job_id"] = info.JobId
			}
			if info.JobName != nil {
				jobListMap["job_name"] = info.JobName
			}
			if info.PayMode != nil {
				jobListMap["pay_mode"] = info.PayMode
			}
			if info.RunMode != nil {
				jobListMap["run_mode"] = info.RunMode
			}
			if info.ExpectRunTime != nil {
				jobListMap["expect_run_time"] = info.ExpectRunTime
			}
			if info.AllActions != nil {
				jobListMap["all_actions"] = info.AllActions
			}
			if info.Actions != nil {
				jobListMap["actions"] = info.Actions
			}
			if info.Options != nil {
				optionsMap := map[string]interface{}{}
				if info.Options.InitType != nil {
					optionsMap["init_type"] = info.Options.InitType
				}
				if info.Options.DealOfExistSameTable != nil {
					optionsMap["deal_of_exist_same_table"] = info.Options.DealOfExistSameTable
				}
				if info.Options.ConflictHandleType != nil {
					optionsMap["conflict_handle_type"] = info.Options.ConflictHandleType
				}
				if info.Options.AddAdditionalColumn != nil {
					optionsMap["add_additional_column"] = info.Options.AddAdditionalColumn
				}
				if info.Options.OpTypes != nil {
					optionsMap["op_types"] = info.Options.OpTypes
				}
				if info.Options.ConflictHandleOption != nil {
					conflictHandleOptionMap := map[string]interface{}{}
					if info.Options.ConflictHandleOption.ConditionColumn != nil {
						conflictHandleOptionMap["condition_column"] = info.Options.ConflictHandleOption.ConditionColumn
					}
					if info.Options.ConflictHandleOption.ConditionOperator != nil {
						conflictHandleOptionMap["condition_operator"] = info.Options.ConflictHandleOption.ConditionOperator
					}
					if info.Options.ConflictHandleOption.ConditionOrderInSrcAndDst != nil {
						conflictHandleOptionMap["condition_order_in_src_and_dst"] = info.Options.ConflictHandleOption.ConditionOrderInSrcAndDst
					}

					optionsMap["conflict_handle_option"] = []interface{}{conflictHandleOptionMap}
				}
				if info.Options.DdlOptions != nil {
					ddlOptionsList := []interface{}{}
					for _, ddlOptions := range info.Options.DdlOptions {
						ddlOptionsMap := map[string]interface{}{}
						if ddlOptions.DdlObject != nil {
							ddlOptionsMap["ddl_object"] = ddlOptions.DdlObject
						}
						if ddlOptions.DdlValue != nil {
							ddlOptionsMap["ddl_value"] = ddlOptions.DdlValue
						}

						ddlOptionsList = append(ddlOptionsList, ddlOptionsMap)
					}
					optionsMap["ddl_options"] = ddlOptionsList
				}

				jobListMap["options"] = []interface{}{optionsMap}
			}
			if info.Objects != nil {
				objectsMap := map[string]interface{}{}
				if info.Objects.Mode != nil {
					objectsMap["mode"] = info.Objects.Mode
				}
				if info.Objects.Databases != nil {
					databasesList := []interface{}{}
					for _, databases := range info.Objects.Databases {
						databasesMap := map[string]interface{}{}
						if databases.DbName != nil {
							databasesMap["db_name"] = databases.DbName
						}
						if databases.NewDbName != nil {
							databasesMap["new_db_name"] = databases.NewDbName
						}
						if databases.DbMode != nil {
							databasesMap["db_mode"] = databases.DbMode
						}
						if databases.SchemaName != nil {
							databasesMap["schema_name"] = databases.SchemaName
						}
						if databases.NewSchemaName != nil {
							databasesMap["new_schema_name"] = databases.NewSchemaName
						}
						if databases.TableMode != nil {
							databasesMap["table_mode"] = databases.TableMode
						}
						if databases.Tables != nil {
							tablesList := []interface{}{}
							for _, tables := range databases.Tables {
								tablesMap := map[string]interface{}{}
								if tables.TableName != nil {
									tablesMap["table_name"] = tables.TableName
								}
								if tables.NewTableName != nil {
									tablesMap["new_table_name"] = tables.NewTableName
								}
								if tables.FilterCondition != nil {
									tablesMap["filter_condition"] = tables.FilterCondition
								}

								tablesList = append(tablesList, tablesMap)
							}
							databasesMap["tables"] = tablesList
						}
						if databases.ViewMode != nil {
							databasesMap["view_mode"] = databases.ViewMode
						}
						if databases.Views != nil {
							viewsList := []interface{}{}
							for _, views := range databases.Views {
								viewsMap := map[string]interface{}{}
								if views.ViewName != nil {
									viewsMap["view_name"] = views.ViewName
								}
								if views.NewViewName != nil {
									viewsMap["new_view_name"] = views.NewViewName
								}

								viewsList = append(viewsList, viewsMap)
							}
							databasesMap["views"] = viewsList
						}
						if databases.FunctionMode != nil {
							databasesMap["function_mode"] = databases.FunctionMode
						}
						if databases.Functions != nil {
							databasesMap["functions"] = databases.Functions
						}
						if databases.ProcedureMode != nil {
							databasesMap["procedure_mode"] = databases.ProcedureMode
						}
						if databases.Procedures != nil {
							databasesMap["procedures"] = databases.Procedures
						}

						databasesList = append(databasesList, databasesMap)
					}
					objectsMap["databases"] = databasesList
				}
				if info.Objects.AdvancedObjects != nil {
					objectsMap["advanced_objects"] = info.Objects.AdvancedObjects
				}

				jobListMap["objects"] = []interface{}{objectsMap}
			}
			if info.Specification != nil {
				jobListMap["specification"] = info.Specification
			}
			if info.ExpireTime != nil {
				jobListMap["expire_time"] = info.ExpireTime
			}
			if info.SrcRegion != nil {
				jobListMap["src_region"] = info.SrcRegion
			}
			if info.SrcDatabaseType != nil {
				jobListMap["src_database_type"] = info.SrcDatabaseType
			}
			if info.SrcAccessType != nil {
				jobListMap["src_access_type"] = info.SrcAccessType
			}
			if info.SrcInfo != nil {
				srcInfoMap := map[string]interface{}{}
				if info.SrcInfo.Region != nil {
					srcInfoMap["region"] = info.SrcInfo.Region
				}
				if info.SrcInfo.DbKernel != nil {
					srcInfoMap["db_kernel"] = info.SrcInfo.DbKernel
				}
				if info.SrcInfo.InstanceId != nil {
					srcInfoMap["instance_id"] = info.SrcInfo.InstanceId
				}
				if info.SrcInfo.Ip != nil {
					srcInfoMap["ip"] = info.SrcInfo.Ip
				}
				if info.SrcInfo.Port != nil {
					srcInfoMap["port"] = info.SrcInfo.Port
				}
				if info.SrcInfo.User != nil {
					srcInfoMap["user"] = info.SrcInfo.User
				}
				if info.SrcInfo.Password != nil {
					srcInfoMap["password"] = info.SrcInfo.Password
				}
				if info.SrcInfo.DbName != nil {
					srcInfoMap["db_name"] = info.SrcInfo.DbName
				}
				if info.SrcInfo.VpcId != nil {
					srcInfoMap["vpc_id"] = info.SrcInfo.VpcId
				}
				if info.SrcInfo.SubnetId != nil {
					srcInfoMap["subnet_id"] = info.SrcInfo.SubnetId
				}
				if info.SrcInfo.CvmInstanceId != nil {
					srcInfoMap["cvm_instance_id"] = info.SrcInfo.CvmInstanceId
				}
				if info.SrcInfo.UniqDcgId != nil {
					srcInfoMap["uniq_dcg_id"] = info.SrcInfo.UniqDcgId
				}
				if info.SrcInfo.UniqVpnGwId != nil {
					srcInfoMap["uniq_vpn_gw_id"] = info.SrcInfo.UniqVpnGwId
				}
				if info.SrcInfo.CcnId != nil {
					srcInfoMap["ccn_id"] = info.SrcInfo.CcnId
				}
				if info.SrcInfo.Supplier != nil {
					srcInfoMap["supplier"] = info.SrcInfo.Supplier
				}
				if info.SrcInfo.EngineVersion != nil {
					srcInfoMap["engine_version"] = info.SrcInfo.EngineVersion
				}
				if info.SrcInfo.AccountMode != nil {
					srcInfoMap["account_mode"] = info.SrcInfo.AccountMode
				}
				if info.SrcInfo.Account != nil {
					srcInfoMap["account"] = info.SrcInfo.Account
				}
				if info.SrcInfo.AccountRole != nil {
					srcInfoMap["account_role"] = info.SrcInfo.AccountRole
				}
				if info.SrcInfo.TmpSecretId != nil {
					srcInfoMap["tmp_secret_id"] = info.SrcInfo.TmpSecretId
				}
				if info.SrcInfo.TmpSecretKey != nil {
					srcInfoMap["tmp_secret_key"] = info.SrcInfo.TmpSecretKey
				}
				if info.SrcInfo.TmpToken != nil {
					srcInfoMap["tmp_token"] = info.SrcInfo.TmpToken
				}

				jobListMap["src_info"] = []interface{}{srcInfoMap}
			}
			if info.DstRegion != nil {
				jobListMap["dst_region"] = info.DstRegion
			}
			if info.DstDatabaseType != nil {
				jobListMap["dst_database_type"] = info.DstDatabaseType
			}
			if info.DstAccessType != nil {
				jobListMap["dst_access_type"] = info.DstAccessType
			}
			if info.DstInfo != nil {
				dstInfoMap := map[string]interface{}{}
				if info.DstInfo.Region != nil {
					dstInfoMap["region"] = info.DstInfo.Region
				}
				if info.DstInfo.DbKernel != nil {
					dstInfoMap["db_kernel"] = info.DstInfo.DbKernel
				}
				if info.DstInfo.InstanceId != nil {
					dstInfoMap["instance_id"] = info.DstInfo.InstanceId
				}
				if info.DstInfo.Ip != nil {
					dstInfoMap["ip"] = info.DstInfo.Ip
				}
				if info.DstInfo.Port != nil {
					dstInfoMap["port"] = info.DstInfo.Port
				}
				if info.DstInfo.User != nil {
					dstInfoMap["user"] = info.DstInfo.User
				}
				if info.DstInfo.Password != nil {
					dstInfoMap["password"] = info.DstInfo.Password
				}
				if info.DstInfo.DbName != nil {
					dstInfoMap["db_name"] = info.DstInfo.DbName
				}
				if info.DstInfo.VpcId != nil {
					dstInfoMap["vpc_id"] = info.DstInfo.VpcId
				}
				if info.DstInfo.SubnetId != nil {
					dstInfoMap["subnet_id"] = info.DstInfo.SubnetId
				}
				if info.DstInfo.CvmInstanceId != nil {
					dstInfoMap["cvm_instance_id"] = info.DstInfo.CvmInstanceId
				}
				if info.DstInfo.UniqDcgId != nil {
					dstInfoMap["uniq_dcg_id"] = info.DstInfo.UniqDcgId
				}
				if info.DstInfo.UniqVpnGwId != nil {
					dstInfoMap["uniq_vpn_gw_id"] = info.DstInfo.UniqVpnGwId
				}
				if info.DstInfo.CcnId != nil {
					dstInfoMap["ccn_id"] = info.DstInfo.CcnId
				}
				if info.DstInfo.Supplier != nil {
					dstInfoMap["supplier"] = info.DstInfo.Supplier
				}
				if info.DstInfo.EngineVersion != nil {
					dstInfoMap["engine_version"] = info.DstInfo.EngineVersion
				}
				if info.DstInfo.AccountMode != nil {
					dstInfoMap["account_mode"] = info.DstInfo.AccountMode
				}
				if info.DstInfo.Account != nil {
					dstInfoMap["account"] = info.DstInfo.Account
				}
				if info.DstInfo.AccountRole != nil {
					dstInfoMap["account_role"] = info.DstInfo.AccountRole
				}
				if info.DstInfo.TmpSecretId != nil {
					dstInfoMap["tmp_secret_id"] = info.DstInfo.TmpSecretId
				}
				if info.DstInfo.TmpSecretKey != nil {
					dstInfoMap["tmp_secret_key"] = info.DstInfo.TmpSecretKey
				}
				if info.DstInfo.TmpToken != nil {
					dstInfoMap["tmp_token"] = info.DstInfo.TmpToken
				}

				jobListMap["dst_info"] = []interface{}{dstInfoMap}
			}
			if info.CreateTime != nil {
				jobListMap["create_time"] = info.CreateTime
			}
			if info.StartTime != nil {
				jobListMap["start_time"] = info.StartTime
			}
			if info.EndTime != nil {
				jobListMap["end_time"] = info.EndTime
			}
			if info.Status != nil {
				jobListMap["status"] = info.Status
			}
			if info.Tags != nil {
				tagsList := []interface{}{}
				for _, tags := range info.Tags {
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
			if info.Detail != nil {
				detailMap := map[string]interface{}{}
				if info.Detail.StepAll != nil {
					detailMap["step_all"] = info.Detail.StepAll
				}
				if info.Detail.StepNow != nil {
					detailMap["step_now"] = info.Detail.StepNow
				}
				if info.Detail.Progress != nil {
					detailMap["progress"] = info.Detail.Progress
				}
				if info.Detail.CurrentStepProgress != nil {
					detailMap["current_step_progress"] = info.Detail.CurrentStepProgress
				}
				if info.Detail.MasterSlaveDistance != nil {
					detailMap["master_slave_distance"] = info.Detail.MasterSlaveDistance
				}
				if info.Detail.SecondsBehindMaster != nil {
					detailMap["seconds_behind_master"] = info.Detail.SecondsBehindMaster
				}
				if info.Detail.Message != nil {
					detailMap["message"] = info.Detail.Message
				}
				if info.Detail.StepInfos != nil {
					stepInfosList := []interface{}{}
					for _, stepInfos := range info.Detail.StepInfos {
						stepInfosMap := map[string]interface{}{}
						if stepInfos.StepNo != nil {
							stepInfosMap["step_no"] = stepInfos.StepNo
						}
						if stepInfos.StepName != nil {
							stepInfosMap["step_name"] = stepInfos.StepName
						}
						if stepInfos.StepId != nil {
							stepInfosMap["step_id"] = stepInfos.StepId
						}
						if stepInfos.Status != nil {
							stepInfosMap["status"] = stepInfos.Status
						}
						if stepInfos.StartTime != nil {
							stepInfosMap["start_time"] = stepInfos.StartTime
						}
						if stepInfos.Errors != nil {
							errorsList := []interface{}{}
							for _, errors := range stepInfos.Errors {
								errorsMap := map[string]interface{}{}
								if errors.Code != nil {
									errorsMap["code"] = errors.Code
								}
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
							stepInfosMap["errors"] = errorsList
						}
						if stepInfos.Warnings != nil {
							warningsList := []interface{}{}
							for _, warnings := range stepInfos.Warnings {
								warningsMap := map[string]interface{}{}
								if warnings.Code != nil {
									warningsMap["code"] = warnings.Code
								}
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
							stepInfosMap["warnings"] = warningsList
						}
						if stepInfos.Progress != nil {
							stepInfosMap["progress"] = stepInfos.Progress
						}

						stepInfosList = append(stepInfosList, stepInfosMap)
					}
					detailMap["step_infos"] = stepInfosList
				}
				jobListMap["detail"] = []interface{}{detailMap}
			}
			ids = append(ids, *info.JobId)
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
