package tmx

const (
  // chunk size in bytes  
  infiniteChunkSize = 16*16*numBytes
  numBytes          = 4 
)

const (
  // string constants
  groupLayer     = "group"
  tileLayer      = "tilelayer"
  objectLayer    = "objectgroup"
  base_64        = "base64"
  csv            = ""
)

type tilemap struct {
  Version         float32    `json:"version"`         // the json format version
  Tiledversion    string     `json:"tiledversion"`    // the tiled version used to save the file
  Type            string     `json:"type"`            // map (since 1.0)
  Backgroundcolor string     `json:"backgroundcolor"` // hex-formatted color (#RRGGBB or #AARRGGBB)
  Orientation     string     `json:"orientation"`     // orthogonal, isometric, staggered or hexagonal
  Renderorder     string     `json:"renderorder"`     // rendering direction (orthogonal maps only)
  StaggerAxis     string     `json:"staggeraxis"`     // x or y (staggered / hexagonal maps only)
  StaggerIndex    string     `json:"staggerindex"`    // odd or even (staggered / hexagonal maps only)
  Width           int        `json:"width"`           // number of tile columns
  Height          int        `json:"height"`          // number of tile rows
  Tilewidth       int        `json:"tilewidth"`       // map grid width
  Tileheight      int        `json:"tileheight"`      // map grid height
  HexSideLength   int        `json:"hexsidelength"`   // length of the side of a hex tile in pixels
  Nextobjectid    int        `json:"nextobjectid"`    // auto-increments for each placed object
  NextLayerId     int        `json:"nextlayerid"`     // auto-increments for each layer
  Infinite        bool       `json:"infinite"`        // whether the map has infinite dimensions
  Layers          []layer    `json:"layers"`          // array of Layers
  Tilesets        []tileset  `json:"tilesets"`        // array of tilesets
  Properties      []property `json:"properties"`      // a list of properties (name, value, type)
}


// processLayers determines what data needs processed for a given map.
func (m *tilemap) processLayers(ls *[]layer) (e error) {
  for i := 0; i < len((*ls)); i++ {
    // peel the layers one by one
    l := &(*ls)[i]

    switch l.Type {
    case groupLayer:
      // a group is a set of layers, recursively call process layers
      e = m.processLayers(&l.Layers)
      if e != nil {
        return
      }
    
    case tileLayer:
      // process the tile data
      e = m.processLayer(l)
      if e != nil {
        return
      }
    
    case objectLayer:
      // load in template files and transform template objects into proper 
      // objects
      e = processTemplates(&(l.Objects))
      if e != nil {
        return
      }
      // find objects that are tiles, adjust gids, and set flags
      processTileObjects(&(l.Objects))
      // adjust the points of the polygons and polylines
      translatePoints(&(l.Objects))
    }
  }
  return
}

// processLayer determines where the tile data is stored in the map and sends 
// it out for processing.
func (m *tilemap) processLayer(l *layer) (e error) {
  if m.Infinite {
    // tile data is in the chunks
    for j := 0; j < len(l.Chunks); j++ {
      c := &l.Chunks[j]
      if e = m.processTileData(&c.Data, *l); e != nil {
          return
      } 
    }
  } else {
    // tile data is in the layer
    if e = m.processTileData(&l.Data, *l); e != nil {
        return
    } 
  }
  return
}

// processTileData sends the tile data out to be decoded and extracted.
func (m *tilemap) processTileData(d *interface{}, l layer) (e error) {
  // make sure the pointer is not nil before a dereference
  if d == nil {
    return nilDataPtr
  } 
  switch l.Encoding {
  case base_64:
    // encoding is base64
    e = decodeBase64(d, l.Compression)
  
  case csv:
    // encoding is csv or xml in both cases the underlying data structure is 
    // identical
    e = decodeCSV(d)
  
  default:
    // encoding is unsupported 
    return unsupportedEncoding
  }
  if e != nil {
    return 
  }  
  // check for flipped tiles
  return m.extractTileData(d)      
}

// extractTileData extracts and correlates information about each tile and 
// repackages it for consumption.
func (m *tilemap) extractTileData(d *interface{}) (e error) {
  // make sure the data is a byte array
  b, ok := (*d).([]byte)
  if !ok {
    return highBitDataMismatch
  }
  // make sure there is enough data for every tile
  if m.Infinite && len(b) != infiniteChunkSize {
    return dataSizeMismatch
  }
  if !m.Infinite && len(b) != (m.Width * m.Height * numBytes) {
    return dataSizeMismatch
  }

  var data []*Tile 
  for i := 0; i < len(b); i += numBytes {
    // shift the bytes back into a variable
    n := compressBytes(b[i:])
    if n == 0 {
      // there isn't a tile at this location
      data = append(data, nilTile)
      continue
    }
    
    // clear the high bits from the gid and get the flip flags
    gid   := clearHighBits(n)
    h,v,d := flipFlags(n)
    
    // verify that the gid is a valid id and add it to the data container
    verified := false
    for j := 0; j < len(m.Tilesets); j++ {
      t := &m.Tilesets[j]
      lastId := (t.Firstgid + t.Tilecount) - 1
      // if the global id is in this tileset
      if int(gid) >= t.Firstgid && int(gid) <= lastId {
        // add the tile into the container
        data = append(data, &Tile{ 
          // set the global and local ids
          gid: gid, lid: localId(gid, t.Firstgid),
          // set pointer to the tileset that this gid belongs to
          tileset: t,
          // set flip flags
          horizontialFlip: h, verticalFlip: v, diagonalFlip: d,
          nil: false })
        verified = true
        break;
      }
    }
    if !verified {
      // could not verify the global id
      return badGlobalId
    }
  }
  // reset the data container
  *d = data
  return
}
