// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/common"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	pty "github.com/33cn/plugin/plugin/dapp/issuance/types"
	tokenE "github.com/33cn/plugin/plugin/dapp/token/executor"
)

// List control
const (
	ListDESC    = int32(0)   // list降序
	ListASC     = int32(1)   // list升序
	DefultCount = int32(20)  // 默认一次取多少条记录
	MaxCount    = int32(100) // 最多取100条
)

const (
	Coin                      = types.Coin      // 1e8
	PriceWarningRate          = 1.3             // 价格提前预警率
	ExpireWarningTime         = 3600 * 24 * 10  // 提前10天超时预警
)

const (
	priceFeedKey = "issuance-price-feed"
	guarantorKey = "issuance-guarantor"
	manageKey    = "issuance-manage"
)

func getManageKey(key string, db dbm.KV) ([]byte, error) {
	manageKey := types.ManageKey(key)
	value, err := db.Get([]byte(manageKey))
	if err != nil {
		return nil, err
	}
	return value, nil
}

func getGuarantorAddr(db dbm.KV) (string, error) {
	value, err := getManageKey(guarantorKey, db)
	if err != nil {
		clog.Error("IssuancePriceFeed", "getGuarantorAddr", err)
		return "", err
	}
	if value == nil {
		clog.Error("IssuancePriceFeed guarantorKey found nil value")
		return "", err
	}

	var item types.ConfigItem
	err = types.Decode(value, &item)
	if err != nil {
		clog.Error("IssuancePriceFeed", "getGuarantorAddr", err)
		return "", err
	}

	return item.GetAddr(), nil
}

func isRightAddr(key string, addr string, db dbm.KV) bool {
	value, err := getManageKey(key, db)
	if err != nil {
		clog.Error("isRightAddr", "Key", key)
		return false
	}
	if value == nil {
		clog.Error("isRightAddr", "key", key, "error", "Found key nil value")
		return false
	}

	var item types.ConfigItem
	err = types.Decode(value, &item)
	if err != nil {
		clog.Error("isRightAddr", "Decode", value)
		return false
	}

	for _, op := range item.GetArr().Value {
		if op == addr {
			return true
		}
	}
	return false

}

// IssuanceDB def
type IssuanceDB struct {
	pty.Issuance
}

// GetKVSet for IssuanceDB
func (issu *IssuanceDB) GetKVSet() (kvset []*types.KeyValue) {
	value := types.Encode(&issu.Issuance)
	kvset = append(kvset, &types.KeyValue{Key: Key(issu.IssuanceId), Value: value})
	return kvset
}

// Save for IssuanceDB
func (issu *IssuanceDB) Save(db dbm.KV) {
	set := issu.GetKVSet()
	for i := 0; i < len(set); i++ {
		db.Set(set[i].GetKey(), set[i].Value)
	}
}

// Key for Issuance
func Key(id string) (key []byte) {
	key = append(key, []byte("mavl-"+pty.IssuanceX+"-")...)
	key = append(key, []byte(id)...)
	return key
}

// Key for IssuanceAddrConfig
func AddrKey() (key []byte) {
	key = append(key, []byte("mavl-"+pty.IssuanceX+"addr")...)
	return key
}

// Action struct
type Action struct {
	coinsAccount *account.DB  // bty账户
	tokenAccount *account.DB  // ccny账户
	db           dbm.KV
	localDB      dbm.Lister
	txhash       []byte
	fromaddr     string
	blocktime    int64
	height       int64
	execaddr     string
	difficulty   uint64
	index        int
	Issuance   *Issuance
}

// NewIssuanceAction generate New Action
func NewIssuanceAction(c *Issuance, tx *types.Transaction, index int) *Action {
	hash := tx.Hash()
	fromaddr := tx.From()
	tokenDb, err := account.NewAccountDB(tokenE.GetName(), pty.CCNYTokenName, c.GetStateDB())
	if err != nil {
		clog.Error("NewIssuanceAction", "Get Account DB error", "err", err)
		return nil
	}

	return &Action{
		coinsAccount: c.GetCoinsAccount(), tokenAccount:tokenDb, db: c.GetStateDB(), localDB:c.GetLocalDB(),
		txhash: hash, fromaddr: fromaddr, blocktime: c.GetBlockTime(), height: c.GetHeight(),
		execaddr: dapp.ExecAddress(string(tx.Execer)), difficulty: c.GetDifficulty(), index: index, Issuance: c}
}

