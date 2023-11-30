package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudOrganizationOrgNode() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudOrganizationOrgNodeRead,
		Create: resourceTencentCloudOrganizationOrgNodeCreate,
		Update: resourceTencentCloudOrganizationOrgNodeUpdate,
		Delete: resourceTencentCloudOrganizationOrgNodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"parent_node_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Parent node ID.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Node name.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Notes.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Node creation time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
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
		response *organization.AddOrganizationNodeResponse
		nodeId   int64
	)

	if v, _ := d.GetOk("parent_node_id"); v != nil {
		request.ParentNodeId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("name"); ok {

		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {

		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOrganizationClient().AddOrganizationNode(request)
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
		log.Printf("[CRITAL]%s create organization orgNode failed, reason:%+v", logId, err)
		return err
	}

	nodeId = *response.Response.NodeId

	d.SetId(helper.Int64ToStr(nodeId))
	return resourceTencentCloudOrganizationOrgNodeRead(d, meta)
}

func resourceTencentCloudOrganizationOrgNodeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_node.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}

	orgNodeId := d.Id()

	orgNode, err := service.DescribeOrganizationOrgNode(ctx, orgNodeId)

	if err != nil {
		return err
	}

	if orgNode == nil {
		d.SetId("")
		return fmt.Errorf("resource `orgNode` %s does not exist", orgNodeId)
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

	request.NodeId = helper.StrToUint64Point(orgNodeId)

	if d.HasChange("parent_node_id") {
		return fmt.Errorf("`parent_node_id` do not support change now.")
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create organization orgNode failed, reason:%+v", logId, err)
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

	if err := service.DeleteOrganizationOrgNodeById(ctx, orgNodeId); err != nil {
		return err
	}

	return nil
}
