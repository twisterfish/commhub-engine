package guids

import (
	"crypto/rand"
	"fmt"
	"io"
)

//////////////////////////////////////////////////////////////////////////////////////////
// Entry into package - just call it 
//////////////////////////////////////////////////////////////////////////////////////////
func GetGUID() (string) {

	uuid, err := genUUID()
	
	if err != nil {
		return "error"
	} else {
		return uuid
	}
	
}

//////////////////////////////////////////////////////////////////////////////////////////
// genUUID generates a random UUID according to RFC 4122 ( without the hyphens )
//////////////////////////////////////////////////////////////////////////////////////////
func genUUID() (string, error) {
	
	uuid := make([]byte, 16)
	n, err := io.ReadFull( rand.Reader, uuid )
	
	if n != len(uuid) || err != nil {
		return "", err
	}
	
	uuid[8] = uuid[8]&^0xc0 | 0x80 // variant bits; see section 4.1.1
	
	uuid[6] = uuid[6]&^0xf0 | 0x40 // version 4 (pseudo-random); see section 4.1.3
	
	return fmt.Sprintf("%x%x%x%x%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
	
}