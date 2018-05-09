package occlient

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	// projectv1 "github.com/openshift/api/project/v1"
	fkprojectclientset "github.com/openshift/client-go/project/clientset/versioned/fake"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// projectclientset "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	fkubernetes "k8s.io/client-go/kubernetes/fake"
)

func Fakenew() (*Client, error) {
	var client Client

	//    client.kubeConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	//    config, err := client.kubeConfig.ClientConfig()
	//    if err != nil {
	//        return nil, err
	//    }

	kubeClient := fkubernetes.NewSimpleClientset()
	client.kubeClient = kubeClient

	//    imageClient, err := imageclientset.NewForConfig(config)
	//    if err != nil {
	//        return nil, err
	//    }
	//    client.imageClient = imageClient
	//
	//    appsClient, err := appsclientset.NewForConfig(config)
	//    if err != nil {
	//        return nil, err
	//    }
	//    client.appsClient = appsClient
	//
	//    buildClient, err := buildclientset.NewForConfig(config)
	//    if err != nil {
	//        return nil, err
	//    }
	//    client.buildClient = buildClient
	//
	//    serviceCatalogClient, err := servicecatalogclienset.NewForConfig(config)
	//    if err != nil {
	//        return nil, err
	//    }
	//    client.serviceCatalogClient = serviceCatalogClient
	//
	// fmt.Println("flag b4 projectClient init")
	projectClient := fkprojectclientset.NewSimpleClientset().Project()
	client.projectClient = projectClient
	//
	//    routeClient, err := routeclientset.NewForConfig(config)
	//    if err != nil {
	//        return nil, err
	//    }
	//    client.routeClient = routeClient
	//
	//    namespace, _, err := client.kubeConfig.Namespace()
	//    if err != nil {
	//        return nil, err
	//    }
	//    client.namespace = namespace
	//
	//    // The following should go away once we're done with complete migration to
	//    // client-go
	//    ocpath, err := getOcBinary()
	//    if err != nil {
	//        return nil, errors.Wrap(err, "unable to get oc binary")
	//    }
	//    client.ocpath = ocpath
	//
	//    if !isServerUp(client.ocpath) {
	//        return nil, errors.New("Unable to connect to OpenShift cluster, is it down?")
	//    }
	//    if !isLoggedIn(client.ocpath) {
	//        return nil, errors.New("Please log in to the cluster")
	//    }

	return &client, nil
}

/// Delete this code
// func checkError(err error, context string, a ...interface{}) {
// 	if err != nil {
// 		//	log.Debugf("Error:\n%v", err)
// 		if context == "" {
// 			fmt.Println("ds") //errors.Cause(err))
// 		} else {
// 			fmt.Printf(fmt.Sprintf("%s\n", context), a...)
// 		}
//
// 		os.Exit(1)
// 	}
// }

/// Delete this code

func getFakeocClient() *Client {
	client, _ := Fakenew()
	//	checkError(err, "")
	return client
}

