package makers

import (
	"context"
	"log"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// NewMakersCosService wraps a COS client scoped to a Makers temporary
// credential.
//
// A fresh client value is built instead of mutating the shared provider
// client: overwriting the shared Credential would leak the temporary key
// into every other API call, and copying the struct would also copy the
// private cached COS connection (built with the main credential), making the
// temporary credential silently ineffective. The fresh value starts with an
// empty connection cache, so UseTencentCosClient builds the connection with
// the temporary credential.
func NewMakersCosService(client *connectivity.TencentCloudClient, region, secretId, secretKey, token string) MakersCosService {
	scoped := &connectivity.TencentCloudClient{
		Credential: common.NewTokenCredential(secretId, secretKey, token),
		Region:     region,
	}
	if client != nil {
		// Respect a custom COS domain when the provider configured one.
		scoped.CosDomain = client.CosDomain
	}
	return MakersCosService{client: scoped}
}

type MakersCosService struct {
	client *connectivity.TencentCloudClient
}

// PutFile uploads the local file at filepath to the given bucket under the
// object name.
func (me *MakersCosService) PutFile(ctx context.Context, bucket, name, filepath string) (serviceResult *cos.CompleteMultipartUploadResult, errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "PutObject", bucket, errRet.Error())
		}
	}()

	result, response, err := me.client.UseTencentCosClient(bucket).Object.MultiUpload(ctx, name, filepath, nil)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s], http status [%s]\n", logId, "PutObject", bucket, result, response.Status)

	serviceResult = result
	return
}
