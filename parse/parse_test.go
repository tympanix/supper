package parse

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testMovieTitles = []string{
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
	"X-Men Origins: Wolverine",
	"Mr. & Mrs. Smith",
	"Don't Think Twice",
	"10 Things I Hate About You",
	"Berlin, I Love You",
	"To All the Boys I've Loved Before",
	"Mamma Mia!",
}

func TestCleanNameMovieTitles(t *testing.T) {
	space := regexp.MustCompile(`[\s\.]+`)

	for _, v := range testMovieTitles {
		s := space.ReplaceAllString(v, ".")
		assert.Equal(t, v, CleanName(strings.ToLower(s)))
	}
}

func TestCleanNameWebsites(t *testing.T) {
	assert.Equal(t, "Inception", CleanName("www.example.com - Inception"))
	assert.Equal(t, "Inception", CleanName("[www.example.com].Inception"))
}

func TestIdentity(t *testing.T) {
	assert.Equal(t, "thisisatest", Identity("thìs is â tést"))
	assert.Equal(t, "vyzkousejtetentoretezec", Identity("vyzkoušejte tento řetězec"))
	assert.Equal(t, "abc123", Identity(`"?=_ä!'<b½c)#1,2...3`))
	assert.Equal(t, "这是一个测试", Identity("这是一个测试"))
}

func TestCleanName(t *testing.T) {
	assert.Equal(t, "This Is a Test", CleanName("this.is.a.test"))
	assert.Equal(t, "This Is a (Test)", CleanName("this?_=is#.a_(test)"))
	assert.Equal(t, "Abc A.B.C. Abc", CleanName("abc.A.B.C.abc"))
	assert.Equal(t, "Abc A.B.C. Abc", CleanName("abc A B C abc"))
	assert.Equal(t, "A Good Day to Die Hard", CleanName("A.Good.Day.To.Die.Hard"))
	assert.Equal(t, "This Is a Test", CleanName("This.Is.A.Test"))
}

func TestFileName(t *testing.T) {
	assert.Equal(t, "filename", Filename("/path/to/filename.extension"))
}
