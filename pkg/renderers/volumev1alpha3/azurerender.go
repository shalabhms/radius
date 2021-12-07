// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package volumev1alpha3

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/radius/pkg/azure/armauth"
	"github.com/Azure/radius/pkg/azure/azresources"
	"github.com/Azure/radius/pkg/azure/radclient"
	"github.com/Azure/radius/pkg/radrp/outputresource"
	"github.com/Azure/radius/pkg/renderers"
)

const (
	VolumeKindEphemeral                  = "ephemeral"
	VolumeKindPersistent                 = "persistent"
	PersistentVolumeKindAzureFileShare   = "azure.com.fileshare"
	PersistentVolumeKindAzureComKeyVault = "azure.com.keyvault"
)

var storageAccountDependency outputresource.Dependency

type AzureRenderer struct {
	Arm             armauth.ArmConfig
	VolumeRenderers map[string]func(ctx context.Context, resource renderers.RendererResource, dependencies map[string]renderers.RendererDependency) (renderers.RendererOutput, error)
}

var SupportedVolumeRenderers = map[string]func(ctx context.Context, resource renderers.RendererResource, dependencies map[string]renderers.RendererDependency) (renderers.RendererOutput, error){
	PersistentVolumeKindAzureFileShare: GetAzureFileShareVolume,
}

var SupportedVolumeMakeSecretsAndValues = map[string]func(name string) (map[string]renderers.ComputedValueReference, map[string]renderers.SecretValueReference){
	PersistentVolumeKindAzureFileShare: MakeSecretsAndValuesForAzureFileShare,
}

func (r *AzureRenderer) GetDependencyIDs(ctx context.Context, resource renderers.RendererResource) ([]azresources.ResourceID, []azresources.ResourceID, error) {
	return nil, nil, nil
}

func (r *AzureRenderer) Render(ctx context.Context, options renderers.RenderOptions) (renderers.RendererOutput, error) {
	properties := radclient.VolumeProperties{}
	err := options.Resource.ConvertDefinition(&properties)
	if err != nil {
		return renderers.RendererOutput{}, err
	}
	if properties.Kind == nil {
		return renderers.RendererOutput{}, errors.New("`kind` property is required")
	} else if !isSupported(*properties.Kind) {
		return renderers.RendererOutput{}, fmt.Errorf("%v is not supported. Supported kind values: %v", properties.Kind, SupportedVolumeRenderers)
	}

	renderOutput, err := r.VolumeRenderers[*properties.Kind](ctx, options.Resource, options.Dependencies)
	if err != nil {
		return renderers.RendererOutput{}, err
	}

	computedValues, secretValues := SupportedVolumeMakeSecretsAndValues[*properties.Kind](storageAccountDependency.LocalID)

	return renderers.RendererOutput{
		Resources:      renderOutput.Resources,
		ComputedValues: computedValues,
		SecretValues:   secretValues,
	}, nil
}

func isSupported(kind string) bool {
	for k := range SupportedVolumeRenderers {
		if kind == k {
			return true
		}
	}
	return false
}