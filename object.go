package tmx

type object struct {
  Name            string     `json:"name"`       // string assigned to name field in editor
  Type            string     `json:"type"`       // string assigned to type field in editor
  Template        string     `json:"template"`   // reference to a template file
  Source          string                         // if the object is a template and its from a tileset 
  Gid             int        `json:"gid"`        // global id
  Id              int        `json:"id"`         // incremental id - unique across all objects
  X               float64    `json:"x"`          // x coordinate in pixels
  Y               float64    `json:"y"`          // y coordinate in pixels
  Width           float64    `json:"width"`      // width in pixels ignored if using a gid
  Height          float64    `json:"height"`     // height in pixels ignored if using a gid
  Rotation        float64    `json:"rotation"`   // angle in degrees clockwise
  Visible         bool       `json:"visible"`    // whether object is shown in editor
  Ellipse         bool       `json:"ellipse"`    // used to mark an object as an ellipse
  Point           bool       `json:"point"`      // used to mark an object as a point
  HorizontialFlip bool                           // is the tile flipped over the x-axis
  VerticalFlip    bool                           // is the tile flipped over the y-axis
  DiagonalFlip    bool                           // is the tile flipped diagonally
  Text            text       `json:"text"`       // string key-value pairs
  Polygon         []point    `json:"polygon"`    // a list of x,y coordinates in pixels
  Polyline        []point    `json:"polyline"`   // a list of x,y coordinates in pixels
  Properties      []property `json:"properties"` // a list of properties (name, value, type)
}

type text struct {
  Text   string `json:"text"`       // the raw text value
  Color  string `json:"color"`      // color of the text
  Font   string `json:"fontfamily"` // text font
  HAlign string `json:"halign"`     // justify, right, and center (left is blank)  
  VAlign string `json:"valign"`     // center and bottom (top is blank)
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