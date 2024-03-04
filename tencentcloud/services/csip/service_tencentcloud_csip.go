package csip

import (
	"context"
	"log"

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
