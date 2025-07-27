# DeFi Portfolio MVP - Real Data Integration

## 📊 Real Data Sources

The application now integrates with real market data APIs:

### CoinGecko API (Price Data)
- **Endpoint**: `https://api.coingecko.com/api/v3`
- **Rate Limit**: 50 requests/minute (free tier)
- **Features**: Token prices, 24h changes, price history
- **No API key required** for basic usage

### DefiLlama API (Yield & TVL Data)  
- **Endpoint**: `https://api.llama.fi`
- **Rate Limit**: 300 requests/minute (free tier)
- **Features**: Yield pools, APR/APY data, protocol TVL
- **No API key required**

## 🔧 Implementation Details

### Frontend Integration
```typescript
// Real-time price updates
const { data: prices } = useTokenPrices(['ETH', 'USDC', 'UNI']);

// Yield pool data with live APR
const { data: pools } = useYieldPools('ethereum');

// Price history for charts
const { data: history } = usePriceHistory('ethereum', 7);
```

### Caching & Fallbacks
- **Client-side caching**: 5-minute cache for price data, 10-minute for yield data
- **Rate limiting**: Built-in throttling to respect API limits
- **Graceful fallbacks**: Returns cached/mock data if APIs fail
- **Auto-retry**: 3 retry attempts with exponential backoff

### Data Refresh Strategy
- **Price data**: Refreshes every 5 minutes
- **Yield data**: Refreshes every 15 minutes  
- **Portfolio calculations**: Real-time based on latest prices
- **Background updates**: Uses React Query for efficient caching

## 🏗️ Backend Architecture (Recommended)

For production, implement a Go backend with:

### API Services
```go
// pkg/coingecko/client.go
type CoinGeckoClient struct {
    httpClient *http.Client
    rateLimit  *rate.Limiter
    cache      *cache.Cache
}

// pkg/defillama/client.go  
type DefiLlamaClient struct {
    httpClient *http.Client
    rateLimit  *rate.Limiter
    cache      *cache.Cache
}
```

### Database Schema
```sql
-- Price history caching
CREATE TABLE price_history (
    id SERIAL PRIMARY KEY,
    token_id INT REFERENCES tokens(id),
    price_usd NUMERIC,
    timestamp TIMESTAMPTZ,
    source TEXT DEFAULT 'coingecko'
);

-- Yield pool data
CREATE TABLE yield_pools (
    id SERIAL PRIMARY KEY,
    protocol_id INT REFERENCES protocols(id),
    pool_address TEXT,
    apr NUMERIC,
    tvl_usd NUMERIC,
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- API rate limiting
CREATE TABLE api_rate_limits (
    id SERIAL PRIMARY KEY,
    service TEXT,
    requests_count INT DEFAULT 0,
    window_start TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (service, window_start)
);
```

### Cron Jobs
```go
// cmd/worker/main.go
func main() {
    // Every 5 minutes - price updates
    c.AddFunc("*/5 * * * *", updateTokenPrices)
    
    // Every 10 minutes - yield pool updates  
    c.AddFunc("*/10 * * * *", updateYieldPools)
    
    // Every hour - cleanup old cache
    c.AddFunc("0 * * * *", cleanupOldData)
}
```

## 🔑 API Keys & Rate Limits

### Current Setup (Frontend Only)
- **No API keys required** - using free tiers
- **Client-side rate limiting** - respects API limits
- **CORS-enabled endpoints** - direct browser requests

### Production Recommendations  
- **CoinGecko Pro**: $129/month for higher limits
- **Proxy through backend**: Hide API keys, better caching
- **Database caching**: Reduce API calls, faster response times

## 🧪 Testing & Monitoring

### Error Handling
```typescript
// Automatic fallbacks
if (apiError) {
    console.warn('API failed, using cached data');
    return getCachedData() || getMockData();
}
```

### Monitoring
- **API response times**: Tracked in React Query
- **Error rates**: Console logging with retry counts  
- **Cache hit rates**: Local cache performance metrics

## 🚀 Deployment Notes

### Environment Variables
```bash
# Optional for enhanced features
VITE_COINGECKO_API_KEY=your_api_key_here
VITE_ENABLE_REAL_DATA=true
VITE_API_CACHE_TTL=300000
```

### Performance Optimizations
- **Request batching**: Multiple tokens in single API call
- **Background refresh**: Updates don't block UI
- **Stale-while-revalidate**: Show cached data while fetching new
- **Error boundaries**: Graceful degradation on API failures

## 📈 Current Data Flow

1. **Component mounts** → React Query fetches data
2. **API client** → Checks cache, respects rate limits  
3. **External API** → CoinGecko/DefiLlama response
4. **Transform data** → Convert to app format
5. **Update cache** → Store for future requests
6. **Re-render** → Components get fresh data

The application now provides real, live market data while maintaining excellent performance and reliability through smart caching and fallback strategies.