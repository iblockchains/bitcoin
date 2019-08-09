<p align="center"><h1>第一章 比特币网络 The Bitcoin Network</h1></p>

<h2>一、基本概念 Basic Concept</h2>

<h3> 1.点对点网络架构 Peer-to-Peer Network Architecture</h3>
<h5>比特币是在互联网上构建的点对点网络架构.网络中没有服务器、没有集中式服务和层次结构.P2P网络互惠性是节点参与的激励因素,节点提供服务的同时也消费服务.P2P网络有天生的弹性、去中心和开放的特点.</h5>
<h5>Bitcoin is structured as a peer-to-peer network architecture on top of the internet.There is no server, no centralized service, and no hierarchy within the network.Nodes in a P2P network both provide and consume services at the same time with reciprocity acting as the incentive for participation.P2P networks are inherently resilient, decentralized, and open</h5>

<h3> 2.节点类型和角色 Node Types and Roles </h3>
<h5>一个比特币节点是：路由、区块链数据库、采矿和钱包服务的功能集合.所有的节点为都有路由、验证和传播交易和区块链和发现节点及保持节点之前的连接。</h5>
<h5>完整的节点，具有所有功能，可以自主和权威验证所有的交易等</h5>
<h5>简易交易验证(spv)节点或轻量级节点，只保存区块链的一个子集，并通过SPV方法验证交易</h5>
<h5>挖矿节点或矿工，运行在特殊的硬件上，通过工作量证明(Proof-of-Work)来竞争新区块的创建权</h5>
<h5>A bitcoin node is a collection of functions: routing, the blockchain database, mining, and wallet services.All nodes routing, validate and propagate transactions and blocks,and discover and maintain connections to peers.</h5>
<h5>Full nodes can autonomously and authoritatively verify any transaction.</h5>
<h5>SPV nodes or lightweight nodes maintian only a subset of the blockchain and verify transactions using a method called simplified payment verification, or SPV.</h5>
<h5>Mining nodes compete to create new blocks by running specialized hardware to solve the Proof-of-Work algorithm</h5>

<h3> 3. 扩展的比特币网络 The Extended Bitcoin Network </h3>
<h5>扩展的比特币网络包含运行在比特币协议上的网络和其它运行在特殊协议上的网络</h5>
<h5>The extended bitcoin network includes the network running the bitcoin P2P protocol, described earlier, as well as nodes running specialized protocols.</h5>

<img src="https://github.com/iblockchains/bitcoin/blob/master/img/008-Differnt-types-of-nodes-on-the-extended-bitcoin-network.png" alt="Source:Andreas M. Antonopoulos">
<h3> 4. 比特币中继网络 Bitcoin Relay Networks </h3>

<h3> 5. 网络发现 Network Discovery </h3>

<h3> 6. 完整节点 Full Nodes </h3>

<h3> 7. 节点间 “库存”交换 Exchanging “Inventory”

<h3> 8. 简易付款验证节点 Simplified Payment Verification (SPV) Nodes</h3>

<h3> 9. Bloom 过滤器 Bloom Filters</h3>

<h3> 10. 加密和验证的连接 Encrypted and Authenticated Connections </h3>
<br/>

<h2> 二、编程 Program </h2>

