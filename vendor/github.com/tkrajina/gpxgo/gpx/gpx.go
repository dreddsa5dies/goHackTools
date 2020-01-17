// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

import (
	"fmt"
	"math"
	"time"
)

// Global thresholds
const (
	defaultStoppedSpeedThreshold = 1.0
	removeExtreemesTreshold      = 10
)

// ----------------------------------------------------------------------------------------------------

// GPXElementInfo implements some basic stats all common GPX elements (GPX, track and segment) must have
type GPXElementInfo interface {
	Length2D() float64
	Length3D() float64
	Bounds() GpxBounds
	MovingData() MovingData
	UphillDownhill() UphillDownhill
	TimeBounds() TimeBounds
	GetTrackPointsNo() int
}

// GetGpxElementInfo pretty prints some basic information about this GPX elements
func GetGpxElementInfo(prefix string, gpxDoc GPXElementInfo) string {
	result := ""
	result += fmt.Sprint(prefix, " Points: ", gpxDoc.GetTrackPointsNo(), "\n")
	result += fmt.Sprint(prefix, " Length 2D: ", gpxDoc.Length2D()/1000.0, "\n")
	result += fmt.Sprint(prefix, " Length 3D: ", gpxDoc.Length3D()/1000.0, "\n")

	bounds := gpxDoc.Bounds()
	result += fmt.Sprintf("%s Bounds: %f, %f, %f, %f\n", prefix, bounds.MinLatitude, bounds.MaxLatitude, bounds.MinLongitude, bounds.MaxLongitude)

	md := gpxDoc.MovingData()
	result += fmt.Sprint(prefix, " Moving time: ", md.MovingTime, "\n")
	result += fmt.Sprint(prefix, " Stopped time: ", md.StoppedTime, "\n")

	result += fmt.Sprintf("%s Max speed: %fm/s = %fkm/h\n", prefix, md.MaxSpeed, md.MaxSpeed*60*60/1000.0)

	updo := gpxDoc.UphillDownhill()
	result += fmt.Sprint(prefix, " Total uphill: ", updo.Uphill, "\n")
	result += fmt.Sprint(prefix, " Total downhill: ", updo.Downhill, "\n")

	timeBounds := gpxDoc.TimeBounds()
	result += fmt.Sprint(prefix, " Started: ", timeBounds.StartTime, "\n")
	result += fmt.Sprint(prefix, " Ended: ", timeBounds.EndTime, "\n")
	return result
}

// ----------------------------------------------------------------------------------------------------

//GPX implements one or multiple GPS tracks that can be written to and parsed
//from a gpx file
type GPX struct {
	XMLNs        string
	XmlNsXsi     string
	XmlSchemaLoc string

	Version          string
	Creator          string
	Name             string
	Description      string
	AuthorName       string
	AuthorEmail      string
	AuthorLink       string
	AuthorLinkText   string
	AuthorLinkType   string
	Copyright        string
	CopyrightYear    string
	CopyrightLicense string
	Link             string
	LinkText         string
	LinkType         string
	Time             *time.Time
	Keywords         string

	// TODO
	//Extensions []byte
	Waypoints []GPXPoint
	Routes    []GPXRoute
	Tracks    []GPXTrack
}

// ToXml converts the object to xml.
// Params are optional, you can set null to use GPXs Version and no indentation.
func (g *GPX) ToXml(params ToXmlParams) ([]byte, error) {
	return ToXml(g, params)
}

// GetGpxInfo pretty prints some basic information about this GPX, its track and segments
func (g *GPX) GetGpxInfo() string {
	result := ""
	result += fmt.Sprint("GPX name: ", g.Name, "\n")
	result += fmt.Sprint("GPX desctiption: ", g.Description, "\n")
	result += fmt.Sprint("GPX version: ", g.Version, "\n")
	result += fmt.Sprint("Author: ", g.AuthorName, "\n")
	result += fmt.Sprint("Email: ", g.AuthorEmail, "\n\n")

	result += fmt.Sprint("\nGlobal stats:", "\n")
	result += GetGpxElementInfo("", g)
	result += "\n"

	for trackNo, track := range g.Tracks {
		result += fmt.Sprintf("\nTrack #%d:\n", 1+trackNo)
		result += GetGpxElementInfo("    ", &track)
		result += "\n"
		for segmentNo, segment := range track.Segments {
			result += fmt.Sprintf("\nTrack #%d, segment #%d:\n", 1+trackNo, 1+segmentNo)
			result += GetGpxElementInfo("        ", &segment)
			result += "\n"
		}
	}
	return result
}

//GetTrackPointsNo returns the amount of track points of all tracks
func (g *GPX) GetTrackPointsNo() int {
	result := 0
	for _, track := range g.Tracks {
		result += track.GetTrackPointsNo()
	}
	return result
}

// Length2D returns the 2D length of all tracks in a Gpx.
func (g *GPX) Length2D() float64 {
	var length2d float64
	for _, trk := range g.Tracks {
		length2d += trk.Length2D()
	}
	return length2d
}

// Length3D returns the 3D length of all tracks,
func (g *GPX) Length3D() float64 {
	var length3d float64
	for _, trk := range g.Tracks {
		length3d += trk.Length3D()
	}
	return length3d
}

// TimeBounds returns the time bounds of all tacks in a Gpx.
func (g *GPX) TimeBounds() TimeBounds {
	var tbGpx TimeBounds
	for i, trk := range g.Tracks {
		tbTrk := trk.TimeBounds()
		if i == 0 {
			tbGpx = trk.TimeBounds()
		} else {
			tbGpx.EndTime = tbTrk.EndTime
		}
	}
	return tbGpx
}

// Bounds returns the bounds of all tracks in a Gpx.
func (g *GPX) Bounds() GpxBounds {
	minmax := getMaximalGpxBounds()
	for _, trk := range g.Tracks {
		bnds := trk.Bounds()
		minmax.MaxLatitude = math.Max(bnds.MaxLatitude, minmax.MaxLatitude)
		minmax.MinLatitude = math.Min(bnds.MinLatitude, minmax.MinLatitude)
		minmax.MaxLongitude = math.Max(bnds.MaxLongitude, minmax.MaxLongitude)
		minmax.MinLongitude = math.Min(bnds.MinLongitude, minmax.MinLongitude)
	}
	if minmax == getMaximalGpxBounds() {
		return GpxBounds{}
	}
	return minmax
}

//ElevationBounds returns the min and max elevation of all tracks
func (g *GPX) ElevationBounds() ElevationBounds {
	minmax := getMaximalElevationBounds()
	for _, trk := range g.Tracks {
		bnds := trk.ElevationBounds()
		minmax.MaxElevation = math.Max(bnds.MaxElevation, minmax.MaxElevation)
		minmax.MinElevation = math.Min(bnds.MinElevation, minmax.MinElevation)
	}
	return minmax
}

