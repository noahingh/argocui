/*

Package tree present a Argo workflow as tree-like format, such as the "argo get" command. It return rows which consist of some columns.
For example if you want the tree of root node:

	rows, err := GetTreeRoot(w)
	if err != nil {
		panic(err)
	}

	for _, r := range rows {
		fmt.Println(r[0])
	}

*/
package tree
