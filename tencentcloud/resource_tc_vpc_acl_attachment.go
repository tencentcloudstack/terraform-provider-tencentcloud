/*
Provide a resource to attach an existing subnet to Network ACL.

Example Usage

```hcl
data "tencentcloud_vpc_instances" "id_instances" {
}
resource "tencentcloud_vpc_acl" "foo" {
    vpc_id  = data.tencentcloud_vpc_instances.id_instances.instance_list.0.vpc_id
    name  	= "test_acl"
	ingress = [
		"ACCEPT#192.168.1.0/24#800#TCP",
		"ACCEPT#192.168.1.0/24#800-900#TCP",
	]
	egress = [
    	"ACCEPT#192.168.1.0/24#800#TCP",
    	"ACCEPT#192.168.1.0/24#800-900#TCP",
	]
}

resource "tencentcloud_vpc_acl_attachment" "attachment"{
		acl_id = tencentcloud_vpc_acl.foo.id
		subnet_ids = data.tencentcloud_vpc_instances.id_instances.instance_list[0].subnet_ids
}
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpcAclAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcAclAttachmentCreate,
		Read:   resourceTencentCloudVpcAclAttachmentRead,
		Delete: resourceTencentCloudVpcAclAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"acl_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the attached ACL.",
			},
			"subnet_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "ID list of the Subnet instance ID.",
			},
		},
	}
}

func resourceTencentCloudVpcAclAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_acl_attachment.create")()
	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		service   = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
		aclId     string
		subnetIds []string
		sub_id    string
	)

	if temp, ok := d.GetOk("acl_id"); ok {
		aclId = temp.(string)
		if len(aclId) < 1 {
			return fmt.Errorf("acl_id should be not empty string")
		}
	}
	if temp, ok := d.GetOk("subnet_ids"); ok {
		subnetIds = helper.InterfacesStrings(temp.([]interface{}))
	}

	err := service.AssociateAclSubnets(ctx, aclId, subnetIds)
	if err != nil {
		return err
	}

	for _, temp_id := range subnetIds {
		sub_id = sub_id + "#" + temp_id
	}
	d.SetId(aclId + "#" + sub_id)

	aclAttachmentId := d.Id()
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		e := service.DescribeByAclId(ctx, aclAttachmentId)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read acl attachment failed, reason:%s\n", logId, err.Error())
		return err
	}
	time.Sleep(10 * time.Second)

	return resourceTencentCloudVpcAclAttachmentRead(d, meta)
}

func resourceTencentCloudVpcAclAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_acl_attachment.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId        = getLogId(contextNil)
		ctx          = context.WithValue(context.TODO(), logIdKey, logId)
		service      = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
		attachmentId = d.Id()
		aclId        string
	)

	if attachmentId == "" {
		return fmt.Errorf("attachmentId does not exist")
	}

	aclId = strings.Split(attachmentId, "#")[0]

	results, err := service.DescribeNetWorkAcls(ctx, aclId, "", "")
	if err != nil {
		return err
	}
	if len(results) < 1 && len(results[0].SubnetSet) < 1 {
		return fmt.Errorf("[TECENT_TERRAFORM_CHECK][ACL attachment][Exists] check: CAM group policy attachment id is not set")
	}
	return nil

}

func resourceTencentCloudVpcAclAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_acl_attachment.delete")()
	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		service       = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
		attachmentAcl = d.Id()
	)

	err := service.DeleteAclAttachment(ctx, attachmentAcl)
	if err != nil {
		log.Printf("[CRITAL]%s delete acl attachment failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil

}