// MovingData returns the moving data for all tracks in a Gpx.
func (g *GPX) MovingData() MovingData {
	var (
		movingTime      float64
		stoppedTime     float64
		movingDistance  float64
		stoppedDistance float64
		maxSpeed        float64
	)

	for _, trk := range g.Tracks {
		md := trk.MovingData()
		movingTime += md.MovingTime
		stoppedTime += md.StoppedTime
		movingDistance += md.MovingDistance
		stoppedDistance += md.StoppedDistance

		if md.MaxSpeed > maxSpeed {
			maxSpeed = md.MaxSpeed
		}
	}
	return MovingData{
		MovingTime:      movingTime,
		MovingDistance:  movingDistance,
		StoppedTime:     stoppedTime,
		StoppedDistance: stoppedDistance,
		MaxSpeed:        maxSpeed,
	}
}

//ReduceTrackPoints reduces the number of track points of all tracks
func (g *GPX) ReduceTrackPoints(maxPointsNo int, minDistanceBetween float64) {
	pointsNo := g.GetTrackPointsNo()

	if pointsNo < maxPointsNo && minDistanceBetween <= 0 {
		return
	}

	length := g.Length3D()

	minDistance := math.Max(float64(minDistanceBetween), math.Ceil(length/float64(maxPointsNo)))

	for _, track := range g.Tracks {
		track.ReduceTrackPoints(minDistance)
	}
}

// SimplifyTracks does Ramer-Douglas-Peucker algorithm for
// simplification of polyline on all tracks
func (g *GPX) SimplifyTracks(maxDistance float64) {
	for _, track := range g.Tracks {
		track.SimplifyTracks(maxDistance)
	}
}

// Split splits the Gpx segment segNo in a given track trackNo at
// pointNo.
func (g *GPX) Split(trackNo, segNo, pointNo int) {
	if trackNo >= len(g.Tracks) {
		return
	}

	track := &g.Tracks[trackNo]

	track.Split(segNo, pointNo)
}

// Duration returns the duration of all tracks in a Gpx in seconds.
func (g *GPX) Duration() float64 {
	if len(g.Tracks) == 0 {
		return 0.0
	}
	var result float64
	for _, trk := range g.Tracks {
		result += trk.Duration()
	}

	return result
}

// UphillDownhill returns uphill and downhill values for all tracks in a
// Gpx.
func (g *GPX) UphillDownhill() UphillDownhill {
	if len(g.Tracks) == 0 {
		return UphillDownhill{
			Uphill:   0.0,
			Downhill: 0.0,
		}
	}

	var (
		uphill   float64
		downhill float64
	)

	for _, trk := range g.Tracks {
		updo := trk.UphillDownhill()

		uphill += updo.Uphill
		downhill += updo.Downhill
	}

	return UphillDownhill{
		Uphill:   uphill,
		Downhill: downhill,
	}
}

// HasTimes Checks if *tracks* and segments have time information. Routes and Waypoints are ignored.
func (g *GPX) HasTimes() bool {
	result := true
	for _, track := range g.Tracks {
		result = result && track.HasTimes()
	}
	return result
}

// PositionAt returns a LocationResultsPair consisting the segment index
// and the GpxWpt at a certain time.
func (g *GPX) PositionAt(t time.Time) []TrackPosition {
	results := make([]TrackPosition, 0)

	for trackNo, trk := range g.Tracks {
		locs := trk.PositionAt(t)
		if len(locs) > 0 {
			for locNo := range locs {
				locs[locNo].TrackNo = trackNo
			}
			results = append(results, locs...)
		}
	}
	return results
}

//StoppedPositions returns the positions where there was a stop
func (g *GPX) StoppedPositions() []TrackPosition {
	result := make([]TrackPosition, 0)
	for trackNo, track := range g.Tracks {
		positions := track.StoppedPositions()
		for _, position := range positions {
			position.TrackNo = trackNo
			result = append(result, position)
		}
	}
	return result
}

func (g *GPX) getDistancesFromStart(distanceBetweenPoints float64) [][][]float64 {
	result := make([][][]float64, len(g.Tracks))
	var fromStart float64
	var lastSampledPoint float64
	for trackNo, track := range g.Tracks {
		result[trackNo] = make([][]float64, len(track.Segments))
		for segmentNo, segment := range track.Segments {
			result[trackNo][segmentNo] = make([]float64, len(segment.Points))
			for pointNo, point := range segment.Points {
				if pointNo > 0 {
					fromStart += point.Distance2D(&segment.Points[pointNo-1])
				}
				if pointNo == 0 || pointNo == len(segment.Points)-1 || fromStart-lastSampledPoint > distanceBetweenPoints {
					result[trackNo][segmentNo][pointNo] = fromStart
					lastSampledPoint = fromStart
				} else {
					result[trackNo][segmentNo][pointNo] = -1
				}
			}
		}
	}
	return result
}

// GetLocationsPositionsOnTrack finds locations candidates where this location
// is on a track. Returns an array of distances from start for every given
// location. Used (for example) for positioning waypoints on the graph.
// The bigger the samples number the more granular the search will be. For
// example if samples is 100 then (cca) every 100th point will be searched.
// This is for tracks with thousands of waypoints -- computing distances for
// each and every point is slow.
func (g *GPX) GetLocationsPositionsOnTrack(samples int, locations ...Location) [][]float64 {
	length2d := g.Length2D()
	distancesFromStart := g.getDistancesFromStart(length2d / float64(samples))
	result := make([][]float64, len(locations))
	for locationNo, location := range locations {
		result[locationNo] = g.getPositionsOnTrackWithPrecomputedDistances(location, distancesFromStart, length2d)
	}
	return result
}

// GetLocationPositionsOnTrack gets position of a single point on the track.
// TODO: remove following restriction by implementing a caching option
// Always use GetLocationsPositionsOnTrack(...) for multiple points, it is
// faster.
func (g *GPX) GetLocationPositionsOnTrack(samples int, location Location) []float64 {
	return g.GetLocationsPositionsOnTrack(samples, location)[0]
}

