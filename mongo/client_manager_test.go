package mongo_test

import (
	// Standard Library Imports
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	// External Imports
	"github.com/google/uuid"
	"github.com/ory/fosite"
	"go.mongodb.org/mongo-driver/bson"

	// Internal Imports
	"github.com/matthewhartstonge/storage"
	"github.com/matthewhartstonge/storage/mongo"
)

func TestClientManager_List(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	// generate our expected data.
	expected := createClient(ctx, t, store, expectedClient())

	publishedClient := storage.Client{
		ID:                  uuid.NewString(),
		CreateTime:          time.Now().Unix(),
		UpdateTime:          time.Now().Unix() + 600,
		AllowedAudiences:    []string{},
		AllowedRegions:      []string{},
		AllowedTenantAccess: []string{},
		GrantTypes:          []string{},
		ResponseTypes:       []string{},
		Scopes:              []string{},
		Name:                "published client",
		RedirectURIs:        []string{},
		Contacts:            []string{},
		Published:           true,
	}
	publishedClient = createNewClient(t, ctx, store, publishedClient)

	type args struct {
		filter storage.ListClientsRequest
	}
	tests := []struct {
		name        string
		args        args
		wantResults []storage.Client
		wantErr     bool
		err         error
	}{
		{
			name: "should filter clients by allowed tenant access",
			args: args{
				filter: storage.ListClientsRequest{
					AllowedTenantAccess: expected.AllowedTenantAccess[1],
				},
			},
			wantResults: []storage.Client{
				expected,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "should return empty if no clients are found by allowed tenant access",
			args: args{
				filter: storage.ListClientsRequest{
					AllowedTenantAccess: "No tenant here",
				},
			},
			wantResults: []storage.Client(nil),
			wantErr:     false,
			err:         nil,
		},
		{
			name: "should filter clients by region",
			args: args{
				filter: storage.ListClientsRequest{
					AllowedRegion: expected.AllowedRegions[0],
				},
			},
			wantResults: []storage.Client{
				expected,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "should return empty if no clients are found by region",
			args: args{
				filter: storage.ListClientsRequest{
					AllowedRegion: "NZL",
				},
			},
			wantResults: []storage.Client(nil),
			wantErr:     false,
			err:         nil,
		},
		{
			name: "should filter clients by Redirect URI",
			args: args{
				filter: storage.ListClientsRequest{
					RedirectURI: expected.RedirectURIs[0],
				},
			},
			wantResults: []storage.Client{
				expected,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "should return empty if no clients are found by Redirect URI",
			args: args{
				filter: storage.ListClientsRequest{
					RedirectURI: "https://example.com/callback",
				},
			},
			wantResults: []storage.Client(nil),
			wantErr:     false,
			err:         nil,
		},
		{
			name: "should filter clients by Grant Type",
			args: args{
				filter: storage.ListClientsRequest{
					GrantType: string(fosite.AuthorizeCode),
				},
			},
			wantResults: []storage.Client{
				expected,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "should return empty if no clients are found by Grant Type",
			args: args{
				filter: storage.ListClientsRequest{
					GrantType: "grant",
				},
			},
			wantResults: []storage.Client(nil),
			wantErr:     false,
			err:         nil,
		},
		{
			name: "should filter clients by Response Type",
			args: args{
				filter: storage.ListClientsRequest{
					ResponseType: "token",
				},
			},
			wantResults: []storage.Client{
				expected,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "should return empty if no clients are found by Response Type",
			args: args{
				filter: storage.ListClientsRequest{
					ResponseType: "status_ok",
				},
			},
			wantResults: []storage.Client(nil),
			wantErr:     false,
			err:         nil,
		},
		{
			name: "should filter clients by having all scopes provided when filtered by Scopes Intersection in order",
			args: args{
				filter: storage.ListClientsRequest{
					ScopesIntersection: []string{
						"urn:test:cats:write",
						"urn:test:dogs:read",
					},
				},
			},
			wantResults: []storage.Client{
				expected,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "should filter clients by having all scopes provided when filtered by Scopes Intersection, where the client scopes are out of order",
			args: args{
				filter: storage.ListClientsRequest{
					ScopesIntersection: []string{
						"urn:test:cats:write",
						"urn:test:dogs:read",
					},
				},
			},
			wantResults: []storage.Client{
				expected,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "should filter clients by Scopes Intersection",
			args: args{
				filter: storage.ListClientsRequest{
					ScopesIntersection: []string{
						"urn:test:cats:write",
					},
				},
			},
			wantResults: []storage.Client{
				expected,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "should return empty if all client scopes don't match when filtering by Scopes Intersection",
			args: args{
				filter: storage.ListClientsRequest{
					ScopesIntersection: []string{
						"urn:test:cats:write",
						"urn:test:dogs",
					},
				},
			},
			wantResults: []storage.Client(nil),
			wantErr:     false,
			err:         nil,
		},
		{
			name: "should return empty if no clients are found by Scopes Intersection",
			args: args{
				filter: storage.ListClientsRequest{
					ScopesIntersection: []string{
						"urn:test:dogs",
					},
				},
			},
			wantResults: []storage.Client(nil),
			wantErr:     false,
			err:         nil,
		},
		{
			name: "should filter clients by Scopes Union #1",
			args: args{
				filter: storage.ListClientsRequest{
					ScopesUnion: []string{
						"urn:test:cats:write",
					},
				},
			},
			wantResults: []storage.Client{
				expected,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "should filter clients by Scopes Union #2",
			args: args{
				filter: storage.ListClientsRequest{
					ScopesUnion: []string{
						"urn:test:dogs:read",
					},
				},
			},
			wantResults: []storage.Client{
				expected,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "should filter clients by Scopes Union #3",
			args: args{
				filter: storage.ListClientsRequest{
					ScopesUnion: []string{
						"urn:test:cats:write",
						"urn:test:dogs:read",
					},
				},
			},
			wantResults: []storage.Client{
				expected,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "should filter clients by having at least one of the provided scopes when filtered by Scopes Union, ",
			args: args{
				filter: storage.ListClientsRequest{
					ScopesUnion: []string{
						"urn:test:dogs:write",
						"urn:test:cats:write",
					},
				},
			},
			wantResults: []storage.Client{
				expected,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "should return empty if no clients are found by Scopes Union",
			args: args{
				filter: storage.ListClientsRequest{
					ScopesIntersection: []string{
						"urn:test:dogs",
						"urn:test:cats",
					},
				},
			},
			wantResults: []storage.Client(nil),
			wantErr:     false,
			err:         nil,
		},
		{
			name: "should filter clients by contact",
			args: args{
				filter: storage.ListClientsRequest{
					Contact: expected.Contacts[0],
				},
			},
			wantResults: []storage.Client{
				expected,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "should return empty if no clients are found by contact",
			args: args{
				filter: storage.ListClientsRequest{
					Contact: "John Doe",
				},
			},
			wantResults: []storage.Client(nil),
			wantErr:     false,
			err:         nil,
		},
		{
			name: "should filter for public clients",
			args: args{
				filter: storage.ListClientsRequest{
					Public: true,
				},
			},
			wantResults: []storage.Client{
				expected,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "should filter for disabled clients",
			args: args{
				filter: storage.ListClientsRequest{
					Disabled: true,
				},
			},
			wantResults: []storage.Client(nil),
			wantErr:     false,
			err:         nil,
		},
		{
			name: "should filter for published clients",
			args: args{
				filter: storage.ListClientsRequest{
					Published: true,
				},
			},
			wantResults: []storage.Client{
				publishedClient,
			},
			wantErr: false,
			err:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResults, err := store.ClientManager.List(ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				AssertError(t, err, tt.err, "list should return an error")
				return
			}

			if !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("List():\ngot:  %#+v\nwant: %#+v\n", gotResults, tt.wantResults)
			}
		})
	}
}

func expectedClient() storage.Client {
	return storage.Client{
		ID:         uuid.NewString(),
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix() + 600,
		AllowedAudiences: []string{
			uuid.NewString(),
			uuid.NewString(),
		},
		AllowedRegions: []string{
			uuid.NewString(),
		},
		AllowedTenantAccess: []string{
			uuid.NewString(),
			uuid.NewString(),
		},
		GrantTypes: []string{
			string(fosite.AccessToken),
			string(fosite.RefreshToken),
			string(fosite.AuthorizeCode),
			string(fosite.IDToken),
		},
		ResponseTypes: []string{
			"code",
			"token",
		},
		Scopes: []string{
			"urn:test:cats:write",
			"urn:test:dogs:read",
		},
		Public:   true,
		Disabled: false,
		Name:     "Test Client",
		Secret:   "foobar",
		RedirectURIs: []string{
			"https://test.example.com",
		},
		Owner:             "Widgets Inc.",
		PolicyURI:         "https://test.example.com/policy",
		TermsOfServiceURI: "https://test.example.com/tos",
		ClientURI:         "https://app.example.com",
		LogoURI:           "https://app.example.com/favicon-128x128.png",
		Contacts: []string{
			"John Doe <j.doe@example.com>",
		},
		Published: false,
	}
}

func createClient(ctx context.Context, t *testing.T, store *mongo.Store, expected storage.Client) storage.Client {
	return createNewClient(t, ctx, store, expected)
}

func createNewClient(t *testing.T, ctx context.Context, store *mongo.Store, expected storage.Client) storage.Client {
	got, err := store.ClientManager.Create(ctx, expected)
	if err != nil {
		AssertError(t, err, nil, "create should return no database errors")
		t.FailNow()
	}

	if got.Secret == "" || got.Secret == expected.Secret {
		AssertError(t, got.Secret, "bcrypt encoded secret", "create should hash the secret")
		t.FailNow()
	}

	expected.ID = got.ID
	expected.CreateTime = got.CreateTime
	expected.UpdateTime = got.UpdateTime
	expected.Secret = got.Secret
	if !reflect.DeepEqual(got, expected) {
		AssertError(t, got, expected, "client not equal")
		t.FailNow()
	}

	return expected
}

func TestClientManager_Create(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	createClient(ctx, t, store, expectedClient())
}

func TestClientManager_Create_ShouldStoreCustomData(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expectedEntity := expectedClient()
	expectedData := expectedCustomData()

	// push in custom data
	if err := expectedEntity.Data.Marshal(expectedData); err != nil {
		AssertError(t, err, nil, "failed to marshal custom data")
	}

	// save client to mongo
	expectedEntity = createClient(ctx, t, store, expectedEntity)

	// extract raw save
	query := bson.M{
		"id": expectedEntity.ID,
	}
	var gotEntity storage.Client
	if err := store.DB.Collection(storage.EntityClients).FindOne(ctx, query).Decode(&gotEntity); err != nil {
		AssertError(t, err, nil, "expected client to exist")
	}

	// Test expectations
	AssertClientCustomData(t, gotEntity, expectedEntity, &TestCustomData{}, &expectedData)
}

func TestClientManager_Create_ShouldConflict(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createClient(ctx, t, store, expectedClient())
	_, err := store.ClientManager.Create(ctx, expected)
	if err == nil {
		AssertError(t, err, nil, "create should return an error on conflict")
	}
	if !errors.Is(err, storage.ErrResourceExists) {
		AssertError(t, err, nil, "create should return conflict")
	}
}

func TestClientManager_Get(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createClient(ctx, t, store, expectedClient())
	got, err := store.ClientManager.Get(ctx, expected.ID)
	if err != nil {
		AssertError(t, err, nil, "get should return no database errors")
	}
	if !reflect.DeepEqual(got, expected) {
		AssertError(t, got, expected, "client not equal")
	}
}

func TestClientManager_Get_ShouldReturnNotFound(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := fosite.ErrNotFound
	got, err := store.ClientManager.Get(ctx, "lolNotFound")
	if !errors.Is(err, expected) {
		AssertError(t, got, expected, "get should return not found")
	}
}

func TestClientManager_Get_ShouldReturnCustomData(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expectedEntity := expectedClient()
	expectedData := expectedCustomData()

	// push in custom data
	if err := expectedEntity.Data.Marshal(expectedData); err != nil {
		AssertError(t, err, nil, "failed to unmarshal custom data")
	}
	// save client to mongo
	expectedEntity = createClient(ctx, t, store, expectedEntity)

	// Get entity
	gotEntity, err := store.ClientManager.Get(ctx, expectedEntity.ID)
	if err != nil {
		AssertError(t, err, nil, "failed to get client")
	}

	// Test expectations
	AssertClientCustomData(t, gotEntity, expectedEntity, &TestCustomData{}, &expectedData)
}

func TestClientManager_Update(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createClient(ctx, t, store, expectedClient())
	// Perform an update...
	expected.Name = "something completely different!"

	got, err := store.ClientManager.Update(ctx, expected.ID, expected)
	if err != nil {
		AssertError(t, err, nil, "update should return no database errors")
	}

	if expected.UpdateTime == 0 {
		AssertError(t, got.UpdateTime, time.Now().Unix(), "update time was not set")
	}

	if expected.Secret != got.Secret {
		AssertError(t, got.Secret, expected.Secret, "secret should not change on update unless explicitly changed")
	}

	// override update time on expected with got. The time stamp received
	// should match time.Now().Unix() but due to the nature of time based
	// testing against time.Now().Unix(), it can fail on crossing over the
	// second boundary.
	expected.UpdateTime = got.UpdateTime
	if !reflect.DeepEqual(got, expected) {
		AssertError(t, got, expected, "client update object not equal")
	}
}

func TestClientManager_Update_ShouldChangePassword(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	newSecret := "s0methingElse!"
	expected := createClient(ctx, t, store, expectedClient())
	oldHash := expected.Secret

	// Perform a password update...
	expected.Secret = newSecret

	got, err := store.ClientManager.Update(ctx, expected.ID, expected)
	if err != nil {
		AssertError(t, err, nil, "update should return no database errors")
	}

	if expected.UpdateTime == 0 {
		AssertError(t, got.UpdateTime, time.Now().Unix(), "update time was not set")
	}

	if got.Secret == oldHash {
		AssertError(t, got.Secret, "new bcrypt hash", "secret was not updated")
	}

	if got.Secret == newSecret {
		AssertError(t, got.Secret, "new bcrypt hash", "secret was not hashed")
	}

	// Should authenticate against the new hash
	if err := store.Hasher.Compare(ctx, got.GetHashedSecret(), []byte(newSecret)); err != nil {
		AssertError(t, got.Secret, "bcrypt authenticate-able hash", "unable to authenticate with updated hash")
	}

	// override update time on expected with got. The time stamp received
	// should match time.Now().Unix() but due to the nature of time based
	// testing against time.Now().Unix(), it can fail on crossing over the
	// second boundary.
	expected.UpdateTime = got.UpdateTime
	// override expected secret as the assertions have passed above.
	expected.Secret = got.Secret

	if !reflect.DeepEqual(got, expected) {
		AssertError(t, got, expected, "client update object not equal")
	}
}

func TestClientManager_Update_ShouldReturnNotFound(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	_, err := store.ClientManager.Update(ctx, uuid.NewString(), expectedClient())
	if err == nil {
		AssertError(t, err, nil, "update should return an error on not found")
	}
	if !errors.Is(err, fosite.ErrNotFound) {
		AssertError(t, err, nil, "update should return not found")
	}
}

func TestClientManager_Update_ShouldUpdateCustomData(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expectedEntity := expectedClient()
	expectedData := expectedCustomData()

	// push in custom data
	if err := expectedEntity.Data.Marshal(expectedData); err != nil {
		AssertError(t, err, nil, "failed to unmarshal custom data")
	}
	// save client to mongo
	expectedEntity = createClient(ctx, t, store, expectedEntity)

	// Update custom data
	expectedData.Contact.Name = "John Doe"
	if err := expectedEntity.Data.Marshal(expectedData); err != nil {
		AssertError(t, err, nil, "failed to unmarshal custom data")
	}

	gotEntity, err := store.ClientManager.Update(ctx, expectedEntity.ID, expectedEntity)
	if err != nil {
		AssertError(t, err, nil, "failed to update client")
	}

	// Test expectations
	// override update time on expected with got. The time stamp received
	// should match time.Now().Unix() but due to the nature of time based
	// testing against time.Now().Unix(), it can fail on crossing over the
	// second boundary.
	expectedEntity.UpdateTime = gotEntity.UpdateTime
	AssertClientCustomData(t, gotEntity, expectedEntity, &TestCustomData{}, &expectedData)

	// Get record directly to verify struct passed back is persisting correctly
	// on update
	query := bson.M{
		"id": expectedEntity.ID,
	}
	gotEntity = storage.Client{}
	if err := store.DB.Collection(storage.EntityClients).FindOne(ctx, query).Decode(&gotEntity); err != nil {
		AssertError(t, err, nil, "expected client to exist")
	}

	// Test expectations
	AssertClientCustomData(t, gotEntity, expectedEntity, &TestCustomData{}, &expectedData)
}

func TestClientManager_Delete(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	expected := createClient(ctx, t, store, expectedClient())

	err := store.ClientManager.Delete(ctx, expected.ID)
	if err != nil {
		AssertError(t, err, nil, "delete should return no database errors")
	}

	// Double check that the original reference was deleted
	expectedErr := fosite.ErrNotFound
	got, err := store.ClientManager.Get(ctx, expected.ID)
	if !errors.Is(expectedErr, err) {
		AssertError(t, got, expectedErr, "get should return not found")
	}
}

func TestClientManager_Delete_ShouldReturnNotFound(t *testing.T) {
	store, ctx, teardown := setup(t)
	defer teardown()

	err := store.ClientManager.Delete(ctx, expectedClient().ID)
	if err == nil {
		AssertError(t, err, nil, "delete should return an error on not found")
	}
	if !errors.Is(err, fosite.ErrNotFound) {
		AssertError(t, err, nil, "delete should return not found")
	}
}