// GetCollCommonRecipt generate logs for Issuance common action
func (action *Action) GetCollCommonRecipt(issuance *pty.Issuance, preStatus int32) *pty.ReceiptIssuance {
	c := &pty.ReceiptIssuance{}
	c.IssuanceId = issuance.IssuanceId
	c.PreStatus = preStatus
	c.Status = issuance.Status
	c.Index = action.GetIndex()
	return c
}

// GetCreateReceiptLog generate logs for Issuance create action
func (action *Action) GetCreateReceiptLog(issuance *pty.Issuance, preStatus int32) *types.ReceiptLog {
	log := &types.ReceiptLog{}
	log.Ty = pty.TyLogIssuanceCreate

	c := action.GetCollCommonRecipt(issuance, preStatus)

	log.Log = types.Encode(c)

	return log
}

// GetDebtReceiptLog generate logs for Issuance debt action
func (action *Action) GetDebtReceiptLog(issuance *pty.Issuance, preStatus int32) *types.ReceiptLog {
	log := &types.ReceiptLog{}
	log.Ty = pty.TyLogIssuanceDebt

	c := action.GetCollCommonRecipt(issuance, preStatus)
	c.AccountAddr = action.fromaddr

	log.Log = types.Encode(c)

	return log
}

// GetRepayReceiptLog generate logs for Issuance Repay action
func (action *Action) GetRepayReceiptLog(issuance *pty.Issuance, preStatus int32) *types.ReceiptLog {
	log := &types.ReceiptLog{}
	log.Ty = pty.TyLogIssuanceRepay

	c := action.GetCollCommonRecipt(issuance, preStatus)

	log.Log = types.Encode(c)

	return log
}

// GetFeedReceiptLog generate logs for Issuance price feed action
func (action *Action) GetFeedReceiptLog(issuance *pty.Issuance, debtRecord *pty.DebtRecord) *types.ReceiptLog {
	log := &types.ReceiptLog{}
	log.Ty = pty.TyLogIssuanceFeed

	c := action.GetCollCommonRecipt(issuance, debtRecord.PreStatus)
	c.AccountAddr = debtRecord.AccountAddr
	c.RecordStatus = debtRecord.Status

	log.Log = types.Encode(c)

	return log
}

// GetCloseReceiptLog generate logs for Issuance close action
func (action *Action) GetCloseReceiptLog(Issuance *pty.Issuance, preStatus int32) *types.ReceiptLog {
	log := &types.ReceiptLog{}
	log.Ty = pty.TyLogIssuanceClose

	c := action.GetCollCommonRecipt(Issuance, preStatus)

	log.Log = types.Encode(c)

	return log
}

// GetIndex returns index in block
func (action *Action) GetIndex() int64 {
	return action.height*types.MaxTxsPerBlock + int64(action.index)
}

func getLatestLiquidationPrice(issu *pty.Issuance) float32 {
	var latest float32
	for _, collRecord := range issu.DebtRecords {
		if collRecord.LiquidationPrice > latest {
			latest = collRecord.LiquidationPrice
		}
	}

	return latest
}

func getLatestExpireTime(issu *pty.Issuance) int64 {
	var latest int64 = 0x7fffffffffffffff

	for _, collRecord := range issu.DebtRecords {
		if collRecord.ExpireTime < latest {
			latest = collRecord.ExpireTime
		}
	}

	return latest
}

