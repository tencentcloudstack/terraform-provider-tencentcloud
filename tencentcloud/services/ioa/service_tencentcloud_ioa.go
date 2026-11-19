package ioa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	ioav20220601 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ioa/v20220601"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewIoaService(client *connectivity.TencentCloudClient) IoaService {
	return IoaService{client: client}
}

type IoaService struct {
	client *connectivity.TencentCloudClient
}

func (me *IoaService) DescribeCompanyDirectoryConfigById(ctx context.Context, id string) (ret *ioav20220601.DirectoryConfigData, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ioav20220601.NewDescribeCompanyDirectoryConfigRequest()
	response := ioav20220601.NewDescribeCompanyDirectoryConfigResponse()
	request.Id = helper.StrToInt64Point(id)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseIoaV20220601Client().DescribeCompanyDirectoryConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe ioa_company_directory_config failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.Data
	return
}
