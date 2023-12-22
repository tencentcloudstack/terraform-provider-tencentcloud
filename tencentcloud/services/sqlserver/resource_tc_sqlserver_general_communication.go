package sqlserver

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverGeneralCommunication() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverGeneralCommunicationCreate,
		Read:   resourceTencentCloudSqlserverGeneralCommunicationRead,
		Delete: resourceTencentCloudSqlserverGeneralCommunicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of instances.",
			},
		},
	}
}

func resourceTencentCloudSqlserverGeneralCommunicationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_communication.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		request     = sqlserver.NewOpenInterCommunicationRequest()
		response    = sqlserver.NewOpenInterCommunicationResponse()
		flowRequest = sqlserver.NewDescribeFlowStatusRequest()
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceIdSet = append(request.InstanceIdSet, helper.String(v.(string)))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().OpenInterCommunication(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("sqlserver generalCommunication not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver generalCommunication failed, reason:%+v", logId, err)
		return err
	}

	instanceId := *response.Response.InterInstanceFlowSet[0].InstanceId
	flowId := *response.Response.InterInstanceFlowSet[0].FlowId
	flowRequest.FlowId = &flowId
	err = resource.Retry(10*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().DescribeFlowStatus(flowRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if *result.Response.Status == SQLSERVER_TASK_SUCCESS {
			return nil
		} else if *result.Response.Status == SQLSERVER_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("sqlserver generalCommunication status is running"))
		} else if *result.Response.Status == int64(SQLSERVER_TASK_FAIL) {
			return resource.NonRetryableError(fmt.Errorf("sqlserver bgeneralCommunication status is fail"))
		} else {
			e = fmt.Errorf("sqlserver generalCommunication status illegal")
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver generalCommunication failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverGeneralCommunicationRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCommunicationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_communication.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	generalCommunication, err := service.DescribeSqlserverGeneralCommunicationById(ctx, instanceId)
	if err != nil {
		return err
	}

	if generalCommunication == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverGeneralCommunication` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if generalCommunication.InstanceId != nil {
		_ = d.Set("instance_id", generalCommunication.InstanceId)
	}

	return nil
}

func resourceTencentCloudSqlserverGeneralCommunicationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_communication.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service     = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		flowRequest = sqlserver.NewDescribeFlowStatusRequest()
		instanceId  = d.Id()
	)

	flowId, err := service.DeleteSqlserverGeneralCommunicationById(ctx, instanceId)
	if err != nil {
		log.Printf("[CRITAL]%s delete sqlserver generalCommunication failed, reason:%+v", logId, err)
		return err
	}

	flowRequest.FlowId = &flowId
	err = resource.Retry(10*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().DescribeFlowStatus(flowRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if *result.Response.Status == SQLSERVER_TASK_SUCCESS {
			return nil
		} else if *result.Response.Status == SQLSERVER_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("sqlserver generalCommunication status is running"))
		} else if *result.Response.Status == int64(SQLSERVER_TASK_FAIL) {
			return resource.NonRetryableError(fmt.Errorf("sqlserver bgeneralCommunication status is fail"))
		} else {
			e = fmt.Errorf("sqlserver generalCommunication status illegal")
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete sqlserver generalCommunication status failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
