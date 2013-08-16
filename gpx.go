package gpx

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"time"
)

//==========================================================

const TIMELAYOUT = "2006-01-02T15:04:05Z"

//==========================================================
type GpxTrkseg struct {
	Points []GpxWpt `xml:"trkpt"`
}

type GpxTrk struct {
	Name     string      `xml:"name"`
	Cmt      string      `xml:"cmt"`
	Desc     string      `xml:"desc"`
	Src      string      `xml:"src"`
	Links    []GpxLink   `xml:"link"`
	Number   int         `xml:"number"`
	Type     string      `xml:"type"`
	Segments []GpxTrkseg `xml:"trkseg"`
}

type GpxWpt struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	// Position info
	Ele         float64 `xml:"ele"`
	Timestamp   string  `xml:"time"`
	MagVar      string  `xml:"magvar"`
	GeoIdHeight string  `xml:"geoidheight"`
	// Description info
	Name  string    `xml:"name"`
	Cmt   string    `xml:"cmt"`
	Desc  string    `xml:"desc"`
	Src   string    `xml:"src"`
	Links []GpxLink `xml:"link"`
	Sym   string    `xml:"sym"`
	Type  string    `xml:"type"`
	// Accuracy info
	Fix          string  `xml:"fix"`
	Sat          int     `xml:"sat"`
	Hdop         float64 `xml:"hdop"`
	Vdop         float64 `xml:"vdop"`
	Pdop         float64 `xml:"pdop"`
	AgeOfGpsData float64 `xml:"ageofgpsdata"`
	DGpsId       int     `xml:"dgpsid"`
}

type GpxRte struct {
	Name        string    `xml:"name"`
	Cmt         string    `xml:"cmt"`
	Desc        string    `xml:"desc"`
	Src         string    `xml:"src"`
	Links       []GpxLink `xml:"link"`
	Number      int       `xml:"number"`
	Type        string    `xml:"type"`
	RoutePoints []GpxWpt  `xml:"rtept"`
}

type GpxLink struct {
	Url  string `xml:"href,attr"`
	Text string `xml:"text"`
	Type string `xml:"type"`
}

type GpxCopyright struct {
	Author  string `xml:"author,attr"`
	Year    string `xml:"year"`
	License string `xml:"license"`
}

type GpxEmail struct {
	Id     string `xml:"id,attr"`
	Domain string `xml:"domain,attr"`
}

type GpxPerson struct {
	Name  string   `xml:"name"`
	Email GpxEmail `xml:"email"`
	Link  GpxLink  `xml:"link"`
}

type GpxMetadata struct {
	Name      string       `xml:"name"`
	Desc      string       `xml:"desc"`
	Author    GpxPerson    `xml:"author"`
	Copyright GpxCopyright `xml:"copyright"`
	Links     []GpxLink    `xml:"link"`
	Timestamp string       `xml:"time"`
	Keywords  string       `xml:"keywords"`
	Bounds    GpxBounds    `xml:"bounds"`
}

type Gpx struct {
	Version   string      `xml:"version,attr"`
	Creator   string      `xml:"creator,attr"`
	Metadata  GpxMetadata `xml:"metadata"`
	Tracks    []GpxTrk    `xml:"trk"`
	Routes    []GpxRte    `xml:"rte"`
	Waypoints []GpxWpt    `xml:"wpt"`
}

type GpxBounds struct {
	MinLat float64 `xml:"minlat,attr"`
	MaxLat float64 `xml:"maxlat,attr"`
	MinLon float64 `xml:"minlon,attr"`
	MaxLon float64 `xml:"maxlon,attr"`
}

//==========================================================

type TimeBounds struct {
	StartTime time.Time
	EndTime   time.Time
}

type UphillDownhill struct {
	Uphill   float64
	Downhill float64
}

//==========================================================
func Parse(gpxPath string) (Gpx, error) {
	gpxFile, err := os.Open(gpxPath)

	if err != nil {
		// fmt.Println("Error opening file: ", err)
		return Gpx{}, err
	}
	defer gpxFile.Close()

	b, err := ioutil.ReadAll(gpxFile)

	if err != nil {
		// fmt.Println("Error reading file: ", err)
		return Gpx{}, err
	}
	var g Gpx
	xml.Unmarshal(b, &g)

	return g, nil
}

//==========================================================

func getTime(timestr string) time.Time {
	t, err := time.Parse(TIMELAYOUT, timestr)
	if err != nil {
		return time.Time{}
	}
	return t
}

//==========================================================
func (tb *TimeBounds) Equals(tb2 TimeBounds) bool {
	if tb.StartTime == tb2.StartTime && tb.EndTime == tb2.EndTime {
		return true
	}
	return false
}

func (tb *TimeBounds) String() string {
	return fmt.Sprintf("%+v, %+v", tb.StartTime, tb.EndTime)
}

func (b *GpxBounds) Equals(b2 GpxBounds) bool {
	if b.MinLon == b2.MinLon && b.MaxLat == b2.MaxLat && b.MinLon == b2.MinLon && b.MaxLon == b.MaxLon {
		return true
	}
	return false
}