// distancesFromStart must have the same tracks, segments and pointsNo as this track.
// if any distance in distancesFromStart is less than zero that point is ignored.
func (g *GPX) getPositionsOnTrackWithPrecomputedDistances(location Location, distancesFromStart [][][]float64, length2d float64) []float64 {
	if len(g.Tracks) == 0 {
		return []float64{}
	}

	// The point must be closer than this value in order to be a candidate location:
	minDistance := 0.01 * length2d
	pointLocations := make([]float64, 0)

	// True when we enter under the minDistance length
	nearerThanMinDistance := false

	var currentCandidate *GPXPoint
	var currentCandidateFromStart float64
	currentCandidateDistance := minDistance

	var fromStart float64
	for trackNo, track := range g.Tracks {
		for segmentNo, segment := range track.Segments {
			for pointNo, point := range segment.Points {
				fromStart = distancesFromStart[trackNo][segmentNo][pointNo]
				if fromStart >= 0 {
					distance := point.Distance2D(location)
					nearerThanMinDistance = distance < minDistance
					if nearerThanMinDistance {
						if distance < currentCandidateDistance {
							currentCandidate = &point
							currentCandidateDistance = distance
							currentCandidateFromStart = fromStart
						}
					} else {
						if currentCandidate != nil {
							pointLocations = append(pointLocations, currentCandidateFromStart)
						}
						currentCandidate = nil
						currentCandidateDistance = minDistance
					}
				}
			}
		}
	}

	if currentCandidate != nil {
		pointLocations = append(pointLocations, currentCandidateFromStart)
	}

	return pointLocations
}

//ExecuteOnAllPoints executes given function on all points
func (g *GPX) ExecuteOnAllPoints(executor func(*GPXPoint)) {
	g.ExecuteOnWaypoints(executor)
	g.ExecuteOnRoutePoints(executor)
	g.ExecuteOnTrackPoints(executor)
}

//ExecuteOnWaypoints executes given function on waypoints
func (g *GPX) ExecuteOnWaypoints(executor func(*GPXPoint)) {
	for waypointNo := range g.Waypoints {
		executor(&g.Waypoints[waypointNo])
	}
}

//ExecuteOnRoutePoints executes given function on route points
func (g *GPX) ExecuteOnRoutePoints(executor func(*GPXPoint)) {
	for _, route := range g.Routes {
		route.ExecuteOnPoints(executor)
	}
}

//ExecuteOnTrackPoints executes given function on track points
func (g *GPX) ExecuteOnTrackPoints(executor func(*GPXPoint)) {
	for _, track := range g.Tracks {
		track.ExecuteOnPoints(executor)
	}
}

//AddElevation adds elevation on all points (pointElevation = pointElevation + elevation)
func (g *GPX) AddElevation(elevation float64) {
	g.ExecuteOnAllPoints(func(point *GPXPoint) {
		fmt.Println("setting elevation if NotNull for:", point.Elevation)
		if point.Elevation.NotNull() {
			fmt.Println("setting elevation")
			point.Elevation.SetValue(point.Elevation.Value() + elevation)
		}
	})
}

//RemoveElevation removes elevation info on all points
func (g *GPX) RemoveElevation() {
	g.ExecuteOnAllPoints(func(point *GPXPoint) {
		point.Elevation.SetNull()
	})
}

//ReduceGpxToSingleTrack combines all tracks to a single track
func (g *GPX) ReduceGpxToSingleTrack() {
	if len(g.Tracks) <= 1 {
		return
	}

	firstTrack := &g.Tracks[0]
	for _, track := range g.Tracks[1:] {
		for _, segment := range track.Segments {
			firstTrack.AppendSegment(&segment)
		}
	}

	g.Tracks = []GPXTrack{*firstTrack}
}

// RemoveEmpty removes all a) segments without points and b) tracks without segments
func (g *GPX) RemoveEmpty() {
	if len(g.Tracks) == 0 {
		return
	}

	for trackNo, track := range g.Tracks {
		nonEmptySegments := make([]GPXTrackSegment, 0)
		for _, segment := range track.Segments {
			if len(segment.Points) > 0 {
				//fmt.Printf("Valid segment, because of %d points!\n", len(segment.Points))
				nonEmptySegments = append(nonEmptySegments, segment)
			}
		}
		g.Tracks[trackNo].Segments = nonEmptySegments
	}

	nonEmptyTracks := make([]GPXTrack, 0)
	for _, track := range g.Tracks {
		if len(track.Segments) > 0 {
			//fmt.Printf("Valid track, baceuse of %d segments!\n", len(track.Segments))
			nonEmptyTracks = append(nonEmptyTracks, track)
		}
	}
	g.Tracks = nonEmptyTracks
}

//SmoothHorizontal smoothes all tracks horizontally
func (g *GPX) SmoothHorizontal() {
	for trackNo := range g.Tracks {
		g.Tracks[trackNo].SmoothHorizontal()
	}
}

//SmoothVertical smoothes all tracks vertically
func (g *GPX) SmoothVertical() {
	for trackNo := range g.Tracks {
		g.Tracks[trackNo].SmoothVertical()
	}
}

//RemoveHorizontalExtremes removes horizontal extremes
func (g *GPX) RemoveHorizontalExtremes() {
	for trackNo := range g.Tracks {
		g.Tracks[trackNo].RemoveHorizontalExtremes()
	}
}

//RemoveVerticalExtremes removes vertical extremes
func (g *GPX) RemoveVerticalExtremes() {
	for trackNo := range g.Tracks {
		g.Tracks[trackNo].RemoveVerticalExtremes()
	}
}

//AddMissingTime adds missing times
func (g *GPX) AddMissingTime() {
	for trackNo := range g.Tracks {
		g.Tracks[trackNo].AddMissingTime()
	}
}

//AppendTrack adds given track
func (g *GPX) AppendTrack(t *GPXTrack) {
	g.Tracks = append(g.Tracks, *t)
}

// AppendSegment appends a segment on the end of the last track.
// If the track does not exist an empty one will be created and added.
func (g *GPX) AppendSegment(s *GPXTrackSegment) {
	if len(g.Tracks) == 0 {
		g.AppendTrack(new(GPXTrack))
	}
	g.Tracks[len(g.Tracks)-1].AppendSegment(s)
}

// AppendPoint appends a point to the end of the last segment of the last track.
// If no track or segment exists empty ones will be added.
func (g *GPX) AppendPoint(p *GPXPoint) {
	if len(g.Tracks) == 0 {
		g.AppendTrack(new(GPXTrack))
	}

	lastTrack := &g.Tracks[len(g.Tracks)-1]
	if len(lastTrack.Segments) == 0 {
		lastTrack.AppendSegment(new(GPXTrackSegment))
	}

	lastSegment := &lastTrack.Segments[len(lastTrack.Segments)-1]

	lastSegment.AppendPoint(p)
}

//AppendRoute adds a route
func (g *GPX) AppendRoute(r *GPXRoute) {
	g.Routes = append(g.Routes, *r)
}

