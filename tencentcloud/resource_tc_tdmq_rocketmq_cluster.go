/*
Provides a resource to create a tdmqRocketmq cluster

Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
  cluster_name = &lt;nil&gt;
  remark = &lt;nil&gt;
                  }
```

Import

tdmqRocketmq cluster can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_rocketmq_cluster.cluster cluster_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmqRocketmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTdmqRocketmqCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqRocketmqClusterCreate,
		Read:   resourceTencentCloudTdmqRocketmqClusterRead,
		Update: resourceTencentCloudTdmqRocketmqClusterUpdate,
		Delete: resourceTencentCloudTdmqRocketmqClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster name, which can contain 3-64 letters, digits, hyphens, and underscores.",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cluster description (up to 128 characters).",
			},

			"cluster_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"region": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Region information.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Creation time in milliseconds.",
			},

			"public_end_point": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Public network access address.",
			},

			"vpc_end_point": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "VPC access address.",
			},

			"support_namespace_endpoint": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether the namespace access point is supported.",
			},

			"vpcs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Vpc list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vpc ID.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID.",
						},
					},
				},
			},

			"is_vip": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether it is an exclusive instance.",
			},

			"rocket_m_q_flag": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Rocketmq cluster identification.",
			},
		},
	}
}

func resourceTencentCloudTdmqRocketmqClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_cluster.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tdmqRocketmq.NewCreateRocketMQClusterRequest()
		response  = tdmqRocketmq.NewCreateRocketMQClusterResponse()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_name"); ok {
		request.ClusterName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqRocketmqClient().CreateRocketMQCluster(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tdmqRocketmq cluster failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	return resourceTencentCloudTdmqRocketmqClusterRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_cluster.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Id()

	cluster, err := service.DescribeTdmqRocketmqClusterById(ctx, clusterId)
	if err != nil {
		return err
	}

	if cluster == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqRocketmqCluster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if cluster.ClusterName != nil {
		_ = d.Set("cluster_name", cluster.ClusterName)
	}

	if cluster.Remark != nil {
		_ = d.Set("remark", cluster.Remark)
	}

	if cluster.ClusterId != nil {
		_ = d.Set("cluster_id", cluster.ClusterId)
	}

	if cluster.Region != nil {
		_ = d.Set("region", cluster.Region)
	}

	if cluster.CreateTime != nil {
		_ = d.Set("create_time", cluster.CreateTime)
	}

	if cluster.PublicEndPoint != nil {
		_ = d.Set("public_end_point", cluster.PublicEndPoint)
	}

	if cluster.VpcEndPoint != nil {
		_ = d.Set("vpc_end_point", cluster.VpcEndPoint)
	}

	if cluster.SupportNamespaceEndpoint != nil {
		_ = d.Set("support_namespace_endpoint", cluster.SupportNamespaceEndpoint)
	}

	if cluster.Vpcs != nil {
		vpcsList := []interface{}{}
		for _, vpcs := range cluster.Vpcs {
			vpcsMap := map[string]interface{}{}

			if cluster.Vpcs.VpcId != nil {
				vpcsMap["vpc_id"] = cluster.Vpcs.VpcId
			}

			if cluster.Vpcs.SubnetId != nil {
				vpcsMap["subnet_id"] = cluster.Vpcs.SubnetId
			}

			vpcsList = append(vpcsList, vpcsMap)
		}

		_ = d.Set("vpcs", vpcsList)

	}

	if cluster.IsVip != nil {
		_ = d.Set("is_vip", cluster.IsVip)
	}

	if cluster.RocketMQFlag != nil {
		_ = d.Set("rocket_m_q_flag", cluster.RocketMQFlag)
	}

	return nil
}

func resourceTencentCloudTdmqRocketmqClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_cluster.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmqRocketmq.NewModifyRocketMQClusterRequest()

	clusterId := d.Id()

	request.ClusterId = &clusterId

	immutableArgs := []string{"cluster_name", "remark", "cluster_id", "region", "create_time", "public_end_point", "vpc_end_point", "support_namespace_endpoint", "vpcs", "is_vip", "rocket_m_q_flag"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("cluster_name") {
		if v, ok := d.GetOk("cluster_name"); ok {
			request.ClusterName = helper.String(v.(string))
		}
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqRocketmqClient().ModifyRocketMQCluster(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tdmqRocketmq cluster failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRocketmqClusterRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_cluster.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}
	clusterId := d.Id()

	if err := service.DeleteTdmqRocketmqClusterById(ctx, clusterId); err != nil {
		return err
	}

	return nil
}
