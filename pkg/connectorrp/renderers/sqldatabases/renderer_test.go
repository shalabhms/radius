// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package sqldatabases

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2015-05-01-preview/sql"
	"github.com/go-logr/logr"
	v1 "github.com/project-radius/radius/pkg/armrpc/api/v1"
	"github.com/project-radius/radius/pkg/azure/clients"
	"github.com/project-radius/radius/pkg/connectorrp/datamodel"
	"github.com/project-radius/radius/pkg/connectorrp/renderers"
	"github.com/project-radius/radius/pkg/providers"
	"github.com/project-radius/radius/pkg/radlogger"
	"github.com/project-radius/radius/pkg/radrp/outputresource"
	"github.com/project-radius/radius/pkg/resourcekinds"
	"github.com/project-radius/radius/pkg/resourcemodel"
	"github.com/stretchr/testify/require"
)

func createContext(t *testing.T) context.Context {
	logger, err := radlogger.NewTestLogger(t)
	if err != nil {
		t.Log("Unable to initialize logger")
		return context.Background()
	}
	return logr.NewContext(context.Background(), logger)
}

func Test_Render_Success(t *testing.T) {
	ctx := createContext(t)
	renderer := Renderer{}

	resource := datamodel.SqlDatabase{
		TrackedResource: v1.TrackedResource{
			ID:   "/subscriptions/testSub/resourceGroups/testGroup/providers/Applications.Connector/sqlDatabases/sql0",
			Name: "sql0",
			Type: "Applications.Connector/sqlDatabases",
		},
		Properties: datamodel.SqlDatabaseProperties{
			Application: "/subscriptions/test-sub/resourceGroups/test-group/providers/Applications.Core/applications/testApplication",
			Resource:    "/subscriptions/test-sub/resourceGroups/test-group/providers/Microsoft.Sql/servers/test-server/databases/test-database",
			Environment: "/subscriptions/test-sub/resourceGroups/test-group/providers/Applications.Core/environments/env0",
		},
	}

	output, err := renderer.Render(ctx, &resource)
	require.NoError(t, err)

	require.Len(t, output.Resources, 2)
	serverResource := output.Resources[0]
	databaseResource := output.Resources[1]

	require.Equal(t, outputresource.LocalIDAzureSqlServer, serverResource.LocalID)
	require.Equal(t, resourcekinds.AzureSqlServer, serverResource.ResourceType.Type)
	require.Equal(t, resourcemodel.NewARMIdentity(
		&resourcemodel.ResourceType{
			Type:     resourcekinds.AzureSqlServer,
			Provider: providers.ProviderAzure,
		},
		"/subscriptions/test-sub/resourceGroups/test-group/providers/Microsoft.Sql/servers/test-server",
		clients.GetAPIVersionFromUserAgent(sql.UserAgent())),
		serverResource.Identity)

	require.Equal(t, outputresource.LocalIDAzureSqlServerDatabase, databaseResource.LocalID)
	require.Equal(t, resourcekinds.AzureSqlServerDatabase, databaseResource.ResourceType.Type)
	require.Equal(t, resourcemodel.NewARMIdentity(
		&resourcemodel.ResourceType{
			Type:     resourcekinds.AzureSqlServerDatabase,
			Provider: providers.ProviderAzure,
		}, "/subscriptions/test-sub/resourceGroups/test-group/providers/Microsoft.Sql/servers/test-server/databases/test-database",
		clients.GetAPIVersionFromUserAgent(sql.UserAgent())),
		databaseResource.Identity)

	expectedComputedValues := map[string]renderers.ComputedValueReference{
		"database": {
			Value: "test-database",
		},
		"server": {
			LocalID:     outputresource.LocalIDAzureSqlServer,
			JSONPointer: "/properties/fullyQualifiedDomainName",
		},
	}
	require.Equal(t, expectedComputedValues, output.ComputedValues)
	require.Empty(t, output.SecretValues)
}

func Test_Render_MissingResource(t *testing.T) {
	ctx := createContext(t)
	renderer := Renderer{}

	resource := datamodel.SqlDatabase{
		TrackedResource: v1.TrackedResource{
			ID:   "/subscriptions/testSub/resourceGroups/testGroup/providers/Applications.Connector/sqlDatabases/sql0",
			Name: "sql0",
			Type: "Applications.Connector/sqlDatabases",
		},
		Properties: datamodel.SqlDatabaseProperties{
			Application: "/subscriptions/test-sub/resourceGroups/test-group/providers/Applications.Core/applications/testApplication",
			Environment: "/subscriptions/test-sub/resourceGroups/test-group/providers/Applications.Core/environments/env0",
		},
	}

	_, err := renderer.Render(ctx, &resource)
	require.Error(t, err)
	require.Equal(t, renderers.ErrorResourceOrServerNameMissingFromResource.Error(), err.Error())
}

func Test_Render_InvalidResourceType(t *testing.T) {
	ctx := createContext(t)
	renderer := Renderer{}
	resource := datamodel.SqlDatabase{
		TrackedResource: v1.TrackedResource{
			ID:   "/subscriptions/testSub/resourceGroups/testGroup/providers/Applications.Connector/sqlDatabases/sql0",
			Name: "sql0",
			Type: "Applications.Connector/sqlDatabases",
		},
		Properties: datamodel.SqlDatabaseProperties{
			Application: "/subscriptions/test-sub/resourceGroups/test-group/providers/Applications.Core/applications/testApplication",
			Resource:    "/subscriptions/test-sub/resourceGroups/test-group/providers/Microsoft.SomethingElse/servers/sqlDatabases/test-database",
			Environment: "/subscriptions/test-sub/resourceGroups/test-group/providers/Applications.Core/environments/env0",
		},
	}

	_, err := renderer.Render(ctx, &resource)
	require.Error(t, err)
	require.Equal(t, "the 'resource' field must refer to a SQL Database", err.Error())
}