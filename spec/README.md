# DeFi Dashboard API Specification

This directory contains the OpenAPI 3.1 specification for the DeFi Dashboard backend API.

## Generating TypeScript Types

You can generate TypeScript types from the OpenAPI specification using `openapi-typescript`.

### Installation

```bash
npm install --save-dev openapi-typescript
```

### Basic Usage

```bash
npx openapi-typescript spec/openapi.yaml -o src/types/api.ts
```

### Advanced Usage with Package.json Script

Add to your `package.json`:

```json
{
  "scripts": {
    "generate:types": "openapi-typescript spec/openapi.yaml -o src/types/api.ts"
  }
}
```

Then run:
```bash
npm run generate:types
```

### Using the Generated Types

```typescript
import type { paths, components } from './types/api';

// Type for a specific endpoint
type GetBalancesResponse = paths['/portfolio/{address}/balances']['get']['responses']['200']['content']['application/json'];

// Component schemas
type Token = components['schemas']['Token'];
type Balance = components['schemas']['Balance'];
type Transaction = components['schemas']['Transaction'];

// With a type-safe API client like openapi-fetch
import createClient from 'openapi-fetch';

const client = createClient<paths>({ baseUrl: 'https://api.defi-dashboard.com/v1' });

// Type-safe API calls
const { data, error } = await client.GET('/portfolio/{address}/balances', {
  params: {
    path: { address: '0x...' },
    query: { chainId: 1 }
  }
});
```

### Recommended API Client

For the best developer experience, use `openapi-fetch` alongside the generated types:

```bash
npm install openapi-fetch
```

This provides:
- Full type safety for requests and responses
- Automatic path parameter interpolation
- Request/response validation
- Built-in error handling

### Keeping Types in Sync

To ensure your types stay in sync with the API:

1. Run type generation in your CI/CD pipeline
2. Add a pre-commit hook to regenerate types
3. Version your OpenAPI spec alongside your code