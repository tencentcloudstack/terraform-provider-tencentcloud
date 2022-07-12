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
		subnet_id = data.tencentcloud_vpc_instances.id_instances.instance_list[0].subnet_ids[0]
}
```

Import

Acl attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpc_acl_attachment.attachment acl-eotx5qsg#subnet-91x0geu6
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "ID of the attached ACL.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Subnet instance ID.",
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
		subnetIds []string
	)

	aclId := d.Get("acl_id").(string)
	subnetId := d.Get("subnet_id").(string)

	subnetIds = append(subnetIds, subnetId)

	err := service.AssociateAclSubnets(ctx, aclId, subnetIds)
	if err != nil {
		return err
	}

	d.SetId(aclId + FILED_SP + subnetId)

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
	)

	idSplit := strings.Split(attachmentId, FILED_SP)
	if len(idSplit) < 2 {
		d.SetId("")
		return nil
	}
	aclId := idSplit[0]
	subnetId := idSplit[1]
	results, err := service.DescribeNetWorkAcls(ctx, aclId, "", "")
	if err != nil {
		return err
	}
	if len(results) < 1 || len(results[0].SubnetSet) < 1 {
		d.SetId("")
		return nil
	}

	_ = d.Set("acl_id", aclId)
	_ = d.Set("subnet_id", subnetId)

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
		log.Printf("[CRITAL]%s delete ACL attachment failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil

}
