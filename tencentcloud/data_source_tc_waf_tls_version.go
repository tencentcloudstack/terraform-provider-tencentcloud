/*
Use this data source to query detailed information of waf tls_version

Example Usage

```hcl
data "tencentcloud_waf_tls_version" "tls_version" {
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudWafTlsVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafTlsVersionRead,
		Schema: map[string]*schema.Schema{
			"t_l_s": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "TLS key value.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TLS version IDNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"version_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tls version nameNote: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudWafTlsVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_waf_tls_version.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := WafService{client: meta.(*TencentCloudClient).apiV3Conn}

	var tLS []*waf.TLSVersion

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafTlsVersionByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		tLS = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(tLS))
	tmpList := make([]map[string]interface{}, 0, len(tLS))

	if tLS != nil {
		for _, tLSVersion := range tLS {
			tLSVersionMap := map[string]interface{}{}

			if tLSVersion.VersionId != nil {
				tLSVersionMap["version_id"] = tLSVersion.VersionId
			}

			if tLSVersion.VersionName != nil {
				tLSVersionMap["version_name"] = tLSVersion.VersionName
			}

			ids = append(ids, *tLSVersion.VersionId)
			tmpList = append(tmpList, tLSVersionMap)
		}

		_ = d.Set("t_l_s", tmpList)
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
