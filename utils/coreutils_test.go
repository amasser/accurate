package utils

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestFirstNonEmpty(t *testing.T) {
	firstElmnt := ""
	sampleMap := make(map[string]string)
	sampleMap["Third"] = "third"
	fourthElmnt := "fourth"
	winnerElmnt := FirstNonEmpty(firstElmnt, sampleMap["second"], sampleMap["Third"], fourthElmnt)
	if winnerElmnt != sampleMap["Third"] {
		t.Error("Wrong elemnt returned: ", winnerElmnt)
	}
}

func TestUUID(t *testing.T) {
	uuid := GenUUID()
	if len(uuid) == 0 {
		t.Fatalf("GenUUID error %s", uuid)
	}
}

func TestRoundByMethodUp1(t *testing.T) {
	result := Round(12.49, 1, ROUNDING_UP)
	expected := 12.5
	if result != expected {
		t.Errorf("Error rounding up: sould be %v was %v", expected, result)
	}
}

func TestRoundByMethodUp2(t *testing.T) {
	result := Round(12.21, 1, ROUNDING_UP)
	expected := 12.3
	if result != expected {
		t.Errorf("Error rounding up: sould be %v was %v", expected, result)
	}
}

func TestRoundByMethodUp3(t *testing.T) {
	result := Round(0.0701, 2, ROUNDING_UP)
	expected := 0.08
	if result != expected {
		t.Errorf("Error rounding up: sould be %v was %v", expected, result)
	}
}

func TestRoundByMethodDown1(t *testing.T) {
	result := Round(12.49, 1, ROUNDING_DOWN)
	expected := 12.4
	if result != expected {
		t.Errorf("Error rounding down: sould be %v was %v", expected, result)
	}
}

func TestRoundByMethodDown2(t *testing.T) {
	result := Round(12.21, 1, ROUNDING_DOWN)
	expected := 12.2
	if result != expected {
		t.Errorf("Error rounding up: sould be %v was %v", expected, result)
	}
}

