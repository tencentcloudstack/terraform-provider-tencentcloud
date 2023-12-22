package kms

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/pkg/errors"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewKmsService(client *connectivity.TencentCloudClient) KmsService {
	return KmsService{client: client}
}

type KmsService struct {
	client *connectivity.TencentCloudClient
}

func (me *KmsService) DescribeKeysByFilter(ctx context.Context, param map[string]interface{}) (keys []*kms.KeyMetadata, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := kms.NewListKeyDetailRequest()

	for k, v := range param {
		if k == "role" {
			request.Role = helper.Uint64(uint64(v.(int)))
		}
		if k == "order_type" {
			request.OrderType = helper.Uint64(uint64(v.(int)))
		}
		if k == "key_state" {
			request.KeyState = helper.Uint64(v.(uint64))
		}
		if k == "search_key_alias" {
			request.SearchKeyAlias = helper.String(v.(string))
		}
		if k == "origin" {
			request.Origin = helper.String(v.(string))
		}
		if k == "key_usage" {
			request.KeyUsage = helper.String(v.(string))
		}
		if k == "tag_filter" {
			tagFilter := v.(map[string]string)
			for tagKey, tagValue := range tagFilter {
				tag := kms.TagFilter{
					TagKey:   helper.String(tagKey),
					TagValue: []*string{helper.String(tagValue)},
				}
				request.TagFilters = append(request.TagFilters, &tag)
			}
		}
	}
	var offset uint64 = 0
	var pageSize = uint64(KMS_PAGE_LIMIT)
	keys = make([]*kms.KeyMetadata, 0)
	for {
		request.Offset = helper.Uint64(offset)
		request.Limit = helper.Uint64(pageSize)
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseKmsClient().ListKeyDetail(request)
		if err != nil {
			errRet = errors.WithStack(err)
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		if response == nil || len(response.Response.KeyMetadatas) < 1 {
			break
		}

		keys = append(keys, response.Response.KeyMetadatas...)

		if uint64(len(response.Response.KeyMetadatas)) < pageSize {
			break
		}
		offset += pageSize
	}
	return
}

func (me *KmsService) DescribeKeyById(ctx context.Context, keyId string) (key *kms.KeyMetadata, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := kms.NewDescribeKeyRequest()
	request.KeyId = helper.String(keyId)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().DescribeKey(request)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	key = response.Response.KeyMetadata
	return
}

func (me *KmsService) CreateKey(ctx context.Context, keyType uint64, alias, description, keyUsage string) (keyId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := kms.NewCreateKeyRequest()
	request.Type = helper.Uint64(keyType)
	request.Alias = helper.String(alias)
	if description != "" {
		request.Description = helper.String(description)
	}
	if keyUsage != "" {
		request.KeyUsage = helper.String(keyUsage)
	}
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().CreateKey(request)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	keyId = *response.Response.KeyId
	return
}

func (me *KmsService) EnableKeyRotation(ctx context.Context, keyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := kms.NewEnableKeyRotationRequest()
	request.KeyId = helper.String(keyId)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().EnableKeyRotation(request)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *KmsService) DisableKeyRotation(ctx context.Context, keyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := kms.NewDisableKeyRotationRequest()
	request.KeyId = helper.String(keyId)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().DisableKeyRotation(request)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *KmsService) EnableKey(ctx context.Context, keyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := kms.NewEnableKeyRequest()
	request.KeyId = helper.String(keyId)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().EnableKey(request)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *KmsService) DisableKey(ctx context.Context, keyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := kms.NewDisableKeyRequest()
	request.KeyId = helper.String(keyId)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().DisableKey(request)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *KmsService) ArchiveKey(ctx context.Context, keyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := kms.NewArchiveKeyRequest()
	request.KeyId = helper.String(keyId)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().ArchiveKey(request)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *KmsService) CancelKeyArchive(ctx context.Context, keyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := kms.NewCancelKeyArchiveRequest()
	request.KeyId = helper.String(keyId)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().CancelKeyArchive(request)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}
func (me *KmsService) CancelKeyDeletion(ctx context.Context, keyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := kms.NewCancelKeyDeletionRequest()
	request.KeyId = helper.String(keyId)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().CancelKeyDeletion(request)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *KmsService) UpdateKeyDescription(ctx context.Context, keyId, description string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := kms.NewUpdateKeyDescriptionRequest()
	request.KeyId = helper.String(keyId)
	request.Description = helper.String(description)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().UpdateKeyDescription(request)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *KmsService) UpdateKeyAlias(ctx context.Context, keyId, alias string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := kms.NewUpdateAliasRequest()
	request.KeyId = helper.String(keyId)
	request.Alias = helper.String(alias)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().UpdateAlias(request)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *KmsService) DeleteKey(ctx context.Context, keyId string, pendingDeleteWindowInDays uint64) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := kms.NewScheduleKeyDeletionRequest()
	request.KeyId = helper.String(keyId)
	request.PendingWindowInDays = helper.Uint64(pendingDeleteWindowInDays)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().ScheduleKeyDeletion(request)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *KmsService) ImportKeyMaterial(ctx context.Context, param map[string]interface{}) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	var keyId, wrappingAlgorithm, wrappingKeySpec, keyMaterialBase64 string
	var validTo uint64

	for k, v := range param {
		if k == "key_id" {
			keyId = v.(string)
		}
		if k == "algorithm" {
			wrappingAlgorithm = v.(string)
		}
		if k == "key_spec" {
			wrappingKeySpec = v.(string)
		}
		if k == "valid_to" {
			validTo = uint64(v.(int))
		}
		if k == "key_material_base64" {
			keyMaterialBase64 = v.(string)
		}
	}
	paramRequest := kms.NewGetParametersForImportRequest()
	paramRequest.KeyId = helper.String(keyId)
	paramRequest.WrappingAlgorithm = helper.String(wrappingAlgorithm)
	paramRequest.WrappingKeySpec = helper.String(wrappingKeySpec)
	ratelimit.Check(paramRequest.GetAction())

	paramResponse, err := me.client.UseKmsClient().GetParametersForImport(paramRequest)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, paramRequest.GetAction(), paramRequest.ToJsonString(), paramResponse.ToJsonString())

	keyMaterial, err := base64.StdEncoding.DecodeString(keyMaterialBase64)
	if err != nil {
		return fmt.Errorf("error Base64 decoding key material: %s", err)
	}
	publicKeyBytes, err := base64.StdEncoding.DecodeString(*paramResponse.Response.PublicKey)
	if err != nil {
		return fmt.Errorf("error Base64 decoding public key: %s", err)
	}
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return fmt.Errorf("error parsing public key: %s", err)
	}

	var encryptedKeyMaterial string
	if wrappingAlgorithm == KMS_WRAPPING_ALGORITHM_RSAES_OAEP_SHA_1 {
		encryptedKeyMaterialBytes, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, publicKey.(*rsa.PublicKey), keyMaterial, []byte{})
		if err != nil {
			return fmt.Errorf("error encrypting key material: %s", err)
		}
		encryptedKeyMaterial = base64.StdEncoding.EncodeToString(encryptedKeyMaterialBytes)
	} else if wrappingAlgorithm == KMS_WRAPPING_ALGORITHM_RSAES_OAEP_SHA_256 {
		encryptedKeyMaterialBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey.(*rsa.PublicKey), keyMaterial, []byte{})
		if err != nil {
			return fmt.Errorf("error encrypting key material: %s", err)
		}
		encryptedKeyMaterial = base64.StdEncoding.EncodeToString(encryptedKeyMaterialBytes)
	} else {
		encryptedKeyMaterialBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), keyMaterial)
		if err != nil {
			return fmt.Errorf("error encrypting key material: %s", err)
		}
		encryptedKeyMaterial = base64.StdEncoding.EncodeToString(encryptedKeyMaterialBytes)
	}

	importRequest := kms.NewImportKeyMaterialRequest()
	importRequest.KeyId = helper.String(keyId)
	importRequest.ValidTo = helper.Uint64(validTo)
	importRequest.ImportToken = paramResponse.Response.ImportToken
	importRequest.EncryptedKeyMaterial = helper.String(encryptedKeyMaterial)
	ratelimit.Check(importRequest.GetAction())

	importResponse, err := me.client.UseKmsClient().ImportKeyMaterial(importRequest)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, importRequest.GetAction(), importRequest.ToJsonString(), importResponse.ToJsonString())

	return nil
}

