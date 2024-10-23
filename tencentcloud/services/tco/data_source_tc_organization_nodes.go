package tco

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOrganizationNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOrganizationNodesRead,
		Schema: map[string]*schema.Schema{
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Department tag search list, with a maximum of 10.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Organization node ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name.",
						},
						"parent_node_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Parent node ID.",
						},
						"remark": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Remarks.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Update time.",
						},
						"tags": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Member tag list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag value.",
									},
								},
							},
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

func dataSourceTencentCloudOrganizationNodesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_organization_nodes.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("tags"); ok {
		tagsSet := v.([]interface{})
		tmpSet := make([]*organization.Tag, 0, len(tagsSet))
		for _, item := range tagsSet {
			tagsMap := item.(map[string]interface{})
			tag := organization.Tag{}
			if v, ok := tagsMap["tag_key"]; ok {
				tag.TagKey = helper.String(v.(string))
			}
			if v, ok := tagsMap["tag_value"]; ok {
				tag.TagValue = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &tag)
		}
		paramMap["Tags"] = tmpSet
	}

	var nodes []*organization.OrgNode
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOrganizationNodesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		nodes = result
		return nil
	})
	if err != nil {
		return err
	}

	itemList := make([]map[string]interface{}, 0, len(nodes))
	ids := make([]string, 0, len(nodes))
	for _, item := range nodes {
		itemMap := map[string]interface{}{}

		if item.NodeId != nil {
			itemMap["node_id"] = item.NodeId
			nodeIdStr := strconv.FormatInt(*item.NodeId, 10)
			ids = append(ids, nodeIdStr)
		}

		if item.Name != nil {
			itemMap["name"] = item.Name
		}

		if item.ParentNodeId != nil {
			itemMap["parent_node_id"] = item.ParentNodeId
		}

		if item.Remark != nil {
			itemMap["remark"] = item.Remark
		}

		if item.CreateTime != nil {
			itemMap["create_time"] = item.CreateTime
		}

		if item.UpdateTime != nil {
			itemMap["update_time"] = item.UpdateTime
		}

		tagsList := make([]map[string]interface{}, 0, len(item.Tags))
		if item.Tags != nil {
			for _, tags := range item.Tags {
				tagsMap := map[string]interface{}{}

				if tags.TagKey != nil {
					tagsMap["tag_key"] = tags.TagKey
				}

				if tags.TagValue != nil {
					tagsMap["tag_value"] = tags.TagValue
				}

				tagsList = append(tagsList, tagsMap)
			}

			itemMap["tags"] = tagsList
		}
		itemList = append(itemList, itemMap)
	}

	_ = d.Set("items", itemList)

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), itemList); e != nil {
			return e
		}
	}

	return nil
}
