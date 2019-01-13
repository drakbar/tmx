package tmx

const (
	// Bits on the far end of the 32-bit global tile ID
	// are used for tile flags
	horizontalFlag = 0x80000000
	verticalFlag   = 0x40000000
	diagonalFlag   = 0x20000000
)

type Tile struct {
	gid             uint32 // the id of the tile in the tile layer
	lid             uint32 // the id of the tile in the tileset
	tileset         string // the tileset this tile is a part of
	horizontialFlip bool   // is the tile flipped over the x-axis
	verticalFlip    bool   // is the tile flipped over the y-axis
	diagonalFlip    bool   // is the tile flipped diagonally
	nil             bool   // if the global id is zero
}

var nilTile = &Tile{nil: true}

func (t Tile) Gid() uint32 {
	return t.gid
}

func (t Tile) Lid() uint32 {
	return t.lid
}

func (t Tile) Tileset() string {
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

// clearHighBits flips bits 31,30,29 to zero and returns a gid.
func clearHighBits(n uint32) uint32 {
	return n &^ (horizontalFlag | verticalFlag | diagonalFlag)
}

// flipFlags returns a bool for each bit (31,30,29).
func flipFlags(n uint32) (h, v, d bool) {
	return (n & horizontalFlag) == horizontalFlag,
		(n & verticalFlag) == verticalFlag,
		(n & diagonalFlag) == diagonalFlag
}

// localId returns the index of a tile inside a tileset.
func localId(g uint32, f int) uint32 {
	return g - uint32(f)
}
