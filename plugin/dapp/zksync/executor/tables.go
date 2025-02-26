package executor

import (
	"fmt"
	"math/big"

	"github.com/33cn/chain33/common/address"

	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/common/db/table"
	"github.com/33cn/chain33/types"
	zt "github.com/33cn/plugin/plugin/dapp/zksync/types"
)

const (
	//KeyPrefixStateDB state db key必须前缀
	KeyPrefixStateDB = "mavl-zksync-"
	//KeyPrefixLocalDB local db的key必须前缀
	KeyPrefixLocalDB = "LODB-zksync"
)

var opt_account_tree = &table.Option{
	Prefix:  KeyPrefixLocalDB,
	Name:    "tree",
	Primary: "address",
	Index:   []string{"chain33_address", "eth_address"},
}

// NewAccountTreeTable ...
func NewAccountTreeTable(kvdb db.KV) *table.Table {
	rowmeta := NewAccountTreeRow()
	table, err := table.NewTable(rowmeta, kvdb, opt_account_tree)
	if err != nil {
		panic(err)
	}
	return table
}

// AccountTreeRow table meta 结构
type AccountTreeRow struct {
	*zt.Leaf
}

func NewAccountTreeRow() *AccountTreeRow {
	return &AccountTreeRow{Leaf: &zt.Leaf{}}
}

//CreateRow 新建数据行
func (r *AccountTreeRow) CreateRow() *table.Row {
	return &table.Row{Data: &zt.Leaf{}}
}

