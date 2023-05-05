/*
Provides a resource to create a tdmqRocketmq namespace

Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
  cluster_name = "test_rocketmq"
  remark = "test recket mq"
}

resource "tencentcloud_tdmq_rocketmq_namespace" "namespace" {
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
  namespace_name = "test_namespace"
  ttl = 65000
  retention_time = 65000
  remark = "test namespace"
}
```
Import

tdmqRocketmq namespace can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rocketmq_namespace.namespace namespace_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmqRocketmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTdmqRocketmqNamespace() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTdmqRocketmqNamespaceRead,
		Create: resourceTencentCloudTdmqRocketmqNamespaceCreate,
		Update: resourceTencentCloudTdmqRocketmqNamespaceUpdate,
		Delete: resourceTencentCloudTdmqRocketmqNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster ID.",
			},

			"namespace_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Namespace name, which can contain 3-64 letters, digits, hyphens, and underscores.",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Retention time of unconsumed messages in milliseconds. Value range: 60 seconds-15 days.",
			},

			"retention_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Retention time of persisted messages in milliseconds.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remarks (up to 128 characters).",
			},

			"public_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public network access point address.",
			},

			"vpc_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC access point address.",
			},
		},
	}
}

func resourceTencentCloudTdmqRocketmqNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmqRocketmq_namespace.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	var (
		request       = tdmqRocketmq.NewCreateRocketMQNamespaceRequest()
		namespaceName string
		clusterId     string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_name"); ok {
		namespaceName = v.(string)
		request.NamespaceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ttl"); ok {
		request.Ttl = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("retention_time"); ok {
		request.RetentionTime = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("remark"); ok {

		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateRocketMQNamespace(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmqRocketmq namespace failed, reason:%+v", logId, err)
		return err
	}
	d.SetId(clusterId + FILED_SP + namespaceName)
	return resourceTencentCloudTdmqRocketmqNamespaceRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmqRocketmq_namespace.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	namespaceName := idSplit[1]

	namespaceList, err := service.DescribeTdmqRocketmqNamespace(ctx, namespaceName, clusterId)

	if err != nil || len(namespaceList) == 0 {
		return err
	}
	namespace := namespaceList[0]
	if namespace == nil {
		d.SetId("")
		return fmt.Errorf("resource `namespace` %s does not exist", namespaceName)
	}

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("namespace_name", *namespace.NamespaceId)
	_ = d.Set("ttl", namespace.Ttl)
	_ = d.Set("retention_time", namespace.RetentionTime)
	_ = d.Set("remark", namespace.Remark)
	_ = d.Set("public_endpoint", namespace.PublicEndpoint)
	_ = d.Set("vpc_endpoint", namespace.VpcEndpoint)

	return nil
}

func resourceTencentCloudTdmqRocketmqNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmqRocketmq_namespace.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmqRocketmq.NewModifyRocketMQNamespaceRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	namespaceName := idSplit[1]

	request.NamespaceId = &namespaceName
	request.ClusterId = &clusterId

	if d.HasChange("cluster_id") {

		return fmt.Errorf("`cluster_id` do not support change now.")

	}

	if d.HasChange("namespace_name") {

		return fmt.Errorf("`namespace_name` do not support change now.")

	}

	if d.HasChange("ttl") {
		if v, ok := d.GetOk("ttl"); ok {
			request.Ttl = helper.IntUint64(v.(int))
		}

	}

	if d.HasChange("retention_time") {
		if v, ok := d.GetOk("retention_time"); ok {
			request.RetentionTime = helper.IntUint64(v.(int))
		}

	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}

	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().ModifyRocketMQNamespace(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmqRocketmq namespace failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRocketmqNamespaceRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmqRocketmq_namespace.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	namespaceName := idSplit[1]

	if err := service.DeleteTdmqRocketmqNamespaceById(ctx, namespaceName, clusterId); err != nil {
		return err
	}

	return nil
}
