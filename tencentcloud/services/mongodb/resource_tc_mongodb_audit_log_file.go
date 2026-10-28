package mongodb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudMongodbAuditLogFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbAuditLogFileCreate,
		Read:   resourceTencentCloudMongodbAuditLogFileRead,
		Update: resourceTencentCloudMongodbAuditLogFileUpdate,
		Delete: resourceTencentCloudMongodbAuditLogFileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID, the format is: cmgo-xfts****.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Start time, format: \"2021-07-12 10:29:20\".",
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "End time, format: \"2021-07-12 10:39:20\".",
			},
			"order": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort order. Valid values: `ASC`, `DESC`.",
			},
			"order_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort field. Valid values: `timestamp`, `affectRows`, `execTime`.",
			},
			"filter": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Filter conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Client addresses.",
						},
						"user": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Usernames.",
						},
						"exec_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Minimum execution time in ms.",
						},
						"affect_rows": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Minimum affected rows.",
						},
						"atype": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Operation types.",
						},
						"result": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Execution results.",
						},
						"param": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Keywords to filter logs.",
						},
					},
				},
			},
			"file_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The generated audit log file name.",
			},
			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Audit log file details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File name.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File status. Valid values: `creating`, `failed`, `success`.",
						},
						"file_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "File size in KB.",
						},
						"download_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Download URL.",
						},
						"err_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Error message.",
						},
						"progress_rate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Download progress.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMongodbAuditLogFileCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_audit_log_file.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = mongodb.NewCreateAuditLogFileRequest()
		response   = mongodb.NewCreateAuditLogFileResponse()
		instanceId string
		fileName   string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		request.Order = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		request.OrderBy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter"); ok {
		filterList := v.([]interface{})
		if len(filterList) > 0 {
			filterMap := filterList[0].(map[string]interface{})
			auditLogFilter := &mongodb.AuditLogFilter{}
			if v, ok := filterMap["host"]; ok {
				hostList := v.([]interface{})
				if len(hostList) > 0 {
					auditLogFilter.Host = helper.InterfacesStringsPoint(hostList)
				}
			}

			if v, ok := filterMap["user"]; ok {
				userList := v.([]interface{})
				if len(userList) > 0 {
					auditLogFilter.User = helper.InterfacesStringsPoint(userList)
				}
			}

			if v, ok := filterMap["exec_time"]; ok && v.(int) > 0 {
				auditLogFilter.ExecTime = helper.IntUint64(v.(int))
			}

			if v, ok := filterMap["affect_rows"]; ok && v.(int) > 0 {
				auditLogFilter.AffectRows = helper.IntUint64(v.(int))
			}

			if v, ok := filterMap["atype"]; ok {
				atypeList := v.([]interface{})
				if len(atypeList) > 0 {
					auditLogFilter.Atype = helper.InterfacesStringsPoint(atypeList)
				}
			}

			if v, ok := filterMap["result"]; ok {
				resultList := v.([]interface{})
				if len(resultList) > 0 {
					auditLogFilter.Result = helper.InterfacesStringsPoint(resultList)
				}
			}

			if v, ok := filterMap["param"]; ok {
				paramList := v.([]interface{})
				if len(paramList) > 0 {
					auditLogFilter.Param = helper.InterfacesStringsPoint(paramList)
				}
			}

			request.Filter = auditLogFilter
		}
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().CreateAuditLogFile(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create mongodb audit log file failed, Response is nil"))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mongodb audit log file failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.FileName == nil || *response.Response.FileName == "" {
		return fmt.Errorf("Create mongodb audit log file failed, FileName is nil or empty")
	}

	fileName = *response.Response.FileName
	d.SetId(strings.Join([]string{instanceId, fileName}, tccommon.FILED_SP))

	// Wait for the audit log file to be generated.
	describeRequest := mongodb.NewDescribeAuditLogFilesRequest()
	describeRequest.InstanceId = helper.String(instanceId)
	describeRequest.FileName = helper.String(fileName)
	describeRequest.Limit = helper.Uint64(uint64(100))
	describeRequest.Offset = helper.Uint64(uint64(0))

	waitErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(describeRequest.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().DescribeAuditLogFiles(describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, describeRequest.GetAction(), describeRequest.ToJsonString(), result.ToJsonString())
		if result == nil || result.Response == nil || len(result.Response.Items) == 0 {
			return resource.RetryableError(fmt.Errorf("mongodb audit log file [%s] is still creating, items is empty", fileName))
		}

		item := result.Response.Items[0]
		if item == nil || item.Status == nil {
			return resource.RetryableError(fmt.Errorf("mongodb audit log file [%s] is still creating, status is nil", fileName))
		}

		status := *item.Status
		switch status {
		case "success":
			return nil
		case "failed":
			errMsg := ""
			if item.ErrMsg != nil {
				errMsg = *item.ErrMsg
			}

			return resource.NonRetryableError(fmt.Errorf("mongodb audit log file [%s] create failed, err_msg: [%s]", fileName, errMsg))
		default:
			return resource.RetryableError(fmt.Errorf("mongodb audit log file [%s] is still creating, current status: [%s]", fileName, status))
		}
	})

	if waitErr != nil {
		log.Printf("[CRITAL]%s wait mongodb audit log file create failed, reason:%+v", logId, waitErr)
		return waitErr
	}

	return resourceTencentCloudMongodbAuditLogFileRead(d, meta)
}

