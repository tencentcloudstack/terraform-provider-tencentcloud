package postgresql

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlInstanceSslConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlInstanceSslConfigCreate,
		Read:   resourceTencentCloudPostgresqlInstanceSslConfigRead,
		Update: resourceTencentCloudPostgresqlInstanceSslConfigUpdate,
		Delete: resourceTencentCloudPostgresqlInstanceSslConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Postgres instance ID.",
			},

			"ssl_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable or disable SSL. true: enable; false: disable.",
			},

			"connect_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique connection address protected by SSL certificate, which can be set as the internal and external IP address if it is the primary instance; If it is a read-only instance, it can be set as the instance IP or read-only group IP. This parameter is mandatory when enabling SSL or modifying SSL protected connection addresses; When SSL is turned off, this parameter will be ignored.",
			},

			"ca_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud root certificate download link.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlInstanceSslConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_instance_ssl_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var dbInsntaceId string
	if v, ok := d.GetOk("db_instance_id"); ok {
		dbInsntaceId = v.(string)
	}

	d.SetId(dbInsntaceId)

	return resourceTencentCloudPostgresqlInstanceSslConfigUpdate(d, meta)
}

func resourceTencentCloudPostgresqlInstanceSslConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_instance_ssl_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service      = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dbInsntaceId = d.Id()
	)

	respData, err := service.DescribePostgresqlInstanceSslConfigById(ctx, dbInsntaceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `postgresql_instance_ssl_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("db_instance_id", dbInsntaceId)

	if respData.SSLEnabled != nil {
		_ = d.Set("ssl_enabled", respData.SSLEnabled)
	}

	if respData.ConnectAddress != nil {
		_ = d.Set("connect_address", respData.ConnectAddress)
	}

	if respData.CAUrl != nil {
		_ = d.Set("ca_url", respData.CAUrl)
	}

	return nil
}

func resourceTencentCloudPostgresqlInstanceSslConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_instance_ssl_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		dbInsntaceId = d.Id()
	)

	request := postgresv20170312.NewModifyDBInstanceSSLConfigRequest()
	response := postgresv20170312.NewModifyDBInstanceSSLConfigResponse()
	request.DBInstanceId = helper.String(dbInsntaceId)

	if v, ok := d.GetOkExists("ssl_enabled"); ok {
		request.SSLEnabled = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("connect_address"); ok {
		request.ConnectAddress = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresV20170312Client().ModifyDBInstanceSSLConfigWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Update postgresql instance ssl config failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update postgresql instance ssl config failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.TaskId == nil {
		return fmt.Errorf("TaksId is nil.")
	}

	// wait
	taskId := *response.Response.TaskId
	taskRequest := postgresv20170312.NewDescribeTasksRequest()
	taskRequest.TaskId = helper.Int64Uint64(taskId)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlV20170312Client().DescribeTasksWithContext(ctx, taskRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, taskRequest.GetAction(), taskRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskSet == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tasks failed, Response is nil."))
		}

		if len(result.Response.TaskSet) == 0 {
			return resource.RetryableError(fmt.Errorf("wait TaskSet init."))
		}

		if result.Response.TaskSet[0].Status != nil && *result.Response.TaskSet[0].Status == "Success" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("postgresql instance ssl config is running, status is %s.", *result.Response.TaskSet[0].Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s update postgresql instance ssl config, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPostgresqlInstanceSslConfigRead(d, meta)
}

func resourceTencentCloudPostgresqlInstanceSslConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_instance_ssl_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
