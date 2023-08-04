package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type SsmService struct {
	client *connectivity.TencentCloudClient
}

type SecretInfo struct {
	secretName       string
	description      string
	kmsKeyId         string
	secretType       int64
	additionalConfig string
	createUin        uint64
	status           string
	deleteTime       uint64
	createTime       uint64
	resourceId       string
}

type SecretVersionInfo struct {
	secretName   string
	versionId    string
	secretBinary string
	secretString string
}

func (me *SsmService) DescribeSecretsByFilter(ctx context.Context, param map[string]interface{}) (secrets []*ssm.SecretMetadata, errRet error) {
	logId := getLogId(ctx)
	request := ssm.NewListSecretsRequest()

	for k, v := range param {
		if k == "order_type" {
			request.OrderType = helper.Uint64(uint64(v.(int)))
		}
		if k == "state" {
			request.State = helper.Uint64(uint64(v.(int)))
		}
		if k == "secret_name" {
			request.SearchSecretName = helper.String(v.(string))
		}
		if k == "tag_filter" {
			tagFilter := v.(map[string]string)
			for tagKey, tagValue := range tagFilter {
				tag := ssm.TagFilter{
					TagKey:   helper.String(tagKey),
					TagValue: []*string{helper.String(tagValue)},
				}
				request.TagFilters = append(request.TagFilters, &tag)
			}
		}
	}
	var offset uint64 = 0
	var pageSize = uint64(SSM_PAGE_LIMIT)
	secrets = make([]*ssm.SecretMetadata, 0)
	for {
		request.Offset = helper.Uint64(offset)
		request.Limit = helper.Uint64(pageSize)
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseSsmClient().ListSecrets(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		if response == nil || len(response.Response.SecretMetadatas) < 1 {
			break
		}

		secrets = append(secrets, response.Response.SecretMetadatas...)

		if uint64(len(response.Response.SecretMetadatas)) < pageSize {
			break
		}
		offset += pageSize
	}
	return
}

func (me *SsmService) DescribeSecretByName(ctx context.Context, secretName string) (secret *SecretInfo, errRet error) {
	logId := getLogId(ctx)
	request := ssm.NewDescribeSecretRequest()
	request.SecretName = helper.String(secretName)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSsmClient().DescribeSecret(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	secret = &SecretInfo{
		secretName:       *response.Response.SecretName,
		description:      *response.Response.Description,
		kmsKeyId:         *response.Response.KmsKeyId,
		secretType:       *response.Response.SecretType,
		additionalConfig: *response.Response.AdditionalConfig,
		createUin:        *response.Response.CreateUin,
		status:           *response.Response.Status,
		deleteTime:       *response.Response.DeleteTime,
		createTime:       *response.Response.CreateTime,
		resourceId:       fmt.Sprintf("creatorUin/%d/%s", *response.Response.CreateUin, *response.Response.SecretName),
	}
	return
}

func (me *SsmService) DescribeSecretVersionIdsByName(ctx context.Context, secretName string) (versionIds []string, errRet error) {
	logId := getLogId(ctx)
	request := ssm.NewListSecretVersionIdsRequest()
	request.SecretName = helper.String(secretName)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSsmClient().ListSecretVersionIds(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	versionIds = make([]string, 0, len(response.Response.Versions))
	for _, versionInfo := range response.Response.Versions {
		versionIds = append(versionIds, *versionInfo.VersionId)
	}
	return
}

func (me *SsmService) DescribeSecretVersion(ctx context.Context, secretName, versionId string) (secretVersion *SecretVersionInfo, errRet error) {
	logId := getLogId(ctx)
	request := ssm.NewGetSecretValueRequest()
	request.SecretName = helper.String(secretName)
	request.VersionId = helper.String(versionId)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSsmClient().GetSecretValue(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	secretVersion = &SecretVersionInfo{
		secretName:   *response.Response.SecretName,
		versionId:    *response.Response.VersionId,
		secretBinary: *response.Response.SecretBinary,
		secretString: *response.Response.SecretString,
	}
	return
}

func (me *SsmService) CreateSecret(ctx context.Context, param map[string]interface{}) (secretName string, errRet error) {
	logId := getLogId(ctx)
	request := ssm.NewCreateSecretRequest()

	for k, v := range param {
		if k == "secret_name" {
			request.SecretName = helper.String(v.(string))
		}
		if k == "version_id" {
			request.VersionId = helper.String(v.(string))
		}
		if k == "description" {
			request.Description = helper.String(v.(string))
		}
		if k == "kms_key_id" {
			request.KmsKeyId = helper.String(v.(string))
		}
		if k == "secret_type" {
			request.SecretType = helper.IntUint64(v.(int))
		}
		if k == "additional_config" {
			request.AdditionalConfig = helper.String(v.(string))
		}
		if k == "secret_binary" {
			request.SecretBinary = helper.String(v.(string))
		}
		if k == "secret_string" {
			request.SecretString = helper.String(v.(string))
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSsmClient().CreateSecret(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	secretName = *response.Response.SecretName
	return
}

func (me *SsmService) UpdateSecretDescription(ctx context.Context, secretName, description string) (errRet error) {
	logId := getLogId(ctx)
	request := ssm.NewUpdateDescriptionRequest()
	request.SecretName = helper.String(secretName)
	request.Description = helper.String(description)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSsmClient().UpdateDescription(request)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SsmService) UpdateSecret(ctx context.Context, param map[string]interface{}) (errRet error) {
	logId := getLogId(ctx)
	request := ssm.NewUpdateSecretRequest()
	for k, v := range param {
		if k == "secret_name" {
			request.SecretName = helper.String(v.(string))
		}
		if k == "version_id" {
			request.VersionId = helper.String(v.(string))
		}
		if k == "secret_binary" {
			request.SecretBinary = helper.String(v.(string))
		}
		if k == "secret_string" {
			request.SecretString = helper.String(v.(string))
		}
	}
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSsmClient().UpdateSecret(request)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SsmService) EnableSecret(ctx context.Context, secretName string) (errRet error) {
	logId := getLogId(ctx)
	request := ssm.NewEnableSecretRequest()
	request.SecretName = helper.String(secretName)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSsmClient().EnableSecret(request)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SsmService) DisableSecret(ctx context.Context, secretName string) (errRet error) {
	logId := getLogId(ctx)
	request := ssm.NewDisableSecretRequest()
	request.SecretName = helper.String(secretName)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSsmClient().DisableSecret(request)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SsmService) PutSecretValue(ctx context.Context, param map[string]interface{}) (secretName, versionId string, errRet error) {
	logId := getLogId(ctx)
	request := ssm.NewPutSecretValueRequest()
	for k, v := range param {
		if k == "secret_name" {
			request.SecretName = helper.String(v.(string))
		}
		if k == "version_id" {
			request.VersionId = helper.String(v.(string))
		}
		if k == "secret_binary" {
			request.SecretBinary = helper.String(v.(string))
		}
		if k == "secret_string" {
			request.SecretString = helper.String(v.(string))
		}
	}
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSsmClient().PutSecretValue(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	secretName = *response.Response.SecretName
	versionId = *response.Response.VersionId
	return
}

func (me *SsmService) DeleteSecretVersion(ctx context.Context, secretName, versionId string) (errRet error) {
	logId := getLogId(ctx)
	request := ssm.NewDeleteSecretVersionRequest()
	request.SecretName = helper.String(secretName)
	request.VersionId = helper.String(versionId)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSsmClient().DeleteSecretVersion(request)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SsmService) DeleteSecret(ctx context.Context, secretName string, recoveryWindowInDays uint64) (errRet error) {
	logId := getLogId(ctx)
	request := ssm.NewDeleteSecretRequest()
	request.SecretName = helper.String(secretName)
	request.RecoveryWindowInDays = helper.Uint64(recoveryWindowInDays)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSsmClient().DeleteSecret(request)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SsmService) DescribeSecretById(ctx context.Context, secretName string, serviceType uint64) (sshKeyPairSecret *ssm.SecretMetadata, errRet error) {
	logId := getLogId(ctx)

	request := ssm.NewListSecretsRequest()
	request.SearchSecretName = &secretName
	request.SecretType = &serviceType

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	instances := make([]*ssm.SecretMetadata, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSsmClient().ListSecrets(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.SecretMetadatas) < 1 {
			break
		}
		instances = append(instances, response.Response.SecretMetadatas...)
		if len(response.Response.SecretMetadatas) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	sshKeyPairSecret = instances[0]
	return
}

func (me *SsmService) DeleteSsmSshKeyPairSecretById(ctx context.Context, secretName string, cleanSSHKey *bool) (errRet error) {
	logId := getLogId(ctx)

	request := ssm.NewDeleteSecretRequest()
	request.SecretName = &secretName

	if cleanSSHKey != nil {
		request.CleanSSHKey = cleanSSHKey
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSsmClient().DeleteSecret(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SsmService) DeleteSsmProductSecretById(ctx context.Context, secretName string) (errRet error) {
	logId := getLogId(ctx)

	request := ssm.NewDeleteSecretRequest()
	request.SecretName = &secretName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSsmClient().DeleteSecret(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SsmService) DescribeAsyncRequestInfo(ctx context.Context, flowID int64) (taskStatus int64, errRet error) {
	logId := getLogId(ctx)

	request := ssm.NewDescribeAsyncRequestInfoRequest()
	request.FlowID = helper.Int64(flowID)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSsmClient().DescribeAsyncRequestInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	taskStatus = *response.Response.TaskStatus
	return
}

func (me *SsmService) SsmProductSecretStateRefreshFunc(flowId int64, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		taskStatus, err := me.DescribeAsyncRequestInfo(ctx, flowId)

		if err != nil {
			return nil, "", err
		}

		return taskStatus, helper.Int64ToStr(taskStatus), nil
	}
}

func (me *SsmService) DescribeSsmProductsByFilter(ctx context.Context) (products []*string, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = ssm.NewDescribeSupportedProductsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	response, err := me.client.UseSsmClient().DescribeSupportedProducts(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	products = response.Response.Products

	return
}
