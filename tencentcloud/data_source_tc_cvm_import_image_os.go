/*
Use this data source to query detailed information of cvm import_image_os

Example Usage

```hcl
data "tencentcloud_cvm_import_image_os" "import_image_os" {
    }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCvmImportImageOs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCvmImportImageOsRead,
		Schema: map[string]*schema.Schema{
			"import_image_os_list_supported": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Supported operating system types of imported images.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"windows": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Supported Windows OS Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"linux": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Supported Linux OS Note: This field may return null, indicating that no valid values can be obtained.",
						},
					},
				},
			},

			"import_image_os_version_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Supported operating system versions of imported images.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"os_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operating system type.",
						},
						"os_versions": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Supported operating system versions.",
						},
						"architecture": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Supported operating system architecture.",
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

func dataSourceTencentCloudCvmImportImageOsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cvm_import_image_os.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	var importImageOsListSupported []*cvm.ImageOsList

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCvmImportImageOsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		importImageOsListSupported = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(importImageOsListSupported))
	if importImageOsListSupported != nil {
		imageOsListMap := map[string]interface{}{}

		if importImageOsListSupported.Windows != nil {
			imageOsListMap["windows"] = importImageOsListSupported.Windows
		}

		if importImageOsListSupported.Linux != nil {
			imageOsListMap["linux"] = importImageOsListSupported.Linux
		}

		ids = append(ids, *importImageOsListSupported.IdsHash)
		_ = d.Set("import_image_os_list_supported", imageOsListMap)
	}

	if importImageOsVersionSet != nil {
		for _, osVersion := range importImageOsVersionSet {
			osVersionMap := map[string]interface{}{}

			if osVersion.OsName != nil {
				osVersionMap["os_name"] = osVersion.OsName
			}

			if osVersion.OsVersions != nil {
				osVersionMap["os_versions"] = osVersion.OsVersions
			}

			if osVersion.Architecture != nil {
				osVersionMap["architecture"] = osVersion.Architecture
			}

			ids = append(ids, *osVersion.IdsHash)
			tmpList = append(tmpList, osVersionMap)
		}

		_ = d.Set("import_image_os_version_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), imageOsListMap); e != nil {
			return e
		}
	}
	return nil
}
