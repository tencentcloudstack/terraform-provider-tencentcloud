package tencentcloud

import (
	"context"
	"log"

	cos "github.com/tencentyun/cos-go-sdk-v5"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type CosService struct {
	client *connectivity.TencentCloudClient
}

func (me *CosService) GetObjectList(ctx context.Context, bucketName, keyPrefix string) (result *cos.BucketGetResult, errRet error) {
	logId := GetLogId(ctx)

	option := &cos.BucketGetOptions{
		Prefix: keyPrefix,
	}
	result, response, err := me.client.UseCosClient(bucketName).Bucket.Get(ctx, option)
	if err != nil {
		errRet = err
		log.Printf("[CRITAL]%s api[%s] fail, bucket name[%s], key prefix[%s], reason[%s]\n",
			logId, "get object list", bucketName, keyPrefix, errRet.Error())
		return
	}

	requestId := response.Header.Get("x-cos-request-id")
	log.Printf("[DEBUG]%s api[%s] success, bucket name[%s], key[%s], return status code[%d], request id[%s]\n",
		logId, "get object list", bucketName, keyPrefix, response.StatusCode, requestId)

	return
}
