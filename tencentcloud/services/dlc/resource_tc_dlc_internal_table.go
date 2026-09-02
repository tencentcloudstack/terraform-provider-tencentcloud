package dlc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcInternalTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcInternalTableCreate,
		Read:   resourceTencentCloudDlcInternalTableRead,
		Update: resourceTencentCloudDlcInternalTableUpdate,
		Delete: resourceTencentCloudDlcInternalTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"database_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the database to which the internal table belongs.",
			},

			"table_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the internal table.",
			},

			"datasource_connection_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the datasource connection to which the table belongs.",
			},

			"table_comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The comment of the table.",
			},

			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the table, table or view.",
			},

			"table_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The format of the table, such as hive, iceberg, etc.",
			},

			"user_alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The alias of the user who creates the table.",
			},

			"user_sub_uin": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The sub UIN of the user who creates the table.",
			},

			"db_govern_policy_is_disable": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether the database data governance is disabled. true: disabled, false: enabled. Deprecated, use smart_policy instead.",
			},

			"govern_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Data governance policy. Deprecated, use smart_policy instead.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Governance rule type. Customize: custom; Intelligence: intelligent governance.",
						},
						"govern_engine": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Governance engine.",
						},
					},
				},
			},

			"smart_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Smart data governance policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"base_info": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Base info of the smart policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uin": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "User UIN.",
									},
									"policy_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Policy type: Catalog/Database/Table.",
									},
									"catalog": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Catalog name.",
									},
									"database": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Database name.",
									},
									"table": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Table name.",
									},
									"app_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "User APP ID.",
									},
								},
							},
						},
						"policy": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Smart optimizer policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"inherit": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether to inherit.",
									},
									"resources": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Data governance resources.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"attribution_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Attribution type.",
												},
												"resource_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Resource type.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Engine name.",
												},
												"instance": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Instance name.",
												},
												"status": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Status.",
												},
												"resource_group_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Standard engine resource group name.",
												},
												"favor": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Affinity info list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"priority": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Priority.",
															},
															"catalog": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Catalog name.",
															},
															"data_base": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Database name.",
															},
															"table": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Table name.",
															},
														},
													},
												},
												"resource_conf": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Resource configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"parallelism": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The parallelism of the optimization task.",
															},
														},
													},
												},
											},
										},
									},
									"written": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Data rewrite policy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"written_enable": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "none/enable/disable/default.",
												},
												"advance_policy": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Advanced rewrite policy.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"compact_enable": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Whether to enable compaction.",
															},
															"delete_enable": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Whether to enable history data cleanup.",
															},
															"min_input_files": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The minimum number of files to trigger compaction.",
															},
															"target_file_size_bytes": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Target file size in bytes for compaction.",
															},
															"retain_last": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Number of snapshots to retain.",
															},
															"before_days": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Snapshot expiration time in days.",
															},
															"expired_snapshots_interval_min": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Snapshot expiration execution interval in minutes.",
															},
															"remove_orphan_interval_min": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Orphan file removal execution interval in minutes.",
															},
															"cow_compact_enable": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Whether to enable COW table compaction.",
															},
															"compact_strategy": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "File compaction strategy.",
															},
															"sort_orders": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Sort compaction strategy rules.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"column": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Column name for sorting.",
																		},
																		"sort_direction": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Sort direction: ascending or descending.",
																		},
																		"null_order": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Null order: at the beginning or end.",
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
									"lifecycle": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Data lifecycle policy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"lifecycle_enable": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Whether lifecycle is enabled.",
												},
												"expiration": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Expiration time.",
												},
												"drop_table": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to drop the table. Deprecated, use table_expiration instead.",
												},
												"expired_field": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Expired field.",
												},
												"expired_field_format": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Expired field format.",
												},
											},
										},
									},
									"index": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Index policy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"index_enable": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Whether to enable index.",
												},
											},
										},
									},
									"change_table": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Change table policy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"data_retention_time": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Data retention time in days for the change table.",
												},
											},
										},
									},
									"table_expiration": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Table expiration policy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Whether the policy is enabled.",
												},
												"expiration": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Table expiration time in days.",
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

			"primary_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Primary keys of the T-ICEBERG table.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"columns": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Column information of the table.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Column name.",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Column type.",
						},
						"comment": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Column comment.",
						},
						"default": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Column default value.",
						},
						"not_null": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the column is not null.",
						},
						"precision": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The precision of the numeric type.",
						},
						"scale": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The scale of the numeric type.",
						},
						"position": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Column position, smaller is earlier.",
						},
						"is_partition": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the column is a partition field.",
						},
						"nullable": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the column is nullable.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Column creation time.",
						},
						"modified_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Column modification time.",
						},
						"data_mask_strategy_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Data mask strategy info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"strategy_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Strategy name.",
									},
									"strategy_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Strategy type.",
									},
									"strategy_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Strategy description.",
									},
									"users": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "User sub UIN list separated by semicolons.",
									},
									"strategy_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Strategy ID.",
									},
								},
							},
						},
						"type_text": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data field description.",
						},
					},
				},
			},

			"partitions": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Partition information of the table.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Partition column name.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Partition type.",
						},
						"comment": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Partition comment.",
						},
						"transform": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Implicit partition transform strategy.",
						},
						"transform_args": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Transform strategy arguments.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Partition creation time.",
						},
					},
				},
			},

			"properties": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Property information of the table.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Property key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Property value.",
						},
					},
				},
			},

			"upsert_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Upsert keys for V2 upsert table.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Data storage path.",
			},

			"modified_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Table modification time, in ms.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Table creation time, in ms.",
			},

			"input_format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Data format.",
			},

			"storage_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Table storage size in bytes.",
			},

			"record_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Table row count.",
			},

			"map_materialized_view_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Materialized view name mapping.",
			},

			"heat_value": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Access heat value.",
			},

			"input_format_short": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Abbreviation of InputFormat.",
			},

			"execution": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The auto-generated SQL statement from create.",
			},

			"is_t_iceberg_sql": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether it is T-ICEBERG SQL.",
			},
		},
	}
}

func resourceTencentCloudDlcInternalTableCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_internal_table.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = dlc.NewGenerateInternalTableRequest()
		response = dlc.NewGenerateInternalTableResponse()
	)

	tableBaseInfo := &dlc.TableBaseInfo{}
	if v, ok := d.GetOk("database_name"); ok {
		tableBaseInfo.DatabaseName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("table_name"); ok {
		tableBaseInfo.TableName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("datasource_connection_name"); ok {
		tableBaseInfo.DatasourceConnectionName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("table_comment"); ok {
		tableBaseInfo.TableComment = helper.String(v.(string))
	}
	if v, ok := d.GetOk("type"); ok {
		tableBaseInfo.Type = helper.String(v.(string))
	}
	if v, ok := d.GetOk("table_format"); ok {
		tableBaseInfo.TableFormat = helper.String(v.(string))
	}
	if v, ok := d.GetOk("user_alias"); ok {
		tableBaseInfo.UserAlias = helper.String(v.(string))
	}
	if v, ok := d.GetOk("user_sub_uin"); ok {
		tableBaseInfo.UserSubUin = helper.String(v.(string))
	}
	if v, ok := d.GetOk("db_govern_policy_is_disable"); ok {
		tableBaseInfo.DbGovernPolicyIsDisable = helper.String(v.(string))
	}

	if v, ok := d.GetOk("govern_policy"); ok {
		for _, item := range v.([]interface{}) {
			governPolicyMap := item.(map[string]interface{})
			dataGovernPolicy := &dlc.DataGovernPolicy{}
			if v, ok := governPolicyMap["rule_type"].(string); ok && v != "" {
				dataGovernPolicy.RuleType = helper.String(v)
			}
			if v, ok := governPolicyMap["govern_engine"].(string); ok && v != "" {
				dataGovernPolicy.GovernEngine = helper.String(v)
			}
			tableBaseInfo.GovernPolicy = dataGovernPolicy
		}
	}

	if v, ok := d.GetOk("smart_policy"); ok {
		tableBaseInfo.SmartPolicy = buildSmartPolicy(v.([]interface{}))
	}

	if v, ok := d.GetOk("primary_keys"); ok {
		primaryKeys := make([]*string, 0, len(v.([]interface{})))
		for _, item := range v.([]interface{}) {
			primaryKeys = append(primaryKeys, helper.String(item.(string)))
		}
		tableBaseInfo.PrimaryKeys = primaryKeys
	}

	request.TableBaseInfo = tableBaseInfo

	if v, ok := d.GetOk("columns"); ok {
		for _, item := range v.([]interface{}) {
			columnMap := item.(map[string]interface{})
			tColumn := &dlc.TColumn{}
			if v, ok := columnMap["name"].(string); ok && v != "" {
				tColumn.Name = helper.String(v)
			}
			if v, ok := columnMap["type"].(string); ok && v != "" {
				tColumn.Type = helper.String(v)
			}
			if v, ok := columnMap["comment"].(string); ok && v != "" {
				tColumn.Comment = helper.String(v)
			}
			if v, ok := columnMap["default"].(string); ok && v != "" {
				tColumn.Default = helper.String(v)
			}
			if v, ok := columnMap["not_null"].(bool); ok {
				tColumn.NotNull = helper.Bool(v)
			}
			if v, ok := columnMap["precision"].(int); ok && v != 0 {
				tColumn.Precision = helper.Int64(int64(v))
			}
			if v, ok := columnMap["scale"].(int); ok && v != 0 {
				tColumn.Scale = helper.Int64(int64(v))
			}
			if v, ok := columnMap["position"].(int); ok && v != 0 {
				tColumn.Position = helper.Int64(int64(v))
			}
			if v, ok := columnMap["is_partition"].(bool); ok {
				tColumn.IsPartition = helper.Bool(v)
			}
			request.Columns = append(request.Columns, tColumn)
		}
	}

	if v, ok := d.GetOk("partitions"); ok {
		for _, item := range v.([]interface{}) {
			partitionMap := item.(map[string]interface{})
			tPartition := &dlc.TPartition{}
			if v, ok := partitionMap["name"].(string); ok && v != "" {
				tPartition.Name = helper.String(v)
			}
			if v, ok := partitionMap["type"].(string); ok && v != "" {
				tPartition.Type = helper.String(v)
			}
			if v, ok := partitionMap["comment"].(string); ok && v != "" {
				tPartition.Comment = helper.String(v)
			}
			if v, ok := partitionMap["transform"].(string); ok && v != "" {
				tPartition.Transform = helper.String(v)
			}
			if v, ok := partitionMap["transform_args"].([]interface{}); ok {
				transformArgs := make([]*string, 0, len(v))
				for _, arg := range v {
					transformArgs = append(transformArgs, helper.String(arg.(string)))
				}
				tPartition.TransformArgs = transformArgs
			}
			request.Partitions = append(request.Partitions, tPartition)
		}
	}

	if v, ok := d.GetOk("properties"); ok {
		for _, item := range v.([]interface{}) {
			propertyMap := item.(map[string]interface{})
			property := &dlc.Property{}
			if v, ok := propertyMap["key"].(string); ok && v != "" {
				property.Key = helper.String(v)
			}
			if v, ok := propertyMap["value"].(string); ok && v != "" {
				property.Value = helper.String(v)
			}
			request.Properties = append(request.Properties, property)
		}
	}

	if v, ok := d.GetOk("upsert_keys"); ok {
		upsertKeys := make([]*string, 0, len(v.([]interface{})))
		for _, item := range v.([]interface{}) {
			upsertKeys = append(upsertKeys, helper.String(item.(string)))
		}
		request.UpsertKeys = upsertKeys
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().GenerateInternalTableWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dlc_internal_table failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dlc_internal_table failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[CRUD]%s api[%s] dlc_internal_table response: %s, logId=%s, d.Id()=%s", logId, request.GetAction(), response.ToJsonString(), logId, d.Id())

	if response.Response.Execution == nil {
		return fmt.Errorf("dlc_internal_table create response Execution is nil")
	}

	databaseName := ""
	tableName := ""
	if v, ok := d.GetOk("database_name"); ok {
		databaseName = v.(string)
	}
	if v, ok := d.GetOk("table_name"); ok {
		tableName = v.(string)
	}
	d.SetId(strings.Join([]string{databaseName, tableName}, tccommon.FILED_SP))

	if response.Response.Execution != nil && response.Response.Execution.SQL != nil {
		_ = d.Set("execution", response.Response.Execution.SQL)
	}
	if response.Response.IsTIcebergSql != nil {
		_ = d.Set("is_t_iceberg_sql", response.Response.IsTIcebergSql)
	}

	return resourceTencentCloudDlcInternalTableRead(d, meta)
}

func resourceTencentCloudDlcInternalTableRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_internal_table.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	databaseName := idSplit[0]
	tableName := idSplit[1]

	datasourceConnectionName := ""
	if v, ok := d.GetOk("datasource_connection_name"); ok {
		datasourceConnectionName = v.(string)
	}

	tableInfo, err := service.DescribeDlcInternalTableById(ctx, databaseName, tableName, datasourceConnectionName)
	if err != nil {
		return err
	}

	if tableInfo == nil {
		log.Printf("[CRUD] dlc_internal_table id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if tableInfo.TableBaseInfo != nil {
		if tableInfo.TableBaseInfo.DatabaseName != nil {
			_ = d.Set("database_name", tableInfo.TableBaseInfo.DatabaseName)
		}
		if tableInfo.TableBaseInfo.TableName != nil {
			_ = d.Set("table_name", tableInfo.TableBaseInfo.TableName)
		}
		if tableInfo.TableBaseInfo.DatasourceConnectionName != nil {
			_ = d.Set("datasource_connection_name", tableInfo.TableBaseInfo.DatasourceConnectionName)
		}
		if tableInfo.TableBaseInfo.TableComment != nil {
			_ = d.Set("table_comment", tableInfo.TableBaseInfo.TableComment)
		}
		if tableInfo.TableBaseInfo.Type != nil {
			_ = d.Set("type", tableInfo.TableBaseInfo.Type)
		}
		if tableInfo.TableBaseInfo.TableFormat != nil {
			_ = d.Set("table_format", tableInfo.TableBaseInfo.TableFormat)
		}
		if tableInfo.TableBaseInfo.UserAlias != nil {
			_ = d.Set("user_alias", tableInfo.TableBaseInfo.UserAlias)
		}
		if tableInfo.TableBaseInfo.UserSubUin != nil {
			_ = d.Set("user_sub_uin", tableInfo.TableBaseInfo.UserSubUin)
		}
		if tableInfo.TableBaseInfo.DbGovernPolicyIsDisable != nil {
			_ = d.Set("db_govern_policy_is_disable", tableInfo.TableBaseInfo.DbGovernPolicyIsDisable)
		}

		if tableInfo.TableBaseInfo.GovernPolicy != nil {
			governPolicyList := make([]map[string]interface{}, 0, 1)
			governPolicyMap := map[string]interface{}{}
			if tableInfo.TableBaseInfo.GovernPolicy.RuleType != nil {
				governPolicyMap["rule_type"] = tableInfo.TableBaseInfo.GovernPolicy.RuleType
			}
			if tableInfo.TableBaseInfo.GovernPolicy.GovernEngine != nil {
				governPolicyMap["govern_engine"] = tableInfo.TableBaseInfo.GovernPolicy.GovernEngine
			}
			governPolicyList = append(governPolicyList, governPolicyMap)
			_ = d.Set("govern_policy", governPolicyList)
		}

		if tableInfo.TableBaseInfo.SmartPolicy != nil {
			smartPolicyList := flattenSmartPolicy(tableInfo.TableBaseInfo.SmartPolicy)
			if len(smartPolicyList) > 0 {
				_ = d.Set("smart_policy", smartPolicyList)
			}
		}

		if tableInfo.TableBaseInfo.PrimaryKeys != nil {
			primaryKeys := make([]string, 0, len(tableInfo.TableBaseInfo.PrimaryKeys))
			for _, pk := range tableInfo.TableBaseInfo.PrimaryKeys {
				if pk != nil {
					primaryKeys = append(primaryKeys, *pk)
				}
			}
			_ = d.Set("primary_keys", primaryKeys)
		}
	}

	if tableInfo.Columns != nil {
		columnsList := make([]map[string]interface{}, 0, len(tableInfo.Columns))
		for _, column := range tableInfo.Columns {
			columnMap := map[string]interface{}{}
			if column.Name != nil {
				columnMap["name"] = column.Name
			}
			if column.Type != nil {
				columnMap["type"] = column.Type
			}
			if column.Comment != nil {
				columnMap["comment"] = column.Comment
			}
			if column.Precision != nil {
				columnMap["precision"] = column.Precision
			}
			if column.Scale != nil {
				columnMap["scale"] = column.Scale
			}
			if column.Nullable != nil {
				columnMap["nullable"] = column.Nullable
			}
			if column.Position != nil {
				columnMap["position"] = column.Position
			}
			if column.CreateTime != nil {
				columnMap["create_time"] = column.CreateTime
			}
			if column.ModifiedTime != nil {
				columnMap["modified_time"] = column.ModifiedTime
			}
			if column.IsPartition != nil {
				columnMap["is_partition"] = column.IsPartition
			}
			if column.TypeText != nil {
				columnMap["type_text"] = column.TypeText
			}
			if column.DataMaskStrategyInfo != nil {
				dataMaskList := make([]map[string]interface{}, 0, 1)
				dataMaskMap := map[string]interface{}{}
				if column.DataMaskStrategyInfo.StrategyName != nil {
					dataMaskMap["strategy_name"] = column.DataMaskStrategyInfo.StrategyName
				}
				if column.DataMaskStrategyInfo.StrategyType != nil {
					dataMaskMap["strategy_type"] = column.DataMaskStrategyInfo.StrategyType
				}
				if column.DataMaskStrategyInfo.StrategyDesc != nil {
					dataMaskMap["strategy_desc"] = column.DataMaskStrategyInfo.StrategyDesc
				}
				if column.DataMaskStrategyInfo.Users != nil {
					dataMaskMap["users"] = column.DataMaskStrategyInfo.Users
				}
				if column.DataMaskStrategyInfo.StrategyId != nil {
					dataMaskMap["strategy_id"] = column.DataMaskStrategyInfo.StrategyId
				}
				dataMaskList = append(dataMaskList, dataMaskMap)
				columnMap["data_mask_strategy_info"] = dataMaskList
			}
			columnsList = append(columnsList, columnMap)
		}
		_ = d.Set("columns", columnsList)
	}

	if tableInfo.Partitions != nil {
		partitionsList := make([]map[string]interface{}, 0, len(tableInfo.Partitions))
		for _, partition := range tableInfo.Partitions {
			partitionMap := map[string]interface{}{}
			if partition.Name != nil {
				partitionMap["name"] = partition.Name
			}
			if partition.Type != nil {
				partitionMap["type"] = partition.Type
			}
			if partition.Comment != nil {
				partitionMap["comment"] = partition.Comment
			}
			if partition.Transform != nil {
				partitionMap["transform"] = partition.Transform
			}
			if partition.TransformArgs != nil {
				transformArgs := make([]string, 0, len(partition.TransformArgs))
				for _, arg := range partition.TransformArgs {
					if arg != nil {
						transformArgs = append(transformArgs, *arg)
					}
				}
				partitionMap["transform_args"] = transformArgs
			}
			if partition.CreateTime != nil {
				partitionMap["create_time"] = partition.CreateTime
			}
			partitionsList = append(partitionsList, partitionMap)
		}
		_ = d.Set("partitions", partitionsList)
	}

	if tableInfo.Location != nil {
		_ = d.Set("location", tableInfo.Location)
	}
	if tableInfo.Properties != nil {
		propertiesList := make([]map[string]interface{}, 0, len(tableInfo.Properties))
		for _, property := range tableInfo.Properties {
			propertyMap := map[string]interface{}{}
			if property.Key != nil {
				propertyMap["key"] = property.Key
			}
			if property.Value != nil {
				propertyMap["value"] = property.Value
			}
			propertiesList = append(propertiesList, propertyMap)
		}
		_ = d.Set("properties", propertiesList)
	}
	if tableInfo.ModifiedTime != nil {
		_ = d.Set("modified_time", tableInfo.ModifiedTime)
	}
	if tableInfo.CreateTime != nil {
		_ = d.Set("create_time", tableInfo.CreateTime)
	}
	if tableInfo.InputFormat != nil {
		_ = d.Set("input_format", tableInfo.InputFormat)
	}
	if tableInfo.StorageSize != nil {
		_ = d.Set("storage_size", tableInfo.StorageSize)
	}
	if tableInfo.RecordCount != nil {
		_ = d.Set("record_count", tableInfo.RecordCount)
	}
	if tableInfo.MapMaterializedViewName != nil {
		_ = d.Set("map_materialized_view_name", tableInfo.MapMaterializedViewName)
	}
	if tableInfo.HeatValue != nil {
		_ = d.Set("heat_value", tableInfo.HeatValue)
	}
	if tableInfo.InputFormatShort != nil {
		_ = d.Set("input_format_short", tableInfo.InputFormatShort)
	}

	return nil
}

func resourceTencentCloudDlcInternalTableUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_internal_table.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	databaseName := idSplit[0]
	tableName := idSplit[1]

	needChange := false
	mutableArgs := []string{"datasource_connection_name", "table_comment", "type", "table_format", "user_alias", "user_sub_uin", "db_govern_policy_is_disable", "govern_policy", "smart_policy"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := dlc.NewAlterTableCommentRequest()
		tableBaseInfo := &dlc.TableBaseInfo{}
		tableBaseInfo.DatabaseName = helper.String(databaseName)
		tableBaseInfo.TableName = helper.String(tableName)

		if v, ok := d.GetOk("datasource_connection_name"); ok {
			tableBaseInfo.DatasourceConnectionName = helper.String(v.(string))
		}
		if v, ok := d.GetOk("table_comment"); ok {
			tableBaseInfo.TableComment = helper.String(v.(string))
		}
		if v, ok := d.GetOk("type"); ok {
			tableBaseInfo.Type = helper.String(v.(string))
		}
		if v, ok := d.GetOk("table_format"); ok {
			tableBaseInfo.TableFormat = helper.String(v.(string))
		}
		if v, ok := d.GetOk("user_alias"); ok {
			tableBaseInfo.UserAlias = helper.String(v.(string))
		}
		if v, ok := d.GetOk("user_sub_uin"); ok {
			tableBaseInfo.UserSubUin = helper.String(v.(string))
		}
		if v, ok := d.GetOk("db_govern_policy_is_disable"); ok {
			tableBaseInfo.DbGovernPolicyIsDisable = helper.String(v.(string))
		}

		if v, ok := d.GetOk("govern_policy"); ok {
			for _, item := range v.([]interface{}) {
				governPolicyMap := item.(map[string]interface{})
				dataGovernPolicy := &dlc.DataGovernPolicy{}
				if v, ok := governPolicyMap["rule_type"].(string); ok && v != "" {
					dataGovernPolicy.RuleType = helper.String(v)
				}
				if v, ok := governPolicyMap["govern_engine"].(string); ok && v != "" {
					dataGovernPolicy.GovernEngine = helper.String(v)
				}
				tableBaseInfo.GovernPolicy = dataGovernPolicy
			}
		}

		if v, ok := d.GetOk("smart_policy"); ok {
			tableBaseInfo.SmartPolicy = buildSmartPolicy(v.([]interface{}))
		}

		request.TableBaseInfo = tableBaseInfo

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().AlterTableCommentWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Update dlc_internal_table failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update dlc_internal_table failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudDlcInternalTableRead(d, meta)
}

func resourceTencentCloudDlcInternalTableDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_internal_table.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	databaseName := idSplit[0]
	tableName := idSplit[1]

	request := dlc.NewDeleteTableRequest()
	tableBaseInfo := &dlc.TableBaseInfo{}
	tableBaseInfo.DatabaseName = helper.String(databaseName)
	tableBaseInfo.TableName = helper.String(tableName)
	if v, ok := d.GetOk("datasource_connection_name"); ok {
		tableBaseInfo.DatasourceConnectionName = helper.String(v.(string))
	}
	request.TableBaseInfo = tableBaseInfo

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DeleteTableWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete dlc_internal_table failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete dlc_internal_table failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func buildSmartPolicy(smartPolicyList []interface{}) *dlc.SmartPolicy {
	if len(smartPolicyList) < 1 {
		return nil
	}
	smartPolicyMap := smartPolicyList[0].(map[string]interface{})
	smartPolicy := &dlc.SmartPolicy{}

	if v, ok := smartPolicyMap["base_info"].([]interface{}); ok && len(v) > 0 {
		baseInfoMap := v[0].(map[string]interface{})
		baseInfo := &dlc.SmartPolicyBaseInfo{}
		if v, ok := baseInfoMap["uin"].(string); ok && v != "" {
			baseInfo.Uin = helper.String(v)
		}
		if v, ok := baseInfoMap["policy_type"].(string); ok && v != "" {
			baseInfo.PolicyType = helper.String(v)
		}
		if v, ok := baseInfoMap["catalog"].(string); ok && v != "" {
			baseInfo.Catalog = helper.String(v)
		}
		if v, ok := baseInfoMap["database"].(string); ok && v != "" {
			baseInfo.Database = helper.String(v)
		}
		if v, ok := baseInfoMap["table"].(string); ok && v != "" {
			baseInfo.Table = helper.String(v)
		}
		if v, ok := baseInfoMap["app_id"].(string); ok && v != "" {
			baseInfo.AppId = helper.String(v)
		}
		smartPolicy.BaseInfo = baseInfo
	}

	if v, ok := smartPolicyMap["policy"].([]interface{}); ok && len(v) > 0 {
		policyMap := v[0].(map[string]interface{})
		policy := &dlc.SmartOptimizerPolicy{}

		if v, ok := policyMap["inherit"].(string); ok && v != "" {
			policy.Inherit = helper.String(v)
		}

		if v, ok := policyMap["resources"].([]interface{}); ok && len(v) > 0 {
			resources := make([]*dlc.ResourceInfo, 0, len(v))
			for _, resItem := range v {
				resMap := resItem.(map[string]interface{})
				resourceInfo := &dlc.ResourceInfo{}
				if v, ok := resMap["attribution_type"].(string); ok && v != "" {
					resourceInfo.AttributionType = helper.String(v)
				}
				if v, ok := resMap["resource_type"].(string); ok && v != "" {
					resourceInfo.ResourceType = helper.String(v)
				}
				if v, ok := resMap["name"].(string); ok && v != "" {
					resourceInfo.Name = helper.String(v)
				}
				if v, ok := resMap["instance"].(string); ok && v != "" {
					resourceInfo.Instance = helper.String(v)
				}
				if v, ok := resMap["status"].(int); ok && v != 0 {
					resourceInfo.Status = helper.Int64(int64(v))
				}
				if v, ok := resMap["resource_group_name"].(string); ok && v != "" {
					resourceInfo.ResourceGroupName = helper.String(v)
				}
				if v, ok := resMap["favor"].([]interface{}); ok && len(v) > 0 {
					favorList := make([]*dlc.FavorInfo, 0, len(v))
					for _, favorItem := range v {
						favorMap := favorItem.(map[string]interface{})
						favorInfo := &dlc.FavorInfo{}
						if v, ok := favorMap["priority"].(int); ok {
							favorInfo.Priority = helper.Int64(int64(v))
						}
						if v, ok := favorMap["catalog"].(string); ok && v != "" {
							favorInfo.Catalog = helper.String(v)
						}
						if v, ok := favorMap["data_base"].(string); ok && v != "" {
							favorInfo.DataBase = helper.String(v)
						}
						if v, ok := favorMap["table"].(string); ok && v != "" {
							favorInfo.Table = helper.String(v)
						}
						favorList = append(favorList, favorInfo)
					}
					resourceInfo.Favor = favorList
				}
				if v, ok := resMap["resource_conf"].([]interface{}); ok && len(v) > 0 {
					resourceConfMap := v[0].(map[string]interface{})
					resourceConf := &dlc.ResourceConf{}
					if v, ok := resourceConfMap["parallelism"].(int); ok && v != 0 {
						resourceConf.Parallelism = helper.Int64(int64(v))
					}
					resourceInfo.ResourceConf = resourceConf
				}
				resources = append(resources, resourceInfo)
			}
			policy.Resources = resources
		}

		if v, ok := policyMap["written"].([]interface{}); ok && len(v) > 0 {
			writtenMap := v[0].(map[string]interface{})
			written := &dlc.SmartOptimizerWrittenPolicy{}
			if v, ok := writtenMap["written_enable"].(string); ok && v != "" {
				written.WrittenEnable = helper.String(v)
			}
			if v, ok := writtenMap["advance_policy"].([]interface{}); ok && len(v) > 0 {
				advancePolicyMap := v[0].(map[string]interface{})
				advancePolicy := &dlc.WrittenAdvancePolicy{}
				if v, ok := advancePolicyMap["compact_enable"].(string); ok && v != "" {
					advancePolicy.CompactEnable = helper.String(v)
				}
				if v, ok := advancePolicyMap["delete_enable"].(string); ok && v != "" {
					advancePolicy.DeleteEnable = helper.String(v)
				}
				if v, ok := advancePolicyMap["min_input_files"].(int); ok && v != 0 {
					advancePolicy.MinInputFiles = helper.Int64(int64(v))
				}
				if v, ok := advancePolicyMap["target_file_size_bytes"].(int); ok && v != 0 {
					advancePolicy.TargetFileSizeBytes = helper.Int64(int64(v))
				}
				if v, ok := advancePolicyMap["retain_last"].(int); ok && v != 0 {
					advancePolicy.RetainLast = helper.Int64(int64(v))
				}
				if v, ok := advancePolicyMap["before_days"].(int); ok && v != 0 {
					advancePolicy.BeforeDays = helper.Int64(int64(v))
				}
				if v, ok := advancePolicyMap["expired_snapshots_interval_min"].(int); ok && v != 0 {
					advancePolicy.ExpiredSnapshotsIntervalMin = helper.Int64(int64(v))
				}
				if v, ok := advancePolicyMap["remove_orphan_interval_min"].(int); ok && v != 0 {
					advancePolicy.RemoveOrphanIntervalMin = helper.Int64(int64(v))
				}
				if v, ok := advancePolicyMap["cow_compact_enable"].(string); ok && v != "" {
					advancePolicy.CowCompactEnable = helper.String(v)
				}
				if v, ok := advancePolicyMap["compact_strategy"].(string); ok && v != "" {
					advancePolicy.CompactStrategy = helper.String(v)
				}
				if v, ok := advancePolicyMap["sort_orders"].([]interface{}); ok && len(v) > 0 {
					sortOrders := make([]*dlc.SortOrder, 0, len(v))
					for _, sortItem := range v {
						sortMap := sortItem.(map[string]interface{})
						sortOrder := &dlc.SortOrder{}
						if v, ok := sortMap["column"].(string); ok && v != "" {
							sortOrder.Column = helper.String(v)
						}
						if v, ok := sortMap["sort_direction"].(string); ok && v != "" {
							sortOrder.SortDirection = helper.String(v)
						}
						if v, ok := sortMap["null_order"].(string); ok && v != "" {
							sortOrder.NullOrder = helper.String(v)
						}
						sortOrders = append(sortOrders, sortOrder)
					}
					advancePolicy.SortOrders = sortOrders
				}
				written.AdvancePolicy = advancePolicy
			}
			policy.Written = written
		}

		if v, ok := policyMap["lifecycle"].([]interface{}); ok && len(v) > 0 {
			lifecycleMap := v[0].(map[string]interface{})
			lifecycle := &dlc.SmartOptimizerLifecyclePolicy{}
			if v, ok := lifecycleMap["lifecycle_enable"].(string); ok && v != "" {
				lifecycle.LifecycleEnable = helper.String(v)
			}
			if v, ok := lifecycleMap["expiration"].(int); ok && v != 0 {
				lifecycle.Expiration = helper.Int64(int64(v))
			}
			if v, ok := lifecycleMap["drop_table"].(bool); ok {
				lifecycle.DropTable = helper.Bool(v)
			}
			if v, ok := lifecycleMap["expired_field"].(string); ok && v != "" {
				lifecycle.ExpiredField = helper.String(v)
			}
			if v, ok := lifecycleMap["expired_field_format"].(string); ok && v != "" {
				lifecycle.ExpiredFieldFormat = helper.String(v)
			}
			policy.Lifecycle = lifecycle
		}

		if v, ok := policyMap["index"].([]interface{}); ok && len(v) > 0 {
			indexMap := v[0].(map[string]interface{})
			index := &dlc.SmartOptimizerIndexPolicy{}
			if v, ok := indexMap["index_enable"].(string); ok && v != "" {
				index.IndexEnable = helper.String(v)
			}
			policy.Index = index
		}

		if v, ok := policyMap["change_table"].([]interface{}); ok && len(v) > 0 {
			changeTableMap := v[0].(map[string]interface{})
			changeTable := &dlc.SmartOptimizerChangeTablePolicy{}
			if v, ok := changeTableMap["data_retention_time"].(int); ok && v != 0 {
				changeTable.DataRetentionTime = helper.Int64(int64(v))
			}
			policy.ChangeTable = changeTable
		}

		if v, ok := policyMap["table_expiration"].([]interface{}); ok && len(v) > 0 {
			tableExpirationMap := v[0].(map[string]interface{})
			tableExpiration := &dlc.TableExpirationPolicy{}
			if v, ok := tableExpirationMap["enabled"].(bool); ok {
				tableExpiration.Enabled = helper.Bool(v)
			}
			if v, ok := tableExpirationMap["expiration"].(int); ok {
				tableExpiration.Expiration = helper.Uint64(uint64(v))
			}
			policy.TableExpiration = tableExpiration
		}

		smartPolicy.Policy = policy
	}

	return smartPolicy
}

func flattenSmartPolicy(smartPolicy *dlc.SmartPolicy) []map[string]interface{} {
	smartPolicyList := make([]map[string]interface{}, 0, 1)
	smartPolicyMap := map[string]interface{}{}

	if smartPolicy.BaseInfo != nil {
		baseInfoList := make([]map[string]interface{}, 0, 1)
		baseInfoMap := map[string]interface{}{}
		if smartPolicy.BaseInfo.Uin != nil {
			baseInfoMap["uin"] = smartPolicy.BaseInfo.Uin
		}
		if smartPolicy.BaseInfo.PolicyType != nil {
			baseInfoMap["policy_type"] = smartPolicy.BaseInfo.PolicyType
		}
		if smartPolicy.BaseInfo.Catalog != nil {
			baseInfoMap["catalog"] = smartPolicy.BaseInfo.Catalog
		}
		if smartPolicy.BaseInfo.Database != nil {
			baseInfoMap["database"] = smartPolicy.BaseInfo.Database
		}
		if smartPolicy.BaseInfo.Table != nil {
			baseInfoMap["table"] = smartPolicy.BaseInfo.Table
		}
		if smartPolicy.BaseInfo.AppId != nil {
			baseInfoMap["app_id"] = smartPolicy.BaseInfo.AppId
		}
		baseInfoList = append(baseInfoList, baseInfoMap)
		smartPolicyMap["base_info"] = baseInfoList
	}

	if smartPolicy.Policy != nil {
		policyList := make([]map[string]interface{}, 0, 1)
		policyMap := map[string]interface{}{}
		if smartPolicy.Policy.Inherit != nil {
			policyMap["inherit"] = smartPolicy.Policy.Inherit
		}

		if smartPolicy.Policy.Resources != nil {
			resourcesList := make([]map[string]interface{}, 0, len(smartPolicy.Policy.Resources))
			for _, res := range smartPolicy.Policy.Resources {
				resMap := map[string]interface{}{}
				if res.AttributionType != nil {
					resMap["attribution_type"] = res.AttributionType
				}
				if res.ResourceType != nil {
					resMap["resource_type"] = res.ResourceType
				}
				if res.Name != nil {
					resMap["name"] = res.Name
				}
				if res.Instance != nil {
					resMap["instance"] = res.Instance
				}
				if res.Status != nil {
					resMap["status"] = res.Status
				}
				if res.ResourceGroupName != nil {
					resMap["resource_group_name"] = res.ResourceGroupName
				}
				if res.Favor != nil {
					favorList := make([]map[string]interface{}, 0, len(res.Favor))
					for _, favor := range res.Favor {
						favorMap := map[string]interface{}{}
						if favor.Priority != nil {
							favorMap["priority"] = favor.Priority
						}
						if favor.Catalog != nil {
							favorMap["catalog"] = favor.Catalog
						}
						if favor.DataBase != nil {
							favorMap["data_base"] = favor.DataBase
						}
						if favor.Table != nil {
							favorMap["table"] = favor.Table
						}
						favorList = append(favorList, favorMap)
					}
					resMap["favor"] = favorList
				}
				if res.ResourceConf != nil {
					resourceConfList := make([]map[string]interface{}, 0, 1)
					resourceConfMap := map[string]interface{}{}
					if res.ResourceConf.Parallelism != nil {
						resourceConfMap["parallelism"] = res.ResourceConf.Parallelism
					}
					resourceConfList = append(resourceConfList, resourceConfMap)
					resMap["resource_conf"] = resourceConfList
				}
				resourcesList = append(resourcesList, resMap)
			}
			policyMap["resources"] = resourcesList
		}

		if smartPolicy.Policy.Written != nil {
			writtenList := make([]map[string]interface{}, 0, 1)
			writtenMap := map[string]interface{}{}
			if smartPolicy.Policy.Written.WrittenEnable != nil {
				writtenMap["written_enable"] = smartPolicy.Policy.Written.WrittenEnable
			}
			if smartPolicy.Policy.Written.AdvancePolicy != nil {
				advancePolicyList := make([]map[string]interface{}, 0, 1)
				advancePolicyMap := map[string]interface{}{}
				ap := smartPolicy.Policy.Written.AdvancePolicy
				if ap.CompactEnable != nil {
					advancePolicyMap["compact_enable"] = ap.CompactEnable
				}
				if ap.DeleteEnable != nil {
					advancePolicyMap["delete_enable"] = ap.DeleteEnable
				}
				if ap.MinInputFiles != nil {
					advancePolicyMap["min_input_files"] = ap.MinInputFiles
				}
				if ap.TargetFileSizeBytes != nil {
					advancePolicyMap["target_file_size_bytes"] = ap.TargetFileSizeBytes
				}
				if ap.RetainLast != nil {
					advancePolicyMap["retain_last"] = ap.RetainLast
				}
				if ap.BeforeDays != nil {
					advancePolicyMap["before_days"] = ap.BeforeDays
				}
				if ap.ExpiredSnapshotsIntervalMin != nil {
					advancePolicyMap["expired_snapshots_interval_min"] = ap.ExpiredSnapshotsIntervalMin
				}
				if ap.RemoveOrphanIntervalMin != nil {
					advancePolicyMap["remove_orphan_interval_min"] = ap.RemoveOrphanIntervalMin
				}
				if ap.CowCompactEnable != nil {
					advancePolicyMap["cow_compact_enable"] = ap.CowCompactEnable
				}
				if ap.CompactStrategy != nil {
					advancePolicyMap["compact_strategy"] = ap.CompactStrategy
				}
				if ap.SortOrders != nil {
					sortOrdersList := make([]map[string]interface{}, 0, len(ap.SortOrders))
					for _, sortOrder := range ap.SortOrders {
						sortOrderMap := map[string]interface{}{}
						if sortOrder.Column != nil {
							sortOrderMap["column"] = sortOrder.Column
						}
						if sortOrder.SortDirection != nil {
							sortOrderMap["sort_direction"] = sortOrder.SortDirection
						}
						if sortOrder.NullOrder != nil {
							sortOrderMap["null_order"] = sortOrder.NullOrder
						}
						sortOrdersList = append(sortOrdersList, sortOrderMap)
					}
					advancePolicyMap["sort_orders"] = sortOrdersList
				}
				advancePolicyList = append(advancePolicyList, advancePolicyMap)
				writtenMap["advance_policy"] = advancePolicyList
			}
			writtenList = append(writtenList, writtenMap)
			policyMap["written"] = writtenList
		}

		if smartPolicy.Policy.Lifecycle != nil {
			lifecycleList := make([]map[string]interface{}, 0, 1)
			lifecycleMap := map[string]interface{}{}
			if smartPolicy.Policy.Lifecycle.LifecycleEnable != nil {
				lifecycleMap["lifecycle_enable"] = smartPolicy.Policy.Lifecycle.LifecycleEnable
			}
			if smartPolicy.Policy.Lifecycle.Expiration != nil {
				lifecycleMap["expiration"] = smartPolicy.Policy.Lifecycle.Expiration
			}
			if smartPolicy.Policy.Lifecycle.DropTable != nil {
				lifecycleMap["drop_table"] = smartPolicy.Policy.Lifecycle.DropTable
			}
			if smartPolicy.Policy.Lifecycle.ExpiredField != nil {
				lifecycleMap["expired_field"] = smartPolicy.Policy.Lifecycle.ExpiredField
			}
			if smartPolicy.Policy.Lifecycle.ExpiredFieldFormat != nil {
				lifecycleMap["expired_field_format"] = smartPolicy.Policy.Lifecycle.ExpiredFieldFormat
			}
			lifecycleList = append(lifecycleList, lifecycleMap)
			policyMap["lifecycle"] = lifecycleList
		}

		if smartPolicy.Policy.Index != nil {
			indexList := make([]map[string]interface{}, 0, 1)
			indexMap := map[string]interface{}{}
			if smartPolicy.Policy.Index.IndexEnable != nil {
				indexMap["index_enable"] = smartPolicy.Policy.Index.IndexEnable
			}
			indexList = append(indexList, indexMap)
			policyMap["index"] = indexList
		}

		if smartPolicy.Policy.ChangeTable != nil {
			changeTableList := make([]map[string]interface{}, 0, 1)
			changeTableMap := map[string]interface{}{}
			if smartPolicy.Policy.ChangeTable.DataRetentionTime != nil {
				changeTableMap["data_retention_time"] = smartPolicy.Policy.ChangeTable.DataRetentionTime
			}
			changeTableList = append(changeTableList, changeTableMap)
			policyMap["change_table"] = changeTableList
		}

		if smartPolicy.Policy.TableExpiration != nil {
			tableExpirationList := make([]map[string]interface{}, 0, 1)
			tableExpirationMap := map[string]interface{}{}
			if smartPolicy.Policy.TableExpiration.Enabled != nil {
				tableExpirationMap["enabled"] = smartPolicy.Policy.TableExpiration.Enabled
			}
			if smartPolicy.Policy.TableExpiration.Expiration != nil {
				tableExpirationMap["expiration"] = smartPolicy.Policy.TableExpiration.Expiration
			}
			tableExpirationList = append(tableExpirationList, tableExpirationMap)
			policyMap["table_expiration"] = tableExpirationList
		}

		policyList = append(policyList, policyMap)
		smartPolicyMap["policy"] = policyList
	}

	smartPolicyList = append(smartPolicyList, smartPolicyMap)
	return smartPolicyList
}
