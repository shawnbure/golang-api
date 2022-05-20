package services

import (
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
)

func GetAllExplorerTokens(args GetAllExplorerTokensArgs) ([]entities.TokenExplorer, int64, float64, float64, error) {
	// Get tokens count by filter
	total, err := storage.GetTokensCountWithCriteria(args.Filter, args.CollectionFilter, args.Attributes)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	//var min, max float64
	//var err1 error
	//if args.IsVerified {
	//	min, max, err1 = storage.GetVerifiedTokensPriceBoundary(args.Filter, args.Attributes)
	//} else {
	min, max, err1 := storage.GetTokensPriceBoundary(args.Filter, args.CollectionFilter, args.Attributes)
	//}

	if err1 != nil {
		return nil, 0, 0, 0, err1
	}

	tokens, err := storage.GetAllTokens(args.LastTimestamp, args.CurrentPage, args.NextPage, args.Limit, args.Filter, args.SortOptions, args.CollectionFilter, args.Attributes)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	if args.NextPage < args.CurrentPage {
		// reversing array
		for i, j := 0, len(tokens)-1; i < j; i, j = i+1, j-1 {
			tokens[i], tokens[j] = tokens[j], tokens[i]
		}
	}

	for index, token := range tokens {
		if token.Token.LastMarketTimestamp == 0 {
			token.Token.LastMarketTimestamp = token.Token.CreatedAt
			tokens[index] = token
		}
	}

	return tokens, total, min, max, nil
}
