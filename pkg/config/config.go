package config

import (
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/NetEase-Media/ngo/pkg/util/file"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type DecoderConfigOption = viper.DecoderConfigOption

func New(configAddr string, watch bool) (*Configuration, error) {
	urlObj, err := url.Parse(configAddr)
	if err != nil {
		return nil, err
	}

	var scheme = urlObj.Scheme
	if scheme == "" {
		scheme = "file"
	}

	v := viper.New()
	var children []*viper.Viper
	switch scheme {
	case "file":
		absolutePath, err := filepath.Abs(configAddr)
		if err != nil {
			return nil, err
		}
		v.SetConfigFile(absolutePath)
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
		imports := v.GetStringSlice("configImports")
		if len(imports) > 0 {
			children = make([]*viper.Viper, 0, len(imports))
			dir := file.CheckAndGetParentDir(absolutePath)
			for _, f := range imports {
				tmp := viper.New()
				tmp.SetConfigFile(dir + string(os.PathSeparator) + f)
				if err := tmp.ReadInConfig(); err != nil {
					continue
				}
				children = append(children, tmp)
				if err := v.MergeConfigMap(tmp.AllSettings()); err != nil {
					log.Errorf("sub config merge error: %s", err)
				}
			}

		}
	}
	cfg := &Configuration{
		Viper:     v,
		source:    scheme,
		children:  children,
		watch:     watch,
		onChanges: make([]func(*Configuration), 0),
	}
	cfg.watchConfig()
	return cfg, nil
}

func Empty() *Configuration {
	return &Configuration{Viper: viper.New()}
}

type Configuration struct {
	*viper.Viper
	source    string
	children  []*viper.Viper
	watch     bool
	onChanges []func(*Configuration)
}

func (c *Configuration) watchConfig() error {
	if c.watch {
		switch c.source {
		case "file":
			c.WatchConfig()
			c.OnConfigChange(func(_ fsnotify.Event) {
				for _, change := range c.onChanges {
					change(c)
				}
			})
			for i := range c.children {
				child := c.children[i]
				child.WatchConfig()
				child.OnConfigChange(func(_ fsnotify.Event) {
					if err := c.MergeConfigMap(child.AllSettings()); err != nil {
						log.Errorf("sub config merge error: %s", err)
						return
					}
					for _, change := range c.onChanges {
						change(c)
					}
				})
			}
		}
	}
	return nil
}

func (c *Configuration) Sub(key string) *Configuration {
	sub := c.Viper.Sub(key)
	if sub == nil {
		return nil
	}
	return &Configuration{Viper: sub}
}

func (c *Configuration) SubSlice(key string) []*Configuration {
	s := c.Get(key)
	if s == nil {
		return nil
	}
	if reflect.TypeOf(s).Kind() == reflect.Slice {
		v := reflect.ValueOf(s)
		r := make([]*Configuration, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			k := key + "." + strconv.Itoa(i)
			sub := Empty()
			sub.MergeConfigMap(c.GetStringMap(k))
			r = append(r, sub)
		}
		return r
	}
	return nil
}

func (c *Configuration) OnChange(fn func(*Configuration)) {
	c.onChanges = append(c.onChanges, fn)
}
