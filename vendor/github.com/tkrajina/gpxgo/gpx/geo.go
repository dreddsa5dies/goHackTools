// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

import (
	"math"
	"sort"
)

const oneDegree = 1000.0 * 10000.8 / 90.0
const earthRadius = 6371 * 1000

//ToRad converts to radial coordinates
func ToRad(x float64) float64 {
	return x / 180. * math.Pi
}

//Location implements an interface for all kinds of lat/long/elevation information
type Location interface {
	GetLatitude() float64
	GetLongitude() float64
	GetElevation() NullableFloat64
}

//MovingData contains moving data
type MovingData struct {
	MovingTime      float64
	StoppedTime     float64
	MovingDistance  float64
	StoppedDistance float64
	MaxSpeed        float64
}

//Equals compares to another MovingData struct
func (md MovingData) Equals(md2 MovingData) bool {
	return md.MovingTime == md2.MovingTime &&
		md.MovingDistance == md2.MovingDistance &&
		md.StoppedTime == md2.StoppedTime &&
		md.StoppedDistance == md2.StoppedDistance &&
		md.MaxSpeed == md.MaxSpeed
}

//SpeedsAndDistances contaings speed/distance information
type SpeedsAndDistances struct {
	Speed    float64
	Distance float64
}

// HaversineDistance returns the haversine distance between two points.
//
// Implemented from http://www.movable-type.co.uk/scripts/latlong.html
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := ToRad(lat1 - lat2)
	dLon := ToRad(lon1 - lon2)
	thisLat1 := ToRad(lat1)
	thisLat2 := ToRad(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(thisLat1)*math.Cos(thisLat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := earthRadius * c

	return d
}

func length(locs []Point, threeD bool) float64 {
	var previousLoc Point
	var res float64
	for k, v := range locs {
		if k > 0 {
			previousLoc = locs[k-1]
			var d float64
			if threeD {
				d = v.Distance3D(&previousLoc)
			} else {
				d = v.Distance2D(&previousLoc)
			}
			res += d
		}
	}
	return res
}

//Length2D calculates the lenght of given points list disregarding elevation
func Length2D(locs []Point) float64 {
	return length(locs, false)
}

//Length3D calculates the lenght of given points list including elevation distance
func Length3D(locs []Point) float64 {
	return length(locs, true)
}

//CalcMaxSpeed returns the maximum speed
func CalcMaxSpeed(speedsDistances []SpeedsAndDistances) float64 {
	lenArrs := len(speedsDistances)

	if len(speedsDistances) < 20 {
		//log.Println("Segment too small to compute speed, size: ", lenArrs)
		return 0.0
	}

	var sumDists float64
	for _, d := range speedsDistances {
		sumDists += d.Distance
	}
	avgDist := sumDists / float64(lenArrs)

	var variance float64
	for i := 0; i < len(speedsDistances); i++ {
		variance += math.Pow(speedsDistances[i].Distance-avgDist, 2)
	}
	stdDeviation := math.Sqrt(variance)

	// ignore items with distance too long
	filteredSD := make([]SpeedsAndDistances, 0)
	for i := 0; i < len(speedsDistances); i++ {
		dist := math.Abs(speedsDistances[i].Distance - avgDist)
		if dist <= stdDeviation*1.5 {
			filteredSD = append(filteredSD, speedsDistances[i])
		}
	}

	speeds := make([]float64, len(filteredSD))
	for i, sd := range filteredSD {
		speeds[i] = sd.Speed
	}

	speedsSorted := sort.Float64Slice(speeds)

	if len(speedsSorted) == 0 {
		return 0
	}

	maxIdx := int(float64(len(speedsSorted)) * 0.95)
	if maxIdx >= len(speedsSorted) {
		maxIdx = len(speedsSorted) - 1
	}
	if maxIdx < 0 {
		maxIdx = 0
	}
	return speedsSorted[maxIdx]
}

//CalcUphillDownhill calculates uphill and downhill from given elevations
func CalcUphillDownhill(elevations []NullableFloat64) (float64, float64) {
	elevsLen := len(elevations)
	if elevsLen == 0 {
		return 0.0, 0.0
	}

	smoothElevations := make([]NullableFloat64, elevsLen)

	for i, elev := range elevations {
		currEle := elev
		if 0 < i && i < elevsLen-1 {
			prevEle := elevations[i-1]
			nextEle := elevations[i+1]
			if prevEle.NotNull() && nextEle.NotNull() && elev.NotNull() {
				currEle = *NewNullableFloat64(prevEle.Value()*0.3 + elev.Value()*0.4 + nextEle.Value()*0.3)
			}
		}
		smoothElevations[i] = currEle
	}

	var uphill float64
	var downhill float64

	for i := 1; i < len(smoothElevations); i++ {
		if smoothElevations[i].NotNull() && smoothElevations[i-1].NotNull() {
			d := smoothElevations[i].Value() - smoothElevations[i-1].Value()
			if d > 0.0 {
				uphill += d
			} else {
				downhill -= d
			}
		}
	}

	return uphill, downhill
}

func distance(lat1, lon1 float64, ele1 NullableFloat64, lat2, lon2 float64, ele2 NullableFloat64, threeD, haversine bool) float64 {
	absLat := math.Abs(lat1 - lat2)
	absLon := math.Abs(lon1 - lon2)
	if haversine || absLat > 0.2 || absLon > 0.2 {
		return HaversineDistance(lat1, lon1, lat2, lon2)
	}

	coef := math.Cos(ToRad(lat1))
	x := lat1 - lat2
	y := (lon1 - lon2) * coef

	distance2d := math.Sqrt(x*x+y*y) * oneDegree

	if !threeD || ele1 == ele2 {
		return distance2d
	}

	eleDiff := 0.0
	if ele1.NotNull() && ele2.NotNull() {
		eleDiff = ele1.Value() - ele2.Value()
	}

	return math.Sqrt(math.Pow(distance2d, 2) + math.Pow(eleDiff, 2))
}

////not used currently
//func distanceBetweenLocations(loc1, loc2 Location, threeD, haversine bool) float64 {
//	lat1 := loc1.GetLatitude()
//	lon1 := loc1.GetLongitude()
//	ele1 := loc1.GetElevation()
//
//	lat2 := loc2.GetLatitude()
//	lon2 := loc2.GetLongitude()
//	ele2 := loc2.GetElevation()
//
//	return distance(lat1, lon1, ele1, lat2, lon2, ele2, threeD, haversine)
//}

//Distance2D calculates the distance of 2 geo coordinates
func Distance2D(lat1, lon1, lat2, lon2 float64, haversine bool) float64 {
	return distance(lat1, lon1, *new(NullableFloat64), lat2, lon2, *new(NullableFloat64), false, haversine)
}

//Distance3D calculates the distance of 2 geo coordinates including elevation distance
func Distance3D(lat1, lon1 float64, ele1 NullableFloat64, lat2, lon2 float64, ele2 NullableFloat64, haversine bool) float64 {
	return distance(lat1, lon1, ele1, lat2, lon2, ele2, true, haversine)
}

//ElevationAngle calculates the elevation angle (steepness) between to points
func ElevationAngle(loc1, loc2 Point, radians bool) float64 {
	if loc1.Elevation.Null() || loc2.Elevation.Null() {
		return 0.0
	}

	b := loc2.Elevation.Value() - loc1.Elevation.Value()
	a := loc2.Distance2D(&loc1)

	if a == 0.0 {
		return 0.0
	}

	angle := math.Atan(b / a)

	if radians {
		return angle
	}

	return 180 * angle / math.Pi
}

// Distance of point from a line given with two points.
func distanceFromLine(point Point, linePoint1, linePoint2 GPXPoint) float64 {
	a := linePoint1.Distance2D(&linePoint2)

	if a == 0 {
		return linePoint1.Distance2D(&point)
	}

	b := linePoint1.Distance2D(&point)
	c := linePoint2.Distance2D(&point)

	s := (a + b + c) / 2.

	return 2.0 * math.Sqrt(math.Abs(s*(s-a)*(s-b)*(s-c))) / a
}

func getLineEquationCoefficients(location1, location2 Point) (float64, float64, float64) {
	if location1.Longitude == location2.Longitude {
		// Vertical line:
		return 0.0, 1.0, -location1.Longitude
	} else {
		a := (location1.Latitude - location2.Latitude) / (location1.Longitude - location2.Longitude)
		b := location1.Latitude - location1.Longitude*a
		return 1.0, -a, -b
	}
}

func simplifyPoints(points []GPXPoint, maxDistance float64) []GPXPoint {
	if len(points) < 3 {
		return points
	}

	begin, end := points[0], points[len(points)-1]

	/*
	   Use a "normal" line just to detect the most distant point (not its real distance)
	   this is because this is faster to compute than calling distance_from_line() for
	   every point.

	   This is an approximation and may have some errors near the poles and if
	   the points are too distant, but it should be good enough for most use
	   cases...
	*/
	a, b, c := getLineEquationCoefficients(begin.Point, end.Point)

	tmpMaxDistance := -1000000000.0
	tmpMaxDistancePosition := 0
	for pointNo, point := range points {
		d := math.Abs(a*point.Latitude + b*point.Longitude + c)
		if d > tmpMaxDistance {
			tmpMaxDistance = d
			tmpMaxDistancePosition = pointNo
		}
	}

	//fmt.Println()

	//fmt.Println("tmpMaxDistancePosition=", tmpMaxDistancePosition, " len(points)=", len(points))

	realMaxDistance := distanceFromLine(points[tmpMaxDistancePosition].Point, begin, end)
	//fmt.Println("realMaxDistance=", realMaxDistance, " len(points)=", len(points))

	if realMaxDistance < maxDistance {
		return []GPXPoint{begin, end}
	}

	points1 := points[:tmpMaxDistancePosition]
	point := points[tmpMaxDistancePosition]
	points2 := points[tmpMaxDistancePosition+1:]

	//fmt.Println("before simplify: len_points=", len(points), " l_points1=", len(points1), " l_points2=", len(points2))

	points1 = simplifyPoints(points1, maxDistance)
	points2 = simplifyPoints(points2, maxDistance)

	//fmt.Println("after simplify: len_points=", len(points), " l_points1=", len(points1), " l_points2=", len(points2))

	result := append(points1, point)
	return append(result, points2...)
}

func smoothHorizontal(originalPoints []GPXPoint) []GPXPoint {
	result := make([]GPXPoint, len(originalPoints))

	for pointNo, point := range originalPoints {
		result[pointNo] = point
		if 1 <= pointNo && pointNo <= len(originalPoints)-2 {
			previousPoint := originalPoints[pointNo-1]
			nextPoint := originalPoints[pointNo+1]
			result[pointNo] = point
			result[pointNo].Latitude = previousPoint.Latitude*0.4 + point.Latitude*0.2 + nextPoint.Latitude*0.4
			result[pointNo].Longitude = previousPoint.Longitude*0.4 + point.Longitude*0.2 + nextPoint.Longitude*0.4
			//log.Println("->(%f, %f)", seg.Points[pointNo].Latitude, seg.Points[pointNo].Longitude)
		}
	}

	return result
}

func smoothVertical(originalPoints []GPXPoint) []GPXPoint {
	result := make([]GPXPoint, len(originalPoints))

	for pointNo, point := range originalPoints {
		result[pointNo] = point
		if 1 <= pointNo && pointNo <= len(originalPoints)-2 {
			previousPointElevation := originalPoints[pointNo-1].Elevation
			nextPointElevation := originalPoints[pointNo+1].Elevation
			if previousPointElevation.NotNull() && point.Elevation.NotNull() && nextPointElevation.NotNull() {
				result[pointNo].Elevation = *NewNullableFloat64(previousPointElevation.Value()*0.4 + point.Elevation.Value()*0.2 + nextPointElevation.Value()*0.4)
				//log.Println("->%f", seg.Points[pointNo].Elevation.Value())
			}
		}
	}

	return result
}
