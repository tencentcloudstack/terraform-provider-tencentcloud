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
			Create: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"database_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the DLC meta database.",
			},
			"datasource_connection_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Datasource connection name, default `DataLakeCatalog`.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the DLC meta database, length 0~2048.",
			},
			"govern_policy": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Data governance config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Govern rule type, `Customize` or `Intelligence`.",
						},
						"govern_engine": {
							Type:        schema.TypeString,
							Optional:    true,
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
			"batch_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Batch ID of the async task.",
			},
			"task_id_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Task ID set.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
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
			"database_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Database ID.",
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
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request = dlc.NewCreateMetaDatabaseRequest()
	)

	databaseName := d.Get("database_name").(string)
	metaDatabaseInfo := &dlc.MetaDatabaseInfo{
		DatabaseName: helper.String(databaseName),
	}
	if v, ok := d.GetOk("comment"); ok {
		metaDatabaseInfo.Comment = helper.String(v.(string))
	}
	request.MetaDatabaseInfo = metaDatabaseInfo

	if v, ok := d.GetOk("datasource_connection_name"); ok {
		request.DatasourceConnectionName = helper.String(v.(string))
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

	response := dlc.NewCreateMetaDatabaseResponse()
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().CreateMetaDatabaseWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dlc_meta_database failed, Response is nil."))
		}

		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dlc_meta_database failed, reason:%+v", logId, err)
		return err
	}

	log.Printf("[DEBUG]%s dlc_meta_database create response, logId=%s, d.Id()=%s", logId, logId, d.Id())

	if response.Response.BatchId == nil {
		return fmt.Errorf("dlc_meta_database BatchId is nil, request body [%s]", request.ToJsonString())
	}

	batchId := *response.Response.BatchId
	log.Printf("[DEBUG]%s dlc_meta_database BatchId=%s", logId, batchId)
	_ = d.Set("batch_id", batchId)
	if len(response.Response.TaskIdSet) > 0 {
		taskIdSet := make([]interface{}, 0, len(response.Response.TaskIdSet))
		for _, taskId := range response.Response.TaskIdSet {
			if taskId != nil {
				taskIdSet = append(taskIdSet, *taskId)
			}
		}
		_ = d.Set("task_id_set", taskIdSet)
	}

	var datasourceConnectionName string
	if v, ok := d.GetOk("datasource_connection_name"); ok {
		datasourceConnectionName = v.(string)
	}
	if datasourceConnectionName != "" {
		d.SetId(strings.Join([]string{datasourceConnectionName, databaseName}, tccommon.FILED_SP))
	} else {
		d.SetId(databaseName)
	}

	// async poll: wait until the database is queryable
	describeRequest := dlc.NewDescribeDatabaseRequest()
	describeRequest.DatabaseName = helper.String(databaseName)
	if datasourceConnectionName != "" {
		describeRequest.DatasourceConnectionName = helper.String(datasourceConnectionName)
	}
	waitErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeDatabaseWithContext(ctx, describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil || result.Response.DatabaseInfo == nil {
			e = fmt.Errorf("[DEBUG]%s api[%s] dlc_meta_database not ready, request body [%s]", logId, describeRequest.GetAction(), describeRequest.ToJsonString())
			log.Println(e)
			return resource.RetryableError(e)
		}
		return nil
	})
	if waitErr != nil {
		log.Printf("[CRITAL]%s create dlc_meta_database wait failed, reason:%+v", logId, waitErr)
		return waitErr
	}

	return resourceTencentCloudDlcMetaDatabaseRead(d, meta)
}

func resourceTencentCloudDlcMetaDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_meta_database.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	var datasourceConnectionName, databaseName string
	if len(idSplit) == 2 {
		datasourceConnectionName = idSplit[0]
		databaseName = idSplit[1]
	} else if len(idSplit) == 1 {
		databaseName = idSplit[0]
	} else {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	request := dlc.NewDescribeDatabaseRequest()
	request.DatabaseName = helper.String(databaseName)
	if datasourceConnectionName != "" {
		request.DatasourceConnectionName = helper.String(datasourceConnectionName)
	}

	response := dlc.NewDescribeDatabaseResponse()
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeDatabaseWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc_meta_database failed, Response is nil."))
		}

		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read dlc_meta_database failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.DatabaseInfo == nil {
		log.Printf("[CRUD] dlc_meta_database id=%s", d.Id())
		d.SetId("")
		return nil
	}

	info := response.Response.DatabaseInfo

	if info.DatabaseName != nil {
		_ = d.Set("database_name", *info.DatabaseName)
	}
	if info.Comment != nil {
		_ = d.Set("comment", *info.Comment)
	}
	if info.Location != nil {
		_ = d.Set("location", *info.Location)
	}
	if info.CreateTime != nil {
		_ = d.Set("create_time", *info.CreateTime)
	}
	if info.ModifiedTime != nil {
		_ = d.Set("modified_time", *info.ModifiedTime)
	}
	if info.UserAlias != nil {
		_ = d.Set("user_alias", *info.UserAlias)
	}
	if info.UserSubUin != nil {
		_ = d.Set("user_sub_uin", *info.UserSubUin)
	}
	if info.DatabaseId != nil {
		_ = d.Set("database_id", *info.DatabaseId)
	}
	if info.CatalogName != nil {
		_ = d.Set("catalog_name", *info.CatalogName)
	}
	if info.CatalogType != nil {
		_ = d.Set("catalog_type", *info.CatalogType)
	}
	if info.IsInformationSchema != nil {
		_ = d.Set("is_information_schema", *info.IsInformationSchema)
	}

	if info.Properties != nil {
		propertiesList := make([]interface{}, 0, len(info.Properties))
		for _, p := range info.Properties {
			m := map[string]interface{}{}
			if p.Key != nil {
				m["key"] = *p.Key
			}
			if p.Value != nil {
				m["value"] = *p.Value
			}
			propertiesList = append(propertiesList, m)
		}
		_ = d.Set("properties", propertiesList)
	}

	if info.GovernPolicy != nil {
		governPolicyList := flattenDlcMetaDatabaseGovernPolicy(info.GovernPolicy)
		_ = d.Set("govern_policy", governPolicyList)
	}

	// datasource_connection_name is not returned by DescribeDatabase, keep the value from config/state
	if datasourceConnectionName != "" {
		_ = d.Set("datasource_connection_name", datasourceConnectionName)
	}

	return nil
}

func resourceTencentCloudDlcMetaDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_meta_database.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"datasource_connection_name", "comment", "govern_policy", "smart_policy"}
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
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	var datasourceConnectionName, databaseName string
	if len(idSplit) == 2 {
		datasourceConnectionName = idSplit[0]
		databaseName = idSplit[1]
	} else if len(idSplit) == 1 {
		databaseName = idSplit[0]
	} else {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	request := dlc.NewDeleteMetaDatabaseRequest()
	request.DatabaseName = helper.String(databaseName)
	if datasourceConnectionName != "" {
		request.DatasourceConnectionName = helper.String(datasourceConnectionName)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DeleteMetaDatabaseWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete dlc_meta_database failed, Response is nil."))
		}

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete dlc_meta_database failed, reason:%+v", logId, err)
		return err
	}

	// async poll: wait until the database is gone
	describeRequest := dlc.NewDescribeDatabaseRequest()
	describeRequest.DatabaseName = helper.String(databaseName)
	if datasourceConnectionName != "" {
		describeRequest.DatasourceConnectionName = helper.String(datasourceConnectionName)
	}
	waitErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeDatabaseWithContext(ctx, describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil || result.Response.DatabaseInfo == nil {
			return nil
		}
		e = fmt.Errorf("[DEBUG]%s api[%s] dlc_meta_database still exists", logId, describeRequest.GetAction())
		log.Println(e)
		return resource.RetryableError(e)
	})
	if waitErr != nil {
		log.Printf("[CRITAL]%s delete dlc_meta_database wait failed, reason:%+v", logId, waitErr)
		return waitErr
	}

	return nil
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
