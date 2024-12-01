package model

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Product struct {
	IProdID           int32   `db:"iProdID" json:"iProdID"`
	IPCatID           int32   `db:"iPCatID" json:"iPCatID"`
	VCategoryName     string  `db:"vCategoryName" json:"vCategoryName"`
	CCode             *string `db:"cCode" json:"cCode"`
	VName             string  `db:"vName" json:"vName"`
	VURLName          string  `db:"vUrlName" json:"vUrlName"`
	VShortDesc        *string `db:"vShortDesc" json:"vShortDesc"`
	VDescription      *string `db:"vDescription" json:"vDescription"`
	FPrice            float64 `db:"fPrice" json:"fPrice"`
	FOPrice           float64 `db:"fOPrice" json:"fOPrice"`
	FActualWeight     float64 `db:"fActualWeight" json:"fActualWeight"`
	FVolumetricWeight float64 `db:"fVolumetricWeight" json:"fVolumetricWeight"`
	VSmallImage       *string `db:"vSmallImage" json:"vSmallImage"`
	VSmallImageAltTag *string `db:"vSmallImage_AltTag" json:"vSmallImage_AltTag"`
	VImage            *string `db:"vImage" json:"vImage"`
	VImageAltTag      *string `db:"vImage_AltTag" json:"vImage_AltTag"`
	CStatus           *string `db:"cStatus" json:"cStatus"`
	VYTID             *string `db:"vYTID" json:"vYTID"`
}

type ProductImage struct {
	IProdImageID int32   `json:"iProdImageID,omitempty" db:"iProdImageID"`
	VName        *string `json:"vName,omitempty" db:"vName"`
	VAltTag      *string `json:"vAlt_Tag,omitempty" db:"vAltTag"`
	CStatus      *string `json:"cStatus,omitempty" db:"cStatus"`
	VType        string  `json:"vType,omitempty" db:"vType"`
}

type CategoryProducts struct {
	Category *Category
	Products *[]Product
}

type ProductAttribute struct {
	VName  string `db:"vName" json:"vName"`
	VValue string `db:"vValue" json:"vValue"`
}

func (m *Model) Product(iProdID int32) (*Product, error) {

	p := Product{}
	qry := `SELECT 
				p.iProdID, 
				p.vName,
				c.vName vCategoryName,
				p.vUrlName, 
				p.vShortDesc,
				p.vDescription,
				p.fPrice,
				p.fOPrice,
				p.cStatus
			FROM 
				product p JOIN
				prodcat c ON p.iPCatID = c.iPCatID
			WHERE p.iProdID = ?`

	if err := m.DbHandle.Get(&p, qry, iProdID); err != nil {
		return nil, fmt.Errorf("error fetching product %d: %w", iProdID, err)
	}

	return &p, nil
}

func (m *Model) ProductByUrl(vUrlName string) (*Product, error) {

	p := Product{}
	qry := `SELECT 
				p.iProdID, 
				p.vName,
				c.vName vCategoryName,
				p.vImage,
				p.vUrlName, 
				p.vShortDesc,
				p.vDescription,
				p.fPrice,
				p.fOPrice,
				p.cStatus
			FROM 
				product p JOIN
				prodcat c ON p.iPCatID = c.iPCatID
			WHERE p.vUrlName = ?`

	if err := m.DbHandle.Get(&p, qry, vUrlName); err != nil {
		return nil, fmt.Errorf("error fetching product %s: %w", vUrlName, err)
	}

	return &p, nil
}

func (m *Model) ProductAttributes(iProdID int32) (*[]ProductAttribute, error) {

	attribs := []ProductAttribute{}
	qry := `SELECT
				a.vName,
				pa.vValue 
			FROM 
				product_attrib pa JOIN
				attribute a ON pa.iAttribID = a.iAttribID
			WHERE pa.iProdID = ?`

	if err := m.DbHandle.Select(&attribs, qry, iProdID); err != nil {
		return nil, fmt.Errorf("error fetching product %d: %w", iProdID, err)
	}

	return &attribs, nil
}

func (m *Model) CategoryProducts(vUrlName string) (*[]CategoryProducts, error) {

	catIDs, err := m.AllSubIDs(vUrlName, true)
	if err != nil {
		return nil, err
	}

	*catIDs = (*catIDs)[:1]
	query := `SELECT
					p.iProdID,
					p.iPCatID,
					p.cCode,
					p.vName,
					p.vUrlName,
					p.vShortDesc,
					p.vDescription,
					p.fPrice,
					p.fOPrice,
					p.fActualWeight,
					p.fVolumetricWeight,
					p.vSmallImage,
					p.vSmallImage_AltTag,
					p.vImage,
					p.vImage_AltTag,
					p.cStatus,
					p.vYTID
				FROM product p
					JOIN prodcat c ON p.iPCatID = c.iPCatID
				WHERE p.iPCatID IN ( ? ) AND p.cStatus = 'A' AND c.cStatus = 'A'
				ORDER BY
					iPCatID,
					cStatus,
					vName`

	query, args, err := sqlx.In(query, *catIDs)
	if err != nil {
		return nil, fmt.Errorf("error binding category products %s: %w", vUrlName, err)
	}
	query = m.DbHandle.Rebind(query)

	products := []Product{}
	if err = m.DbHandle.Select(&products, query, args...); err != nil {
		return nil, err
	}

	// Make the map
	var cID int32
	id := products[0].IPCatID
	cps := new([]Product)
	cpMap := make(map[int32]*[]Product, 0)
	for _, p := range products {
		cID = p.IPCatID
		if cID != id {
			cpMap[id] = cps

			cps = new([]Product)
			id = cID
		}
		*cps = append(*cps, p)
	}
	cpMap[id] = cps

	catMap, err := m.CategoryMapByID()
	if err != nil {
		return nil, err
	}

	cpList := make([]CategoryProducts, 0)
	for _, iPCatID := range *catIDs {

		category := catMap[iPCatID]
		products, found := cpMap[iPCatID]
		if !found {
			continue
		}
		cpList = append(cpList, CategoryProducts{category, products})
	}

	return &cpList, nil
}

func (m *Model) ProductImages(productID int32) (*[]ProductImage, error) {

	pImages := []ProductImage{}

	query := `SELECT
					0             as iProdImageID,
					vImage        as vName,
					vImage_AltTag as vAltTag,
					'A'           as cStatus,
					'main'        as vType
				FROM product
				WHERE iProdID = ?
			UNION
				SELECT
					iProdImageID,
					vPic         as vName,
					vTitle       as vAltTag,
					cStatus,
					'additional' as vType
				FROM product_images
				WHERE iProdID = ?
				ORDER BY 1`

	if err := m.DbHandle.Select(&pImages, query, productID, productID); err != nil {
		return nil, err
	}

	return &pImages, nil
}
