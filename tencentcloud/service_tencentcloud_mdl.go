package tencentcloud

import (
	"context"
	"log"

	mdl "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/mdl/v20200326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type MdlService struct {
	client *connectivity.TencentCloudClient
}

func (me *MdlService) DescribeMdlStreamLiveInputById(ctx context.Context, id string) (streamliveInput *mdl.InputInfo, errRet error) {
	logId := getLogId(ctx)

	request := mdl.NewDescribeStreamLiveInputRequest()
	request.Id = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMdlClient().DescribeStreamLiveInput(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Info != nil {
		streamliveInput = response.Response.Info
	}

	return
}

func (me *MdlService) DeleteMdlStreamLiveInputById(ctx context.Context, id string) (errRet error) {
	logId := getLogId(ctx)

	request := mdl.NewDeleteStreamLiveInputRequest()
	request.Id = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMdlClient().DeleteStreamLiveInput(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
