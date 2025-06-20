package go_ora

import (
	"github.com/overshinin/go-ora/network"
)

type RefCursor struct {
	defaultStmt
	len        uint8
	MaxRowSize int
	parent     *defaultStmt
	// ID         int
	// scnFromExe []int
	// connection *Connection
	// noOfRowsToFetch int
	// hasMoreRows bool
}

func (cursor *RefCursor) load(session *network.Session) error {
	// initialize ref cursor object
	cursor.text = ""
	cursor._hasLONG = false
	cursor._hasBLOB = false
	cursor._hasReturnClause = false
	cursor.disableCompression = true
	cursor.arrayBindCount = 1
	cursor.scnFromExe = make([]int, 2)
	cursor.stmtType = SELECT
	var err error
	cursor.len, err = session.GetByte()
	if err != nil {
		return err
	}
	cursor.MaxRowSize, err = session.GetInt(4, true, true)
	if err != nil {
		return err
	}
	columnCount, err := session.GetInt(4, true, true)
	if err != nil {
		return err
	}
	if columnCount > 0 {
		cursor.columns = make([]ParameterInfo, columnCount)
		_, err = session.GetByte()
		if err != nil {
			return err
		}
		for x := 0; x < len(cursor.columns); x++ {
			err = cursor.columns[x].load(session)
			if err != nil {
				return err
			}
		}
		//for _, col := range cursor.Cols {
		//	err = col.load(session)
		//	fmt.Println(col)
		//	if err != nil {
		//		return err
		//	}
		//}
	}
	_, err = session.GetDlc()
	if err != nil {
		return err
	}
	if session.TTCVersion >= 3 {
		_, err = session.GetInt(4, true, true)
		if err != nil {
			return err
		}
		_, err = session.GetInt(4, true, true)
		if err != nil {
			return err
		}
	}
	if session.TTCVersion >= 4 {
		_, err = session.GetInt(4, true, true)
		if err != nil {
			return err
		}
		_, err = session.GetInt(4, true, true)
		if err != nil {
			return err
		}
	}
	if session.TTCVersion >= 5 {
		_, err = session.GetDlc()
		if err != nil {
			return err
		}
	}
	cursor.cursorID, err = session.GetInt(4, true, true)
	if err != nil {
		return err
	}
	_, err = session.GetInt(2, true, true)
	if err != nil {
		return err
	}
	return nil
}

func (cursor *RefCursor) getExeOptions() int {
	return 0x8040
}

func (cursor *RefCursor) Query() (*DataSet, error) {
	cursor.connection.connOption.Tracer.Printf("Query RefCursor: %d", cursor.cursorID)
	cursor._noOfRowsToFetch = cursor.connection.connOption.PrefetchRows
	cursor._hasMoreRows = true
	if len(cursor.parent.scnFromExe) > 0 {
		copy(cursor.scnFromExe, cursor.parent.scnFromExe)
	}
	session := cursor.connection.session
	session.ResetBuffer()
	err := cursor.write()
	if err != nil {
		return nil, err
	}
	dataSet := new(DataSet)
	err = cursor.read(dataSet)
	if err != nil {
		return nil, err
	}
	return dataSet, nil
}

func (cursor *RefCursor) write() error {
	err := cursor.basicWrite(cursor.getExeOptions(), false, false)
	if err != nil {
		return err
	}
	return cursor.connection.session.Write()
}

func (cursor *RefCursor) Close() error {
	return nil
}

//func (cursor *RefCursor) Exec(args []driver.Value) (driver.Result, error) {
//	return nil, nil
//}
//
//func (cursor *RefCursor) NumInput() int {
//	return -1;
//}
//func (cursor *RefCursor) readQ() (*DataSet, error) {
//	dataSet := new(DataSet)
//	return dataSet, nil
//}
