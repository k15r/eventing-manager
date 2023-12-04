package eventing

import (
	"context"
	"testing"

	eventingcontrollermocks "github.com/kyma-project/eventing-manager/internal/controller/operator/eventing/mocks"

	kapiclientsetfake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"

	"github.com/kyma-project/eventing-manager/pkg/k8s"

	kadmissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	kcorev1 "k8s.io/api/core/v1"

	"github.com/kyma-project/eventing-manager/pkg/env"

	"github.com/kyma-project/eventing-manager/options"

	"github.com/kyma-project/eventing-manager/pkg/logger"

	"k8s.io/apimachinery/pkg/runtime"

	natsv1alpha1 "github.com/kyma-project/nats-manager/api/v1alpha1"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	operatorv1alpha1 "github.com/kyma-project/eventing-manager/api/operator/v1alpha1"
	eventingmocks "github.com/kyma-project/eventing-manager/pkg/eventing/mocks"

	kdynamicfake "k8s.io/client-go/dynamic/fake"
)

// MockedUnitTestEnvironment provides mocked resources for unit tests.
type MockedUnitTestEnvironment struct {
	Context         context.Context
	Client          client.Client
	kubeClient      *k8s.Client
	eventingManager *eventingmocks.Manager
	ctrlManager     *eventingcontrollermocks.Manager
	Reconciler      *Reconciler
	Logger          *logger.Logger
	Recorder        *record.FakeRecorder
}

func NewMockedUnitTestEnvironment(t *testing.T, objs ...client.Object) *MockedUnitTestEnvironment {
	// setup context
	ctx := context.Background()

	// setup logger
	ctrLogger, err := logger.New("json", "info")
	require.NoError(t, err)

	// setup fake client for k8s
	newScheme := runtime.NewScheme()
	err = natsv1alpha1.AddToScheme(newScheme)
	require.NoError(t, err)
	err = operatorv1alpha1.AddToScheme(newScheme)
	require.NoError(t, err)
	err = kcorev1.AddToScheme(newScheme)
	require.NoError(t, err)
	err = kadmissionregistrationv1.AddToScheme(newScheme)
	require.NoError(t, err)

	// Create a fake dynamic client
	fakeDynamicClient := kdynamicfake.NewSimpleDynamicClient(newScheme)

	fakeClientBuilder := fake.NewClientBuilder().WithScheme(newScheme)
	fakeClient := fakeClientBuilder.WithObjects(objs...).WithStatusSubresource(objs...).Build()
	fakeClientSet := kapiclientsetfake.NewSimpleClientset()
	recorder := &record.FakeRecorder{}
	kubeClient := k8s.NewKubeClient(fakeClient, fakeClientSet, "eventing-manager", fakeDynamicClient)

	// setup custom mocks
	eventingManager := new(eventingmocks.Manager)
	mockManager := new(eventingcontrollermocks.Manager)

	opts := options.New()

	// get backend configs.
	backendConfig := env.BackendConfig{}

	// setup reconciler
	reconciler := NewReconciler(
		fakeClient,
		kubeClient,
		fakeDynamicClient,
		newScheme,
		ctrLogger,
		recorder,
		eventingManager,
		backendConfig,
		nil,
		opts,
		nil,
	)
	reconciler.ctrlManager = mockManager

	return &MockedUnitTestEnvironment{
		Context:         ctx,
		Client:          fakeClient,
		kubeClient:      &kubeClient,
		Reconciler:      reconciler,
		Logger:          ctrLogger,
		Recorder:        recorder,
		eventingManager: eventingManager,
		ctrlManager:     mockManager,
	}
}

func (testEnv *MockedUnitTestEnvironment) GetEventing(name, namespace string) (operatorv1alpha1.Eventing, error) {
	var evnt operatorv1alpha1.Eventing
	err := testEnv.Client.Get(testEnv.Context, types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, &evnt)
	return evnt, err
}
