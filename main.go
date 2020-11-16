package main

import (
	"fmt"
	"github.com/loffa/dynmap-area-generator/dynmap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "dynmap-area-generator",
	Short: "dynmap-area-generator creates areas for Minecraft Dynmaps",
	Run: func(cmd *cobra.Command, args []string) {
		handleImage()
	},
}

func configure() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func init() {
	cobra.OnInitialize(configure)

	rootCmd.Flags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.Flags().StringP("image", "i", "", "SVG image for zone data")
	rootCmd.Flags().StringP("map-layer", "l", "", "The layer in the Dynmap config to update/create")
	rootCmd.Flags().StringP("dynmap-config", "d", "", "The Dynmap config file to read from and write to")
	rootCmd.Flags().Float64("offset-x", 0, "Offset in X from SVG image")
	rootCmd.Flags().Float64("offset-y", 0, "Offset in Y from SVG image")
	rootCmd.Flags().Float64("scale-x", 0, "Scale the resulting coordinate in X from SVG")
	rootCmd.Flags().Float64("scale-y", 0, "Scale the resulting coordinates in Y from SVG")

	if cfgFile == "" {
		_ = rootCmd.MarkFlagRequired("image")
		_ = rootCmd.MarkFlagRequired("map-layer")
		_ = rootCmd.MarkFlagRequired("dynmap-config")
	}

	_ = viper.BindPFlag("image", rootCmd.Flags().Lookup("image"))
	_ = viper.BindPFlag("map-layer", rootCmd.Flags().Lookup("map-layer"))
	_ = viper.BindPFlag("dynmap-config", rootCmd.Flags().Lookup("dynmap-config"))
	_ = viper.BindPFlag("offset-x", rootCmd.Flags().Lookup("offset-x"))
	_ = viper.BindPFlag("offset-y", rootCmd.Flags().Lookup("offset-y"))
	_ = viper.BindPFlag("scale-x", rootCmd.Flags().Lookup("scale-x"))
	_ = viper.BindPFlag("scale-y", rootCmd.Flags().Lookup("scale-y"))

	viper.SetDefault("offset_x", 0)
	viper.SetDefault("offset_y", 0)
	viper.SetDefault("scale_x", 0.0)
	viper.SetDefault("scale_y", 0.0)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func handleImage() {
	dyn, err := dynmap.FromFile(viper.GetString("dynmap-config"))
	if err != nil {
		log.Println("Could not read dynmap file:", err)
		return
	}

	layer := viper.GetString("map-layer")
	set, ok := dyn.Sets[layer]
	if !ok {
		set = &dynmap.Set{
			Hide:      true,
			Circles:   nil,
			DefIcon:   "default",
			Areas:     make(map[string]*dynmap.Area),
			Label:     layer,
			Markers:   make(map[string]*dynmap.Marker),
			Lines:     nil,
			LayerPrio: 0,
		}
		dyn.Sets[layer] = set
	}
	newAreas := make(map[string]*dynmap.Area)
	paths, err := getPaths(viper.GetString("image"))
	if err != nil {
		log.Println("Could not get path from SVG:", err)
		return
	}

	for i, path := range paths {
		area := &dynmap.Area{
			FillColor:     0xFFFFFF,
			World:         "world",
			Markup:        false,
			YTop:          64,
			YBottom:       64,
			FillOpacity:   0.35,
			StrokeWeight:  3,
			Label:         path.Name,
			StrokeColor:   0xFFFFFF,
			StrokeOpacity: 0.8,
			X:             make([]float64, 0, len(path.Coordinates)),
			Z:             make([]float64, 0, len(path.Coordinates)),
		}

		for _, c := range path.Coordinates {
			area.X = append(area.X, c.X)
			area.Z = append(area.Z, c.Y)
		}

		areaName := fmt.Sprintf("area_%d", i)
		newAreas[areaName] = area
	}

	set.Areas = newAreas
	err = dynmap.WriteToFile(viper.GetString("dynmap-config"), dyn)
	if err != nil {
		log.Println("Could not write config after editing:", err)
	}
}
