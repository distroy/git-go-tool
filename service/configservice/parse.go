/*
 * Copyright (C) distroy
 */

package configservice

import (
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/distroy/git-go-tool/3rd/yaml"
	"github.com/distroy/git-go-tool/core/flagcore"
	"github.com/distroy/git-go-tool/core/git"
	"github.com/distroy/git-go-tool/core/mergecore"
)

type any interface{}

func MustGetConfigPath() string {
	gitRoot := git.MustGetRootDir()
	return path.Join(gitRoot, ".git-go-tool/config.yaml")
}

func MustParse(cfg any, fieldName string) {
	typ := reflect.TypeOf(cfg)

	flags := reflect.New(typ.Elem()).Interface()
	flagcore.EnableDefault(false)
	flagcore.MustParse(flags)

	tmp := reflect.New(typ.Elem()).Interface()
	ok := mustUnmarshalFileWithField(tmp, fieldName)
	if ok {
		mergecore.Merge(cfg, tmp)
	}
	mergecore.Merge(cfg, flags)
	// log.Printf("config: %s", jsoncore.MustMarshal(cfg))
}

func mustUnmarshalFileWithField(res any, fieldName string) bool {
	cfgPath := MustGetConfigPath()
	f, err := os.Open(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(fmt.Errorf("open config file fail. file:%s, err:%v", cfgPath, err))
	}
	defer f.Close()

	m := make(map[string]interface{})
	d := yaml.NewDecoder(f)
	if err := d.Decode(&m); err != nil {
		panic(fmt.Errorf("decode config file fail. file:%s, err:%v", cfgPath, err))
	}

	v, ok := m[fieldName]
	if !ok {
		return false
	}

	s, _ := yaml.Marshal(v)

	if err := yaml.Unmarshal(s, res); err != nil {
		panic(fmt.Errorf("decode config file fail. file:%s, field:%s, err:%v",
			cfgPath, fieldName, err))
	}

	return true
}