// IssuanceConfig 设置全局借贷参数（管理员权限）
func (action *Action) IssuanceManage(manage *pty.IssuanceManage) (*types.Receipt, error) {
	var kv []*types.KeyValue
	var receipt *types.Receipt

	// 是否配置管理用户
	if !isRightAddr(manageKey, action.fromaddr, action.db) {
		clog.Error("IssuanceManage", "addr", action.fromaddr, "error", "Address has no permission to config")
		return nil, pty.ErrPermissionDeny
	}

	// 添加大户地址
	data, err := action.db.Get(AddrKey())
	if err != nil {
		if err != types.ErrNotFound {
			clog.Error("IssuanceManage", "error", err)
			return nil, err
		}
		value := types.Encode(manage)
		action.db.Set(AddrKey(), value)
		kv = append(kv, &types.KeyValue{Key:AddrKey(), Value: value})
	} else {
		var addrStore pty.IssuanceManage
		err = types.Decode(data, &addrStore)
		if err != nil {
			clog.Debug("IssuanceManage", "decode", err)
			return nil, err
		}
		addrStore.SuperAddrs = append(addrStore.SuperAddrs, manage.SuperAddrs...)
		value := types.Encode(&addrStore)
		action.db.Set(AddrKey(), value)
		kv = append(kv, &types.KeyValue{Key:AddrKey(), Value: value})
	}

	receipt = &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: nil}
	return receipt, nil
}

func (action *Action) getSuperAddr() []string {
	data, err := action.db.Get(AddrKey())
	if err != nil {
		clog.Error("getSuperAddr", "error", err)
		return nil
	}

	var addrStore pty.IssuanceManage
	err = types.Decode(data, &addrStore)
	if err != nil {
		clog.Debug("getSuperAddr", "decode", err)
		return nil
	}

	return addrStore.SuperAddrs
}

func isSuperAddr(super []string, addr string) bool {
	if super == nil || len(super) == 0 {
		return false
	}

	for _, superAddr := range super {
		if superAddr == addr {
			return true
		}
	}

	return false
}

