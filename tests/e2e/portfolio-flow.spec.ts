import { test, expect } from '@playwright/test';

// Test configuration
const API_URL = process.env.API_URL || 'http://localhost:3000';
const WEB_URL = process.env.WEB_URL || 'http://localhost:8080';

test.describe('Portfolio Management Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the app
    await page.goto(WEB_URL);
  });

  test('Connect wallet and view portfolio', async ({ page }) => {
    // Click connect wallet button
    await page.click('text=Connect Wallet');
    
    // For testing purposes, we'll mock the wallet connection
    // In a real test, you'd use a test wallet provider
    await page.evaluate(() => {
      // Mock wallet connection
      window.localStorage.setItem('wallet_connected', 'true');
      window.localStorage.setItem('wallet_address', '0x742d35Cc6634C0532925a3b844Bc9e7095Ed6aA2');
    });
    
    // Refresh to apply mocked state
    await page.reload();
    
    // Wait for dashboard to load
    await page.waitForSelector('[data-testid="portfolio-overview"]', { timeout: 10000 });
    
    // Verify portfolio overview is visible
    expect(await page.isVisible('[data-testid="total-value"]')).toBeTruthy();
    expect(await page.isVisible('[data-testid="24h-change"]')).toBeTruthy();
    
    // Verify token table is loaded
    await page.waitForSelector('[data-testid="token-table"]');
    const tokenRows = await page.$$('[data-testid="token-row"]');
    expect(tokenRows.length).toBeGreaterThan(0);
  });

  test('Create and manage alerts', async ({ page, context }) => {
    // Setup authentication
    await context.addCookies([{
      name: 'auth_token',
      value: 'test-jwt-token',
      domain: 'localhost',
      path: '/'
    }]);
    
    // Navigate to alerts page
    await page.goto(`${WEB_URL}/alerts`);
    
    // Click create alert button
    await page.click('[data-testid="create-alert-btn"]');
    
    // Fill alert form
    await page.selectOption('[data-testid="alert-type"]', 'price');
    await page.fill('[data-testid="alert-token"]', 'ETH');
    await page.selectOption('[data-testid="alert-condition"]', 'above');
    await page.fill('[data-testid="alert-threshold"]', '3000');
    await page.selectOption('[data-testid="alert-channel"]', 'email');
    
    // Submit form
    await page.click('[data-testid="submit-alert"]');
    
    // Verify alert was created
    await page.waitForSelector('[data-testid="alert-item"]');
    const alertText = await page.textContent('[data-testid="alert-item"]');
    expect(alertText).toContain('ETH');
    expect(alertText).toContain('above');
    expect(alertText).toContain('3000');
    
    // Toggle alert status
    await page.click('[data-testid="toggle-alert"]');
    await page.waitForTimeout(500); // Wait for API call
    
    // Delete alert
    await page.click('[data-testid="delete-alert"]');
    await page.click('[data-testid="confirm-delete"]');
    
    // Verify alert was deleted
    await page.waitForSelector('[data-testid="no-alerts-message"]');
  });

  test('Export P&L report', async ({ page, context }) => {
    // Setup authentication
    await context.addCookies([{
      name: 'auth_token',
      value: 'test-jwt-token',
      domain: 'localhost',
      path: '/'
    }]);
    
    // Navigate to analytics page
    await page.goto(`${WEB_URL}/analytics`);
    
    // Select date range
    await page.click('[data-testid="date-range-picker"]');
    await page.click('[data-testid="last-30-days"]');
    
    // Wait for data to load
    await page.waitForSelector('[data-testid="pnl-chart"]');
    
    // Click export button
    const [download] = await Promise.all([
      page.waitForEvent('download'),
      page.click('[data-testid="export-csv-btn"]')
    ]);
    
    // Verify download
    expect(download.suggestedFilename()).toContain('pnl-report');
    expect(download.suggestedFilename()).toContain('.csv');
  });

  test('Manage watchlist', async ({ page, context }) => {
    // Setup authentication
    await context.addCookies([{
      name: 'auth_token',
      value: 'test-jwt-token',
      domain: 'localhost',
      path: '/'
    }]);
    
    // Navigate to watchlist page
    await page.goto(`${WEB_URL}/watchlist`);
    
    // Add token to watchlist
    await page.click('[data-testid="add-to-watchlist"]');
    await page.fill('[data-testid="search-token"]', 'UNI');
    await page.click('[data-testid="token-option-UNI"]');
    await page.click('[data-testid="confirm-add"]');
    
    // Verify token was added
    await page.waitForSelector('[data-testid="watchlist-item-UNI"]');
    
    // Remove from watchlist
    await page.hover('[data-testid="watchlist-item-UNI"]');
    await page.click('[data-testid="remove-from-watchlist"]');
    
    // Verify token was removed
    await expect(page.locator('[data-testid="watchlist-item-UNI"]')).not.toBeVisible();
  });
});

test.describe('API Integration Tests', () => {
  test('Portfolio API returns correct data structure', async ({ request }) => {
    const address = '0x742d35Cc6634C0532925a3b844Bc9e7095Ed6aA2';
    
    // Get auth token first
    const nonceResponse = await request.post(`${API_URL}/api/v1/auth/siwe/nonce`, {
      data: { address }
    });
    expect(nonceResponse.ok()).toBeTruthy();
    const { nonce, message } = await nonceResponse.json();
    expect(nonce).toBeTruthy();
    expect(message).toBeTruthy();
    
    // Mock signature verification for testing
    const authToken = 'test-jwt-token';
    
    // Test portfolio balances endpoint
    const balancesResponse = await request.get(`${API_URL}/api/v1/portfolio/${address}/balances`, {
      headers: {
        'Authorization': `Bearer ${authToken}`
      }
    });
    
    expect(balancesResponse.ok()).toBeTruthy();
    const balancesData = await balancesResponse.json();
    
    // Verify response structure
    expect(balancesData).toHaveProperty('totalValue');
    expect(balancesData).toHaveProperty('balances');
    expect(Array.isArray(balancesData.balances)).toBeTruthy();
    
    if (balancesData.balances.length > 0) {
      const balance = balancesData.balances[0];
      expect(balance).toHaveProperty('token');
      expect(balance.token).toHaveProperty('address');
      expect(balance.token).toHaveProperty('symbol');
      expect(balance.token).toHaveProperty('decimals');
      expect(balance).toHaveProperty('balance');
      expect(balance).toHaveProperty('balanceUsd');
    }
  });

  test('Swap quote API returns multiple routes', async ({ request }) => {
    const swapRequest = {
      chainId: 1,
      fromToken: '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2', // WETH
      toToken: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48', // USDC
      fromAmount: '1000000000000000000', // 1 ETH
      userAddress: '0x742d35Cc6634C0532925a3b844Bc9e7095Ed6aA2'
    };
    
    const response = await request.post(`${API_URL}/api/v1/swap/quote`, {
      headers: {
        'Authorization': 'Bearer test-jwt-token',
        'Content-Type': 'application/json'
      },
      data: swapRequest
    });
    
    expect(response.ok()).toBeTruthy();
    const routes = await response.json();
    
    expect(Array.isArray(routes)).toBeTruthy();
    expect(routes.length).toBeGreaterThan(0);
    
    const route = routes[0];
    expect(route).toHaveProperty('id');
    expect(route).toHaveProperty('fromToken');
    expect(route).toHaveProperty('toToken');
    expect(route).toHaveProperty('fromAmount');
    expect(route).toHaveProperty('toAmount');
    expect(route).toHaveProperty('provider');
    expect(route).toHaveProperty('dex');
  });
});