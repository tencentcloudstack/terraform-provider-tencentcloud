package cvm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func resourceTencentCloudKeyPairReadPostHandleResponse0(ctx context.Context, resp *cvm.KeyPair) error {
	d := tccommon.ResourceDataFromContext(ctx)

	if resp.PublicKey != nil {
		publicKey := *resp.PublicKey
		split := strings.Split(publicKey, " ")
		if len(split) > 2 {
			publicKey = strings.Join(split[0:2], " ")
		}
		_ = d.Set("public_key", publicKey)
	}
	return nil
}

func resourceTencentCloudKeyPairDeletePostFillRequest0(ctx context.Context, req *cvm.DeleteKeyPairsRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	cvmService := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	keyId := d.Id()

	var keyPair *cvm.KeyPair
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		keyPair, errRet = cvmService.DescribeKeyPairById(ctx, keyId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if keyPair == nil {
		d.SetId("")
		return nil
	}

	if len(keyPair.AssociatedInstanceIds) > 0 {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			errRet := cvmService.UnbindKeyPair(ctx, []*string{&keyId}, keyPair.AssociatedInstanceIds)
			if errRet != nil {
				if sdkErr, ok := errRet.(*errors.TencentCloudSDKError); ok {
					if sdkErr.Code == CVM_NOT_FOUND_ERROR {
						return nil
					}
				}
				return tccommon.RetryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func publicKeyStateFunc(v interface{}) string {
	switch value := v.(type) {
	case string:
		publicKey := value
		split := strings.Split(value, " ")
		if len(split) > 2 {
			publicKey = strings.Join(split[0:2], " ")
		}
		return strings.TrimSpace(publicKey)
	default:
		return ""
	}
}

func resourceTencentCloudKeyPairDeleteRequestOnError0(ctx context.Context, e error) *resource.RetryError {
	return tccommon.RetryError(e, KYE_PAIR_INVALID_ERROR, KEY_PAIR_NOT_SUPPORT_ERROR)
}

func resourceTencentCloudKeyPairCreateOnExit(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	var (
		keyId string
		err   error
	)
	if _, ok := d.GetOk("public_key"); ok {
		keyId, err = cvmCreateKeyPairByImportPublicKey(ctx, d, meta)
	} else {
		keyId, err = cvmCreateKeyPair(ctx, d, meta)
	}
	if err != nil {
		return err
	}
	d.SetId(keyId)
	return nil
}

func cvmCreateKeyPair(ctx context.Context, d *schema.ResourceData, meta interface{}) (keyId string, err error) {
	logId := tccommon.GetLogId(ctx)
	request := cvm.NewCreateKeyPairRequest()
	response := cvm.NewCreateKeyPairResponse()
	request.KeyName = helper.String(d.Get("key_name").(string))
	request.ProjectId = helper.IntInt64(d.Get("project_id").(int))

	innerErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().CreateKeyPair(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if innerErr != nil {
		log.Printf("[CRITAL]%s create cvm keyPair by import failed, reason:%+v", logId, err)
		err = innerErr
		return
	}
	if response == nil || response.Response == nil || response.Response.KeyPair == nil {
		err = fmt.Errorf("Response is nil")
		return
	}

	keyId = *response.Response.KeyPair.KeyId
	return
}

func cvmCreateKeyPairByImportPublicKey(ctx context.Context, d *schema.ResourceData, meta interface{}) (keyId string, err error) {
	logId := tccommon.GetLogId(ctx)
	request := cvm.NewImportKeyPairRequest()
	response := cvm.NewImportKeyPairResponse()
	request.KeyName = helper.String(d.Get("key_name").(string))
	request.ProjectId = helper.IntInt64(d.Get("project_id").(int))
	request.PublicKey = helper.String(d.Get("public_key").(string))

	innerErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ImportKeyPair(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if innerErr != nil {
		log.Printf("[CRITAL]%s create cvm keyPair by import failed, reason:%+v", logId, err)
		err = innerErr
		return
	}
	if response == nil || response.Response == nil {
		err = fmt.Errorf("Response is nil")
		return
	}

	keyId = *response.Response.KeyId
	return
}
