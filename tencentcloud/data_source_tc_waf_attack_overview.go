/*
Use this data source to query detailed information of waf attack_overview

Example Usage

```hcl
data "tencentcloud_waf_attack_overview" "attack_overview" {
  from_time = ""
  to_time = ""
  appid =
  domain = ""
  edition = ""
  instance_i_d = ""
              }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudWafAttackOverview() *schema.Resource {
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
				Description: "Appid.",
			},

			"domain": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},

			"edition": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Use &amp;amp;#39;sparta-waf&amp;amp;#39;,&amp;amp;#39;clb-waf&amp;amp;#39;, otherwise not filter.",
			},

			"instance_i_d": {
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

			"a_c_l_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Access control count.",
			},

			"c_c_count": {
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
				Description: "Api risk event count attention：the field may return null，Indicates that no valid value can be obtained.",
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
	defer logElapsed("data_source.tencentcloud_waf_attack_overview.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("from_time"); ok {
		paramMap["FromTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("to_time"); ok {
		paramMap["ToTime"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("appid"); v != nil {
		paramMap["Appid"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("edition"); ok {
		paramMap["Edition"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_i_d"); ok {
		paramMap["InstanceID"] = helper.String(v.(string))
	}

	service := WafService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafAttackOverviewByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		accessCount = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(accessCount))
	if accessCount != nil {
		_ = d.Set("access_count", accessCount)
	}

	if attackCount != nil {
		_ = d.Set("attack_count", attackCount)
	}

	if aCLCount != nil {
		_ = d.Set("a_c_l_count", aCLCount)
	}

	if cCCount != nil {
		_ = d.Set("c_c_count", cCCount)
	}

	if botCount != nil {
		_ = d.Set("bot_count", botCount)
	}

	if apiAssetsCount != nil {
		_ = d.Set("api_assets_count", apiAssetsCount)
	}

	if apiRiskEventCount != nil {
		_ = d.Set("api_risk_event_count", apiRiskEventCount)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
