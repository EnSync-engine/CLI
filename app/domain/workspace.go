package domain

import "time"

type Workspace struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	ParentID  string       `json:"parentId"`
	Path      string       `json:"path"`
	CreatedAt time.Time    `json:"createdAt"`
	Children  []*Workspace `json:"children,omitempty"`
}

type WorkspaceList struct {
	ResultsLength int          `json:"resultsLength"`
	Results       []*Workspace `json:"results"`
}

type CreateWorkspaceRequest struct {
	Name string `json:"name"`
}
