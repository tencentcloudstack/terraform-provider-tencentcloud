package gaap

import (
	"context"
	"errors"
	"fmt"
	"net"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGaapRealserver() *schema.Resource {
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
				ValidateFunc:  tccommon.ValidateIp,
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
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 30),
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
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_realserver.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

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
	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	realservers, err := service.DescribeRealservers(ctx, &address, nil, nil, -1)
	if err != nil {
		return err
	}
	if len(realservers) > 0 {
		return fmt.Errorf("the realserver with ip/domain %s already exists", address)
	}

	id, err := service.CreateRealserver(ctx, address, name, projectId)
	if err != nil {
		return err
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagClient := m.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := &TagService{client: tagClient}
		resourceName := tccommon.BuildTagResourceName("gaap", "realServer", tagClient.Region, id)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	d.SetId(id)

	return resourceTencentCloudGaapRealserverRead(d, m)
}

func resourceTencentCloudGaapRealserverRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_realserver.read")()
	defer tccommon.InconsistentCheck(d, m)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

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

	tagClient := m.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := TagService{client: tagClient}
	tags, err := tagService.DescribeResourceTags(ctx, "gaap", "realServer", tagClient.Region, id)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudGaapRealserverUpdate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_realserver.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	d.Partial(true)

	if d.HasChange("name") {
		newName := d.Get("name").(string)

		gaapService := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

		if err := gaapService.ModifyRealserverName(ctx, id, newName); err != nil {
			return err
		}

	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		tagService := TagService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

		region := m.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := tccommon.BuildTagResourceName("gaap", "realServer", region, id)

		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudGaapRealserverRead(d, m)
}

func resourceTencentCloudGaapRealserverDelete(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_realserver.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	return service.DeleteRealserver(ctx, id)
}
