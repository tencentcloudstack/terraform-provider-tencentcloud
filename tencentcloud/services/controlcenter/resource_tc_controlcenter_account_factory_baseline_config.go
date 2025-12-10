package controlcenter

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	controlcenterv20230110 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/controlcenter/v20230110"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudControlcenterAccountFactoryBaselineConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudControlcenterAccountFactoryBaselineConfigCreate,
		Read:   resourceTencentCloudControlcenterAccountFactoryBaselineConfigRead,
		Update: resourceTencentCloudControlcenterAccountFactoryBaselineConfigUpdate,
		Delete: resourceTencentCloudControlcenterAccountFactoryBaselineConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Baseline name, which must be unique. Supports only English letters, numbers, Chinese characters, and symbols @, &, _, [], -. Combination of 1-25 Chinese or English characters.",
			},

			"baseline_config_items": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Baseline configuration, overwrite update. You can query existing baseline configurations via controlcenter:GetAccountFactoryBaseline. You can query supported baseline lists via controlcenter:ListAccountFactoryBaselineItems.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the unique identifier for account factory baseline item, can only contain `english letters`, `digits`, and `@,._[]-:()()[]+=.`, with a length of 2-128 characters.",
						},
						"configuration": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Account factory baseline item configuration, different baseline items have different configuration parameters.",
						},
						"apply_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the number of accounts for baseline applications.",
						},
					},
				},
			},

			// computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time.",
			},
		},
	}
}

func resourceTencentCloudControlcenterAccountFactoryBaselineConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_controlcenter_account_factory_baseline_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(helper.BuildToken())

	return resourceTencentCloudControlcenterAccountFactoryBaselineConfigUpdate(d, meta)
}

func resourceTencentCloudControlcenterAccountFactoryBaselineConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_controlcenter_account_factory_baseline_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ControlcenterService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeControlcenterAccountFactoryBaselineConfigById(ctx)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_controlcenter_account_factory_baseline_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.BaselineConfigItems != nil && len(respData.BaselineConfigItems) > 0 {
		baselineConfigItemsList := make([]map[string]interface{}, 0, len(respData.BaselineConfigItems))
		for _, baselineConfigItems := range respData.BaselineConfigItems {
			baselineConfigItemsMap := map[string]interface{}{}
			if baselineConfigItems.Identifier != nil {
				baselineConfigItemsMap["identifier"] = baselineConfigItems.Identifier
			}

			if baselineConfigItems.Configuration != nil {
				baselineConfigItemsMap["configuration"] = baselineConfigItems.Configuration
			}

			if baselineConfigItems.ApplyCount != nil {
				baselineConfigItemsMap["apply_count"] = baselineConfigItems.ApplyCount
			}

			baselineConfigItemsList = append(baselineConfigItemsList, baselineConfigItemsMap)
		}

		_ = d.Set("baseline_config_items", baselineConfigItemsList)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	return nil
}

func resourceTencentCloudControlcenterAccountFactoryBaselineConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_controlcenter_account_factory_baseline_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = controlcenterv20230110.NewUpdateAccountFactoryBaselineRequest()
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("baseline_config_items"); ok {
		for _, item := range v.(*schema.Set).List() {
			baselineConfigItemsMap := item.(map[string]interface{})
			baselineConfigItem := controlcenterv20230110.BaselineConfigItem{}
			if v, ok := baselineConfigItemsMap["identifier"].(string); ok && v != "" {
				baselineConfigItem.Identifier = helper.String(v)
			}

			if v, ok := baselineConfigItemsMap["configuration"].(string); ok && v != "" {
				baselineConfigItem.Configuration = helper.String(v)
			}

			request.BaselineConfigItems = append(request.BaselineConfigItems, &baselineConfigItem)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseControlcenterV20230110Client().UpdateAccountFactoryBaselineWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update controlcenter account factory baseline config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudControlcenterAccountFactoryBaselineConfigRead(d, meta)
}

func resourceTencentCloudControlcenterAccountFactoryBaselineConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_controlcenter_account_factory_baseline_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
