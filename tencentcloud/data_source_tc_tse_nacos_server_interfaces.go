/*
Use this data source to query detailed information of tse nacos_server_interfaces

Example Usage

```hcl
data "tencentcloud_tse_nacos_server_interfaces" "nacos_server_interfaces" {
  instance_id = "ins-xxxxxx"
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTseNacosServerInterfaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseNacosServerInterfacesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Engine instance ID.",
			},

			"content": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Interface list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Interface nameNote: This field may return null, indicating that a valid value is not available.",
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

func dataSourceTencentCloudTseNacosServerInterfacesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tse_nacos_server_interfaces.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var content []*tse.NacosServerInterface

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTseNacosServerInterfacesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		content = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(content))
	tmpList := make([]map[string]interface{}, 0, len(content))

	if content != nil {
		for _, nacosServerInterface := range content {
			nacosServerInterfaceMap := map[string]interface{}{}

			if nacosServerInterface.Interface != nil {
				nacosServerInterfaceMap["interface"] = nacosServerInterface.Interface
			}

			ids = append(ids, *nacosServerInterface.InstanceId)
			tmpList = append(tmpList, nacosServerInterfaceMap)
		}

		_ = d.Set("content", tmpList)
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
