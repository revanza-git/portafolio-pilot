# Backend Setup for Bridge & Swap Routes

Since Lovable is a frontend-only platform, you'll need to implement backend proxy endpoints using **Supabase Edge Functions** to handle API calls to LI.FI, Socket, 0x, and 1inch APIs.

## Required Edge Functions

### 1. Bridge Routes Function (`/bridge/routes`)

Create a Supabase Edge Function to proxy bridge quotes:

```typescript
// supabase/functions/bridge-routes/index.ts
import { serve } from "https://deno.land/std@0.168.0/http/server.ts"

const LIFI_API_KEY = Deno.env.get('LIFI_API_KEY')
const SOCKET_API_KEY = Deno.env.get('SOCKET_API_KEY')

serve(async (req) => {
  if (req.method !== 'POST') {
    return new Response('Method not allowed', { status: 405 })
  }

  try {
    const { fromChain, toChain, fromToken, toToken, fromAmount, userAddress, slippage } = await req.json()

    // Fetch from LI.FI
    const lifiResponse = await fetch('https://li.quest/v1/quote', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'x-lifi-api-key': LIFI_API_KEY || ''
      },
      body: JSON.stringify({
        fromChain,
        toChain,
        fromToken,
        toToken,
        fromAmount,
        fromAddress: userAddress,
        slippage: slippage || 0.5
      })
    })

    // Fetch from Socket
    const socketResponse = await fetch('https://api.socket.tech/v2/quote', {
      method: 'GET',
      headers: {
        'API-KEY': SOCKET_API_KEY || '',
        'Content-Type': 'application/json'
      }
    })

    const routes = []
    
    if (lifiResponse.ok) {
      const lifiData = await lifiResponse.json()
      routes.push(...transformLifiRoutes(lifiData))
    }

    if (socketResponse.ok) {
      const socketData = await socketResponse.json()
      routes.push(...transformSocketRoutes(socketData))
    }

    return new Response(JSON.stringify(routes), {
      headers: { 'Content-Type': 'application/json' }
    })
  } catch (error) {
    return new Response(JSON.stringify({ error: error.message }), {
      status: 500,
      headers: { 'Content-Type': 'application/json' }
    })
  }
})
```

### 2. Swap Quotes Function (`/swap/quote`)

```typescript
// supabase/functions/swap-quote/index.ts
import { serve } from "https://deno.land/std@0.168.0/http/server.ts"

const ZEROX_API_KEY = Deno.env.get('ZEROX_API_KEY')
const ONEINCH_API_KEY = Deno.env.get('ONEINCH_API_KEY')

serve(async (req) => {
  if (req.method !== 'POST') {
    return new Response('Method not allowed', { status: 405 })
  }

  try {
    const { chainId, fromToken, toToken, fromAmount, userAddress, slippage } = await req.json()

    // Fetch from 0x API
    const zeroXUrl = `https://api.0x.org/swap/v1/quote?sellToken=${fromToken}&buyToken=${toToken}&sellAmount=${fromAmount}&takerAddress=${userAddress}&slippagePercentage=${slippage || 0.5}`
    
    const zeroXResponse = await fetch(zeroXUrl, {
      headers: {
        '0x-api-key': ZEROX_API_KEY || ''
      }
    })

    // Fetch from 1inch API
    const oneInchUrl = `https://api.1inch.io/v5.0/${chainId}/quote?fromTokenAddress=${fromToken}&toTokenAddress=${toToken}&amount=${fromAmount}&slippage=${slippage || 0.5}`
    
    const oneInchResponse = await fetch(oneInchUrl, {
      headers: {
        'Authorization': `Bearer ${ONEINCH_API_KEY || ''}`
      }
    })

    const routes = []
    
    if (zeroXResponse.ok) {
      const zeroXData = await zeroXResponse.json()
      routes.push(transform0xRoute(zeroXData))
    }

    if (oneInchResponse.ok) {
      const oneInchData = await oneInchResponse.json()
      routes.push(transform1inchRoute(oneInchData))
    }

    return new Response(JSON.stringify(routes), {
      headers: { 'Content-Type': 'application/json' }
    })
  } catch (error) {
    return new Response(JSON.stringify({ error: error.message }), {
      status: 500,
      headers: { 'Content-Type': 'application/json' }
    })
  }
})
```

## Setting up API Keys

1. **Connect to Supabase**: Use the Supabase integration in Lovable
2. **Add Secrets**: Use the Supabase dashboard to add these environment variables:
   - `LIFI_API_KEY` - Get from https://li.quest/
   - `SOCKET_API_KEY` - Get from https://socket.tech/
   - `ZEROX_API_KEY` - Get from https://0x.org/
   - `ONEINCH_API_KEY` - Get from https://1inch.io/

## Database Schema for Analytics

```sql
-- Create tables for route analytics
CREATE TABLE bridge_routes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_address TEXT NOT NULL,
  from_chain INTEGER NOT NULL,
  to_chain INTEGER NOT NULL,
  from_token TEXT NOT NULL,
  to_token TEXT NOT NULL,
  from_amount TEXT NOT NULL,
  to_amount TEXT NOT NULL,
  provider TEXT NOT NULL,
  fees_total DECIMAL,
  estimated_time INTEGER,
  tx_hash TEXT,
  status TEXT DEFAULT 'pending',
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE swap_routes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_address TEXT NOT NULL,
  chain_id INTEGER NOT NULL,
  from_token TEXT NOT NULL,
  to_token TEXT NOT NULL,
  from_amount TEXT NOT NULL,
  to_amount TEXT NOT NULL,
  provider TEXT NOT NULL,
  dex TEXT,
  price_impact DECIMAL,
  fees_total DECIMAL,
  tx_hash TEXT,
  status TEXT DEFAULT 'pending',
  created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_bridge_routes_user ON bridge_routes(user_address);
CREATE INDEX idx_swap_routes_user ON swap_routes(user_address);
CREATE INDEX idx_bridge_routes_created ON bridge_routes(created_at);
CREATE INDEX idx_swap_routes_created ON swap_routes(created_at);
```

## Frontend Configuration

Update your API clients to use the Supabase Edge Function URLs:

```typescript
// In src/lib/api/bridge-client.ts and swap-client.ts
const baseUrl = 'https://your-project.supabase.co/functions/v1';
```

## Rate Limiting & Caching

Implement caching in your Edge Functions to reduce API calls:

```typescript
// Add Redis or in-memory caching
const cacheKey = `quote_${fromToken}_${toToken}_${fromAmount}`
const cachedResult = await redis.get(cacheKey)

if (cachedResult) {
  return new Response(cachedResult, {
    headers: { 'Content-Type': 'application/json' }
  })
}

// ... fetch and cache result for 30 seconds
await redis.setex(cacheKey, 30, JSON.stringify(routes))
```

## Next Steps

1. Set up Supabase project and connect to Lovable
2. Deploy the Edge Functions above
3. Configure API keys in Supabase secrets
4. Update frontend API URLs to use your Edge Functions
5. Test the integration with real API calls