package utils

type PaginationInfo struct {
	Count            int  `json:"count"`
	TotalCount       int  `json:"total_count"`
	TotalPages       int  `json:"total_pages"`
	CurrentPage      int  `json:"current_page"`
	NextPage         *int `json:"next_page"`
	PrevPage         *int `json:"prev_page"`
	IsFirstPage      bool `json:"is_first_page"`
	IsLastPage       bool `json:"is_last_page"`
	IsOutOfRangePage bool `json:"is_out_of_range_page"`
}

type PaginatedResponse struct {
	Data           interface{}    `json:"data"`
	PaginationInfo PaginationInfo `json:"pagination_info"`
}

func BuildPaginatedResponse(responseRecords []any, totalCount int64, paginationParams *PaginationParams) *PaginatedResponse {
	pageSize := paginationParams.PerPage
	pageNumber := paginationParams.PageNo

	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize != 0 || totalPages == 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Data: responseRecords,
		PaginationInfo: PaginationInfo{
			Count:            len(responseRecords),
			TotalCount:       int(totalCount),
			TotalPages:       totalPages,
			CurrentPage:      pageNumber,
			NextPage:         nextPage(pageNumber, totalPages),
			PrevPage:         prevPage(pageNumber, totalPages),
			IsFirstPage:      isFirstPage(pageNumber),
			IsLastPage:       isLastPage(pageNumber, totalPages),
			IsOutOfRangePage: isOutOfRangePage(pageNumber, totalPages),
		},
	}
}

func nextPage(pageNumber int, totalPages int) *int {
	if pageNumber >= totalPages || pageNumber < 0 {
		return nil
	}
	nextPage := pageNumber + 1
	return &nextPage
}

func prevPage(pageNumber int, totalPages int) *int {
	if pageNumber <= 1 || pageNumber > totalPages+1 {
		return nil
	}
	prevPage := pageNumber - 1
	return &prevPage
}

func isFirstPage(pageNumber int) bool {
	return pageNumber == 1
}

func isLastPage(pageNumber int, totalPages int) bool {
	return pageNumber == totalPages
}

func isOutOfRangePage(pageNumber int, totalPages int) bool {
	return pageNumber > totalPages || pageNumber < 1
}
