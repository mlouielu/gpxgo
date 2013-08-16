package gpx

import (
	"math"
)

const oneDegree = 1000.0 * 10000.8 / 90.0
const earthRadius = 6371 * 1000

func ToRad(x float64) float64 {
	return x / 180. * math.Pi
}

// Haversine distance between two points.
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

func length(locs []GpxWpt, threeD bool) float64 {
	var previousLoc GpxWpt
	var res float64
	for k, v := range locs {
		if k > 0 {
			previousLoc = locs[k-1]
			var d float64
			if threeD {
				d = v.Distance3D(previousLoc)
			} else {
				d = v.Distance2D(previousLoc)
			}
			res += d
		}
	}
	return res
}

func Length2D(locs []GpxWpt) float64 {
	return length(locs, false)
}

func Length3D(locs []GpxWpt) float64 {
	return length(locs, true)
}

func CalcMaxSpeed(speeds, dists []float64) float64 {
	// TODO
	return 0.0
}

func CalcUphillDownhill(elevations []float64) (float64, float64) {
	elevsLen := len(elevations)
	if elevsLen == 0 {
		return 0.0, 0.0
	}

	smooth_elevations := make([]float64, elevsLen)

	for i, elev := range elevations {
		var currEle float64
		if 0 < i && i < elevsLen-1 {
			prevEle := elevations[i-1]
			nextEle := elevations[i+1]
			currEle = prevEle*0.3 + elev*0.4 + nextEle*0.3
		} else {
			currEle = elev
		}
		smooth_elevations[i] = currEle
	}

	var uphill float64
	var downhill float64

	for i := 1; i < len(smooth_elevations); i++ {
		d := smooth_elevations[i] - smooth_elevations[i-1]
		if d > 0.0 {
			uphill += d
		} else {
			downhill -= d
		}
	}

	return uphill, downhill
}

func distance(lat1, lon1, ele1, lat2, lon2, ele2 float64, threeD, haversine bool) float64 {

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

	return math.Sqrt(math.Pow(distance2d, 2) + math.Pow((ele1-ele2), 2))
}
func Distance2D(lat1, lon1, lat2, lon2 float64, haversine bool) float64 {
	return distance(lat1, lon1, 0.0, lat2, lon2, 0.0, false, haversine)
}

func Distance3D(lat1, lon1, ele1, lat2, lon2, ele2 float64, haversine bool) float64 {
	return distance(lat1, lon1, ele1, lat2, lon2, ele2, true, haversine)
}

func ElevationAngle(loc1, loc2 GpxWpt, radians bool) float64 {
	// TODO
	return 0.0
}