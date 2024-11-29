package model

import "fmt"

type Category struct {
	IPCatID           int32       `db:"iPCatID" json:"iPCatID" schema:"iPCatID,required"`
	VName             string      `db:"vName" json:"vName" schema:"vName,required"`
	VURLName          string      `db:"vUrlName" json:"vUrlName" schema:"vUrlName"`
	IParentID         int32       `db:"iParentID" json:"iParentID" schema:"iParentID"`
	VShortDesc        *string     `db:"vShortDesc" json:"vShortDesc" schema:"vShortDesc"`
	VMenuImage        *string     `db:"vMenuImage" json:"vMenuImage" schema:"vMenuImage"`
	VMenuImage_AltTag *string     `db:"vMenuImage_AltTag" json:"vMenuImage_AltTag" schema:"vMenuImage_AltTag"`
	CStatus           string      `db:"cStatus" json:"cStatus" schema:"cStatus,required"`
	Children          []*Category `db:"-" json:"children,omitempty"`
	IProductCount     int32       `db:"iProductCount" json:"iProductCount,omitempty"`
}

type CategoryAttribute struct {
	IAttribDatID int64   `db:"iAttribDatID" json:"iAttribDatID" diff:"iAttribDatID"`
	IPCatID      int64   `db:"iPCatID" json:"iPCatID" diff:"iPCatID"`
	IAttribID    int64   `db:"iAttribID" json:"iAttribID" diff:"iAttribID"`
	VAttribName  *string `db:"vAttribName" json:"vAttribName" diff:"-"`
	VName        *string `db:"vName" json:"vName" diff:"vName"`
	IRank        int     `db:"iRank" json:"iRank" diff:"iRank"`
}

type categoryMap map[int32][]*Category

// All three vars tie up the same category structs
var subCatMap categoryMap            // maps id to child categories
var catMapByURL map[string]*Category // maps vUrlName to category
var catMapByID map[int32]*Category   // map category ID to category
var catRoot *Category                // root of the category tree

func (m *Model) CategoryTree() (*Category, error) {
	// TODO: Fix crashes if it finds a category with iPCatID = 0
	query := `SELECT
					prodcat.iPCatID,
					prodcat.iParentID,
					prodcat.vName,
					prodcat.vUrlName,
					prodcat.vMenuImage,
					prodcat.vMenuImage_AltTag,
					prodcat.cStatus,
					count(distinct iProdID) as iProductCount
				FROM prodcat
				LEFT JOIN product ON prodcat.iPCatID = product.iPCatID
				GROUP BY
					prodcat.iPCatID,
					prodcat.iParentID,
					prodcat.vName,
					prodcat.vUrlName,
					prodcat.vMenuImage,
					prodcat.vMenuImage_AltTag,
					prodcat.cStatus`

	rows, err := m.DbHandle.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Build the catMap and the subCatMap
	subCatMap = make(categoryMap, 0)
	catMapByURL = make(map[string]*Category, 0)
	catMapByID = make(map[int32]*Category)

	for rows.Next() {
		var cat Category
		if err := rows.StructScan(&cat); err != nil {
			return nil, err
		}

		subCatMap[cat.IParentID] = append(subCatMap[cat.IParentID], &cat)

		catMapByURL[cat.VURLName] = &cat
		catMapByID[cat.IPCatID] = &cat
	}

	// Build the tree starting from an "arbitrarily" defined root
	catRoot = &Category{
		IPCatID:  0,
		VName:    "All Categories",
		VURLName: "category-root",
		CStatus:  "A",
	}

	catMapByURL["category-root"] = catRoot
	catMapByID[0] = catRoot

	subCatMap.mapTraverse(catRoot, func(c *Category) {
		c.Children = subCatMap[c.IPCatID]
	})

	return catRoot, nil
}

func (m *Model) CategoryMapByURL() (map[string]*Category, error) {

	if len(catMapByURL) == 0 {
		var err error
		catRoot, err = m.CategoryTree()
		if err != nil {
			return map[string]*Category{}, err
		}
	}

	return catMapByURL, nil
}

func (m *Model) CategoryMapByID() (map[int32]*Category, error) {

	if len(catMapByURL) == 0 {
		var err error
		catRoot, err = m.CategoryTree()
		if err != nil {
			return map[int32]*Category{}, err
		}
	}

	return catMapByID, nil
}

const MAX_ITERATIONS = 5

func (m *Model) ReverseCatTree(c *Category) ([]*Category, error) {

	revTree := []*Category{c}
	parentID := c.IParentID
	for iter := 1; iter < MAX_ITERATIONS; iter++ {

		if parentID == 0 {
			return revTree, nil
		}

		cp, found := catMapByID[parentID]
		if !found {
			return nil, fmt.Errorf("cannot find parent category of category ID %d", parentID)
		}

		revTree = append([]*Category{cp}, revTree...)
		c = cp
		parentID = c.IParentID
	}

	return nil, fmt.Errorf("exceeded max iterations finding top category")
}

func (m *Model) AllSubIDs(vUrlName string, withProducts bool) (*[]int32, error) {

	if catRoot == nil {
		var err error
		catRoot, err = m.CategoryTree()
		if err != nil {
			return nil, err
		}
	}

	var categories = make([]int32, 0)
	start, found := catMapByURL[vUrlName]
	if !found {
		return nil, fmt.Errorf("no category found with URL: %s", vUrlName)
	}

	start.treeTraverse(func(c *Category) {
		if withProducts {
			if c.IProductCount > 0 {
				categories = append(categories, c.IPCatID)
				return
			}
		} else {
			categories = append(categories, c.IPCatID)
		}
	})

	return &categories, nil
}

func (m *Model) AllSubCategories(root *Category, withProducts bool) ([]*Category, error) {

	if catRoot == nil {
		var err error
		catRoot, err = m.CategoryTree()
		if err != nil {
			return nil, err
		}
	}

	var categories = make([]*Category, 0)
	start, found := catMapByURL[root.VURLName]
	if !found {
		return nil, fmt.Errorf("no category found with URL: %s", root.VURLName)
	}

	start.treeTraverse(func(c *Category) {
		if withProducts {
			if c.IProductCount > 0 {
				categories = append(categories, c)
				return
			}
		} else {
			categories = append(categories, c)
		}
	})

	return categories, nil

}

// The traversal function "walks" through the tree executing a function
// for each node that is passed to us. This function has two parameters
// (1) the value of the node and (2) the depth of the node that is being
// visited relative to the starting node.

func (cMap *categoryMap) mapTraverse(start *Category, f func(*Category)) {

	f(start)
	for _, c := range (*cMap)[start.IPCatID] {
		cMap.mapTraverse(c, f)
	}
}

func (c *Category) treeTraverse(f func(*Category)) {

	f(c)
	for _, subC := range c.Children {
		subC.treeTraverse(f)
	}
}
