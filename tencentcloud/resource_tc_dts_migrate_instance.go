/*
Provides a resource to create a dts migrate_instance

Example Usage

```hcl
resource "tencentcloud_dts_migrate_instance" "migrate_instance" {
  src_database_type = &lt;nil&gt;
  dst_database_type = &lt;nil&gt;
  src_region = &lt;nil&gt;
  dst_region = &lt;nil&gt;
  instance_class = &lt;nil&gt;
  count = &lt;nil&gt;
  job_name = &lt;nil&gt;
  tags {
		tag_key = &lt;nil&gt;
		tag_value = &lt;nil&gt;

  }
              complete_mode = &lt;nil&gt;
}
```

Import

dts migrate_instance can be imported using the id, e.g.

```
terraform import tencentcloud_dts_migrate_instance.migrate_instance migrate_instance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudDtsMigrateInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsMigrateInstanceCreate,
		Read:   resourceTencentCloudDtsMigrateInstanceRead,
		Update: resourceTencentCloudDtsMigrateInstanceUpdate,
		Delete: resourceTencentCloudDtsMigrateInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"src_database_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Source database type, optional value is mysql/redis/percona/mongodb/postgresql/sqlserver/mariadb.",
			},

			"dst_database_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Destination database type, optional value is mysql/redis/percona/mongodb/postgresql/sqlserver/mariadb.",
			},

			"src_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Source region.",
			},

			"dst_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Destination region.",
			},

			"instance_class": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance class, optional value is micro/small/medium/large/xlarge/2xlarge.",
			},

			"count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Count.",
			},

			"job_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Job name.",
			},

			"tags": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"job_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Job id.",
			},

			"run_mode": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Run mode.",
			},

			"migrate_option": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Migrate option.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_table": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Database table.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Object mode.",
									},
									"databases": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Database list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"db_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Database name.",
												},
												"new_db_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "New database name.",
												},
												"schema_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Schema name.",
												},
												"new_schema_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "New schema name.",
												},
												"d_b_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Database mode.",
												},
												"schema_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Schema mode.",
												},
												"table_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Table mode.",
												},
												"tables": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Table list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"table_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Table name.",
															},
															"new_table_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "New table name.",
															},
															"tmp_tables": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Computed:    true,
																Description: "Temporary tables.",
															},
															"table_edit_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Table edit mode.",
															},
														},
													},
												},
												"view_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "View mode.",
												},
												"views": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Views.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"view_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "View name.",
															},
															"new_view_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "New view name.",
															},
														},
													},
												},
												"role_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Role mode.",
												},
												"roles": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Role list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"role_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Role name.",
															},
															"new_role_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "New role name.",
															},
														},
													},
												},
												"function_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Function mode.",
												},
												"trigger_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Trigger mode.",
												},
												"event_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Event mode.",
												},
												"procedure_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Procedure mode.",
												},
												"functions": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Function list.",
												},
												"procedures": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Procedure list.",
												},
												"events": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Event list.",
												},
												"triggers": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Trigger list.",
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
							Description: "Migrate type.",
						},
						"consistency": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Consistency option.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Mode, optional value is full/noCheck/notConfigure.",
									},
								},
							},
						},
						"is_migrate_account": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Migrate account.",
						},
						"is_override_root": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Override root destination by source database.",
						},
						"is_dst_read_only": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Destination readonly set.",
						},
						"extra_attr": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Extra attributes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Value.",
									},
								},
							},
						},
					},
				},
			},

			"src_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Source info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"access_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access type.",
						},
						"database_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database type.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node type.",
						},
						"info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Databse info list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Node role.",
									},
									"db_kernel": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database kernel.",
									},
									"host": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Host.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port.",
									},
									"user": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "User.",
									},
									"password": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password.",
									},
									"cvm_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cvm instance id.",
									},
									"uniq_vpn_gw_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Vpn gateway id.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance id.",
									},
									"ccn_gw_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Ccn gateway id.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Vpc id.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet id.",
									},
									"engine_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Enginer version.",
									},
									"account": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Account.",
									},
									"account_role": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Account role.",
									},
									"account_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Account mode.",
									},
									"tmp_secret_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Temporary secret id.",
									},
									"tmp_secret_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Temporary secret key.",
									},
									"tmp_token": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Temporary token.",
									},
								},
							},
						},
						"supplier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Supplier.",
						},
						"extra_attr": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Extra attributes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Value.",
									},
								},
							},
						},
					},
				},
			},

			"dst_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Destination info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"access_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access type.",
						},
						"database_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database type.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node type.",
						},
						"info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Databse info list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Node role.",
									},
									"d_b_kernel": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database kernel.",
									},
									"host": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Host.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port.",
									},
									"user": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "User.",
									},
									"password": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password.",
									},
									"cvm_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cvm instance id.",
									},
									"uniq_vpn_gw_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Vpn gateway id.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance id.",
									},
									"ccn_gw_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Ccn gateway id.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Vpc id.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet id.",
									},
									"engine_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Enginer version.",
									},
									"account": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Account.",
									},
									"account_role": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Account role.",
									},
									"account_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Account mode.",
									},
									"tmp_secret_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Temporary secret id.",
									},
									"tmp_secret_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Temporary secret key.",
									},
									"tmp_token": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Temporary token.",
									},
								},
							},
						},
						"supplier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Supplier.",
						},
						"extra_attr": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Extra attributes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Value.",
									},
								},
							},
						},
					},
				},
			},

			"expect_run_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Expected run time, such as 2006-01-02 15:04:05.",
			},

			"complete_mode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Complete mode, optional value is waitForSync or immediately.",
			},
		},
	}
}

func resourceTencentCloudDtsMigrateInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = dts.NewCreateMigrationServiceRequest()
		response = dts.NewCreateMigrationServiceResponse()
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

	if v, ok := d.GetOkExists("count"); ok {
		request.Count = helper.IntInt64(v.(int))
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

	if v, ok := d.GetOk("complete_mode"); ok {
		request.CompleteMode = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().CreateMigrationService(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dts migrateInstance failed, reason:%+v", logId, err)
		return err
	}

	jobId = *response.Response.JobId
	d.SetId(jobId)

	return resourceTencentCloudDtsMigrateInstanceRead(d, meta)
}

func resourceTencentCloudDtsMigrateInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	migrateInstanceId := d.Id()

	migrateInstance, err := service.DescribeDtsMigrateInstanceById(ctx, jobId)
	if err != nil {
		return err
	}

	if migrateInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DtsMigrateInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if migrateInstance.SrcDatabaseType != nil {
		_ = d.Set("src_database_type", migrateInstance.SrcDatabaseType)
	}

	if migrateInstance.DstDatabaseType != nil {
		_ = d.Set("dst_database_type", migrateInstance.DstDatabaseType)
	}

	if migrateInstance.SrcRegion != nil {
		_ = d.Set("src_region", migrateInstance.SrcRegion)
	}

	if migrateInstance.DstRegion != nil {
		_ = d.Set("dst_region", migrateInstance.DstRegion)
	}

	if migrateInstance.InstanceClass != nil {
		_ = d.Set("instance_class", migrateInstance.InstanceClass)
	}

	if migrateInstance.Count != nil {
		_ = d.Set("count", migrateInstance.Count)
	}

	if migrateInstance.JobName != nil {
		_ = d.Set("job_name", migrateInstance.JobName)
	}

	if migrateInstance.Tags != nil {
		tagsList := []interface{}{}
		for _, tags := range migrateInstance.Tags {
			tagsMap := map[string]interface{}{}

			if migrateInstance.Tags.TagKey != nil {
				tagsMap["tag_key"] = migrateInstance.Tags.TagKey
			}

			if migrateInstance.Tags.TagValue != nil {
				tagsMap["tag_value"] = migrateInstance.Tags.TagValue
			}

			tagsList = append(tagsList, tagsMap)
		}

		_ = d.Set("tags", tagsList)

	}

	if migrateInstance.JobId != nil {
		_ = d.Set("job_id", migrateInstance.JobId)
	}

	if migrateInstance.RunMode != nil {
		_ = d.Set("run_mode", migrateInstance.RunMode)
	}

	if migrateInstance.MigrateOption != nil {
		migrateOptionMap := map[string]interface{}{}

		if migrateInstance.MigrateOption.DatabaseTable != nil {
			databaseTableMap := map[string]interface{}{}

			if migrateInstance.MigrateOption.DatabaseTable.ObjectMode != nil {
				databaseTableMap["object_mode"] = migrateInstance.MigrateOption.DatabaseTable.ObjectMode
			}

			if migrateInstance.MigrateOption.DatabaseTable.Databases != nil {
				databasesList := []interface{}{}
				for _, databases := range migrateInstance.MigrateOption.DatabaseTable.Databases {
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

			migrateOptionMap["database_table"] = []interface{}{databaseTableMap}
		}

		if migrateInstance.MigrateOption.MigrateType != nil {
			migrateOptionMap["migrate_type"] = migrateInstance.MigrateOption.MigrateType
		}

		if migrateInstance.MigrateOption.Consistency != nil {
			consistencyMap := map[string]interface{}{}

			if migrateInstance.MigrateOption.Consistency.Mode != nil {
				consistencyMap["mode"] = migrateInstance.MigrateOption.Consistency.Mode
			}

			migrateOptionMap["consistency"] = []interface{}{consistencyMap}
		}

		if migrateInstance.MigrateOption.IsMigrateAccount != nil {
			migrateOptionMap["is_migrate_account"] = migrateInstance.MigrateOption.IsMigrateAccount
		}

		if migrateInstance.MigrateOption.IsOverrideRoot != nil {
			migrateOptionMap["is_override_root"] = migrateInstance.MigrateOption.IsOverrideRoot
		}

		if migrateInstance.MigrateOption.IsDstReadOnly != nil {
			migrateOptionMap["is_dst_read_only"] = migrateInstance.MigrateOption.IsDstReadOnly
		}

		if migrateInstance.MigrateOption.ExtraAttr != nil {
			extraAttrList := []interface{}{}
			for _, extraAttr := range migrateInstance.MigrateOption.ExtraAttr {
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

	if migrateInstance.SrcInfo != nil {
		srcInfoMap := map[string]interface{}{}

		if migrateInstance.SrcInfo.Region != nil {
			srcInfoMap["region"] = migrateInstance.SrcInfo.Region
		}

		if migrateInstance.SrcInfo.AccessType != nil {
			srcInfoMap["access_type"] = migrateInstance.SrcInfo.AccessType
		}

		if migrateInstance.SrcInfo.DatabaseType != nil {
			srcInfoMap["database_type"] = migrateInstance.SrcInfo.DatabaseType
		}

		if migrateInstance.SrcInfo.NodeType != nil {
			srcInfoMap["node_type"] = migrateInstance.SrcInfo.NodeType
		}

		if migrateInstance.SrcInfo.Info != nil {
			infoList := []interface{}{}
			for _, info := range migrateInstance.SrcInfo.Info {
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

			srcInfoMap["info"] = []interface{}{infoList}
		}

		if migrateInstance.SrcInfo.Supplier != nil {
			srcInfoMap["supplier"] = migrateInstance.SrcInfo.Supplier
		}

		if migrateInstance.SrcInfo.ExtraAttr != nil {
			extraAttrList := []interface{}{}
			for _, extraAttr := range migrateInstance.SrcInfo.ExtraAttr {
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

	if migrateInstance.DstInfo != nil {
		dstInfoMap := map[string]interface{}{}

		if migrateInstance.DstInfo.Region != nil {
			dstInfoMap["region"] = migrateInstance.DstInfo.Region
		}

		if migrateInstance.DstInfo.AccessType != nil {
			dstInfoMap["access_type"] = migrateInstance.DstInfo.AccessType
		}

		if migrateInstance.DstInfo.DatabaseType != nil {
			dstInfoMap["database_type"] = migrateInstance.DstInfo.DatabaseType
		}

		if migrateInstance.DstInfo.NodeType != nil {
			dstInfoMap["node_type"] = migrateInstance.DstInfo.NodeType
		}

		if migrateInstance.DstInfo.Info != nil {
			infoList := []interface{}{}
			for _, info := range migrateInstance.DstInfo.Info {
				infoMap := map[string]interface{}{}

				if info.Role != nil {
					infoMap["role"] = info.Role
				}

				if info.DBKernel != nil {
					infoMap["d_b_kernel"] = info.DBKernel
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

			dstInfoMap["info"] = []interface{}{infoList}
		}

		if migrateInstance.DstInfo.Supplier != nil {
			dstInfoMap["supplier"] = migrateInstance.DstInfo.Supplier
		}

		if migrateInstance.DstInfo.ExtraAttr != nil {
			extraAttrList := []interface{}{}
			for _, extraAttr := range migrateInstance.DstInfo.ExtraAttr {
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

	if migrateInstance.ExpectRunTime != nil {
		_ = d.Set("expect_run_time", migrateInstance.ExpectRunTime)
	}

	if migrateInstance.CompleteMode != nil {
		_ = d.Set("complete_mode", migrateInstance.CompleteMode)
	}

	return nil
}

func resourceTencentCloudDtsMigrateInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"src_database_type", "dst_database_type", "src_region", "dst_region", "instance_class", "count", "job_name", "tags", "job_id", "run_mode", "migrate_option", "src_info", "dst_info", "expect_run_time", "complete_mode"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudDtsMigrateInstanceRead(d, meta)
}

func resourceTencentCloudDtsMigrateInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}
	migrateInstanceId := d.Id()

	if err := service.DeleteDtsMigrateInstanceById(ctx, jobId); err != nil {
		return err
	}

	return nil
}
