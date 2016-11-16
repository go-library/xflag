package xflag

import (
	"flag"
	"testing"
	"time"
)

type Custom struct {
	index int
}

func (c Custom) Set(string) error {
	return nil
}

func (c Custom) String() string {
	return "......."
}

type Opt struct {
	Verbose  bool          `xflag-default:"true"`
	Duration time.Duration `xflag-default:"1m30s"`
	Float    float64       `xflag-default:"1.23"`
	MaxConn  int           `xflag-name:"max-conn" xflag-default:"100" xflag-usage:"usage...max-conn"`
	IPAddr   string        `xflag-name:"ipaddr" xflag-default:"0.0.0.0:80" xflag-usage:"usage...ipaddr"`
	Custom   Custom
}

func TestFlag(t *testing.T) {

	opt := &Opt{}

	fs, err := NewFlagSetFromStruct(opt)
	if err != nil {
		t.Fatal(err)
	}

	err = fs.Parse([]string{
		"-ipaddr", "0.0.0.0:8080",
		"-max-conn", "500",
	})

	if err != nil {
		t.Error(err)
	}

	fs.VisitAll(func(f *flag.Flag) {
		t.Logf("- %+v\n", f)
	})

}
