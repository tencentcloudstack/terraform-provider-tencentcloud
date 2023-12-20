package cvm

import (
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCvmImportImageOs() *schema.Resource {
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
	defer tccommon.LogElapsed("data_source.tencentcloud_cvm_import_image_os.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		request  = cvm.NewDescribeImportImageOsRequest()
		response = cvm.NewDescribeImportImageOsResponse()
	)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		resule, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().DescribeImportImageOs(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = resule
		return nil
	})
	if err != nil {
		return err
	}

	if response == nil || response.Response == nil {
		d.SetId("")
		return fmt.Errorf("Response is null")
	}
	imageOsList := response.Response.ImportImageOsListSupported
	importImageOsVersionSet := response.Response.ImportImageOsVersionSet
	result := make(map[string]interface{})
	if imageOsList != nil {
		imageOsListMap := make(map[string]interface{})

		if len(imageOsList.Windows) != 0 {
			windowsList := make([]string, 0, len(imageOsList.Windows))
			for _, v := range imageOsList.Windows {
				windowsList = append(windowsList, *v)
			}
			imageOsListMap["windows"] = windowsList
		}

		if len(imageOsList.Linux) != 0 {
			linuxList := make([]string, 0, len(imageOsList.Linux))
			for _, v := range imageOsList.Linux {
				linuxList = append(linuxList, *v)
			}
			imageOsListMap["linux"] = linuxList
		}

		result["import_image_os_list_supported"] = imageOsListMap
		_ = d.Set("import_image_os_list_supported", []map[string]interface{}{imageOsListMap})
	}
	if importImageOsVersionSet != nil {
		tmpList := make([]map[string]interface{}, 0)
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

			tmpList = append(tmpList, osVersionMap)
		}
		result["import_image_os_version_set"] = tmpList
		_ = d.Set("import_image_os_version_set", tmpList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