// IssuanceCreate 创建借贷，持有一定数量ccny的用户可创建借贷，提供给其他用户借贷
func (action *Action) IssuanceCreate(create *pty.IssuanceCreate) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var receipt *types.Receipt

	// 是否配置管理用户
	if !isRightAddr(manageKey, action.fromaddr, action.db) {
		clog.Error("IssuanceCreate", "addr", action.fromaddr, "error", "Address has no permission to create")
		return nil, pty.ErrPermissionDeny
	}

	// 检查ccny余额
	if !action.CheckExecTokenAccount(action.fromaddr, create.TotalBalance, false) {
		return nil, types.ErrInsufficientBalance
	}

	// 查找ID是否重复
	issuanceID := common.ToHex(action.txhash)
	_, err := queryIssuanceByID(action.db, issuanceID)
	if err != types.ErrNotFound {
		clog.Error("IssuanceCreate", "IssuanceCreate repeated", issuanceID)
		return nil, pty.ErrIssuanceRepeatHash
	}

	// 冻结ccny
	receipt, err = action.tokenAccount.ExecFrozen(action.fromaddr, action.execaddr, create.TotalBalance)
	if err != nil {
		clog.Error("IssuanceCreate.Frozen", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", create.TotalBalance)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	// 构造coll结构
	issu := &IssuanceDB{}
	issu.IssuanceId = issuanceID
	issu.LiquidationRatio = create.LiquidationRatio
	issu.TotalBalance = create.TotalBalance
	issu.Balance = create.TotalBalance
	issu.DebtCeiling = create.DebtCeiling
	issu.Period = create.Period
	issu.CreateTime = action.blocktime
	issu.IssuerAddr = action.fromaddr
	issu.Status = pty.IssuanceActionCreate

	clog.Debug("IssuanceCreate created", "IssuanceID", issuanceID, "TotalBalance", issu.TotalBalance)

	// 保存
	issu.Save(action.db)
	kv = append(kv, issu.GetKVSet()...)

	receiptLog := action.GetCreateReceiptLog(&issu.Issuance, 0)
	logs = append(logs, receiptLog)

	receipt = &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

// 根据最近抵押物价格计算需要冻结的BTY数量
func getBtyNumToFrozen(value int64, price float32, ratio float32) (int64,error) {
	if price == 0 {
		clog.Error("Bty price should greate to 0")
		return 0, pty.ErrPriceInvalid
	}

    btyValue := float32(value)/ratio
    btyNum := int64(btyValue/price) + 1

    return btyNum, nil
}

// 获取最近抵押物价格
func (action *Action)getLatestPrice(db dbm.KV) (float32, error) {
	data, err := db.Get(calcIssuanceLatestPriceKey())
	if err != nil {
		clog.Debug("getLatestPrice", "get", err)
		return -1, err
	}
	var price pty.AssetPriceRecord
	//decode
	err = types.Decode(data, &price)
	if err != nil {
		clog.Debug("getLatestPrice", "decode", err)
		return -1, err
	}

	return price.BtyPrice, nil
}

// CheckExecAccountBalance 检查账户抵押物余额
func (action *Action) CheckExecAccountBalance(fromAddr string, ToFrozen, ToActive int64) bool {
	acc := action.coinsAccount.LoadExecAccount(fromAddr, action.execaddr)
	if acc.GetBalance() >= ToFrozen && acc.GetFrozen() >= ToActive {
		return true
	}
	return false
}

// CheckExecAccount 检查账户token余额
func (action *Action) CheckExecTokenAccount(addr string, amount int64, isFrozen bool) bool {
	acc := action.tokenAccount.LoadExecAccount(addr, action.execaddr)
	if isFrozen {
		if acc.GetFrozen() >= amount {
			return true
		}
	} else {
		if acc.GetBalance() >= amount {
			return true
		}
	}
	return false
}

// IssuanceDebt 大户质押bty借出ccny
func (action *Action) IssuanceDebt(debt *pty.IssuanceDebt) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	// 查找对应的借贷ID
	issuance, err := queryIssuanceByID(action.db, debt.IssuanceId)
	if err != nil {
		clog.Error("IssuanceDebt", "IssuanceId", debt.IssuanceId, "err", err)
		return nil, err
	}

	// 状态检查
	if issuance.Status == pty.IssuanceStatusClose {
		clog.Error("IssuanceDebt", "CollID", issuance.IssuanceId, "addr", action.fromaddr, "execaddr", action.execaddr, "status", issuance.Status, "err", pty.ErrIssuanceStatus)
		return nil, pty.ErrIssuanceStatus
	}

	// 一个地址在一期借贷中只允许借出一次
	for _, record := range issuance.DebtRecords {
		if record.AccountAddr == action.fromaddr {
			clog.Error("IssuanceDebt","IssuanceId", debt.IssuanceId, action.fromaddr, "execaddr", action.execaddr, "err", pty.ErrIssuanceAccountExist)
			return nil, err
		}
	}
	issu := &IssuanceDB{*issuance}
	preStatus := issu.Status

	// 借贷金额检查
	if debt.GetValue() <= 0 {
		clog.Error("IssuanceDebt", "CollID", issu.IssuanceId, "addr", action.fromaddr, "execaddr", action.execaddr, "debt value", debt.GetValue(), "err", types.ErrInvalidParam)
		return  nil, types.ErrInvalidParam
	}

	// 借贷金额不超过个人限额
	if debt.GetValue() > issu.DebtCeiling {
		clog.Error("IssuanceDebt", "CollID", issu.IssuanceId, "addr", action.fromaddr, "execaddr", action.execaddr, "debt value", debt.GetValue(), "err", pty.ErrIssuanceExceedDebtCeiling)
		return nil, pty.ErrIssuanceExceedDebtCeiling
	}

	// 借贷金额不超过当前可借贷金额
	if debt.GetValue() > issu.Balance {
		clog.Error("IssuanceDebt", "CollID", issu.IssuanceId, "addr", action.fromaddr, "execaddr", action.execaddr, "debt value", debt.GetValue(), "err", pty.ErrIssuanceLowBalance)
		return nil, pty.ErrIssuanceLowBalance
	}
	clog.Debug("IssuanceDebt", "value", debt.GetValue())

	// 获取抵押物价格
	lastPrice, err := action.getLatestPrice(action.db)
	if err != nil {
		clog.Error("IssuanceDebt", "CollID", issu.IssuanceId, "addr", action.fromaddr, "execaddr", action.execaddr, "err", err)
		return nil, err
	}

	// 根据价格和需要借贷的金额，计算需要质押的抵押物数量
	btyFrozen, err := getBtyNumToFrozen(debt.Value, lastPrice, issu.LiquidationRatio)
	if err != nil {
		clog.Error("IssuanceDebt", "CollID", issu.IssuanceId, "addr", action.fromaddr, "execaddr", action.execaddr, "err", err)
		return nil, err
	}

	// 检查抵押物账户余额
	if !action.CheckExecAccountBalance(action.fromaddr, btyFrozen, 0) {
		clog.Error("IssuanceDebt", "CollID", issu.IssuanceId, "addr", action.fromaddr, "execaddr", action.execaddr, "err", types.ErrNoBalance)
		return nil, types.ErrNoBalance
	}

	// 抵押物转账
	receipt, err := action.coinsAccount.ExecTransfer(action.fromaddr, issu.IssuerAddr, action.execaddr, btyFrozen*Coin)
	if err != nil {
		clog.Error("IssuanceDebt.ExecTransfer", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", btyFrozen)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	// 抵押物冻结
	receipt, err = action.coinsAccount.ExecFrozen(issu.IssuerAddr, action.execaddr, btyFrozen)
	if err != nil {
		clog.Error("IssuanceDebt.Frozen", "addr", issu.IssuerAddr, "execaddr", action.execaddr, "amount", btyFrozen)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	// 借出ccny
	receipt, err = action.tokenAccount.ExecTransfer(issu.IssuerAddr, action.fromaddr, action.execaddr, debt.Value)
	if err != nil {
		clog.Error("IssuanceDebt.ExecTokenTransfer", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", debt.Value)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	// 构造借出记录
	debtRecord := &pty.DebtRecord{}
	debtRecord.AccountAddr = action.fromaddr
	debtRecord.CollateralValue = btyFrozen
	debtRecord.StartTime = action.blocktime
	debtRecord.CollateralPrice = lastPrice
	debtRecord.DebtValue = debt.Value
	debtRecord.LiquidationPrice = issu.LiquidationRatio * lastPrice * pty.IssuancePreLiquidationRatio
	debtRecord.Status = pty.IssuanceUserStatusCreate
	debtRecord.ExpireTime = action.blocktime + issu.Period

	// 记录当前借贷的最高自动清算价格
	if issu.LatestLiquidationPrice < debtRecord.LiquidationPrice {
		issu.LatestLiquidationPrice = debtRecord.LiquidationPrice
	}

	// 保存
	issu.DebtRecords = append(issu.DebtRecords, debtRecord)
	issu.CollateralValue += btyFrozen
	issu.DebtValue += debt.Value
	issu.Balance -= debt.Value
	issu.Save(action.db)
	kv = append(kv, issu.GetKVSet()...)

	receiptLog := action.GetDebtReceiptLog(&issu.Issuance, preStatus)
	logs = append(logs, receiptLog)

	receipt = &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

// IssuanceRepay 用户主动清算
func (action *Action) IssuanceRepay(repay *pty.IssuanceRepay) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var receipt *types.Receipt

	// 找到相应的借贷
	issuance, err := queryIssuanceByID(action.db, repay.IssuanceId)
	if err != nil {
		clog.Error("IssuanceRepay", "CollID", repay.IssuanceId, "err", err)
		return nil, err
	}

	issu := &IssuanceDB{*issuance}
	preStatus := issu.Status

	// 状态检查
	if issu.Status != pty.IssuanceStatusCreated {
		clog.Error("IssuanceRepay", "CollID", repay.IssuanceId, "addr", action.fromaddr, "execaddr", action.execaddr, "err", "status error", "Status", issu.Status)
		return nil, pty.ErrIssuanceStatus
	}

	// 查找借出记录
	var debtRecord *pty.DebtRecord
	var index int
	for i, record := range issu.DebtRecords {
		if record.AccountAddr == action.fromaddr {
			debtRecord = record
			index = i
			break
		}
	}

	if debtRecord == nil {
		clog.Error("IssuanceRepay", "CollID", repay.IssuanceId, "addr", action.fromaddr, "execaddr", action.execaddr, "err", "Can not find debt record")
		return nil, pty.ErrRecordNotExist
	}

	// 检查
	if !action.CheckExecTokenAccount(action.fromaddr, debtRecord.DebtValue, false) {
		clog.Error("IssuanceRepay", "CollID", issu.IssuanceId, "addr", action.fromaddr, "execaddr", action.execaddr, "err", types.ErrInsufficientBalance)
		return nil, types.ErrNoBalance
	}

	// ccny转移
	receipt, err = action.tokenAccount.ExecTransfer(action.fromaddr, issu.IssuerAddr, action.execaddr, debtRecord.DebtValue)
	if err != nil {
		clog.Error("IssuanceRepay.ExecTokenTransfer", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", debtRecord.DebtValue)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	// 抵押物归还
	receipt, err = action.coinsAccount.ExecTransferFrozen(issu.IssuerAddr, action.execaddr, action.execaddr, debtRecord.CollateralValue)
	if err != nil {
		clog.Error("IssuanceRepay.ExecTransferFrozen", "addr", issu.IssuerAddr, "execaddr", action.execaddr, "amount", debtRecord.CollateralValue)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	// 借贷记录关闭
	debtRecord.PreStatus = debtRecord.Status
	debtRecord.Status = pty.IssuanceUserStatusClose

	// 保存
	issu.Balance += debtRecord.DebtValue
	issu.CollateralValue -= debtRecord.CollateralValue
	issu.DebtValue -= debtRecord.DebtValue
	issu.DebtRecords = append(issu.DebtRecords[:index], issu.DebtRecords[index+1:]...)
	issu.InvalidRecords = append(issu.InvalidRecords, debtRecord)
	issu.LatestLiquidationPrice = getLatestLiquidationPrice(&issu.Issuance)
	issu.LatestExpireTime = getLatestExpireTime(&issu.Issuance)
	issu.Save(action.db)
	kv = append(kv, issu.GetKVSet()...)

	receiptLog := action.GetRepayReceiptLog(&issu.Issuance, preStatus)
	logs = append(logs, receiptLog)

	receipt = &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

// 系统清算
func (action *Action) systemLiquidation(issu *pty.Issuance, price float32) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	collDB := &IssuanceDB{*issu}
	for index, debtRecord := range issu.DebtRecords {
		if debtRecord.LiquidationPrice * PriceWarningRate < price {
			if debtRecord.Status == pty.IssuanceUserStatusSystemLiquidate {
				debtRecord.Status = debtRecord.PreStatus
				debtRecord.PreStatus = pty.IssuanceUserStatusSystemLiquidate
			}
			continue
		}

		if debtRecord.LiquidationPrice >= price {
			getGuarantorAddr, err := getGuarantorAddr(action.db)
			if err != nil {
				if err != nil {
					clog.Error("systemLiquidation", "getGuarantorAddr", err)
					continue
				}
			}

			// 抵押物转移
			receipt, err := action.coinsAccount.ExecTransferFrozen(action.fromaddr, getGuarantorAddr, action.execaddr, debtRecord.CollateralValue)
			if err != nil {
				clog.Error("systemLiquidation", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", debtRecord.CollateralValue, "err", err)
				continue
			}
			logs = append(logs, receipt.Logs...)
			kv = append(kv, receipt.KV...)

			// 借贷记录清算
			debtRecord.LiquidateTime = action.blocktime
			debtRecord.PreStatus = debtRecord.Status
			debtRecord.Status = pty.IssuanceUserStatusSystemLiquidate
			issu.DebtRecords = append(issu.DebtRecords[:index], issu.DebtRecords[index+1:]...)
			issu.InvalidRecords = append(issu.InvalidRecords, debtRecord)
		} else {
			debtRecord.PreStatus = debtRecord.Status
			debtRecord.Status = pty.IssuanceUserStatusWarning
		}

		log := action.GetFeedReceiptLog(issu, debtRecord)
		logs = append(logs, log)
	}

	// 保存
	issu.LatestLiquidationPrice = getLatestLiquidationPrice(issu)
	issu.LatestExpireTime = getLatestExpireTime(issu)
	collDB.Save(action.db)
	kv = append(kv, collDB.GetKVSet()...)

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

// 超时清算
func (action *Action) expireLiquidation(issu *pty.Issuance) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	collDB := &IssuanceDB{*issu}
	for index, debtRecord := range issu.DebtRecords {
		if debtRecord.ExpireTime - ExpireWarningTime > action.blocktime {
			continue
		}

		if debtRecord.ExpireTime >= action.blocktime {
			getGuarantorAddr, err := getGuarantorAddr(action.db)
			if err != nil {
				if err != nil {
					clog.Error("systemLiquidation", "getGuarantorAddr", err)
					continue
				}
			}

			// 抵押物转移
			receipt, err := action.coinsAccount.ExecTransferFrozen(action.fromaddr, getGuarantorAddr, action.execaddr, debtRecord.CollateralValue)
			if err != nil {
				clog.Error("systemLiquidation", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", debtRecord.CollateralValue, "err", err)
				continue
			}
			logs = append(logs, receipt.Logs...)
			kv = append(kv, receipt.KV...)

			// 借贷记录清算
			debtRecord.LiquidateTime = action.blocktime
			debtRecord.PreStatus = debtRecord.Status
			debtRecord.Status = pty.IssuanceUserStatusExpireLiquidate
			issu.DebtRecords = append(issu.DebtRecords[:index], issu.DebtRecords[index+1:]...)
			issu.InvalidRecords = append(issu.InvalidRecords, debtRecord)
		} else {
			debtRecord.PreStatus = debtRecord.Status
			debtRecord.Status = pty.IssuanceUserStatusExpire
		}

		log := action.GetFeedReceiptLog(issu, debtRecord)
		logs = append(logs, log)
	}

	// 保存
	issu.LatestLiquidationPrice = getLatestLiquidationPrice(issu)
	issu.LatestExpireTime = getLatestExpireTime(issu)
	collDB.Save(action.db)
	kv = append(kv, collDB.GetKVSet()...)

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

// 价格计算策略
func pricePolicy(feed *pty.IssuanceFeed) float32 {
	var totalPrice float32
	var totalVolume int64
	for _, volume := range feed.Volume {
		totalVolume += volume
	}

	for i, price := range feed.Price {
		totalPrice += price * float32(float64(feed.Volume[i])/float64(totalVolume))
	}

	return totalPrice
}

// IssuanceFeed 喂价
func (action *Action) IssuanceFeed(feed *pty.IssuanceFeed) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	if feed == nil || len(feed.Price) == 0 || len(feed.Price) != len(feed.Volume) {
		clog.Error("IssuancePriceFeed", types.ErrInvalidParam)
		return nil, types.ErrInvalidParam
	}

	// 是否后台管理用户
	if !isRightAddr(priceFeedKey, action.fromaddr, action.db) {
		clog.Error("IssuancePriceFeed", "addr", action.fromaddr, "error", "Address has no permission to feed price")
		return nil, pty.ErrPermissionDeny
	}

	price := pricePolicy(feed)
	if price == 0 || price == -1 {
		clog.Error("IssuancePriceFeed", "price", price, "err", pty.ErrPriceInvalid)
		return nil, pty.ErrPriceInvalid
	}

	collIDRecords, err := queryIssuanceByStatus(action.localDB, pty.IssuanceStatusCreated)
	if err != nil {
		clog.Error("IssuancePriceFeed", "get issuance record error", err)
		return nil, err
	}

	for _, collID := range collIDRecords {
		issu, err := queryIssuanceByID(action.db, collID.IssuanceId)
		if err != nil {
			clog.Error("IssuancePriceFeed", "Issuance ID", issu.IssuanceId, "get issuance record by id error", err)
			continue
		}

		// 超时清算判断
		if issu.LatestExpireTime - ExpireWarningTime <= action.blocktime {
			receipt, err := action.expireLiquidation(issu)
			if err != nil {
				clog.Error("IssuancePriceFeed", "Issuance ID", issu.IssuanceId, "expire liquidation error", err)
				continue
			}
			logs = append(logs, receipt.Logs...)
			kv = append(kv, receipt.KV...)
		}

		// 系统清算判断
		receipt, err := action.systemLiquidation(issu, price)
		if err != nil {
			clog.Error("IssuancePriceFeed", "Issuance ID", issu.IssuanceId, "system liquidation error", err)
			continue
		}
		logs = append(logs, receipt.Logs...)
		kv = append(kv, receipt.KV...)
	}

	var priceRecord pty.AssetPriceRecord
	priceRecord.BtyPrice = price
	priceRecord.RecordTime = action.blocktime

	// 喂价记录
	pricekv := &types.KeyValue{Key: calcIssuancePriceKey(string(action.blocktime)), Value: types.Encode(&priceRecord)}
	action.db.Set(pricekv.Key, pricekv.Value)
	kv = append(kv, pricekv)

	// 最近喂价记录
	pricekv = &types.KeyValue{Key: calcIssuanceLatestPriceKey(), Value: types.Encode(&priceRecord)}
	action.db.Set(pricekv.Key, pricekv.Value)
	kv = append(kv, pricekv)

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

// IssuanceClose 终止借贷
func (action *Action) IssuanceClose(close *pty.IssuanceClose) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	issuance, err := queryIssuanceByID(action.db, close.IssuanceId)
	if err != nil {
		clog.Error("IssuanceClose", "IssuanceId", close.IssuanceId, "err", err)
		return nil, err
	}

	for _, debtRecord := range issuance.DebtRecords {
		if debtRecord.Status != pty.IssuanceUserStatusClose {
			clog.Error("IssuanceClose", "IssuanceId", close.IssuanceId, "addr", action.fromaddr, "execaddr", action.execaddr, "err", pty.ErrIssuanceRecordNotEmpty)
			return nil, pty.ErrIssuanceRecordNotEmpty
		}
	}

	clog.Debug("IssuanceClose", "ID", close.IssuanceId)

	issu := &IssuanceDB{*issuance}
	preStatus := issu.Status
	issu.Status = pty.IssuanceStatusClose
	issu.Save(action.db)
	kv = append(kv, issu.GetKVSet()...)

	receiptLog := action.GetCloseReceiptLog(&issu.Issuance, preStatus)
	logs = append(logs, receiptLog)

	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

// 查找借贷
func queryIssuanceByID(db dbm.KV, IssuanceID string) (*pty.Issuance, error) {
	data, err := db.Get(Key(IssuanceID))
	if err != nil {
		clog.Debug("queryIssuanceByID", "error", err)
		return nil, err
	}

	var issu pty.Issuance
	err = types.Decode(data, &issu)
	if err != nil {
		clog.Debug("queryIssuanceByID", "decode", err)
		return nil, err
	}
	return &issu, nil
}

func queryIssuanceByStatus(localdb dbm.Lister, status int32) ([]*pty.IssuanceRecord, error) {
	data, err := localdb.List(calcIssuanceStatusPrefix(status), nil, DefultCount, ListDESC)
	if err != nil {
		clog.Debug("queryIssuancesByStatus", "error", err)
		return nil, err
	}

	var colls []*pty.IssuanceRecord
	var issu pty.IssuanceRecord
	for _, collBytes := range data {
		err = types.Decode(collBytes, &issu)
		if err != nil {
			clog.Debug("queryIssuancesByStatus", "decode", err)
			return nil, err
		}
		colls = append(colls, &issu)
	}

	return colls, nil
}

func queryIssuanceByAddr(localdb dbm.Lister, addr string) ([]*pty.IssuanceRecord, error) {
	data, err := localdb.List(calcIssuanceAddrPrefix(addr), nil, DefultCount, ListDESC)
	if err != nil {
		clog.Debug("queryIssuancesByAddr", "error", err)
		return nil, err
	}

	var colls []*pty.IssuanceRecord
	var issu pty.IssuanceRecord
	for _, collBytes := range data {
		err = types.Decode(collBytes, &issu)
		if err != nil {
			clog.Debug("queryIssuancesByAddr", "decode", err)
			return nil, err
		}
		colls = append(colls, &issu)
	}

	return colls, nil
}

func queryIssuanceRecordByStatus(localdb dbm.Lister, status int32) ([]*pty.IssuanceRecord, error) {
	data, err := localdb.List(calcIssuanceRecordStatusPrefix(status), nil, DefultCount, ListDESC)
	if err != nil {
		clog.Debug("queryIssuanceRecordByStatus", "error", err)
		return nil, err
	}

	var colls []*pty.IssuanceRecord
	var issu pty.IssuanceRecord
	for _, collBytes := range data {
		err = types.Decode(collBytes, &issu)
		if err != nil {
			clog.Debug("queryIssuancesByStatus", "decode", err)
			return nil, err
		}
		colls = append(colls, &issu)
	}

	return colls, nil
}
