package cls

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsDataTransform() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsDataTransformCreate,
		Read:   resourceTencentCloudClsDataTransformRead,
		Update: resourceTencentCloudClsDataTransformUpdate,
		Delete: resourceTencentCloudClsDataTransformDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"func_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Task type. `1`: Specify the theme; `2`: Dynamic creation.",
			},

			"src_topic_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Source topic ID.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Task name.",
			},

			"etl_content": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Data transform content. If `func_type` is `2`, must use `log_auto_output`.",
			},

			"task_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Task type. `1`: Use random data from the source log theme for processing preview; `2`: Use user-defined test data for processing preview; `3`: Create real machining tasks.",
			},

			"enable_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Task enable flag. `1`: enable, `2`: disable, Default is `1`.",
			},

			"dst_resources": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Data transform des resources. If `func_type` is `1`, this parameter is required. If `func_type` is `2`, this parameter does not need to be filled in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Dst topic ID.",
						},
						"alias": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Alias.",
						},
					},
				},
			},

			"backup_give_up_data": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "When `func_type` is `2`, whether to discard data when the number of dynamically created logsets and topics exceeds the product specification limit. Default is `false`. `false`: Create backup logset and topic and write logs to the backup topic; `true`: Discard log data.",
			},

			"has_services_log": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Whether to enable service log delivery. `1`: disable; `2`: enable.",
			},

			"data_transform_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Data transform type. `0`: Standard data transform task; `1`: Pre-processing data transform task (process collected logs before writing to the log topic).",
			},

			"keep_failure_log": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Keep failure log status. `1`: do not keep (default); `2`: keep.",
			},

			"failure_log_key": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Field name for failure logs.",
			},

			"process_from_timestamp": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Specify the start time of processing data, in seconds-level timestamp. Any time range within the log topic lifecycle. If it exceeds the lifecycle, only the part with data within the lifecycle is processed.",
			},

			"process_to_timestamp": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Specify the end time of processing data, in seconds-level timestamp. Cannot specify a future time. If not filled, it means continuous execution.",
			},

			"data_transform_sql_data_sources": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Associated data source information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_source": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Data source type. `1`: MySQL; `2`: Self-built MySQL; `3`: PostgreSQL.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "InstanceId region. For example: ap-guangzhou.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance ID. When DataSource is `1`, it represents the cloud database MySQL instance ID, such as: cdb-zxcvbnm.",
						},
						"user": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "MySQL access username.",
						},
						"alias_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Alias. Used in data transform statements.",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "MySQL access password.",
						},
					},
				},
			},

			"env_infos": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Set environment variables.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Environment variable name.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Environment variable value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClsDataTransformCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_data_transform.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = cls.NewCreateDataTransformRequest()
		response = cls.NewCreateDataTransformResponse()
		taskId   string
	)

	if v, ok := d.GetOkExists("func_type"); ok {
		request.FuncType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("src_topic_id"); ok {
		request.SrcTopicId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("etl_content"); ok {
		request.EtlContent = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("task_type"); ok {
		request.TaskType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("enable_flag"); ok {
		request.EnableFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("dst_resources"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dataTransformResouceInfo := cls.DataTransformResouceInfo{}
			if v, ok := dMap["topic_id"]; ok {
				dataTransformResouceInfo.TopicId = helper.String(v.(string))
			}

			if v, ok := dMap["alias"]; ok {
				dataTransformResouceInfo.Alias = helper.String(v.(string))
			}

			request.DstResources = append(request.DstResources, &dataTransformResouceInfo)
		}
	}

	if v, ok := d.GetOkExists("backup_give_up_data"); ok {
		request.BackupGiveUpData = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("has_services_log"); ok {
		request.HasServicesLog = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("data_transform_type"); ok {
		request.DataTransformType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("keep_failure_log"); ok {
		request.KeepFailureLog = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("failure_log_key"); ok {
		request.FailureLogKey = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("process_from_timestamp"); ok {
		request.ProcessFromTimestamp = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("process_to_timestamp"); ok {
		request.ProcessToTimestamp = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("data_transform_sql_data_sources"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dataTransformSqlDataSource := cls.DataTransformSqlDataSource{}
			if v, ok := dMap["data_source"]; ok {
				dataTransformSqlDataSource.DataSource = helper.IntUint64(v.(int))
			}

			if v, ok := dMap["region"]; ok {
				dataTransformSqlDataSource.Region = helper.String(v.(string))
			}

			if v, ok := dMap["instance_id"]; ok {
				dataTransformSqlDataSource.InstanceId = helper.String(v.(string))
			}

			if v, ok := dMap["user"]; ok {
				dataTransformSqlDataSource.User = helper.String(v.(string))
			}

			if v, ok := dMap["alias_name"]; ok {
				dataTransformSqlDataSource.AliasName = helper.String(v.(string))
			}

			if v, ok := dMap["password"]; ok {
				dataTransformSqlDataSource.Password = helper.String(v.(string))
			}

			request.DataTransformSqlDataSources = append(request.DataTransformSqlDataSources, &dataTransformSqlDataSource)
		}
	}

	if v, ok := d.GetOk("env_infos"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			envInfo := cls.EnvInfo{}
			if v, ok := dMap["key"]; ok {
				envInfo.Key = helper.String(v.(string))
			}

			if v, ok := dMap["value"]; ok {
				envInfo.Value = helper.String(v.(string))
			}

			request.EnvInfos = append(request.EnvInfos, &envInfo)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateDataTransform(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls dataTransform failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	return resourceTencentCloudClsDataTransformRead(d, meta)
}

func resourceTencentCloudClsDataTransformRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_data_transform.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		ctx                 = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service             = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dataTransformTaskId = d.Id()
	)

	dataTransform, err := service.DescribeClsDataTransformById(ctx, dataTransformTaskId)
	if err != nil {
		return err
	}

	if dataTransform == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cls_data_transform` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if dataTransform.SrcTopicId != nil {
		_ = d.Set("src_topic_id", dataTransform.SrcTopicId)
	}

	if dataTransform.Name != nil {
		_ = d.Set("name", dataTransform.Name)
	}

	if dataTransform.EtlContent != nil {
		_ = d.Set("etl_content", dataTransform.EtlContent)
	}

	if dataTransform.EnableFlag != nil {
		_ = d.Set("enable_flag", dataTransform.EnableFlag)
	}

	if dataTransform.DstResources != nil {
		var dstResourcesList []interface{}
		for _, dstResources := range dataTransform.DstResources {
			dstResourcesMap := map[string]interface{}{}

			if dstResources.TopicId != nil {
				dstResourcesMap["topic_id"] = dstResources.TopicId
			}

			if dstResources.Alias != nil {
				dstResourcesMap["alias"] = dstResources.Alias
			}

			dstResourcesList = append(dstResourcesList, dstResourcesMap)
		}

		_ = d.Set("dst_resources", dstResourcesList)
	}

	if dataTransform.BackupGiveUpData != nil {
		_ = d.Set("backup_give_up_data", dataTransform.BackupGiveUpData)
	}

	if dataTransform.HasServicesLog != nil {
		_ = d.Set("has_services_log", dataTransform.HasServicesLog)
	}

	if dataTransform.DataTransformType != nil {
		_ = d.Set("data_transform_type", dataTransform.DataTransformType)
	}

	if dataTransform.KeepFailureLog != nil {
		_ = d.Set("keep_failure_log", dataTransform.KeepFailureLog)
	}

	if dataTransform.FailureLogKey != nil {
		_ = d.Set("failure_log_key", dataTransform.FailureLogKey)
	}

	if dataTransform.ProcessFromTimestamp != nil {
		_ = d.Set("process_from_timestamp", dataTransform.ProcessFromTimestamp)
	}

	if dataTransform.ProcessToTimestamp != nil {
		_ = d.Set("process_to_timestamp", dataTransform.ProcessToTimestamp)
	}

	if dataTransform.DataTransformSqlDataSources != nil {
		var dataSourcesList []interface{}
		for _, dataSource := range dataTransform.DataTransformSqlDataSources {
			dataSourceMap := map[string]interface{}{}

			if dataSource.DataSource != nil {
				dataSourceMap["data_source"] = dataSource.DataSource
			}

			if dataSource.Region != nil {
				dataSourceMap["region"] = dataSource.Region
			}

			if dataSource.InstanceId != nil {
				dataSourceMap["instance_id"] = dataSource.InstanceId
			}

			if dataSource.User != nil {
				dataSourceMap["user"] = dataSource.User
			}

			if dataSource.AliasName != nil {
				dataSourceMap["alias_name"] = dataSource.AliasName
			}

			if dataSource.Password != nil {
				dataSourceMap["password"] = dataSource.Password
			}

			dataSourcesList = append(dataSourcesList, dataSourceMap)
		}

		_ = d.Set("data_transform_sql_data_sources", dataSourcesList)
	}

	if dataTransform.EnvInfos != nil {
		var envInfosList []interface{}
		for _, envInfo := range dataTransform.EnvInfos {
			envInfoMap := map[string]interface{}{}

			if envInfo.Key != nil {
				envInfoMap["key"] = envInfo.Key
			}

			if envInfo.Value != nil {
				envInfoMap["value"] = envInfo.Value
			}

			envInfosList = append(envInfosList, envInfoMap)
		}

		_ = d.Set("env_infos", envInfosList)
	}

	return nil
}

func resourceTencentCloudClsDataTransformUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_data_transform.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		request             = cls.NewModifyDataTransformRequest()
		dataTransformTaskId = d.Id()
	)

	immutableArgs := []string{"src_topic_id", "preview_log_statistics", "data_transform_type", "process_from_timestamp", "process_to_timestamp"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.TaskId = &dataTransformTaskId

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("etl_content") {
		if v, ok := d.GetOk("etl_content"); ok {
			request.EtlContent = helper.String(v.(string))
		}
	}

	if d.HasChange("enable_flag") {
		if v, ok := d.GetOkExists("enable_flag"); ok {
			request.EnableFlag = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("dst_resources") {
		if v, ok := d.GetOk("dst_resources"); ok {
			for _, item := range v.([]interface{}) {
				dataTransformResouceInfo := cls.DataTransformResouceInfo{}
				dMap := item.(map[string]interface{})
				if v, ok := dMap["topic_id"]; ok {
					dataTransformResouceInfo.TopicId = helper.String(v.(string))
				}

				if v, ok := dMap["alias"]; ok {
					dataTransformResouceInfo.Alias = helper.String(v.(string))
				}

				request.DstResources = append(request.DstResources, &dataTransformResouceInfo)
			}
		}
	}

	if d.HasChange("backup_give_up_data") {
		if v, ok := d.GetOkExists("backup_give_up_data"); ok {
			request.BackupGiveUpData = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("has_services_log") {
		if v, ok := d.GetOkExists("has_services_log"); ok {
			request.HasServicesLog = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("keep_failure_log") {
		if v, ok := d.GetOkExists("keep_failure_log"); ok {
			request.KeepFailureLog = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("failure_log_key") {
		if v, ok := d.GetOk("failure_log_key"); ok {
			request.FailureLogKey = helper.String(v.(string))
		}
	}

	if d.HasChange("data_transform_sql_data_sources") {
		if v, ok := d.GetOk("data_transform_sql_data_sources"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				dataTransformSqlDataSource := cls.DataTransformSqlDataSource{}
				if v, ok := dMap["data_source"]; ok {
					dataTransformSqlDataSource.DataSource = helper.IntUint64(v.(int))
				}

				if v, ok := dMap["region"]; ok {
					dataTransformSqlDataSource.Region = helper.String(v.(string))
				}

				if v, ok := dMap["instance_id"]; ok {
					dataTransformSqlDataSource.InstanceId = helper.String(v.(string))
				}

				if v, ok := dMap["user"]; ok {
					dataTransformSqlDataSource.User = helper.String(v.(string))
				}

				if v, ok := dMap["alias_name"]; ok {
					dataTransformSqlDataSource.AliasName = helper.String(v.(string))
				}

				if v, ok := dMap["password"]; ok {
					dataTransformSqlDataSource.Password = helper.String(v.(string))
				}

				request.DataTransformSqlDataSources = append(request.DataTransformSqlDataSources, &dataTransformSqlDataSource)
			}
		}
	}

	if d.HasChange("env_infos") {
		if v, ok := d.GetOk("env_infos"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				envInfo := cls.EnvInfo{}
				if v, ok := dMap["key"]; ok {
					envInfo.Key = helper.String(v.(string))
				}

				if v, ok := dMap["value"]; ok {
					envInfo.Value = helper.String(v.(string))
				}

				request.EnvInfos = append(request.EnvInfos, &envInfo)
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyDataTransform(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cls dataTransform failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudClsDataTransformRead(d, meta)
}

func resourceTencentCloudClsDataTransformDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_data_transform.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		ctx                 = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service             = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dataTransformTaskId = d.Id()
	)

	if err := service.DeleteClsDataTransformById(ctx, dataTransformTaskId); err != nil {
		return err
	}

	return nil
}
