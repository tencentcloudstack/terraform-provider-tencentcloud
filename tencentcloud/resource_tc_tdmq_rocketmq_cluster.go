package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmqRocketmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTdmqRocketmqCluster() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTdmqRocketmqClusterRead,
		Create: resourceTencentCloudTdmqRocketmqClusterCreate,
		Update: resourceTencentCloudTdmqRocketmqClusterUpdate,
		Delete: resourceTencentCloudTdmqRocketmqClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster name, which can contain 3-64 letters, digits, hyphens, and underscores.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cluster description (up to 128 characters).",
			},

			"cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster ID.",
			},

			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region information.",
			},

			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Creation time in milliseconds.",
			},

			"public_end_point": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public network access address.",
			},

			"vpc_end_point": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC access address.",
			},

			"support_namespace_endpoint": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the namespace access point is supported.",
			},

			"vpcs": {
				Type:        schema.TypeList,
				Computed:    true,
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
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether it is an exclusive instance.",
			},

			"rocket_m_q_flag": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Rocketmq cluster identification.",
			},
		},
	}
}

func resourceTencentCloudTdmqRocketmqClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmqRocketmq_cluster.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request   = tdmqRocketmq.NewCreateRocketMQClusterRequest()
		response  *tdmqRocketmq.CreateRocketMQClusterResponse
		clusterId string
	)

	if v, ok := d.GetOk("cluster_name"); ok {

		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {

		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateRocketMQCluster(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmqRocketmq cluster failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}
		_, innerErr := service.DescribeTdmqRocketmqCluster(ctx, clusterId)
		if innerErr != nil {
			return resource.RetryableError(innerErr)
		}
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(clusterId)
	return resourceTencentCloudTdmqRocketmqClusterRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmqRocketmq_cluster.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Id()

	cluster, err := service.DescribeTdmqRocketmqCluster(ctx, clusterId)

	if err != nil {
		return err
	}

	if cluster == nil {
		d.SetId("")
		return fmt.Errorf("resource `cluster` %s does not exist", clusterId)
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
			if vpcs.VpcId != nil {
				vpcsMap["vpc_id"] = vpcs.VpcId
			}
			if vpcs.SubnetId != nil {
				vpcsMap["subnet_id"] = vpcs.SubnetId
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
	defer logElapsed("resource.tencentcloud_tdmqRocketmq_cluster.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmqRocketmq.NewModifyRocketMQClusterRequest()

	clusterId := d.Id()

	request.ClusterId = &clusterId

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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().ModifyRocketMQCluster(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmqRocketmq cluster failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRocketmqClusterRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmqRocketmq_cluster.delete")()
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