//AppendWaypoint adds a waypoint
func (g *GPX) AppendWaypoint(w *GPXPoint) {
	g.Waypoints = append(g.Waypoints, *w)
}

// ----------------------------------------------------------------------------------------------------

//ElevationBounds contains max/min elevation
type ElevationBounds struct {
	MinElevation float64
	MaxElevation float64
}

// Equals returns true if two Bounds objects are equal
func (b ElevationBounds) Equals(b2 ElevationBounds) bool {
	return b.MinElevation == b2.MinElevation && b.MaxElevation == b2.MaxElevation
}

//String implements Stringer interface
func (b *ElevationBounds) String() string {
	return fmt.Sprintf("Max: %+v Min: %+v", b.MinElevation, b.MaxElevation)
}

// ----------------------------------------------------------------------------------------------------

//GpxBounds contains min/max latitude and longitude
type GpxBounds struct {
	MinLatitude  float64
	MaxLatitude  float64
	MinLongitude float64
	MaxLongitude float64
}

// Equals returns true if two Bounds objects are equal
func (b GpxBounds) Equals(b2 GpxBounds) bool {
	return b.MinLatitude == b2.MinLatitude && b.MaxLatitude == b2.MaxLatitude && b.MinLongitude == b2.MinLongitude && b.MaxLongitude == b2.MaxLongitude
}

//String implements Stringer interface
func (b *GpxBounds) String() string {
	return fmt.Sprintf("Max: %+v, %+v Min: %+v, %+v", b.MinLatitude, b.MinLongitude, b.MaxLatitude, b.MaxLongitude)
}

// ----------------------------------------------------------------------------------------------------

// Point represents generic point data and implements the Location interface
type Point struct {
	Latitude  float64
	Longitude float64
	Elevation NullableFloat64
}

//GetLatitude returns the latitude
func (pt *Point) GetLatitude() float64 {
	return pt.Latitude
}

//GetLongitude returns the longititude
func (pt *Point) GetLongitude() float64 {
	return pt.Longitude
}

//GetElevation returns the elevation
func (pt *Point) GetElevation() NullableFloat64 {
	return pt.Elevation
}

// Distance2D returns the 2D distance of two GpxWpts.
func (pt *Point) Distance2D(pt2 Location) float64 {
	return Distance2D(pt.GetLatitude(), pt.GetLongitude(), pt2.GetLatitude(), pt2.GetLongitude(), false)
}

// Distance3D returns the 3D distance of two GpxWpts.
func (pt *Point) Distance3D(pt2 Location) float64 {
	return Distance3D(pt.GetLatitude(), pt.GetLongitude(), pt.GetElevation(), pt2.GetLatitude(), pt2.GetLongitude(), pt2.GetElevation(), false)
}

// ----------------------------------------------------------------------------------------------------

//TimeBounds contains min/max time
type TimeBounds struct {
	StartTime time.Time
	EndTime   time.Time
}

//Equals compares to another TimeBounds struct
func (tb TimeBounds) Equals(tb2 TimeBounds) bool {
	if tb.StartTime == tb2.StartTime && tb.EndTime == tb2.EndTime {
		return true
	}
	return false
}

//String implements Stringer interface
func (tb *TimeBounds) String() string {
	return fmt.Sprintf("%+v, %+v", tb.StartTime, tb.EndTime)
}

// ----------------------------------------------------------------------------------------------------

//UphillDownhill contains uphill/downhill information
type UphillDownhill struct {
	Uphill   float64
	Downhill float64
}

//Equals compares to another UphillDownhill struct
func (ud UphillDownhill) Equals(ud2 UphillDownhill) bool {
	if ud.Uphill == ud2.Uphill && ud.Downhill == ud2.Downhill {
		return true
	}
	return false
}

// ----------------------------------------------------------------------------------------------------

// TrackPosition implements the position of a point on the track
type TrackPosition struct {
	Point
	TrackNo   int
	SegmentNo int
	PointNo   int
}

// ----------------------------------------------------------------------------------------------------

//GPXPoint represents a point of the gpx file
type GPXPoint struct {
	Point
	// TODO
	Timestamp time.Time
	// TODO: Type
	MagneticVariation string
	// TODO: Type
	GeoidHeight string
	// Description info
	Name        string
	Comment     string
	Description string
	Source      string
	// TODO
	// Links       []GpxLink
	Symbol string
	Type   string
	// Accuracy info
	TypeOfGpsFix       string
	Satellites         NullableInt
	HorizontalDilution NullableFloat64
	VerticalDilution   NullableFloat64
	PositionalDilution NullableFloat64
	AgeOfDGpsData      NullableFloat64
	DGpsId             NullableInt
}

// SpeedBetween calculates the speed between two GpxWpts.
func (pt *GPXPoint) SpeedBetween(pt2 *GPXPoint, threeD bool) float64 {
	seconds := pt.TimeDiff(pt2)
	var distLen float64
	if threeD {
		distLen = pt.Distance3D(pt2)
	} else {
		distLen = pt.Distance2D(pt2)
	}

	return distLen / seconds
}

// TimeDiff returns the time difference of two GpxWpts in seconds.
func (pt *GPXPoint) TimeDiff(pt2 *GPXPoint) float64 {
	t1 := pt.Timestamp
	t2 := pt2.Timestamp

	if t1.Equal(t2) {
		return 0.0
	}

	var delta time.Duration
	if t1.After(t2) {
		delta = t1.Sub(t2)
	} else {
		delta = t2.Sub(t1)
	}

	return delta.Seconds()
}

// MaxDilutionOfPrecision returns the dilution precision of a GpxWpt.
func (pt *GPXPoint) MaxDilutionOfPrecision() float64 {
	return math.Max(pt.HorizontalDilution.Value(), math.Max(pt.VerticalDilution.Value(), pt.PositionalDilution.Value()))
}

// ----------------------------------------------------------------------------------------------------

//GPXRoute implements a gpx route
type GPXRoute struct {
	Name        string
	Comment     string
	Description string
	Source      string
	// TODO
	//Links       []Link
	Number NullableInt
	Type   string
	// TODO
	Points []GPXPoint
}

// Length returns the length of a GPX route.
func (rte *GPXRoute) Length() float64 {
	// TODO: npe check
	points := make([]Point, len(rte.Points))
	for pointNo, point := range rte.Points {
		points[pointNo] = point.Point
	}
	return Length2D(points)
}

// Center returns the center of a GPX route.
func (rte *GPXRoute) Center() (float64, float64) {
	lenRtePts := len(rte.Points)
	if lenRtePts == 0 {
		return 0.0, 0.0
	}

	var (
		sumLat float64
		sumLon float64
	)

	for _, pt := range rte.Points {
		sumLat += pt.Latitude
		sumLon += pt.Longitude
	}

	n := float64(lenRtePts)
	return sumLat / n, sumLon / n
}

