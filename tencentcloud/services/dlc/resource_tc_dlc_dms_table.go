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

func ResourceTencentCloudDlcDmsTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcDmsTableCreate,
		Read:   resourceTencentCloudDlcDmsTableRead,
		Update: resourceTencentCloudDlcDmsTableUpdate,
		Delete: resourceTencentCloudDlcDmsTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"asset": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Basic object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Primary key.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name.",
						},
						"guid": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Object GUID value.",
						},
						"catalog": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Data catalog.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description.",
						},
						"owner": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Object owner.",
						},
						"owner_account": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Object owner account.",
						},
						"perm_values": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Permissions.",
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
						"params": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional attributes.",
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
						"biz_params": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional business attributes.",
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
						"data_version": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Data version.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Create time.",
						},
						"modified_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Modified time.",
						},
						"datasource_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Data source primary key.",
						},
					},
				},
			},

			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Table type: EXTERNAL_TABLE, VIRTUAL_VIEW, MATERIALIZED_VIEW.",
			},

			"db_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Database name.",
			},

			"storage_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Storage size.",
			},

			"record_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Record count.",
			},

			"life_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Life cycle.",
			},

			"data_update_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Data update time.",
			},

			"struct_update_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Structure update time.",
			},

			"last_access_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Last access time.",
			},

			"sds": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Storage object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Storage location.",
						},
						"input_format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Input format.",
						},
						"output_format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Output format.",
						},
						"num_buckets": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Bucket count.",
						},
						"compressed": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether compressed.",
						},
						"stored_as_sub_directories": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether has sub directories.",
						},
						"serde_lib": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Serde lib.",
						},
						"serde_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Serde name.",
						},
						"bucket_cols": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Bucket columns.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"serde_params": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Serde parameters.",
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
						"params": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional parameters.",
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
						"sort_cols": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Column sort (Expired).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"col": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Column name.",
									},
									"order": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Order.",
									},
								},
							},
						},
						"dms_cols": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Columns.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Name.",
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Description.",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Type.",
									},
									"position": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Position.",
									},
									"params": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Additional parameters.",
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
									"biz_params": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Business parameters.",
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
									"is_partition": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether partition.",
									},
								},
							},
						},
						"sort_columns": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Column sort fields.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"col": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Column name.",
									},
									"order": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Order.",
									},
								},
							},
						},
					},
				},
			},

			"columns": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Columns.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type.",
						},
						"position": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Position.",
						},
						"params": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional parameters.",
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
						"biz_params": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Business parameters.",
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
						"is_partition": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether partition.",
						},
					},
				},
			},

			"partition_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Partition keys.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type.",
						},
						"position": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Position.",
						},
						"params": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional parameters.",
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
						"biz_params": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Business parameters.",
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
						"is_partition": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether partition.",
						},
					},
				},
			},

			"view_original_text": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "View original text.",
			},

			"view_expanded_text": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "View expanded text.",
			},

			"partitions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Partitions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database name.",
						},
						"schema_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Schema name.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Table name.",
						},
						"data_version": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Data version.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Partition name.",
						},
						"values": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Value list.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"storage_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Storage size.",
						},
						"record_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Record count.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Create time.",
						},
						"modified_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Modified time.",
						},
						"last_access_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last access time.",
						},
						"params": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional attributes.",
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
						"sds": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Storage object.",
							Elem: &schema.Resource{
								Schema: dmsSdsSchema(),
							},
						},
						"datasource_connection_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data source connection name.",
						},
					},
				},
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Table name.",
			},

			"datasource_connection_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Data source connection name.",
			},

			"delete_data": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to delete data.",
			},

			"env_props": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Environment attributes.",
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

			"schema_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Schema name.",
			},

			"retention": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Hive retention version.",
			},
		},
	}
}

// dmsSdsSchema returns the schema definition for DMSSds structure, reused by sds and partition.sds.
func dmsSdsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"location": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Storage location.",
		},
		"input_format": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Input format.",
		},
		"output_format": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Output format.",
		},
		"num_buckets": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Bucket count.",
		},
		"compressed": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether compressed.",
		},
		"stored_as_sub_directories": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether has sub directories.",
		},
		"serde_lib": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Serde lib.",
		},
		"serde_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Serde name.",
		},
		"bucket_cols": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Bucket columns.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"serde_params": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Serde parameters.",
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
		"params": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Additional parameters.",
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
		"sort_cols": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Column sort (Expired).",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"col": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Column name.",
					},
					"order": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Order.",
					},
				},
			},
		},
		"dms_cols": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Columns.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name.",
					},
					"description": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Description.",
					},
					"type": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Type.",
					},
					"position": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Position.",
					},
					"params": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Additional parameters.",
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
					"biz_params": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Business parameters.",
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
					"is_partition": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Whether partition.",
					},
				},
			},
		},
		"sort_columns": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Column sort fields.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"col": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Column name.",
					},
					"order": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Order.",
					},
				},
			},
		},
	}
}

func resourceTencentCloudDlcDmsTableCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_dms_table.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = dlc.NewCreateDMSTableRequest()
		response = dlc.NewCreateDMSTableResponse()
		dbName   string
		name     string
	)

	if v, ok := d.GetOk("db_name"); ok {
		request.DbName = helper.String(v.(string))
		dbName = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
		name = v.(string)
	}

	if v, ok := d.GetOk("asset"); ok {
		for _, item := range v.([]interface{}) {
			assetMap := item.(map[string]interface{})
			asset := dlc.Asset{}
			if v, ok := assetMap["id"].(int); ok && v != 0 {
				asset.Id = helper.Int64(int64(v))
			}
			if v, ok := assetMap["name"].(string); ok && v != "" {
				asset.Name = helper.String(v)
			}
			if v, ok := assetMap["guid"].(string); ok && v != "" {
				asset.Guid = helper.String(v)
			}
			if v, ok := assetMap["catalog"].(string); ok && v != "" {
				asset.Catalog = helper.String(v)
			}
			if v, ok := assetMap["description"].(string); ok && v != "" {
				asset.Description = helper.String(v)
			}
			if v, ok := assetMap["owner"].(string); ok && v != "" {
				asset.Owner = helper.String(v)
			}
			if v, ok := assetMap["owner_account"].(string); ok && v != "" {
				asset.OwnerAccount = helper.String(v)
			}
			if v, ok := assetMap["perm_values"].([]interface{}); ok && len(v) > 0 {
				for _, kvItem := range v {
					kvMap := kvItem.(map[string]interface{})
					kvPair := dlc.KVPair{}
					if v, ok := kvMap["key"].(string); ok && v != "" {
						kvPair.Key = helper.String(v)
					}
					if v, ok := kvMap["value"].(string); ok && v != "" {
						kvPair.Value = helper.String(v)
					}
					asset.PermValues = append(asset.PermValues, &kvPair)
				}
			}
			if v, ok := assetMap["params"].([]interface{}); ok && len(v) > 0 {
				for _, kvItem := range v {
					kvMap := kvItem.(map[string]interface{})
					kvPair := dlc.KVPair{}
					if v, ok := kvMap["key"].(string); ok && v != "" {
						kvPair.Key = helper.String(v)
					}
					if v, ok := kvMap["value"].(string); ok && v != "" {
						kvPair.Value = helper.String(v)
					}
					asset.Params = append(asset.Params, &kvPair)
				}
			}
			if v, ok := assetMap["biz_params"].([]interface{}); ok && len(v) > 0 {
				for _, kvItem := range v {
					kvMap := kvItem.(map[string]interface{})
					kvPair := dlc.KVPair{}
					if v, ok := kvMap["key"].(string); ok && v != "" {
						kvPair.Key = helper.String(v)
					}
					if v, ok := kvMap["value"].(string); ok && v != "" {
						kvPair.Value = helper.String(v)
					}
					asset.BizParams = append(asset.BizParams, &kvPair)
				}
			}
			if v, ok := assetMap["data_version"].(int); ok && v != 0 {
				asset.DataVersion = helper.Int64(int64(v))
			}
			if v, ok := assetMap["create_time"].(string); ok && v != "" {
				asset.CreateTime = helper.String(v)
			}
			if v, ok := assetMap["modified_time"].(string); ok && v != "" {
				asset.ModifiedTime = helper.String(v)
			}
			if v, ok := assetMap["datasource_id"].(int); ok && v != 0 {
				asset.DatasourceId = helper.Int64(int64(v))
			}
			request.Asset = &asset
		}
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("storage_size"); ok {
		request.StorageSize = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("record_count"); ok {
		request.RecordCount = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("life_time"); ok {
		request.LifeTime = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("data_update_time"); ok {
		request.DataUpdateTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("struct_update_time"); ok {
		request.StructUpdateTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("last_access_time"); ok {
		request.LastAccessTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sds"); ok {
		for _, item := range v.([]interface{}) {
			sdsMap := item.(map[string]interface{})
			sds := buildDlcDmsSds(sdsMap)
			request.Sds = &sds
		}
	}

	if v, ok := d.GetOk("columns"); ok {
		for _, item := range v.([]interface{}) {
			columnMap := item.(map[string]interface{})
			column := buildDlcDmsColumn(columnMap)
			request.Columns = append(request.Columns, &column)
		}
	}

	if v, ok := d.GetOk("partition_keys"); ok {
		for _, item := range v.([]interface{}) {
			columnMap := item.(map[string]interface{})
			column := buildDlcDmsColumn(columnMap)
			request.PartitionKeys = append(request.PartitionKeys, &column)
		}
	}

	if v, ok := d.GetOk("view_original_text"); ok {
		request.ViewOriginalText = helper.String(v.(string))
	}

	if v, ok := d.GetOk("view_expanded_text"); ok {
		request.ViewExpandedText = helper.String(v.(string))
	}

	if v, ok := d.GetOk("partitions"); ok {
		for _, item := range v.([]interface{}) {
			partitionMap := item.(map[string]interface{})
			partition := buildDlcDmsPartition(partitionMap)
			request.Partitions = append(request.Partitions, &partition)
		}
	}

	if v, ok := d.GetOk("datasource_connection_name"); ok {
		request.DatasourceConnectionName = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().CreateDMSTableWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dlc_dms_table failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dlc_dms_table failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[INFO]%s create dlc_dms_table success, db_name=%s, name=%s, RequestId=%s", logId, dbName, name, *response.Response.RequestId)

	if dbName == "" || name == "" {
		return fmt.Errorf("dlc_dms_table db_name or name is empty, db_name=%s, name=%s", dbName, name)
	}

	d.SetId(strings.Join([]string{dbName, name}, tccommon.FILED_SP))
	return resourceTencentCloudDlcDmsTableRead(d, meta)
}

func resourceTencentCloudDlcDmsTableRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_dms_table.read")()
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

	dbName := idSplit[0]
	name := idSplit[1]

	respData, err := service.DescribeDmsTableById(ctx, dbName, name)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] dlc_dms_table id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.Asset != nil {
		assetMap := map[string]interface{}{}
		if respData.Asset.Id != nil {
			assetMap["id"] = respData.Asset.Id
		}
		if respData.Asset.Name != nil {
			assetMap["name"] = respData.Asset.Name
		}
		if respData.Asset.Guid != nil {
			assetMap["guid"] = respData.Asset.Guid
		}
		if respData.Asset.Catalog != nil {
			assetMap["catalog"] = respData.Asset.Catalog
		}
		if respData.Asset.Description != nil {
			assetMap["description"] = respData.Asset.Description
		}
		if respData.Asset.Owner != nil {
			assetMap["owner"] = respData.Asset.Owner
		}
		if respData.Asset.OwnerAccount != nil {
			assetMap["owner_account"] = respData.Asset.OwnerAccount
		}
		if respData.Asset.PermValues != nil {
			permValuesList := make([]map[string]interface{}, 0, len(respData.Asset.PermValues))
			for _, kvPair := range respData.Asset.PermValues {
				kvMap := map[string]interface{}{}
				if kvPair.Key != nil {
					kvMap["key"] = kvPair.Key
				}
				if kvPair.Value != nil {
					kvMap["value"] = kvPair.Value
				}
				permValuesList = append(permValuesList, kvMap)
			}
			assetMap["perm_values"] = permValuesList
		}
		if respData.Asset.Params != nil {
			paramsList := make([]map[string]interface{}, 0, len(respData.Asset.Params))
			for _, kvPair := range respData.Asset.Params {
				kvMap := map[string]interface{}{}
				if kvPair.Key != nil {
					kvMap["key"] = kvPair.Key
				}
				if kvPair.Value != nil {
					kvMap["value"] = kvPair.Value
				}
				paramsList = append(paramsList, kvMap)
			}
			assetMap["params"] = paramsList
		}
		if respData.Asset.BizParams != nil {
			bizParamsList := make([]map[string]interface{}, 0, len(respData.Asset.BizParams))
			for _, kvPair := range respData.Asset.BizParams {
				kvMap := map[string]interface{}{}
				if kvPair.Key != nil {
					kvMap["key"] = kvPair.Key
				}
				if kvPair.Value != nil {
					kvMap["value"] = kvPair.Value
				}
				bizParamsList = append(bizParamsList, kvMap)
			}
			assetMap["biz_params"] = bizParamsList
		}
		if respData.Asset.DataVersion != nil {
			assetMap["data_version"] = respData.Asset.DataVersion
		}
		if respData.Asset.CreateTime != nil {
			assetMap["create_time"] = respData.Asset.CreateTime
		}
		if respData.Asset.ModifiedTime != nil {
			assetMap["modified_time"] = respData.Asset.ModifiedTime
		}
		if respData.Asset.DatasourceId != nil {
			assetMap["datasource_id"] = respData.Asset.DatasourceId
		}
		_ = d.Set("asset", []interface{}{assetMap})
	}

	if respData.Type != nil {
		_ = d.Set("type", respData.Type)
	}

	if respData.DbName != nil {
		_ = d.Set("db_name", respData.DbName)
	}

	if respData.StorageSize != nil {
		_ = d.Set("storage_size", respData.StorageSize)
	}

	if respData.RecordCount != nil {
		_ = d.Set("record_count", respData.RecordCount)
	}

	if respData.LifeTime != nil {
		_ = d.Set("life_time", respData.LifeTime)
	}

	if respData.DataUpdateTime != nil {
		_ = d.Set("data_update_time", respData.DataUpdateTime)
	}

	if respData.StructUpdateTime != nil {
		_ = d.Set("struct_update_time", respData.StructUpdateTime)
	}

	if respData.LastAccessTime != nil {
		_ = d.Set("last_access_time", respData.LastAccessTime)
	}

	if respData.Sds != nil {
		sdsList := []interface{}{flattenDlcDmsSds(respData.Sds)}
		_ = d.Set("sds", sdsList)
	}

	if respData.Columns != nil {
		columnsList := make([]interface{}, 0, len(respData.Columns))
		for _, column := range respData.Columns {
			columnsList = append(columnsList, flattenDlcDmsColumn(column))
		}
		_ = d.Set("columns", columnsList)
	}

	if respData.PartitionKeys != nil {
		partitionKeysList := make([]interface{}, 0, len(respData.PartitionKeys))
		for _, column := range respData.PartitionKeys {
			partitionKeysList = append(partitionKeysList, flattenDlcDmsColumn(column))
		}
		_ = d.Set("partition_keys", partitionKeysList)
	}

	if respData.ViewOriginalText != nil {
		_ = d.Set("view_original_text", respData.ViewOriginalText)
	}

	if respData.ViewExpandedText != nil {
		_ = d.Set("view_expanded_text", respData.ViewExpandedText)
	}

	if respData.Partitions != nil {
		partitionsList := make([]interface{}, 0, len(respData.Partitions))
		for _, partition := range respData.Partitions {
			partitionsList = append(partitionsList, flattenDlcDmsPartition(partition))
		}
		_ = d.Set("partitions", partitionsList)
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.SchemaName != nil {
		_ = d.Set("schema_name", respData.SchemaName)
	}

	if respData.Retention != nil {
		_ = d.Set("retention", respData.Retention)
	}

	return nil
}

func resourceTencentCloudDlcDmsTableUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_dms_table.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	currentDbName := idSplit[0]
	currentName := idSplit[1]

	request := dlc.NewAlterDMSTableRequest()

	if d.HasChange("db_name") {
		_, newDbName := d.GetChange("db_name")
		request.DbName = helper.String(newDbName.(string))
		request.CurrentDbName = helper.String(currentDbName)
	}

	if d.HasChange("name") {
		_, newName := d.GetChange("name")
		request.Name = helper.String(newName.(string))
		request.CurrentName = helper.String(currentName)
	}

	if v, ok := d.GetOk("asset"); ok {
		for _, item := range v.([]interface{}) {
			assetMap := item.(map[string]interface{})
			asset := dlc.Asset{}
			if v, ok := assetMap["id"].(int); ok && v != 0 {
				asset.Id = helper.Int64(int64(v))
			}
			if v, ok := assetMap["name"].(string); ok && v != "" {
				asset.Name = helper.String(v)
			}
			if v, ok := assetMap["guid"].(string); ok && v != "" {
				asset.Guid = helper.String(v)
			}
			if v, ok := assetMap["catalog"].(string); ok && v != "" {
				asset.Catalog = helper.String(v)
			}
			if v, ok := assetMap["description"].(string); ok && v != "" {
				asset.Description = helper.String(v)
			}
			if v, ok := assetMap["owner"].(string); ok && v != "" {
				asset.Owner = helper.String(v)
			}
			if v, ok := assetMap["owner_account"].(string); ok && v != "" {
				asset.OwnerAccount = helper.String(v)
			}
			if v, ok := assetMap["perm_values"].([]interface{}); ok && len(v) > 0 {
				for _, kvItem := range v {
					kvMap := kvItem.(map[string]interface{})
					kvPair := dlc.KVPair{}
					if v, ok := kvMap["key"].(string); ok && v != "" {
						kvPair.Key = helper.String(v)
					}
					if v, ok := kvMap["value"].(string); ok && v != "" {
						kvPair.Value = helper.String(v)
					}
					asset.PermValues = append(asset.PermValues, &kvPair)
				}
			}
			if v, ok := assetMap["params"].([]interface{}); ok && len(v) > 0 {
				for _, kvItem := range v {
					kvMap := kvItem.(map[string]interface{})
					kvPair := dlc.KVPair{}
					if v, ok := kvMap["key"].(string); ok && v != "" {
						kvPair.Key = helper.String(v)
					}
					if v, ok := kvMap["value"].(string); ok && v != "" {
						kvPair.Value = helper.String(v)
					}
					asset.Params = append(asset.Params, &kvPair)
				}
			}
			if v, ok := assetMap["biz_params"].([]interface{}); ok && len(v) > 0 {
				for _, kvItem := range v {
					kvMap := kvItem.(map[string]interface{})
					kvPair := dlc.KVPair{}
					if v, ok := kvMap["key"].(string); ok && v != "" {
						kvPair.Key = helper.String(v)
					}
					if v, ok := kvMap["value"].(string); ok && v != "" {
						kvPair.Value = helper.String(v)
					}
					asset.BizParams = append(asset.BizParams, &kvPair)
				}
			}
			if v, ok := assetMap["data_version"].(int); ok && v != 0 {
				asset.DataVersion = helper.Int64(int64(v))
			}
			if v, ok := assetMap["create_time"].(string); ok && v != "" {
				asset.CreateTime = helper.String(v)
			}
			if v, ok := assetMap["modified_time"].(string); ok && v != "" {
				asset.ModifiedTime = helper.String(v)
			}
			if v, ok := assetMap["datasource_id"].(int); ok && v != 0 {
				asset.DatasourceId = helper.Int64(int64(v))
			}
			request.Asset = &asset
		}
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if !d.HasChange("db_name") {
		request.DbName = helper.String(currentDbName)
	}

	if !d.HasChange("name") {
		request.Name = helper.String(currentName)
	}

	if v, ok := d.GetOkExists("storage_size"); ok {
		request.StorageSize = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("record_count"); ok {
		request.RecordCount = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("life_time"); ok {
		request.LifeTime = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("data_update_time"); ok {
		request.DataUpdateTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("struct_update_time"); ok {
		request.StructUpdateTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("last_access_time"); ok {
		request.LastAccessTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sds"); ok {
		for _, item := range v.([]interface{}) {
			sdsMap := item.(map[string]interface{})
			sds := buildDlcDmsSds(sdsMap)
			request.Sds = &sds
		}
	}

	if v, ok := d.GetOk("columns"); ok {
		for _, item := range v.([]interface{}) {
			columnMap := item.(map[string]interface{})
			column := buildDlcDmsColumn(columnMap)
			request.Columns = append(request.Columns, &column)
		}
	}

	if v, ok := d.GetOk("partition_keys"); ok {
		for _, item := range v.([]interface{}) {
			columnMap := item.(map[string]interface{})
			column := buildDlcDmsColumn(columnMap)
			request.PartitionKeys = append(request.PartitionKeys, &column)
		}
	}

	if v, ok := d.GetOk("view_original_text"); ok {
		request.ViewOriginalText = helper.String(v.(string))
	}

	if v, ok := d.GetOk("view_expanded_text"); ok {
		request.ViewExpandedText = helper.String(v.(string))
	}

	if v, ok := d.GetOk("partitions"); ok {
		for _, item := range v.([]interface{}) {
			partitionMap := item.(map[string]interface{})
			partition := buildDlcDmsPartition(partitionMap)
			request.Partitions = append(request.Partitions, &partition)
		}
	}

	if v, ok := d.GetOk("datasource_connection_name"); ok {
		request.DatasourceConnectionName = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().AlterDMSTableWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Update dlc_dms_table failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update dlc_dms_table failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if request.DbName != nil && request.Name != nil {
		newDbName := *request.DbName
		newName := *request.Name
		if newDbName != currentDbName || newName != currentName {
			d.SetId(strings.Join([]string{newDbName, newName}, tccommon.FILED_SP))
		}
	}

	return resourceTencentCloudDlcDmsTableRead(d, meta)
}

func resourceTencentCloudDlcDmsTableDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_dms_table.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = dlc.NewDropDMSTableRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	dbName := idSplit[0]
	name := idSplit[1]

	request.DbName = helper.String(dbName)
	request.Name = helper.String(name)

	if v, ok := d.GetOkExists("delete_data"); ok {
		request.DeleteData = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("env_props"); ok {
		for _, item := range v.([]interface{}) {
			kvMap := item.(map[string]interface{})
			kvPair := dlc.KVPair{}
			if v, ok := kvMap["key"].(string); ok && v != "" {
				kvPair.Key = helper.String(v)
			}
			if v, ok := kvMap["value"].(string); ok && v != "" {
				kvPair.Value = helper.String(v)
			}
			request.EnvProps = &kvPair
		}
	}

	if v, ok := d.GetOk("datasource_connection_name"); ok {
		request.DatasourceConnectionName = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DropDMSTableWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete dlc_dms_table failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete dlc_dms_table failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

// buildDlcDmsSds builds a DMSSds struct from a schema map.
func buildDlcDmsSds(sdsMap map[string]interface{}) dlc.DMSSds {
	sds := dlc.DMSSds{}
	if v, ok := sdsMap["location"].(string); ok && v != "" {
		sds.Location = helper.String(v)
	}
	if v, ok := sdsMap["input_format"].(string); ok && v != "" {
		sds.InputFormat = helper.String(v)
	}
	if v, ok := sdsMap["output_format"].(string); ok && v != "" {
		sds.OutputFormat = helper.String(v)
	}
	if v, ok := sdsMap["num_buckets"].(int); ok && v != 0 {
		sds.NumBuckets = helper.Int64(int64(v))
	}
	if v, ok := sdsMap["compressed"].(bool); ok {
		sds.Compressed = helper.Bool(v)
	}
	if v, ok := sdsMap["stored_as_sub_directories"].(bool); ok {
		sds.StoredAsSubDirectories = helper.Bool(v)
	}
	if v, ok := sdsMap["serde_lib"].(string); ok && v != "" {
		sds.SerdeLib = helper.String(v)
	}
	if v, ok := sdsMap["serde_name"].(string); ok && v != "" {
		sds.SerdeName = helper.String(v)
	}
	if v, ok := sdsMap["bucket_cols"].([]interface{}); ok && len(v) > 0 {
		for _, col := range v {
			if colStr, ok := col.(string); ok && colStr != "" {
				sds.BucketCols = append(sds.BucketCols, helper.String(colStr))
			}
		}
	}
	if v, ok := sdsMap["serde_params"].([]interface{}); ok && len(v) > 0 {
		for _, kvItem := range v {
			kvMap := kvItem.(map[string]interface{})
			kvPair := dlc.KVPair{}
			if v, ok := kvMap["key"].(string); ok && v != "" {
				kvPair.Key = helper.String(v)
			}
			if v, ok := kvMap["value"].(string); ok && v != "" {
				kvPair.Value = helper.String(v)
			}
			sds.SerdeParams = append(sds.SerdeParams, &kvPair)
		}
	}
	if v, ok := sdsMap["params"].([]interface{}); ok && len(v) > 0 {
		for _, kvItem := range v {
			kvMap := kvItem.(map[string]interface{})
			kvPair := dlc.KVPair{}
			if v, ok := kvMap["key"].(string); ok && v != "" {
				kvPair.Key = helper.String(v)
			}
			if v, ok := kvMap["value"].(string); ok && v != "" {
				kvPair.Value = helper.String(v)
			}
			sds.Params = append(sds.Params, &kvPair)
		}
	}
	if v, ok := sdsMap["sort_cols"].([]interface{}); ok && len(v) > 0 {
		for _, sortItem := range v {
			sortMap := sortItem.(map[string]interface{})
			sortCol := dlc.DMSColumnOrder{}
			if v, ok := sortMap["col"].(string); ok && v != "" {
				sortCol.Col = helper.String(v)
			}
			if v, ok := sortMap["order"].(int); ok && v != 0 {
				sortCol.Order = helper.Int64(int64(v))
			}
			sds.SortCols = &sortCol
		}
	}
	if v, ok := sdsMap["dms_cols"].([]interface{}); ok && len(v) > 0 {
		for _, colItem := range v {
			colMap := colItem.(map[string]interface{})
			col := buildDlcDmsColumn(colMap)
			sds.Cols = append(sds.Cols, &col)
		}
	}
	if v, ok := sdsMap["sort_columns"].([]interface{}); ok && len(v) > 0 {
		for _, sortItem := range v {
			sortMap := sortItem.(map[string]interface{})
			sortCol := dlc.DMSColumnOrder{}
			if v, ok := sortMap["col"].(string); ok && v != "" {
				sortCol.Col = helper.String(v)
			}
			if v, ok := sortMap["order"].(int); ok && v != 0 {
				sortCol.Order = helper.Int64(int64(v))
			}
			sds.SortColumns = append(sds.SortColumns, &sortCol)
		}
	}
	return sds
}

// buildDlcDmsColumn builds a DMSColumn struct from a schema map.
func buildDlcDmsColumn(columnMap map[string]interface{}) dlc.DMSColumn {
	column := dlc.DMSColumn{}
	if v, ok := columnMap["name"].(string); ok && v != "" {
		column.Name = helper.String(v)
	}
	if v, ok := columnMap["description"].(string); ok && v != "" {
		column.Description = helper.String(v)
	}
	if v, ok := columnMap["type"].(string); ok && v != "" {
		column.Type = helper.String(v)
	}
	if v, ok := columnMap["position"].(int); ok && v != 0 {
		column.Position = helper.Int64(int64(v))
	}
	if v, ok := columnMap["params"].([]interface{}); ok && len(v) > 0 {
		for _, kvItem := range v {
			kvMap := kvItem.(map[string]interface{})
			kvPair := dlc.KVPair{}
			if v, ok := kvMap["key"].(string); ok && v != "" {
				kvPair.Key = helper.String(v)
			}
			if v, ok := kvMap["value"].(string); ok && v != "" {
				kvPair.Value = helper.String(v)
			}
			column.Params = append(column.Params, &kvPair)
		}
	}
	if v, ok := columnMap["biz_params"].([]interface{}); ok && len(v) > 0 {
		for _, kvItem := range v {
			kvMap := kvItem.(map[string]interface{})
			kvPair := dlc.KVPair{}
			if v, ok := kvMap["key"].(string); ok && v != "" {
				kvPair.Key = helper.String(v)
			}
			if v, ok := kvMap["value"].(string); ok && v != "" {
				kvPair.Value = helper.String(v)
			}
			column.BizParams = append(column.BizParams, &kvPair)
		}
	}
	if v, ok := columnMap["is_partition"].(bool); ok {
		column.IsPartition = helper.Bool(v)
	}
	return column
}

// buildDlcDmsPartition builds a DMSPartition struct from a schema map.
func buildDlcDmsPartition(partitionMap map[string]interface{}) dlc.DMSPartition {
	partition := dlc.DMSPartition{}
	if v, ok := partitionMap["database_name"].(string); ok && v != "" {
		partition.DatabaseName = helper.String(v)
	}
	if v, ok := partitionMap["schema_name"].(string); ok && v != "" {
		partition.SchemaName = helper.String(v)
	}
	if v, ok := partitionMap["table_name"].(string); ok && v != "" {
		partition.TableName = helper.String(v)
	}
	if v, ok := partitionMap["data_version"].(int); ok && v != 0 {
		partition.DataVersion = helper.Int64(int64(v))
	}
	if v, ok := partitionMap["name"].(string); ok && v != "" {
		partition.Name = helper.String(v)
	}
	if v, ok := partitionMap["values"].([]interface{}); ok && len(v) > 0 {
		for _, val := range v {
			if valStr, ok := val.(string); ok && valStr != "" {
				partition.Values = append(partition.Values, helper.String(valStr))
			}
		}
	}
	if v, ok := partitionMap["storage_size"].(int); ok && v != 0 {
		partition.StorageSize = helper.Int64(int64(v))
	}
	if v, ok := partitionMap["record_count"].(int); ok && v != 0 {
		partition.RecordCount = helper.Int64(int64(v))
	}
	if v, ok := partitionMap["create_time"].(string); ok && v != "" {
		partition.CreateTime = helper.String(v)
	}
	if v, ok := partitionMap["modified_time"].(string); ok && v != "" {
		partition.ModifiedTime = helper.String(v)
	}
	if v, ok := partitionMap["last_access_time"].(string); ok && v != "" {
		partition.LastAccessTime = helper.String(v)
	}
	if v, ok := partitionMap["params"].([]interface{}); ok && len(v) > 0 {
		for _, kvItem := range v {
			kvMap := kvItem.(map[string]interface{})
			kvPair := dlc.KVPair{}
			if v, ok := kvMap["key"].(string); ok && v != "" {
				kvPair.Key = helper.String(v)
			}
			if v, ok := kvMap["value"].(string); ok && v != "" {
				kvPair.Value = helper.String(v)
			}
			partition.Params = append(partition.Params, &kvPair)
		}
	}
	if v, ok := partitionMap["sds"].([]interface{}); ok && len(v) > 0 {
		for _, sdsItem := range v {
			sdsMap := sdsItem.(map[string]interface{})
			sds := buildDlcDmsSds(sdsMap)
			partition.Sds = &sds
		}
	}
	if v, ok := partitionMap["datasource_connection_name"].(string); ok && v != "" {
		partition.DatasourceConnectionName = helper.String(v)
	}
	return partition
}

// flattenDlcDmsSds converts a DMSSds struct to a schema map.
func flattenDlcDmsSds(sds *dlc.DMSSds) map[string]interface{} {
	sdsMap := map[string]interface{}{}
	if sds.Location != nil {
		sdsMap["location"] = sds.Location
	}
	if sds.InputFormat != nil {
		sdsMap["input_format"] = sds.InputFormat
	}
	if sds.OutputFormat != nil {
		sdsMap["output_format"] = sds.OutputFormat
	}
	if sds.NumBuckets != nil {
		sdsMap["num_buckets"] = sds.NumBuckets
	}
	if sds.Compressed != nil {
		sdsMap["compressed"] = sds.Compressed
	}
	if sds.StoredAsSubDirectories != nil {
		sdsMap["stored_as_sub_directories"] = sds.StoredAsSubDirectories
	}
	if sds.SerdeLib != nil {
		sdsMap["serde_lib"] = sds.SerdeLib
	}
	if sds.SerdeName != nil {
		sdsMap["serde_name"] = sds.SerdeName
	}
	if sds.BucketCols != nil {
		bucketColsList := make([]interface{}, 0, len(sds.BucketCols))
		for _, col := range sds.BucketCols {
			bucketColsList = append(bucketColsList, col)
		}
		sdsMap["bucket_cols"] = bucketColsList
	}
	if sds.SerdeParams != nil {
		serdeParamsList := make([]map[string]interface{}, 0, len(sds.SerdeParams))
		for _, kvPair := range sds.SerdeParams {
			kvMap := map[string]interface{}{}
			if kvPair.Key != nil {
				kvMap["key"] = kvPair.Key
			}
			if kvPair.Value != nil {
				kvMap["value"] = kvPair.Value
			}
			serdeParamsList = append(serdeParamsList, kvMap)
		}
		sdsMap["serde_params"] = serdeParamsList
	}
	if sds.Params != nil {
		paramsList := make([]map[string]interface{}, 0, len(sds.Params))
		for _, kvPair := range sds.Params {
			kvMap := map[string]interface{}{}
			if kvPair.Key != nil {
				kvMap["key"] = kvPair.Key
			}
			if kvPair.Value != nil {
				kvMap["value"] = kvPair.Value
			}
			paramsList = append(paramsList, kvMap)
		}
		sdsMap["params"] = paramsList
	}
	if sds.SortCols != nil {
		sortColsMap := map[string]interface{}{}
		if sds.SortCols.Col != nil {
			sortColsMap["col"] = sds.SortCols.Col
		}
		if sds.SortCols.Order != nil {
			sortColsMap["order"] = sds.SortCols.Order
		}
		sdsMap["sort_cols"] = []interface{}{sortColsMap}
	}
	if sds.Cols != nil {
		colsList := make([]map[string]interface{}, 0, len(sds.Cols))
		for _, col := range sds.Cols {
			colsList = append(colsList, flattenDlcDmsColumn(col))
		}
		sdsMap["dms_cols"] = colsList
	}
	if sds.SortColumns != nil {
		sortColumnsList := make([]map[string]interface{}, 0, len(sds.SortColumns))
		for _, sortCol := range sds.SortColumns {
			sortColMap := map[string]interface{}{}
			if sortCol.Col != nil {
				sortColMap["col"] = sortCol.Col
			}
			if sortCol.Order != nil {
				sortColMap["order"] = sortCol.Order
			}
			sortColumnsList = append(sortColumnsList, sortColMap)
		}
		sdsMap["sort_columns"] = sortColumnsList
	}
	return sdsMap
}

// flattenDlcDmsColumn converts a DMSColumn struct to a schema map.
func flattenDlcDmsColumn(column *dlc.DMSColumn) map[string]interface{} {
	columnMap := map[string]interface{}{}
	if column.Name != nil {
		columnMap["name"] = column.Name
	}
	if column.Description != nil {
		columnMap["description"] = column.Description
	}
	if column.Type != nil {
		columnMap["type"] = column.Type
	}
	if column.Position != nil {
		columnMap["position"] = column.Position
	}
	if column.Params != nil {
		paramsList := make([]map[string]interface{}, 0, len(column.Params))
		for _, kvPair := range column.Params {
			kvMap := map[string]interface{}{}
			if kvPair.Key != nil {
				kvMap["key"] = kvPair.Key
			}
			if kvPair.Value != nil {
				kvMap["value"] = kvPair.Value
			}
			paramsList = append(paramsList, kvMap)
		}
		columnMap["params"] = paramsList
	}
	if column.BizParams != nil {
		bizParamsList := make([]map[string]interface{}, 0, len(column.BizParams))
		for _, kvPair := range column.BizParams {
			kvMap := map[string]interface{}{}
			if kvPair.Key != nil {
				kvMap["key"] = kvPair.Key
			}
			if kvPair.Value != nil {
				kvMap["value"] = kvPair.Value
			}
			bizParamsList = append(bizParamsList, kvMap)
		}
		columnMap["biz_params"] = bizParamsList
	}
	if column.IsPartition != nil {
		columnMap["is_partition"] = column.IsPartition
	}
	return columnMap
}

// flattenDlcDmsPartition converts a DMSPartition struct to a schema map.
func flattenDlcDmsPartition(partition *dlc.DMSPartition) map[string]interface{} {
	partitionMap := map[string]interface{}{}
	if partition.DatabaseName != nil {
		partitionMap["database_name"] = partition.DatabaseName
	}
	if partition.SchemaName != nil {
		partitionMap["schema_name"] = partition.SchemaName
	}
	if partition.TableName != nil {
		partitionMap["table_name"] = partition.TableName
	}
	if partition.DataVersion != nil {
		partitionMap["data_version"] = partition.DataVersion
	}
	if partition.Name != nil {
		partitionMap["name"] = partition.Name
	}
	if partition.Values != nil {
		valuesList := make([]interface{}, 0, len(partition.Values))
		for _, val := range partition.Values {
			valuesList = append(valuesList, val)
		}
		partitionMap["values"] = valuesList
	}
	if partition.StorageSize != nil {
		partitionMap["storage_size"] = partition.StorageSize
	}
	if partition.RecordCount != nil {
		partitionMap["record_count"] = partition.RecordCount
	}
	if partition.CreateTime != nil {
		partitionMap["create_time"] = partition.CreateTime
	}
	if partition.ModifiedTime != nil {
		partitionMap["modified_time"] = partition.ModifiedTime
	}
	if partition.LastAccessTime != nil {
		partitionMap["last_access_time"] = partition.LastAccessTime
	}
	if partition.Params != nil {
		paramsList := make([]map[string]interface{}, 0, len(partition.Params))
		for _, kvPair := range partition.Params {
			kvMap := map[string]interface{}{}
			if kvPair.Key != nil {
				kvMap["key"] = kvPair.Key
			}
			if kvPair.Value != nil {
				kvMap["value"] = kvPair.Value
			}
			paramsList = append(paramsList, kvMap)
		}
		partitionMap["params"] = paramsList
	}
	if partition.Sds != nil {
		partitionMap["sds"] = []interface{}{flattenDlcDmsSds(partition.Sds)}
	}
	if partition.DatasourceConnectionName != nil {
		partitionMap["datasource_connection_name"] = partition.DatasourceConnectionName
	}
	return partitionMap
}
