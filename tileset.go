package tmx

type tileset struct {
  Name             string     `json:"name"`             // name given to this tileset
  Type             string     `json:"type"`             // tileset (for tileset files, since 1.0)
  Source           string     `json:"source"`           // external tileset file
  Image            string     `json:"image"`            // image used for tiles in this set
  TransparentColor string     `json:"transparentcolor"` // hex-formatted color (#rrggbb)
  Firstgid         int        `json:"firstgid"`         // gid corresponding to the first tile in the set
  Tilewidth        int        `json:"tilewidth"`        // maximum width of tiles in this set   
  Tileheight       int        `json:"tileheight"`       // maximum height of tiles in this set
  Spacing          int        `json:"spacing"`          // spacing between adjacent tiles in image (pixels)
  Margin           int        `json:"margin"`           // buffer between image edge and first tile (pixels)
  Tilecount        int        `json:"tilecount"`        // the number of tiles in this tileset
  Columns          int        `json:"columns"`          // the number of tile columns in the tileset
  Imagewidth       int        `json:"imagewidth"`       // width of source image in pixels
  Imageheight      int        `json:"imageheight"`      // height of source image in pixels
  Grid             grid       `json:"grid"`             // see <grid>
  TileOffsets      offset     `json:"tileoffset"`       // see <tileoffset>
  TerrianTypes     []terrian  `json:"terrains"`         // array of terrains 
  Tiles            []tile     `json:"tiles"`            // array of tiles
  Wangsets         []wangset  `json:"wangsets"`         // array of wang sets
  Properties       []property `json:"properties"`       // a list of properties (name, value, type)     
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
  Tile       int        `json:"tile"`       // local id of tile representing terrain
  Properties []property `json:"properties"` // a list of properties (name, value, type)
}

type offset struct {
  X int `json:"x"` // horizontal offset in pixels
  Y int `json:"y"` // vertical offset in pixels (positive is down)
}

type tile struct {
  Type         string     `json:"type"`        // the type of the tile
  Image        string     `json:"image"`       // image representing this tile
  ImageWidth   int        `json:"imagewidth"`  // imagewidth  int width of the tile image in pixels                                                      
  ImageHeight  int        `json:"imageheight"` // imageheight int height of the tile image in pixels
  Id           int        `json:"id"`          // local id of the tile
  ObjectGroup  layer      `json:"objectgroup"` // layer with type objectgroup
  Terrian      []int      `json:"terrain"`     // array index of terrain for each corner of tile
  Animation    []frame    `json:"animation"`   // array of frames
  Properties   []property `json:"properties"`  // a list of properties (name, value, type)
}

type frame struct {
  TileId   int `json:"tileid"`   // local tile id representing this frame
  Duration int `json:"duration"` // frame duration in milliseconds
}