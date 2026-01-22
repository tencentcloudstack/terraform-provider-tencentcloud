package bh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudBhAccountGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudBhAccountGroupsRead,
		Schema: map[string]*schema.Schema{
			"deep_in": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to recursively query, 0 for non-recursive, 1 for recursive.",
			},

			"parent_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Parent account group ID, default 0, query all groups under the root account group.",
			},

			"group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Account group name, fuzzy query.",
			},

			"page_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Get data from which page.",
			},

			"account_group_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Account group information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Account group ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account group name.",
						},
						"id_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account group ID path.",
						},
						"name_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account group name path.",
						},
						"parent_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Parent account group ID.",
						},
						"source": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Account group source.",
						},
						"user_total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of users under the account group.",
						},
						"is_leaf": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether it is a leaf node.",
						},
						"import_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account group import type.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account group description.",
						},
						"parent_org_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parent source account organization ID. When using third-party import user sources, record the group ID of this group in the source organization structure.",
						},
						"org_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source account organization ID. When using third-party import user sources, record the group ID of this group in the source organization structure.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the account group has been connected, 0 means not connected, 1 means connected.",
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

func dataSourceTencentCloudBhAccountGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_bh_account_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = BhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOkExists("deep_in"); ok {
		paramMap["DeepIn"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("parent_id"); ok {
		paramMap["ParentId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("group_name"); ok {
		paramMap["GroupName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("page_num"); ok {
		paramMap["PageNum"] = helper.IntInt64(v.(int))
	}

	var respData []*bhv20230418.AccountGroup
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeBhAccountGroupsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	accountGroupSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, accountGroupSet := range respData {
			accountGroupSetMap := map[string]interface{}{}
			if accountGroupSet.Id != nil {
				accountGroupSetMap["id"] = accountGroupSet.Id
			}

			if accountGroupSet.Name != nil {
				accountGroupSetMap["name"] = accountGroupSet.Name
			}

			if accountGroupSet.IdPath != nil {
				accountGroupSetMap["id_path"] = accountGroupSet.IdPath
			}

			if accountGroupSet.NamePath != nil {
				accountGroupSetMap["name_path"] = accountGroupSet.NamePath
			}

			if accountGroupSet.ParentId != nil {
				accountGroupSetMap["parent_id"] = accountGroupSet.ParentId
			}

			if accountGroupSet.Source != nil {
				accountGroupSetMap["source"] = accountGroupSet.Source
			}

			if accountGroupSet.UserTotal != nil {
				accountGroupSetMap["user_total"] = accountGroupSet.UserTotal
			}

			if accountGroupSet.IsLeaf != nil {
				accountGroupSetMap["is_leaf"] = accountGroupSet.IsLeaf
			}

			if accountGroupSet.ImportType != nil {
				accountGroupSetMap["import_type"] = accountGroupSet.ImportType
			}

			if accountGroupSet.Description != nil {
				accountGroupSetMap["description"] = accountGroupSet.Description
			}

			if accountGroupSet.ParentOrgId != nil {
				accountGroupSetMap["parent_org_id"] = accountGroupSet.ParentOrgId
			}

			if accountGroupSet.OrgId != nil {
				accountGroupSetMap["org_id"] = accountGroupSet.OrgId
			}

			if accountGroupSet.Status != nil {
				accountGroupSetMap["status"] = accountGroupSet.Status
			}

			accountGroupSetList = append(accountGroupSetList, accountGroupSetMap)
		}

		_ = d.Set("account_group_set", accountGroupSetList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), accountGroupSetList); e != nil {
			return e
		}
	}

	return nil
}
