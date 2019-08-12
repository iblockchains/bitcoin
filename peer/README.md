<h1>对等节点 Peer</h1>
<h5>包peer为创建和管理比特币网络提供了共同的基础</h5>
<h5>这个包是有意设计的，这样它就可以做为独立的包给任何需要完整bitcoin peer特性依赖的项目使用</h5>
<h5>Package peer provides a common base for creating and managing bitcoin network peers.</h5>
<h5>This package has intentionally been designed so it can be used as a standalone package for any projects needing a full featured bitcoin peer base to build on.</h5>

<h2>概览 Overview</h2>
<h5>这个包建立在wire包基础上，为了简化生成完整功能对等结节的流程 wire 包提供了表述 bitcoin wire protocol 的基本的必要原语。</h5>
<h5>peer包提供的主要功能如下：</h5>
<ul>
<li>1. 为处理基于点对点协议的比特币通信提供一个基本的并发安全的节点。</li>
<li>2. 完整的比特币协议消息读写</li>
<li></li>
<li></li>
</ul>
<h5>1. 为处理基于点对点协议的比特币通信提供一个基本的并发安全的节点。</h5>
<h5>2. 完整的比特币协议消息读写</h5>
<h5>3. 自动处理包含协议版本沟通的握手流程的启动</h5>
<h5>4. 向外的异步消息队列，附带一个可选的已发消息提示通道</h5>
<h5>This package builds upon the wire package, which provides the fundamental primitives necessary to speak the bitcoin wire protocol, in order to simplify the process of creating fully functional peers. </h5>
<h5>A quick overview of the major features peer provides are as follows:</h5>
<h5>1. Provides a basic concurrent safe bitcoin peer for handling bitcoin communications via the peer-to-peer protocol.</h5>
<h5>2. Full duplex reading and writing of bitcoin protocol messages</h5>
<h5>3. Automatic handling of the initial handshake process including protocol version negotiation</h5>
<h5>4. Asynchronous message queueing of outbound messages with optional channel for notification when the message is actually sent</h5>