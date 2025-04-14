package domain

type Product struct {
    ID         string
    Name       string
    Price      float64
    Stock      int
}


type FilterParams struct {
    Name     string
    MinPrice float64
    MaxPrice float64
}

type PaginationParams struct {
    Page    int
    PerPage int
}