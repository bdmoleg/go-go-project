package app

type AsanaProjectsResponse struct {
	AsanaProjects []AsanaProject `json:"data"`
}

type AsanaProject struct {
	Gid          string `json:"gid"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
}

type AsanaUsersResponse struct {
	AsanaUsers []AsanaUser `json:"data"`
}

type AsanaUser struct {
	Gid          string `json:"gid"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
}
