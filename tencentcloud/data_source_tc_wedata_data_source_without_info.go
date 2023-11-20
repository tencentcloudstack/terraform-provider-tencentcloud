/*
Use this data source to query detailed information of wedata data_source_without_info

Example Usage

```hcl
data "tencentcloud_wedata_data_source_without_info" "example" {
  filters {
    name   = "ownerProjectId"
    values = ["1612982498218618880"]
  }

  order_fields {
    name      = "create_time"
    direction = "DESC"
  }
}
```
*/
package tencentcloud

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudWedataDataSourceWithoutInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataDataSourceWithoutInfoRead,
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
			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Data.",
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

func dataSourceTencentCloudWedataDataSourceWithoutInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_wedata_data_source_without_info.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
		data    []*wedata.DataSourceInfo
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
		result, e := service.DescribeWedataDataSourceWithoutInfoByFilter(ctx, paramMap)
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
	tmpList := make([]map[string]interface{}, 0, len(data))

	if data != nil {
		for _, dataSourceInfo := range data {
			dataSourceInfoMap := map[string]interface{}{}

			if dataSourceInfo.DatabaseName != nil {
				dataSourceInfoMap["database_name"] = dataSourceInfo.DatabaseName
			}

			if dataSourceInfo.Description != nil {
				dataSourceInfoMap["description"] = dataSourceInfo.Description
			}

			if dataSourceInfo.ID != nil {
				dataSourceInfoMap["id"] = dataSourceInfo.ID
			}

			if dataSourceInfo.Instance != nil {
				dataSourceInfoMap["instance"] = dataSourceInfo.Instance
			}

			if dataSourceInfo.Name != nil {
				dataSourceInfoMap["name"] = dataSourceInfo.Name
			}

			if dataSourceInfo.Region != nil {
				dataSourceInfoMap["region"] = dataSourceInfo.Region
			}

			if dataSourceInfo.Type != nil {
				dataSourceInfoMap["type"] = dataSourceInfo.Type
			}

			if dataSourceInfo.ClusterId != nil {
				dataSourceInfoMap["cluster_id"] = dataSourceInfo.ClusterId
			}

			if dataSourceInfo.AppId != nil {
				dataSourceInfoMap["app_id"] = dataSourceInfo.AppId
			}

			if dataSourceInfo.BizParams != nil {
				dataSourceInfoMap["biz_params"] = dataSourceInfo.BizParams
			}

			if dataSourceInfo.Category != nil {
				dataSourceInfoMap["category"] = dataSourceInfo.Category
			}

			if dataSourceInfo.Display != nil {
				dataSourceInfoMap["display"] = dataSourceInfo.Display
			}

			if dataSourceInfo.OwnerAccount != nil {
				dataSourceInfoMap["owner_account"] = dataSourceInfo.OwnerAccount
			}

			if dataSourceInfo.Params != nil {
				dataSourceInfoMap["params"] = dataSourceInfo.Params
			}

			if dataSourceInfo.Status != nil {
				dataSourceInfoMap["status"] = dataSourceInfo.Status
			}

			if dataSourceInfo.OwnerAccountName != nil {
				dataSourceInfoMap["owner_account_name"] = dataSourceInfo.OwnerAccountName
			}

			if dataSourceInfo.ClusterName != nil {
				dataSourceInfoMap["cluster_name"] = dataSourceInfo.ClusterName
			}

			if dataSourceInfo.OwnerProjectId != nil {
				dataSourceInfoMap["owner_project_id"] = dataSourceInfo.OwnerProjectId
			}

			if dataSourceInfo.OwnerProjectName != nil {
				dataSourceInfoMap["owner_project_name"] = dataSourceInfo.OwnerProjectName
			}

			if dataSourceInfo.OwnerProjectIdent != nil {
				dataSourceInfoMap["owner_project_ident"] = dataSourceInfo.OwnerProjectIdent
			}

			if dataSourceInfo.AuthorityProjectName != nil {
				dataSourceInfoMap["authority_project_name"] = dataSourceInfo.AuthorityProjectName
			}

			if dataSourceInfo.AuthorityUserName != nil {
				dataSourceInfoMap["authority_user_name"] = dataSourceInfo.AuthorityUserName
			}

			if dataSourceInfo.Edit != nil {
				dataSourceInfoMap["edit"] = dataSourceInfo.Edit
			}

			if dataSourceInfo.Author != nil {
				dataSourceInfoMap["author"] = dataSourceInfo.Author
			}

			if dataSourceInfo.Deliver != nil {
				dataSourceInfoMap["deliver"] = dataSourceInfo.Deliver
			}

			if dataSourceInfo.DataSourceStatus != nil {
				dataSourceInfoMap["data_source_status"] = dataSourceInfo.DataSourceStatus
			}

			if dataSourceInfo.CreateTime != nil {
				dataSourceInfoMap["create_time"] = dataSourceInfo.CreateTime
			}

			if dataSourceInfo.ParamsString != nil {
				dataSourceInfoMap["params_string"] = dataSourceInfo.ParamsString
			}

			if dataSourceInfo.BizParamsString != nil {
				dataSourceInfoMap["biz_params_string"] = dataSourceInfo.BizParamsString
			}

			if dataSourceInfo.ModifiedTime != nil {
				dataSourceInfoMap["modified_time"] = dataSourceInfo.ModifiedTime
			}

			if dataSourceInfo.ShowType != nil {
				dataSourceInfoMap["show_type"] = dataSourceInfo.ShowType
			}

			idInt := *dataSourceInfo.ID
			id := strconv.FormatUint(idInt, 10)
			ids = append(ids, id)
			tmpList = append(tmpList, dataSourceInfoMap)
		}

		_ = d.Set("data", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
