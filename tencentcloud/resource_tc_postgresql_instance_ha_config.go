package tencentcloud

import (
	"context"
	"log"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudPostgresqlInstanceHAConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlInstanceHAConfigCreate,
		Read:   resourceTencentCloudPostgresqlInstanceHAConfigRead,
		Update: resourceTencentCloudPostgresqlInstanceHAConfigUpdate,
		Delete: resourceTencentCloudPostgresqlInstanceHAConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},
			"sync_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(SYNC_MODE),
				Description:  "Master slave synchronization method, Semi-sync: Semi synchronous; Async: Asynchronous. Main instance default value: Semi-sync, Read-only instance default value: Async.",
			},
			"max_standby_latency": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(1073741824, 322122547200),
				Description:  "Maximum latency data volume for highly available backup machines. When the delay data amount of the backup node is less than or equal to this value, and the delay time of the backup node is less than or equal to MaxStandbyLag, it can switch to the main node. Unit: byte; Parameter range: [1073741824, 322122547200].",
			},
			"max_standby_lag": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(5, 10),
				Description:  "Maximum latency of highly available backup machines. When the delay time of the backup node is less than or equal to this value, and the amount of delay data of the backup node is less than or equal to MaxStandbyLatency, the primary node can be switched. Unit: s; Parameter range: [5, 10].",
			},
			"max_sync_standby_latency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum latency data for synchronous backup. When the amount of data delayed by the backup machine is less than or equal to this value, and the delay time of the backup machine is less than or equal to MaxSyncStandbyLag, then the backup machine adopts synchronous replication; Otherwise, adopt asynchronous replication. This parameter value is valid for instances where SyncMode is set to Semi sync. When a semi synchronous instance prohibits degradation to asynchronous replication, MaxSyncStandbyLatency and MaxSyncStandbyLag are not set. When semi synchronous instances allow degenerate asynchronous replication, PostgreSQL version 9 instances must have MaxSyncStandbyLatency set and MaxSyncStandbyLag not set, while PostgreSQL version 10 and above instances must have MaxSyncStandbyLatency and MaxSyncStandbyLag set.",
			},
			"max_sync_standby_lag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum delay time for synchronous backup. When the delay time of the standby machine is less than or equal to this value, and the amount of delay data of the standby machine is less than or equal to MaxSyncStandbyLatency, then the standby machine adopts synchronous replication; Otherwise, adopt asynchronous replication. This parameter value is valid for instances where SyncMode is set to Semi sync. When a semi synchronous instance prohibits degradation to asynchronous replication, MaxSyncStandbyLatency and MaxSyncStandbyLag are not set. When semi synchronous instances allow degenerate asynchronous replication, PostgreSQL version 9 instances must have MaxSyncStandbyLatency set and MaxSyncStandbyLag not set, while PostgreSQL version 10 and above instances must have MaxSyncStandbyLatency and MaxSyncStandbyLag set.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlInstanceHAConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_instance_ha_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudPostgresqlInstanceHAConfigUpdate(d, meta)
}

func resourceTencentCloudPostgresqlInstanceHAConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_instance_ha_config.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId = d.Id()
	)

	haConfig, err := service.DescribePostgresqlInstanceHAConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if haConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresqlInstanceHAConfig` [%s] not found.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if haConfig.SyncMode != nil {
		_ = d.Set("sync_mode", haConfig.SyncMode)
	}

	if haConfig.MaxStandbyLatency != nil {
		_ = d.Set("max_standby_latency", haConfig.MaxStandbyLatency)
	}

	if haConfig.MaxStandbyLag != nil {
		_ = d.Set("max_standby_lag", haConfig.MaxStandbyLag)
	}

	if haConfig.MaxSyncStandbyLatency != nil {
		_ = d.Set("max_sync_standby_latency", haConfig.MaxSyncStandbyLatency)
	}

	if haConfig.MaxSyncStandbyLag != nil {
		_ = d.Set("max_sync_standby_lag", haConfig.MaxSyncStandbyLag)
	}

	return nil
}

func resourceTencentCloudPostgresqlInstanceHAConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_instance_ha_config.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		request           = postgresql.NewModifyDBInstanceHAConfigRequest()
		instanceId        = d.Id()
		syncMode          = d.Get("sync_mode").(string)
		maxStandbyLatency = d.Get("max_standby_latency").(int)
		maxStandbyLag     = d.Get("max_standby_lag").(int)
	)

	request.DBInstanceId = &instanceId
	request.SyncMode = &syncMode
	request.MaxStandbyLatency = helper.IntUint64(maxStandbyLatency)
	request.MaxStandbyLag = helper.IntUint64(maxStandbyLag)
	if v, ok := d.GetOkExists("max_sync_standby_latency"); ok {
		request.MaxSyncStandbyLatency = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("max_sync_standby_lag"); ok {
		request.MaxSyncStandbyLag = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().ModifyDBInstanceHAConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate postgresql ModifyDBInstanceHAConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPostgresqlInstanceHAConfigRead(d, meta)
}

func resourceTencentCloudPostgresqlInstanceHAConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_instance_ha_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
