package waf

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	wafv20180125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudWafIpAccessControlV2ReadPostHandleResponse0(ctx context.Context, resp *wafv20180125.DescribeIpAccessControlResponseParams) error {
	logId := tccommon.GetLogId(ctx)
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	if resp.Data == nil || len(resp.Data.Res) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `waf_ip_access_control_v2` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	res := resp.Data.Res[0]
	if res.IpList != nil {
		tmpList := make([]string, 0, len(res.IpList))
		for _, v := range res.IpList {
			tmpList = append(tmpList, *v)
		}

		_ = d.Set("ip_list", tmpList)
	}

	if res.ActionType != nil {
		_ = d.Set("action_type", res.ActionType)
	}

	if res.Note != nil {
		_ = d.Set("note", res.Note)
	}

	if res.JobType != nil {
		_ = d.Set("job_type", res.JobType)
	}

	if res.JobDateTime != nil {
		jobDateTimeMap := map[string]interface{}{}
		if res.JobDateTime.Timed != nil {
			tmpList := make([]map[string]interface{}, 0, len(res.JobDateTime.Timed))
			for _, timed := range res.JobDateTime.Timed {
				timedMap := map[string]interface{}{}
				if timed.StartDateTime != nil {
					timedMap["start_date_time"] = timed.StartDateTime
				}

				if timed.EndDateTime != nil {
					timedMap["end_date_time"] = timed.EndDateTime
				}

				tmpList = append(tmpList, timedMap)
			}

			jobDateTimeMap["timed"] = tmpList
		}

		if res.JobDateTime.Cron != nil {
			tmpList := make([]map[string]interface{}, 0, len(res.JobDateTime.Cron))
			for _, cron := range res.JobDateTime.Cron {
				cronMap := map[string]interface{}{}
				if cron.Days != nil {
					cronMap["days"] = cron.Days
				}

				if cron.WDays != nil {
					cronMap["w_days"] = cron.WDays
				}

				if cron.StartTime != nil {
					cronMap["start_time"] = cron.StartTime
				}

				if cron.EndTime != nil {
					cronMap["end_time"] = cron.EndTime
				}

				tmpList = append(tmpList, cronMap)
			}

			jobDateTimeMap["cron"] = tmpList
		}

		if res.JobDateTime.TimeTZone != nil {
			jobDateTimeMap["time_t_zone"] = res.JobDateTime.TimeTZone
		}

		_ = d.Set("job_date_time", []interface{}{jobDateTimeMap})
	}

	return nil
}

func resourceTencentCloudWafIpAccessControlV2UpdatePreRequest0(ctx context.Context, req *wafv20180125.ModifyIpAccessControlRequest) *resource.RetryError {
	d := tccommon.ResourceDataFromContext(ctx)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return resource.NonRetryableError(fmt.Errorf("id is broken,%s", d.Id()))
	}

	instanceId := idSplit[0]
	domain := idSplit[1]
	ruleId := idSplit[2]

	req.InstanceId = &instanceId
	req.Domain = &domain
	ruleIdUint64, _ := strconv.ParseUint(ruleId, 10, 64)
	req.RuleId = &ruleIdUint64

	if v, ok := d.GetOkExists("action_type"); ok {
		req.ActionType = helper.IntInt64(v.(int))
	}

	sourceType := "custom"
	req.SourceType = &sourceType

	return nil
}

func resourceTencentCloudWafIpAccessControlV2DeletePreRequest0(ctx context.Context, req *wafv20180125.DeleteIpAccessControlV2Request) *resource.RetryError {
	d := tccommon.ResourceDataFromContext(ctx)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return resource.NonRetryableError(fmt.Errorf("id is broken,%s", d.Id()))
	}

	ruleId := idSplit[2]
	ruleIdUint64, _ := strconv.ParseUint(ruleId, 10, 64)
	req.RuleIds = []*uint64{&ruleIdUint64}

	return nil
}
