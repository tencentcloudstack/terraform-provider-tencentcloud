/*
Use this data source to query elastic kubernetes cluster resource.

Example Usage

```
data "tencentcloud_eks_clusters" "foo" {
  cluster_id = "cls-xxxxxxxx"
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func dataSourceTencentCloudEKSClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEKSClustersRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"cluster_name"},
				Description:   "ID of the cluster. Conflict with cluster_name, can not be set at the same time.",
				Optional:      true,
			},
			"cluster_name": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"cluster_id"},
				Optional:      true,
				Description:   "Name of the cluster. Conflict with cluster_id, can not be set at the same time.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "EKS cluster list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the cluster.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the cluster.",
						},
						"cluster_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the cluster.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vpc id.",
						},
						"subnet_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Subnet id list.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"k8s_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "EKS cluster kubernetes version.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "EKS status.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the clusters.",
						},
						"service_subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet id of service.",
						},
						"dns_servers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of cluster custom DNS Server info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DNS Server domain. Empty indicates all domain.",
									},
									"servers": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of DNS Server IP address.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"need_delete_cbs": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether to delete CBS after EKS cluster remove.",
						},
						"enable_vpc_core_dns": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether to enable dns in user cluster, default value is `true`.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags of EKS cluster.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudEKSClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_eks_clusters.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := EksService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var (
		id   string
		name string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		id = v.(string)
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		name = v.(string)
	}

	tags := helper.GetTags(d, "tags")

	infos, err := service.DescribeEKSClusters(ctx, id, name)
	if err != nil && id == "" {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			infos, err = service.DescribeEKSClusters(ctx, id, name)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(infos))
	ids := make([]string, 0, len(infos))

LOOP:
	for _, info := range infos {
		if len(tags) > 0 {
			for k, v := range tags {
				if info.Tags[k] != v {
					continue LOOP
				}
			}
		}
		var infoMap = map[string]interface{}{
			"cluster_id":          info.ClusterId,
			"cluster_name":        info.ClusterName,
			"cluster_desc":        info.ClusterDesc,
			"vpc_id":              info.VpcId,
			"subnet_ids":          info.SubnetIds,
			"dns_servers":         info.DnsServers,
			"k8s_version":         info.K8SVersion,
			"status":              info.Status,
			"created_time":        info.CreatedTime,
			"service_subnet_id":   info.ServiceSubnetId,
			"need_delete_cbs":     info.NeedDeleteCbs,
			"enable_vpc_core_dns": info.EnableVpcCoreDNS,
		}

		list = append(list, infoMap)
		ids = append(ids, info.ClusterId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("list", list)
	if err != nil {
		log.Printf("[CRITAL]%s provider set tencentcloud_eks_clusters list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err = writeToFile(output.(string), list); err != nil {
			return err
		}
	}
	return nil
}
