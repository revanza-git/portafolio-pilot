// Temporary mock API client to bypass package issues
interface MockAPIClient {
  setAuthToken: (token: string) => void;
  clearAuthToken: () => void;
  getCurrentUser: () => Promise<any>;
  getNonce: (address: string) => Promise<any>;
  verifySiwe: (message: string, signature: string) => Promise<any>;
  getBalances: (address: string, params?: any) => Promise<any>;
  getTransactions: (address: string, params?: any) => Promise<any>;
}

class MockAPIClientImpl implements MockAPIClient {
  setAuthToken(token: string) {
    console.log('Mock API: Setting auth token');
  }
  
  clearAuthToken() {
    console.log('Mock API: Clearing auth token');
  }

  async getCurrentUser() {
    console.log('Mock API: Getting current user');
    return { id: 'mock-user', address: '0x123' };
  }

  async getNonce(address: string) {
    console.log('Mock API: Getting nonce for', address);
    return { nonce: 'mock-nonce' };
  }

  async verifySiwe(message: string, signature: string) {
    console.log('Mock API: Verifying SIWE');
    return { token: 'mock-token', user: { id: 'mock-user', address: '0x123' } };
  }

  async getBalances(address: string, params?: any) {
    console.log('Mock API: Getting balances for', address);
    return { balances: [], totalValue: 0 };
  }

  async getTransactions(address: string, params?: any) {
    console.log('Mock API: Getting transactions for', address);
    return { transactions: [], meta: { total: 0, page: 1, limit: 20 } };
  }
}

let apiClient: MockAPIClient | null = null;

export function getAPIClient(): MockAPIClient {
  if (!apiClient) {
    apiClient = new MockAPIClientImpl();
  }

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
export function useAPIClient(): MockAPIClient {
  return getAPIClient();
}