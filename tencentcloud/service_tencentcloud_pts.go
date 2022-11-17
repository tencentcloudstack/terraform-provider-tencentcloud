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

func (me *PtsService) DescribePtsAlertChannel(ctx context.Context, noticeId, projectId string) (alertChannel *pts.AlertChannelRecord, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = pts.NewDescribeAlertChannelsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.NoticeIds = []*string{&noticeId}
	request.ProjectIds = []*string{&projectId}

	response, err := me.client.UsePtsClient().DescribeAlertChannels(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.AlertChannelSet) < 1 {
		return
	}
	alertChannel = response.Response.AlertChannelSet[0]
	return
}

func (me *PtsService) DeletePtsAlertChannelById(ctx context.Context, noticeId, projectId string) (errRet error) {
	logId := getLogId(ctx)

	request := pts.NewDeleteAlertChannelRequest()

	request.NoticeId = &noticeId
	request.ProjectId = &projectId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UsePtsClient().DeleteAlertChannel(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *PtsService) DescribePtsScenario(ctx context.Context, projectId, scenarioId string) (scenario *pts.Scenario, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = pts.NewDescribeScenariosRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ProjectIds = []*string{&projectId}
	request.ScenarioIds = []*string{&scenarioId}

	response, err := me.client.UsePtsClient().DescribeScenarios(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.ScenarioSet) < 1 {
		return
	}
	scenario = response.Response.ScenarioSet[0]
	return
}

func (me *PtsService) DeletePtsScenarioById(ctx context.Context, projectId, scenarioId string) (errRet error) {
	logId := getLogId(ctx)

	request := pts.NewDeleteScenariosRequest()

	request.ProjectId = &projectId
	request.ScenarioIds = []*string{&scenarioId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UsePtsClient().DeleteScenarios(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *PtsService) DescribePtsFile(ctx context.Context, projectId, fileIds string) (file *pts.File, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = pts.NewDescribeFilesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ProjectIds = []*string{&projectId}
	request.FileIds = []*string{&fileIds}

	response, err := me.client.UsePtsClient().DescribeFiles(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.FileSet) < 1 {
		return
	}
	file = response.Response.FileSet[0]
	return
}

func (me *PtsService) DeletePtsFileById(ctx context.Context, projectId, fileIds string) (errRet error) {
	logId := getLogId(ctx)

	request := pts.NewDeleteFilesRequest()

	request.ProjectId = &projectId
	request.FileIds = []*string{&fileIds}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UsePtsClient().DeleteFiles(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *PtsService) DescribePtsJob(ctx context.Context, projectId, scenarioId, jobId string) (job *pts.Job, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = pts.NewDescribeJobsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ProjectIds = []*string{&projectId}
	request.ScenarioIds = []*string{&scenarioId}
	request.JobIds = []*string{&jobId}

	response, err := me.client.UsePtsClient().DescribeJobs(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.JobSet) < 1 {
		return
	}
	job = response.Response.JobSet[0]
	return
}

func (me *PtsService) DeletePtsJobById(ctx context.Context, projectId, scenarioId, jobId string) (errRet error) {
	logId := getLogId(ctx)

	request := pts.NewDeleteJobsRequest()

	request.ProjectId = &projectId
	request.ScenarioIds = []*string{&scenarioId}
	request.JobIds = []*string{&jobId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UsePtsClient().DeleteJobs(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *PtsService) DescribePtsCronJob(ctx context.Context, cronJobId, projectId string) (cronJob *pts.CronJob, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = pts.NewDescribeCronJobsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.CronJobIds = []*string{&cronJobId}
	request.ProjectIds = []*string{&projectId}

	response, err := me.client.UsePtsClient().DescribeCronJobs(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.CronJobSet) < 1 {
		return
	}
	cronJob = response.Response.CronJobSet[0]
	return
}

func (me *PtsService) DeletePtsCronJobById(ctx context.Context, cronJobId, projectId string) (errRet error) {
	logId := getLogId(ctx)

	request := pts.NewDeleteCronJobsRequest()

	request.CronJobIds = []*string{&cronJobId}
	request.ProjectId = &projectId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UsePtsClient().DeleteCronJobs(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
