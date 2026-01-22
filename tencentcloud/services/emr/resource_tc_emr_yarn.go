package emr

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudEmrYarn() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEmrYarnCreate,
		Read:   resourceTencentCloudEmrYarnRead,
		Update: resourceTencentCloudEmrYarnUpdate,
		Delete: resourceTencentCloudEmrYarnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "EMR Instance ID.",
			},

			"enable_resource_schedule": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the resource scheduling function is enabled.",
			},
			"disable_resource_schedule_sync": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to synchronize when closing.",
			},
			"active_scheduler": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource dispatchers currently in effect.",
			},
			"capacity_global_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Information about capacity scheduler.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_label": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether label scheduling is turned on.",
						},
						"label_dir": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "If label scheduling is enabled, the path where label information is stored.",
						},
						"queue_mapping_override": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether to overwrite the user-specified queue. True means overwriting.",
						},
						"default_settings": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Description: "Advanced settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Name, as the key for the input parameter.",
									},
									"desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Described.",
									},
									"prompt": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Prompt.",
									},
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Key, used for display. This configuration corresponds to the configuration item in the configuration file.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Value corresponding to tame.",
									},
								},
							},
						},
					},
				},
			},
			"fair_global_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Information about fair scheduler.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_max_apps_default": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Corresponding to the page procedural upper limit.",
						},
					},
				},
			},
			"scheduler": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The latest resource scheduler.",
			},
		},
	}
}

func resourceTencentCloudEmrYarnCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_yarn.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId string
	)

	instanceId = d.Get("instance_id").(string)

	d.SetId(instanceId)

	return resourceTencentCloudEmrYarnUpdate(d, meta)
}

func resourceTencentCloudEmrYarnRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_yarn.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()
	_ = d.Set("instance_id", instanceId)
	respData, err := service.DescribeEmrYarnById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `emr_yarn` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if respData.EnableResourceSchedule != nil {
		_ = d.Set("enable_resource_schedule", respData.EnableResourceSchedule)
	}

	if respData.ActiveScheduler != nil {
		_ = d.Set("active_scheduler", respData.ActiveScheduler)
	}

	if respData.CapacityGlobalConfig != nil {
		capacityGlobalConfigMap := map[string]interface{}{}
		capacityGlobalConfig := respData.CapacityGlobalConfig
		if capacityGlobalConfig.EnableLabel != nil {
			capacityGlobalConfigMap["enable_label"] = capacityGlobalConfig.EnableLabel
		}

		if capacityGlobalConfig.LabelDir != nil {
			capacityGlobalConfigMap["label_dir"] = capacityGlobalConfig.LabelDir
		}

		if capacityGlobalConfig.QueueMappingOverride != nil {
			capacityGlobalConfigMap["queue_mapping_override"] = capacityGlobalConfig.QueueMappingOverride
		}

		defaultSettingsList := make([]map[string]interface{}, 0, len(capacityGlobalConfig.DefaultSettings))
		if capacityGlobalConfig.DefaultSettings != nil {
			for _, defaultSettings := range capacityGlobalConfig.DefaultSettings {
				defaultSettingsMap := map[string]interface{}{}

				if defaultSettings.Name != nil {
					defaultSettingsMap["name"] = defaultSettings.Name
				}

				if defaultSettings.Desc != nil {
					defaultSettingsMap["desc"] = defaultSettings.Desc
				}

				if defaultSettings.Prompt != nil {
					defaultSettingsMap["prompt"] = defaultSettings.Prompt
				}

				if defaultSettings.Key != nil {
					defaultSettingsMap["key"] = defaultSettings.Key
				}

				if defaultSettings.Value != nil {
					defaultSettingsMap["value"] = defaultSettings.Value
				}

				defaultSettingsList = append(defaultSettingsList, defaultSettingsMap)
			}

			capacityGlobalConfigMap["default_settings"] = defaultSettingsList
		}
		_ = d.Set("capacity_global_config", []interface{}{capacityGlobalConfigMap})
	}
	if respData.Scheduler != nil {
		_ = d.Set("scheduler", respData.Scheduler)
	}
	if respData.FairGlobalConfig != nil {
		fairGlobalConfig := respData.FairGlobalConfig
		fairGlobalConfigMap := map[string]interface{}{}
		if fairGlobalConfig.UserMaxAppsDefault != nil {
			fairGlobalConfigMap["user_max_apps_default"] = fairGlobalConfig.UserMaxAppsDefault
		}
		_ = d.Set("fair_global_config", []interface{}{fairGlobalConfigMap})
	}

	return nil
}

func resourceTencentCloudEmrYarnUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_yarn.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	instanceId := d.Id()

	request := emr.NewModifyGlobalConfigRequest()
	request.InstanceId = helper.String(instanceId)

	hasChange := false
	if d.HasChange("enable_resource_schedule") {
		if v, ok := d.GetOkExists("enable_resource_schedule"); ok {
			enableResourceScheduleValue := v.(bool)
			request.Items = append(request.Items, &emr.Item{
				Key:   helper.String("enableResourceSchedule"),
				Value: helper.String(strconv.FormatBool(enableResourceScheduleValue)),
			})
			hasChange = true
			if !enableResourceScheduleValue {
				if v1, ok1 := d.GetOk("disable_resource_schedule_sync"); ok1 {
					request.Items = append(request.Items, &emr.Item{
						Key:   helper.String("sync"),
						Value: helper.String(strconv.FormatBool(v1.(bool))),
					})
				} else {
					return fmt.Errorf("disable_resource_schedule_sync must be set when you change enable_resource_schedule")
				}
			}
		}
	}

	if d.HasChange("scheduler") {
		if v, ok := d.GetOk("scheduler"); ok {
			request.Items = append(request.Items, &emr.Item{
				Key:   helper.String("scheduler"),
				Value: helper.String(v.(string)),
			})
			hasChange = true
		}
	}
	if d.HasChange("capacity_global_config.0.enable_label") {
		if v, ok := d.GetOkExists("capacity_global_config.0.enable_label"); ok {
			request.Items = append(request.Items, &emr.Item{
				Key:   helper.String("enableLabel"),
				Value: helper.String(strconv.FormatBool(v.(bool))),
			})
			hasChange = true
		}
	}
	if d.HasChange("capacity_global_config.0.label_dir") {
		if v, ok := d.GetOk("capacity_global_config.0.label_dir"); ok {
			request.Items = append(request.Items, &emr.Item{
				Key:   helper.String("labelDir"),
				Value: helper.String(v.(string)),
			})
			hasChange = true
		}
	}
	if d.HasChange("capacity_global_config.0.queue_mapping_override") {
		if v, ok := d.GetOkExists("capacity_global_config.0.queue_mapping_override"); ok {
			request.Items = append(request.Items, &emr.Item{
				Key:   helper.String("queueMappingOverride"),
				Value: helper.String(strconv.FormatBool(v.(bool))),
			})
			hasChange = true
		}
	}
	if d.HasChange("capacity_global_config.0.default_settings") {
		if v, ok := d.GetOk("capacity_global_config.0.default_settings"); ok {
			defaultSettings := v.(*schema.Set).List()
			for _, defaultSetting := range defaultSettings {
				defaultSettingMap := defaultSetting.(map[string]interface{})
				request.Items = append(request.Items, &emr.Item{
					Key:   helper.String(defaultSettingMap["name"].(string)),
					Value: helper.String(defaultSettingMap["value"].(string)),
				})
				hasChange = true
			}

		}
	}
	if d.HasChange("fair_global_config.0.user_max_apps_default") {
		if v, ok := d.GetOk("fair_global_config.0.user_max_apps_default"); ok {
			request.Items = append(request.Items, &emr.Item{
				Key:   helper.String("userMaxAppsDefault"),
				Value: helper.String(helper.IntToStr(v.(int))),
			})
			hasChange = true
		}
	}

	if hasChange {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().ModifyGlobalConfigWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update emr yarn failed, reason:%+v", logId, err)
			return err
		}

		// deployRequest := emr.NewDeployYarnConfRequest()
		// deployRequest.InstanceId = helper.String(instanceId)
		// var flowId *uint64
		// err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		// 	result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().DeployYarnConfWithContext(ctx, deployRequest)
		// 	if e != nil {
		// 		return tccommon.RetryError(e)
		// 	} else {
		// 		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, deployRequest.GetAction(), deployRequest.ToJsonString(), result.ToJsonString())
		// 	}
		// 	flowId = result.Response.FlowId
		// 	return nil
		// })
		// if err != nil {
		// 	log.Printf("[CRITAL]%s update emr yarn failed, reason:%+v", logId, err)
		// 	return err
		// }

		// if flowId != nil {
		// 	emrService := EMRService{
		// 		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
		// 	}
		// 	conf := tccommon.BuildStateChangeConf([]string{}, []string{"2"}, 10*tccommon.ReadRetryTimeout, time.Second, emrService.FlowStatusRefreshFunc(instanceId, strconv.FormatUint(*flowId, 10), F_KEY_FLOW_ID, []string{}))
		// 	if _, e := conf.WaitForState(); e != nil {
		// 		return e
		// 	}
		// }
	}

	return resourceTencentCloudEmrYarnRead(d, meta)
}

func resourceTencentCloudEmrYarnDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_yarn.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
