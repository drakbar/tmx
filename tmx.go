package tmx

import (
  "os"
  "io"
  "reflect"
  "errors"
  "bytes"
  "path"
  "path/filepath"
  "encoding/json"
  "encoding/base64"
  "encoding/binary"
  "compress/gzip"
  "compress/zlib"
  "io/ioutil"
)

const (
  // string constants
  empty        = "" 
  gZip         = "gzip"
  zLib         = "zlib"
  uncompressed = ""
)

const (
  // Bits on the far end of the 32-bit global tile ID 
  // are used for tile flags
  horizontalFlag = 0x80000000
  verticalFlag   = 0x40000000
  diagonalFlag   = 0x20000000
)

var (
  // runtime errors
  nilDataPtr             = errors.New("data pointer is nil")
  missingData            = errors.New("base64 data is an empty string")
  unsupportedEncoding    = errors.New("the encoding type is unsupported")
  unsupportedCompression = errors.New("the compression type is unsupported")
  dataStringMismatch     = errors.New("the data is not of type string")
  dataSizeMismatch       = errors.New("tile data and map size do not match")
  csvDataMismatch        = errors.New("csv data structure incorrect")
  highBitDataMismatch    = errors.New("tile data is not a byte array")
  reflectionBothWrong    = errors.New("both the src and dst are not a structures")
  reflectionSrcWrong     = errors.New("the src is not a structures")
  reflectionDstWrong     = errors.New("the dst is not a structures")
  badGlobalId            = errors.New("global id could not be found in any tile set")
  templateNotLoaded      = errors.New("the template wasn't loaded")
)

var mapDirectory string

// LoadTileMap reads in a tilemap from disk, sends the data out to be 
// processed, and finally returns a tilemap.
func LoadTileMap(fp string) (m tilemap, e error) {
  // get path to the map directory so relative paths can be resolved from there
  if pwd, e := os.Getwd(); e == nil {
    mapDirectory = filepath.FromSlash(path.Join(pwd, fp))
  } 

  // read in tile map from disk
  var b []byte
  if b, e = ioutil.ReadFile(fp); e == nil {
    // store the json data into tilemap
    if e = json.Unmarshal(b, &m); e == nil {
      // determine if there are external tilesets and load them if necessary
      if e = processTilesets(&m.Tilesets); e != nil {
        return
      }
      // decode and if necessary decompress all layer data to a workable format
      if e = m.processLayers(&m.Layers); e != nil {
        return
      }
    }
  }
  return
}

// loadTileset reads in a tileset from disk, and returns a external tileset.
func loadTileset(fp string) (ts external, e error) {
  // reslove path
  fp = path.Join(mapDirectory, fp)
  // read in tile set from disk
  var b []byte
  if b, e = ioutil.ReadFile(fp); e == nil {
    e = json.Unmarshal(b, &ts)
    return 
  }  
  return
}

// loadTemplate reads in a template from disk, and returns an external tileset.
func loadTemplate(fp string) (t template, e error) {
  // reslove path
  fp = path.Join(mapDirectory, fp)
  // read in template from disk
  var b []byte
  if b, e = ioutil.ReadFile(fp); e == nil {
    e = json.Unmarshal(b, &t)
    return 
  }  
  return
}

// getTemplates minimizes the numbers of reads from the disk by calling load 
// template only once for each template file.
func getTemplates(objs *[]object) (tmp map[string]template, e error) {
  tmp = make(map[string]template) 
  for i := 0; i < len(*(objs)); i++ {
    o := &(*objs)[i]
    if o.Template != empty {
      if t, loaded := tmp[o.Template]; !loaded {
         if t, e = loadTemplate(o.Template); e != nil {
          return
        }
        tmp[o.Template] = t
      }     
    }
  }
  return
}