func resourceTencentCloudMongodbAuditLogFileRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_audit_log_file.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		_        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = mongodb.NewDescribeAuditLogFilesRequest()
		response = mongodb.NewDescribeAuditLogFilesResponse()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	instanceId := idSplit[0]
	fileName := idSplit[1]

	request.InstanceId = helper.String(instanceId)
	request.FileName = helper.String(fileName)
	request.Limit = helper.Uint64(uint64(100))
	request.Offset = helper.Uint64(uint64(0))

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().DescribeAuditLogFiles(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe mongodb audit log files failed, Response is nil"))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read mongodb audit log file failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.Items == nil || len(response.Response.Items) == 0 {
		log.Printf("[WARN]%s resource `tencentcloud_mongodb_audit_log_file` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("file_name", fileName)

	itemsList := make([]map[string]interface{}, 0, len(response.Response.Items))
	for _, item := range response.Response.Items {
		itemMap := map[string]interface{}{}
		if item.FileName != nil {
			itemMap["file_name"] = item.FileName
		}

		if item.CreateTime != nil {
			itemMap["create_time"] = item.CreateTime
		}

		if item.Status != nil {
			itemMap["status"] = item.Status
		}

		if item.FileSize != nil {
			itemMap["file_size"] = item.FileSize
		}

		if item.DownloadUrl != nil {
			itemMap["download_url"] = item.DownloadUrl
		}

		if item.ErrMsg != nil {
			itemMap["err_msg"] = item.ErrMsg
		}

		if item.ProgressRate != nil {
			itemMap["progress_rate"] = item.ProgressRate
		}

		itemsList = append(itemsList, itemMap)
	}

	_ = d.Set("items", itemsList)

	return nil
}

func resourceTencentCloudMongodbAuditLogFileUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_audit_log_file.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"start_time", "end_time", "order", "order_by", "filter"}
	if err := helper.ImmutableArgsChek(d, immutableArgs...); err != nil {
		return err
	}

	return resourceTencentCloudMongodbAuditLogFileRead(d, meta)
}

func resourceTencentCloudMongodbAuditLogFileDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_audit_log_file.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = mongodb.NewDeleteAuditLogFileRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	instanceId := idSplit[0]
	fileName := idSplit[1]

	request.InstanceId = helper.String(instanceId)
	request.FileName = helper.String(fileName)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().DeleteAuditLogFile(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete mongodb audit log file failed, reason:%+v", logId, err)
		return err
	}

	// Wait for the audit log file to be generated.
	describeRequest := mongodb.NewDescribeAuditLogFilesRequest()
	describeRequest.InstanceId = helper.String(instanceId)
	describeRequest.FileName = helper.String(fileName)
	describeRequest.Limit = helper.Uint64(uint64(100))
	describeRequest.Offset = helper.Uint64(uint64(0))

	waitErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(describeRequest.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().DescribeAuditLogFiles(describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, describeRequest.GetAction(), describeRequest.ToJsonString(), result.ToJsonString())
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe mongodb audit log file filed, Response is nil"))
		}

		if len(result.Response.Items) == 0 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("mongodb audit log file [%s] is still deleting", fileName))
	})

	if waitErr != nil {
		log.Printf("[CRITAL]%s wait mongodb audit log file delete failed, reason:%+v", logId, waitErr)
		return waitErr
	}

	return nil
}
