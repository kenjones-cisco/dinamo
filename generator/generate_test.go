package generator

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func uniqueID() string {
	b := make([]byte, 36)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func makeTmp(name string) string {
	tmp := path.Join(os.TempDir(), uniqueID())
	if err := os.Mkdir(tmp, os.ModePerm); err != nil {
		return path.Join(os.TempDir(), name)
	}
	return path.Join(tmp, name)
}

func readOutFile(name string) string {
	rawData, err := ioutil.ReadFile(name)
	if err != nil {
		return ""
	}
	return string(rawData)
}

func TestGenerate_withSourceData(t *testing.T) {
	type args struct {
		inputTemplate string
		outfile       string
		ds            *DataSources
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Generate Kubconfig successfully",
			args: args{
				inputTemplate: "fixtures/config.tmpl",
				outfile:       makeTmp("config.yaml"),
				ds:            &DataSources{Data: []string{"NAMESPACE=cnt-refapp-dev", "USERNAME=branlyon", "INSTANCE=", "INSTANCE_NAME="}},
			},
			wantErr: false,
		},
		{
			name: "Generate basic yaml file successfully",
			args: args{
				inputTemplate: "fixtures/basic.tmpl",
				outfile:       makeTmp("basic.yaml"),
				ds:            &DataSources{Data: []string{"name=ui-alln", "containerPort=8080", "image_name=redis", "namespace=testing", "kind=v1"}},
			},
			wantErr: false,
		},
		{
			name: "Generate basic file unsuccessfully due to incorrect data",
			args: args{
				inputTemplate: "fixtures/basic.tmpl",
				outfile:       makeTmp("wrongdata.yaml"),
				ds:            &DataSources{Data: []string{"firstName=Brandon"}},
			},
			wantErr: true,
		},
		{
			name: "Generate basic file successfully with too much data",
			args: args{
				inputTemplate: "fixtures/basic.tmpl",
				outfile:       makeTmp("toomuchdata.yaml"),
				ds:            &DataSources{Data: []string{"name=ui-alln", "containerPort=8080", "image_name=redis", "namespace=testing", "kind=v1", "Test=DoesNotFail"}},
			},
			wantErr: false,
		},
		{
			name: "Generate basic file unsuccessfully due to empty data",
			args: args{
				inputTemplate: "fixtures/nestedmap.tmpl",
				outfile:       makeTmp("emptydata.yaml"),
				ds:            &DataSources{},
			},
			wantErr: true,
		},
		{
			name: "Generate basic file unsuccessfully due to invalid outfile file path",
			args: args{
				inputTemplate: "fixtures/nestedmap.tmpl",
				outfile:       "./templates/InvalidPath.yaml",
				ds:            &DataSources{Data: []string{"name=ui-alln", "containerPort=8080", "image_name=redis", "namespace=testing", "kind=v1"}},
			},
			wantErr: true,
		},
		{
			name: "Generate basic file unsuccessfully due to invalid inputTemplate file path",
			args: args{
				inputTemplate: "template/nestedmap.tmpl",
				outfile:       makeTmp("invalidpath.yaml"),
				ds:            &DataSources{Data: []string{"name=ui-alln", "containerPort=8080", "image_name=redis", "namespace=testing", "kind=v1"}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		if err := Generate(tt.args.inputTemplate, tt.args.outfile, tt.args.ds); (err != nil) != tt.wantErr {
			t.Errorf("%q. Generate() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		// Checks if Output file exist
		if _, err := os.Stat(tt.args.outfile); os.IsNotExist(err) {
			if (err != nil) != tt.wantErr {
				// path does not exist
				t.Errorf("File path does not exist: %v ", err)
			}
		}

	}
}

func TestGenerate_withSourceDataFile(t *testing.T) {
	type args struct {
		inputTemplate string
		outfile       string
		ds            *DataSources
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Generate Resource successfully - json file",
			args: args{
				inputTemplate: "fixtures/complex.tmpl",
				outfile:       makeTmp("resources1j.yaml"),
				ds:            &DataSources{DataFile: "fixtures/source1.json"},
			},
			wantErr: false,
		},

		{
			name: "Generate Resource successfully - yaml file",
			args: args{
				inputTemplate: "fixtures/complex.tmpl",
				outfile:       makeTmp("resources1y.yaml"),
				ds:            &DataSources{DataFile: "fixtures/source1.yaml"},
			},
			wantErr: false,
		},
		{
			name: "Generate file that has nested maps in data successfully - json file",
			args: args{
				inputTemplate: "fixtures/nestedmap.tmpl",
				outfile:       makeTmp("nestedmap1j.yaml"),
				ds:            &DataSources{DataFile: "fixtures/source2.json"},
			},
			wantErr: false,
		},
		{
			name: "Generate file that has nested maps in data successfully - yml file",
			args: args{
				inputTemplate: "fixtures/nestedmap.tmpl",
				outfile:       makeTmp("nestedmap1y.yaml"),
				ds:            &DataSources{DataFile: "fixtures/source2.yml"},
			},
			wantErr: false,
		},
		{
			name: "Generate basic file unsuccessfully due to incorrect data file type",
			args: args{
				inputTemplate: "fixtures/basic.tmpl",
				outfile:       makeTmp("wrongdata.yaml"),
				ds:            &DataSources{DataFile: "fixtures/source1.toml"},
			},
			wantErr: true,
		},
		{
			name: "Generate basic file unsuccessfully due to unknown data file",
			args: args{
				inputTemplate: "fixtures/basic.tmpl",
				outfile:       makeTmp("wrongdata.yaml"),
				ds:            &DataSources{DataFile: "random/fakefilethatdoesnotexist.yaml"},
			},
			wantErr: true,
		},
		{
			name: "Generate basic file unsuccessfully due to invalid yaml data file",
			args: args{
				inputTemplate: "fixtures/complex.tmpl",
				outfile:       makeTmp("resourcesxy.yaml"),
				ds:            &DataSources{DataFile: "fixtures/badsource.yaml"},
			},
			wantErr: true,
		},
		{
			name: "Generate basic file unsuccessfully due to invalid json data file",
			args: args{
				inputTemplate: "fixtures/complex.tmpl",
				outfile:       makeTmp("resourcesxj.yaml"),
				ds:            &DataSources{DataFile: "fixtures/badsource.json"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		if err := Generate(tt.args.inputTemplate, tt.args.outfile, tt.args.ds); (err != nil) != tt.wantErr {
			t.Errorf("%q. Generate() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		// Checks if Output file exist
		if _, err := os.Stat(tt.args.outfile); os.IsNotExist(err) {
			if (err != nil) != tt.wantErr {
				// path does not exist
				t.Errorf("File path does not exist: %v ", err)
			}
		}

	}
}

var wantConfig = `apiVersion: v1
clusters:
- cluster:
    insecure-skip-tls-verify: true
    server: https://cae-alln.cisco.com:443
  name: cae-alln-cisco-com:443
- cluster:
    insecure-skip-tls-verify: true
    server: https://cae-rcdn.cisco.com:443
  name: cae-rcdn-cisco-com:443
- cluster:
    insecure-skip-tls-verify: true
    server: https://cae-rtp.cisco.com:443
  name: cae-rtp-cisco-com:443
- cluster:
    insecure-skip-tls-verify: true
    server: https://localhost:443
  name: localhost:443
contexts:
- context:
    cluster: cae-alln-cisco-com:443
    namespace: mynamespace
    user: fakeuser/cae-alln-cisco-com:443
  name: mynamespace/cae-alln-cisco-com:443/fakeuser
- context:
    cluster: cae-rcdn-cisco-com:443
    namespace: mynamespace
    user: fakeuser/cae-rcdn-cisco-com:443
  name: mynamespace/cae-rcdn-cisco-com:443/fakeuser
- context:
    cluster: cae-rtp-cisco-com:443
    namespace: mynamespace
    user: fakeuser/cae-rtp-cisco-com:443
  name: mynamespace/cae-rtp-cisco-com:443/fakeuser
- context:
    cluster: localhost:443
    namespace: mynamespace
    user: fakeuser/localhost:443
  name: mynamespace/localhost:443/fakeuser
`

var wantSecret = `apiVersion: v1
kind: Secret
metadata:
  name: test-app
  namespace: mynamespace
  annotations:
    lifecycle: dev
data:
  common.db.password: NUpCUzB0UURWSHB3RjJFTVVhdks=
  common.internal.api.key: VGx2bzNoMTlTZ1lYdFZ6YnpVQ21taE5RNUhSU1hkTDdjbE9tRUFlYXNIeVBWQndieE4zdFBo
  common.client.password: b1NPQlorU0FDa0hSNjlaNWQ1QWI=
  common.client.secret: eUR3UkFZY0RIeGRCekQvN3pTUUNCYnV2WkMvag==
  common.openstack.client.password: MnlkcUprb3QwTzJJMjRlNEtZcFQ=
  common.openstack.client.secret: Ti9XdjFzZU1Bb1p6SEtSMS9LbkVJVWJueWF5ZXIyTzVrN1N4ZzladzhvRDdrbCtrdUtOZ2RBPT0=
  common.iam.ldap.password: cUU4d09yQmpMKzJvejdJVExwRHI=
  common.synthetic.influx: aW5mbHV4PWh0dHA6Ly9kYnVzZXJuYW1lOlNac1pPS2pHYkhNdDFkZE5XQ2U5QHNlcnZlci1hcGkuZXhhbXBsZS5jb206ODIwMC9kYm5hbWU=
`

func TestGenerate_withSourceEnv(t *testing.T) {
	envVars := []string{"INSTANCE", "INSTANCE_NAME", "NAMESPACE", "USERNAME", "APP_NAME", "LIFECYCLE",
		"SECRET_DB_PASSWORD", "SECRET_INTERNAL_API_KEY", "SECRET_CLIENT_PASSWORD", "SECRET_CLIENT_SECRET",
		"SECRET_OPENSTACK_CLIENT_PASSWORD", "SECRET_OPENSTACK_CLIENT_SECRET", "SECRET_IAM_LDAP_PASSWORD",
		"SECRET_SYNTHETIC_INFLUX",
	}
	defer func() {
		for _, k := range envVars {
			_ = os.Unsetenv(k)
		}
	}()

	_ = os.Setenv("INSTANCE", "https://localhost:443")
	_ = os.Setenv("INSTANCE_NAME", "localhost:443")
	_ = os.Setenv("NAMESPACE", "mynamespace")
	_ = os.Setenv("USERNAME", "fakeuser")
	_ = os.Setenv("APP_NAME", "test-app")
	_ = os.Setenv("LIFECYCLE", "dev")
	_ = os.Setenv("SECRET_DB_PASSWORD", "5JBS0tQDVHpwF2EMUavK")
	_ = os.Setenv("SECRET_INTERNAL_API_KEY", "Tlvo3h19SgYXtVzbzUCmmhNQ5HRSXdL7clOmEAeasHyPVBwbxN3tPh")
	_ = os.Setenv("SECRET_CLIENT_PASSWORD", "oSOBZ+SACkHR69Z5d5Ab")
	_ = os.Setenv("SECRET_CLIENT_SECRET", "yDwRAYcDHxdBzD/7zSQCBbuvZC/j")
	_ = os.Setenv("SECRET_OPENSTACK_CLIENT_PASSWORD", "2ydqJkot0O2I24e4KYpT")
	_ = os.Setenv("SECRET_OPENSTACK_CLIENT_SECRET", "N/Wv1seMAoZzHKR1/KnEIUbnyayer2O5k7Sxg9Zw8oD7kl+kuKNgdA==")
	_ = os.Setenv("SECRET_IAM_LDAP_PASSWORD", "qE8wOrBjL+2oz7ITLpDr")
	_ = os.Setenv("SECRET_SYNTHETIC_INFLUX", "influx=http://dbusername:SZsZOKjGbHMt1ddNWCe9@server-api.example.com:8200/dbname")

	type args struct {
		inputTemplate string
		outfile       string
		ds            *DataSources
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Generate Kubconfig successfully - use env",
			args: args{
				inputTemplate: "fixtures/config.tmpl",
				outfile:       makeTmp("config.yaml"),
				ds:            &DataSources{UseEnv: true},
			},
			want:    wantConfig,
			wantErr: false,
		},
		{
			name: "Generate k8s secrets successfully - use env",
			args: args{
				inputTemplate: "fixtures/secrets.tmpl",
				outfile:       makeTmp("secrets.yaml"),
				ds:            &DataSources{UseEnv: true},
			},
			want:    wantSecret,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		if err := Generate(tt.args.inputTemplate, tt.args.outfile, tt.args.ds); (err != nil) != tt.wantErr {
			t.Errorf("%q. Generate() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		// Checks if Output file exist
		if _, err := os.Stat(tt.args.outfile); os.IsNotExist(err) {
			if (err != nil) != tt.wantErr {
				// path does not exist
				t.Errorf("File path does not exist: %v ", err)
			}
		}

		got := readOutFile(tt.args.outfile)
		if got != tt.want {
			t.Errorf("expected: %v, got: %v", tt.want, got)
		}
	}
}

var wantConfigMulti = `apiVersion: v1
clusters:
- cluster:
    insecure-skip-tls-verify: true
    server: https://cae-alln.cisco.com:443
  name: cae-alln-cisco-com:443
- cluster:
    insecure-skip-tls-verify: true
    server: https://cae-rcdn.cisco.com:443
  name: cae-rcdn-cisco-com:443
- cluster:
    insecure-skip-tls-verify: true
    server: https://cae-rtp.cisco.com:443
  name: cae-rtp-cisco-com:443
- cluster:
    insecure-skip-tls-verify: true
    server: https://localhost:443
  name: localhost:443
contexts:
- context:
    cluster: cae-alln-cisco-com:443
    namespace: testns
    user: otheruser.gen/cae-alln-cisco-com:443
  name: testns/cae-alln-cisco-com:443/otheruser.gen
- context:
    cluster: cae-rcdn-cisco-com:443
    namespace: testns
    user: otheruser.gen/cae-rcdn-cisco-com:443
  name: testns/cae-rcdn-cisco-com:443/otheruser.gen
- context:
    cluster: cae-rtp-cisco-com:443
    namespace: testns
    user: otheruser.gen/cae-rtp-cisco-com:443
  name: testns/cae-rtp-cisco-com:443/otheruser.gen
- context:
    cluster: localhost:443
    namespace: testns
    user: otheruser.gen/localhost:443
  name: testns/localhost:443/otheruser.gen
`

var wantSecretMulti = `apiVersion: v1
kind: Secret
metadata:
  name: test-app
  namespace: mynamespace
  annotations:
    lifecycle: prod
data:
  common.db.password: NUpCUzB0UURWSHB3RjJFTVVhdks=
  common.internal.api.key: VGx2bzNoMTlTZ1lYdFZ6YnpVQ21taE5RNUhSU1hkTDdjbE9tRUFlYXNIeVBWQndieE4zdFBo
  common.client.password: b1NPQlorU0FDa0hSNjlaNWQ1QWI=
  common.client.secret: eUR3UkFZY0RIeGRCekQvN3pTUUNCYnV2WkMvag==
  common.openstack.client.password: MnlkcUprb3QwTzJJMjRlNEtZcFQ=
  common.openstack.client.secret: Ti9XdjFzZU1Bb1p6SEtSMS9LbkVJVWJueWF5ZXIyTzVrN1N4ZzladzhvRDdrbCtrdUtOZ2RBPT0=
  common.iam.ldap.password: cUU4d09yQmpMKzJvejdJVExwRHI=
  common.synthetic.influx: aW5mbHV4PWh0dHA6Ly9kYnVzZXJuYW1lOlNac1pPS2pHYkhNdDFkZE5XQ2U5QHNlcnZlci1hcGkuZXhhbXBsZS5jb206ODIwMC9kYm5hbWU=
`

func TestGenerate_withMultipleSources(t *testing.T) {
	envVars := []string{"INSTANCE", "INSTANCE_NAME", "NAMESPACE", "APP_NAME"}
	defer func() {
		for _, k := range envVars {
			_ = os.Unsetenv(k)
		}
	}()

	_ = os.Setenv("INSTANCE", "https://localhost:443")
	_ = os.Setenv("INSTANCE_NAME", "localhost:443")
	_ = os.Setenv("NAMESPACE", "mynamespace")
	_ = os.Setenv("APP_NAME", "test-app")

	type args struct {
		inputTemplate string
		outfile       string
		ds            *DataSources
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Generate Kubconfig successfully - use multiple sources",
			args: args{
				inputTemplate: "fixtures/config.tmpl",
				outfile:       makeTmp("config.yaml"),
				ds: &DataSources{
					Data:     []string{"NAMESPACE=testns"},
					DataFile: "fixtures/source_app.yaml",
					UseEnv:   true,
				},
			},
			want:    wantConfigMulti,
			wantErr: false,
		},
		{
			name: "Generate k8s secrets successfully - use multiple sources",
			args: args{
				inputTemplate: "fixtures/secrets.tmpl",
				outfile:       makeTmp("secrets.yaml"),
				ds: &DataSources{
					Data:     []string{"LIFECYCLE=prod"},
					DataFile: "fixtures/source_secret.yaml",
					UseEnv:   true,
				},
			},
			want:    wantSecretMulti,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		if err := Generate(tt.args.inputTemplate, tt.args.outfile, tt.args.ds); (err != nil) != tt.wantErr {
			t.Errorf("%q. Generate() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		// Checks if Output file exist
		if _, err := os.Stat(tt.args.outfile); os.IsNotExist(err) {
			if (err != nil) != tt.wantErr {
				// path does not exist
				t.Errorf("File path does not exist: %v ", err)
			}
		}

		got := readOutFile(tt.args.outfile)
		if got != tt.want {
			t.Errorf("expected: %v, got: %v", tt.want, got)
		}
	}
}
