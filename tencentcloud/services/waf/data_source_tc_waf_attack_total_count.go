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

func DataSourceTencentCloudWafAttackTotalCount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafAttackTotalCountRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Begin time.",
			},
			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time.",
			},
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Query domain name, all domain use all.",
			},
			"query_string": {
				Optional:    true,
				Type:        schema.TypeString,
				Default:     "",
				Description: "Query conditions.",
			},
			"total_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total number of attacks.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudWafAttackTotalCountRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_waf_attack_total_count.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service          = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		attackTotalCount *waf.GetAttackTotalCountResponseParams
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("query_string"); ok {
		paramMap["QueryString"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafAttackTotalCountByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		attackTotalCount = result
		return nil
	})

	if err != nil {
		return err
	}

	if attackTotalCount.TotalCount != nil {
		_ = d.Set("total_count", attackTotalCount.TotalCount)
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
