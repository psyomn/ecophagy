package tinystory

// The visistor shall provide read only access to possible stories.
type Visitor struct {
	Documents []Document
}

func VisitorNew(documents []Document) *Visitor {
	return &Visitor{
		Documents: documents,
	}
}

type IndexListing struct {
	Index int
	Title string
}

func (s *Visitor) GetIndexListing() []IndexListing {
	listing := make([]IndexListing, 0, len(s.Documents))
	for index := range s.Documents {
		listing = append(listing, IndexListing{index, s.Documents[index].Title})
	}

	return listing
}
