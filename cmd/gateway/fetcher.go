package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/clusterit/orca/cmd"
	"github.com/clusterit/orca/etcd"
	"github.com/clusterit/orca/users"

	"github.com/jmcvetta/napping"
)

type UserFetcher interface {
	UserByKey(key string) (*users.User, error)
	CheckToken(uid, token string, maxtime int) error
}

type httpFetcher struct {
	managerConfig etcd.Configurator
}

func NewHttpFetcher(cc *etcd.Cluster) (UserFetcher, error) {
	cfg, err := cc.NewManager()
	if err != nil {
		return nil, err
	}
	return &httpFetcher{managerConfig: cfg}, nil
}

func (hf *httpFetcher) CheckToken(uid, token string, maxtime int) error {
	urls, err := hf.managerConfig.GetValues("/" + cmd.ManagerService)
	if err != nil {
		return err
	}
	if len(urls) < 1 {
		return fmt.Errorf("no managers registered in configuration")
	}
	var res string
	for _, url := range urls {
		serviceUrl := fmt.Sprintf("%s/users/%s/%s/check?maxtime=%d", url, uid, token, maxtime)
		if strings.HasSuffix(url, "/") {
			serviceUrl = fmt.Sprintf("%susers/%s/%s/check?maxtime=%d", url, uid, token, maxtime)
		}
		r := napping.Request{
			Url:    serviceUrl,
			Method: "GET",
			Result: &res,
			Header: &http.Header{},
		}
		r.Header.Add("Content-Type", "application/json")
		r.Header.Add("Accept", "application/json")
		resp, err := napping.Send(&r)
		if err != nil {
			continue
		}
		if resp.Status()/100 == 5 {
			// the rest call returned an unknown error
			continue
		}
		if resp.Status() != 200 {
			return fmt.Errorf("HTTP %d: %s", resp.Status(), resp.RawText())
		}
		return nil
	}
	return fmt.Errorf("no working manager found in configuration")
}

func (hf *httpFetcher) UserByKey(key string) (*users.User, error) {
	var u users.User
	urls, err := hf.managerConfig.GetValues("/" + cmd.ManagerService)
	if err != nil {
		return nil, err
	}
	if len(urls) < 1 {
		return nil, fmt.Errorf("no managers registered in configuration")
	}
	for _, url := range urls {
		serviceUrl := fmt.Sprintf("%s/users/pubkey", url)
		if strings.HasSuffix(url, "/") {
			serviceUrl = fmt.Sprintf("%susers/pubkey", url)
		}
		r := napping.Request{
			Url:     serviceUrl,
			Method:  "POST",
			Payload: key,
			Result:  &u,
			Header:  &http.Header{},
		}
		r.Header.Add("Content-Type", "application/json")
		r.Header.Add("Accept", "application/json")
		resp, err := napping.Send(&r)
		if err != nil {
			continue
		}
		if resp.Status()/100 == 5 {
			// the rest call returned an unknown error
			continue
		}
		if resp.Status() != 200 {
			return nil, fmt.Errorf("%s: HTTP %d: %s", serviceUrl, resp.Status(), resp.RawText())
		}
		return &u, nil
	}
	return nil, fmt.Errorf("no working manager found in configuration")
}
