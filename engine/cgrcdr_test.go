package engine

import (
	"reflect"
	"testing"
	"time"

	"github.com/accurateproject/accurate/dec"
	"github.com/accurateproject/accurate/utils"
)

/*
curl --data "OriginID=asbfdsaf&OriginHost=192.168.1.1&RequestType=rated&direction=*out&tenant=cgrates.org&tor=call&account=1001&subject=1001&destination=1002&time_answer=1383813746&duration=10&field_extr1=val_extr1&fieldextr2=valextr2" http://ipbxdev:2080/cgr
*/

func TestCgrCdrInterfaces(t *testing.T) {
	var _ RawCdr = make(CgrCdr)
}

func TestCgrCdrAsCDR(t *testing.T) {
	cgrCdr := CgrCdr{utils.TOR: utils.VOICE, utils.ACCID: "dsafdsaf", utils.CDRHOST: "192.168.1.1", utils.CDRSOURCE: "internal_test", utils.REQTYPE: utils.META_RATED,
		utils.DIRECTION: utils.OUT,
		utils.TENANT:    "cgrates.org", utils.CATEGORY: "call",
		utils.ACCOUNT: "1001", utils.SUBJECT: "1001", utils.DESTINATION: "1002", utils.SETUP_TIME: "2013-11-07T08:42:20Z", utils.ANSWER_TIME: "2013-11-07T08:42:26Z",
		utils.USAGE: "10", utils.SUPPLIER: "SUPPL1", "field_extr1": "val_extr1", "fieldextr2": "valextr2"}
	setupTime, _ := utils.ParseTimeDetectLayout(cgrCdr[utils.SETUP_TIME], "")
	expctRtCdr := &CDR{UniqueID: utils.Sha1(cgrCdr[utils.ACCID], setupTime.String()), ToR: utils.VOICE, OriginID: cgrCdr[utils.ACCID], OriginHost: cgrCdr[utils.CDRHOST],
		Source:      cgrCdr[utils.CDRSOURCE],
		RequestType: cgrCdr[utils.REQTYPE],
		Direction:   cgrCdr[utils.DIRECTION], Tenant: cgrCdr[utils.TENANT], Category: cgrCdr[utils.CATEGORY], Account: cgrCdr[utils.ACCOUNT], Subject: cgrCdr[utils.SUBJECT],
		Destination: cgrCdr[utils.DESTINATION], SetupTime: time.Date(2013, 11, 7, 8, 42, 20, 0, time.UTC), AnswerTime: time.Date(2013, 11, 7, 8, 42, 26, 0, time.UTC),
		Usage: time.Duration(10) * time.Second, Supplier: "SUPPL1",
		ExtraFields: map[string]string{"field_extr1": "val_extr1", "fieldextr2": "valextr2"}, Cost: dec.NewFloat(-1)}
	if CDR := cgrCdr.AsStoredCdr(""); !reflect.DeepEqual(expctRtCdr, CDR) {
		t.Errorf("Expecting %v, received: %v", expctRtCdr, CDR)
	}
}

// Make sure the replicated CDR matches the expected CDR
func TestReplicatedCgrCdrAsCDR(t *testing.T) {
	cgrCdr := CgrCdr{utils.UniqueID: "164b0422fdc6a5117031b427439482c6a4f90e41", utils.TOR: utils.VOICE, utils.ACCID: "dsafdsaf", utils.CDRHOST: "192.168.1.1",
		utils.CDRSOURCE: "internal_test", utils.REQTYPE: utils.META_RATED,
		utils.DIRECTION: utils.OUT, utils.TENANT: "cgrates.org", utils.CATEGORY: "call",
		utils.ACCOUNT: "1001", utils.SUBJECT: "1001", utils.DESTINATION: "1002", utils.SETUP_TIME: "2013-11-07T08:42:20Z", utils.PDD: "0.200", utils.ANSWER_TIME: "2013-11-07T08:42:26Z",
		utils.USAGE: "10", utils.SUPPLIER: "SUPPL1", utils.DISCONNECT_CAUSE: "NORMAL_CLEARING", utils.COST: "0.12", utils.RATED: "true", "field_extr1": "val_extr1", "fieldextr2": "valextr2"}
	expctRtCdr := &CDR{UniqueID: cgrCdr[utils.UniqueID],
		ToR:             cgrCdr[utils.TOR],
		OriginID:        cgrCdr[utils.ACCID],
		OriginHost:      cgrCdr[utils.CDRHOST],
		Source:          cgrCdr[utils.CDRSOURCE],
		RequestType:     cgrCdr[utils.REQTYPE],
		Direction:       cgrCdr[utils.DIRECTION],
		Tenant:          cgrCdr[utils.TENANT],
		Category:        cgrCdr[utils.CATEGORY],
		Account:         cgrCdr[utils.ACCOUNT],
		Subject:         cgrCdr[utils.SUBJECT],
		Destination:     cgrCdr[utils.DESTINATION],
		SetupTime:       time.Date(2013, 11, 7, 8, 42, 20, 0, time.UTC),
		PDD:             time.Duration(200) * time.Millisecond,
		AnswerTime:      time.Date(2013, 11, 7, 8, 42, 26, 0, time.UTC),
		Usage:           time.Duration(10) * time.Second,
		Supplier:        cgrCdr[utils.SUPPLIER],
		DisconnectCause: cgrCdr[utils.DISCONNECT_CAUSE],
		ExtraFields:     map[string]string{"field_extr1": "val_extr1", "fieldextr2": "valextr2"},
		Cost:            dec.NewFloat(0.12),
		Rated:           true,
	}
	if CDR := cgrCdr.AsStoredCdr(""); !reflect.DeepEqual(expctRtCdr, CDR) {
		t.Errorf("Expecting %v, received: %v", expctRtCdr, CDR)
	}
}