func TestParseTimeDetectLayout(t *testing.T) {
	tmStr := "2013-12-30T15:00:01Z"
	expectedTime := time.Date(2013, 12, 30, 15, 0, 1, 0, time.UTC)
	tm, err := ParseTimeDetectLayout(tmStr, "")
	if err != nil {
		t.Error(err)
	} else if !tm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", tm, expectedTime)
	}
	_, err = ParseTimeDetectLayout(tmStr[1:], "")
	if err == nil {
		t.Errorf("Expecting error")
	}
	tmStr = "2016-04-01T02:00:00+02:00"
	expectedTime = time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC)
	tm, err = ParseTimeDetectLayout(tmStr, "")
	if err != nil {
		t.Error(err)
	} else if !tm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", tm, expectedTime)
	}
	_, err = ParseTimeDetectLayout(tmStr[1:], "")
	if err == nil {
		t.Errorf("Expecting error")
	}
	sqlTmStr := "2013-12-30 15:00:01"
	expectedTime = time.Date(2013, 12, 30, 15, 0, 1, 0, time.UTC)
	sqlTm, err := ParseTimeDetectLayout(sqlTmStr, "")
	if err != nil {
		t.Error(err)
	} else if !sqlTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", sqlTm, expectedTime)
	}
	_, err = ParseTimeDetectLayout(sqlTmStr[1:], "")
	if err == nil {
		t.Errorf("Expecting error")
	}
	unixTmStr := "1388415601"
	unixTm, err := ParseTimeDetectLayout(unixTmStr, "")
	if err != nil {
		t.Error(err)
	} else if !unixTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", unixTm, expectedTime)
	}
	_, err = ParseTimeDetectLayout(unixTmStr[1:], "")
	if err == nil {
		t.Errorf("Expecting error")
	}
	goTmStr := "2013-12-30 15:00:01 +0000 UTC"
	goTm, err := ParseTimeDetectLayout(goTmStr, "")
	if err != nil {
		t.Error(err)
	} else if !goTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", goTm, expectedTime)
	}
	_, err = ParseTimeDetectLayout(goTmStr[1:], "")
	if err == nil {
		t.Errorf("Expecting error")
	}
	goTmStr = "2013-12-30 15:00:01.000000000 +0000 UTC"
	goTm, err = ParseTimeDetectLayout(goTmStr, "")
	if err != nil {
		t.Error(err)
	} else if !goTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", goTm, expectedTime)
	}
	_, err = ParseTimeDetectLayout(goTmStr[1:], "")
	if err == nil {
		t.Errorf("Expecting error")
	}
	fsTmstampStr := "1394291049287234"
	fsTm, err := ParseTimeDetectLayout(fsTmstampStr, "")
	expectedTime = time.Date(2014, 3, 8, 15, 4, 9, 287234000, time.UTC)
	if err != nil {
		t.Error(err)
	} else if !fsTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", fsTm, expectedTime)
	}
	fsTmstampStr = "0"
	fsTm, err = ParseTimeDetectLayout(fsTmstampStr, "")
	expectedTime = time.Time{}
	if err != nil {
		t.Error(err)
	} else if !fsTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", fsTm, expectedTime)
	}
	onelineTmstampStr := "20131023215149"
	olTm, err := ParseTimeDetectLayout(onelineTmstampStr, "")
	expectedTime = time.Date(2013, 10, 23, 21, 51, 49, 0, time.UTC)
	if err != nil {
		t.Error(err)
	} else if !olTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", olTm, expectedTime)
	}
	oneSpaceTmStr := "08.04.2014 22:14:29"
	tsTm, err := ParseTimeDetectLayout(oneSpaceTmStr, "")
	expectedTime = time.Date(2014, 4, 8, 22, 14, 29, 0, time.UTC)
	if err != nil {
		t.Error(err)
	} else if !tsTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", tsTm, expectedTime)
	}
	if nowTm, err := ParseTimeDetectLayout(META_NOW, ""); err != nil {
		t.Error(err)
	} else if time.Now().Sub(nowTm) > time.Duration(10)*time.Millisecond {
		t.Errorf("Unexpected time parsed: %v", nowTm)
	}
	eamonTmStr := "31/05/2015 14:46:00"
	eamonTmS, err := ParseTimeDetectLayout(eamonTmStr, "")
	expectedTime = time.Date(2015, 5, 31, 14, 46, 0, 0, time.UTC)
	if err != nil {
		t.Error(err)
	} else if !eamonTmS.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", eamonTmS, expectedTime)
	}
	broadSoftTmStr := "20160419210007.037"
	broadTmS, err := ParseTimeDetectLayout(broadSoftTmStr, "")
	expectedTime = time.Date(2016, 4, 19, 21, 0, 7, 37000000, time.UTC)
	if err != nil {
		t.Error(err)
	} else if !broadTmS.Equal(expectedTime) {
		t.Errorf("Expecting: %v, received: %v", expectedTime, broadTmS)
	}
	astTimestamp := "2016-09-14T19:37:43.665+0000"
	expectedTime = time.Date(2016, 9, 14, 19, 37, 43, 665000000, time.UTC)
	astTMS, err := ParseTimeDetectLayout(astTimestamp, "")
	if err != nil {
		t.Error(err)
	} else if !astTMS.Equal(expectedTime) {
		t.Errorf("Expecting: %v, received: %v", expectedTime, astTMS)
	}
}

func TestParseDateUnix(t *testing.T) {
	date, err := ParseDate("1375212790")
	expected := time.Date(2013, 7, 30, 19, 33, 10, 0, time.UTC)
	if err != nil || !date.Equal(expected) {
		t.Error("error parsing date: ", expected.Sub(date))
	}
}

func TestParseDateUnlimited(t *testing.T) {
	date, err := ParseDate("*unlimited")
	if err != nil || !date.IsZero() {
		t.Error("error parsing unlimited date!: ")
	}
}

func TestParseDateEmpty(t *testing.T) {
	date, err := ParseDate("")
	if err != nil || !date.IsZero() {
		t.Error("error parsing unlimited date!: ")
	}
}

