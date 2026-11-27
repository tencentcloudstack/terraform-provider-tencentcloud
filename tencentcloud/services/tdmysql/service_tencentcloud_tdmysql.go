package tdmysql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewTdmysqlService(client *connectivity.TencentCloudClient) TdmysqlService {
	return TdmysqlService{client: client}
}

type TdmysqlService struct {
	client *connectivity.TencentCloudClient
}

func (me *TdmysqlService) DescribeTdmysqlDbInstanceById(ctx context.Context, instanceId string) (ret *tdmysqlv20211122.DescribeDBInstanceDetailResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tdmysqlv20211122.NewDescribeDBInstanceDetailRequest()
	response := tdmysqlv20211122.NewDescribeDBInstanceDetailResponse()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DescribeDBInstanceDetail(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tdmysql db instance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response
	return
}

func (me *TdmysqlService) DescribeTdmysqlDBInstanceSSLStatusById(ctx context.Context, instanceId string) (ret *tdmysqlv20211122.DescribeInstanceSSLStatusResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tdmysqlv20211122.NewDescribeInstanceSSLStatusRequest()
	response := tdmysqlv20211122.NewDescribeInstanceSSLStatusResponse()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DescribeInstanceSSLStatus(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tdmysql db instance ssl status failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response
	return
}

func (me *TdmysqlService) DescribeTdmysqlDBInstanceDataReservedSpaceById(ctx context.Context, instanceId string) (ret *tdmysqlv20211122.DescribeInstanceDataReservedSpaceResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tdmysqlv20211122.NewDescribeInstanceDataReservedSpaceRequest()
	response := tdmysqlv20211122.NewDescribeInstanceDataReservedSpaceResponse()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DescribeInstanceDataReservedSpace(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tdmysql db instance data reserved space failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response
	return
}

func (me *TdmysqlService) DescribeSaleInfo(ctx context.Context, request *tdmysqlv20211122.DescribeSaleInfoRequest) (regionList []*tdmysqlv20211122.DescribeSaleRegionInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	response := tdmysqlv20211122.NewDescribeSaleInfoResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DescribeSaleInfo(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tdmysql sale info failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	regionList = response.Response.RegionList
	return
}

func (me *TdmysqlService) DescribeTdmysqlDBSecurityGroupsById(ctx context.Context, instanceId string) (ret *tdmysqlv20211122.DescribeDBSecurityGroupsResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tdmysqlv20211122.NewDescribeDBSecurityGroupsRequest()
	response := tdmysqlv20211122.NewDescribeDBSecurityGroupsResponse()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DescribeDBSecurityGroups(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tdmysql db security groups failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response
	return
}

func (me *TdmysqlService) ModifyTdmysqlUserPrivileges(ctx context.Context, request *tdmysqlv20211122.ModifyUserPrivilegesRequest) (err error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), err.Error())
		}
	}()

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().ModifyUserPrivileges(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify tdmysql user privileges failed, Response is nil."))
		}

		return nil
	})

	return
}

func (me *TdmysqlService) DescribeTdmysqlUserPrivileges(ctx context.Context, request *tdmysqlv20211122.DescribeUserPrivilegesRequest) (ret *tdmysqlv20211122.DescribeUserPrivilegesResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DescribeUserPrivileges(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tdmysql user privileges failed, Response is nil."))
		}

		ret = result.Response
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	// Singleton describe shape: an answer without any privilege item means
	// "confirmed absent" — report the shared sentinel so the caller can
	// remove the resource from state (never a per-resource errors.New).
	if len(ret.Privileges) < 1 {
		errRet = fmt.Errorf("describe tdmysql user privileges: %w", fwHelper.ErrNotFound)

		return
	}

	return
}

// CreateTdmysqlUsers wraps CreateUsers: it performs the retried write call and
// returns the created resource id (instance_id#user_name#host, parent-first
// join of the request addressing fields) together with the async flow id the
// wait helper below polls. The flow id field name and Go type come from the
// vendored SDK response struct (FlowId *int64).
func (me *TdmysqlService) CreateTdmysqlUsers(ctx context.Context, request *tdmysqlv20211122.CreateUsersRequest) (usersId string, flowId int64, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *tdmysqlv20211122.CreateUsersResponse
	)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().CreateUsers(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create tdmysql users failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmysql users failed, reason:%+v", logId, err)
		return
	}

	if request.InstanceId == nil || len(request.Users) < 1 || request.Users[0].UserName == nil || request.Users[0].Host == nil {
		err = fmt.Errorf("api %s returns no user addressing fields", request.GetAction())
		return
	}

	usersId = helper.IdFormat(*request.InstanceId, *request.Users[0].UserName, *request.Users[0].Host)
	if response.Response.FlowId != nil {
		flowId = *response.Response.FlowId
	}

	return
}

