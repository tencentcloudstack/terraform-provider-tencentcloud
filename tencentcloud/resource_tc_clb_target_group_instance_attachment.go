/*
Provides a resource to create a CLB target group instance attachment.

Example Usage

```hcl
data "tencentcloud_images" "my_favorite_image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_instance_types" "my_favorite_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S3"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

data "tencentcloud_availability_zones" "default" {
}

resource "tencentcloud_vpc" "app" {
  cidr_block = "10.0.0.0/16"
  name       = "awesome_app_vpc"
}

resource "tencentcloud_subnet" "app" {
  vpc_id            = tencentcloud_vpc.app.id
  availability_zone = data.tencentcloud_availability_zones.default.zones.0.name
  name              = "awesome_app_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_instance" "my_awesome_app" {
  instance_name              = "awesome_app"
  availability_zone          = data.tencentcloud_availability_zones.default.zones.0.name
  image_id                   = data.tencentcloud_images.my_favorite_image.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.my_favorite_instance_types.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  hostname                   = "user"
  project_id                 = 0
  vpc_id                     = tencentcloud_vpc.app.id
  subnet_id                  = tencentcloud_subnet.app.id
  internet_max_bandwidth_out = 20

  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
	encrypt = false
  }

  tags = {
    tagKey = "tagValue"
  }
}

data "tencentcloud_instances" "foo" {
  instance_id = tencentcloud_instance.my_awesome_app.id
}

resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "test"
    vpc_id            = tencentcloud_vpc.app.id
    port              = 33
}

resource "tencentcloud_clb_target_group_instance_attachment" "test"{
    target_group_id = tencentcloud_clb_targetgroup.test.id
    bind_ip         = data.tencentcloud_instances.foo.instance_list[0].private_ip
    port            = 222
    weight          = 3
}

```

Import

CLB target group instance attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_target_group_instance_attachment.test lbtg-3k3io0i0#172.16.48.18#222
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func resourceTencentCloudClbTGAttachmentInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbTGAttachmentInstanceCreate,
		Read:   resourceTencentCloudClbTGAttachmentInstanceRead,
		Update: resourceTencentCloudClbTGAttachmentInstanceUpdate,
		Delete: resourceTencentCloudClbTGAttachmentInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"target_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				ForceNew:     true,
				Description:  "Target group ID.",
			},
			"bind_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				ForceNew:     true,
				Description:  "The Intranet IP of the target group instance.",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Port of the target group instance.",
			},
			"weight": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The weight of the target group instance.",
			},
		},
	}
}

func resourceTencentCloudClbTGAttachmentInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group_instance_attachment.create")()

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		clbService    = ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
		targetGroupId = d.Get("target_group_id").(string)
		bindIp        = d.Get("bind_ip").(string)
		port          = d.Get("port").(int)
		weight        = d.Get("weight").(int)
		err           error
	)

	err = clbService.RegisterTargetInstances(ctx, targetGroupId, bindIp, uint64(port), uint64(weight))

	if err != nil {
		return err
	}
	time.Sleep(time.Duration(3) * time.Second)

	d.SetId(strings.Join([]string{targetGroupId, bindIp, strconv.Itoa(port)}, FILED_SP))

	return resourceTencentCloudClbTGAttachmentInstanceRead(d, meta)
}

func resourceTencentCloudClbTGAttachmentInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group_instance_attachment.read")()

	var (
		logId                = getLogId(contextNil)
		ctx                  = context.WithValue(context.TODO(), logIdKey, logId)
		clbService           = ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                   = d.Id()
		targetGroupInstances []*clb.TargetGroupBackend
	)
	idSplit := strings.Split(id, FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("target group instance attachment id is not set")
	}
	targetGroupId := idSplit[0]
	bindIp := idSplit[1]
	port, err := strconv.ParseUint(idSplit[2], 0, 64)
	if err != nil {
		return err
	}

	filters := make(map[string]string)
	filters["TargetGroupId"] = targetGroupId
	filters["BindIP"] = bindIp
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		targetGroupInstances, err = clbService.DescribeTargetGroupInstances(ctx, filters)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, tgInstance := range targetGroupInstances {
		if *tgInstance.Port == port {
			_ = d.Set("target_group_id", idSplit[0])
			_ = d.Set("bind_ip", idSplit[1])
			_ = d.Set("port", idSplit[2])
			return nil
		}
	}
	d.SetId("")
	return nil
}

func resourceTencentCloudClbTGAttachmentInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group_instance_attachment.update")()

	var (
		logId                 = getLogId(contextNil)
		ctx                   = context.WithValue(context.TODO(), logIdKey, logId)
		clbService            = ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                    = d.Id()
		port                  int
		bindIp, targetGroupId string
		err                   error
	)
	idSplit := strings.Split(id, FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("target group instance attachment id is not set")
	}
	targetGroupId = idSplit[0]
	bindIp = idSplit[1]
	port, err = strconv.Atoi(idSplit[2])
	if err != nil {
		return err
	}

	if d.HasChange("weight") {
		newWeight := d.Get("weight").(int)
		err := clbService.ModifyTargetGroupInstancesWeight(ctx, targetGroupId, bindIp, uint64(port), uint64(newWeight))
		if err != nil {
			return nil
		}
	}
	return resourceTencentCloudClbTGAttachmentInstanceRead(d, meta)
}

func resourceTencentCloudClbTGAttachmentInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group_instance_attachment.delete")()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		clbService = ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
		id         = d.Id()
	)
	idSplit := strings.Split(id, FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("target group instance attachment id is not set")
	}
	targetGroupId := idSplit[0]
	bindIp := idSplit[1]
	port, err := strconv.ParseUint(idSplit[2], 0, 64)
	if err != nil {
		return err
	}

	err = clbService.DeregisterTargetInstances(ctx, targetGroupId, bindIp, port)

	if err != nil {
		return err
	}
	return nil
}
