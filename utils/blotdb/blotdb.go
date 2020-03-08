package blotdb

import (
	"bytes"
	"github.com/astaxie/beego/logs"
	"github.com/boltdb/bolt"
	"log"
)

var (
	//DbUrl = "C:\\iceblue\\go\\src\\ds-yibasuo\\web\\ds-yibasuo.db"
	DbUrl = "ds-yibasuo.db"
	Db    *BoltManager
)

func BlotInit() {
	logs.Info("Init blotdb")
	Db, _ = NewBoltManager(DbUrl, []string{"host", "cluster", "config", "user", "devops", "init"})
}

//BlotDB的管理类
type BoltManager struct {
	db *bolt.DB
}

//创建库管理,并生成Bucket
func NewBoltManager(dbPath string, bucket []string) (*BoltManager, error) {
	//  bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second,ReadOnly: true})
	db, err := bolt.Open(dbPath, 0644, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		for _, v := range bucket {
			_, err := tx.CreateBucketIfNotExists([]byte(v))
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &BoltManager{db}, nil
}

//关库数据库
func (m *BoltManager) Close() error {
	return m.db.Close()
}

//移除Bucket
func (m *BoltManager) RemoveBucket(bucketName string) (err error) {
	err = m.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucketName))
	})
	return err
}

//组Bucket增加值
//func (m *BoltManager) Add(bucketName string, val []byte) (err error) {
//	err = m.db.Update(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte(bucketName))
//		id, _ = b.NextSequence() //sequence uint64
//		bBuf := fmt.Sprintf("%d", id)
//		return b.Put([]byte(bBuf), val)
//	})
//	return
//}

func (m *BoltManager) Add(bucketName string, id, val []byte) (err error) {
	err = m.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		return b.Put(id, val)
	})
	return
}

// 插入指定id来代替更新功能
func (m *BoltManager) Update(bucketName string, id, val []byte) error {
	err := m.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		return b.Put(id, val)
	})
	return err
}

//遍历bucket 返回 map
func (m *BoltManager) SelectAll2Map(bucketName string) (res map[string]string, err error) {
	res = make(map[string]string)
	_ = m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		_ = b.ForEach(func(k, v []byte) error {
			//log.Printf("key=%s, value=%s\n", string(k), v)
			res[string(k)] = string(v)
			return nil
		})
		return nil
	})
	return res, nil
}

//遍历bucket 返回 slice
func (m *BoltManager) SelectAll2List(bucketName string) (res []map[string]string, err error) {
	tmp := make(map[string]string)
	_ = m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		_ = b.ForEach(func(k, v []byte) error {
			//log.Printf("key=%s, value=%s\n", string(k), v)
			tmp[string(k)] = string(v)
			res = append(res, tmp)
			return nil
		})
		return nil
	})
	return res, nil
}

// 查询某个topic，且只返回value
func (m *BoltManager) SelectValues(bucketName string) (res [][]byte, err error) {
	_ = m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		_ = b.ForEach(func(k, v []byte) error {
			res = append(res, v)
			return nil
		})
		return nil
	})
	return res, nil
}

//移除指定Bucket中指定ID
func (m *BoltManager) RemoveID(bucketName string, id []byte) error {
	err := m.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		return b.Delete(id)
	})
	return err
}

//移除指定Bucket中指定Val
func (m *BoltManager) RemoveVal(bucketName string, val []byte) (err error) {
	var arrID []string
	arrID = make([]string, 1)
	err = m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			log.Printf("key=%s, value=%s\n", k, string(v))
			if bytes.Compare(v, val) == 0 {
				arrID = append(arrID, string(k))
			}
		}
		return nil
	})

	err = m.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		for _, v := range arrID {
			_ = b.Delete([]byte(v))
			log.Println("Del k:", v)
		}
		return nil
	})

	return err
}

//查找指定值
//func (m *BoltManager) SelectVal(bucketName string, val []byte) (arr []string, err error) {
//  arr = make([]string, 0, 1)
//  err = m.db.View(func(tx *bolt.Tx) error {
//    c := tx.Bucket([]byte(bucketName)).Cursor()
//    for k, v := c.First(); k != nil; k, v = c.Next() {
//      if bytes.Compare(v, val) == 0 {
//        arr = append(arr, string(k))
//      }
//    }
//    return nil
//  })
//  return arr, err
//}

func (m *BoltManager) SelectVal(bucketName string, id []byte) (arr []string, err error) {
	arr = make([]string, 0, 1)
	err = m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		get := b.Get(id)
		if string(get) != "" {
			arr = append(arr, string(get))
		}
		return nil
	})
	return arr, err
}

//在事务中，移除指定Bucket中指定Val
func (m *BoltManager) RemoveValTransaction(bucketName string, val []byte) (count int, err error) {
	arrID, err1 := m.SelectVal(bucketName, val)
	if err1 != nil {
		return 0, err1
	}
	count = len(arrID)
	if count == 0 {
		return count, nil
	}

	tx, err1 := m.db.Begin(true)
	if err1 != nil {
		return count, err1
	}
	b := tx.Bucket([]byte(bucketName))
	for _, v := range arrID {
		if err = b.Delete([]byte(v)); err != nil {
			log.Printf("删除ID(%s)失败! 执行回滚. err:%s \r\n", v, err)
			_ = tx.Rollback()
			return
		}
		log.Println("删除ID(", v, ")成功!")
	}
	err = tx.Commit()
	return
}
