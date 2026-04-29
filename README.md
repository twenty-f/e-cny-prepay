# 🏦 e-CNY Prepay Smart Contract | 数字人民币预付资金防跑路合约

![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)
![Smart Contract](https://img.shields.io/badge/Smart_Contract-State_Machine-blueviolet)
![Status](https://img.shields.io/badge/Status-Completed-success)

## 📖 项目背景
针对传统服务业（如健身房、教培机构）频发的“商家卷款跑路”痛点，本项目基于数字人民币（e-CNY）的智能合约特性，设计了一套**预付资金定向元管控制**。通过代码强制约束资金流向与拨付条件，重塑消费者与商家之间的信任机制。

## 🏗️ 核心架构与亮点
- **有限状态机 (FSM) 流转**：严格控制资金池状态在 `Locked`（建仓锁定） -> `Consuming`（按次核销消费） -> `Refunded`（突发退款清算） 之间单向流转，杜绝非法状态跃迁。
- **前置业务防御校验**：在合约初始化阶段引入严格的取模运算（`totalAmount % perUseAmount == 0`），从根源拦截脏数据入链，保障账本逻辑的绝对自洽。
- **并发安全与原子性保证**：在本地网关/模拟层引入 `sync.Mutex` 互斥锁，有效防范高并发场景下的竞态条件（Race Condition），确保每一次扣费操作的资金原子性。

## 🚀 业务生命周期演示
1. **注资锁定**：消费者打入 500 元，合约状态变更为 `Locked`。
2. **按次核销**：每次消费触发智能扣费（如 100 元），资金释放给商家，余额精准递减。
3. **退款清算**：触发商家违约或退款事件，剩余可用余额强制原路退回，彻底阻断资金卷逃风险。