package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataListDatabase() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataListDatabaseRead,
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

			"keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Database name search keyword.",
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Database record list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"guid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database GUID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database name.",
						},
						"catalog_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database catalog.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database description.",
						},
						"location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database location.",
						},
						"storage_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Database storage size.",
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

func dataSourceTencentCloudWedataListDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_list_database.read")()
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

	if v, ok := d.GetOk("keyword"); ok {
		paramMap["Keyword"] = helper.String(v.(string))
	}

	var respData []*wedatav20250806.DatabaseInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataListDatabaseByFilter(ctx, paramMap)
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

		if items.CatalogName != nil {
			itemsMap["catalog_name"] = items.CatalogName
		}

		if items.Description != nil {
			itemsMap["description"] = items.Description
		}

		if items.Location != nil {
			itemsMap["location"] = items.Location
		}

		if items.StorageSize != nil {
			itemsMap["storage_size"] = items.StorageSize
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
