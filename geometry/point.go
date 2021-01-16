package geometry

type Point struct {
	X, Y, Z int
}

type PointWithColour struct {
	Point  Point
	Colour byte
}

type PointF struct {
	X, Y, Z float64
}

type Bounds struct {
	Min, Max Point
}

func (b *Bounds) GetSize() Point {
	return Point{X: b.Max.X - b.Min.X, Y: b.Max.Y - b.Min.Y, Z: b.Max.Z - b.Min.Z}
}

// IsInBounds returns true if the point is within the supplied bounds
func (p *Point) IsInBounds(b Bounds) bool {
	return p.X >= b.Min.X && p.X < b.Max.X &&
		p.Y >= b.Min.Y && p.Y < b.Max.Y &&
		p.Z >= b.Min.Z && p.Z < b.Max.Z
}

// NewPoint is a shortcut to create a new point
func NewPoint(x,y,z int) Point {
	return Point{
		X: x,
		Y: y,
		Z: z,
	}
}

// NewBounds is a shortcut to create a new bounds
func NewBounds(x1,y1,z1,x2,y2,z2 int) Bounds {
	return Bounds{
		Min: NewPoint(x1,y1,z1),
		Max: NewPoint(x2,y2,z2),
	}
}