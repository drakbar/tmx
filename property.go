package tmx

type property struct {
  Name  string      `json:"name"`  // name of the property
  Type  string      `json:"type"`  // string, int, float, bool, color or file
  Value interface{} `json:"value"` // value of the property
}
