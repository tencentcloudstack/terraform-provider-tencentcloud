package makers

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	makerssdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewMakersService(client *connectivity.TencentCloudClient) MakersService {
	return MakersService{client: client}
}

type MakersService struct {
	client *connectivity.TencentCloudClient
}

// DescribeDeploymentLog returns the build log of the deployment identified by
// zoneId, projectId and deploymentId. The three keys form an exact query, so
// no pagination is involved. A nil LogUrl with a nil error is a valid outcome:
// the cloud side documents that LogUrl may be null when no log is available.
func (me *MakersService) DescribeDeploymentLog(ctx context.Context, request *makerssdk.DescribeDeploymentLogRequest) (logUrl *string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makerssdk.DescribeDeploymentLogResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersClient().DescribeDeploymentLogWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DescribeDeploymentLog returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	logUrl = response.Response.LogUrl
	return
}
