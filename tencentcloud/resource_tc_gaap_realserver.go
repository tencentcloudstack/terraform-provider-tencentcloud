/*
Provides a resource to create a GAAP realserver.

Example Usage

```hcl
resource "tencentcloud_gaap_realserver" "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"

  tags = {
    test = "test"
  }
}
```

Import

GAAP realserver can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_realserver.foo rs-4ftghy6
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudGaapRealserver() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapRealserverCreate,
		Read:   resourceTencentCloudGaapRealserverRead,
		Update: resourceTencentCloudGaapRealserverUpdate,
		Delete: resourceTencentCloudGaapRealserverDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validateIp,
				ConflictsWith: []string{"domain"},
				ForceNew:      true,
				Description:   "IP of the GAAP realserver, conflict with `domain`.",
			},
			"domain": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ip"},
				ForceNew:      true,
				Description:   "Domain of the GAAP realserver, conflict with `ip`.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 30),
				Description:  "Name of the GAAP realserver, the maximum length is 30.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				ForceNew:    true,
				Description: "ID of the project within the GAAP realserver, '0' means is default project.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the GAAP realserver.",
			},
		},
	}
}

func resourceTencentCloudGaapRealserverCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_realserver.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		addressIsSet bool
		address      string
	)

	if ip, ok := d.GetOk("ip"); ok {
		addressIsSet = true
		address = ip.(string)
	}

	if domain, ok := d.GetOk("domain"); ok {
		addressIsSet = true
		address = domain.(string)
	}

	if !addressIsSet {
		return errors.New("ip or domain must be set")
	}

	name := d.Get("name").(string)
	projectId := d.Get("project_id").(int)

	tags := helper.GetTags(d, "tags")

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	realservers, err := service.DescribeRealservers(ctx, &address, nil, nil, -1)
	if err != nil {
		return err
	}
	if len(realservers) > 0 {
		return fmt.Errorf("the realserver with ip/domain %s already exists", address)
	}

	id, err := service.CreateRealserver(ctx, address, name, projectId, tags)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceTencentCloudGaapRealserverRead(d, m)
}

func resourceTencentCloudGaapRealserverRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_realserver.read")()
	defer inconsistentCheck(d, m)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	realservers, err := service.DescribeRealservers(ctx, nil, nil, nil, -1)
	if err != nil {
		return err
	}

	var realserver *gaap.BindRealServerInfo
	for _, rs := range realservers {
		if rs.RealServerId == nil {
			return errors.New("realserver id is nil")
		}

		if *rs.RealServerId == id {
			realserver = rs
			break
		}
	}

	if realserver == nil {
		d.SetId("")
		return nil
	}

	if realserver.RealServerIP == nil {
		return errors.New("realserver ip or domain is nil")
	}
	if net.ParseIP(*realserver.RealServerIP) != nil {
		_ = d.Set("ip", realserver.RealServerIP)
	} else {
		_ = d.Set("domain", realserver.RealServerIP)
	}

	if realserver.RealServerName == nil {
		return errors.New("realserver name is nil")
	}
	_ = d.Set("name", realserver.RealServerName)

	if realserver.ProjectId == nil {
		return errors.New("realserver project id is nil")
	}
	_ = d.Set("project_id", realserver.ProjectId)

	respTags := make(map[string]string, len(realserver.TagSet))
	for _, tag := range realserver.TagSet {
		if tag.TagKey == nil || tag.TagValue == nil {
			return errors.New("realserver tag key or value is nil")
		}
		respTags[*tag.TagKey] = *tag.TagValue
	}
	_ = d.Set("tags", respTags)

	return nil
}

func resourceTencentCloudGaapRealserverUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_realserver.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	d.Partial(true)

	if d.HasChange("name") {
		newName := d.Get("name").(string)

		gaapService := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

		if err := gaapService.ModifyRealserverName(ctx, id, newName); err != nil {
			return err
		}

		d.SetPartial("name")
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		tagService := TagService{client: m.(*TencentCloudClient).apiV3Conn}

		region := m.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::gaap:%s:uin/:realserver/%s", region, id)

		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceTencentCloudGaapRealserverRead(d, m)
}

func resourceTencentCloudGaapRealserverDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_realserver.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DeleteRealserver(ctx, id)
}
