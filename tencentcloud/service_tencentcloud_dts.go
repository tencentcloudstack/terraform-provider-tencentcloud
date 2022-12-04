package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

const (
	DTSJobStatus   = "Status"
	DTSTradeStatus = "TradeStatus"
)

type DtsService struct {
	client *connectivity.TencentCloudClient
}

// sync job
func (me *DtsService) DescribeDtsSyncJob(ctx context.Context, jobId *string) (jobInfo *dts.SyncJobInfo, errRet error) {
	logId := getLogId(ctx)
	params := map[string]interface{}{}

	if jobId != nil {
		params["job_id"] = jobId
	}

	ret, err := me.DescribeDtsSyncJobsByFilter(ctx, params)
	if err != nil {
		errRet = err
		return
	}
	if len(ret) == 0 {
		log.Printf("[CRITAL]%s DescribeDtsSyncJob fail, reason[%s]\n", logId, "the result DescribeDtsSyncJobsByFilter is nil!")
		errRet = err
		return
	}

	jobInfo = ret[0]
	return
}

func (me *DtsService) DescribeDtsSyncJobsByFilter(ctx context.Context, param map[string]interface{}) (syncJobs []*dts.SyncJobInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dts.NewDescribeSyncJobsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "job_id" {
			request.JobId = v.(*string)
		}

		if k == "job_name" {
			request.JobName = v.(*string)
		}

		if k == "order" {
			request.Order = v.(*string)
		}

		if k == "order_seq" {
			request.OrderSeq = v.(*string)
		}

		if k == "status" {
			request.Status = helper.Strings(v.([]string))
		}

		if k == "run_mode" {
			request.RunMode = v.(*string)
		}

		if k == "job_type" {
			request.JobType = v.(*string)
		}

		if k == "pay_mode" {
			request.PayMode = v.(*string)
		}

		if k == "tag_filters" {
			request.TagFilters = v.([]*dts.TagFilter)
		}

	}
	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 20

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDtsClient().DescribeSyncJobs(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.JobList) < 1 {
			break
		}
		syncJobs = append(syncJobs, response.Response.JobList...)
		if len(response.Response.JobList) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *DtsService) IsolateDtsSyncJobById(ctx context.Context, jobId string) (errRet error) {
	logId := getLogId(ctx)

	request := dts.NewIsolateSyncJobRequest()
	request.JobId = helper.String(jobId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "isolate object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDtsClient().IsolateSyncJob(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DtsService) DeleteDtsSyncJobById(ctx context.Context, jobId string) (errRet error) {
	var (
		logId    = getLogId(ctx)
		request  = dts.NewDestroySyncJobRequest()
		response = dts.NewDestroySyncJobResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	err := me.IsolateDtsSyncJobById(ctx, jobId)
	if err != nil {
		errRet = err
		return
	}

	// err = me.PollingStatusUntil(ctx, jobId, "Isolated")
	// if err != nil {
	// 	return err
	// }

	ratelimit.Check(request.GetAction())
	err = resource.Retry(3*writeRetryTimeout, func() *resource.RetryError {
		request.JobId = helper.String(jobId)
		_, err := me.client.UseDtsClient().DestroySyncJob(request)
		if err != nil {
			time.Sleep(10 * time.Second)
			return resource.RetryableError(fmt.Errorf("destroy failed, retry... %v", err))
		}
		return nil
	})
	if err != nil {
		return err
	}

	// err = me.PollingStatusUntil(ctx, jobId, "Deleted")
	// if err != nil {
	// 	return err
	// }
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DtsService) PollingSyncJobStatusUntil(ctx context.Context, jobId string, targetStatus string) error {
	logId := getLogId(ctx)

	err := resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		ret, err := me.DescribeDtsSyncJob(ctx, helper.String(jobId))
		if err != nil {
			return retryError(err)
		}

		if ret != nil && ret.Status != nil {
			status := *ret.Status
			if strings.Contains(targetStatus, status) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("DTS sync job[%s] status is still on [%s], retry...", jobId, status))
		}

		log.Printf("[DEBUG]%s sync job[%s] doesn't exist, exit retry...\n", logId, jobId)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// compare task
func (me *DtsService) DescribeDtsCompareTasksByFilter(ctx context.Context, param map[string]interface{}) (compareTasks []*dts.CompareTaskItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dts.NewDescribeCompareTasksRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query objects", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "job_id" {
			request.JobId = v.(*string)
		}
	}
	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 20

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDtsClient().DescribeCompareTasks(request)
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
		compareTasks = append(compareTasks, response.Response.Items...)
		if len(response.Response.Items) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *DtsService) DescribeDtsCompareTask(ctx context.Context, jobId, compareTaskId *string) (tasks []*dts.CompareTaskItem, errRet error) {
	logId := getLogId(ctx)
	param := map[string]interface{}{
		"job_id": jobId,
	}

	ret, err := me.DescribeDtsCompareTasksByFilter(ctx, param)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, "DescribeDtsCompareTask", err.Error())
		errRet = err
		return
	}
	if ret == nil {
		errRet = fmt.Errorf("DescribeDtsCompareTask failed, ret is nil. jobId:[%s], compareTaskId:[%s]", *jobId, *compareTaskId)
		return
	}

	tasks = ret

	// exactly search for specific compare task
	if compareTaskId != nil {
		for _, t := range tasks {
			if *compareTaskId == *t.CompareTaskId {
				tasks = []*dts.CompareTaskItem{t}
				break
			}
		}
	}

	log.Printf("[DEBUG]%s api[%s] success, tasks num:[%v]\n", logId, "DescribeDtsCompareTask", len(tasks))
	return
}