func TestParseDatePlus(t *testing.T) {
	date, err := ParseDate("+20s")
	expected := time.Now()
	if err != nil || date.Sub(expected).Seconds() > 20 || date.Sub(expected).Seconds() < 19 {
		t.Error("error parsing date: ", date.Sub(expected).Seconds())
	}
}

func TestParseDateMonthly(t *testing.T) {
	date, err := ParseDate("*monthly")
	expected := time.Now().AddDate(0, 1, 0)
	if err != nil || expected.Sub(date).Seconds() > 1 {
		t.Error("error parsing date: ", expected.Sub(date).Seconds())
	}
}

func TestParseDateRFC3339(t *testing.T) {
	date, err := ParseDate("2013-07-30T19:33:10Z")
	expected := time.Date(2013, 7, 30, 19, 33, 10, 0, time.UTC)
	if err != nil || !date.Equal(expected) {
		t.Error("error parsing date: ", expected.Sub(date))
	}
	date, err = ParseDate("2016-04-01T02:00:00+02:00")
	expected = time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC)
	if err != nil || !date.Equal(expected) {
		t.Errorf("Expecting: %v, received: %v", expected, date)
	}
}

func TestRoundDuration(t *testing.T) {
	minute := time.Minute
	result := RoundDuration(minute, 0*time.Second)
	expected := 0 * time.Second
	if result != expected {
		t.Errorf("Error rounding to minute1: expected %v was %v", expected, result)
	}
	result = RoundDuration(time.Second, 1*time.Second+500*time.Millisecond)
	expected = 2 * time.Second
	if result != expected {
		t.Errorf("Error rounding to minute1: expected %v was %v", expected, result)
	}
	result = RoundDuration(minute, 1*time.Second)
	expected = minute
	if result != expected {
		t.Errorf("Error rounding to minute2: expected %v was %v", expected, result)
	}
	result = RoundDuration(minute, 5*time.Second)
	expected = minute
	if result != expected {
		t.Errorf("Error rounding to minute3: expected %v was %v", expected, result)
	}
	result = RoundDuration(minute, minute)
	expected = minute
	if result != expected {
		t.Errorf("Error rounding to minute4: expected %v was %v", expected, result)
	}
	result = RoundDuration(minute, 90*time.Second)
	expected = 120 * time.Second
	if result != expected {
		t.Errorf("Error rounding to minute5: expected %v was %v", expected, result)
	}
	result = RoundDuration(60, 120)
	expected = 120.0
	if result != expected {
		t.Errorf("Error rounding to minute5: expected %v was %v", expected, result)
	}
}

func TestRoundAlredyHavingPrecision(t *testing.T) {
	x := 0.07
	if y := Round(x, 2, ROUNDING_UP); y != x {
		t.Error("Error rounding when already has desired precision: ", y)
	}
	if y := Round(x, 2, ROUNDING_MIDDLE); y != x {
		t.Error("Error rounding when already has desired precision: ", y)
	}
	if y := Round(x, 2, ROUNDING_DOWN); y != x {
		t.Error("Error rounding when already has desired precision: ", y)
	}
}

func TestSplitPrefix(t *testing.T) {
	a := SplitPrefix("0123456789", 1)
	if len(a) != 10 {
		t.Error("Error splitting prefix: ", a)
	}
}

func TestSplitPrefixFive(t *testing.T) {
	a := SplitPrefix("0123456789", 5)
	if len(a) != 6 {
		t.Error("Error splitting prefix: ", a)
	}
}

func TestSplitPrefixEmpty(t *testing.T) {
	a := SplitPrefix("", 1)
	if len(a) != 0 {
		t.Error("Error splitting prefix: ", a)
	}
}

