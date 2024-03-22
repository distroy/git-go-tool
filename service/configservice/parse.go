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

const (
	fieldName_Push = "push"
)

func MustGetConfigPath() string {
	gitRoot := git.MustGetRootDir()
	return path.Join(gitRoot, ".git-go-tool/config.yaml")
}

func MustParse(outCfg interface{}, fieldName string) {
	typ := reflect.TypeOf(outCfg)

	flagCfg := reflect.New(typ.Elem()).Interface()
	flagcore.EnableDefault(false)
	flagcore.MustParse(flagCfg)

	fileCfg := reflect.New(typ.Elem()).Interface()
	cfgFilePath := MustGetConfigPath()
	ok := mustUnmarshalFileWithField(fileCfg, cfgFilePath, fieldName)
	if ok {
		mergecore.Merge(outCfg, fileCfg)
	}
	mergecore.Merge(outCfg, flagCfg)
	// log.Printf("config: %s", jsoncore.MustMarshal(cfg))
}

func mustUnmarshalFileWithField(res interface{}, cfgPath, fieldName string) bool {
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

	v0, ok0 := m[fieldName]
	v1, ok1 := m[fieldName_Push]
	if !ok0 && !ok1 {
		return false
	}

	if v0 == nil {
		v0 = map[string]interface{}{}
		m[fieldName] = v0
	}

	m0 := v0.(map[string]interface{})
	delete(m0, fieldName_Push)
	if v1 != nil {
		m0[fieldName_Push] = v1
	}

	s, _ := yaml.Marshal(v0)

	if err := yaml.Unmarshal(s, res); err != nil {
		panic(fmt.Errorf("decode config file fail. file:%s, field:%s, err:%v",
			cfgPath, fieldName, err))
	}

	return true
}
