/*
Provides a resource to create a vpc classic_link_attachment

Example Usage

```hcl
resource "tencentcloud_vpc_classic_link_attachment" "classic_link_attachment" {
  vpc_id       = "vpc-hdvfe0g1"
  instance_ids = ["ins-ceynqvnu"]
}
```

Import

vpc classic_link_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_classic_link_attachment.classic_link_attachment classic_link_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpcClassicLinkAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcClassicLinkAttachmentCreate,
		Read:   resourceTencentCloudVpcClassicLinkAttachmentRead,
		Delete: resourceTencentCloudVpcClassicLinkAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "VPC instance ID.",
			},

			"instance_ids": {
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "CVM instance ID. It only support set one instance now.",
			},
		},
	}
}

func resourceTencentCloudVpcClassicLinkAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_classic_link_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = vpc.NewAttachClassicLinkVpcRequest()
		vpcId      string
		instanceId string
	)
	if v, ok := d.GetOk("vpc_id"); ok {
		vpcId = v.(string)
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceId = instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceId)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().AttachClassicLinkVpc(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc ClassicLinkAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(vpcId + FILED_SP + instanceId)

	return resourceTencentCloudVpcClassicLinkAttachmentRead(d, meta)
}

func resourceTencentCloudVpcClassicLinkAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_classic_link_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	vpcId := idSplit[0]
	instanceId := idSplit[1]

	classicLinkAttachment, err := service.DescribeVpcClassicLinkAttachmentById(ctx, vpcId, instanceId)
	if err != nil {
		return err
	}

	if classicLinkAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcClassicLinkAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if classicLinkAttachment.VpcId != nil {
		_ = d.Set("vpc_id", classicLinkAttachment.VpcId)
	}

	if classicLinkAttachment.InstanceId != nil {
		_ = d.Set("instance_ids", []*string{classicLinkAttachment.InstanceId})
	}

	return nil
}

func resourceTencentCloudVpcClassicLinkAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_classic_link_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	vpcId := idSplit[0]
	instanceId := idSplit[1]

	if err := service.DeleteVpcClassicLinkAttachmentById(ctx, vpcId, instanceId); err != nil {
		return err
	}

	return nil
}