func TestParseDurationWithSecs(t *testing.T) {
	durStr := "2"
	durExpected := time.Duration(2) * time.Second
	if parsed, err := ParseDurationWithSecs(durStr); err != nil {
		t.Error(err)
	} else if parsed != durExpected {
		t.Error("Parsed different than expected")
	}
	durStr = "2s"
	if parsed, err := ParseDurationWithSecs(durStr); err != nil {
		t.Error(err)
	} else if parsed != durExpected {
		t.Error("Parsed different than expected")
	}
	durStr = "2ms"
	durExpected = time.Duration(2) * time.Millisecond
	if parsed, err := ParseDurationWithSecs(durStr); err != nil {
		t.Error(err)
	} else if parsed != durExpected {
		t.Error("Parsed different than expected")
	}
	durStr = "0.002"
	durExpected = time.Duration(2) * time.Millisecond
	if parsed, err := ParseDurationWithSecs(durStr); err != nil {
		t.Error(err)
	} else if parsed != durExpected {
		t.Error("Parsed different than expected")
	}
	durStr = "1.002"
	durExpected = time.Duration(1002) * time.Millisecond
	if parsed, err := ParseDurationWithSecs(durStr); err != nil {
		t.Error(err)
	} else if parsed != durExpected {
		t.Error("Parsed different than expected")
	}
}

func TestMinDuration(t *testing.T) {
	d1, _ := time.ParseDuration("1m")
	d2, _ := time.ParseDuration("59s")
	minD1 := MinDuration(d1, d2)
	minD2 := MinDuration(d2, d1)
	if minD1 != d2 || minD2 != d2 {
		t.Error("Error getting min duration: ", minD1, minD2)
	}
}

func TestParseZeroRatingSubject(t *testing.T) {
	subj := []string{"", "*zero1s", "*zero5m", "*zero10h"}
	dur := []time.Duration{time.Second, time.Second, 5 * time.Minute, 10 * time.Hour}
	for i, s := range subj {
		if d, err := ParseZeroRatingSubject(s); err != nil || d != dur[i] {
			t.Error("Error parsing rating subject: ", s, d, err)
		}
	}
}

func TestConcatKey(t *testing.T) {
	if key := ConcatKey("a"); key != "a" {
		t.Error("Unexpected key value received: ", key)
	}
	if key := ConcatKey("a", "b"); key != fmt.Sprintf("a%sb", CONCATENATED_KEY_SEP) {
		t.Error("Unexpected key value received: ", key)
	}
	if key := ConcatKey("a", "b", "c"); key != fmt.Sprintf("a%sb%sc", CONCATENATED_KEY_SEP, CONCATENATED_KEY_SEP) {
		t.Error("Unexpected key value received: ", key)
	}
}

func TestConvertIfaceToString(t *testing.T) {
	val := interface{}("string1")
	if resVal, converted := ConvertIfaceToString(val); !converted || resVal != "string1" {
		t.Error(resVal, converted)
	}
	val = interface{}(123)
	if resVal, converted := ConvertIfaceToString(val); !converted || resVal != "123" {
		t.Error(resVal, converted)
	}
	val = interface{}([]byte("byte_val"))
	if resVal, converted := ConvertIfaceToString(val); !converted || resVal != "byte_val" {
		t.Error(resVal, converted)
	}
	val = interface{}(true)
	if resVal, converted := ConvertIfaceToString(val); !converted || resVal != "true" {
		t.Error(resVal, converted)
	}
}

func TestMandatory(t *testing.T) {
	_, err := FmtFieldWidth("", 0, "", "", true)
	if err == nil {
		t.Errorf("Failed to detect mandatory value")
	}
}

func TestMaxLen(t *testing.T) {
	result, err := FmtFieldWidth("test", 4, "", "", false)
	expected := "test"
	if err != nil || result != expected {
		t.Errorf("Expected \"test\" was \"%s\"", result)
	}
}

func TestRPadding(t *testing.T) {
	result, err := FmtFieldWidth("test", 8, "", "right", false)
	expected := "test    "
	if err != nil || result != expected {
		t.Errorf("Expected \"%s \" was \"%s\"", expected, result)
	}
}

func TestPaddingFiller(t *testing.T) {
	result, err := FmtFieldWidth("", 8, "", "right", false)
	expected := "        "
	if err != nil || result != expected {
		t.Errorf("Expected \"%s \" was \"%s\"", expected, result)
	}
}

