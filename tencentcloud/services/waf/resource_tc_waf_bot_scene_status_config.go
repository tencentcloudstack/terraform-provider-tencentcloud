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

func ResourceTencentCloudWafBotSceneStatusConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafBotSceneStatusConfigCreate,
		Read:   resourceTencentCloudWafBotSceneStatusConfigRead,
		Update: resourceTencentCloudWafBotSceneStatusConfigUpdate,
		Delete: resourceTencentCloudWafBotSceneStatusConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain.",
			},

			"scene_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Scene ID.",
			},

			"status": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Bot status. true - enable; false - disable.",
			},

			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Scene type, default: Default scenario, custom: Non default scenario.",
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
		},
	}
}

func resourceTencentCloudWafBotSceneStatusConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_bot_scene_status_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		domain  string
		sceneId string
	)

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}

	if v, ok := d.GetOk("scene_id"); ok {
		sceneId = v.(string)
	}

	d.SetId(strings.Join([]string{domain, sceneId}, tccommon.FILED_SP))

	return resourceTencentCloudWafBotSceneStatusConfigUpdate(d, meta)
}

func resourceTencentCloudWafBotSceneStatusConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_bot_scene_status_config.read")()
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

	domain := idSplit[0]
	sceneId := idSplit[1]

	respData, err := service.DescribeWafBotSceneStatusConfigById(ctx, domain, sceneId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `waf_bot_scene_status_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("domain", domain)
	_ = d.Set("scene_id", sceneId)

	if respData.SceneStatus != nil {
		_ = d.Set("status", respData.SceneStatus)
	}

	if respData.Type != nil {
		_ = d.Set("type", respData.Type)
	}

	if respData.SceneName != nil {
		_ = d.Set("scene_name", respData.SceneName)
	}

	if respData.Priority != nil {
		_ = d.Set("priority", respData.Priority)
	}

	return nil
}

func resourceTencentCloudWafBotSceneStatusConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_bot_scene_status_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = wafv20180125.NewModifyBotSceneStatusRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	domain := idSplit[0]
	sceneId := idSplit[1]

	if v, ok := d.GetOkExists("status"); ok {
		request.Status = helper.Bool(v.(bool))
	}

	request.Domain = &domain
	request.SceneId = &sceneId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyBotSceneStatusWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update waf bot scene status config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudWafBotSceneStatusConfigRead(d, meta)
}

func resourceTencentCloudWafBotSceneStatusConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_bot_scene_status_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
