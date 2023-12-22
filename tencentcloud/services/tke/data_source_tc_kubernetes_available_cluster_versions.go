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

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		versions []*string
		clusters []*tke.ClusterVersion
		id       *string
		ids      []string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		id = helper.String(v.(string))
		paramMap["cluster_id"] = id
	}

	if v, ok := d.GetOk("cluster_ids"); ok {
		clusterIdsSet := v.(*schema.Set).List()
		ids = helper.InterfacesStrings(clusterIdsSet)
		paramMap["cluster_ids"] = helper.InterfacesStringsPoint(clusterIdsSet)
	}

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesAvailableClusterVersionsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result != nil {
			versions = result.Versions
			clusters = result.Clusters
		}
		return nil
	})
	if err != nil {
		return err
	}

	if versions != nil {
		_ = d.Set("versions", versions)
	}

	tmpList := make([]map[string]interface{}, 0, len(clusters))

	if clusters != nil {
		for _, clusterVersion := range clusters {
			clusterVersionMap := map[string]interface{}{}

			if clusterVersion.ClusterId != nil {
				clusterVersionMap["cluster_id"] = clusterVersion.ClusterId
			}

			if clusterVersion.Versions != nil {
				clusterVersionMap["versions"] = clusterVersion.Versions
			}

			tmpList = append(tmpList, clusterVersionMap)
		}

		_ = d.Set("clusters", tmpList)
	}

	var clusterIds []string
	if id != nil {
		clusterIds = []string{*id}
	} else {
		clusterIds = ids
	}

	d.SetId(helper.DataResourceIdsHash(clusterIds))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
