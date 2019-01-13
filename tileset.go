package tmx

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"io"
	"io/ioutil"
	"reflect"
)

const (
	// string constants
	empty        = ""
	gZip         = "gzip"
	zLib         = "zlib"
	uncompressed = ""
)

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
	Type        string     `json:"type"`        // type of the tile
	Image       string     `json:"image"`       // image representing this tile
	ImageWidth  int        `json:"imagewidth"`  // width of the tile image
	ImageHeight int        `json:"imageheight"` // height of the tile image
	Id          int        `json:"id"`          // local id of the tile
	ObjectGroup layer      `json:"objectgroup"` // layer with type objectgroup
	Terrian     []int      `json:"terrain"`     // index of each terrain corner
	Animation   []frame    `json:"animation"`   // array of frames
	Properties  []property `json:"properties"`  // a list of properties
}

type frame struct {
	TileId   int `json:"tileid"`   // local tile id representing this frame
	Duration int `json:"duration"` // frame duration in milliseconds
}

// decodeCSV splits up the global ids and saves them  into a byte array.
func decodeCSV(d *interface{}) (e error) {
	// make sure the underlying data structure is correct
	if c, ok := (*d).([]interface{}); ok {
		// make a byte array that can hold all four bytes per tile
		b := make([]byte, len(c)*numBytes)
		// break each tile id into the appropriate number of bytes and store them
		// in the byte array
		for i, v := range c {
			binary.LittleEndian.PutUint32(b[i*numBytes:], uint32(v.(float64)))
		}
		// reset the data container
		*d = b
	} else {
		return csvDataMismatch
	}
	return
}

// decodeBase64 strips off the base64 encoding, uncompresses (if necessary),
// and saves the data into a byte array.
func decodeBase64(d *interface{}, c string) (e error) {
	// make sure the data type is a string
	if _, ok := (*d).(string); !ok {
		return dataStringMismatch
	}
	// make sure it isn't just an empty string
	if *d == empty {
		return missingData
	}
	// trim off any additional white spaces
	b := bytes.TrimSpace([]byte((*d).(string)))
	r := bytes.NewReader(b)

	// setup base64 decoder
	enc := base64.NewDecoder(base64.StdEncoding, r)

	var dec io.Reader
	// switch based on compression type
	switch c {
	case gZip:
		dec, e = gzip.NewReader(enc)
		if e != nil {
			return
		}

	case zLib:
		dec, e = zlib.NewReader(enc)
		if e != nil {
			return
		}

	case uncompressed:
		dec = enc

	default:
		return unsupportedCompression
	}
	// reset data container
	*d, e = ioutil.ReadAll(dec)
	return
}

// compressBytes converts a byte array into a unsigned int32.
func compressBytes(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

// processTilesets determines if a tileset needs to be imported from an
// external tileset file.
func processTilesets(s *[]tileset) (e error) {
	for i := 0; i < len(*s); i++ {
		ts := &(*s)[i]
		// determine if this is a external tileset
		if ts.Source != empty {
			// load in the external tileset from file
			var ex external
			if ex, e = loadTileset(ts.Source); e != nil {
				return
			}
			// get the reflect value of the external tileset and the corresponding
			// tileset
			src, dst := reflect.ValueOf(ex), reflect.ValueOf(ts).Elem()
			// copy the fields of the src into the dst
			if e = copyFields(&src, &dst); e != nil {
				return
			}
		}
	}
	return
}
