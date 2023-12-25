/*
Provides a resource to create a vpc classic_link_attachment

# Example Usage

```hcl
data "tencentcloud_availability_zones" "zones" {}

	data "tencentcloud_images" "image" {
	  image_type       = ["PUBLIC_IMAGE"]
	  image_name_regex = "Final"
	}

	data "tencentcloud_instance_types" "instance_types" {
	  filter {
	    name   = "zone"
	    values = [data.tencentcloud_availability_zones.zones.zones.0.name]
	  }

	  filter {
	    name   = "instance-family"
	    values = ["S5"]
	  }

	  cpu_core_count   = 2
	  exclude_sold_out = true
	}

	resource "tencentcloud_vpc" "vpc" {
	  name       = "vpc-example"
	  cidr_block = "10.0.0.0/16"
	}

	resource "tencentcloud_subnet" "subnet" {
	  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
	  name              = "subnet-example"
	  vpc_id            = tencentcloud_vpc.vpc.id
	  cidr_block        = "10.0.0.0/16"
	  is_multicast      = false
	}

	resource "tencentcloud_instance" "example" {
	  instance_name            = "tf-example"
	  availability_zone        = data.tencentcloud_availability_zones.zones.zones.0.name
	  image_id                 = data.tencentcloud_images.image.images.0.image_id
	  instance_type            = data.tencentcloud_instance_types.instance_types.instance_types.0.instance_type
	  system_disk_type         = "CLOUD_PREMIUM"
	  disable_security_service = true
	  disable_monitor_service  = true
	  vpc_id                   = tencentcloud_vpc.vpc.id
	  subnet_id                = tencentcloud_subnet.subnet.id
	}

	resource "tencentcloud_vpc_classic_link_attachment" "classic_link_attachment" {
	  vpc_id       = tencentcloud_vpc.vpc.id
	  instance_ids = [tencentcloud_instance.example.id]
	}

```

# Import

vpc classic_link_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_classic_link_attachment.classic_link_attachment classic_link_attachment_id
```
*/
package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcClassicLinkAttachment() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_classic_link_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AttachClassicLinkVpc(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc ClassicLinkAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(vpcId + tccommon.FILED_SP + instanceId)

	return resourceTencentCloudVpcClassicLinkAttachmentRead(d, meta)
}

func resourceTencentCloudVpcClassicLinkAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_classic_link_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_classic_link_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
