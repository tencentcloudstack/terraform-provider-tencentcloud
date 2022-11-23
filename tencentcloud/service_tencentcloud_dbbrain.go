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
	ret, errRet := me.DescribeDbbrainSqlFilters(ctx, instanceId, []*int64{helper.StrToInt64Point(*filterId)})
	if errRet != nil {
		return
	}
	if ret != nil {
		return ret[0], nil
	}
	return
}

func (me *DbbrainService) DescribeDbbrainSqlFilters(ctx context.Context, instanceId *string, filterIds []*int64) (sqlFilters []*dbbrain.SQLFilter, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dbbrain.NewDescribeSqlFiltersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = instanceId
	if filterIds != nil {
		request.FilterIds = filterIds
	}

	response, err := me.client.UseDbbrainClient().DescribeSqlFilters(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.Items) < 1 {
		return
	}
	sqlFilters = response.Response.Items
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
	ret, errRet := me.DescribeDbbrainSecurityAuditLogExportTasks(ctx, secAuditGroupId, []*string{asyncRequestId}, product)
	if errRet != nil {
		return
	}
	if ret != nil {
		return ret.Tasks[0], nil
	}
	return
}

func (me *DbbrainService) DescribeDbbrainSecurityAuditLogExportTasks(ctx context.Context, secAuditGroupId *string, asyncRequestId []*string, product *string) (params *dbbrain.DescribeSecurityAuditLogExportTasksResponseParams, errRet error) {
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

	request.SecAuditGroupId = secAuditGroupId

	if asyncRequestId != nil {
		request.AsyncRequestIds = helper.StringsToUint64Pointer(asyncRequestId)
	}

	if product != nil {
		request.Product = product
	} else {
		request.Product = helper.String("mysql")
	}

	response, err := me.client.UseDbbrainClient().DescribeSecurityAuditLogExportTasks(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.Tasks) < 1 {
		return
	}
	params = response.Response
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
