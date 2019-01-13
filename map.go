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
  Version         float32    `json:"version"`         // json format version
  Tiledversion    string     `json:"tiledversion"`    // tiled version
  Type            string     `json:"type"`            // "map"
  Backgroundcolor string     `json:"backgroundcolor"` // hex color (#AARRGGBB)
  Orientation     string     `json:"orientation"`     // map type
  Renderorder     string     `json:"renderorder"`     // rendering direction
  StaggerAxis     string     `json:"staggeraxis"`     // x or y 
  StaggerIndex    string     `json:"staggerindex"`    // odd or even 
  Width           int        `json:"width"`           // number of tile columns
  Height          int        `json:"height"`          // number of tile rows
  Tilewidth       int        `json:"tilewidth"`       // map grid width
  Tileheight      int        `json:"tileheight"`      // map grid height
  HexSideLength   int        `json:"hexsidelength"`   // side length of hex
  Nextobjectid    int        `json:"nextobjectid"`    // unique for each object
  NextLayerId     int        `json:"nextlayerid"`     // unique for each layer
  Infinite        bool       `json:"infinite"`        // is map infinite
  Layers          []layer    `json:"layers"`          // layers
  Tilesets        []tileset  `json:"tilesets"`        // tilesets
  Properties      []property `json:"properties"`      // a list of properties
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
      e = m.processTileObjects(&(l.Objects))
      if e != nil {
        return
      }
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
    
    // verify that the gid is a valid id
    var t *tileset
    t, e = m.verifyGid(gid)
    if e != nil {
      return
    }

    // add the tile into the container
    data = append(data, &Tile{ 
      // set the global and local ids
      gid: gid, lid: localId(gid, t.Firstgid),
      // set pointer to the tileset that this gid belongs to
      tileset: t.Source,
      // set flip flags
      horizontialFlip: h, verticalFlip: v, diagonalFlip: d,
      nil: false })  
  }
  // reset the data container
  *d = data
  return
}

 // verifyGid confirms a gid is a valid id for a tile in one of the tilesets.
func (m *tilemap) verifyGid(gid uint32) (t *tileset, e error) {
  for i := 0; i < len(m.Tilesets); i++ {
    t := &m.Tilesets[i]
    lastId := (t.Firstgid + t.Tilecount) - 1
    // if the global id is in this tileset
    if int(gid) >= t.Firstgid && int(gid) <= lastId {
      return t, nil
    }
  }
  return nil, badGlobalId
}

// processTileObjects checks for objects that are from a tileset and extracts 
// the gid and flip flags of the tile and saves them to the object. 
func (m *tilemap) processTileObjects(objs *[]object) (e error) {
  for i := 0; i < len(*objs); i++ {
    o := &(*objs)[i]

    if o.Gid != 0 {
      // get flags
      h,v,d := flipFlags(uint32(o.Gid))
      // set flags
      (*o).HorizontialFlip = h
      (*o).VerticalFlip    = v
      (*o).DiagonalFlip    = d
      // strip out the flags
      (*o).Gid = int(clearHighBits(uint32(o.Gid)))
      // verify that there is a matching tileset
      if e = m.matchTileset(o); e != nil {
        return
      }
    }
  }
  return
}

// matchTileset compares the tileset of a template with the list of loaded
// tileset to verify that there is a tileset loaded for the template. It sets
// the local and global ids of the object.
func (m *tilemap) matchTileset(o *object) (e error) {
  for i := 0; i < len(m.Tilesets); i++ {
    t := m.Tilesets[i] 
    if matchTilesetName(t.Source, o.Source) {
      (*o).Lid = o.Gid
      (*o).Gid += (t.Firstgid - 1)
      _, e = m.verifyGid(uint32(o.Gid))
      return 
    } 
  }
  return noMatchingTileset
}
