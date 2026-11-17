package cat

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCatProbeMetricTagValues() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCatProbeMetricTagValuesRead,
		Schema: map[string]*schema.Schema{
			"analyze_task_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Analysis of task type, supported types: `AnalyzeTaskType_Network`: network quality, `AnalyzeTaskType_Browse`: page performance, `AnalyzeTaskType_Transport`: port performance, `AnalyzeTaskType_UploadDownload`: file transport, `AnalyzeTaskType_MediaStream`: audiovisual experience.",
			},

			"key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Dimension tag value, reference: `host`: task domain, `errorInfo`: status type, `area`: probe point area, `operator`: probe point operator, `taskId`: task ID.",
			},

			"filter": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter conditions can be passed as a single filter or multiple parameters concatenated together, support regular matching.",
			},

			"filters": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Multiple condition filtering, supports combining multiple filtering conditions for query.",
			},

			"time_range": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Time range.",
			},

			"tag_value_set": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Tag value serialized string.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCatProbeMetricTagValuesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cat_probe_metric_tag_values.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("analyze_task_type"); ok {
		paramMap["AnalyzeTaskType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("key"); ok {
		paramMap["Key"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter"); ok {
		paramMap["Filter"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.(*schema.Set).List()
		paramMap["Filters"] = helper.InterfacesStringsPoint(filtersSet)
	}

	if v, ok := d.GetOk("time_range"); ok {
		paramMap["TimeRange"] = helper.String(v.(string))
	}

	service := CatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var tagValues *cat.DescribeProbeMetricTagValuesResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCatProbeMetricTagValuesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil {
			log.Printf("[DATASOURCE] read empty, skip SetId")
			return resource.NonRetryableError(e)
		}
		tagValues = result
		return nil
	})
	if err != nil {
		return err
	}

	var tagValueSet string
	if tagValues != nil && tagValues.TagValueSet != nil {
		tagValueSet = *tagValues.TagValueSet
		_ = d.Set("tag_value_set", tagValues.TagValueSet)
	}

	d.SetId(helper.DataResourceIdsHash([]string{tagValueSet}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tagValueSet); e != nil {
			return e
		}
	}
	return nil
}