//ExecuteOnPoints executes given function on all points of the route
func (rte *GPXRoute) ExecuteOnPoints(executor func(*GPXPoint)) {
	for pointNo := range rte.Points {
		executor(&rte.Points[pointNo])
	}
}

// ----------------------------------------------------------------------------------------------------

//GPXTrackSegment represents a segment of a track
type GPXTrackSegment struct {
	Points []GPXPoint
	// TODO extensions
}

// Length2D returns the 2D length of a GPX segment.
func (seg *GPXTrackSegment) Length2D() float64 {
	// TODO: There should be a better way to do this:
	points := make([]Point, len(seg.Points))
	for pointNo, point := range seg.Points {
		points[pointNo] = point.Point
	}
	return Length2D(points)
}

// Length3D returns the 3D length of a GPX segment.
func (seg *GPXTrackSegment) Length3D() float64 {
	// TODO: There should be a better way to do this:
	points := make([]Point, len(seg.Points))
	for pointNo, point := range seg.Points {
		points[pointNo] = point.Point
	}
	return Length3D(points)
}

//GetTrackPointsNo returns the amount of points of the segment
func (seg *GPXTrackSegment) GetTrackPointsNo() int {
	return len(seg.Points)
}

// TimeBounds returns the time bounds of a GPX segment.
func (seg *GPXTrackSegment) TimeBounds() TimeBounds {
	timeTuple := make([]time.Time, 0)

	for _, trkpt := range seg.Points {
		if len(timeTuple) < 2 {
			timeTuple = append(timeTuple, trkpt.Timestamp)
		} else {
			timeTuple[1] = trkpt.Timestamp
		}
	}

	if len(timeTuple) == 2 {
		return TimeBounds{StartTime: timeTuple[0], EndTime: timeTuple[1]}
	}

	return TimeBounds{StartTime: time.Time{}, EndTime: time.Time{}}
}

// Bounds returns the bounds of a GPX segment.
func (seg *GPXTrackSegment) Bounds() GpxBounds {
	minmax := getMaximalGpxBounds()
	for _, pt := range seg.Points {
		minmax.MaxLatitude = math.Max(pt.Latitude, minmax.MaxLatitude)
		minmax.MinLatitude = math.Min(pt.Latitude, minmax.MinLatitude)
		minmax.MaxLongitude = math.Max(pt.Longitude, minmax.MaxLongitude)
		minmax.MinLongitude = math.Min(pt.Longitude, minmax.MinLongitude)
	}
	if minmax == getMaximalGpxBounds() {
		return GpxBounds{}
	}
	return minmax
}

//ElevationBounds returns the min and max elevation of the segment
func (seg *GPXTrackSegment) ElevationBounds() ElevationBounds {
	minmax := getMaximalElevationBounds()
	for _, pt := range seg.Points {
		if pt.Elevation.NotNull() {
			minmax.MaxElevation = math.Max(pt.Elevation.Value(), minmax.MaxElevation)
			minmax.MinElevation = math.Min(pt.Elevation.Value(), minmax.MinElevation)
		}
	}
	return minmax
}

//HasTimes checks if there are times
//WARNING: currently not implemented!
func (seg *GPXTrackSegment) HasTimes() bool {
	return false
	/*
	   withTimes := 0
	   for _, point := range seg.Points {
	       if point.Timestamp != nil {
	           withTimes += 1
	       }
	   }
	   return withTimes / len(seg.Points) >= 0.75
	*/
}

// Speed returns the speed at point number in a GPX segment.
func (seg *GPXTrackSegment) Speed(pointIdx int) float64 {
	trkptsLen := len(seg.Points)
	if pointIdx >= trkptsLen {
		pointIdx = trkptsLen - 1
	}

	point := seg.Points[pointIdx]

	var prevPt *GPXPoint
	var nextPt *GPXPoint

	havePrev := false
	haveNext := false
	if 0 < pointIdx && pointIdx < trkptsLen {
		prevPt = &seg.Points[pointIdx-1]
		havePrev = true
	}

	if 0 < pointIdx && pointIdx < trkptsLen-1 {
		nextPt = &seg.Points[pointIdx+1]
		haveNext = true
	}

	haveSpeed1 := false
	haveSpeed2 := false

	var speed1 float64
	var speed2 float64
	if havePrev {
		speed1 = math.Abs(point.SpeedBetween(prevPt, true))
		haveSpeed1 = true
	}
	if haveNext {
		speed2 = math.Abs(point.SpeedBetween(nextPt, true))
		haveSpeed2 = true
	}

	if haveSpeed1 && haveSpeed2 {
		return (speed1 + speed2) / 2.0
	}

	if haveSpeed1 {
		return speed1
	}
	return speed2
}

// Duration returns the duration in seconds in a GPX segment.
func (seg *GPXTrackSegment) Duration() float64 {
	trksLen := len(seg.Points)
	if trksLen == 0 {
		return 0.0
	}

	first := seg.Points[0]
	last := seg.Points[trksLen-1]

	firstTimestamp := first.Timestamp
	lastTimestamp := last.Timestamp

	if firstTimestamp.Equal(lastTimestamp) {
		return 0.0
	}

	if lastTimestamp.Before(firstTimestamp) {
		return 0.0
	}
	dur := lastTimestamp.Sub(firstTimestamp)

	return dur.Seconds()
}

// Elevations returns a slice with the elevations in a GPX segment.
func (seg *GPXTrackSegment) Elevations() []NullableFloat64 {
	elevations := make([]NullableFloat64, len(seg.Points))
	for i, trkpt := range seg.Points {
		elevations[i] = trkpt.Elevation
	}
	return elevations
}

// UphillDownhill returns uphill and dowhill in a GPX segment.
func (seg *GPXTrackSegment) UphillDownhill() UphillDownhill {
	if len(seg.Points) == 0 {
		return UphillDownhill{Uphill: 0.0, Downhill: 0.0}
	}

	elevations := seg.Elevations()

	uphill, downhill := CalcUphillDownhill(elevations)

	return UphillDownhill{Uphill: uphill, Downhill: downhill}
}

//ExecuteOnPoints executes given function on segment points
func (seg *GPXTrackSegment) ExecuteOnPoints(executor func(*GPXPoint)) {
	for pointNo := range seg.Points {
		executor(&seg.Points[pointNo])
	}
}

