package waf

import (
	"context"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWafAttackOverview() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafAttackOverviewRead,
		Schema: map[string]*schema.Schema{
			"from_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Begin time.",
			},
			"to_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time.",
			},
			"appid": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "App id.",
			},
			"domain": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},
			"edition": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue(EDITION_TYPE),
				Description:  "support `sparta-waf`, `clb-waf`, otherwise not filter.",
			},
			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Waf instanceId, otherwise not filter.",
			},
			"access_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Access count.",
			},
			"attack_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Attack count.",
			},
			"acl_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Access control count.",
			},
			"cc_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "CC attack count.",
			},
			"bot_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Bot attack count.",
			},
			"api_assets_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Api asset count.",
			},
			"api_risk_event_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of API risk events.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudWafAttackOverviewRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_waf_attack_overview.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service        = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		attackOverview *waf.DescribeAttackOverviewResponseParams
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("from_time"); ok {
		paramMap["FromTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("to_time"); ok {
		paramMap["ToTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("appid"); ok {
		paramMap["Appid"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("edition"); ok {
		paramMap["Edition"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceID"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafAttackOverviewByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		attackOverview = result
		return nil
	})

	if err != nil {
		return err
	}

	if attackOverview.AccessCount != nil {
		_ = d.Set("access_count", attackOverview.AccessCount)
	}

	if attackOverview.AttackCount != nil {
		_ = d.Set("attack_count", attackOverview.AttackCount)
	}

	if attackOverview.ACLCount != nil {
		_ = d.Set("acl_count", attackOverview.ACLCount)
	}

	if attackOverview.CCCount != nil {
		_ = d.Set("cc_count", attackOverview.CCCount)
	}

	if attackOverview.BotCount != nil {
		_ = d.Set("bot_count", attackOverview.BotCount)
	}

	if attackOverview.ApiAssetsCount != nil {
		_ = d.Set("api_assets_count", attackOverview.ApiAssetsCount)
	}

	if attackOverview.ApiRiskEventCount != nil {
		_ = d.Set("api_risk_event_count", attackOverview.ApiRiskEventCount)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
