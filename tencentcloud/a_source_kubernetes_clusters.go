package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceTencentCloudKubernetesClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKubernetesClustersRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"cluster_name"},
				Optional:      true,
			},
			"cluster_name": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"cluster_id"},
				Optional:      true,
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_os": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_runtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_deploy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_ipvs": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cluster_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ignore_cluster_cidr_conflict": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"cluster_max_pod_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cluster_max_service_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cluster_node_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"worker_instances_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: TkeCvmState(),
							},
						},
					},
				},
			},
		},
	}

}
func dataSourceTencentCloudKubernetesClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_kubernetes_clusters.read")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := TkeService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	id := ""
	name := ""

	if v, ok := d.GetOk("cluster_id"); ok {
		id = v.(string)
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		name = v.(string)
	}

	infos, err := service.DescribeClusters(ctx, id, name)

	if err != nil && id=="" {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			infos, err = service.DescribeClusters(ctx, id, name)
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

	for _, info := range infos {
		var infoMap = map[string]interface{}{}
		infoMap["cluster_name"] = info.ClusterName
		infoMap["cluster_desc"] = info.ClusterDescription
		infoMap["cluster_os"] = info.ClusterOs
		infoMap["cluster_deploy_type"] = info.DeployType
		infoMap["cluster_version"] = info.ClusterVersion
		infoMap["cluster_ipvs"] = info.Ipvs
		infoMap["vpc_id"] = info.VpcId
		infoMap["project_id"] = info.ProjectId
		infoMap["cluster_cidr"] = info.ClusterCidr
		infoMap["ignore_cluster_cidr_conflict"] = info.IgnoreClusterCidrConflict
		infoMap["cluster_max_pod_num"] = info.MaxClusterServiceNum
		infoMap["cluster_max_service_num"] = info.MaxClusterServiceNum
		infoMap["cluster_node_num"] = info.ClusterNodeNum

		_, workers, err := service.DescribeClusterInstances(ctx, info.ClusterId)
		if err != nil {
			_, workers, err = service.DescribeClusterInstances(ctx, info.ClusterId)

		}
		if err != nil {
			log.Printf("[CRITAL]%s tencentcloud_kubernetes_clusters DescribeClusterInstances fail, reason:%s\n ", logId, err.Error())
			return  err
		} else {

			workerInstancesList := make([]map[string]interface{}, 0, len(workers))
			for _, cvm := range workers {
				tempMap := make(map[string]interface{})
				tempMap["instance_id"] = cvm.InstanceId
				tempMap["instance_role"] = cvm.InstanceRole
				tempMap["instance_state"] = cvm.InstanceState
				tempMap["failed_reason"] = cvm.FailedReason
				workerInstancesList = append(workerInstancesList, tempMap)
			}

			infoMap["worker_instances_list"] = workerInstancesList

		}

		list = append(list, infoMap)
	}

	d.SetId("KubernetesClusters" + name + id)
	err = d.Set("list", list)
	if err != nil {
		log.Printf("[CRITAL]%s provider set tencentcloud_kubernetes_clusters list fail, reason:%s\n ", logId, err.Error())
		return  err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err = writeToFile(output.(string), list); err != nil {
			return err
		}
	}
	return nil
}
