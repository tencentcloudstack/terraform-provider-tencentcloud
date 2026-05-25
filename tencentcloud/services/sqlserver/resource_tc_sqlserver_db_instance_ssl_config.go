package sqlserver

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverDbInstanceSslConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverDbInstanceSslConfigCreate,
		Read:   resourceTencentCloudSqlserverDbInstanceSslConfigRead,
		Update: resourceTencentCloudSqlserverDbInstanceSslConfigUpdate,
		Delete: resourceTencentCloudSqlserverDbInstanceSslConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SQL Server instance ID.",
			},

			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(SSL_TYPE, false),
				Description:  "SSL operation type. Valid values: enable, disable, renew.",
			},

			"wait_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Execution timing. 0: execute immediately, 1: execute during maintenance window. Default is 0.",
			},

			"is_kms": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Whether to enable KMS encryption protection. 0: no, 1: yes. Default is 0.",
			},

			"key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "KMS CMK key ID, required when IsKMS is 1.",
			},

			"key_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CMK region, required when IsKMS is 1.",
			},

			"encryption": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SSL encryption status. Valid values: enable, disable, enable_doing, disable_doing, renew_doing, wait_doing.",
			},

			"ssl_validity_period": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SSL certificate validity period, format: YYYY-MM-DD HH:MM:SS.",
			},

			"ssl_validity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "SSL certificate validity. 0: invalid, 1: valid.",
			},
		},
	}
}

func resourceTencentCloudSqlserverDbInstanceSslConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_db_instance_ssl_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverDbInstanceSslConfigUpdate(d, meta)
}

func resourceTencentCloudSqlserverDbInstanceSslConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_db_instance_ssl_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		sqlserverService = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId       = d.Id()
	)

	resp, err := sqlserverService.DescribeSqlserverInstanceSslById(ctx, instanceId)
	if err != nil {
		return err
	}

	if resp == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `sqlserver_db_instance_ssl_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if resp.SSLConfig != nil {
		sslConfig := resp.SSLConfig
		if sslConfig.Encryption != nil {
			_ = d.Set("encryption", sslConfig.Encryption)
		}
		if sslConfig.SSLValidityPeriod != nil {
			_ = d.Set("ssl_validity_period", sslConfig.SSLValidityPeriod)
		}
		if sslConfig.SSLValidity != nil {
			_ = d.Set("ssl_validity", int(*sslConfig.SSLValidity))
		}
	}

	return nil
}

func resourceTencentCloudSqlserverDbInstanceSslConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_db_instance_ssl_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		instanceId = d.Id()
	)

	request := sqlserver.NewModifyDBInstanceSSLRequest()
	request.InstanceId = helper.String(instanceId)

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("wait_switch"); ok {
		request.WaitSwitch = helper.Int64Uint64(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("is_kms"); ok {
		request.IsKMS = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("key_id"); ok {
		request.KeyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("key_region"); ok {
		request.KeyRegion = helper.String(v.(string))
	}

	var flowId uint64
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().ModifyDBInstanceSSLWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Update sqlserver db instance ssl config failed, Response is nil."))
		}

		if result.Response.FlowId == nil {
			return resource.NonRetryableError(fmt.Errorf("Update sqlserver db instance ssl config failed, FlowId is nil."))
		}

		flowId = *result.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver db instance ssl config failed, reason:%+v", logId, err)
		return err
	}

	// wait for async flow to complete
	flowRequest := sqlserver.NewDescribeFlowStatusRequest()
	flowRequest.FlowId = helper.UInt64Int64(flowId)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().DescribeFlowStatus(flowRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe FlowStatus failed, Response is nil."))
		}

		if *result.Response.Status == SQLSERVER_TASK_SUCCESS {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Operate sqlserver db instance ssl config status is %d.", *result.Response.Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate sqlserver db instance ssl config failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverDbInstanceSslConfigRead(d, meta)
}

func resourceTencentCloudSqlserverDbInstanceSslConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_db_instance_ssl_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
