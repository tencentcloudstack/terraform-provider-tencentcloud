package waf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wafv20180125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafBotStatusConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafBotStatusConfigCreate,
		Read:   resourceTencentCloudWafBotStatusConfigRead,
		Update: resourceTencentCloudWafBotStatusConfigUpdate,
		Delete: resourceTencentCloudWafBotStatusConfigDelete,
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

			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain.",
			},

			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Bot status. 1 - enable; 0 - disable.",
			},

			"scene_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Scene total count.",
			},

			"valid_scene_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of effective scenarios.",
			},

			"current_global_scene": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The currently enabled scenario with a global matching range and the highest priority.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scene_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scene ID.",
						},
						"scene_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scene name.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Priority.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Update time.",
						},
					},
				},
			},

			"custom_rule_nums": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of custom rules, excluding BOT whitelist.",
			},
		},
	}
}

func resourceTencentCloudWafBotStatusConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_bot_status_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId string
		domain     string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}

	d.SetId(strings.Join([]string{instanceId, domain}, tccommon.FILED_SP))

	return resourceTencentCloudWafBotStatusConfigUpdate(d, meta)
}

func resourceTencentCloudWafBotStatusConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_bot_status_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	domain := idSplit[1]

	respData, err := service.DescribeWafBotStatusConfigById(ctx, domain)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `waf_bot_status_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("domain", domain)
	_ = d.Set("instance_id", instanceId)

	if respData.Status != nil {
		if *respData.Status == true {
			_ = d.Set("status", "1")
		} else {
			_ = d.Set("status", "0")
		}
	}

	if respData.SceneCount != nil {
		_ = d.Set("scene_count", respData.SceneCount)
	}

	if respData.ValidSceneCount != nil {
		_ = d.Set("valid_scene_count", respData.ValidSceneCount)
	}

	if respData.CurrentGlobalScene != nil {
		tmpList := make([]map[string]interface{}, 0, 1)
		dMap := make(map[string]interface{})
		if respData.CurrentGlobalScene.SceneId != nil {
			dMap["scene_id"] = respData.CurrentGlobalScene.SceneId
		}

		if respData.CurrentGlobalScene.SceneName != nil {
			dMap["scene_name"] = respData.CurrentGlobalScene.SceneName
		}

		if respData.CurrentGlobalScene.Priority != nil {
			dMap["priority"] = respData.CurrentGlobalScene.Priority
		}

		if respData.CurrentGlobalScene.UpdateTime != nil {
			dMap["update_time"] = respData.CurrentGlobalScene.UpdateTime
		}

		tmpList = append(tmpList, dMap)
		_ = d.Set("current_global_scene", tmpList)
	}

	if respData.CustomRuleNums != nil {
		_ = d.Set("custom_rule_nums", respData.CustomRuleNums)
	}

	return nil
}

func resourceTencentCloudWafBotStatusConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_bot_status_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = wafv20180125.NewModifyBotStatusRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	domain := idSplit[1]

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	request.InstanceID = &instanceId
	request.Domain = &domain
	request.Category = helper.String("bot")
	request.IsVersionFour = helper.Bool(true)
	request.BotVersion = helper.String("4.1.0")
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyBotStatusWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update waf bot status config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudWafBotStatusConfigRead(d, meta)
}

func resourceTencentCloudWafBotStatusConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_bot_status_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
