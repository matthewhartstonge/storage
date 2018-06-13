package mongo

//
//import (
//	"context"
//	"reflect"
//	"testing"
//	"time"
//
//	"github.com/globalsign/mgo"
//	"github.com/matthewhartstonge/storage"
//	"github.com/pborman/uuid"
//)
//
//var cacheExpected = storage.SessionCache{
//	ID:         uuid.New(),
//	CreateTime: time.Now().Unix(),
//	UpdateTime: time.Now().Unix() + 600,
//	Signature:  "Yhte@ensa#ei!+suu$re%sta^viik&oss*aha(joaisiaut)ta-is+ie%to_n==",
//}
//
//func Test_cacheMongoManager_Create(t *testing.T) {
//	type fields struct {
//		db *mgo.Database
//	}
//	type args struct {
//		ctx         context.Context
//		entityName  string
//		cacheObject storage.SessionCache
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    storage.SessionCache
//		wantErr bool
//	}{
//		{
//			name: "Should create a cache object",
//			fields: fields{
//				db: mongoStore.db,
//			},
//			args: args{
//				ctx:         context.Background(),
//				entityName:  storage.EntityCacheAccessTokens,
//				cacheObject: cacheExpected,
//			},
//			want:    cacheExpected,
//			wantErr: false,
//		},
//		{
//			name: "Should conflict on create",
//			fields: fields{
//				db: mongoStore.db,
//			},
//			args: args{
//				ctx:         context.Background(),
//				entityName:  storage.EntityCacheAccessTokens,
//				cacheObject: cacheExpected,
//			},
//			want:    storage.SessionCache{},
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &cacheMongoManager{
//				db: tt.fields.db,
//			}
//			got, err := c.Create(tt.args.ctx, tt.args.entityName, tt.args.cacheObject)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("cacheMongoManager.Create() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("cacheMongoManager.Create() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_cacheMongoManager_Get(t *testing.T) {
//	type fields struct {
//		db *mgo.Database
//	}
//	type args struct {
//		ctx        context.Context
//		entityName string
//		key        string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    storage.SessionCache
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &cacheMongoManager{
//				db: tt.fields.db,
//			}
//			got, err := c.Get(tt.args.ctx, tt.args.entityName, tt.args.key)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("cacheMongoManager.Get() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("cacheMongoManager.Get() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_cacheMongoManager_Update(t *testing.T) {
//	type fields struct {
//		db *mgo.Database
//	}
//	type args struct {
//		ctx                context.Context
//		entityName         string
//		updatedCacheObject storage.SessionCache
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    storage.SessionCache
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &cacheMongoManager{
//				db: tt.fields.db,
//			}
//			got, err := c.Update(tt.args.ctx, tt.args.entityName, tt.args.updatedCacheObject)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("cacheMongoManager.Update() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("cacheMongoManager.Update() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_cacheMongoManager_Delete(t *testing.T) {
//	type fields struct {
//		db *mgo.Database
//	}
//	type args struct {
//		ctx        context.Context
//		entityName string
//		key        string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &cacheMongoManager{
//				db: tt.fields.db,
//			}
//			if err := c.Delete(tt.args.ctx, tt.args.entityName, tt.args.key); (err != nil) != tt.wantErr {
//				t.Errorf("cacheMongoManager.Delete() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func Test_cacheMongoManager_DeleteByValue(t *testing.T) {
//	type fields struct {
//		db *mgo.Database
//	}
//	type args struct {
//		ctx        context.Context
//		entityName string
//		value      string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &cacheMongoManager{
//				db: tt.fields.db,
//			}
//			if err := c.DeleteByValue(tt.args.ctx, tt.args.entityName, tt.args.value); (err != nil) != tt.wantErr {
//				t.Errorf("cacheMongoManager.DeleteByValue() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
