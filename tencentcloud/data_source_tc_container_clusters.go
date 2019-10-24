package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
)

func dataSourceTencentCloudContainerClusters() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source has been deprecated in Terraform TencentCloud provider version 1.16.0. Please use 'tencentcloud_kubernetes_clusters' instead.",
		Read:               dataSourceTencentCloudContainerClustersRead,

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
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_certification_authority": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_cluster_external_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kubernetes_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nodes_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"nodes_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_mem": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudContainerClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_container_clusters.read")()

	logId := getLogId(contextNil)

	request := tke.NewDescribeClustersRequest()
	if clusterId, ok := d.GetOkExists("cluster_id"); ok {
		request.ClusterIds = []*string{common.StringPtr(clusterId.(string))}
	}

	if limit, ok := d.GetOkExists("limit"); ok {
		request.Limit = common.Int64Ptr(limit.(int64))
	}

	var response *tke.DescribeClustersResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTkeClient().DescribeClusters(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s DescribeClusters failed, reason:%s\n ", logId, err.Error())
		return err
	}

	ids := make([]string, 0, *response.Response.TotalCount)
	clustersList := make([]map[string]interface{}, 0, *response.Response.TotalCount)
	for _, cluster := range response.Response.Clusters {
		ids = append(ids, *cluster.ClusterId)

		clusterInfo := make(map[string]interface{}, 1)
		clusterInfo["cluster_id"] = *cluster.ClusterId
		clusterInfo["cluster_name"] = *cluster.ClusterName
		clusterInfo["description"] = *cluster.ClusterDescription
		clusterInfo["kubernetes_version"] = *cluster.ClusterVersion
		clusterInfo["nodes_num"] = *cluster.ClusterNodeNum
		clusterInfo["nodes_status"] = *cluster.ClusterStatus
		clusterInfo["total_cpu"] = int64(0)
		clusterInfo["total_mem"] = int64(0)

		describeClusterInstancesreq := tke.NewDescribeClusterInstancesRequest()
		describeClusterInstancesreq.ClusterId = cluster.ClusterId
		describeClusterInstancesreq.Limit = common.Int64Ptr(100)
		var describeClusterInstancesResponse *tke.DescribeClusterInstancesResponse
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseTkeClient().DescribeClusterInstances(describeClusterInstancesreq)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, describeClusterInstancesreq.GetAction(), describeClusterInstancesreq.ToJsonString(), e.Error())
				return retryError(e)
			}
			describeClusterInstancesResponse = result
			return nil
		})
		if err != nil {
			continue
		}

		instanceIds := []*string{}
		for _, v := range describeClusterInstancesResponse.Response.InstanceSet {
			instanceIds = append(instanceIds, v.InstanceId)
		}

		if len(instanceIds) > 0 {
			describeInstancesreq := cvm.NewDescribeInstancesRequest()
			describeInstancesreq.InstanceIds = instanceIds
			var describeInstancesResponse *cvm.DescribeInstancesResponse
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().DescribeInstances(describeInstancesreq)
				if e != nil {
					log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
						logId, describeInstancesreq.GetAction(), describeInstancesreq.ToJsonString(), e.Error())
					return retryError(e)
				}
				describeInstancesResponse = result
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s DescribeInstances failed, reason:%s\n ", logId, err.Error())
				return err
			}

			for _, v := range describeInstancesResponse.Response.InstanceSet {
				clusterInfo["total_cpu"] = clusterInfo["total_cpu"].(int64) + *v.CPU
				clusterInfo["total_mem"] = clusterInfo["total_mem"].(int64) + *v.Memory
			}
		}

		describeClusterSecurityreq := tke.NewDescribeClusterSecurityRequest()
		describeClusterSecurityreq.ClusterId = cluster.ClusterId
		var securityResponse *tke.DescribeClusterSecurityResponse
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseTkeClient().DescribeClusterSecurity(describeClusterSecurityreq)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, describeClusterSecurityreq.GetAction(), describeClusterSecurityreq.ToJsonString(), e.Error())
				return retryError(e)
			}
			securityResponse = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s DescribeClusterSecurity failed, reason:%s\n ", logId, err.Error())
			return err
		}

		clusterInfo["security_certification_authority"] = *securityResponse.Response.CertificationAuthority
		clusterInfo["security_cluster_external_endpoint"] = *securityResponse.Response.ClusterExternalEndpoint
		clusterInfo["security_username"] = *securityResponse.Response.UserName
		clusterInfo["security_password"] = *securityResponse.Response.Password
		clustersList = append(clustersList, clusterInfo)
	}

	d.SetId(dataResourceIdsHash(ids))
	d.Set("clusters", clustersList)
	d.Set("total_count", *response.Response.TotalCount)

	return nil
}
