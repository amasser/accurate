package sessionmanager

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/accurateproject/accurate/config"
	"github.com/accurateproject/accurate/engine"
	"github.com/accurateproject/accurate/utils"
	"github.com/accurateproject/kamevapi"
	"github.com/accurateproject/rpcclient"
)

func NewKamailioSessionManager(smKamCfg *config.SmKamailio, rater, cdrsrv rpcclient.RpcClientConnection, timezone string) (*KamailioSessionManager, error) {
	ksm := &KamailioSessionManager{cfg: smKamCfg, rater: rater, cdrsrv: cdrsrv, timezone: timezone, conns: make(map[string]*kamevapi.KamEvapi), sessions: NewSessions()}
	return ksm, nil
}

type KamailioSessionManager struct {
	cfg      *config.SmKamailio
	rater    rpcclient.RpcClientConnection
	cdrsrv   rpcclient.RpcClientConnection
	timezone string
	conns    map[string]*kamevapi.KamEvapi
	sessions *Sessions
}

func (self *KamailioSessionManager) onCgrAuth(evData []byte, connId string) {
	kev, err := NewKamEvent(evData)
	if err != nil {
		utils.Logger.Info(fmt.Sprintf("<SM-Kamailio> ERROR unmarshalling event: %s, error: %s", evData, err.Error()))
		return
	}
	if kev.GetReqType(utils.META_DEFAULT) == utils.META_NONE { // Do not process this request
		return
	}
	if kev.MissingParameter(self.timezone) {
		if kar, err := kev.AsKamAuthReply(0.0, "", utils.ErrMandatoryIeMissing); err != nil {
			utils.Logger.Error(fmt.Sprintf("<SM-Kamailio> Failed building auth reply %s", err.Error()))
		} else if err = self.conns[connId].Send(kar.String()); err != nil {
			utils.Logger.Error(fmt.Sprintf("<SM-Kamailio> Failed sending auth reply %s", err.Error()))
		}
		return
	}
	var remainingDuration float64
	var errMaxSession error
	if errMaxSession = self.rater.Call("Responder.GetDerivedMaxSessionTime", kev.AsStoredCdr(self.Timezone()), &remainingDuration); errMaxSession != nil {
		utils.Logger.Error(fmt.Sprintf("<SM-Kamailio> Could not get max session time, error: %s", errMaxSession.Error()))
	}
	var supplStr string
	var errSuppl error
	if kev.ComputeLcr() {
		if supplStr, errSuppl = self.getSuppliers(kev); errSuppl != nil {
			utils.Logger.Error(fmt.Sprintf("<SM-Kamailio> Could not get suppliers, error: %s", errSuppl.Error()))
		}
	}
	if errMaxSession == nil { // Overwrite the error from maxSessionTime with the one from suppliers if nil
		errMaxSession = errSuppl
	}
	if kar, err := kev.AsKamAuthReply(remainingDuration, supplStr, errMaxSession); err != nil {
		utils.Logger.Error(fmt.Sprintf("<SM-Kamailio> Failed building auth reply %s", err.Error()))
	} else if err = self.conns[connId].Send(kar.String()); err != nil {
		utils.Logger.Error(fmt.Sprintf("<SM-Kamailio> Failed sending auth reply %s", err.Error()))
	}
}

func (self *KamailioSessionManager) onCgrLcrReq(evData []byte, connId string) {
	kev, err := NewKamEvent(evData)
	if err != nil {
		utils.Logger.Info(fmt.Sprintf("<SM-Kamailio> ERROR unmarshalling event: %s, error: %s", string(evData), err.Error()))
		return
	}
	supplStr, err := self.getSuppliers(kev)
	kamLcrReply, errReply := kev.AsKamAuthReply(-1.0, supplStr, err)
	kamLcrReply.Event = CGR_LCR_REPLY // Hit the CGR_LCR_REPLY event route on Kamailio side
	if errReply != nil {
		utils.Logger.Error(fmt.Sprintf("<SM-Kamailio> Failed building auth reply %s", errReply.Error()))
	} else if err = self.conns[connId].Send(kamLcrReply.String()); err != nil {
		utils.Logger.Error(fmt.Sprintf("<SM-Kamailio> Failed sending lcr reply %s", err.Error()))
	}
}

func (self *KamailioSessionManager) getSuppliers(kev KamEvent) (string, error) {
	cd, err := kev.AsCallDescriptor()
	cd.UniqueID = kev.GetUniqueID(self.timezone)
	if err != nil {
		utils.Logger.Info(fmt.Sprintf("<SM-Kamailio> LCR_PREPROCESS_ERROR error: %s", err.Error()))
		return "", errors.New("LCR_PREPROCESS_ERROR")
	}
	var lcr engine.LCRCost
	if err = self.Rater().Call("Responder.GetLCR", &engine.AttrGetLcr{CallDescriptor: cd}, &lcr); err != nil {
		utils.Logger.Info(fmt.Sprintf("<SM-Kamailio> LCR_API_ERROR error: %s", err.Error()))
		return "", errors.New("LCR_API_ERROR")
	}
	if lcr.HasErrors() {
		lcr.LogErrors()
		return "", errors.New("LCR_COMPUTE_ERROR")
	}
	return lcr.SuppliersString()
}

