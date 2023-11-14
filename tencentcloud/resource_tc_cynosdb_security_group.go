/*
Provides a resource to create a cynosdb security_group

Example Usage

```hcl
resource "tencentcloud_cynosdb_security_group" "security_group" {
  instance_ids = &lt;nil&gt;
  security_group_ids = &lt;nil&gt;
  zone = &lt;nil&gt;
}
```

Import

cynosdb security_group can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_security_group.security_group security_group_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudCynosdbSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbSecurityGroupCreate,
		Read:   resourceTencentCloudCynosdbSecurityGroupRead,
		Delete: resourceTencentCloudCynosdbSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The id list of instances.",
			},

			"security_group_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of security group IDs to be modified, an array of one or more security group IDs.",
			},

			"zone": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Availability zone.",
			},
		},
	}
}

func resourceTencentCloudCynosdbSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_security_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = cynosdb.NewAssociateSecurityGroupsRequest()
		response        = cynosdb.NewAssociateSecurityGroupsResponse()
		instanceId      string
		securityGroupId string
	)
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().AssociateSecurityGroups(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb securityGroup failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(strings.Join([]string{instanceId, securityGroupId}, FILED_SP))

	return resourceTencentCloudCynosdbSecurityGroupRead(d, meta)
}

func resourceTencentCloudCynosdbSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_security_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]

	securityGroup, err := service.DescribeCynosdbSecurityGroupById(ctx, instanceId, securityGroupId)
	if err != nil {
		return err
	}

	if securityGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbSecurityGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if securityGroup.InstanceIds != nil {
		_ = d.Set("instance_ids", securityGroup.InstanceIds)
	}

	if securityGroup.SecurityGroupIds != nil {
		_ = d.Set("security_group_ids", securityGroup.SecurityGroupIds)
	}

	if securityGroup.Zone != nil {
		_ = d.Set("zone", securityGroup.Zone)
	}

	return nil
}

func resourceTencentCloudCynosdbSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_security_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]

	if err := service.DeleteCynosdbSecurityGroupById(ctx, instanceId, securityGroupId); err != nil {
		return err
	}

	return nil
}
