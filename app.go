package main

import (
	"context"

	"oreno-tools/backend"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ConvertUnit(size float64, unit string) (backend.ConversionResult, error) {
	return backend.ConvertUnit(size, unit)
}

func (a *App) ConvertDataTransferRate(speed float64, unit string) (backend.DataTransferRateResult, error) {
	return backend.ConvertDataTransferRate(speed, unit)
}

func (a *App) CalculateTransferTime(dataSize float64, dataUnit string, speed float64, speedUnit string) (backend.TransferTimeResult, error) {
	return backend.CalculateTransferTime(dataSize, dataUnit, speed, speedUnit)
}

func (a *App) Base64Encode(input string, urlSafe bool) string {
	return backend.Base64Encode(input, urlSafe)
}

func (a *App) Base64Decode(input string, urlSafe bool) (string, error) {
	return backend.Base64Decode(input, urlSafe)
}

func (a *App) JWTEncode(payload string, secret string, algorithm string) (string, error) {
	return backend.JWTEncode(payload, secret, algorithm)
}

func (a *App) JWTDecode(token string) (backend.JWTDecodeResult, error) {
	return backend.JWTDecode(token)
}

func (a *App) JWTVerify(tokenString string, secret string) (backend.JWTDecodeResult, error) {
	return backend.JWTVerify(tokenString, secret)
}

func (a *App) ExtractPublicKey(privateKey string, algorithm string) (string, error) {
	return backend.ExtractPublicKey(privateKey, algorithm)
}

// GenerateMTLSCertificates は mTLS 用の CA、サーバー、クライアント証明書を生成する
func (a *App) GenerateMTLSCertificates(
	caCommonName string,
	serverCommonName string,
	clientCommonName string,
	organization string,
	serverSanEntries []string,
	clientSanEntries []string,
	validDays int,
	algorithm string,
) (backend.MTLSCertificateResult, error) {
	return backend.GenerateMTLSCertificates(caCommonName, serverCommonName, clientCommonName, organization, serverSanEntries, clientSanEntries, validDays, algorithm)
}

func (a *App) VerifyKeyPair(privateKey string, publicKey string, algorithm string) (bool, error) {
	return backend.VerifyKeyPair(privateKey, publicKey, algorithm)
}

func (a *App) GeneratePrivateKey(algorithm string) (string, error) {
	return backend.GeneratePrivateKey(algorithm)
}

func (a *App) GenerateSelfSignedCertificate(commonName string, organization string, sanEntries []string, validDays int, algorithm string) (backend.SelfSignedCertificateResult, error) {
	return backend.GenerateSelfSignedCertificate(commonName, organization, sanEntries, validDays, algorithm)
}

func (a *App) GenerateMTLSCertificatesMultiClient(
	caCommonName string,
	serverCommonName string,
	clientCommonNamePrefix string,
	clientCount int,
	organization string,
	serverSanEntries []string,
	clientSanEntries []string,
	validDays int,
	algorithm string,
) (backend.MTLSCertificatesMultiClientResult, error) {
	return backend.GenerateMTLSCertificatesMultiClient(caCommonName, serverCommonName, clientCommonNamePrefix, clientCount, organization, serverSanEntries, clientSanEntries, validDays, algorithm)
}

func (a *App) SaveTextFile(content string, filename string) (string, error) {
	return backend.SaveTextFile(content, filename)
}

func (a *App) AddCertificatesToCRL(
	caCertificatePEM string,
	caPrivateKeyPEM string,
	existingCRLPEM string,
	revokedCertificatePEMs []string,
	nextUpdateDays int,
) (backend.CRLUpdateResult, error) {
	return backend.AddCertificatesToCRL(caCertificatePEM, caPrivateKeyPEM, existingCRLPEM, revokedCertificatePEMs, nextUpdateDays)
}

func (a *App) CalculateCIDR(cidr string) (backend.CIDRCalculationResult, error) {
	return backend.CalculateCIDR(cidr)
}

func (a *App) ConvertBaseValue(input string, base int, bitWidth int, signed bool) (backend.BaseConversionResult, error) {
	return backend.ConvertBaseValue(input, base, bitWidth, signed)
}

func (a *App) CalculateBaseExpression(aInput string, aBase int, operator string, bInput string, bBase int) (backend.BaseCalculationResult, error) {
	return backend.CalculateBaseExpression(aInput, aBase, operator, bInput, bBase)
}
