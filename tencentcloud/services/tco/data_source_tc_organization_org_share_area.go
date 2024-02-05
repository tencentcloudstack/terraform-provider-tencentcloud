package tco

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOrganizationOrgShareArea() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOrganizationOrgShareAreaRead,
		Schema: map[string]*schema.Schema{
			"lang": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Language.default zh.\nValid values:\n  - `zh`: Chinese.\n  - `en`: English.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Area list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region name.",
						},
						"area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region identifier.",
						},
						"area_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Region ID.",
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

func dataSourceTencentCloudOrganizationOrgShareAreaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_organization_org_share_area.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("lang"); ok {
		paramMap["Lang"] = helper.String(v.(string))
	}

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var items []*organization.ShareArea
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOrganizationOrgShareAreaByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(items))
	tmpList := make([]map[string]interface{}, 0, len(items))

	if items != nil {
		for _, shareArea := range items {
			shareAreaMap := map[string]interface{}{}

			if shareArea.Name != nil {
				shareAreaMap["name"] = shareArea.Name
			}

			if shareArea.Area != nil {
				shareAreaMap["area"] = shareArea.Area
			}

			if shareArea.AreaId != nil {
				shareAreaMap["area_id"] = shareArea.AreaId
			}

			ids = append(ids, *shareArea.Area)
			tmpList = append(tmpList, shareAreaMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
