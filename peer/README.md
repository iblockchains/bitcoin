<h1>对等节点 Peer</h1>
<h5>包peer为创建和管理比特币网络提供了共同的基础</h5>
<h5>这个包是有意设计的，这样它就可以做为独立的包给任何需要完整bitcoin peer特性依赖的项目使用</h5>
<h5>Package peer provides a common base for creating and managing bitcoin network peers.</h5>
<h5>This package has intentionally been designed so it can be used as a standalone package for any projects needing a full featured bitcoin peer base to build on.</h5>

<h2>概览 Overview</h2>
<h5>这个包建立在wire包基础上，为了简化生成完整功能对等结节的流程 wire 包提供了表述 bitcoin wire protocol 的基本的必要原语。</h5>
<h5>peer包提供的主要功能如下：</h5>
<ul>
<li><h5>1. 为处理基于点对点协议的比特币通信提供一个基本的并发安全的节点。</h5></li>
<li><h5>2. 完整的比特币协议消息读写</h5></li>
<li><h5>3. 自动处理包含协议版本沟通的握手流程的启动</h5></li>
<li><h5>4. 向外的异步消息队列，附带一个可选的已发消息提示通道</h5></li>
<li><h5>弹性化节点设置:</h5>
  <ul>
    <li><h5>调用方（Caller）负责创建传出连接（outgoing connections）并监听传入连接（incoming connections），因此他们可以根据自己的需要灵活地建立连接。</h5></li>
    <li><h5>用户代理名称和版本</h5></li>
    <li><h5>比特币网络</h5></li>
    <li><h5>发信支持服务</h5></li>
    <li><h5>支持的最大协议版本</h5></li>
    <li><h5>能够注册处理比特币协议消息的回调</h5></li>
  </ul>
</li>
<li><h5>库存消息( inventory message )的批处理和发送</h5></li>
<li><h5>定期自动发起和响应节点是否存活（keep-alive）</h5></li>
<li><h5>随机数生成和自连接检测</h5></li>
<li><h5>合理处理跟命令行相关的 Bloom 过滤器，当调用方(Caller)没有特别发信指定相关的标志时</h5>
  <ul>
    <li><h5>当被连节点的协议太高时，与之断开</h5></li>
    <li><h5>不调用旧协议版本的相关回调</h5></li>
  </ul>
</li>
<li><h5>节点数据分析快照表，如读和写的总字节数、远程地址、用户代理和协商协议版本</h5></li>
<li><h5>帮助程序函数pushing address、getblocks、getheaders和reject消息</h5>
  <ul>
    <li><h5>这些都可以通过标准消息输出函数进行发送，但是帮助函数提供了额外更优的功能，如双重过滤和地址随机</h5></li>
  </ul>
</li>
<li><h5>等待关机/断开连接的能力</h5></li>
<li><h5>综合测试覆盖率</h5></li>
</ul>
<h5>This package builds upon the wire package, which provides the fundamental primitives necessary to speak the bitcoin wire protocol, in order to simplify the process of creating fully functional peers. </h5>
<h5>A quick overview of the major features peer provides are as follows:</h5>
<ul>
  <li><h5>1. Provides a basic concurrent safe bitcoin peer for handling bitcoin communications via the peer-to-peer protocol.</h5></li>
  <li><h5>2. Full duplex reading and writing of bitcoin protocol messages</h5></li>
  <li><h5>3. Automatic handling of the initial handshake process including protocol version negotiation</h5></li>
  <li><h5>4. Asynchronous message queueing of outbound messages with optional channel for notification when the message is actually sent</h5></li>
  <li><h5>Flexible peer configuration:</h5>
  <ul>
    <li><h5>Caller is responsible for creating outgoing connections and listening for incoming connections so they have flexibility to establish connections as they see fit (proxies, etc)</h5></li>
    <li><h5>User agent name and version</h5></li>
    <li><h5>Bitcoin network</h5></li>
    <li><h5>Service support signalling</h5></li>
    <li><h5>Maximum supported protocol version</h5></li>
    <li><h5>Ability to register callbacks for handling bitcoin protocol messages</h5></li>
  </ul>
</li>
<li><h5>Inventory message batching and send trickling with known inventory detection and avoidance</h5></li>
<li><h5>Automatic periodic keep-alive pinging and pong responses</h5></li>
<li><h5>Random nonce generation and self connection detection</h5></li>
<li><h5>Proper handling of bloom filter related commands when the caller does not specify the related flag to signal support</h5>
<ul>
    <li><h5>Disconnects the peer when the protocol version is high enough</h5></li>
    <li><h5>Does not invoke the related callbacks for older protocol versions</h5></li>
  </ul>
</li>
<li><h5>Snapshottable peer statistics such as the total number of bytes read and written, the remote address, user agent, and negotiated protocol version</h5></li>
<li><h5>Helper functions pushing addresses, getblocks, getheaders, and reject messages</h5>
  <ul>
    <li><h5>These could all be sent manually via the standard message output function, but the helpers provide additional nice functionality such as duplicate filtering and address randomization</h5></li>
  </ul>
</li>
<li><h5>Ability to wait forshutdown/disconnect</h5></li>
<li><h5>Comprehensive test coverage</h5></li>
</ul>
<h2>案例 Examples</h2>
<h5><a href="https://godoc.org/github.com/btcsuite/btcd/peer#example-package--NewOutboundPeer">New Outbound Peer Example</a></h5>
<h5>说明了基本的初始和创建外部节点的流程。节点间的协议交流和协商。为了说明，节点还附上了简单的消息处理程序</h5>
<h5>Demonstrates the basic process for initializing and creating an outbound peer. Peers negotiate by exchanging version and verack messages. For demonstration, a simple handler for the version message is attached to the peer.</h5>