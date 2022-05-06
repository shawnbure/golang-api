package services

import (
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
)

func GetAllExplorerTokens(args GetAllExplorerTokensArgs) ([]entities.TokenExplorer, error) {
	tokens, err := storage.GetAllTokens(args.LastTimestamp, args.CurrentPage, args.NextPage, args.Limit, args.Filter, args.SortOptions)
	if err != nil {
		return nil, err
	}

	if args.NextPage < args.CurrentPage {
		// reversing array
		for i, j := 0, len(tokens)-1; i < j; i, j = i+1, j-1 {
			tokens[i], tokens[j] = tokens[j], tokens[i]
		}
	}

	return tokens, nil
}
