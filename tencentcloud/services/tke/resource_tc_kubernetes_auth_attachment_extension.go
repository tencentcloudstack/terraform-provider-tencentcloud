package tke

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudKubernetesAuthAttachmentCreateRequestOnError0(ctx context.Context, req *tke.ModifyClusterAuthenticationOptionsRequest, e error) *resource.RetryError {
	return tccommon.RetryError(e, tke.RESOURCEUNAVAILABLE_CLUSTERSTATE)
}

func resourceTencentCloudKubernetesAuthAttachmentCreatePreRequest0(ctx context.Context, req *tke.ModifyClusterAuthenticationOptionsRequest) *resource.RetryError {
	d := tccommon.ResourceDataFromContext(ctx)
	tmpReqServiceAccount := tke.ServiceAccountAuthenticationOptions{}
	if v, ok := d.GetOk("use_tke_default"); ok && v.(bool) {
		req.ServiceAccounts.Issuer = tmpReqServiceAccount.Issuer
		req.ServiceAccounts.JWKSURI = tmpReqServiceAccount.JWKSURI
	}

	return nil
}

func resourceTencentCloudKubernetesAuthAttachmentReadRequestOnSuccess0(ctx context.Context, resp *tke.DescribeClusterAuthenticationOptionsResponseParams) *resource.RetryError {
	d := tccommon.ResourceDataFromContext(ctx)

	if resp != nil && resp.ServiceAccounts != nil && resp.ServiceAccounts.UseTKEDefault != nil && *resp.ServiceAccounts.UseTKEDefault {
		_ = d.Set("tke_default_issuer", resp.ServiceAccounts.Issuer)
		_ = d.Set("tke_default_jwks_uri", resp.ServiceAccounts.JWKSURI)

		// if true, set params nil
		resp.ServiceAccounts.Issuer = nil
		resp.ServiceAccounts.JWKSURI = nil
	}

	return nil
}

func resourceTencentCloudKubernetesAuthAttachmentUpdatePreRequest0(ctx context.Context, req *tke.ModifyClusterAuthenticationOptionsRequest) *resource.RetryError {
	d := tccommon.ResourceDataFromContext(ctx)

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

func resourceTencentCloudKubernetesAuthAttachmentUpdateRequestOnError0(ctx context.Context, req *tke.ModifyClusterAuthenticationOptionsRequest, e error) *resource.RetryError {
	return tccommon.RetryError(e, tke.RESOURCEUNAVAILABLE_CLUSTERSTATE)
}

func resourceTencentCloudKubernetesAuthAttachmentReadPostFillRequest0(ctx context.Context, req *tke.DescribeClusterAuthenticationOptionsRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)

	meta := tccommon.ProviderMetaFromContext(ctx)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	_, _, err := service.WaitForAuthenticationOptionsUpdateSuccess(ctx, d.Id())

	if err != nil {
		d.SetId("")
		return err
	}
	return nil
}

func resourceTencentCloudKubernetesAuthAttachmentDeletePreRequest0(ctx context.Context, req *tke.ModifyClusterAuthenticationOptionsRequest) *resource.RetryError {
	req.ServiceAccounts = &tke.ServiceAccountAuthenticationOptions{
		JWKSURI: helper.String(""),
		Issuer:  helper.String(DefaultAuthenticationOptionsIssuer),
	}
	return nil
}

func resourceTencentCloudKubernetesAuthAttachmentDeletePostHandleResponse0(ctx context.Context, resp *tke.ModifyClusterAuthenticationOptionsResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)

	meta := tccommon.ProviderMetaFromContext(ctx)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	_, _, err := service.WaitForAuthenticationOptionsUpdateSuccess(ctx, d.Id())

	if err != nil {
		return err
	}
	return nil
}