func (me *KmsService) DeleteImportKeyMaterial(ctx context.Context, keyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := kms.NewDeleteImportedKeyMaterialRequest()
	request.KeyId = helper.String(keyId)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().DeleteImportedKeyMaterial(request)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *KmsService) DescribeKmsPublicKeyByFilter(ctx context.Context, param map[string]interface{}) (publicKey *kms.GetPublicKeyResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = kms.NewGetPublicKeyRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "KeyId" {
			request.KeyId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().GetPublicKey(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	publicKey = response.Response
	return
}

func (me *KmsService) DescribeKmsGetParametersForImportByFilter(ctx context.Context, param map[string]interface{}) (getParametersForImport *kms.GetParametersForImportResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = kms.NewGetParametersForImportRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "KeyId" {
			request.KeyId = v.(*string)
		}

		if k == "WrappingAlgorithm" {
			request.WrappingAlgorithm = v.(*string)
		}

		if k == "WrappingKeySpec" {
			request.WrappingKeySpec = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().GetParametersForImport(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	getParametersForImport = response.Response
	return
}

func (me *KmsService) DescribeKmsKeyListsByFilter(ctx context.Context, param map[string]interface{}) (KeyLists []*kms.KeyMetadata, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = kms.NewDescribeKeysRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "KeyIds" {
			request.KeyIds = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().DescribeKeys(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.KeyMetadatas) < 1 {
		return
	}

	KeyLists = response.Response.KeyMetadatas

	return
}

func (me *KmsService) DescribeKmsWhiteBoxKeyDetailsByFilter(ctx context.Context, param map[string]interface{}) (whiteBoxKeyInfo []*kms.WhiteboxKeyInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = kms.NewDescribeWhiteBoxKeyDetailsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "KeyStatus" {
			request.KeyStatus = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseKmsClient().DescribeWhiteBoxKeyDetails(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.KeyInfos) < 1 {
			break
		}

		whiteBoxKeyInfo = append(whiteBoxKeyInfo, response.Response.KeyInfos...)
		if len(response.Response.KeyInfos) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *KmsService) DescribeKmsListKeysByFilter(ctx context.Context, param map[string]interface{}) (listKeys []*kms.Key, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = kms.NewListKeysRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Role" {
			request.Role = v.(*uint64)
		}

		if k == "HsmClusterId" {
			request.HsmClusterId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseKmsClient().ListKeys(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Keys) < 1 {
			break
		}

		listKeys = append(listKeys, response.Response.Keys...)
		if len(response.Response.Keys) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *KmsService) DescribeKmsWhiteBoxKeyById(ctx context.Context, keyId string) (whiteBoxKey *kms.WhiteboxKeyInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := kms.NewDescribeWhiteBoxKeyRequest()
	request.KeyId = &keyId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().DescribeWhiteBoxKey(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.KeyInfo == nil {
		return
	}

	whiteBoxKey = response.Response.KeyInfo
	return
}

func (me *KmsService) DeleteKmsWhiteBoxKeyById(ctx context.Context, keyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := kms.NewDeleteWhiteBoxKeyRequest()
	request.KeyId = &keyId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().DeleteWhiteBoxKey(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *KmsService) DescribeKmsCloudResourceAttachmentById(ctx context.Context, keyId string) (keyMetadata *kms.KeyMetadata, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := kms.NewDescribeKeyRequest()
	request.KeyId = &keyId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().DescribeKey(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	keyMetadata = response.Response.KeyMetadata
	return
}

func (me *KmsService) DeleteKmsCloudResourceAttachmentById(ctx context.Context, keyId, productId, resourceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := kms.NewUnbindCloudResourceRequest()
	request.KeyId = &keyId
	request.ProductId = &productId
	request.ResourceId = &resourceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().UnbindCloudResource(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *KmsService) DescribeKmsWhiteBoxDecryptKeyByFilter(ctx context.Context, param map[string]interface{}) (whiteBoxDecryptKey *kms.DescribeWhiteBoxDecryptKeyResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = kms.NewDescribeWhiteBoxDecryptKeyRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "KeyId" {
			request.KeyId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().DescribeWhiteBoxDecryptKey(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	whiteBoxDecryptKey = response.Response
	return
}

func (me *KmsService) DescribeKmsWhiteBoxDeviceFingerprintsByFilter(ctx context.Context, param map[string]interface{}) (whiteBoxDeviceFingerprints []*kms.DeviceFingerprint, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = kms.NewDescribeWhiteBoxDeviceFingerprintsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "KeyId" {
			request.KeyId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().DescribeWhiteBoxDeviceFingerprints(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	whiteBoxDeviceFingerprints = response.Response.DeviceFingerprints
	return
}

func (me *KmsService) DescribeKmsListAlgorithmsByFilter(ctx context.Context) (listAlgorithms *kms.ListAlgorithmsResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = kms.NewListAlgorithmsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().ListAlgorithms(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	listAlgorithms = response.Response
	return
}
