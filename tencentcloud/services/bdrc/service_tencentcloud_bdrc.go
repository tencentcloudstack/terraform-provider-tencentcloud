package bdrc

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	bdrcv20260330 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewBdrcService(client *connectivity.TencentCloudClient) BdrcService {
	return BdrcService{client: client}
}

type BdrcService struct {
	client *connectivity.TencentCloudClient
}

func (me *BdrcService) RunCopyPairTasks(ctx context.Context, copyPairIds []*string, copyPairType string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bdrcv20260330.NewRunCopyPairTasksRequest()
	request.CopyPairIds = copyPairIds
	request.CopyPairType = &copyPairType

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBdrcV20260330Client().RunCopyPairTasksWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		errRet = err
		return err
	}

	return
}
