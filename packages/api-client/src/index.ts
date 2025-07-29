import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { paths, components } from './types';

export * from './types';

// Type aliases for cleaner code
export type AuthResponse = components['schemas']['AuthResponse'];
export type Balance = components['schemas']['Balance'];
export type Transaction = components['schemas']['Transaction'];
export type YieldPool = components['schemas']['YieldPool'];
export type YieldPosition = components['schemas']['YieldPosition'];
export type Alert = components['schemas']['Alert'];
export type SwapRoute = components['schemas']['SwapRoute'];
export type BridgeRoute = components['schemas']['BridgeRoute'];
export type Watchlist = components['schemas']['Watchlist'];
export type PnLExport = components['schemas']['PnLExport'];

export interface DefiAPIConfig {
  baseURL: string;
  authToken?: string;
  timeout?: number;
  headers?: Record<string, string>;
}

export class DefiAPIClient {
  private client: AxiosInstance;
  private authToken?: string;

  constructor(config: DefiAPIConfig) {
    this.client = axios.create({
      baseURL: config.baseURL,
      timeout: config.timeout || 30000,
      headers: {
        'Content-Type': 'application/json',
        ...config.headers,
      },
    });

    if (config.authToken) {
      this.setAuthToken(config.authToken);
    }

    // Request interceptor to add auth token
    this.client.interceptors.request.use((config) => {
      if (this.authToken) {
        config.headers.Authorization = `Bearer ${this.authToken}`;
      }
      return config;
    });

    // Response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          // Handle unauthorized
          this.authToken = undefined;
        }
        return Promise.reject(error);
      }
    );
  }

  setAuthToken(token: string) {
    this.authToken = token;
  }

  clearAuthToken() {
    this.authToken = undefined;
  }

  // Auth endpoints
  async getNonce(address: string) {
    const response = await this.client.post<components['schemas']['NonceResponse']>(
      '/api/v1/auth/siwe/nonce',
      { address }
    );
    return response.data;
  }

  async verifySiwe(message: string, signature: string) {
    const response = await this.client.post<AuthResponse>(
      '/api/v1/auth/siwe/verify',
      { message, signature }
    );
    this.setAuthToken(response.data.token);
    return response.data;
  }

  async getCurrentUser() {
    const response = await this.client.get<components['schemas']['UserProfileResponse']>(
      '/api/v1/auth/me'
    );
    return response.data;
  }

  // Portfolio endpoints
  async getBalances(address: string, params?: { chainId?: number; hideSmall?: boolean }) {
    const response = await this.client.get<{ balances: Balance[]; totalValue: number }>(
      `/api/v1/portfolio/${address}/balances`,
      { params }
    );
    return response.data;
  }

  async getPortfolioHistory(
    address: string,
    params?: { chainId?: number; period?: string; interval?: string }
  ) {
    const response = await this.client.get<{ history: components['schemas']['PortfolioHistory'][] }>(
      `/api/v1/portfolio/${address}/history`,
      { params }
    );
    return response.data;
  }

  // Transaction endpoints
  async getTransactions(
    address: string,
    params?: { page?: number; limit?: number; chainId?: number; type?: string }
  ) {
    const response = await this.client.get<{
      transactions: Transaction[];
      meta: components['schemas']['PaginationMeta'];
    }>(`/api/v1/transactions/${address}`, { params });
    return response.data;
  }

  async getApprovals(address: string, params?: { chainId?: number }) {
    const response = await this.client.get<{ approvals: components['schemas']['TokenApproval'][] }>(
      `/api/v1/transactions/${address}/approvals`,
      { params }
    );
    return response.data;
  }

  async revokeApproval(address: string, token: string) {
    const response = await this.client.delete(
      `/api/v1/transactions/${address}/approvals/${token}`
    );
    return response.data;
  }

  // Yield endpoints
  async getYieldPools(params?: {
    chainId?: number;
    protocol?: string;
    page?: number;
    limit?: number;
  }) {
    const response = await this.client.get<{
      pools: YieldPool[];
      meta: components['schemas']['PaginationMeta'];
    }>('/api/v1/yield/pools', { params });
    return response.data;
  }

  async getYieldPositions(address: string) {
    const response = await this.client.get<{ positions: YieldPosition[] }>(
      `/api/v1/yield/positions/${address}`
    );
    return response.data;
  }

  // Bridge endpoints
  async getBridgeRoutes(request: components['schemas']['BridgeRoute']) {
    const response = await this.client.post<{ routes: BridgeRoute[] }>(
      '/api/v1/bridge/routes',
      request
    );
    return response.data;
  }

  // Swap endpoints
  async getSwapQuote(request: components['schemas']['SwapQuoteRequest']) {
    const response = await this.client.post<SwapRoute[]>('/api/v1/swap/quote', request);
    return response.data;
  }

  async executeSwap(routeId: string, userAddress: string) {
    const response = await this.client.post<{ txHash: string }>('/api/v1/swap/execute', {
      routeId,
      userAddress,
    });
    return response.data;
  }

  // Alert endpoints
  async getAlerts(params?: { page?: number; limit?: number; status?: string }) {
    const response = await this.client.get<{
      alerts: Alert[];
      meta: components['schemas']['PaginationMeta'];
    }>('/api/v1/alerts', { params });
    return response.data;
  }

  async createAlert(alert: Omit<Alert, 'id' | 'createdAt'>) {
    const response = await this.client.post<Alert>('/api/v1/alerts', alert);
    return response.data;
  }

  async updateAlert(alertId: string, updates: Partial<Alert>) {
    const response = await this.client.patch<Alert>(`/api/v1/alerts/${alertId}`, updates);
    return response.data;
  }

  async deleteAlert(alertId: string) {
    await this.client.delete(`/api/v1/alerts/${alertId}`);
  }

  // Watchlist endpoints
  async getWatchlist() {
    const response = await this.client.get<{ items: Watchlist[] }>('/api/v1/watchlist');
    return response.data;
  }

  async addToWatchlist(item: Omit<Watchlist, 'id' | 'createdAt' | 'updatedAt'>) {
    const response = await this.client.post<Watchlist>('/api/v1/watchlist', item);
    return response.data;
  }

  async removeFromWatchlist(id: string) {
    await this.client.delete(`/api/v1/watchlist/${id}`);
  }

  // Analytics endpoints
  async exportPnL(
    address: string,
    params: { startDate: string; endDate: string; format?: 'json' | 'csv' }
  ) {
    const response = await this.client.get<PnLExport>(
      `/api/v1/analytics/pnl/${address}`,
      { params }
    );
    return response.data;
  }

  // Admin endpoints
  async getFeatureFlags() {
    const response = await this.client.get<components['schemas']['FeatureFlag'][]>(
      '/api/v1/admin/feature-flags'
    );
    return response.data;
  }

  async getSystemBanners() {
    const response = await this.client.get<components['schemas']['SystemBanner'][]>(
      '/api/v1/admin/banners'
    );
    return response.data;
  }
}

// Factory function for convenience
export function createDefiAPIClient(config: DefiAPIConfig): DefiAPIClient {
  return new DefiAPIClient(config);
}