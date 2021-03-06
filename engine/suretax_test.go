package engine

import (
	"reflect"
	"testing"
	"time"

	"github.com/accurateproject/accurate/config"
	"github.com/accurateproject/accurate/dec"
	"github.com/accurateproject/accurate/utils"
)

func TestNewSureTaxRequest(t *testing.T) {
	CCID := utils.Sha1("dsafdsaf", time.Date(2013, 11, 7, 8, 42, 20, 0, time.UTC).String())
	cdr := &CDR{UniqueID: CCID, OrderID: 123, ToR: utils.VOICE,
		OriginID: "dsafdsaf", OriginHost: "192.168.1.1", Source: utils.UNIT_TEST, RequestType: utils.META_RATED, Direction: "*out",
		Tenant: "cgrates.org", Category: "call", Account: "1001", Subject: "1001", Destination: "1002", Supplier: "SUPPL1",
		SetupTime: time.Date(2013, 11, 7, 8, 42, 20, 0, time.UTC), AnswerTime: time.Date(2013, 11, 7, 8, 42, 26, 0, time.UTC), RunID: utils.DEFAULT_RUNID,
		Usage: time.Duration(12) * time.Second, PDD: time.Duration(7) * time.Second,
		ExtraFields: map[string]string{"field_extr1": "val_extr1", "fieldextr2": "valextr2"}, Cost: dec.NewFloat(1.01), Rated: true,
	}
	config.Reset()
	stCfg := config.Get().SureTax
	stCfg.ClientNumber = utils.StringPointer("000000000")
	stCfg.ValidationKey = utils.StringPointer("19491161-F004-4F44-BDB3-E976D6739A64")
	stCfg.Timezone = utils.StringPointer(time.UTC.String())
	eSTRequest := &STRequest{
		ClientNumber:   "000000000",
		ValidationKey:  "19491161-F004-4F44-BDB3-E976D6739A64",
		DataYear:       "2013",
		DataMonth:      "11",
		TotalRevenue:   1.01,
		ReturnFileCode: "0",
		ClientTracking: CCID,
		ResponseGroup:  "03",
		ResponseType:   "D4",
		ItemList: []*STRequestItem{
			&STRequestItem{
				CustomerNumber:       "1001",
				OrigNumber:           "1001",
				TermNumber:           "1002",
				BillToNumber:         "",
				TransDate:            "2013-11-07T08:42:26",
				Revenue:              1.01,
				Units:                1,
				UnitType:             "00",
				Seconds:              12,
				TaxIncludedCode:      "0",
				TaxSitusRule:         "04",
				TransTypeCode:        "010101",
				SalesTypeCode:        "R",
				RegulatoryCode:       "03",
				TaxExemptionCodeList: []string{},
			},
		},
	}
	jsnReq := utils.ToJSON(eSTRequest)
	eSureTaxRequest := &SureTaxRequest{Request: string(jsnReq)}
	if stReq, err := NewSureTaxRequest(cdr, stCfg); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(eSureTaxRequest, stReq) {
		t.Errorf("Expecting: %s, received: %s", string(eSureTaxRequest.Request), string(stReq.Request))
	}
}
