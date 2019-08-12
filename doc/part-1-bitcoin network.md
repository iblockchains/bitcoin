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
<h5>比特币矿工参与的是工作量证明和扩展区块链，一个对时间特别敏感的竞争。</h5>
<h5>比特币中继网络是一个试图最小化矿工之间传递区块的延迟的网络，像FIBRE和Falcon</h5>
<h5>Bitcoin miners are engaged in a time-sensitive competition to solve the Proof-of-Work problem and extend the blockchain</h5>
<h5>A Bitcoin Relay Network is a network that attempts to minimize the latency in the transmission of blocks between miners, like FIBRE and Falcon</h5>
<h3> 5. 网络发现 Network Discovery </h3>
<h5>当一个新节点启动时，它需要发现网络中的其它节点才能加入。当建立一个连接时，节点会发发送一个传递版本信息的握手包</h5>
<h5>版本信息 Version Message:</h5>
<h5>nVersion: 比特币P2P协议版本 The bitcoin P2P protocol</h5>
<h5>nLocalServices: 一个本地节点所支持服务的列表 A list of local services supported by the node</h5>
<h5>nTime: 当前时间 The current time</h5>
<h5>addrYou: 远程节点IP The IP address of the remote node as seen from this node</h5>
<h5>addrMe: 本地节点IP The IP address of the local node, as discovered by the local node</h5>
<h5>subver: 该节点运行的软件类型的子版本 A sub-version showing the type of software running on this node</h5>
<h5>BestHeight: 该节点区块链的区块高度 The block height of this node’s blockchain</h5>
<h5>本地节点在接收到版本信息时，会先检测远程节点的nVersion是否兼容，如果兼容本地节点会建立连接并发送verack确认信息</h5>
<img src="https://github.com/iblockchains/bitcoin/blob/master/img/008-handshake.png">
<h5>新节点如何发现其它节点？一个方法是通过DNS服务器提供的包含其它比特币节点IP的DNS种子列表来查询DNS。比特币核心客户端包含五个不同的DNS种子。</h5>
<h5>另一个方法是提供至少一个比特币节点的IP给新节点，并与它们建立连接。新节点给它的邻居们发送addr消息，邻居们将会将addr信息转发给他们的邻居，以确保加入的节点变得众所周知和有效连接;然后，新节点给邻居们发送getaddr消息，请求一个其它节点的IP地址列表。这样，新节点就可以找到其它节点去链接，并在网络中广播它的存在以便其它节点可以找到它。</h5>
<img src="https://github.com/iblockchains/bitcoin/blob/master/img/008-address-propagation-and-discovery.png">
<h5>When a new node boots up, it must discover other bitcoin nodes on the network in order to participate.Upon establishing a connection, the node will start a "handshake" by transmitting a  version message</h5>
<h5>The local peer receiving a  version message will examine the remote peer’s reported  nVersion and decide if the remote peer is compatible. If the remote peer is compatible, the local peer will acknowledge the  version message and establish a connection by sending a  verack </h5>
<h5>How does a new node find peers? The first method is to query DNS using a number of "DNS seeds", which are DNS servers that provide a list of IP addresses of bitcoin nodes. The Bitcoin Core client contains the names of five different DNS seeds.</h5>
<h5>Alternatively, give a bootstrapping node the IP address of at least one bitcoin node.The new node will send an addr message to its neighbors.Once one or more connections are established, the new node will send an  addr mes‐
sage containing its own IP address to its neighbors. The neighbors will, in turn, forward the  addr message to their neighbors, ensuring that the newly connected node becomes well known and better connected. Additionally, the newly connected node can send  getaddr to the neighbors, asking them to return a list of IP addresses of other peers. That way, a node can find peers to connect to and advertise its existence on the network for other nodes to find it</h5>
<h3> 6. 完整节点 Full Nodes </h3>
<h5>完整的节点保存一个完整的区块链和所有的交易</h5>
<h5>Full nodes are nodes that maintain a full blockchain with all transactions.</h5>
<h3> 7. 节点间 “库存”交换 Exchanging “Inventory”
<h5>当一个完整节点连接上其它节点时，它做的第一件事就是尝试构建一个完整的区块链。区块链的同步过程从版本信息开始，因为它包含当前节点的高度BestHeight。一个节点可以看到其对等节点上的版本信息，并比较相互之间所拥有的区块数量。节点间通过getblocks消息交换本地区块链最上面的区块，并凭此推断自身的区块链是否长于其它节点。</h5>
<h5>拥有较长区块链的节点，会通过inv 发送前500个区块的哈希数给其它节点。缺少这些区块的节点会通过使用inv发送一系列的getdata包含区块哈希数的消息来获取它们</h5>
<h5>假设，例如，一个节点只有创世区块。它将从对等节点上接收一个包含下500个区块哈希数的 inv 消息。它将从所有相连的节点请求这些区块，分散负载并保证请求不会大压垮任一节点。节点会追踪每个连接的对等节点上的区块请求数，检查它是否超过上限(MAX_BLOCKS_IN_TRANSIT_PER_PEER)</h5>
<img src="https://github.com/iblockchains/bitcoin/blob/master/img/008-node-synchronizing-the-blockchain-by-retrieving-blocks-from-a-peer.png">
<h5>The first thing a full node will do once it connects to peers is try to construct a complete blockchain.The process of syncing the blockchain starts with the  version message, because that contains  BestHeight , a node’s current blockchain height.A node will see the version message from its perrs, and be able to compare to how many blocks it has in its own blockchain.Peered nodes will exchange a getblocks message that contains the hash of the top block on their local blockchain, from this peers will be to deduce whether our blockchain is longer then others.</h5>
<h5>The peer that has the longer blockchain has more blocks than the other node will identify the first 500 blocks to share and transmit their hashes using an inv message.The node missing these blocks will then retrieve them, by issuing a series of getdata message requesting the full block data and identifying the requeted blocks using hashes from the inv message.</h5>
<h5>Let’s assume, for example, that a node only has the genesis block. It will then receive an  inv message from its peers containing the hashes of the next 500 blocks in the chain. It will start requesting blocks from all of its connected peers, spreading the load and ensuring that it doesn’t overwhelm any peer with requests.The node keeps track of how many blocks are “in transit” per peer connection, meaning blocks that it has requested but not received, checking that it does not exceed a limit ( MAX_BLOCKS_IN_TRANSIT_PER_PEER ).</h5>
<h3> 8. 简易付款验证节点 Simplified Payment Verification (SPV) Nodes</h3>
<h5>简易付款验证（SPV）允许设备无需储存完整的区块链也能操作。</h5>
<h5>SPV节点只下载区块的头，这样的区块链只有完整的千分之一。</h5>
<h5>完整的节点会检索所有的区块来生成UTXO（unspent transaction output）数据库。通过确认UTXO是否被使用，来确认交易的有效性。SPV节点则不能验证UTXO是否还未被支付。相反，SPV节点通过Merkle路径，在交易和包含交易的区块之间建立链接，通过检查在其上面的区块将它压在下面的深度来验证交易。</h5>
<h5>SPV节点通过getheaders消息来获取区块头。响应的节点会通过headers消息发送最多2000个区块的头信息。</h5>
<h5>A simplified payment verfication (SPV) method is used to allow devices to operate without storing the full blockchain.</h5>
<h5>SPV nodes download only the block headers, the resulting chain of blocks is 1000 times smaller then the full blockchain.</h5>
<h5>A full node will go through all blocks and builds a full database of UTXO (unspent transaction output), establishing the validity of the transaction by confirming that the UTXO remains unspent.An SPV node cannot validate whether the UTXO is unspent.Instead, the SPV node will establish a link between the transaction and the block that contains it, using a merkle path, and checks how deep the block is buried by a handful of blocks above it.</h5>
<h5>To get the block headers, SPV nodes use a  getheaders message.The responding peer will send up to 2,000 block headers using a single  headers message.</h5>
<h3> 9. Bloom 过滤器 Bloom Filters</h3>
<h5>Bloom过滤器提供了一种有效的方式来表达搜索模式，同时保护隐私。SPV节点使用它们向对等节点请求匹配特定模式的交易，同时不披露它们查询的具体的地址、密钥或交易</h5>
<h4>Bloom过滤器工作原理</h4>
<h5>Bloom过滤器的实现是由一个可变长度（N）的二进制数组（N位二进制数构成一个位域）和数量可变（M）的一组哈希函数组成。这些函数为确定性函数，也就是说任何一个使用相同Bloom过滤器的节点通过该函数同样的输入都能得到同一个的结果。Bloom过滤器的准确性和私密性能通过改变长度（N）和哈希函数的数量（M）来调节。</h5>
<img src="https://github.com/iblockchains/bitcoin/blob/master/img/008-an-exapmle-of-a-simplistic-bloom-filter.png">
<h5>Bloom filters offer an efficient way to express a search pattern while protecting privacy.They are used by SPV nodes to ask their peers for transactions matching a specific pattern, without revealing exactly which addresses,keys, or transactions they are searching for.</h5>
<h4>How Bloom Filters Work</h4>
<h5>Bloom filters are implemented as a variable-size array of N binary digits (a bit field) and a variable number of M hash functions.The hash functions are generated deterministically, so that any node implementing a bloom filter will always use the same hash functions and get the same results for a specific input.By choosing different length (N) bloom filters and a different number (M) of hash functions, the bloom filter can be tuned, varying the level of accuracy and therefore privacy.</h5>
<h3> 10. 加密和验证的连接 Encrypted and Authenticated Connections </h3>
<br/>

<h2> 二、编程 Program </h2>

