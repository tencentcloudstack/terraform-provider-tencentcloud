package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataGetTableColumns() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataGetTableColumnsRead,
		Schema: map[string]*schema.Schema{
			"table_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Table GUID.",
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Table column list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Field type.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Field name.",
						},
						"length": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Field length.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Field description.",
						},
						"position": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Field position.",
						},
						"is_partition": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether it is a partition field.",
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

func dataSourceTencentCloudWedataGetTableColumnsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_get_table_columns.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(nil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		tableGuid string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("table_guid"); ok {
		paramMap["TableGuid"] = helper.String(v.(string))
		tableGuid = v.(string)
	}

	var respData []*wedatav20250806.ColumnInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataGetTableColumnsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	dataList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, data := range respData {
			dataMap := map[string]interface{}{}
			if data.Type != nil {
				dataMap["type"] = data.Type
			}

			if data.Name != nil {
				dataMap["name"] = data.Name
			}

			if data.Length != nil {
				dataMap["length"] = data.Length
			}

			if data.Description != nil {
				dataMap["description"] = data.Description
			}

			if data.Position != nil {
				dataMap["position"] = data.Position
			}

			if data.IsPartition != nil {
				dataMap["is_partition"] = data.IsPartition
			}

			dataList = append(dataList, dataMap)
		}

		_ = d.Set("data", dataList)
	}

	d.SetId(tableGuid)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataList); e != nil {
			return e
		}
	}

	return nil
}
