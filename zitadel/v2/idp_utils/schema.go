package idp_utils

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const (
	IdpIDVar             = "id"
	NameVar              = "name"
	ClientIDVar          = "client_id"
	ClientSecretVar      = "client_secret"
	ScopesVar            = "scopes"
	IsLinkingAllowedVar  = "is_linking_allowed"
	IsCreationAllowedVar = "is_creation_allowed"
	IsAutoCreationVar    = "is_auto_creation"
	IsAutoUpdateVar      = "is_auto_update"
)

var (
	IdPIDDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The ID of this resource.",
	}
	NameResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of the IDP",
	}
	NameDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the IDP",
	}
	ClientIDResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "client id generated by the identity provider",
	}
	ClientIDDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "client id generated by the identity provider",
	}
	ClientSecretResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "client secret generated by the identity provider",
		Sensitive:   true,
	}
	ClientSecretDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "client secret generated by the identity provider",
		Sensitive:   true,
	}
	ScopesResourceField = &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Required:    true,
		Description: "the scopes requested by ZITADEL during the request on the identity provider",
	}
	ScopesDataSourceField = &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Computed:    true,
		Description: "the scopes requested by ZITADEL during the request on the identity provider",
	}
	IsLinkingAllowedResourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Required:    true,
		Description: "enable if users should be able to link an existing ZITADEL user with an external account",
	}
	IsLinkingAllowedDataSourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "enabled if users are able to link an existing ZITADEL user with an external account",
	}
	IsCreationAllowedResourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Required:    true,
		Description: "enable if users should be able to create a new account in ZITADEL when using an external account",
	}
	IsCreationAllowedDataSourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "enabled if users are able to create a new account in ZITADEL when using an external account",
	}
	IsAutoCreationResourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Required:    true,
		Description: "enable if a new account in ZITADEL should be created automatically on login with an external account",
	}
	IsAutoCreationDataSourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "enabled if a new account in ZITADEL are created automatically on login with an external account",
	}
	IsAutoUpdateResourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Required:    true,
		Description: "enable if a the ZITADEL account fields should be updated automatically on each login",
	}
	IsAutoUpdateDataSourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "enabled if a the ZITADEL account fields are updated automatically on each login",
	}
)
