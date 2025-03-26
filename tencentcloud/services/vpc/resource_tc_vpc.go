package vpc

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"fmt"
	"log"
	"time"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcInstanceCreate,
		Read:   resourceTencentCloudVpcInstanceRead,
		Update: resourceTencentCloudVpcInstanceUpdate,
		Delete: resourceTencentCloudVpcInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "The name of the VPC.",
			},
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateCIDRNetworkAddress,
				Description:  "A network address block which should be a subnet of the three internal network segments (10.0.0.0/16, 172.16.0.0/12 and 192.168.0.0/16).",
			},
			"dns_servers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set: func(v interface{}) int {
					return helper.HashString(v.(string))
				},
				Description: "The DNS server list of the VPC. And you can specify 0 to 5 servers to this list.",
			},
			"is_multicast": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether VPC multicast is enabled. The default value is 'true'.",
			},
			"assistant_cidrs": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of Assistant CIDR, NOTE: Only `NORMAL` typed CIDRs included, check the Docker CIDR by readonly `assistant_docker_cidrs`.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"docker_assistant_cidrs": {
				Type:        schema.TypeList,
				Description: "List of Docker Assistant CIDR.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the VPC.",
			},

			// Computed values
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether it is the default VPC for this region.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of VPC.",
			},
			"default_route_table_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default route table id, which created automatically after VPC create.",
			},
		},
	}
}

func resourceTencentCloudVpcInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	vpcService := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		name        string
		cidrBlock   string
		dnsServers  = make([]string, 0, 4)
		isMulticast bool
		tags        map[string]string
	)
	if temp, ok := d.GetOk("name"); ok {
		name = temp.(string)
	}
	if temp, ok := d.GetOk("cidr_block"); ok {
		cidrBlock = temp.(string)
	}
	if temp, ok := d.GetOk("dns_servers"); ok {

		slice := temp.(*schema.Set).List()
		dnsServers = make([]string, 0, len(slice))
		for _, v := range slice {
			dnsServers = append(dnsServers, v.(string))
		}
		if len(dnsServers) < 1 {
			return fmt.Errorf("If dns_servers is set, then len(dns_servers) should be [1:4]")
		}
		if len(dnsServers) > 4 {
			return fmt.Errorf("If dns_servers is set, then len(dns_servers) should be [1:4]")
		}

	}
	isMulticast = d.Get("is_multicast").(bool)

	if temp := helper.GetTags(d, "tags"); len(temp) > 0 {
		tags = temp
	}
	vpcId, _, err := vpcService.CreateVpc(ctx, name, cidrBlock, isMulticast, dnsServers, tags)
	if err != nil {
		return err
	}

	d.SetId(vpcId)

	if v, ok := d.GetOk("assistant_cidrs"); ok {
		assistantCidrs := v.(*schema.Set).List()
		tmpList := make([]*string, 0, len(assistantCidrs))
		for i := range assistantCidrs {
			userIdSet := assistantCidrs[i].(string)
			tmpList = append(tmpList, &userIdSet)
		}

		request := vpc.NewCreateAssistantCidrRequest()
		request.VpcId = &vpcId
		request.CidrBlocks = tmpList
		_, err = vpcService.CreateAssistantCidr(ctx, request)
		if err != nil {
			return err
		}
	}

	// protected while tag is not ready, default is 1s
	time.Sleep(tccommon.WaitReadTimeout)

	return resourceTencentCloudVpcInstanceRead(d, meta)
}

func resourceTencentCloudVpcInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	id := d.Id()
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, e := service.DescribeVpc(ctx, id, "", "")
		if e != nil {
			return tccommon.RetryError(e)
		}

		// deleted
		if has == 0 {
			log.Printf("[WARN]%s %s\n", logId, "vpc has been delete")
			d.SetId("")
			return nil
		}

		if has != 1 {
			errRet := fmt.Errorf("one vpc_id read get %d vpc info", has)
			log.Printf("[CRITAL]%s %s\n", logId, errRet.Error())
			return resource.NonRetryableError(errRet)
		}

		routeTables, err := service.DescribeRouteTables(ctx, "", "", d.Id(), nil, helper.Bool(true), "")

		if err != nil {
			log.Printf("[WARN] Describe default Route Table error: %s", err.Error())
		}

		for i := range routeTables {
			routeTable := routeTables[i]
			if routeTable.isDefault {
				_ = d.Set("default_route_table_id", routeTable.routeTableId)
				break
			}
		}

		tags := make(map[string]string, len(info.tags))
		for _, tag := range info.tags {
			if tag.Key == nil {
				return resource.NonRetryableError(fmt.Errorf("vpc %s tag key is nil", id))
			}
			if tag.Value == nil {
				return resource.NonRetryableError(fmt.Errorf("vpc %s tag value is nil", id))
			}

			tags[*tag.Key] = *tag.Value
		}

		_ = d.Set("name", info.name)
		_ = d.Set("cidr_block", info.cidr)
		_ = d.Set("dns_servers", info.dnsServers)
		_ = d.Set("is_multicast", info.isMulticast)
		_ = d.Set("create_time", info.createTime)
		_ = d.Set("is_default", info.isDefault)
		_ = d.Set("assistant_cidrs", info.assistantCidrs)
		_ = d.Set("docker_assistant_cidrs", info.dockerAssistantCidrs)
		_ = d.Set("tags", tags)

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudVpcInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	vpcService := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	d.Partial(true)

	var (
		name        string
		dnsServers  = make([]string, 0, 4)
		slice       []interface{}
		isMulticast bool
	)

	old, now := d.GetChange("name")
	if d.HasChange("name") {
		name = now.(string)
	} else {
		name = old.(string)
	}

	old, now = d.GetChange("dns_servers")
	if d.HasChange("dns_servers") {
		slice = now.(*schema.Set).List()
		if len(slice) < 1 {
			return fmt.Errorf("If dns_servers is set, then len(dns_servers) should be [1:4]")
		}
		if len(slice) > 4 {
			return fmt.Errorf("If dns_servers is set, then len(dns_servers) should be [1:4]")
		}
	} else {
		slice = old.(*schema.Set).List()
	}

	if len(slice) > 0 {
		for _, v := range slice {
			dnsServers = append(dnsServers, v.(string))
		}
	}

	old, now = d.GetChange("is_multicast")
	if d.HasChange("is_multicast") {
		isMulticast = now.(bool)
	} else {
		isMulticast = old.(bool)
	}

	if err := vpcService.ModifyVpcAttribute(ctx, id, name, isMulticast, dnsServers); err != nil {
		return err
	}

	if d.HasChange("assistant_cidrs") {
		old, now := d.GetChange("assistant_cidrs")
		oldSet := old.(*schema.Set)
		nowSet := now.(*schema.Set)
		add := nowSet.Difference(oldSet).List()
		remove := oldSet.Difference(nowSet).List()
		addLen := len(add)
		removeLen := len(remove)

		request := vpc.NewModifyAssistantCidrRequest()
		request.VpcId = &id
		if removeLen > 0 {
			request.OldCidrBlocks = helper.InterfacesStringsPoint(remove)
		}

		if addLen > 0 {
			request.NewCidrBlocks = helper.InterfacesStringsPoint(add)
		}

		if err := vpcService.ModifyAssistantCidr(ctx, request); err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::vpc:%s:uin/:vpc/%s", region, id)

		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudVpcInstanceRead(d, meta)
}

func resourceTencentCloudVpcInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		if err := service.DeleteVpc(ctx, d.Id()); err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == VPCNotFound {
					return nil
				}
			}
			return resource.RetryableError(err)
		}
		return nil
	})

	return err
}
