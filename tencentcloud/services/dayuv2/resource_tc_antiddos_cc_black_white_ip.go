package dayuv2

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAntiddosCcBlackWhiteIp() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAntiddosCcBlackWhiteIpCreate,
		Read:   resourceTencentCloudAntiddosCcBlackWhiteIpRead,
		Delete: resourceTencentCloudAntiddosCcBlackWhiteIpDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"black_white_ip": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Black white ip.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ip address.",
						},
						"mask": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "ip mask.",
						},
					},
				},
			},

			"type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "IP type, value [black(blacklist IP), white(whitelist IP)].",
			},

			"ip": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ip address.",
			},

			"domain": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "domain.",
			},

			"protocol": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "protocol.",
			},
		},
	}
}

func resourceTencentCloudAntiddosCcBlackWhiteIpCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_cc_black_white_ip.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request          = antiddos.NewCreateCcBlackWhiteIpListRequest()
		instanceId       string
		domain           string
		protocol         string
		blackWhiteIpType string
		blackWhiteIp     string
		ip               string
		mask             int
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "black_white_ip"); ok {
		ipSegment := antiddos.IpSegment{}
		if v, ok := dMap["ip"]; ok {
			blackWhiteIp = v.(string)
			ipSegment.Ip = helper.String(blackWhiteIp)
		}
		if v, ok := dMap["mask"]; ok {
			mask = v.(int)
			ipSegment.Mask = helper.IntUint64(mask)
		}
		request.IpList = []*antiddos.IpSegment{&ipSegment}
	}

	if v, ok := d.GetOk("type"); ok {
		blackWhiteIpType = v.(string)
		request.Type = helper.String(blackWhiteIpType)
	}

	if v, ok := d.GetOk("ip"); ok {
		ip = v.(string)
		request.Ip = helper.String(ip)
	}

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		request.Domain = helper.String(domain)
	}

	if v, ok := d.GetOk("protocol"); ok {
		protocol = v.(string)
		request.Protocol = helper.String(protocol)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAntiddosClient().CreateCcBlackWhiteIpList(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create antiddos ccBlackWhiteIp failed, reason:%+v", logId, err)
		return err
	}

	service := AntiddosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	ccBlackWhiteIps, err := service.DescribeAntiddosCcBlackWhiteIpById(ctx, "bgpip", instanceId, ip, domain, protocol)
	if err != nil {
		return err
	}
	var ccBlackWhiteIpPolicy *antiddos.CcBlackWhiteIpPolicy
	for _, ccBlackWhiteIp := range ccBlackWhiteIps {
		if *ccBlackWhiteIp.Domain != domain {
			continue
		}
		if *ccBlackWhiteIp.Protocol != protocol {
			continue
		}
		if *ccBlackWhiteIp.Type != blackWhiteIpType {
			continue
		}
		if *ccBlackWhiteIp.BlackWhiteIp != blackWhiteIp {
			continue
		}
		if int(*ccBlackWhiteIp.Mask) != mask {
			continue
		}
		if *ccBlackWhiteIp.Ip != ip {
			continue
		}
		ccBlackWhiteIpPolicy = ccBlackWhiteIp
	}

	if ccBlackWhiteIpPolicy == nil {
		d.SetId("")
		return fmt.Errorf("can not find cc black white ip policy")
	}

	d.SetId(strings.Join([]string{instanceId, *ccBlackWhiteIpPolicy.PolicyId, ip, domain, protocol}, tccommon.FILED_SP))

	return resourceTencentCloudAntiddosCcBlackWhiteIpRead(d, meta)
}

func resourceTencentCloudAntiddosCcBlackWhiteIpRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_cc_black_white_ip.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := AntiddosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 5 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	policyId := idSplit[1]
	ip := idSplit[2]
	domain := idSplit[3]
	protocol := idSplit[4]

	ccBlackWhiteIps, err := service.DescribeAntiddosCcBlackWhiteIpById(ctx, "bgpip", instanceId, ip, domain, protocol)
	if err != nil {
		return err
	}

	var ccBlackWhiteIp *antiddos.CcBlackWhiteIpPolicy
	for _, item := range ccBlackWhiteIps {
		if *item.PolicyId == policyId {
			ccBlackWhiteIp = item
			break
		}
	}
	if ccBlackWhiteIp == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `AntiddosCcBlackWhiteIp` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	ipListMap := map[string]interface{}{}

	if ccBlackWhiteIp.BlackWhiteIp != nil {
		ipListMap["ip"] = ccBlackWhiteIp.BlackWhiteIp
	}

	if ccBlackWhiteIp.Mask != nil {
		ipListMap["mask"] = ccBlackWhiteIp.Mask
	}

	_ = d.Set("black_white_ip", []interface{}{ipListMap})

	if ccBlackWhiteIp.Type != nil {
		_ = d.Set("type", ccBlackWhiteIp.Type)
	}

	if ccBlackWhiteIp.Ip != nil {
		_ = d.Set("ip", ccBlackWhiteIp.Ip)
	}

	if ccBlackWhiteIp.Domain != nil {
		_ = d.Set("domain", ccBlackWhiteIp.Domain)
	}

	if ccBlackWhiteIp.Protocol != nil {
		_ = d.Set("protocol", ccBlackWhiteIp.Protocol)
	}

	return nil
}

func resourceTencentCloudAntiddosCcBlackWhiteIpDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_cc_black_white_ip.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := AntiddosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 5 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	policyId := idSplit[1]

	if err := service.DeleteAntiddosCcBlackWhiteIpById(ctx, instanceId, policyId); err != nil {
		return err
	}

	return nil
}
