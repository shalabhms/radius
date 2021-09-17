// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package cosmosdbsqlv1alpha1

import (
	"context"
	"fmt"

	"github.com/Azure/radius/pkg/azure/armauth"
	"github.com/Azure/radius/pkg/azure/azresources"
	"github.com/Azure/radius/pkg/azure/clients"
	"github.com/Azure/radius/pkg/handlers"
	"github.com/Azure/radius/pkg/model/components"
	"github.com/Azure/radius/pkg/radlogger"
	"github.com/Azure/radius/pkg/radrp/outputresource"
	"github.com/Azure/radius/pkg/renderers"
	"github.com/Azure/radius/pkg/resourcekinds"
	"github.com/Azure/radius/pkg/workloads"
)

// Renderer WorkloadRenderer implementation for the CosmosDB for SQL workload.
type Renderer struct {
	Arm armauth.ArmConfig
}

// AllocateBindings WorkloadRenderer implementation for CosmosDB for SQL workload.
func (r Renderer) AllocateBindings(ctx context.Context, workload workloads.InstantiatedWorkload, resources []workloads.WorkloadResourceProperties) (map[string]components.BindingState, error) {
	logger := radlogger.GetLogger(ctx)
	if len(workload.Workload.Bindings) > 0 {
		return nil, fmt.Errorf("component of kind %s does not support user-defined bindings", Kind)
	}

	if len(resources) != 1 || resources[0].Type != resourcekinds.AzureCosmosDBSQL {
		return nil, fmt.Errorf("cannot fulfill service - expected properties for %s", resourcekinds.AzureCosmosDBSQL)
	}

	properties := resources[0].Properties
	accountname := properties[handlers.CosmosDBAccountNameKey]
	dbname := properties[handlers.CosmosDBDatabaseNameKey]

	logger.Info(fmt.Sprintf("fulfilling service for account: %v, db: %v", accountname, dbname))

	cosmosDBClient := clients.NewDatabaseAccountsClient(r.Arm.SubscriptionID, r.Arm.Auth)

	connectionStrings, err := cosmosDBClient.ListConnectionStrings(ctx, r.Arm.ResourceGroup, accountname)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve connection strings: %w", err)
	}

	if connectionStrings.ConnectionStrings == nil || len(*connectionStrings.ConnectionStrings) == 0 {
		return nil, fmt.Errorf("no connection strings found for cosmos db account: %s", accountname)
	}

	bindings := map[string]components.BindingState{
		"cosmos": {
			Component: workload.Name,
			Binding:   "cosmos",
			Kind:      "azure.com/CosmosDBSQL",
			Properties: map[string]interface{}{
				"connectionString": *(*connectionStrings.ConnectionStrings)[0].ConnectionString,
				"database":         dbname,
			},
		},
		"sql": {
			Component: workload.Name,
			Binding:   "sql",
			Kind:      "microsoft.com/SQL",
			Properties: map[string]interface{}{
				"connectionString": *(*connectionStrings.ConnectionStrings)[0].ConnectionString,
				"database":         dbname,
			},
		},
	}

	return bindings, nil
}

// Render WorkloadRenderer implementation for CosmosDB for SQL workload.
func (r Renderer) Render(ctx context.Context, w workloads.InstantiatedWorkload) ([]outputresource.OutputResource, error) {
	component := CosmosDBSQLComponent{}
	err := w.Workload.AsRequired(Kind, &component)
	if err != nil {
		return []outputresource.OutputResource{}, err
	}

	if component.Config.Managed {
		if component.Config.Resource != "" {
			return nil, renderers.ErrResourceSpecifiedForManagedResource
		}

		// generate data we can use to manage a cosmosdb instance
		resource := outputresource.OutputResource{
			Kind:    resourcekinds.AzureCosmosDBSQL,
			Type:    outputresource.TypeARM,
			LocalID: outputresource.LocalIDAzureCosmosDBSQL,
			Resource: map[string]string{
				handlers.ManagedKey:              "true",
				handlers.CosmosDBAccountBaseName: w.Workload.Name,
				handlers.CosmosDBDatabaseNameKey: w.Workload.Name,
			},
		}

		return []outputresource.OutputResource{resource}, nil
	}

	if component.Config.Resource == "" {
		return nil, renderers.ErrResourceMissingForUnmanagedResource
	}

	databaseID, err := renderers.ValidateResourceID(component.Config.Resource, SQLResourceType, "CosmosDB SQL Database")
	if err != nil {
		return nil, err
	}

	resource := outputresource.OutputResource{
		Kind:    resourcekinds.AzureCosmosDBSQL,
		Type:    outputresource.TypeARM,
		LocalID: outputresource.LocalIDAzureCosmosDBSQL,
		Resource: map[string]string{
			handlers.ManagedKey: "false",

			// Truncate the database part of the ID to make an ID for the account
			handlers.CosmosDBAccountIDKey:    azresources.MakeID(databaseID.SubscriptionID, databaseID.ResourceGroup, databaseID.Types[0]),
			handlers.CosmosDBDatabaseIDKey:   databaseID.ID,
			handlers.CosmosDBAccountNameKey:  databaseID.Types[0].Name,
			handlers.CosmosDBDatabaseNameKey: databaseID.Types[1].Name,
		},
	}
	return []outputresource.OutputResource{resource}, nil
}