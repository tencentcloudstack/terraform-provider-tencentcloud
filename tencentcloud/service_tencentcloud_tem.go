package tencentcloud

import (
	"context"
	"log"

	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TemService struct {
	client *connectivity.TencentCloudClient
}

func (me *TemService) DescribeTemApplication(ctx context.Context, applicationId string) (application *tem.DescribeApplicationsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tem.NewDescribeApplicationsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ApplicationId = &applicationId

	response, err := me.client.UseTemClient().DescribeApplications(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	application = response.Response
	return
}

func (me *TemService) DeleteTemApplicationById(ctx context.Context, applicationId string) (errRet error) {
	logId := getLogId(ctx)

	request := tem.NewDeleteApplicationRequest()
	request.ApplicationId = &applicationId
	request.EnvironmentId = helper.String("")

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTemClient().DeleteApplication(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TemService) DescribeTemWorkload(ctx context.Context, environmentId string, applicationId string) (workload *tem.DescribeApplicationInfoResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tem.NewDescribeApplicationInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId

	response, err := me.client.UseTemClient().DescribeApplicationInfo(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	workload = response.Response
	return
}

func (me *TemService) DeleteTemWorkloadById(ctx context.Context, environmentId string, applicationId string) (errRet error) {
	logId := getLogId(ctx)

	request := tem.NewDeleteApplicationRequest()
	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTemClient().DeleteApplication(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
