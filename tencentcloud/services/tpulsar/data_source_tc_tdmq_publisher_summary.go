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

func DataSourceTencentCloudTdmqPublisherSummary() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqPublisherSummaryRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
			"namespace": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "namespace name.",
			},
			"topic": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "subject name.",
			},
			// computed
			"msg_rate_in": {
				Computed:    true,
				Type:        schema.TypeFloat,
				Description: "Production rate (units per second)Note: This field may return null, indicating that no valid value can be obtained.",
			},
			"msg_throughput_in": {
				Computed:    true,
				Type:        schema.TypeFloat,
				Description: "Production rate (bytes per second)Note: This field may return null, indicating that no valid value can be obtained.",
			},
			"publisher_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "number of producersNote: This field may return null, indicating that no valid value can be obtained.",
			},
			"storage_size": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Message store size in bytesNote: This field may return null, indicating that no valid value can be obtained.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTdmqPublisherSummaryRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tdmq_publisher_summary.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service          = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		publisherSummary *tdmq.DescribePublisherSummaryResponseParams
		clusterId        string
		Namespace        string
		Topic            string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
		Namespace = v.(string)
	}

	if v, ok := d.GetOk("topic"); ok {
		paramMap["Topic"] = helper.String(v.(string))
		Topic = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqPublisherSummaryByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		publisherSummary = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0)

	if publisherSummary.MsgRateIn != nil {
		_ = d.Set("msg_rate_in", publisherSummary.MsgRateIn)
	}

	if publisherSummary.MsgThroughputIn != nil {
		_ = d.Set("msg_throughput_in", publisherSummary.MsgThroughputIn)
	}

	if publisherSummary.PublisherCount != nil {
		_ = d.Set("publisher_count", publisherSummary.PublisherCount)
	}

	if publisherSummary.StorageSize != nil {
		_ = d.Set("storage_size", publisherSummary.StorageSize)
	}

	ids = append(ids, clusterId)
	ids = append(ids, Namespace)
	ids = append(ids, Topic)
	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
