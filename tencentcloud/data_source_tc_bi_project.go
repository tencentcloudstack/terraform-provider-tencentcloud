package tencentcloud

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bi/v20220105"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudBiProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudBiProjectRead,
		Schema: map[string]*schema.Schema{
			"page_no": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Page number.",
			},

			"keyword": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Retrieve fuzzy fields.",
			},

			"all_page": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to display all, if true, ignore paging.",
			},

			"module_collection": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Role information, can be ignored.",
			},

			"extra": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Additional information(Note: This field may return null, indicating that no valid value can be obtained).",
			},

			"msg": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Interface information(Note: This field may return null, indicating that no valid value can be obtained).",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array(Note: This field may return null, indicating that no valid value can be obtained).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project id.",
						},
						"logo": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project logo(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project name(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"color_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Logo colour(Note: This field may return null, indicating that no valid value can be obtained).",
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
						"member_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Member count(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"page_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Page count(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"last_modify_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modified report and presentation names(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Interface call source(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"apply": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Apply(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"updated_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updated by(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updated by(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"corp_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enterprise id(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"mark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"seed": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Obfuscated field(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"auth_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "List of permissions within the project(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"panel_scope": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default kanban(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"is_external_manage": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Determine whether it is hosted(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"manage_platform": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Hosting platform name(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"config_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Customized parameters, this parameter can be ignored(Note: This field may return null, indicating that no valid value can be obtained).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"module_group": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Module group(Note: This field may return null, indicating that no valid value can be obtained).",
									},
									"components": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Components(Note: This field may return null, indicating that no valid value can be obtained).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"module_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Module id(Note: This field may return null, indicating that no valid value can be obtained).",
												},
												"include_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Include type(Note: This field may return null, indicating that no valid value can be obtained).",
												},
												"params": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Extra parameters(Note: This field may return null, indicating that no valid value can be obtained).",
												},
											},
										},
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

func dataSourceTencentCloudBiProjectRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_bi_project.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOkExists("page_no"); ok {
		paramMap["PageNo"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("keyword"); ok {
		paramMap["Keyword"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("all_page"); ok {
		paramMap["AllPage"] = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("module_collection"); ok {
		paramMap["ModuleCollection"] = helper.String(v.(string))
	}

	service := BiService{client: meta.(*TencentCloudClient).apiV3Conn}

	var project []*bi.Project
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeBiProjectByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		project = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(project))
	listList := []interface{}{}
	if project != nil {
		for _, list := range project {
			listMap := map[string]interface{}{}

			if list.Id != nil {
				listMap["id"] = list.Id
			}

			if list.Logo != nil {
				listMap["logo"] = list.Logo
			}

			if list.Name != nil {
				listMap["name"] = list.Name
			}

			if list.ColorCode != nil {
				listMap["color_code"] = list.ColorCode
			}

			if list.CreatedUser != nil {
				listMap["created_user"] = list.CreatedUser
			}

			if list.CreatedAt != nil {
				listMap["created_at"] = list.CreatedAt
			}

			if list.MemberCount != nil {
				listMap["member_count"] = list.MemberCount
			}

			if list.PageCount != nil {
				listMap["page_count"] = list.PageCount
			}

			if list.LastModifyName != nil {
				listMap["last_modify_name"] = list.LastModifyName
			}

			if list.Source != nil {
				listMap["source"] = list.Source
			}

			if list.Apply != nil {
				listMap["apply"] = list.Apply
			}

			if list.UpdatedUser != nil {
				listMap["updated_user"] = list.UpdatedUser
			}

			if list.UpdatedAt != nil {
				listMap["updated_at"] = list.UpdatedAt
			}

			if list.CorpId != nil {
				listMap["corp_id"] = list.CorpId
			}

			if list.Mark != nil {
				listMap["mark"] = list.Mark
			}

			if list.Seed != nil {
				listMap["seed"] = list.Seed
			}

			if list.AuthList != nil {
				listMap["auth_list"] = list.AuthList
			}

			if list.PanelScope != nil {
				listMap["panel_scope"] = list.PanelScope
			}

			if list.IsExternalManage != nil {
				listMap["is_external_manage"] = list.IsExternalManage
			}

			if list.ManagePlatform != nil {
				listMap["manage_platform"] = list.ManagePlatform
			}

			if list.ConfigList != nil {
				configListList := []interface{}{}
				for _, configList := range list.ConfigList {
					configListMap := map[string]interface{}{}

					if configList.ModuleGroup != nil {
						configListMap["module_group"] = configList.ModuleGroup
					}

					if configList.Components != nil {
						componentsList := []interface{}{}
						for _, components := range configList.Components {
							componentsMap := map[string]interface{}{}

							if components.ModuleId != nil {
								componentsMap["module_id"] = components.ModuleId
							}

							if components.IncludeType != nil {
								componentsMap["include_type"] = components.IncludeType
							}

							if components.Params != nil {
								componentsMap["params"] = components.Params
							}

							componentsList = append(componentsList, componentsMap)
						}

						configListMap["components"] = componentsList
					}

					configListList = append(configListList, configListMap)
				}

				listMap["config_list"] = configListList
			}

			listList = append(listList, listMap)
			ids = append(ids, strconv.FormatUint(*list.Id, 10))
		}

		_ = d.Set("list", listList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), listList); e != nil {
			return e
		}
	}
	return nil
}
