/*
Provides a resource to create a tdcpg cluster

Example Usage

```hcl
resource "tencentcloud_tdcpg_cluster" "cluster" {
  zone = &lt;nil&gt;
  master_user_password = &lt;nil&gt;
  c_p_u = &lt;nil&gt;
  memory = &lt;nil&gt;
  vpc_id = &lt;nil&gt;
  subnet_id = &lt;nil&gt;
  pay_mode = &lt;nil&gt;
  cluster_name = &lt;nil&gt;
  d_b_version = &lt;nil&gt;
  instance_count = &lt;nil&gt;
  period = &lt;nil&gt;
  storage = &lt;nil&gt;
  project_id = &lt;nil&gt;
}
```

Import

tdcpg cluster can be imported using the id, e.g.

```
terraform import tencentcloud_tdcpg_cluster.cluster cluster_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdcpg "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdcpg/v20211118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTdcpgCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdcpgClusterCreate,
		Read:   resourceTencentCloudTdcpgClusterRead,
		Update: resourceTencentCloudTdcpgClusterUpdate,
		Delete: resourceTencentCloudTdcpgClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Available zone.",
			},

			"master_user_password": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "User password.",
			},

			"c_p_u": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Cpu cores.",
			},

			"memory": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Memory size.",
			},

			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Vpc id.",
			},

			"subnet_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Subnet id.",
			},

			"pay_mode": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Pay mode, the value is either PREPAID or POSTPAID_BY_HOUR.",
			},

			"cluster_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cluster name.",
			},

			"d_b_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Community version number, default to 10.17.",
			},

			"instance_count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Instance count.",
			},

			"period": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Purchase time, required when PayMode is PREPAID, the value range is 1~60, default to 1.",
			},

			"storage": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Max storage, the unit is GB.",
			},

			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Project id, default to 0, means default project.",
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
		response  = tdcpg.NewCreateClusterResponse()
		clusterId string
	)
	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("master_user_password"); ok {
		request.MasterUserPassword = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("c_p_u"); ok {
		request.CPU = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("memory"); ok {
		request.Memory = helper.IntInt64(v.(int))
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

	if v, ok := d.GetOk("d_b_version"); ok {
		request.DBVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("instance_count"); ok {
		request.InstanceCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("storage"); ok {
		request.Storage = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdcpgClient().CreateCluster(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tdcpg cluster failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	return resourceTencentCloudTdcpgClusterRead(d, meta)
}

func resourceTencentCloudTdcpgClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdcpg_cluster.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdcpgService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Id()

	cluster, err := service.DescribeTdcpgClusterById(ctx, clusterId)
	if err != nil {
		return err
	}

	if cluster == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdcpgCluster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if cluster.Zone != nil {
		_ = d.Set("zone", cluster.Zone)
	}

	if cluster.MasterUserPassword != nil {
		_ = d.Set("master_user_password", cluster.MasterUserPassword)
	}

	if cluster.CPU != nil {
		_ = d.Set("c_p_u", cluster.CPU)
	}

	if cluster.Memory != nil {
		_ = d.Set("memory", cluster.Memory)
	}

	if cluster.VpcId != nil {
		_ = d.Set("vpc_id", cluster.VpcId)
	}

	if cluster.SubnetId != nil {
		_ = d.Set("subnet_id", cluster.SubnetId)
	}

	if cluster.PayMode != nil {
		_ = d.Set("pay_mode", cluster.PayMode)
	}

	if cluster.ClusterName != nil {
		_ = d.Set("cluster_name", cluster.ClusterName)
	}

	if cluster.DBVersion != nil {
		_ = d.Set("d_b_version", cluster.DBVersion)
	}

	if cluster.InstanceCount != nil {
		_ = d.Set("instance_count", cluster.InstanceCount)
	}

	if cluster.Period != nil {
		_ = d.Set("period", cluster.Period)
	}

	if cluster.Storage != nil {
		_ = d.Set("storage", cluster.Storage)
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

	immutableArgs := []string{"zone", "master_user_password", "c_p_u", "memory", "vpc_id", "subnet_id", "pay_mode", "cluster_name", "d_b_version", "instance_count", "period", "storage", "project_id"}

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdcpgClient().ModifyClusterName(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tdcpg cluster failed, reason:%+v", logId, err)
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

	if err := service.DeleteTdcpgClusterById(ctx, clusterId); err != nil {
		return err
	}

	return nil
}
