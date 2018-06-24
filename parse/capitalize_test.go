package parse

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapitalize(t *testing.T) {
	for _, v := range []string{
		"Gone with the Wind",
		"The Shawshank Redemption",
		"The Godfather: Part II",
		"Schindler's List",
		"The Lord of the Rings: The Return of the King",
		"The Good, the Bad and the Ugly",
		"12 Angry Men",
		"Avengers: Infinity War",
		"The Lord of the Rings: The Fellowship of the Ring",
		"Star Wars: Episode V - The Empire Strikes Back",
		"One Flew Over the Cuckoo's Nest",
		"The Silence of the Lambs",
		"Léon: The Professional",
		"Se7en",
		"Star Wars: Episode IV - A New Hope",
		"City of God",
		"Life Is Beautiful",
		"Once Upon a Time in America",
		"21 and Over",
		"2001: A Space Odyssey",
		"To Kill a Mockingbird",
		"Monty Python and the Holy Grail",
		"L.A. Confidential",
		"Lock, Stock and Two Smoking Barrels",
		"Mr. Smith Goes to Washington",
		"V for Vendetta",
		"Kill Bill: Vol. 1",
		"Agents of S.H.I.E.L.D.",
		"Dr. Strangelove or: How I Learned to Stop Worrying and Love the Bomb",
	} {
		assert.Equal(t, v, Capitalize(strings.ToLower(v)))
	}
}

func TestUpperLowerRatio(t *testing.T) {
	assert.Equal(t, true, isUpper("ABC"))
	assert.Equal(t, true, isUpper("ABC123"))
	assert.Equal(t, true, isUpper("ABC_ABC:,'%/(!¤%))"))
	assert.Equal(t, true, isUpper("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))

	assert.Equal(t, false, isUpper("abc"))
	assert.Equal(t, false, isUpper("ABCDEFGHIJKLmNOPQRSTUVWXYZ"))
	assert.Equal(t, false, isUpper("123"))
}
