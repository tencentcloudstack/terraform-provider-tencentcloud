package cdwpg

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	cdwpg "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwpg/v20201230"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CdwpgService struct {
	client *connectivity.TencentCloudClient
}

func (me *CdwpgService) DescribeCdwpgInstanceById(ctx context.Context, instanceId string) (instance *cdwpg.SimpleInstanceInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cdwpg.NewDescribeInstanceInfoRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwpgClient().DescribeInstanceInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	instance = response.Response.SimpleInstanceInfo
	return
}

func (me *CdwpgService) DeleteCdwpgInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cdwpg.NewDestroyInstanceByApiRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwpgClient().DestroyInstanceByApi(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CdwpgService) InstanceStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		request := cdwpg.NewDescribeInstanceStateRequest()
		request.InstanceId = &instanceId
		ratelimit.Check(request.GetAction())
		object, err := me.client.UseCdwpgClient().DescribeInstanceState(request)

		if err != nil {
			return nil, "", err
		}
		if object == nil || object.Response == nil || object.Response.InstanceState == nil {
			return nil, "", nil
		}

		return object, *object.Response.InstanceState, nil
	}
}

func NewCdwpgService(client *connectivity.TencentCloudClient) CdwpgService {
	return CdwpgService{client: client}
}

func (me *CdwpgService) DescribeCdwpgInstancesByFilter(ctx context.Context, param map[string]interface{}) (instances []*cdwpg.InstanceSimpleInfoNew, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cdwpg.NewDescribeSimpleInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "SearchInstanceId" {
			request.SearchInstanceId = v.(*string)
		}
		if k == "SearchInstanceName" {
			request.SearchInstanceName = v.(*string)
		}
		if k == "SearchTags" {
			request.SearchTags = v.([]*string)
		}
	}

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	instances = make([]*cdwpg.InstanceSimpleInfoNew, 0)
	for {
		request.Offset = helper.Int64(offset)
		request.Limit = helper.Int64(limit)
		ratelimit.Check(request.GetAction())
		var response *cdwpg.DescribeSimpleInstancesResponse
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseCdwpgV20201230Client().DescribeSimpleInstances(request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			response = result
			return nil
		})
		if err != nil {
			errRet = err
			return
		}
		if response != nil && response.Response != nil {
			if len(response.Response.InstancesList) > 0 {
				instances = append(instances, response.Response.InstancesList...)
			}
			if len(response.Response.InstancesList) < int(limit) {
				break
			}
			offset += limit
		}

	}
	return
}

func (me *CdwpgService) DescribeCdwpgLogByFilter(ctx context.Context, param map[string]interface{}) (ret *cdwpg.DescribeSlowLogResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cdwpg.NewDescribeSlowLogRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "Database" {
			request.Database = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "OrderByType" {
			request.OrderByType = v.(*string)
		}
		if k == "Duration" {
			request.Duration = v.(*float64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwpgV20201230Client().DescribeSlowLog(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	ret = response.Response
	return
}

func (me *CdwpgService) DescribeCdwpgLogByFilter1(ctx context.Context, param map[string]interface{}) (ret []*cdwpg.ErrorLogDetail, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cdwpg.NewDescribeErrorLogRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdwpgV20201230Client().DescribeErrorLog(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ErrorLogDetails) < 1 {
			break
		}
		ret = append(ret, response.Response.ErrorLogDetails...)
		if len(response.Response.ErrorLogDetails) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdwpgService) DescribeCdwpgNodesByFilter(ctx context.Context, param map[string]interface{}) (ret *cdwpg.DescribeInstanceNodesResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cdwpg.NewDescribeInstanceNodesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwpgV20201230Client().DescribeInstanceNodes(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	ret = response.Response
	return
}

func (me *CdwpgService) DescribeCdwpgDbconfigByFilter(ctx context.Context, param map[string]interface{}) (ret []*cdwpg.ParamItem, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cdwpg.NewDescribeDBParamsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "NodeTypes" {
			request.NodeTypes = v.([]*string)
		}
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdwpgV20201230Client().DescribeDBParams(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		ret = append(ret, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdwpgService) DescribeCdwpgAccountById(ctx context.Context, instanceId string) (ret *cdwpg.AccountInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cdwpg.NewDescribeAccountsRequest()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwpgV20201230Client().DescribeAccounts(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Accounts) < 1 {
		return
	}

	ret = response.Response.Accounts[0]
	return
}

func (me *CdwpgService) DescribeCdwpgUserhbaById(ctx context.Context, instanceId string) (ret *cdwpg.DescribeUserHbaConfigResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cdwpg.NewDescribeUserHbaConfigRequest()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwpgV20201230Client().DescribeUserHbaConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}
