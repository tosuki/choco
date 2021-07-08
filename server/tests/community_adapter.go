package tests

import (
	"choco/server/internals/adapters"
	"choco/server/internals/models"
	"testing"
)

func testCommunityAdapterCreate(t *testing.T, adapter adapters.CommunityAdapter, object *models.Community) *models.Community {
	err := adapter.Add(object)

	if err != nil {
		t.Errorf("Not expected an error to create a community in the database: %s", err)
	}

	return object
}

func testCommunityAdapterGet(t *testing.T, adapter adapters.CommunityAdapter, id string) *models.Community {
	community, err := adapter.Get(id)

	if err != nil {
		t.Errorf("Not expected an error to get a community that should exists in the database: %s", err)
	}

	return community
}

func testCommunityAdapterInvalidGet(t *testing.T, adapter adapters.CommunityAdapter, id string) {
	_, err := adapter.Get(id)

	if err == nil {
		t.Errorf("Expected an error to get a community that shouldn't exists in the database: %s", err)
	}
}

func testCommunityAdapterName(t *testing.T, adapter adapters.CommunityAdapter, name string) *models.Community {
	community, err := adapter.Name(name)

	if err != nil {
		t.Errorf("Not expected an error to get a community by its name that should exists in the database: %s", err)
	}

	return community
}

func testCommunityAdapterInvalidName(t *testing.T, adapter adapters.CommunityAdapter, name string) {
	_, err := adapter.Name(name)
	
	if err == nil {
		t.Errorf("Expected an error to get a community by its name that shouldn't exists in the database: %s", err)
	}
}

func testCommunityAdapter(t *testing.T, adapter adapters.CommunityAdapter) {
	user, userErr := models.NewUser("choco", "choco@choco", []byte("choco"), models.ROOT_PERMISSION)
	
	if userErr != nil {
		t.Errorf("Not expected an error to create the user model to test the community adapter: %s", userErr)
	}

	community, commErr := models.NewCommunity("Choco", "Group to test things a.a", user.ID, true, false)
	
	
	if commErr != nil {
		t.Errorf("Not expected an error to create the community model to test the community adapter: %s", commErr)
	}

	testCommunityAdapterCreate(t, adapter, community)
	testCommunityAdapterGet(t, adapter, community.ID)
	testCommunityAdapterInvalidGet(t, adapter, "pdkawpodkad")
	testCommunityAdapterName(t, adapter, community.Name)
	testCommunityAdapterInvalidName(t, adapter, "dpkapkapapapapapapa")
}