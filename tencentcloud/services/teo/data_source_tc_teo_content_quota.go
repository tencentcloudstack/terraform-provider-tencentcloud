package teo

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoContentQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoContentQuotaRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"purge_quota": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cache purge quota list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Quota type. Valid values: `purge_prefix`, `purge_url`, `purge_host`, `purge_all`, `purge_cache_tag`.",
						},
						"batch": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Single batch submission quota limit.",
						},
						"daily": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Daily submission quota limit.",
						},
						"daily_available": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Daily remaining available quota.",
						},
					},
				},
			},

			"prefetch_quota": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cache prefetch quota list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Quota type. Valid values: `prefetch_url`.",
						},
						"batch": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Single batch submission quota limit.",
						},
						"daily": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Daily submission quota limit.",
						},
						"daily_available": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Daily remaining available quota.",
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

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["zone_id"] = v.(string)
	}

	var purgeQuotaList, prefetchQuotaList []map[string]interface{}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		purgeData, prefetchData, e := service.DescribeTeoContentQuotaByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if purgeData != nil {
			for _, item := range purgeData {
				quotaMap := map[string]interface{}{}
				if item.Type != nil {
					quotaMap["type"] = item.Type
				}
				if item.Batch != nil {
					quotaMap["batch"] = item.Batch
				}
				if item.Daily != nil {
					quotaMap["daily"] = item.Daily
				}
				if item.DailyAvailable != nil {
					quotaMap["daily_available"] = item.DailyAvailable
				}
				purgeQuotaList = append(purgeQuotaList, quotaMap)
			}
		}
		if prefetchData != nil {
			for _, item := range prefetchData {
				quotaMap := map[string]interface{}{}
				if item.Type != nil {
					quotaMap["type"] = item.Type
				}
				if item.Batch != nil {
					quotaMap["batch"] = item.Batch
				}
				if item.Daily != nil {
					quotaMap["daily"] = item.Daily
				}
				if item.DailyAvailable != nil {
					quotaMap["daily_available"] = item.DailyAvailable
				}
				prefetchQuotaList = append(prefetchQuotaList, quotaMap)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read tencentcloud_teo_content_quota failed, reason:%s\n", logId, err.Error())
		return err
	}

	_ = d.Set("purge_quota", purgeQuotaList)
	_ = d.Set("prefetch_quota", prefetchQuotaList)

	ids := make([]string, 0)
	if zoneId, ok := paramMap["zone_id"]; ok {
		ids = append(ids, zoneId.(string))
	}
	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		allData := map[string]interface{}{
			"purge_quota":    purgeQuotaList,
			"prefetch_quota": prefetchQuotaList,
		}
		if e := tccommon.WriteToFile(output.(string), allData); e != nil {
			return e
		}
	}

	return nil
}
