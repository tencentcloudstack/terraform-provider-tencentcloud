package sts

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	sts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type StsService struct {
	client *connectivity.TencentCloudClient
}

func (me *StsService) DescribeStsCallerIdentityByFilter(ctx context.Context) (callerIdentity *sts.GetCallerIdentityResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = sts.NewGetCallerIdentityRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseStsClient().GetCallerIdentity(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	callerIdentity = response.Response
	return
}
