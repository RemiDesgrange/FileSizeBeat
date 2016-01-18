// author: Rémi Desgrange
// date : 18/01/16

package beat

import (
  "time"
  "os"
  "path/filepath"
  "github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
)

type FileSizeBeat struct {
  period time.Duration
  paths []string
  config ConfigSettings
  events publisher.Client
  done chan struct{}
}

//TODO maybe need a isDir flag 
type DirSize struct {
  size int64
  nbFile int64
  nbFolder int64
}

func New() *FileSizeBeat {
  return &FileSizeBeat{}
}

func (fs *FileSizeBeat) Config(b *beat.Beat) error {
  err := cfgfile.Read(&fs.config, "")
  if err != nil {
    logp.Err("Error reading configuration file: %v", err)
		return err
  }

  if fs.config.Input.Period != nil {
    fs.period = time.Duration(*fs.config.Input.Period) * time.Second
  } else {
    fs.period = 10 * time.Second
  }
  logp.Debug("filesizebeat", "Period %v\n", fs.period)

  if fs.config.Input.Paths != nil {
    fs.paths = make([]string, len(*fs.config.Input.Paths))
		for _, path := range *fs.config.Input.Paths {
				fs.AddPath(path)
		}
    logp.Debug("filesizebeat", "Paths : %v\n", fs.paths)
	} else {
		logp.Critical("Error: no paths specified, cannot continue!")
		os.Exit(1)
	}
  return nil
}

// Setup performs boilerplate Beats setup
func (fs *FileSizeBeat) Setup(b *beat.Beat) error {
	fs.events = b.Events
	fs.done = make(chan struct{})
	return nil
}

func (fs *FileSizeBeat) Run(b *beat.Beat) error {
  ticker := time.NewTicker(fs.period)
  logp.Info("ticker %v", fs.period)
  defer ticker.Stop()
  for {
    select {
      case <- fs.done:
  			return nil
      case <- ticker.C: {
        //for _, path := range fs.paths {
          var ds DirSize
          filepath.Walk(fs.paths[0], ds.WalkFn)
          event := common.MapStr{
            "@timestamp": common.Time(time.Now()),
            "folder" : fs.paths[0],
            "nbFolder": ds.nbFolder,
            "nbFile": ds.nbFile,
            "size": ds.size,
          }
          logp.Info("%v", event)
          fs.events.PublishEvent(event)
        //}
      }
    }
  }
  return nil
}

func (fs *FileSizeBeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (fs *FileSizeBeat) Stop() {
	close(fs.done)
}


func (fs *FileSizeBeat) AddPath(target string) {
  fs.paths = append(fs.paths, target)
}


func (ds *DirSize) WalkFn(path string, info os.FileInfo, err error) error {
  if info.IsDir() {
    ds.nbFolder += 1
  } else {
    ds.nbFile += 1
    ds.size += info.Size()
  }
  return nil
}
