package service

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zhanbu/internal/model"
)

type chatStarterRecordRepo struct {
	nextID  uint
	records []*model.DivinationRecord
}

func (r *chatStarterRecordRepo) Create(record *model.DivinationRecord) error {
	r.nextID++
	record.ID = r.nextID
	r.records = append(r.records, record)
	return nil
}

func TestChatModeStarterStartsBaZiRecord(t *testing.T) {
	repo := &chatStarterRecordRepo{}
	starter := NewChatModeStarter(repo, nil, nil, NewBaZiService(), nil, nil)

	record, err := starter.Start(7, "bazi", "1990-05-12 08:30 女，看看事业")

	require.NoError(t, err)
	require.NotNil(t, record)
	assert.Equal(t, uint(1), record.ID)
	assert.Equal(t, uint(7), record.UserID)
	assert.Equal(t, "bazi", record.Type)
	assert.Contains(t, record.Question, "看看事业")

	var result model.BaZiResult
	require.NoError(t, json.Unmarshal([]byte(record.Result), &result))
	assert.Equal(t, "1990-05-12 08:30", result.Birth.Solar)
	assert.Len(t, repo.records, 1)
}

func TestChatModeStarterRejectsBaZiWithoutBirthTime(t *testing.T) {
	starter := NewChatModeStarter(&chatStarterRecordRepo{}, nil, nil, NewBaZiService(), nil, nil)

	record, err := starter.Start(7, "bazi", "看看我的事业")

	require.Error(t, err)
	assert.Nil(t, record)
	assert.Contains(t, err.Error(), "出生日期和时间")
}
