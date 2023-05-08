/*
Provides a resource to create a tdcpg cluster.

~> **NOTE:** This resource is still in internal testing. To experience its functions, you need to apply for a whitelist from Tencent Cloud.

Example Usage

```hcl
resource "tencentcloud_tdcpg_cluster" "cluster" {
  zone = "ap-guangzhou-3"
  master_user_password = ""
  cpu = 1
  memory = 1
  vpc_id = "vpc_id"
  subnet_id = "subnet_id"
  pay_mode = "POSTPAID_BY_HOUR"
  cluster_name = "cluster_name"
  db_version = "10.17"
  instance_count = 1
  period = 1
  project_id = 0
}

```
Import

tdcpg cluster can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdcpg_cluster.cluster cluster_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdcpg "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdcpg/v20211118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTdcpgCluster() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTdcpgClusterRead,
		Create: resourceTencentCloudTdcpgClusterCreate,
		Update: resourceTencentCloudTdcpgClusterUpdate,
		Delete: resourceTencentCloudTdcpgClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "available zone.",
			},

			"master_user_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "user password.",
			},

			"cpu": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "cpu cores.",
			},

			"memory": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "memory size.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "vpc id.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "subnet id.",
			},

			"pay_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "pay mode, the value is either PREPAID or POSTPAID_BY_HOUR.",
			},

			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "cluster name.",
			},

			"db_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "community version number, default to 10.17.",
			},

			"instance_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "instance count.",
			},

			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "purchase time, required when PayMode is PREPAID, the value range is 1~60, default to 1.",
			},

			"storage": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "max storage, the unit is GB.",
			},

			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "project id, default to 0, means default project.",
			},
		},
	}
}

func resourceTencentCloudTdcpgClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdcpg_cluster.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tdcpg.NewCreateClusterRequest()
		response  *tdcpg.CreateClusterResponse
		service   = TdcpgService{client: meta.(*TencentCloudClient).apiV3Conn}
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		clusterId string
		dealNames []*string
	)

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("master_user_password"); ok {
		request.MasterUserPassword = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cpu"); ok {
		request.CPU = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("memory"); ok {
		request.Memory = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pay_mode"); ok {
		request.PayMode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		request.ClusterName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_version"); ok {
		request.DBVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_count"); ok {
		request.InstanceCount = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("period"); ok {
		request.Period = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("storage"); ok {
		request.Storage = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdcpgClient().CreateCluster(request)
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
		log.Printf("[CRITICAL]%s create tdcpg cluster failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		dealNames = response.Response.DealNameSet
		resources, e := service.DescribeTdcpgResourceByDealName(ctx, dealNames)

		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s call api[%s] success, request body [%s], resources [%v]\n",
				logId, "DescribeTdcpgResourceByDealName", request.ToJsonString(), resources)
		}
		clusterId = *resources[0].ClusterId
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s query tdcpg cluster resource by deal name:[%v] failed, reason:%+v", logId, dealNames, err)
		return err
	}

	d.SetId(clusterId)
	return resourceTencentCloudTdcpgClusterRead(d, meta)
}

func resourceTencentCloudTdcpgClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdcpg_cluster.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId           = getLogId(contextNil)
		ctx             = context.WithValue(context.TODO(), logIdKey, logId)
		cluster         *tdcpg.Cluster
		clusterInstance *tdcpg.Instance
		service         = TdcpgService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	// query the cluster
	clusterId := d.Id()
	err := resource.Retry(5*readRetryTimeout, func() *resource.RetryError {

		result, err := service.DescribeTdcpgCluster(ctx, &clusterId)
		if err != nil {
			return retryError(err)
		}

		if result != nil && result.ClusterSet[0] != nil {
			currStatus := *result.ClusterSet[0].Status

			if currStatus == "running" {
				cluster = result.ClusterSet[0]
				return nil
			}

			if currStatus == "creating" || currStatus == "recovering" {
				return resource.RetryableError(fmt.Errorf("cluster[%s] status is still creating or recovering, retry...", clusterId))
			}
			return resource.NonRetryableError(fmt.Errorf("cluster[%s] status is invalid, exit!", clusterId))
		}
		return resource.RetryableError(fmt.Errorf("can not get cluster[%s] status, retry...", clusterId))
	})
	if err != nil {
		d.SetId("")
		return fmt.Errorf("resource `cluster` %s does not exist. error reason:[%s]", clusterId, err)
	}

	// query the instance of cluster
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instances, e := service.DescribeTdcpgInstancesByFilter(ctx, &clusterId, nil)
		if e != nil {
			return retryError(e)
		}

		if instances[0] != nil {
			clusterInstance = instances[0]
			return nil
		}
		return resource.RetryableError(fmt.Errorf("instances is nil, retry..."))
	})
	if err != nil {
		return err
	}

	// set attr from the instance of cluster
	if clusterInstance.CPU != nil {
		_ = d.Set("cpu", clusterInstance.CPU)
	}

	if clusterInstance.Memory != nil {
		_ = d.Set("memory", clusterInstance.Memory)
	}

	// set attr from endpoint
	if cluster.EndpointSet[0] != nil {
		endpoint := cluster.EndpointSet[0]

		if endpoint.VpcId != nil {
			_ = d.Set("vpc_id", endpoint.VpcId)
		}

		if endpoint.SubnetId != nil {
			_ = d.Set("subnet_id", endpoint.SubnetId)
		}
	}

	// set rest of attributes
	if cluster.Zone != nil {
		_ = d.Set("zone", cluster.Zone)
	}

	if cluster.PayPeriodEndTime != nil && cluster.CreateTime != nil && *cluster.PayMode == "PREPAID" {
		_ = d.Set("period", monthBetweenTwoDates(*cluster.CreateTime, *cluster.PayPeriodEndTime))
	}

	if cluster.StorageLimit != nil {
		_ = d.Set("storage", *cluster.StorageLimit)
	}

	if cluster.PayMode != nil {
		_ = d.Set("pay_mode", *cluster.PayMode)
	}

	if cluster.ClusterName != nil {
		_ = d.Set("cluster_name", cluster.ClusterName)
	}

	if cluster.DBVersion != nil {
		_ = d.Set("db_version", cluster.DBVersion)
	}

	if cluster.InstanceCount != nil {
		_ = d.Set("instance_count", cluster.InstanceCount)
	}

	if cluster.ProjectId != nil {
		_ = d.Set("project_id", cluster.ProjectId)
	}

	return nil
}

func resourceTencentCloudTdcpgClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdcpg_cluster.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdcpg.NewModifyClusterNameRequest()

	clusterId := d.Id()

	request.ClusterId = &clusterId

	if d.HasChange("zone") {
		return fmt.Errorf("`zone` do not support change now.")
	}

	if d.HasChange("master_user_password") {
		return fmt.Errorf("`master_user_password` do not support change now.")
	}

	if d.HasChange("cpu") {
		return fmt.Errorf("`cpu` do not support change now.")
	}

	if d.HasChange("memory") {
		return fmt.Errorf("`memory` do not support change now.")
	}

	if d.HasChange("vpc_id") {
		return fmt.Errorf("`vpc_id` do not support change now.")
	}

	if d.HasChange("subnet_id") {
		return fmt.Errorf("`subnet_id` do not support change now.")
	}

	if d.HasChange("pay_mode") {
		return fmt.Errorf("`pay_mode` do not support change now.")
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		request.ClusterName = helper.String(v.(string))
	}

	if d.HasChange("db_version") {
		return fmt.Errorf("`db_version` do not support change now.")
	}

	if d.HasChange("instance_count") {
		return fmt.Errorf("`instance_count` do not support change now.")
	}

	if d.HasChange("period") {
		return fmt.Errorf("`period` do not support change now.")
	}

	if d.HasChange("storage") {
		return fmt.Errorf("`storage` do not support change now.")
	}

	if d.HasChange("project_id") {
		return fmt.Errorf("`project_id` do not support change now.")
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdcpgClient().ModifyClusterName(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s create tdcpg cluster failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdcpgClusterRead(d, meta)
}

func resourceTencentCloudTdcpgClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdcpg_cluster.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdcpgService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Id()

	if err := service.DeleteTdcpgClusterById(ctx, &clusterId); err != nil {
		return err
	}

	return nil
}
