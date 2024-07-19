package cvm

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func resourceTencentCloudCvmChcConfigCreateOnExit(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	logId := ctx.Value(tccommon.LogIdKey).(string)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		assistChange  bool
		deployChange  bool
		chcId         string
		vpcId         string
		assistRequest = cvm.NewConfigureChcAssistVpcRequest()
		deployRequest = cvm.NewConfigureChcDeployVpcRequest()
	)
	if v, ok := d.GetOk("chc_id"); ok {
		chcId = v.(string)
	}

	if v, ok := d.GetOk("instance_name"); ok {
		attributeRequest := cvm.NewModifyChcAttributeRequest()
		attributeRequest.InstanceName = helper.String(v.(string))
		attributeRequest.ChcIds = []*string{&chcId}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ModifyChcAttribute(attributeRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, attributeRequest.GetAction(), attributeRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate cvm chcAttribute failed, reason:%+v", logId, err)
			return err
		}
	}

	if v, ok := d.GetOk("device_type"); ok {
		attributeRequest := cvm.NewModifyChcAttributeRequest()
		attributeRequest.DeviceType = helper.String(v.(string))
		attributeRequest.ChcIds = []*string{&chcId}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ModifyChcAttribute(attributeRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, attributeRequest.GetAction(), attributeRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate cvm chcAttribute failed, reason:%+v", logId, err)
			return err
		}
	}
	bmcUser, bmcUserok := d.GetOk("bmc_user")
	password, passwordOk := d.GetOk("password")
	if bmcUserok && passwordOk {
		attributeRequest := cvm.NewModifyChcAttributeRequest()
		attributeRequest.BmcUser = helper.String(bmcUser.(string))
		attributeRequest.Password = helper.String(password.(string))
		attributeRequest.ChcIds = []*string{&chcId}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ModifyChcAttribute(attributeRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, attributeRequest.GetAction(), attributeRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate cvm chcAttribute failed, reason:%+v", logId, err)
			return err
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "bmc_virtual_private_cloud"); ok {
		virtualPrivateCloud := cvm.VirtualPrivateCloud{}
		if v, ok := dMap["vpc_id"]; ok {
			virtualPrivateCloud.VpcId = helper.String(v.(string))
		}
		if v, ok := dMap["subnet_id"]; ok {
			virtualPrivateCloud.SubnetId = helper.String(v.(string))
		}
		if v, ok := dMap["as_vpc_gateway"]; ok {
			virtualPrivateCloud.AsVpcGateway = helper.Bool(v.(bool))
		}
		if v, ok := dMap["private_ip_addresses"]; ok {
			privateIpAddresses := v.([]interface{})
			for i := range privateIpAddresses {
				privateIpAddresses := privateIpAddresses[i].(string)
				virtualPrivateCloud.PrivateIpAddresses = append(virtualPrivateCloud.PrivateIpAddresses, &privateIpAddresses)
			}
		}
		if v, ok := dMap["ipv6_address_count"]; ok {
			virtualPrivateCloud.Ipv6AddressCount = helper.IntUint64(v.(int))
		}
		assistChange = true
		assistRequest.BmcVirtualPrivateCloud = &virtualPrivateCloud
	}

	if v, ok := d.GetOk("bmc_security_group_ids"); ok {
		bmcSecurityGroupIds := v.([]interface{})
		for i := range bmcSecurityGroupIds {
			bmcSecurityGroupIds := bmcSecurityGroupIds[i].(string)
			assistRequest.BmcSecurityGroupIds = append(assistRequest.BmcSecurityGroupIds, &bmcSecurityGroupIds)
		}
		assistChange = true
	}

	if assistChange {
		assistRequest.ChcIds = []*string{&chcId}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ConfigureChcAssistVpc(assistRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, assistRequest.GetAction(), assistRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create cvm chcAssistVpc failed, reason:%+v", logId, err)
			return err
		}
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"READY"}, 20*tccommon.ReadRetryTimeout, time.Second, service.CvmChcInstanceStateRefreshFunc(chcId, []string{}))

		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "deploy_virtual_private_cloud"); ok {
		virtualPrivateCloud := cvm.VirtualPrivateCloud{}
		if v, ok := dMap["vpc_id"]; ok {
			vpcId = v.(string)
			virtualPrivateCloud.VpcId = helper.String(vpcId)
		}
		if v, ok := dMap["subnet_id"]; ok {
			virtualPrivateCloud.SubnetId = helper.String(v.(string))
		}
		if v, ok := dMap["as_vpc_gateway"]; ok {
			virtualPrivateCloud.AsVpcGateway = helper.Bool(v.(bool))
		}
		if v, ok := dMap["private_ip_addresses"]; ok {
			privateIpAddresses := v.([]interface{})
			for i := range privateIpAddresses {
				privateIpAddresses := privateIpAddresses[i].(string)
				virtualPrivateCloud.PrivateIpAddresses = append(virtualPrivateCloud.PrivateIpAddresses, &privateIpAddresses)
			}
		}
		if v, ok := dMap["ipv6_address_count"]; ok {
			virtualPrivateCloud.Ipv6AddressCount = helper.IntUint64(v.(int))
		}
		deployRequest.DeployVirtualPrivateCloud = &virtualPrivateCloud
		deployChange = true
	}

	if v, ok := d.GetOk("deploy_security_group_ids"); ok {
		deploySecurityGroupIds := v.([]interface{})
		for i := range deploySecurityGroupIds {
			deploySecurityGroupIds := deploySecurityGroupIds[i].(string)
			deployRequest.DeploySecurityGroupIds = append(deployRequest.DeploySecurityGroupIds, &deploySecurityGroupIds)
		}
		deployChange = true
	}

	if deployChange {
		deployRequest.ChcIds = []*string{&chcId}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ConfigureChcDeployVpc(deployRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, deployRequest.GetAction(), deployRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create cvm chcDeployVpc failed, reason:%+v", logId, err)
			return err
		}

		conf := tccommon.BuildStateChangeConf([]string{}, []string{vpcId}, 10*tccommon.ReadRetryTimeout, time.Second, service.CvmChcInstanceDeployVpcStateRefreshFunc(chcId, []string{}))

		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	d.SetId(chcId)

	return nil
}

func resourceTencentCloudCvmChcConfigDeletePostHandleResponse0(ctx context.Context, resp *cvm.RemoveChcDeployVpcResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	chcId := d.Id()

	conf := tccommon.BuildStateChangeConf([]string{}, []string{""}, 5*tccommon.ReadRetryTimeout, time.Second, service.CvmChcInstanceDeployVpcStateRefreshFunc(chcId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	params := map[string]interface{}{
		"chc_ids": []string{chcId},
	}
	chcHosts, err := service.DescribeCvmChcHostsByFilter(ctx, params)
	if err != nil {
		return err
	}
	if len(chcHosts) > 0 && *chcHosts[0].InstanceState == "INIT" {
		return nil
	}

	return nil
}

func resourceTencentCloudCvmChcConfigDeletePostHandleResponse1(ctx context.Context, resp *cvm.RemoveChcAssistVpcResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"INIT"}, 10*tccommon.ReadRetryTimeout, time.Second, service.CvmChcInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
