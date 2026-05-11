package crs

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudRedisAuditLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisAuditLogCreate,
		Read:   resourceTencentCloudRedisAuditLogRead,
		Update: resourceTencentCloudRedisAuditLogUpdate,
		Delete: resourceTencentCloudRedisAuditLogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID, such as: crs-xjhsdj****, which can be copied from the instance list in the Redis console.",
			},

			"log_sub_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Log sub-type. Valid values: `write` (write commands), `read` (read commands), `all` (read and write commands).",
			},

			"log_expire_day": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Log retention period in days. Valid values: `7` (7 days), `30` (30 days).",
			},

			"high_log_expire_day": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "High-frequency log retention period in days. Valid values: `7` (7 days).",
			},

			"degrade_strategy": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     500,
				Description: "Degradation strategy threshold in milliseconds. When the instance P99 latency reaches this threshold, audit logs will be automatically discarded to ensure business availability. Value range: [300, 1000]. Default value: `500`.",
			},
		},
	}
}

func resourceTencentCloudRedisAuditLogCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_audit_log.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = redis.NewOpenLogRequest()
		response   = redis.NewOpenLogResponse()
		instanceId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	request.LogType = helper.String("auditLog")

	if v, ok := d.GetOk("log_sub_type"); ok {
		request.LogSubType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("log_expire_day"); ok {
		request.LogExpireDay = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("high_log_expire_day"); ok {
		request.HighLogExpireDay = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("degrade_strategy"); ok {
		request.DegradeStrategy = helper.IntInt64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().OpenLogWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create redis audit log failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create redis audit log failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	_ = response
	d.SetId(instanceId)

	// Poll DescribeLogInstanceList until Status is "normal", indicating the async task has completed.
	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	pollErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, e := service.DescribeRedisAuditLogById(ctx, instanceId)
		if e != nil {
			return resource.NonRetryableError(fmt.Errorf("polling DescribeLogInstanceList failed: %w", e))
		}

		if instance == nil {
			return resource.RetryableError(fmt.Errorf("audit log for instance [%s] not found yet, retrying", instanceId))
		}

		if instance.Status == nil || *instance.Status != "normal" {
			return resource.RetryableError(fmt.Errorf("audit log for instance [%s] is not ready yet, Status: %v", instanceId, instance.Status))
		}

		return nil
	})

	if pollErr != nil {
		log.Printf("[CRITAL]%s polling audit log status failed, reason:%+v", logId, pollErr)
		return pollErr
	}

	return resourceTencentCloudRedisAuditLogRead(d, meta)
}

func resourceTencentCloudRedisAuditLogRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_audit_log.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	instanceId := d.Id()

	respData, err := service.DescribeRedisAuditLogById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_redis_audit_log` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if respData.LogSubType != nil {
		_ = d.Set("log_sub_type", respData.LogSubType)
	}

	if respData.LogExpireDay != nil {
		_ = d.Set("log_expire_day", respData.LogExpireDay)
	}

	if respData.HighLogExpireDay != nil {
		_ = d.Set("high_log_expire_day", respData.HighLogExpireDay)
	}

	return nil
}

func resourceTencentCloudRedisAuditLogUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_audit_log.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	instanceId := d.Id()

	needChange := false
	mutableArgs := []string{"log_sub_type", "log_expire_day", "high_log_expire_day", "degrade_strategy"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := redis.NewModifyLogRequest()
		request.InstanceId = &instanceId
		request.LogType = helper.String("auditLog")

		if v, ok := d.GetOk("log_sub_type"); ok {
			request.LogSubType = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("log_expire_day"); ok {
			request.LogExpireDay = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("high_log_expire_day"); ok {
			request.HighLogExpireDay = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("degrade_strategy"); ok {
			request.DegradeStrategy = helper.IntInt64(v.(int))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().ModifyLogWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify redis audit log failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update redis audit log failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// Poll DescribeLogInstanceList until Status is "normal", indicating the async task has completed.
		service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		pollErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			instance, e := service.DescribeRedisAuditLogById(ctx, instanceId)
			if e != nil {
				return resource.NonRetryableError(fmt.Errorf("polling DescribeLogInstanceList failed: %w", e))
			}

			if instance == nil {
				return resource.RetryableError(fmt.Errorf("audit log for instance [%s] not found yet, retrying", instanceId))
			}

			if instance.Status == nil || *instance.Status != "normal" {
				return resource.RetryableError(fmt.Errorf("audit log for instance [%s] is not ready yet, Status: %v", instanceId, instance.Status))
			}

			return nil
		})

		if pollErr != nil {
			log.Printf("[CRITAL]%s polling audit log status failed, reason:%+v", logId, pollErr)
			return pollErr
		}
	}

	return resourceTencentCloudRedisAuditLogRead(d, meta)
}

func resourceTencentCloudRedisAuditLogDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_audit_log.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = redis.NewCloseLogRequest()
	)

	instanceId := d.Id()

	request.InstanceId = &instanceId
	request.LogType = helper.String("auditLog")

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().CloseLogWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete redis audit log failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete redis audit log failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Poll DescribeLogInstanceList until no data is returned, indicating the audit log has been closed.
	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	pollErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, e := service.DescribeRedisAuditLogById(ctx, instanceId)
		if e != nil {
			return resource.NonRetryableError(fmt.Errorf("polling DescribeLogInstanceList failed: %w", e))
		}

		if instance != nil {
			return resource.RetryableError(fmt.Errorf("audit log for instance [%s] is still closing, retrying", instanceId))
		}

		return nil
	})

	if pollErr != nil {
		log.Printf("[CRITAL]%s polling audit log close status failed, reason:%+v", logId, pollErr)
		return pollErr
	}

	return nil
}
