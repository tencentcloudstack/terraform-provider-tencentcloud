package mongodb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodbv20190725 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMongodbScaleUpDbInstanceCpu() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbScaleUpDbInstanceCpuCreate,
		Read:   resourceTencentCloudMongodbScaleUpDbInstanceCpuRead,
		Delete: resourceTencentCloudMongodbScaleUpDbInstanceCpuDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "MongoDB instance ID, for example: cmgo-p8vnipr5.",
			},

			"extra_cpu": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Number of CPU cores to scale up (increment). Each node will increase by this number of cores. It cannot be changed after creation.",
			},

			// computed
			"scalable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether elastic CPU scale-up is available.",
			},

			"is_scaled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the instance is currently in scaled state.",
			},

			"is_locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the instance is locked (a flow is in progress).",
			},

			"reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reason why elastic CPU scale-up is not available.",
			},

			"max_extra_cpu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum number of CPU cores that can be scaled up, MIN(minimum shard CPU cores, 24).",
			},

			"trigger_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Scale-up trigger type. Valid values: 1-manual, 2-periodic, 3-one-time time period, 4-monitor.",
			},

			"scale_up_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Scale-up time.",
			},
		},
	}
}

func resourceTencentCloudMongodbScaleUpDbInstanceCpuCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_scale_up_db_instance_cpu.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request    = mongodbv20190725.NewScaleUpDBInstanceCpuRequest()
		response   = mongodbv20190725.NewScaleUpDBInstanceCpuResponse()
		instanceId string
		flowId     string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOkExists("extra_cpu"); ok {
		request.ExtraCpu = helper.IntInt64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().ScaleUpDBInstanceCpuWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("ScaleUpDBInstanceCpu failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create mongodb scale_up_db_instance_cpu failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.FlowId != nil {
		flowId = helper.Int64ToStr(*response.Response.FlowId)
	}

	d.SetId(instanceId)

	// ScaleUpDBInstanceCpu is asynchronous: the returned flow ID is polled until
	// the task finishes, so the resource is only reported as created once the
	// instance is really scaled up.
	if err := service.DescribeAsyncRequestInfo(ctx, flowId, d.Timeout(schema.TimeoutCreate)); err != nil {
		log.Printf("[CRITAL]%s wait mongodb scale_up_db_instance_cpu scale up complete failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMongodbScaleUpDbInstanceCpuRead(d, meta)
}

func resourceTencentCloudMongodbScaleUpDbInstanceCpuRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_scale_up_db_instance_cpu.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	respData, err := service.DescribeMongodbScaleUpDbInstanceCpuById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_mongodb_scale_up_db_instance_cpu` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if respData.Scalable != nil {
		_ = d.Set("scalable", respData.Scalable)
	}

	if respData.IsScaled != nil {
		_ = d.Set("is_scaled", respData.IsScaled)
	}

	if respData.IsLocked != nil {
		_ = d.Set("is_locked", respData.IsLocked)
	}

	if respData.Reason != nil {
		_ = d.Set("reason", respData.Reason)
	}

	if respData.MaxExtraCpu != nil {
		_ = d.Set("max_extra_cpu", respData.MaxExtraCpu)
	}

	if respData.ExtraCpu != nil {
		_ = d.Set("extra_cpu", respData.ExtraCpu)
	}

	if respData.TriggerType != nil {
		_ = d.Set("trigger_type", respData.TriggerType)
	}

	if respData.ScaleUpTime != nil {
		_ = d.Set("scale_up_time", respData.ScaleUpTime)
	}

	return nil
}

func resourceTencentCloudMongodbScaleUpDbInstanceCpuDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_scale_up_db_instance_cpu.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request    = mongodbv20190725.NewScaleDownDBInstanceCpuRequest()
		response   = mongodbv20190725.NewScaleDownDBInstanceCpuResponse()
		instanceId = d.Id()
		flowId     string
	)

	request.InstanceId = helper.String(instanceId)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().ScaleDownDBInstanceCpuWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("ScaleDownDBInstanceCpu failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		// The instance may already be out of the scaled state, for example when the
		// scale-up was disabled manually; destroying the resource is then a no-op.
		if strings.Contains(reqErr.Error(), "not scaled") || strings.Contains(reqErr.Error(), "not in scaled") {
			log.Printf("[INFO]%s mongodb scale_up_db_instance_cpu instance [%s] is not in scaled state, treat as deleted, reason:%+v", logId, instanceId, reqErr)
			return nil
		}

		log.Printf("[CRITAL]%s delete mongodb scale_up_db_instance_cpu failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.FlowId != nil {
		flowId = helper.Int64ToStr(*response.Response.FlowId)
	}

	// ScaleDownDBInstanceCpu is asynchronous: the returned flow ID is polled until
	// the task finishes, so the resource is only reported as destroyed once the
	// instance is really scaled down.
	if err := service.DescribeAsyncRequestInfo(ctx, flowId, d.Timeout(schema.TimeoutDelete)); err != nil {
		log.Printf("[CRITAL]%s wait mongodb scale_up_db_instance_cpu scale down complete failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
