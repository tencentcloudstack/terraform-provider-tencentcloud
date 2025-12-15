package bh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudBhSourceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudBhSourceTypesRead,
		Schema: map[string]*schema.Schema{
			"source_type_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Authentication source information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Account group source.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account group source type.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account group source name.",
						},
						"target": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Distinguish between ioa original and iam-mini.",
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

func dataSourceTencentCloudBhSourceTypesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_bh_source_types.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = BhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	var respData []*bhv20230418.SourceType
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeBhSourceTypesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	sourceTypeSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, sourceTypeSet := range respData {
			sourceTypeSetMap := map[string]interface{}{}
			if sourceTypeSet.Source != nil {
				sourceTypeSetMap["source"] = sourceTypeSet.Source
			}

			if sourceTypeSet.Type != nil {
				sourceTypeSetMap["type"] = sourceTypeSet.Type
			}

			if sourceTypeSet.Name != nil {
				sourceTypeSetMap["name"] = sourceTypeSet.Name
			}

			if sourceTypeSet.Target != nil {
				sourceTypeSetMap["target"] = sourceTypeSet.Target
			}

			sourceTypeSetList = append(sourceTypeSetList, sourceTypeSetMap)
		}

		_ = d.Set("source_type_set", sourceTypeSetList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), sourceTypeSetList); e != nil {
			return e
		}
	}

	return nil
}
