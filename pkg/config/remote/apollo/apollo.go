package apollo

import (
	"bytes"
	"strings"

	"github.com/NetEase-Media/ngo/pkg/config/remote"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/philchia/agollo/v4"
	"go.uber.org/zap"
)

const (
	propertiesSuffix = ".properties"
)

var (
	replacer = strings.NewReplacer(
		"\n", "\\n ",
		"\r", "\\r",
		"\t", "\\t",
	)
)

type apolloDataSource struct {
	nameSpaceNames []string
	client         agollo.Client
	changed        chan struct{}
}

// NewDataSource creates an apolloDataSource
func NewDataSource(conf *agollo.Conf, nameSpaceNames []string) remote.DataSource {
	client := agollo.NewClient(conf, agollo.WithLogger(&agolloLogger{
		logger: log.DefaultLogger().(*log.NgoLogger).WithOptions(zap.AddCallerSkip(1)),
	}))
	ap := &apolloDataSource{
		nameSpaceNames: nameSpaceNames,
		client:         client,
		changed:        make(chan struct{}, 1),
	}
	_ = ap.client.Start()
	ap.client.OnUpdate(
		func(event *agollo.ChangeEvent) {
			ap.changed <- struct{}{}
		})
	return ap
}

// ReadConfig reads config content from apollo
func (ap *apolloDataSource) ReadConfig() ([]byte, error) {
	var buf bytes.Buffer
	for i := range ap.nameSpaceNames {
		op := agollo.WithNamespace(ap.nameSpaceNames[i])
		if strings.HasSuffix(ap.nameSpaceNames[i], propertiesSuffix) {
			keys := ap.client.GetAllKeys(op)
			for j := range keys {
				v := ap.client.GetString(keys[j], op)
				buf.WriteString(keys[j])
				buf.WriteString("=")
				buf.WriteString(replacer.Replace(v))
				buf.WriteByte('\n')
			}
		} else {
			buf.WriteString(ap.client.GetContent(op))
			buf.WriteByte('\n')
		}
	}
	return buf.Bytes(), nil
}

// IsConfigChanged returns a chanel for notification when the config changed
func (ap *apolloDataSource) IsConfigChanged() <-chan struct{} {
	return ap.changed
}

// Close stops watching the config changed
func (ap *apolloDataSource) Close() error {
	_ = ap.client.Stop()
	close(ap.changed)
	return nil
}

type agolloLogger struct {
	logger log.Logger
}

// Infof ...
func (l *agolloLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Errorf ...
func (l *agolloLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}