//ReduceTrackPoints reduces the number of track points of the segment
func (seg *GPXTrackSegment) ReduceTrackPoints(minDistance float64) {
	if minDistance <= 0 {
		return
	}

	if len(seg.Points) <= 1 {
		return
	}

	newPoints := make([]GPXPoint, 0)
	newPoints = append(newPoints, seg.Points[0])

	for _, point := range seg.Points {
		previousPoint := newPoints[len(newPoints)-1]
		if point.Distance3D(&previousPoint) >= minDistance {
			newPoints = append(newPoints, point)
		}
	}

	seg.Points = newPoints
}

// SimplifyTracks does Ramer-Douglas-Peucker algorithm for simplification of polyline
func (seg *GPXTrackSegment) SimplifyTracks(maxDistance float64) {
	seg.Points = simplifyPoints(seg.Points, maxDistance)
}

//AddElevation adds elevation on segment points (pointElevation = pointElevation + elevation)
func (seg *GPXTrackSegment) AddElevation(elevation float64) {
	for _, point := range seg.Points {
		if point.Elevation.NotNull() {
			point.Elevation.SetValue(point.Elevation.Value() + elevation)
		}
	}
}

// Split splits a GPX segment at point index pt. Point pt remains in
// first part.
func (seg *GPXTrackSegment) Split(pt int) (*GPXTrackSegment, *GPXTrackSegment) {
	pts1 := seg.Points[:pt+1]
	pts2 := seg.Points[pt+1:]

	return &GPXTrackSegment{Points: pts1}, &GPXTrackSegment{Points: pts2}
}

// Join concatenates to GPX segments.
func (seg *GPXTrackSegment) Join(seg2 *GPXTrackSegment) {
	seg.Points = append(seg.Points, seg2.Points...)
}

// PositionAt returns the GpxWpt at a given time.
func (seg *GPXTrackSegment) PositionAt(t time.Time) int {
	lenPts := len(seg.Points)
	if lenPts == 0 {
		return -1
	}
	first := seg.Points[0]
	last := seg.Points[lenPts-1]

	firstTimestamp := first.Timestamp
	lastTimestamp := last.Timestamp

	if firstTimestamp.Equal(lastTimestamp) || firstTimestamp.After(lastTimestamp) {
		return -1
	}

	for i := 0; i < len(seg.Points); i++ {
		pt := seg.Points[i]
		if t.Before(pt.Timestamp) {
			return i
		}
	}

	return -1
}

//StoppedPositions returns the positions where there was a stop
func (seg *GPXTrackSegment) StoppedPositions() []TrackPosition {
	result := make([]TrackPosition, 0)
	for pointNo, point := range seg.Points {
		if pointNo > 0 {
			previousPoint := seg.Points[pointNo-1]
			if point.SpeedBetween(&previousPoint, true) < defaultStoppedSpeedThreshold {
				var trackPos TrackPosition
				trackPos.Point = point.Point
				trackPos.PointNo = pointNo
				trackPos.SegmentNo = -1
				trackPos.TrackNo = -1
				result = append(result, trackPos)
			}
		}
	}
	return result
}

// MovingData returns the moving data of a GPX segment.
func (seg *GPXTrackSegment) MovingData() MovingData {
	var (
		movingTime      float64
		stoppedTime     float64
		movingDistance  float64
		stoppedDistance float64
	)

	speedsDistances := make([]SpeedsAndDistances, 0)

	for i := 1; i < len(seg.Points); i++ {
		prev := seg.Points[i-1]
		pt := seg.Points[i]

		dist := pt.Distance3D(&prev)

		timedelta := pt.Timestamp.Sub(prev.Timestamp)
		seconds := timedelta.Seconds()
		var speedKmh float64

		if seconds > 0 {
			speedKmh = (dist / 1000.0) / (timedelta.Seconds() / math.Pow(60, 2))
		}

		if speedKmh <= defaultStoppedSpeedThreshold {
			stoppedTime += timedelta.Seconds()
			stoppedDistance += dist
		} else {
			movingTime += timedelta.Seconds()
			movingDistance += dist

			sd := SpeedsAndDistances{dist / timedelta.Seconds(), dist}
			speedsDistances = append(speedsDistances, sd)
		}
	}

	var maxSpeed float64
	if len(speedsDistances) > 0 {
		maxSpeed = CalcMaxSpeed(speedsDistances)
		if math.IsNaN(maxSpeed) {
			maxSpeed = 0
		}
	}

	return MovingData{
		movingTime,
		stoppedTime,
		movingDistance,
		stoppedDistance,
		maxSpeed,
	}
}

//AppendPoint adds a point to the segment
func (seg *GPXTrackSegment) AppendPoint(p *GPXPoint) {
	seg.Points = append(seg.Points, *p)
}

//SmoothVertical smoothes the segment vertically
func (seg *GPXTrackSegment) SmoothVertical() {
	seg.Points = smoothVertical(seg.Points)
}

//SmoothHorizontal smoothes the segment horizontally
func (seg *GPXTrackSegment) SmoothHorizontal() {
	seg.Points = smoothHorizontal(seg.Points)
}

//RemoveVerticalExtremes removes vertical extremes from the segment
func (seg *GPXTrackSegment) RemoveVerticalExtremes() {
	if len(seg.Points) < removeExtreemesTreshold {
		return
	}

	elevationDeltaSum := 0.0
	elevationDeltaNo := 0
	for pointNo, point := range seg.Points {
		if pointNo > 0 && point.Elevation.NotNull() && seg.Points[pointNo-1].Elevation.NotNull() {
			elevationDeltaSum += math.Abs(point.Elevation.Value() - seg.Points[pointNo-1].Elevation.Value())
			elevationDeltaNo += 1
		}
	}
	avgElevationDelta := elevationDeltaSum / float64(elevationDeltaNo)
	removeElevationExtremesThreshold := avgElevationDelta * 5.0

	smoothedPoints := smoothVertical(seg.Points)
	originalPoints := seg.Points

	newPoints := make([]GPXPoint, 0)
	for pointNo, point := range originalPoints {
		smoothedPoint := smoothedPoints[pointNo]
		if 0 < pointNo && pointNo < len(originalPoints)-1 && point.Elevation.NotNull() && smoothedPoints[pointNo].Elevation.NotNull() {
			d := originalPoints[pointNo-1].Distance3D(&originalPoints[pointNo+1])
			d1 := originalPoints[pointNo].Distance3D(&originalPoints[pointNo-1])
			d2 := originalPoints[pointNo].Distance3D(&originalPoints[pointNo+1])
			if d1+d2 > d*1.5 {
				if math.Abs(point.Elevation.Value()-smoothedPoint.Elevation.Value()) < removeElevationExtremesThreshold {
					newPoints = append(newPoints, point)
				}
			} else {
				newPoints = append(newPoints, point)
			}
		} else {
			newPoints = append(newPoints, point)
		}
	}
	seg.Points = newPoints
}

