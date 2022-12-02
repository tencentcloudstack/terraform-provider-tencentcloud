package tencentcloud

import (
	"context"
	"log"

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
