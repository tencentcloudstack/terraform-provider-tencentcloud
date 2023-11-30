package tencentcloud

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceTencentCloudKubernetesClusterLevels() *schema.Resource {
	return &schema.Resource{
		Read: datasourceTencentCloudKubernetesClusterLevelsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify cluster Id, if set will only query current cluster's available levels.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used for save result.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of level information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alias used for pass to cluster level arguments.",
						},
						"crd_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of CRDs.",
						},
						"config_map_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of ConfigMaps.",
						},
						"enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the current level enabled.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Level name.",
						},
						"node_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of nodes.",
						},
						"other_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of others.",
						},
						"pod_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of pods.",
						},
					},
				},
			},
		},
	}
}

func datasourceTencentCloudKubernetesClusterLevelsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("datasource.tencentcloud_kubernetes_cluster_levels.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TkeService{meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Get("cluster_id").(string)
	result, err := service.DescribeClusterLevelAttribute(ctx, clusterId)
	if err != nil {
		d.SetId("")
		return err
	}
	d.SetId(fmt.Sprintf("cluster-level-%s__%d", clusterId, rand.Intn(4)))

	list := make([]interface{}, 0, len(result))

	for i := range result {
		item := result[i]
		level := map[string]interface{}{
			"name":             item.Name,
			"alias":            item.Alias,
			"crd_count":        item.CRDCount,
			"config_map_count": item.ConfigMapCount,
			"enable":           item.Enable,
			"node_count":       item.NodeCount,
			"other_count":      item.OtherCount,
			"pod_count":        item.PodCount,
		}
		list = append(list, level)
	}

	_ = d.Set("list", list)

	if output, ok := d.GetOk("result_output_file"); ok {
		return writeToFile(output.(string), list)
	}

	return nil
}
