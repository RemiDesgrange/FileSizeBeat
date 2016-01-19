// author: Rémi Desgrange
// date : 18/01/16

package beat

import (
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
	"os"
	"path/filepath"
	"time"
)

type FileSizeBeat struct {
	period time.Duration
	paths  []Path
	config ConfigSettings
	events publisher.Client
	done   chan struct{}
}

type Path struct{
	path string
	isDir bool
}

type DirSize struct {
	size     int64
	nbFile   int64
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
		//fs.paths = make([]Path, len(*fs.config.Input.Paths))
		for _, path := range *fs.config.Input.Paths {
			err := fs.AddPath(path)
			if err != nil {
				logp.Critical("Error: %v", err)
				os.Exit(1)
			}
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

	for _, onepath := range fs.paths {

		go func(onepath Path){
			ticker := time.NewTicker(fs.period)
			defer ticker.Stop()
			for {
				select {
				case <-fs.done:{
					logp.Debug("filesizebeat", "done in %v path ", onepath.path)
				}
				case <-ticker.C: {
					ds := DirSize{}
					err := filepath.Walk(onepath.path, func(path string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}
						if info.IsDir() {
							ds.nbFolder += 1
						} else {
							ds.nbFile += 1
							ds.size += info.Size()
						}
						return nil
					})
					if (err == nil) && (ds != DirSize{}) {
						event := common.MapStr{
							"@timestamp": common.Time(time.Now()),
							"type": "filesizebeat",
							"path": onepath.path,
							"nbFolder": ds.nbFolder,
							"nbFile": ds.nbFile,
							"isDir": onepath.isDir,
							"size": ds.size,
						}
						fs.events.PublishEvent(event)
					} else {
						logp.Err("Error while filepathWalk %v, Path %v, (DirSize %v)\n", err, onepath, ds)
					}
				}
				}
			}
		}(onepath)
	}
	<- fs.done
	return nil
}

func (fs *FileSizeBeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (fs *FileSizeBeat) Stop() {
	close(fs.done)
}

func (fs *FileSizeBeat) AddPath(target string) error {
	newPath := Path{target, false}
	//check that the target is a dir or a regular file
	fdir, err := os.Open(target)
	if err != nil {
			logp.Err("%v\n", err)
			return err
	}
	defer fdir.Close()

	finfo, err := fdir.Stat()
	if err != nil {
		logp.Err("%v\n", err)
		return err
	}

	mode := finfo.Mode()
	if mode.IsDir() {
		newPath.isDir = true
	}
	fs.paths = append(fs.paths, newPath)
	logp.Debug("filesizebeat", "Append %v to the folder to monitor", newPath)
	return nil
}