//SetPayload 设置数据
func (r *AccountTreeRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*zt.Leaf); ok {
		r.Leaf = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get 按照indexName 查询 indexValue
func (r *AccountTreeRow) Get(key string) ([]byte, error) {
	if key == "address" {
		return GetLocalChain33EthPrimaryKey(r.GetChain33Addr(), r.GetEthAddress()), nil
	} else if key == "chain33_address" {
		return []byte(fmt.Sprintf("%s", address.FormatAddrKey(r.GetChain33Addr()))), nil
	} else if key == "eth_address" {
		return []byte(fmt.Sprintf("%s", address.FormatAddrKey(r.GetEthAddress()))), nil
	}
	return nil, types.ErrNotFound
}

var opt_zksync_info = &table.Option{
	Prefix:  KeyPrefixLocalDB,
	Name:    "zksync",
	Primary: "account_token_id",
	Index:   []string{"accountID", "tokenID"},
}

// NewZksyncInfoTable ...
func NewZksyncInfoTable(kvdb db.KV) *table.Table {
	rowmeta := NewZksyncInfoRow()
	table, err := table.NewTable(rowmeta, kvdb, opt_zksync_info)
	if err != nil {
		panic(err)
	}
	return table
}

// AccountTreeRow table meta 结构
type ZksyncInfoRow struct {
	*zt.AccountTokenBalanceReceipt
}

func NewZksyncInfoRow() *ZksyncInfoRow {
	return &ZksyncInfoRow{AccountTokenBalanceReceipt: &zt.AccountTokenBalanceReceipt{}}
}

//CreateRow 新建数据行
func (r *ZksyncInfoRow) CreateRow() *table.Row {
	return &table.Row{Data: &zt.AccountTokenBalanceReceipt{}}
}

//SetPayload 设置数据
func (r *ZksyncInfoRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*zt.AccountTokenBalanceReceipt); ok {
		r.AccountTokenBalanceReceipt = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get 按照indexName 查询 indexValue
func (r *ZksyncInfoRow) Get(key string) ([]byte, error) {
	if key == "account_token_id" {
		return []byte(fmt.Sprintf("%016d.%016d", r.GetAccountId(), r.GetTokenId())), nil
	} else if key == "accountID" {
		return []byte(fmt.Sprintf("%016d", r.GetAccountId())), nil
	} else if key == "tokenID" {
		return []byte(fmt.Sprintf("%016d", r.GetTokenId())), nil
	}
	return nil, types.ErrNotFound
}

var opt_commit_proof = &table.Option{
	Prefix:  KeyPrefixLocalDB,
	Name:    "proof",
	Primary: "proofId",
	Index:   []string{"endHeight", "root", "commitHeight", "onChainId"},
}

// NewCommitProofTable ...
func NewCommitProofTable(kvdb db.KV) *table.Table {
	rowmeta := NewCommitProofRow()
	table, err := table.NewTable(rowmeta, kvdb, opt_commit_proof)
	if err != nil {
		panic(err)
	}
	return table
}

// CommitProofRow table meta 结构
type CommitProofRow struct {
	*zt.ZkCommitProof
}

func NewCommitProofRow() *CommitProofRow {
	return &CommitProofRow{ZkCommitProof: &zt.ZkCommitProof{}}
}

//CreateRow 新建数据行
func (r *CommitProofRow) CreateRow() *table.Row {
	return &table.Row{Data: &zt.ZkCommitProof{}}
}

//SetPayload 设置数据
func (r *CommitProofRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*zt.ZkCommitProof); ok {
		r.ZkCommitProof = txdata
		return nil
	}
	return types.ErrTypeAsset
}

func (r *CommitProofRow) isProofNeedOnChain() int {
	if len(r.GetOnChainPubDatas()) > 0 {
		return 1
	}
	return 0
}

//Get 按照indexName 查询 indexValue
func (r *CommitProofRow) Get(key string) ([]byte, error) {
	chainTitleId := new(big.Int).SetUint64(r.ChainTitleId).String()
	if key == "proofId" {
		return []byte(fmt.Sprintf("%016d-%s", r.GetProofId(), chainTitleId)), nil
	} else if key == "root" {
		return []byte(fmt.Sprintf("%s-%s", chainTitleId, r.GetNewTreeRoot())), nil
	} else if key == "endHeight" {
		return []byte(fmt.Sprintf("%s-%016d", chainTitleId, r.GetBlockEnd())), nil
	} else if key == "commitHeight" {
		return []byte(fmt.Sprintf("%s-%016d", chainTitleId, r.GetCommitBlockHeight())), nil
	} else if key == "onChainId" {
		return []byte(fmt.Sprintf("%s-%016d", chainTitleId, r.GetOnChainProofId())), nil
	}
	return nil, types.ErrNotFound
}

var opt_history_account_tree = &table.Option{
	Prefix:  KeyPrefixLocalDB,
	Name:    "historyTree",
	Primary: "proofId_accountId",
	Index:   []string{"proofId"},
}

// NewHistoryAccountTreeTable ...
func NewHistoryAccountTreeTable(kvdb db.KV) *table.Table {
	rowmeta := NewHistoryAccountTreeRow()
	table, err := table.NewTable(rowmeta, kvdb, opt_history_account_tree)
	if err != nil {
		panic(err)
	}
	return table
}

// HistoryAccountTreeRow table meta 结构
type HistoryAccountTreeRow struct {
	*zt.HistoryLeaf
}

func NewHistoryAccountTreeRow() *HistoryAccountTreeRow {
	return &HistoryAccountTreeRow{HistoryLeaf: &zt.HistoryLeaf{}}
}

//CreateRow 新建数据行
func (r *HistoryAccountTreeRow) CreateRow() *table.Row {
	return &table.Row{Data: &zt.HistoryLeaf{}}
}

//SetPayload 设置数据
func (r *HistoryAccountTreeRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*zt.HistoryLeaf); ok {
		r.HistoryLeaf = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get 按照indexName 查询 indexValue
func (r *HistoryAccountTreeRow) Get(key string) ([]byte, error) {
	if key == "proofId_accountId" {
		return []byte(fmt.Sprintf("%016d.%16d", r.GetProofId(), r.GetAccountId())), nil
	} else if key == "proofId" {
		return []byte(fmt.Sprintf("%016d", r.GetProofId())), nil
	}
	return nil, types.ErrNotFound
}