//RemoveHorizontalExtremes removes horizontal extremes from the segment
func (seg *GPXTrackSegment) RemoveHorizontalExtremes() {
	// Dont't remove extreemes if segment too small
	if len(seg.Points) < removeExtreemesTreshold {
		return
	}

	var sum float64
	for pointNo, point := range seg.Points {
		if pointNo > 0 {
			sum += point.Distance2D(&seg.Points[pointNo-1])
		}
	}
	// Division by zero not a problems since this is not computed on zero-length segments:
	avgDistanceBetweenPoints := float64(sum) / float64(len(seg.Points)-1)

	remove2dExtremesThreshold := 1.75 * avgDistanceBetweenPoints

	smoothedPoints := smoothHorizontal(seg.Points)
	originalPoints := seg.Points

	newPoints := make([]GPXPoint, 0)
	for pointNo, point := range originalPoints {
		if 0 < pointNo && pointNo < len(originalPoints)-1 {
			d := originalPoints[pointNo-1].Distance2D(&originalPoints[pointNo+1])
			d1 := originalPoints[pointNo].Distance2D(&originalPoints[pointNo-1])
			d2 := originalPoints[pointNo].Distance2D(&originalPoints[pointNo+1])
			if d1+d2 > d*1.5 {
				pointMovedBy := smoothedPoints[pointNo].Distance2D(&point)
				if pointMovedBy < remove2dExtremesThreshold {
					newPoints = append(newPoints, point)
				} else {
					// Removed!
				}
			} else {
				newPoints = append(newPoints, point)
			}
		} else {
			newPoints = append(newPoints, point)
		}
	}
	seg.Points = newPoints
}

//AddMissingTime adds missing times in the segment
func (seg *GPXTrackSegment) AddMissingTime() {
	emptySegmentStart := -1
	for pointNo := range seg.Points {
		timestampEmpty := seg.Points[pointNo].Timestamp.Year() <= 1
		if timestampEmpty {
			if emptySegmentStart == -1 {
				emptySegmentStart = pointNo
			}
		} else {
			if 0 < emptySegmentStart && pointNo < len(seg.Points) {
				seg.addMissingTimeInSegment(emptySegmentStart, pointNo-1)
			}
			emptySegmentStart = -1
		}
	}
}

func (seg *GPXTrackSegment) addMissingTimeInSegment(start, end int) {
	if start <= 0 {
		return
	}
	if end >= len(seg.Points)-1 {
		return
	}
	startTime, endTime := seg.Points[start-1].Timestamp, seg.Points[end+1].Timestamp
	ratios := make([]float64, end-start+1)

	length := 0.0
	for i := start; i <= end; i++ {
		length += seg.Points[i].Point.Distance2D(&seg.Points[i-1])
		ratios[i-start] = length
	}
	length += seg.Points[end].Point.Distance2D(&seg.Points[end+1])
	for i := start; i <= end; i++ {
		ratios[i-start] = ratios[i-start] / length
	}

	for i := start; i <= end; i++ {
		d := int64(ratios[i-start] * float64(endTime.Sub(startTime).Nanoseconds()))
		seg.Points[i].Timestamp = startTime.Add(time.Duration(d))
	}
}

// ----------------------------------------------------------------------------------------------------

//GPXTrack implements a gpx track
type GPXTrack struct {
	Name        string
	Comment     string
	Description string
	Source      string
	// TODO
	//Links    []Link
	Number   NullableInt
	Type     string
	Segments []GPXTrackSegment
}

// Length2D returns the 2D length of a GPX track.
func (trk *GPXTrack) Length2D() float64 {
	var l float64
	for _, seg := range trk.Segments {
		d := seg.Length2D()
		l += d
	}
	return l
}

// Length3D returns the 3D length of a GPX track.
func (trk *GPXTrack) Length3D() float64 {
	var l float64
	for _, seg := range trk.Segments {
		d := seg.Length3D()
		l += d
	}
	return l
}

//GetTrackPointsNo returns the amount of points on the track
func (trk *GPXTrack) GetTrackPointsNo() int {
	result := 0
	for _, segment := range trk.Segments {
		result += segment.GetTrackPointsNo()
	}
	return result
}

// TimeBounds returns the time bounds of a GPX track.
func (trk *GPXTrack) TimeBounds() TimeBounds {
	var tbTrk TimeBounds

	for i, seg := range trk.Segments {
		tbSeg := seg.TimeBounds()
		if i == 0 {
			tbTrk = tbSeg
		} else {
			tbTrk.EndTime = tbSeg.EndTime
		}
	}
	return tbTrk
}

// Bounds returns the bounds of a GPX track.
func (trk *GPXTrack) Bounds() GpxBounds {
	minmax := getMaximalGpxBounds()
	for _, seg := range trk.Segments {
		bnds := seg.Bounds()
		minmax.MaxLatitude = math.Max(bnds.MaxLatitude, minmax.MaxLatitude)
		minmax.MinLatitude = math.Min(bnds.MinLatitude, minmax.MinLatitude)
		minmax.MaxLongitude = math.Max(bnds.MaxLongitude, minmax.MaxLongitude)
		minmax.MinLongitude = math.Min(bnds.MinLongitude, minmax.MinLongitude)
	}
	if minmax == getMaximalGpxBounds() {
		return GpxBounds{}
	}
	return minmax
}

//ElevationBounds returns the elevation bouds of the track
func (trk *GPXTrack) ElevationBounds() ElevationBounds {
	minmax := getMaximalElevationBounds()
	for _, seg := range trk.Segments {
		bnds := seg.ElevationBounds()
		minmax.MaxElevation = math.Max(bnds.MaxElevation, minmax.MaxElevation)
		minmax.MinElevation = math.Min(bnds.MinElevation, minmax.MinElevation)
	}
	return minmax
}

//HasTimes checks if the track has times
func (trk *GPXTrack) HasTimes() bool {
	result := true
	for _, segment := range trk.Segments {
		result = result && segment.HasTimes()
	}
	return result
}

//ReduceTrackPoints reduces the number of points on the track
func (trk *GPXTrack) ReduceTrackPoints(minDistance float64) {
	for segmentNo := range trk.Segments {
		trk.Segments[segmentNo].ReduceTrackPoints(minDistance)
	}
}

// SimplifyTracks does Ramer-Douglas-Peucker algorithm for simplification of polyline
func (trk *GPXTrack) SimplifyTracks(maxDistance float64) {
	for segmentNo := range trk.Segments {
		trk.Segments[segmentNo].SimplifyTracks(maxDistance)
	}
}

