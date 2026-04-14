package config

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudConfigRecorderConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudConfigRecorderConfigCreate,
		Read:   resourceTencentCloudConfigRecorderConfigRead,
		Update: resourceTencentCloudConfigRecorderConfigUpdate,
		Delete: resourceTencentCloudConfigRecorderConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable resource monitoring. true: enable (OpenConfigRecorder), false: disable (CloseConfigRecorder).",
			},

			"resource_types": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Resource type list to monitor (e.g. QCS::CAM::Group, QCS::CVM::Instance).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// Computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Recorder creation time.",
			},

			"trigger_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of snapshots taken today.",
			},

			"open_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of times monitoring was opened today.",
			},

			"update_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of monitoring range updates today.",
			},
		},
	}
}

func resourceTencentCloudConfigRecorderConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_recorder_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(helper.BuildToken())

	return resourceTencentCloudConfigRecorderConfigUpdate(d, meta)
}

func resourceTencentCloudConfigRecorderConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_recorder_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeConfigRecorder(ctx)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_config_recorder_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Status != nil {
		_ = d.Set("status", *respData.Status == 1)
	}

	if respData.Items != nil {
		resourceTypes := make([]string, 0, len(respData.Items))
		for _, item := range respData.Items {
			if item.ResourceType != nil {
				resourceTypes = append(resourceTypes, *item.ResourceType)
			}
		}

		_ = d.Set("resource_types", resourceTypes)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.TriggerCount != nil {
		_ = d.Set("trigger_count", int(*respData.TriggerCount))
	}

	if respData.OpenCount != nil {
		_ = d.Set("open_count", int(*respData.OpenCount))
	}

	if respData.UpdateCount != nil {
		_ = d.Set("update_count", int(*respData.UpdateCount))
	}

	return nil
}

func resourceTencentCloudConfigRecorderConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_recorder_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	if d.HasChange("status") || d.IsNewResource() {
		enable := d.Get("status").(bool)
		if enable {
			request := configv20220802.NewOpenConfigRecorderRequest()
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().OpenConfigRecorderWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				}

				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s open config recorder failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		} else {
			request := configv20220802.NewCloseConfigRecorderRequest()
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().CloseConfigRecorderWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				}

				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s close config recorder failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		}
	}

	if d.HasChange("resource_types") || d.IsNewResource() {
		request := configv20220802.NewUpdateConfigRecorderRequest()
		if v, ok := d.GetOk("resource_types"); ok {
			rawList := v.([]interface{})
			resourceTypes := make([]*string, 0, len(rawList))
			for _, item := range rawList {
				val := item.(string)
				resourceTypes = append(resourceTypes, &val)
			}

			request.ResourceTypes = resourceTypes
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().UpdateConfigRecorderWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update config recorder failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudConfigRecorderConfigRead(d, meta)
}

func resourceTencentCloudConfigRecorderConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_recorder_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
