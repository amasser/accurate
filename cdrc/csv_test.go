package cdrc

import (
	"reflect"
	"testing"
	"time"

	"github.com/accurateproject/accurate/config"
	"github.com/accurateproject/accurate/dec"
	"github.com/accurateproject/accurate/engine"
	"github.com/accurateproject/accurate/utils"
)

func TestCsvRecordForkCdr(t *testing.T) {
	config.Reset()
	cgrConfig := config.Get()
	cdrcConfig := *cgrConfig.CdrcProfiles()["/var/spool/accurate/cdrc/in"][0]
	cdrcConfig.CdrSourceID = utils.StringPointer("TEST_CDRC")
	cdrcConfig.ContentFields = append(cdrcConfig.ContentFields, &config.CdrField{Tag: "SupplierTest", Type: utils.META_COMPOSED, FieldID: utils.SUPPLIER, Value: []*utils.RSRField{&utils.RSRField{Id: "14"}}})
	cdrcConfig.ContentFields = append(cdrcConfig.ContentFields, &config.CdrField{Tag: "DisconnectCauseTest", Type: utils.META_COMPOSED, FieldID: utils.DISCONNECT_CAUSE,
		Value: []*utils.RSRField{&utils.RSRField{Id: "16"}}})
	//
	csvProcessor := &CsvRecordsProcessor{dfltCdrcCfg: &cdrcConfig, cdrcCfgs: []*config.Cdrc{&cdrcConfig}}
	cdrRow := []string{"firstField", "secondField"}
	_, err := csvProcessor.recordToStoredCdr(cdrRow, &cdrcConfig)
	if err == nil {
		t.Error("Failed to corectly detect missing fields from record")
	}
	cdrRow = []string{"ignored", "ignored", utils.VOICE, "acc1", utils.META_PREPAID, "*out", "cgrates.org", "call", "1001", "1001", "+4986517174963",
		"2013-02-03 19:50:00", "2013-02-03 19:54:00", "62", "supplier1", "172.16.1.1", "NORMAL_DISCONNECT"}
	rtCdr, err := csvProcessor.recordToStoredCdr(cdrRow, &cdrcConfig)
	if err != nil {
		t.Error("Failed to parse CDR in rated cdr", err)
	}
	expectedCdr := &engine.CDR{
		UniqueID:        utils.Sha1(cdrRow[3], time.Date(2013, 2, 3, 19, 50, 0, 0, time.UTC).String()),
		ToR:             cdrRow[2],
		OriginID:        cdrRow[3],
		OriginHost:      "0.0.0.0", // Got it over internal interface
		Source:          "TEST_CDRC",
		RequestType:     cdrRow[4],
		Direction:       cdrRow[5],
		Tenant:          cdrRow[6],
		Category:        cdrRow[7],
		Account:         cdrRow[8],
		Subject:         cdrRow[9],
		Destination:     cdrRow[10],
		SetupTime:       time.Date(2013, 2, 3, 19, 50, 0, 0, time.UTC),
		AnswerTime:      time.Date(2013, 2, 3, 19, 54, 0, 0, time.UTC),
		Usage:           time.Duration(62) * time.Second,
		Supplier:        "supplier1",
		DisconnectCause: "NORMAL_DISCONNECT",
		ExtraFields:     map[string]string{},
		Cost:            dec.NewVal(-1, 0),
	}
	if !reflect.DeepEqual(expectedCdr, rtCdr) {
		t.Errorf("Expected: \n%v, \nreceived: \n%v", expectedCdr, rtCdr)
	}
}

