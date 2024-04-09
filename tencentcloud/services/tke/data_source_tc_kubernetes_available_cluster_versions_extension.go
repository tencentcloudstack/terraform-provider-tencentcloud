package tke

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudKubernetesAvailableClusterVersionsReadPostRequest0(ctx context.Context, req *tke.DescribeAvailableClusterVersionRequest, resp *tke.DescribeAvailableClusterVersionResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)

	var (
		versions   []*string
		clusters   []*tke.ClusterVersion
		clusterIds []string
		id         *string
		ids        []string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		id = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_ids"); ok {
		clusterIdsSet := v.(*schema.Set).List()
		ids = helper.InterfacesStrings(clusterIdsSet)
	}

	if resp != nil {
		versions = resp.Response.Versions
		clusters = resp.Response.Clusters
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

	if id != nil {
		clusterIds = []string{*id}
	} else {
		clusterIds = ids
	}

	d.SetId(helper.DataResourceIdsHash(clusterIds))
	return nil
}
