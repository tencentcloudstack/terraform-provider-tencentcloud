package common

//import (
//	audit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
//	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
//	"time"
//)
//
//// CloudAuditDescribeEvents query operation audit log
//func CloudAuditDescribeEvents(client *connectivity.TencentCloudClient, resourceId, eventName string) {
//	startTime, endTime := GetTimestamp()
//
//	request := audit.NewDescribeEventsRequest()
//
//	client.UseAuditClient().DescribeEvents(request)
//}
//
//// GetTimestamp get the current timestamp and the timestamps of the past 90 days
//func GetTimestamp() (int64, int64) {
//	now := time.Now()
//	currentTimestamp := now.Unix()
//
//	past90Days := now.AddDate(0, 0, -90)
//	past90DaysTimestamp := past90Days.Unix()
//
//	return currentTimestamp, past90DaysTimestamp
//}
