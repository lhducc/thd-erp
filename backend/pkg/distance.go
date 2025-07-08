package utils

import "math"

func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const EarthRadius = 6371000 //metter
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return EarthRadius * c
}

func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

func IsWithinScope(lat1, lon1, lat2, lon2 float64, scopeMeters int) bool {
	const EarthRadius = 6371000

	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := EarthRadius * c

	return distance <= float64(scopeMeters)
}