func (b *GpxBounds) String() string {
	return fmt.Sprintf("%+v, %+v, %+v, %+v", b.MinLat, b.MaxLat, b.MinLat, b.MaxLon)
}

//==========================================================
func (g *Gpx) Length2D() float64 {
	var length2d float64
	for _, trk := range g.Tracks {
		length2d += trk.Length2D()
	}
	return length2d
}

func (g *Gpx) Length3D() float64 {
	var length3d float64
	for _, trk := range g.Tracks {
		length3d += trk.Length3D()
	}
	return length3d
}

//==========================================================

func (trk *GpxTrk) Length2D() float64 {
	var l float64
	for _, seg := range trk.Segments {
		d := seg.Length2D()
		l += d
	}
	return l
}
func (trk *GpxTrk) Length3D() float64 {
	var l float64
	for _, seg := range trk.Segments {
		d := seg.Length3D()
		l += d
	}
	return l
}

func (seg *GpxTrkseg) Length2D() float64 {
	return Length2D(seg.Points)
}

func (seg *GpxTrkseg) Length3D() float64 {
	return Length3D(seg.Points)
}

func (seg *GpxTrkseg) TimeBounds() TimeBounds {
	timeTuple := make([]time.Time, 0)

	for _, trkpt := range seg.Points {
		if trkpt.Timestamp != "" {
			if len(timeTuple) < 2 {
				timeTuple = append(timeTuple, trkpt.Time())
			} else {
				timeTuple[1] = trkpt.Time()
			}
		}
	}
	if len(timeTuple) == 2 {
		return TimeBounds{StartTime: timeTuple[0], EndTime: timeTuple[1]}
	}
	return TimeBounds{}
}

func (seg *GpxTrkseg) Bounds() GpxBounds {

	maxLat := -math.MaxFloat64
	minLat := math.MaxFloat64
	maxLon := -math.MaxFloat64
	minLon := math.MaxFloat64

	for _, trkpt := range seg.Points {
		maxLat = math.Max(trkpt.Lat, maxLat)
		minLat = math.Min(trkpt.Lat, minLat)
		maxLon = math.Max(trkpt.Lon, maxLon)
		minLon = math.Min(trkpt.Lon, minLon)
	}

	return GpxBounds{
		MaxLat: maxLat, MinLat: minLat,
		MaxLon: maxLon, MinLon: minLon,
	}
}

// Get speed at point
func (seg *GpxTrkseg) Speed(pointIdx int) float64 {
	trkptsLen := len(seg.Points)
	if pointIdx >= trkptsLen {
		pointIdx = trkptsLen - 1
	}

	point := seg.Points[pointIdx]

	var prevPt GpxWpt
	var nextPt GpxWpt

	havePrev := false
	haveNext := false
	if 0 < pointIdx && pointIdx < trkptsLen {
		prevPt = seg.Points[pointIdx-1]
		havePrev = true
	}

	if 0 < pointIdx && pointIdx < trkptsLen-1 {
		nextPt = seg.Points[pointIdx+1]
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

// Duration in seconds
func (seg *GpxTrkseg) Duration() float64 {
	trksLen := len(seg.Points)
	if trksLen == 0 {
		return 0.0
	}

	first := seg.Points[0]
	last := seg.Points[trksLen-1]

	if first.Time().Equal(last.Time()) {
		return 0.0
	}

	if last.Time().Before(first.Time()) {
		return 0.0
	}
	dur := last.Time().Sub(first.Time())

	return dur.Seconds()
}

func (seg *GpxTrkseg) Elevations() []float64 {
	elevations := make([]float64, len(seg.Points))
	for i, trkpt := range seg.Points {
		elevations[i] = trkpt.Ele
	}
	return elevations
}

// Return uphill and dowhill
func (seg *GpxTrkseg) UphillDownhill() UphillDownhill {
	if len(seg.Points) == 0 {
		return UphillDownhill{}
	}

	elevations := seg.Elevations()

	uphill, downhill := CalcUphillDownhill(elevations)

	return UphillDownhill{Uphill: uphill, Downhill: downhill}
}

//==========================================================

// Return Timestamp string as Time object
func (pt *GpxWpt) Time() time.Time {
	return getTime(pt.Timestamp)
}

// Time difference of two GpxWpts in seconds
func (pt *GpxWpt) TimeDiff(pt2 GpxWpt) float64 {
	t1 := pt.Time()
	t2 := pt2.Time()

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

func (pt *GpxWpt) SpeedBetween(pt2 GpxWpt, threeD bool) float64 {
	seconds := pt.TimeDiff(pt2)
	var distLen float64
	if threeD {
		distLen = pt.Distance3D(pt2)
	} else {
		distLen = pt.Distance2D(pt2)
	}

	return distLen / seconds
}

func (pt *GpxWpt) Distance2D(pt2 GpxWpt) float64 {
	return Distance2D(pt.Lat, pt.Lon, pt2.Lat, pt2.Lon, false)
}
func (pt *GpxWpt) Distance3D(pt2 GpxWpt) float64 {
	return Distance3D(pt.Lat, pt.Lon, pt.Ele, pt2.Lat, pt2.Lon, pt2.Ele, false)
}

//==========================================================