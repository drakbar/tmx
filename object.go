package tmx

type object struct {
  Name            string     `json:"name"`       // name field in editor
  Type            string     `json:"type"`       // type field in editor
  Template        string     `json:"template"`   // path to a template file
  Source          string                         // path to a template tileset
  Gid             int        `json:"gid"`        // global id
  Lid             int                            // local id to tileset
  Id              int        `json:"id"`         // incremental id
  X               float64    `json:"x"`          // x coordinate in pixels
  Y               float64    `json:"y"`          // y coordinate in pixels
  Width           float64    `json:"width"`      // width ignored if using gid
  Height          float64    `json:"height"`     // height ignored if using gid
  Rotation        float64    `json:"rotation"`   // angle in degrees clockwise
  Visible         bool       `json:"visible"`    // is object shown in editor
  Ellipse         bool       `json:"ellipse"`    // is object an ellipse
  Point           bool       `json:"point"`      // is object a point
  HorizontialFlip bool                           // is flipped over the x-axis
  VerticalFlip    bool                           // is flipped over the y-axis
  DiagonalFlip    bool                           // is flipped diagonally
  Text            text       `json:"text"`       // raw string of text object
  Polygon         []point    `json:"polygon"`    // list of points x/y coords
  Polyline        []point    `json:"polyline"`   // list of points x/y coords
  Properties      []property `json:"properties"` // list of custom properties
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