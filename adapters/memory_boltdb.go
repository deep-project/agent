package adapters

import (
	"encoding/json"
	"fmt"

	"github.com/deep-project/agent/pkg/ability"
	"github.com/deep-project/agent/pkg/message"

	"go.etcd.io/bbolt"
)

type MemoryBoltDBAdapter struct {
	client *bbolt.DB
}

func NewMemoryBoltDBAdapter(client *bbolt.DB) *MemoryBoltDBAdapter {
	return &MemoryBoltDBAdapter{client: client}
}

func (m *MemoryBoltDBAdapter) GetMeta(sessionID string) (ability.Meta, error) {
	return ability.NewMeta(), nil
}

// message
func (m *MemoryBoltDBAdapter) getMessageBucketName(sessionID string) []byte {
	return []byte("messages-" + sessionID)
}

func (m *MemoryBoltDBAdapter) HasMessageSession(sessionID string) (exists bool, err error) {
	err = m.client.View(func(tx *bbolt.Tx) error {
		if bucket := tx.Bucket(m.getMessageBucketName(sessionID)); bucket != nil {
			exists = true
		}
		return nil
	})
	return
}

func (m *MemoryBoltDBAdapter) AddMessage(sessionID string, msg *message.Message) error {
	return m.client.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(m.getMessageBucketName(sessionID))
		if err != nil {
			return err
		}
		data, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		id, err := bucket.NextSequence() // 自增 key 作为数组索引
		if err != nil {
			return err
		}
		return bucket.Put(fmt.Appendf(nil, "%d", id), data)
	})
}

func (m *MemoryBoltDBAdapter) ListMessages(sessionID string, limit int) (res []message.Message, err error) {
	err = m.client.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(m.getMessageBucketName(sessionID))
		if bucket == nil {
			return fmt.Errorf("MemoryBoltDB bucket not found")
		}
		count := 0
		cursor := bucket.Cursor()
		for k, v := cursor.Last(); k != nil && count < limit; k, v = cursor.Prev() {
			var data message.Message
			if err := json.Unmarshal(v, &data); err == nil {
				res = append(res, data)
			}
			count++
		}
		m.reverseMessages(res) // 由于读取出来的数据是倒序的，所以这里要反转一下数组
		return nil
	})
	return
}

func (m *MemoryBoltDBAdapter) reverseMessages(s []message.Message) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
