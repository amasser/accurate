// +build integration

package general_tests

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/accurateproject/accurate/config"
	"github.com/accurateproject/accurate/engine"
	"github.com/accurateproject/accurate/utils"
	"github.com/accurateproject/rpcclient"
)

var cdrsMasterCfgPath, cdrsSlaveCfgPath string
var cdrsMasterCfg, cdrsSlaveCfg *config.Config
var cdrsMasterRpc *rpcclient.RpcClient

var testIntegration = flag.Bool("integration", false, "Perform the tests in integration mode, not by default.") // This flag will be passed here via "go test -local" args

func TestCdrsInitConfig(t *testing.T) {
	config.Reset()
	var err error
	cdrsMasterCfgPath = path.Join(*dataDir, "conf", "samples", "cdrsreplicationmaster")
	if err = config.LoadPath(cdrsMasterCfgPath); err != nil {
		t.Fatal("Got config error: ", err.Error())
	}
	cdrsMasterCfg = &config.Config{}
	if err := utils.Clone(config.Get(), cdrsMasterCfg); err != nil {
		t.Fatal("error cloning config: ", err)
	}

	config.Reset()
	cdrsSlaveCfgPath = path.Join(*dataDir, "conf", "samples", "cdrsreplicationslave")
	if err = config.LoadPath(cdrsSlaveCfgPath); err != nil {
		t.Fatal("Got config error: ", err.Error())
	}
	cdrsSlaveCfg = config.Get()
}

// InitDb so we can rely on count
func TestCdrsInitCdrDb(t *testing.T) {
	if err := engine.InitStorDb(cdrsMasterCfg); err != nil {
		t.Fatal(err)
	}
	if err := engine.InitStorDb(cdrsSlaveCfg); err != nil {
		t.Fatal(err)
	}
	/*
		if err := os.Mkdir(cdrsMasterCfg.HttpFailedDir, 0700); err != nil {
			t.Error(err)
		}
	*/

}

func TestCdrsStartMasterEngine(t *testing.T) {
	if _, err := engine.StopStartEngine(cdrsMasterCfgPath, *waitRater); err != nil {
		t.Fatal(err)
	}
}

func TestCdrsStartSlaveEngine(t *testing.T) {
	if _, err := engine.StartEngine(cdrsSlaveCfgPath, *waitRater); err != nil {
		t.Fatal(err)
	}
}

// Connect rpc client to rater
func TestCdrsHttpCdrReplication(t *testing.T) {
	cdrsMasterRpc, err = rpcclient.NewRpcClient("tcp", *cdrsMasterCfg.Listen.RpcJson, 1, 1, time.Duration(1*time.Second), time.Duration(2*time.Second), "json", nil, false)
	if err != nil {
		t.Fatal("Could not connect to rater: ", err.Error())
	}
	testCdr1 := &engine.CDR{UniqueID: utils.Sha1("httpjsonrpc1", time.Date(2013, 12, 7, 8, 42, 24, 0, time.UTC).String()),
		ToR: utils.VOICE, OriginID: "httpjsonrpc1", OriginHost: "192.168.1.1", Source: "UNKNOWN", RequestType: utils.META_PSEUDOPREPAID,
		Direction: "*out", Tenant: "cgrates.org", Category: "call", Account: "1001", Subject: "1001", Destination: "1002",
		SetupTime: time.Date(2013, 12, 7, 8, 42, 24, 0, time.UTC), AnswerTime: time.Date(2013, 12, 7, 8, 42, 26, 0, time.UTC),
		Usage: time.Duration(10) * time.Second, ExtraFields: map[string]string{"field_extr1": "val_extr1", "fieldextr2": "valextr2"},
		RunID: utils.DEFAULT_RUNID, Cost: 1.201, Rated: true}
	var reply string
	if err := cdrsMasterRpc.Call("CdrsV2.ProcessCdr", testCdr1, &reply); err != nil {
		t.Error("Unexpected error: ", err.Error())
	} else if reply != utils.OK {
		t.Error("Unexpected reply received: ", reply)
	}
	time.Sleep(time.Duration(*waitRater) * time.Millisecond)
	cdrsSlaveRpc, err := rpcclient.NewRpcClient("tcp", "127.0.0.1:12012", 1, 1, time.Duration(1*time.Second), time.Duration(2*time.Second), "json", nil)
	if err != nil {
		t.Fatal("Could not connect to rater: ", err.Error())
	}
	// ToDo: Fix cdr_http to be compatible with rest of processCdr methods
	var rcvedCdrs []*engine.ExternalCDR
	if err := cdrsSlaveRpc.Call("ApierV2.GetCdrs", utils.RPCCDRsFilter{UniqueIDs: []string{testCdr1.UniqueID}, RunIDs: []string{utils.META_DEFAULT}}, &rcvedCdrs); err != nil {
		t.Error("Unexpected error: ", err.Error())
	} else if len(rcvedCdrs) != 1 {
		t.Error("Unexpected number of CDRs returned: ", len(rcvedCdrs))
	} else {
		rcvSetupTime, _ := utils.ParseTimeDetectLayout(rcvedCdrs[0].SetupTime, "")
		rcvAnswerTime, _ := utils.ParseTimeDetectLayout(rcvedCdrs[0].AnswerTime, "")
		//rcvUsage, _ := utils.ParseDurationWithSecs(rcvedCdrs[0].Usage)
		if rcvedCdrs[0].UniqueID != testCdr1.UniqueID ||
			rcvedCdrs[0].ToR != testCdr1.ToR ||
			rcvedCdrs[0].OriginHost != testCdr1.OriginHost ||
			rcvedCdrs[0].RequestType != testCdr1.RequestType ||
			rcvedCdrs[0].Direction != testCdr1.Direction ||
			rcvedCdrs[0].Tenant != testCdr1.Tenant ||
			rcvedCdrs[0].Category != testCdr1.Category ||
			rcvedCdrs[0].Account != testCdr1.Account ||
			rcvedCdrs[0].Subject != testCdr1.Subject ||
			rcvedCdrs[0].Destination != testCdr1.Destination ||
			!rcvSetupTime.Equal(testCdr1.SetupTime) ||
			!rcvAnswerTime.Equal(testCdr1.AnswerTime) ||
			//rcvUsage != 10 ||
			rcvedCdrs[0].RunID != testCdr1.RunID {
			//rcvedCdrs[0].Cost != testCdr1.Cost ||
			//!reflect.DeepEqual(rcvedCdrs[0].ExtraFields, testCdr1.ExtraFields) {
			t.Errorf("Expected: %+v, received: %+v", testCdr1, rcvedCdrs[0])
		}
	}
}

