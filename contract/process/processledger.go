package process

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strings"
	"unicode/utf8"
)

type ProcessLedger interface {
	AddProcess(key string, process *Process) error
	GetProcess(key string) (*Process, error)
	UpdateProcess(key string, process *Process) error
}

func newLedger(ctx contractapi.TransactionContextInterface) ProcessLedger {
	return &processLedger{ctx: ctx}
}

type processLedger struct {
	ctx contractapi.TransactionContextInterface
}

func (p *processLedger) AddProcess(key string, process *Process) error {
	bytes, err := process.Serialize()
	if err != nil {
		return err
	}
	return p.ctx.GetStub().PutState(key, bytes)
}

func (p *processLedger) GetProcess(key string) (*Process, error) {
	bytes, err := p.ctx.GetStub().GetState(key)
	if err != nil {
		return nil, err
	}
	process := &Process{}
	err = Deserialize(bytes, process)
	if err != nil {
		return nil, err
	}
	return process, nil
}

func (p *processLedger) UpdateProcess(key string, process *Process) error {
	return p.AddProcess(key,process)
}

const KeySplit = ":"
const maxUnicodeRuneValue   = utf8.MaxRune


func CreateCompositeKey(org string, localId string) string {
	return strings.Join([]string{org, localId}, KeySplit)
}

func SplitCompositeKey(key string) []string {
	return strings.Split(key, KeySplit)
}
