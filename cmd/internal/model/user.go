package model

type Role string

const (
	RoleDiretoria  Role = "diretoria"
	RoleFinanceiro Role = "financeiro"
	RoleProfessor  Role = "professor"
	RoleAluno      Role = "aluno"
)

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Role         Role   `json:"role"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
