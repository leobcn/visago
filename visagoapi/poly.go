package visagoapi

// BoundingPoly is used to store the
// vertexes marking the postition of the face.
type BoundingPoly struct {
	Vertices []*Vertex `json:"vertices,omitempty"`
}

// Vertex is the x and y coordinates of a vertex
type Vertex struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}
