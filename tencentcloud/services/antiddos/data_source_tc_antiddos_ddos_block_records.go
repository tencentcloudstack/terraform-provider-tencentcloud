package antiddos

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddosv20250903 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20250903"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAntiddosDDoSBlockRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAntiddosDDoSBlockRecordsRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Query start time. Parameter format: 2026-02-04T11:30:00+08:00.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Query end time. (EndTime - StartTime) must be <= 31 days. Parameter format: 2026-03-04T11:30:00+08:00.",
			},
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions. The upper limit of Filters.Values is 20. If not filled in, it returns the list of all blocked resources under the current appid.\n- Resource: filter by blocked IP or resource six-segment style.\n- Status: filter by block status (Blocked/Unblocking/Unblocked).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter field name. Supported: Resource, Status.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Filter field value.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"block_records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Block/unblock records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Blocked resource, public network IP.",
						},
						"block_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Block time.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Block/unblock status. Enum: Blocked, Unblocking, Unblocked.",
						},
					},
				},
			},

			"unblock_quota_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Unblock quota info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total unblock quota.",
						},
						"used_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Used unblock quota.",
						},
						"quota_start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Quota effective start time.",
						},
						"quota_end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Quota effective end time.",
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

func dataSourceTencentCloudAntiddosDDoSBlockRecordsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_antiddos_ddos_block_records.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = AntiddosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = v.(string)
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = v.(string)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*antiddosv20250903.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := antiddosv20250903.Filter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				filter.Name = helper.String(v)
			}

			if v, ok := filtersMap["values"]; ok {
				valueSet := v.(*schema.Set).List()
				for i := range valueSet {
					value := valueSet[i].(string)
					filter.Values = append(filter.Values, helper.String(value))
				}
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	var (
		respData         []*antiddosv20250903.DDoSBlockRecord
		unblockQuotaInfo *antiddosv20250903.DDoSUnblockQuota
	)
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		blockRecords, quotaInfo, e := service.DescribeAntiddosDDoSBlockRecordsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = blockRecords
		unblockQuotaInfo = quotaInfo
		return nil
	})

	if reqErr != nil {
		log.Printf("[DATASOURCE] read antiddos_ddos_block_records empty, skip SetId")
		return reqErr
	}

	blockRecordsList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, blockRecord := range respData {
			blockRecordMap := map[string]interface{}{}
			if blockRecord.Resource != nil {
				blockRecordMap["resource"] = blockRecord.Resource
			}

			if blockRecord.BlockTime != nil {
				blockRecordMap["block_time"] = blockRecord.BlockTime
			}

			if blockRecord.Status != nil {
				blockRecordMap["status"] = blockRecord.Status
			}

			blockRecordsList = append(blockRecordsList, blockRecordMap)
		}

		_ = d.Set("block_records", blockRecordsList)
	}

	if unblockQuotaInfo != nil {
		unblockQuotaInfoMap := map[string]interface{}{}
		if unblockQuotaInfo.TotalQuota != nil {
			unblockQuotaInfoMap["total_quota"] = unblockQuotaInfo.TotalQuota
		}

		if unblockQuotaInfo.UsedQuota != nil {
			unblockQuotaInfoMap["used_quota"] = unblockQuotaInfo.UsedQuota
		}

		if unblockQuotaInfo.QuotaStartTime != nil {
			unblockQuotaInfoMap["quota_start_time"] = unblockQuotaInfo.QuotaStartTime
		}

		if unblockQuotaInfo.QuotaEndTime != nil {
			unblockQuotaInfoMap["quota_end_time"] = unblockQuotaInfo.QuotaEndTime
		}

		_ = d.Set("unblock_quota_info", []interface{}{unblockQuotaInfoMap})
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
