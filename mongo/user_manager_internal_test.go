package mongo

import (
	"testing"

	"github.com/matthewhartstonge/storage"
)

func TestUserMongoManager_ImplementsStorageConfigurer(t *testing.T) {
	u := &userMongoManager{}

	var i interface{} = u
	if _, ok := i.(storage.Configurer); !ok {
		t.Error("userMongoManager does not implement interface storage.Configurer")
	}
}

func TestUserMongoManager_ImplementsStorageAuthUserMigrator(t *testing.T) {
	u := &userMongoManager{}

	var i interface{} = u
	if _, ok := i.(storage.AuthUserMigrator); !ok {
		t.Error("userMongoManager does not implement interface storage.AuthUserMigrator")
	}
}

func TestUserMongoManager_ImplementsStorageUserStorer(t *testing.T) {
	u := &userMongoManager{}

	var i interface{} = u
	if _, ok := i.(storage.UserStorer); !ok {
		t.Error("userMongoManager does not implement interface storage.UserStorer")
	}
}

func TestUserMongoManager_ImplementsStorageUserManager(t *testing.T) {
	u := &userMongoManager{}

	var i interface{} = u
	if _, ok := i.(storage.UserManager); !ok {
		t.Error("userMongoManager does not implement interface storage.UserManager")
	}
}
