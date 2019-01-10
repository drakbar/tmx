package tmx

type wangset struct {
  Name         string         `json:"name"`         // name of the wang set
  Tile         int            `json:"tile"`         // local id of tile
  CornerColors []wangcolor    `json:"cornercolors"` // array of wang colors
  EdgeColors   []wangcolor    `json:"edgecolors"`   // array of wang colors
  WangTiles    []wangtile     `json:"wangtiles"`    // array of wang tiles
}

type wangcolor struct {
  Color       string  `json:"color"`       // hex color (#rrggbb or #aarrggbb)
  Name        string  `json:"name"`        // name of the wang color
  Probability float64 `json:"probability"` // probability used when randomizing
  Tile        int     `json:"tile"`        // local tile id of the wang color
}

type wangtile struct {
  TileId int   `json:"tileid"` // local id of tile
  DFlip  bool  `json:"dflip"`  // tile is flipped diagonally
  HFlip  bool  `json:"hflip"`  // tile is flipped horizontally
  VFlip  bool  `json:"vflip"`  // tile is flipped vertically
  WangId []int `json:"wangid"` // array of wang color indexes
}
