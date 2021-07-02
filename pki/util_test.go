package pki

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

func TestReadCertsAndKeysFromDir(t *testing.T) {

	comments := regexp.MustCompile(`(subject|issuer)=.+[\r\n]+`)
	certDir, err := ioutil.TempDir("", "certs-")
	if err != nil {
		t.Fatal(err)
	}
	for name, pem := range testPemFiles {
		if err := ioutil.WriteFile(fmt.Sprintf("%s/%s.pem", certDir, name), []byte(pem), 0644); err != nil {
			t.Fatal(err)
		}
	}
	certs, err := ReadCertsAndKeysFromDir(certDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range certs {
		received := c.CertificatePEM
		expected := comments.ReplaceAllString(testPemFiles[c.Name], "")
		assert.Equal(t, expected, received)
	}
	defer os.Remove(certDir) // clean up
}

var testPemFiles = map[string]string{
	"sub-ca": `-----BEGIN CERTIFICATE-----
MIIDqTCCApGgAwIBAgIUPDmvvihp1mPMKGBYvrF8bOZe7h0wDQYJKoZIhvcNAQEL
BQAwWjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAk5ZMQwwCgYDVQQHEwNMSUMxDDAK
BgNVBAoTA0s4UzEMMAoGA1UECxMDUktFMRQwEgYDVQQDEwtSS0UgUm9vdCBDQTAg
Fw0yMTA2MzAyMjE1MDBaGA8yMDUxMDYyMzIyMTUwMFowWTELMAkGA1UEBhMCVVMx
CzAJBgNVBAgTAk5ZMQwwCgYDVQQHEwNMSUMxDDAKBgNVBAoTA0s4UzEMMAoGA1UE
CxMDUktFMRMwEQYDVQQDEwpSS0UgU3ViIENBMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEAwOHE6HJXPNZEHOtlH3rM4SJlS7pFtHjoSsc3+9VKZbmSwbap
lL9I07q3WIVdVaW/Rnua3xg7FsvdeebhrTD3n8VKsduhovStPDreq2cNe51k93Mw
m6L8yxBEIHIv/up/TmRBlkyXp93hWJMDrA+K9Oe/yGMh7m3AgEMAVcK+AexNuyTV
tanHiqMlMYLLYwEZTI91ZFWJBkZgmNg/dk5LHAt7EKzeXH3DhFMGwqlneIcaogFq
R4lW4xPKAF9gv0LbR/trjUKr+K+0qairut4EFDcYjohP2L15GUYa9ZmyiTH0OMCh
uAjF9Dxe+pae7wYMQl36dGt8eYT1UEhscBLzfwIDAQABo2YwZDAOBgNVHQ8BAf8E
BAMCAYYwEgYDVR0TAQH/BAgwBgEB/wIBADAdBgNVHQ4EFgQUkXe0wifbup4B8zFm
u6WPUHkIPsEwHwYDVR0jBBgwFoAU3a3eVyFyQxfIuuokDxSWNiFEbQcwDQYJKoZI
hvcNAQELBQADggEBAJ4H0SHXCQRBtTGFwpGeCco8m4Tfnew2/64+Ua7BTzEPP1Zm
y7n5Xw40wJ5iuC87s8Ngs1vPjN9/peDrcmd8uN0sjW0pnpQ8fDl4Ks7fHlasBLiU
gTy01uux+J4vW+FM5rs/BJUxwA/KP8JEStpSQxrUEKxffGNeb5LH+qQHl7f9+wxE
E/ePm+2k0e8tgimHpno4ZUxScXTOSciN5I1pTfO7LxgQUGdhKXlXGQzZsO/EwqpG
n7/oqtviNudU+inIGaeKhvjTVBD9pWu57oAEieTd3quLgeKiOpCEvXQp9vFRky4R
m392kbtzI0HWYLjZAnn0B6pEKo4iDoKvdOpSln0=
-----END CERTIFICATE-----
`,
	"root-ca": `-----BEGIN CERTIFICATE-----
MIIDhDCCAmygAwIBAgIUKRrWt11PR7LF/tNEVq+nowy9Hh8wDQYJKoZIhvcNAQEL
BQAwWjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAk5ZMQwwCgYDVQQHEwNMSUMxDDAK
BgNVBAoTA0s4UzEMMAoGA1UECxMDUktFMRQwEgYDVQQDEwtSS0UgUm9vdCBDQTAe
Fw0yMTA2MzAyMjE1MDBaFw0zMTA2MjgyMjE1MDBaMFoxCzAJBgNVBAYTAlVTMQsw
CQYDVQQIEwJOWTEMMAoGA1UEBxMDTElDMQwwCgYDVQQKEwNLOFMxDDAKBgNVBAsT
A1JLRTEUMBIGA1UEAxMLUktFIFJvb3QgQ0EwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQCrHoz1iRA0ZekM3+G8M/mHCtTQMOVeZQvea0E9v47VPJ7iSGKp
/nyK4RPmbhch29nhtWRnQaqteJYVc3LaPk4MENac/UKYjwGlhhjYYumE/u7hTriv
7DPbrWfiDHyKHeBIiodZ6j45SaewWpU506dom0Yak09ltBgYipRAyJdw9nsnRcLA
y9tsne3Zj+mqC9myCKCEwfYq+i7IfkAgoU0iXrC6CBL9RjoV7o35SSS5Y8gEJqZH
gpVNbSYxA1hek+INW0Dx2lSferafsOdT8rYrV954cMu0fpQzU4FvUgbTb8lWQMx7
Uni0NzO8YneK4ZcACisTDFMMlXH3Y+9cS7E/AgMBAAGjQjBAMA4GA1UdDwEB/wQE
AwIBBjAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBTdrd5XIXJDF8i66iQPFJY2
IURtBzANBgkqhkiG9w0BAQsFAAOCAQEACRUl4VO1lRRN01lNf3RlN1mquNLBATxc
lRCoawrWwkYG5yXgUr41TbxrNVVHU1AnDhPRa2HdbT0zJ8afC8HMFexSjdhOPtnQ
R6JxREbuMQvWe/7OsTPR2yVHtw2kcn3FzbOL2/yb1DDzTEihSUoLCLr3Jnop6rPD
6iRAnqUegKckqKrf0ZztsIDS+a/azPpaUxfqmfQrRJzuZ6AWVQYWb4hDYVGg7gQK
5BPc30ireJ/juICjIhsQp/IJXVwQgf7N075x65rwWahHhRDYd/T5GG/xQoGw+fXB
QPGQKWLa3Z4YYqGAIZEYrOGQPty2VDlvzTe3XvlXf+wWhaTaZnl4pA==
-----END CERTIFICATE-----
`,
	"ca-bundle": `-----BEGIN CERTIFICATE-----
MIIDqTCCApGgAwIBAgIUPDmvvihp1mPMKGBYvrF8bOZe7h0wDQYJKoZIhvcNAQEL
BQAwWjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAk5ZMQwwCgYDVQQHEwNMSUMxDDAK
BgNVBAoTA0s4UzEMMAoGA1UECxMDUktFMRQwEgYDVQQDEwtSS0UgUm9vdCBDQTAg
Fw0yMTA2MzAyMjE1MDBaGA8yMDUxMDYyMzIyMTUwMFowWTELMAkGA1UEBhMCVVMx
CzAJBgNVBAgTAk5ZMQwwCgYDVQQHEwNMSUMxDDAKBgNVBAoTA0s4UzEMMAoGA1UE
CxMDUktFMRMwEQYDVQQDEwpSS0UgU3ViIENBMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEAwOHE6HJXPNZEHOtlH3rM4SJlS7pFtHjoSsc3+9VKZbmSwbap
lL9I07q3WIVdVaW/Rnua3xg7FsvdeebhrTD3n8VKsduhovStPDreq2cNe51k93Mw
m6L8yxBEIHIv/up/TmRBlkyXp93hWJMDrA+K9Oe/yGMh7m3AgEMAVcK+AexNuyTV
tanHiqMlMYLLYwEZTI91ZFWJBkZgmNg/dk5LHAt7EKzeXH3DhFMGwqlneIcaogFq
R4lW4xPKAF9gv0LbR/trjUKr+K+0qairut4EFDcYjohP2L15GUYa9ZmyiTH0OMCh
uAjF9Dxe+pae7wYMQl36dGt8eYT1UEhscBLzfwIDAQABo2YwZDAOBgNVHQ8BAf8E
BAMCAYYwEgYDVR0TAQH/BAgwBgEB/wIBADAdBgNVHQ4EFgQUkXe0wifbup4B8zFm
u6WPUHkIPsEwHwYDVR0jBBgwFoAU3a3eVyFyQxfIuuokDxSWNiFEbQcwDQYJKoZI
hvcNAQELBQADggEBAJ4H0SHXCQRBtTGFwpGeCco8m4Tfnew2/64+Ua7BTzEPP1Zm
y7n5Xw40wJ5iuC87s8Ngs1vPjN9/peDrcmd8uN0sjW0pnpQ8fDl4Ks7fHlasBLiU
gTy01uux+J4vW+FM5rs/BJUxwA/KP8JEStpSQxrUEKxffGNeb5LH+qQHl7f9+wxE
E/ePm+2k0e8tgimHpno4ZUxScXTOSciN5I1pTfO7LxgQUGdhKXlXGQzZsO/EwqpG
n7/oqtviNudU+inIGaeKhvjTVBD9pWu57oAEieTd3quLgeKiOpCEvXQp9vFRky4R
m392kbtzI0HWYLjZAnn0B6pEKo4iDoKvdOpSln0=
-----END CERTIFICATE-----

-----BEGIN CERTIFICATE-----
MIIDhDCCAmygAwIBAgIUKRrWt11PR7LF/tNEVq+nowy9Hh8wDQYJKoZIhvcNAQEL
BQAwWjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAk5ZMQwwCgYDVQQHEwNMSUMxDDAK
BgNVBAoTA0s4UzEMMAoGA1UECxMDUktFMRQwEgYDVQQDEwtSS0UgUm9vdCBDQTAe
Fw0yMTA2MzAyMjE1MDBaFw0zMTA2MjgyMjE1MDBaMFoxCzAJBgNVBAYTAlVTMQsw
CQYDVQQIEwJOWTEMMAoGA1UEBxMDTElDMQwwCgYDVQQKEwNLOFMxDDAKBgNVBAsT
A1JLRTEUMBIGA1UEAxMLUktFIFJvb3QgQ0EwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQCrHoz1iRA0ZekM3+G8M/mHCtTQMOVeZQvea0E9v47VPJ7iSGKp
/nyK4RPmbhch29nhtWRnQaqteJYVc3LaPk4MENac/UKYjwGlhhjYYumE/u7hTriv
7DPbrWfiDHyKHeBIiodZ6j45SaewWpU506dom0Yak09ltBgYipRAyJdw9nsnRcLA
y9tsne3Zj+mqC9myCKCEwfYq+i7IfkAgoU0iXrC6CBL9RjoV7o35SSS5Y8gEJqZH
gpVNbSYxA1hek+INW0Dx2lSferafsOdT8rYrV954cMu0fpQzU4FvUgbTb8lWQMx7
Uni0NzO8YneK4ZcACisTDFMMlXH3Y+9cS7E/AgMBAAGjQjBAMA4GA1UdDwEB/wQE
AwIBBjAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBTdrd5XIXJDF8i66iQPFJY2
IURtBzANBgkqhkiG9w0BAQsFAAOCAQEACRUl4VO1lRRN01lNf3RlN1mquNLBATxc
lRCoawrWwkYG5yXgUr41TbxrNVVHU1AnDhPRa2HdbT0zJ8afC8HMFexSjdhOPtnQ
R6JxREbuMQvWe/7OsTPR2yVHtw2kcn3FzbOL2/yb1DDzTEihSUoLCLr3Jnop6rPD
6iRAnqUegKckqKrf0ZztsIDS+a/azPpaUxfqmfQrRJzuZ6AWVQYWb4hDYVGg7gQK
5BPc30ireJ/juICjIhsQp/IJXVwQgf7N075x65rwWahHhRDYd/T5GG/xQoGw+fXB
QPGQKWLa3Z4YYqGAIZEYrOGQPty2VDlvzTe3XvlXf+wWhaTaZnl4pA==
-----END CERTIFICATE-----
`,
	"ca-bundle-with-comments": `subject=CN=RKE Sub CA, OU=RKE, O=K8S, L=LIC, ST=NY, C=US
subject=CN=RKE Root CA, OU=RKE, O=K8S, L=LIC, ST=NY, C=US
-----BEGIN CERTIFICATE-----
MIIDqTCCApGgAwIBAgIUPDmvvihp1mPMKGBYvrF8bOZe7h0wDQYJKoZIhvcNAQEL
BQAwWjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAk5ZMQwwCgYDVQQHEwNMSUMxDDAK
BgNVBAoTA0s4UzEMMAoGA1UECxMDUktFMRQwEgYDVQQDEwtSS0UgUm9vdCBDQTAg
Fw0yMTA2MzAyMjE1MDBaGA8yMDUxMDYyMzIyMTUwMFowWTELMAkGA1UEBhMCVVMx
CzAJBgNVBAgTAk5ZMQwwCgYDVQQHEwNMSUMxDDAKBgNVBAoTA0s4UzEMMAoGA1UE
CxMDUktFMRMwEQYDVQQDEwpSS0UgU3ViIENBMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEAwOHE6HJXPNZEHOtlH3rM4SJlS7pFtHjoSsc3+9VKZbmSwbap
lL9I07q3WIVdVaW/Rnua3xg7FsvdeebhrTD3n8VKsduhovStPDreq2cNe51k93Mw
m6L8yxBEIHIv/up/TmRBlkyXp93hWJMDrA+K9Oe/yGMh7m3AgEMAVcK+AexNuyTV
tanHiqMlMYLLYwEZTI91ZFWJBkZgmNg/dk5LHAt7EKzeXH3DhFMGwqlneIcaogFq
R4lW4xPKAF9gv0LbR/trjUKr+K+0qairut4EFDcYjohP2L15GUYa9ZmyiTH0OMCh
uAjF9Dxe+pae7wYMQl36dGt8eYT1UEhscBLzfwIDAQABo2YwZDAOBgNVHQ8BAf8E
BAMCAYYwEgYDVR0TAQH/BAgwBgEB/wIBADAdBgNVHQ4EFgQUkXe0wifbup4B8zFm
u6WPUHkIPsEwHwYDVR0jBBgwFoAU3a3eVyFyQxfIuuokDxSWNiFEbQcwDQYJKoZI
hvcNAQELBQADggEBAJ4H0SHXCQRBtTGFwpGeCco8m4Tfnew2/64+Ua7BTzEPP1Zm
y7n5Xw40wJ5iuC87s8Ngs1vPjN9/peDrcmd8uN0sjW0pnpQ8fDl4Ks7fHlasBLiU
gTy01uux+J4vW+FM5rs/BJUxwA/KP8JEStpSQxrUEKxffGNeb5LH+qQHl7f9+wxE
E/ePm+2k0e8tgimHpno4ZUxScXTOSciN5I1pTfO7LxgQUGdhKXlXGQzZsO/EwqpG
n7/oqtviNudU+inIGaeKhvjTVBD9pWu57oAEieTd3quLgeKiOpCEvXQp9vFRky4R
m392kbtzI0HWYLjZAnn0B6pEKo4iDoKvdOpSln0=
-----END CERTIFICATE-----

subject=CN=RKE Root CA, OU=RKE, O=K8S, L=LIC, ST=NY, C=US
-----BEGIN CERTIFICATE-----
MIIDhDCCAmygAwIBAgIUKRrWt11PR7LF/tNEVq+nowy9Hh8wDQYJKoZIhvcNAQEL
BQAwWjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAk5ZMQwwCgYDVQQHEwNMSUMxDDAK
BgNVBAoTA0s4UzEMMAoGA1UECxMDUktFMRQwEgYDVQQDEwtSS0UgUm9vdCBDQTAe
Fw0yMTA2MzAyMjE1MDBaFw0zMTA2MjgyMjE1MDBaMFoxCzAJBgNVBAYTAlVTMQsw
CQYDVQQIEwJOWTEMMAoGA1UEBxMDTElDMQwwCgYDVQQKEwNLOFMxDDAKBgNVBAsT
A1JLRTEUMBIGA1UEAxMLUktFIFJvb3QgQ0EwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQCrHoz1iRA0ZekM3+G8M/mHCtTQMOVeZQvea0E9v47VPJ7iSGKp
/nyK4RPmbhch29nhtWRnQaqteJYVc3LaPk4MENac/UKYjwGlhhjYYumE/u7hTriv
7DPbrWfiDHyKHeBIiodZ6j45SaewWpU506dom0Yak09ltBgYipRAyJdw9nsnRcLA
y9tsne3Zj+mqC9myCKCEwfYq+i7IfkAgoU0iXrC6CBL9RjoV7o35SSS5Y8gEJqZH
gpVNbSYxA1hek+INW0Dx2lSferafsOdT8rYrV954cMu0fpQzU4FvUgbTb8lWQMx7
Uni0NzO8YneK4ZcACisTDFMMlXH3Y+9cS7E/AgMBAAGjQjBAMA4GA1UdDwEB/wQE
AwIBBjAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBTdrd5XIXJDF8i66iQPFJY2
IURtBzANBgkqhkiG9w0BAQsFAAOCAQEACRUl4VO1lRRN01lNf3RlN1mquNLBATxc
lRCoawrWwkYG5yXgUr41TbxrNVVHU1AnDhPRa2HdbT0zJ8afC8HMFexSjdhOPtnQ
R6JxREbuMQvWe/7OsTPR2yVHtw2kcn3FzbOL2/yb1DDzTEihSUoLCLr3Jnop6rPD
6iRAnqUegKckqKrf0ZztsIDS+a/azPpaUxfqmfQrRJzuZ6AWVQYWb4hDYVGg7gQK
5BPc30ireJ/juICjIhsQp/IJXVwQgf7N075x65rwWahHhRDYd/T5GG/xQoGw+fXB
QPGQKWLa3Z4YYqGAIZEYrOGQPty2VDlvzTe3XvlXf+wWhaTaZnl4pA==
-----END CERTIFICATE-----
`,
}
