## Uniswap V2

### 1. 什么是 AMM？它和订单簿有什么区别？

AMM：自动做市商

传统的订单簿则把买卖盘按报价排列，由互相匹配的订单成交。

AMM 是 需要 LP 先往 Pair （Pool）合约里放入两种代币；交易者把一种代币交给池子，从池子取走另一种。池子按储备量和公式决定兑换数量，不用等另一位交易者挂出相反方向的订单。

它们两者都能提供流动性，但交易者面对的直接对手和定价方式不同。

一个 Pair 只对应一组无序的 ERC-20 代币，Factory 负责创建和索引。外部市场价格不会自动写进 Pair；若池内价格偏离外部市场，需要**套利交易**才可能把池内价格推回去。

### 2. AMM 的核心算法是什么？

V2 的核心是恒定乘积：`x × y = k` ,设交易前两种代币储备为 `x`、`y`，在不收手续费的理想情况下，一笔交易须满足 `(x + Δx)(y - Δy) = x × y`。因此买走越多 `y`，剩余 `y` 越少，后续单位 `y` 越贵；池内边际价格可从储备比 `y/x` 理解。

输入的需要收 0.30% 交易费，也就是只有扣除手续费的token数量参与兑换。所以上面的是理想不收手续费的情况  `k` 值不会变，有手续费的情况，`k` 值还是会增长的，合约最后会校验扣除手续费后的余额乘积不能小于交易之前的余额乘积。

### 3. 什么是滑点？交易时怎样避免过高？

**滑点**主要指看到报价到上链成交之间，池子状态变化使实际结果偏离报价。

实际交易：通过页面点的滑点设置，自动计算当前交易的`amountInMax`（精确输出），也就是最小能接受这个数量，低于这个数量就会失败，但是会消耗gas。这也是不能忽视的磨损。

### 4. 什么是三明治攻击？怎样降低风险？

攻击者看到一笔待执行的交易，先买入一笔同一池且同向的交易，先把池内的 `token` 价格推高，让受害者在其允许的滑点范围内以更差价格成交，然后再卖回获利。这要求攻击者能影响交易前后的排序，并且利润足以覆盖手续费、gas 等成本。如果用户把 `amountOutMin` 设得过低，相当于给了攻击者更大的可获利空间。

收缩滑点，也就是`amountOutMin`，可以使用交易入口提供的抗 MEV 路由或者私有提交来降低风险。

### 5. LP 的作用是什么？

LP是流动性提供者提供资产份额的证明，也是参与手续费分成的凭证。

手续费是保存在pool中的，撤出流动性时，按照持有的 LP 占总LP的比例，取出相应比例的两种 `token` ,包含自己抵押的资产和手续费，但是这里取出的数量不是原先的数量，存在“无常损失”。

### 6. 什么是无常损失？如何计算？

无常损失是指，同样的资产 做和不做 LP 的价差。因为 Pair 使用恒定公式保持两种 `token` 的比例，也就是价格，

它是“做 LP”相对“原样持有最初两种资产”的价值差，不必表示以计价货币计算的绝对亏损。价格变动后，套利使恒定乘积池自动卖出上涨资产、买入下跌资产，所以资产组合与单纯持有不同。若价格回到原点，这一模型下的相对差额也会消失；若此时退出，差额就成为已经实现的相对结果。

设入池时两种资产等值、池子遵循 V2 恒定乘积、外部价格充分套利，且**忽略手续费、协议费、gas、税费和整数取整**。若资产 A 相对 B 的价格变为原来的 `r` 倍：

`IL = V_LP / V_hold - 1 = 2√r / (1+r) - 1`。

例如起初持有 `1 A + 100 B`，价格为 `1 A = 100 B`，总值 `200 B`。把它们做 LP 后，若 A 涨到 `400 B`（`r=4`），理想模型下池内份额变成 `0.5 A + 200 B`，值 `400 B`；原样持有则值 `500 B`，相对差 `400/500 - 1 = -20%`。实际 LP 净收益还要加收到的手续费、减各项成本；V3 集中流动性不直接套用这一条全区间公式

