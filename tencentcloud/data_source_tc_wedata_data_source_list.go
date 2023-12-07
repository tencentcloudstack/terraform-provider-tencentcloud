package tencentcloud

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudWedataDataSourceList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataDataSourceListRead,
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
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "Filter value.",
						},
					},
				},
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
						"id": {
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
							Description: "Datasource can Edit.",
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
							Description: "DatasourceDataSourceStatus.",
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
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudWedataDataSourceListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_wedata_data_source_list.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId          = getLogId(contextNil)
		ctx            = context.WithValue(context.TODO(), logIdKey, logId)
		service        = WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
		dataSourceList []*wedata.DataSourceInfo
	)

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

		paramMap["OrderFields"] = tmpSet
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

		paramMap["Filters"] = tmpSet
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataDataSourceListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		dataSourceList = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(dataSourceList))
	tmpList := make([]map[string]interface{}, 0, len(dataSourceList))
	if dataSourceList != nil {

		for _, rows := range dataSourceList {
			rowsMap := map[string]interface{}{}

			if rows.DatabaseName != nil {
				rowsMap["database_name"] = rows.DatabaseName
			}

			if rows.Description != nil {
				rowsMap["description"] = rows.Description
			}

			if rows.ID != nil {
				rowsMap["id"] = rows.ID
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

			rowIdInt := *rows.ID
			rowId := strconv.FormatUint(rowIdInt, 10)
			ids = append(ids, rowId)
			tmpList = append(tmpList, rowsMap)
		}

		_ = d.Set("rows", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
