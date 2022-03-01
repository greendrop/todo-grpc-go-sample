package pagination

func ParsePageParam(page *int32) int {
	parsedPage := 1
	if page != nil {
		parsedPage = int(*page)
	}
	return parsedPage
}

func ParsePerPageParam(perPage *int32) int {
	parsedPerPage := 25
	if perPage != nil {
		parsedPerPage = int(*perPage)
	}
	return parsedPerPage
}
