/*
Provides a resource to create a dts migrate_job

Example Usage

```hcl
resource "tencentcloud_dts_migrate_job" "migrate_job" {
  job_id = ""
  compare_task_id = ""
  run_mode = ""
  resume_option = ""
  status = ""
  migrate_option {
		database_table {
			object_mode = ""
			databases {
				db_name = ""
				new_db_name = ""
				schema_name = ""
				new_schema_name = ""
				d_b_mode = ""
				schema_mode = ""
				table_mode = ""
				tables {
					table_name = &lt;nil&gt;
					new_table_name = &lt;nil&gt;
					tmp_tables = &lt;nil&gt;
					table_edit_mode = &lt;nil&gt;
				}
				view_mode = ""
				views {
					view_name = ""
					new_view_name = ""
				}
				role_mode = ""
				roles {
					role_name = ""
					new_role_name = ""
				}
				function_mode = ""
				trigger_mode = ""
				event_mode = ""
				procedure_mode = ""
				functions =
				procedures =
				events =
				triggers =
			}
			advanced_objects =
		}
		migrate_type = ""
		consistency {
			mode = ""
		}
		is_migrate_account =
		is_override_root =
		is_dst_read_only =
		extra_attr {
			key = ""
			value = ""
		}

  }
  src_info {
		region = ""
		access_type = ""
		database_type = ""
		node_type = ""
		info {
			role = ""
			db_kernel = ""
			host = ""
			port =
			user = ""
			password = ""
			cvm_instance_id = ""
			uniq_vpn_gw_id = ""
			uniq_dcg_id = ""
			instance_id = ""
			ccn_gw_id = ""
			vpc_id = ""
			subnet_id = ""
			engine_version = ""
			account = ""
			account_role = ""
			account_mode = ""
			tmp_secret_id = ""
			tmp_secret_key = ""
			tmp_token = ""
		}
		supplier = ""
		extra_attr {
			key = ""
			value = ""
		}

  }
  dst_info {
		region = ""
		access_type = ""
		database_type = ""
		node_type = ""
		info {
			role = ""
			db_kernel = ""
			host = ""
			port =
			user = ""
			password = ""
			cvm_instance_id = ""
			uniq_vpn_gw_id = ""
			uniq_dcg_id = ""
			instance_id = ""
			ccn_gw_id = ""
			vpc_id = ""
			subnet_id = ""
			engine_version = ""
			account = ""
			account_role = ""
			account_mode = ""
			tmp_secret_id = ""
			tmp_secret_key = ""
			tmp_token = ""
		}
		supplier = ""
		extra_attr {
			key = ""
			value = ""
		}

  }
  job_name = ""
  expect_run_time = ""
  auto_retry_time_range_minutes =
}
```

Import

dts migrate_job can be imported using the id, e.g.

```
terraform import tencentcloud_dts_migrate_job.migrate_job migrate_job_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDtsMigrateJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsMigrateJobCreate,
		Read:   resourceTencentCloudDtsMigrateJobRead,
		Update: resourceTencentCloudDtsMigrateJobUpdate,
		// Delete: resourceTencentCloudDtsMigrateJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Job Id.",
			},

			"compare_task_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Compare taskId.",
			},

			"status": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Change status to specify the migrate step. The valid values:config/startMigrate/resume/startCompare/complete/stop.",
			},

			"run_mode": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Run Mode. eg:immediate,timed.",
			},

			"resume_option": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The mode of the recovery task, the valid values: `clearData`: clears the target instance data. `overwrite`: executes the task in an overwriting way. `normal`: the normal process, no additional action is performed.",
			},

			"complete_mode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The way to complete the task, only support the old version of MySQL migration task, the valid values: waitForSync,immediately.",
			},

			"migrate_option": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Migration job configuration options, used to describe how the task performs migration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_table": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Migration object option, you need to tell the migration service which library table objects to migrate.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Object mode. eg:all,partial.",
									},
									"databases": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The database list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"db_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "database name.",
												},
												"new_db_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "New database name.",
												},
												"schema_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "schema name.",
												},
												"new_schema_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "schema name after migration or synchronization.",
												},
												"d_b_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "DB selection mode:all (for all objects under the current object), partial (partial objects), when the ObjectMode is partial, this item is required.",
												},
												"schema_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "schema mode: all,partial.",
												},
												"table_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "table mode: all,partial.",
												},
												"tables": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "tables list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"table_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "table name.",
															},
															"new_table_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "new table name.",
															},
															"tmp_tables": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Optional:    true,
																Description: "temporary tables.",
															},
															"table_edit_mode": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "table edit mode.",
															},
														},
													},
												},
												"view_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "ViewMode.",
												},
												"views": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Views.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"view_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "ViewName.",
															},
															"new_view_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "NewViewName.",
															},
														},
													},
												},
												"role_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "RoleMode.",
												},
												"roles": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Roles.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"role_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "RoleName.",
															},
															"new_role_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "NewRoleName.",
															},
														},
													},
												},
												"function_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "FunctionMode.",
												},
												"trigger_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "TriggerMode.",
												},
												"event_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "EventMode.",
												},
												"procedure_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "ProcedureMode.",
												},
												"functions": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Description: "Functions.",
												},
												"procedures": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Description: "Procedures.",
												},
												"events": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Description: "Events.",
												},
												"triggers": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Description: "Triggers.",
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
										Description: "AdvancedObjects.",
									},
								},
							},
						},
						"migrate_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "MigrateType.",
						},
						"consistency": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Consistency.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ConsistencyOption.",
									},
								},
							},
						},
						"is_migrate_account": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "IsMigrateAccount.",
						},
						"is_override_root": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "IsOverrideRoot.",
						},
						"is_dst_read_only": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "IsDstReadOnly.",
						},
						"extra_attr": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "ExtraAttr.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Value.",
									},
								},
							},
						},
					},
				},
			},

			"src_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "SrcInfo.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region.",
						},
						"access_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "AccessType.",
						},
						"database_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "DatabaseType.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "NodeType.",
						},
						"info": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Role.",
									},
									"db_kernel": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "DbKernel.",
									},
									"host": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Host.",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Port.",
									},
									"user": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "User.",
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Password.",
									},
									"cvm_instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "CvmInstanceId.",
									},
									"uniq_vpn_gw_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "UniqVpnGwId.",
									},
									"uniq_dcg_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "UniqDcgId.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "InstanceId.",
									},
									"ccn_gw_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "CcnGwId.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "VpcId.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "SubnetId.",
									},
									"engine_version": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "EngineVersion.",
									},
									"account": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Account.",
									},
									"account_role": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "AccountRole.",
									},
									"account_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "AccountMode.",
									},
									"tmp_secret_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TmpSecretId.",
									},
									"tmp_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TmpSecretKey.",
									},
									"tmp_token": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TmpToken.",
									},
								},
							},
						},
						"supplier": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Supplier.",
						},
						"extra_attr": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "ExtraAttr.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Value.",
									},
								},
							},
						},
					},
				},
			},

			"dst_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "DstInfo.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region.",
						},
						"access_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "AccessType.",
						},
						"database_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "DatabaseType.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "NodeType.",
						},
						"info": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Role.",
									},
									"db_kernel": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "DbKernel.",
									},
									"host": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Host.",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Port.",
									},
									"user": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "User.",
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Password.",
									},
									"cvm_instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "CvmInstanceId.",
									},
									"uniq_vpn_gw_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "UniqVpnGwId.",
									},
									"uniq_dcg_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "UniqDcgId.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "InstanceId.",
									},
									"ccn_gw_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "CcnGwId.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "VpcId.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "SubnetId.",
									},
									"engine_version": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Engine Version.",
									},
									"account": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Account.",
									},
									"account_role": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Account Role.",
									},
									"account_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Account Mode.",
									},
									"tmp_secret_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Tmp SecretId.",
									},
									"tmp_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Tmp SecretKey.",
									},
									"tmp_token": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Tmp Token.",
									},
								},
							},
						},
						"supplier": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Supplier.",
						},
						"extra_attr": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "ExtraAttr.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Value.",
									},
								},
							},
						},
					},
				},
			},

			"job_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "JobName.",
			},

			"expect_run_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ExpectRunTime.",
			},

			"auto_retry_time_range_minutes": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "AutoRetryTimeRangeMinutes.",
			},
		},
	}
}

func resourceTencentCloudDtsMigrateJobCreate2(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_service.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	// ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := dts.NewModifyMigrationJobRequest()
	request.JobId = helper.String(d.Id())

	if d.HasChange("src_database_type") {
		return fmt.Errorf("`src_database_type` do not support change now.")
	}

	if d.HasChange("dst_database_type") {
		return fmt.Errorf("`dst_database_type` do not support change now.")
	}

	if d.HasChange("src_region") {
		return fmt.Errorf("`src_region` do not support change now.")
	}

	if d.HasChange("dst_region") {
		return fmt.Errorf("`dst_region` do not support change now.")
	}

	if d.HasChange("instance_class") {
		return fmt.Errorf("`instance_class` do not support change now.")
	}

	if d.HasChange("job_name") {
		if v, ok := d.GetOk("job_name"); ok {
			request.JobName = helper.String(v.(string))
		}
	}

	if d.HasChange("tags") {
		return fmt.Errorf("`tags` do not support change now.")
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().ModifyMigrationJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update dts migrateJob failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDtsMigrateServiceRead(d, meta)
}

func resourceTencentCloudDtsMigrateJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = dts.NewModifyMigrationJobRequest()
		// response = dts.NewModifyMigrationJobResponse()
		jobId string
	)
	if v, ok := d.GetOk("job_id"); ok {
		request.JobId = helper.String(v.(string))
		jobId = *request.JobId
	}

	if v, ok := d.GetOk("run_mode"); ok {
		request.RunMode = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "migrate_option"); ok {
		migrateOption := dts.MigrateOption{}
		if databaseTableMap, ok := helper.InterfaceToMap(dMap, "database_table"); ok {
			databaseTableObject := dts.DatabaseTableObject{}
			if v, ok := databaseTableMap["object_mode"]; ok {
				databaseTableObject.ObjectMode = helper.String(v.(string))
			}
			if v, ok := databaseTableMap["databases"]; ok {
				for _, item := range v.([]interface{}) {
					databasesMap := item.(map[string]interface{})
					dBItem := dts.DBItem{}
					if v, ok := databasesMap["db_name"]; ok {
						dBItem.DbName = helper.String(v.(string))
					}
					if v, ok := databasesMap["new_db_name"]; ok {
						dBItem.NewDbName = helper.String(v.(string))
					}
					if v, ok := databasesMap["schema_name"]; ok {
						dBItem.SchemaName = helper.String(v.(string))
					}
					if v, ok := databasesMap["new_schema_name"]; ok {
						dBItem.NewSchemaName = helper.String(v.(string))
					}
					if v, ok := databasesMap["d_b_mode"]; ok {
						dBItem.DBMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["schema_mode"]; ok {
						dBItem.SchemaMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["table_mode"]; ok {
						dBItem.TableMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["tables"]; ok {
						for _, item := range v.([]interface{}) {
							tablesMap := item.(map[string]interface{})
							tableItem := dts.TableItem{}
							if v, ok := tablesMap["table_name"]; ok {
								tableItem.TableName = helper.String(v.(string))
							}
							if v, ok := tablesMap["new_table_name"]; ok {
								tableItem.NewTableName = helper.String(v.(string))
							}
							if v, ok := tablesMap["tmp_tables"]; ok {
								tmpTablesSet := v.(*schema.Set).List()
								for i := range tmpTablesSet {
									tmpTables := tmpTablesSet[i].(string)
									tableItem.TmpTables = append(tableItem.TmpTables, &tmpTables)
								}
							}
							if v, ok := tablesMap["table_edit_mode"]; ok {
								tableItem.TableEditMode = helper.String(v.(string))
							}
							dBItem.Tables = append(dBItem.Tables, &tableItem)
						}
					}
					if v, ok := databasesMap["view_mode"]; ok {
						dBItem.ViewMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["views"]; ok {
						for _, item := range v.([]interface{}) {
							viewsMap := item.(map[string]interface{})
							viewItem := dts.ViewItem{}
							if v, ok := viewsMap["view_name"]; ok {
								viewItem.ViewName = helper.String(v.(string))
							}
							if v, ok := viewsMap["new_view_name"]; ok {
								viewItem.NewViewName = helper.String(v.(string))
							}
							dBItem.Views = append(dBItem.Views, &viewItem)
						}
					}
					if v, ok := databasesMap["role_mode"]; ok {
						dBItem.RoleMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["roles"]; ok {
						for _, item := range v.([]interface{}) {
							rolesMap := item.(map[string]interface{})
							roleItem := dts.RoleItem{}
							if v, ok := rolesMap["role_name"]; ok {
								roleItem.RoleName = helper.String(v.(string))
							}
							if v, ok := rolesMap["new_role_name"]; ok {
								roleItem.NewRoleName = helper.String(v.(string))
							}
							dBItem.Roles = append(dBItem.Roles, &roleItem)
						}
					}
					if v, ok := databasesMap["function_mode"]; ok {
						dBItem.FunctionMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["trigger_mode"]; ok {
						dBItem.TriggerMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["event_mode"]; ok {
						dBItem.EventMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["procedure_mode"]; ok {
						dBItem.ProcedureMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["functions"]; ok {
						functionsSet := v.(*schema.Set).List()
						for i := range functionsSet {
							functions := functionsSet[i].(string)
							dBItem.Functions = append(dBItem.Functions, &functions)
						}
					}
					if v, ok := databasesMap["procedures"]; ok {
						proceduresSet := v.(*schema.Set).List()
						for i := range proceduresSet {
							procedures := proceduresSet[i].(string)
							dBItem.Procedures = append(dBItem.Procedures, &procedures)
						}
					}
					if v, ok := databasesMap["events"]; ok {
						eventsSet := v.(*schema.Set).List()
						for i := range eventsSet {
							events := eventsSet[i].(string)
							dBItem.Events = append(dBItem.Events, &events)
						}
					}
					if v, ok := databasesMap["triggers"]; ok {
						triggersSet := v.(*schema.Set).List()
						for i := range triggersSet {
							triggers := triggersSet[i].(string)
							dBItem.Triggers = append(dBItem.Triggers, &triggers)
						}
					}
					databaseTableObject.Databases = append(databaseTableObject.Databases, &dBItem)
				}
			}
			if v, ok := databaseTableMap["advanced_objects"]; ok {
				advancedObjectsSet := v.(*schema.Set).List()
				for i := range advancedObjectsSet {
					advancedObjects := advancedObjectsSet[i].(string)
					databaseTableObject.AdvancedObjects = append(databaseTableObject.AdvancedObjects, &advancedObjects)
				}
			}
			migrateOption.DatabaseTable = &databaseTableObject
		}
		if v, ok := dMap["migrate_type"]; ok {
			migrateOption.MigrateType = helper.String(v.(string))
		}
		if consistencyMap, ok := helper.InterfaceToMap(dMap, "consistency"); ok {
			consistencyOption := dts.ConsistencyOption{}
			if v, ok := consistencyMap["mode"]; ok {
				consistencyOption.Mode = helper.String(v.(string))
			}
			migrateOption.Consistency = &consistencyOption
		}
		if v, ok := dMap["is_migrate_account"]; ok {
			migrateOption.IsMigrateAccount = helper.Bool(v.(bool))
		}
		if v, ok := dMap["is_override_root"]; ok {
			migrateOption.IsOverrideRoot = helper.Bool(v.(bool))
		}
		if v, ok := dMap["is_dst_read_only"]; ok {
			migrateOption.IsDstReadOnly = helper.Bool(v.(bool))
		}
		if v, ok := dMap["extra_attr"]; ok {
			for _, item := range v.([]interface{}) {
				extraAttrMap := item.(map[string]interface{})
				keyValuePairOption := dts.KeyValuePairOption{}
				if v, ok := extraAttrMap["key"]; ok {
					keyValuePairOption.Key = helper.String(v.(string))
				}
				if v, ok := extraAttrMap["value"]; ok {
					keyValuePairOption.Value = helper.String(v.(string))
				}
				migrateOption.ExtraAttr = append(migrateOption.ExtraAttr, &keyValuePairOption)
			}
		}
		request.MigrateOption = &migrateOption
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "src_info"); ok {
		dBEndpointInfo := dts.DBEndpointInfo{}
		if v, ok := dMap["region"]; ok {
			dBEndpointInfo.Region = helper.String(v.(string))
		}
		if v, ok := dMap["access_type"]; ok {
			dBEndpointInfo.AccessType = helper.String(v.(string))
		}
		if v, ok := dMap["database_type"]; ok {
			dBEndpointInfo.DatabaseType = helper.String(v.(string))
		}
		if v, ok := dMap["node_type"]; ok {
			dBEndpointInfo.NodeType = helper.String(v.(string))
		}
		if v, ok := dMap["info"]; ok {
			for _, item := range v.([]interface{}) {
				infoMap := item.(map[string]interface{})
				dBInfo := dts.DBInfo{}
				if v, ok := infoMap["role"]; ok {
					dBInfo.Role = helper.String(v.(string))
				}
				if v, ok := infoMap["db_kernel"]; ok {
					dBInfo.DbKernel = helper.String(v.(string))
				}
				if v, ok := infoMap["host"]; ok {
					dBInfo.Host = helper.String(v.(string))
				}
				if v, ok := infoMap["port"]; ok {
					dBInfo.Port = helper.IntUint64(v.(int))
				}
				if v, ok := infoMap["user"]; ok {
					dBInfo.User = helper.String(v.(string))
				}
				if v, ok := infoMap["password"]; ok {
					dBInfo.Password = helper.String(v.(string))
				}
				if v, ok := infoMap["cvm_instance_id"]; ok {
					dBInfo.CvmInstanceId = helper.String(v.(string))
				}
				if v, ok := infoMap["uniq_vpn_gw_id"]; ok {
					dBInfo.UniqVpnGwId = helper.String(v.(string))
				}
				if v, ok := infoMap["uniq_dcg_id"]; ok {
					dBInfo.UniqDcgId = helper.String(v.(string))
				}
				if v, ok := infoMap["instance_id"]; ok {
					dBInfo.InstanceId = helper.String(v.(string))
				}
				if v, ok := infoMap["ccn_gw_id"]; ok {
					dBInfo.CcnGwId = helper.String(v.(string))
				}
				if v, ok := infoMap["vpc_id"]; ok {
					dBInfo.VpcId = helper.String(v.(string))
				}
				if v, ok := infoMap["subnet_id"]; ok {
					dBInfo.SubnetId = helper.String(v.(string))
				}
				if v, ok := infoMap["engine_version"]; ok {
					dBInfo.EngineVersion = helper.String(v.(string))
				}
				if v, ok := infoMap["account"]; ok {
					dBInfo.Account = helper.String(v.(string))
				}
				if v, ok := infoMap["account_role"]; ok {
					dBInfo.AccountRole = helper.String(v.(string))
				}
				if v, ok := infoMap["account_mode"]; ok {
					dBInfo.AccountMode = helper.String(v.(string))
				}
				if v, ok := infoMap["tmp_secret_id"]; ok {
					dBInfo.TmpSecretId = helper.String(v.(string))
				}
				if v, ok := infoMap["tmp_secret_key"]; ok {
					dBInfo.TmpSecretKey = helper.String(v.(string))
				}
				if v, ok := infoMap["tmp_token"]; ok {
					dBInfo.TmpToken = helper.String(v.(string))
				}
				dBEndpointInfo.Info = append(dBEndpointInfo.Info, &dBInfo)
			}
		}
		if v, ok := dMap["supplier"]; ok {
			dBEndpointInfo.Supplier = helper.String(v.(string))
		}
		if v, ok := dMap["extra_attr"]; ok {
			for _, item := range v.([]interface{}) {
				extraAttrMap := item.(map[string]interface{})
				keyValuePairOption := dts.KeyValuePairOption{}
				if v, ok := extraAttrMap["key"]; ok {
					keyValuePairOption.Key = helper.String(v.(string))
				}
				if v, ok := extraAttrMap["value"]; ok {
					keyValuePairOption.Value = helper.String(v.(string))
				}
				dBEndpointInfo.ExtraAttr = append(dBEndpointInfo.ExtraAttr, &keyValuePairOption)
			}
		}
		request.SrcInfo = &dBEndpointInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "dst_info"); ok {
		dBEndpointInfo := dts.DBEndpointInfo{}
		if v, ok := dMap["region"]; ok {
			dBEndpointInfo.Region = helper.String(v.(string))
		}
		if v, ok := dMap["access_type"]; ok {
			dBEndpointInfo.AccessType = helper.String(v.(string))
		}
		if v, ok := dMap["database_type"]; ok {
			dBEndpointInfo.DatabaseType = helper.String(v.(string))
		}
		if v, ok := dMap["node_type"]; ok {
			dBEndpointInfo.NodeType = helper.String(v.(string))
		}
		if v, ok := dMap["info"]; ok {
			for _, item := range v.([]interface{}) {
				infoMap := item.(map[string]interface{})
				dBInfo := dts.DBInfo{}
				if v, ok := infoMap["role"]; ok {
					dBInfo.Role = helper.String(v.(string))
				}
				if v, ok := infoMap["db_kernel"]; ok {
					dBInfo.DbKernel = helper.String(v.(string))
				}
				if v, ok := infoMap["host"]; ok {
					dBInfo.Host = helper.String(v.(string))
				}
				if v, ok := infoMap["port"]; ok {
					dBInfo.Port = helper.IntUint64(v.(int))
				}
				if v, ok := infoMap["user"]; ok {
					dBInfo.User = helper.String(v.(string))
				}
				if v, ok := infoMap["password"]; ok {
					dBInfo.Password = helper.String(v.(string))
				}
				if v, ok := infoMap["cvm_instance_id"]; ok {
					dBInfo.CvmInstanceId = helper.String(v.(string))
				}
				if v, ok := infoMap["uniq_vpn_gw_id"]; ok {
					dBInfo.UniqVpnGwId = helper.String(v.(string))
				}
				if v, ok := infoMap["uniq_dcg_id"]; ok {
					dBInfo.UniqDcgId = helper.String(v.(string))
				}
				if v, ok := infoMap["instance_id"]; ok {
					dBInfo.InstanceId = helper.String(v.(string))
				}
				if v, ok := infoMap["ccn_gw_id"]; ok {
					dBInfo.CcnGwId = helper.String(v.(string))
				}
				if v, ok := infoMap["vpc_id"]; ok {
					dBInfo.VpcId = helper.String(v.(string))
				}
				if v, ok := infoMap["subnet_id"]; ok {
					dBInfo.SubnetId = helper.String(v.(string))
				}
				if v, ok := infoMap["engine_version"]; ok {
					dBInfo.EngineVersion = helper.String(v.(string))
				}
				if v, ok := infoMap["account"]; ok {
					dBInfo.Account = helper.String(v.(string))
				}
				if v, ok := infoMap["account_role"]; ok {
					dBInfo.AccountRole = helper.String(v.(string))
				}
				if v, ok := infoMap["account_mode"]; ok {
					dBInfo.AccountMode = helper.String(v.(string))
				}
				if v, ok := infoMap["tmp_secret_id"]; ok {
					dBInfo.TmpSecretId = helper.String(v.(string))
				}
				if v, ok := infoMap["tmp_secret_key"]; ok {
					dBInfo.TmpSecretKey = helper.String(v.(string))
				}
				if v, ok := infoMap["tmp_token"]; ok {
					dBInfo.TmpToken = helper.String(v.(string))
				}
				dBEndpointInfo.Info = append(dBEndpointInfo.Info, &dBInfo)
			}
		}
		if v, ok := dMap["supplier"]; ok {
			dBEndpointInfo.Supplier = helper.String(v.(string))
		}
		if v, ok := dMap["extra_attr"]; ok {
			for _, item := range v.([]interface{}) {
				extraAttrMap := item.(map[string]interface{})
				keyValuePairOption := dts.KeyValuePairOption{}
				if v, ok := extraAttrMap["key"]; ok {
					keyValuePairOption.Key = helper.String(v.(string))
				}
				if v, ok := extraAttrMap["value"]; ok {
					keyValuePairOption.Value = helper.String(v.(string))
				}
				dBEndpointInfo.ExtraAttr = append(dBEndpointInfo.ExtraAttr, &keyValuePairOption)
			}
		}
		request.DstInfo = &dBEndpointInfo
	}

	if v, ok := d.GetOk("job_name"); ok {
		request.JobName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("expect_run_time"); ok {
		request.ExpectRunTime = helper.String(v.(string))
	}

	if v, _ := d.GetOk("auto_retry_time_range_minutes"); v != nil {
		request.AutoRetryTimeRangeMinutes = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().ModifyMigrationJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dts migrateJob failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(jobId)

	return resourceTencentCloudDtsMigrateJobRead(d, meta)
}

func resourceTencentCloudDtsMigrateJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	jobId := d.Id()

	migrateJob, err := service.DescribeDtsMigrateJobById(ctx, jobId)
	if err != nil {
		return err
	}

	if migrateJob == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	if migrateJob.JobId != nil {
		_ = d.Set("job_id", migrateJob.JobId)
	}

	if migrateJob.RunMode != nil {
		_ = d.Set("run_mode", migrateJob.RunMode)
	}

	// if migrateJob.CompareTaskId != nil {
	// 	_ = d.Set("compare_task_id", migrateJob.CompareTaskId)
	// }

	// if migrateJob.ResumeOption != nil {
	// 	_ = d.Set("resume_option", migrateJob.ResumeOption)
	// }

	// if migrateJob.CompleteMode != nil {
	// 	_ = d.Set("complete_mode", migrateJob.CompleteMode)
	// }

	if migrateJob.Status != nil {
		_ = d.Set("status", migrateJob.Status)
	}

	if migrateJob.MigrateOption != nil {
		migrateOptionMap := map[string]interface{}{}

		if migrateJob.MigrateOption.DatabaseTable != nil {
			databaseTableMap := map[string]interface{}{}

			if migrateJob.MigrateOption.DatabaseTable.ObjectMode != nil {
				databaseTableMap["object_mode"] = migrateJob.MigrateOption.DatabaseTable.ObjectMode
			}

			if migrateJob.MigrateOption.DatabaseTable.Databases != nil {
				databasesList := []interface{}{}
				for _, databases := range migrateJob.MigrateOption.DatabaseTable.Databases {
					databasesMap := map[string]interface{}{}

					if databases.DbName != nil {
						databasesMap["db_name"] = databases.DbName
					}

					if databases.NewDbName != nil {
						databasesMap["new_db_name"] = databases.NewDbName
					}

					if databases.SchemaName != nil {
						databasesMap["schema_name"] = databases.SchemaName
					}

					if databases.NewSchemaName != nil {
						databasesMap["new_schema_name"] = databases.NewSchemaName
					}

					if databases.DBMode != nil {
						databasesMap["d_b_mode"] = databases.DBMode
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

							if tables.TmpTables != nil {
								tablesMap["tmp_tables"] = tables.TmpTables
							}

							if tables.TableEditMode != nil {
								tablesMap["table_edit_mode"] = tables.TableEditMode
							}

							tablesList = append(tablesList, tablesMap)
						}

						databasesMap["tables"] = []interface{}{tablesList}
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

						databasesMap["views"] = []interface{}{viewsList}
					}

					if databases.RoleMode != nil {
						databasesMap["role_mode"] = databases.RoleMode
					}

					if databases.Roles != nil {
						rolesList := []interface{}{}
						for _, roles := range databases.Roles {
							rolesMap := map[string]interface{}{}

							if roles.RoleName != nil {
								rolesMap["role_name"] = roles.RoleName
							}

							if roles.NewRoleName != nil {
								rolesMap["new_role_name"] = roles.NewRoleName
							}

							rolesList = append(rolesList, rolesMap)
						}

						databasesMap["roles"] = []interface{}{rolesList}
					}

					if databases.FunctionMode != nil {
						databasesMap["function_mode"] = databases.FunctionMode
					}

					if databases.TriggerMode != nil {
						databasesMap["trigger_mode"] = databases.TriggerMode
					}

					if databases.EventMode != nil {
						databasesMap["event_mode"] = databases.EventMode
					}

					if databases.ProcedureMode != nil {
						databasesMap["procedure_mode"] = databases.ProcedureMode
					}

					if databases.Functions != nil {
						databasesMap["functions"] = databases.Functions
					}

					if databases.Procedures != nil {
						databasesMap["procedures"] = databases.Procedures
					}

					if databases.Events != nil {
						databasesMap["events"] = databases.Events
					}

					if databases.Triggers != nil {
						databasesMap["triggers"] = databases.Triggers
					}

					databasesList = append(databasesList, databasesMap)
				}

				databaseTableMap["databases"] = []interface{}{databasesList}
			}

			if migrateJob.MigrateOption.DatabaseTable.AdvancedObjects != nil {
				databaseTableMap["advanced_objects"] = migrateJob.MigrateOption.DatabaseTable.AdvancedObjects
			}

			migrateOptionMap["database_table"] = []interface{}{databaseTableMap}
		}

		if migrateJob.MigrateOption.MigrateType != nil {
			migrateOptionMap["migrate_type"] = migrateJob.MigrateOption.MigrateType
		}

		if migrateJob.MigrateOption.Consistency != nil {
			consistencyMap := map[string]interface{}{}

			if migrateJob.MigrateOption.Consistency.Mode != nil {
				consistencyMap["mode"] = migrateJob.MigrateOption.Consistency.Mode
			}

			migrateOptionMap["consistency"] = []interface{}{consistencyMap}
		}

		if migrateJob.MigrateOption.IsMigrateAccount != nil {
			migrateOptionMap["is_migrate_account"] = migrateJob.MigrateOption.IsMigrateAccount
		}

		if migrateJob.MigrateOption.IsOverrideRoot != nil {
			migrateOptionMap["is_override_root"] = migrateJob.MigrateOption.IsOverrideRoot
		}

		if migrateJob.MigrateOption.IsDstReadOnly != nil {
			migrateOptionMap["is_dst_read_only"] = migrateJob.MigrateOption.IsDstReadOnly
		}

		if migrateJob.MigrateOption.ExtraAttr != nil {
			extraAttrList := []interface{}{}
			for _, extraAttr := range migrateJob.MigrateOption.ExtraAttr {
				extraAttrMap := map[string]interface{}{}

				if extraAttr.Key != nil {
					extraAttrMap["key"] = extraAttr.Key
				}

				if extraAttr.Value != nil {
					extraAttrMap["value"] = extraAttr.Value
				}

				extraAttrList = append(extraAttrList, extraAttrMap)
			}

			migrateOptionMap["extra_attr"] = []interface{}{extraAttrList}
		}

		_ = d.Set("migrate_option", []interface{}{migrateOptionMap})
	}

	if migrateJob.SrcInfo != nil {
		srcInfoMap := map[string]interface{}{}

		if migrateJob.SrcInfo.Region != nil {
			srcInfoMap["region"] = migrateJob.SrcInfo.Region
		}

		if migrateJob.SrcInfo.AccessType != nil {
			srcInfoMap["access_type"] = migrateJob.SrcInfo.AccessType
		}

		if migrateJob.SrcInfo.DatabaseType != nil {
			srcInfoMap["database_type"] = migrateJob.SrcInfo.DatabaseType
		}

		if migrateJob.SrcInfo.NodeType != nil {
			srcInfoMap["node_type"] = migrateJob.SrcInfo.NodeType
		}

		if migrateJob.SrcInfo.Info != nil {
			infoList := []interface{}{}
			for _, info := range migrateJob.SrcInfo.Info {
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

				if info.UniqDcgId != nil {
					infoMap["uniq_dcg_id"] = info.UniqDcgId
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

			srcInfoMap["info"] = []interface{}{infoList}
		}

		if migrateJob.SrcInfo.Supplier != nil {
			srcInfoMap["supplier"] = migrateJob.SrcInfo.Supplier
		}

		if migrateJob.SrcInfo.ExtraAttr != nil {
			extraAttrList := []interface{}{}
			for _, extraAttr := range migrateJob.SrcInfo.ExtraAttr {
				extraAttrMap := map[string]interface{}{}

				if extraAttr.Key != nil {
					extraAttrMap["key"] = extraAttr.Key
				}

				if extraAttr.Value != nil {
					extraAttrMap["value"] = extraAttr.Value
				}

				extraAttrList = append(extraAttrList, extraAttrMap)
			}

			srcInfoMap["extra_attr"] = []interface{}{extraAttrList}
		}

		_ = d.Set("src_info", []interface{}{srcInfoMap})
	}

	if migrateJob.DstInfo != nil {
		dstInfoMap := map[string]interface{}{}

		if migrateJob.DstInfo.Region != nil {
			dstInfoMap["region"] = migrateJob.DstInfo.Region
		}

		if migrateJob.DstInfo.AccessType != nil {
			dstInfoMap["access_type"] = migrateJob.DstInfo.AccessType
		}

		if migrateJob.DstInfo.DatabaseType != nil {
			dstInfoMap["database_type"] = migrateJob.DstInfo.DatabaseType
		}

		if migrateJob.DstInfo.NodeType != nil {
			dstInfoMap["node_type"] = migrateJob.DstInfo.NodeType
		}

		if migrateJob.DstInfo.Info != nil {
			infoList := []interface{}{}
			for _, info := range migrateJob.DstInfo.Info {
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

				if info.UniqDcgId != nil {
					infoMap["uniq_dcg_id"] = info.UniqDcgId
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

			dstInfoMap["info"] = []interface{}{infoList}
		}

		if migrateJob.DstInfo.Supplier != nil {
			dstInfoMap["supplier"] = migrateJob.DstInfo.Supplier
		}

		if migrateJob.DstInfo.ExtraAttr != nil {
			extraAttrList := []interface{}{}
			for _, extraAttr := range migrateJob.DstInfo.ExtraAttr {
				extraAttrMap := map[string]interface{}{}

				if extraAttr.Key != nil {
					extraAttrMap["key"] = extraAttr.Key
				}

				if extraAttr.Value != nil {
					extraAttrMap["value"] = extraAttr.Value
				}

				extraAttrList = append(extraAttrList, extraAttrMap)
			}

			dstInfoMap["extra_attr"] = []interface{}{extraAttrList}
		}

		_ = d.Set("dst_info", []interface{}{dstInfoMap})
	}

	if migrateJob.JobName != nil {
		_ = d.Set("job_name", migrateJob.JobName)
	}

	if migrateJob.ExpectRunTime != nil {
		_ = d.Set("expect_run_time", migrateJob.ExpectRunTime)
	}

	return nil
}

func resourceTencentCloudDtsMigrateJobUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		configMigrationJobRequest = dts.NewModifyMigrationJobRequest()
		startMigrateJobRequest    = dts.NewStartMigrateJobRequest()
		completeMigrateJobRequest = dts.NewCompleteMigrateJobRequest()
		resumeMigrateJobRequest   = dts.NewResumeMigrateJobRequest()
		stopMigrateJobRequest     = dts.NewStopMigrateJobRequest()
		startCompareRequest       = dts.NewStartCompareRequest()
		status                    *string
	)

	startMigrateJobRequest.JobId = helper.String(d.Id())
	if d.HasChange("job_id") {
		if v, ok := d.GetOk("job_id"); ok {
			configMigrationJobRequest.JobId = helper.String(v.(string))
			startMigrateJobRequest.JobId = helper.String(v.(string))
			completeMigrateJobRequest.JobId = helper.String(v.(string))
			resumeMigrateJobRequest.JobId = helper.String(v.(string))
			stopMigrateJobRequest.JobId = helper.String(v.(string))
			startCompareRequest.JobId = helper.String(v.(string))
		}
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			status = helper.String(v.(string))
		}
	}

	if status != nil {

	}

	if d.HasChange("run_mode") {
		if v, ok := d.GetOk("run_mode"); ok {
			configMigrationJobRequest.RunMode = helper.String(v.(string))
		}
	}

	if d.HasChange("resume_option") {
		if v, ok := d.GetOk("resume_option"); ok {
			resumeMigrateJobRequest.ResumeOption = helper.String(v.(string))
		}
	}

	if d.HasChange("complete_mode") {
		if v, ok := d.GetOk("complete_mode"); ok {
			completeMigrateJobRequest.CompleteMode = helper.String(v.(string))
		}
	}

	if d.HasChange("compare_task_id") {
		if v, ok := d.GetOk("compare_task_id"); ok {
			startCompareRequest.CompareTaskId = helper.String(v.(string))
		}
	}

	if d.HasChange("migrate_option") {
		if dMap, ok := helper.InterfacesHeadMap(d, "migrate_option"); ok {
			migrateOption := dts.MigrateOption{}
			if databaseTableMap, ok := helper.InterfaceToMap(dMap, "database_table"); ok {
				databaseTableObject := dts.DatabaseTableObject{}
				if v, ok := databaseTableMap["object_mode"]; ok {
					databaseTableObject.ObjectMode = helper.String(v.(string))
				}
				if v, ok := databaseTableMap["databases"]; ok {
					for _, item := range v.([]interface{}) {
						databasesMap := item.(map[string]interface{})
						dBItem := dts.DBItem{}
						if v, ok := databasesMap["db_name"]; ok {
							dBItem.DbName = helper.String(v.(string))
						}
						if v, ok := databasesMap["new_db_name"]; ok {
							dBItem.NewDbName = helper.String(v.(string))
						}
						if v, ok := databasesMap["schema_name"]; ok {
							dBItem.SchemaName = helper.String(v.(string))
						}
						if v, ok := databasesMap["new_schema_name"]; ok {
							dBItem.NewSchemaName = helper.String(v.(string))
						}
						if v, ok := databasesMap["d_b_mode"]; ok {
							dBItem.DBMode = helper.String(v.(string))
						}
						if v, ok := databasesMap["schema_mode"]; ok {
							dBItem.SchemaMode = helper.String(v.(string))
						}
						if v, ok := databasesMap["table_mode"]; ok {
							dBItem.TableMode = helper.String(v.(string))
						}
						if v, ok := databasesMap["tables"]; ok {
							for _, item := range v.([]interface{}) {
								tablesMap := item.(map[string]interface{})
								tableItem := dts.TableItem{}
								if v, ok := tablesMap["table_name"]; ok {
									tableItem.TableName = helper.String(v.(string))
								}
								if v, ok := tablesMap["new_table_name"]; ok {
									tableItem.NewTableName = helper.String(v.(string))
								}
								if v, ok := tablesMap["tmp_tables"]; ok {
									tmpTablesSet := v.(*schema.Set).List()
									for i := range tmpTablesSet {
										tmpTables := tmpTablesSet[i].(string)
										tableItem.TmpTables = append(tableItem.TmpTables, &tmpTables)
									}
								}
								if v, ok := tablesMap["table_edit_mode"]; ok {
									tableItem.TableEditMode = helper.String(v.(string))
								}
								dBItem.Tables = append(dBItem.Tables, &tableItem)
							}
						}
						if v, ok := databasesMap["view_mode"]; ok {
							dBItem.ViewMode = helper.String(v.(string))
						}
						if v, ok := databasesMap["views"]; ok {
							for _, item := range v.([]interface{}) {
								viewsMap := item.(map[string]interface{})
								viewItem := dts.ViewItem{}
								if v, ok := viewsMap["view_name"]; ok {
									viewItem.ViewName = helper.String(v.(string))
								}
								if v, ok := viewsMap["new_view_name"]; ok {
									viewItem.NewViewName = helper.String(v.(string))
								}
								dBItem.Views = append(dBItem.Views, &viewItem)
							}
						}
						if v, ok := databasesMap["role_mode"]; ok {
							dBItem.RoleMode = helper.String(v.(string))
						}
						if v, ok := databasesMap["roles"]; ok {
							for _, item := range v.([]interface{}) {
								rolesMap := item.(map[string]interface{})
								roleItem := dts.RoleItem{}
								if v, ok := rolesMap["role_name"]; ok {
									roleItem.RoleName = helper.String(v.(string))
								}
								if v, ok := rolesMap["new_role_name"]; ok {
									roleItem.NewRoleName = helper.String(v.(string))
								}
								dBItem.Roles = append(dBItem.Roles, &roleItem)
							}
						}
						if v, ok := databasesMap["function_mode"]; ok {
							dBItem.FunctionMode = helper.String(v.(string))
						}
						if v, ok := databasesMap["trigger_mode"]; ok {
							dBItem.TriggerMode = helper.String(v.(string))
						}
						if v, ok := databasesMap["event_mode"]; ok {
							dBItem.EventMode = helper.String(v.(string))
						}
						if v, ok := databasesMap["procedure_mode"]; ok {
							dBItem.ProcedureMode = helper.String(v.(string))
						}
						if v, ok := databasesMap["functions"]; ok {
							functionsSet := v.(*schema.Set).List()
							for i := range functionsSet {
								functions := functionsSet[i].(string)
								dBItem.Functions = append(dBItem.Functions, &functions)
							}
						}
						if v, ok := databasesMap["procedures"]; ok {
							proceduresSet := v.(*schema.Set).List()
							for i := range proceduresSet {
								procedures := proceduresSet[i].(string)
								dBItem.Procedures = append(dBItem.Procedures, &procedures)
							}
						}
						if v, ok := databasesMap["events"]; ok {
							eventsSet := v.(*schema.Set).List()
							for i := range eventsSet {
								events := eventsSet[i].(string)
								dBItem.Events = append(dBItem.Events, &events)
							}
						}
						if v, ok := databasesMap["triggers"]; ok {
							triggersSet := v.(*schema.Set).List()
							for i := range triggersSet {
								triggers := triggersSet[i].(string)
								dBItem.Triggers = append(dBItem.Triggers, &triggers)
							}
						}
						databaseTableObject.Databases = append(databaseTableObject.Databases, &dBItem)
					}
				}
				if v, ok := databaseTableMap["advanced_objects"]; ok {
					advancedObjectsSet := v.(*schema.Set).List()
					for i := range advancedObjectsSet {
						advancedObjects := advancedObjectsSet[i].(string)
						databaseTableObject.AdvancedObjects = append(databaseTableObject.AdvancedObjects, &advancedObjects)
					}
				}
				migrateOption.DatabaseTable = &databaseTableObject
			}
			if v, ok := dMap["migrate_type"]; ok {
				migrateOption.MigrateType = helper.String(v.(string))
			}
			if consistencyMap, ok := helper.InterfaceToMap(dMap, "consistency"); ok {
				consistencyOption := dts.ConsistencyOption{}
				if v, ok := consistencyMap["mode"]; ok {
					consistencyOption.Mode = helper.String(v.(string))
				}
				migrateOption.Consistency = &consistencyOption
			}
			if v, ok := dMap["is_migrate_account"]; ok {
				migrateOption.IsMigrateAccount = helper.Bool(v.(bool))
			}
			if v, ok := dMap["is_override_root"]; ok {
				migrateOption.IsOverrideRoot = helper.Bool(v.(bool))
			}
			if v, ok := dMap["is_dst_read_only"]; ok {
				migrateOption.IsDstReadOnly = helper.Bool(v.(bool))
			}
			if v, ok := dMap["extra_attr"]; ok {
				for _, item := range v.([]interface{}) {
					extraAttrMap := item.(map[string]interface{})
					keyValuePairOption := dts.KeyValuePairOption{}
					if v, ok := extraAttrMap["key"]; ok {
						keyValuePairOption.Key = helper.String(v.(string))
					}
					if v, ok := extraAttrMap["value"]; ok {
						keyValuePairOption.Value = helper.String(v.(string))
					}
					migrateOption.ExtraAttr = append(migrateOption.ExtraAttr, &keyValuePairOption)
				}
			}
			configMigrationJobRequest.MigrateOption = &migrateOption
		}
	}

	if d.HasChange("src_info") {
		if dMap, ok := helper.InterfacesHeadMap(d, "src_info"); ok {
			dBEndpointInfo := dts.DBEndpointInfo{}
			if v, ok := dMap["region"]; ok {
				dBEndpointInfo.Region = helper.String(v.(string))
			}
			if v, ok := dMap["access_type"]; ok {
				dBEndpointInfo.AccessType = helper.String(v.(string))
			}
			if v, ok := dMap["database_type"]; ok {
				dBEndpointInfo.DatabaseType = helper.String(v.(string))
			}
			if v, ok := dMap["node_type"]; ok {
				dBEndpointInfo.NodeType = helper.String(v.(string))
			}
			if v, ok := dMap["info"]; ok {
				for _, item := range v.([]interface{}) {
					infoMap := item.(map[string]interface{})
					dBInfo := dts.DBInfo{}
					if v, ok := infoMap["role"]; ok {
						dBInfo.Role = helper.String(v.(string))
					}
					if v, ok := infoMap["db_kernel"]; ok {
						dBInfo.DbKernel = helper.String(v.(string))
					}
					if v, ok := infoMap["host"]; ok {
						dBInfo.Host = helper.String(v.(string))
					}
					if v, ok := infoMap["port"]; ok {
						dBInfo.Port = helper.IntUint64(v.(int))
					}
					if v, ok := infoMap["user"]; ok {
						dBInfo.User = helper.String(v.(string))
					}
					if v, ok := infoMap["password"]; ok {
						dBInfo.Password = helper.String(v.(string))
					}
					if v, ok := infoMap["cvm_instance_id"]; ok {
						dBInfo.CvmInstanceId = helper.String(v.(string))
					}
					if v, ok := infoMap["uniq_vpn_gw_id"]; ok {
						dBInfo.UniqVpnGwId = helper.String(v.(string))
					}
					if v, ok := infoMap["uniq_dcg_id"]; ok {
						dBInfo.UniqDcgId = helper.String(v.(string))
					}
					if v, ok := infoMap["instance_id"]; ok {
						dBInfo.InstanceId = helper.String(v.(string))
					}
					if v, ok := infoMap["ccn_gw_id"]; ok {
						dBInfo.CcnGwId = helper.String(v.(string))
					}
					if v, ok := infoMap["vpc_id"]; ok {
						dBInfo.VpcId = helper.String(v.(string))
					}
					if v, ok := infoMap["subnet_id"]; ok {
						dBInfo.SubnetId = helper.String(v.(string))
					}
					if v, ok := infoMap["engine_version"]; ok {
						dBInfo.EngineVersion = helper.String(v.(string))
					}
					if v, ok := infoMap["account"]; ok {
						dBInfo.Account = helper.String(v.(string))
					}
					if v, ok := infoMap["account_role"]; ok {
						dBInfo.AccountRole = helper.String(v.(string))
					}
					if v, ok := infoMap["account_mode"]; ok {
						dBInfo.AccountMode = helper.String(v.(string))
					}
					if v, ok := infoMap["tmp_secret_id"]; ok {
						dBInfo.TmpSecretId = helper.String(v.(string))
					}
					if v, ok := infoMap["tmp_secret_key"]; ok {
						dBInfo.TmpSecretKey = helper.String(v.(string))
					}
					if v, ok := infoMap["tmp_token"]; ok {
						dBInfo.TmpToken = helper.String(v.(string))
					}
					dBEndpointInfo.Info = append(dBEndpointInfo.Info, &dBInfo)
				}
			}
			if v, ok := dMap["supplier"]; ok {
				dBEndpointInfo.Supplier = helper.String(v.(string))
			}
			if v, ok := dMap["extra_attr"]; ok {
				for _, item := range v.([]interface{}) {
					extraAttrMap := item.(map[string]interface{})
					keyValuePairOption := dts.KeyValuePairOption{}
					if v, ok := extraAttrMap["key"]; ok {
						keyValuePairOption.Key = helper.String(v.(string))
					}
					if v, ok := extraAttrMap["value"]; ok {
						keyValuePairOption.Value = helper.String(v.(string))
					}
					dBEndpointInfo.ExtraAttr = append(dBEndpointInfo.ExtraAttr, &keyValuePairOption)
				}
			}
			configMigrationJobRequest.SrcInfo = &dBEndpointInfo
		}
	}

	if d.HasChange("dst_info") {
		if dMap, ok := helper.InterfacesHeadMap(d, "dst_info"); ok {
			dBEndpointInfo := dts.DBEndpointInfo{}
			if v, ok := dMap["region"]; ok {
				dBEndpointInfo.Region = helper.String(v.(string))
			}
			if v, ok := dMap["access_type"]; ok {
				dBEndpointInfo.AccessType = helper.String(v.(string))
			}
			if v, ok := dMap["database_type"]; ok {
				dBEndpointInfo.DatabaseType = helper.String(v.(string))
			}
			if v, ok := dMap["node_type"]; ok {
				dBEndpointInfo.NodeType = helper.String(v.(string))
			}
			if v, ok := dMap["info"]; ok {
				for _, item := range v.([]interface{}) {
					infoMap := item.(map[string]interface{})
					dBInfo := dts.DBInfo{}
					if v, ok := infoMap["role"]; ok {
						dBInfo.Role = helper.String(v.(string))
					}
					if v, ok := infoMap["db_kernel"]; ok {
						dBInfo.DbKernel = helper.String(v.(string))
					}
					if v, ok := infoMap["host"]; ok {
						dBInfo.Host = helper.String(v.(string))
					}
					if v, ok := infoMap["port"]; ok {
						dBInfo.Port = helper.IntUint64(v.(int))
					}
					if v, ok := infoMap["user"]; ok {
						dBInfo.User = helper.String(v.(string))
					}
					if v, ok := infoMap["password"]; ok {
						dBInfo.Password = helper.String(v.(string))
					}
					if v, ok := infoMap["cvm_instance_id"]; ok {
						dBInfo.CvmInstanceId = helper.String(v.(string))
					}
					if v, ok := infoMap["uniq_vpn_gw_id"]; ok {
						dBInfo.UniqVpnGwId = helper.String(v.(string))
					}
					if v, ok := infoMap["uniq_dcg_id"]; ok {
						dBInfo.UniqDcgId = helper.String(v.(string))
					}
					if v, ok := infoMap["instance_id"]; ok {
						dBInfo.InstanceId = helper.String(v.(string))
					}
					if v, ok := infoMap["ccn_gw_id"]; ok {
						dBInfo.CcnGwId = helper.String(v.(string))
					}
					if v, ok := infoMap["vpc_id"]; ok {
						dBInfo.VpcId = helper.String(v.(string))
					}
					if v, ok := infoMap["subnet_id"]; ok {
						dBInfo.SubnetId = helper.String(v.(string))
					}
					if v, ok := infoMap["engine_version"]; ok {
						dBInfo.EngineVersion = helper.String(v.(string))
					}
					if v, ok := infoMap["account"]; ok {
						dBInfo.Account = helper.String(v.(string))
					}
					if v, ok := infoMap["account_role"]; ok {
						dBInfo.AccountRole = helper.String(v.(string))
					}
					if v, ok := infoMap["account_mode"]; ok {
						dBInfo.AccountMode = helper.String(v.(string))
					}
					if v, ok := infoMap["tmp_secret_id"]; ok {
						dBInfo.TmpSecretId = helper.String(v.(string))
					}
					if v, ok := infoMap["tmp_secret_key"]; ok {
						dBInfo.TmpSecretKey = helper.String(v.(string))
					}
					if v, ok := infoMap["tmp_token"]; ok {
						dBInfo.TmpToken = helper.String(v.(string))
					}
					dBEndpointInfo.Info = append(dBEndpointInfo.Info, &dBInfo)
				}
			}
			if v, ok := dMap["supplier"]; ok {
				dBEndpointInfo.Supplier = helper.String(v.(string))
			}
			if v, ok := dMap["extra_attr"]; ok {
				for _, item := range v.([]interface{}) {
					extraAttrMap := item.(map[string]interface{})
					keyValuePairOption := dts.KeyValuePairOption{}
					if v, ok := extraAttrMap["key"]; ok {
						keyValuePairOption.Key = helper.String(v.(string))
					}
					if v, ok := extraAttrMap["value"]; ok {
						keyValuePairOption.Value = helper.String(v.(string))
					}
					dBEndpointInfo.ExtraAttr = append(dBEndpointInfo.ExtraAttr, &keyValuePairOption)
				}
			}
			configMigrationJobRequest.DstInfo = &dBEndpointInfo
		}
	}

	if d.HasChange("job_name") {
		if v, ok := d.GetOk("job_name"); ok {
			configMigrationJobRequest.JobName = helper.String(v.(string))
		}
	}

	if d.HasChange("expect_run_time") {
		if v, ok := d.GetOk("expect_run_time"); ok {
			configMigrationJobRequest.ExpectRunTime = helper.String(v.(string))
		}
	}

	if d.HasChange("auto_retry_time_range_minutes") {
		if v, _ := d.GetOk("auto_retry_time_range_minutes"); v != nil {
			configMigrationJobRequest.AutoRetryTimeRangeMinutes = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().StartMigrateJob(startMigrateJobRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, startMigrateJobRequest.GetAction(), startMigrateJobRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dts migrateJob failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().ModifyMigrationJob(configMigrationJobRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, configMigrationJobRequest.GetAction(), configMigrationJobRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dts migrateJob failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().CompleteMigrateJob(completeMigrateJobRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, completeMigrateJobRequest.GetAction(), completeMigrateJobRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dts migrateJob failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().ResumeMigrateJob(resumeMigrateJobRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, resumeMigrateJobRequest.GetAction(), resumeMigrateJobRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dts migrateJob failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().StopMigrateJob(stopMigrateJobRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, stopMigrateJobRequest.GetAction(), stopMigrateJobRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dts migrateJob failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().StartCompare(startCompareRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, startCompareRequest.GetAction(), startCompareRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dts migrateJob failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDtsMigrateJobRead(d, meta)
}
