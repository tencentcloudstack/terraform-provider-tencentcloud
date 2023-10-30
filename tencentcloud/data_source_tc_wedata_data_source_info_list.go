/*
Use this data source to query detailed information of wedata data_source_info_list

Example Usage

```hcl
data "tencentcloud_wedata_data_source_info_list" "example" {
  project_id = "1927766435649077248"
  filters {
    name   = "Name"
    values = ["tf_example"]
  }

  order_fields {
    name      = "CreateTime"
    direction = "DESC"
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

func dataSourceTencentCloudWedataDataSourceInfoList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataDataSourceInfoListRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project ID.",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
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
							Description: "Filter values.",
						},
					},
				},
			},
			"order_fields": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
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
			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Datasource type.",
			},
			"datasource_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "DatasourceName.",
			},
			"datasource_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "DatasourceSet.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_names": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "DatabaseNames.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description.",
						},
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Id.",
						},
						"instance": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ClusterId.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Desc.",
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

func dataSourceTencentCloudWedataDataSourceInfoListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_wedata_data_source_info_list.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		service       = WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
		datasourceSet []*wedata.DatasourceBaseInfo
		projectId     string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("project_id"); ok {
		paramMap["ProjectId"] = helper.String(v.(string))
		projectId = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "filters"); ok {
		filter := wedata.Filter{}
		if v, ok := dMap["name"]; ok {
			filter.Name = helper.String(v.(string))
		}

		if v, ok := dMap["values"]; ok {
			valuesSet := v.(*schema.Set).List()
			filter.Values = helper.InterfacesStringsPoint(valuesSet)
		}

		paramMap["Filters"] = &filter
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "order_fields"); ok {
		orderField := wedata.OrderField{}
		if v, ok := dMap["name"]; ok {
			orderField.Name = helper.String(v.(string))
		}

		if v, ok := dMap["direction"]; ok {
			orderField.Direction = helper.String(v.(string))
		}

		paramMap["OrderFields"] = &orderField
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("datasource_name"); ok {
		paramMap["DatasourceName"] = helper.String(v.(string))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataDataSourceInfoListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		datasourceSet = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(datasourceSet))

	if datasourceSet != nil {
		for _, datasourceBaseInfo := range datasourceSet {
			datasourceBaseInfoMap := map[string]interface{}{}

			if datasourceBaseInfo.DatabaseNames != nil {
				nameList := make([]string, 0, len(datasourceBaseInfo.DatabaseNames))
				for _, databaseName := range datasourceBaseInfo.DatabaseNames {
					nameList = append(nameList, *databaseName)
				}

				datasourceBaseInfoMap["database_names"] = nameList
			}

			if datasourceBaseInfo.Description != nil {
				datasourceBaseInfoMap["description"] = datasourceBaseInfo.Description
			}

			if datasourceBaseInfo.ID != nil {
				datasourceBaseInfoMap["id"] = datasourceBaseInfo.ID
			}

			if datasourceBaseInfo.Instance != nil {
				datasourceBaseInfoMap["instance"] = datasourceBaseInfo.Instance
			}

			if datasourceBaseInfo.Name != nil {
				datasourceBaseInfoMap["name"] = datasourceBaseInfo.Name
			}

			if datasourceBaseInfo.Region != nil {
				datasourceBaseInfoMap["region"] = datasourceBaseInfo.Region
			}

			if datasourceBaseInfo.Type != nil {
				datasourceBaseInfoMap["type"] = datasourceBaseInfo.Type
			}

			if datasourceBaseInfo.ClusterId != nil {
				datasourceBaseInfoMap["cluster_id"] = datasourceBaseInfo.ClusterId
			}

			if datasourceBaseInfo.Version != nil {
				datasourceBaseInfoMap["version"] = datasourceBaseInfo.Version
			}

			tmpList = append(tmpList, datasourceBaseInfoMap)
		}

		_ = d.Set("datasource_set", tmpList)
	}

	d.SetId(projectId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
