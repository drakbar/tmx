package tmx

type tileset struct {
  Name             string     `json:"name"`             // name of tileset
  Type             string     `json:"type"`             // "tileset"
  Source           string     `json:"source"`           // path to tileset file
  Image            string     `json:"image"`            // path to image file
  TransparentColor string     `json:"transparentcolor"` // hex color (#rrggbb)
  Firstgid         int        `json:"firstgid"`         // first tile in a set
  Tilewidth        int        `json:"tilewidth"`        // width of tiles   
  Tileheight       int        `json:"tileheight"`       // height of tiles
  Spacing          int        `json:"spacing"`          // space between tiles
  Margin           int        `json:"margin"`           // space around edge
  Tilecount        int        `json:"tilecount"`        // number of tiles 
  Columns          int        `json:"columns"`          // number of columns
  Imagewidth       int        `json:"imagewidth"`       // width of image 
  Imageheight      int        `json:"imageheight"`      // height of image
  Grid             grid       `json:"grid"`             // see <grid>
  TileOffsets      offset     `json:"tileoffset"`       // see <tileoffset>
  TerrianTypes     []terrian  `json:"terrains"`         // array of terrains 
  Tiles            []tile     `json:"tiles"`            // array of tiles
  Wangsets         []wangset  `json:"wangsets"`         // array of wang sets
  Properties       []property `json:"properties"`       // a list of properties     
}

type external struct {
  Name         string  `json:"name"`         // same as tileset
  Image        string  `json:"image"`        // same as tileset
  Type         string  `json:"type"`         // same as tileset
  Tilewidth    int     `json:"tilewidth"`    // same as tileset
  Tileheight   int     `json:"tileheight"`   // same as tileset
  Spacing      int     `json:"spacing"`      // same as tileset
  Margin       int     `json:"margin"`       // same as tileset
  Tilecount    int     `json:"tilecount"`    // same as tileset
  Imagewidth   int     `json:"imagewidth"`   // same as tileset
  Imageheight  int     `json:"imageheight"`  // same as tileset
  Columns      int     `json:"columns"`      // same as tileset
  Tiledversion string  `json:"tiledversion"` // external only
  Version      float64 `json:"version"`      // external only
}

type grid struct {
  Orientation string `json:"orientation"` // orthogonal or isometric
  Width       int    `json:"width"`       // width of a grid cell
  Height      int    `json:"height"`      // height of a grid cell
}

type terrian struct {
  Name       string     `json:"name"`       // name of terrain
  Tile       int        `json:"tile"`       // local id of terrain tile
  Properties []property `json:"properties"` // a list of properties
}

type offset struct {
  X int `json:"x"` // horizontal offset in pixels
  Y int `json:"y"` // vertical offset in pixels (positive is down)
}

type tile struct {
  Type         string     `json:"type"`        // type of the tile
  Image        string     `json:"image"`       // image representing this tile
  ImageWidth   int        `json:"imagewidth"`  // width of the tile image                                                     
  ImageHeight  int        `json:"imageheight"` // height of the tile image
  Id           int        `json:"id"`          // local id of the tile
  ObjectGroup  layer      `json:"objectgroup"` // layer with type objectgroup
  Terrian      []int      `json:"terrain"`     // index of each terrain corner
  Animation    []frame    `json:"animation"`   // array of frames
  Properties   []property `json:"properties"`  // a list of properties
}

type frame struct {
  TileId   int `json:"tileid"`   // local tile id representing this frame
  Duration int `json:"duration"` // frame duration in milliseconds
}