package vcube

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	vcubev20220410 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vcube/v20220410"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewVcubeService(client *connectivity.TencentCloudClient) VcubeService {
	return VcubeService{client: client}
}

type VcubeService struct {
	client *connectivity.TencentCloudClient
}

func (me *VcubeService) DescribeVcubeApplicationAndVideoById(ctx context.Context, licenseId string) (ret *vcubev20220410.ApplicationInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vcubev20220410.NewDescribeVcubeApplicationAndPlayListRequest()
	response := vcubev20220410.NewDescribeVcubeApplicationAndPlayListResponse()
	request.LicenseId = helper.StrToUint64Point(licenseId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseVcubeV20220410Client().DescribeVcubeApplicationAndPlayList(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.ApplicationInfoList == nil || len(result.Response.ApplicationInfoList) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe vcube application and play list failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.ApplicationInfoList[0]
	return
}

func (me *VcubeService) DescribeVcubeApplicationAndWebPlayerLicenseById(ctx context.Context, licenseId string) (ret *vcubev20220410.ApplicationInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vcubev20220410.NewDescribeVcubeApplicationAndPlayListRequest()
	response := vcubev20220410.NewDescribeVcubeApplicationAndPlayListResponse()
	request.LicenseId = helper.StrToUint64Point(licenseId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseVcubeV20220410Client().DescribeVcubeApplicationAndPlayList(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.ApplicationInfoList == nil || len(result.Response.ApplicationInfoList) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe vcube application and play list failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.ApplicationInfoList[0]
	return
}
