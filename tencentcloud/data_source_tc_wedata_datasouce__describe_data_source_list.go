/*
Use this data source to query detailed information of wedata datasouce__describe_data_source_list

Example Usage

```hcl
data "tencentcloud_wedata_datasouce__describe_data_source_list" "datasouce__describe_data_source_list" {
  order_fields {
		name = ""
		direction = ""

  }
  filters {
		name = ""
		values =

  }
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudWedataDatasouce_DescribeDataSourceList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataDatasouce_DescribeDataSourceListRead,
		Schema: map[string]*schema.Schema{
			"order_fields": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "OrderFields.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "OrderFields name.",
						},
						"direction": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "OrderFields rule.",
						},
					},
				},
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter name.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Filter value.",
						},
					},
				},
			},

			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"page_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Page number.",
						},
						"page_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Page size.",
						},
						"rows": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Data rows.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DatabaseName.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description.",
									},
									"i_d": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID.",
									},
									"instance": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource name.",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource engin cluster region.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource type.",
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource cluster id.",
									},
									"app_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Appid.",
									},
									"biz_params": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Biz params.",
									},
									"category": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource category.",
									},
									"display": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource display name.",
									},
									"owner_account": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource owner account.",
									},
									"params": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource params.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Datasource status.",
									},
									"owner_account_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource owner account name.",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource cluster name.",
									},
									"owner_project_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource owner project id.",
									},
									"owner_project_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource OwnerProjectName.",
									},
									"owner_project_ident": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource OwnerProjectIdent.",
									},
									"authority_project_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource AuthorityProjectName.",
									},
									"authority_user_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource AuthorityUserName.",
									},
									"edit": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Datasource can Edit .",
									},
									"author": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Has Author.",
									},
									"deliver": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Can Deliver.",
									},
									"data_source_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DatasourceDataSourceStatus .",
									},
									"create_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "CreateTime.",
									},
									"params_string": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Params json string.",
									},
									"biz_params_string": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Biz params json string.",
									},
									"modified_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Datasource ModifiedTime.",
									},
									"show_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource show type.",
									},
								},
							},
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TotalCount.",
						},
						"total_page_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TotalPageNumber.",
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

func dataSourceTencentCloudWedataDatasouce_DescribeDataSourceListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_wedata_datasouce__describe_data_source_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("order_fields"); ok {
		orderFieldsSet := v.([]interface{})
		tmpSet := make([]*wedata.OrderField, 0, len(orderFieldsSet))

		for _, item := range orderFieldsSet {
			orderField := wedata.OrderField{}
			orderFieldMap := item.(map[string]interface{})

			if v, ok := orderFieldMap["name"]; ok {
				orderField.Name = helper.String(v.(string))
			}
			if v, ok := orderFieldMap["direction"]; ok {
				orderField.Direction = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &orderField)
		}
		paramMap["order_fields"] = tmpSet
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*wedata.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := wedata.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}

	var data []*wedata.DataSourceInfoPage

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataDatasouce_DescribeDataSourceListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		data = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(data))
	if data != nil {
		dataSourceInfoPageMap := map[string]interface{}{}

		if data.PageNumber != nil {
			dataSourceInfoPageMap["page_number"] = data.PageNumber
		}

		if data.PageSize != nil {
			dataSourceInfoPageMap["page_size"] = data.PageSize
		}

		if data.Rows != nil {
			rowsList := []interface{}{}
			for _, rows := range data.Rows {
				rowsMap := map[string]interface{}{}

				if rows.DatabaseName != nil {
					rowsMap["database_name"] = rows.DatabaseName
				}

				if rows.Description != nil {
					rowsMap["description"] = rows.Description
				}

				if rows.ID != nil {
					rowsMap["i_d"] = rows.ID
				}

				if rows.Instance != nil {
					rowsMap["instance"] = rows.Instance
				}

				if rows.Name != nil {
					rowsMap["name"] = rows.Name
				}

				if rows.Region != nil {
					rowsMap["region"] = rows.Region
				}

				if rows.Type != nil {
					rowsMap["type"] = rows.Type
				}

				if rows.ClusterId != nil {
					rowsMap["cluster_id"] = rows.ClusterId
				}

				if rows.AppId != nil {
					rowsMap["app_id"] = rows.AppId
				}

				if rows.BizParams != nil {
					rowsMap["biz_params"] = rows.BizParams
				}

				if rows.Category != nil {
					rowsMap["category"] = rows.Category
				}

				if rows.Display != nil {
					rowsMap["display"] = rows.Display
				}

				if rows.OwnerAccount != nil {
					rowsMap["owner_account"] = rows.OwnerAccount
				}

				if rows.Params != nil {
					rowsMap["params"] = rows.Params
				}

				if rows.Status != nil {
					rowsMap["status"] = rows.Status
				}

				if rows.OwnerAccountName != nil {
					rowsMap["owner_account_name"] = rows.OwnerAccountName
				}

				if rows.ClusterName != nil {
					rowsMap["cluster_name"] = rows.ClusterName
				}

				if rows.OwnerProjectId != nil {
					rowsMap["owner_project_id"] = rows.OwnerProjectId
				}

				if rows.OwnerProjectName != nil {
					rowsMap["owner_project_name"] = rows.OwnerProjectName
				}

				if rows.OwnerProjectIdent != nil {
					rowsMap["owner_project_ident"] = rows.OwnerProjectIdent
				}

				if rows.AuthorityProjectName != nil {
					rowsMap["authority_project_name"] = rows.AuthorityProjectName
				}

				if rows.AuthorityUserName != nil {
					rowsMap["authority_user_name"] = rows.AuthorityUserName
				}

				if rows.Edit != nil {
					rowsMap["edit"] = rows.Edit
				}

				if rows.Author != nil {
					rowsMap["author"] = rows.Author
				}

				if rows.Deliver != nil {
					rowsMap["deliver"] = rows.Deliver
				}

				if rows.DataSourceStatus != nil {
					rowsMap["data_source_status"] = rows.DataSourceStatus
				}

				if rows.CreateTime != nil {
					rowsMap["create_time"] = rows.CreateTime
				}

				if rows.ParamsString != nil {
					rowsMap["params_string"] = rows.ParamsString
				}

				if rows.BizParamsString != nil {
					rowsMap["biz_params_string"] = rows.BizParamsString
				}

				if rows.ModifiedTime != nil {
					rowsMap["modified_time"] = rows.ModifiedTime
				}

				if rows.ShowType != nil {
					rowsMap["show_type"] = rows.ShowType
				}

				rowsList = append(rowsList, rowsMap)
			}

			dataSourceInfoPageMap["rows"] = []interface{}{rowsList}
		}

		if data.TotalCount != nil {
			dataSourceInfoPageMap["total_count"] = data.TotalCount
		}

		if data.TotalPageNumber != nil {
			dataSourceInfoPageMap["total_page_number"] = data.TotalPageNumber
		}

		ids = append(ids, *data.ID)
		_ = d.Set("data", dataSourceInfoPageMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), dataSourceInfoPageMap); e != nil {
			return e
		}
	}
	return nil
}
