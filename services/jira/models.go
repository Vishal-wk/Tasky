package jira

type Project struct {
	IsLast     bool    `json:"isLast"`
	MaxResults int64   `json:"maxResults"`
	NextPage   string  `json:"nextPage"`
	Self       string  `json:"self"`
	StartAt    int64   `json:"startAt"`
	Total      int64   `json:"total"`
	Values     []Value `json:"values"`
}

type Value struct {
	AvatarUrls      AvatarUrls      `json:"avatarUrls"`
	ID              string          `json:"id"`
	Insight         Insight         `json:"insight"`
	Key             string          `json:"key"`
	Name            string          `json:"name"`
	ProjectCategory ProjectCategory `json:"projectCategory"`
	Self            string          `json:"self"`
	Simplified      bool            `json:"simplified"`
	Style           string          `json:"style"`
}

type AvatarUrls struct {
	The16X16 string `json:"16x16"`
	The24X24 string `json:"24x24"`
	The32X32 string `json:"32x32"`
	The48X48 string `json:"48x48"`
}

type Insight struct {
	LastIssueUpdateTime string `json:"lastIssueUpdateTime"`
	TotalIssueCount     int64  `json:"totalIssueCount"`
}

type ProjectCategory struct {
	Description string `json:"description"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Self        string `json:"self"`
}
