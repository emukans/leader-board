package model

import "testing"

func TestScoreRangeCalculation(test *testing.T) {
	below, above := calculateScoreRange(10, 10, 10)
	expectedBelow := 4
	expectedAbove := 5
	if below != expectedBelow || above != expectedAbove {
		test.Errorf("score bounds calculation failure: should return (%d, %d), but returned (%d, %d)", expectedBelow, expectedAbove, below, above)
	}

	below, above = calculateScoreRange(2, 10, 10)
	expectedBelow = 2
	expectedAbove = 7
	if below != expectedBelow || above != expectedAbove {
		test.Errorf("score bounds calculation failure: should return (%d, %d), but returned (%d, %d)", expectedBelow, expectedAbove, below, above)
	}

	below, above = calculateScoreRange(4, 10, 10)
	expectedBelow = 4
	expectedAbove = 5
	if below != expectedBelow || above != expectedAbove {
		test.Errorf("score bounds calculation failure: should return (%d, %d), but returned (%d, %d)", expectedBelow, expectedAbove, below, above)
	}

	below, above = calculateScoreRange(10, 2, 10)
	expectedBelow = 7
	expectedAbove = 2
	if below != expectedBelow || above != expectedAbove {
		test.Errorf("score bounds calculation failure: should return (%d, %d), but returned (%d, %d)", expectedBelow, expectedAbove, below, above)
	}

	below, above = calculateScoreRange(10, 4, 10)
	expectedBelow = 5
	expectedAbove = 4
	if below != expectedBelow || above != expectedAbove {
		test.Errorf("score bounds calculation failure: should return (%d, %d), but returned (%d, %d)", expectedBelow, expectedAbove, below, above)
	}

	below, above = calculateScoreRange(4, 4, 10)
	expectedBelow = 4
	expectedAbove = 4
	if below != expectedBelow || above != expectedAbove {
		test.Errorf("score bounds calculation failure: should return (%d, %d), but returned (%d, %d)", expectedBelow, expectedAbove, below, above)
	}

	below, above = calculateScoreRange(2, 2, 10)
	expectedBelow = 2
	expectedAbove = 2
	if below != expectedBelow || above != expectedAbove {
		test.Errorf("score bounds calculation failure: should return (%d, %d), but returned (%d, %d)", expectedBelow, expectedAbove, below, above)
	}

	below, above = calculateScoreRange(10, 10, 7)
	expectedBelow = 3
	expectedAbove = 3
	if below != expectedBelow || above != expectedAbove {
		test.Errorf("score bounds calculation failure: should return (%d, %d), but returned (%d, %d)", expectedBelow, expectedAbove, below, above)
	}

	below, above = calculateScoreRange(2, 10, 7)
	expectedBelow = 2
	expectedAbove = 4
	if below != expectedBelow || above != expectedAbove {
		test.Errorf("score bounds calculation failure: should return (%d, %d), but returned (%d, %d)", expectedBelow, expectedAbove, below, above)
	}

	below, above = calculateScoreRange(3, 10, 7)
	expectedBelow = 3
	expectedAbove = 3
	if below != expectedBelow || above != expectedAbove {
		test.Errorf("score bounds calculation failure: should return (%d, %d), but returned (%d, %d)", expectedBelow, expectedAbove, below, above)
	}

	below, above = calculateScoreRange(10, 2, 7)
	expectedBelow = 4
	expectedAbove = 2
	if below != expectedBelow || above != expectedAbove {
		test.Errorf("score bounds calculation failure: should return (%d, %d), but returned (%d, %d)", expectedBelow, expectedAbove, below, above)
	}

	below, above = calculateScoreRange(10, 4, 7)
	expectedBelow = 3
	expectedAbove = 3
	if below != expectedBelow || above != expectedAbove {
		test.Errorf("score bounds calculation failure: should return (%d, %d), but returned (%d, %d)", expectedBelow, expectedAbove, below, above)
	}

	below, above = calculateScoreRange(3, 3, 7)
	expectedBelow = 3
	expectedAbove = 3
	if below != expectedBelow || above != expectedAbove {
		test.Errorf("score bounds calculation failure: should return (%d, %d), but returned (%d, %d)", expectedBelow, expectedAbove, below, above)
	}

	below, above = calculateScoreRange(2, 2, 7)
	expectedBelow = 2
	expectedAbove = 2
	if below != expectedBelow || above != expectedAbove {
		test.Errorf("score bounds calculation failure: should return (%d, %d), but returned (%d, %d)", expectedBelow, expectedAbove, below, above)
	}
}
