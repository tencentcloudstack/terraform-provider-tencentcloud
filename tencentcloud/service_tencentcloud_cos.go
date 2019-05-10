package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type CosService struct {
	client *connectivity.TencentCloudClient
}

func (me *CosService) HeadObject(ctx context.Context, bucket, key string) (info *s3.HeadObjectOutput, errRet error) {
	logId := GetLogId(ctx)

	request := s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
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

func (me *CosService) DeleteObject(ctx context.Context, bucket, key string) (errRet error) {
	logId := GetLogId(ctx)

	request := s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.String(), errRet.Error())
		}
	}()
	response, err := me.client.UseCosClient().DeleteObject(&request)
	if err != nil {
		errRet = fmt.Errorf("cos delete object error: %s, bucket: %s, object: %s", err.Error(), bucket, key)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "delete object", request.String(), response.String())

	return nil
}

func (me *CosService) PutObjectAcl(ctx context.Context, bucket, key, acl string) (errRet error) {
	logId := GetLogId(ctx)

	request := s3.PutObjectAclInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		ACL:    aws.String(acl),
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "put object acl", request.String(), errRet.Error())
		}
	}()
	response, err := me.client.UseCosClient().PutObjectAcl(&request)
	if err != nil {
		errRet = fmt.Errorf("cos put object acl error: %s, bucket: %s, object: %s", err.Error(), bucket, key)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "put object acl", request.String(), response.String())

	return nil
}

func (me *CosService) PutBucket(ctx context.Context, bucket, acl string) (errRet error) {
	logId := GetLogId(ctx)

	request := s3.CreateBucketInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String(acl),
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "put bucket", request.String(), errRet.Error())
		}
	}()
	response, err := me.client.UseCosClient().CreateBucket(&request)
	if err != nil {
		errRet = fmt.Errorf("cos put bucket error: %s, bucket: %s", err.Error(), bucket)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "put bucket", request.String(), response.String())

	return nil
}

func (me *CosService) HeadBucket(ctx context.Context, bucket string) (errRet error) {
	logId := GetLogId(ctx)

	request := s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}
	response, err := me.client.UseCosClient().HeadBucket(&request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "head bucket", request.String(), err.Error())
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "head bucket", request.String(), response.String())

	return nil
}

func (me *CosService) DeleteBucket(ctx context.Context, bucket string) (errRet error) {
	logId := GetLogId(ctx)

	request := s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	}
	response, err := me.client.UseCosClient().DeleteBucket(&request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "delete bucket", request.String(), err.Error())
		return fmt.Errorf("cos delete bucket error: %s, bucket: %s", err.Error(), bucket)
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "delete bucket", request.String(), response.String())

	return nil
}

func (me *CosService) GetBucketCors(ctx context.Context, bucket string) (corsRules []map[string]interface{}, errRet error) {
	logId := GetLogId(ctx)

	request := s3.GetBucketCorsInput{
		Bucket: aws.String(bucket),
	}
	response, err := me.client.UseCosClient().GetBucketCors(&request)
	if err != nil {
		awsError, ok := err.(awserr.Error)
		if !ok || awsError.Code() != "NoSuchCORSConfiguration" {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "get bucket cors", request.String(), err.Error())
			errRet = fmt.Errorf("cos get bucket cors error: %s, bucket: %s", err.Error(), bucket)
			return
		}
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get bucket cors", request.String(), response.String())

	corsRules = make([]map[string]interface{}, 0, len(response.CORSRules))
	if len(response.CORSRules) > 0 {
		for _, value := range response.CORSRules {
			rule := make(map[string]interface{})
			rule["allowed_origins"] = flattenStringList(value.AllowedOrigins)
			rule["allowed_methods"] = flattenStringList(value.AllowedMethods)
			rule["allowed_headers"] = flattenStringList(value.AllowedHeaders)

			if value.ExposeHeaders != nil {
				rule["expose_headers"] = flattenStringList(value.ExposeHeaders)
			}
			if value.MaxAgeSeconds != nil {
				rule["max_age_seconds"] = int(*value.MaxAgeSeconds)
			}

			corsRules = append(corsRules, rule)
		}
	}
	return
}

func (me *CosService) GetBucketLifecycle(ctx context.Context, bucket string) (lifecycleRules []map[string]interface{}, errRet error) {
	logId := GetLogId(ctx)

	request := s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
	}
	response, err := me.client.UseCosClient().GetBucketLifecycleConfiguration(&request)
	if err != nil {
		awsError, ok := err.(awserr.Error)
		if !ok || awsError.Code() != "NoSuchLifecycleConfiguration" {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "get bucket lifecycle", request.String(), err.Error())
			errRet = fmt.Errorf("cos get bucket cors error: %s, bucket: %s", err.Error(), bucket)
			return
		}
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get bucket lifecycle", request.String(), response.String())

	lifecycleRules = make([]map[string]interface{}, 0, len(response.Rules))
	if len(response.Rules) > 0 {
		for _, value := range response.Rules {
			rule := make(map[string]interface{})

			// filter_prefix
			if value.Filter != nil {
				if value.Filter.And != nil && value.Filter.And.Prefix != nil &&
					*value.Filter.And.Prefix != "" {
					rule["filter_prefix"] = *value.Filter.And.Prefix
				} else if value.Filter.Prefix != nil && *value.Filter.Prefix != "" {
					rule["filter_prefix"] = *value.Filter.Prefix
				}
			}
			// transition
			if len(value.Transitions) > 0 {
				transitions := make([]interface{}, 0, len(value.Transitions))
				for _, v := range value.Transitions {
					t := make(map[string]interface{})
					if v.Date != nil {
						t["date"] = (*v.Date).Format("2006-01-02")
					}
					if v.Days != nil {
						t["days"] = int(*v.Days)
					}
					if v.StorageClass != nil {
						t["storage_class"] = *v.StorageClass
					}
					transitions = append(transitions, t)
				}
				rule["transition"] = schema.NewSet(transitionHash, transitions)
			}
			// expiration
			if value.Expiration != nil {
				e := make(map[string]interface{})
				if value.Expiration.Date != nil {
					e["date"] = (*value.Expiration.Date).Format("2006-01-02")
				}
				if value.Expiration.Days != nil {
					e["days"] = int(*value.Expiration.Days)
				}
				rule["expiration"] = schema.NewSet(expirationHash, []interface{}{e})
			}

			lifecycleRules = append(lifecycleRules, rule)
		}
	}
	return
}