func TestCsvDataMultiplyFactor(t *testing.T) {
	config.Reset()
	cgrConfig := config.Get()
	cdrcConfig := cgrConfig.CdrcProfiles()["/var/spool/accurate/cdrc/in"][0]
	cdrcConfig.CdrSourceID = utils.StringPointer("TEST_CDRC")
	cdrcConfig.ContentFields = []*config.CdrField{&config.CdrField{Tag: "TORField", Type: utils.META_COMPOSED, FieldID: utils.TOR, Value: []*utils.RSRField{&utils.RSRField{Id: "0"}}},
		&config.CdrField{Tag: "UsageField", Type: utils.META_COMPOSED, FieldID: utils.USAGE, Value: []*utils.RSRField{&utils.RSRField{Id: "1"}}}}
	csvProcessor := &CsvRecordsProcessor{dfltCdrcCfg: cdrcConfig, cdrcCfgs: []*config.Cdrc{cdrcConfig}}
	*csvProcessor.cdrcCfgs[0].DataUsageMultiplyFactor = 0
	cdrRow := []string{"*data", "1"}
	rtCdr, err := csvProcessor.recordToStoredCdr(cdrRow, cdrcConfig)
	if err != nil {
		t.Error("Failed to parse CDR in rated cdr", err)
	}
	var sTime time.Time
	expectedCdr := &engine.CDR{
		UniqueID:    utils.Sha1("", sTime.String()),
		ToR:         cdrRow[0],
		OriginHost:  "0.0.0.0",
		Source:      "TEST_CDRC",
		Usage:       time.Duration(1) * time.Second,
		ExtraFields: map[string]string{},
		Cost:        dec.NewVal(-1, 0),
	}
	if !reflect.DeepEqual(expectedCdr, rtCdr) {
		t.Errorf("Expected: \n%v, \nreceived: \n%v", expectedCdr, rtCdr)
	}
	*csvProcessor.cdrcCfgs[0].DataUsageMultiplyFactor = 1024
	expectedCdr = &engine.CDR{
		UniqueID:    utils.Sha1("", sTime.String()),
		ToR:         cdrRow[0],
		OriginHost:  "0.0.0.0",
		Source:      "TEST_CDRC",
		Usage:       time.Duration(1024) * time.Second,
		ExtraFields: map[string]string{},
		Cost:        dec.NewVal(-1, 0),
	}
	if rtCdr, _ := csvProcessor.recordToStoredCdr(cdrRow, cdrcConfig); !reflect.DeepEqual(expectedCdr, rtCdr) {
		t.Errorf("Expected: \n%v, \nreceived: \n%v", expectedCdr, rtCdr)
	}
	cdrRow = []string{"*voice", "1"}
	expectedCdr = &engine.CDR{
		UniqueID:    utils.Sha1("", sTime.String()),
		ToR:         cdrRow[0],
		OriginHost:  "0.0.0.0",
		Source:      "TEST_CDRC",
		Usage:       time.Duration(1) * time.Second,
		ExtraFields: map[string]string{},
		Cost:        dec.NewVal(-1, 0),
	}
	if rtCdr, _ := csvProcessor.recordToStoredCdr(cdrRow, cdrcConfig); !reflect.DeepEqual(expectedCdr, rtCdr) {
		t.Errorf("Expected: \n%v, \nreceived: \n%v", expectedCdr, rtCdr)
	}
}

func TestCsvPairToRecord(t *testing.T) {
	eRecord := []string{"INVITE", "2daec40c", "548625ac", "dd0c4c617a9919d29a6175cdff223a9e@0:0:0:0:0:0:0:0", "200", "OK", "1436454408", "*prepaid", "1001", "1002", "", "3401:2069362475", "2"}
	invPr := &UnpairedRecord{Method: "INVITE", Timestamp: time.Date(2015, 7, 9, 15, 6, 48, 0, time.UTC),
		Values: []string{"INVITE", "2daec40c", "548625ac", "dd0c4c617a9919d29a6175cdff223a9e@0:0:0:0:0:0:0:0", "200", "OK", "1436454408", "*prepaid", "1001", "1002", "", "3401:2069362475"}}
	byePr := &UnpairedRecord{Method: "BYE", Timestamp: time.Date(2015, 7, 9, 15, 6, 50, 0, time.UTC),
		Values: []string{"BYE", "2daec40c", "548625ac", "dd0c4c617a9919d29a6175cdff223a9e@0:0:0:0:0:0:0:0", "200", "OK", "1436454410", "", "", "", "", "3401:2069362475"}}
	if rec, err := pairToRecord(invPr, byePr); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(eRecord, rec) {
		t.Errorf("Expected: %+v, received: %+v", eRecord, rec)
	}
	if rec, err := pairToRecord(byePr, invPr); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(eRecord, rec) {
		t.Errorf("Expected: %+v, received: %+v", eRecord, rec)
	}
	if _, err := pairToRecord(byePr, byePr); err == nil || err.Error() != "MISSING_INVITE" {
		t.Error(err)
	}
	if _, err := pairToRecord(invPr, invPr); err == nil || err.Error() != "MISSING_BYE" {
		t.Error(err)
	}
	byePr.Values = []string{"BYE", "2daec40c", "548625ac", "dd0c4c617a9919d29a6175cdff223a9e@0:0:0:0:0:0:0:0", "200", "OK", "1436454410", "", "", "", "3401:2069362475"} // Took one value out
	if _, err := pairToRecord(invPr, byePr); err == nil || err.Error() != "INCONSISTENT_VALUES_LENGTH" {
		t.Error(err)
	}
}
