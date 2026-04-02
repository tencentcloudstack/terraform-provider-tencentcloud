package cls

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsDlcDeliver() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsDlcDeliverCreate,
		Read:   resourceTencentCloudClsDlcDeliverRead,
		Update: resourceTencentCloudClsDlcDeliverUpdate,
		Delete: resourceTencentCloudClsDlcDeliverDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Log topic ID.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Task name. Length does not exceed 64 characters, starts with a letter, accepts 0-9, a-z, A-Z, _, -, Chinese characters.",
			},

			"deliver_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Delivery type. `0`: batch delivery, `1`: real-time delivery.",
			},

			"start_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Start time of the delivery time range (Unix timestamp).",
			},

			"dlc_info": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "DLC configuration information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"table_info": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "DLC table information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data_directory": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Data directory.",
									},
									"database_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Database name.",
									},
									"table_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Table name.",
									},
								},
							},
						},

						"field_infos": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "DLC data field information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cls_field": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Field name in CLS log.",
									},
									"dlc_field": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Column name in DLC table.",
									},
									"dlc_field_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "DLC field type, e.g. `string`, `int`, `struct`.",
									},
									"fill_field": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Fill field when parsing fails.",
									},
									"disable": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to disable this field.",
									},
								},
							},
						},

						"partition_infos": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "DLC partition information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cls_field": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Field name in CLS log.",
									},
									"dlc_field": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Column name in DLC table.",
									},
									"dlc_field_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "DLC field type.",
									},
								},
							},
						},

						"partition_extra": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "DLC partition extra information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time_format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Time format, e.g. `/%Y/%m/%d/%H`.",
									},
									"time_zone": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Time zone, e.g. `UTC+08:00`.",
									},
								},
							},
						},
					},
				},
			},

			"max_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Delivery file size in MB. Required when `deliver_type=0`. Range: 5 <= MaxSize <= 256.",
			},

			"interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Delivery interval in seconds. Required when `deliver_type=0`. Range: 300 <= Interval <= 900.",
			},

			"end_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "End time of the delivery time range (Unix timestamp). If empty, no time limit. Must be greater than `start_time` when set.",
			},

			"has_services_log": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to enable delivery service logs. `1`: disabled, `2`: enabled. Default is enabled.",
			},

			"status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Task status. `1`: running, `2`: stopped.",
			},

			// computed
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Delivery task ID.",
			},
		},
	}
}

func resourceTencentCloudClsDlcDeliverCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_dlc_deliver.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = cls.NewCreateDlcDeliverRequest()
		response = cls.NewCreateDlcDeliverResponse()
		topicId  string
	)

	if v, ok := d.GetOk("topic_id"); ok {
		request.TopicId = helper.String(v.(string))
		topicId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("deliver_type"); ok {
		request.DeliverType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("start_time"); ok {
		request.StartTime = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("dlc_info"); ok {
		dlcInfoList := v.([]interface{})
		if len(dlcInfoList) > 0 {
			dlcInfoMap := dlcInfoList[0].(map[string]interface{})
			dlcInfo := &cls.DlcInfo{}

			if tableInfoList, ok := dlcInfoMap["table_info"].([]interface{}); ok && len(tableInfoList) > 0 {
				tableInfoMap := tableInfoList[0].(map[string]interface{})
				tableInfo := &cls.DlcTableInfo{}
				if v, ok := tableInfoMap["data_directory"].(string); ok && v != "" {
					tableInfo.DataDirectory = helper.String(v)
				}
				if v, ok := tableInfoMap["database_name"].(string); ok && v != "" {
					tableInfo.DatabaseName = helper.String(v)
				}
				if v, ok := tableInfoMap["table_name"].(string); ok && v != "" {
					tableInfo.TableName = helper.String(v)
				}
				dlcInfo.TableInfo = tableInfo
			}

			if fieldInfosList, ok := dlcInfoMap["field_infos"].([]interface{}); ok {
				for _, item := range fieldInfosList {
					fieldMap := item.(map[string]interface{})
					fieldInfo := &cls.DlcFiledInfo{}
					if v, ok := fieldMap["cls_field"].(string); ok && v != "" {
						fieldInfo.ClsField = helper.String(v)
					}
					if v, ok := fieldMap["dlc_field"].(string); ok && v != "" {
						fieldInfo.DlcField = helper.String(v)
					}
					if v, ok := fieldMap["dlc_field_type"].(string); ok && v != "" {
						fieldInfo.DlcFieldType = helper.String(v)
					}
					if v, ok := fieldMap["fill_field"].(string); ok && v != "" {
						fieldInfo.FillField = helper.String(v)
					}
					if v, ok := fieldMap["disable"].(bool); ok {
						fieldInfo.Disable = helper.Bool(v)
					}
					dlcInfo.FieldInfos = append(dlcInfo.FieldInfos, fieldInfo)
				}
			}

			if partitionInfosList, ok := dlcInfoMap["partition_infos"].([]interface{}); ok {
				for _, item := range partitionInfosList {
					partMap := item.(map[string]interface{})
					partInfo := &cls.DlcPartitionInfo{}
					if v, ok := partMap["cls_field"].(string); ok && v != "" {
						partInfo.ClsField = helper.String(v)
					}
					if v, ok := partMap["dlc_field"].(string); ok && v != "" {
						partInfo.DlcField = helper.String(v)
					}
					if v, ok := partMap["dlc_field_type"].(string); ok && v != "" {
						partInfo.DlcFieldType = helper.String(v)
					}
					dlcInfo.PartitionInfos = append(dlcInfo.PartitionInfos, partInfo)
				}
			}

			if partitionExtraList, ok := dlcInfoMap["partition_extra"].([]interface{}); ok && len(partitionExtraList) > 0 {
				extraMap := partitionExtraList[0].(map[string]interface{})
				partitionExtra := &cls.DlcPartitionExtra{}
				if v, ok := extraMap["time_format"].(string); ok && v != "" {
					partitionExtra.TimeFormat = helper.String(v)
				}
				if v, ok := extraMap["time_zone"].(string); ok && v != "" {
					partitionExtra.TimeZone = helper.String(v)
				}
				dlcInfo.PartitionExtra = partitionExtra
			}

			request.DlcInfo = dlcInfo
		}
	}

	if v, ok := d.GetOkExists("max_size"); ok {
		request.MaxSize = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("interval"); ok {
		request.Interval = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("end_time"); ok {
		request.EndTime = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("has_services_log"); ok {
		request.HasServicesLog = helper.IntUint64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateDlcDeliverWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create cls dlc deliver failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create cls dlc deliver failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.TaskId == nil {
		return fmt.Errorf("TaskId is nil.")
	}

	d.SetId(strings.Join([]string{topicId, *response.Response.TaskId}, tccommon.FILED_SP))
	return resourceTencentCloudClsDlcDeliverRead(d, meta)
}

func resourceTencentCloudClsDlcDeliverRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_dlc_deliver.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	topicId := idSplit[0]
	taskId := idSplit[1]

	respData, err := service.DescribeClsDlcDeliverById(ctx, topicId, taskId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cls_dlc_deliver` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("topic_id", topicId)

	if respData.TaskId != nil {
		_ = d.Set("task_id", respData.TaskId)
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.DeliverType != nil {
		_ = d.Set("deliver_type", int(*respData.DeliverType))
	}

	if respData.StartTime != nil {
		_ = d.Set("start_time", int(*respData.StartTime))
	}

	if respData.EndTime != nil {
		_ = d.Set("end_time", int(*respData.EndTime))
	}

	if respData.MaxSize != nil {
		_ = d.Set("max_size", int(*respData.MaxSize))
	}

	if respData.Interval != nil {
		_ = d.Set("interval", int(*respData.Interval))
	}

	if respData.HasServicesLog != nil {
		_ = d.Set("has_services_log", int(*respData.HasServicesLog))
	}

	if respData.Status != nil {
		_ = d.Set("status", int(*respData.Status))
	}

	if respData.DlcInfo != nil {
		dlcInfoMap := map[string]interface{}{}

		if respData.DlcInfo.TableInfo != nil {
			tableInfoMap := map[string]interface{}{}
			if respData.DlcInfo.TableInfo.DataDirectory != nil {
				tableInfoMap["data_directory"] = respData.DlcInfo.TableInfo.DataDirectory
			}
			if respData.DlcInfo.TableInfo.DatabaseName != nil {
				tableInfoMap["database_name"] = respData.DlcInfo.TableInfo.DatabaseName
			}
			if respData.DlcInfo.TableInfo.TableName != nil {
				tableInfoMap["table_name"] = respData.DlcInfo.TableInfo.TableName
			}
			dlcInfoMap["table_info"] = []interface{}{tableInfoMap}
		}

		if len(respData.DlcInfo.FieldInfos) > 0 {
			fieldInfosList := make([]interface{}, 0, len(respData.DlcInfo.FieldInfos))
			for _, fi := range respData.DlcInfo.FieldInfos {
				fiMap := map[string]interface{}{}
				if fi.ClsField != nil {
					fiMap["cls_field"] = fi.ClsField
				}
				if fi.DlcField != nil {
					fiMap["dlc_field"] = fi.DlcField
				}
				if fi.DlcFieldType != nil {
					fiMap["dlc_field_type"] = fi.DlcFieldType
				}
				if fi.FillField != nil {
					fiMap["fill_field"] = fi.FillField
				}
				if fi.Disable != nil {
					fiMap["disable"] = fi.Disable
				}
				fieldInfosList = append(fieldInfosList, fiMap)
			}
			dlcInfoMap["field_infos"] = fieldInfosList
		}

		if len(respData.DlcInfo.PartitionInfos) > 0 {
			partInfosList := make([]interface{}, 0, len(respData.DlcInfo.PartitionInfos))
			for _, pi := range respData.DlcInfo.PartitionInfos {
				piMap := map[string]interface{}{}
				if pi.ClsField != nil {
					piMap["cls_field"] = pi.ClsField
				}
				if pi.DlcField != nil {
					piMap["dlc_field"] = pi.DlcField
				}
				if pi.DlcFieldType != nil {
					piMap["dlc_field_type"] = pi.DlcFieldType
				}
				partInfosList = append(partInfosList, piMap)
			}
			dlcInfoMap["partition_infos"] = partInfosList
		}

		if respData.DlcInfo.PartitionExtra != nil {
			extraMap := map[string]interface{}{}
			if respData.DlcInfo.PartitionExtra.TimeFormat != nil {
				extraMap["time_format"] = respData.DlcInfo.PartitionExtra.TimeFormat
			}
			if respData.DlcInfo.PartitionExtra.TimeZone != nil {
				extraMap["time_zone"] = respData.DlcInfo.PartitionExtra.TimeZone
			}
			dlcInfoMap["partition_extra"] = []interface{}{extraMap}
		}

		_ = d.Set("dlc_info", []interface{}{dlcInfoMap})
	}

	return nil
}

func resourceTencentCloudClsDlcDeliverUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_dlc_deliver.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	needChange := false
	mutableArgs := []string{
		"name", "deliver_type", "start_time", "end_time",
		"max_size", "interval", "dlc_info", "has_services_log", "status",
	}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := cls.NewModifyDlcDeliverRequest()
		request.TopicId = &idSplit[0]
		request.TaskId = &idSplit[1]

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("deliver_type"); ok {
			request.DeliverType = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("start_time"); ok {
			request.StartTime = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("end_time"); ok {
			request.EndTime = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("max_size"); ok {
			request.MaxSize = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("interval"); ok {
			request.Interval = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("has_services_log"); ok {
			request.HasServicesLog = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("status"); ok {
			request.Status = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("dlc_info"); ok {
			dlcInfoList := v.([]interface{})
			if len(dlcInfoList) > 0 {
				dlcInfoMap := dlcInfoList[0].(map[string]interface{})
				dlcInfo := &cls.DlcInfo{}

				if tableInfoList, ok := dlcInfoMap["table_info"].([]interface{}); ok && len(tableInfoList) > 0 {
					tableInfoMap := tableInfoList[0].(map[string]interface{})
					tableInfo := &cls.DlcTableInfo{}
					if v, ok := tableInfoMap["data_directory"].(string); ok && v != "" {
						tableInfo.DataDirectory = helper.String(v)
					}
					if v, ok := tableInfoMap["database_name"].(string); ok && v != "" {
						tableInfo.DatabaseName = helper.String(v)
					}
					if v, ok := tableInfoMap["table_name"].(string); ok && v != "" {
						tableInfo.TableName = helper.String(v)
					}
					dlcInfo.TableInfo = tableInfo
				}

				if fieldInfosList, ok := dlcInfoMap["field_infos"].([]interface{}); ok {
					for _, item := range fieldInfosList {
						fieldMap := item.(map[string]interface{})
						fieldInfo := &cls.DlcFiledInfo{}
						if v, ok := fieldMap["cls_field"].(string); ok && v != "" {
							fieldInfo.ClsField = helper.String(v)
						}
						if v, ok := fieldMap["dlc_field"].(string); ok && v != "" {
							fieldInfo.DlcField = helper.String(v)
						}
						if v, ok := fieldMap["dlc_field_type"].(string); ok && v != "" {
							fieldInfo.DlcFieldType = helper.String(v)
						}
						if v, ok := fieldMap["fill_field"].(string); ok && v != "" {
							fieldInfo.FillField = helper.String(v)
						}
						if v, ok := fieldMap["disable"].(bool); ok {
							fieldInfo.Disable = helper.Bool(v)
						}
						dlcInfo.FieldInfos = append(dlcInfo.FieldInfos, fieldInfo)
					}
				}

				if partitionInfosList, ok := dlcInfoMap["partition_infos"].([]interface{}); ok {
					for _, item := range partitionInfosList {
						partMap := item.(map[string]interface{})
						partInfo := &cls.DlcPartitionInfo{}
						if v, ok := partMap["cls_field"].(string); ok && v != "" {
							partInfo.ClsField = helper.String(v)
						}
						if v, ok := partMap["dlc_field"].(string); ok && v != "" {
							partInfo.DlcField = helper.String(v)
						}
						if v, ok := partMap["dlc_field_type"].(string); ok && v != "" {
							partInfo.DlcFieldType = helper.String(v)
						}
						dlcInfo.PartitionInfos = append(dlcInfo.PartitionInfos, partInfo)
					}
				}

				if partitionExtraList, ok := dlcInfoMap["partition_extra"].([]interface{}); ok && len(partitionExtraList) > 0 {
					extraMap := partitionExtraList[0].(map[string]interface{})
					partitionExtra := &cls.DlcPartitionExtra{}
					if v, ok := extraMap["time_format"].(string); ok && v != "" {
						partitionExtra.TimeFormat = helper.String(v)
					}
					if v, ok := extraMap["time_zone"].(string); ok && v != "" {
						partitionExtra.TimeZone = helper.String(v)
					}
					dlcInfo.PartitionExtra = partitionExtra
				}

				request.DlcInfo = dlcInfo
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyDlcDeliverWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update cls dlc deliver failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudClsDlcDeliverRead(d, meta)
}

func resourceTencentCloudClsDlcDeliverDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_dlc_deliver.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = cls.NewDeleteDlcDeliverRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	request.TopicId = &idSplit[0]
	request.TaskId = &idSplit[1]

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().DeleteDlcDeliverWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete cls dlc deliver failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
