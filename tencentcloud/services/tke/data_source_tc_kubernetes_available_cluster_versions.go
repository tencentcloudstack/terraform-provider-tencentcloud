package tke

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudKubernetesAvailableClusterVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKubernetesAvailableClusterVersionsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cluster Id.",
			},

			"cluster_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "list of cluster IDs.",
			},

			"versions": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Upgradable cluster version number. Note: This field may return null, indicating that no valid value can be obtained.",
			},

			"clusters": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "cluster information. Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster ID.",
						},
						"versions": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "list of cluster major version numbers, for example 1.18.4.",
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

func dataSourceTencentCloudKubernetesAvailableClusterVersionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_kubernetes_available_cluster_versions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_ids"); ok {
		clusterIdsList := []*string{}
		clusterIdsSet := v.(*schema.Set).List()
		for i := range clusterIdsSet {
			clusterIds := clusterIdsSet[i].(string)
			clusterIdsList = append(clusterIdsList, helper.String(clusterIds))
		}
		paramMap["ClusterIds"] = clusterIdsList
	}

	var respData *tke.DescribeAvailableClusterVersionResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesAvailableClusterVersionsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	if respData.Versions != nil {
		_ = d.Set("versions", respData.Versions)
	}

	clustersList := make([]map[string]interface{}, 0, len(respData.Clusters))
	if respData.Clusters != nil {
		for _, clusters := range respData.Clusters {
			clustersMap := map[string]interface{}{}

			if clusters.ClusterId != nil {
				clustersMap["cluster_id"] = clusters.ClusterId
			}

			if clusters.Versions != nil {
				clustersMap["versions"] = clusters.Versions
			}

			clustersList = append(clustersList, clustersMap)
		}

		_ = d.Set("clusters", clustersList)
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataSourceTencentCloudKubernetesAvailableClusterVersionsReadOutputContent(ctx)); e != nil {
			return e
		}
	}

	return nil
}
