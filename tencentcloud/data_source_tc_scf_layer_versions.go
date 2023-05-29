/*
Use this data source to query detailed information of scf layer_versions

Example Usage

```hcl
data "tencentcloud_scf_layer_versions" "layer_versions" {
  layer_name = "tf-test"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudScfLayerVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfLayerVersionsRead,
		Schema: map[string]*schema.Schema{
			"layer_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Layer name.",
			},

			"compatible_runtime": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Compatible runtimes.",
			},

			"layer_versions": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Layer version list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compatible_runtimes": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Runtime applicable to a versionNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version descriptionNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"license_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "License informationNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"layer_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Version number.",
						},
						"layer_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Layer name.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current status of specific layer version. For valid values, please see [here](https://intl.cloud.tencent.com/document/product/583/47175?from_cn_redirect=1#.E5.B1.82.EF.BC.88layer.EF.BC.89.E7.8A.B6.E6.80.81).",
						},
						"stamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "StampNote: This field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudScfLayerVersionsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_scf_layer_versions.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("layer_name"); ok {
		paramMap["LayerName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("compatible_runtime"); ok {
		compatibleRuntimeSet := v.(*schema.Set).List()
		paramMap["CompatibleRuntime"] = helper.InterfacesStringsPoint(compatibleRuntimeSet)
	}

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var layerVersions []*scf.LayerVersionInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfLayerVersions(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		layerVersions = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(layerVersions))
	tmpList := make([]map[string]interface{}, 0, len(layerVersions))

	if layerVersions != nil {
		for _, layerVersionInfo := range layerVersions {
			layerVersionInfoMap := map[string]interface{}{}

			if layerVersionInfo.CompatibleRuntimes != nil {
				layerVersionInfoMap["compatible_runtimes"] = layerVersionInfo.CompatibleRuntimes
			}

			if layerVersionInfo.AddTime != nil {
				layerVersionInfoMap["add_time"] = layerVersionInfo.AddTime
			}

			if layerVersionInfo.Description != nil {
				layerVersionInfoMap["description"] = layerVersionInfo.Description
			}

			if layerVersionInfo.LicenseInfo != nil {
				layerVersionInfoMap["license_info"] = layerVersionInfo.LicenseInfo
			}

			if layerVersionInfo.LayerVersion != nil {
				layerVersionInfoMap["layer_version"] = layerVersionInfo.LayerVersion
			}

			if layerVersionInfo.LayerName != nil {
				layerVersionInfoMap["layer_name"] = layerVersionInfo.LayerName
			}

			if layerVersionInfo.Status != nil {
				layerVersionInfoMap["status"] = layerVersionInfo.Status
			}

			if layerVersionInfo.Stamp != nil {
				layerVersionInfoMap["stamp"] = layerVersionInfo.Stamp
			}

			ids = append(ids, *layerVersionInfo.LayerName)
			tmpList = append(tmpList, layerVersionInfoMap)
		}

		_ = d.Set("layer_versions", tmpList)
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
