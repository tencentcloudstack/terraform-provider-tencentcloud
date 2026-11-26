package csip

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	csip "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/csip/v20221121"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewCsipService(client *connectivity.TencentCloudClient) CsipService {
	return CsipService{client: client}
}

type CsipService struct {
	client *connectivity.TencentCloudClient
}

const csipCheckViewRisksPageSize = uint64(100)

// DescribeCheckViewRisks queries the cloud asset configuration risks under the
// check item view. The API pages (Offset/Limit), so this method carries the
// paging loop internally and returns the complete lists; the TotalCount is an
// internal paging detail and never part of the return value.
func (me *CsipService) DescribeCheckViewRisks(ctx context.Context, request *csip.DescribeCheckViewRisksRequest) (
	checkViewRiskList []*csip.CheckViewRiskItem,
	standardNameList []*csip.StandardItem,
	assetTypeList []*csip.AttributeOptionSet,
	providerList []*csip.AttributeOptionSet,
	errRet error,
) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset uint64 = 0
		limit  uint64 = csipCheckViewRisksPageSize
	)

	checkViewRiskList = make([]*csip.CheckViewRiskItem, 0)
	standardNameList = make([]*csip.StandardItem, 0)
	assetTypeList = make([]*csip.AttributeOptionSet, 0)
	providerList = make([]*csip.AttributeOptionSet, 0)

	for {
		request.Offset = &offset
		request.Limit = &limit

		var response *csip.DescribeCheckViewRisksResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseCsipClient().DescribeCheckViewRisks(request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("DescribeCheckViewRisks returned an empty response")
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		checkViewRiskList = append(checkViewRiskList, response.Response.CheckViewRiskList...)
		if len(response.Response.StandardNameList) > 0 {
			standardNameList = append(standardNameList, response.Response.StandardNameList...)
		}
		if len(response.Response.AssetTypeList) > 0 {
			assetTypeList = append(assetTypeList, response.Response.AssetTypeList...)
		}
		if len(response.Response.ProviderList) > 0 {
			providerList = append(providerList, response.Response.ProviderList...)
		}

		if len(response.Response.CheckViewRiskList) < 1 {
			break
		}

		if uint64(len(response.Response.CheckViewRiskList)) < limit {
			break
		}

		offset += limit
	}

	return
}

func (me *CsipService) DescribeCsipRiskCenterById(ctx context.Context, taskId string) (riskCenter *csip.ScanTaskInfoList, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := csip.NewDescribeScanTaskListRequest()
	request.Filter = &csip.Filter{
		Filters: []*csip.WhereFilter{
			{
				Name:         common.StringPtr("TaskId"),
				Values:       common.StringPtrs([]string{taskId}),
				OperatorType: common.Int64Ptr(1),
			},
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCsipClient().DescribeScanTaskList(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Data) != 1 {
		return
	}

	riskCenter = response.Response.Data[0]
	return
}

func (me *CsipService) StopCsipRiskCenterById(ctx context.Context, taskId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := csip.NewStopRiskCenterTaskRequest()
	request.TaskIdList = []*csip.TaskIdListKey{
		{
			TaskId: common.StringPtr(taskId),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCsipClient().StopRiskCenterTask(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CsipService) DeleteCsipRiskCenterById(ctx context.Context, taskId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := csip.NewDeleteRiskScanTaskRequest()
	request.TaskIdList = []*csip.TaskIdListKey{
		{
			TaskId: common.StringPtr(taskId),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCsipClient().DeleteRiskScanTask(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
