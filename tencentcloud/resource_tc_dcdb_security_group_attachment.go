/*
Provides a resource to create a dcdb security_group_attachment

Example Usage

```hcl
resource "tencentcloud_dcdb_security_group_attachment" "security_group_attachment" {
  security_group_id = &lt;nil&gt;
  instance_id = &lt;nil&gt;
}
```

Import

dcdb security_group_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_security_group_attachment.security_group_attachment security_group_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudDcdbSecurityGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbSecurityGroupAttachmentCreate,
		Read:   resourceTencentCloudDcdbSecurityGroupAttachmentRead,
		Delete: resourceTencentCloudDcdbSecurityGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Security group id.",
			},

			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Attached instance id.",
			},
		},
	}
}

func resourceTencentCloudDcdbSecurityGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_security_group_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = dcdb.NewAssociateSecurityGroupsRequest()
		response        = dcdb.NewAssociateSecurityGroupsResponse()
		instanceId      string
		securityGroupId string
	)
	if v, ok := d.GetOk("security_group_id"); ok {
		securityGroupId = v.(string)
		request.SecurityGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().AssociateSecurityGroups(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dcdb securityGroupAttachment failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(strings.Join([]string{instanceId, securityGroupId}, FILED_SP))

	return resourceTencentCloudDcdbSecurityGroupAttachmentRead(d, meta)
}

func resourceTencentCloudDcdbSecurityGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_security_group_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]

	securityGroupAttachment, err := service.DescribeDcdbSecurityGroupAttachmentById(ctx, instanceId, securityGroupId)
	if err != nil {
		return err
	}

	if securityGroupAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DcdbSecurityGroupAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if securityGroupAttachment.SecurityGroupId != nil {
		_ = d.Set("security_group_id", securityGroupAttachment.SecurityGroupId)
	}

	if securityGroupAttachment.InstanceId != nil {
		_ = d.Set("instance_id", securityGroupAttachment.InstanceId)
	}

	return nil
}

func resourceTencentCloudDcdbSecurityGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_security_group_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]

	if err := service.DeleteDcdbSecurityGroupAttachmentById(ctx, instanceId, securityGroupId); err != nil {
		return err
	}

	return nil
}
