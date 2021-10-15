/*
Provide a resource to create a TDMQ instance.

Example Usage

```hcl
resource "tencentcloud_tdmq_instance" "foo" {
  cluster_name = "example"
  remark = "this is description."
}
```

Import

Tdmq instance can be imported, e.g.

```
$ terraform import tencentcloud_tdmq_instance.test tdmq_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func resourceTencentCloudTdmqInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqCreate,
		Read:   resourceTencentCloudTdmqRead,
		Update: resourceTencentCloudTdmqUpdate,
		Delete: resourceTencentCloudTdmqDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of tdmq cluster to be created.",
			},
			"bind_cluster_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The Dedicated Cluster Id.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the tdmq cluster.",
			},
		},
	}
}

func resourceTencentCloudTdmqCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_instance.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tdmqService := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		clusterName   string
		bindClusterId uint64
		remark        string
	)
	if temp, ok := d.GetOk("cluster_name"); ok {
		clusterName = temp.(string)
		if len(clusterName) < 1 {
			return fmt.Errorf("cluster_name should be not empty string")
		}
	}

	bindClusterId = uint64(d.Get("bind_cluster_id").(int))

	if temp, ok := d.GetOk("remark"); ok {
		remark = temp.(string)
	}

	clusterId, err := tdmqService.CreateTdmqInstance(ctx, clusterName, bindClusterId, remark)
	if err != nil {
		return err
	}
	d.SetId(clusterId)

	return resourceTencentCloudTdmqRead(d, meta)
}

func resourceTencentCloudTdmqRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	tdmqService := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		info, has, e := tdmqService.DescribeTdmqInstanceById(ctx, id)
		if e != nil {
			return retryError(e)
		}
		if !has {
			d.SetId("")
			return nil
		}

		_ = d.Set("cluster_name", info.ClusterName)
		_ = d.Set("remark", info.Remark)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudTdmqUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_instance.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		clusterName string
		remark      string
	)
	old, now := d.GetChange("cluster_name")
	if d.HasChange("cluster_name") {
		clusterName = now.(string)
	} else {
		clusterName = old.(string)
	}

	old, now = d.GetChange("remark")
	if d.HasChange("remark") {
		remark = now.(string)
	} else {
		remark = old.(string)
	}

	d.Partial(true)

	if err := service.ModifyTdmqInstanceAttribute(ctx, id, clusterName, remark); err != nil {
		return err
	}
	d.SetPartial("cluster_name")
	d.SetPartial("remark")

	d.Partial(false)
	return resourceTencentCloudTdmqRead(d, meta)
}

func resourceTencentCloudTdmqDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_instance.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if err := service.DeleteTdmqInstance(ctx, d.Id()); err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == VPCNotFound {
					return nil
				}
			}
			return resource.RetryableError(err)
		}
		return nil
	})

	return err
}
