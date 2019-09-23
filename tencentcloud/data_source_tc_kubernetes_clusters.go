/*
Use this data source to query detailed information of kubernetes clusters.

Example Usage

```hcl
data "tencentcloud_kubernetes_clusters"  "name" {
    cluster_name ="terraform"
}

data "tencentcloud_kubernetes_clusters"  "id" {
    cluster_id ="cls-godovr32"
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func tkeClusterInfo() map[string]*schema.Schema {
	schemaBody := map[string]*schema.Schema{
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
		"cluster_os": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Operating system of the cluster.",
		},
		"container_runtime": {
			Type:        schema.TypeString,
			Deprecated:  "It has been deprecated from version 1.18.1.",
			Computed:    true,
			Description: "Container runtime of the cluster.",
		},
		"cluster_deploy_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Deployment type of the cluster.",
		},
		"cluster_version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: " Version of the cluster.",
		},
		"cluster_ipvs": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: " Indicates whether ipvs is enabled.",
		},
		"vpc_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Vpc Id of the cluster.",
		},
		"project_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Project Id of the cluster.",
		},
		"cluster_cidr": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "A network address block of the cluster. Different from vpc cidr and cidr of other clusters within this vpc.",
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
		"cluster_node_num": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of nodes in the  cluster.",
		},
		"worker_instances_list": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "An information list of cvm within the WORKER clusters. Each element contains the following attributes.",
			Elem: &schema.Resource{
				Schema: tkeCvmState(),
			},
		},
	}

	for k, v := range tkeSecurityInfo() {
		schemaBody[k] = v
	}

	return schemaBody
}

func dataSourceTencentCloudKubernetesClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKubernetesClustersRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"cluster_name"},
				Description:   " ID of the cluster. Conflict with cluster_name, can not be set at the same time.",
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
				Description: "An information  list of kubernetes clusters . Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: tkeClusterInfo(),
				},
			},
		},
	}

}
func dataSourceTencentCloudKubernetesClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_kubernetes_clusters.read")()

	logId := getLogId(contextNil)
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

	if err != nil && id == "" {
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

	var emptyStrFunc = func(ptr *string) string {
		if ptr == nil {
			return ""
		} else {
			return *ptr
		}
	}

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
			return err
		}

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

		securityRet, err := service.DescribeClusterSecurity(ctx, info.ClusterId)
		if err != nil {
			securityRet, err = service.DescribeClusterSecurity(ctx, info.ClusterId)
		}

		if err != nil {
			log.Printf("[CRITAL]%s tencentcloud_kubernetes_clusters DescribeClusterSecurity fail, reason:%s\n ", logId, err.Error())
			return err
		}

		policies := make([]string, 0, len(securityRet.Response.SecurityPolicy))
		for _, v := range securityRet.Response.SecurityPolicy {
			policies = append(policies, *v)
		}

		infoMap["user_name"] = emptyStrFunc(securityRet.Response.UserName)
		infoMap["password"] = emptyStrFunc(securityRet.Response.Password)
		infoMap["certification_authority"] = emptyStrFunc(securityRet.Response.CertificationAuthority)
		infoMap["cluster_external_endpoint"] = emptyStrFunc(securityRet.Response.ClusterExternalEndpoint)
		infoMap["domain"] = emptyStrFunc(securityRet.Response.Domain)
		infoMap["pgw_endpoint"] = emptyStrFunc(securityRet.Response.PgwEndpoint)
		infoMap["security_policy"] = policies

		list = append(list, infoMap)
	}

	d.SetId("KubernetesClusters" + name + id)
	err = d.Set("list", list)
	if err != nil {
		log.Printf("[CRITAL]%s provider set tencentcloud_kubernetes_clusters list fail, reason:%s\n ", logId, err.Error())
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
