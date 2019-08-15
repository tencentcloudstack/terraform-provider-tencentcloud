package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
)

func resourceTencentCloudGaapRealserver() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapRealserverCreate,
		Read:   resourceTencentCloudGaapRealserverRead,
		Update: resourceTencentCloudGaapRealserverUpdate,
		Delete: resourceTencentCloudGaapRealserverDelete,
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validateIp,
				ConflictsWith: []string{"domain"},
				ForceNew:      true,
			},
			"domain": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ip"},
				ForceNew:      true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 30),
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceTencentCloudGaapRealserverCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_realserver.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

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

	name := d.Get("name").(string)
	projectId := d.Get("project_id").(int)

	if !addressIsSet {
		return errors.New("ip or domain must be set")
	}

	tags := getTags(d, "tags")

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	id, err := service.CreateRealserver(ctx, address, name, projectId, tags)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceTencentCloudGaapRealserverRead(d, m)
}

func resourceTencentCloudGaapRealserverRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_realserver.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	name := d.Get("name").(string)
	tags := getTags(d, "tags")
	projectId := d.Get("project_id").(int)

	var address *string
	if ip, ok := d.GetOk("ip"); ok {
		address = stringToPointer(ip.(string))
	}
	if domain, ok := d.GetOk("domain"); ok {
		address = stringToPointer(domain.(string))
	}

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	realservers, err := service.DescribeRealservers(ctx, address, &name, tags, projectId)
	if err != nil {
		return err
	}

	if len(realservers) == 0 {
		d.SetId("")
		return nil
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
		err := fmt.Errorf("realserver %s ip or domain is nil", *realserver.RealServerId)
		log.Printf("[CRITAL]%s %v", logId, err)
		return err
	}
	if _, ok := d.GetOk("ip"); ok {
		d.Set("ip", realserver.RealServerIP)
	}
	if _, ok := d.GetOk("domain"); ok {
		d.Set("domain", realserver.RealServerIP)
	}

	if realserver.RealServerName == nil {
		err := fmt.Errorf("realserver %s name is nil", *realserver.RealServerId)
		log.Printf("[CRITAL]%s %v", logId, err)
		return err
	}
	d.Set("name", realserver.RealServerName)

	if realserver.ProjectId == nil {
		err := fmt.Errorf("realserver %s project id is nil", *realserver.RealServerId)
		log.Printf("[CRITAL]%s %v", logId, err)
		return err
	}
	d.Set("project_id", realserver.ProjectId)

	respTags := make(map[string]string, len(realserver.TagSet))
	for _, tag := range realserver.TagSet {
		if tag.TagKey == nil || tag.TagValue == nil {
			err := fmt.Errorf("one of realserver %s tag key or value is nil", *realserver.RealServerId)
			log.Printf("[CRITAL]%s %v", logId, err)
			return err
		}
		respTags[*tag.TagKey] = *tag.TagValue
	}
	d.Set("tags", respTags)

	return nil
}

func resourceTencentCloudGaapRealserverUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_realserver.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	newName := d.Get("name").(string)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	if err := service.ModifyRealserverName(ctx, id, newName); err != nil {
		return err
	}

	return resourceTencentCloudGaapRealserverRead(d, m)
}

func resourceTencentCloudGaapRealserverDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_realserver.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DeleteRealserver(ctx, id)
}
