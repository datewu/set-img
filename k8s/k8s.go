package k8s

// Bio wrap container bio
type Bio struct {
	Name       string    `json:"name"`
	Containers []*ConBio `json:"containers"`
}

// ConBio container bio
type ConBio struct {
	Name  string `json:"name"`
	Image string `json:"img"`
	Pull  string `json:"pull"`
}

// ContainerPath ...
type ContainerPath struct {
	Ns    string `json:"namespace"`
	Kind  string `json:"kind" binding:"required"`
	Name  string `json:"name" binding:"required"`
	CName string `json:"container_name" binding:"required"`
	Img   string `json:"img" binding:"required"`
}
