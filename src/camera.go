package main

import "math"

type OrbitCamera struct {
	Target Vec3
	Distance float64
	Yaw float64
	Pitch float64
	OpPresVAngle  float64
	Near  float64
	Far   float64
}

func NewOrbitCamera() OrbitCamera {
	return OrbitCamera{
		Target: Vec3{0, 0, 0},
		Distance: 2.0,
		Yaw: 0.0,
		Pitch: 0.0,
		OpPresVAngle: 60.0 * math.Pi / 180.0,
		Near: 0.1,
		Far: 100.1,
	}
}

func (c *OrbitCamera) Position() Vec3 {
	cp := math.Cos(c.Pitch)
	sp := math.Sin(c.Pitch)
	sy := math.Sin(c.Yaw)
	cy := math.Cos(c.Yaw)
	
	return Vec3{
		X: c.Target.X + c.Distance*cp*sy,
		Y: c.Target.Y + c.Distance*sp,
		Z: c.Target.Z + c.Distance*cp*cy,
	}
}

func (c *OrbitCamera) ViewMatrix() Mat4 {
	return LookAt(c.Position(), c.Target, Vec3{0, 1, 0})
}

func (c *OrbitCamera) ProjectionMatrix(width, height int) Mat4 {
	aspect := float64(width) / float64(height)
	return Perspective(c.OpPresVAngle, aspect, c.Near, c.Far)
}

func (c *OrbitCamera) Reset() {
	*c = NewOrbitCamera()
}

func clamp(val, minVal, maxVal float64) float64 {
	if val < minVal {
		return minVal
	}
	if val > maxVal {
		return maxVal
	}
	return val
}

func (c *OrbitCamera) MoveCamera(yaw, pitch, zoom, panV, panH int) {
	const rotStep = 0.03
	const zoomStep = 0.10
	const panStep = 0.05
	const pitchLimit = math.Pi/2 - 0.05
	
	if yaw > 0 { c.Yaw += rotStep }
	if yaw < 0 { c.Yaw -= rotStep }
	if pitch > 0 { c.Pitch += rotStep }
	if pitch < 0 { c.Pitch -= rotStep }

	c.Pitch = clamp(c.Pitch, -pitchLimit, pitchLimit)
	
	if zoom > 0 { c.Distance += zoomStep }
	if zoom < 0 { c.Distance -= zoomStep }
	
	c.Distance = clamp(c.Distance, 1, c.Far)
	
	if panV > 0 { c.Target.Y += panStep }
	if panV < 0 { c.Target.Y -= panStep }
	if panH > 0 { c.Target.X += panStep }
	if panH < 0 { c.Target.X -= panStep }
}

