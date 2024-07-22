// Code generated by iacg; DO NOT EDIT.
package tke

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudKubernetesClusterCommonNames() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKubernetesClusterCommonNamesRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cluster ID.",
			},

			"subaccount_uins": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of sub-account. Up to 50 sub-accounts can be passed in at a time.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"role_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of Role ID. Up to 50 sub-accounts can be passed in at a time.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of the CommonName in the certificate of the client corresponding to the sub-account UIN.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subaccount_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User UIN.",
						},
						"common_names": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CommonName in the certificate of the client corresponding to the sub-account.",
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

func dataSourceTencentCloudKubernetesClusterCommonNamesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_kubernetes_cluster_common_names.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		clusterId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	var respData []*tke.CommonName
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesClusterCommonNamesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	cns := make([]string, 0, len(respData))
	commonNamesList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, commonNames := range respData {
			commonNamesMap := map[string]interface{}{}

			if commonNames.SubaccountUin != nil {
				commonNamesMap["subaccount_uin"] = commonNames.SubaccountUin
			}

			var cN string
			if commonNames.CN != nil {
				commonNamesMap["common_names"] = commonNames.CN
				cN = *commonNames.CN
			}

			cns = append(cns, cN)
			commonNamesList = append(commonNamesList, commonNamesMap)
		}

		_ = d.Set("list", commonNamesList)
	}

	d.SetId(strings.Join([]string{clusterId, helper.DataResourceIdsHash(cns)}, tccommon.FILED_SP))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), commonNamesList); e != nil {
			return e
		}
	}

	return nil
}
