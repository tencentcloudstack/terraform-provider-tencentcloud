package tencentcloud

import (
	"context"
	"log"

	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type PtsService struct {
	client *connectivity.TencentCloudClient
}

func (me *PtsService) DescribePtsProject(ctx context.Context, projectId string) (project *pts.Project, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = pts.NewDescribeProjectsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ProjectIds = []*string{&projectId}

	response, err := me.client.UsePtsClient().DescribeProjects(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.ProjectSet) < 1 {
		return
	}
	project = response.Response.ProjectSet[0]
	return
}

func (me *PtsService) DeletePtsProjectById(ctx context.Context, projectId string) (errRet error) {
	logId := getLogId(ctx)

	request := pts.NewDeleteProjectsRequest()

	request.ProjectIds = []*string{&projectId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UsePtsClient().DeleteProjects(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
