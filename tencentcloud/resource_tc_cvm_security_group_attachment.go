/*
Provides a resource to create a cvm security_group_attachment

Example Usage

```hcl
resource "tencentcloud_cvm_security_group_attachment" "security_group_attachment" {
  security_group_ids =
  instance_ids =
}
```

Import

cvm security_group_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_security_group_attachment.security_group_attachment security_group_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"log"
	"strings"
)

func resourceTencentCloudCvmSecurityGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmSecurityGroupAttachmentCreate,
		Read:   resourceTencentCloudCvmSecurityGroupAttachmentRead,
		Delete: resourceTencentCloudCvmSecurityGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"security_group_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "ID of the security group to be associated, such as sg-efil73jd. Only one security group can be associated.",
			},

			"instance_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "ID of the instance bound in the format of ins-lesecurk. You can specify up to 100 instances in each request.",
			},
		},
	}
}

func resourceTencentCloudCvmSecurityGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_security_group_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request          = cvm.NewAssociateSecurityGroupsRequest()
		response         = cvm.NewAssociateSecurityGroupsResponse()
		securityGroupIds string
		instanceIds      string
	)
	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().AssociateSecurityGroups(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cvm securityGroupAttachment failed, reason:%+v", logId, err)
		return err
	}

	securityGroupIds = *response.Response.SecurityGroupIds
	d.SetId(strings.Join([]string{securityGroupIds, instanceIds}, FILED_SP))

	return resourceTencentCloudCvmSecurityGroupAttachmentRead(d, meta)
}

func resourceTencentCloudCvmSecurityGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_security_group_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	securityGroupIds := idSplit[0]
	instanceIds := idSplit[1]

	securityGroupAttachment, err := service.DescribeCvmSecurityGroupAttachmentById(ctx, securityGroupIds, instanceIds)
	if err != nil {
		return err
	}

	if securityGroupAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CvmSecurityGroupAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if securityGroupAttachment.SecurityGroupIds != nil {
		_ = d.Set("security_group_ids", securityGroupAttachment.SecurityGroupIds)
	}

	if securityGroupAttachment.InstanceIds != nil {
		_ = d.Set("instance_ids", securityGroupAttachment.InstanceIds)
	}

	return nil
}

func resourceTencentCloudCvmSecurityGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_security_group_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	securityGroupIds := idSplit[0]
	instanceIds := idSplit[1]

	if err := service.DeleteCvmSecurityGroupAttachmentById(ctx, securityGroupIds, instanceIds); err != nil {
		return err
	}

	return nil
}
