package tpulsar

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"

	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTdmqEnvironmentAttributes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqEnvironmentAttributesRead,
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Environment (namespace) name.",
			},
			"cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID of the Pulsar cluster.",
			},
			// computed
			"msg_ttl": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Expiration time of unconsumed messages, unit second, maximum 1296000 (15 days).",
			},
			"rate_in_byte": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Consumption rate limit, unit byte/second, 0 unlimited rate.",
			},
			"rate_in_size": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Consumption rate limit, unit number/second, 0 is unlimited.",
			},
			"retention_hours": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Consumed message storage policy, unit hour, 0 will be deleted immediately after consumption.",
			},
			"retention_size": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Consumed message storage strategy, unit G, 0 Delete immediately after consumption.",
			},
			"replicas": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Duplicate number.",
			},
			"remark": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Remark.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTdmqEnvironmentAttributesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tdmq_environment_attributes.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service       = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		tdmqEnv       *tdmq.DescribeEnvironmentAttributesResponseParams
		environmentId string
		clusterId     string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("environment_id"); ok {
		paramMap["EnvironmentId"] = helper.String(v.(string))
		environmentId = v.(string)
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
		clusterId = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqEnvironmentAttributesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		tdmqEnv = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0)
	if tdmqEnv.EnvironmentId != nil {
		_ = d.Set("environment_id", tdmqEnv.EnvironmentId)
	}

	if tdmqEnv.MsgTTL != nil {
		_ = d.Set("msg_ttl", tdmqEnv.MsgTTL)
	}

	if tdmqEnv.RateInByte != nil {
		_ = d.Set("rate_in_byte", tdmqEnv.RateInByte)
	}

	if tdmqEnv.RateInSize != nil {
		_ = d.Set("rate_in_size", tdmqEnv.RateInSize)
	}

	if tdmqEnv.RetentionHours != nil {
		_ = d.Set("retention_hours", tdmqEnv.RetentionHours)
	}

	if tdmqEnv.RetentionSize != nil {
		_ = d.Set("retention_size", tdmqEnv.RetentionSize)
	}

	if tdmqEnv.Replicas != nil {
		_ = d.Set("replicas", tdmqEnv.Replicas)
	}

	if tdmqEnv.Remark != nil {
		_ = d.Set("remark", tdmqEnv.Remark)
	}

	ids = append(ids, environmentId)
	ids = append(ids, clusterId)
	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