func TestLPadding(t *testing.T) {
	result, err := FmtFieldWidth("test", 8, "", "left", false)
	expected := "    test"
	if err != nil || result != expected {
		t.Errorf("Expected \"%s \" was \"%s\"", expected, result)
	}
}

func TestZeroLPadding(t *testing.T) {
	result, err := FmtFieldWidth("test", 8, "", "zeroleft", false)
	expected := "0000test"
	if err != nil || result != expected {
		t.Errorf("Expected \"%s \" was \"%s\"", expected, result)
	}
}

func TestRStrip(t *testing.T) {
	result, err := FmtFieldWidth("test", 2, "right", "", false)
	expected := "te"
	if err != nil || result != expected {
		t.Errorf("Expected \"%s \" was \"%s\"", expected, result)
	}
}

func TestXRStrip(t *testing.T) {
	result, err := FmtFieldWidth("test", 3, "xright", "", false)
	expected := "tex"
	if err != nil || result != expected {
		t.Errorf("Expected \"%s \" was \"%s\"", expected, result)
	}
}

func TestLStrip(t *testing.T) {
	result, err := FmtFieldWidth("test", 2, "left", "", false)
	expected := "st"
	if err != nil || result != expected {
		t.Errorf("Expected \"%s \" was \"%s\"", expected, result)
	}
}

func TestXLStrip(t *testing.T) {
	result, err := FmtFieldWidth("test", 3, "xleft", "", false)
	expected := "xst"
	if err != nil || result != expected {
		t.Errorf("Expected \"%s \" was \"%s\"", expected, result)
	}
}

func TestStripNotAllowed(t *testing.T) {
	_, err := FmtFieldWidth("test", 3, "", "", false)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestPaddingNotAllowed(t *testing.T) {
	_, err := FmtFieldWidth("test", 5, "", "", false)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestCastIfToString(t *testing.T) {
	v := interface{}("somestr")
	if sOut, casts := CastIfToString(v); !casts {
		t.Error("Does not cast")
	} else if sOut != "somestr" {
		t.Errorf("Received: %+v", sOut)
	}
	v = interface{}(1)
	if sOut, casts := CastIfToString(v); !casts {
		t.Error("Does not cast")
	} else if sOut != "1" {
		t.Errorf("Received: %+v", sOut)
	}
	v = interface{}(1.2)
	if sOut, casts := CastIfToString(v); !casts {
		t.Error("Does not cast")
	} else if sOut != "1.2" {
		t.Errorf("Received: %+v", sOut)
	}
}

func TestEndOfMonth(t *testing.T) {
	eom := GetEndOfMonth(time.Date(2016, time.February, 5, 10, 1, 2, 3, time.UTC))
	expected := time.Date(2016, time.February, 29, 23, 59, 59, 0, time.UTC)
	if !eom.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, eom)
	}
	eom = GetEndOfMonth(time.Date(2015, time.February, 5, 10, 1, 2, 3, time.UTC))
	expected = time.Date(2015, time.February, 28, 23, 59, 59, 0, time.UTC)
	if !eom.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, eom)
	}
	eom = GetEndOfMonth(time.Date(2016, time.January, 31, 10, 1, 2, 3, time.UTC))
	expected = time.Date(2016, time.January, 31, 23, 59, 59, 0, time.UTC)
	if !eom.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, eom)
	}
	eom = GetEndOfMonth(time.Date(2016, time.December, 31, 10, 1, 2, 3, time.UTC))
	expected = time.Date(2016, time.December, 31, 23, 59, 59, 0, time.UTC)
	if !eom.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, eom)
	}
	eom = GetEndOfMonth(time.Date(2016, time.July, 31, 23, 59, 59, 0, time.UTC))
	expected = time.Date(2016, time.July, 31, 23, 59, 59, 0, time.UTC)
	if !eom.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, eom)
	}
}

