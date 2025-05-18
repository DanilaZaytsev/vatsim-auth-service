package middleware

import (
	"context"
)

type Claims struct {
	CID          uint64
	Email        string
	Roles        string
	CountryName  string
	DivisionName string
}

// Получение значений из контекста и сборка Claims
func GetClaimsFromContext(ctx context.Context) *Claims {
	cidVal := ctx.Value(ContextKeyCID)
	emailVal := ctx.Value(ContextKeyEmail)
	rolesVal := ctx.Value(ContextKeyRoles)
	countryVal := ctx.Value(ContextKey("country_name"))
	divisionVal := ctx.Value(ContextKey("division_name"))

	if cidVal == nil || emailVal == nil || rolesVal == nil {
		return nil
	}

	cid, ok1 := cidVal.(uint64)
	email, ok2 := emailVal.(string)
	roles, ok3 := rolesVal.(string)
	country, ok4 := countryVal.(string)
	division, ok5 := divisionVal.(string)

	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return nil
	}

	return &Claims{
		CID:          cid,
		Email:        email,
		Roles:        roles,
		CountryName:  country,
		DivisionName: division,
	}
}
