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
func (c *Contract) StartProcess(ctx TransactionContextInterface, processLocalId string, ownerOrg string, optionName string, startTime int64, startPosition string, preKey []string) (string, error) {
	if !ctx.CheckOrgValid(ownerOrg) {
		return "", perror.New("Org check failed")
	}
	process := &Process{}
	process.ProcessLocalId=processLocalId
	process.Class = ClassName
	process.State = InProcess
	process.OwnerOrg = ownerOrg
	process.StartTime = startTime
	process.StartPosition = startPosition
	process.OptionName = optionName
	preKeyField, err := c.createPreKeySlice(ctx, preKey)
	if err!=nil{
		return "",err
	}else {
		process.PreKey=preKeyField
	}
	key := CreateCompositeKey(ownerOrg, processLocalId)
	process.Key=key
	return key, ctx.GetProcessLedger().AddProcess(key, process)
}

//--------Complete a new process
func (c *Contract) CompleteProcess(ctx TransactionContextInterface, key string, completeTime int64, completePosition string) error {
	process, err := c.QueryProcess(ctx,key)
	if err!=nil {
		return err
	}
	if !ctx.CheckOrgValid(process.OwnerOrg) {
		return perror.New("Org check failed")
	}
	process.State=Completed
	process.CompleteTime=completeTime
	process.CompletePosition=completePosition
	return ctx.GetProcessLedger().UpdateProcess(key, process)
}

//--------Link an existing process to its previous ones, WILL OVERWRITE if preKey field already exists
func (c *Contract) LinkProcess(ctx TransactionContextInterface, key string, preKey []string) error {
	// current process
	process,err := c.QueryProcess(ctx, key)
	if err!=nil{
		return perror.Errorf("Current process not exist. %s", err)
	}
	if !ctx.CheckOrgValid(process.OwnerOrg) {
		return perror.New("Org check failed")
	}
	preKeyField, err:= c.createPreKeySlice(ctx, preKey)
	if err!=nil{
		return err
	}
	process.PreKey=preKeyField
	return ctx.GetProcessLedger().UpdateProcess(key, process)
}

//--------Link an existing process to its previous ones, WILL APPEND if preKey field already exists
func (c *Contract) AddLinkedProcess(ctx TransactionContextInterface, key string, preKey []string) error {
	// current process
	process,err := c.QueryProcess(ctx, key)
	if err!=nil{
		return perror.Errorf("Current process not exist. %s", err)
	}
	if !ctx.CheckOrgValid(process.OwnerOrg) {
		return perror.New("Org check failed")
	}
	preKeyField, err:= c.createPreKeySlice(ctx, preKey)
	if err!=nil{
		return err
	}
	if len(process.PreKey) == 0{
		process.PreKey=preKeyField
	}else{
		process.PreKey = append(process.PreKey, preKeyField...)
	}

	return ctx.GetProcessLedger().UpdateProcess(key, process)
}

func (c *Contract) QueryProcess(ctx TransactionContextInterface, key string) (*Process, error) {
	return ctx.GetProcessLedger().GetProcess(key)
}

//---------Get the previous processes within depth of 1
func (c *Contract) PrevProcess(ctx TransactionContextInterface, key string) ([]*Process, error) {
	process,err := c.QueryProcess(ctx, key)
	if err!=nil{
		return nil,perror.Errorf("Key not exist. %s", err)
	}
	preKey := process.PreKey
	preProcess:=make([]*Process, 0, 2)
	for _,k:=range preKey{
		p,err:=c.QueryProcess(ctx, k)
		if err!=nil{
			return nil,err
		}
		preProcess=append(preProcess,p)
	}
	return preProcess,nil
}

//---------Get the sourcing chain of on one branch within given depth (actually position 0 of preKey array
func (c *Contract) DigProcess(ctx TransactionContextInterface, key string, depth int) ([]*Process, error) {
	process,err := c.QueryProcess(ctx, key)
	if err!=nil{
		return nil,perror.Errorf("Key not exist. %s", err)
	}
	preProcess:=make([]*Process, 0, MaxSourcingDepth)
	for d:=0;d<=depth&&d<=MaxSourcingDepth;d++{
		preKey:=process.PreKey
		if len(preKey)==0{
			break
		}
		pk:=preKey[0]
		process,err = c.QueryProcess(ctx,pk)
		if err!=nil{
			return nil,perror.Errorf("Error digging source preKey %s. %s",pk,err)
		}
		preProcess=append(preProcess,process)
	}
	return preProcess,nil
}

func (c *Contract) createPreKeySlice(ctx TransactionContextInterface, preKey []string ) ([]string, error){
	if len(preKey)!=0 {
		preKeyField := make([]string,0,1)
		for _,k := range preKey{
			_, err := c.QueryProcess(ctx, k)
			if err != nil {
				return nil, perror.Errorf("PreKey %s not exist. %s", k, err)
			}
			preKeyField=append(preKeyField,k)
		}
		return preKeyField,nil
	}
	return nil, nil
}