// Split splits a GPX segment at a point number ptNo in a GPX track.
func (trk *GPXTrack) Split(segNo, ptNo int) {
	lenSegs := len(trk.Segments)
	if segNo >= lenSegs {
		return
	}

	newSegs := make([]GPXTrackSegment, 0)
	for i := 0; i < lenSegs; i++ {
		seg := trk.Segments[i]

		if i == segNo && ptNo < len(seg.Points) {
			seg1, seg2 := seg.Split(ptNo)
			newSegs = append(newSegs, *seg1, *seg2)
		} else {
			newSegs = append(newSegs, seg)
		}
	}
	trk.Segments = newSegs
}

//ExecuteOnPoints executes given function on track points
func (trk *GPXTrack) ExecuteOnPoints(executor func(*GPXPoint)) {
	for segmentNo := range trk.Segments {
		trk.Segments[segmentNo].ExecuteOnPoints(executor)
	}
}

//AddElevation adds elevation on all points of the track (pointElevation = pointElevation + elevation)
func (trk *GPXTrack) AddElevation(elevation float64) {
	for segmentNo := range trk.Segments {
		trk.Segments[segmentNo].AddElevation(elevation)
	}
}

// Join joins two GPX segments in a GPX track.
func (trk *GPXTrack) Join(segNo, segNo2 int) {
	lenSegs := len(trk.Segments)
	if segNo >= lenSegs && segNo2 >= lenSegs {
		return
	}
	newSegs := make([]GPXTrackSegment, 0)
	for i := 0; i < lenSegs; i++ {
		seg := trk.Segments[i]
		if i == segNo {
			secondSeg := trk.Segments[segNo2]
			seg.Join(&secondSeg)
			newSegs = append(newSegs, seg)
		} else if i == segNo2 {
			// do nothing, its already joined
		} else {
			newSegs = append(newSegs, seg)
		}
	}
	trk.Segments = newSegs
}

// JoinNext joins a GPX segment with the next segment in the current GPX
// track.
func (trk *GPXTrack) JoinNext(segNo int) {
	trk.Join(segNo, segNo+1)
}

// MovingData returns the moving data of a GPX track.
func (trk *GPXTrack) MovingData() MovingData {
	var (
		movingTime      float64
		stoppedTime     float64
		movingDistance  float64
		stoppedDistance float64
		maxSpeed        float64
	)

	for _, seg := range trk.Segments {
		md := seg.MovingData()
		movingTime += md.MovingTime
		stoppedTime += md.StoppedTime
		movingDistance += md.MovingDistance
		stoppedDistance += md.StoppedDistance

		if md.MaxSpeed > maxSpeed {
			maxSpeed = md.MaxSpeed
		}
	}
	return MovingData{
		MovingTime:      movingTime,
		MovingDistance:  movingDistance,
		StoppedTime:     stoppedTime,
		StoppedDistance: stoppedDistance,
		MaxSpeed:        maxSpeed,
	}
}

// Duration returns the duration of a GPX track.
func (trk *GPXTrack) Duration() float64 {
	if len(trk.Segments) == 0 {
		return 0.0
	}

	var result float64
	for _, seg := range trk.Segments {
		result += seg.Duration()
	}
	return result
}

// UphillDownhill return the uphill and downhill values of a GPX track.
func (trk *GPXTrack) UphillDownhill() UphillDownhill {
	if len(trk.Segments) == 0 {
		return UphillDownhill{
			Uphill:   0,
			Downhill: 0,
		}
	}

	var (
		uphill   float64
		downhill float64
	)

	for _, seg := range trk.Segments {
		updo := seg.UphillDownhill()

		uphill += updo.Uphill
		downhill += updo.Downhill
	}

	return UphillDownhill{
		Uphill:   uphill,
		Downhill: downhill,
	}
}

// PositionAt returns a LocationResultsPair for a given time.
func (trk *GPXTrack) PositionAt(t time.Time) []TrackPosition {
	results := make([]TrackPosition, 0)

	for i := 0; i < len(trk.Segments); i++ {
		seg := trk.Segments[i]
		loc := seg.PositionAt(t)
		if loc != -1 {
			results = append(results, TrackPosition{SegmentNo: i, PointNo: loc, Point: seg.Points[loc].Point})
		}
	}
	return results
}

//StoppedPositions returns the positions where there was a stop
func (trk *GPXTrack) StoppedPositions() []TrackPosition {
	result := make([]TrackPosition, 0)
	for segmentNo, segment := range trk.Segments {
		positions := segment.StoppedPositions()
		for _, position := range positions {
			position.SegmentNo = segmentNo
			result = append(result, position)
		}
	}
	return result
}

//AppendSegment adds a segment to the track
func (trk *GPXTrack) AppendSegment(s *GPXTrackSegment) {
	trk.Segments = append(trk.Segments, *s)
}

//SmoothVertical smoothes the track vertically
func (trk *GPXTrack) SmoothVertical() {
	for segmentNo := range trk.Segments {
		trk.Segments[segmentNo].SmoothVertical()
	}
}

//SmoothHorizontal smoothes the track horizontally
func (trk *GPXTrack) SmoothHorizontal() {
	for segmentNo := range trk.Segments {
		trk.Segments[segmentNo].SmoothHorizontal()
	}
}

//RemoveVerticalExtremes removes vertical extremes
func (trk *GPXTrack) RemoveVerticalExtremes() {
	for segmentNo := range trk.Segments {
		trk.Segments[segmentNo].RemoveVerticalExtremes()
	}
}

//RemoveHorizontalExtremes removes horizontal extremes
func (trk *GPXTrack) RemoveHorizontalExtremes() {
	for segmentNo := range trk.Segments {
		trk.Segments[segmentNo].RemoveHorizontalExtremes()
	}
}

//AddMissingTime adds missing times
func (trk *GPXTrack) AddMissingTime() {
	for segmentNo := range trk.Segments {
		trk.Segments[segmentNo].AddMissingTime()
	}
}

// ----------------------------------------------------------------------------------------------------

const minLatitude, maxLatitude = -90, +90
const minLongitude, maxLongitude = -180, +180

/**
 * Useful when looking for smaller bounds
 *
 * TODO does it work is region is between 179E and 179W?
 */
func getMaximalGpxBounds() GpxBounds {
	return GpxBounds{
		MinLatitude:  maxLatitude,
		MaxLatitude:  minLatitude,
		MinLongitude: maxLongitude,
		MaxLongitude: minLongitude,
	}
}

func getMaximalElevationBounds() ElevationBounds {
	return ElevationBounds{
		MaxElevation: -math.MaxFloat64,
		MinElevation: math.MaxFloat64,
	}
}
