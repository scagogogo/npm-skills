package models

// Dist 表示 NPM 包的分发信息
//
// 包含下载包所需的信息，包括下载链接、完整性校验和签名等
// 这些信息用于确保下载的包是真实、完整和未被篡改的
//
// 主要字段说明:
//   - Shasum: 包的 SHA1 校验和
//   - Tarball: 包的下载 URL
//   - Integrity: 完整性校验值，通常为 SRI 格式（子资源完整性）
//   - Signatures: 包的签名信息列表
//   - FileCount: 包中文件的数量
//   - UnpackedSize: 解压后的大小（字节）
type Dist struct {
	Shasum       string       `json:"shasum"`       // 包的 SHA1 校验和
	Tarball      string       `json:"tarball"`      // 包的下载 URL
	Integrity    string       `json:"integrity"`    // 完整性校验值
	Signatures   []*Signature `json:"signatures"`   // 签名信息列表
	FileCount    int          `json:"fileCount"`    // 包中文件的数量
	UnpackedSize int64        `json:"unpackedSize"` // 解压后的大小（字节）
	NpmSignature string       `json:"npm-signature"` // NPM 签名
}

// Signature 表示 NPM 包的签名信息
//
// NPM 使用签名来验证包的发布者身份，确保包来源可信
//
// 主要字段说明:
//   - Keyid: 签名密钥的 ID
//   - Sig: 签名内容
type Signature struct {
	Keyid string `json:"keyid"` // 签名密钥的 ID
	Sig   string `json:"sig"`   // 签名内容
}
