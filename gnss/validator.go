package gnss

import (
        //"strconv"
        //"strings"
        "regexp"
)

var coordinate = regexp.MustCompile(`^-?\d{1,3}(?:\.\d+)?$`)


func IsDottedDecimal(lat, lon string) bool {
        return isDottedDecimal(lat) && isDottedDecimal(lon)
}

func isDottedDecimal(s string) bool {
        //if s == "" {
        //      return false
        //}

        // when no position data is coming (unlikely to be exactly there :) )
        if s == "0.0" {
                return false
        }


        return coordinate.MatchString(s)

        // must contain a decimal point
        //if !strings.Contains(s, ".") {
        //      return false
        //}

        // must be a valid number
        //_, err := strconv.ParseFloat(s, 64)
        //return err == nil
}
