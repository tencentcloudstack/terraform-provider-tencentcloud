package tco

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOrganizationOrgShareUnitNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOrganizationOrgShareUnitNodesRead,
		Schema: map[string]*schema.Schema{
			"unit_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Shared unit ID.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Offset, default is 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Limit, range 1-50, default is 10.",
			},

			"search_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search key, supports searching by department ID.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of share unit nodes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"share_node_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Department ID.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time.",
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

func dataSourceTencentCloudOrganizationOrgShareUnitNodesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_organization_org_share_unit_nodes.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("unit_id"); ok {
		paramMap["UnitId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("offset"); ok {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("limit"); ok {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("search_key"); ok {
		paramMap["SearchKey"] = helper.String(v.(string))
	}

	orgShareUnitNodes, err := service.DescribeOrganizationOrgShareUnitNodesByFilter(ctx, paramMap)
	if err != nil {
		return err
	}

	tmpList := make([]interface{}, 0, len(orgShareUnitNodes))
	if orgShareUnitNodes != nil {
		for _, shareUnitNode := range orgShareUnitNodes {
			shareUnitNodeMap := map[string]interface{}{}

			if shareUnitNode.ShareNodeId != nil {
				shareUnitNodeMap["share_node_id"] = int(*shareUnitNode.ShareNodeId)
			}

			if shareUnitNode.CreateTime != nil {
				shareUnitNodeMap["create_time"] = *shareUnitNode.CreateTime
			}

			tmpList = append(tmpList, shareUnitNodeMap)
		}

		_ = d.Set("items", tmpList)
	}

	idParam := ""
	if v, ok := d.GetOk("unit_id"); ok {
		idParam = v.(string)
	}

	d.SetId(helper.DataResourceIdsHash([]string{idParam}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n", logId, output.(string), e.Error())
			return e
		}
	}

	return nil
}