func TestParseHierarchyPath(t *testing.T) {
	eHP := HierarchyPath([]string{"Root", "CGRateS"})
	if hp := ParseHierarchyPath("Root>CGRateS", ""); !reflect.DeepEqual(hp, eHP) {
		t.Errorf("Expecting: %+v, received: %+v", eHP, hp)
	}
	if hp := ParseHierarchyPath("/Root/CGRateS/", ""); !reflect.DeepEqual(hp, eHP) {
		t.Errorf("Expecting: %+v, received: %+v", eHP, hp)
	}
}

func TestHierarchyPathAsString(t *testing.T) {
	eStr := "/Root/CGRateS"
	hp := HierarchyPath([]string{"Root", "CGRateS"})
	if hpStr := hp.AsString("/", true); hpStr != eStr {
		t.Errorf("Expecting: %q, received: %q", eStr, hpStr)
	}
}

func TestMaskSuffix(t *testing.T) {
	dest := "+4986517174963"
	if destMasked := MaskSuffix(dest, 3); destMasked != "+4986517174***" {
		t.Error("Unexpected mask applied", destMasked)
	}
	if destMasked := MaskSuffix(dest, -1); destMasked != dest {
		t.Error("Negative maskLen should not modify destination", destMasked)
	}
	if destMasked := MaskSuffix(dest, 0); destMasked != dest {
		t.Error("Zero maskLen should not modify destination", destMasked)
	}
	if destMasked := MaskSuffix(dest, 100); destMasked != "**************" {
		t.Error("High maskLen should return complete mask", destMasked)
	}
}

func TestTimeIs0h(t *testing.T) {
	t1, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	if err != nil {
		t.Error("time parsing error")
	}
	result := TimeIs0h(t1)
	if result != false {
		t.Error("time is 0 when it's supposed to be", t1)
	}
}

func TestToJSON(t *testing.T) {
	if outNilObj := ToJSON(nil); outNilObj != "null" {
		t.Errorf("Expecting null, received: <%q>", outNilObj)
	}
}

func TestClone(t *testing.T) {
	a := 15
	var b int
	err := Clone(a, &b)
	if err != nil {
		t.Error("Cloning failed")
	}
	if b != a {
		t.Error("Expected:", a, ", received:", b)
	}
}

func TestIntPointer(t *testing.T) {
	t1 := 14
	result := IntPointer(t1)
	expected := &t1
	if *expected != *result {
		t.Error("Expected:", expected, ", received:", result)
	}
}

func TestInt64Pointer(t *testing.T) {
	var t1 int64 = 19
	result := Int64Pointer(t1)
	expected := &t1
	if *expected != *result {
		t.Error("Expected:", expected, ", received:", result)
	}
}

func TestFloat64Pointer(t *testing.T) {
	var t1 float64 = 11.5
	result := Float64Pointer(t1)
	expected := &t1
	if *expected != *result {
		t.Error("Expected:", expected, ", received:", result)
	}
}

func TestBoolPointer(t *testing.T) {
	t1 := true
	result := BoolPointer(t1)
	expected := &t1
	if *expected != *result {
		t.Error("Expected:", expected, ", received:", result)
	}
}

func TestStringSlicePointer(t *testing.T) {
	t1 := []string{"CGR", "CGR", "CGR", "CGR"}
	expected := &t1
	result := StringSlicePointer(t1)
	if *result == nil {
		t.Error("Expected:", expected, ", received: nil")
	}
}

func TestFloat64SlicePointer(t *testing.T) {
	t1 := []float64{1.2, 12.3, 123.4, 1234.5}
	expected := &t1
	result := Float64SlicePointer(t1)
	if *result == nil {
		t.Error("Expected:", expected, ", received: nil")
	}
}

func TestStringMapPointer(t *testing.T) {
	t1 := map[string]bool{"cgr1": true, "cgr2": true}
	expected := &t1
	result := StringMapPointer(t1)
	if *result == nil {
		t.Error("Expected:", expected, ", received: nil")
	}
}

func TestTimePointer(t *testing.T) {
	t1, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	if err != nil {
		t.Error("time parsing error")
	}
	result := TimePointer(t1)
	expected := &t1
	if *expected != *result {
		t.Error("Expected:", expected, ", received:", result)
	}
}