// DescribeTdmysqlUsers wraps DescribeUsers: the caller builds the filter
// request (instance id); this method returns the FULL user list. The API has
// no pagination fields, so a single retried call is kept.
func (me *TdmysqlService) DescribeTdmysqlUsers(ctx context.Context, request *tdmysqlv20211122.DescribeUsersRequest) (users []*tdmysqlv20211122.UserInfo, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *tdmysqlv20211122.DescribeUsersResponse
	)

	users = make([]*tdmysqlv20211122.UserInfo, 0)

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DescribeUsers(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("describe tdmysql users failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s describe tdmysql users failed, reason:%+v", logId, err)
		return nil, err
	}

	users = append(users, response.Response.Users...)

	return
}

// DeleteTdmysqlUsers wraps DeleteUsers: the caller builds the delete request;
// this method only performs the retried call and returns the async flow id
// (FlowId *int64 in the vendored SDK response struct) the wait helper polls.
func (me *TdmysqlService) DeleteTdmysqlUsers(ctx context.Context, request *tdmysqlv20211122.DeleteUsersRequest) (flowId int64, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *tdmysqlv20211122.DeleteUsersResponse
	)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DeleteUsers(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("delete tdmysql users failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete tdmysql users failed, reason:%+v", logId, err)
		return
	}

	if response.Response.FlowId != nil {
		flowId = *response.Response.FlowId
	}

	return
}

// ResetTdmysqlUsersPassword wraps ResetUsersPassword: the caller builds the
// reset request; this method only performs the retried call and returns the
// async flow id (FlowId *int64 in the vendored SDK response struct) the wait
// helper polls. EncryptedPassword is intentionally not handled — the provider
// always sends the plain-text Password.
func (me *TdmysqlService) ResetTdmysqlUsersPassword(ctx context.Context, request *tdmysqlv20211122.ResetUsersPasswordRequest) (flowId int64, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *tdmysqlv20211122.ResetUsersPasswordResponse
	)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().ResetUsersPassword(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("reset tdmysql users password failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s reset tdmysql users password failed, reason:%+v", logId, err)
		return
	}

	if response.Response.FlowId != nil {
		flowId = *response.Response.FlowId
	}

	return
}

// WaitTdmysqlUsersResult is the async FLOW wait for tdmysql user tasks: it
// polls DescribeFlow by the flow id returned by CreateUsers / DeleteUsers
// until the flow reaches a terminal state. The terminal status values follow
// the shipped SDKv2 reference (FLOW_STATUS_SUCCESS / FAILED / PAUSED).
func (me *TdmysqlService) WaitTdmysqlUsersResult(ctx context.Context, flowId int64, timeout time.Duration) (err error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tdmysqlv20211122.NewDescribeFlowRequest()
	)

	request.FlowId = helper.Int64(flowId)
	err = resource.Retry(timeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseTdmysqlV20211122Client().DescribeFlow(request)
		if e != nil {
			// Error codes get ZERO second-guessing: retry classification is
			// tccommon.RetryError's job.
			return tccommon.RetryError(e)
		}

		if response == nil || response.Response == nil || response.Response.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeFlow task %d returns empty status", flowId))
		}

		switch *response.Response.Status {
		case FLOW_STATUS_SUCCESS:
			return nil
		case FLOW_STATUS_FAILED, FLOW_STATUS_PAUSED:
			return resource.NonRetryableError(fmt.Errorf("tdmysql user flow %d failed, status is %s", flowId, *response.Response.Status))
		default:
			return resource.RetryableError(fmt.Errorf("tdmysql user flow %d status is %s", flowId, *response.Response.Status))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, flow[%d], reason[%s]\n", logId, request.GetAction(), flowId, err.Error())
	}

	return
}

func (me *TdmysqlService) DescribeDBEngines(ctx context.Context, request *tdmysqlv20211122.DescribeDBEnginesRequest) (ret []*tdmysqlv20211122.DBEngineInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DescribeDBEngines(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tdmysql db engines failed, Response is nil."))
		}

		ret = result.Response.Items
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *TdmysqlService) DescribeSpecs(ctx context.Context, request *tdmysqlv20211122.DescribeSpecsRequest) (hybridNodeSpecs []*tdmysqlv20211122.StorageNodeSpec, serverlessCcuSpec []*tdmysqlv20211122.ServerlessCcu, serverlessNodeNumSpec *tdmysqlv20211122.ServerlessNodeNumSpec, errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DescribeSpecs(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tdmysql specs failed, Response is nil."))
		}

		hybridNodeSpecs = result.Response.HybridNodeSpecs
		serverlessCcuSpec = result.Response.ServerlessCcuSpec
		serverlessNodeNumSpec = result.Response.ServerlessNodeNumSpec
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}
