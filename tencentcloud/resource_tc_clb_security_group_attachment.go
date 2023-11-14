/*
Provides a resource to create a clb security_group_attachment

Example Usage

```hcl
resource "tencentcloud_clb_security_group_attachment" "security_group_attachment" {
  security_group = "esg-12345678"
  load_balancer_ids =
}
```

Import

clb security_group_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_clb_security_group_attachment.security_group_attachment security_group_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudClbSecurityGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbSecurityGroupAttachmentCreate,
		Read:   resourceTencentCloudClbSecurityGroupAttachmentRead,
		Update: resourceTencentCloudClbSecurityGroupAttachmentUpdate,
		Delete: resourceTencentCloudClbSecurityGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"security_group": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Security group ID, such as esg-12345678.",
			},

			"load_balancer_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Array of CLB instance IDs.",
			},
		},
	}
}

func resourceTencentCloudClbSecurityGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_security_group_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request        = clb.NewSetSecurityGroupForLoadbalancersRequest()
		response       = clb.NewSetSecurityGroupForLoadbalancersResponse()
		securityGroup  string
		loadBalancerId string
	)
	if v, ok := d.GetOk("security_group"); ok {
		securityGroup = v.(string)
		request.SecurityGroup = helper.String(v.(string))
	}

	if v, ok := d.GetOk("load_balancer_ids"); ok {
		loadBalancerIdsSet := v.(*schema.Set).List()
		for i := range loadBalancerIdsSet {
			loadBalancerIds := loadBalancerIdsSet[i].(string)
			request.LoadBalancerIds = append(request.LoadBalancerIds, &loadBalancerIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().SetSecurityGroupForLoadbalancers(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create clb securityGroupAttachment failed, reason:%+v", logId, err)
		return err
	}

	securityGroup = *response.Response.SecurityGroup
	d.SetId(strings.Join([]string{securityGroup, loadBalancerId}, FILED_SP))

	return resourceTencentCloudClbSecurityGroupAttachmentRead(d, meta)
}

func resourceTencentCloudClbSecurityGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_security_group_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	securityGroup := idSplit[0]
	loadBalancerId := idSplit[1]

	securityGroupAttachment, err := service.DescribeClbSecurityGroupAttachmentById(ctx, securityGroup, loadBalancerId)
	if err != nil {
		return err
	}

	if securityGroupAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClbSecurityGroupAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if securityGroupAttachment.SecurityGroup != nil {
		_ = d.Set("security_group", securityGroupAttachment.SecurityGroup)
	}

	if securityGroupAttachment.LoadBalancerIds != nil {
		_ = d.Set("load_balancer_ids", securityGroupAttachment.LoadBalancerIds)
	}

	return nil
}

func resourceTencentCloudClbSecurityGroupAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_security_group_attachment.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := clb.NewSetSecurityGroupForLoadbalancersRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	securityGroup := idSplit[0]
	loadBalancerId := idSplit[1]

	request.SecurityGroup = &securityGroup
	request.LoadBalancerId = &loadBalancerId

	immutableArgs := []string{"security_group", "load_balancer_ids"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().SetSecurityGroupForLoadbalancers(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update clb securityGroupAttachment failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudClbSecurityGroupAttachmentRead(d, meta)
}

func resourceTencentCloudClbSecurityGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_security_group_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	securityGroup := idSplit[0]
	loadBalancerId := idSplit[1]

	if err := service.DeleteClbSecurityGroupAttachmentById(ctx, securityGroup, loadBalancerId); err != nil {
		return err
	}

	return nil
}
