package occlient

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	routev1 "github.com/openshift/api/route/v1"
	componentlabels "github.com/redhat-developer/odo/pkg/component/labels"
	"github.com/redhat-developer/odo/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	ktesting "k8s.io/client-go/testing"
)

//func TestAddPVCToDeploymentConfig(t *testing.T) {
//	type args struct {
//		dc   *appsv1.DeploymentConfig
//		pvc  string
//		path string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		{
//			name: "first test",
//			args: args{
//				dc:  "ssa",
//				pvc: "storage1",
//				path: "/mnt"
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			fkclient, fkclientset := FakeNew()
//
//			_, err := fkclient.AddPVCToDeploymentConfig(tt.args.dc, tt.args.pvc, tt.args.path)
//
//			// Checks for error in positive cases
//			if !tt.wantErr == (err != nil) {
//				t.Errorf(" client.CreatePVC(name, size, labels) unexpected error %v, wantErr %v", err, tt.wantErr)
//			}
//
//			// Check for validating actions performed
//			if (len(fkclientset.Kubernetes.Actions()) != 1) && (tt.wantErr != true) {
//				t.Errorf("expected 1 action in CreatePVC got: %v", fkclientset.RouteClientset.Actions())
//			}
//			// Checks for return values in positive cases
//			if err == nil {
//				createdPVC := fkclientset.Kubernetes.Actions()[0].(ktesting.CreateAction).GetObject().(*corev1.PersistentVolumeClaim)
//				quantity, err := resource.ParseQuantity(tt.size)
//				if err != nil {
//					t.Errorf("failed to create quantity by calling resource.ParseQuantity(%v)", tt.size)
//				}
//
//				// created PVC should be labeled with labels passed to CreatePVC
//				if !reflect.DeepEqual(createdPVC.Labels, tt.labels) {
//					t.Errorf("labels in created route is not matching expected labels, expected: %v, got: %v", tt.labels, createdPVC.Labels)
//				}
//				// name, size of createdPVC should be matching to size, name passed to CreatePVC
//				if !reflect.DeepEqual(createdPVC.Spec.Resources.Requests["storage"], quantity) {
//					t.Errorf("size of PVC is not matching to expected size, expected: %v, got %v", quantity, createdPVC.Spec.Resources.Requests["storage"])
//				}
//				if createdPVC.Name != tt.name {
//					t.Errorf("PVC name is not matching to expected name, expected: %s, got %s", tt.name, createdPVC.Name)
//				}
//			}
//		})
//	}
//}

func TestCreatePVC(t *testing.T) {
	tests := []struct {
		name    string
		size    string
		labels  map[string]string
		wantErr bool
	}{
		{
			name: "storage 10Gi",
			size: "10Gi",
			labels: map[string]string{
				"name":      "mongodb",
				"namespace": "blog",
			},
			wantErr: false,
		},
		{
			name: "storage 1024",
			size: "1024",
			labels: map[string]string{
				"name":      "PostgreSQL",
				"namespace": "backend",
			},
			wantErr: false,
		},
		{
			name: "storage invalid size",
			size: "4#0",
			labels: map[string]string{
				"name":      "MySQL",
				"namespace": "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fkclient, fkclientset := FakeNew()

			_, err := fkclient.CreatePVC(tt.name, tt.size, tt.labels)

			componentLabels := componentlabels.GetLabels("nodejs", "nodejs", false)
			componentSelector := util.ConvertLabelsToSelector(componentLabels)
			dc, _ := fkclient.GetOneDeploymentConfigFromSelector(componentSelector)

			fmt.Println(dc)
			fmt.Println(componentLabels)
			fmt.Println(componentSelector)

			// Checks for error in positive cases
			if !tt.wantErr == (err != nil) {
				t.Errorf(" client.CreatePVC(name, size, labels) unexpected error %v, wantErr %v", err, tt.wantErr)
			}

			// Check for validating actions performed
			if (len(fkclientset.Kubernetes.Actions()) != 1) && (tt.wantErr != true) {
				t.Errorf("expected 1 action in CreatePVC got: %v", fkclientset.RouteClientset.Actions())
			}
			// Checks for return values in positive cases
			if err == nil {
				createdPVC := fkclientset.Kubernetes.Actions()[0].(ktesting.CreateAction).GetObject().(*corev1.PersistentVolumeClaim)
				quantity, err := resource.ParseQuantity(tt.size)
				if err != nil {
					t.Errorf("failed to create quantity by calling resource.ParseQuantity(%v)", tt.size)
				}

				// created PVC should be labeled with labels passed to CreatePVC
				if !reflect.DeepEqual(createdPVC.Labels, tt.labels) {
					t.Errorf("labels in created route is not matching expected labels, expected: %v, got: %v", tt.labels, createdPVC.Labels)
				}
				// name, size of createdPVC should be matching to size, name passed to CreatePVC
				if !reflect.DeepEqual(createdPVC.Spec.Resources.Requests["storage"], quantity) {
					t.Errorf("size of PVC is not matching to expected size, expected: %v, got %v", quantity, createdPVC.Spec.Resources.Requests["storage"])
				}
				if createdPVC.Name != tt.name {
					t.Errorf("PVC name is not matching to expected name, expected: %s, got %s", tt.name, createdPVC.Name)
				}
			}
		})
	}
}

func TestCreateRoute(t *testing.T) {
	tests := []struct {
		name    string
		service string
		labels  map[string]string
		wantErr bool
	}{
		{
			name:    "Case : mailserver",
			service: "mailserver",
			labels: map[string]string{
				"SLA": "High",
				"app.kubernetes.io/component-name": "backend",
				"app.kubernetes.io/component-type": "python",
			},
			wantErr: false,
		},

		{
			name:    "Case : blog",
			service: "blog",
			labels: map[string]string{
				"SLA": "High",
				"app.kubernetes.io/component-name": "backend",
				"app.kubernetes.io/component-type": "golang",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// initialising the fakeclient
			fkclient, fkclientset := FakeNew()

			_, err := fkclient.CreateRoute(tt.service, tt.labels)

			// Checks for error in positive cases
			if !tt.wantErr == (err != nil) {
				t.Errorf(" client.CreateRoute(string, labels) unexpected error %v, wantErr %v", err, tt.wantErr)
			}

			// Check for validating actions performed
			if len(fkclientset.RouteClientset.Actions()) != 1 {
				t.Errorf("expected 1 action in CreateRoute got: %v", fkclientset.RouteClientset.Actions())
			}
			// Checks for return values in positive cases
			if err == nil {
				createdRoute := fkclientset.RouteClientset.Actions()[0].(ktesting.CreateAction).GetObject().(*routev1.Route)
				// created route should be labeled with labels passed to CreateRoute
				if !reflect.DeepEqual(createdRoute.Labels, tt.labels) {
					t.Errorf("labels in created route is not matching expected labels, expected: %v, got: %v", tt.labels, createdRoute.Labels)
				}
				// route name and service that route is pointg to should match
				if createdRoute.Spec.To.Name != tt.service {
					t.Errorf("route is not matching to expected service name, expected: %s, got %s", tt.service, createdRoute)
				}
				if createdRoute.Name != tt.service {
					t.Errorf("route name is not matching to expected name, expected: %s, got %s", tt.service, createdRoute.Name)

				}
			}
		})
	}
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