// decodeCSV splits up the global ids and saves them  into a byte array.
func decodeCSV(d *interface{}) (e error) {
  // make sure the underlying data structure is correct
  if c, ok := (*d).([]interface{}); ok {
    // make a byte array that can hold all four bytes per tile
    b := make([]byte, len(c)*numBytes)
    // break each tile id into the appropriate number of bytes and store them 
    // in the byte array
    for i,v := range c {
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
  if _,ok := (*d).(string); !ok {
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

// clearHighBits flips bits 31,30,29 to zero and returns a gid.
func clearHighBits(n uint32) uint32 {
  return n&^(horizontalFlag|verticalFlag|diagonalFlag)
}

// flipFlags returns a bool for each bit (31,30,29).
func flipFlags(n uint32) (h,v,d bool) {
  return (n&horizontalFlag) == horizontalFlag,
  (n&verticalFlag) == verticalFlag,
  (n&diagonalFlag) == diagonalFlag
}

// localId returns the index of a tile inside a tileset.
func localId(g uint32, f int) uint32 {
  return g - uint32(f)
}

// compressBytes converts a byte array into a unsigned int32.
func compressBytes(b []byte) uint32 {
  return binary.LittleEndian.Uint32(b)
}

// processTemplates determines if there are templates, if so it loads the 
// template, and applies it to the object.
func processTemplates(objs *[]object) (e error) {
  // load in the templates
  tmp, er := getTemplates(objs)
  // there was an error or none of the object are templates
  if er != nil || len(tmp) == 0 {
    return er
  }
  for i := 0; i < len(*objs); i++ {
    o := &(*objs)[i]
    if o.Template != empty {
      // get the template
      t, loaded := tmp[o.Template]
      if !loaded {
        // this is not the template your looking for
        return templateNotLoaded
      }
      to := t.Object
      // get the reflect value of the object and template object 
      src, dst := reflect.ValueOf(*o), reflect.ValueOf(&to).Elem()
      // copy the fields of the src into the dst
      if e = copyFields(&src, &dst); e != nil {
        return
      }
      // insert new and overridden properties
      to.Properties = overrideProperties(o.Properties, to.Properties)
      // insert overridden points in polygons and polylines
      to.Polygon  = overridePoints(o.Polygon, to.Polygon)
      to.Polyline = overridePoints(o.Polyline, to.Polyline)
      // assign the tileset to the object reference
      (*o).Source = t.Tileset.Source
      // place the fully constructed object into the set of objects
      *o = to
    }
  }
  return
}

// processTileObjects checks for objects that are from a tileset and extracts 
// the gid and flip flags of the tile and saves them to the object. 
func processTileObjects(objs *[]object) {
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
    }
  }
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

// copyFields copies the fields of one structure over to another. It does not
// copy slices or structure however.
func copyFields(src, dst *reflect.Value) (e error) {
  // verify that both are of type struct
  if e = checkStruct(*src, *dst); e != nil {
    return
  }
  for i := 0; i < src.NumField(); i++ {
    // get the src field name and value
    n, v := src.Type().Field(i).Name, src.Field(i)
    // get the dst field 
    f := dst.FieldByName(n)
    // if the field exists and it can be assigned a value
    if f.IsValid() && f.CanSet() {
      // assign the field a value based on its type
      switch v.Type().Kind() {
      case reflect.String:
        f.SetString(v.String())    
      case reflect.Int:
        // nothing has a gid of zero
        if n == "Gid" && v.Int() == 0 {
          continue
        }
        f.SetInt(v.Int())
      case reflect.Float64:
        f.SetFloat(v.Float())
      case reflect.Bool:
        f.SetBool(v.Bool())
      default:
        continue
      }  
    }
  }
  return
}

// checkStruct verifies if both the source an destination reflect values are of
// the type struct.
func checkStruct(src, dst reflect.Value) (e error) {
  if src.Kind() != reflect.Struct && dst.Kind() != reflect.Struct {
    // src and dst values are not structs
    return reflectionBothWrong
  }
  if src.Kind() != reflect.Struct {
    // src value is not a struct 
    return reflectionSrcWrong
  }
  if dst.Kind() != reflect.Struct {
    // dst value is not a struct
    return reflectionDstWrong
  }
  return
}


// overrideProperties combines the overridden properties of a template and 
// object.  
func overrideProperties(o, n []property) []property {
  for i := 0; i < len(o); i++ {
    present := false
    for j := 0; j < len(n); j++ {
      if o[i].Name == n[j].Name {
        // if property exists overwrite it
        n[j].Value = o[i].Value
        present = true
      }
    }
    if !present {
      // if the property didn't exist, add it
      n = append(n, o[i])
    }
  }
  return n
}

// overridePoints determines if there are points for polygons and ploylines 
// that need to be overridden and handles them accordingly.
func overridePoints(o, n []point) []point {
  if len(o) > 0 {
    n = o
  }
  return n
}

// translatePoints adjusts the coordinates of polygons and polylines from being
// relative coordinates to being global coordinates.
func translatePoints(objs *[]object) {
  for i := 0; i < len(*(objs)); i++ {
    o := &(*objs)[i]
    // if there are polygons, fix their points
    for j := 0; j < len(o.Polygon); j++ {
      p := &o.Polygon[j]
      p.X += o.X 
      p.Y += o.Y
    }
    // if there are polylines, fix their points
    for j := 0; j < len(o.Polyline); j++ {
      p := &o.Polyline[j]
      p.X += o.X 
      p.Y += o.Y
    }
  }
}
