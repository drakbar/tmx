package tmx

type layer struct {
  Name             string      `json:"name"`             // name of the layer
  Type             string      `json:"type"`             // type of layer
  DrawOrder        string      `json:"draworder"`        // topdown (default)
  Compression      string      `json:"compression"`      // zlib, gzip or empty
  Encoding         string      `json:"encoding"`         // csv or base64 
  Image            string      `json:"image"`            // imagelayer only
  TransparentColor string      `json:"transparentcolor"` // hex color (#rrggbb)
  Id               int         `json:"id"`               // incremental id
  X                int         `json:"x"`                // tile offset x-axis
  Y                int         `json:"y"`                // tile offset y-axis 
  Width            int         `json:"width"`            // column count 
  Height           int         `json:"height"`           // row count 
  Opacity          float64     `json:"opacity"`          // between 0 and 1
  Offsetx          float64     `json:"offsetx"`          // pixel offset x-axis
  Offsety          float64     `json:"offsety"`          // pixel offset y-axis
  Visible          bool        `json:"visible"`          // is shown in editor
  Data             interface{} `json:"data"`             // array gids
  Layers           []layer     `json:"layers"`           // group of layers
  Chunks           []chunk     `json:"chunks"`           // infinte map gids
  Objects          []object    `json:"objects"`          // array of objects
  Properties       []property  `json:"properties"`       // list of properties
}

type chunk struct {
  X      int         `json:"x"`      // x coordinate in tiles
  Y      int         `json:"y"`      // y coordinate in tiles
  Width  int         `json:"width"`  // width in tiles
  Height int         `json:"height"` // height in tiles
  Data   interface{} `json:"data"`   // unsigned int (gids) or base64-encoded
}