func (me *DtsService) StopDtsCompareById(ctx context.Context, jobId, compareTaskId string) (errRet error) {
	logId := getLogId(ctx)

	request := dts.NewStopCompareRequest()
	request.JobId = helper.String(jobId)
	request.CompareTaskId = helper.String(compareTaskId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "isolate object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDtsClient().StopCompare(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DtsService) DeleteDtsCompareTaskById(ctx context.Context, jobId, compareTaskId string) (errRet error) {
	var (
		logId    = getLogId(ctx)
		request  = dts.NewDeleteCompareTaskRequest()
		response = dts.NewDeleteCompareTaskResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	// ***canceled task cannot be deleted!!!***
	// err := me.StopDtsCompareById(ctx, jobId, compareTaskId)
	// if err != nil {
	// 	errRet = err
	// 	return
	// }

	// wait success or failed
	err := me.PollingCompareTaskStatusUntil(ctx, jobId, compareTaskId, "success,failed")
	if err != nil {
		return err
	}

	request.JobId = helper.String(jobId)
	request.CompareTaskId = helper.String(compareTaskId)

	ratelimit.Check(request.GetAction())
	response, err = me.client.UseDtsClient().DeleteCompareTask(request)
	if err != nil {
		errRet = err
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DtsService) PollingCompareTaskStatusUntil(ctx context.Context, jobId, compareTaskId, targetStatus string) error {
	logId := getLogId(ctx)

	err := resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		ret, err := me.DescribeDtsCompareTask(ctx, helper.String(jobId), helper.String(compareTaskId))
		if err != nil {
			return retryError(err)
		}

		if ret != nil && ret[0].Status != nil {
			status := *ret[0].Status
			if strings.Contains(targetStatus, status) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("DTS compare task [%s,%s] status is still on [%s], retry...", jobId, compareTaskId, status))
		}

		log.Printf("[DEBUG]%s compare task[%s,%s] doesn't exist, exit retry...\n", logId, jobId, compareTaskId)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// migration job
func (me *DtsService) DescribeDtsMigrateJobsByFilter(ctx context.Context, param map[string]interface{}) (migrateJobs []*dts.JobItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dts.NewDescribeMigrationJobsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "job_id" {
			request.JobId = v.(*string)
		}

		if k == "job_name" {
			request.JobName = v.(*string)
		}

		if k == "status" {
			request.Status = v.([]*string)
		}

		if k == "src_instance_id" {
			request.SrcInstanceId = v.(*string)
		}

		if k == "src_region" {
			request.SrcRegion = v.(*string)
		}

		if k == "src_database_type" {
			request.SrcDatabaseType = v.([]*string)
		}

		if k == "src_access_type" {
			request.SrcAccessType = v.([]*string)
		}

		if k == "dst_instance_id" {
			request.DstInstanceId = v.(*string)
		}

		if k == "dst_region" {
			request.DstRegion = v.(*string)
		}

		if k == "dst_database_type" {
			request.DstDatabaseType = v.([]*string)
		}

		if k == "dst_access_type" {
			request.DstAccessType = v.([]*string)
		}

		if k == "run_mode" {
			request.RunMode = v.(*string)
		}

		if k == "order_seq" {
			request.OrderSeq = v.(*string)
		}

		if k == "tag_filters" {
			request.TagFilters = v.([]*dts.TagFilter)
		}

	}
	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 20

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDtsClient().DescribeMigrationJobs(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.JobList) < 1 {
			break
		}
		migrateJobs = append(migrateJobs, response.Response.JobList...)
		if len(response.Response.JobList) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *DtsService) DescribeDtsMigrateJob(ctx context.Context, jobId string) (migrateJob *dts.DescribeMigrationDetailResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dts.NewDescribeMigrationDetailRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.JobId = helper.String(jobId)

	response, err := me.client.UseDtsClient().DescribeMigrationDetail(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	migrateJob = response.Response
	return
}

func (me *DtsService) IsolateDtsMigrateJobById(ctx context.Context, jobId string) (errRet error) {
	logId := getLogId(ctx)

	request := dts.NewIsolateMigrateJobRequest()
	request.JobId = helper.String(jobId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "isolate object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDtsClient().IsolateMigrateJob(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DtsService) DeleteDtsMigrateJobById(ctx context.Context, jobId string) (errRet error) {
	var (
		logId    = getLogId(ctx)
		request  = dts.NewDestroyMigrateJobRequest()
		response = dts.NewDestroyMigrateJobResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	err := me.IsolateDtsMigrateJobById(ctx, jobId)
	if err != nil {
		errRet = err
		return
	}

	err = me.PollingMigrateJobStatusUntil(ctx, jobId, DTSTradeStatus, "isolated")
	if err != nil {
		return err
	}

	ratelimit.Check(request.GetAction())
	err = resource.Retry(3*writeRetryTimeout, func() *resource.RetryError {
		request.JobId = helper.String(jobId)
		_, err := me.client.UseDtsClient().DestroyMigrateJob(request)
		if err != nil {
			time.Sleep(10 * time.Second)
			return resource.RetryableError(fmt.Errorf("destroy failed, retry... %v", err))
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = me.PollingMigrateJobStatusUntil(ctx, jobId, DTSTradeStatus, "offlined")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DtsService) PollingMigrateJobStatusUntil(ctx context.Context, jobId, statusType, targetStatus string) error {
	logId := getLogId(ctx)

	err := resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		ret, err := me.DescribeDtsMigrateJob(ctx, jobId)
		if err != nil {
			return retryError(err)
		}

		if statusType == DTSJobStatus {
			if ret != nil && ret.Status != nil {
				status := *ret.Status
				if strings.Contains(targetStatus, status) {
					return nil
				}
				return resource.RetryableError(fmt.Errorf("DTS migrate job[%s] Status is still on [%s], retry...", jobId, status))
			}
		}
		if statusType == DTSTradeStatus {
			if ret != nil && ret.TradeInfo.TradeStatus != nil {
				status := *ret.TradeInfo.TradeStatus
				if strings.Contains(targetStatus, status) {
					return nil
				}
				return resource.RetryableError(fmt.Errorf("DTS migrate job[%s] TradeStatus is still on [%s], retry...", jobId, status))
			}
		}

		log.Printf("[DEBUG]%s migrate job[%s] doesn't exist, exit retry...\n", logId, jobId)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
