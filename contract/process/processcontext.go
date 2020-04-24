package process

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type TransactionContextInterface interface {
	contractapi.TransactionContextInterface
	CheckOrgValid(org string) bool
	GetProcessLedger() ProcessLedger
}

type TransactionContext struct {
	contractapi.TransactionContext
	ledger ProcessLedger
}

//todo
func (t *TransactionContext) CheckOrgValid(org string) bool {
	id, _:=t.GetClientIdentity().GetMSPID()
	fmt.Println("MSPID: ",id)
	creator, _ := t.GetStub().GetCreator()
	fmt.Println("creator: ", string(creator))
	return true
}

func (t *TransactionContext) GetProcessLedger() ProcessLedger {
	if t.ledger == nil {
		t.ledger = newLedger(&t.TransactionContext)
	}

	return t.ledger
}
