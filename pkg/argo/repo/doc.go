/*
Package repo is the repository of Argo. It is composed of two components, controller and storage,
the controller always reflect Workflow CRD to storage by the slice of keys. When you create the repository
it start to syncronize at background.

Create a new repository:

	// start to syncronize workflows from api-server
	r := NewArgoRepository(argoClientset, argoInformer, kubeClientset)

*/
package repo
