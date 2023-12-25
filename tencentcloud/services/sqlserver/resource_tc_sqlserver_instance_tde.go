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

func ResourceTencentCloudSqlserverInstanceTDE() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverInstanceTDECreate,
		Read:   resourceTencentCloudSqlserverInstanceTDERead,
		Update: resourceTencentCloudSqlserverInstanceTDEUpdate,
		Delete: resourceTencentCloudSqlserverInstanceTDEDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"certificate_attribution": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Certificate attribution. self- means to use the account's own certificate, others- means to refer to the certificate of other accounts, and the default is self.",
			},
			"quote_uin": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Other referenced main account IDs, required when CertificateAttribute is others.",
			},
		},
	}
}

func resourceTencentCloudSqlserverInstanceTDECreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_instance_tde.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverInstanceTDEUpdate(d, meta)
}

func resourceTencentCloudSqlserverInstanceTDERead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_instance_tde.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	instanceTDE, err := service.DescribeSqlserverInstanceTDEById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instanceTDE == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverInstanceTDE` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instanceTDE.InstanceId != nil {
		_ = d.Set("instance_id", instanceTDE.InstanceId)
	}

	if instanceTDE.TDEConfig.CertificateAttribution != nil {
		_ = d.Set("certificate_attribution", instanceTDE.TDEConfig.CertificateAttribution)
	}

	if instanceTDE.TDEConfig.QuoteUin != nil {
		_ = d.Set("quote_uin", instanceTDE.TDEConfig.QuoteUin)
	}

	return nil
}

func resourceTencentCloudSqlserverInstanceTDEUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_instance_tde.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request    = sqlserver.NewModifyInstanceEncryptAttributesRequest()
		instanceId = d.Id()
		flowId     int64
	)

	if v, ok := d.GetOk("certificate_attribution"); ok {
		request.CertificateAttribution = helper.String(v.(string))
	}

	if v, ok := d.GetOk("quote_uin"); ok {
		request.QuoteUin = helper.String(v.(string))
	}

	request.InstanceId = &instanceId

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().ModifyInstanceEncryptAttributes(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("sqlserver instanceTDE not exists")
			return resource.NonRetryableError(e)
		}

		flowId = *result.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver instanceTDE failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(10*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCloneStatusByFlowId(ctx, flowId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil {
			e = fmt.Errorf("update sqlserver instanceTDE instanceId %s flowId %d not exists", instanceId, flowId)
			return resource.NonRetryableError(e)
		}

		if *result.Status == SQLSERVER_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("update sqlserver instanceTDE task status is running"))
		}

		if *result.Status == SQLSERVER_TASK_SUCCESS {
			return nil
		}

		if *result.Status == SQLSERVER_TASK_FAIL {
			return resource.NonRetryableError(fmt.Errorf("update sqlserver instanceTDE task status is failed"))
		}

		e = fmt.Errorf("update sqlserver instanceTDE task status is %v, we won't wait for it finish", *result.Status)
		return resource.NonRetryableError(e)
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver instanceTDE task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudSqlserverInstanceTDERead(d, meta)
}

func resourceTencentCloudSqlserverInstanceTDEDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_instance_tde.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
