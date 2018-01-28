package configaction

import (
	"io/ioutil"
	"regexp"

	"github.com/windler/godepg/action"
	yaml "gopkg.in/yaml.v2"
)

//Config represents the config file
type Config struct {
	Language   string
	Filter     []string
	Depth      int
	Output     string
	Edgestyle  map[string]map[string]string
	Nodestyle  map[string]string
	Graphstyle map[string]string

	//PHP
	Project string
	Exclude []string
	StopAt  []string

	//GO
	Package        string
	NoGoPackages   bool
	MyPackagesOnly bool
}

//CreateContext creates a Config
func CreateContext(file string, context action.Context) *Config {
	cfg := &Config{}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal(data, cfg)

	context.SetStringSliceFlag("f", cfg.Filter)
	context.SetStringSliceFlag("s", cfg.StopAt)
	context.SetStringSliceFlag("e", cfg.Exclude)

	context.SetBoolFlag("n", cfg.NoGoPackages)
	context.SetBoolFlag("m", cfg.MyPackagesOnly)

	re := regexp.MustCompile("(.*)\\/(.+)")
	projectRoot := re.FindStringSubmatch(file)[1]
	context.SetStringFlag("p", projectRoot)
	if cfg.Project != "" {
		context.SetStringFlag("p", projectRoot+"/"+cfg.Project)
	}
	if cfg.Package != "" {
		context.SetStringFlag("p", cfg.Package)
	}

	context.SetIntFlag("d", cfg.Depth)
	if cfg.Depth == 0 {
		context.SetIntFlag("d", -1)
	}

	if cfg.Output != "" {
		context.SetStringFlag("o", cfg.Output)
	}

	return cfg
}
