package mongodb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func ResourceTencentCloudMongodbScaleUpDbInstanceCpu() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbScaleUpDbInstanceCpuCreate,
		Read:   resourceTencentCloudMongodbScaleUpDbInstanceCpuRead,
		Update: resourceTencentCloudMongodbScaleUpDbInstanceCpuUpdate,
		Delete: resourceTencentCloudMongodbScaleUpDbInstanceCpuDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "MongoDB instance ID, for example: cmgo-p8vnipr5.",
			},

			"extra_cpu": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of CPU cores to scale up (increment). Each node will increase by this number of cores.",
			},

			"flow_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Task flow ID.",
			},

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
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId string
		extraCpu   int64
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	} else {
		return fmt.Errorf("instance_id is required")
	}

	if v, ok := d.GetOk("extra_cpu"); ok {
		extraCpu = int64(v.(int))
	}

	flowId, err := service.ScaleUpDBInstanceCpu(ctx, instanceId, extraCpu)
	if err != nil {
		log.Printf("[CRITAL]%s create mongodb scale_up_db_instance_cpu failed, reason: %v", logId, err)
		return err
	}

	if flowId == nil {
		return fmt.Errorf("ScaleUpDBInstanceCpu FlowId is nil, logId=%s", logId)
	}

	log.Printf("[INFO]%s mongodb scale_up_db_instance_cpu created, instanceId=%s, flowId=%d", logId, instanceId, *flowId)
	d.SetId(instanceId)
	_ = d.Set("flow_id", *flowId)

	if err = service.WaitMongodbScaleUpDbInstanceCpuComplete(ctx, instanceId, true); err != nil {
		log.Printf("[CRITAL]%s wait mongodb scale_up_db_instance_cpu scale up complete failed, reason: %v", logId, err)
		return err
	}

	return resourceTencentCloudMongodbScaleUpDbInstanceCpuRead(d, meta)
}

func resourceTencentCloudMongodbScaleUpDbInstanceCpuRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_scale_up_db_instance_cpu.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	response, err := service.CheckDBInstanceElasticCpuScalable(ctx, instanceId)
	if err != nil {
		return err
	}

	if response == nil || response.Response == nil {
		log.Printf("[CRUD] tencentcloud_mongodb_scale_up_db_instance_cpu id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if response.Response.Scalable != nil {
		_ = d.Set("scalable", *response.Response.Scalable)
	}

	if response.Response.IsScaled != nil {
		_ = d.Set("is_scaled", *response.Response.IsScaled)
	}

	if response.Response.IsLocked != nil {
		_ = d.Set("is_locked", *response.Response.IsLocked)
	}

	if response.Response.Reason != nil {
		_ = d.Set("reason", *response.Response.Reason)
	}

	if response.Response.MaxExtraCpu != nil {
		_ = d.Set("max_extra_cpu", *response.Response.MaxExtraCpu)
	}

	if response.Response.ExtraCpu != nil {
		_ = d.Set("extra_cpu", *response.Response.ExtraCpu)
	}

	if response.Response.TriggerType != nil {
		_ = d.Set("trigger_type", *response.Response.TriggerType)
	}

	if response.Response.ScaleUpTime != nil {
		_ = d.Set("scale_up_time", *response.Response.ScaleUpTime)
	}

	return nil
}

func resourceTencentCloudMongodbScaleUpDbInstanceCpuUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_scale_up_db_instance_cpu.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	if d.HasChange("extra_cpu") {
		var extraCpu int64
		if v, ok := d.GetOk("extra_cpu"); ok {
			extraCpu = int64(v.(int))
		}

		flowId, err := service.ScaleUpDBInstanceCpu(ctx, instanceId, extraCpu)
		if err != nil {
			log.Printf("[CRITAL]%s update mongodb scale_up_db_instance_cpu failed, reason: %v", logId, err)
			return err
		}

		if flowId != nil {
			_ = d.Set("flow_id", *flowId)
		}

		if err = service.WaitMongodbScaleUpDbInstanceCpuComplete(ctx, instanceId, true); err != nil {
			log.Printf("[CRITAL]%s wait mongodb scale_up_db_instance_cpu scale up complete failed, reason: %v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMongodbScaleUpDbInstanceCpuRead(d, meta)
}

func resourceTencentCloudMongodbScaleUpDbInstanceCpuDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_scale_up_db_instance_cpu.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	_, err := service.ScaleDownDBInstanceCpu(ctx, instanceId)
	if err != nil {
		if strings.Contains(err.Error(), "not") || strings.Contains(err.Error(), "NoChange") ||
			strings.Contains(err.Error(), "not scaled") || strings.Contains(err.Error(), "not in scaled") {
			log.Printf("[INFO]%s mongodb scale_up_db_instance_cpu instance %s not in scaled state, treat as deleted, reason: %v", logId, instanceId, err)
			return nil
		}
		log.Printf("[CRITAL]%s delete mongodb scale_up_db_instance_cpu failed, reason: %v", logId, err)
		return err
	}

	if err = service.WaitMongodbScaleUpDbInstanceCpuComplete(ctx, instanceId, false); err != nil {
		log.Printf("[CRITAL]%s wait mongodb scale_up_db_instance_cpu scale down complete failed, reason: %v", logId, err)
		return err
	}

	return nil
}
