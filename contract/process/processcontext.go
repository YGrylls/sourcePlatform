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

const readonlyOrgMSP = "QueryMSP"
//todo
//to check if the client is from a query org or if it is operating a process belonging to others
func (t *TransactionContext) CheckOrgValid(org string) bool {
	id, _:=t.GetClientIdentity().GetMSPID()
	fmt.Println("MSPID: ",id)
	return id==org || id==readonlyOrgMSP
}

func (t *TransactionContext) GetProcessLedger() ProcessLedger {
	if t.ledger == nil {
		t.ledger = newLedger(&t.TransactionContext)
	}

	return t.ledger
}
