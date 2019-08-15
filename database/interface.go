package database

// DB 提供了一个通用接口用来存储比特币区块和相关的元数据.
// 此接口目的是隐藏后台数据存储的实际实现.
// RegisterDriver 函数可以用来添加新的后台数据存储方法
//
// 此接口分为两种不同的功能
//
// 第一类是支持 bucket (存储桶)的原子元数据存储。这通过
// 数据库事务来实现
//
// 第二类是通用区块存储.功能的隔离是因为区块存储和元数据存储
// 的实现机制可能一样也可能不一样.
//
// provides a generic interface that is used to store bitcoin blocks and
// related metadata.  This interface is intended to be agnostic to the actual
// mechanism used for backend data storage.  The RegisterDriver function can be
// used to add a new backend data storage method.
//
// This interface is divided into two distinct categories of functionality.
//
// The first category is atomic metadata storage with bucket support.  This is
// accomplished through the use of database transactions.
//
// The second category is generic block storage.  This functionality is
// intentionally separate because the mechanism used for block storage may or
// may not be the same mechanism used for metadata storage.  For example, it is
// often more efficient to store the block data as flat files while the metadata
// is kept in a database.  However, this interface aims to be generic enough to
// support blocks in the database too, if needed by a particular backend.
type DB interface {
	// Close 完全关闭数据库并同步所有数据.
	// 它会阻塞直到数据库所有事务都完成
	// cleanly shuts down the database and syncs all data.  It will
	// block until all database transactions have been finalized (rolled
	// back or committed).
	Close() error
}
