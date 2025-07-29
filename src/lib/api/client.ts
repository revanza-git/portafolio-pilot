import { createDefiAPIClient, DefiAPIClient } from '@defip/api-client';
import { API_CONFIG } from './config';

let apiClient: DefiAPIClient | null = null;

export function getAPIClient(): DefiAPIClient {
  if (!apiClient) {
    apiClient = createDefiAPIClient({
      baseURL: API_CONFIG.baseUrl,
      timeout: 30000,
    });
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
export function useAPIClient(): DefiAPIClient {
  return getAPIClient();
}