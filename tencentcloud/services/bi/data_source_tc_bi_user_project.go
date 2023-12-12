package bi

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bi/v20220105"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudBiUserProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudBiUserProjectRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Project id.",
			},

			"all_page": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to display all, if true, ignore paging.",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array(Note: This field may return null, indicating that no valid value can be obtained).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User id.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Username.",
						},
						"corp_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enterprise id(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "E-mail(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"last_login": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last login time, public cloud unrelated fields(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disabled state(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"first_modify": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "First login to change password, public cloud unrelated fields(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"phone_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Phone number(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"area_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mobile area code(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"created_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created by(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created at(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"updated_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updated by(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updated at(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"global_user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Global role name(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"mobile": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mobile number, public cloud unrelated fields(Note: This field may return null, indicating that no valid value can be obtained).",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudBiUserProjectRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_bi_user_project.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOkExists("project_id"); ok {
		paramMap["ProjectId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("all_page"); ok {
		paramMap["AllPage"] = helper.Bool(v.(bool))
	}

	service := BiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var data []*bi.UserIdAndUserName
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeBiUserProjectByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		data = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(data))
	dataList := []interface{}{}
	if data != nil {
		for _, list := range data {
			listMap := map[string]interface{}{}

			if list.UserId != nil {
				listMap["user_id"] = list.UserId
			}

			if list.UserName != nil {
				listMap["user_name"] = list.UserName
			}

			if list.CorpId != nil {
				listMap["corp_id"] = list.CorpId
			}

			if list.Email != nil {
				listMap["email"] = list.Email
			}

			if list.LastLogin != nil {
				listMap["last_login"] = list.LastLogin
			}

			if list.Status != nil {
				listMap["status"] = list.Status
			}

			if list.FirstModify != nil {
				listMap["first_modify"] = list.FirstModify
			}

			if list.PhoneNumber != nil {
				listMap["phone_number"] = list.PhoneNumber
			}

			if list.AreaCode != nil {
				listMap["area_code"] = list.AreaCode
			}

			if list.CreatedUser != nil {
				listMap["created_user"] = list.CreatedUser
			}

			if list.CreatedAt != nil {
				listMap["created_at"] = list.CreatedAt
			}

			if list.UpdatedUser != nil {
				listMap["updated_user"] = list.UpdatedUser
			}

			if list.UpdatedAt != nil {
				listMap["updated_at"] = list.UpdatedAt
			}

			if list.GlobalUserName != nil {
				listMap["global_user_name"] = list.GlobalUserName
			}

			if list.Mobile != nil {
				listMap["mobile"] = list.Mobile
			}

			dataList = append(dataList, listMap)
		}
		_ = d.Set("list", dataList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataList); e != nil {
			return e
		}
	}
	return nil
}
