package common

import (
	"log"
	"strconv"
	"time"

	billing "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

const (
	BillingStartMonth = "2018-05"
)

// CreatorAccountInfo 创建者账号信息
type CreatorAccountInfo struct {
	UserId   string // 用户ID
	UserName string // 用户名
}

// CloudDescribeBillResourceSummary get billing resource summary data
func CloudDescribeBillResourceSummary(client *connectivity.TencentCloudClient, resources []*ResourceInstance) map[string]*CreatorAccountInfo {
	resourceIdToSubAccountInfoMap := make(map[string]*CreatorAccountInfo)

	request := billing.NewDescribeBillResourceSummaryRequest()
	request.Offset = helper.Uint64(0)
	request.Limit = helper.Uint64(1)
	request.NeedRecordNum = helper.Int64(1)

	currentMonth := GetCurrentMonth()
	for _, r := range resources {
		if r.Id == "" {
			continue
		}

		request.ResourceId = helper.String(r.Id)
		prevMonth := currentMonth
		for {
			if prevMonth == BillingStartMonth {
				break
			}

			request.Month = helper.String(prevMonth)
			response, err := client.UseBillingClient().DescribeBillResourceSummary(request)
			if err != nil {
				log.Printf("[CRITAL] get billing resource[%v] summary data error: %v", r.Id, err.Error())
				break
			}
			if response == nil || response.Response == nil {
				log.Printf("[CRITAL] get billing resource[%v] summary data response is nil", r.Id)
				break
			}
			if *response.Response.Total == 1 {
				billResourceSummary := response.Response.ResourceSummarySet[0]
				userName := CloudDescribeSubAccounts(client, *billResourceSummary.OperateUin)

				resourceIdToSubAccountInfoMap[r.Id] = &CreatorAccountInfo{
					UserId:   *billResourceSummary.OperateUin,
					UserName: userName,
				}
				break
			}

			prevMonth = PrevMonth(prevMonth)
		}
	}
	return resourceIdToSubAccountInfoMap
}

// CloudDescribeSubAccounts get sub account data
func CloudDescribeSubAccounts(client *connectivity.TencentCloudClient, uin string) string {
	uinNum, err := strconv.ParseUint(uin, 10, 64)
	if err != nil {
		log.Printf("[CRITAL] parse uin[%v] to uint64 type error: %v", uin, err.Error())
		return ""
	}

	request := cam.NewDescribeSubAccountsRequest()

	uinArray := []*uint64{helper.Uint64(uinNum)}
	request.FilterSubAccountUin = uinArray

	response, err := client.UseCamClient().DescribeSubAccounts(request)
	if err != nil {
		log.Printf("[CRITAL] get sub account[%v] data error: %v", uin, err.Error())
		return ""
	}
	if response == nil || response.Response == nil {
		log.Printf("[CRITAL] get sub account[%v] data response is nil", uin)
		return ""
	}

	name := response.Response.SubAccounts[0].Name
	return *name
}

// GetCurrentMonth get current month
func GetCurrentMonth() string {
	currentTime := time.Now()
	formattedMonth := currentTime.Format("2006-01")
	return formattedMonth
}

// PrevMonth get prev month
func PrevMonth(date string) string {
	t, _ := time.Parse("2006-01", date)
	t = t.AddDate(0, -1, 0)
	if t.Before(time.Date(2018, 5, 1, 0, 0, 0, 0, time.UTC)) {
		return BillingStartMonth
	}
	return t.Format("2006-01")
}
