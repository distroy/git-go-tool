/*
 * Copyright (C) distroy
 */

package resultservice

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/distroy/git-go-tool/core/jsoncore"
	"github.com/distroy/git-go-tool/obj/resultobj"
	"github.com/distroy/git-go-tool/service/resultservice/internal/validate"
)

func Push(url string, result *resultobj.Result) error {
	err := push(url, result)
	if err != nil {
		log.Printf("push result fail. url:%s, err:%v", url, err)
		return err
	}
	return nil
}

func push(url string, result *resultobj.Result) error {
	if url == "" {
		return nil
	}

	if !validate.Result(result) {
		return nil
	}

	raw := jsoncore.MustMarshal(result)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(raw))
	if err != nil {
		return err
		// return fmt.Errorf("new http request for push result fail. url:%s, err:%v", url, err)
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
		// return fmt.Errorf("call http request fail. url:%s, err:%v", url, err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("read http response body fail. url:%s, err:%v", url, err)
	}

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
		// return fmt.Errorf("read http response body fail. url:%s, err:%v", url, err)
	}

	log.Printf("push result succ. url:%s, status:%d, body:%s", url, rsp.StatusCode, rspBody)
	return nil
}
