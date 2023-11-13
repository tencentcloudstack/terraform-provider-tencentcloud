/*
Provides a resource to create a mariadb security_groups

Example Usage

```hcl
resource "tencentcloud_mariadb_security_groups" "security_groups" {
  instance_ids = &lt;nil&gt;
  security_group_id = &lt;nil&gt;
  product = &lt;nil&gt;
}
```

Import

mariadb security_groups can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_security_groups.security_groups security_groups_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudMariadbSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbSecurityGroupsCreate,
		Read:   resourceTencentCloudMariadbSecurityGroupsRead,
		Update: resourceTencentCloudMariadbSecurityGroupsUpdate,
		Delete: resourceTencentCloudMariadbSecurityGroupsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Instance ids.",
			},

			"security_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Security group id.",
			},

			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Product name, fixed to mariadb.",
			},
		},
	}
}

func resourceTencentCloudMariadbSecurityGroupsCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_security_groups.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = mariadb.NewAssociateSecurityGroupsRequest()
		response        = mariadb.NewAssociateSecurityGroupsResponse()
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

	if v, ok := d.GetOk("security_group_id"); ok {
		securityGroupId = v.(string)
		request.SecurityGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().AssociateSecurityGroups(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mariadb securityGroups failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(strings.Join([]string{instanceId, securityGroupId}, FILED_SP))

	return resourceTencentCloudMariadbSecurityGroupsRead(d, meta)
}

func resourceTencentCloudMariadbSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_security_groups.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]

	securityGroups, err := service.DescribeMariadbSecurityGroupsById(ctx, instanceId, securityGroupId)
	if err != nil {
		return err
	}

	if securityGroups == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MariadbSecurityGroups` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if securityGroups.InstanceIds != nil {
		_ = d.Set("instance_ids", securityGroups.InstanceIds)
	}

	if securityGroups.SecurityGroupId != nil {
		_ = d.Set("security_group_id", securityGroups.SecurityGroupId)
	}

	if securityGroups.Product != nil {
		_ = d.Set("product", securityGroups.Product)
	}

	return nil
}

func resourceTencentCloudMariadbSecurityGroupsUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_security_groups.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mariadb.NewModifyDBInstanceSecurityGroupsRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]

	request.InstanceId = &instanceId
	request.SecurityGroupId = &securityGroupId

	immutableArgs := []string{"instance_ids", "security_group_id", "product"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyDBInstanceSecurityGroups(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mariadb securityGroups failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMariadbSecurityGroupsRead(d, meta)
}

func resourceTencentCloudMariadbSecurityGroupsDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_security_groups.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]

	if err := service.DeleteMariadbSecurityGroupsById(ctx, instanceId, securityGroupId); err != nil {
		return err
	}

	return nil
}
