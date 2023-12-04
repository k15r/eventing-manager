package object

import (
	"fmt"
	"reflect"
	"testing"

	apigateway "github.com/kyma-incubator/api-gateway/api/v1beta1"
	"github.com/stretchr/testify/require"
	kmeta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"

	eventingv1alpha2 "github.com/kyma-project/eventing-manager/api/eventing/v1alpha2"
)

func TestApplyExistingAPIRuleAttributes(t *testing.T) {
	// given
	const (
		name            = "name-0"
		generateName    = "0123"
		resourceVersion = "4567"
	)

	var (
		host   = ptr.To("some.host")
		status = apigateway.APIRuleStatus{
			LastProcessedTime:    ptr.To(kmeta.Time{}),
			ObservedGeneration:   512,
			APIRuleStatus:        nil,
			VirtualServiceStatus: nil,
			AccessRuleStatus:     nil,
		}
	)

	type args struct {
		givenSrc *apigateway.APIRule
		givenDst *apigateway.APIRule
		wantDst  *apigateway.APIRule
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ApiRule attributes are applied from src to dst",
			args: args{
				givenSrc: &apigateway.APIRule{
					ObjectMeta: kmeta.ObjectMeta{
						Name:            name,
						GenerateName:    generateName,
						ResourceVersion: resourceVersion,
					},
					Spec:   apigateway.APIRuleSpec{Host: host},
					Status: status,
				},
				givenDst: &apigateway.APIRule{
					ObjectMeta: kmeta.ObjectMeta{
						Name:            name,
						GenerateName:    generateName,
						ResourceVersion: resourceVersion,
					},
					Spec:   apigateway.APIRuleSpec{Host: host},
					Status: status,
				},
				wantDst: &apigateway.APIRule{
					ObjectMeta: kmeta.ObjectMeta{
						Name:            name,
						GenerateName:    "",
						ResourceVersion: resourceVersion,
					},
					Spec:   apigateway.APIRuleSpec{Host: host},
					Status: status,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			ApplyExistingAPIRuleAttributes(tt.args.givenSrc, tt.args.givenDst)

			// then
			require.Equal(t, tt.args.wantDst.Name, tt.args.givenDst.Name)
			require.Equal(t, tt.args.wantDst.GenerateName, tt.args.givenDst.GenerateName)
			require.Equal(t, tt.args.wantDst.ResourceVersion, tt.args.givenDst.ResourceVersion)
			require.Equal(t, tt.args.wantDst.Spec, tt.args.givenDst.Spec)
			require.Equal(t, tt.args.wantDst.Status, tt.args.givenDst.Status)
		})
	}
}

