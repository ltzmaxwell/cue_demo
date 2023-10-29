package calc

// extract field type, construct `check.cue`
// parse this through reflect
func Add(x, y int8, isSet bool, prefix string)int8{
	return x + y
}