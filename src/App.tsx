import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { WagmiProvider } from 'wagmi';
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { ThemeProvider } from "@/components/providers/theme-provider";
import { AuthProvider } from "@/contexts/auth-context";
import { FeatureFlagProvider } from "@/contexts/feature-flag-context";
import { config } from "@/lib/wagmi-config";
import { SystemBanners } from "@/components/shared/system-banner";
import Index from "./pages/Index";
import Dashboard from "./pages/Dashboard";
import Yield from "./pages/Yield";
import Bridge from "./pages/Bridge";
import Analytics from "./pages/Analytics";
import Transactions from "./pages/Transactions";
import Approvals from "./pages/Approvals";
import Swap from "./pages/Swap";
import Alerts from "./pages/Alerts";
import Watchlist from "./pages/Watchlist";
import SignIn from "./pages/SignIn";
import EmailVerification from "./pages/EmailVerification";
import AdminDashboard from "./pages/admin/AdminDashboard";
import FeatureFlags from "./pages/admin/FeatureFlags";
import AdminBanners from "./pages/admin/SystemBanners";
import NotFound from "./pages/NotFound";
import { AdminGuard } from "@/components/admin/admin-guard";

const queryClient = new QueryClient();

const App = () => (
  <WagmiProvider config={config}>
    <QueryClientProvider client={queryClient}>
      {/* AuthProvider temporarily bypassed */}
      {/* <AuthProvider> */}
        <FeatureFlagProvider>
          <ThemeProvider>
            <TooltipProvider>
              <Toaster />
              <Sonner />
              <BrowserRouter>
                <div className="min-h-screen">
                  <SystemBanners />
                  <Routes>
                    <Route path="/" element={<Index />} />
                    <Route path="/dashboard" element={<Dashboard />} />
                    <Route path="/yield" element={<Yield />} />
                    <Route path="/bridge" element={<Bridge />} />
                    <Route path="/analytics" element={<Analytics />} />
                    <Route path="/transactions" element={<Transactions />} />
                    <Route path="/approvals" element={<Approvals />} />
                    <Route path="/swap" element={<Swap />} />
                    <Route path="/alerts" element={<Alerts />} />
                    <Route path="/watchlist" element={<Watchlist />} />
                    <Route path="/auth/signin" element={<SignIn />} />
                    <Route path="/auth/verify-email" element={<EmailVerification />} />
                    {/* Admin routes temporarily disabled */}
                    {/* <Route path="/admin" element={<AdminGuard><AdminDashboard /></AdminGuard>} />
                    <Route path="/admin/feature-flags" element={<AdminGuard><FeatureFlags /></AdminGuard>} />
                    <Route path="/admin/banners" element={<AdminGuard><AdminBanners /></AdminGuard>} /> */}
                    {/* ADD ALL CUSTOM ROUTES ABOVE THE CATCH-ALL "*" ROUTE */}
                    <Route path="*" element={<NotFound />} />
                  </Routes>
                </div>
              </BrowserRouter>
            </TooltipProvider>
          </ThemeProvider>
        </FeatureFlagProvider>
      {/* </AuthProvider> */}
    </QueryClientProvider>
  </WagmiProvider>
);

export default App;
