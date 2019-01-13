package tmx

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// path to the tile map from the current working directory
var mapDirectory string

// LoadTileMap reads in a tilemap from disk, sends the data out to be
// processed, and finally returns a tilemap.
func LoadTileMap(fp string) (m tilemap, e error) {
	// get path to the map directory so relative paths can be resolved from there
	if e = resolveMapPath(fp); e != nil {
		return
	}
	// read in tile map from disk
	var b []byte
	if b, e = read(fp); e == nil {
		// store the json data into tilemap
		if e = decode(b, &m); e == nil {
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
	fp = externalFilePath(fp)
	// read in tile set from disk
	var b []byte
	if b, e = read(fp); e == nil {
		e = decode(b, &ts)
		return
	}
	return
}

// loadTemplate reads in a template from disk, and returns an external tileset.
func loadTemplate(fp string) (t template, e error) {
	// reslove path
	fp = externalFilePath(fp)
	// read in template from disk
	var b []byte
	if b, e = read(fp); e == nil {
		e = decode(b, &t)
		return
	}
	return
}

// read loads a file from the disk and reads it into a byte array.
func read(fp string) ([]byte, error) {
	return ioutil.ReadFile(fp)
}

// decode takes an array of bytes and places the data inside the provided
// structure.
func decode(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}

// resolveMapPath saves the path to where the map file is located so that
// external files can be resolved later on.
func resolveMapPath(fp string) error {
	pwd, e := os.Getwd()
	if e != nil {
		return e
	}
	dir, _ := filepath.Split(fp)
	mapDirectory = filepath.FromSlash(path.Join(pwd, dir))
	return nil
}

// externalFilePath get the path of a file relative to the working directory.
func externalFilePath(fp string) string {
	return filepath.FromSlash(path.Join(mapDirectory, fp))
}

// matchTilesetName returns whether the names of the tilesets are equal.
func matchTilesetName(fp1, fp2 string) bool {
	return filename(fp1) == filename(fp2)
}

// filename strips off the path from a filepath, and returns the filename.
func filename(fp string) string {
	_, f := filepath.Split(fp)
	return f
}
