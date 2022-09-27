package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNgoLogger_Info(t *testing.T) {
	opt := NewDefaultOptions()
	logger, err := New(opt)
	assert.NoError(t, err)
	logger.Info("hello world")
}

func TestNgoLogger_InfoWithPackageLevel(t *testing.T) {
	opt := NewDefaultOptions()
	opt.PackageLevel = map[string]string{
		"github.com/NetEase-Media/ngo/pkg/log": "error",
	}
	logger, err := New(opt)
	assert.NoError(t, err)
	logger.Info("hello world")
}

func TestNgoLogger_Infof(t *testing.T) {
	opt := NewDefaultOptions()
	logger, err := New(opt)
	assert.NoError(t, err)
	logger.Infof("hello world: %s", "ngo")
}

func TestNgoLogger_Infol(t *testing.T) {
	opt := NewDefaultOptions()
	logger, err := New(opt)
	assert.NoError(t, err)
	logger.Infol("hello world", String("app", "ngo"))
}

func TestNgoLogger_InfoWithField(t *testing.T) {
	opt := NewDefaultOptions()
	//opt.Format = formatJSON
	logger, err := New(opt)
	assert.NoError(t, err)
	logger.WithField("app", "ngo").Info("hello world")
}

func TestNgoLogger_InfoWithFields(t *testing.T) {
	opt := NewDefaultOptions()
	//opt.Format = formatJSON
	logger, err := New(opt)
	assert.NoError(t, err)
	logger.WithFields("app", "ngo", "version", "v1.0.0").Info("hello world")
}

func TestNgoLogger_Infow(t *testing.T) {
	opt := NewDefaultOptions()
	//opt.Format = formatJSON
	logger, err := New(opt)
	assert.NoError(t, err)
	logger.Infow("hello world", "app", "ngo")
}

func TestNgoLogger_Error(t *testing.T) {
	opt := NewDefaultOptions()
	logger, err := New(opt)
	assert.NoError(t, err)
	logger.Error("hello world")
}

func TestNgoLogger_Errorf(t *testing.T) {
	opt := NewDefaultOptions()
	logger, err := New(opt)
	assert.NoError(t, err)
	logger.Errorf("hello world: %s", "ngo")
}

func TestNgoLogger_Errorl(t *testing.T) {
	opt := NewDefaultOptions()
	logger, err := New(opt)
	assert.NoError(t, err)
	logger.Errorl("hello world", String("app", "ngo"))
}

func TestNgoLogger_Errorw(t *testing.T) {
	opt := NewDefaultOptions()
	opt.Format = formatJSON
	logger, err := New(opt)
	assert.NoError(t, err)
	logger.Errorw("hello world", "app", "ngo")
}

func BenchmarkNgoLogger_Info(b *testing.B) {
	opt := NewDefaultOptions()
	logger, _ := New(opt)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("hello world")
	}
}

func BenchmarkNgoLogger_Infof(b *testing.B) {
	opt := NewDefaultOptions()
	logger, _ := New(opt)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("hello world %s:%s", "app", "ngo")
	}
}

func BenchmarkNgoLogger_Infow(b *testing.B) {
	opt := NewDefaultOptions()
	logger, _ := New(opt)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infow("hello world", "app", "ngo")
	}
}

func BenchmarkNgoLogger_Infol(b *testing.B) {
	opt := NewDefaultOptions()
	logger, _ := New(opt)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infol("hello world", String("app", "ngo"))
	}
}