func (me *CosService) GetDataSourceBucketLifecycle(ctx context.Context, bucket string) (lifecycleRules []map[string]interface{}, errRet error) {
	logId := GetLogId(ctx)

	request := s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
	}
	response, err := me.client.UseCosClient().GetBucketLifecycleConfiguration(&request)
	if err != nil {
		awsError, ok := err.(awserr.Error)
		if !ok || awsError.Code() != "NoSuchLifecycleConfiguration" {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "get bucket lifecycle", request.String(), err.Error())
			errRet = fmt.Errorf("cos get bucket cors error: %s, bucket: %s", err.Error(), bucket)
			return
		}
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get bucket lifecycle", request.String(), response.String())

	lifecycleRules = make([]map[string]interface{}, 0, len(response.Rules))
	if len(response.Rules) > 0 {
		for _, value := range response.Rules {
			rule := make(map[string]interface{})

			// filter_prefix
			if value.Filter != nil {
				if value.Filter.And != nil && value.Filter.And.Prefix != nil &&
					*value.Filter.And.Prefix != "" {
					rule["filter_prefix"] = *value.Filter.And.Prefix
				} else if value.Filter.Prefix != nil && *value.Filter.Prefix != "" {
					rule["filter_prefix"] = *value.Filter.Prefix
				}
			}
			// transition
			if len(value.Transitions) > 0 {
				transitions := make([]interface{}, 0, len(value.Transitions))
				for _, v := range value.Transitions {
					t := make(map[string]interface{})
					if v.Date != nil {
						t["date"] = (*v.Date).Format("2006-01-02")
					}
					if v.Days != nil {
						t["days"] = int(*v.Days)
					}
					if v.StorageClass != nil {
						t["storage_class"] = *v.StorageClass
					}
					transitions = append(transitions, t)
				}
				rule["transition"] = transitions
			}
			// expiration
			if value.Expiration != nil {
				e := make(map[string]interface{})
				if value.Expiration.Date != nil {
					e["date"] = (*value.Expiration.Date).Format("2006-01-02")
				}
				if value.Expiration.Days != nil {
					e["days"] = int(*value.Expiration.Days)
				}
				rule["expiration"] = []interface{}{e}
			}

			lifecycleRules = append(lifecycleRules, rule)
		}
	}
	return
}

func (me *CosService) GetBucketWebsite(ctx context.Context, bucket string) (websites []map[string]interface{}, errRet error) {
	logId := GetLogId(ctx)

	request := s3.GetBucketWebsiteInput{
		Bucket: aws.String(bucket),
	}
	response, err := me.client.UseCosClient().GetBucketWebsite(&request)
	if err != nil {
		awsError, ok := err.(awserr.Error)
		if !ok || awsError.Code() != "NoSuchWebsiteConfiguration" {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "get bucket website", request.String(), err.Error())
			errRet = fmt.Errorf("cos get bucket website error: %s, bucket: %s", err.Error(), bucket)
			return
		}
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get bucket website", request.String(), response.String())

	websites = make([]map[string]interface{}, 0, 1)
	website := make(map[string]interface{})
	if response.IndexDocument != nil {
		website["index_document"] = *response.IndexDocument.Suffix
	}
	if response.ErrorDocument != nil {
		website["error_document"] = *response.ErrorDocument.Key
	}
	if len(website) > 0 {
		websites = append(websites, website)
	}

	return
}

func (me *CosService) ListBuckets(ctx context.Context) (buckets []*s3.Bucket, errRet error) {
	logId := GetLogId(ctx)

	request := s3.ListBucketsInput{}
	response, err := me.client.UseCosClient().ListBuckets(&request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "get bucket list", request.String(), err.Error())
		errRet = fmt.Errorf("cos get bucket list error: %s", err.Error())
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get bucket list", request.String(), response.String())

	buckets = response.Buckets
	return
}

func (me *CosService) ListObjects(ctx context.Context, bucket string) (objects []*s3.Object, errRet error) {
	logId := GetLogId(ctx)

	request := s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	}
	response, err := me.client.UseCosClient().ListObjects(&request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "get object list", request.String(), err.Error())
		errRet = fmt.Errorf("cos get object list error: %s", err.Error())
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get object list", request.String(), response.String())

	objects = response.Contents
	return
}
