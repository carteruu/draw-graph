# draw-graph

English | [简体中文](readme-cn.md)

`draw-gragh` is a draw DAG util

## Project layout and command line tools

* cmd
    * main
      Read `.json` file and generate picture

        * `-d` reads all `json` files in the specified directory
        * `-i` read the specified `json` files, splitting multiple files by commas
        * `-o` specify the output directory

## Data format

The data format can be provided via `json` files, refer to the format in [example](cmd/data/example.json).

## Acknowledgements

gridder - github.com/shomali11/gridder