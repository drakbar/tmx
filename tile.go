package tmx

type Tile struct {
  gid             uint32   // the id of the tile in the tile layer
  lid             uint32   // the id of the tile in the tileset
  tileset         *tileset // the tileset this tile is a part of
  horizontialFlip bool     // is the tile flipped over the x-axis
  verticalFlip    bool     // is the tile flipped over the y-axis
  diagonalFlip    bool     // is the tile flipped diagonally
  nil             bool     // if the global id is zero
}

var nilTile = &Tile{ nil:true }

func (t Tile) Gid() uint32 {
  return t.gid
}

func (t Tile) Lid() uint32 {
  return t.lid
}

func (t Tile) Tileset() *tileset {
  return t.tileset
}

func (t Tile) HorizontialFlip() bool {
  return t.horizontialFlip
}

func (t Tile) VerticalFlip() bool {
  return t.verticalFlip
}

func (t Tile) DiagonalFlip() bool {
  return t.diagonalFlip
}

func (t Tile) Nil() bool {
  return t.nil
}