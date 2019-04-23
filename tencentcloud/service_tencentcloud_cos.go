package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type CosService struct {
	client *connectivity.TencentCloudClient
}

func (me *CosService) HeadObject(ctx context.Context, bucket, key, versionId string) (info *s3.HeadObjectOutput, errRet error) {
	logId := GetLogId(ctx)

	request := s3.HeadObjectInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(key),
		VersionId: aws.String(versionId),
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "head object", request.String(), errRet.Error())
		}
	}()
	response, err := me.client.UseCosClient().HeadObject(&request)
	if err != nil {
		errRet = fmt.Errorf("cos head object error: %s, bucket: %s, object: %s", err.Error(), bucket, key)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "head object", request.String(), response.String())

	return response, nil
}

func (me *CosService) GetObjectTags(ctx context.Context, bucket, key string) (tags map[string]string, errRet error) {
	logId := GetLogId(ctx)

	request := s3.GetObjectTaggingInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "get object tags", request.String(), errRet.Error())
		}
	}()
	response, err := me.client.UseCosClient().GetObjectTagging(&request)
	if err != nil {
		errRet = fmt.Errorf("cos get object tags error: %s, bucket: %s, object: %s", err.Error(), bucket, key)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get object tags", request.String(), response.String())

	tags = make(map[string]string, len(response.TagSet))
	for _, t := range response.TagSet {
		tags[*t.Key] = *t.Value
	}

	return
}
