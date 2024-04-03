package account

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	config "starland-account/configs"
	"starland-account/internal/biz"
	"time"

	bin "github.com/gagliardetto/binary"

	solana "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserPoints struct {
	Authority     solana.PublicKey
	Points        uint64
	LastSignature [96]uint8
	ClaimCount    uint64
}

func (s *AccountService) Auth(ctx context.Context, req *AccountRequest) error {
	var err error
	zap.S().Infof("Auth: req:%+v", req)
	_, err = s.account.QueryAccount(ctx, req.AccountID, req.Email, req.Provider)
	if err != nil {
		zap.S().Infof("Auth: register: %s query account err: %w ", req.AccountID, err)
		accoutID := req.AccountID
		if accoutID == "" {
			accoutID = uuid.NewString()
		}
		if req.Name == "" {
			if len(req.AccountID) > 6 {
				req.Name = req.AccountID[:6]
			} else {
				req.Name = req.AccountID
			}
		}

		bizAct := &biz.AccountRequest{
			AccountID: accoutID,
			Email:     req.Email,
			Name:      req.Name,
			Provider:  req.Provider,
			AvatarURL: req.AvatarURL,
		}
		if err = s.account.SaveAccount(ctx, bizAct); err != nil {
			return fmt.Errorf("Auth: save accout err: %w ", err)
		}
	}
	return nil
}

func (s *AccountService) QueryAccount(ctx context.Context, accountID string) (*AccountResponse, error) {
	account, err := s.account.QueryAccount(ctx, accountID, "", "")
	if err != nil {
		return nil, fmt.Errorf("QueryAccount: query accout by addr(%s) err: %w", accountID, err)
	}
	return makeBizToAccountResponse(account), nil
}

func (s *AccountService) ClaimPoints(ctx context.Context, req *ClaimPointsRequest) (string, error) {
	zap.S().Infof("ClaimPoints: req: %+v", *req)
	account, err := s.QueryAccount(ctx, req.AccountID)
	if err != nil {
		return "", fmt.Errorf("ClaimPoints: query account err: %w", err)
	}
	received := account.Received + req.Points
	if account.Integral < received {
		return "", fmt.Errorf("Not enough points")
	}

	zap.S().Infof("ClaimPoints: account: %+v", *account)

	if req.IsOK {
		err = s.account.UpdateClaimPoints(ctx, req.AccountID, account.Integral, received)
		if err != nil {
			return "", fmt.Errorf("ClaimPoints: save account err: %w", err)
		}
	}

	return signature(account.AccountID, account.ClaimCount)
}

func makeBizToAccountResponse(req *biz.AccountResponse) *AccountResponse {
	return &AccountResponse{
		AccountID:  req.AccountID,
		Integral:   req.Integral,
		Received:   req.Received,
		Email:      req.Email,
		Name:       req.Name,
		Provider:   req.Provider,
		AvatarURL:  req.AvatarURL,
		SolanaAddr: req.SolanaAddr,
	}
}

func (s *AccountService) solanaChainDataCheckTask() {
	defer func() {
		if p := recover(); p != nil {
			zap.S().Errorf("")
		}
		s.solanaChainDataCheckTask()
	}()

	t := time.NewTicker(time.Second * 24)
	for range t.C {
		ar, err := s.account.QueryAccounts(context.Background())
		if err != nil {
			zap.S().Errorf("solanaChainDataCheckTask: query accounts err: %w", err)
			continue
		}
		for i := range ar {
			go s.solanaTask(ar[i])
		}
	}
}

func (s *AccountService) solanaTask(ar *biz.AccountResponse) {
	defer func() {
		if p := recover(); p != nil {
			zap.S().Infof("solanaTask: panic: %v", p)
		}
	}()
	if ar.SolanaAddr == "" {
		return
	}
	endpoint := rpc.DevNet_RPC
	if config.GetConfig().Env == "pro" {
		endpoint = rpc.MainNetBeta_RPC
	} else if config.GetConfig().Env == "test" {
		endpoint = rpc.TestNet_RPC
	}
	client := rpc.New(endpoint)
	pubKey := solana.MustPublicKeyFromBase58(ar.SolanaAddr) // serum token

	resp, err := client.GetAccountInfo(
		context.TODO(),
		pubKey,
	)
	if err != nil {
		zap.S().Error(err)
	}

	var meta UserPoints
	// Account{}.Data.GetBinary() returns the *decoded* binary data
	// regardless the original encoding (it can handle them all).
	err = bin.NewBorshDecoder(resp.GetBinary()).Decode(&meta)
	if err != nil {
		zap.S().Error(err)
		return
	}
	if int(meta.Points) == ar.Received {
		return
	} else {
		lastSignature := string(meta.LastSignature[:])
		if _, bol := signatureVerify(ar.AccountID, lastSignature, int(meta.ClaimCount)); !bol {
			err = s.account.SaveAccount(context.Background(), &biz.AccountRequest{
				AccountID:  ar.AccountID,
				Email:      ar.Email,
				Name:       ar.Name,
				Provider:   ar.Provider,
				AvatarURL:  ar.AvatarURL,
				ClaimCount: ar.ClaimCount,
				SolanaAddr: ar.SolanaAddr,
				State:      -1,
			})
			if err != nil {
				zap.S().Error(fmt.Errorf("solanaTask: save account err: %w", err))
			}
		}
		return
	}
}

func signatureVerify(user, lastSignature string, claimCount int) (string, bool) {

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

	message := fmt.Sprintf("%s-%d", user, claimCount)

	hash := sha256.Sum256([]byte(message))

	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		panic(err)
	}

	base64Sig := base64.StdEncoding.EncodeToString(sig)
	fmt.Printf("Signature: (%x) -- len %d\n", base64Sig, len(base64Sig))

	publicKey := privateKey.Public().(*ecdsa.PublicKey)

	sig, err = base64.StdEncoding.DecodeString(lastSignature)
	if err != nil {
		panic(err)
	}

	valid := ecdsa.VerifyASN1(publicKey, hash[:], sig)
	if valid {
		fmt.Println("Signature is valid.")
	} else {
		fmt.Println("Signature is invalid.")
		return "", false
	}
	return base64Sig, true
}

func signature(user string, claimCount int) (string, error) {

	privateKeyFile, err := os.Open(config.GetConfig().PrivatePath)
	if err != nil {
		return "", fmt.Errorf("signature: open privateKeyFile err: %w", err)
	}
	defer privateKeyFile.Close()

	privateKeyBytes, err := ioutil.ReadAll(privateKeyFile)
	if err != nil {
		return "", fmt.Errorf("signature: read privateKeyFile err: %w", err)
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
				return "", fmt.Errorf("signature: ParseECPrivateKey err: %w", err)
			}
			privateKey = parsedPrivateKey
			break
		}
	}

	if privateKey == nil {
		return "", fmt.Errorf("signature: privateKey is nil")
	}

	fmt.Println("Private key loaded successfully.")

	message := fmt.Sprintf("%s-%d", user, claimCount)

	hash := sha256.Sum256([]byte(message))

	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		panic(err)
	}

	base64Sig := base64.StdEncoding.EncodeToString(sig)
	fmt.Printf("Signature: (%x) -- len %d\n", base64Sig, len(base64Sig))

	return base64Sig, nil
}

func (s *AccountService) SavePointsAddr(ctx context.Context, account, addr string) error {
	if err := s.account.UpdateAddr(ctx, account, addr); err != nil {
		return fmt.Errorf("SavePointsAddr: save to db err: %w", err)
	}
	return nil
}
