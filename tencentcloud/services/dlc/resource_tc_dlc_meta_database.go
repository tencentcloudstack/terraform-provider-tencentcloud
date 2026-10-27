package dlc

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcMetaDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcMetaDatabaseCreate,
		Read:   resourceTencentCloudDlcMetaDatabaseRead,
		Update: resourceTencentCloudDlcMetaDatabaseUpdate,
		Delete: resourceTencentCloudDlcMetaDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"datasource_connection_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Datasource connection name, default `DataLakeCatalog`.",
			},
			"meta_database_info": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Meta database basic information, which corresponds to the `MetaDatabaseInfo` structure of the `CreateMetaDatabase` API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Name of the DLC meta database, length 0~128, digits, letters and underscores are supported, cannot start with a digit, and is converted to lowercase.",
						},
						"comment": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of the DLC meta database, length 0~2048.",
						},
					},
				},
			},
			"govern_policy": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Data governance config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Govern rule type, `Customize` or `Intelligence`.",
						},
						"govern_engine": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Govern engine.",
						},
					},
				},
			},
			"smart_policy": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Smart data governance config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"base_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Base info of smart policy.",
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
										Description: "Policy type.",
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
										Description: "User AppID.",
									},
								},
							},
						},
						"policy": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Policy description of smart policy.",
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
												"favor": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Affinity info.",
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
												"resource_conf": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Resource config info.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"parallelism": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Parallelism of the optimize task.",
															},
														},
													},
												},
											},
										},
									},
									"written": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Data rewrite policy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"written_enable": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Whether to enable, `none`/`enable`/`disable`/`default`.",
												},
												"advance_policy": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "User custom advance params.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"compact_enable": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Whether to enable compact.",
															},
															"delete_enable": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Whether to enable history data cleanup.",
															},
															"cow_compact_enable": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Whether to enable COW table compact.",
															},
															"compact_strategy": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "File compact strategy.",
															},
															"min_input_files": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Min input files to compact.",
															},
															"target_file_size_bytes": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Target file size after compact.",
															},
															"retain_last": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Retain last snapshot count.",
															},
															"before_days": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Snapshot expire before days.",
															},
															"expired_snapshots_interval_min": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Snapshot expire interval in minutes.",
															},
															"remove_orphan_interval_min": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Remove orphan files interval in minutes.",
															},
															"sort_orders": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Sort compact strategy rules.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"column": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Sort column name.",
																		},
																		"sort_direction": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Sort direction, `asc` or `desc`.",
																		},
																		"null_order": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Null order, `first` or `last`.",
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
										MaxItems:    1,
										Optional:    true,
										Description: "Data lifecycle policy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"lifecycle_enable": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Whether to enable lifecycle.",
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
												"expiration": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Expiration time.",
												},
												"drop_table": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to drop table (deprecated, use table_expiration instead).",
												},
											},
										},
									},
									"index": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
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
										MaxItems:    1,
										Optional:    true,
										Description: "Change table policy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"data_retention_time": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Data retention time of change table, in days.",
												},
											},
										},
									},
									"table_expiration": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Table expiration policy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Whether to enable the policy.",
												},
												"expiration": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Table expiration time, in days.",
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
			// computed
			"properties": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Properties of the database.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Property key.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Property value.",
						},
					},
				},
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Database create time, in seconds.",
			},
			"modified_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Database modified time, in seconds.",
			},
			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "COS storage path.",
			},
			"user_alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User alias who created the database.",
			},
			"user_sub_uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User sub UIN who created the database.",
			},
			"catalog_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Catalog name.",
			},
			"catalog_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Catalog type.",
			},
			"is_information_schema": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether it is InformationSchema.",
			},
		},
	}
}

func resourceTencentCloudDlcMetaDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_meta_database.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                    = tccommon.GetLogId(tccommon.ContextNil)
		ctx                      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request                  = dlc.NewCreateMetaDatabaseRequest()
		datasourceConnectionName string
		databaseName             string
	)

	if v, ok := d.GetOk("datasource_connection_name"); ok {
		request.DatasourceConnectionName = helper.String(v.(string))
		datasourceConnectionName = v.(string)
	}

	if v, ok := d.GetOk("meta_database_info"); ok {
		request.MetaDatabaseInfo = buildDlcMetaDatabaseInfo(v.([]interface{}))
	}

	if request.MetaDatabaseInfo == nil || request.MetaDatabaseInfo.DatabaseName == nil {
		return fmt.Errorf("meta_database_info.database_name is required.")
	} else {
		databaseName = *request.MetaDatabaseInfo.DatabaseName
	}

	if v, ok := d.GetOk("govern_policy"); ok {
		governPolicy := buildDlcMetaDatabaseGovernPolicy(v.([]interface{}))
		if governPolicy != nil {
			request.GovernPolicy = governPolicy
		}
	}

	if v, ok := d.GetOk("smart_policy"); ok {
		smartPolicy := buildDlcMetaDatabaseSmartPolicy(v.([]interface{}))
		if smartPolicy != nil {
			request.SmartPolicy = smartPolicy
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().CreateMetaDatabaseWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dlc meta database failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dlc meta database failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join([]string{datasourceConnectionName, databaseName}, tccommon.FILED_SP))

	// wait 10 seconds for database sync
	time.Sleep(10 * time.Second)

	// The create interface is asynchronous, wait until the meta database is
	// visible before reading it back.
	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	waitErr := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		respData, e := service.DescribeDlcMetaDatabaseByName(ctx, datasourceConnectionName, databaseName)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if respData == nil {
			return resource.RetryableError(fmt.Errorf("waiting for dlc meta database %s to be created.", d.Id()))
		}

		return nil
	})

	if waitErr != nil {
		log.Printf("[CRITAL]%s create dlc meta database failed, reason:%+v", logId, waitErr)
		return waitErr
	}

	return resourceTencentCloudDlcMetaDatabaseRead(d, meta)
}

func resourceTencentCloudDlcMetaDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_meta_database.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	datasourceConnectionName, databaseName, err := parseDlcMetaDatabaseId(d.Id())
	if err != nil {
		return err
	}

	respData, err := service.DescribeDlcMetaDatabaseByName(ctx, datasourceConnectionName, databaseName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_dlc_meta_database` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("datasource_connection_name", datasourceConnectionName)

	comment := d.Get("meta_database_info.0.comment").(string)
	if respData.Comment != nil {
		comment = *respData.Comment
	}

	metaDatabaseInfo := map[string]interface{}{
		"database_name": databaseName,
		"comment":       comment,
	}

	_ = d.Set("meta_database_info", []interface{}{metaDatabaseInfo})

	if respData.GovernPolicy != nil {
		_ = d.Set("govern_policy", flattenDlcMetaDatabaseGovernPolicy(respData.GovernPolicy))
	}

	if respData.Properties != nil {
		propertiesList := make([]interface{}, 0, len(respData.Properties))
		for _, property := range respData.Properties {
			propertiesMap := map[string]interface{}{}
			if property.Key != nil {
				propertiesMap["key"] = *property.Key
			}

			if property.Value != nil {
				propertiesMap["value"] = *property.Value
			}

			propertiesList = append(propertiesList, propertiesMap)
		}

		_ = d.Set("properties", propertiesList)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.ModifiedTime != nil {
		_ = d.Set("modified_time", respData.ModifiedTime)
	}

	if respData.Location != nil {
		_ = d.Set("location", respData.Location)
	}

	if respData.UserAlias != nil {
		_ = d.Set("user_alias", respData.UserAlias)
	}

	if respData.UserSubUin != nil {
		_ = d.Set("user_sub_uin", respData.UserSubUin)
	}

	if respData.CatalogName != nil {
		_ = d.Set("catalog_name", respData.CatalogName)
	}

	if respData.CatalogType != nil {
		_ = d.Set("catalog_type", respData.CatalogType)
	}

	if respData.IsInformationSchema != nil {
		_ = d.Set("is_information_schema", respData.IsInformationSchema)
	}

	return nil
}

func resourceTencentCloudDlcMetaDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_meta_database.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	// The DLC meta database only provides create/describe/delete interfaces and
	// the resource id is composed of the datasource connection name and the
	// database name, so all the arguments are immutable.
	immutableArgs := []string{"govern_policy", "smart_policy"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	return resourceTencentCloudDlcMetaDatabaseRead(d, meta)
}

func resourceTencentCloudDlcMetaDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_meta_database.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = dlc.NewDeleteMetaDatabaseRequest()
	)

	datasourceConnectionName, databaseName, err := parseDlcMetaDatabaseId(d.Id())
	if err != nil {
		return err
	}

	request.DatabaseName = helper.String(databaseName)
	request.DatasourceConnectionName = helper.String(datasourceConnectionName)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DeleteMetaDatabaseWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete dlc meta database failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete dlc meta database failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// The delete interface is asynchronous, wait until the meta database is gone.
	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	waitErr := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		respData, e := service.DescribeDlcMetaDatabaseByName(ctx, datasourceConnectionName, databaseName)
		if e != nil {
			if tccommon.IsExpectError(e, []string{"ResourceNotFound", "FailedOperation.DatabaseNotExist"}) {
				return nil
			}

			return tccommon.RetryError(e)
		}

		if respData != nil {
			return resource.RetryableError(fmt.Errorf("waiting for dlc meta database %s to be deleted.", d.Id()))
		}

		return nil
	})

	if waitErr != nil {
		log.Printf("[CRITAL]%s delete dlc meta database failed, reason:%+v", logId, waitErr)
		return waitErr
	}

	return nil
}

// parseDlcMetaDatabaseId splits the resource id `datasourceConnectionName#databaseName`
// into the datasource connection name and the meta database name.
func parseDlcMetaDatabaseId(id string) (datasourceConnectionName, databaseName string, err error) {
	idSplit := strings.Split(id, tccommon.FILED_SP)
	if len(idSplit) != 2 || idSplit[0] == "" || idSplit[1] == "" {
		err = fmt.Errorf("id is broken,%s", id)
		return
	}

	datasourceConnectionName = idSplit[0]
	databaseName = idSplit[1]
	return
}

// buildDlcMetaDatabaseInfo maps the `meta_database_info` block onto the
// `MetaDatabaseInfo` structure of the CreateMetaDatabase API.
func buildDlcMetaDatabaseInfo(list []interface{}) *dlc.MetaDatabaseInfo {
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	m := list[0].(map[string]interface{})
	info := &dlc.MetaDatabaseInfo{}
	if v, ok := m["database_name"]; ok && v.(string) != "" {
		info.DatabaseName = helper.String(v.(string))
	}
	if v, ok := m["comment"]; ok && v.(string) != "" {
		info.Comment = helper.String(v.(string))
	}
	return info
}

func buildDlcMetaDatabaseGovernPolicy(list []interface{}) *dlc.DataGovernPolicy {
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	m := list[0].(map[string]interface{})
	policy := &dlc.DataGovernPolicy{}
	if v, ok := m["rule_type"]; ok && v.(string) != "" {
		policy.RuleType = helper.String(v.(string))
	}
	if v, ok := m["govern_engine"]; ok && v.(string) != "" {
		policy.GovernEngine = helper.String(v.(string))
	}
	return policy
}

func flattenDlcMetaDatabaseGovernPolicy(policy *dlc.DataGovernPolicy) []interface{} {
	if policy == nil {
		return []interface{}{}
	}
	m := map[string]interface{}{}
	if policy.RuleType != nil {
		m["rule_type"] = *policy.RuleType
	}
	if policy.GovernEngine != nil {
		m["govern_engine"] = *policy.GovernEngine
	}
	return []interface{}{m}
}

func buildDlcMetaDatabaseSmartPolicy(list []interface{}) *dlc.SmartPolicy {
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	m := list[0].(map[string]interface{})
	policy := &dlc.SmartPolicy{}

	if v, ok := m["base_info"]; ok {
		baseInfo := buildDlcMetaDatabaseSmartPolicyBaseInfo(v.([]interface{}))
		if baseInfo != nil {
			policy.BaseInfo = baseInfo
		}
	}
	if v, ok := m["policy"]; ok {
		p := buildDlcMetaDatabaseSmartOptimizerPolicy(v.([]interface{}))
		if p != nil {
			policy.Policy = p
		}
	}
	return policy
}

func buildDlcMetaDatabaseSmartPolicyBaseInfo(list []interface{}) *dlc.SmartPolicyBaseInfo {
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	m := list[0].(map[string]interface{})
	info := &dlc.SmartPolicyBaseInfo{}
	if v, ok := m["uin"]; ok {
		info.Uin = helper.String(v.(string))
	}
	if v, ok := m["policy_type"]; ok && v.(string) != "" {
		info.PolicyType = helper.String(v.(string))
	}
	if v, ok := m["catalog"]; ok && v.(string) != "" {
		info.Catalog = helper.String(v.(string))
	}
	if v, ok := m["database"]; ok && v.(string) != "" {
		info.Database = helper.String(v.(string))
	}
	if v, ok := m["table"]; ok && v.(string) != "" {
		info.Table = helper.String(v.(string))
	}
	if v, ok := m["app_id"]; ok && v.(string) != "" {
		info.AppId = helper.String(v.(string))
	}
	return info
}

func buildDlcMetaDatabaseSmartOptimizerPolicy(list []interface{}) *dlc.SmartOptimizerPolicy {
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	m := list[0].(map[string]interface{})
	policy := &dlc.SmartOptimizerPolicy{}
	if v, ok := m["inherit"]; ok && v.(string) != "" {
		policy.Inherit = helper.String(v.(string))
	}
	if v, ok := m["resources"]; ok {
		policy.Resources = buildDlcMetaDatabaseResourceInfoList(v.([]interface{}))
	}
	if v, ok := m["written"]; ok {
		written := buildDlcMetaDatabaseSmartOptimizerWrittenPolicy(v.([]interface{}))
		if written != nil {
			policy.Written = written
		}
	}
	if v, ok := m["lifecycle"]; ok {
		lifecycle := buildDlcMetaDatabaseSmartOptimizerLifecyclePolicy(v.([]interface{}))
		if lifecycle != nil {
			policy.Lifecycle = lifecycle
		}
	}
	if v, ok := m["index"]; ok {
		index := buildDlcMetaDatabaseSmartOptimizerIndexPolicy(v.([]interface{}))
		if index != nil {
			policy.Index = index
		}
	}
	if v, ok := m["change_table"]; ok {
		changeTable := buildDlcMetaDatabaseSmartOptimizerChangeTablePolicy(v.([]interface{}))
		if changeTable != nil {
			policy.ChangeTable = changeTable
		}
	}
	if v, ok := m["table_expiration"]; ok {
		tableExpiration := buildDlcMetaDatabaseTableExpirationPolicy(v.([]interface{}))
		if tableExpiration != nil {
			policy.TableExpiration = tableExpiration
		}
	}
	return policy
}

func buildDlcMetaDatabaseResourceInfoList(list []interface{}) []*dlc.ResourceInfo {
	if len(list) == 0 {
		return nil
	}
	result := make([]*dlc.ResourceInfo, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		m := item.(map[string]interface{})
		info := &dlc.ResourceInfo{}
		if v, ok := m["attribution_type"]; ok && v.(string) != "" {
			info.AttributionType = helper.String(v.(string))
		}
		if v, ok := m["resource_type"]; ok && v.(string) != "" {
			info.ResourceType = helper.String(v.(string))
		}
		if v, ok := m["name"]; ok && v.(string) != "" {
			info.Name = helper.String(v.(string))
		}
		if v, ok := m["instance"]; ok && v.(string) != "" {
			info.Instance = helper.String(v.(string))
		}
		if v, ok := m["favor"]; ok {
			info.Favor = buildDlcMetaDatabaseFavorInfoList(v.([]interface{}))
		}
		if v, ok := m["status"]; ok {
			info.Status = helper.IntInt64(v.(int))
		}
		if v, ok := m["resource_group_name"]; ok && v.(string) != "" {
			info.ResourceGroupName = helper.String(v.(string))
		}
		if v, ok := m["resource_conf"]; ok {
			conf := buildDlcMetaDatabaseResourceConf(v.([]interface{}))
			if conf != nil {
				info.ResourceConf = conf
			}
		}
		result = append(result, info)
	}
	return result
}

func buildDlcMetaDatabaseFavorInfoList(list []interface{}) []*dlc.FavorInfo {
	if len(list) == 0 {
		return nil
	}
	result := make([]*dlc.FavorInfo, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		m := item.(map[string]interface{})
		info := &dlc.FavorInfo{}
		if v, ok := m["priority"]; ok {
			info.Priority = helper.IntInt64(v.(int))
		}
		if v, ok := m["catalog"]; ok && v.(string) != "" {
			info.Catalog = helper.String(v.(string))
		}
		if v, ok := m["data_base"]; ok && v.(string) != "" {
			info.DataBase = helper.String(v.(string))
		}
		if v, ok := m["table"]; ok && v.(string) != "" {
			info.Table = helper.String(v.(string))
		}
		result = append(result, info)
	}
	return result
}

func buildDlcMetaDatabaseResourceConf(list []interface{}) *dlc.ResourceConf {
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	m := list[0].(map[string]interface{})
	conf := &dlc.ResourceConf{}
	if v, ok := m["parallelism"]; ok {
		conf.Parallelism = helper.IntInt64(v.(int))
	}
	return conf
}

func buildDlcMetaDatabaseSmartOptimizerWrittenPolicy(list []interface{}) *dlc.SmartOptimizerWrittenPolicy {
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	m := list[0].(map[string]interface{})
	policy := &dlc.SmartOptimizerWrittenPolicy{}
	if v, ok := m["written_enable"]; ok && v.(string) != "" {
		policy.WrittenEnable = helper.String(v.(string))
	}
	if v, ok := m["advance_policy"]; ok {
		advance := buildDlcMetaDatabaseWrittenAdvancePolicy(v.([]interface{}))
		if advance != nil {
			policy.AdvancePolicy = advance
		}
	}
	return policy
}

func buildDlcMetaDatabaseWrittenAdvancePolicy(list []interface{}) *dlc.WrittenAdvancePolicy {
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	m := list[0].(map[string]interface{})
	policy := &dlc.WrittenAdvancePolicy{}
	if v, ok := m["compact_enable"]; ok && v.(string) != "" {
		policy.CompactEnable = helper.String(v.(string))
	}
	if v, ok := m["delete_enable"]; ok && v.(string) != "" {
		policy.DeleteEnable = helper.String(v.(string))
	}
	if v, ok := m["cow_compact_enable"]; ok && v.(string) != "" {
		policy.CowCompactEnable = helper.String(v.(string))
	}
	if v, ok := m["compact_strategy"]; ok && v.(string) != "" {
		policy.CompactStrategy = helper.String(v.(string))
	}
	if v, ok := m["min_input_files"]; ok {
		policy.MinInputFiles = helper.IntInt64(v.(int))
	}
	if v, ok := m["target_file_size_bytes"]; ok {
		policy.TargetFileSizeBytes = helper.IntInt64(v.(int))
	}
	if v, ok := m["retain_last"]; ok {
		policy.RetainLast = helper.IntInt64(v.(int))
	}
	if v, ok := m["before_days"]; ok {
		policy.BeforeDays = helper.IntInt64(v.(int))
	}
	if v, ok := m["expired_snapshots_interval_min"]; ok {
		policy.ExpiredSnapshotsIntervalMin = helper.IntInt64(v.(int))
	}
	if v, ok := m["remove_orphan_interval_min"]; ok {
		policy.RemoveOrphanIntervalMin = helper.IntInt64(v.(int))
	}
	if v, ok := m["sort_orders"]; ok {
		policy.SortOrders = buildDlcMetaDatabaseSortOrderList(v.([]interface{}))
	}
	return policy
}

func buildDlcMetaDatabaseSortOrderList(list []interface{}) []*dlc.SortOrder {
	if len(list) == 0 {
		return nil
	}
	result := make([]*dlc.SortOrder, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		m := item.(map[string]interface{})
		order := &dlc.SortOrder{}
		if v, ok := m["column"]; ok && v.(string) != "" {
			order.Column = helper.String(v.(string))
		}
		if v, ok := m["sort_direction"]; ok && v.(string) != "" {
			order.SortDirection = helper.String(v.(string))
		}
		if v, ok := m["null_order"]; ok && v.(string) != "" {
			order.NullOrder = helper.String(v.(string))
		}
		result = append(result, order)
	}
	return result
}

func buildDlcMetaDatabaseSmartOptimizerLifecyclePolicy(list []interface{}) *dlc.SmartOptimizerLifecyclePolicy {
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	m := list[0].(map[string]interface{})
	policy := &dlc.SmartOptimizerLifecyclePolicy{}
	if v, ok := m["lifecycle_enable"]; ok && v.(string) != "" {
		policy.LifecycleEnable = helper.String(v.(string))
	}
	if v, ok := m["expired_field"]; ok && v.(string) != "" {
		policy.ExpiredField = helper.String(v.(string))
	}
	if v, ok := m["expired_field_format"]; ok && v.(string) != "" {
		policy.ExpiredFieldFormat = helper.String(v.(string))
	}
	if v, ok := m["expiration"]; ok {
		policy.Expiration = helper.IntInt64(v.(int))
	}
	if v, ok := m["drop_table"]; ok {
		policy.DropTable = helper.Bool(v.(bool))
	}
	return policy
}

func buildDlcMetaDatabaseSmartOptimizerIndexPolicy(list []interface{}) *dlc.SmartOptimizerIndexPolicy {
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	m := list[0].(map[string]interface{})
	policy := &dlc.SmartOptimizerIndexPolicy{}
	if v, ok := m["index_enable"]; ok && v.(string) != "" {
		policy.IndexEnable = helper.String(v.(string))
	}
	return policy
}

func buildDlcMetaDatabaseSmartOptimizerChangeTablePolicy(list []interface{}) *dlc.SmartOptimizerChangeTablePolicy {
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	m := list[0].(map[string]interface{})
	policy := &dlc.SmartOptimizerChangeTablePolicy{}
	if v, ok := m["data_retention_time"]; ok {
		policy.DataRetentionTime = helper.IntInt64(v.(int))
	}
	return policy
}

func buildDlcMetaDatabaseTableExpirationPolicy(list []interface{}) *dlc.TableExpirationPolicy {
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	m := list[0].(map[string]interface{})
	policy := &dlc.TableExpirationPolicy{}
	if v, ok := m["enabled"]; ok {
		policy.Enabled = helper.Bool(v.(bool))
	}
	if v, ok := m["expiration"]; ok {
		policy.Expiration = helper.Uint64(uint64(v.(int)))
	}
	return policy
}
