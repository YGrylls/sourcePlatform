package process

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	perror "github.com/pkg/errors"
)

type Contract struct {
	contractapi.Contract
}

func (c *Contract) Instantiate() {
	fmt.Println("Init")
}

//--------Create a new process
func (c *Contract) StartProcess(ctx TransactionContextInterface, processLocalId string, ownerOrg string, optionName string, startTime int64, startPosition string, preKey string) (string, error) {
	if !ctx.CheckOrgValid(ownerOrg) {
		return "", perror.New("Org check failed")
	}
	process := &Process{}
	process.Class = ClassName
	process.State = InProcess
	process.OwnerOrg = ownerOrg
	process.StartTime = startTime
	process.StartPosition = startPosition
	if preKey != "" {
		_, err := c.QueryProcess(ctx, preKey)
		if err != nil {
			return "", perror.Errorf("PreKey not exist. %s", err)
		}
	}
	process.PreKey = preKey
	key := CreateCompositeKey(ownerOrg, processLocalId)
	process.OptionName = optionName
	return key, ctx.GetProcessLedger().AddProcess(key, process)
}

//--------Complete a new process
func (c *Contract) CompleteProcess(ctx TransactionContextInterface, key string) error {
	return nil
}

//--------Complete the previous process and start a new one next to it
func (c *Contract) TransferProcess(ctx TransactionContextInterface) error {
	return nil
}

func (c *Contract) QueryProcess(ctx TransactionContextInterface, key string) (*Process, error) {
	return ctx.GetProcessLedger().GetProcess(key)
}

func (c *Contract) SourceProcess(ctx TransactionContextInterface, key string, depth int) ([]*Process, error) {
	return nil, nil
}
