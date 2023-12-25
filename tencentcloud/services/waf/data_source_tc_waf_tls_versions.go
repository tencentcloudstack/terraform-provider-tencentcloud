package waf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWafTlsVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafTlsVersionsRead,
		Schema: map[string]*schema.Schema{
			"tls": {
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

func dataSourceTencentCloudWafTlsVersionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_waf_tls_versions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		tLS     []*waf.TLSVersion
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafTlsVersionsByFilter(ctx)
		if e != nil {
			return tccommon.RetryError(e)
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

			ids = append(ids, *tLSVersion.VersionName)
			tmpList = append(tmpList, tLSVersionMap)
		}

		_ = d.Set("tls", tmpList)
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
