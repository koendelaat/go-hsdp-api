package stl

import (
	"context"
	"github.com/hasura/go-graphql-client"
)

type AppsService struct {
	client *Client
}

type AppResource struct {
	ID       int64
	DeviceID int64
	Name     string
	Content  string
}

type CreateApplicationResourceInput struct {
	DeviceId     int64
	SerialNumber string
	GroupId      string
	Name         string
	Content      string
	IsLocked     bool
}

type DeleteApplicationResourceInput struct {
	ID           int64
	Name         string
	SerialNumber string
	DeviceID     int64
	GroupID      string
}

func (a *AppsService) GetAppResourceByDeviceIDAndName(ctx context.Context, deviceID int64, name string) (*AppResource, error) {
	var query struct {
		App AppResource `graphql:"applicationResource(id: $id, name: $name)"`
	}
	err := a.client.gql.Query(ctx, &query, map[string]interface{}{
		"id":   graphql.Int(deviceID),
		"name": graphql.String(name),
	})
	if err != nil {
		return nil, err
	}
	return &query.App, nil
}

func (a *AppsService) GetAppResourcesBySerial(ctx context.Context, serial string) (*[]AppResource, error) {
	var query struct {
		Resources struct {
			Edges []struct {
				Node AppResource
			}
		} `graphql:"applicationResources(serialNumber: $serial, first: 10000)"`
	}
	err := a.client.gql.Query(ctx, &query, map[string]interface{}{
		"serial": graphql.String(serial),
	})
	if err != nil {
		return nil, err
	}
	appResources := make([]AppResource, 0)
	for _, a := range query.Resources.Edges {
		appResources = append(appResources, a.Node)
	}
	return &appResources, nil
}

func (a *AppsService) CreateAppResource(ctx context.Context, input CreateApplicationResourceInput) (*AppResource, error) {
	var mutation struct {
		CreateApplicationResource struct {
			Success             string
			Message             string
			StatusCode          int
			RequestID           string
			ApplicationResource AppResource
		} `graphql:"createApplicationResource(input: $input)"`
	}
	err := a.client.gql.Query(ctx, &mutation, map[string]interface{}{
		"input": input,
	})
	if err != nil {
		return nil, err
	}
	return &mutation.CreateApplicationResource.ApplicationResource, nil
}

func (a *AppsService) DeleteAppResource(ctx context.Context, input DeleteApplicationResourceInput) (bool, error) {
	var mutation struct {
		DeleteApplicationResource struct {
			Success    string
			Message    string
			StatusCode int
			RequestID  string
		} `graphql:"deleteApplicationResource(input: $input)"`
	}
	err := a.client.gql.Query(ctx, &mutation, map[string]interface{}{
		"input": input,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
