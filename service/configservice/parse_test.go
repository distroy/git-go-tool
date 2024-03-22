/*
 * Copyright (C) distroy
 */

package configservice

import (
	"reflect"
	"testing"

	"github.com/distroy/git-go-tool/core/ptrcore"
)

func Test_mustUnmarshalFileWithField(t *testing.T) {
	type PushConfig struct {
		PushUrl    string `yaml:"push-url"`
		ProjectUrl string `yaml:"project-url"`
	}

	type Config struct {
		Rate *float64    `yaml:"rate"  flag:"default:0.65; usage:the lowest coverage rate. range: [0, 1.0)"`
		Top  *int        `yaml:"top"   flag:"meta:N; default:10; usage:show the top <N> least coverage rage file only"`
		File *string     `yaml:"file"  flag:"meta:file; usage:the coverage file path, cannot be empty"`
		Push *PushConfig `yaml:"push"`
	}

	fieldName := "go-coverage"

	type args struct {
		res       interface{}
		cfgPath   string
		fieldName string
	}
	type want struct {
		ok  bool
		res interface{}
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "file not exist",
			args: args{
				res:       &Config{},
				cfgPath:   "./985c1399b90fa9feee4096951a67e042.yaml",
				fieldName: fieldName,
			},
			want: want{
				ok:  false,
				res: &Config{},
			},
		},
		{
			name: "empty file",
			args: args{
				res:       &Config{},
				cfgPath:   "./config_for_test_empty.yaml",
				fieldName: fieldName,
			},
			want: want{
				ok:  false,
				res: &Config{},
			},
		},
		{
			name: "no push domain",
			args: args{
				res:       &Config{},
				cfgPath:   "./config_for_test_no_push.yaml",
				fieldName: fieldName,
			},
			want: want{
				ok: true,
				res: &Config{
					Rate: ptrcore.NewFloat64(0.5),
					Top:  ptrcore.NewInt(100),
				},
			},
		},
		{
			name: "no go-coverage domain",
			args: args{
				res:       &Config{},
				cfgPath:   "./config_for_test_no_coverage.yaml",
				fieldName: fieldName,
			},
			want: want{
				ok: true,
				res: &Config{
					Push: &PushConfig{
						PushUrl: "https://github.com",
					},
				},
			},
		},
		{
			name: "both",
			args: args{
				res:       &Config{},
				cfgPath:   "./config_for_test_both.yaml",
				fieldName: fieldName,
			},
			want: want{
				ok: true,
				res: &Config{
					Rate: ptrcore.NewFloat64(0.5),
					Top:  ptrcore.NewInt(100),
					Push: &PushConfig{
						PushUrl: "https://github.com",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustUnmarshalFileWithField(tt.args.res, tt.args.cfgPath, tt.args.fieldName)
			if got != tt.want.ok {
				t.Errorf("mustUnmarshalFileWithField() = %v, want %v", got, tt.want.ok)
			}
			if !reflect.DeepEqual(tt.args.res, tt.want.res) {
				t.Errorf("mustUnmarshalFileWithField() output res = %v, want %v", tt.args.res, tt.want.res)
			}
		})
	}
}