// Connect rpc client to rater
func TestCdrsFileFailover(t *testing.T) {
	time.Sleep(time.Duration(1 * time.Second))
	failoverContent := []byte(`Account=1001&AnswerTime=2013-12-07T08%3A42%3A26Z&Category=call&Destination=1002&Direction=%2Aout&DisconnectCause=&OriginHost=192.168.1.1&OriginID=httpjsonrpc1&PDD=0&RequestType=%2Apseudoprepaid&SetupTime=2013-12-07T08%3A42%3A24Z&Source=UNKNOWN&Subject=1001&Supplier=&Tenant=cgrates.org&ToR=%2Avoice&Usage=10&field_extr1=val_extr1&fieldextr2=valextr2`)
	var rplCfg *config.CdrReplication
	for _, rplCfg = range cdrsMasterCfg.Cdrs.CdrReplication {
		if strings.HasSuffix(rplCfg.Address, "invalid") { // Find the config which shold generate the failoback
			break
		}
	}
	filesInDir, _ := ioutil.ReadDir(*cdrsMasterCfg.General.HttpFailedDir)
	if len(filesInDir) == 0 {
		t.Fatalf("No files in directory: %s", *cdrsMasterCfg.General.HttpFailedDir)
	}
	var fileName string
	for _, file := range filesInDir { // First file in directory is the one we need, harder to find it's name out of config
		fileName = file.Name()
		break
	}
	filePath := path.Join(*cdrsMasterCfg.General.HttpFailedDir, fileName)
	if readBytes, err := ioutil.ReadFile(filePath); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(failoverContent[11], readBytes[11]) { // Checking just the prefix should do since some content is dynamic
		t.Errorf("Expecting: %q, received: %q", string(failoverContent[11]), string(readBytes[11]))
	}
	if err := os.Remove(filePath); err != nil {
		t.Error("Failed removing file: ", filePath)
	}
}

/*
// Performance test, check `lsof -a -p 8427 | wc -l`

func TestCdrsHttpCdrReplication2(t *testing.T) {
	if !*testIntegration {
		return
	}
	cdrs := make([]*engine.CDR, 0)
	for i := 0; i < 10000; i++ {
		cdr := &engine.CDR{OriginID: fmt.Sprintf("httpjsonrpc_%d", i),
			ToR: utils.VOICE, OriginHost: "192.168.1.1", Source: "UNKNOWN", RequestType: utils.META_PSEUDOPREPAID,
			Direction: "*out", Tenant: "cgrates.org", Category: "call", Account: "1001", Subject: "1001", Destination: "1002",
			SetupTime: time.Date(2013, 12, 7, 8, 42, 24, 0, time.UTC), AnswerTime: time.Date(2013, 12, 7, 8, 42, 26, 0, time.UTC),
			Usage: time.Duration(10) * time.Second, ExtraFields: map[string]string{"field_extr1": "val_extr1", "fieldextr2": "valextr2"}}
		cdrs = append(cdrs, cdr)
	}
	var reply string
	for _, cdr := range cdrs {
		if err := cdrsMasterRpc.Call("CdrsV2.ProcessCdr", cdr, &reply); err != nil {
			t.Error("Unexpected error: ", err.Error())
		} else if reply != utils.OK {
			t.Error("Unexpected reply received: ", reply)
		}
	}
}
*/
