package mariadb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMariadbSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMariadbSecurityGroupsRead,
		Create: resourceTencentCloudMariadbSecurityGroupsCreate,
		Update: resourceTencentCloudMariadbSecurityGroupsUpdate,
		Delete: resourceTencentCloudMariadbSecurityGroupsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id.",
			},

			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "security group id.",
			},

			"product": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "product name, fixed to mariadb.",
			},
		},
	}
}

func resourceTencentCloudMariadbSecurityGroupsCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_security_groups.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request         = mariadb.NewAssociateSecurityGroupsRequest()
		instanceId      string
		securityGroupId string
		product         string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceIds = []*string{helper.String(v.(string))}
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		securityGroupId = v.(string)
		request.SecurityGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		product = v.(string)
		request.Product = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMariadbClient().AssociateSecurityGroups(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mariadb securityGroups failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + tccommon.FILED_SP + securityGroupId + tccommon.FILED_SP + product)
	return resourceTencentCloudMariadbSecurityGroupsRead(d, meta)
}

func resourceTencentCloudMariadbSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_security_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]
	product := idSplit[2]

	securityGroup, err := service.DescribeMariadbSecurityGroup(ctx, instanceId, securityGroupId, product)

	if err != nil {
		return err
	}

	if securityGroup == nil {
		d.SetId("")
		return fmt.Errorf("resource `securityGroups` %s does not exist", securityGroupId)
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("security_group_id", securityGroupId)
	_ = d.Set("product", product)

	return nil
}

func resourceTencentCloudMariadbSecurityGroupsUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_security_groups.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mariadb.NewModifyDBInstanceSecurityGroupsRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]
	product := idSplit[2]

	request.InstanceId = &instanceId
	request.SecurityGroupIds = []*string{&securityGroupId}
	request.Product = &product

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}

	if d.HasChange("security_group_id") {
		return fmt.Errorf("`security_group_id` do not support change now.")
	}

	if d.HasChange("product") {
		return fmt.Errorf("`product` do not support change now.")
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMariadbClient().ModifyDBInstanceSecurityGroups(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mariadb securityGroups failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMariadbSecurityGroupsRead(d, meta)
}

func resourceTencentCloudMariadbSecurityGroupsDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_security_groups.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]
	product := idSplit[2]

	if err := service.DeleteMariadbSecurityGroupsById(ctx, instanceId, securityGroupId, product); err != nil {
		return err
	}

	return nil
}
