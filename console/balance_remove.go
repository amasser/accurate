package console

import (
	"github.com/accurateproject/accurate/api/v1"
	"github.com/accurateproject/accurate/utils"
)

func init() {
	c := &CmdRemoveBalance{
		name:      "balance_remove",
		rpcMethod: "ApiV1.RemoveBalances",
	}
	commands[c.Name()] = c
	c.CommandExecuter = &CommandExecuter{c}
}

// Commander implementation
type CmdRemoveBalance struct {
	name      string
	rpcMethod string
	rpcParams *v1.AttrSetBalance
	*CommandExecuter
}

func (self *CmdRemoveBalance) Name() string {
	return self.name
}

func (self *CmdRemoveBalance) RpcMethod() string {
	return self.rpcMethod
}

func (self *CmdRemoveBalance) RpcParams(reset bool) interface{} {
	if reset || self.rpcParams == nil {
		self.rpcParams = &v1.AttrSetBalance{TOR: utils.MONETARY, Overwrite: false}
	}
	return self.rpcParams
}

func (self *CmdRemoveBalance) PostprocessRpcParams() error {
	return nil
}

func (self *CmdRemoveBalance) RpcResult() interface{} {
	var s string
	return &s
}
