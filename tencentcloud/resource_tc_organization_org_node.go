/*
Provides a resource to create a organization org_node

Example Usage

```hcl
resource "tencentcloud_organization_org_node" "org_node" {
  node_id = &lt;nil&gt;
  parent_node_id = &lt;nil&gt;
  name = &lt;nil&gt;
  remark = &lt;nil&gt;
  create_time = &lt;nil&gt;
  update_time = &lt;nil&gt;
}
```

Import

organization org_node can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_node.org_node org_node_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudOrganizationOrgNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrgNodeCreate,
		Read:   resourceTencentCloudOrganizationOrgNodeRead,
		Update: resourceTencentCloudOrganizationOrgNodeUpdate,
		Delete: resourceTencentCloudOrganizationOrgNodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"node_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Node ID.",
			},

			"parent_node_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Parent node ID.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Node name.",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Notes.",
			},

			"create_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Node creation time.",
			},

			"update_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Node update time.",
			},
		},
	}
}

func resourceTencentCloudOrganizationOrgNodeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_node.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = organization.NewAddOrganizationNodeRequest()
		response = organization.NewAddOrganizationNodeResponse()
		nodeId   int
	)
	if v, ok := d.GetOkExists("node_id"); ok {
		nodeId = v.(int)
		request.NodeId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("parent_node_id"); ok {
		request.ParentNodeId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("create_time"); ok {
		request.CreateTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("update_time"); ok {
		request.UpdateTime = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOrganizationClient().AddOrganizationNode(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create organization orgNode failed, reason:%+v", logId, err)
		return err
	}

	nodeId = *response.Response.NodeId
	d.SetId(helper.Int64ToStr(int64(nodeId)))

	return resourceTencentCloudOrganizationOrgNodeRead(d, meta)
}

func resourceTencentCloudOrganizationOrgNodeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_node.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}

	orgNodeId := d.Id()

	orgNode, err := service.DescribeOrganizationOrgNodeById(ctx, nodeId)
	if err != nil {
		return err
	}

	if orgNode == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OrganizationOrgNode` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if orgNode.NodeId != nil {
		_ = d.Set("node_id", orgNode.NodeId)
	}

	if orgNode.ParentNodeId != nil {
		_ = d.Set("parent_node_id", orgNode.ParentNodeId)
	}

	if orgNode.Name != nil {
		_ = d.Set("name", orgNode.Name)
	}

	if orgNode.Remark != nil {
		_ = d.Set("remark", orgNode.Remark)
	}

	if orgNode.CreateTime != nil {
		_ = d.Set("create_time", orgNode.CreateTime)
	}

	if orgNode.UpdateTime != nil {
		_ = d.Set("update_time", orgNode.UpdateTime)
	}

	return nil
}

func resourceTencentCloudOrganizationOrgNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_node.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := organization.NewUpdateOrganizationNodeRequest()

	orgNodeId := d.Id()

	request.NodeId = &nodeId

	immutableArgs := []string{"node_id", "parent_node_id", "name", "remark", "create_time", "update_time"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOrganizationClient().UpdateOrganizationNode(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update organization orgNode failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudOrganizationOrgNodeRead(d, meta)
}

func resourceTencentCloudOrganizationOrgNodeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_node.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}
	orgNodeId := d.Id()

	if err := service.DeleteOrganizationOrgNodeById(ctx, nodeId); err != nil {
		return err
	}

	return nil
}
