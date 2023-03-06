package test

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"reflect"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

var privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDYKDhr0jA+/Fvcf/Fcv7Xn8PC2lkHCUjqE2LdMXAN+4RVve4PC
hh1dllgECC4doxK5ukEUijJeCwUrzLW2hzoxgN//Z9B6+LuW6twjF4vAVy1M6Abo
R69KNp2o34tDXm+Vo+4dgcHgBrH9s2ANqY++ImV0se+y6cd4Eo6bqfop/QIDAQAB
AoGBAMJBuxri5WrlfmS2MqIoxACy3pEojeZl4aNb47bjBl0zSQFMXkgmISPnJihR
dag60mxJP42G+ObdPoNzUGa+NoN5Gn8ro8esrBrOiUBdFIl7e24xQKd6skxoelIX
1dc2SzN5vOKFLAbu7r6GCohIpKsi0Icsldgh3REWJhmvkKKhAkEA+JSNOO2kxX5g
9kbZ8rfhUYN/xWWsCRXBcewocVlEUuit7QlZZeYYhAzoF2UlX7YjbPhWKEJIGPua
vLiMg+9+iQJBAN6b6tE59CRHd3JeSfs2x4Czudeb0hdOxdR2wnecOexP/LncWrRI
HOGqqdZFkAsrFyhZ3k2+Eff3fuV7GEg+UtUCQBUtGobB/+pvJLV2PbTmo0Q9bpIT
Yj934f3hf2SAlUh21/I8fKgonOgK7W6oyDFKI+Rxl21gkCHItVrkYdwPd/kCQQCU
fWTRU9srKBDhVUv8KrpBe6GH1QT7TyxfYSivKKLqoyBtyjMm9sNtNK49pAFFseSs
oeXL7fGGeq1G3imAZzJRAkEA7rIZYUjRmmqiNTFILkFo6OZ9cYd7AyZWg2KckUQ4
pwdWR1prOW0BORYozJc3ZloSN4JZP5Lu2gCerRNGRdenDw==
-----END RSA PRIVATE KEY-----
`

func GetTestPrivateKey() (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("get private key error")
	}
	
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return privateKey, nil
	}
	
	privateKey2, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey2.(*rsa.PrivateKey), nil
}

func TestT808_0x0A00_EncodeDecode(t *testing.T) {
	privateKey, err := GetTestPrivateKey()
	if err != nil {
		t.Error(err)
	}
	
	message := protocol.T808_0x0A00{
		PublicKey: &privateKey.PublicKey,
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x0A00
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
