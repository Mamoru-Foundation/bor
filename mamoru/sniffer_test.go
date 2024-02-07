package mamoru

import (
	"fmt"
	"os"
	"testing"

	"github.com/Mamoru-Foundation/mamoru-sniffer-go/mamoru_sniffer"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/internal/testlog"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
)

type StatusProgressMock struct {
	status ethereum.SyncProgress
}

func CreateProgress(currentBlock, highestBlock uint64) *StatusProgressMock {
	return &StatusProgressMock{
		status: ethereum.SyncProgress{CurrentBlock: currentBlock, HighestBlock: highestBlock},
	}
}

func (s *StatusProgressMock) Progress() ethereum.SyncProgress {
	return s.status
}

func TestSniffer_CheckRequirements(t *testing.T) {
	tests := []struct {
		name   string
		status statusProgress
		want   bool
	}{
		{
			name:   "FALSE - currentBlock == 0 && highestBlock == 0",
			status: CreateProgress(0, 0),
			want:   false,
		},
		{
			name:   "FALSE - currentBlock < highestBlock",
			status: CreateProgress(1, 100),
			want:   false,
		},
		{
			name:   "TRUE - currentBlock == highestBlock",
			status: CreateProgress(10, 10),
			want:   true,
		},
		{
			name:   "TRUE - currentBlock > highestBlock",
			status: CreateProgress(200, 10),
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sniffer{
				status: tt.status,
			}
			if got := s.checkSynced(); got != tt.want {
				t.Errorf("CheckRequirements() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSniffer_isSnifferEnable(t *testing.T) {
	t.Run("TRUE env is set 1", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "1")
		defer unsetEnvSnifferEnable()
		got := isSnifferEnable()
		assert.True(t, got)
	})
	t.Run("TRUE env is set true", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "true")
		defer unsetEnvSnifferEnable()
		got := isSnifferEnable()
		assert.True(t, got)
	})
	t.Run("FALSE env is set 0", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "0")
		defer unsetEnvSnifferEnable()
		got := isSnifferEnable()
		assert.False(t, got)
	})
	t.Run("FALSE env is set 0", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "false")
		defer unsetEnvSnifferEnable()
		got := isSnifferEnable()
		assert.False(t, got)
	})
	t.Run("FALSE env is not set", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "")
		defer unsetEnvSnifferEnable()
		got := isSnifferEnable()
		assert.False(t, got)
	})
}

func unsetEnvSnifferEnable() {
	_ = os.Unsetenv("MAMORU_SNIFFER_ENABLE")
}

func TestSniffer_connect(t *testing.T) {
	t.Run("TRUE ", func(t *testing.T) {
		SnifferConnectFunc = func() (*mamoru_sniffer.Sniffer, error) { return nil, nil }
		got := connect()
		assert.True(t, got)
	})
	t.Run("FALSE connect have error", func(t *testing.T) {
		SnifferConnectFunc = func() (*mamoru_sniffer.Sniffer, error) { return nil, fmt.Errorf("Some err") }
		got := connect()
		assert.False(t, got)
	})
}

func TestSniffer_CheckRequirements1(t *testing.T) {
	t.Run("TRUE ", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "true")
		defer unsetEnvSnifferEnable()
		SnifferConnectFunc = func() (*mamoru_sniffer.Sniffer, error) { return nil, nil }
		s := &Sniffer{
			status: CreateProgress(100, 5),
			synced: true,
		}
		assert.True(t, s.CheckRequirements())
	})
	t.Run("FALSE chain not sync ", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "true")
		defer unsetEnvSnifferEnable()
		SnifferConnectFunc = func() (*mamoru_sniffer.Sniffer, error) { return nil, nil }
		s := &Sniffer{
			status: CreateProgress(5, 100),
			synced: true,
		}
		assert.False(t, s.CheckRequirements())
	})
	t.Run("FALSE connect error", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "true")
		defer unsetEnvSnifferEnable()
		SnifferConnectFunc = func() (*mamoru_sniffer.Sniffer, error) { return nil, fmt.Errorf("Some err") }
		s := &Sniffer{
			status: CreateProgress(100, 5),
			synced: true,
		}
		assert.False(t, s.CheckRequirements())
	})
	t.Run("FALSE env not set", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "0")
		defer unsetEnvSnifferEnable()
		SnifferConnectFunc = func() (*mamoru_sniffer.Sniffer, error) { return nil, nil }
		s := &Sniffer{
			status: CreateProgress(100, 5),
			synced: true,
		}
		assert.False(t, s.CheckRequirements())
	})
}

func Test_getDeltaBlocks(t *testing.T) {
	t.Run("Success return delta from env ", func(t *testing.T) {
		want := int64(100)
		_ = os.Setenv("MAMORU_SNIFFER_DELTA", fmt.Sprintf("%d", want))
		defer os.Unsetenv("MAMORU_SNIFFER_DELTA")
		got := getDelta()
		assert.Equal(t, want, got)
	})
	t.Run("Success return delta from env ", func(t *testing.T) {
		want := DefaultDelta
		got := getDelta()
		assert.Equal(t, want, got)
		defer os.Unsetenv("MAMORU_SNIFFER_DELTA")
	})
}

func TestLogger(t *testing.T) {
	ctx := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	kvs := mapToInterfaceSlice(ctx)
	logger := testlog.Logger(t, log.LvlInfo)
	logger.Info("Message", kvs...)
}
