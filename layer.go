package tmx

type layer struct {
  Name             string      `json:"name"`             // name assigned to this layer
  Type             string      `json:"type"`             // tilelayer, objectgroup, imagelayer or group
  DrawOrder        string      `json:"draworder"`        // topdown (default) or index objectgroup only
  Compression      string      `json:"compression"`      // zlib, gzip or empty tilelayer only
  Encoding         string      `json:"encoding"`         // csv (default) or base64 tilelayer only
  Image            string      `json:"image"`            // image used by this layer imagelayer only
  TransparentColor string      `json:"transparentcolor"` // hex-formatted color (#rrggbb) imagelayer only
  Id               int         `json:"id"`               // incremental id - unique across all layers
  X                int         `json:"x"`                // int horizontal layer offset in tiles always 0
  Y                int         `json:"y"`                // int vertical layer offset in tiles always 0
  Width            int         `json:"width"`            // column count same as map width for fixed-size maps
  Height           int         `json:"height"`           // row count same as map height for fixed-size maps
  Opacity          float64     `json:"opacity"`          // value between 0 and 1
  Offsetx          float64     `json:"offsetx"`          // horizontal layer offset in pixels (default: 0)
  Offsety          float64     `json:"offsety"`          // vertical layer offset in pixels (default: 0)
  Visible          bool        `json:"visible"`          // whether layer is shown or hidden in editor
  Data             interface{} `json:"data"`             // array of unsigned int (gids) or base64-encoded data tilelayer only
  Layers           []layer     `json:"layers"`           // array of layers group only
  Chunks           []chunk     `json:"chunks"`           // array of chunks infinte map only
  Objects          []object    `json:"objects"`          // array of objects objectgroup only
  Properties       []property  `json:"properties"`       // list of properties (name, value, type)
}

type chunk struct {
  X      int         `json:"x"`      // x coordinate in tiles
  Y      int         `json:"y"`      // y coordinate in tiles
  Width  int         `json:"width"`  // width in tiles
  Height int         `json:"height"` // height in tiles
  Data   interface{} `json:"data"`   // array of unsigned int (gids) or base64-encoded data
}