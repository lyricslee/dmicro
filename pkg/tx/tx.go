package tx

import (
	"github.com/go-xorm/xorm"
)

type TxFn func(*xorm.Session) error

// 事务函数
func ExecWithTransaction(session *xorm.Session, fn TxFn) (err error) {
	if err = session.Begin(); err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			session.Rollback()
			panic(p)
		} else if err != nil {
			session.Rollback()
		} else {
			err = session.Commit()
		}
	}()

	err = fn(session)
	return
}

type PipelineStmt struct {
	sqlorArgs []interface{}
}

func NewPipelineStmt(sqlorArgs ...interface{}) *PipelineStmt {
	return &PipelineStmt{sqlorArgs}
}

func RunPipeline(session *xorm.Session, stmts ...*PipelineStmt) (err error) {
	for _, st := range stmts {
		if _, err = session.Exec(st.sqlorArgs...); err != nil {
			return
		}
	}

	return
}
