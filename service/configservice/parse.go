/*
 * Copyright (C) distroy
 */

package configservice

import (
	"fmt"
	"log"
	"os"
	"path"
	"reflect"

	"github.com/distroy/git-go-tool/3rd/yaml"
	"github.com/distroy/git-go-tool/core/clonecore"
	"github.com/distroy/git-go-tool/core/flagcore"
	"github.com/distroy/git-go-tool/core/git"
	"github.com/distroy/git-go-tool/core/jsoncore"
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

	printf := func(fmt string, args ...interface{}) {
		if false {
			log.Printf(fmt, args...)
		}
	}

	defCfg := reflect.New(typ.Elem()).Interface()
	mergecore.Merge(defCfg, outCfg)

	merge := func(target, source interface{}) {
		tmp := clonecore.DeepClone(target)
		mergecore.Merge(tmp, source)
		reflect.ValueOf(source).Elem().Set(reflect.ValueOf(tmp).Elem())
	}

	fileCfg := reflect.New(typ.Elem()).Interface()
	cfgFilePath := MustGetConfigPath()
	mustUnmarshalFileWithField(fileCfg, cfgFilePath, fieldName)
	printf("file config: %s", jsoncore.MustMarshal(fileCfg))

	flagCfg := reflect.New(typ.Elem()).Interface()
	flagcore.EnableDefault(false)
	flagcore.MustParse(flagCfg)
	printf("flag config: %s", jsoncore.MustMarshal(flagCfg))

	merge(fileCfg, defCfg)
	printf("merge file config: %s", jsoncore.MustMarshal(defCfg))

	merge(flagCfg, defCfg)
	printf("merge flag config: %s", jsoncore.MustMarshal(defCfg))

	reflect.ValueOf(outCfg).Elem().Set(reflect.ValueOf(defCfg).Elem())
	printf("out config: %s", jsoncore.MustMarshal(outCfg))
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
