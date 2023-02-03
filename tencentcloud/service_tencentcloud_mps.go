package tencentcloud

import (
	"context"
	"log"

	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type MpsService struct {
	client *connectivity.TencentCloudClient
}

func (me *MpsService) DescribeMpsWorkflowById(ctx context.Context, workflowId string) (workflow *mps.WorkflowInfo, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeWorkflowsRequest()
	request.WorkflowIds = []*int64{helper.Int64(helper.StrToInt64(workflowId))}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeWorkflows(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.WorkflowInfoSet) < 1 {
		return
	}

	workflow = response.Response.WorkflowInfoSet[0]
	return
}

func (me *MpsService) DeleteMpsWorkflowById(ctx context.Context, workflowId string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeleteWorkflowRequest()
	request.WorkflowId = helper.Int64(helper.StrToInt64(workflowId))

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeleteWorkflow(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) EnableWorkflow(ctx context.Context, workflowId int64) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewEnableWorkflowRequest()
	request.WorkflowId = helper.Int64(workflowId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().EnableWorkflow(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DisableWorkflow(ctx context.Context, workflowId int64) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDisableWorkflowRequest()
	request.WorkflowId = helper.Int64(workflowId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DisableWorkflow(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
