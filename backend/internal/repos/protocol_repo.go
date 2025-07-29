package repos

import (
	"context"
	"encoding/json"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type protocolRepository struct {
	db *pgxpool.Pool
}

// NewProtocolRepository creates a new protocol repository
func NewProtocolRepository(db *pgxpool.Pool) ProtocolRepository {
	return &protocolRepository{db: db}
}

func (r *protocolRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Protocol, error) {
	query := `
		SELECT id, name, slug, description, website_url, logo_uri, 
		       category, total_tvl_usd, chains, is_active, risk_level, 
		       created_at, updated_at
		FROM protocols 
		WHERE id = $1 AND is_active = true
	`
	
	var protocol models.Protocol
	var chainsJSON []byte
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&protocol.ID, &protocol.Name, &protocol.Slug, &protocol.Description,
		&protocol.WebsiteURL, &protocol.LogoURI, &protocol.Category,
		&protocol.TotalTVLUSD, &chainsJSON, &protocol.IsActive,
		&protocol.RiskLevel, &protocol.CreatedAt, &protocol.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Parse chains JSON
	if chainsJSON != nil {
		if err := json.Unmarshal(chainsJSON, &protocol.Chains); err != nil {
			return nil, err
		}
	}

	return &protocol, nil
}

func (r *protocolRepository) GetBySlug(ctx context.Context, slug string) (*models.Protocol, error) {
	query := `
		SELECT id, name, slug, description, website_url, logo_uri, 
		       category, total_tvl_usd, chains, is_active, risk_level, 
		       created_at, updated_at
		FROM protocols 
		WHERE slug = $1 AND is_active = true
	`
	
	var protocol models.Protocol
	var chainsJSON []byte
	
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&protocol.ID, &protocol.Name, &protocol.Slug, &protocol.Description,
		&protocol.WebsiteURL, &protocol.LogoURI, &protocol.Category,
		&protocol.TotalTVLUSD, &chainsJSON, &protocol.IsActive,
		&protocol.RiskLevel, &protocol.CreatedAt, &protocol.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Parse chains JSON
	if chainsJSON != nil {
		if err := json.Unmarshal(chainsJSON, &protocol.Chains); err != nil {
			return nil, err
		}
	}

	return &protocol, nil
}

func (r *protocolRepository) GetAll(ctx context.Context, filters ProtocolFilters) ([]*models.Protocol, error) {
	query := `
		SELECT id, name, slug, description, website_url, logo_uri, 
		       category, total_tvl_usd, chains, is_active, risk_level, 
		       created_at, updated_at
		FROM protocols
		WHERE ($1::varchar IS NULL OR category = $1)
		  AND ($2::boolean IS NULL OR is_active = $2)
		  AND ($3::varchar IS NULL OR risk_level = $3)
		ORDER BY 
		  CASE WHEN $4 = 'name' THEN name END ASC,
		  CASE WHEN $4 = 'tvl' THEN total_tvl_usd END DESC,
		  CASE WHEN $4 = 'category' THEN category END ASC,
		  created_at DESC
		LIMIT $5 OFFSET $6
	`
	
	rows, err := r.db.Query(ctx, query, 
		filters.Category, filters.IsActive, filters.RiskLevel,
		filters.SortBy, filters.Limit, filters.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var protocols []*models.Protocol
	for rows.Next() {
		var protocol models.Protocol
		var chainsJSON []byte
		
		err := rows.Scan(
			&protocol.ID, &protocol.Name, &protocol.Slug, &protocol.Description,
			&protocol.WebsiteURL, &protocol.LogoURI, &protocol.Category,
			&protocol.TotalTVLUSD, &chainsJSON, &protocol.IsActive,
			&protocol.RiskLevel, &protocol.CreatedAt, &protocol.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse chains JSON
		if chainsJSON != nil {
			if err := json.Unmarshal(chainsJSON, &protocol.Chains); err != nil {
				return nil, err
			}
		}

		protocols = append(protocols, &protocol)
	}

	return protocols, nil
}

func (r *protocolRepository) Count(ctx context.Context, filters ProtocolFilters) (int64, error) {
	query := `
		SELECT COUNT(*) FROM protocols
		WHERE ($1::varchar IS NULL OR category = $1)
		  AND ($2::boolean IS NULL OR is_active = $2)
		  AND ($3::varchar IS NULL OR risk_level = $3)
	`
	
	var count int64
	err := r.db.QueryRow(ctx, query, 
		filters.Category, filters.IsActive, filters.RiskLevel).Scan(&count)
	return count, err
}

func (r *protocolRepository) Create(ctx context.Context, protocol *models.Protocol) (*models.Protocol, error) {
	// Serialize chains to JSON
	chainsJSON, err := json.Marshal(protocol.Chains)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO protocols (
			name, slug, description, website_url, logo_uri, 
			category, total_tvl_usd, chains, is_active, risk_level
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	
	err = r.db.QueryRow(ctx, query,
		protocol.Name, protocol.Slug, protocol.Description,
		protocol.WebsiteURL, protocol.LogoURI, protocol.Category,
		protocol.TotalTVLUSD, chainsJSON, protocol.IsActive,
		protocol.RiskLevel,
	).Scan(&protocol.ID, &protocol.CreatedAt, &protocol.UpdatedAt)
	
	return protocol, err
}

func (r *protocolRepository) Update(ctx context.Context, protocol *models.Protocol) (*models.Protocol, error) {
	// Serialize chains to JSON
	chainsJSON, err := json.Marshal(protocol.Chains)
	if err != nil {
		return nil, err
	}

	query := `
		UPDATE protocols 
		SET name = $2,
		    description = $3,
		    website_url = $4,
		    logo_uri = $5,
		    category = $6,
		    total_tvl_usd = $7,
		    chains = $8,
		    is_active = $9,
		    risk_level = $10,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`
	
	err = r.db.QueryRow(ctx, query,
		protocol.ID, protocol.Name, protocol.Description,
		protocol.WebsiteURL, protocol.LogoURI, protocol.Category,
		protocol.TotalTVLUSD, chainsJSON, protocol.IsActive,
		protocol.RiskLevel,
	).Scan(&protocol.UpdatedAt)
	
	return protocol, err
}

func (r *protocolRepository) UpdateTVL(ctx context.Context, id uuid.UUID, tvlUSD float64) error {
	query := `
		UPDATE protocols 
		SET total_tvl_usd = $2,
		    updated_at = NOW()
		WHERE id = $1
	`
	
	_, err := r.db.Exec(ctx, query, id, tvlUSD)
	return err
}

func (r *protocolRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE protocols 
		SET is_active = false,
		    updated_at = NOW()
		WHERE id = $1
	`
	
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *protocolRepository) GetByChain(ctx context.Context, chainID int) ([]*models.Protocol, error) {
	query := `
		SELECT id, name, slug, description, website_url, logo_uri, 
		       category, total_tvl_usd, chains, is_active, risk_level, 
		       created_at, updated_at
		FROM protocols
		WHERE chains ? $1::text
		  AND is_active = true
		ORDER BY total_tvl_usd DESC
	`
	
	rows, err := r.db.Query(ctx, query, chainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var protocols []*models.Protocol
	for rows.Next() {
		var protocol models.Protocol
		var chainsJSON []byte
		
		err := rows.Scan(
			&protocol.ID, &protocol.Name, &protocol.Slug, &protocol.Description,
			&protocol.WebsiteURL, &protocol.LogoURI, &protocol.Category,
			&protocol.TotalTVLUSD, &chainsJSON, &protocol.IsActive,
			&protocol.RiskLevel, &protocol.CreatedAt, &protocol.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse chains JSON
		if chainsJSON != nil {
			if err := json.Unmarshal(chainsJSON, &protocol.Chains); err != nil {
				return nil, err
			}
		}

		protocols = append(protocols, &protocol)
	}

	return protocols, nil
}

func (r *protocolRepository) GetWithPoolCount(ctx context.Context, limit, offset int) ([]*models.Protocol, error) {
	query := `
		SELECT p.id, p.name, p.slug, p.description, p.website_url, p.logo_uri, 
		       p.category, p.total_tvl_usd, p.chains, p.is_active, p.risk_level, 
		       p.created_at, p.updated_at,
		       COUNT(yp.id) as pool_count,
		       COALESCE(SUM(yp.tvl_usd), 0) as total_pools_tvl
		FROM protocols p
		LEFT JOIN yield_pools yp ON p.id = yp.protocol_id AND yp.is_active = true
		WHERE p.is_active = true
		GROUP BY p.id
		ORDER BY total_pools_tvl DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var protocols []*models.Protocol
	for rows.Next() {
		var protocol models.Protocol
		var chainsJSON []byte
		var poolCount int64
		var totalPoolsTVL float64
		
		err := rows.Scan(
			&protocol.ID, &protocol.Name, &protocol.Slug, &protocol.Description,
			&protocol.WebsiteURL, &protocol.LogoURI, &protocol.Category,
			&protocol.TotalTVLUSD, &chainsJSON, &protocol.IsActive,
			&protocol.RiskLevel, &protocol.CreatedAt, &protocol.UpdatedAt,
			&poolCount, &totalPoolsTVL,
		)
		if err != nil {
			return nil, err
		}

		// Parse chains JSON
		if chainsJSON != nil {
			if err := json.Unmarshal(chainsJSON, &protocol.Chains); err != nil {
				return nil, err
			}
		}

		protocols = append(protocols, &protocol)
	}

	return protocols, nil
}