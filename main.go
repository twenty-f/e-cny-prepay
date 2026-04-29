package main

import (
	"errors"
	"fmt"
)

// ContractStatus 表示预付资金合约所处状态。
type ContractStatus int

const (
	Locked ContractStatus = iota
	Consuming
	Completed
	Refunded
)

func (s ContractStatus) String() string {
	switch s {
	case Locked:
		return "Locked"
	case Consuming:
		return "Consuming"
	case Completed:
		return "Completed"
	case Refunded:
		return "Refunded"
	default:
		return "Unknown"
	}
}

// PrepayContract 模拟“数字人民币预付资金管理”合约。
type PrepayContract struct {
	ContractID   string
	Consumer     string
	Merchant     string
	TotalAmount  int
	Balance      int
	PerUseAmount int
	Status       ContractStatus
}

// InitContract 初始化合约，写入总金额与单次扣费，并将状态设置为 Locked 后进入 Consuming。
func (pc *PrepayContract) InitContract(totalAmount, perUseAmount int) error {
	if totalAmount <= 0 {
		return errors.New("totalAmount 必须大于 0")
	}
	if perUseAmount <= 0 {
		return errors.New("perUseAmount 必须大于 0")
	}
	if perUseAmount > totalAmount {
		return errors.New("perUseAmount 不能大于 totalAmount")
	}
	if totalAmount%perUseAmount != 0 {
		return errors.New("totalAmount 必须是 perUseAmount 的整数倍")
	}

	pc.TotalAmount = totalAmount
	pc.Balance = totalAmount
	pc.PerUseAmount = perUseAmount
	pc.Status = Locked

	fmt.Printf("[Init] 合约 %s 初始化完成，锁定资金=%d，状态=%s\n", pc.ContractID, pc.TotalAmount, pc.Status)

	// 模拟资金从“锁定”切换到“可消费中”
	pc.Status = Consuming
	fmt.Printf("[Init] 合约 %s 进入消费状态，状态=%s\n", pc.ContractID, pc.Status)
	return nil
}

// ConsumeOnce 模拟一次到店消费扣费。
func (pc *PrepayContract) ConsumeOnce() error {
	if pc.Status != Consuming {
		return fmt.Errorf("当前状态=%s，不允许消费", pc.Status)
	}

	deduct := pc.PerUseAmount
	if pc.Balance < pc.PerUseAmount {
		// 最后一次不足单次额度时，扣掉剩余全部余额，确保不会出现负数。
		deduct = pc.Balance
	}

	pc.Balance -= deduct
	fmt.Printf("[Consume] 合约 %s 本次扣费=%d，剩余余额=%d，状态=%s\n", pc.ContractID, deduct, pc.Balance, pc.Status)

	if pc.Balance == 0 {
		pc.Status = Completed
		fmt.Printf("[Consume] 合约 %s 余额已清零，状态切换为=%s\n", pc.ContractID, pc.Status)
	}

	return nil
}

// Refund 模拟商家跑路或用户申请退款，将未消费余额退回。
func (pc *PrepayContract) Refund() error {
	if pc.Status == Completed {
		return errors.New("合约已完成，无可退款余额")
	}
	if pc.Status == Refunded {
		return errors.New("合约已退款，请勿重复操作")
	}

	refundAmount := pc.Balance
	pc.Balance = 0
	pc.Status = Refunded

	fmt.Printf("[Refund] 合约 %s 触发退款，退款金额=%d，状态=%s\n", pc.ContractID, refundAmount, pc.Status)
	return nil
}

func main() {
	contract := &PrepayContract{
		ContractID: "CNY-PREPAY-0001",
		Consumer:   "Alice",
		Merchant:   "FitGym",
	}

	fmt.Println("==== 数字人民币预付资金管理模拟开始 ====")

	// 用户充值 500 元，单次扣费 100 元。
	if err := contract.InitContract(500, 100); err != nil {
		fmt.Printf("初始化失败: %v\n", err)
		return
	}

	// 消费 2 次。
	if err := contract.ConsumeOnce(); err != nil {
		fmt.Printf("第 1 次消费失败: %v\n", err)
		return
	}
	if err := contract.ConsumeOnce(); err != nil {
		fmt.Printf("第 2 次消费失败: %v\n", err)
		return
	}

	// 触发退款（例如商家跑路）。
	if err := contract.Refund(); err != nil {
		fmt.Printf("退款失败: %v\n", err)
		return
	}

	fmt.Printf("==== 生命周期结束：最终状态=%s，最终余额=%d ====\n", contract.Status, contract.Balance)
}
