package org_idp_azure_ad

import (
	"context"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_azure_ad"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org_idp_utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo, d.Get(org_idp_utils.OrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.AddAzureADProvider(ctx, &management.AddAzureADProviderRequest{
		Name:         d.Get(idp_utils.NameVar).(string),
		ClientId:     d.Get(idp_utils.ClientIDVar).(string),
		ClientSecret: d.Get(idp_utils.ClientSecretVar).(string),
		Scopes:       helper.GetOkSetToStringSlice(d, idp_utils.ScopesVar),
		ProviderOptions: &idp.Options{
			IsLinkingAllowed:  d.Get(idp_utils.IsLinkingAllowedVar).(bool),
			IsCreationAllowed: d.Get(idp_utils.IsCreationAllowedVar).(bool),
			IsAutoUpdate:      d.Get(idp_utils.IsAutoUpdateVar).(bool),
			IsAutoCreation:    d.Get(idp_utils.IsAutoCreationVar).(bool),
		},
		Tenant:        idp_azure_ad.ConstructTenant(d),
		EmailVerified: d.Get(idp_azure_ad.EmailVerifiedVar).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to create idp: %v", err)
	}
	d.SetId(resp.GetId())
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo, d.Get(org_idp_utils.OrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChangesExcept(idp_utils.IdpIDVar, org_idp_utils.OrgIDVar) {
		_, err = client.UpdateAzureADProvider(ctx, &management.UpdateAzureADProviderRequest{
			Id:           d.Id(),
			Name:         d.Get(idp_utils.NameVar).(string),
			ClientId:     d.Get(idp_utils.ClientIDVar).(string),
			ClientSecret: d.Get(idp_utils.ClientSecretVar).(string),
			Scopes:       helper.GetOkSetToStringSlice(d, idp_utils.ScopesVar),
			ProviderOptions: &idp.Options{
				IsLinkingAllowed:  d.Get(idp_utils.IsLinkingAllowedVar).(bool),
				IsCreationAllowed: d.Get(idp_utils.IsCreationAllowedVar).(bool),
				IsAutoCreation:    d.Get(idp_utils.IsAutoCreationVar).(bool),
				IsAutoUpdate:      d.Get(idp_utils.IsAutoUpdateVar).(bool),
			},
			Tenant:        idp_azure_ad.ConstructTenant(d),
			EmailVerified: d.Get(idp_azure_ad.EmailVerifiedVar).(bool),
		})
		if err != nil {
			return diag.Errorf("failed to update idp: %v", err)
		}
	}
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo, d.Get(org_idp_utils.OrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.GetProviderByID(ctx, &management.GetProviderByIDRequest{Id: helper.GetID(d, idp_utils.IdpIDVar)})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get idp")
	}
	respIdp := resp.GetIdp()
	cfg := respIdp.GetConfig()
	specificCfg := cfg.GetAzureAd()
	generalCfg := cfg.GetOptions()
	set := map[string]interface{}{
		org_idp_utils.OrgIDVar:         respIdp.GetDetails().GetResourceOwner(),
		idp_utils.NameVar:              respIdp.GetName(),
		idp_utils.ClientIDVar:          specificCfg.GetClientId(),
		idp_utils.ClientSecretVar:      d.Get(idp_utils.ClientSecretVar).(string),
		idp_utils.ScopesVar:            specificCfg.GetScopes(),
		idp_utils.IsLinkingAllowedVar:  generalCfg.GetIsLinkingAllowed(),
		idp_utils.IsCreationAllowedVar: generalCfg.GetIsCreationAllowed(),
		idp_utils.IsAutoCreationVar:    generalCfg.GetIsAutoCreation(),
		idp_utils.IsAutoUpdateVar:      generalCfg.GetIsAutoUpdate(),
		idp_azure_ad.EmailVerifiedVar:  specificCfg.GetEmailVerified(),
		idp_azure_ad.TenantTypeVar:     idp.AzureADTenantType_name[int32(specificCfg.GetTenant().GetTenantType())],
		idp_azure_ad.TenantIDVar:       specificCfg.GetTenant().GetTenantId(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of oidc idp: %v", k, err)
		}
	}
	d.SetId(respIdp.Id)
	return nil
}
