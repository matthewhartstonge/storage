package utils

import (
	// Standard Library Imports
	"reflect"
	"testing"
)

func TestAppendToStringSet(t *testing.T) {
	type args struct {
		strings []string
		items   []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Should return nil if nil slice given, and no items passed in",
			args: args{
				strings: nil,
				items:   nil,
			},
			want: nil,
		},
		{
			name: "Should return empty if no items given, and no items passed in",
			args: args{
				strings: []string{},
				items:   nil,
			},
			want: []string{},
		},
		{
			name: "Should add a value, given a nil slice",
			args: args{
				strings: nil,
				items:   []string{"cat"},
			},
			want: []string{"cat"},
		},
		{
			name: "Should add a value, given an empty slice",
			args: args{
				strings: []string{},
				items:   []string{"cat"},
			},
			want: []string{"cat"},
		},
		{
			name: "Should not add a value already in the slice",
			args: args{
				strings: []string{"cat"},
				items:   []string{"cat"},
			},
			want: []string{"cat"},
		},
		{
			name: "Should not add items already in the slice",
			args: args{
				strings: []string{"cat", "dog"},
				items:   []string{"cat", "dog"},
			},
			want: []string{"cat", "dog"},
		},
		{
			name: "Should add multiple items",
			args: args{
				strings: []string{},
				items:   []string{"cat", "dog", "gorilla"},
			},
			want: []string{"cat", "dog", "gorilla"},
		},
		{
			name: "Should ignore items already in the set, but add new additions",
			args: args{
				strings: []string{"cat", "dog"},
				items:   []string{"cat", "gorilla", "dog", "fish"},
			},
			want: []string{"cat", "dog", "gorilla", "fish"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppendToStringSet(tt.args.strings, tt.args.items...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendToStringSet()\ngot:  %#+v\nwant: %#+v\n", got, tt.want)
			}
		})
	}
}

func TestRemoveFromStringSet(t *testing.T) {
	type args struct {
		strings []string
		items   []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Should return nil if nil slice given, and no items passed in",
			args: args{
				strings: nil,
				items:   nil,
			},
			want: nil,
		},
		{
			name: "Should return empty if no items given, and no items passed in",
			args: args{
				strings: []string{},
				items:   nil,
			},
			want: []string{},
		},
		{
			name: "Should not remove a value, given a nil slice",
			args: args{
				strings: nil,
				items:   []string{"cat"},
			},
			want: nil,
		},
		{
			name: "Should not remove a value, given an empty slice",
			args: args{
				strings: []string{},
				items:   []string{"cat"},
			},
			want: []string{},
		},
		{
			name: "Should remove a value",
			args: args{
				strings: []string{"cat"},
				items:   []string{"cat"},
			},
			want: []string{},
		},
		{
			name: "Should not remove a value, if it doesn't exist",
			args: args{
				strings: []string{"cat"},
				items:   []string{"dog"},
			},
			want: []string{"cat"},
		},
		{
			name: "Should remove multiple items",
			args: args{
				strings: []string{"cat", "dog"},
				items:   []string{"cat", "dog"},
			},
			want: []string{},
		},
		{
			name: "Should remove all instances of an item",
			args: args{
				strings: []string{"cat", "dog", "cat"},
				items:   []string{"cat"},
			},
			want: []string{"dog"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveFromStringSet(tt.args.strings, tt.args.items...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveFromStringSet()\ngot:  %#+v\nwant: %#+v\n", got, tt.want)
			}
		})
	}
}
