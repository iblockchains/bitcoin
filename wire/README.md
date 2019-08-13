<h1>联网 Wire</h1>
<br>
<h5>wire包是对比特币联网协议的实现。提供了一套全面的测试，测试覆盖率为100%，以确保正确的功能。</h5>
<h5>Package wire implements the bitcoin wire protocol. A comprehensive suite of tests with 100% test coverage is provided to ensure proper functionality.
This package has intentionally been designed so it can be used as a standalone package for any projects needing to interface with bitcoin peers at the wire protocol level.</h5>
<h3>比特币消息概览 Bitcoin Message Overview</h3>
<h5>比特币协议包含节点间的消息交换。每条消息之前都有一个头（header)，用于标识有关消息的信息，例如它的是比特币网络的哪个部分、类型、大小和验证有效性的较验和。</h5>
<h5>The bitcoin protocol consists of exchanging messages between peers. Each message is preceded by a header which identifies information about it such as which bitcoin network it is a part of, its type, how big it is, and a checksum to verify validity. All encoding and decoding of message headers is handled by this package.</h5>