func TestGetService(t *testing.T) {
	// given
	const (
		name       = "name-0"
		port       = uint32(9080)
		isExternal = true
	)

	type args struct {
		svcName string
		port    uint32
	}
	tests := []struct {
		name string
		args args
		want apigateway.Service
	}{
		{
			name: "get service with the given properties",
			args: args{
				svcName: name,
				port:    port,
			},
			want: apigateway.Service{
				Name:       ptr.To(name),
				Port:       ptr.To(port),
				IsExternal: ptr.To(isExternal),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			if got := GetService(tt.args.svcName, tt.args.port); !reflect.DeepEqual(got, tt.want) {
				// then
				t.Errorf("GetService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAPIRule(t *testing.T) {
	// given
	const (
		namespace  = "namespace-0"
		namePrefix = "name-0"
	)

	type args struct {
		ns         string
		namePrefix string
		opts       []Option
	}
	tests := []struct {
		name string
		args args
		want *apigateway.APIRule
	}{
		{
			name: "get APIRule with the given properties",
			args: args{
				ns:         namespace,
				namePrefix: namePrefix,
				opts:       nil,
			},
			want: &apigateway.APIRule{
				TypeMeta: kmeta.TypeMeta{},
				ObjectMeta: kmeta.ObjectMeta{
					Namespace:    namespace,
					GenerateName: namePrefix,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			if got := NewAPIRule(tt.args.ns, tt.args.namePrefix, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				// then
				t.Errorf("NewAPIRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDuplicateValues(t *testing.T) {
	// given
	type args struct {
		values []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "list without duplicates",
			args: args{
				values: []string{
					"1", "2", "3",
				},
			},
			want: []string{
				"1", "2", "3",
			},
		},
		{
			name: "list with duplicates",
			args: args{
				values: []string{
					"1", "2", "3",
					"3", "2", "1",
				},
			},
			want: []string{
				"1", "2", "3",
			},
		},
		{
			name: "empty list",
			args: args{
				values: []string{},
			},
			want: []string{},
		},
		{
			name: "nil list",
			args: args{
				values: nil,
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			if got := RemoveDuplicateValues(tt.args.values); !reflect.DeepEqual(got, tt.want) {
				// then
				t.Errorf("RemoveDuplicateValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithGateway(t *testing.T) {
	// given
	const (
		gateway = "some.gateway"
	)

	type args struct {
		givenGateway string
		givenObject  *apigateway.APIRule
	}
	tests := []struct {
		name       string
		args       args
		wantObject *apigateway.APIRule
	}{
		{
			name: "apply gateway to object",
			args: args{
				givenGateway: gateway,
				givenObject:  &apigateway.APIRule{},
			},
			wantObject: &apigateway.APIRule{
				Spec: apigateway.APIRuleSpec{
					Gateway: ptr.To(gateway),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			WithGateway(tt.args.givenGateway)(tt.args.givenObject)

			// then
			require.Equal(t, tt.wantObject.Spec.Gateway, tt.args.givenObject.Spec.Gateway)
		})
	}
}

func TestWithLabels(t *testing.T) {
	// given
	type args struct {
		givenLabels map[string]string
		givenObject *apigateway.APIRule
	}
	tests := []struct {
		name       string
		args       args
		wantObject *apigateway.APIRule
	}{
		{
			name: "object with nil labels",
			args: args{
				givenLabels: map[string]string{
					"key-0": "val-0",
				},
				givenObject: &apigateway.APIRule{
					ObjectMeta: kmeta.ObjectMeta{
						Labels: nil,
					},
				},
			},
			wantObject: &apigateway.APIRule{
				ObjectMeta: kmeta.ObjectMeta{
					Labels: map[string]string{
						"key-0": "val-0",
					},
				},
			},
		},
		{
			name: "object with empty labels",
			args: args{
				givenLabels: map[string]string{
					"key-0": "val-0",
				},
				givenObject: &apigateway.APIRule{
					ObjectMeta: kmeta.ObjectMeta{
						Labels: map[string]string{},
					},
				},
			},
			wantObject: &apigateway.APIRule{
				ObjectMeta: kmeta.ObjectMeta{
					Labels: map[string]string{
						"key-0": "val-0",
					},
				},
			},
		},
		{
			name: "object with labels",
			args: args{
				givenLabels: map[string]string{
					"key-0": "val-0",
				},
				givenObject: &apigateway.APIRule{
					ObjectMeta: kmeta.ObjectMeta{
						Labels: map[string]string{
							"key-1": "val-1",
							"key-2": "val-2",
						},
					},
				},
			},
			wantObject: &apigateway.APIRule{
				ObjectMeta: kmeta.ObjectMeta{
					Labels: map[string]string{
						"key-0": "val-0",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			WithLabels(tt.args.givenLabels)(tt.args.givenObject)

			// then
			require.Equal(t, tt.wantObject.GetLabels(), tt.args.givenObject.GetLabels())
		})
	}
}

func TestWithOwnerReference(t *testing.T) {
	// given
	const (
		kind0 = "kind-0"
		kind1 = "kind-1"
		kind2 = "kind-2"

		apiVersion0 = "version-0"
		apiVersion1 = "version-1"
		apiVersion2 = "version-2"

		name0 = "name-0"
		name1 = "name-1"
		name2 = "name-2"

		uid0 = "000000"
		uid1 = "111111"
		uid2 = "222222"

		blockOwnerDeletion = true
	)

	var (
		sub0 = eventingv1alpha2.Subscription{
			TypeMeta:   kmeta.TypeMeta{Kind: kind0, APIVersion: apiVersion0},
			ObjectMeta: kmeta.ObjectMeta{Name: name0, UID: uid0},
		}
		sub1 = eventingv1alpha2.Subscription{
			TypeMeta:   kmeta.TypeMeta{Kind: kind1, APIVersion: apiVersion1},
			ObjectMeta: kmeta.ObjectMeta{Name: name1, UID: uid1},
		}
	)

	type args struct {
		givenSubs   []eventingv1alpha2.Subscription
		givenObject *apigateway.APIRule
	}
	tests := []struct {
		name       string
		args       args
		wantObject *apigateway.APIRule
	}{
		{
			name: "nil Subscriptions",
			args: args{
				givenSubs: nil,
				givenObject: &apigateway.APIRule{
					ObjectMeta: kmeta.ObjectMeta{
						OwnerReferences: nil,
					},
				},
			},
			wantObject: &apigateway.APIRule{
				ObjectMeta: kmeta.ObjectMeta{
					OwnerReferences: []kmeta.OwnerReference{},
				},
			},
		},
		{
			name: "empty Subscriptions",
			args: args{
				givenSubs: []eventingv1alpha2.Subscription{},
				givenObject: &apigateway.APIRule{
					ObjectMeta: kmeta.ObjectMeta{
						OwnerReferences: nil,
					},
				},
			},
			wantObject: &apigateway.APIRule{
				ObjectMeta: kmeta.ObjectMeta{
					OwnerReferences: []kmeta.OwnerReference{},
				},
			},
		},
		{
			name: "object with nil OwnerReferences",
			args: args{
				givenSubs: []eventingv1alpha2.Subscription{
					sub0,
					sub1,
				},
				givenObject: &apigateway.APIRule{
					ObjectMeta: kmeta.ObjectMeta{
						OwnerReferences: nil,
					},
				},
			},
			wantObject: &apigateway.APIRule{
				ObjectMeta: kmeta.ObjectMeta{
					OwnerReferences: []kmeta.OwnerReference{
						{
							APIVersion:         apiVersion0,
							Kind:               kind0,
							Name:               name0,
							UID:                uid0,
							BlockOwnerDeletion: ptr.To(blockOwnerDeletion),
						},
						{
							APIVersion:         apiVersion1,
							Kind:               kind1,
							Name:               name1,
							UID:                uid1,
							BlockOwnerDeletion: ptr.To(blockOwnerDeletion),
						},
					},
				},
			},
		},
		{
			name: "object with empty OwnerReferences",
			args: args{
				givenSubs: []eventingv1alpha2.Subscription{
					sub0,
					sub1,
				},
				givenObject: &apigateway.APIRule{
					ObjectMeta: kmeta.ObjectMeta{
						OwnerReferences: []kmeta.OwnerReference{},
					},
				},
			},
			wantObject: &apigateway.APIRule{
				ObjectMeta: kmeta.ObjectMeta{
					OwnerReferences: []kmeta.OwnerReference{
						{
							APIVersion:         apiVersion0,
							Kind:               kind0,
							Name:               name0,
							UID:                uid0,
							BlockOwnerDeletion: ptr.To(blockOwnerDeletion),
						},
						{
							APIVersion:         apiVersion1,
							Kind:               kind1,
							Name:               name1,
							UID:                uid1,
							BlockOwnerDeletion: ptr.To(blockOwnerDeletion),
						},
					},
				},
			},
		},
		{
			name: "object with OwnerReferences",
			args: args{
				givenSubs: []eventingv1alpha2.Subscription{
					sub0,
					sub1,
				},
				givenObject: &apigateway.APIRule{
					ObjectMeta: kmeta.ObjectMeta{
						OwnerReferences: []kmeta.OwnerReference{
							{
								APIVersion:         apiVersion2,
								Kind:               kind2,
								Name:               name2,
								UID:                uid2,
								BlockOwnerDeletion: ptr.To(blockOwnerDeletion),
							},
						},
					},
				},
			},
			wantObject: &apigateway.APIRule{
				ObjectMeta: kmeta.ObjectMeta{
					OwnerReferences: []kmeta.OwnerReference{
						{
							APIVersion:         apiVersion0,
							Kind:               kind0,
							Name:               name0,
							UID:                uid0,
							BlockOwnerDeletion: ptr.To(blockOwnerDeletion),
						},
						{
							APIVersion:         apiVersion1,
							Kind:               kind1,
							Name:               name1,
							UID:                uid1,
							BlockOwnerDeletion: ptr.To(blockOwnerDeletion),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			WithOwnerReference(tt.args.givenSubs)(tt.args.givenObject)
			gotOwnerReferences := tt.args.givenObject.GetOwnerReferences()

			// then
			require.Equal(t, tt.wantObject.GetOwnerReferences(), gotOwnerReferences)
		})
	}
}

func TestWithRules(t *testing.T) {
	// given
	const (
		endpoint0 = "/endpoint0"
		endpoint1 = "/endpoint1"

		sink0 = "https://sink0.com" + endpoint0
		sink1 = "https://sink1.com" + endpoint1

		certsURL = "some.url"
		name     = "name-0"
		port     = uint32(9999)
		external = true
	)

	var (
		sub0 = eventingv1alpha2.Subscription{
			Spec: eventingv1alpha2.SubscriptionSpec{Sink: sink0},
		}
		sub1 = eventingv1alpha2.Subscription{
			Spec: eventingv1alpha2.SubscriptionSpec{Sink: sink1},
		}

		methods = []string{"some.method"}
	)

	type args struct {
		givenCertsURL string
		givenSubs     []eventingv1alpha2.Subscription
		givenSvc      apigateway.Service
		givenMethods  []string
		givenObject   *apigateway.APIRule
	}
	tests := []struct {
		name       string
		args       args
		wantObject *apigateway.APIRule
	}{
		{
			name: "apply properties to object",
			args: args{
				givenCertsURL: certsURL,
				givenSubs: []eventingv1alpha2.Subscription{
					sub0,
					sub1,
				},
				givenSvc: apigateway.Service{
					Name:       ptr.To(name),
					Port:       ptr.To(port),
					IsExternal: ptr.To(external),
				},
				givenMethods: methods,
				givenObject:  &apigateway.APIRule{},
			},
			wantObject: &apigateway.APIRule{
				Spec: apigateway.APIRuleSpec{
					Rules: []apigateway.Rule{
						{
							Path: endpoint0,
							Service: &apigateway.Service{
								Name:       ptr.To(name),
								Port:       ptr.To(port),
								IsExternal: ptr.To(external),
							},
							Methods: methods,
							AccessStrategies: []*apigateway.Authenticator{
								{
									Handler: &apigateway.Handler{
										Name: OAuthHandlerNameJWT,
										Config: &runtime.RawExtension{
											Raw: []byte(fmt.Sprintf(JWKSURLFormat, certsURL)),
										},
									},
								},
							},
						},
						{
							Path: endpoint1,
							Service: &apigateway.Service{
								Name:       ptr.To(name),
								Port:       ptr.To(port),
								IsExternal: ptr.To(external),
							},
							Methods: methods,
							AccessStrategies: []*apigateway.Authenticator{
								{
									Handler: &apigateway.Handler{
										Name: OAuthHandlerNameJWT,
										Config: &runtime.RawExtension{
											Raw: []byte(fmt.Sprintf(JWKSURLFormat, certsURL)),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			WithRules(tt.args.givenCertsURL, tt.args.givenSubs, tt.args.givenSvc, tt.args.givenMethods...)(tt.args.givenObject)

			// then
			require.Equal(t, tt.wantObject.Spec.Rules, tt.args.givenObject.Spec.Rules)
		})
	}
}

func TestWithService(t *testing.T) {
	// given
	const (
		host     = "host0"
		name     = "name-0"
		port     = uint32(9999)
		external = true
	)

	type args struct {
		givenHost    string
		givenSvcName string
		givenPort    uint32
		givenObject  *apigateway.APIRule
	}
	tests := []struct {
		name       string
		args       args
		wantObject *apigateway.APIRule
	}{
		{
			name: "apply properties to object",
			args: args{
				givenHost:    host,
				givenSvcName: name,
				givenPort:    port,
				givenObject:  &apigateway.APIRule{},
			},
			wantObject: &apigateway.APIRule{
				Spec: apigateway.APIRuleSpec{
					Host: ptr.To(host),
					Service: &apigateway.Service{
						Name:       ptr.To(name),
						Port:       ptr.To(port),
						IsExternal: ptr.To(external),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			WithService(tt.args.givenHost, tt.args.givenSvcName, tt.args.givenPort)(tt.args.givenObject)

			// then
			require.Equal(t, tt.wantObject, tt.args.givenObject)
		})
	}
}
