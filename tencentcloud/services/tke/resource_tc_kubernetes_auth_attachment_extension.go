package tke

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudKubernetesAuthAttachmentCreateRequestOnError0(d *schema.ResourceData, meta interface{}, req *tke.ModifyClusterAuthenticationOptionsRequest, e error) *resource.RetryError {
	return tccommon.RetryError(e, tke.RESOURCEUNAVAILABLE_CLUSTERSTATE)
}

func resourceTencentCloudKubernetesAuthAttachmentCreatePreRequest0(d *schema.ResourceData, meta interface{}, req *tke.ModifyClusterAuthenticationOptionsRequest) error {
	tmpReqServiceAccount := tke.ServiceAccountAuthenticationOptions{}
	if v, ok := d.GetOk("use_tke_default"); ok && v.(bool) {
		req.ServiceAccounts.Issuer = tmpReqServiceAccount.Issuer
		req.ServiceAccounts.JWKSURI = tmpReqServiceAccount.JWKSURI
	}

	return nil
}
func resourceTencentCloudKubernetesAuthAttachmentReadRequestOnSuccess0(d *schema.ResourceData, meta interface{}, resp *tke.DescribeClusterAuthenticationOptionsResponseParams) *resource.RetryError {
	tmpRespServiceAccount := tke.ServiceAccountAuthenticationOptions{}

	if resp != nil && resp.ServiceAccounts != nil {
		if v, ok := d.GetOk("use_tke_default"); ok && v.(bool) {
			resp.ServiceAccounts.Issuer = tmpRespServiceAccount.Issuer
			resp.ServiceAccounts.JWKSURI = tmpRespServiceAccount.JWKSURI
		}
		resp.ServiceAccounts.UseTKEDefault = tmpRespServiceAccount.UseTKEDefault
		resp.ServiceAccounts.AutoCreateDiscoveryAnonymousAuth = tmpRespServiceAccount.AutoCreateDiscoveryAnonymousAuth
	}

	return nil
}
func resourceTencentCloudKubernetesAuthAttachmentReadPostRequest0(ctx context.Context, d *schema.ResourceData, meta interface{}, resp *tke.DescribeClusterAuthenticationOptionsResponseParams) error {
	if resp != nil && resp.ServiceAccounts != nil {
		if v, ok := d.GetOk("use_tke_default"); ok && v.(bool) {
			if resp.ServiceAccounts.Issuer != nil {
				_ = d.Set("tke_default_issuer", resp.ServiceAccounts.Issuer)
			}
			if resp.ServiceAccounts.JWKSURI != nil {
				_ = d.Set("tke_default_jwks_uri", resp.ServiceAccounts.JWKSURI)
			}
		}
	}
	return nil
}
func resourceTencentCloudKubernetesAuthAttachmentUpdatePreRequest0(ctx context.Context, d *schema.ResourceData, meta interface{}, req *tke.ModifyClusterAuthenticationOptionsRequest) error {
	useTkeDefault := false
	tmpReqServiceAccount := tke.ServiceAccountAuthenticationOptions{}
	req.ServiceAccounts.JWKSURI = tmpReqServiceAccount.JWKSURI
	req.ServiceAccounts.Issuer = tmpReqServiceAccount.Issuer
	req.ServiceAccounts.UseTKEDefault = tmpReqServiceAccount.UseTKEDefault

	if v, ok := d.GetOk("use_tke_default"); ok {
		req.ServiceAccounts.UseTKEDefault = helper.Bool(v.(bool))
		useTkeDefault = v.(bool)
	}

	if !useTkeDefault {
		if d.HasChange("jwks_uri") {
			req.ServiceAccounts.JWKSURI = helper.String(d.Get("jwks_uri").(string))
		}
		if d.HasChange("issuer") {
			issuer := d.Get("issuer").(string)
			req.ServiceAccounts.Issuer = helper.String(issuer)
		}
	}
	return nil
}

func resourceTencentCloudKubernetesAuthAttachmentUpdateRequestOnError0(d *schema.ResourceData, meta interface{}, req *tke.ModifyClusterAuthenticationOptionsRequest, e error) *resource.RetryError {
	return tccommon.RetryError(e, tke.RESOURCEUNAVAILABLE_CLUSTERSTATE)
}
