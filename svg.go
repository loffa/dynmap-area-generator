package main

import (
	"github.com/spf13/viper"
	"golang.org/x/net/html"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Coordinate represents the final map coordinate
type Coordinate struct {
	X float64
	Y float64
}

// Path represent the resulting path from the SVG
type Path struct {
	Coordinates []*Coordinate
	Name        string
}

var pathCommands = regexp.MustCompile("[MmLlHhVvZz]")

func getPaths(filePath string) ([]*Path, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(f)
	if err != nil {
		return nil, err
	}

	res := make([]*Path, 0, 50)

	var crawler func(n *html.Node)

	currentGroup := ""
	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "g" {
			// n is a group
			for _, attr := range n.Attr {
				if attr.Key == "id" {
					currentGroup = attr.Val
					break
				}
			}
		} else if n.Type == html.ElementNode && n.Data == "path" {
			// n is a path
			if currentGroup == "areas" {
				pathInfo, err := getPathInfo(n)
				if err != nil {
					log.Println("Could not read path, skipping:", err)
					return
				}
				res = append(res, pathInfo)
				return
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}
	crawler(doc)

	return res, nil
}

func getPathInfo(n *html.Node) (*Path, error) {
	coordinates := make([]*Coordinate, 0, 10)
	name := ""

	for _, attr := range n.Attr {
		if attr.Key == "d" {
			fields := strings.Fields(attr.Val)

			lastX := 0.0
			lastY := 0.0
			lastCmd := ""
			for _, f := range fields {
				if pathCommands.MatchString(f) {
					if f == "z" || f == "Z" {
						break
					}
					lastCmd = f
				} else {
					x := lastX
					y := lastY
					switch lastCmd {
					case "V":
						// Vertical line (absolute)
						y = getSingleCoordinate(f)
					case "v":
						// Vertical line (relative)
						relY := getSingleCoordinate(f)
						y += relY
					case "H":
						// Horizontal line (absolute)
						x = getSingleCoordinate(f)
					case "h":
						// Horizontal line (relative)
						relX := getSingleCoordinate(f)
						x += relX
					case "L":
						// Straight Line (absolute)
						x, y = getDualCoordinate(f)
					case "l":
						// Straight Line (relative)
						relX, relY := getDualCoordinate(f)
						x += relX
						y += relY
					case "M":
						// MoveTo (absolute)
						x, y = getDualCoordinate(f)
					case "m":
						// MoveTo (relative)
						relX, relY := getDualCoordinate(f)
						x += relX
						y += relY
					}

					coordinates = append(coordinates, calculatedCoordinate(x, y))
					lastX = x
					lastY = y

				}
			}
		} else if attr.Key == "id" {
			name = attr.Val
		}
	}

	return &Path{
		Coordinates: coordinates,
		Name:        name,
	}, nil
}

func calculatedCoordinate(x, y float64) *Coordinate {
	offsetX := viper.GetFloat64("offset-x")
	offsetY := viper.GetFloat64("offset-y")
	scaleX := viper.GetFloat64("scale-x")
	scaleY := viper.GetFloat64("scale-y")

	return &Coordinate{
		X: math.Floor((x + offsetX) * scaleX),
		Y: math.Floor((y + offsetY) * scaleY),
	}
}

func getDualCoordinate(s string) (x float64, y float64) {
	if strings.Contains(s, ",") {
		xy := strings.Split(s, ",")
		x, _ = strconv.ParseFloat(xy[0], 64)
		y, _ = strconv.ParseFloat(xy[1], 64)
	}

	return x, y
}

func getSingleCoordinate(s string) float64 {
	c, _ := strconv.ParseFloat(s, 64)
	return c
}
