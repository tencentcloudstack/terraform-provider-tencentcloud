package teo

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func DataSourceTencentCloudTeoContentQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoContentQuotaRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Zone ID.",
			},

			"purge_quota": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Purge quota list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"batch": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Upper limit of quota for each batch submission.",
						},
						"daily": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Upper limit of daily submission quota.",
						},
						"daily_available": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Remaining daily submission quota.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cache purge type. Values: `purge_prefix` (prefix purge), `purge_url` (URL purge), `purge_host` (hostname purge), `purge_all` (purge all cache), `purge_cache_tag` (cache tag purge), `prefetch_url` (URL prefetch).",
						},
					},
				},
			},

			"prefetch_quota": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Prefetch quota list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"batch": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Upper limit of quota for each batch submission.",
						},
						"daily": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Upper limit of daily submission quota.",
						},
						"daily_available": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Remaining daily submission quota.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cache purge type. Values: `purge_prefix` (prefix purge), `purge_url` (URL purge), `purge_host` (hostname purge), `purge_all` (purge all cache), `purge_cache_tag` (cache tag purge), `prefetch_url` (URL prefetch).",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTeoContentQuotaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_content_quota.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	_ = tccommon.GetLogId(tccommon.ContextNil)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient()

	request := teov20220901.NewDescribeContentQuotaRequest()
	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	ratelimit.Check(request.GetAction())

	var response *teov20220901.DescribeContentQuotaResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := client.DescribeContentQuota(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}

	if response == nil || response.Response == nil {
		return nil
	}

	if response.Response.PurgeQuota != nil {
		purgeQuotaList := make([]map[string]interface{}, 0, len(response.Response.PurgeQuota))
		for _, quota := range response.Response.PurgeQuota {
			quotaMap := map[string]interface{}{}
			if quota.Batch != nil {
				quotaMap["batch"] = *quota.Batch
			}
			if quota.Daily != nil {
				quotaMap["daily"] = *quota.Daily
			}
			if quota.DailyAvailable != nil {
				quotaMap["daily_available"] = *quota.DailyAvailable
			}
			if quota.Type != nil {
				quotaMap["type"] = *quota.Type
			}
			purgeQuotaList = append(purgeQuotaList, quotaMap)
		}
		_ = d.Set("purge_quota", purgeQuotaList)
	}

	if response.Response.PrefetchQuota != nil {
		prefetchQuotaList := make([]map[string]interface{}, 0, len(response.Response.PrefetchQuota))
		for _, quota := range response.Response.PrefetchQuota {
			quotaMap := map[string]interface{}{}
			if quota.Batch != nil {
				quotaMap["batch"] = *quota.Batch
			}
			if quota.Daily != nil {
				quotaMap["daily"] = *quota.Daily
			}
			if quota.DailyAvailable != nil {
				quotaMap["daily_available"] = *quota.DailyAvailable
			}
			if quota.Type != nil {
				quotaMap["type"] = *quota.Type
			}
			prefetchQuotaList = append(prefetchQuotaList, quotaMap)
		}
		_ = d.Set("prefetch_quota", prefetchQuotaList)
	}

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
