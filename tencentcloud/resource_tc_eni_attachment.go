package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudEniAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEniAttachmentCreate,
		Read:   resourceTencentCloudEniAttachmentRead,
		Delete: resourceTencentCloudEniAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"eni_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the ENI.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the instance which bind the ENI.",
			},
		},
	}
}

func resourceTencentCloudEniAttachmentCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_eni_attachment.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	eniId := d.Get("eni_id").(string)
	cvmId := d.Get("instance_id").(string)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	if err := service.AttachEniToCvm(ctx, eniId, cvmId); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s+%s", eniId, cvmId))

	return resourceTencentCloudEniAttachmentRead(d, m)
}

func resourceTencentCloudEniAttachmentRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_eni_attachment.read")()
	defer inconsistentCheck(d, m)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()
	split := strings.Split(id, "+")
	if len(split) != 2 {
		log.Printf("[CRITAL]%s id %s is invalid", logId, id)
		d.SetId("")
		return nil
	}

	eniId := split[0]

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	enis, err := service.DescribeEniById(ctx, []string{eniId})
	if err != nil {
		return err
	}

	if len(enis) < 1 {
		d.SetId("")
		return nil
	}

	eni := enis[0]

	if eni.Attachment == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("eni_id", eni.NetworkInterfaceId)
	_ = d.Set("instance_id", eni.Attachment.InstanceId)

	return nil
}

func resourceTencentCloudEniAttachmentDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_eni_attachment.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()
	split := strings.Split(id, "+")
	if len(split) != 2 {
		log.Printf("[CRITAL]%s id %s is invalid", logId, id)
		d.SetId("")
		return nil
	}

	eniId, cvmId := split[0], split[1]

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DetachEniFromCvm(ctx, eniId, cvmId)
}
