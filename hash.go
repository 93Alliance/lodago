package lodago

// Hash 哈希值
func Hash(str string) uint64 {
	bytes := String2Bytes(str)
	var hashValue uint64 = 0
	for i := 0; i < len(bytes); i++ {
		hashValue = (hashValue * 17) ^ (uint64)(bytes[i])
	}
	return hashValue
}