func TestGetOcBinary(t *testing.T) {

	// test setup
	// test shouldn't have external dependency, so we are faking oc binary with empty tmpfile
	tmpfile, err := ioutil.TempFile("", "fake-oc")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile1, err := ioutil.TempFile("", "fake-oc")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	defer os.Remove(tmpfile1.Name())

	type args struct {
		oc string
	}
	tests := []struct {
		name    string
		envs    map[string]string
		want    string
		wantErr bool
	}{
		{
			name: "set via KUBECTL_PLUGINS_CALLER exists",
			envs: map[string]string{
				"KUBECTL_PLUGINS_CALLER": tmpfile.Name(),
			},
			want:    tmpfile.Name(),
			wantErr: false,
		},
		{
			name: "set via KUBECTL_PLUGINS_CALLER (invalid file)",
			envs: map[string]string{
				"KUBECTL_PLUGINS_CALLER": "invalid",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "set via OC_BIN exists",
			envs: map[string]string{
				"OC_BIN": tmpfile.Name(),
			},
			want:    tmpfile.Name(),
			wantErr: false,
		},
		{
			name: "set via OC_BIN (invalid file)",
			envs: map[string]string{
				"OC_BIN": "invalid",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "bot OC_BIN and KUBECTL_PLUGINS_CALLER set",
			envs: map[string]string{
				"OC_BIN":                 tmpfile.Name(),
				"KUBECTL_PLUGINS_CALLER": tmpfile1.Name(),
			},
			want:    tmpfile1.Name(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// cleanup variables before running test
			os.Unsetenv("OC_BIN")
			os.Unsetenv("KUBECTL_PLUGINS_CALLER")

			for k, v := range tt.envs {
				if err := os.Setenv(k, v); err != nil {
					t.Error(err)
				}
			}
			got, err := getOcBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("getOcBinary() unexpected error \n%v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOcBinary() \ngot: %v \nwant: %v", got, tt.want)
			}
		})
	}
}

func TestAddLabelsToArgs(t *testing.T) {
	tests := []struct {
		name     string
		argsIn   []string
		labels   map[string]string
		argsOut1 []string
		argsOut2 []string
	}{
		{
			name:   "one label in empty args",
			argsIn: []string{},
			labels: map[string]string{
				"label1": "value1",
			},
			argsOut1: []string{
				"--labels", "label1=value1",
			},
		},
		{
			name: "one label with existing args",
			argsIn: []string{
				"--foo", "bar",
			},
			labels: map[string]string{
				"label1": "value1",
			},
			argsOut1: []string{
				"--foo", "bar",
				"--labels", "label1=value1",
			},
		},
		{
			name: "multiple label with existing args",
			argsIn: []string{
				"--foo", "bar",
			},
			labels: map[string]string{
				"label1": "value1",
				"label2": "value2",
			},
			argsOut1: []string{
				"--foo", "bar",
				"--labels", "label1=value1,label2=value2",
			},
			argsOut2: []string{
				"--foo", "bar",
				"--labels", "label2=value2,label1=value1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			argsGot := addLabelsToArgs(tt.labels, tt.argsIn)

			if !reflect.DeepEqual(argsGot, tt.argsOut1) && !reflect.DeepEqual(argsGot, tt.argsOut2) {
				t.Errorf("addLabelsToArgs() \ngot:  %#v \nwant: %#v or %#v", argsGot, tt.argsOut1, tt.argsOut2)
			}
		})
	}
}

func Test_parseImageName(t *testing.T) {

	tests := []struct {
		arg     string
		want1   string
		want2   string
		want3   string
		wantErr bool
	}{
		{
			arg:     "nodejs:8",
			want1:   "nodejs",
			want2:   "8",
			want3:   "",
			wantErr: false,
		},
		{
			arg:     "nodejs@sha256:7e56ca37d1db225ebff79dd6d9fd2a9b8f646007c2afc26c67962b85dd591eb2",
			want1:   "nodejs",
			want2:   "",
			want3:   "sha256:7e56ca37d1db225ebff79dd6d9fd2a9b8f646007c2afc26c67962b85dd591eb2",
			wantErr: false,
		},
		{
			arg:     "nodejs@sha256:asdf@",
			wantErr: true,
		},
		{
			arg:     "nodejs@@",
			wantErr: true,
		},
		{
			arg:     "nodejs::",
			wantErr: true,
		},
		{
			arg:     "nodejs",
			want1:   "nodejs",
			want2:   "latest",
			want3:   "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("image name: %s", tt.arg)
		t.Run(name, func(t *testing.T) {
			got1, got2, got3, err := parseImageName(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseImageName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got1 != tt.want1 {
				t.Errorf("parseImageName() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("parseImageName() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("parseImageName() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}

// func TestCreateNewProject(t *testing.T) {
// 	//client := getFakeocClient()
// 	obj := &projectv1.ProjectRequest{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: "junkfoo",
// 		},
// 	}
// 	client := fkprojectclientset.NewSimpleClientset(obj)
//
// 	_, err := client.Project().ProjectRequests().Create(obj)
// 	if err != nil {
// 		t.Errorf("some thing went wrong")
// 	}
// }

func TestCreateNewProject(t *testing.T) {
	// type fields struct {
	// 	projectClient        projectclientset.ProjectV1Interface
	// }
	// type args struct {
	// 	name string
	// }
	// tests := []struct {
	// 	name    string
	// 	fields  fields
	// 	args    args
	// 	wantErr bool
	// }{
	// // TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		c := &Client{
	// 			ocpath:               tt.fields.ocpath,
	// 			kubeClient:           tt.fields.kubeClient,
	// 			imageClient:          tt.fields.imageClient,
	// 			appsClient:           tt.fields.appsClient,
	// 			buildClient:          tt.fields.buildClient,
	// 			serviceCatalogClient: tt.fields.serviceCatalogClient,
	// 			projectClient:        tt.fields.projectClient,
	// 			routeClient:          tt.fields.routeClient,
	// 			kubeConfig:           tt.fields.kubeConfig,
	// 			namespace:            tt.fields.namespace,
	// 		}
	// 		if err := c.CreateNewProject(tt.args.name); (err != nil) != tt.wantErr {
	// 			t.Errorf("Client.CreateNewProject() error = %v, wantErr %v", err, tt.wantErr)
	// 		}
	// 	})
	// }

	//newFakeClient := func() *Client {
	//	return &Client{
	//		projectClient: fkprojectclientset.NewSimpleClientset().Project(),
	//	}
	//}
	c, _ := Fakenew()
	fmt.Println("reached")
	err := c.CreateNewProject("foo")
	if err != nil {
		t.Error(err)
	}

}
