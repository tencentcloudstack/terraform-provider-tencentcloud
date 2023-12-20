package clb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClbSecurityGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbSecurityGroupAttachmentCreate,
		Read:   resourceTencentCloudClbSecurityGroupAttachmentRead,
		Delete: resourceTencentCloudClbSecurityGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"security_group": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Security group ID, such as esg-12345678.",
			},

			"load_balancer_ids": {
				Required: true,
				Type:     schema.TypeSet,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Array of CLB instance IDs. Only support set one security group now.",
			},
		},
	}
}

func resourceTencentCloudClbSecurityGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_security_group_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request        = clb.NewSetSecurityGroupForLoadbalancersRequest()
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
			loadBalancerId = loadBalancerIdsSet[i].(string)
			request.LoadBalancerIds = append(request.LoadBalancerIds, &loadBalancerId)
		}
	}

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := service.SetClbSecurityGroup(ctx, securityGroup, loadBalancerId, "ADD")
	if err != nil {
		return err
	}

	d.SetId(securityGroup + tccommon.FILED_SP + loadBalancerId)

	return resourceTencentCloudClbSecurityGroupAttachmentRead(d, meta)
}

func resourceTencentCloudClbSecurityGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_security_group_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	securityGroup := idSplit[0]
	loadBalancerId := idSplit[1]

	instance, err := service.DescribeLoadBalancerById(ctx, loadBalancerId)
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClbSecurityGroupAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	exist := false
	for _, sg := range instance.SecureGroups {
		if *sg == securityGroup {
			exist = true
			break
		}
	}
	if !exist {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClbSecurityGroupAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("security_group", securityGroup)

	_ = d.Set("load_balancer_ids", []*string{&loadBalancerId})

	return nil
}

func resourceTencentCloudClbSecurityGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_security_group_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	securityGroup := idSplit[0]
	loadBalancerId := idSplit[1]

	if err := service.SetClbSecurityGroup(ctx, securityGroup, loadBalancerId, "DEL"); err != nil {
		return err
	}

	return nil
}
