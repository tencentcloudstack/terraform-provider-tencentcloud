package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataListSchema() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataListSchemaRead,
		Schema: map[string]*schema.Schema{
			"catalog_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Catalog name.",
			},

			"datasource_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Data source ID.",
			},

			"database_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Database name.",
			},

			"keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Database schema search keyword.",
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Schema record list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"guid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Schema GUID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Schema name.",
						},
						"database_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database name.",
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

func dataSourceTencentCloudWedataListSchemaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_list_schema.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("catalog_name"); ok {
		paramMap["CatalogName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("datasource_id"); ok {
		paramMap["DatasourceId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("database_name"); ok {
		paramMap["DatabaseName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("keyword"); ok {
		paramMap["Keyword"] = helper.String(v.(string))
	}

	var respData []*wedatav20250806.SchemaInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataListSchemaByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	itemsList := make([]map[string]interface{}, 0, len(respData))
	for _, items := range respData {
		itemsMap := map[string]interface{}{}
		if items.Guid != nil {
			itemsMap["guid"] = items.Guid
		}

		if items.Name != nil {
			itemsMap["name"] = items.Name
		}

		if items.DatabaseName != nil {
			itemsMap["database_name"] = items.DatabaseName
		}

		itemsList = append(itemsList, itemsMap)
	}

	_ = d.Set("items", itemsList)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), itemsList); e != nil {
			return e
		}
	}

	return nil
}
