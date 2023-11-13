/*
Provides a resource to create a cdb security_groups_attachment

Example Usage

```hcl
resource "tencentcloud_cdb_security_groups_attachment" "security_groups_attachment" {
  security_group_id = &lt;nil&gt;
  instance_ids = &lt;nil&gt;
  for_readonly_instance = false
}
```

Import

cdb security_groups_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_security_groups_attachment.security_groups_attachment security_groups_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudCdbSecurityGroupsAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbSecurityGroupsAttachmentCreate,
		Read:   resourceTencentCloudCdbSecurityGroupsAttachmentRead,
		Delete: resourceTencentCloudCdbSecurityGroupsAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of security group.",
			},

			"instance_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The id list of instances.",
			},

			"for_readonly_instance": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "When the read-only instance ID is passed in, the default operation is the security group corresponding to the read-only group. If you need to operate the security group of the read-only instance ID, you need to set this parameter to True.",
			},
		},
	}
}

func resourceTencentCloudCdbSecurityGroupsAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_security_groups_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = cdb.NewAssociateSecurityGroupsRequest()
		response        = cdb.NewAssociateSecurityGroupsResponse()
		securityGroupId string
		instanceId      string
	)
	if v, ok := d.GetOk("security_group_id"); ok {
		securityGroupId = v.(string)
		request.SecurityGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	if v, ok := d.GetOkExists("for_readonly_instance"); ok {
		request.ForReadonlyInstance = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().AssociateSecurityGroups(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cdb securityGroupsAttachment failed, reason:%+v", logId, err)
		return err
	}

	securityGroupId = *response.Response.SecurityGroupId
	d.SetId(strings.Join([]string{securityGroupId, instanceId}, FILED_SP))

	return resourceTencentCloudCdbSecurityGroupsAttachmentRead(d, meta)
}

func resourceTencentCloudCdbSecurityGroupsAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_security_groups_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	securityGroupId := idSplit[0]
	instanceId := idSplit[1]

	securityGroupsAttachment, err := service.DescribeCdbSecurityGroupsAttachmentById(ctx, securityGroupId, instanceId)
	if err != nil {
		return err
	}

	if securityGroupsAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdbSecurityGroupsAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if securityGroupsAttachment.SecurityGroupId != nil {
		_ = d.Set("security_group_id", securityGroupsAttachment.SecurityGroupId)
	}

	if securityGroupsAttachment.InstanceIds != nil {
		_ = d.Set("instance_ids", securityGroupsAttachment.InstanceIds)
	}

	if securityGroupsAttachment.ForReadonlyInstance != nil {
		_ = d.Set("for_readonly_instance", securityGroupsAttachment.ForReadonlyInstance)
	}

	return nil
}

func resourceTencentCloudCdbSecurityGroupsAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_security_groups_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	securityGroupId := idSplit[0]
	instanceId := idSplit[1]

	if err := service.DeleteCdbSecurityGroupsAttachmentById(ctx, securityGroupId, instanceId); err != nil {
		return err
	}

	return nil
}
