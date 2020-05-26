/*

Use this data source to query tcaplus clusters

Example Usage

```hcl
data "tencentcloud_tcaplus_clusters" "name" {
  cluster_name = "cluster"
}
data "tencentcloud_tcaplus_clusters" "id" {
  cluster_id = tencentcloud_tcaplus_cluster.test.id
}
data "tencentcloud_tcaplus_clusters" "idname" {
  cluster_id   = tencentcloud_tcaplus_cluster.test.id
  cluster_name = "cluster"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceTencentCloudTcaplusClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcaplusClustersRead,
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the tcaplus cluster to be query.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the tcaplus cluster to be query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of tcaplus cluster. Each element contains the following attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the tcaplus cluster.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the tcaplus cluster.",
						},
						"idl_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Idl type of the tcaplus cluster.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC id of the tcaplus cluster.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet id of the tcaplus cluster.",
						},
						"password": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Password of the tcaplus cluster.",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network type of the tcaplus cluster.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the tcaplus cluster.",
						},
						"password_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Password status of the tcaplus cluster.`unmodifiable` means:can not change password now,`modifiable` means:can change password now.",
						},
						"api_access_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access id of the tcaplus cluster.For TcaplusDB SDK connect.",
						},
						"api_access_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access ip of the tcaplus cluster.For TcaplusDB SDK connect.",
						},
						"api_access_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Access port of the tcaplus cluster.For TcaplusDB SDK connect.",
						},
						"old_password_expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "This field will display the old password expiration time,if password_status is `unmodifiable` means the old password has not yet expired, otherwise `-`.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTcaplusClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcaplus_clusters.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcaplusService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	clusterId := d.Get("cluster_id").(string)
	clusterName := d.Get("cluster_name").(string)

	clusters, err := service.DescribeClusters(ctx, clusterId, clusterName)
	if err != nil {
		clusters, err = service.DescribeClusters(ctx, clusterId, clusterName)
	}

	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(clusters))

	for _, cluster := range clusters {
		listItem := make(map[string]interface{})
		listItem["cluster_name"] = cluster.ClusterName
		listItem["cluster_id"] = cluster.ClusterId
		listItem["idl_type"] = cluster.IdlType
		listItem["vpc_id"] = cluster.VpcId
		listItem["subnet_id"] = cluster.SubnetId
		listItem["password"] = cluster.Password
		listItem["network_type"] = cluster.NetworkType
		listItem["create_time"] = cluster.CreatedTime
		listItem["password_status"] = cluster.PasswordStatus
		listItem["api_access_id"] = cluster.ApiAccessId
		listItem["api_access_ip"] = cluster.ApiAccessIp
		listItem["api_access_port"] = cluster.ApiAccessPort
		listItem["old_password_expire_time"] = cluster.OldPasswordExpireTime
		list = append(list, listItem)
	}

	d.SetId("cluster." + clusterId + "." + clusterName)
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}
	return nil

}
