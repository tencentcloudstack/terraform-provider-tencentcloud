package cdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlClsLogAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlClsLogAttachmentCreate,
		Read:   resourceTencentCloudMysqlClsLogAttachmentRead,
		Delete: resourceTencentCloudMysqlClsLogAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The id of instance.",
			},
			"log_type": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue(MYSQL_LOG_TO_CLS_TYPE),
				Description:  "Log type. Support `error` or `slowlog`.",
			},
			"create_log_set": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to create log set.",
			},
			"create_log_topic": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to create log topic.",
			},
			"log_set": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "If `create_log_set` is `true`, use log set name, Else use log set Id.",
			},
			"log_topic": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "If `create_log_topic` is `true`, use log topic name, Else use log topic Id.",
			},
			"period": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The validity period of the log theme is 30 days by default when not filled in.",
			},
			"create_index": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to create index.",
			},
			"cls_region": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Cls region.",
			},
			"log_set_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Log set Id.",
			},
			"log_topic_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Log topic Id.",
			},
			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Log Status.",
			},
		},
	}
}

func resourceTencentCloudMysqlClsLogAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_cls_log_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = mysql.NewModifyDBInstanceLogToCLSRequest()
		instanceId string
		logType    string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_type"); ok {
		logType = v.(string)
		request.LogType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("create_log_set"); ok {
		request.CreateLogset = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("create_log_topic"); ok {
		request.CreateLogTopic = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("log_set"); ok {
		request.Logset = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_topic"); ok {
		request.LogTopic = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("create_index"); ok {
		request.CreateIndex = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("cls_region"); ok {
		request.ClsRegion = helper.String(v.(string))
	}

	request.Status = helper.String("ON")
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().ModifyDBInstanceLogToCLS(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mysql ModifyDBInstanceLogToCLS failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{instanceId, logType}, tccommon.FILED_SP))

	return resourceTencentCloudMysqlClsLogAttachmentRead(d, meta)
}

func resourceTencentCloudMysqlClsLogAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_cls_log_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	logType := idSplit[1]

	logToCLSResponseParam, err := service.DescribeMysqlInstanceLogToCLSById(ctx, instanceId)
	if err != nil {
		return err
	}

	if logToCLSResponseParam == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DescribeDBInstanceLogToCLS` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("log_type", logType)

	if logType == MYSQL_LOG_TO_CLS_TYPE_ERROR {
		if logToCLSResponseParam.ErrorLog.LogSetId != nil {
			_ = d.Set("log_set_id", *logToCLSResponseParam.ErrorLog.LogSetId)
		}

		if logToCLSResponseParam.ErrorLog.LogTopicId != nil {
			_ = d.Set("log_topic_id", *logToCLSResponseParam.ErrorLog.LogTopicId)
		}

		if logToCLSResponseParam.ErrorLog.Status != nil {
			_ = d.Set("status", *logToCLSResponseParam.ErrorLog.Status)
		}

		if logToCLSResponseParam.ErrorLog.ClsRegion != nil {
			_ = d.Set("cls_region", *logToCLSResponseParam.ErrorLog.ClsRegion)
		}
	} else {
		if logToCLSResponseParam.SlowLog.LogSetId != nil {
			_ = d.Set("log_set_id", *logToCLSResponseParam.SlowLog.LogSetId)
		}

		if logToCLSResponseParam.SlowLog.LogTopicId != nil {
			_ = d.Set("log_topic_id", *logToCLSResponseParam.SlowLog.LogTopicId)
		}

		if logToCLSResponseParam.SlowLog.Status != nil {
			_ = d.Set("status", *logToCLSResponseParam.SlowLog.Status)
		}

		if logToCLSResponseParam.SlowLog.ClsRegion != nil {
			_ = d.Set("cls_region", *logToCLSResponseParam.SlowLog.ClsRegion)
		}
	}

	return nil
}

func resourceTencentCloudMysqlClsLogAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_cls_log_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	logType := idSplit[1]

	if err := service.DeleteMysqlInstanceLogToCLSById(ctx, instanceId, logType); err != nil {
		return err
	}

	return nil
}
