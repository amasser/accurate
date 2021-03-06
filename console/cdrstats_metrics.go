package console

import "github.com/accurateproject/accurate/utils"

func init() {
	c := &CmdCdrStatsMetrics{
		name:      "cdrstats_metrics",
		rpcMethod: "CDRStatsV1.GetMetrics",
	}
	commands[c.Name()] = c
	c.CommandExecuter = &CommandExecuter{c}
}

// Commander implementation
type CmdCdrStatsMetrics struct {
	name      string
	rpcMethod string
	rpcParams *utils.AttrStatsQueueID
	*CommandExecuter
}

func (self *CmdCdrStatsMetrics) Name() string {
	return self.name
}

func (self *CmdCdrStatsMetrics) RpcMethod() string {
	return self.rpcMethod
}

func (self *CmdCdrStatsMetrics) RpcParams(reset bool) interface{} {
	if reset || self.rpcParams == nil {
		self.rpcParams = &utils.AttrStatsQueueID{}
	}
	return self.rpcParams
}

func (self *CmdCdrStatsMetrics) PostprocessRpcParams() error {
	return nil
}

func (self *CmdCdrStatsMetrics) RpcResult() interface{} {
	return &map[string]float64{}
}
