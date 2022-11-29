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

type DtsService struct {
	client *connectivity.TencentCloudClient
}

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
	request.JobId = &jobId

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

func (me *DtsService) PollingStatusUntil(ctx context.Context, jobId string, targetStatus string) error {
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

		log.Printf("[DEBUG]%s api[%s] has been destroyed, return with nil\n", logId, "DescribeDtsSyncJob")
		return nil
		// return resource.RetryableError(fmt.Errorf("DTS sync job[%s] is nil, retry...", jobId))
	})
	if err != nil {
		return err
	}
	return nil
}

func (me *DtsService) DeleteDtsSyncJobById(ctx context.Context, jobId string) (errRet error) {
	var (
		logId    = getLogId(ctx)
		request  = dts.NewDestroySyncJobRequest()
		response = dts.NewDestroyMigrateJobResponse()
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
		request.JobId = &jobId
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

func (me *DtsService) DescribeDtsCompareTask(ctx context.Context, jobId, compareTaskId string) (compareTask *dts.DescribeCompareReportResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dts.NewDescribeCompareReportRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.JobId = &jobId
	request.CompareTaskId = &compareTaskId

	response, err := me.client.UseDtsClient().DescribeCompareReport(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	compareTask = response.Response
	return
}

func (me *DtsService) DeleteDtsCompareTaskById(ctx context.Context, jobId, compareTaskId string) (errRet error) {
	logId := getLogId(ctx)

	request := dts.NewDeleteCompareTaskRequest()

	request.JobId = &jobId
	request.CompareTaskId = &compareTaskId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDtsClient().DeleteCompareTask(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
