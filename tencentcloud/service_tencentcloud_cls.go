package tencentcloud

import (
	"context"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"log"
)

type ClsService struct {
	client *connectivity.TencentCloudClient
}

func (me *ClsService) DescribeClsLogSetById(ctx context.Context, logSetId string) (logset *cls.LogsetInfo, errRet error) {
	logId := getLogId(ctx)
	request := cls.NewDescribeLogsetsRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = append(
		request.Filters,
		&cls.Filter{
			Key: helper.String("logsetId"),
			Values: []*string{&logSetId},
		},
	)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClsClient().DescribeLogsets(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || len(response.Response.Logsets) == 0 {
		return
	}

	logset = response.Response.Logsets[0]

	return
}


func (me *ClsService) UpdateClsLogSet(ctx context.Context, request *cls.ModifyLogsetRequest) (errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClsClient().ModifyLogset(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ClsService) DeleteClsLogSet(ctx context.Context, id string) (errRet error) {
	logId := getLogId(ctx)

	request := cls.NewDeleteLogsetRequest()
	request.LogsetId = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClsClient().DeleteLogset(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}