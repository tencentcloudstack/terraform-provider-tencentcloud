/*
Provides a resource to create a dts migrate_job

Example Usage

```hcl
resource "tencentcloud_dts_migrate_job" "migrate_job" {
  src_database_type = "mysql"
  dst_database_type = "cynosdbmysql"
  src_region = "ap-guangzhou"
  dst_region = "ap-guangzhou"
  instance_class = "small"
  job_name = "tf_test_migration_job"
  tags {
	tag_key = "aaa"
	tag_value = "bbb"
  }
}

```
Import

dts migrate_job can be imported using the id, e.g.
```
$ terraform import tencentcloud_dts_migrate_job.migrate_job migrateJob_id
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
		Read:   resourceTencentCloudDtsMigrateJobRead,
		Create: resourceTencentCloudDtsMigrateJobCreate,
		Update: resourceTencentCloudDtsMigrateJobUpdate,
		Delete: resourceTencentCloudDtsMigrateJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"src_database_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "source database type, optional value is mysql/redis/percona/mongodb/postgresql/sqlserver/mariadb.",
			},

			"dst_database_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "destination database type, optional value is mysql/redis/percona/mongodb/postgresql/sqlserver/mariadb.",
			},

			"src_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "source region.",
			},

			"dst_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "destination region.",
			},

			"instance_class": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance class, optional value is small/medium/large/xlarge/2xlarge.",
			},

			"job_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "job name.",
			},

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "tags.",
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
			}, //tags

			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "job id.",
			},

			"run_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "run mode.",
			},

			"migrate_option": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "migrate option.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_table": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "database table.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object_mode": {
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
												"d_b_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "database mode.",
												},
												"schema_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "schema mode.",
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
															"tmp_tables": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Computed:    true,
																Description: "temporary tables.",
															},
															"table_edit_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "table edit mode.",
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
													Description: "views.",
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
												"role_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "role mode.",
												},
												"roles": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "role list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"role_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "role name.",
															},
															"new_role_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "new role name.",
															},
														},
													},
												},
												"function_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "function mode.",
												},
												"trigger_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "trigger mode.",
												},
												"event_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "event mode.",
												},
												"procedure_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "procedure mode.",
												},
												"functions": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "function list.",
												},
												"procedures": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "procedure list.",
												},
												"events": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "event list.",
												},
												"triggers": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "trigger list.",
												},
											},
										},
									},
								},
							},
						},
						"migrate_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "migrate type.",
						},
						"consistency": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "consistency option.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "mode, optional value is full/noCheck/notConfigure.",
									},
								},
							},
						},
						"is_migrate_account": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "migrate account.",
						},
						"is_override_root": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "override root destination by source database.",
						},
						"is_dst_read_only": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "destination readonly set.",
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
							Description: "databse info list.",
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
							Description: "databse info list.",
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

			"expect_run_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "expected run time, such as 2006-01-02 15:04:05.",
			},
		},
	}
}

func resourceTencentCloudDtsMigrateJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		request  = dts.NewCreateMigrationServiceRequest()
		response *dts.CreateMigrationServiceResponse
		ctx      = context.WithValue(context.TODO(), logIdKey, logId)
		service  = DtsService{client: meta.(*TencentCloudClient).apiV3Conn}
		jobId    string
	)

	if v, ok := d.GetOk("src_database_type"); ok {
		request.SrcDatabaseType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_database_type"); ok {
		request.DstDatabaseType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("src_region"); ok {
		request.SrcRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_region"); ok {
		request.DstRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_class"); ok {
		request.InstanceClass = helper.String(v.(string))
	}

	if v, ok := d.GetOk("job_name"); ok {
		request.JobName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			tagItem := dts.TagItem{}
			if v, ok := dMap["tag_key"]; ok {
				tagItem.TagKey = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				tagItem.TagValue = helper.String(v.(string))
			}

			request.Tags = append(request.Tags, &tagItem)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().CreateMigrationService(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dts migrateJob failed, reason:%+v", logId, err)
		return err
	}
	if response.Response == nil || response.Response.JobIds == nil {
		return fmt.Errorf("%s create dts migrateJob failed, response is nil!", logId)
	}

	jobId = *response.Response.JobIds[0]
	// wait created
	if err = service.PollingMigrateJobStatusUntil(ctx, jobId, DTSJobStatus, "created"); err != nil {
		return err
	}

	d.SetId(jobId)
	return resourceTencentCloudDtsMigrateJobRead(d, meta)
}

func resourceTencentCloudDtsMigrateJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = DtsService{client: meta.(*TencentCloudClient).apiV3Conn}
		jobId   = d.Id()
	)

	migrateJob, err := service.DescribeDtsMigrateJob(ctx, jobId)

	if err != nil {
		return err
	}

	if migrateJob == nil {
		d.SetId("")
		return fmt.Errorf("resource `migrateJob` %s does not exist", jobId)
	}

	if migrateJob.SrcInfo != nil {
		srcInfo := migrateJob.SrcInfo
		if srcInfo.DatabaseType != nil {
			_ = d.Set("src_database_type", srcInfo.DatabaseType)
		}

		if srcInfo.Region != nil {
			_ = d.Set("src_region", srcInfo.Region)
		}
	}

	if migrateJob.DstInfo != nil {
		destInfo := migrateJob.DstInfo
		if destInfo.DatabaseType != nil {
			_ = d.Set("dst_database_type", destInfo.DatabaseType)
		}

		if destInfo.Region != nil {
			_ = d.Set("dst_region", destInfo.Region)
		}
	}

	if migrateJob.TradeInfo != nil && migrateJob.TradeInfo.InstanceClass != nil {
		_ = d.Set("instance_class", migrateJob.TradeInfo.InstanceClass)
	}

	if migrateJob.JobName != nil {
		_ = d.Set("job_name", migrateJob.JobName)
	}

	if migrateJob.Tags != nil {
		tagsList := []interface{}{}
		for _, tags := range migrateJob.Tags {
			tagsMap := map[string]interface{}{}
			if tags.TagKey != nil {
				tagsMap["tag_key"] = tags.TagKey
			}
			if tags.TagValue != nil {
				tagsMap["tag_value"] = tags.TagValue
			}

			tagsList = append(tagsList, tagsMap)
		}
		_ = d.Set("tags", tagsList)
	}

	if migrateJob.JobId != nil {
		_ = d.Set("job_id", migrateJob.JobId)
	}

	if migrateJob.RunMode != nil {
		_ = d.Set("run_mode", migrateJob.RunMode)
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
						databasesMap["roles"] = rolesList
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
				databaseTableMap["databases"] = databasesList
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
			migrateOptionMap["extra_attr"] = extraAttrList
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
			srcInfoMap["extra_attr"] = extraAttrList
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
			dstInfoMap["extra_attr"] = extraAttrList
		}

		_ = d.Set("dst_info", []interface{}{dstInfoMap})
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

	return resourceTencentCloudDtsMigrateJobRead(d, meta)
}

func resourceTencentCloudDtsMigrateJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	migrateJobId := d.Id()

	if err := service.DeleteDtsMigrateJobById(ctx, migrateJobId); err != nil {
		return err
	}

	return nil
}
