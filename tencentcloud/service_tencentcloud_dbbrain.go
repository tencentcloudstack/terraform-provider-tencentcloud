package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type DbbrainService struct {
	client *connectivity.TencentCloudClient
}

func (me *DbbrainService) DescribeDbbrainSqlFilter(ctx context.Context, instanceId, filterId *string) (sqlFilter *dbbrain.SQLFilter, errRet error) {
	param := make(map[string]interface{})
	if instanceId != nil {
		param["instance_id"] = instanceId
	}
	if filterId != nil {
		param["filter_ids"] = []*int64{helper.StrToInt64Point(*filterId)}
	}

	ret, errRet := me.DescribeDbbrainSqlFiltersByFilter(ctx, param)
	if errRet != nil {
		return
	}
	if ret != nil {
		return ret[0], nil
	}
	return
}

func (me *DbbrainService) DescribeDbbrainSqlFiltersByFilter(ctx context.Context, param map[string]interface{}) (sqlFilters []*dbbrain.SQLFilter, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dbbrain.NewDescribeSqlFiltersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query objects", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}

		if k == "filter_ids" {
			request.FilterIds = v.([]*int64)
		}

		if k == "statuses" {
			request.Statuses = v.([]*string)
		}
	}
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 20

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDbbrainClient().DescribeSqlFilters(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		sqlFilters = append(sqlFilters, response.Response.Items...)
		if len(response.Response.Items) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *DbbrainService) getSessionToken(ctx context.Context, instanceId, user, pw, product *string) (sessionToken *string, errRet error) {
	logId := getLogId(ctx)
	request := dbbrain.NewVerifyUserAccountRequest()

	request.InstanceId = instanceId
	request.User = user
	request.Password = pw
	if product != nil {
		request.Product = product
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "VerifyUserAccount", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDbbrainClient().VerifyUserAccount(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	sessionToken = response.Response.SessionToken
	return
}

func (me *DbbrainService) DeleteDbbrainSqlFilterById(ctx context.Context, instanceId, filterId, sessionToken *string) (errRet error) {
	logId := getLogId(ctx)

	request := dbbrain.NewDeleteSqlFiltersRequest()

	request.InstanceId = instanceId
	request.FilterIds = []*int64{helper.StrToInt64Point(*filterId)}
	request.SessionToken = sessionToken

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDbbrainClient().DeleteSqlFilters(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DbbrainService) DescribeDbbrainSecurityAuditLogExportTask(ctx context.Context, secAuditGroupId, asyncRequestId, product *string) (task *dbbrain.SecLogExportTaskInfo, errRet error) {
	param := make(map[string]interface{})
	if secAuditGroupId != nil {
		param["sec_audit_group_id"] = secAuditGroupId
	}
	if asyncRequestId != nil {
		param["async_request_ids"] = []*uint64{helper.StrToUint64Point(*asyncRequestId)}
	}
	if product != nil {
		param["product"] = product
	} else {
		param["product"] = helper.String("mysql")
	}

	ret, errRet := me.DescribeDbbrainSecurityAuditLogExportTasksByFilter(ctx, param)
	if errRet != nil {
		return
	}
	if ret != nil {
		return ret[0], nil
	}
	return
}

func (me *DbbrainService) DescribeDbbrainSecurityAuditLogExportTasksByFilter(ctx context.Context, param map[string]interface{}) (securityAuditLogExportTasks []*dbbrain.SecLogExportTaskInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dbbrain.NewDescribeSecurityAuditLogExportTasksRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query objects", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "sec_audit_group_id" {
			request.SecAuditGroupId = v.(*string)
		}

		if k == "product" {
			request.Product = v.(*string)
		}

		if k == "async_request_ids" {
			request.AsyncRequestIds = v.([]*uint64)
		}
	}
	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 20

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDbbrainClient().DescribeSecurityAuditLogExportTasks(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Tasks) < 1 {
			break
		}
		securityAuditLogExportTasks = append(securityAuditLogExportTasks, response.Response.Tasks...)
		if len(response.Response.Tasks) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *DbbrainService) DeleteDbbrainSecurityAuditLogExportTaskById(ctx context.Context, secAuditGroupId, asyncRequestId, product *string) (errRet error) {
	logId := getLogId(ctx)

	request := dbbrain.NewDeleteSecurityAuditLogExportTasksRequest()

	request.SecAuditGroupId = secAuditGroupId
	request.AsyncRequestIds = []*uint64{helper.StrToUint64Point(*asyncRequestId)}
	if product != nil {
		request.Product = product
	} else {
		request.Product = helper.String("mysql")
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDbbrainClient().DeleteSecurityAuditLogExportTasks(request)
	if err != nil {
		errRet = err
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DbbrainService) DescribeDbbrainDiagEventsByFilter(ctx context.Context, param map[string]interface{}) (diagEvents []*dbbrain.DiagHistoryEventItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dbbrain.NewDescribeDBDiagEventsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_ids" {
			request.InstanceIds = v.([]*string)
		}
		if k == "start_time" {
			request.StartTime = v.(*string)
		}
		if k == "end_time" {
			request.EndTime = v.(*string)
		}
		if k == "severities" {
			request.Severities = v.([]*int64)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseDbbrainClient().DescribeDBDiagEvents(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		diagEvents = append(diagEvents, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *DbbrainService) DescribeDbbrainDiagEventByFilter(ctx context.Context, param map[string]interface{}) (diagEvent *dbbrain.DescribeDBDiagEventResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dbbrain.NewDescribeDBDiagEventRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "event_id" {
			request.EventId = v.(*int64)
		}
		if k == "product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeDBDiagEvent(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil {
		diagEvent = response.Response
	}

	return
}

func (me *DbbrainService) DescribeDbbrainDiagHistoryByFilter(ctx context.Context, param map[string]interface{}) (diagHistory []*dbbrain.DiagHistoryEventItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dbbrain.NewDescribeDBDiagHistoryRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "start_time" {
			request.StartTime = v.(*string)
		}
		if k == "end_time" {
			request.EndTime = v.(*string)
		}
		if k == "product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeDBDiagHistory(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil {
		diagHistory = response.Response.Events
	}

	return
}

func (me *DbbrainService) DescribeDbbrainSecurityAuditLogDownloadUrlsByFilter(ctx context.Context, param map[string]interface{}) (securityAuditLogDownloadUrls []*string, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dbbrain.NewDescribeSecurityAuditLogDownloadUrlsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "sec_audit_group_id" {
			request.SecAuditGroupId = v.(*string)
		}
		if k == "async_request_id" {
			request.AsyncRequestId = v.(*uint64)
		}
		if k == "product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeSecurityAuditLogDownloadUrls(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil {
		securityAuditLogDownloadUrls = response.Response.Urls
	}

	return
}

func (me *DbbrainService) DescribeDbbrainSlowLogTimeSeriesStatsByFilter(ctx context.Context, param map[string]interface{}) (slowLogTimeSeriesStats *dbbrain.DescribeSlowLogTimeSeriesStatsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dbbrain.NewDescribeSlowLogTimeSeriesStatsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "start_time" {
			request.StartTime = v.(*string)
		}
		if k == "end_time" {
			request.EndTime = v.(*string)
		}
		if k == "product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeSlowLogTimeSeriesStats(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil {
		slowLogTimeSeriesStats = response.Response
	}

	return
}

func (me *DbbrainService) DescribeDbbrainSlowLogTopSqlsByFilter(ctx context.Context, param map[string]interface{}) (slowLogTopSqls []*dbbrain.SlowLogTopSqlItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dbbrain.NewDescribeSlowLogTopSqlsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "start_time" {
			request.StartTime = v.(*string)
		}
		if k == "end_time" {
			request.EndTime = v.(*string)
		}
		if k == "sort_by" {
			request.SortBy = v.(*string)
		}
		if k == "order_by" {
			request.OrderBy = v.(*string)
		}
		if k == "limit" {
			request.Limit = v.(*int64)
		}
		if k == "offset" {
			request.Offset = v.(*int64)
		}
		if k == "schema_list" {
			request.SchemaList = v.([]*dbbrain.SchemaItem)
		}
		if k == "product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseDbbrainClient().DescribeSlowLogTopSqls(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Rows) < 1 {
			break
		}
		slowLogTopSqls = append(slowLogTopSqls, response.Response.Rows...)
		if len(response.Response.Rows) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *DbbrainService) DescribeDbbrainSlowLogUserHostStatsByFilter(ctx context.Context, param map[string]interface{}) (slowLogUserHostStats []*dbbrain.SlowLogHost, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dbbrain.NewDescribeSlowLogUserHostStatsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "start_time" {
			request.StartTime = v.(*string)
		}
		if k == "end_time" {
			request.EndTime = v.(*string)
		}
		if k == "product" {
			request.Product = v.(*string)
		}
		if k == "md5" {
			request.Md5 = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeSlowLogUserHostStats(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil {
		slowLogUserHostStats = response.Response.Items
	}

	return
}

func (me *DbbrainService) DescribeDbbrainSlowLogUserSqlAdviceByFilter(ctx context.Context, param map[string]interface{}) (slowLogUserSqlAdvice *dbbrain.DescribeUserSqlAdviceResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dbbrain.NewDescribeUserSqlAdviceRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "sql_text" {
			request.SqlText = v.(*string)
		}
		if k == "schema" {
			request.Schema = v.(*string)
		}
		if k == "product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeUserSqlAdvice(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil {
		slowLogUserSqlAdvice = response.Response
	}

	return
}

func (me *DbbrainService) DescribeDbbrainDbDiagReportTaskById(ctx context.Context, asyncRequestId *int64, instanceId string, product string) (dbDiagReportTask *dbbrain.HealthReportTask, errRet error) {
	logId := getLogId(ctx)

	request := dbbrain.NewDescribeDBDiagReportTasksRequest()
	request.InstanceIds = []*string{helper.String(instanceId)}
	request.Product = &product

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeDBDiagReportTasks(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if asyncRequestId != nil {
		for _, task := range response.Response.Tasks {
			if *task.AsyncRequestId == *asyncRequestId {
				dbDiagReportTask = task
				return
			}
		}
		return nil, fmt.Errorf("[ERROR]%sThe asyncRequestId[%v] not found in the qurey results. \n", logId, *asyncRequestId)
	}

	dbDiagReportTask = response.Response.Tasks[0]
	return
}

func (me *DbbrainService) DeleteDbbrainDbDiagReportTaskById(ctx context.Context, asyncRequestId int64, instanceId string, product string) (errRet error) {
	logId := getLogId(ctx)

	request := dbbrain.NewDeleteDBDiagReportTasksRequest()
	request.AsyncRequestIds = []*int64{helper.Int64(asyncRequestId)}
	request.InstanceId = &instanceId
	request.Product = &product

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DeleteDBDiagReportTasks(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DbbrainService) DbbrainDbDiagReportTaskStateRefreshFunc(asyncRequestId *int64, instanceId string, product string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeDbbrainDbDiagReportTaskById(ctx, asyncRequestId, instanceId, product)

		if err != nil {
			return nil, "", err
		}

		return object, helper.Int64ToStr(*object.Progress), nil
	}
}
