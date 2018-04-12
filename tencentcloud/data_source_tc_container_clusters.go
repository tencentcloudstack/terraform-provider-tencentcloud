package tencentcloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	ccs "github.com/zqfan/tencentcloud-sdk-go/services/ccs/unversioned"
)

func dataSourceTencentCloudContainerClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudContainerClustersRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_certification_authority": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_cluster_external_endpoint": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_username": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_password": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"kubernetes_version": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"nodes_num": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"nodes_status": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_cpu": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_mem": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudContainerClustersRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).ccsConn
	describeClustersReq := ccs.NewDescribeClusterRequest()
	if clusterId, ok := d.GetOkExists("cluster_id"); ok {
		describeClustersReq.ClusterIds = []*string{common.StringPtr(clusterId.(string))}
	}
	if limit, ok := d.GetOkExists("limit"); ok {
		describeClustersReq.Limit = common.IntPtr(limit.(int))
	}

	response, err := client.DescribeCluster(describeClustersReq)
	if err != nil {
		return err
	}

	if response.Code == nil {
		return fmt.Errorf("data_source_tencent_cloud_container_clusters got error, no code response")
	}

	if *response.Code != 0 {
		return fmt.Errorf("data_source_tencent_cloud_container_clusters got error, code %v , message %v", *response.Code, *response.CodeDesc)
	}

	id := fmt.Sprintf("%d", time.Now().Unix())

	clustersList := make([]map[string]interface{}, 0)
	for _, cluster := range response.Data.Clusters {
		clusterInfo := make(map[string]interface{}, 1)
		//basic info
		if cluster.ClusterId != nil {
			clusterInfo["cluster_id"] = *cluster.ClusterId
		}
		if cluster.Description != nil {
			clusterInfo["description"] = *cluster.Description
		}

		if cluster.K8sVersion != nil {
			clusterInfo["kubernetes_version"] = *cluster.K8sVersion
		}

		if cluster.NodeNum != nil {
			clusterInfo["nodes_num"] = *cluster.NodeNum
		}

		if cluster.NodeStatus != nil {
			clusterInfo["nodes_status"] = *cluster.NodeStatus
		}

		if cluster.TotalCPU != nil {
			clusterInfo["total_cpu"] = *cluster.TotalCPU
		}

		if cluster.TotalMem != nil {
			clusterInfo["total_mem"] = *cluster.TotalMem
		}

		//security info

		describeClusterSecurityInfoReq := ccs.NewDescribeClusterSecurityInfoRequest()
		describeClusterSecurityInfoReq.ClusterId = cluster.ClusterId

		securityResponse, err := client.DescribeClusterSecurityInfo(describeClusterSecurityInfoReq)

		if err != nil {
			continue
		}

		if securityResponse.Code == nil {
			continue
		}

		if *securityResponse.Code != 0 {
			continue
		}

		if securityResponse.Data.CertificationAuthority != nil {
			clusterInfo["security_certification_authority"] = *securityResponse.Data.CertificationAuthority
		}

		if securityResponse.Data.ClusterExternalEndpoint != nil {
			clusterInfo["security_cluster_external_endpoint"] = *securityResponse.Data.ClusterExternalEndpoint
		}

		if securityResponse.Data.Password != nil {
			clusterInfo["security_password"] = *securityResponse.Data.Password
		}

		if securityResponse.Data.UserName != nil {
			clusterInfo["security_username"] = *securityResponse.Data.UserName
		}

		clustersList = append(clustersList, clusterInfo)

	}
	d.Set("clusters", clustersList)
	d.SetId(id)
	if response.Data.TotalCount != nil {
		d.Set("total_count", *response.Data.TotalCount)
	} else {
		d.Set("total_count", 0)
	}

	return nil
}
