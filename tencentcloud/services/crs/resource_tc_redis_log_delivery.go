package crs

import (
	"context"
	"log"

	svcCls "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudRedisLogDelivery() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisLogDeliveryCreate,
		Read:   resourceTencentCloudRedisLogDeliveryRead,
		Update: resourceTencentCloudRedisLogDeliveryUpdate,
		Delete: resourceTencentCloudRedisLogDeliveryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID.",
			},

			"logset_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"logset_name"},
				Description:   "The ID of the log set being delivered.",
			},

			"topic_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"topic_name"},
				Description:   "The ID of the topic being delivered.",
			},

			"logset_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"logset_id"},
				Computed:      true,
				Description:   "Log set name. If LogsetId does not specify a specific log set ID, please configure this parameter to set the log set name, and the system will automatically create a new log set with the specified name.",
			},

			"topic_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"topic_id"},
				Computed:      true,
				Description:   "Log topic name, required when TopicId is empty, a new log topic will be automatically created.",
			},

			"log_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the log set is located; if not specified, the region where the instance is located will be used by default.",
			},

			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Log storage time, defaults to 30 days, with an optional range of 1-3600 days.",
			},

			"create_index": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to create an index when creating a log topic.",
			},
		},
	}
}

func resourceTencentCloudRedisLogDeliveryCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_log_delivery.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		instanceId string
		request    = redis.NewModifyInstanceLogDeliveryRequest()
		response   = redis.NewModifyInstanceLogDeliveryResponse()
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	request.LogType = helper.String("slowlog")
	request.Enabled = helper.Bool(true)

	if v, ok := d.GetOk("logset_id"); ok {
		request.LogsetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_id"); ok {
		request.TopicId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("logset_name"); ok {
		request.LogsetName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_name"); ok {
		request.TopicName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_region"); ok {
		request.LogRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("create_index"); ok {
		request.CreateIndex = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().ModifyInstanceLogDeliveryWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create RedisLogDelivery failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	d.SetId(instanceId)

	return resourceTencentCloudRedisLogDeliveryRead(d, meta)
}

func resourceTencentCloudRedisLogDeliveryRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_log_delivery.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	clsService := svcCls.NewClsService(client)

	instanceId := d.Id()

	respData, err := service.DescribeRedisLogDeliveryById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisLogDelivery` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if respData.LogsetId != nil {
		logsetId := respData.LogsetId
		_ = d.Set("logset_id", logsetId)

		logset, err := clsService.DescribeClsLogset(ctx, *logsetId)
		if err != nil {
			return err
		}

		if logset != nil {
			_ = d.Set("logset_name", logset.LogsetName)
		}

	}

	if respData.TopicId != nil {
		topicId := respData.TopicId
		_ = d.Set("topic_id", topicId)

		topic, err := clsService.DescribeClsTopicById(ctx, *topicId)
		if err != nil {
			return err
		}

		if topic != nil {
			if topic.TopicName != nil {
				_ = d.Set("topic_name", topic.TopicName)
			}
			if topic.TopicName != nil {
				_ = d.Set("topic_name", topic.TopicName)
			}
			if topic.Period != nil {
				_ = d.Set("period", topic.Period)
			}
			if topic.Index != nil {
				_ = d.Set("create_index", topic.Index)
			}
		}
	}

	if respData.LogRegion != nil {
		_ = d.Set("log_region", respData.LogRegion)
	}

	_ = instanceId
	return nil
}

func resourceTencentCloudRedisLogDeliveryUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_log_delivery.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	instanceId := d.Id()

	needChange := false
	mutableArgs := []string{"logset_id", "topic_id", "logset_name", "topic_name"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		var (
			request = redis.NewModifyInstanceLogDeliveryRequest()
		)

		request.InstanceId = helper.String(instanceId)
		request.LogType = helper.String("slowlog")

		request.Enabled = helper.Bool(true)

		if d.HasChange("topic_id") || d.HasChange("logset_id") {
			if v, ok := d.GetOk("topic_id"); ok {
				request.TopicId = helper.String(v.(string))
			}

			if v, ok := d.GetOk("logset_id"); ok {
				request.LogsetId = helper.String(v.(string))
			}

		}

		if d.HasChange("logset_name") || d.HasChange("topic_name") {

			if v, ok := d.GetOk("logset_name"); ok {
				request.LogsetName = helper.String(v.(string))
			}

			if v, ok := d.GetOk("topic_name"); ok {
				request.TopicName = helper.String(v.(string))
			}

			if d.HasChange("topic_name") && !d.HasChange("logset_name") && !d.HasChange("logset_id") {
				if v, ok := d.GetOk("logset_id"); ok {
					request.LogsetId = helper.String(v.(string))
				}
			}

			if v, ok := d.GetOk("log_region"); ok {
				request.LogRegion = helper.String(v.(string))
			}

			if v, ok := d.GetOkExists("period"); ok {
				request.Period = helper.IntInt64(v.(int))
			}

			if v, ok := d.GetOkExists("create_index"); ok {
				request.CreateIndex = helper.Bool(v.(bool))
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().ModifyInstanceLogDeliveryWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			_ = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create RedisLogDelivery failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = instanceId
	return resourceTencentCloudRedisLogDeliveryRead(d, meta)
}

func resourceTencentCloudRedisLogDeliveryDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_log_delivery.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	instanceId := d.Id()

	var (
		request  = redis.NewModifyInstanceLogDeliveryRequest()
		response = redis.NewModifyInstanceLogDeliveryResponse()
	)

	request.InstanceId = helper.String(instanceId)
	request.LogType = helper.String("slowlog")
	request.Enabled = helper.Bool(false)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().ModifyInstanceLogDeliveryWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete RedisLogDelivery failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	_ = instanceId
	return nil
}
