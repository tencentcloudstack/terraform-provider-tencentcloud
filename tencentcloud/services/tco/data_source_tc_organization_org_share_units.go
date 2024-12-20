package tco

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOrganizationOrgShareUnits() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOrganizationOrgShareUnitsRead,
		Schema: map[string]*schema.Schema{
			"area": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Shared unit area.",
			},

			"search_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search for keywords. Support UnitId and Name searches.",
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Shared unit list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unit_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Shared unit ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Shared unit name.",
						},
						"uin": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Shared unit manager Uin.",
						},
						"owner_uin": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Shared unit manager OwnerUin.",
						},
						"area": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Shared unit area.",
						},
						"description": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "description.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Create time.",
						},
						"share_resource_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of shared unit resources.",
						},
						"share_member_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of shared unit members.",
						},
						"share_scope": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Shared scope. Value: 1-Only allowed to share within the group organization 2-allowed to share to any account.",
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

func dataSourceTencentCloudOrganizationOrgShareUnitsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_organization_org_share_units.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("area"); ok {
		paramMap["Area"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_key"); ok {
		paramMap["SearchKey"] = helper.String(v.(string))
	}

	var respData []*organization.ManagerShareUnit
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOrganizationOrgShareUnitsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	var ids []string
	itemsList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, items := range respData {
			itemsMap := map[string]interface{}{}

			var unitId string
			if items.UnitId != nil {
				itemsMap["unit_id"] = items.UnitId
				unitId = *items.UnitId
			}

			if items.Name != nil {
				itemsMap["name"] = items.Name
			}

			if items.Uin != nil {
				itemsMap["uin"] = items.Uin
			}

			if items.OwnerUin != nil {
				itemsMap["owner_uin"] = items.OwnerUin
			}

			if items.Area != nil {
				itemsMap["area"] = items.Area
			}

			if items.Description != nil {
				itemsMap["description"] = items.Description
			}

			if items.CreateTime != nil {
				itemsMap["create_time"] = items.CreateTime
			}

			if items.ShareResourceNum != nil {
				itemsMap["share_resource_num"] = items.ShareResourceNum
			}

			if items.ShareMemberNum != nil {
				itemsMap["share_member_num"] = items.ShareMemberNum
			}

			if items.ShareScope != nil {
				itemsMap["share_scope"] = items.ShareScope
			}

			ids = append(ids, unitId)
			itemsList = append(itemsList, itemsMap)
		}

		_ = d.Set("items", itemsList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), itemsList); e != nil {
			return e
		}
	}

	return nil
}
