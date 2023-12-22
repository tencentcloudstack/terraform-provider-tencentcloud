package scf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudScfFunctionVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfFunctionVersionsRead,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function Name.",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The namespace where the function locates.",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "It specifies whether to return the results in ascending or descending order. The value is `ASC` or `DESC`.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "It specifies the sorting order of the results according to a specified field, such as `AddTime`, `ModTime`.",
			},

			"versions": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Function version listNote: This field may return null, indicating that no valid values is found.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Function version name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version descriptionNote: This field may return null, indicating that no valid values is found.",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation timeNote: This field may return null, indicating that no valid value was found.",
						},
						"mod_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update timeNote: This field may return null, indicating that no valid value was found.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version statusNote: this field may return `null`, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudScfFunctionVersionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_scf_function_versions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("function_name"); ok {
		paramMap["FunctionName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	service := ScfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var versions []*scf.FunctionVersion

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfFunctionVersionsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		versions = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(versions))
	tmpList := make([]map[string]interface{}, 0, len(versions))

	if versions != nil {
		for _, functionVersion := range versions {
			functionVersionMap := map[string]interface{}{}

			if functionVersion.Version != nil {
				functionVersionMap["version"] = functionVersion.Version
			}

			if functionVersion.Description != nil {
				functionVersionMap["description"] = functionVersion.Description
			}

			if functionVersion.AddTime != nil {
				functionVersionMap["add_time"] = functionVersion.AddTime
			}

			if functionVersion.ModTime != nil {
				functionVersionMap["mod_time"] = functionVersion.ModTime
			}

			if functionVersion.Status != nil {
				functionVersionMap["status"] = functionVersion.Status
			}

			ids = append(ids, *functionVersion.Version)
			tmpList = append(tmpList, functionVersionMap)
		}

		_ = d.Set("versions", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
