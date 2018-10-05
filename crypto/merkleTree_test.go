package crypto

import (
	"testing"
	."github.com/icloudland/starchain/common"
	"crypto/sha256"
	"fmt"
)

func TestMerkletree(t *testing.T){
	var data []Uint256
	data = append(data,Uint256(sha256.Sum256([]byte("k"))))
	data = append(data,Uint256(sha256.Sum256([]byte("i"))))
	data = append(data,Uint256(sha256.Sum256([]byte("l"))))
	data = append(data,Uint256(sha256.Sum256([]byte("l"))))
	res,_ := ComputeRoot(data)
	fmt.Println(res)
}