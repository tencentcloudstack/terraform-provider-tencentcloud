// Code generated by iacg; DO NOT EDIT.
package tke

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudKubernetesClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKubernetesClustersRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"cluster_name"},
				Description:   "ID of the cluster. Conflict with cluster_name, can not be set at the same time.",
			},

			"cluster_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"cluster_id"},
				Description:   "Name of the cluster. Conflict with cluster_id, can not be set at the same time.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the cluster.",
			},

			"kube_config_file_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path prefix of kube config. You can store KubeConfig in a specified directory by specifying this field, such as ~/.kube/k8s, then public network access will use ~/.kube/k8s-clusterID-kubeconfig naming, and intranet access will use ~/.kube /k8s-clusterID-kubeconfig-intranet naming. If this field is not set, the KubeConfig will not be exported.",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of kubernetes clusters. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of cluster.",
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
						"cdc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CDC ID.",
						},
						"cluster_os": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operating system of the cluster.",
						},
						"container_runtime": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Container runtime of the cluster.",
							Deprecated:  "It has been deprecated from version 1.18.1.",
						},
						"cluster_deploy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment type of the cluster.",
						},
						"cluster_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version of the cluster.",
						},
						"cluster_ipvs": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether ipvs is enabled.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vpc ID of the cluster.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID of the cluster.",
						},
						"cluster_cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A network address block of the cluster. Different from vpc cidr and cidr of other clusters within this VPC.",
						},
						"ignore_cluster_cidr_conflict": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether to ignore the cluster cidr conflict error.",
						},
						"cluster_max_pod_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum number of Pods per node in the cluster.",
						},
						"cluster_max_service_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum number of services in the cluster.",
						},
						"cluster_as_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether to enable cluster node auto scaler.",
						},
						"node_name_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node name type of cluster.",
						},
						"cluster_extra_args": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Customized parameters for master component.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kube_apiserver": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The customized parameters for kube-apiserver.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"kube_controller_manager": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The customized parameters for kube-controller-manager.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"kube_scheduler": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The customized parameters for kube-scheduler.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster network type.",
						},
						"is_non_static_ip_mode": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether non-static ip mode is enabled.",
						},
						"kube_proxy_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster kube-proxy mode.",
						},
						"vpc_cni_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Distinguish between shared network card multi-IP mode and independent network card mode.",
						},
						"service_cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network address block of the cluster.",
						},
						"eni_subnet_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Subnet IDs for cluster with VPC-CNI network mode.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"claim_expired_seconds": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The expired seconds to recycle ENI.",
						},
						"deletion_protection": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether cluster deletion protection is enabled.",
						},
						"cluster_node_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of nodes in the cluster.",
						},
						"worker_instances_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "An information list of cvm within the WORKER clusters. Each element contains the following attributes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the cvm.",
									},
									"instance_role": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Role of the cvm.",
									},
									"instance_state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "State of the cvm.",
									},
									"failed_reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Information of the cvm when it is failed.",
									},
									"lan_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "LAN IP of the cvm.",
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags of the cluster.",
						},
						"kube_config": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "Kubernetes config.",
						},
						"kube_config_intranet": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "Kubernetes config of private network.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User name of account.",
						},
						"password": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Password of account.",
						},
						"certification_authority": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The certificate used for access.",
						},
						"cluster_external_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "External network address to access.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name for access.",
						},
						"pgw_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Intranet address used for access.",
						},
						"security_policy": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Access policy.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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

func dataSourceTencentCloudKubernetesClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_kubernetes_clusters.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	var respData []*tkev20180525.Cluster
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesClustersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	if err := dataSourceTencentCloudKubernetesClustersReadPreHandleResponse0(ctx, paramMap, &respData); err != nil {
		return err
	}

	var ids []string
	clustersList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, clusters := range respData {
			clustersMap := map[string]interface{}{}

			var clusterID string
			if clusters.ClusterId != nil {
				clustersMap["cluster_id"] = clusters.ClusterId
				clusterID = *clusters.ClusterId
			}

			if clusters.ClusterName != nil {
				clustersMap["cluster_name"] = clusters.ClusterName
			}

			if clusters.ClusterDescription != nil {
				clustersMap["cluster_desc"] = clusters.ClusterDescription
			}

			if clusters.CdcId != nil {
				clustersMap["cdc_id"] = clusters.CdcId
			}

			if clusters.ClusterOs != nil {
				clustersMap["cluster_os"] = clusters.ClusterOs
			}

			if clusters.ClusterType != nil {
				clustersMap["cluster_deploy_type"] = clusters.ClusterType
			}

			if clusters.ClusterVersion != nil {
				clustersMap["cluster_version"] = clusters.ClusterVersion
			}

			if clusters.ClusterNetworkSettings != nil {
				if clusters.ClusterNetworkSettings.Ipvs != nil {
					clustersMap["cluster_ipvs"] = clusters.ClusterNetworkSettings.Ipvs
				}

				if clusters.ClusterNetworkSettings.KubeProxyMode != nil {
					clustersMap["kube_proxy_mode"] = clusters.ClusterNetworkSettings.KubeProxyMode
				}

				if clusters.ClusterNetworkSettings.VpcId != nil {
					clustersMap["vpc_id"] = clusters.ClusterNetworkSettings.VpcId
				}

				if clusters.ClusterNetworkSettings.ClusterCIDR != nil {
					clustersMap["cluster_cidr"] = clusters.ClusterNetworkSettings.ClusterCIDR
				}

				if clusters.ClusterNetworkSettings.IgnoreClusterCIDRConflict != nil {
					clustersMap["ignore_cluster_cidr_conflict"] = clusters.ClusterNetworkSettings.IgnoreClusterCIDRConflict
				}

				if clusters.ClusterNetworkSettings.MaxNodePodNum != nil {
					clustersMap["cluster_max_pod_num"] = clusters.ClusterNetworkSettings.MaxNodePodNum
				}

				if clusters.ClusterNetworkSettings.MaxClusterServiceNum != nil {
					clustersMap["cluster_max_service_num"] = clusters.ClusterNetworkSettings.MaxClusterServiceNum
				}

				if clusters.ClusterNetworkSettings.ServiceCIDR != nil {
					clustersMap["service_cidr"] = clusters.ClusterNetworkSettings.ServiceCIDR
				}

				if clusters.ClusterNetworkSettings.Subnets != nil {
					clustersMap["eni_subnet_ids"] = clusters.ClusterNetworkSettings.Subnets
				}

			}

			if clusters.DeletionProtection != nil {
				clustersMap["deletion_protection"] = clusters.DeletionProtection
			}

			if clusters.ProjectId != nil {
				clustersMap["project_id"] = clusters.ProjectId
			}

			if clusters.ClusterNodeNum != nil {
				clustersMap["cluster_node_num"] = clusters.ClusterNodeNum
			}

			ids = append(ids, clusterID)
			clustersList = append(clustersList, clustersMap)
		}

		_ = d.Set("list", clustersList)
	}

	if err := dataSourceTencentCloudKubernetesClustersReadPostHandleResponse0(ctx, paramMap, &respData); err != nil {
		return err
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataSourceTencentCloudKubernetesClustersReadOutputContent(ctx)); e != nil {
			return e
		}
	}

	return nil
}
