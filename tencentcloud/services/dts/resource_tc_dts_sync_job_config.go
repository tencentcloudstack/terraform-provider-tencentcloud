package dts

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ddlOptionSchemaMap() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ddl_object": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Ddl type, such as Database, Table, View, Index, etc. Note: This field may return null, indicating that no valid value can be obtained.",
		},
		"ddl_value": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional:    true,
			Computed:    true,
			Description: "The specific value of ddl, the possible values for Database [Create,Drop,Alter].The possible values for Table [Create,Drop,Alter,Truncate,Rename].The possible values for View[Create,Drop].For the possible values of Index [Create, Drop]. Note: This field may return null, indicating that no valid value can be obtained.",
		},
	}
}

func opFilterSchemaMap() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"op_types": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional:    true,
			Computed:    true,
			Description: "DML white list.",
		},
		"ddl_options": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "DDL white list.",
			Elem: &schema.Resource{
				Schema: ddlOptionSchemaMap(),
			},
		},
	}
}

func opFilterListSchema(computed bool, desc string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    computed,
		MaxItems:    1,
		Description: desc,
		Elem: &schema.Resource{
			Schema: opFilterSchemaMap(),
		},
	}
}

func ResourceTencentCloudDtsSyncJobConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsSyncJobConfigCreate,
		Read:   resourceTencentCloudDtsSyncJobConfigRead,
		Update: resourceTencentCloudDtsSyncJobConfigUpdate,
		Delete: resourceTencentCloudDtsSyncJobConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Synchronization instance id (i.e. identifies a synchronization job).",
			},

			"dynamic_objects": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Modified sync objects, mirrors SDK Objects struct.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Migration object type Partial (partial object). Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"databases": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Synchronization object, not null when Mode is Partial. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The name of the library that needs to be migrated or synchronized. This item is required when the ObjectMode is Partial. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"new_db_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The name of the library after migration or synchronization, which is the same as the source library by default. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"db_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "DB selection mode: All (for all objects under the current object), Partial (for some objects), when the Mode is Partial, this item is required. Note that synchronization of advanced objects does not depend on this value. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"schema_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Migrated or synchronized schema. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"new_schema_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Schema name after migration or synchronization. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"schema_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Schema selection mode, used by PostgreSQL and SQL Server sync links. Valid values: `All` (all objects under the current object), `Partial` (some objects). Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"table_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Table selection mode: All (for all objects under the current object), Partial (for some objects), this item is required when the DBMode is Partial. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"tables": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "A collection of table graph objects, when TableMode is Partial, this item needs to be filled in. Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"table_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Table name. Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"new_table_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "New table name. Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"filter_condition": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Filter condition. Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"column_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Whether to synchronize all columns in the table, All: all columns under the current table, Partial: part of the columns under the current table, fill in the detailed column information through the columns field. Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"columns": {
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "Synchronized column information, required when ColumnMode is Partial. Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"column_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Column name. Note: This field may return null, indicating that no valid value can be obtained.",
															},
															"new_column_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "New column name. Note: This field may return null, indicating that no valid value can be obtained.",
															},
														},
													},
												},
												"tmp_tables": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Computed:    true,
													Description: "Temporary table generated by pt-osc/gh-ost. Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"table_edit_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Edit table type, rename (table mapping), pt (synchronize attached tables). Note: This field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"view_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "View selection mode: All is all view objects under the current object, Partial is part of the view objects. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"views": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "View object collection, when ViewMode is Partial, this item needs to be filled in. Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"view_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "View name. Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"new_view_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "New view name. Note: This field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"function_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Select the mode to be synchronized, Partial is a part, all is an entire selection. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"functions": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Computed:    true,
										Description: "Required when the FunctionMode value is Partial. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"procedure_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Select the mode to be synchronized, Partial is part, All is the whole selection. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"procedures": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Computed:    true,
										Description: "Required when the value of ProcedureMode is Partial. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"trigger_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Trigger migration mode, all (for all objects under the current object), partial (partial objects). Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"triggers": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Computed:    true,
										Description: "When TriggerMode is partial, specify the name of the trigger to be migrated. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"event_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Event migration mode, all (for all objects under the current object), partial (partial objects). Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"events": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Computed:    true,
										Description: "When EventMode is partial, specify the name of the event to be migrated. Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"advanced_objects": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							Description: "For advanced object types, such as function and procedure, when an advanced object needs to be synchronized, the initialization type must include the structure initialization type, that is, the value of the Options.InitType field is Structure or Full. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"online_ddl": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Deprecated:  "It has been deprecated.",
							Description: "OnlineDDL type. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:        schema.TypeString,
										Required:    true,
										Deprecated:  "It has been deprecated.",
										Description: "Status: ON - Enabled, OFF - Disabled.",
									},
								},
							},
						},
						"databases_op_filter": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Library/table/view level DML/DDL white list. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Database name that the rule takes effect on.",
									},
									"op_filter": opFilterListSchema(true, "Database level DDL/DML filter rule."),
									"tables": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Table level DDL/DML filter info.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"table_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Table name that the rule takes effect on.",
												},
												"op_filter": opFilterListSchema(true, "Table level DDL/DML filter rule."),
											},
										},
									},
									"views": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "View level DDL/DML filter info.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"view_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "View name that the rule takes effect on.",
												},
												"op_filter": opFilterListSchema(true, "View level DDL/DML filter rule."),
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"dynamic_options": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Modified sync task options, mirrors SDK DynamicOptions struct.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"op_types": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "DML and DDL options to be synchronized, Insert (insert operation), Update (update operation), Delete (delete operation), DDL (structure synchronization), PartialDDL (custom, work with DdlOptions). Required field.",
						},
						"ddl_options": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "DDL synchronization options, specifically describe which DDLs to synchronize. Required field.",
							Elem: &schema.Resource{
								Schema: ddlOptionSchemaMap(),
							},
						},
						"conflict_handle_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Conflict handling options, ReportError (error report, the default value), Ignore (ignore), Cover (cover), ConditionCover (condition coverage). Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"conflict_handle_option": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Detailed options for conflict handling, such as conditional rows and conditional actions in conditional overrides. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"condition_column": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Columns covered by the condition. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"condition_operator": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Conditional Override Operation. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"condition_order_in_src_and_dst": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Conditional Override Priority Processing. Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"kafka_option": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Kafka config for sync to kafka link. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Data type delivered to kafka, such as Avro, Json, canal-pb, canal-json, debezium.",
									},
									"topic_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Sync topic strategy, such as Single (centralized delivery to single topic), Multi (custom topic name).",
									},
									"ddl_topic_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Topic used to store ddl.",
									},
									"topic_rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Description of single topic and custom topic.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"topic_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Topic name. When single topic, all TopicName must be the same.",
												},
												"partition_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Topic partition strategy.",
												},
												"db_match_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Database name matching rule, such as Regular, Default.",
												},
												"db_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Database name, takes effect when DbMatchMode=Regular.",
												},
												"table_match_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Table name matching rule, such as Regular, Default.",
												},
												"table_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Table name, only takes effect when TableMatchMode=Regular.",
												},
												"columns": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Computed:    true,
													Description: "Column names for partition by columns, multiple columns can be selected.",
												},
											},
										},
									},
									"data_option": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Other additional information, for specific data types extra parameters can be set.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Option key.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Option value.",
												},
											},
										},
									},
								},
							},
						},
						"filter_begin_commit": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether to filter out begin and commit messages for sync to kafka link. Currently only mysql2kafka link supports it.",
						},
						"filter_checkpoint": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether to filter out checkpoint messages for sync to kafka link. Currently only mysql2kafka link supports it.",
						},
						"deal_of_exist_same_table": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The processing of the table with the same name, ReportErrorAfterCheck (pre-check and report error, default), ExecuteAfterIgnore (ignore and continue to execute). Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"start_position": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Only incremental tasks reset the designated site.",
						},
					},
				},
			},

			"job_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Sync job name.",
			},
			"pay_mode": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Payment method, PostPay (pay as you go), PrePay (monthly subscription).",
			},
			"run_mode": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Operation mode, such as: Immediate (indicates immediate operation, the default value is this value), Timed (indicates scheduled operation).",
			},
			"expect_run_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Expected start time, when the value of RunMode is Timed, this value is required, such as: 2006-01-02 15:04:05.",
			},
			"all_actions": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "All supported actions.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"actions": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Actions that can be performed in the current state.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"specification": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Task specification.",
			},
			"expire_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Expiration time, format is yyyy-mm-dd hh:mm:ss.",
			},
			"src_region": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Source region, such as: ap-guangzhou.",
			},
			"src_database_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Source database type, such as mysql, tdsqlmysql, mariadb, cynosdbmysql, etc.",
			},
			"src_access_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Source access type, cdb (cloud database), cvm (cloud host self-built), vpc (private network), etc.",
			},
			"src_node_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Source node type, cluster or single.",
			},
			"dst_region": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Target region, such as: ap-guangzhou.",
			},
			"dst_database_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Target database type, such as mysql, tdsqlmysql, mariadb, cynosdbmysql, etc.",
			},
			"dst_access_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Target access type, cdb (cloud database), cvm (cloud host self-built), vpc (private network), etc.",
			},
			"dst_node_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Target node type, cluster or single.",
			},
			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Creation time.",
			},
			"start_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Start time.",
			},
			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Task status.",
			},
			"end_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "End time.",
			},
			"trade_status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Status used for billing.",
			},
			"instance_class": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Sync link specification, such as micro, small, medium, large.",
			},
			"auto_renew": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Auto renewal flag, 1 (auto renewal), 0 (no auto renewal).",
			},
			"offline_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Offline time, format is yyyy-mm-dd hh:mm:ss.",
			},
			"opt_obj_status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Dynamic modification object status.",
			},
			"auto_retry_time_range_minutes": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Auto retry time period setting.",
			},
			"dumper_resume_ctrl": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Full export reentrant flag: enum yes/no.",
			},
		},
	}
}

func resourceTencentCloudDtsSyncJobConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_sync_job_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
	)

	var jobId string
	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
	}

	// job config only support job status is running
	request := dts.NewDescribeSyncJobsRequest()
	request.JobId = &jobId
	err := resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDtsClient().DescribeSyncJobs(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dts syncJob failed, Response is nil."))
		}

		if result.Response.JobList == nil || len(result.Response.JobList) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe dts syncJob failed, JobList is nil."))
		}

		jobObj := result.Response.JobList[0]
		if jobObj.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dts syncJob failed, JobList status is nil."))
		}

		if *jobObj.Status == "Running" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Current job %s status is %s, not running...", jobId, *jobObj.Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s Describe dts sync job config failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(jobId)

	return resourceTencentCloudDtsSyncJobConfigUpdate(d, meta)
}

func resourceTencentCloudDtsSyncJobConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_sync_job_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		jobId   = d.Id()
	)

	syncJobConfig, err := service.DescribeDtsSyncJobConfigById(ctx, jobId)
	if err != nil {
		return err
	}

	if syncJobConfig == nil {
		log.Printf("[WARN]%s resource `tencentcloud_dts_sync_job_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("job_id", jobId)

	if syncJobConfig.Objects != nil {
		objectsMap := map[string]interface{}{}

		if syncJobConfig.Objects.Mode != nil {
			objectsMap["mode"] = syncJobConfig.Objects.Mode
		}

		if syncJobConfig.Objects.Databases != nil {
			databasesList := []interface{}{}
			for _, databases := range syncJobConfig.Objects.Databases {
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
				if databases.SchemaMode != nil {
					databasesMap["schema_mode"] = databases.SchemaMode
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
						if tables.ColumnMode != nil {
							tablesMap["column_mode"] = tables.ColumnMode
						}
						if tables.Columns != nil && len(tables.Columns) > 0 {
							columnsList := []interface{}{}
							for _, columns := range tables.Columns {
								columnsMap := map[string]interface{}{}

								if columns.ColumnName != nil {
									columnsMap["column_name"] = columns.ColumnName
								}
								if columns.NewColumnName != nil {
									columnsMap["new_column_name"] = columns.NewColumnName
								}

								columnsList = append(columnsList, columnsMap)
							}

							tablesMap["columns"] = columnsList
						}
						if tables.TmpTables != nil && len(tables.TmpTables) > 0 {
							tablesMap["tmp_tables"] = helper.StringsInterfaces(tables.TmpTables)
						}
						if tables.TableEditMode != nil {
							tablesMap["table_edit_mode"] = tables.TableEditMode
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
					databasesMap["functions"] = helper.StringsInterfaces(databases.Functions)
				}
				if databases.ProcedureMode != nil {
					databasesMap["procedure_mode"] = databases.ProcedureMode
				}
				if databases.Procedures != nil {
					databasesMap["procedures"] = helper.StringsInterfaces(databases.Procedures)
				}
				if databases.TriggerMode != nil {
					databasesMap["trigger_mode"] = databases.TriggerMode
				}
				if databases.Triggers != nil {
					databasesMap["triggers"] = helper.StringsInterfaces(databases.Triggers)
				}
				if databases.EventMode != nil {
					databasesMap["event_mode"] = databases.EventMode
				}
				if databases.Events != nil {
					databasesMap["events"] = helper.StringsInterfaces(databases.Events)
				}

				databasesList = append(databasesList, databasesMap)
			}

			objectsMap["databases"] = databasesList
		}

		if syncJobConfig.Objects.AdvancedObjects != nil {
			objectsMap["advanced_objects"] = helper.StringsInterfaces(syncJobConfig.Objects.AdvancedObjects)
		}

		if syncJobConfig.Objects.OnlineDDL != nil {
			onlineDDLMap := map[string]interface{}{}

			if syncJobConfig.Objects.OnlineDDL.Status != nil {
				onlineDDLMap["status"] = syncJobConfig.Objects.OnlineDDL.Status
			}

			objectsMap["online_ddl"] = []interface{}{onlineDDLMap}
		}

		if syncJobConfig.Objects.DatabasesOpFilter != nil {
			databasesOpFilterList := []interface{}{}
			for _, dbOpFilter := range syncJobConfig.Objects.DatabasesOpFilter {
				dbOpFilterMap := map[string]interface{}{}

				if dbOpFilter.DbName != nil {
					dbOpFilterMap["db_name"] = dbOpFilter.DbName
				}
				if dbOpFilter.OpFilter != nil {
					opFilterMap := map[string]interface{}{}
					if dbOpFilter.OpFilter.OpTypes != nil {
						opFilterMap["op_types"] = helper.StringsInterfaces(dbOpFilter.OpFilter.OpTypes)
					}
					if dbOpFilter.OpFilter.DdlOptions != nil {
						ddlOptionsList := []interface{}{}
						for _, ddlOption := range dbOpFilter.OpFilter.DdlOptions {
							ddlOptionMap := map[string]interface{}{}
							if ddlOption.DdlObject != nil {
								ddlOptionMap["ddl_object"] = ddlOption.DdlObject
							}
							if ddlOption.DdlValue != nil {
								ddlOptionMap["ddl_value"] = ddlOption.DdlValue
							}
							ddlOptionsList = append(ddlOptionsList, ddlOptionMap)
						}
						opFilterMap["ddl_options"] = ddlOptionsList
					}
					dbOpFilterMap["op_filter"] = []interface{}{opFilterMap}
				}
				if dbOpFilter.Tables != nil {
					tablesList := []interface{}{}
					for _, tableFilter := range dbOpFilter.Tables {
						tableFilterMap := map[string]interface{}{}
						if tableFilter.TableName != nil {
							tableFilterMap["table_name"] = tableFilter.TableName
						}
						if tableFilter.OpFilter != nil {
							opFilterMap := map[string]interface{}{}
							if tableFilter.OpFilter.OpTypes != nil {
								opFilterMap["op_types"] = helper.StringsInterfaces(tableFilter.OpFilter.OpTypes)
							}
							if tableFilter.OpFilter.DdlOptions != nil {
								ddlOptionsList := []interface{}{}
								for _, ddlOption := range tableFilter.OpFilter.DdlOptions {
									ddlOptionMap := map[string]interface{}{}
									if ddlOption.DdlObject != nil {
										ddlOptionMap["ddl_object"] = ddlOption.DdlObject
									}
									if ddlOption.DdlValue != nil {
										ddlOptionMap["ddl_value"] = ddlOption.DdlValue
									}
									ddlOptionsList = append(ddlOptionsList, ddlOptionMap)
								}
								opFilterMap["ddl_options"] = ddlOptionsList
							}
							tableFilterMap["op_filter"] = []interface{}{opFilterMap}
						}
						tablesList = append(tablesList, tableFilterMap)
					}
					dbOpFilterMap["tables"] = tablesList
				}
				if dbOpFilter.Views != nil {
					viewsList := []interface{}{}
					for _, viewFilter := range dbOpFilter.Views {
						viewFilterMap := map[string]interface{}{}
						if viewFilter.ViewName != nil {
							viewFilterMap["view_name"] = viewFilter.ViewName
						}
						if viewFilter.OpFilter != nil {
							opFilterMap := map[string]interface{}{}
							if viewFilter.OpFilter.OpTypes != nil {
								opFilterMap["op_types"] = helper.StringsInterfaces(viewFilter.OpFilter.OpTypes)
							}
							if viewFilter.OpFilter.DdlOptions != nil {
								ddlOptionsList := []interface{}{}
								for _, ddlOption := range viewFilter.OpFilter.DdlOptions {
									ddlOptionMap := map[string]interface{}{}
									if ddlOption.DdlObject != nil {
										ddlOptionMap["ddl_object"] = ddlOption.DdlObject
									}
									if ddlOption.DdlValue != nil {
										ddlOptionMap["ddl_value"] = ddlOption.DdlValue
									}
									ddlOptionsList = append(ddlOptionsList, ddlOptionMap)
								}
								opFilterMap["ddl_options"] = ddlOptionsList
							}
							viewFilterMap["op_filter"] = []interface{}{opFilterMap}
						}
						viewsList = append(viewsList, viewFilterMap)
					}
					dbOpFilterMap["views"] = viewsList
				}

				databasesOpFilterList = append(databasesOpFilterList, dbOpFilterMap)
			}

			objectsMap["databases_op_filter"] = databasesOpFilterList
		}

		_ = d.Set("dynamic_objects", []interface{}{objectsMap})
	}

	if syncJobConfig.Options != nil {
		dynamicOptionsMap := map[string]interface{}{}

		if syncJobConfig.Options.OpTypes != nil {
			dynamicOptionsMap["op_types"] = helper.StringsInterfaces(syncJobConfig.Options.OpTypes)
		}

		if syncJobConfig.Options.DdlOptions != nil {
			ddlOptionsList := []interface{}{}
			for _, ddlOptions := range syncJobConfig.Options.DdlOptions {
				ddlOptionsMap := map[string]interface{}{}

				if ddlOptions.DdlObject != nil {
					ddlOptionsMap["ddl_object"] = ddlOptions.DdlObject
				}

				if ddlOptions.DdlValue != nil {
					ddlOptionsMap["ddl_value"] = ddlOptions.DdlValue
				}

				ddlOptionsList = append(ddlOptionsList, ddlOptionsMap)
			}

			dynamicOptionsMap["ddl_options"] = ddlOptionsList
		}

		if syncJobConfig.Options.ConflictHandleType != nil {
			dynamicOptionsMap["conflict_handle_type"] = syncJobConfig.Options.ConflictHandleType
		}

		if syncJobConfig.Options.ConflictHandleOption != nil {
			conflictHandleOptionMap := map[string]interface{}{}

			if syncJobConfig.Options.ConflictHandleOption.ConditionColumn != nil {
				conflictHandleOptionMap["condition_column"] = syncJobConfig.Options.ConflictHandleOption.ConditionColumn
			}

			if syncJobConfig.Options.ConflictHandleOption.ConditionOperator != nil {
				conflictHandleOptionMap["condition_operator"] = syncJobConfig.Options.ConflictHandleOption.ConditionOperator
			}

			if syncJobConfig.Options.ConflictHandleOption.ConditionOrderInSrcAndDst != nil {
				conflictHandleOptionMap["condition_order_in_src_and_dst"] = syncJobConfig.Options.ConflictHandleOption.ConditionOrderInSrcAndDst
			}

			dynamicOptionsMap["conflict_handle_option"] = []interface{}{conflictHandleOptionMap}
		}

		if syncJobConfig.Options.KafkaOption != nil {
			kafkaOptionMap := map[string]interface{}{}

			if syncJobConfig.Options.KafkaOption.DataType != nil {
				kafkaOptionMap["data_type"] = syncJobConfig.Options.KafkaOption.DataType
			}

			if syncJobConfig.Options.KafkaOption.TopicType != nil {
				kafkaOptionMap["topic_type"] = syncJobConfig.Options.KafkaOption.TopicType
			}

			if syncJobConfig.Options.KafkaOption.DDLTopicName != nil {
				kafkaOptionMap["ddl_topic_name"] = syncJobConfig.Options.KafkaOption.DDLTopicName
			}

			if syncJobConfig.Options.KafkaOption.TopicRules != nil {
				topicRulesList := []interface{}{}
				for _, topicRule := range syncJobConfig.Options.KafkaOption.TopicRules {
					topicRuleMap := map[string]interface{}{}

					if topicRule.TopicName != nil {
						topicRuleMap["topic_name"] = topicRule.TopicName
					}
					if topicRule.PartitionType != nil {
						topicRuleMap["partition_type"] = topicRule.PartitionType
					}
					if topicRule.DbMatchMode != nil {
						topicRuleMap["db_match_mode"] = topicRule.DbMatchMode
					}
					if topicRule.DbName != nil {
						topicRuleMap["db_name"] = topicRule.DbName
					}
					if topicRule.TableMatchMode != nil {
						topicRuleMap["table_match_mode"] = topicRule.TableMatchMode
					}
					if topicRule.TableName != nil {
						topicRuleMap["table_name"] = topicRule.TableName
					}
					if topicRule.Columns != nil {
						topicRuleMap["columns"] = helper.StringsInterfaces(topicRule.Columns)
					}

					topicRulesList = append(topicRulesList, topicRuleMap)
				}

				kafkaOptionMap["topic_rules"] = topicRulesList
			}

			if syncJobConfig.Options.KafkaOption.DataOption != nil {
				dataOptionList := []interface{}{}
				for _, kvOption := range syncJobConfig.Options.KafkaOption.DataOption {
					kvOptionMap := map[string]interface{}{}

					if kvOption.Key != nil {
						kvOptionMap["key"] = kvOption.Key
					}
					if kvOption.Value != nil {
						kvOptionMap["value"] = kvOption.Value
					}

					dataOptionList = append(dataOptionList, kvOptionMap)
				}

				kafkaOptionMap["data_option"] = dataOptionList
			}

			dynamicOptionsMap["kafka_option"] = []interface{}{kafkaOptionMap}
		}

		if syncJobConfig.Options.FilterBeginCommit != nil {
			dynamicOptionsMap["filter_begin_commit"] = syncJobConfig.Options.FilterBeginCommit
		}

		if syncJobConfig.Options.FilterCheckpoint != nil {
			dynamicOptionsMap["filter_checkpoint"] = syncJobConfig.Options.FilterCheckpoint
		}

		if syncJobConfig.Options.DealOfExistSameTable != nil {
			dynamicOptionsMap["deal_of_exist_same_table"] = syncJobConfig.Options.DealOfExistSameTable
		}

		if syncJobConfig.Options.StartPosition != nil {
			dynamicOptionsMap["start_position"] = syncJobConfig.Options.StartPosition
		}

		_ = d.Set("dynamic_options", []interface{}{dynamicOptionsMap})
	}

	if syncJobConfig.JobName != nil {
		_ = d.Set("job_name", syncJobConfig.JobName)
	}

	if syncJobConfig.PayMode != nil {
		_ = d.Set("pay_mode", syncJobConfig.PayMode)
	}

	if syncJobConfig.RunMode != nil {
		_ = d.Set("run_mode", syncJobConfig.RunMode)
	}

	if syncJobConfig.ExpectRunTime != nil {
		_ = d.Set("expect_run_time", syncJobConfig.ExpectRunTime)
	}

	if syncJobConfig.AllActions != nil {
		_ = d.Set("all_actions", helper.StringsInterfaces(syncJobConfig.AllActions))
	}

	if syncJobConfig.Actions != nil {
		_ = d.Set("actions", helper.StringsInterfaces(syncJobConfig.Actions))
	}

	if syncJobConfig.Specification != nil {
		_ = d.Set("specification", syncJobConfig.Specification)
	}

	if syncJobConfig.ExpireTime != nil {
		_ = d.Set("expire_time", syncJobConfig.ExpireTime)
	}

	if syncJobConfig.SrcRegion != nil {
		_ = d.Set("src_region", syncJobConfig.SrcRegion)
	}

	if syncJobConfig.SrcDatabaseType != nil {
		_ = d.Set("src_database_type", syncJobConfig.SrcDatabaseType)
	}

	if syncJobConfig.SrcAccessType != nil {
		_ = d.Set("src_access_type", syncJobConfig.SrcAccessType)
	}

	if syncJobConfig.SrcNodeType != nil {
		_ = d.Set("src_node_type", syncJobConfig.SrcNodeType)
	}

	if syncJobConfig.DstRegion != nil {
		_ = d.Set("dst_region", syncJobConfig.DstRegion)
	}

	if syncJobConfig.DstDatabaseType != nil {
		_ = d.Set("dst_database_type", syncJobConfig.DstDatabaseType)
	}

	if syncJobConfig.DstAccessType != nil {
		_ = d.Set("dst_access_type", syncJobConfig.DstAccessType)
	}

	if syncJobConfig.DstNodeType != nil {
		_ = d.Set("dst_node_type", syncJobConfig.DstNodeType)
	}

	if syncJobConfig.CreateTime != nil {
		_ = d.Set("create_time", syncJobConfig.CreateTime)
	}

	if syncJobConfig.StartTime != nil {
		_ = d.Set("start_time", syncJobConfig.StartTime)
	}

	if syncJobConfig.Status != nil {
		_ = d.Set("status", syncJobConfig.Status)
	}

	if syncJobConfig.EndTime != nil {
		_ = d.Set("end_time", syncJobConfig.EndTime)
	}

	if syncJobConfig.TradeStatus != nil {
		_ = d.Set("trade_status", syncJobConfig.TradeStatus)
	}

	if syncJobConfig.InstanceClass != nil {
		_ = d.Set("instance_class", syncJobConfig.InstanceClass)
	}

	if syncJobConfig.AutoRenew != nil {
		_ = d.Set("auto_renew", syncJobConfig.AutoRenew)
	}

	if syncJobConfig.OfflineTime != nil {
		_ = d.Set("offline_time", syncJobConfig.OfflineTime)
	}

	if syncJobConfig.OptObjStatus != nil {
		_ = d.Set("opt_obj_status", syncJobConfig.OptObjStatus)
	}

	if syncJobConfig.AutoRetryTimeRangeMinutes != nil {
		_ = d.Set("auto_retry_time_range_minutes", syncJobConfig.AutoRetryTimeRangeMinutes)
	}

	if syncJobConfig.DumperResumeCtrl != nil {
		_ = d.Set("dumper_resume_ctrl", syncJobConfig.DumperResumeCtrl)
	}

	return nil
}

func flattenDtsSyncJobConfigDdlOptions(ddlOptionsList []interface{}) []*dts.DdlOption {
	var ddlOptions []*dts.DdlOption
	for _, item := range ddlOptionsList {
		ddlOptionsMap := item.(map[string]interface{})
		ddlOption := dts.DdlOption{}
		if v, ok := ddlOptionsMap["ddl_object"]; ok {
			ddlOption.DdlObject = helper.String(v.(string))
		}
		if v, ok := ddlOptionsMap["ddl_value"]; ok {
			ddlValueSet := v.(*schema.Set).List()
			for i := range ddlValueSet {
				if ddlValueSet[i] != nil {
					ddlValue := ddlValueSet[i].(string)
					ddlOption.DdlValue = append(ddlOption.DdlValue, &ddlValue)
				}
			}
		}
		ddlOptions = append(ddlOptions, &ddlOption)
	}
	return ddlOptions
}

func flattenDtsSyncJobConfigOpFilter(opFilterMap map[string]interface{}) *dts.OpFilter {
	opFilter := dts.OpFilter{}
	if v, ok := opFilterMap["op_types"]; ok {
		opTypesSet := v.(*schema.Set).List()
		for i := range opTypesSet {
			if opTypesSet[i] != nil {
				opType := opTypesSet[i].(string)
				opFilter.OpTypes = append(opFilter.OpTypes, &opType)
			}
		}
	}
	if v, ok := opFilterMap["ddl_options"]; ok {
		opFilter.DdlOptions = flattenDtsSyncJobConfigDdlOptions(v.([]interface{}))
	}
	return &opFilter
}

func resourceTencentCloudDtsSyncJobConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_sync_job_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request = dts.NewModifySyncJobConfigRequest()
		jobId   = d.Id()
	)

	if dMap, ok := helper.InterfacesHeadMap(d, "dynamic_objects"); ok {
		objects := dts.Objects{}
		if v, ok := dMap["mode"]; ok {
			objects.Mode = helper.String(v.(string))
		}
		if v, ok := dMap["databases"]; ok {
			for _, item := range v.([]interface{}) {
				databasesMap := item.(map[string]interface{})
				database := dts.Database{}
				if v, ok := databasesMap["db_name"]; ok {
					database.DbName = helper.String(v.(string))
				}
				if v, ok := databasesMap["new_db_name"]; ok {
					database.NewDbName = helper.String(v.(string))
				}
				if v, ok := databasesMap["db_mode"]; ok {
					database.DbMode = helper.String(v.(string))
				}
				if v, ok := databasesMap["schema_name"]; ok {
					database.SchemaName = helper.String(v.(string))
				}
				if v, ok := databasesMap["new_schema_name"]; ok {
					database.NewSchemaName = helper.String(v.(string))
				}
				if v, ok := databasesMap["schema_mode"]; ok {
					database.SchemaMode = helper.String(v.(string))
				}
				if v, ok := databasesMap["table_mode"]; ok {
					database.TableMode = helper.String(v.(string))
				}
				if v, ok := databasesMap["tables"]; ok {
					for _, item := range v.([]interface{}) {
						tablesMap := item.(map[string]interface{})
						table := dts.Table{}
						if v, ok := tablesMap["table_name"]; ok {
							table.TableName = helper.String(v.(string))
						}
						if v, ok := tablesMap["new_table_name"]; ok {
							table.NewTableName = helper.String(v.(string))
						}
						if v, ok := tablesMap["filter_condition"]; ok {
							table.FilterCondition = helper.String(v.(string))
						}
						if v, ok := tablesMap["column_mode"]; ok {
							table.ColumnMode = helper.String(v.(string))
						}
						if v, ok := tablesMap["columns"]; ok {
							for _, item := range v.([]interface{}) {
								columnsMap := item.(map[string]interface{})
								column := dts.Column{}
								if v, ok := columnsMap["column_name"]; ok {
									column.ColumnName = helper.String(v.(string))
								}
								if v, ok := columnsMap["new_column_name"]; ok {
									column.NewColumnName = helper.String(v.(string))
								}
								table.Columns = append(table.Columns, &column)
							}
						}
						if v, ok := tablesMap["tmp_tables"]; ok {
							tmpTablesSet := v.(*schema.Set).List()
							for i := range tmpTablesSet {
								if tmpTablesSet[i] != nil {
									tmpTables := tmpTablesSet[i].(string)
									table.TmpTables = append(table.TmpTables, &tmpTables)
								}
							}
						}
						if v, ok := tablesMap["table_edit_mode"]; ok {
							table.TableEditMode = helper.String(v.(string))
						}
						database.Tables = append(database.Tables, &table)
					}
				}
				if v, ok := databasesMap["view_mode"]; ok {
					database.ViewMode = helper.String(v.(string))
				}
				if v, ok := databasesMap["views"]; ok {
					for _, item := range v.([]interface{}) {
						viewsMap := item.(map[string]interface{})
						view := dts.View{}
						if v, ok := viewsMap["view_name"]; ok {
							view.ViewName = helper.String(v.(string))
						}
						if v, ok := viewsMap["new_view_name"]; ok {
							view.NewViewName = helper.String(v.(string))
						}
						database.Views = append(database.Views, &view)
					}
				}
				if v, ok := databasesMap["function_mode"]; ok {
					database.FunctionMode = helper.String(v.(string))
				}
				if v, ok := databasesMap["functions"]; ok {
					functionsSet := v.(*schema.Set).List()
					for i := range functionsSet {
						if functionsSet[i] != nil {
							functions := functionsSet[i].(string)
							database.Functions = append(database.Functions, &functions)
						}
					}
				}
				if v, ok := databasesMap["procedure_mode"]; ok {
					database.ProcedureMode = helper.String(v.(string))
				}
				if v, ok := databasesMap["procedures"]; ok {
					proceduresSet := v.(*schema.Set).List()
					for i := range proceduresSet {
						if proceduresSet[i] != nil {
							procedures := proceduresSet[i].(string)
							database.Procedures = append(database.Procedures, &procedures)
						}
					}
				}
				if v, ok := databasesMap["trigger_mode"]; ok {
					database.TriggerMode = helper.String(v.(string))
				}
				if v, ok := databasesMap["triggers"]; ok {
					triggersSet := v.(*schema.Set).List()
					for i := range triggersSet {
						if triggersSet[i] != nil {
							triggers := triggersSet[i].(string)
							database.Triggers = append(database.Triggers, &triggers)
						}
					}
				}
				if v, ok := databasesMap["event_mode"]; ok {
					database.EventMode = helper.String(v.(string))
				}
				if v, ok := databasesMap["events"]; ok {
					eventsSet := v.(*schema.Set).List()
					for i := range eventsSet {
						if eventsSet[i] != nil {
							events := eventsSet[i].(string)
							database.Events = append(database.Events, &events)
						}
					}
				}
				objects.Databases = append(objects.Databases, &database)
			}
		}
		if v, ok := dMap["advanced_objects"]; ok {
			advancedObjectsSet := v.(*schema.Set).List()
			for i := range advancedObjectsSet {
				if advancedObjectsSet[i] != nil {
					advancedObjects := advancedObjectsSet[i].(string)
					objects.AdvancedObjects = append(objects.AdvancedObjects, &advancedObjects)
				}
			}
		}
		if onlineDDLMap, ok := helper.InterfaceToMap(dMap, "online_ddl"); ok {
			onlineDDL := dts.OnlineDDL{}
			if v, ok := onlineDDLMap["status"]; ok {
				onlineDDL.Status = helper.String(v.(string))
			}
			objects.OnlineDDL = &onlineDDL
		}
		if v, ok := dMap["databases_op_filter"]; ok {
			for _, item := range v.([]interface{}) {
				dbOpFilterMap := item.(map[string]interface{})
				dbOpFilter := dts.DBOpFilter{}
				if v, ok := dbOpFilterMap["db_name"]; ok {
					dbOpFilter.DbName = helper.String(v.(string))
				}
				if opFilterMap, ok := helper.InterfaceToMap(dbOpFilterMap, "op_filter"); ok {
					dbOpFilter.OpFilter = flattenDtsSyncJobConfigOpFilter(opFilterMap)
				}
				if v, ok := dbOpFilterMap["tables"]; ok {
					for _, tableItem := range v.([]interface{}) {
						tableFilterMap := tableItem.(map[string]interface{})
						tableFilter := dts.TableFilter{}
						if v, ok := tableFilterMap["table_name"]; ok {
							tableFilter.TableName = helper.String(v.(string))
						}
						if opFilterMap, ok := helper.InterfaceToMap(tableFilterMap, "op_filter"); ok {
							tableFilter.OpFilter = flattenDtsSyncJobConfigOpFilter(opFilterMap)
						}
						dbOpFilter.Tables = append(dbOpFilter.Tables, &tableFilter)
					}
				}
				if v, ok := dbOpFilterMap["views"]; ok {
					for _, viewItem := range v.([]interface{}) {
						viewFilterMap := viewItem.(map[string]interface{})
						viewFilter := dts.ViewFilter{}
						if v, ok := viewFilterMap["view_name"]; ok {
							viewFilter.ViewName = helper.String(v.(string))
						}
						if opFilterMap, ok := helper.InterfaceToMap(viewFilterMap, "op_filter"); ok {
							viewFilter.OpFilter = flattenDtsSyncJobConfigOpFilter(opFilterMap)
						}
						dbOpFilter.Views = append(dbOpFilter.Views, &viewFilter)
					}
				}
				objects.DatabasesOpFilter = append(objects.DatabasesOpFilter, &dbOpFilter)
			}
		}
		request.DynamicObjects = &objects
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "dynamic_options"); ok {
		dynamicOptions := dts.DynamicOptions{}
		if v, ok := dMap["op_types"]; ok {
			opTypesSet := v.(*schema.Set).List()
			for i := range opTypesSet {
				if opTypesSet[i] != nil {
					opType := opTypesSet[i].(string)
					dynamicOptions.OpTypes = append(dynamicOptions.OpTypes, &opType)
				}
			}
		}
		if v, ok := dMap["ddl_options"]; ok {
			dynamicOptions.DdlOptions = flattenDtsSyncJobConfigDdlOptions(v.([]interface{}))
		}
		if v, ok := dMap["conflict_handle_type"]; ok {
			dynamicOptions.ConflictHandleType = helper.String(v.(string))
		}
		if conflictHandleOptionMap, ok := helper.InterfaceToMap(dMap, "conflict_handle_option"); ok {
			conflictHandleOption := dts.ConflictHandleOption{}
			if v, ok := conflictHandleOptionMap["condition_column"]; ok {
				conflictHandleOption.ConditionColumn = helper.String(v.(string))
			}
			if v, ok := conflictHandleOptionMap["condition_operator"]; ok {
				conflictHandleOption.ConditionOperator = helper.String(v.(string))
			}
			if v, ok := conflictHandleOptionMap["condition_order_in_src_and_dst"]; ok {
				conflictHandleOption.ConditionOrderInSrcAndDst = helper.String(v.(string))
			}
			dynamicOptions.ConflictHandleOption = &conflictHandleOption
		}
		if kafkaOptionMap, ok := helper.InterfaceToMap(dMap, "kafka_option"); ok {
			kafkaOption := dts.KafkaOption{}
			if v, ok := kafkaOptionMap["data_type"]; ok {
				kafkaOption.DataType = helper.String(v.(string))
			}
			if v, ok := kafkaOptionMap["topic_type"]; ok {
				kafkaOption.TopicType = helper.String(v.(string))
			}
			if v, ok := kafkaOptionMap["ddl_topic_name"]; ok {
				kafkaOption.DDLTopicName = helper.String(v.(string))
			}
			if v, ok := kafkaOptionMap["topic_rules"]; ok {
				for _, topicItem := range v.([]interface{}) {
					topicRuleMap := topicItem.(map[string]interface{})
					topicRule := dts.TopicRule{}
					if v, ok := topicRuleMap["topic_name"]; ok {
						topicRule.TopicName = helper.String(v.(string))
					}
					if v, ok := topicRuleMap["partition_type"]; ok {
						topicRule.PartitionType = helper.String(v.(string))
					}
					if v, ok := topicRuleMap["db_match_mode"]; ok {
						topicRule.DbMatchMode = helper.String(v.(string))
					}
					if v, ok := topicRuleMap["db_name"]; ok {
						topicRule.DbName = helper.String(v.(string))
					}
					if v, ok := topicRuleMap["table_match_mode"]; ok {
						topicRule.TableMatchMode = helper.String(v.(string))
					}
					if v, ok := topicRuleMap["table_name"]; ok {
						topicRule.TableName = helper.String(v.(string))
					}
					if v, ok := topicRuleMap["columns"]; ok {
						columnsSet := v.(*schema.Set).List()
						for i := range columnsSet {
							if columnsSet[i] != nil {
								column := columnsSet[i].(string)
								topicRule.Columns = append(topicRule.Columns, &column)
							}
						}
					}
					kafkaOption.TopicRules = append(kafkaOption.TopicRules, &topicRule)
				}
			}
			if v, ok := kafkaOptionMap["data_option"]; ok {
				for _, kvItem := range v.([]interface{}) {
					kvOptionMap := kvItem.(map[string]interface{})
					kvOption := dts.KeyValuePairOption{}
					if v, ok := kvOptionMap["key"]; ok {
						kvOption.Key = helper.String(v.(string))
					}
					if v, ok := kvOptionMap["value"]; ok {
						kvOption.Value = helper.String(v.(string))
					}
					kafkaOption.DataOption = append(kafkaOption.DataOption, &kvOption)
				}
			}
			dynamicOptions.KafkaOption = &kafkaOption
		}
		if v, ok := dMap["filter_begin_commit"]; ok {
			dynamicOptions.FilterBeginCommit = helper.Bool(v.(bool))
		}
		if v, ok := dMap["filter_checkpoint"]; ok {
			dynamicOptions.FilterCheckpoint = helper.Bool(v.(bool))
		}
		if v, ok := dMap["deal_of_exist_same_table"]; ok {
			dynamicOptions.DealOfExistSameTable = helper.String(v.(string))
		}
		if v, ok := dMap["start_position"]; ok {
			dynamicOptions.StartPosition = helper.String(v.(string))
		}
		request.DynamicOptions = &dynamicOptions
	}

	request.JobId = &jobId
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := service.ModifyDtsSyncJobConfig(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update dts_sync_job_config failed, reason:%+v", logId, err)
		return err
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"Running", "Paused"}, tccommon.ReadRetryTimeout, time.Second, service.DtsSyncJobStateRefreshFunc(d.Id(), "", []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDtsSyncJobConfigRead(d, meta)
}

func resourceTencentCloudDtsSyncJobConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_sync_job_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
