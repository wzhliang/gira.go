package gitee

// User ...
type User struct {
	ID   uint32 `json:"id"`
	Name string `json:"login"`
	Type string `json:"type"`
}

// PR ...
type PR struct {
	ID        uint32 `json:"id"`
	URL       string `json:"html_url"`
	Title     string `json:"title"`
	St        string `json:"state"`
	Number    uint32 `json:"number"`
	Mergable  bool   `json:"mergeable"`
	Assignees []User `json:"assignees"`
	Testers   []User `json:"testers"`
}
