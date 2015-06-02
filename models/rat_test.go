package models

import "testing"

func TestGetAverageRating(t *testing.T) {
	rat := &Rat{}
	for _, v := range []int{700, 701} {
		val, err := rat.GetAverageRating(v)
		if err != nil {
			t.Error(err)
		}

		switch v {
		case 700:
			if val != float64(4) {
				t.Error("Value not equals! Expected 4, get", val)
			}
		case 701:
			if val != float64(3.5) {
				t.Error("Value not equals! Expected 3.5, get", val)
			}

		}
	}
}

func TestGetGroupedRatings(t *testing.T) {
	rat := &Rat{}

	for _, id := range []int{700, 701} {
		groups, err := rat.GetGroupedRatings(id)

		if err != nil {
			t.Error(err)
		}

		for _, group := range groups {
			switch id {
			case 700:
				if group.Rating != float64(5) && group.Rating != float64(3) {
					t.Error("Wrong values! Expected 5 or 3, got:", group.Rating)
				}
			case 701:
				if group.Rating != float64(3.5) {
					t.Error("Wrong values! Expected 3.5, got:", group.Rating)
				}

			}
		}
	}
}

func TestGetDetailedRatings(t *testing.T) {
	rat := &Rat{}

	for _, id := range []int{700, 701} {
		ratings, err := rat.GetDetailedRatings(id)

		if err != nil {
			t.Error(err)
		}

		for _, rating := range ratings {
			switch id {
			case 700:
				if rating.Rating != float64(5) && rating.Rating != float64(3) {
					t.Error("Wrong values! Expected 5 or 3, got:", rating.Rating)
				}
			case 701:
				if rating.Rating != float64(3) && rating.Rating != float64(4) {
					t.Error("Wrong values! Expected 3 or 4, got:", rating.Rating)
				}

			}
		}
	}
}
