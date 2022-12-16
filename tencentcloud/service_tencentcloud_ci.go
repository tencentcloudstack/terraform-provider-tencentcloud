package tencentcloud

import (
	"context"
	"log"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	ci "github.com/tencentyun/cos-go-sdk-v5"
)

type CiService struct {
	client *connectivity.TencentCloudClient
}

func (me *CiService) DescribeCiBucketById(ctx context.Context, bucket string) (serviceResult *ci.CIServiceResult, errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "GetCIService", bucket, errRet.Error())
		}
	}()

	result, response, err := me.client.UseCiClient(bucket).CI.GetCIService(ctx)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s], http status [%s]\n", logId, "GetCIService", bucket, result, response.Status)

	serviceResult = result
	return
}

func (me *CiService) DeleteCiBucketById(ctx context.Context, bucket string) (errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "CloseCIService", bucket, errRet.Error())
		}
	}()

	response, err := me.client.UseCiClient(bucket).CI.CloseCIService(ctx)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response status [%s]\n", logId, "CloseCIService", bucket, response.Status)

	return
}
