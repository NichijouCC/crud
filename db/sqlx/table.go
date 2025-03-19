package sqlx

// ITable 定义表操作的基本接口
type ITable interface {
	TableName() string               // 获取表名
	ColumnsMap() map[string]struct{} // 获取允许过滤的字段列表
	Columns() []string               // 获取表的所有列名
	GetId() int64                    // 获取主键ID
}

// ITableUpdate 定义表更新操作的接口
type ITableUpdate interface {
	TableName() string // 获取表名
	Columns() []string // 获取表的所有列名
	GetId() int64      // 获取主键ID
}
