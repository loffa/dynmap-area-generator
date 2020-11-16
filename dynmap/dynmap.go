package dynmap

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Marker struct {
	World  string  `yaml:"world"`
	Markup bool    `yaml:"markup"`
	Icon   string  `yaml:"icon"`
	Label  string  `yaml:"label"`
	X      float64 `yaml:"x"`
	Y      float64 `yaml:"y"`
	Z      float64 `yaml:"z"`
}

type Area struct {
	FillColor     int64     `yaml:"fillColor"`
	World         string    `yaml:"world"`
	Markup        bool      `yaml:"markup"`
	YTop          float64   `yaml:"ytop"`
	YBottom       float64   `yaml:"ybottom"`
	FillOpacity   float64   `yaml:"fillOPacity"`
	StrokeWeight  int       `yaml:"strokeWeight"`
	Label         string    `yaml:"label"`
	StrokeColor   int64     `yaml:"strokeColor"`
	StrokeOpacity float64   `yaml:"strokeOpacity"`
	X             []float64 `yaml:"x"`
	Z             []float64 `yaml:"z"`
}

type Set struct {
	Hide      bool                   `yaml:"hide"`
	Circles   map[string]interface{} `yaml:"circles"`
	DefIcon   string                 `yaml:"deficon"`
	Areas     map[string]*Area       `yaml:"areas"`
	Label     string                 `yaml:"label"`
	Markers   map[string]*Marker     `yaml:"markers"`
	Lines     map[string]interface{} `yaml:"lines"`
	LayerPrio int                    `yaml:"layerprio"`
}

type Dynmap struct {
	Icons      map[string]interface{} `yaml:"icons"`
	Sets       map[string]*Set        `yaml:"sets"`
	PlayerSets map[string]interface{} `yaml:"playersets"`
}

func FromFile(path string) (*Dynmap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	res := &Dynmap{}
	err = yaml.NewDecoder(f).Decode(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func WriteToFile(path string, dyn *Dynmap) error {
	tmpFilePath := path + ".tmp"
	f, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = yaml.NewEncoder(f).Encode(dyn)
	if err != nil {
		return err
	}

	err = os.Rename(tmpFilePath, path)
	if err != nil {
		return err
	}

	return nil
}
