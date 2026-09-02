package cam

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCamAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamAccountsRead,
		Schema: map[string]*schema.Schema{
			"max_items": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum number of accounts to return per request. Valid range: [1, 100]. When the returned result is truncated due to reaching MaxItems, the output `is_truncated` will be true.",
			},
			"marker": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "When the returned result is truncated, use Marker to fetch the content after the current truncation position. Output `marker` carries the next page marker when `is_truncated` is true.",
			},
			"user_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Account type to filter by. Valid values: `Owner` (master account), `SubUser` (sub user), `CICUser` (CIC sub user), `WechatCorpUser` (WeCom sub user), `AgentIdentity` (AgentIdentity sub user), `Collaborator` (collaborator), `MessageReceiver` (message receiver).",
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of CAM accounts. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Account ID (Uin).",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account name.",
						},
						"uid": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Account UID.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account remark.",
						},
						"console_login": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the account can log in to the console. Returned as the raw int value from the API.",
						},
						"phone_num": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Phone number.",
						},
						"country_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Country code.",
						},
						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Email.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time. Format: YYYY-MM-DD hh:mm:ss.",
						},
						"user_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account type. Valid values: `Owner`, `SubUser`, `CICUser`, `WechatCorpUser`, `AgentIdentity`, `Collaborator`, `MessageReceiver`, `Unknown`.",
						},
					},
				},
			},
			"is_truncated": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the returned result is truncated.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCamAccountsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cam_accounts.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CamService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("max_items"); ok {
		paramMap["MaxItems"] = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("marker"); ok {
		paramMap["Marker"] = helper.String(v.(string))
	}
	if v, ok := d.GetOk("user_type"); ok {
		paramMap["UserType"] = helper.String(v.(string))
	}

	var (
		respData    []*cam.ListAllUser
		respMarker  string
		isTruncated bool
	)
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		users, marker, truncated, e := service.DescribeCamAccountsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if len(users) == 0 {
			return resource.NonRetryableError(fmt.Errorf("cam accounts is empty"))
		}

		respData = users
		respMarker = marker
		isTruncated = truncated
		return nil
	})

	if reqErr != nil {
		log.Printf("[DATASOURCE] read empty, skip SetId, cam_accounts id=%s", d.Id())
		return reqErr
	}

	ids := make([]string, 0, len(respData))
	usersList := make([]map[string]interface{}, 0, len(respData))
	for _, user := range respData {
		userMap := map[string]interface{}{}
		if user.Uin != nil {
			userMap["uin"] = user.Uin
			ids = append(ids, helper.Int64ToStr(*user.Uin))
		}
		if user.Name != nil {
			userMap["name"] = user.Name
		}
		if user.Uid != nil {
			userMap["uid"] = user.Uid
		}
		if user.Remark != nil {
			userMap["remark"] = user.Remark
		}
		if user.ConsoleLogin != nil {
			userMap["console_login"] = user.ConsoleLogin
		}
		if user.PhoneNum != nil {
			userMap["phone_num"] = user.PhoneNum
		}
		if user.CountryCode != nil {
			userMap["country_code"] = user.CountryCode
		}
		if user.Email != nil {
			userMap["email"] = user.Email
		}
		if user.CreateTime != nil {
			userMap["create_time"] = user.CreateTime
		}
		if user.UserType != nil {
			userMap["user_type"] = user.UserType
		}
		usersList = append(usersList, userMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("users", usersList)
	_ = d.Set("marker", respMarker)
	_ = d.Set("is_truncated", isTruncated)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), usersList); e != nil {
			return e
		}
	}

	return nil
}
