package dayuv2

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAntiddosOverviewIndex() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAntiddosOverviewIndexRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "StartTime.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "EndTime.",
			},

			"all_ip_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "ip count.",
			},

			"antiddos_ip_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total number of advanced defense IPs (including advanced defense packets and advanced defense IPs).",
			},

			"attack_ip_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "AttackIpCount.",
			},

			"block_ip_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "BlockIpCount.",
			},

			"antiddos_domain_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "AntiddosDomainCount.",
			},

			"attack_domain_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "AttackDomainCount.",
			},

			"max_attack_flow": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "MaxAttackFlow.",
			},

			"new_attack_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The time in the most recent attack.",
			},

			"new_attack_ip": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The IP address in the most recent attack.",
			},

			"new_attack_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The type in the most recent attack.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudAntiddosOverviewIndexRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_antiddos_overview_index.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	service := AntiddosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var describeOverviewIndexResponseParams *antiddos.DescribeOverviewIndexResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAntiddosOverviewIndexByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		describeOverviewIndexResponseParams = result
		return nil
	})
	if err != nil {
		return err
	}

	resultMap := make(map[string]interface{})
	if describeOverviewIndexResponseParams.AllIpCount != nil {
		resultMap["all_ip_count"] = describeOverviewIndexResponseParams.AllIpCount
		_ = d.Set("all_ip_count", describeOverviewIndexResponseParams.AllIpCount)
	}

	if describeOverviewIndexResponseParams.AntiddosIpCount != nil {
		resultMap["antiddos_ip_count"] = describeOverviewIndexResponseParams.AntiddosIpCount
		_ = d.Set("antiddos_ip_count", describeOverviewIndexResponseParams.AntiddosIpCount)
	}

	if describeOverviewIndexResponseParams.AttackIpCount != nil {
		resultMap["attack_ip_count"] = describeOverviewIndexResponseParams.AttackIpCount
		_ = d.Set("attack_ip_count", describeOverviewIndexResponseParams.AttackIpCount)
	}

	if describeOverviewIndexResponseParams.BlockIpCount != nil {
		resultMap["block_ip_count"] = describeOverviewIndexResponseParams.BlockIpCount
		_ = d.Set("block_ip_count", describeOverviewIndexResponseParams.BlockIpCount)
	}

	if describeOverviewIndexResponseParams.AntiddosDomainCount != nil {
		resultMap["antiddos_domain_count"] = describeOverviewIndexResponseParams.AntiddosDomainCount
		_ = d.Set("antiddos_domain_count", describeOverviewIndexResponseParams.AntiddosDomainCount)
	}

	if describeOverviewIndexResponseParams.AttackDomainCount != nil {
		resultMap["attack_domain_count"] = describeOverviewIndexResponseParams.AttackDomainCount
		_ = d.Set("attack_domain_count", describeOverviewIndexResponseParams.AttackDomainCount)
	}

	if describeOverviewIndexResponseParams.MaxAttackFlow != nil {
		resultMap["max_attack_flow"] = describeOverviewIndexResponseParams.MaxAttackFlow
		_ = d.Set("max_attack_flow", describeOverviewIndexResponseParams.MaxAttackFlow)
	}

	if describeOverviewIndexResponseParams.NewAttackTime != nil {
		resultMap["new_attack_time"] = describeOverviewIndexResponseParams.NewAttackTime
		_ = d.Set("new_attack_time", describeOverviewIndexResponseParams.NewAttackTime)
	}

	if describeOverviewIndexResponseParams.NewAttackIp != nil {
		resultMap["new_attack_ip"] = describeOverviewIndexResponseParams.NewAttackIp
		_ = d.Set("new_attack_ip", describeOverviewIndexResponseParams.NewAttackIp)
	}

	if describeOverviewIndexResponseParams.NewAttackType != nil {
		resultMap["new_attack_type"] = describeOverviewIndexResponseParams.NewAttackType
		_ = d.Set("new_attack_type", describeOverviewIndexResponseParams.NewAttackType)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), resultMap); e != nil {
			return e
		}
	}
	return nil
}
