package cam

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"

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

	ratelimit.Check(request.GetAction())
	response, err := client.UseCamClient().GetUserAppId(request)

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if err != nil {
		return err
	}

	result := response.Response

	if result == nil {
		return fmt.Errorf("get user appid error: empty response")
	}

	appId := strconv.FormatUint(*result.AppId, 10)
	uin := *result.Uin
	ownerUin := *result.OwnerUin
	accountInfoRequest := cam.NewDescribeSubAccountsRequest()
	accountInfoRequest.FilterSubAccountUin = []*uint64{helper.Uint64(helper.StrToUInt64(uin))}
	accountInfoResponse, err := client.UseCamClient().DescribeSubAccounts(accountInfoRequest)
	if err != nil {
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
