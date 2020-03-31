package lodago

import uuid "github.com/satori/go.uuid"

// UUID 生成唯一的uuid
func UUID() string {
	uuid := uuid.NewV4()
	return uuid.String()
}
