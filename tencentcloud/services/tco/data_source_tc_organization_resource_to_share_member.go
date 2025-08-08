package tco

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organizationv20210331 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOrganizationResourceToShareMember() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOrganizationResourceToShareMemberRead,
		Schema: map[string]*schema.Schema{
			"area": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Area.",
			},

			"search_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search keywords, support business resource ID search.",
			},

			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource Type.",
			},

			"product_resource_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Business resource ID. Maximum 50.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Resource ID.",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Resource type.",
						},
						"unit_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Shared unit ID.",
						},
						"unit_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Shared unit name.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Create time.",
						},
						"product_resource_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Business resource ID.",
						},
						"share_manager_uin": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Shared Administrator uin.",
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

func dataSourceTencentCloudOrganizationResourceToShareMemberRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_organization_resource_to_share_member.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("area"); ok {
		paramMap["Area"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_key"); ok {
		paramMap["SearchKey"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product_resource_ids"); ok {
		productResourceIdsList := []*string{}
		productResourceIdsSet := v.(*schema.Set).List()
		for i := range productResourceIdsSet {
			productResourceIds := productResourceIdsSet[i].(string)
			productResourceIdsList = append(productResourceIdsList, helper.String(productResourceIds))
		}

		paramMap["ProductResourceIds"] = productResourceIdsList
	}

	var respData []*organizationv20210331.ShareResourceToMember
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOrganizationResourceToShareMemberByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	ids := make([]string, 0, len(respData))
	itemsList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, items := range respData {
			itemsMap := map[string]interface{}{}
			if items.ResourceId != nil {
				itemsMap["resource_id"] = items.ResourceId
			}

			if items.Type != nil {
				itemsMap["type"] = items.Type
			}

			if items.UnitId != nil {
				itemsMap["unit_id"] = items.UnitId
			}

			if items.UnitName != nil {
				itemsMap["unit_name"] = items.UnitName
			}

			if items.CreateTime != nil {
				itemsMap["create_time"] = items.CreateTime
			}

			if items.ProductResourceId != nil {
				itemsMap["product_resource_id"] = items.ProductResourceId
			}

			if items.ShareManagerUin != nil {
				itemsMap["share_manager_uin"] = items.ShareManagerUin
			}

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
