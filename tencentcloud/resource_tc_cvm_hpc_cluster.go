/*
Provides a resource to create a cvm hpc_cluster

Example Usage

```hcl
resource "tencentcloud_cvm_hpc_cluster" "hpc_cluster" {
  zone = "ap-beijing-6"
  name = "terraform-test"
  remark = "create for test"
}
```

Import

cvm hpc_cluster can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_hpc_cluster.hpc_cluster hpc_cluster_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCvmHpcCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmHpcClusterCreate,
		Read:   resourceTencentCloudCvmHpcClusterRead,
		Update: resourceTencentCloudCvmHpcClusterUpdate,
		Delete: resourceTencentCloudCvmHpcClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Available zone.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Name of Hpc Cluster.",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Remark of Hpc Cluster.",
			},
		},
	}
}

func resourceTencentCloudCvmHpcClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_hpc_cluster.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = cvm.NewCreateHpcClusterRequest()
		response     = cvm.NewCreateHpcClusterResponse()
		hpcClusterId string
	)
	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().CreateHpcCluster(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cvm hpcCluster failed, reason:%+v", logId, err)
		return err
	}

	if len(response.Response.HpcClusterSet) < 1 {
		return fmt.Errorf("resource `tencentcloud_cvm_hpc_cluster` create failed.")
	}

	hpcClusterId = *response.Response.HpcClusterSet[0].HpcClusterId
	d.SetId(hpcClusterId)

	return resourceTencentCloudCvmHpcClusterRead(d, meta)
}

func resourceTencentCloudCvmHpcClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_hpc_cluster.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	hpcClusterId := d.Id()

	hpcCluster, err := service.DescribeCvmHpcClusterById(ctx, hpcClusterId)
	if err != nil {
		return err
	}

	if hpcCluster == nil {
		d.SetId("")
		return fmt.Errorf("resource `tencentcloud_cvm_hpc_cluster` %s does not exist", d.Id())
	}

	if hpcCluster.Zone != nil {
		_ = d.Set("zone", hpcCluster.Zone)
	}

	if hpcCluster.Name != nil {
		_ = d.Set("name", hpcCluster.Name)
	}

	if hpcCluster.Remark != nil {
		_ = d.Set("remark", hpcCluster.Remark)
	}

	return nil
}

func resourceTencentCloudCvmHpcClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_hpc_cluster.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cvm.NewModifyHpcClusterAttributeRequest()

	hpcClusterId := d.Id()

	request.HpcClusterId = &hpcClusterId

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ModifyHpcClusterAttribute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cvm hpcCluster failed, reason:%+v", logId, err)
		return nil
	}

	return resourceTencentCloudCvmHpcClusterRead(d, meta)
}

func resourceTencentCloudCvmHpcClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_hpc_cluster.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}
	hpcClusterId := d.Id()

	if err := service.DeleteCvmHpcClusterById(ctx, hpcClusterId); err != nil {
		return nil
	}

	return nil
}
