package manager

import (
	"fmt"
	"path/filepath"

	"github.com/lumosolutions/yerd/internal/constants"
	"github.com/lumosolutions/yerd/internal/utils"
)

type CertificateManager struct {
}

func NewCertificateManager() *CertificateManager {
	return &CertificateManager{}
}

func (certManager *CertificateManager) GenerateCaCertificate(name string) error {
	caPath := filepath.Join(constants.CertsDir, "ca")
	keyName := name + ".key"
	certName := name + ".crt"

	utils.CreateDirectory(caPath)

	if !certManager.generateKey(caPath, keyName, "4096") {
		return fmt.Errorf("unable to generate ca key")
	}

	params := []string{"req", "-new", "-x509", "-days", "3650", "-key", keyName, "-out", certName, "-subj", "/C=US/CN=YERD/O=YERD/OU=YERD"}
	output, success := utils.ExecuteCommandInDir(caPath, "openssl", params...)
	if !success {
		utils.LogInfo("cacert", "openssl command failed")
		utils.LogInfo("cacert", "output: %s", output)
		return fmt.Errorf("failed to generate ca cert")
	}

	depMan, _ := NewDependencyManager()
	depMan.TrustCertificate(filepath.Join(caPath, certName), name)

	params = []string{"-A", "-n", "YERD CA", "-t", "TCu,Cu,Tu", "-i", filepath.Join(caPath, certName), "-d", "sql:$HOME/.pki/nssdb"}
	utils.ExecuteCommand("certutil", params...)

	return nil
}

func (certManager *CertificateManager) GenerateCert(domain, caName string) (string, string, error) {
	certPath := filepath.Join(constants.CertsDir, "sites")
	keyName := domain + ".key"
	csrName := domain + ".csr"
	certName := domain + ".crt"

	caCert := filepath.Join(constants.CertsDir, "ca", caName+".crt")
	caKey := filepath.Join(constants.CertsDir, "ca", caName+".key")

	utils.CreateDirectory(certPath)

	if !certManager.generateKey(certPath, keyName, "2048") {
		return "", "", fmt.Errorf("unable to generate site key")
	}

	if !certManager.generateCsr(certPath, keyName, csrName, "/C=GB/CN="+domain) {
		return "", "", fmt.Errorf("unable to generate site csr")
	}

	if !certManager.generateSiteCertificate(certPath, domain, csrName, certName, caCert, caKey) {
		return "", "", fmt.Errorf("unable to generate site cert")
	}

	return filepath.Join(certPath, keyName), filepath.Join(certPath, certName), nil
}

func (certManager *CertificateManager) generateKey(folder, keyFile, keySize string) bool {
	output, success := utils.ExecuteCommandInDir(folder, "openssl", "genrsa", "-out", keyFile, keySize)
	if !success {
		utils.LogInfo("generateKey", "openssl command failed")
		utils.LogInfo("generateKey", "output: %s", output)
		return false
	}

	return true
}

func (certManager *CertificateManager) generateCsr(folder, keyFile, csrFile, subject string) bool {
	params := []string{"req", "-new", "-key", keyFile, "-out", csrFile, "-subj", subject}
	output, success := utils.ExecuteCommandInDir(folder, "openssl", params...)
	if !success {
		utils.LogInfo("createcerts", "openssl command failed")
		utils.LogInfo("createcerts", "output: %s", output)
		return false
	}

	return true
}

func (certManager *CertificateManager) generateSiteCertificate(certPath, domain, csrFileName, certFileName, caCertPath, caKeyPath string) bool {
	content, err := utils.FetchFromGitHub("ssl", "ext.conf")
	if err != nil {
		utils.LogError(err, "createcerts")
		utils.LogInfo("createcerts", "failed to ext.conf")
		return false
	}

	content = utils.Template(content, utils.TemplateData{
		"domain": domain,
	})

	extFile := domain + ".ext"
	extPath := filepath.Join(certPath, extFile)
	utils.WriteStringToFile(extPath, content, constants.FilePermissions)

	params := []string{"x509", "-req", "-in", csrFileName, "-CA", caCertPath, "-CAkey", caKeyPath, "-CAcreateserial", "-out", certFileName, "-days", "3650", "-extfile", extFile}
	output, success := utils.ExecuteCommandInDir(certPath, "openssl", params...)
	if !success {
		utils.LogInfo("createcerts", "openssl command failed")
		utils.LogInfo("createcerts", "output: %s", output)
		return false
	}

	utils.RemoveFile(extPath)
	utils.RemoveFile(filepath.Join(constants.CertsDir, "sites", csrFileName))

	return true
}
