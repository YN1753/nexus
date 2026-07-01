package request

type SaveNodeReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	SshUser     string `json:"ssh_user"`
	AuthType    int    `json:"auth_type"`
	Password    string `json:"password"`
	PrivateKey  string `json:"privateKey"`
}

type GetNodeReq struct {
	ID uint `json:"id"`
}
