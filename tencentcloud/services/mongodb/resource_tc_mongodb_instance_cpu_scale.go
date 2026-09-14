package mongodb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudMongodbInstanceCpuScaleConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbInstanceCpuScaleConfigCreate,
		Read:   resourceTencentCloudMongodbInstanceCpuScaleConfigRead,
		Update: resourceTencentCloudMongodbInstanceCpuScaleConfigUpdate,
		Delete: resourceTencentCloudMongodbInstanceCpuScaleConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "MongoDB instance ID, for example: cmgo-xxxxxxxx.",
			},

			"extra_cpu": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Extra CPU cores for elastic scaling. Each node will be scaled up by this number of cores. Setting this value enables CPU elastic scaling. Removing this field or setting it to 0 disables CPU elastic scaling.",
			},

			"flow_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Task flow ID returned by the last operation.",
			},
		},
	}
}

func resourceTencentCloudMongodbInstanceCpuScaleConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_cpu_scale.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	if v, ok := d.GetOk("extra_cpu"); ok {
		extraCpu := int64(v.(int))
		flowId, err := service.ScaleUpMongodbDBInstanceCpu(ctx, instanceId, extraCpu)
		if err != nil {
			log.Printf("[CRITAL]%s mongodb_instance_cpu_scale create scale up cpu failed, reason:%+v", logId, err)
			return err
		}

		_ = d.Set("flow_id", flowId)

		// Wait for async task completion
		flowIdString := helper.Int64ToStr(flowId)
		if err = service.DescribeAsyncRequestInfo(ctx, flowIdString, 3*tccommon.ReadRetryTimeout); err != nil {
			log.Printf("[CRITAL]%s mongodb_instance_cpu_scale create wait async task failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMongodbInstanceCpuScaleConfigRead(d, meta)
}

func resourceTencentCloudMongodbInstanceCpuScaleConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_cpu_scale.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	instance, has, err := service.DescribeMongodbInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if !has || instance == nil {
		log.Printf("[CRUD] mongodb_instance_cpu_scale id=%s", d.Id())
		d.SetId("")
		log.Printf("[WARN]%s resource `mongodb_instance_cpu_scale` [%s] not found, please check if it has been deleted.\n", logId, instanceId)
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	return nil
}

func resourceTencentCloudMongodbInstanceCpuScaleConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_cpu_scale.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	if d.HasChange("extra_cpu") {
		oldValue, newValue := d.GetChange("extra_cpu")
		oldExtraCpu := int64(oldValue.(int))
		newExtraCpu := int64(newValue.(int))

		if newExtraCpu > 0 {
			// Scale up CPU
			flowId, err := service.ScaleUpMongodbDBInstanceCpu(ctx, instanceId, newExtraCpu)
			if err != nil {
				log.Printf("[CRITAL]%s mongodb_instance_cpu_scale update scale up cpu failed, reason:%+v", logId, err)
				return err
			}

			_ = d.Set("flow_id", flowId)

			flowIdString := helper.Int64ToStr(flowId)
			if err = service.DescribeAsyncRequestInfo(ctx, flowIdString, 3*tccommon.ReadRetryTimeout); err != nil {
				log.Printf("[CRITAL]%s mongodb_instance_cpu_scale update wait async task failed, reason:%+v", logId, err)
				return err
			}
		} else if oldExtraCpu > 0 && newExtraCpu == 0 {
			// Scale down CPU (close elastic scaling)
			flowId, err := service.ScaleDownMongodbDBInstanceCpu(ctx, instanceId)
			if err != nil {
				log.Printf("[CRITAL]%s mongodb_instance_cpu_scale update scale down cpu failed, reason:%+v", logId, err)
				return err
			}

			_ = d.Set("flow_id", flowId)

			flowIdString := helper.Int64ToStr(flowId)
			if err = service.DescribeAsyncRequestInfo(ctx, flowIdString, 3*tccommon.ReadRetryTimeout); err != nil {
				log.Printf("[CRITAL]%s mongodb_instance_cpu_scale update wait async task failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudMongodbInstanceCpuScaleConfigRead(d, meta)
}

func resourceTencentCloudMongodbInstanceCpuScaleConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_cpu_scale.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	// Scale down CPU when deleting the resource
	flowId, err := service.ScaleDownMongodbDBInstanceCpu(ctx, instanceId)
	if err != nil {
		log.Printf("[CRITAL]%s mongodb_instance_cpu_scale delete scale down cpu failed, reason:%+v", logId, err)
		return err
	}

	_ = d.Set("flow_id", flowId)

	flowIdString := helper.Int64ToStr(flowId)
	if err = service.DescribeAsyncRequestInfo(ctx, flowIdString, 3*tccommon.ReadRetryTimeout); err != nil {
		log.Printf("[CRITAL]%s mongodb_instance_cpu_scale delete wait async task failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