func (self *KamailioSessionManager) onCallStart(evData []byte, connId string) {
	kamEv, err := NewKamEvent(evData)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("<SM-Kamailio> ERROR unmarshalling event: %s, error: %s", evData, err.Error()))
		return
	}
	if kamEv.GetReqType(utils.META_DEFAULT) == utils.META_NONE { // Do not process this request
		return
	}
	if kamEv.MissingParameter(self.timezone) {
		self.DisconnectSession(kamEv, connId, utils.ErrMandatoryIeMissing.Error())
		return
	}
	s := NewSession(kamEv, connId, self)
	if s != nil {
		self.sessions.indexSession(s)
	}
}

func (self *KamailioSessionManager) onCallEnd(evData []byte, connId string) {
	kev, err := NewKamEvent(evData)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("<SM-Kamailio> ERROR unmarshalling event: %s, error: %s", evData, err.Error()))
		return
	}
	if kev.GetReqType(utils.META_DEFAULT) == utils.META_NONE { // Do not process this request
		return
	}
	if kev.MissingParameter(self.timezone) {
		utils.Logger.Error(fmt.Sprintf("<SM-Kamailio> Mandatory IE missing out of event: %+v", kev))
	}
	go self.ProcessCdr(kev.AsStoredCdr(self.Timezone()))
	s := self.sessions.getSession(kev.GetUUID())
	if s == nil { // Not handled by us
		return
	}
	if err := self.sessions.removeSession(s, kev); err != nil {
		utils.Logger.Error(err.Error())
	}
}

func (self *KamailioSessionManager) Connect() error {
	var err error
	eventHandlers := map[*regexp.Regexp][]func([]byte, string){
		regexp.MustCompile("CGR_AUTH_REQUEST"): []func([]byte, string){self.onCgrAuth},
		regexp.MustCompile("CGR_LCR_REQUEST"):  []func([]byte, string){self.onCgrLcrReq},
		regexp.MustCompile("CGR_CALL_START"):   []func([]byte, string){self.onCallStart},
		regexp.MustCompile("CGR_CALL_END"):     []func([]byte, string){self.onCallEnd},
	}
	errChan := make(chan error)
	for _, connCfg := range self.cfg.EvapiConns {
		connID := utils.GenUUID()
		if self.conns[connID], err = kamevapi.NewKamEvapi(connCfg.Address, connID, connCfg.Reconnects, eventHandlers); err != nil {
			return err
		}
		go func() { // Start reading in own goroutine, return on error
			if err := self.conns[connID].ReadEvents(); err != nil {
				errChan <- err
			}
		}()
	}
	err = <-errChan // Will keep the Connect locked until the first error in one of the connections
	return err
}

func (self *KamailioSessionManager) DisconnectSession(ev engine.Event, connId, notify string) error {
	sessionIds := ev.GetSessionIds()
	disconnectEv := &KamSessionDisconnect{Event: CGR_SESSION_DISCONNECT, HashEntry: sessionIds[0], HashId: sessionIds[1], Reason: notify}
	if err := self.conns[connId].Send(disconnectEv.String()); err != nil {
		utils.Logger.Error(fmt.Sprintf("<SM-Kamailio> Failed sending disconnect request, error %s, connection id: %s", err.Error(), connId))
		return err
	}
	return nil
}

func (self *KamailioSessionManager) DebitInterval() time.Duration {
	return self.cfg.DebitInterval.D()
}
func (self *KamailioSessionManager) CdrSrv() rpcclient.RpcClientConnection {
	return self.cdrsrv
}
func (self *KamailioSessionManager) Rater() rpcclient.RpcClientConnection {
	return self.rater
}

func (self *KamailioSessionManager) ProcessCdr(cdr *engine.CDR) error {
	if !*self.cfg.CreateCdr {
		return nil
	}
	var reply string
	if err := self.cdrsrv.Call("CdrsV1.ProcessCDR", cdr, &reply); err != nil {
		utils.Logger.Error(fmt.Sprintf("<SM-Kamailio> Failed processing CDR, uniqueid: %s, accid: %s, error: <%s>", cdr.UniqueID, cdr.OriginID, err.Error()))
	}
	return nil
}

func (sm *KamailioSessionManager) WarnSessionMinDuration(sessionUuid, connId string) {
}

func (self *KamailioSessionManager) Shutdown() error {
	return nil
}

func (self *KamailioSessionManager) Sessions() []*Session {
	return self.sessions.getSessions()
}

func (self *KamailioSessionManager) SyncSessions() error {
	return nil
}

func (self *KamailioSessionManager) Timezone() string {
	return self.timezone
}
