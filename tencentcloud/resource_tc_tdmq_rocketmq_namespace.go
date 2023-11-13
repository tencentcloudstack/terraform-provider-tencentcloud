/*
Provides a resource to create a tdmqRocketmq namespace

Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_namespace" "namespace" {
  cluster_id = &lt;nil&gt;
  namespace_id = &lt;nil&gt;
  ttl = &lt;nil&gt;
  retention_time = &lt;nil&gt;
  remark = &lt;nil&gt;
    }
```

Import

tdmqRocketmq namespace can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_rocketmq_namespace.namespace namespace_id
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
	"strings"
)

func resourceTencentCloudTdmqRocketmqNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqRocketmqNamespaceCreate,
		Read:   resourceTencentCloudTdmqRocketmqNamespaceRead,
		Update: resourceTencentCloudTdmqRocketmqNamespaceUpdate,
		Delete: resourceTencentCloudTdmqRocketmqNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"namespace_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Namespace name, which can contain 3-64 letters, digits, hyphens, and underscores.",
			},

			"ttl": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Retention time of unconsumed messages in milliseconds. Value range: 60 seconds-15 days.",
			},

			"retention_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Retention time of persisted messages in milliseconds.",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Remarks (up to 128 characters).",
			},

			"public_endpoint": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Public network access point address.",
			},

			"vpc_endpoint": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "VPC access point address.",
			},
		},
	}
}

func resourceTencentCloudTdmqRocketmqNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_namespace.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = tdmqRocketmq.NewCreateRocketMQNamespaceRequest()
		response    = tdmqRocketmq.NewCreateRocketMQNamespaceResponse()
		namespaceId string
		clusterId   string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_id"); ok {
		namespaceId = v.(string)
		request.NamespaceId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("ttl"); ok {
		request.Ttl = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("retention_time"); ok {
		request.RetentionTime = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqRocketmqClient().CreateRocketMQNamespace(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tdmqRocketmq namespace failed, reason:%+v", logId, err)
		return err
	}

	namespaceId = *response.Response.NamespaceId
	d.SetId(strings.Join([]string{namespaceId, clusterId}, FILED_SP))

	return resourceTencentCloudTdmqRocketmqNamespaceRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_namespace.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	namespaceId := idSplit[0]
	clusterId := idSplit[1]

	namespace, err := service.DescribeTdmqRocketmqNamespaceById(ctx, namespaceId, clusterId)
	if err != nil {
		return err
	}

	if namespace == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqRocketmqNamespace` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if namespace.ClusterId != nil {
		_ = d.Set("cluster_id", namespace.ClusterId)
	}

	if namespace.NamespaceId != nil {
		_ = d.Set("namespace_id", namespace.NamespaceId)
	}

	if namespace.Ttl != nil {
		_ = d.Set("ttl", namespace.Ttl)
	}

	if namespace.RetentionTime != nil {
		_ = d.Set("retention_time", namespace.RetentionTime)
	}

	if namespace.Remark != nil {
		_ = d.Set("remark", namespace.Remark)
	}

	if namespace.PublicEndpoint != nil {
		_ = d.Set("public_endpoint", namespace.PublicEndpoint)
	}

	if namespace.VpcEndpoint != nil {
		_ = d.Set("vpc_endpoint", namespace.VpcEndpoint)
	}

	return nil
}

func resourceTencentCloudTdmqRocketmqNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_namespace.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmqRocketmq.NewModifyRocketMQNamespaceRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	namespaceId := idSplit[0]
	clusterId := idSplit[1]

	request.NamespaceId = &namespaceId
	request.ClusterId = &clusterId

	immutableArgs := []string{"cluster_id", "namespace_id", "ttl", "retention_time", "remark", "public_endpoint", "vpc_endpoint"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("ttl") {
		if v, ok := d.GetOkExists("ttl"); ok {
			request.Ttl = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("retention_time") {
		if v, ok := d.GetOkExists("retention_time"); ok {
			request.RetentionTime = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqRocketmqClient().ModifyRocketMQNamespace(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tdmqRocketmq namespace failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRocketmqNamespaceRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_namespace.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	namespaceId := idSplit[0]
	clusterId := idSplit[1]

	if err := service.DeleteTdmqRocketmqNamespaceById(ctx, namespaceId, clusterId); err != nil {
		return err
	}

	return nil
}
