package buffer

import "testing"

const testString = `Here's some...
NewJeans for testing UTF-8:

누가 내게 뭐라든
남들과는 달라 넌
Maybe you could be the one

Hype boy 내가 전해
`

func TestFromString(t *testing.T) {
	b := FromString(testString)
	content := b.String()
	if content != testString {
		t.Fatalf("content does not match:\n\n%v\n", content)
	}
}

func helperTestIndexing(t *testing.T, b *Buffer, expected string) {
	// We cannot use the index in the range because we need the index of the
	// rune and not the byte index.
	i := 0
	for _, c := range expected {
		bc, err := b.Get(i)
		if err != nil {
			t.Fatalf("erroed on get: %v", err)
		}
		if bc != c {
			t.Fatalf(
				"chars are different at %v: %v (buffer) %v (string)",
				i,
				string(bc),
				string(c),
			)
		}
		i++
	}
}

func TestIndexingVirginBuffer(t *testing.T) {
	b := FromString(testString)
	helperTestIndexing(t, b, testString)
}

func helperInsertBeggining(b *Buffer, insert string) {
	i := 0
	for _, c := range insert {
		b.Insert(i, c)
		i++
	}
}

func helperTestContent(t *testing.T, b *Buffer, expected string, pieces int, edits int) {
	content := b.String()
	if content != expected {
		t.Fatalf(
			"content does not match:\n\n%v\n\nexpected:\n\n%v\n",
			content,
			expected,
		)
	}

	if len(b.pieces) != pieces {
		t.Fatalf(
			"buffer does not have %v pieces, got %v",
			pieces,
			len(b.pieces),
		)
	}

	if len(b.edits) != edits {
		t.Fatalf("buffer does not have %v edits, got %v", edits, len(b.edits))
	}
}

func TestInsertionBeggining(t *testing.T) {
	insert := "빠져버리는 daydream\n"

	b := FromString(testString)
	helperInsertBeggining(b, insert)

	expected := insert + testString
	helperTestContent(t, b, expected, 2, 1)
	helperTestIndexing(t, b, expected)
}

func helperInsertEnd(b *Buffer, insert string) {
	for _, c := range insert {
		b.Insert(b.Size(), c)
	}
}

func TestInsertionEnd(t *testing.T) {
	b := FromString(testString)
	insert := "빠져버리는 daydream\n"

	helperInsertEnd(b, insert)

	expected := testString + insert
	helperTestContent(t, b, expected, 2, 1)
	helperTestIndexing(t, b, expected)
}

func helperInsertMiddle(b *Buffer, insert string, i int) {
	for _, c := range insert {
		b.Insert(i, c)
		i++
	}
}

func TestInsertionMiddle(t *testing.T) {
	b := FromString(testString)
	insert := "빠져버리는 daydream\n"
	left := testString[:11]
	right := testString[11:]

	helperInsertMiddle(b, insert, 11)

	expected := left + insert + right
	helperTestContent(t, b, expected, 3, 1)
	helperTestIndexing(t, b, expected)
}

func TestInsertionsBeginAndEnd(t *testing.T) {
	b := FromString(testString)

	insertBegin := "빠져버리는"
	insertEnd := "daydream\n"

	helperInsertBeggining(b, insertBegin)
	helperInsertEnd(b, insertEnd)

	expected := insertBegin + testString + insertEnd
	helperTestContent(t, b, expected, 3, 2)
	helperTestIndexing(t, b, expected)
}

func TestTwoInsertionsMiddle(t *testing.T) {
	b := FromString(testString)

	insertAt11 := "빠져버리는"
	insertAt47 := "daydream\n"

	helperInsertMiddle(b, insertAt11, 11)
	helperInsertMiddle(b, insertAt47, 47)

	left := testString[:11]
	middle := testString[11:42]
	right := testString[42:]
	expected := left + insertAt11 + middle + insertAt47 + right

	helperTestContent(t, b, expected, 5, 2)
	helperTestIndexing(t, b, expected)
}

func TestThreeInsertions(t *testing.T) {
	// We already tested UTF-8 enough and I'm lazy.
	b := FromString("hello")

	insertAt1 := "123" // "h123ello"
	insertAt6 := "ABC" // "h123elABClo"
	insertAt5 := "!@#" // "h123e!@#lABClo"

	helperInsertMiddle(b, insertAt1, 1)
	helperInsertMiddle(b, insertAt6, 6)
	helperInsertMiddle(b, insertAt5, 5)

	expected := "h123e!@#lABClo"
	helperTestContent(t, b, expected, 7, 3)
	helperTestIndexing(t, b, expected)
}

func TestThreeInsertionsWithPosAppending(t *testing.T) {
	b := FromString("hello")
	insertAt1 := "123" // "h123ello"
	insertAt6 := "ABC" // "h123elABClo"
	insertAt4 := "!@#" // "h123!@#elABClo"

	helperInsertMiddle(b, insertAt1, 1)
	helperInsertMiddle(b, insertAt6, 6)
	helperInsertMiddle(b, insertAt4, 4)

	expected := "h123!@#elABClo"
	helperTestContent(t, b, expected, 6, 3)
	helperTestIndexing(t, b, expected)
}

func TestSplitLastInsertion(t *testing.T) {
	b := FromString("hello")
	insertAt3 := "1234" // "hel1234lo"
	insertAt5 := "ABC"  // "hel12ABC34lo"

	helperInsertMiddle(b, insertAt3, 3)
	helperInsertMiddle(b, insertAt5, 5)

	expected := "hel12ABC34lo"
	helperTestContent(t, b, expected, 5, 2)
	helperTestIndexing(t, b, expected)
}
