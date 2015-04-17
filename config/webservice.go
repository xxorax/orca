package config

import (
	"github.com/clusterit/orca/auth"
	"github.com/clusterit/orca/rest"
	"github.com/clusterit/orca/users"
	"gopkg.in/emicklei/go-restful.v1"
)

type ConfigService struct {
	Zone   string
	Auth   auth.Auther
	Users  users.Users
	Config Configer
}

func (t *ConfigService) Shutdown() error {
	return nil
}

func (t *ConfigService) Register(root string, c *restful.Container) {
	ws := new(restful.WebService)

	mgr := users.CheckUser(t.Auth, t.Users, users.ManagerRoles)

	ws.
		Path(root + "configuration").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PUT("/{zone}/mgrConfig").To(mgr(t.putManagerConfig)).
		Doc("Store the ManagerConfig config for a given zone").
		Param(ws.PathParameter("zone", "the zone to put the ManagerConfig config to").DataType("string")).
		Operation("putManagerConfig").
		Reads(ManagerConfig{}))
	ws.Route(ws.GET("/{zone}/mgrConfig").To(mgr(t.getManagerConfig)).
		Doc("Get the ManagerConfig config for a given zone").
		Param(ws.PathParameter("zone", "the zone to read the ManagerConfig config from").DataType("string")).
		Operation("getManagerConfig").
		Writes(ManagerConfig{}))
	ws.Route(ws.PUT("/{zone}/gateway").To(mgr(t.putGateway)).
		Doc("Store the Gateway config for a given zone").
		Param(ws.PathParameter("zone", "the stzoneage to put the Gateway config to").DataType("string")).
		Operation("putGateway").
		Reads(Gateway{}))
	ws.Route(ws.GET("/{zone}/gateway").To(mgr(t.getGateway)).
		Doc("Get the Gateway config for a given zone").
		Param(ws.PathParameter("zone", "the zone to reamd the JWT from").DataType("string")).
		Operation("getGateway").
		Writes(Gateway{}))
	ws.Route(ws.GET("/zones").To(mgr(t.getZones)).
		Doc("Get all current configured zones").
		Operation("getZones").
		Writes([]string{}))
	ws.Route(ws.GET("/zone").To(t.getZone).
		Doc("Get the current zone").
		Operation("getZone").
		Writes(""))

	c.Add(ws)
}

func (t *ConfigService) getZones(u *users.User, rq *restful.Request, rsp *restful.Response) {
	rest.HandleEntity(t.Config.Zones())(rq, rsp)
}

func (t *ConfigService) getZone(rq *restful.Request, rsp *restful.Response) {
	rest.HandleEntity(t.Zone, nil)(rq, rsp)
}

func (t *ConfigService) putManagerConfig(u *users.User, rq *restful.Request, rsp *restful.Response) {
	var mc ManagerConfig
	z := rq.PathParameter("zone")
	err := rq.ReadEntity(&mc)
	if err != nil {
		rest.HandleError(err, rsp)
		return
	}
	if err := t.Config.PutManagerConfig(z, mc); err != nil {
		rest.HandleError(err, rsp)
		return
	}
	rsp.WriteEntity(mc)
}

func (t *ConfigService) getManagerConfig(u *users.User, rq *restful.Request, rsp *restful.Response) {
	z := rq.PathParameter("zone")
	rest.HandleEntity(t.Config.GetManagerConfig(z))(rq, rsp)
}

func (t *ConfigService) putGateway(u *users.User, rq *restful.Request, rsp *restful.Response) {
	var gw Gateway
	z := rq.PathParameter("zone")
	err := rq.ReadEntity(&gw)
	if err != nil {
		rest.HandleError(err, rsp)
		return
	}
	if err := t.Config.PutGateway(z, gw); err != nil {
		rest.HandleError(err, rsp)
		return
	}
	rsp.WriteEntity(gw)
}

func (t *ConfigService) getGateway(u *users.User, rq *restful.Request, rsp *restful.Response) {
	z := rq.PathParameter("zone")
	rest.HandleEntity(t.Config.GetGateway(z))(rq, rsp)
}
