package tmx

import "reflect"

type object struct {
	Name            string     `json:"name"`       // name field in editor
	Type            string     `json:"type"`       // type field in editor
	Template        string     `json:"template"`   // path to a template file
	Gid             int        `json:"gid"`        // global id
	Id              int        `json:"id"`         // incremental id
	X               float64    `json:"x"`          // x coordinate in pixels
	Y               float64    `json:"y"`          // y coordinate in pixels
	Width           float64    `json:"width"`      // width ignored if using gid
	Height          float64    `json:"height"`     // height ignored if using gid
	Rotation        float64    `json:"rotation"`   // angle in degrees clockwise
	Visible         bool       `json:"visible"`    // is object shown in editor
	Ellipse         bool       `json:"ellipse"`    // is object an ellipse
	Point           bool       `json:"point"`      // is object a point
	Text            text       `json:"text"`       // raw string of text object
	Polygon         []point    `json:"polygon"`    // list of points x/y coords
	Polyline        []point    `json:"polyline"`   // list of points x/y coords
	Properties      []property `json:"properties"` // list of custom properties
	Lid             int
	Source          string
	HorizontialFlip bool
	VerticalFlip    bool
	DiagonalFlip    bool
}

type text struct {
	Text   string `json:"text"`       // the raw text value
	Color  string `json:"color"`      // color of the text
	Font   string `json:"fontfamily"` // text font
	HAlign string `json:"halign"`     // justify, right, and center (left blank)
	VAlign string `json:"valign"`     // center and bottom (top blank)
	Wrap   bool   `json:"wrap"`       // whether to wrap the text
}

type point struct {
	X float64 `json:"x"` // x pixel coordinate
	Y float64 `json:"y"` // y pixel coordinate
}

type template struct {
	Type    string  `json:"type"`    // "template"
	Object  object  `json:"object"`  // all the same fields as an object
	Tileset tileset `json:"tileset"` // abbreviated tile set
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
			to.Polygon = overridePoints(o.Polygon, to.Polygon)
			to.Polyline = overridePoints(o.Polyline, to.Polyline)
			// place the fully constructed object into the set of objects
			*o = to
			// assign the tileset to the object reference
			(*o).Source = t.Tileset.Source
		}
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
