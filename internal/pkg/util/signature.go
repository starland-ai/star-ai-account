package util
import (
   "crypto/ecdsa"
   "crypto/rand"
   "crypto/sha256"
   "crypto/x509"
   "encoding/base64"
   "encoding/pem"
   "fmt"
   "io/ioutil"
   "os"
)

func gen() {

   privateKeyFile, err := os.Open("private_key.pem")
   if err != nil {
      panic(err)
   }
   defer privateKeyFile.Close()

   privateKeyBytes, err := ioutil.ReadAll(privateKeyFile)
   if err != nil {
      panic(err)
   }

   var privateKey *ecdsa.PrivateKey
   for {
      var block *pem.Block
      block, privateKeyBytes = pem.Decode(privateKeyBytes)
      if block == nil {
         break
      }

      if block.Type == "EC PRIVATE KEY" {
         var parsedPrivateKey *ecdsa.PrivateKey
         parsedPrivateKey, err = x509.ParseECPrivateKey(block.Bytes)
         if err != nil {
            panic(err)
         }
         privateKey = parsedPrivateKey
         break
      }
   }

   if privateKey == nil {
      panic("failed to parse private key")
   }

   fmt.Println("Private key loaded successfully.")

   publicKey := privateKey.Public().(*ecdsa.PublicKey)

   message := fmt.Sprintf("%s-%d", "user", 1)

   hash := sha256.Sum256([]byte(message))

   sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
   if err != nil {
      panic(err)
   }

   base64Sig := base64.StdEncoding.EncodeToString(sig)
   fmt.Printf("Signature: (%x) -- len %d\n", base64Sig, len(base64Sig))

   sig, err = base64.StdEncoding.DecodeString(base64Sig)
   if err != nil {
      panic(err)
   }
   valid := ecdsa.VerifyASN1(publicKey, hash[:], sig)
   if valid {
      fmt.Println("Signature is valid.")
   } else {
      fmt.Println("Signature is invalid.")
   }
}