### 7. V2 Core 和 Periphery 的核心接口怎么读？
分工：Core 的 Factory 创建 Pair，Pair 持有资产并校验不变量；Periphery 的 Library 做地址与数量计算，Router02 把转账、路径和最差成交条件组合起来。用户入口是 Router02；Pair 的 `mint`、`burn`、`swap` 是低层接口，直接调用前要自己做好转账、报价和边界检查。

| 模块 | 核心接口 | 用途与阅读要点 |
| --- | --- | --- |
| Core / Factory | `createPair(tokenA, tokenB)`；`getPair(tokenA, tokenB)`；`allPairs(i)`、`allPairsLength()`；| 为一对代币部署唯一 Pair；`getPair` 查询**已部署**地址。两个 setter 仅授权地址可调用。`PairCreated` 记录创建事件。 |
| Core / Pair：状态 | `factory()`、`token0()`、`token1()`；`getReserves()`；`price0CumulativeLast()`、`price1CumulativeLast()`；`kLast()`； | 储备顺序按代币地址排序；累计价格用于时间加权价格资料；`kLast` 是最近一次流动性事件后的乘积基准，不是每次 swap 后自动更新。 |
| Core / Pair：资金 | `mint(to)`；`burn(to)`；`swap(amount0Out, amount1Out, to, data)`；`skim(to)`；`sync()` | `mint` 按刚转入的代币铸 LP；`burn` 销毁已转到 Pair 的 LP 并返还资产；`swap` 转出资产并核对实际转入及手续费调整后的不变量。`data` 非空可触发 flash swap 回调；`skim` 转走余额超出储备的部分，`sync` 把储备同步到余额。 |
| Core / Pair：回调 | `IUniswapV2Callee.uniswapV2Call(sender, amount0, amount1, data)` | `swap` 的 `data` 非空时，Pair 在转出代币后调用接收合约的回调，交易结束前仍须满足还款及手续费约束。 |
| Periphery / Router02：流动性 | `addLiquidity()`、`addLiquidityETH()`；`removeLiquidity()`、`removeLiquidityETH()`；`removeLiquidityWithPermit()`、`removeLiquidityETHWithPermit()` | `Desired` 是用户愿意投入的上限，`Min` 是逐资产下限，`deadline` 限制有效时间；ETH 入口由 Router 包装/解包 WETH。移除需先授权 LP，或使用 Permit 版本。 |
| Periphery / Router02：交易 | `swapExactTokensForTokens()`、`swapTokensForExactTokens()`；`swapExactETHForTokens()`、`swapTokensForExactETH()`、`swapExactTokensForETH()`、`swapETHForExactTokens()`；`getAmountsOut()`、`getAmountsIn()` | `Exact...For...` 固定输入、设置最少输出；`...ForExact...` 固定输出、设置最多输入。`path` 列出逐跳代币，ETH 入口用 `WETH` 对接 Pair。Router 用 Library 报价并依次调用 Pair。报价只反映读取时的储备，不保证执行结果。 |
| Periphery / Router02：特殊代币 | `swapExactTokensForTokensSupportingFeeOnTransferTokens()` 及 ETH 变体；`removeLiquidityETHSupportingFeeOnTransferTokens()` 及 Permit 变体 | 固定输入 swap 按各 Pair 实际收到的输入计算，最终用接收方余额变化检查 `amountOutMin`；不提供固定输出版本。移除流动性版本把 Router 实际收到的 token 转给用户，但 `amountTokenMin` 是移除时的检查，若最后一次转账也收费，不能据此保证用户最终到账量。 |

加流动性则是 `Router02._addLiquidity` 先决定实际投入量，再转入 Pair，最后 `Pair.mint()` 铸造 LP。

`amountADesired/amountBDesired` 是意愿上限，`amountAMin/amountBMin` 是已有储备时实际投入的下限；空池首次加流动性直接按 `Desired` 建立初始比例，不能靠 `Min` 校正初始价格。

交易接口里的 `amountOutMin/amountInMax` 保护最终成交，`path[0]` 是输入代币，`path` 末项是输出代币。

普通精确输入交易可以建议的串成 `Router02.getAmountsOut() → 转入首个 Pair → 各 Pair.swap() → 最终接收地址`。