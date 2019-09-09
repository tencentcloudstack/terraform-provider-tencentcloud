package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func resourceTencentCloudEniAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEniAttachmentCreate,
		Read:   resourceTencentCloudEniAttachmentRead,
		Delete: resourceTencentCloudEniAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"eni_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceTencentCloudEniAttachmentCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_eni_attachment.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

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
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	split := strings.Split(id, "+")
	if len(split) != 2 {
		log.Printf("[CRITAL]%s id %s is invalid", logId, id)
		d.SetId("")
		return nil
	}

	eniId := split[0]

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	enis, err := service.DescribeEniById(ctx, eniId)
	if err != nil {
		return err
	}

	var eni *vpc.NetworkInterface
	for _, e := range enis {
		if e.NetworkInterfaceId == nil {
			return errors.New("eni id is nil")
		}

		if *e.NetworkInterfaceId == eniId {
			eni = e
			break
		}
	}

	if eni == nil {
		d.SetId("")
		return nil
	}

	if nilFields := CheckNil(eni, map[string]string{
		"NetworkInterfaceId": "eni id",
	}); len(nilFields) > 0 {
		return fmt.Errorf("eni %v are nil", nilFields)
	}

	if eni.Attachment == nil {
		d.SetId("")
		return nil
	}

	if eni.Attachment.InstanceId == nil {
		return errors.New("eni attach instance id is nil")
	}

	d.Set("eni_id", eni.NetworkInterfaceId)
	d.Set("cvm_id", eni.Attachment.InstanceId)

	return nil
}

func resourceTencentCloudEniAttachmentDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_eni_attachment.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

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
