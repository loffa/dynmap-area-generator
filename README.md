# Dynmap Area Generator
This tool will generate areas in a Dynmap config file for Minecraft. It uses a special SVG image as base for the areas.
All paths from the SVG in the layer with id `areas` is added as areas on the map in the specified layer. The id of the
Path is used as the area-name on the map.

## Usage
```
Usage:
  dynmap-area-generator [flags]

Flags:
      --config string          config file
  -d, --dynmap-config string   The Dynmap config file to read from and write to
  -h, --help                   help for dynmap-area-generator
  -i, --image string           SVG image for zone data
  -l, --map-layer string       The layer in the Dynmap config to update/create
      --offset-x float         Offset in X from SVG image
      --offset-y float         Offset in Y from SVG image
      --scale-x float          Scale the resulting coordinate in X from SVG
      --scale-y float          Scale the resulting coordinates in Y from SVG
```

## Example
```bash
./dynmap-area-generator --image ./zones.svg --map-layer "Custom areas" --dynmap-config /path/to/dynmap/config.yaml \
    --offset-x 0 --offset-y 0 --scale-x 0 --scale-y 0
```

## Creating your own SVG

The SVG needs to follow a special format to work. It needs to have a layer with id `areas`. All paths for export should
reside in that layer.

An empty base image is provided in the repo.
