// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	API "github.com/kraudcloud/cli/api"
)

// Ensure KraudeProvider satisfies various provider interfaces.
var _ provider.Provider = &KraudeProvider{}

// KraudeProvider defines the provider implementation.
type KraudeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// KraudeProviderModel describes the provider data model.
type KraudeProviderModel struct {
	AuthToken types.String `tfsdk:"authToken"`
}

func (p *KraudeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kraude"
	resp.Version = p.version
}
func (p *KraudeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *KraudeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config KraudeProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.AuthToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("authToken"),
			"Unknown Kraude Auth Token",
			"The provider cannot create the Kraude API client as there is an unknown configuration value for the Kraude API auth token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the AUTH_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	authToken := os.Getenv("AUTH_TOKEN")

	if !config.AuthToken.IsNull() {
		authToken = config.AuthToken.ValueString()
	}

	if authToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("authToken"),
			"Missing Auth Token",
			"The provider cannot create the Kraude API client as there is a missing or empty value for the Kraude auth Token. "+
				"Set the authToken value in the configuration or use the AUTH_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client := API.NewClient(authToken)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *KraudeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *KraudeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &KraudeProvider{
			version: version,
		}
	}
}
