package cam

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudUserInfo() *schema.Resource {
	return &schema.Resource{
		Read: datasourceTencentCloudUserInfoRead,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current account App ID.",
			},

			"uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current account UIN.",
			},

			"owner_uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current account OwnerUIN.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current account Name. NOTE: only support subaccount.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used for save results.",
			},
		},
	}
}

func datasourceTencentCloudUserInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("datasource.tencentcloud_user_info.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()

	logId = tccommon.GetLogId(ctx)
	request := cam.NewGetUserAppIdRequest()
	response := cam.NewGetUserAppIdResponse()

	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := client.UseCamClient().GetUserAppId(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return fmt.Errorf("get user appid error: empty response")
	}
	var appId, uin, ownerUin string
	accountInfoRequest := cam.NewDescribeSubAccountsRequest()
	accountInfoResponse := cam.NewDescribeSubAccountsResponse()

	if response.Response.AppId != nil {
		appId = strconv.FormatUint(*response.Response.AppId, 10)
	}
	if response.Response.Uin != nil {
		uin = *response.Response.Uin
	}
	if response.Response.OwnerUin != nil {
		ownerUin = *response.Response.OwnerUin
	}
	accountInfoRequest.FilterSubAccountUin = []*uint64{helper.Uint64(helper.StrToUInt64(uin))}

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		accountInfoResult, e := client.UseCamClient().DescribeSubAccounts(accountInfoRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}
		accountInfoResponse = accountInfoResult
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM users failed, reason:%s\n", logId, err.Error())
		return err
	}

	subAccounts := accountInfoResponse.Response.SubAccounts
	var name string
	if len(subAccounts) > 0 {
		name = *subAccounts[0].Name
	}
	d.SetId(fmt.Sprintf("user-%s-%s-%d", uin, appId, rand.Intn(10000)))

	_ = d.Set("app_id", appId)
	_ = d.Set("uin", uin)
	_ = d.Set("owner_uin", ownerUin)
	_ = d.Set("name", name)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), map[string]interface{}{
			"app_id":   appId,
			"uin":      uin,
			"ownerUin": ownerUin,
			"name":     name,
		}); e != nil {
			return e
		}
	}

	return nil
}
