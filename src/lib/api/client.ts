// Real API client that connects to the backend
import { API_CONFIG } from './config';

interface APIClient {
  setAuthToken: (token: string) => void;
  clearAuthToken: () => void;
  getCurrentUser: () => Promise<any>;
  getNonce: (address: string) => Promise<any>;
  verifySiwe: (message: string, signature: string) => Promise<any>;
  getBalances: (address: string, params?: any) => Promise<any>;
  getTransactions: (address: string, params?: any) => Promise<any>;
}

class RealAPIClientImpl implements APIClient {
  private authToken?: string;
  private baseURL: string;

  constructor() {
    this.baseURL = API_CONFIG.baseUrl;
    console.log('Real API Client initialized with baseURL:', this.baseURL);
  }

  setAuthToken(token: string) {
    this.authToken = token;
    console.log('Real API: Setting auth token');
  }
  
  clearAuthToken() {
    this.authToken = undefined;
    console.log('Real API: Clearing auth token');
  }

  private getApiKeysFromStorage() {
    try {
      const savedKeys = localStorage.getItem('defi_api_keys');
      if (savedKeys) {
        return JSON.parse(savedKeys);
      }
    } catch (error) {
      console.error('Failed to load API keys from localStorage:', error);
    }
    return {};
  }

  private async makeRequest(endpoint: string, options: RequestInit = {}) {
    const headers = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (this.authToken) {
      headers['Authorization'] = `Bearer ${this.authToken}`;
    }

    // Add API keys from localStorage to headers
    const apiKeys = this.getApiKeysFromStorage();
    if (apiKeys.alchemyApiKey) {
      headers['X-Alchemy-API-Key'] = apiKeys.alchemyApiKey;
    }
    if (apiKeys.coingeckoApiKey) {
      headers['X-CoinGecko-API-Key'] = apiKeys.coingeckoApiKey;
    }
    if (apiKeys.etherscanApiKey) {
      headers['X-Etherscan-API-Key'] = apiKeys.etherscanApiKey;
    }
    if (apiKeys.infuraApiKey) {
      headers['X-Infura-API-Key'] = apiKeys.infuraApiKey;
    }

    console.log('API Request with keys:', {
      endpoint,
      hasAlchemy: !!apiKeys.alchemyApiKey,
      hasCoinGecko: !!apiKeys.coingeckoApiKey,
      hasEtherscan: !!apiKeys.etherscanApiKey,
      hasInfura: !!apiKeys.infuraApiKey
    });

    const response = await fetch(`${this.baseURL}${endpoint}`, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const errorText = await response.text();
      console.error('API Error:', response.status, errorText);
      throw new Error(`API Error: ${response.status} ${errorText}`);
    }

    return response.json();
  }

  async getCurrentUser() {
    console.log('Real API: Getting current user');
    return this.makeRequest('/api/v1/auth/me');
  }

  async getNonce(address: string) {
    console.log('Real API: Getting nonce for', address);
    return this.makeRequest('/api/v1/auth/siwe/nonce', {
      method: 'POST',
      body: JSON.stringify({ address }),
    });
  }

  async verifySiwe(message: string, signature: string) {
    console.log('Real API: Verifying SIWE');
    return this.makeRequest('/api/v1/auth/siwe/verify', {
      method: 'POST',
      body: JSON.stringify({ message, signature }),
    });
  }

  async getBalances(address: string, params?: any) {
    console.log('Real API: Getting balances for', address, 'with params:', params);
    
    const queryParams = new URLSearchParams();
    if (params?.chainId) queryParams.append('chainId', params.chainId.toString());
    if (params?.hideSmall) queryParams.append('hideSmall', params.hideSmall.toString());
    
    const queryString = queryParams.toString();
    const endpoint = `/api/v1/portfolio/${address}/balances${queryString ? '?' + queryString : ''}`;
    
    try {
      const result = await this.makeRequest(endpoint);
      console.log('Real API: Balance response:', result);
      return result;
    } catch (error) {
      console.error('Failed to fetch balances:', error);
      // Fallback to empty data on error
      return { balances: [], total_value: 0 };
    }
  }

  async getTransactions(address: string, params?: any) {
    console.log('Real API: Getting transactions for', address, 'with params:', params);
    
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.chainId) queryParams.append('chainId', params.chainId.toString());
    if (params?.type) queryParams.append('type', params.type);
    
    const queryString = queryParams.toString();
    const endpoint = `/api/v1/transactions/${address}${queryString ? '?' + queryString : ''}`;
    
    try {
      const result = await this.makeRequest(endpoint);
      console.log('Real API: Transaction response:', result);
      return result;
    } catch (error) {
      console.error('Failed to fetch transactions:', error);
      // Fallback to empty data on error
      return { data: [], meta: { total: 0, page: 1, limit: 20 } };
    }
  }
}

let apiClient: APIClient | null = null;

export function getAPIClient(): APIClient {
  if (!apiClient) {
    apiClient = new RealAPIClientImpl();
  }

  // Always check for token on every call to ensure it's up to date
  const authToken = localStorage.getItem('auth_token');
  console.log('API Client: Checking auth token:', !!authToken, 'length:', authToken?.length);
  
  if (authToken) {
    console.log('API Client: Setting auth token on client');
    apiClient.setAuthToken(authToken);
  } else {
    console.log('API Client: No auth token found, clearing');
    apiClient.clearAuthToken();
  }

  return apiClient;
}

// Hook to use in React components
export function useAPIClient(): APIClient {
  return getAPIClient();
}