/*

Package tree present a Argo workflow as tree-like format, such as the "argo get" command. 
The type of return value is the double dimension of slice.

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
